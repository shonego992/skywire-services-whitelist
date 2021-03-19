package whitelist

import (
	"net/http"
	"time"

	"github.com/SkycoinPro/skywire-services-util/src/rpc/authorization"
)

//stub implements minerStore and userStore interfaces for mocking purposes
type stub struct {
	data []interface{}
	err  []error
}

func newStub(simulatedErrors []error, data []interface{}) *stub {
	return &stub{
		data: data,
		err:  simulatedErrors,
	}
}
func (cs *stub) returnRecord() interface{} {
	response := cs.data[0]
	cs.data = cs.data[1:]
	return response
}

func (cs *stub) returnError() error {
	err := cs.err[0]
	cs.err = cs.err[1:]
	return err
}

//datastore implementation
func (cs *stub) createApplication(application *Application) error {
	return cs.returnError()
}
func (cs *stub) createChangeHistory(change *ChangeHistory) error {
	return cs.returnError()
}
func (cs *stub) getUserApplicationWithHistory(applicationId uint) (Application, error) {
	return cs.returnRecord().(Application), cs.returnError()
}
func (cs *stub) updateChangeHistory(change *ChangeHistory) error {
	return cs.returnError()
}
func (cs *stub) getWhitelist(id string) (Application, error) {
	return cs.returnRecord().(Application), cs.returnError()
}
func (cs *stub) updateApplication(application *Application, oldChangeID uint) error {
	return cs.returnError()
}
func (cs *stub) findApplications(email string) ([]Application, error) {
	return cs.returnRecord().([]Application), cs.returnError()
}
func (cs *stub) appUpdateCreatedAt(application *Application, createdAtImport string) error {
	return cs.returnError()
}
func (cs *stub) findActiveNodeByKey(key string) (Node, error) {
	return cs.returnRecord().(Node), cs.returnError()
}
func (cs *stub) transferMiner(request transferMinerReq, currentUser string) error {
	return cs.returnError()
}
func (cs *stub) updateMiner(miner *Miner, updatedIDs []uint) error {
	return cs.returnError()
}
func (cs *stub) updateShopInfo(shopData *ShopData) error {
	return cs.returnError()
}
func (cs *stub) findMinerForApplication(applicationID uint) (Miner, error) {
	return cs.returnRecord().(Miner), cs.returnError()
}
func (cs *stub) findPendingApplication(email string) (Application, error) {
	return cs.returnRecord().(Application), cs.returnError()
}
func (cs *stub) getApplicationForNode(node Node) (Application, error) {
	return cs.returnRecord().(Application), cs.returnError()
}
func (cs *stub) findApplicationForMiner(miner Miner) (Application, error) {
	return cs.returnRecord().(Application), cs.returnError()
}
func (cs *stub) getAllWhitelists() ([]Application, error) {
	return cs.returnRecord().([]Application), cs.returnError()
}
func (cs *stub) getApplicationApprovedNodesCount(appID uint) (numOfApprovedNodes int) {
	return cs.returnRecord().(int)
}
func (cs *stub) findAllApplicationIDsForMiner(minerID uint) ([]uint, error) {
	return cs.returnRecord().([]uint), cs.returnError()
}
func (cs *stub) findImagesByHash(hashed string, minId uint) (image Image, err error) {
	return cs.returnRecord().(Image), cs.returnError()
}
func (cs *stub) getAllImages() ([]Image, error) {
	return cs.returnRecord().([]Image), cs.returnError()
}
func (cs *stub) getImageRecordsForUser(username string) ([]Image, error) {
	return cs.returnRecord().([]Image), cs.returnError()
}

func (cs *stub) updateHash(imgID uint, hashed string) error {
	return cs.returnError()
}
func (cs *stub) getMinerForApplication(application *Application) Miner {
	return cs.returnRecord().(Miner)
}
func (cs *stub) getWhitelistWithoutPreload(id string) (Application, error) {
	return cs.returnRecord().(Application), cs.returnError()
}
func (cs *stub) removeNode(node *Node, currTime time.Time) error {
	return cs.returnError()
}

//userStore implementation
func (cs *stub) create(newUser *User) error {
	return cs.returnError()
}
func (cs *stub) findBy(email string) (User, error) {
	return cs.returnRecord().(User), cs.returnError()
}
func (cs *stub) findUserById(id uint) (User, error) {
	return cs.returnRecord().(User), cs.returnError()
}
func (cs *stub) updateUser(user *User) error {
	return cs.returnError()
}
func (cs *stub) getUsers() ([]User, error) {
	return cs.returnRecord().([]User), cs.returnError()
}
func (cs *stub) getAdmins() ([]User, error) {
	return cs.returnRecord().([]User), cs.returnError()
}
func (cs *stub) findUserWithAPIKeys(email string) (User, error) {
	return cs.returnRecord().(User), cs.returnError()
}
func (cs *stub) removeAPIKey(key *ApiKey) error {
	return cs.returnError()
}
func (cs *stub) findPayoutAddress(username string) (string, error) {
	return cs.returnRecord().(string), cs.returnError()
}
func (cs *stub) findUsernameByApiKey(key string) (string, error) {
	return cs.returnRecord().(string), cs.returnError()
}
func (cs *stub) getUsersWithMiners() ([]User, error) {
	return cs.returnRecord().([]User), cs.returnError()
}
func (cs *stub) createExportRecord(record *ExportRecord) error {
	return cs.returnError()
}
func (cs *stub) createNewActiveAddress(address *Address, hasOld bool) error {
	return cs.returnError()
}
func (cs *stub) getUsersWithAddressesAndMiners() ([]User, error) {
	return cs.returnRecord().([]User), cs.returnError()
}
func (cs *stub) userUpdateCreatedAt(username string, createdAt time.Time) error {
	return cs.returnError()
}
func (cs *stub) updateUserWithoutChildEntities(user *User) error {
	return cs.returnError()
}

func (cs *stub) findAddressRecord(addr string) (Address, error) {
	return cs.returnRecord().(Address), cs.returnError()
}

//minerStore implementation
func (cs *stub) createMiner(miner *Miner) error {
	return cs.returnError()
}
func (cs *stub) findMiners(email string) ([]Miner, error) {
	return cs.returnRecord().([]Miner), cs.returnError()
}
func (cs *stub) findSpecificMiner(id string) (Miner, error) {
	return cs.returnRecord().(Miner), cs.returnError()
}

func (cs *stub) getAllMiners() ([]Miner, error) {
	return cs.returnRecord().([]Miner), cs.returnError()
}
func (cs *stub) findImportData() ([]MinerImport, error) {
	return cs.returnRecord().([]MinerImport), cs.returnError()
}
func (cs *stub) saveImportData(importedData []MinerImport) ([]MinerImport, error) {
	return cs.returnRecord().([]MinerImport), cs.returnError()
}
func (cs *stub) removeImportData(record MinerImport) error {
	return cs.returnError()
}
func (cs *stub) createShopInfo(shopData *ShopData) error {
	return cs.returnError()
}
func (cs *stub) getShopInfo(id uint64) (ShopData, error) {
	return cs.returnRecord().(ShopData), cs.returnError()
}
func (cs *stub) exportMiners(request exportMinersReq, startDate time.Time, endDate time.Time, useFilter bool) ([]Miner, error) {
	return cs.returnRecord().([]Miner), cs.returnError()
}
func (cs *stub) addImagesToMiner(minerID uint, images []uint) error {
	return cs.returnError()
}
func (cs *stub) createTransferMinerRecord(transferRecord MinerTransfer) error {
	return cs.returnError()
}
func (cs *stub) findMinersNotDeleted(email string) ([]Miner, error) {
	return cs.returnRecord().([]Miner), cs.returnError()
}
func (cs *stub) getMiners() ([]Miner, error) {
	return cs.returnRecord().([]Miner), cs.returnError()
}
func (cs *stub) removeMiner(id string, currTime time.Time) error {
	return cs.returnError()
}

func (cs *stub) findMinerByNodeKey(nodeKey string) (Miner, error) {
	return cs.returnRecord().(Miner), cs.returnError()
}

func (cs *stub) findSpecificMinerWithApplications(id string) (Miner, error) {
	return cs.returnRecord().(Miner), cs.returnError()
}

func (cs *stub) activateMiner(minerId string) error {
	return cs.returnError()
}

func (cs *stub) activateNode(nodeID uint) error {
	return cs.returnError()
}

func (cs *stub) findSpecificDisabledMiner(id string, deletionTime *time.Time) (Miner, error) {
	return cs.returnRecord().(Miner), cs.returnError()
}

//minerservice func stub
func (cs *stub) createMinerEntitiesFromImport(importedRecord MinerImport) error {
	return cs.returnError()
}

//service func stub
func (cs *stub) extractDataFromApplicationRequest(req *http.Request) (ApplicationReq, error) {
	return cs.returnRecord().(ApplicationReq), cs.returnError()
}

//some predefined users/miners/apps for reusing in testing

var stbMiner1 = Miner{ID: 42, Username: "Miner1", Type: DIY}
var stbMiner2 = Miner{ID: 43, Username: "Miner2", Type: DIY}
var stbMiner3 = Miner{ID: 44, Username: "Miner3", Type: DIY}
var stbMiner4 = Miner{ID: 1, Nodes: []Node{stbValidNode1, stbValidNode2, stbValidNode3}, Username: "test1@mail.com", Type: DIY}

var stbUsername1 = "test1@mail.com"
var stbUsername2 = "test2@mail.com"

var stbMiner1DiffUsername = Miner{ID: 42, Username: "Miner3", Type: DIY}
var stbValidChangeHistory = ChangeHistory{Nodes: []*Node{{MinerID: 42, Key: "DummyKey", ID: 42}}}

var stbValidImportedData = []MinerImport{{ID: 42, Username: "ImportMiner1", NumberOfMiners: 42}, {ID: 43, Username: "ImportMiner2", NumberOfMiners: 43}}
var stbValidMinerImport = MinerImport{ID: 42, Username: "ImportMiner1", NumberOfMiners: 42}

var stbValidWhiteListApplications = []Application{stbValidApp1, stbValidApp2}
var stbValidApp1 = Application{ID: 442, CurrentStatus: PENDING, Username: "test1@mail.com", ChangeHistory: []ChangeHistory{{ID: 45, Status: PENDING, Images: stbImages}}}
var stbValidApp2 = Application{ID: 443, CurrentStatus: PENDING, Username: "test2@mail.com"}
var stbValidAppDenied = Application{ID: 442, CurrentStatus: DECLINED, Username: "test1@mail.com", ChangeHistory: []ChangeHistory{{ID: 452, Status: PENDING, Images: stbImages}, {ID: 453, Status: DECLINED, Images: stbImages}}}
var stbValidPendingApps = []Application{stbValidApp2, stbValidApp1}
var stbValidAppApproved1 = Application{ID: 442, CurrentStatus: APPROVED, Username: "test1@mail.com"}
var stbValidAppApproved2 = Application{ID: 443, CurrentStatus: APPROVED, Username: "test2@mail.com"}
var stbValidAppListLastApproved = []Application{stbValidAppApproved1, stbValidAppApproved2}

var stbValidAppApprovedWithHistory = Application{ID: 442, CurrentStatus: APPROVED, Username: "test1@mail.com", ChangeHistory: []ChangeHistory{{ID: 45, Status: PENDING, Images: stbImages}, {ID: 46, Status: APPROVED, Images: stbImages}}}
var stbValidListForStatusUpdate = []Application{stbValidAppApprovedWithHistory, stbValidAppDenied}

var stbValidApplicationChangeStatusApproved = ChangeApplicationStatus{ApplicationId: 442, Status: APPROVED, AdminComment: "Change to approved"}
var stbValidApplicationChangeStatusDenied = ChangeApplicationStatus{ApplicationId: 442, Status: DECLINED, AdminComment: "Change to denied"}

var skyCoinAdress = "2mj62xkVrB24z1rbVwBfLtGPe4FkFP1eKaW"
var stbValidUser1 = User{Username: "test1@mail.com", Addresses: []Address{{Username: "test1@mail.com", SkycoinAddress: skyCoinAdress}}, Status: 1}
var stbValidUser1a = User{Username: "test1@mail.com", Status: 1}
var stbValidUser2 = User{Username: "test2@mail.com", Status: 1}
var stbValidUserList = []User{stbValidUser1, stbValidUser2}
var stbValidAdmin1 = User{Username: "test1@mail.com", Status: 65}
var stbValidAdmin2 = User{Username: "test2@mail.com", Status: 65}
var stbValidAdminList = []User{stbValidAdmin1, stbValidAdmin2}

var stbValidNode1 = Node{ID: 42, Key: "038a3a1e1a7c34c979887bfc6e20e7ae09c754ef6ab57a8485563fc6345c054138", MinerID: 142}
var stbValidNode2 = Node{ID: 43, Key: "02cf5cb4fb42ee00c128383fb0cb1f1670c64789d5525a98412b4adb0c6be123fb", MinerID: 143}
var stbValidNode3 = Node{ID: 44, Key: "02af759eb74b26463a7b2780c4fe6d76c160f18ce36e564650ce7d0e808d0e7f7b", MinerID: 144}
var stbValidNode4 = Node{ID: 500, Key: "02af759eb74b26463a7b2780c4fe6d76c160f18ce36e964650ce7d0e808d0e7f71"}
var stbValidAppRequest = ApplicationReq{Description: "Description1", Location: "Location1", Nodes: []Node{stbValidNode1, stbValidNode2, stbValidNode3}}
var stbValidChangeHistory1 = ChangeHistory{Nodes: []*Node{&stbValidNode1, &stbValidNode2, &stbValidNode3}, Description: "Description1", Location: "Location1"}
var stbUpdateRights = User{Username: "test1@mail.com", Rights: []authorization.Right{{Name: "create-user", Value: true}, {Name: "disable-user", Value: true}}}
var stbAdminWithRightsResponse = User{Username: "test1@mail.com", Status: 192, Rights: []authorization.Right{{Name: "create-user", Value: true}, {Name: "disable-user", Value: true}}}
var stbEnableAllRights = []authorization.Right{{Name: "create-user", Value: true}, {Name: "disable-user", Value: true}}
var stbValidUserWithAPIKeys = User{Username: "test1@mail.com", Status: 1, ApiKeys: []ApiKey{{Key: "O6VYbvcFUWWHyHGuSaAYJRk5PSpp6ZkImDxgmXLC"}, {Key: "CVVJ9YOB2uivEMKC8KCUGfhzxNGJRhqWTpiIThXI"}}}
var stbValidUserWithAPIKeysOneRemoved = User{Username: "test1@mail.com", Status: 1, ApiKeys: []ApiKey{{Key: "O6VYbvcFUWWHyHGuSaAYJRk5PSpp6ZkImDxgmXLC"}}}
var stbValidAppReqResponse = ApplicationReq{Description: "Description1", Location: "Location1", Nodes: []Node{Node{Key: "node1"}}, OldImages: []Image{Image{Path: "oldImage1"}}}
var stbFiles = []string{"file1", "file2"}
var stbImages = []*Image{{ID: 1, Path: "pathToImageOne"}, {ID: 2, Path: "pathToImageTwo"}}
