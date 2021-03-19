package whitelist

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	log "github.com/sirupsen/logrus"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/spf13/viper"
)

// Service provides access to User related data
type Service struct {
	db store
}

// DefaultService prepares new instance of Service
func DefaultService() Service {
	return NewService(DefaultData())
}

// NewService prepares new instance of Service
func NewService(whitelistStore store) Service {
	return Service{
		db: whitelistStore,
	}
}

func (us *Service) GetAllWhitelists() ([]Application, error) {
	list, err := us.db.getAllWhitelists()
	if err != nil {
		return nil, errCannotLoadWhitelists
	}
	return list, nil
}

func (us *Service) GetWhitelist(id string) (Application, error) {
	application, err := us.db.getWhitelist(id)
	if err != nil {
		return Application{}, ErrCannotFindWhitelist
	}
	return application, nil
}

func (us *Service) GetWhitelistsForUser(username string) ([]Application, error) {
	applications, err := us.db.findApplications(username)
	if err != nil {
		return []Application{}, ErrCannotFindWhitelist
	}
	return applications, nil
}

func (us *Service) UpdateApplicationStatus(req ChangeApplicationStatus) (string, error) {
	app, err := us.db.getUserApplicationWithHistory(req.ApplicationId)
	if err != nil {
		return "", ErrCannotFindWhitelist
	}
	if app.CurrentStatus != PENDING {
		return "", errNoApplicationInProgressForUser
	}
	latestChange := app.GetLatestChangeHistory()
	change := ChangeHistory{
		Description:   latestChange.Description,
		Location:      latestChange.Location,
		Status:        ApplicationStatus(req.Status),
		ApplicationId: latestChange.ApplicationId,
	}
	change.AdminComment = req.AdminComment
	change.UserComment = req.UserComment

	for _, node := range latestChange.Nodes {
		newNode := &Node{Key: node.Key}
		newNode.ID = node.ID
		newNode.CreatedAt = node.CreatedAt
		change.Nodes = append(change.Nodes, newNode)
	}
	for _, image := range latestChange.Images {
		change.Images = append(change.Images, &Image{ID: image.ID, Path: image.Path, CreatedAt: image.CreatedAt, ImgHash: image.ImgHash})
	}
	us.db.createChangeHistory(&change)
	app.ChangeHistory = append(app.ChangeHistory, change)
	app.CurrentStatus = ApplicationStatus(req.Status)
	err = us.db.updateApplication(&app, 0)
	if err != nil {
		return "", errCannotUpdateApplication
	}
	if req.Status == APPROVED {
		applications, err := us.db.findApplications(app.Username)
		for _, application := range applications {
			if application.CurrentStatus == AUTO_DISABLED {
				if latestCH := application.GetLatestChangeHistory(); latestCH.AdminComment == autoDeclineMessage {
					newChange := ChangeHistory{
						Description: latestCH.Description,
						Location:    latestCH.Location,
						Status:      PENDING,
					}
					newChange.AdminComment = autoPendingMessage
					newChange.UserComment = autoPendingMessageUser

					for _, node := range latestCH.Nodes {
						newNode := &Node{Key: node.Key}
						newNode.ID = node.ID
						newNode.CreatedAt = node.CreatedAt
						newChange.Nodes = append(newChange.Nodes, newNode)
					}
					for _, image := range latestCH.Images {
						newChange.Images = append(newChange.Images, &Image{ID: image.ID, Path: image.Path, CreatedAt: image.CreatedAt})
					}
					us.db.createChangeHistory(&newChange)
					application.ChangeHistory = append(application.ChangeHistory, newChange)
					application.CurrentStatus = PENDING
					err = us.db.updateApplication(&application, 0)
					if err != nil {
						log.Error("Unable to reactivate previously declined application due to error ", err)
					}
					break
				}
			}
		}
	} else if req.Status == DISABLED {
		var currTime = time.Now()
		for _, node := range change.Nodes {
			if err := us.RemoveNode(node, currTime); err != nil {
				log.Errorf("Unable to delete node with key %v due to error: %v ", node.Key, err)
			}
		}
	}
	return app.Username, nil
}

func (us *Service) RestartApplication(newNodes []Node, minerOrigin Application, location, desctiption string) error {
	newChangeHistory := ChangeHistory{
		ApplicationId: minerOrigin.ID,
		Location:      location,
		Description:   desctiption,
		Status:        PENDING,
		AdminComment:  autoUpdateMessage,
		UserComment:   autoUpdateMessageUser,
		Nodes:         []*Node{},
	}

	for _, n := range newNodes {
		newChangeHistory.Nodes = append(newChangeHistory.Nodes, &Node{ID: n.ID, Key: n.Key, MinerID: n.MinerID, CreatedAt: n.CreatedAt})
	}
	log.Debug("Creating change history for adding nodes to the miner ", newChangeHistory)

	if err := us.db.createChangeHistory(&newChangeHistory); err != nil {
		return ErrUnableToRestartApplication
	}
	minerOrigin.ChangeHistory = append(minerOrigin.ChangeHistory, newChangeHistory)
	minerOrigin.CurrentStatus = PENDING
	if err := us.db.updateApplication(&minerOrigin, 0); err != nil { //TODO check if preloaded Miner can cause problems here
		return ErrUnableToRestartApplication
	}

	return nil
}

func (us *Service) UpdateApplication(appReq ApplicationReq, username string) (applicationError ApplicationError) {
	latestApplication, err := us.getActiveApplication(username)
	if err != nil {
		applicationError.Error = err
		return
	}
	if latestApplication.CurrentStatus != PENDING && latestApplication.CurrentStatus != DECLINED {
		applicationError.Error = errApplicationInProgress
		return
	}

	failedImages, err := us.saveFiles(&appReq, username)
	applicationError.FailedImages = failedImages
	if err != nil {
		applicationError.Error = err
		return
	}

	latestChange := latestApplication.GetLatestChangeHistory()
	var hasChanged bool
	for _, img := range appReq.Files {
		if img.Path != "" {
			hasChanged = true
			break
		}
	}

	if !hasChanged {
		if sameAsPrevious := us.compareRequestWithLatestChangeHistory(appReq, latestChange); sameAsPrevious == true {
			applicationError.Error = errIdenticalAsPrevious
			return
		}
	}

	//validate keys
	us.ValidateNodeKeys(appReq.Nodes, latestApplication.ID, 0, &applicationError)
	if len(applicationError.AlreadyTakenKeys) != 0 || len(applicationError.DuplicateKeys) != 0 || len(applicationError.WrongKeys) != 0 {
		applicationError.Error = errWrongNodeKeys
		return
	}

	newChange := createChangeHistoryForApplication(appReq)
	newChange.ApplicationId = latestApplication.ID
	for _, node := range appReq.Nodes {
		newNode := &Node{Key: node.Key}
		if node.ID > 0 {
			newNode.ID = node.ID
			newNode.CreatedAt = node.CreatedAt
		}
		newChange.Nodes = append(newChange.Nodes, newNode)
	}
	images := createImagesForChangeHistory(appReq.Files)
	for _, oldImage := range appReq.OldImages {
		newChange.Images = append(newChange.Images, &Image{Path: oldImage.Path, ImgHash: oldImage.ImgHash, ID: oldImage.ID, CreatedAt: oldImage.CreatedAt})
	}
	for _, newImage := range images {
		newChange.Images = append(newChange.Images, &Image{Path: newImage.Path, ImgHash: newImage.ImgHash})
	}

	latestApplication.ChangeHistory = append(latestApplication.ChangeHistory, newChange)
	latestApplication.CurrentStatus = PENDING
	us.db.updateApplication(&latestApplication, latestChange.ID)

	return
}

func (us *Service) compareRequestWithLatestChangeHistory(req ApplicationReq, latestChange ChangeHistory) bool {
	if req.Description != latestChange.Description {
		return false
	} else if req.Location != latestChange.Location {
		return false
	} else if len(req.Nodes) != len(latestChange.Nodes) {
		return false
	} else if len(req.OldImages) != len(latestChange.Images) {
		return false
	}
	for _, newNode := range req.Nodes {
		foundMatch := false
		for _, oldNode := range latestChange.Nodes {
			if newNode.Key == oldNode.Key {
				foundMatch = true
				break
			}
		}
		if foundMatch == false {
			return false
		}
	}
	return true
}

func (us *Service) getActiveApplication(username string) (Application, error) {
	applications, err := us.db.findApplications(username)
	if err != nil {
		return Application{}, err
	}
	if len(applications) == 0 {
		return Application{}, errNoApplicationInProgressForUser
	}

	latestApplication := applications[len(applications)-1]

	if latestApplication.CurrentStatus != PENDING {
		for _, application := range applications {
			if application.CurrentStatus == PENDING || application.CurrentStatus == DECLINED {
				latestApplication = application
				break
			}
		}
	}

	if latestApplication.CurrentStatus != PENDING && latestApplication.CurrentStatus != DECLINED {
		return Application{}, errNoApplicationInProgressForUser
	}
	applicationWithHistory, err := us.db.getUserApplicationWithHistory(latestApplication.ID)
	if err != nil {
		return Application{}, err
	}
	return applicationWithHistory, nil
}

func (us *Service) CreateApplication(appReq ApplicationReq, username string) (newApplication Application, applicationError ApplicationError) {
	applications, err := us.db.findApplications(username)
	if err != nil {
		applicationError.Error = err
		return
	}

	for _, application := range applications {
		if application.CurrentStatus == PENDING || application.CurrentStatus == DECLINED {
			applicationError.Error = errApplicationInProgress
			return
		}
	}

	//validate keys
	us.ValidateNodeKeys(appReq.Nodes, 0, 0, &applicationError)
	if len(applicationError.AlreadyTakenKeys) != 0 || len(applicationError.DuplicateKeys) != 0 || len(applicationError.WrongKeys) != 0 {
		applicationError.Error = errWrongNodeKeys
		return
	}

	// save files into application change history
	failedImages, err := us.saveFiles(&appReq, username)
	applicationError.FailedImages = failedImages
	if err != nil {
		applicationError.Error = err
		return
	}

	newApplication = createApplicationFromRequest(appReq, username)
	if err = us.db.createApplication(&newApplication); err != nil {
		applicationError.Error = errUnableToSave
		return
	}
	changeHistory := newApplication.GetLatestChangeHistory()

	for _, node := range appReq.Nodes {
		changeHistory.Nodes = append(changeHistory.Nodes, &Node{ID: node.ID, Key: node.Key})
	}

	for _, image := range appReq.Files {
		if image.Path != "" {
			changeHistory.Images = append(changeHistory.Images, &Image{Path: image.Path, ImgHash: image.Hashed})
		}
	}
	us.db.updateChangeHistory(&changeHistory)

	return
}

func createImagesForChangeHistory(appImages []ApplicationImage) []Image {
	var images []Image
	for _, element := range appImages {
		if element.Path != "" {
			images = append(images, Image{
				Path:    element.Path,
				ImgHash: element.Hashed,
			})
		}
	}
	return images
}

// create pending application
func createApplicationFromRequest(req ApplicationReq, username string) Application {
	changeHistory := createChangeHistoryForApplication(req)
	return Application{
		Username:      username,
		CurrentStatus: PENDING,
		ChangeHistory: []ChangeHistory{
			changeHistory,
		},
	}
}

//create a pending change history for application
func createChangeHistoryForApplication(req ApplicationReq) ChangeHistory {
	return ChangeHistory{
		Description: req.Description,
		Location:    req.Location,
		Status:      PENDING,
	}
}

func (us *Service) ValidateNodeKeys(nodes []Node, applicationID uint, currentMinerID uint, appError *ApplicationError) {
	var uniqueKeys = make(map[string]bool)

	for i := 0; i < len(nodes); i++ {
		nodes[i].Key = strings.Trim(nodes[i].Key, " ")
		if len(nodes[i].Key) == 0 {
			continue
		}
		if len(nodes[i].Key) > 66 {
			appError.WrongKeys = append(appError.WrongKeys, nodes[i].Key)
			continue
		}
		if node, err := us.db.findActiveNodeByKey(nodes[i].Key); err == nil && node.MinerID != 0 &&
			(node.MinerID != currentMinerID || (applicationID == 0 && currentMinerID == 0)) {
			appError.AlreadyTakenKeys = append(appError.AlreadyTakenKeys, nodes[i].Key)
			continue
		} else if err == nil && currentMinerID == 0 {
			if app, err := us.db.getApplicationForNode(node); err != nil || (app.ID != applicationID && app.CurrentStatus != 3) {
				appError.AlreadyTakenKeys = append(appError.AlreadyTakenKeys, nodes[i].Key)
				continue
			}
		}

		if err := checkIsKeyWrong(nodes[i].Key); err != nil {
			appError.WrongKeys = append(appError.WrongKeys, nodes[i].Key)
			continue
		}
		if uniqueKeys[nodes[i].Key] == true {
			appError.DuplicateKeys = append(appError.DuplicateKeys, nodes[i].Key)
			continue
		}
		uniqueKeys[nodes[i].Key] = true
	}
}

//  TODO: check if to move out the path for the files outside to config- no need since it has to be inside project
func (us *Service) saveFiles(appReq *ApplicationReq, path string) ([]string, error) {
	var (
		//status int
		err          error
		failedImages []string
	)
	session, err := us.openS3Connection()
	if err != nil {
		log.Error("Unable to connect to S3 storage", err)
		return failedImages, err
	}
	storedImages, err := us.GetImagesForUser(path)
	if err != nil {
		log.Error("Unable to fetch images to check if there is a picture with existing name due to ", err)
		return failedImages, err
	}

	for i, file := range appReq.Files {
		var savedPath = path + "/" + file.Name
		//check is file name taken, if it is append a unique number to it.
		if checkIsImageNameTaken(storedImages, savedPath) {
			savedPath = path + "/" + strconv.Itoa(len(storedImages)) + file.Name
		}
		// open uploaded
		if err := us.saveToS3(session, savedPath, file.File); err != nil {
			log.Errorf("Unable to persist image %v to database using credentials %v with error %v", savedPath, viper.GetString("bucket.role"), err)
			failedImages = append(failedImages, savedPath)
		} else {
			appReq.Files[i].Path = savedPath
			storedImages = append(storedImages, Image{Path: savedPath})
		}

	}
	var errFailedImages error = nil
	if len(failedImages) > 0 {
		errFailedImages = errUploadingImages
	}
	return failedImages, errFailedImages
}

func (us *Service) openS3Connection() (*session.Session, error) {
	return session.NewSession(&aws.Config{Region: aws.String(viper.GetString("bucket.region"))})
}

func (us *Service) saveToS3(session *session.Session, path string, file []byte) error {
	if viper.GetBool("bucket.disable-image-upload") {
		return nil
	}
	_, err := s3.New(session).PutObject((&s3.PutObjectInput{}).SetBucket(viper.GetString("bucket.name")).
		SetKey(path).
		SetBody(bytes.NewReader(file)).
		SetContentType(http.DetectContentType(file)).
		SetContentDisposition("attachment"),
	)
	return err
}
func (us *Service) UpdateCreatedAt(app Application, createdAtImport string) error {
	err := us.db.appUpdateCreatedAt(&app, createdAtImport)
	if err != nil {
		log.Errorf("Error during updating created at field %v", err)
	}
	return err
}

// GetMinerForApplication returns miner connected to application. This function should be used instead of preloading miner.
func (us *Service) GetMinerForApplication(app Application) Miner {
	return us.db.getMinerForApplication(&app)
}

func (us *Service) RemoveNode(node *Node, currTime time.Time) error {
	if err := us.db.removeNode(node, currTime); err != nil {
		return errTechnicalError
	}
	return nil
}

func checkIsKeyWrong(key string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Error("Panic during node validation ", r)
			err = errWrongDecodedLength
		}
	}()
	if _, err := cipher.PubKeyFromHex(key); err != nil {
		log.Error("Error during checking for validity of key", key)
		return err
	}
	return
}

// TODO remove this after fix
func (us *Service) GetApplicationApprovedNodesCount(appID uint) (numOfApprovedNodes int) {
	return us.db.getApplicationApprovedNodesCount(appID)
}

// GetAllApplicationsForMiner - returns all applications for miner in case miner was transferred
// This is made for purpose of filling new table miner_applications and can be removed after migration is done
func (us *Service) GetAllApplicationsForMiner(minerID uint) (applications []*Application, err error) {
	appIDs, err := us.db.findAllApplicationIDsForMiner(minerID)
	if err != nil {
		return
	}

	for _, appID := range appIDs {
		app, err := us.db.getWhitelistWithoutPreload(fmt.Sprint(appID))
		if err != nil {
			log.Errorf("Unable to find application with id: %v - %v", appID, err)
			continue
		}
		applications = append(applications, &app)
	}

	return
}
func (us *Service) checkIfHashIsUnique(hashed string, minId uint) (isUnique bool, err error) {
	_, err = us.db.findImagesByHash(hashed, minId)
	if err != nil && err == errCannotImageByHash {
		isUnique = true
		err = nil
		return
	}
	return

}
func (us *Service) CheckForDuplicates(appReq *ApplicationReq, newApp bool, username string) ([]int, error) {
	var hash [32]byte
	var IDs []int
	var minerID uint
	if !newApp {
		app, err := us.getActiveApplication(username)
		if err != nil {
			log.Error("Error while retrieving application during image validation: ", err)
			return IDs, err
		}
		minerID = app.Miner.ID
	}

	for i, image := range appReq.Files {
		hash = sha256.Sum256(image.File)
		hashedString := hex.EncodeToString(hash[:])
		unique, err := us.checkIfHashIsUnique(hashedString, minerID)
		appReq.Files[i].Hashed = hashedString
		if err != nil {
			log.Errorf("Unable to check if picture with position %v is unique due to error %v", i, err)
			continue
		}
		if !unique {
			log.Debugf("Picture number %v is duplicate", i+1)
			IDs = append(IDs, i+1)

		}

	}
	return IDs, nil
}
func (us *Service) GetAllImages() ([]Image, error) {
	list, err := us.db.getAllImages()
	if err != nil {
		return nil, errCannotLoadImages
	}
	return list, nil
}

func (us *Service) GetImagesForUser(username string) ([]Image, error) {
	list, err := us.db.getImageRecordsForUser(username)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (us *Service) AddHashToImages(imgID uint, hashed string) error {
	log.Debugf("Adding %v hash to image %v", hashed, imgID)
	return us.db.updateHash(imgID, hashed)
}

func checkIsImageNameTaken(images []Image, imgName string) bool {
	for _, image := range images {
		if image.Path == imgName {
			return true
		}
	}
	return false
}
