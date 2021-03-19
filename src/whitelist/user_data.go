package whitelist

import (
	"time"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"github.com/SkycoinPro/skywire-services-whitelist/src/database/postgres"
)

type userStore interface {
	create(newUser *User) error
	findBy(email string) (User, error)
	findUserById(id uint) (User, error)
	updateUser(user *User) error
	updateUserWithoutChildEntities(user *User) error
	getUsers() ([]User, error)
	getAdmins() ([]User, error)
	findUserWithAPIKeys(email string) (User, error)
	removeAPIKey(key *ApiKey) error
	findPayoutAddress(username string) (string, error)
	findUsernameByApiKey(key string) (string, error)
	getUsersWithMiners() ([]User, error)
	getUsersWithAddressesAndMiners() ([]User, error)
	userUpdateCreatedAt(username string, createdAt time.Time) error
	createNewActiveAddress(address *Address, hasOld bool) error
	createExportRecord(record *ExportRecord) error
}

type userData struct {
	db *gorm.DB
}

func DefaultUserData() userData {
	return NewUserData(postgres.DB)
}

func NewUserData(database *gorm.DB) userData {
	return userData{
		db: database,
	}
}

func (u userData) findUsernameByApiKey(key string) (string, error) {
	var apiKey ApiKey
	record := u.db.Where("key = ?", key).Find(&apiKey)
	if record.RecordNotFound() {
		return "", errCannotFindUserByKey
	}
	if errs := record.GetErrors(); len(errs) > 0 {
		for err := range errs {
			log.Errorf("Error occurred while fetching user by key %v - %v", key, err)
		}
		return "", errUnableToRead
	}
	return apiKey.Username, nil
}

func (u userData) findPayoutAddress(username string) (string, error) {
	var user User
	record := u.db.Where("username = ?", NormalizeMail(username)).Preload("Addresses", func(db *gorm.DB) *gorm.DB {
		return db.Order("addresses.id DESC")
	}).Find(&user)
	if record.RecordNotFound() {
		return "", errCannotLoadPayoutAddress
	}
	if errs := record.GetErrors(); len(errs) > 0 {
		for err := range errs {
			log.Error("Error while finding payout address", err)
		}
		return "", errUnableToRead
	}
	if len(user.Addresses) == 0 {
		return "", errNoAddressSet
	}
	return user.Addresses[0].SkycoinAddress, nil
}

// const adminStatusStart uint8 = 16
func (u userData) getUsers() ([]User, error) {
	var users []User
	record := u.db.Where("status < ?", adminStatusStart).Find(&users)
	if record.RecordNotFound() {
		return nil, errCannotLoadUsers
	}
	if errs := record.GetErrors(); len(errs) > 0 {
		for err := range errs {
			log.Error("Error while finding users", err)
		}
		return nil, errUnableToRead
	}
	return users, nil
}

// const adminStatusStart uint8 = 16
func (u userData) getUsersWithMiners() ([]User, error) {
	var users []User
	record := u.db.Where("status < ?", adminStatusStart).Preload("Miners").Find(&users)
	if record.RecordNotFound() {
		return nil, errCannotLoadUsers
	}
	if errs := record.GetErrors(); len(errs) > 0 {
		for err := range errs {
			log.Error("Error while finding users", err)
		}
		return nil, errUnableToRead
	}
	return users, nil
}
func (u userData) getUsersWithAddressesAndMiners() ([]User, error) {
	var users []User
	record := u.db.Where("status < ?", adminStatusStart).Preload("Addresses").Preload("Miners").Find(&users)
	if record.RecordNotFound() {
		return nil, errCannotLoadUsers
	}
	if errs := record.GetErrors(); len(errs) > 0 {
		for err := range errs {
			log.Errorf("Error while finding users %v", err)
		}
		return nil, errUnableToRead
	}
	return users, nil
}
func (u userData) getAdmins() ([]User, error) {
	var users []User
	record := u.db.Where("status >= ?", adminStatusStart).Find(&users)
	if record.RecordNotFound() {
		return nil, errCannotLoadUsers
	}
	if errs := record.GetErrors(); len(errs) > 0 {
		for err := range errs {
			log.Error("Error while qu", err)
		}
		return nil, errUnableToRead
	}
	return users, nil
}

func (u userData) createExportRecord(record *ExportRecord) error {
	db := u.db.Begin()
	var dbError error
	for _, err := range db.Create(record).GetErrors() {
		dbError = err
		log.Error("Error while persisting new export record in DB: ", err)
	}
	if dbError != nil {
		db.Rollback()
		return dbError
	}
	db.Commit()

	return nil
}

func (u userData) create(newUser *User) error {
	db := u.db.Begin()
	var dbError error
	newUser.Username = NormalizeMail(newUser.Username)
	for _, err := range db.Create(newUser).GetErrors() {
		dbError = err
		log.Errorf("Error while persisting new user %v in DB %v", newUser.Username, err)
	}
	if dbError != nil {
		db.Rollback()
		return dbError
	}
	db.Commit()

	return nil
}

func (u userData) removeUser(user *User) error {
	if errs := u.db.Delete(user).GetErrors(); len(errs) > 0 {
		for err := range errs {
			log.Errorf("Error occurred while removing user %v - %v", user.Username, err)
		}
		return errUnableToSave
	}

	return nil
}

func (u userData) findBy(email string) (User, error) {
	var user User
	record := u.db.Where("username = ?", NormalizeMail(email)).Preload("Addresses", func(db *gorm.DB) *gorm.DB {
		return db.Unscoped().Order("addresses.id DESC")
	}).Preload("Applications", func(db *gorm.DB) *gorm.DB {
		return db.Order("applications.id ASC")
	}).Find(&user)
	if record.RecordNotFound() {
		return User{}, ErrCannotFindUser
	}
	if errs := record.GetErrors(); len(errs) > 0 {
		for err := range errs {
			log.Errorf("Error occurred while fetching user by email %v - %v", email, err)
		}
		return User{}, errUnableToRead
	}
	return user, nil
}

func (u userData) findUserById(value uint) (User, error) {
	var user User
	record := u.db.Where("username = ?", value).Find(&user)
	if record.RecordNotFound() {
		return User{}, nil
	}
	if errs := record.GetErrors(); len(errs) > 0 {
		for err := range errs {
			log.Errorf("Error occurred while fetching user by email %v - %v", value, err)
		}
		return User{}, errUnableToRead
	}
	return user, nil
}

func (u userData) updateUser(user *User) error {
	db := u.db.Begin()
	var dbError error
	for _, err := range db.Save(&user).GetErrors() {
		dbError = err
		log.Error("Error while updating user in DB", err)
	}
	if dbError != nil {
		db.Rollback()
		return dbError
	}
	db.Commit()

	return nil
}

func (u userData) updateUserWithoutChildEntities(user *User) error {
	db := u.db.Begin()
	var dbError error
	for _, err := range db.Set("gorm:association_autoupdate", false).Save(&user).GetErrors() {
		dbError = err
		log.Error("Error while updating user in DB")
	}
	if dbError != nil {
		db.Rollback()
		return dbError
	}
	db.Commit()

	return nil
}
func (u userData) findUserWithAPIKeys(email string) (User, error) {
	var user User
	record := u.db.Where("username = ?", email).Preload("ApiKeys").Find(&user)
	if record.RecordNotFound() {
		return User{}, ErrCannotFindUser
	}
	if errs := record.GetErrors(); len(errs) > 0 {
		for err := range errs {
			log.Errorf("Error occurred while fetching user by email %v - %v", email, err)
		}
		return User{}, errUnableToRead
	}
	return user, nil
}

func (u userData) removeAPIKey(key *ApiKey) error {
	if errs := u.db.Delete(key).GetErrors(); len(errs) > 0 {
		for err := range errs {
			log.Errorf("Error occurred while removing api key %v - %v", key, err)
		}
		return errUnableToSave
	}
	return nil
}
func (u userData) userUpdateCreatedAt(username string, createdAt time.Time) error {
	db := u.db.Begin()
	var dbError error
	for _, err := range db.Unscoped().Exec("UPDATE users SET created_at= ? WHERE username = ?", createdAt, username).GetErrors() {
		dbError = err
		log.Errorf("Error while updating  user %v - %v", username, err)
		break
	}

	if dbError != nil {
		db.Rollback()
		return dbError
	}
	db.Commit()
	return nil
}
func (u userData) createNewActiveAddress(address *Address, hasOld bool) error {
	db := u.db.Begin()
	var dbError error
	//deleting only occurs when there is old addresses present
	if hasOld {
		for _, err := range db.Where("username=?", address.Username).Delete(Address{}).GetErrors() {
			dbError = err
			log.Errorf("Error while updating imported miner in DB %v", err)
			break
		}
	}
	if dbError == nil {
		for _, err := range db.Create(address).GetErrors() {
			dbError = err
			log.Error("Error while persisting address record in DB ", err)
			break
		}
	}

	if dbError != nil {
		db.Rollback()
		return dbError
	}
	db.Commit()

	return nil
}
