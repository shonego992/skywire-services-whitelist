package whitelist

import (
	"net/rpc"
	"time"

	"fmt"
	"strings"

	"github.com/dchest/uniuri"
	log "github.com/sirupsen/logrus"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/spf13/viper"
	"github.com/SkycoinPro/skywire-services-util/src/rpc/authentication"
	"github.com/SkycoinPro/skywire-services-util/src/rpc/authorization"
	"github.com/SkycoinPro/skywire-services-whitelist/src/rpc_client"
)

// UserService provides access to User related data
type UserService struct {
	db userStore
}

// DefaultUserService prepares new instance of UserService
func DefaultUserService() UserService {
	return NewUserService(DefaultUserData())
}

// NewUserService prepares new instance of UserService
func NewUserService(whitelistStore userStore) UserService {
	return UserService{
		db: whitelistStore,
	}
}

func (us *UserService) FindUserByApiKey(key string) (string, error) {
	if len(key) == 0 {
		return "", errMissingMandatoryFields
	}

	username, err := us.db.findUsernameByApiKey(key)
	if err != nil {
		log.Errorf("Unable to find User by api key %v due to error %v", key, err)
		return "", err
	}

	return username, nil
}

func (us *UserService) FindPayoutAddress(username string) (string, error) {
	address, err := us.db.findPayoutAddress(username)
	if err != nil {
		if err == errNoAddressSet {
			log.Debugf("User %v has no skycoin address set: ", username)
			return "", nil
		}
		return "", err
	}
	return address, nil
}

func (us *UserService) GetUsers() ([]User, error) {
	users, err := us.db.getUsersWithMiners()
	if err != nil {
		return nil, errCannotLoadUsers
	}
	return users, nil
}

func (us *UserService) GetUsersWithAddressesAndMiners() ([]User, error) {
	users, err := us.db.getUsersWithAddressesAndMiners()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (us *UserService) CreateExportRecord(record ExportRecord) error {
	err := us.db.createExportRecord(&record)
	if err != nil {
		return errCreatingExportRecord
	}
	return nil
}

func (us *UserService) GetAdmins() ([]User, error) {
	users, err := us.db.getAdmins()
	if err != nil {
		return nil, errCannotLoadUsers
	}
	return users, nil
}

func (us *UserService) UpdateRights(dbUser *User, rights []authorization.Right) error {
	dbUser.SetCanReviewWhitelsit(false)
	dbUser.SetCanFlagUserAsVIP(false)
	for _, right := range rights {
		if right.Name == "review_whitelist" {
			dbUser.SetCanReviewWhitelsit(right.Value)
		} else if right.Name == "flag_vip" {
			dbUser.SetCanFlagUserAsVIP(right.Value)
		}
	}
	us.db.updateUser(dbUser)

	return nil
}

func (us *UserService) FindBy(username string) (User, error) {
	if len(username) == 0 {
		return User{}, errMissingMandatoryFields
	}

	u, err := us.db.findBy(username)
	if err != nil {
		log.Debugf("Unable to find User by username %v due to error %v", username, err)
		return User{}, err
	}
	if u.IsAdmin() {
		u.Rights = append(u.Rights, authorization.Right{Name: "flag-vip", Label: "Flag VIP", Value: u.CanFlagUserAsVIP()})
		u.Rights = append(u.Rights, authorization.Right{Name: "review-whitelist", Label: "Review Whitelist", Value: u.CanReviewWhitelsit()})

	}

	return u, nil
}

// Change users right to submit whitelists - if allowWhitelistSubmission is false then status is set to 2, and user
// cannot submit new whitelists
func (us *UserService) ChangeUserWhitelistSubmissionPrivilege(username string, allowWhitelistSubmission bool) error {
	if len(username) == 0 {
		return errMissingMandatoryFields
	}

	user, err := us.db.findBy(username)
	if err != nil {
		log.Infof("Unable to find User by username %v due to error %v", username, err)
		return err
	}

	var status uint8 = 2
	if allowWhitelistSubmission {
		status = 0
	}
	user.Status = status
	err = us.db.updateUserWithoutChildEntities(&user)
	if err != nil {
		log.Infof("Unable to update User by username %v due to error %v", username, err)
		return err
	}
	return nil
}

func (us *UserService) FindUserInfo(username string) (User, error) {
	if len(username) == 0 {
		return User{}, errMissingMandatoryFields
	}

	u, err := us.db.findBy(username)
	if err == ErrCannotFindUser {
		u.Username = username
		if err := us.db.create(&u); err != nil {
			return User{}, errUnableToSave
		}
		createdAt, err := rpc_client.FetchCreatedAt(username)
		if err != nil {
			log.Debug("Error while fetching created_at for username ", username)
		}
		err = us.UpdateCreatedAtForUser(username, createdAt)
		if err != nil {
			log.Debug("Error while updating created_at for username ", username)
		}

	} else if err != nil {
		log.Infof("Unable to find User by username %v due to error %v", username, err)
		return User{}, err
	}

	if u.IsAdmin() {
		u.Rights = append(u.Rights, authorization.Right{Name: "flag-vip", Label: "Flag VIP", Value: u.CanFlagUserAsVIP()})
		u.Rights = append(u.Rights, authorization.Right{Name: "review-whitelist", Label: "Review Whitelist", Value: u.CanReviewWhitelsit()})

	}

	return u, nil
}

func (us *UserService) Create(username string) error {
	var newUser User

	if len(username) == 0 {
		return errMissingMandatoryFields
	}
	newUser.Username = username
	if err := us.db.create(&newUser); err != nil {
		return errUnableToSave
	}
	return nil
}

// Return true if provided username didn't exist and User is new to the entire Skywire ecosystem
func (us *UserService) ImportUser(username string) (bool, error) {
	if _, err := us.FindBy(username); err != nil {
		if err != ErrCannotFindUser {
			return false, err
		}

		// check in remote service here
		client, err := rpc.DialHTTP(viper.GetString("rpc.user.protocol"), viper.GetString("rpc.user.address"))
		if err != nil {
			log.Error("dialing:", err)
			return false, err
		}
		args := &authentication.Request{Username: username}
		var reply authentication.Response
		err = client.Call("Handler.Create", args, &reply)
		if err != nil {
			log.Error("Unable to import new user in remote service with error: ", err)
			return false, err
		}

		if err := us.Create(username); err != nil {
			log.Error("Unable to persist newly imported user")
			return false, err
		}

		return reply.Success, nil
	}

	return false, nil
}

func containsAccessRight(rights []string, right string) bool {
	for _, r := range rights {
		if r == right {
			return true
		}
	}
	return false
}

func (us *UserService) GetKeys(email string) ([]string, error) {
	user, err := us.db.findUserWithAPIKeys(email)
	if err != nil {
		log.Errorf("Unable to find User by email %v due to error %v", email, err)
		return []string{""}, ErrCannotFindUser
	}

	keys := make([]string, 0)
	for _, val := range user.ApiKeys {
		keys = append(keys, val.Key)
	}

	return keys, nil
}

func (us *UserService) AddKey(email string) (string, error) {
	usr, err := us.db.findUserWithAPIKeys(email)
	if err != nil {
		log.Errorf("Unable to find User by email %v due to error %v", email, err)
		return "", ErrCannotFindUser
	}

	apiKey := ApiKey{
		Key: uniuri.NewLen(40),
	}
	usr.ApiKeys = append(usr.ApiKeys, apiKey)

	err = us.db.updateUser(&usr)
	if err != nil {
		log.Errorf("Error while updating user on api key creation")
		return "", errUnableToSave
	}

	return apiKey.Key, nil
}

func (us *UserService) RemoveKey(email, key string) error {
	usr, err := us.db.findUserWithAPIKeys(email)
	if err != nil {
		log.Errorf("Unable to find User by email %v due to error %v", email, err)
		return ErrCannotFindUser
	}

	var keyToBeRemoved ApiKey
	for _, k := range usr.ApiKeys {
		if k.Key == key {
			keyToBeRemoved = k
		}
	}
	if len(keyToBeRemoved.Key) == 0 {
		log.Infof("Unable to find API Key %v for user %v", key, email)
		return errCannotFindAPIKey
	}

	if err := us.db.removeAPIKey(&keyToBeRemoved); err != nil {
		return err
	}

	return nil
}

func (us *UserService) UpdateAddress(username, address string) (User, error) {
	if len(address) == 0 {
		return User{}, errMissingMandatoryFields
	} else if _, err := cipher.DecodeBase58Address(address); err != nil {
		return User{}, errSkycoinAddressNotValid
	} else if len(address) > 35 {
		return User{}, errSkycoinAddressNotValid
	}

	// TODO confirm we're not going to use this feature and remove CheckIsAddressFree as well
	// isAddressFree, err := us.CheckIsAddressFree(address)
	// if err != nil {
	// 	log.Error("Error while checking availability of address", err)
	// 	return User{}, err
	// }
	// if !isAddressFree {Error while persisting new user in DB
	// 	return User{}, errAddressAlreadyTaken
	// }
	dbUser, err := us.FindBy(username)
	if err != nil {
		return User{}, err
	}

	var newAddress Address
	newAddress.SkycoinAddress = address
	newAddress.Username = dbUser.Username
	err = us.db.createNewActiveAddress(&newAddress, len(dbUser.Addresses) > 0)
	if err != nil {
		log.Error("Error while perstiting new address", err)
		return User{}, errUnableToSave
	}
	tempAddr := make([]Address, len(dbUser.Addresses)+1)
	tempAddr[0] = newAddress
	copy(tempAddr[1:], dbUser.Addresses)
	dbUser.Addresses = tempAddr

	return dbUser, nil
}

func NormalizeMail(mail string) string {
	// normalize and check mail
	mail = strings.Replace(strings.ToLower(mail), " ", "", -1) //to lower case and remove space
	if strings.Contains(mail, "@gmail.com") {
		//if gmail remove dots from username part of email
		splinter := strings.Split(mail, "@")
		dropDot := strings.Replace(splinter[0], ".", "", -1)
		mail = fmt.Sprintf("%s@%s", dropDot, splinter[1])
	}
	return mail
}

func (us *UserService) UpdateCreatedAtForUser(username string, createdAt time.Time) error {
	err := us.db.userUpdateCreatedAt(username, createdAt)
	if err != nil {
		log.Errorf("Unable to update created at for %v because of error %v", username, err)
		return errUnableToUpdateCreatedAt
	}
	return nil
}

// func (us *UserService) CheckIsAddressFree(addr string) (isFree bool, err error) {
// 	_, err = us.db.findAddressRecord(addr)
// 	if err != nil && err == errCannotFindActiveAddress {
// 		isFree = true
// 		err = nil
// 		return
// 	}
// 	if err != nil {
// 		log.Error("Error while checking availability of address due to ", err)
// 	}
// 	return

// }

func (us *UserService) CheckIfUserHasThirdBatchMiners(username string) (hasThird bool, err error) {
	user, err := us.FindBy(username)
	if err != nil {
		log.Error("Unable to fetch user to check if he has third batch miners due to: ", err)
		return
	}
	for _, miner := range user.Miners {
		if miner.BatchLabel == "Third" {
			hasThird = true
			return
		}
	}
	return
}
