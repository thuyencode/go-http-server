package users

import (
	"errors"
	"fmt"
	"log/slog"
	"net/mail"
	"time"
)

var ErrNoUserFound = errors.New("no such user found")
var ErrEmptyNames = errors.New("you cannot leave first name or last name empty")
var ErrDuplicatedName = errors.New("duplicated name")

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
		return ErrEmptyNames
	}

	existedUser, err := m.GetUserByName(firstName, lastName)

	if err != nil && !errors.Is(err, ErrNoUserFound) {
		return fmt.Errorf("error checking if user is already existed: %v", err)
	}

	if existedUser != nil {
		return ErrDuplicatedName
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

	return nil, ErrNoUserFound
}

func (m *Manager) Shutdown() {
	slog.Info("User manager shutting down")
	time.Sleep(2 * time.Second)
	slog.Info("User manager shutdown completed")
}
