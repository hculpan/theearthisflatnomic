package entity

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// NextTurnAction represents the next action
// the acting user needs to take
type NextTurnAction int

// A list of actions the user takes in
// their actions
const (
	Propose NextTurnAction = iota
	SubmitForVote
)

// User contains data for a user of the game
type User struct {
	gorm.Model
	FullName    string
	DisplayName string
	Username    string
	Password    []byte
	LastLogin   time.Time
	Points      int
	IsTurn      bool
	EndOfTurn   time.Time
	Action      NextTurnAction
}

// FindUserByUsername finds the user based on the username.
// It will return true if the user was found; false otherwise.
func FindUserByUsername(username string) (*User, bool) {
	var user *User = &User{Username: username}
	db.Where(user).First(user)
	return user, user.ID != 0
}

// FindUserByDisplayName returns a user based on the
// display name.  It will return true if the user was found;
// false otherwise.
func FindUserByDisplayName(name string) (*User, bool) {
	var user *User = &User{DisplayName: name}
	db.Where(user).First(user)
	return user, user.ID != 0
}

// AddNewUser creates a new user and inserts it
// into the DB
func AddNewUser(fullname, displayname, username, password string) (*User, error) {
	if _, found := FindUserByUsername(username); found {
		return nil, fmt.Errorf("User already registered for userid '%s'", username)
	} else if _, found := FindUserByDisplayName(displayname); found {
		return nil, fmt.Errorf("Display name '%s' already taken", displayname)
	}

	result := &User{
		Username:    username,
		FullName:    fullname,
		DisplayName: displayname,
		LastLogin:   time.Now(),
		Points:      0,
		IsTurn:      false,
	}
	if err := result.HashAndSalt(password); err != nil {
		return nil, err
	}

	return result, result.Save()
}

// SetFirstTurn sets the turn to the first user
// by display name
func SetFirstTurn() (*User, error) {
	var result *User = nil
	users := []User{}
	db.Order("display_name").Find(&users)
	if len(users) > 0 {
		result = &users[0]
	} else {
		return nil, fmt.Errorf("No users")
	}

	db.Model(&User{}).Update("is_turn", false)
	result.StartTurn()
	result.Save()

	return result, nil
}

// SetNextTurn sets the turn indicator to true
// for the next user and false for all others
func SetNextTurn() (*User, error) {
	activeUser := User{IsTurn: true}
	db.Where(&activeUser).First(&activeUser)
	if activeUser.ID == 0 {
		return nil, fmt.Errorf("No active user")
	}

	users := []User{}
	db.Where("is_turn = ?", false).Order("display_name").Find(&users)
	var nextActiveUser *User = nil
	for _, v := range users {
		if v.DisplayName > activeUser.DisplayName {
			nextActiveUser = &v
			break
		}
	}

	if nextActiveUser == nil && len(users) > 0 {
		nextActiveUser = &users[0]
	} else if len(users) == 0 {
		return nil, fmt.Errorf("No non-active users.  Maybe only one user?")
	}

	db.Model(&User{}).Where("is_turn = true").Update("is_turn", false)
	nextActiveUser.StartTurn()
	nextActiveUser.Save()

	return nextActiveUser, nil
}

// StartTurn will set the turn to active for
// this user and set the EndOfTurn appropriately,
// ending at 11:59 PM EST four days in the future.
// Note this doesn't update any other user, so caller
// is responsible for making sure no other user is
// set to active
func (u *User) StartTurn() {
	u.IsTurn = true
	loc, _ := time.LoadLocation("EST")
	n := time.Now().In(loc)
	u.EndOfTurn = time.Date(n.Year(), n.Month(), n.Day(), 23, 59, 0, 0, loc)
	u.EndOfTurn = u.EndOfTurn.Add((4 * 24) * time.Hour)
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
func (u *User) Insert() error {
	return db.Create(&u).Error
}

// Save will insert or update the record,
// based on whether it has a primary key or
// not
func (u *User) Save() error {
	return db.Save(u).Error
}
