package whitelist

import (
	"time"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"github.com/SkycoinPro/skywire-services-whitelist/src/database/postgres"
)

type minerStore interface {
	createMiner(miner *Miner) error
	findMiners(email string) ([]Miner, error)

	findSpecificMiner(id string) (Miner, error)
	findSpecificDisabledMiner(id string, deletionTime *time.Time) (Miner, error)
	findSpecificMinerWithApplications(id string) (Miner, error)
	updateMiner(miner *Miner, updatedIDs []uint) error
	addImagesToMiner(minerID uint, images []uint) error
	findActiveNodeByKey(key string) (Node, error)

	findImportData() ([]MinerImport, error)
	saveImportData(importedData []MinerImport) ([]MinerImport, error)
	removeImportData(record MinerImport) error
	exportMiners(request exportMinersReq, startDate time.Time, endDate time.Time, useFilter bool) ([]Miner, error)
	getShopInfo(id uint64) (ShopData, error)
	createShopInfo(shopData *ShopData) error
	updateShopInfo(shopData *ShopData) error
	createTransferMinerRecord(transferRecord MinerTransfer) error
	getAllMiners() ([]Miner, error)
	findMinerByNodeKey(nodeKey string) (Miner, error)
	removeMiner(id string, currTime time.Time) error
	activateMiner(minerId string) error
	activateNode(nodeID uint) error
	insertCreatedAtForOldMiners(dateFrom string) error
}

type minerData struct {
	db *gorm.DB
}

func DefaultMinerData() minerData {
	return NewMinerData(postgres.DB)
}

func NewMinerData(database *gorm.DB) minerData {
	return minerData{
		db: database,
	}
}

func (u minerData) removeMiner(id string, currTime time.Time) error {
	db := u.db.Begin()
	var dbError error
	for _, err := range db.Exec("UPDATE miners SET deleted_at=? WHERE id=?", currTime, id).GetErrors() {
		dbError = err
		log.Errorf("Error while removing miner with id: %v  due to %v ", id, err)
	}
	if dbError != nil {
		db.Rollback()
		return dbError
	}
	db.Commit()
	return nil
}

//activateMiner reenables disabled (soft deleted) miner by setting via raw query deleted_at field to NULL and update_at to current time
func (u minerData) activateMiner(minerID string) error {
	db := u.db.Begin()
	var dbError error
	for _, err := range db.Unscoped().Exec("UPDATE miners SET deleted_at=NULL, updated_at= ? WHERE id = ?", time.Now(), minerID).GetErrors() {
		dbError = err
		log.Errorf("Error while reenabling miner %v - %v", minerID, err)
		break
	}

	if dbError != nil {
		db.Rollback()
		return dbError
	}
	db.Commit()
	return nil
}

//activateNode reenables disabled node  by setting via raw query deleted_at field to NULL and update_at to current time
func (u minerData) activateNode(nodeID uint) error {
	db := u.db.Begin()
	var dbError error
	for _, err := range db.Unscoped().Exec("UPDATE nodes SET deleted_at=NULL, updated_at= ? WHERE id = ?", time.Now(), nodeID).GetErrors() {
		dbError = err
		log.Errorf("Error while reactivating node %v - %v", nodeID, err)
		break
	}

	if dbError != nil {
		db.Rollback()
		return dbError
	}
	db.Commit()
	return nil
}

func (u minerData) findMinerByNodeKey(nodeKey string) (Miner, error) {
	var node Node
	record := u.db.Where("key = ?", nodeKey).First(&node)
	if record.RecordNotFound() {
		return Miner{}, errCannotFindActiveNode
	}
	if errs := record.GetErrors(); len(errs) > 0 {
		for err := range errs {
			log.Errorf("Error fetching node by key %v, %v ", nodeKey, err)
		}
		return Miner{}, errCannotFindActiveNode
	}

	if node.MinerID <= 0 {
		return Miner{}, nil
	}
	var miner Miner
	record = u.db.Where("id = ?", node.MinerID).Find(&miner)
	if record.RecordNotFound() {
		return Miner{}, errCannotFindMiner
	}
	if errs := record.GetErrors(); len(errs) > 0 {
		for err := range errs {
			log.Errorf("Error occurred while fetching miner id %v - %v", node.MinerID, err)
		}
		return Miner{}, errCannotFindMiner
	}
	return miner, nil
}

func (u minerData) getShopInfo(id uint64) (ShopData, error) {
	var shopData ShopData
	record := u.db.Where("id = ?", id).Find(&shopData)
	if record.RecordNotFound() {
		return ShopData{}, errCannotFindShopRecord
	}
	if errs := record.GetErrors(); len(errs) > 0 {
		for err := range errs {
			log.Errorf("Error occurred while fetching shop record by id %v - %v", id, err)
		}
		return ShopData{}, errUnableToRead
	}
	return shopData, nil
}

func (u minerData) createShopInfo(shopData *ShopData) error {
	db := u.db.Begin()
	var dbError error
	for _, err := range db.Create(shopData).GetErrors() {
		dbError = err
		log.Error("Error while persisting new shop data in DB", err)
	}
	if dbError != nil {
		db.Rollback()
		return dbError
	}
	db.Commit()

	return nil
}

func (u minerData) updateShopInfo(shopData *ShopData) error {
	db := u.db.Begin()
	var dbError error
	for _, err := range db.Save(&shopData).GetErrors() {
		dbError = err
		log.Error("Error while updating shop data in DB", err)
	}
	if dbError != nil {
		db.Rollback()
		return dbError
	}
	db.Commit()

	return nil
}

func (u minerData) createTransferMinerRecord(transferRecord MinerTransfer) error {
	db := u.db.Begin()
	var dbError error
	for _, err := range db.Save(&transferRecord).GetErrors() {
		dbError = err
		log.Error("Error while creating miner transfer record in DB", err)
	}
	if dbError != nil {
		db.Rollback()
		return dbError
	}
	db.Commit()

	return nil
}

func (u minerData) exportMiners(request exportMinersReq, startDate time.Time, endDate time.Time, useFilter bool) ([]Miner, error) {
	var miners []Miner
	var record *gorm.DB
	if useFilter {
		record = u.db.Where("created_at BETWEEN ? AND ?", startDate, endDate).Preload("Nodes").Order("username").Find(&miners)
	} else {
		record = u.db.Preload("Nodes").Order("username").Find(&miners)
	}
	if record.RecordNotFound() {
		return nil, ErrCannotFindUser
	}
	if errs := record.GetErrors(); len(errs) > 0 {
		for err := range errs {
			log.Errorf("Error occurred while fetching miners- %v", err)
		}
		return nil, errUnableToRead
	}
	return miners, nil
}

func (u minerData) insertCreatedAtForOldMiners(dateFrom string) error {
	db := u.db.Begin()
	defaultData := DefaultData()
	var miners []Miner
	var applicationIDs []uint
	var dbError error
	records := u.db.Raw("SELECT * FROM miners WHERE created_at < ? AND type = ?", dateFrom, 1).Scan(&miners)
	if records.RecordNotFound() {
		dbError = errUnableToRead
		log.Errorf("Cant find miners made before this date %v", dateFrom)
		return errUnableToRead
	}
	for _, miner := range miners {

		applicationIDs = []uint{}
		recordApplicationID := u.db.Where("miner_id = ?", miner.ID).Table("miner_applications").Pluck("application_id", &applicationIDs)
		if recordApplicationID.RecordNotFound() {
			log.Errorf("Cant find application_id from miner_applications table where miner_id is %v", miner.ID)
			continue
		}
		application, err := defaultData.getUserApplicationWithHistory(applicationIDs[0])
		if err != nil {
			log.Errorf("Cant find Application where id is %v, %v", applicationIDs[0], err)
		}
		log.Infof("Updating two oldest change histories from application where id is %v", application.ID)
		if application.ChangeHistory[0].AdminComment == "IMPORTED - successfully" || application.ChangeHistory[1].AdminComment == "IMPORTED - successfully" {
			changeHistoriesForUpdate := application.ChangeHistory[:2]
			for _, changeHistoryForUpdate := range changeHistoriesForUpdate {
				for _, err := range db.Exec("UPDATE change_history SET created_at = ? WHERE id = ?", application.CreatedAt, changeHistoryForUpdate.ID).GetErrors() {
					dbError = err
					log.Error("Error while updating changeHistory in DB ", err)
				}
			}
		}
		log.Infof("Updating miner with id %v from old created at %v to new created at %v", miner.ID, miner.CreatedAt, application.CreatedAt)
		for _, err := range u.db.Exec("UPDATE miners SET created_at = ? WHERE id = ?", application.CreatedAt, miner.ID).GetErrors() {
			dbError = err
			log.Error("Error while updating miner", err)
		}
	}
	if dbError != nil {
		db.Rollback()
		return dbError
	}
	db.Commit()

	return nil
}

func (u minerData) createMiner(miner *Miner) error {
	db := u.db.Begin()
	var dbError error
	miner.Username = NormalizeMail(miner.Username)
	for _, err := range db.Create(miner).GetErrors() {
		dbError = err
		log.Error("Error while persisting new miner in DB ", err)
	}
	if dbError != nil {
		db.Rollback()
		return dbError
	}
	db.Commit()

	return nil
}

func (u minerData) findMiners(email string) ([]Miner, error) {
	var miners []Miner
	record := u.db.Unscoped().Where("username = ?", NormalizeMail(email)).Preload("Nodes").Find(&miners)
	if record.RecordNotFound() {
		return nil, ErrCannotFindUser
	}
	if errs := record.GetErrors(); len(errs) > 0 {
		for err := range errs {
			log.Errorf("Error occurred while fetching miners by user email %v - %v", email, err)
		}
		return nil, errUnableToRead
	}
	return miners, nil
}

func (u minerData) findSpecificMiner(id string) (Miner, error) {
	var miner Miner
	record := u.db.Unscoped().Where("id = ?", id).Preload("Images").Preload("Applications").Preload("Nodes").Find(&miner)
	if record.RecordNotFound() {
		return Miner{}, errCannotFindMiner
	}
	if errs := record.GetErrors(); len(errs) > 0 {
		for err := range errs {
			log.Errorf("Error occurred while fetching miner by id %v - %v", id, err)
		}
		return Miner{}, errUnableToRead
	}
	return miner, nil
}
func (u minerData) findSpecificDisabledMiner(id string, deletionTime *time.Time) (Miner, error) {
	var miner Miner
	record := u.db.Unscoped().Where("id = ?", id).Preload("Images").Preload("Applications").Preload("Nodes", func(db *gorm.DB) *gorm.DB {
		return db.Unscoped().Where("nodes.deleted_at=?", deletionTime)
	}).Find(&miner)
	if record.RecordNotFound() {
		return Miner{}, errCannotFindMiner
	}
	if errs := record.GetErrors(); len(errs) > 0 {
		for err := range errs {
			log.Errorf("Error occurred while fetching miner by id %v - %v", id, err)
		}
		return Miner{}, errUnableToRead
	}
	return miner, nil
}

func (u minerData) findSpecificMinerWithApplications(id string) (Miner, error) {
	var miner Miner
	record := u.db.Preload("Images").Preload("Nodes").Preload("Applications", func(db *gorm.DB) *gorm.DB {
		return db.Order("applications.id DESC")
	}).Find(&miner, "id = ?", id)
	if record.RecordNotFound() {
		return Miner{}, errCannotFindMiner
	}
	if errs := record.GetErrors(); len(errs) > 0 {
		for err := range errs {
			log.Errorf("Error occurred while fetching miner by id %v - %v", id, err)
		}
		return Miner{}, errUnableToRead
	}
	return miner, nil
}

func (u minerData) addImagesToMiner(minerID uint, images []uint) error {
	db := u.db.Begin()
	var dbError error
	for _, err := range db.Exec("UPDATE images SET miner_id=? WHERE id IN (?)", minerID, images).GetErrors() {
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

func (u minerData) updateMiner(miner *Miner, updatedIDs []uint) error {
	db := u.db.Begin()
	var dbError error
	if len(updatedIDs) > 0 {
		for _, err := range db.Exec("UPDATE nodes SET deleted_at=? WHERE deleted_at IS NULL AND miner_id = ? AND id NOT IN (?)", time.Now(), miner.ID, updatedIDs).GetErrors() {
			dbError = err
			log.Error("Error while updating miner nodes in DB", err)
			break
		}
	} else {
		for _, err := range db.Exec("UPDATE nodes SET deleted_at=? WHERE deleted_at IS NULL AND miner_id = ?", time.Now(), miner.ID).GetErrors() {
			dbError = err
			log.Error("Error while updating miner nodes in DB", err)
			break
		}
	}
	if dbError != nil {
		db.Rollback()
		return dbError
	}
	for _, err := range db.Save(&miner).GetErrors() {
		dbError = err
		log.Error("Error while updating miner in DB ", err)
	}
	if dbError != nil {
		db.Rollback()
		return dbError
	}
	db.Commit()

	return nil
}

func (u minerData) findImportData() ([]MinerImport, error) {
	var importData []MinerImport
	record := u.db.Preload("Nodes").Find(&importData)
	if record.RecordNotFound() || len(importData) == 0 {
		return nil, errCannotFindMinerImportData
	}
	if errs := record.GetErrors(); len(errs) > 0 {
		for err := range errs {
			log.Error("Error occurred while fetching miner import ", err)
		}
		return nil, errUnableToRead
	}
	return importData, nil
}

func (u minerData) saveImportData(importedData []MinerImport) ([]MinerImport, error) {
	db := u.db.Begin()
	successfullySaved := []uint{}
	var dbError error
	for _, importRecord := range importedData {
		if importRecord.ID > 0 {
			for _, err := range db.Model(&importRecord).Update("number_of_miners", importRecord.NumberOfMiners).GetErrors() {
				dbError = err
				log.Error("Error while updating imported miner in DB", err)
				break
			}
		} else {
			for _, err := range db.Save(&importRecord).GetErrors() {
				dbError = err
				log.Error("Error while updating import record in DB", err)
				break
			}
		}
		if dbError == nil {
			successfullySaved = append(successfullySaved, importRecord.ID)
		} else {
			break
		}
	}
	if dbError == nil {
		if len(successfullySaved) > 0 {
			for _, err := range db.Exec("UPDATE miner_imports SET deleted_at=? WHERE deleted_at IS NULL AND id NOT IN (?)", time.Now(), successfullySaved).GetErrors() {
				dbError = err
				log.Error("Error while updating imported miner in DB", err)
				break
			}
		} else {
			for _, err := range db.Exec("UPDATE miner_imports SET deleted_at=? WHERE deleted_at IS NULL", time.Now()).GetErrors() {
				dbError = err
				log.Error("Error while updating imported miner in DB", err)
				break
			}
		}
	}
	if dbError != nil {
		db.Rollback()
		return importedData, dbError
	}
	db.Commit()

	return importedData, nil
}

func (u minerData) removeImportData(record MinerImport) error {
	db := u.db.Begin()
	var dbError error
	for _, err := range db.Delete(&record).GetErrors() {
		dbError = err
		log.Error("Error while updating imported miner in DB", err)
		break
	}
	if dbError != nil {
		db.Rollback()
		log.Errorf("Unable to remove processed import minerData %v due to error %v", record, dbError)
		return dbError
	}
	db.Commit()
	return nil
}

//fetches all miners including deleted ones
func (u minerData) getAllMiners() ([]Miner, error) {
	var miners []Miner
	record := u.db.Unscoped().Find(&miners)
	if record.RecordNotFound() {
		return nil, ErrCannotLoadMiners
	}
	if errs := record.GetErrors(); len(errs) > 0 {
		for err := range errs {
			log.Error("Error while fetching miners", err)
		}
		return nil, errUnableToRead
	}
	return miners, nil

}
func (u minerData) findActiveNodeByKey(key string) (node Node, err error) {
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
