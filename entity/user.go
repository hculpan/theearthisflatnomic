package entity

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User contains data for a user of the game
type User struct {
	gorm.Model
	FullName    string
	DisplayName string
	Username    string
	Password    []byte
	LastLogin   time.Time
}

// FindUserByUsername finds the user based on the username
func FindUserByUsername(username string) *User {
	var user *User = &User{Username: username}
	db.Where(user).First(user)
	return user
}

// AddNewUser creates a new user and inserts it
// into the DB
func AddNewUser(fullname, displayname, username, password string) (*User, error) {
	if FindUserByUsername(username).ID > 0 {
		return nil, fmt.Errorf("User already registered for userid '%s'", username)
	}

	result := &User{
		Username:    username,
		FullName:    fullname,
		DisplayName: displayname,
		LastLogin:   time.Now(),
	}
	if err := result.HashAndSalt(password); err != nil {
		return nil, err
	}

	return result, result.Save()
}

// HashAndSalt converts password to hashed value
func (u *User) HashAndSalt(value string) error {
	var err error
	u.Password, err = bcrypt.GenerateFromPassword([]byte(value), bcrypt.MinCost)
	return err
}

// VerifyPassword will validate the given password
func (u User) VerifyPassword(password string) error {
	return bcrypt.CompareHashAndPassword(u.Password, []byte(password))
}

// Insert inserts the record into the db
// This does not check to see if the records
// already exists; it assumes it does not.
func (u User) Insert() error {
	return db.Create(&u).Error
}

// Save will insert or update the record,
// based on whether it has a primary key or
// not
func (u User) Save() error {
	return db.Save(&u).Error
}
