package whitelist

import (
	"time"

	"github.com/SkycoinPro/skywire-services-whitelist/src/database/postgres"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

// store is user related interface for dealing with database operations
type store interface {
	createApplication(application *Application) error
	createChangeHistory(change *ChangeHistory) error
	getUserApplicationWithHistory(applicationId uint) (Application, error)
	updateChangeHistory(change *ChangeHistory) error
	getWhitelist(id string) (Application, error)
	getWhitelistWithoutPreload(id string) (Application, error)
	updateApplication(application *Application, oldChangeID uint) error
	findApplications(email string) ([]Application, error)
	findPendingApplication(email string) (Application, error)
	findActiveNodeByKey(key string) (Node, error)
	getApplicationForNode(node Node) (Application, error)
	appUpdateCreatedAt(app *Application, createdAtImport string) error
	getAllWhitelists() ([]Application, error)
	getApplicationApprovedNodesCount(appID uint) (numOfApprovedNodes int)
	removeNode(node *Node, currTime time.Time) error
	findAllApplicationIDsForMiner(minerID uint) ([]uint, error)
	getMinerForApplication(application *Application) Miner
	findImagesByHash(hashed string, minId uint) (image Image, err error)
	getAllImages() ([]Image, error)
	getImageRecordsForUser(username string) ([]Image, error)
	updateHash(imgID uint, hashed string) error
}

// data implements store interface which uses GORM library
type data struct {
	db *gorm.DB
}

func DefaultData() data {
	return NewData(postgres.DB)
}

func NewData(database *gorm.DB) data {
	return data{
		db: database,
	}
}

const adminStatusStart uint8 = 16

func (u data) getWhitelist(id string) (Application, error) {
	var application Application
	record := u.db.Preload("Miner").Preload("Miner.Nodes").Preload("ChangeHistory", func(db *gorm.DB) *gorm.DB {
		return db.Order("change_history.id ASC")
	}).Preload("ChangeHistory.Nodes", func(db *gorm.DB) *gorm.DB {
		return db.Unscoped()
	}).Preload("ChangeHistory.Images").Where("id = ?", id).Find(&application)
	if record.RecordNotFound() {
		return Application{}, errCannotLoadWhitelists
	}
	if errs := record.GetErrors(); len(errs) > 0 {
		for err := range errs {
			log.Errorf("Error occurred while fetching whitelist by id %v - %v", id, err)
		}
		return Application{}, errUnableToRead
	}
	return application, nil
}

func (u data) getWhitelistWithoutPreload(id string) (application Application, err error) {
	record := u.db.Find(&application, "id = ?", id)

	if record.RecordNotFound() {
		err = errCannotLoadWhitelists
		return
	}

	if errs := record.GetErrors(); len(errs) > 0 {
		for _, err = range errs {
			log.Errorf("Error occurred while fetching whitelist by id %v - %v", id, err)
		}
		err = errUnableToRead
		return
	}

	return
}

func (u data) updateApplication(application *Application, oldChangeID uint) error {
	db := u.db.Begin()
	var dbError error
	for _, err := range db.Save(&application).GetErrors() {
		dbError = err
		log.Errorf("Error while updating application %v in DB due to error %v", application, err)
	}
	if dbError != nil {
		db.Rollback()
		return dbError
	}
	db.Commit()

	//TODO drop this workaround for miner_id NULL value in DB once proper solution is found
	for _, err := range u.db.Exec("UPDATE nodes SET miner_id=null WHERE miner_id = 0").GetErrors() {
		dbError = err
		log.Error("Error while cleaning up miner IDs for the new change history in DB", err)
	}

	newChangeID := application.GetLatestChangeHistory().ID

	// workaround for Nodes deleted during whitelisting process (they get removed from pivot table but don't get deleted_at attribute set)
	if oldChangeID > 0 && newChangeID > 0 && oldChangeID != newChangeID {
		for _, err := range u.db.Exec("UPDATE nodes SET deleted_at=? WHERE nodes.id IN (SELECT node_id FROM change_nodes WHERE change_history_id = ?) AND nodes.id NOT IN (SELECT node_id FROM change_nodes WHERE change_history_id = ?)", time.Now(), oldChangeID, newChangeID).GetErrors() {
			dbError = err
			log.Error("Error while cleaning up miner IDs for the new change history in DB", err)
		}
	}

	return dbError
}

func (u data) createApplication(application *Application) error {
	db := u.db.Begin()
	var dbError error
	for _, err := range db.Create(application).GetErrors() {
		dbError = err
		log.Error("Error while creating new application in DB ", err)
	}
	if dbError != nil {
		db.Rollback()
		return dbError
	}
	db.Commit()

	return nil
}

func (u data) createChangeHistory(change *ChangeHistory) error {
	db := u.db.Begin()
	var dbError error
	for _, err := range db.Create(change).GetErrors() {
		dbError = err
		log.Error("Failed while persisting nodes:")
		for _, n := range change.Nodes {
			log.Error(n.Key)
		}
		log.Errorf("Error while creating new change %v for application in DB due to error %v", change, err)
	}
	if dbError != nil {
		db.Rollback()
		return dbError
	}
	db.Commit()

	return nil
}

func (u data) getUserApplicationWithHistory(applicationId uint) (Application, error) {
	var application Application
	record := u.db.Preload("ChangeHistory", func(db *gorm.DB) *gorm.DB {
		return db.Order("change_history.id ASC")
	}).Preload("ChangeHistory.Nodes").Preload("ChangeHistory.Images").Where("id = ?", applicationId).Find(&application)
	if record.RecordNotFound() {
		return Application{}, ErrCannotFindUser
	}
	if errs := record.GetErrors(); len(errs) > 0 {
		for err := range errs {
			log.Errorf("Error occurred while fetching application by id %v - %v", applicationId, err)
		}
		return Application{}, errUnableToRead
	}
	return application, nil
}

func (u data) updateChangeHistory(change *ChangeHistory) error {
	db := u.db.Begin()
	var dbError error
	for _, err := range db.Save(&change).GetErrors() {
		dbError = err
		log.Error("Error while updating change history in DB", err)
	}
	if dbError != nil {
		db.Rollback()
		return dbError
	}
	db.Commit()

	return nil
}

func (u data) findApplications(email string) ([]Application, error) {
	var applications []Application
	record := u.db.Order("id").Where("username = ?", NormalizeMail(email)).Preload("ChangeHistory", func(db *gorm.DB) *gorm.DB {
		return db.Order("change_history.id ASC")
	}).Preload("ChangeHistory.Nodes").Preload("ChangeHistory.Images").Find(&applications)
	if record.RecordNotFound() { //|| len(applications) == 0
		return nil, ErrCannotFindUser
	}
	if errs := record.GetErrors(); len(errs) > 0 {
		for err := range errs {
			log.Errorf("Error occurred while fetching applications by user email %v - %v", email, err)
		}
		return nil, errUnableToRead
	}

	return applications, nil
}

//TODO (improvement point) consider adding parameter for status and then find applications for user based on status,
//not hardcoding the query only to pending
func (u data) findPendingApplication(email string) (application Application, err error) {
	record := u.db.Preload("ChangeHistory", func(db *gorm.DB) *gorm.DB {
		return db.Order("change_history.id ASC")
	}).Preload("ChangeHistory.Nodes").Preload("ChangeHistory.Images").Where("username = ? and current_status = 0", NormalizeMail(email)).Find(&application)
	if record.RecordNotFound() {
		err = errCannotFindPendingApp
		return
	}
	if errs := record.GetErrors(); len(errs) > 0 {
		for err := range errs {
			log.Errorf("Error occurred while fetching pending application by user email %v - %v", email, err)
		}
		err = errUnableToRead
		return
	}

	return
}

func (u data) findActiveNodeByKey(key string) (node Node, err error) {
	record := u.db.Where("key = ?", key).Find(&node)
	if record.RecordNotFound() {
		err = errCannotFindActiveNode
		return
	}

	if errs := record.GetErrors(); len(errs) > 0 {
		for err := range errs {
			log.Errorf("Error occurred while fetching node by key %v - %v", key, err)
		}
		err = errUnableToRead
		return
	}

	return
}

func (u data) getApplicationForNode(node Node) (application Application, err error) {
	record := u.db.Joins("JOIN change_history ch ON applications.id = ch.application_id").Joins("JOIN change_nodes cn on ch.id = cn.change_history_id").Where("cn.node_id = ?", node.ID).Preload("ChangeHistory", func(db *gorm.DB) *gorm.DB {
		return db.Order("change_history.id ASC")
	}).Last(&application)
	if record.RecordNotFound() {
		err = errNoApplicationForNode
	}

	if errs := record.GetErrors(); len(errs) > 0 {
		for err := range errs {
			log.Errorf("Error occurred while fetching application by node with id %v - %v", node.ID, err)
		}
		err = errUnableToRead
		return
	}

	return
}

func (u data) appUpdateCreatedAt(app *Application, createdAtImport string) error {
	db := u.db.Begin()
	var dbError error
	for _, err := range db.Exec("UPDATE applications SET created_at=? WHERE id=?", createdAtImport, &app.ID).GetErrors() {
		dbError = err
		log.Error("Error while updating application created_at time", err)
	}
	if dbError != nil {
		db.Rollback()
		return dbError
	}
	db.Commit()
	return nil
}

//fetches all whitelists including deleted ones
func (u data) getAllWhitelists() ([]Application, error) {
	var applications []Application
	record := u.db.Unscoped().Find(&applications)
	if record.RecordNotFound() {
		return nil, errCannotLoadWhitelists
	}
	if errs := record.GetErrors(); len(errs) > 0 {
		for err := range errs {
			log.Error("Error while fetching applications ", err)
		}
		return nil, errUnableToRead
	}
	return applications, nil

}

func (u data) getApplicationApprovedNodesCount(appID uint) (numOfApprovedNodes int) {
	u.db.Table("change_nodes").Where("change_history_id IN (SELECT id FROM change_history WHERE application_id = ? and status = 1)", appID).Count(&numOfApprovedNodes)
	return
}

func (u data) removeNode(node *Node, currTime time.Time) error {
	db := u.db.Begin()
	var dbError error
	for _, err := range db.Exec("UPDATE nodes SET deleted_at=? WHERE id=?", currTime, &node.ID).GetErrors() {
		dbError = err
		log.Errorf("Error while removing node with id: %v  due to %v ", node.ID, err)
	}
	if dbError != nil {
		db.Rollback()
		return dbError
	}
	db.Commit()
	return nil
}

func (u data) findAllApplicationIDsForMiner(minerID uint) (applicationIDs []uint, err error) {

	record := u.db.Raw(`SELECT DISTINCT (ap.id) from applications ap
		JOIN change_history ch ON ap.id = ch.application_id 
		JOIN change_nodes cn on ch.id = cn.change_history_id 
		WHERE cn.node_id IN (select nodes.id from nodes where miner_id = ?)`, minerID)

	if record.RecordNotFound() {
		return
	}

	if errs := record.GetErrors(); len(errs) > 0 {
		for err := range errs {
			log.Errorf("Error occurred while fetching applications by miner with id %v - %v", minerID, err)
		}
		err = errUnableToRead
		return
	}

	rows, err := record.Rows()
	defer rows.Close()

	if err != nil {
		log.Error(err)
	}

	for rows.Next() {
		var appID uint
		rows.Scan(&appID)
		applicationIDs = append(applicationIDs, appID)
	}

	return
}

func (u data) getMinerForApplication(application *Application) (miner Miner) {
	u.db.Model(application).Preload("Nodes").Preload("Images").Related(&miner)
	return
}
func (u data) findImagesByHash(hashed string, minId uint) (image Image, err error) {
	record := u.db.Where("img_hash = ? and miner_id <> ?", hashed, minId).Find(&image)
	if record.RecordNotFound() {
		err = errCannotImageByHash
		return
	}

	if errs := record.GetErrors(); len(errs) > 0 {
		for err := range errs {
			log.Errorf("Error occurred while fetching image by hash %v - %v", hashed, err)
		}
		err = errUnableToRead
		return
	}

	return
}

//fetches all images
func (u data) getAllImages() ([]Image, error) {
	var images []Image
	record := u.db.Find(&images)
	if record.RecordNotFound() {
		return nil, errCannotLoadImages
	}
	if errs := record.GetErrors(); len(errs) > 0 {
		for err := range errs {
			log.Error("Error while fetching images ", err)
		}
		return nil, errUnableToRead
	}
	return images, nil

}

//fetches all images
func (u data) getImageRecordsForUser(username string) ([]Image, error) {
	var images []Image
	record := u.db.Where("path like ?", username+"%").Find(&images)
	if record.RecordNotFound() {
		return nil, errCannotLoadImagesForUser
	}
	if errs := record.GetErrors(); len(errs) > 0 {
		for err := range errs {
			log.Errorf("Error while fetching images for user %v  ", err)
		}
		return nil, errUnableToRead
	}
	return images, nil

}

func (u data) updateHash(imgID uint, hashed string) error {
	db := u.db.Begin()
	var dbError error
	for _, err := range db.Exec("UPDATE images SET img_hash=? WHERE id=?", hashed, imgID).GetErrors() {
		dbError = err
		log.Error("Error while updating images in DB", err)
		break
	}

	if dbError != nil {
		db.Rollback()
		return dbError
	}
	db.Commit()

	return nil
}
