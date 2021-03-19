package whitelist

import (
	"errors"
	"reflect"
	"testing"

	"github.com/SkycoinPro/skywire-services-util/src/rpc/authorization"
)

func TestGetUsers(t *testing.T) {
	tests := []struct {
		name         string
		storeRecords []interface{}
		storeErrors  []error
		err          error
		expected     []User
	}{
		{
			name:         "Sucessful retrieval of users ",
			storeRecords: []interface{}{stbValidUserList},
			storeErrors:  []error{nil},
			expected:     stbValidUserList,
		},
		{
			name:         "Users not found",
			storeRecords: []interface{}{[]User{}},
			storeErrors:  []error{errCannotLoadUsers},
			err:          errCannotLoadUsers,
		},
		{
			name:         "Error while fetching users",
			storeRecords: []interface{}{[]User{}},
			storeErrors:  []error{errUnableToRead},
			err:          errCannotLoadUsers,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			svc := NewUserService(newStub(test.storeErrors, test.storeRecords))
			resp, err := svc.GetUsers()
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
func TestGetAdmins(t *testing.T) {
	tests := []struct {
		name         string
		storeRecords []interface{}
		storeErrors  []error
		err          error
		expected     []User
	}{
		{
			name:         "Sucessful retrieval of admins ",
			storeRecords: []interface{}{stbValidAdminList},
			storeErrors:  []error{nil},
			expected:     stbValidAdminList,
		},
		{
			name:         "Admins not found",
			storeRecords: []interface{}{[]User{}},
			storeErrors:  []error{errCannotLoadUsers},
			err:          errCannotLoadUsers,
		},
		{
			name:         "Error while fetching admins",
			storeRecords: []interface{}{[]User{}},
			storeErrors:  []error{errUnableToRead},
			err:          errCannotLoadUsers,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			svc := NewUserService(newStub(test.storeErrors, test.storeRecords))
			resp, err := svc.GetAdmins()
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
func TestUpdateRights(t *testing.T) {
	tests := []struct {
		name         string
		storeRecords []interface{}
		storeErrors  []error
		user         *User
		rights       []authorization.Right
		err          error
	}{
		{
			name:         "Sucessful update of rights ",
			storeRecords: []interface{}{},
			storeErrors:  []error{nil},
			user:         &stbValidUser1,
			rights:       stbEnableAllRights,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			svc := NewUserService(newStub(test.storeErrors, test.storeRecords))
			err := svc.UpdateRights(test.user, test.rights)
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
func TestFindBy(t *testing.T) {
	tests := []struct {
		name         string
		storeRecords []interface{}
		storeErrors  []error
		userName     string
		err          error
		expected     User
	}{
		{
			name:         "Sucessful retrieval of user",
			storeRecords: []interface{}{stbValidUser1},
			storeErrors:  []error{nil},
			userName:     "test1@gmail.com",
			expected:     stbValidUser1,
		},
		{
			name:         "Empty username given",
			storeRecords: []interface{}{},
			storeErrors:  []error{nil},
			userName:     "",
			err:          errMissingMandatoryFields,
		},
		{
			name:         "User not found",
			storeRecords: []interface{}{User{}},
			storeErrors:  []error{ErrCannotFindUser},
			userName:     "test1@gmail.com",
			err:          ErrCannotFindUser,
		},
		{
			name:         "Error while fetching user",
			storeRecords: []interface{}{User{}},
			storeErrors:  []error{errUnableToRead},
			userName:     "test1@gmail.com",
			err:          errUnableToRead,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			svc := NewUserService(newStub(test.storeErrors, test.storeRecords))
			resp, err := svc.FindBy(test.userName)
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

func TestUpdateAddress(t *testing.T) {
	tests := []struct {
		name         string
		storeRecords []interface{}
		storeErrors  []error
		userName     string
		address      string
		err          error
		expected     User
	}{
		{
			name:         "Sucessful updating of address",
			storeRecords: []interface{}{stbValidUser1a},
			storeErrors:  []error{nil, nil},
			userName:     "test1@gmail.com",
			address:      skyCoinAdress,
			expected:     stbValidUser1,
		},
		{
			name:     "Zero length adress",
			userName: "test1@gmail.com",
			address:  "",
			err:      errMissingMandatoryFields,
		},
		{
			name:     "Error while decoding",
			userName: "test1@gmail.com",
			address:  "invalid",
			err:      errSkycoinAddressNotValid,
		},
		{
			name:         "Error while fetching user",
			storeRecords: []interface{}{User{}},
			storeErrors:  []error{errUnableToRead},
			userName:     "test1@gmail.com",
			address:      skyCoinAdress,
			err:          errUnableToRead,
		},
		{
			name:         "User not found",
			storeRecords: []interface{}{User{}},
			storeErrors:  []error{ErrCannotFindUser},
			userName:     "test2@gmail.com",
			address:      skyCoinAdress,
			err:          ErrCannotFindUser,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			svc := NewUserService(newStub(test.storeErrors, test.storeRecords))
			resp, err := svc.UpdateAddress(test.userName, test.address)
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
func TestImportUser(t *testing.T) {
	tests := []struct {
		name         string
		storeRecords []interface{}
		storeErrors  []error
		userName     string
		err          error
		expected     bool
	}{ //TODO Mock rpc client call
		/*
			{
				name:         "Sucessful import",
				storeRecords: []interface{}{User{}, &rpc.Client{}},
				storeErrors:  []error{ErrCannotFindUser, nil, nil},
				userName:     "test1@gmail.com",
				expected:     true,
			},*/
		{
			name:         "User with given username already exists ",
			storeRecords: []interface{}{stbValidUser1},
			storeErrors:  []error{nil},
			userName:     "test1@gmail.com",
			expected:     false,
		},

		{
			name:         "Error while fetching username- user already exists with given username",
			storeRecords: []interface{}{User{}},
			storeErrors:  []error{errUnableToRead},
			userName:     "test1@gmail.com",
			err:          errUnableToRead,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			svc := NewUserService(newStub(test.storeErrors, test.storeRecords))
			resp, err := svc.ImportUser(test.userName)
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

func TestCreate(t *testing.T) {
	tests := []struct {
		name         string
		storeRecords []interface{}
		storeErrors  []error
		userName     string
		err          error
	}{
		{
			name:        "Sucessfully created new user",
			storeErrors: []error{nil},
			userName:    "test1@gmail.com",
		},
		{
			name:     "Empty username error",
			userName: "",
			err:      errMissingMandatoryFields,
		},
		{
			name:        "Empty username error",
			storeErrors: []error{errors.New("Database saving error")},
			userName:    "test1@gmail.com",
			err:         errUnableToSave,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			svc := NewUserService(newStub(test.storeErrors, test.storeRecords))
			err := svc.Create(test.userName)
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
func TestGetKeys(t *testing.T) {
	tests := []struct {
		name         string
		storeRecords []interface{}
		storeErrors  []error
		userName     string
		err          error
		expected     []string
	}{
		{
			name:         "Sucessfully retrieved api keys",
			storeRecords: []interface{}{stbValidUserWithAPIKeys},
			storeErrors:  []error{nil},
			userName:     "test1@gmail.com",
			expected:     []string{"O6VYbvcFUWWHyHGuSaAYJRk5PSpp6ZkImDxgmXLC", "CVVJ9YOB2uivEMKC8KCUGfhzxNGJRhqWTpiIThXI"},
		},
		{
			name:         "Can not find user for given email",
			storeRecords: []interface{}{User{}},
			storeErrors:  []error{ErrCannotFindUser},
			userName:     "test1@gmail.com",
			err:          ErrCannotFindUser,
		},
		{
			name:         "Error while fetching the user",
			storeRecords: []interface{}{User{}},
			storeErrors:  []error{errUnableToRead},
			userName:     "test1@gmail.com",
			err:          ErrCannotFindUser,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			svc := NewUserService(newStub(test.storeErrors, test.storeRecords))
			resp, err := svc.GetKeys(test.userName)
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
func TestAddKey(t *testing.T) {
	tests := []struct {
		name         string
		storeRecords []interface{}
		storeErrors  []error
		userName     string
		err          error
	}{
		{
			name:         "Sucessfully retrieved api keys",
			storeRecords: []interface{}{stbValidUserWithAPIKeys},
			storeErrors:  []error{nil, nil},
			userName:     "test1@gmail.com",
		},
		{
			name:         "Can not find user for given email",
			storeRecords: []interface{}{User{}},
			storeErrors:  []error{ErrCannotFindUser},
			userName:     "test1@gmail.com",
			err:          ErrCannotFindUser,
		},
		{
			name:         "Error while fetching the user",
			storeRecords: []interface{}{User{}},
			storeErrors:  []error{errUnableToRead},
			userName:     "test1@gmail.com",
			err:          ErrCannotFindUser,
		},
		{
			name:         "Error while updating user after key creation",
			storeRecords: []interface{}{stbValidUserWithAPIKeys},
			storeErrors:  []error{nil, errors.New("Database saving error")},
			userName:     "test1@gmail.com",
			err:          errUnableToSave,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			svc := NewUserService(newStub(test.storeErrors, test.storeRecords))
			_, err := svc.AddKey(test.userName)
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
func TestRemoveKey(t *testing.T) {
	tests := []struct {
		name         string
		storeRecords []interface{}
		storeErrors  []error
		userName     string
		key          string
		err          error
		expected     []string
	}{
		{
			name:         "Sucessful removal of key",
			storeRecords: []interface{}{stbValidUserWithAPIKeys},
			storeErrors:  []error{nil, nil},
			userName:     "test1@gmail.com",
			key:          "CVVJ9YOB2uivEMKC8KCUGfhzxNGJRhqWTpiIThXI",
		},
		{
			name:         "Can not find user for given email",
			storeRecords: []interface{}{User{}},
			storeErrors:  []error{ErrCannotFindUser},
			userName:     "test1@gmail.com",
			key:          "CVVJ9YOB2uivEMKC8KCUGfhzxNGJRhqWTpiIThXI",
			err:          ErrCannotFindUser,
		},
		{
			name:         "Error while fetching user",
			storeRecords: []interface{}{User{}},
			storeErrors:  []error{errUnableToRead},
			userName:     "test1@gmail.com",
			key:          "CVVJ9YOB2uivEMKC8KCUGfhzxNGJRhqWTpiIThXI",
			err:          ErrCannotFindUser,
		},
		{
			name:         "Error while fetching user",
			storeRecords: []interface{}{stbValidUserWithAPIKeys},
			storeErrors:  []error{nil},
			userName:     "test1@gmail.com",
			key:          "Key that does exist",
			err:          errCannotFindAPIKey,
		},
		{
			name:         "Error while removing key from the found user",
			storeRecords: []interface{}{stbValidUserWithAPIKeys},
			storeErrors:  []error{nil, errUnableToSave},
			userName:     "test1@gmail.com",
			key:          "CVVJ9YOB2uivEMKC8KCUGfhzxNGJRhqWTpiIThXI",
			err:          errUnableToSave,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			svc := NewUserService(newStub(test.storeErrors, test.storeRecords))
			err := svc.RemoveKey(test.userName, test.key)
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
func TestFindUserInfo(t *testing.T) {
	tests := []struct {
		name         string
		storeRecords []interface{}
		storeErrors  []error
		userName     string
		err          error
		expected     User
	}{
		{
			name:         "Sucessfully retrieved user info",
			storeRecords: []interface{}{stbValidUser1},
			storeErrors:  []error{nil},
			userName:     "test1@gmail.com",
			expected:     stbValidUser1,
		},
		{
			name:     "Zero length username",
			userName: "",
			err:      errMissingMandatoryFields,
		},
		{
			name:         "Can not find user",
			storeRecords: []interface{}{User{}},
			storeErrors:  []error{ErrCannotFindUser, errors.New("Creating new user error")},
			userName:     "test1@gmail.com",
			err:          errUnableToSave,
		},
		{
			name:         "Error while fetching user",
			storeRecords: []interface{}{User{}},
			storeErrors:  []error{errUnableToRead},
			userName:     "test1@gmail.com",
			err:          errUnableToRead,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			svc := NewUserService(newStub(test.storeErrors, test.storeRecords))
			resp, err := svc.FindUserInfo(test.userName)
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
