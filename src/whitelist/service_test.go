package whitelist

import (
	"reflect"
	"testing"
)

func TestGetAllWhitelists(t *testing.T) {
	tests := []struct {
		name         string
		storeRecords []interface{}
		storeErrors  []error
		err          error
		expected     []Application
	}{
		{
			name:         "Can not find whitelist in database ",
			storeRecords: []interface{}{stbValidWhiteListApplications},
			storeErrors:  []error{nil},
			expected:     stbValidWhiteListApplications,
		},

		{
			name:         "Can not find whitelist in database ",
			storeRecords: []interface{}{[]Application{}},
			storeErrors:  []error{errCannotLoadWhitelists},
			err:          errCannotLoadWhitelists,
		},
		{
			name:         "Error while fetching the whitelist from database ",
			storeRecords: []interface{}{[]Application{}},
			storeErrors:  []error{errUnableToRead},
			err:          errCannotLoadWhitelists,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			svc := NewService(newStub(test.storeErrors, test.storeRecords))
			resp, err := svc.GetAllWhitelists()
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
func TestGetWhitelist(t *testing.T) {
	tests := []struct {
		name         string
		storeRecords []interface{}
		storeErrors  []error
		id           string
		err          error
		expected     Application
	}{
		{
			name:         "Whitelist sucessfuly retrieved",
			storeRecords: []interface{}{stbValidApp1},
			storeErrors:  []error{nil},
			id:           "442",
			expected:     stbValidApp1,
		},

		{
			name:         "Can not find whitelist in database ",
			storeRecords: []interface{}{Application{}},
			storeErrors:  []error{errCannotLoadWhitelists},
			id:           "442",
			err:          ErrCannotFindWhitelist,
		},
		{
			name:         "Error while fetching the whitelist from database ",
			storeRecords: []interface{}{Application{}},
			storeErrors:  []error{errUnableToRead},
			id:           "442",
			err:          ErrCannotFindWhitelist,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			svc := NewService(newStub(test.storeErrors, test.storeRecords))
			resp, err := svc.GetWhitelist(test.id)
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

func TestUpdateApplicationStatus(t *testing.T) {
	tests := []struct {
		name         string
		storeRecords []interface{}
		storeErrors  []error
		req          ChangeApplicationStatus
		err          error
		expected     string
	}{
		{
			name:         "Sucessful status update, not approved",
			storeRecords: []interface{}{stbValidApp1},
			storeErrors:  []error{nil, nil, nil},
			req:          stbValidApplicationChangeStatusDenied,
			err:          nil,
			expected:     stbUsername1,
		},
		{
			name:         "Sucessful status update, approved, no previous auto declined",
			storeRecords: []interface{}{stbValidApp1, stbValidListForStatusUpdate},
			storeErrors:  []error{nil, nil, nil, nil},
			req:          stbValidApplicationChangeStatusApproved,
			err:          nil,
			expected:     stbUsername1,
		},
		{
			name:         "Can not find whitelist in database",
			storeRecords: []interface{}{Application{}},
			storeErrors:  []error{ErrCannotFindUser},
			req:          stbValidApplicationChangeStatusApproved,
			err:          ErrCannotFindWhitelist,
		},
		{
			name:         "Application is not curently in any progress",
			storeRecords: []interface{}{stbValidAppApproved1},
			storeErrors:  []error{nil},
			req:          stbValidApplicationChangeStatusApproved,
			err:          errNoApplicationInProgressForUser,
		},
		{
			name:         "Can not find whitelist in database",
			storeRecords: []interface{}{Application{}},
			storeErrors:  []error{ErrCannotFindUser},
			req:          stbValidApplicationChangeStatusApproved,
			err:          ErrCannotFindWhitelist,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			svc := NewService(newStub(test.storeErrors, test.storeRecords))
			resp, err := svc.UpdateApplicationStatus(test.req)
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

// Commented while service changes forms

func TestUpdateApplication(t *testing.T) {
	tests := []struct {
		name         string
		storeRecords []interface{}
		storeErrors  []error
		req          ApplicationReq
		userName     string
		err          error
		expected     ApplicationError
	}{

		{
			name:         "Cannot find applications",
			storeRecords: []interface{}{[]Application{}},
			storeErrors:  []error{ErrCannotFindUser},
			err:          ErrCannotFindUser,
		},
		{
			name:         "Error while fetching applications",
			storeRecords: []interface{}{[]Application{}},
			storeErrors:  []error{errUnableToRead},
			err:          errUnableToRead,
		},
		{
			name:         "Application already active",
			storeRecords: []interface{}{stbValidAppListLastApproved},
			storeErrors:  []error{nil},
			err:          errApplicationInProgress,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			svc := NewService(newStub(test.storeErrors, test.storeRecords))
			resp := svc.UpdateApplication(test.req, test.userName)

			if resp.Error != nil && test.err == nil {
				t.Errorf("%s failed, expected no error but received error %v", test.name, resp.Error)
			} else if test.err != nil && resp.Error == nil {
				t.Errorf("%s failed, expected error %v but no error was received", test.name, test.err)
			} else if test.err != nil && resp.Error != nil {
				if resp.Error != test.err {
					t.Errorf("%s failed, expected error %v but received error %v", test.name, test.err, resp.Error)
				}
			} else {
				if !reflect.DeepEqual(test.expected, resp) {
					t.Errorf("%s failed, expected: %#v - actual %#v", test.name, test.expected, resp)
				}

			}
		})
	}
}
func TestCompareRequestWithLatestChangeHistory(t *testing.T) {
	tests := []struct {
		name         string
		storeRecords []interface{}
		storeErrors  []error
		req          ApplicationReq
		latestChange ChangeHistory
		expected     bool
	}{

		{
			name:         "Found matches",
			req:          stbValidAppRequest,
			latestChange: stbValidChangeHistory1,
			expected:     true,
		},
		{
			name:         "Different description error",
			req:          stbValidAppRequest,
			latestChange: ChangeHistory{Nodes: []*Node{&stbValidNode1, &stbValidNode2, &stbValidNode3}, Description: "Description12323", Location: "Location1"},
			expected:     false,
		},
		{
			name:         "Different location error",
			req:          stbValidAppRequest,
			latestChange: ChangeHistory{Nodes: []*Node{&stbValidNode1, &stbValidNode2, &stbValidNode3}, Description: "Description1", Location: "Location123"},
			expected:     false,
		},
		{
			name:         "Different node array length error",
			req:          stbValidAppRequest,
			latestChange: ChangeHistory{Nodes: []*Node{&stbValidNode1, &stbValidNode2}, Description: "Description1", Location: "Location123"},
			expected:     false,
		},
		{
			name:         "Different image array length error",
			req:          stbValidAppRequest,
			latestChange: ChangeHistory{Nodes: []*Node{&stbValidNode1, &stbValidNode2, &stbValidNode3}, Description: "Description1", Location: "Location1", Images: []*Image{{ID: 42}}},
			expected:     false,
		},
		{
			name:         "No matched nodes",
			req:          ApplicationReq{Description: "Description1", Location: "Location1", Nodes: []Node{stbValidNode1}},
			latestChange: ChangeHistory{Nodes: []*Node{&stbValidNode2}, Description: "Description1", Location: "Location1"},
			expected:     false,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			svc := NewService(newStub(test.storeErrors, test.storeRecords))
			resp := svc.compareRequestWithLatestChangeHistory(test.req, test.latestChange)

			if !reflect.DeepEqual(test.expected, resp) {
				t.Errorf("%s failed, expected: %#v - actual %#v", test.name, test.expected, resp)
			}

		})
	}
}
func TestGetActiveApplication(t *testing.T) {
	tests := []struct {
		name         string
		storeRecords []interface{}
		storeErrors  []error
		userName     string
		err          error
		expected     Application
	}{
		{
			name:         "Sucessfuly returned application",
			storeRecords: []interface{}{[]Application{stbValidApp2, stbValidApp1}, stbValidApp1},
			storeErrors:  []error{nil, nil},
			userName:     "test1@gmail.com",
			expected:     stbValidApp1,
		},

		{
			name:         "Cannot find applications for given user",
			storeRecords: []interface{}{[]Application{}},
			storeErrors:  []error{ErrCannotFindUser},
			userName:     "test1@gmail.com",
			err:          ErrCannotFindUser,
		},
		{
			name:         "Error while fetching the applications for given user",
			storeRecords: []interface{}{[]Application{}},
			storeErrors:  []error{errUnableToRead},
			userName:     "test1@gmail.com",
			err:          errUnableToRead,
		},
		{
			name:         "No applications for user",
			storeRecords: []interface{}{[]Application{}},
			storeErrors:  []error{nil},
			userName:     "test1@gmail.com",
			err:          errNoApplicationInProgressForUser,
		},
		{
			name:         "No active applications in progress for user",
			storeRecords: []interface{}{stbValidAppListLastApproved},
			storeErrors:  []error{nil},
			userName:     "test1@gmail.com",
			err:          errNoApplicationInProgressForUser,
		},
		{
			name:         "Could not find application with history for given id",
			storeRecords: []interface{}{[]Application{stbValidApp2, stbValidApp1}, Application{}},
			storeErrors:  []error{nil, ErrCannotFindUser},
			userName:     "test1@gmail.com",
			err:          ErrCannotFindUser,
		},
		{
			name:         "Error while fetching application for given id",
			storeRecords: []interface{}{[]Application{stbValidApp2, stbValidApp1}, Application{}},
			storeErrors:  []error{nil, errUnableToRead},
			userName:     "test1@gmail.com",
			err:          errUnableToRead,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			svc := NewService(newStub(test.storeErrors, test.storeRecords))
			resp, err := svc.getActiveApplication(test.userName)
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
