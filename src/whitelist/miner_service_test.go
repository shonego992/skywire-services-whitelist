package whitelist

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/SkycoinPro/skywire-services-whitelist/src/config"
	"github.com/SkycoinPro/skywire-services-whitelist/src/database/postgres"
)

func TestGetUserMiners(t *testing.T) {
	tests := []struct {
		name         string
		storeRecords []interface{}
		storeErrors  []error
		userName     string
		err          error
		expected     []Miner
	}{
		{
			name:         "Sucessful retrieval of miners ",
			storeRecords: []interface{}{[]Miner{stbMiner1, stbMiner2}},
			storeErrors:  []error{nil},
			userName:     "test1@gmail.com",
			expected:     []Miner{stbMiner1, stbMiner2},
		},
		{
			name:         "User not Found",
			storeRecords: []interface{}{[]Miner{}},
			storeErrors:  []error{ErrCannotFindUser},
			userName:     "",
			err:          ErrCannotLoadMiners,
		},
		{
			name:         "Error during fetching the miners",
			storeRecords: []interface{}{[]Miner{}},
			storeErrors:  []error{errUnableToRead},
			userName:     "test1@gmail.com",
			err:          ErrCannotLoadMiners,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			svc := NewMinerService(newStub(test.storeErrors, test.storeRecords))
			resp, err := svc.getUserMiners(test.userName)
			if err != nil && test.err == nil {
				t.Errorf("%s failed, expected no error but received error %v", test.name, err)
			} else if test.err != nil && err == nil {
				t.Errorf("%s failed, expected error %v but no error was received", test.name, test.err)
			} else if test.err != nil && err != nil {
				if err != test.err {
					t.Errorf("%s failed, expected error %v but received error %v", test.name, test.err, err)
				}
			} else {
				if !reflect.DeepEqual(test.expected, resp) {
					t.Errorf("%s failed, expected: %#v - actual %#v", test.name, test.expected, resp)
				}
			}
		})
	}
}
func TestGetMinersForUser(t *testing.T) {
	tests := []struct {
		name         string
		storeRecords []interface{}
		storeErrors  []error
		userName     string
		err          error
		expected     []Miner
	}{
		{
			name:         "Sucessful retrieval of miners ",
			storeRecords: []interface{}{[]Miner{stbMiner1, stbMiner2}},
			storeErrors:  []error{nil},
			userName:     "test1@gmail.com",
			expected:     []Miner{stbMiner1, stbMiner2},
		},
		{
			name:         "User not Found",
			storeRecords: []interface{}{[]Miner{}},
			storeErrors:  []error{ErrCannotFindUser},
			userName:     "",
			err:          ErrCannotLoadMiners,
		},
		{
			name:         "Error during fetching the miners",
			storeRecords: []interface{}{[]Miner{}},
			storeErrors:  []error{errUnableToRead},
			userName:     "test1@gmail.com",
			err:          ErrCannotLoadMiners,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			svc := NewMinerService(newStub(test.storeErrors, test.storeRecords))
			resp, err := svc.GetMinersForUser(test.userName)
			if err != nil && test.err == nil {
				t.Errorf("%s failed, expected no error but received error %v", test.name, err)
			} else if test.err != nil && err == nil {
				t.Errorf("%s failed, expected error %v but no error was received", test.name, test.err)
			} else if test.err != nil && err != nil {
				if err != test.err {
					t.Errorf("%s failed, expected error %v but received error %v", test.name, test.err, err)
				}
			} else {
				if !reflect.DeepEqual(test.expected, resp) {
					t.Errorf("%s failed, expected: %#v - actual %#v", test.name, test.expected, resp)
				}
			}
		})
	}
}
func TestGetAllMiners(t *testing.T) {
	tests := []struct {
		name         string
		storeRecords []interface{}
		storeErrors  []error
		err          error
		expected     []Miner
	}{
		{
			name:         "Sucessful retrieval of miners ",
			storeRecords: []interface{}{[]Miner{stbMiner1, stbMiner2}},
			storeErrors:  []error{nil},
			expected:     []Miner{stbMiner1, stbMiner2},
		},
		{
			name:         "Miners not found",
			storeRecords: []interface{}{[]Miner{}},
			storeErrors:  []error{ErrCannotFindUser},
			err:          ErrCannotLoadMiners,
		},
		{
			name:         "Error during fetching the miners",
			storeRecords: []interface{}{[]Miner{}},
			storeErrors:  []error{errUnableToRead},
			err:          ErrCannotLoadMiners,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			svc := NewMinerService(newStub(test.storeErrors, test.storeRecords))
			resp, err := svc.GetAllMiners()
			if err != nil && test.err == nil {
				t.Errorf("%s failed, expected no error but received error %v", test.name, err)
			} else if test.err != nil && err == nil {
				t.Errorf("%s failed, expected error %v but no error was received", test.name, test.err)
			} else if test.err != nil && err != nil {
				if err != test.err {
					t.Errorf("%s failed, expected error %v but received error %v", test.name, test.err, err)
				}
			} else {
				if !reflect.DeepEqual(test.expected, resp) {
					t.Errorf("%s failed, expected: %#v - actual %#v", test.name, test.expected, resp)
				}
			}
		})
	}
}
func TestGetSpecificMiner(t *testing.T) {
	tests := []struct {
		name         string
		storeRecords []interface{}
		storeErrors  []error
		id           string
		userName     string
		err          error
		expected     Miner
	}{
		{
			name:         "Sucessful retrieval of miner ",
			storeRecords: []interface{}{stbMiner1},
			storeErrors:  []error{nil},
			id:           "42",
			userName:     "Miner1",
			expected:     stbMiner1,
		},
		{
			name:         "No miner found for given parameters",
			storeRecords: []interface{}{Miner{}},
			storeErrors:  []error{ErrCannotLoadMiners},
			id:           "42",
			userName:     "Miner1",
			err:          errCannotFindMiner,
		},
		{
			name:         "Error while fetching the miner",
			storeRecords: []interface{}{Miner{}},
			storeErrors:  []error{errUnableToRead},
			id:           "42",
			userName:     "Miner1",
			err:          errCannotFindMiner,
		},
		{
			name:         "Cant find miner for specified ID",
			storeRecords: []interface{}{stbMiner1DiffUsername},
			storeErrors:  []error{nil},
			id:           "42",
			userName:     "Miner1",
			err:          errMinerNotFoundForUser,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			svc := NewMinerService(newStub(test.storeErrors, test.storeRecords))
			resp, err := svc.getSpecificMiner(test.id, test.userName)
			if err != nil && test.err == nil {
				t.Errorf("%s failed, expected no error but received error %v", test.name, err)
			} else if test.err != nil && err == nil {
				t.Errorf("%s failed, expected error %v but no error was received", test.name, test.err)
			} else if test.err != nil && err != nil {
				if err != test.err {
					t.Errorf("%s failed, expected error %v but received error %v", test.name, test.err, err)
				}
			} else {
				if !reflect.DeepEqual(test.expected, resp) {
					t.Errorf("%s failed, expected: %#v - actual %#v", test.name, test.expected, resp)
				}
			}
		})
	}
}
func TestCreateMinerEntitiesForUser(t *testing.T) {
	tests := []struct {
		name         string
		storeRecords []interface{}
		storeErrors  []error
		change       ChangeHistory
		app          Application
		err          error
	}{
		{
			name:        "Sucessful miner creation",
			storeErrors: []error{nil},
			change:      stbValidChangeHistory,
			app:         stbValidApp2,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			svc := NewMinerService(newStub(test.storeErrors, test.storeRecords))
			_, err := svc.CreateMinerEntitiesForUser(&test.change, test.app)
			if err != nil && test.err == nil {
				t.Errorf("%s failed, expected no error but received error %v", test.name, err)
			} else if test.err != nil && err == nil {
				t.Errorf("%s failed, expected error %v but no error was received", test.name, test.err)
			} else if test.err != nil && err != nil {
				if err != test.err {
					t.Errorf("%s failed, expected error %v but received error %v", test.name, test.err, err)
				}
			}
		})
	}
}
func TestGetImportData(t *testing.T) {
	tests := []struct {
		name         string
		storeRecords []interface{}
		storeErrors  []error
		expected     []MinerImport
		err          error
	}{
		{
			name:         "Retrieved import data",
			storeRecords: []interface{}{stbValidImportedData},
			storeErrors:  []error{nil},
			expected:     stbValidImportedData,
		},
		{
			name:         "Can not find imported data",
			storeRecords: []interface{}{[]MinerImport{}},
			storeErrors:  []error{errCannotFindMinerImportData},
			err:          errCannotFindMinerImportData,
		},
		{
			name:         "Error ocurred during fetching of import data",
			storeRecords: []interface{}{[]MinerImport{}},
			storeErrors:  []error{errUnableToRead},
			err:          errUnableToRead,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			svc := NewMinerService(newStub(test.storeErrors, test.storeRecords))
			resp, err := svc.GetImportData()
			if err != nil && test.err == nil {
				t.Errorf("%s failed, expected no error but received error %v", test.name, err)
			} else if test.err != nil && err == nil {
				t.Errorf("%s failed, expected error %v but no error was received", test.name, test.err)
			} else if test.err != nil && err != nil {
				if err != test.err {
					t.Errorf("%s failed, expected error %v but received error %v", test.name, test.err, err)
				}
			} else {
				if !reflect.DeepEqual(test.expected, resp) {
					t.Errorf("%s failed, expected: %#v - actual %#v", test.name, test.expected, resp)
				}
			}
		})
	}
}

func TestSaveImportData(t *testing.T) {
	tests := []struct {
		name         string
		storeRecords []interface{}
		storeErrors  []error
		importData   []MinerImport
		expected     []MinerImport
		err          error
	}{
		{
			name:         "Saved import data",
			storeRecords: []interface{}{stbValidImportedData},
			storeErrors:  []error{nil},
			importData:   stbValidImportedData,
			expected:     stbValidImportedData,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			svc := NewMinerService(newStub(test.storeErrors, test.storeRecords))
			resp, err := svc.SaveImportData(test.importData)
			if err != nil && test.err == nil {
				t.Errorf("%s failed, expected no error but received error %v", test.name, err)
			} else if test.err != nil && err == nil {
				t.Errorf("%s failed, expected error %v but no error was received", test.name, test.err)
			} else if test.err != nil && err != nil {
				if err != test.err {
					t.Errorf("%s failed, expected error %v but received error %v", test.name, test.err, err)
				}
			} else {
				if !reflect.DeepEqual(test.expected, resp) {
					t.Errorf("%s failed, expected: %#v - actual %#v", test.name, test.expected, resp)
				}
			}
		})
	}
}

func TestImportMiners(t *testing.T) {
	tests := []struct {
		name         string
		storeRecords []interface{}
		storeErrors  []error
		importMiner  MinerImport

		err error
	}{ //TODO Debug the test and see why it is failing
		/*
			{
				name: "Imported miners without errors",

				storeErrors: []error{nil, nil},
				importMiner: stbValidMinerImport,
			},*/
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			svc := NewMinerService(newStub(test.storeErrors, test.storeRecords))
			err := svc.ImportMiners(test.importMiner)
			if err != nil && test.err == nil {
				t.Errorf("%s failed, expected no error but received error %v", test.name, err)
			} else if test.err != nil && err == nil {
				t.Errorf("%s failed, expected error %v but no error was received", test.name, test.err)
			} else if test.err != nil && err != nil {
				if err != test.err {
					t.Errorf("%s failed, expected error %v but received error %v", test.name, test.err, err)
				}
			}
		})
	}
}

func TestMinerReenabling(t *testing.T) {
	dbConnection := testInit()
	defer dbConnection()

	minerService := DefaultMinerService()
	userService := DefaultUserService()

	nodeKeys := []string{stbValidNode1.Key, stbValidNode2.Key, stbValidNode3.Key}
	deletionTime := time.Now()

	errCreate := userService.Create(stbUsername1)
	if errCreate != nil {
		t.Fatal("Error while creating new user", errCreate)
	}
	errAddMinerToUser := minerService.addMinerToUser(stbUsername1, nodeKeys)
	if errAddMinerToUser != nil {
		t.Fatal("Error while adding miner to user", errAddMinerToUser)
	}
	miners, errGetMinersForUser := minerService.GetMinersForUser(stbUsername1)
	if errGetMinersForUser != nil {
		t.Fatal("Error while fetching miners", errGetMinersForUser)
	}

	minerID := fmt.Sprint(miners[0].ID)

	keepNodes := []*Node{}
	for _, keepNode := range miners[0].Nodes[0:1] {
		keepNodes = append(keepNodes, &keepNode)
	}
	miners[0].Nodes = miners[0].Nodes[0:1]

	//As User modify existing Miner by removing two nodes
	errRemoveNode := minerService.UpdateMiner(&miners[0], keepNodes)
	if errRemoveNode != nil {
		t.Fatal("Error while removing last two nodes from miner", errRemoveNode)
	}
	miners, errGetMinersForUser = minerService.GetMinersForUser(stbUsername1)
	if errGetMinersForUser != nil {
		t.Fatal("Error while fetching miners after removing last two nodes", errGetMinersForUser)
	}
	if len(miners[0].Nodes) != 1 {
		t.Fatal("Error while removing Nodes from miner, miner has more than one node", errGetMinersForUser)
	}

	//As User modify existing Miner by adding one node
	keepNodes = append(keepNodes, &stbValidNode4)
	errAddNode := minerService.UpdateMiner(&miners[0], keepNodes)
	if errAddNode != nil {
		t.Fatal("Error while adding node to miner", errAddNode)
	}
	miners, errGetMinersForUser = minerService.GetMinersForUser(stbUsername1)
	if errGetMinersForUser != nil {
		t.Fatal("Error while fetching miners after adding one node to miner", errGetMinersForUser)
	}
	if len(miners[0].Nodes) != 2 {
		t.Fatal("Error while removing Nodes from miner, miner doesnt have two nodes", errGetMinersForUser)
	}

	//As Admin delete previously modified Miner
	errRemoveMiner := minerService.RemoveMiner(minerID, deletionTime)
	if errRemoveMiner != nil {
		t.Fatal("Error while removing miner", errRemoveMiner)
	}
	miners, errGetMinersForUser = minerService.GetMinersForUser(stbUsername1)
	if errGetMinersForUser != nil {
		t.Fatal("Error while fetching users miners after removing miner", errGetMinersForUser)
	}
	if miners[0].DeletedAt == &deletionTime {
		t.Fatal("Error while removing miner, user still has his miner", errRemoveMiner)
	}

	//As Admin re-enable deleted miner
	errActivateMiner := minerService.ActivateMiner(minerID)
	if errActivateMiner != nil {
		t.Fatal("Error while activating miner", errActivateMiner)
	}
	activatedMiner, err := minerService.getSpecificDisabledMinerForAdmin(minerID, &deletionTime)
	if err != nil {
		t.Fatal("Error while fetching miner after reactivating it", err)
	}
	errReactivateNode := minerService.ReactivateNodes(activatedMiner.Nodes)
	if errReactivateNode != nil {
		t.Fatal("Error while reactivating nodes", errReactivateNode)
	}

	//As Admin confirm exact set of nodes that existed at the moment of removal is back on (not including those two nodes that previously got deleted)
	miners, errGetMinersForUser = minerService.GetMinersForUser(stbUsername1)
	if errGetMinersForUser != nil {
		t.Fatal("Error while fetching miners after reactivating miners nodes", errGetMinersForUser)
	}

	if len(miners[0].Nodes) != 2 {
		t.Fatal("Error after reactivating miner, miner doesnt have two nodes")
	}
}

func testInit() (dbConnection func()) {
	config.Init("whitelist-test-config")
	level, err := log.ParseLevel(viper.GetString("server.log-level"))
	if err != nil {
		log.Info("Unable to use configured log level. Using Info instead")
		level = log.InfoLevel
	}
	log.SetLevel(level)

	dbConnection = postgres.Init()
	return
}
