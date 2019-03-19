package users

import (
	"fmt"
	"time"
	"github.com/nehay100/decepticons/servers/gateway/indexes"
)

type FakeMockStore struct {
	User *User
	UserSignIn *UserSignIn
	expectedErr	bool
}


func makeUser() *User {
	newUser := &NewUser {
		Email:			"nehay100@gmail.com",
		Password:		"thisPasswordIsSecure",
		PasswordConf:	"thisPasswordIsSecure",
		UserName:		"nehay100",
		FirstName:		"Neha",
		LastName:		"Yadav",
	}
	user, err := newUser.ToUser()
	if err != nil {
		fmt.Errorf("error creating user, err: %s", err)
	}
	return user
}

func makeSignInUser() *UserSignIn {
	return &UserSignIn {
		ID:	0,
		DateTime: time.Now(),
		IPAddress: "0.0.0.0",
	}
}


func NewFakeMockStore (expectedErr bool) *FakeMockStore {
	return &FakeMockStore {
		User:	makeUser(),
		UserSignIn:		makeSignInUser(),
		expectedErr:	expectedErr,
	}
}

// GetByID returns the User with the given ID
func (fms *FakeMockStore) GetByID (id int64) (*User, error) {
	if fms.expectedErr {
		return &User{}, fmt.Errorf("error when receiving fake mock store user by id")
	}
	return fms.User, nil
}

// GetByEmail returns the User with the given email
func (fms *FakeMockStore) GetByEmail(email string) (*User, error) {
	if fms.expectedErr {
		return &User{}, fmt.Errorf("error when receiving fake mock store user by id")
	}
	return fms.User, nil
}


// GetByUserName returns the User with the given username
func (fms *FakeMockStore) GetByUserName (username string) (*User, error) {
	if fms.expectedErr {
		return &User{}, fmt.Errorf("error when receiving fake mock store user by id")
	}
	return fms.User, nil
}


// Insert inserts the user into the database, and returns
// the newly-inserted User, complete with the DBMS-assigned ID
func (fms *FakeMockStore) Insert(user *User) (*User, error) {
	if fms.expectedErr {
		return &User{}, fmt.Errorf("error when receiving fake mock store user by id")
	}
	return fms.User, nil
}

// Insert inserts the user into the database, and returns
// the newly-inserted User, complete with the DBMS-assigned ID
func (fms *FakeMockStore) SignInInsert(userSignIn *UserSignIn) (*UserSignIn, error) {
	if fms.expectedErr {
		return &UserSignIn{}, fmt.Errorf("error when receiving fake mock store user by id")
	}
	return fms.UserSignIn, nil
}

// Update applies UserUpdates to the given user ID
// and returns the newly-updated user
func (fms *FakeMockStore) Update(id int64, updates *Updates) (*User, error) {
	if fms.expectedErr {
		return &User{}, fmt.Errorf("error when receiving fake mock store user by id")
	}

	err := fms.User.ApplyUpdates(updates)
	if err != nil {
		return &User{}, fmt.Errorf("error when applying fake mock store user updates")
	}
	return fms.User, nil
}

// Delete deletes the user with the given ID
func (fms *FakeMockStore) Delete(id int64) error {
	if fms.expectedErr {
		return fmt.Errorf("error when receiving fake mock store user by id")
	}
	return nil
}
func (fms *FakeMockStore) AddAllToTrie() (*indexes.Trie, error) {
	return nil, nil
}	
