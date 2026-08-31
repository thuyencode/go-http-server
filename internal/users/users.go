package users

import (
	"errors"
	"fmt"
	"net/mail"
)

var NoUserFoundError = errors.New("no such user found")
var EmptyNamesError = errors.New("you cannot leave first name or last name empty")
var DuplicatedNameError = errors.New("duplicated name")

type User struct {
	FirstName string
	LastName  string
	Email     mail.Address
}

type Manager struct {
	users []User
}

func NewManager() *Manager {
	return &Manager{
		users: nil,
	}
}

func (m *Manager) AddUser(firstName, lastName, email string) error {
	parsedAddress, err := mail.ParseAddress(email)

	if err != nil {
		return fmt.Errorf("invalid email: %s", email)
	}

	if firstName == "" || lastName == "" {
		return EmptyNamesError
	}

	existedUser, _ := m.GetUserByName(firstName, lastName)

	if existedUser != nil {
		return DuplicatedNameError
	}

	newUser := User{firstName, lastName, *parsedAddress}
	m.users = append(m.users, newUser)

	return nil
}

func (m *Manager) GetUserByName(firstName, lastName string) (*User, error) {
	for _, user := range m.users {
		if user.FirstName == firstName && user.LastName == lastName {
			return &user, nil
		}
	}

	return nil, NoUserFoundError
}
