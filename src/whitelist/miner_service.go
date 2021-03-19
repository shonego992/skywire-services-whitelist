package whitelist

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	autoDeclineMessage     = "AUTOMATIC UPDATE - application is declined due to User's request to append new nodes to existing miner"
	autoUpdateMessage      = "AUTOMATIC UPDATE - user requested new nodes for the existing miner"
	autoUpdateMessageUser  = "You requested new nodes for the existing miner"
	autoPendingMessage     = "AUTOMATIC UPDATE - application returned back to PENDING state"
	autoPendingMessageUser = "Application returned back to PENDING state after you requested new nodes to be added to the existing miner"
)

// MinerService provides access to User related data
type MinerService struct {
	db minerStore
}

// DefaultMinerService prepares new instance of MinerService
func DefaultMinerService() MinerService {
	return NewMinerService(DefaultMinerData())
}

// NewMinerService prepares new instance of MinerService
func NewMinerService(whitelistStore minerStore) MinerService {
	return MinerService{
		db: whitelistStore,
	}
}

func (us *MinerService) getShopInfo(id uint64) (ShopData, error) {
	return us.db.getShopInfo(id)
}

func (us *MinerService) updateStoreInfo(order MinerShopOrder, dbRecord ShopData) error {
	if order.FinancialStatus != "paid" {
		return nil
	}
	if dbRecord.ID == 0 {
		shopInfo := ShopData{ID: order.Id, Status: order.FinancialStatus}
		err := us.db.createShopInfo(&shopInfo)
		if err != nil {
			return err
		}
	} else {
		dbRecord.Status = order.FinancialStatus
		err := us.db.updateShopInfo(&dbRecord)
		if err != nil {
			return err
		}
	}
	return nil
}

func getShopifyData() (*MinerShopResponse, error) {
	var shopifyUrl = viper.GetString("shopify.shop-url")
	var shopifyUser = viper.GetString("shopify.shop-user")
	var shopifyPassword = viper.GetString("shopify.shop-password")
	apiString := fmt.Sprintf("https://%v:%v@%v", shopifyUser, shopifyPassword, shopifyUrl)
	response, err := http.Get(apiString)
	if err != nil {
		return nil, errCannotGetShopifyData
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errCannotGetShopifyData
	}

	shopRecords, err := extractShopRecordsFromURL(contents)
	if err != nil {
		return nil, err
	}
	return shopRecords, nil
}

func extractShopRecordsFromURL(body []byte) (*MinerShopResponse, error) {
	var s = new(MinerShopResponse)
	err := json.Unmarshal(body, &s)
	if err != nil {
		return nil, errCannotGetShopifyData
	}
	return s, nil
}

func (us *MinerService) addMinerToUser(username string, nodes []string) error {
	m := Miner{}
	m.Username = username
	m.Type = DIY
	for _, node := range nodes {
		m.Nodes = append(m.Nodes, Node{Key: node})
	}
	err := us.db.createMiner(&m)
	if err != nil {
		return errCannotCreateMiners
	}
	return nil
}

func (us *MinerService) RemoveMiner(id string, currTime time.Time) error {
	if err := us.db.removeMiner(id, currTime); err != nil {
		return errTechnicalError
	}

	return nil
}

//ActivateMiner reenables the miner, using id provided from the controller
func (us *MinerService) ActivateMiner(id string) error {
	if err := us.db.activateMiner(id); err != nil {
		return errUnableToSave
	}
	return nil
}
func (us *MinerService) ReactivateNodes(nodes []Node) error {

	for i := 0; i < len(nodes); i++ {
		free, err := us.CheckIsNodeFree(nodes[i].Key)
		if err != nil {
			log.Errorf("Unable to check node with key %v due to error: %v ", nodes[i].Key, err)
		} else if free {
			err := us.db.activateNode(nodes[i].ID)
			if err == nil {
				log.Debug("Successfuly reactivated node with key ", nodes[i].Key)

			} else {
				log.Errorf("Unable to reactivate node with key %v due to error: %v ", nodes[i].Key, err)
			}

		}
	}
	return nil
}
func (us *MinerService) CheckIsNodeFree(nodeKey string) (bool, error) {
	_, err := us.db.findActiveNodeByKey(nodeKey)
	if err != nil && err == errCannotFindActiveNode {
		return true, nil
	}
	return false, err

}

func (us *MinerService) transferMiner(request transferMinerReq, currentUser string) (err error) {
	var miner Miner
	if miner, err = us.db.findSpecificMiner(request.MinerId); err != nil {
		return
	}

	//TODO allow admins to perform this action
	if miner.Username != currentUser {
		log.Errorf("User %v not owner of miner %v ", currentUser, miner.ID)
		return errNotOwnerOfMiner
	}

	miner.ApplicationID = 0 // reset application id

	miner.Username = NormalizeMail(request.TransferTo)
	var nodeIds []uint
	for _, node := range miner.Nodes {
		nodeIds = append(nodeIds, node.ID)
	}
	err = us.db.updateMiner(&miner, nodeIds)
	if err != nil {
		log.Errorf("Error transferring miner to %v - %v", request.TransferTo, err)
		return errCannotTransferMiner
	}

	transferRecord := MinerTransfer{MinerID: miner.ID, OldUsername: currentUser, NewUsername: request.TransferTo}
	err = us.db.createTransferMinerRecord(transferRecord)

	return
}

func (us *MinerService) exportMiners(request exportMinersReq, startDate time.Time, endDate time.Time) ([]Miner, error) {
	useFilter := false
	if startDate.Before(endDate) {
		useFilter = true
	}
	miners, err := us.db.exportMiners(request, startDate, endDate, useFilter)

	if err != nil {
		return nil, ErrCannotLoadMiners
	}
	return miners, nil
}

func (us *MinerService) getUserMiners(username string) ([]Miner, error) {
	miners, err := us.db.findMiners(username)
	if err != nil {
		return nil, ErrCannotLoadMiners
	}
	return miners, nil
}

func (us *MinerService) GetMinersForUser(username string) ([]Miner, error) {
	miners, err := us.db.findMiners(username)
	if err != nil {
		return nil, err
	}
	return miners, nil

}

func (us *MinerService) GetAllMiners() ([]Miner, error) {
	miners, err := us.db.getAllMiners()
	if err != nil {
		return nil, ErrCannotLoadMiners
	}
	return miners, nil
}

func (us *MinerService) getSpecificMiner(id string, username string) (Miner, error) {
	miner, err := us.db.findSpecificMiner(id)
	if err != nil {
		return Miner{}, errCannotFindMiner
	}
	if miner.Username != username {
		return Miner{}, errMinerNotFoundForUser
	}
	return miner, nil
}
func (us *MinerService) getSpecificDisabledMiner(id string, username string, deletionTime *time.Time) (Miner, error) {
	miner, err := us.db.findSpecificDisabledMiner(id, deletionTime)
	if err != nil {
		return Miner{}, errCannotFindMiner
	}
	if miner.Username != username {
		return Miner{}, errMinerNotFoundForUser
	}
	return miner, nil
}
func (us *MinerService) getSpecificMinerWithApplications(id string, username string) (Miner, error) {
	miner, err := us.db.findSpecificMinerWithApplications(id)
	if err != nil {
		return Miner{}, errCannotFindMiner
	}
	if miner.Username != username {
		return Miner{}, errMinerNotFoundForUser
	}
	return miner, nil
}

func (us *MinerService) getSpecificMinerForAdmin(id string) (Miner, error) {
	miner, err := us.db.findSpecificMiner(id)
	if err != nil {
		return Miner{}, err
	}
	return miner, nil
}
func (us *MinerService) getSpecificDisabledMinerForAdmin(id string, deletionTime *time.Time) (Miner, error) {
	miner, err := us.db.findSpecificDisabledMiner(id, deletionTime)
	if err != nil {
		return Miner{}, errCannotFindMiner
	}
	return miner, nil
}

func (us *MinerService) UpdateMiner(miner *Miner, nodes []*Node) error {
	var updatedIDs []uint
	for _, node := range nodes {
		foundNode := false
		if node.ID > 0 {
			for _, oldNode := range miner.Nodes {
				if miner.Type == OFFICIAL || (oldNode.ID == node.ID && oldNode.Key == node.Key) {
					updatedIDs = append(updatedIDs, node.ID)
					foundNode = true
					break
				}
			}
		}

		if !foundNode {
			appendingNode := Node{Key: node.Key, MinerID: miner.ID}
			if node.ID > 0 {
				appendingNode.ID = node.ID
				appendingNode.CreatedAt = node.CreatedAt

			}
			miner.Nodes = append(miner.Nodes, appendingNode)
		}
	}
	if err := us.db.updateMiner(miner, updatedIDs); err != nil {
		return errCannotUpdateMiner
	}

	return nil
}

func (us *MinerService) InsertCreatedAtForOldMiners() error {
	if err := us.db.insertCreatedAtForOldMiners("2019-01-12 15:30:21.130612"); err != nil {
		return err
	}
	return nil
}

func (us *MinerService) CreateMinerEntitiesForUser(change *ChangeHistory, app Application) (m Miner, err error) {
	m.Username = app.Username
	m.Type = DIY
	for _, node := range change.Nodes {
		m.Nodes = append(m.Nodes, *node)
	}
	m.ApplicationID = app.ID
	m.ApprovedNodesCount = len(m.Nodes)
	m.Applications = append(m.Applications, &app)
	err = us.db.createMiner(&m)
	if err != nil {
		err = errCannotCreateMiners
		return
	}

	return
}

func (us *MinerService) AddImagesToMiner(minerID uint, images []*Image) error {
	log.Debugf("Adding %v images to miner %v", len(images), minerID)
	var imagesExtracted []uint
	for _, img := range images {
		imagesExtracted = append(imagesExtracted, img.ID)
	}
	return us.db.addImagesToMiner(minerID, imagesExtracted)
}

//Appends uptime to nodes from uptime service
//return errNoKeysToGetDataFor if all the keys are empty
func (us *MinerService) GetUptimeNoDate(nodes []Node, forExport bool) ([]NodeUptimeResponse, error) {
	endpoint := viper.GetString("uptime.uptime-getuptime-endpoint")
	if forExport {
		endpoint = viper.GetString("uptime.uptime-getuptime-export")
	}
	var apiString = viper.GetString("uptime.uptime-service") + endpoint
	var results []NodeUptimeResponse
	var stringForRequest string

	for ind, node := range nodes {
		if node.Key != "" {
			stringForRequest += strings.Trim(node.Key, " ")
			if ind == len(nodes)-1 {
				continue
			}
			stringForRequest += ","
		}
	}
	//Making a call to uptime service api
	if stringForRequest == "" {
		return results, errNoKeysToGetDataFor
	}
	log.Debug("Calling uptime service for nodes: ", stringForRequest)
	response, err := http.Get(apiString + stringForRequest)
	if err != nil {
		log.Error("Error while fetching data from uptime service. Error: ", err)
		return results, errCannotGetUpTimeData
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error("Error while reading data from uptime service. Error: ", err)
		return results, errCannotGetUpTimeData
	}
	err = json.Unmarshal(body, &results)
	if err != nil {
		bodyString := string(body)
		log.Errorf("Error while unmarshalling data from uptime service.\nBody:\n%v\nError:\n%v ", bodyString, err)
		return nil, errCannotGetUpTimeData
	}
	log.Debugf("Successfully  called uptime service and received answer %+v", results)
	return results, nil
}

//Appends uptime to nodes from uptime service
//return errNoKeysToGetDataFor if all the keys are empty
func (us *MinerService) GetUptime(nodes []Node, forExport bool, startDate int64, endDate int64) ([]NodeUptimeResponse, error) {
	log.Debug("Calling uptime service for dates: ", startDate, endDate)
	startDateString := strconv.FormatInt(startDate, 10)
	endDateString := strconv.FormatInt(endDate, 10)
	endpoint := viper.GetString("uptime.uptime-getuptime-endpoint")
	if forExport {
		endpoint = viper.GetString("uptime.uptime-getuptime-export")
	}
	var apiString = viper.GetString("uptime.uptime-service") + endpoint
	var results []NodeUptimeResponse
	var stringForRequest string

	for ind, node := range nodes {
		if node.Key != "" {
			stringForRequest += strings.Trim(node.Key, " ")
			if ind == len(nodes)-1 {
				continue
			}
			stringForRequest += ","
		}
	}
	//Making a call to uptime service api
	if stringForRequest == "" {
		return results, errNoKeysToGetDataFor
	}
	log.Debug("Calling uptime service for nodes: ", stringForRequest)
	var requestUrl string
	requestUrl = apiString + stringForRequest
	if forExport {
		requestUrl += "&startDate=" + startDateString + "&endDate=" + endDateString
	}
	response, err := http.Get(requestUrl)
	if err != nil {
		log.Error("Error while fetching data from uptime service. Error: ", err)
		return results, errCannotGetUpTimeData
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error("Error while reading data from uptime service. Error: ", err)
		return results, errCannotGetUpTimeData
	}
	log.Debugf("Successfully  called uptime service and received answer %+v", results)
	err = json.Unmarshal(body, &results)
	if err != nil {
		bodyString := string(body)
		log.Errorf("Error while unmarshalling data from uptime service.\nBody:\n%v\nError:\n%v ", bodyString, err)
		return nil, errCannotGetUpTimeData
	}
	return results, nil
}

func (us *MinerService) GetAllUptimeRecords(startDate int64, endDate int64) ([]NodeUptimeResponse, error) {
	endpoint := viper.GetString("uptime.uptime-get-all-uptimes-endpoint")
	var apiString = viper.GetString("uptime.uptime-service") + endpoint
	var results []NodeUptimeResponse
	var stringForRequest string

	//Making a call to uptime service api
	log.Debug("Calling uptime service for dates: ", startDate, endDate)
	startDateString := strconv.FormatInt(startDate, 10)
	endDateString := strconv.FormatInt(endDate, 10)
	response, err := http.Get(apiString + stringForRequest + "?startDate=" + startDateString + "&endDate=" + endDateString)
	if err != nil {
		log.Error("Error while fetching data from uptime service. Error: ", err)
		return results, errCannotGetUpTimeData
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error("Error while reading data from uptime service. Error: ", err)
		return results, errCannotGetUpTimeData
	}
	log.Debugf("Successfully  called uptime service and received answer %+v", results)
	err = json.Unmarshal(body, &results)
	if err != nil {
		bodyString := string(body)
		log.Errorf("Error while unmarshalling data from uptime service.\nBody:\n%v\nError:\n%v ", bodyString, err)
		return nil, errCannotGetUpTimeData
	}
	return results, nil
}

func (us *MinerService) GetUserByNodeKey(nodeKey string) (Miner, error) {
	return us.db.findMinerByNodeKey(nodeKey)
}

func (us *MinerService) GetImportData() ([]MinerImport, error) {
	return us.db.findImportData()
}

func (us *MinerService) SaveImportData(importedData []MinerImport) ([]MinerImport, error) {
	return us.db.saveImportData(importedData)
}

func (us *MinerService) ImportMiners(importedData MinerImport) error {
	if err := us.createMinerEntitiesFromImport(importedData); err != nil {
		return err
	}
	if err := us.db.removeImportData(importedData); err != nil {
		return err
	}

	return nil
}

func (us *MinerService) ImportMinersNoRemoval(importedData MinerImport) error {
	if err := us.createMinerEntitiesFromImport(importedData); err != nil {
		return err
	}

	return nil
}

func (us *MinerService) createMinerEntitiesFromImport(importedRecord MinerImport) error {
	keyToBeUsed := 0
	for i := 0; i < importedRecord.NumberOfMiners; i++ {
		miner := &Miner{
			Username:   importedRecord.Username,
			BatchLabel: importedRecord.BatchLabel,
			Gifted:     importedRecord.Gifted,
			Type:       OFFICIAL,
		}

		for i := 0; i < 8; i++ {
			if keyToBeUsed < len(importedRecord.Nodes) {
				miner.Nodes = append(miner.Nodes, importedRecord.Nodes[keyToBeUsed])
				keyToBeUsed++
			} else {
				miner.Nodes = append(miner.Nodes, Node{Key: ""})
			}
		}

		if err := us.db.createMiner(miner); err != nil {
			return err
		}
	}
	return nil
}

type MinerShopOrder struct {
	Id              uint64
	Email           string
	ProcessedAt     string               `json:"processed_at"`
	LineItems       []MinerShopLineItems `json:"line_items"`
	FinancialStatus string               `json:"financial_status"`
}

type MinerShopLineItems struct {
	Sku      string
	Quantity int
}

type MinerShopResponse struct {
	Orders []MinerShopOrder
}
