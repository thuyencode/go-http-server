package users

import (
	"errors"
	"fmt"
	"net/mail"
	"reflect"
	"testing"
)

func TestUnit_AddUser(t *testing.T) {
	testManager := NewManager()
	firstName, lastName := "Go", "Gopher"
	email, err := mail.ParseAddress("gopher@example.com")

	if err != nil {
		t.Fatalf("error parsing email address: %v", err)
	}

	err = testManager.AddUser(firstName, lastName, email.String())

	if err != nil {
		t.Fatalf("error creating user: %v", err)
	}

	expectedUserCount := 1
	actualUserCount := len(testManager.users)

	if actualUserCount != expectedUserCount {
		t.Fatalf(`bad user count, expected %d, got %d`, expectedUserCount, actualUserCount)
	}

	expectedUser := User{firstName, lastName, *email}
	actualUser := testManager.users[0]

	if !reflect.DeepEqual(expectedUser, actualUser) {
		t.Errorf("bad added user,\nexpected: %v,\ngot: %v", expectedUser, actualUser)
	}
}

func TestUnit_AddUser_InvalidEmail(t *testing.T) {
	testManager := NewManager()
	firstName, lastName := "Go", "Gopher"
	invalidEmail := "gopher"

	err := testManager.AddUser(firstName, lastName, invalidEmail)

	if err == nil {
		t.Fatalf(`no error returned for invalid email: "%s"`, invalidEmail)
	}

	expectedErrorMessage := fmt.Sprintf("invalid email: %s", invalidEmail)
	actualErrorMessage := err.Error()

	if actualErrorMessage != expectedErrorMessage {
		t.Errorf(`bad error message, expected: "%s", got: "%s"`, expectedErrorMessage, actualErrorMessage)
	}

	expectedUserCount := 0
	actualUserCount := len(testManager.users)

	if actualUserCount != expectedUserCount {
		t.Errorf(`bad user count, expected %d, got %d`, expectedUserCount, actualUserCount)
	}
}

func TestUnit_AddUser_EmptyName(t *testing.T) {
	testManager := NewManager()
	firstName, lastName := "", ""
	email, err := mail.ParseAddress("gopher@example.com")

	if err != nil {
		t.Fatalf("error parsing email address: %v", err)
	}

	actualErr := testManager.AddUser(firstName, lastName, email.String())

	if actualErr == nil {
		t.Fatal("bad return, expected an error")
	}

	expectedErr := ErrEmptyNames

	if !errors.Is(actualErr, expectedErr) {
		t.Errorf(`bad error, expected "%v", got "%v"`, expectedErr, actualErr)
	}

	expectedUserCount := 0
	actualUserCount := len(testManager.users)

	if actualUserCount != expectedUserCount {
		t.Errorf(`bad user count, expected %d, got %d`, expectedUserCount, actualUserCount)
	}
}

func TestUnit_AddUser_DuplicatedName(t *testing.T) {
	testManager := NewManager()
	firstName, lastName := "Go", "Gopher"
	email, err := mail.ParseAddress("gopher@example.com")

	if err != nil {
		t.Fatalf("error parsing email address: %v", err)
	}

	err = testManager.AddUser(firstName, lastName, email.String())

	if err != nil {
		t.Fatalf("error creating user: %v", err)
	}

	err = testManager.AddUser(firstName, lastName, email.String())
	expectedErr := ErrDuplicatedName

	if err == nil {
		t.Fatal("bad return, expected an error")
	}

	if !errors.Is(err, expectedErr) {
		t.Errorf(`bad error, expected "%v", got "%v"`, expectedErr, err)
	}

	expectedUserCount := 1
	actualUserCount := len(testManager.users)

	if actualUserCount != expectedUserCount {
		t.Errorf(`bad user count, expected %d, got %d`, expectedUserCount, actualUserCount)
	}
}

func TestUnit_GetUserByName(t *testing.T) {
	testManager := NewManager()
	addUserParams := []struct{ firstName, lastName, email string }{
		{firstName: "foo", lastName: "bar", email: "f.bar@example.com"},
		{firstName: "bar", lastName: "baz", email: "b.baz@example.com"},
		{firstName: "foo", lastName: "baz", email: "f.baz@example.com"},
		{firstName: "baz", lastName: "foo", email: "baz.f@example.com"},
	}

	for _, params := range addUserParams {
		err := testManager.AddUser(params.firstName, params.lastName, params.email)

		if err != nil {
			t.Fatalf("error creating user: %v", err)
		}
	}

	tests := map[string]struct {
		firstName, lastName string
		expected            *User
		expectedErr         error
	}{
		"simple lookup": {
			firstName:   addUserParams[0].firstName,
			lastName:    addUserParams[0].lastName,
			expected:    &testManager.users[0],
			expectedErr: nil,
		},
		"last element lookup": {
			firstName:   addUserParams[3].firstName,
			lastName:    addUserParams[3].lastName,
			expected:    &testManager.users[3],
			expectedErr: nil,
		},
		"no match lookup": {
			firstName:   "Go",
			lastName:    "Gopher",
			expected:    nil,
			expectedErr: ErrNoUserFound,
		},
		"partial match lookup": {
			firstName:   "foo",
			lastName:    "foo",
			expected:    nil,
			expectedErr: ErrNoUserFound,
		},
		"empty names": {
			firstName:   "",
			lastName:    "",
			expected:    nil,
			expectedErr: ErrNoUserFound,
		},
	}

	for name, test := range tests {
		result, err := testManager.GetUserByName(test.firstName, test.lastName)

		if !errors.Is(err, test.expectedErr) {
			t.Errorf("%s: bad error, expected \"%v\"\n, got \"%v\"\n", name, test.expectedErr, err)
		}

		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("%s: bad result, expected %v\n, got %v\n", name, test.expected, result)
		}
	}
}
