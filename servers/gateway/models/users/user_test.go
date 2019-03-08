package users

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
	"testing"
	"reflect"
)

//TODO: add tests for the various functions in user.go, as described in the assignment.
//use `go test -cover` to ensure that you are covering all or nearly all of your code paths.

func setUserPassHashAndPhoto(nu *NewUser, u *User) {
	u.SetPassword(nu.Password)
	trimEmail := strings.TrimSpace(nu.Email)
	trimEmail = strings.ToLower(trimEmail)
	hash := md5.New()
	hash.Write([]byte(trimEmail))
	u.PhotoURL = gravatarBasePhotoURL + hex.EncodeToString(hash.Sum(nil))
}

func TestValidate(t *testing.T) {
	cases := []struct {
		caseMessage   string
		newUser       NewUser
		expectedError bool
	}{
		{
			"Good Case Scenario",
			NewUser{
				Email:        "hello@example.com",
				Password:     "password",
				PasswordConf: "password",
				UserName:     "nehay100",
				FirstName:    "Neha",
				LastName:     "Yadav",
			},
			false,
		},
		{
			"Bad Email Scenario",
			NewUser{
				Email:        "hello",
				Password:     "password",
				PasswordConf: "password",
				UserName:     "nehay100",
				FirstName:    "Neha",
				LastName:     "Yadav",
			},
			true,
		},
		{
			"Bad Password Scenario",
			NewUser{
				Email:        "hello@example.com",
				Password:     "pass",
				PasswordConf: "pass",
				UserName:     "nehay100",
				FirstName:    "Neha",
				LastName:     "Yadav",
			},
			true,
		},
		{
			"Bad Password Confirm Scenario",
			NewUser{
				Email:        "hello@example.com",
				Password:     "password",
				PasswordConf: "pass",
				UserName:     "nehay100",
				FirstName:    "Neha",
				LastName:     "Yadav",
			},
			true,
		},
		{
			"Bad Username Scenario",
			NewUser{
				Email:        "hello@example.com",
				Password:     "password",
				PasswordConf: "password",
				UserName:     "",
				FirstName:    "Neha",
				LastName:     "Yadav",
			},
			true,
		},
		{
			"Bad Username 2 Scenario", 
			NewUser{
				Email:        "hello@example.com",
				Password:     "pass",
				PasswordConf: "password",
				UserName:     "nehay 100",
				FirstName:    "Neha",
				LastName:     "Yadav",
			},
			true,
		},
		{
			"Bad Username 3 Scenario", 
			NewUser{
				Email:        "hello@example.com",
				Password:     "pass",
				PasswordConf: "password",
				UserName:     " nehay100",
				FirstName:    "Neha",
				LastName:     "Yadav",
			},
			true,
		},
		{
			"Bad Username 4 Scenario", 
			NewUser{
				Email:        "hello@example.com",
				Password:     "pass",
				PasswordConf: "password",
				UserName:     "nehay100 ",
				FirstName:    "Neha",
				LastName:     "Yadav",
			},
			true,
		},
		{
			"Bad Username 5 Scenario", 
			NewUser{
				Email:        "hello@example.com",
				Password:     "pass",
				PasswordConf: "password",
				UserName:     " ",
				FirstName:    "Neha",
				LastName:     "Yadav",
			},
			true,
		},
	}

	for _, c := range cases {
		err := c.newUser.Validate()

		if !c.expectedError && err != nil {
			t.Errorf("%s: expected: no error, got: %v\n", c.caseMessage, err)
		}
	}
}

func TestToUser(t *testing.T) {
	// Good Case Scenario
	newUser := &NewUser{
		Email:        "hello@example.com",
		Password:     "password",
		PasswordConf: "password",
		UserName:     "nehay100",
		FirstName:    "Neha",
		LastName:     "Yadav",
	}

	user := &User{
		ID:        0,
		Email:     "hello@example.com",
		PassHash:  []byte{},
		UserName:  "nehay100",
		FirstName: "Neha",
		LastName:  "Yadav",
		PhotoURL:  "",
	}
	setUserPassHashAndPhoto(newUser, user)

	cases := []struct {
		caseMessage   string
		newUser       *NewUser
		expectedError bool
	}{
		{
			"Good Case Scenario",
			newUser,
			false,
		},
	}

	for _, c := range cases {
		output, err := c.newUser.ToUser()

		if !c.expectedError && err != nil {
			t.Errorf("%s: expected: no error, got: %v, newUser instance:%v, user instance:%v\n", c.caseMessage, err, c.newUser, user)
		}
		if newUser.FirstName != output.FirstName {
			t.Errorf("%s: the user's first name was save incorrectly", c.caseMessage)
		}
		if newUser.LastName != output.LastName {
			t.Errorf("%s: the user's last name was save incorrectly", c.caseMessage)
		}
		if newUser.Email != output.Email {
			t.Errorf("%s: the user's email was save incorrectly", c.caseMessage)
		}
		if newUser.UserName != output.UserName {
			t.Errorf("%s: the user's username was save incorrectly", c.caseMessage)
		}
		if !reflect.DeepEqual(user.PhotoURL, output.PhotoURL) {
			t.Errorf("%s: the user's photoURL was save incorrectly", c.caseMessage)
		}
	}
}

func TestFullName(t *testing.T) {
	cases := []struct {
		caseMessage   string
		user          User
		expectedError bool
		expectedAnswer string
	}{
		{
			"Good Case Scenario",
			User{
				ID:        0,
				Email:     "hello@example.com",
				PassHash:  []byte{},
				UserName:  "nehay100",
				FirstName: "Neha",
				LastName:  "Yadav",
				PhotoURL:  "hash",
			},
			false,
			"Neha Yadav",
		},
		{
			"Bad First Name Scenario",
			User{
				ID:        0,
				Email:     "hello@example.com",
				PassHash:  []byte{},
				UserName:  "nehay100",
				FirstName: "",
				LastName:  "Yadav",
				PhotoURL:  "hash",
			},
			true,
			"Yadav",
		},
		{
			"Bad Last Name Scenario",
			User{
				ID:        0,
				Email:     "hello@example.com",
				PassHash:  []byte{},
				UserName:  "nehay100",
				FirstName: "Neha",
				LastName:  "",
				PhotoURL:  "hash",
			},
			true,
			"Neha",
		},
		{
			"Bad First and Last Name Scenario",
			User{
				ID:        0,
				Email:     "hello@example.com",
				PassHash:  []byte{},
				UserName:  "nehay100",
				FirstName: "",
				LastName:  "",
				PhotoURL:  "hash",
			},
			true,
			"",
		},
	}

	for _, c := range cases {
		output := c.user.FullName()
		if c.expectedAnswer != output {
			t.Errorf("%s: expected string(%s) didn't match actual string(%s)", c.caseMessage, c.expectedAnswer, output)
		}
		if (len(c.user.FirstName) == 0 && len(c.user.LastName) == 0) && output != ""  {
			t.Errorf("%s: full name string is supposed to be empty", c.caseMessage)
		}
	}
}

func TestApplyUpdates(t *testing.T) {
	goodUser := &User{
		ID: 0,
		Email: "nehay100@gmail.com",
		PassHash: []byte{},
		UserName: "nehay100",
		FirstName: "Neha",
		LastName:"Yadav",
		PhotoURL: "hash",
	}

	goodUpdate := &Updates{
		"Neha",
		"Yadav",
	}

	badFirst := &Updates{
		"",
		"Yadav",
	}

	badLast := &Updates{
		"Neha",
		"",
	}

	badBoth := &Updates{
		"",
		"",
	}

	cases := []struct {
		caseMessage   string
		update		  *Updates
		user		  *User
		expectedError bool
	}{
		{
			"Good Case Scenario",
			goodUpdate,
			goodUser,
			false,
		},
		{
			"Bad First Name Update",
			badFirst,
			goodUser,
			true,
		},
		{
			"Bad Last Name Update",
			badLast,
			goodUser,
			true,
		},
		{
			"Bad First and Last Name Update",
			badBoth,
			goodUser,
			true,
		},
	}

	for _, c := range cases {
		err := c.user.ApplyUpdates(c.update)

		if !c.expectedError && err != nil {
			t.Errorf("%s: expected: no error, got: %v, user instance:%v\n", c.caseMessage, err, c.user)
		}
	}
}

func TestAuthenticate(t *testing.T) {
	// Good Case Scenario
	goodNewUser := &NewUser{
		Email:        "hello@example.com",
		Password:     "password",
		PasswordConf: "password",
		UserName:     "nehay100",
		FirstName:    "Neha",
		LastName:     "Yadav",
	}
	goodUser, _ := goodNewUser.ToUser()

	badHashUser := &User{
		ID:        0,
		Email:     "hello@example.com",
		PassHash:  []byte{},
		UserName:  "nehay100",
		FirstName: "Neha",
		LastName:  "Yadav",
		PhotoURL:  "",
	}
	
	cases := []struct {
		caseMessage   string
		newUser		  *NewUser
		user		  *User
		expectedError bool
	}{
		{
			"Good Case Scenario",
			goodNewUser,
			goodUser,
			false,
		},
		{
			"No Matching Hash Scenario",
			goodNewUser,
			badHashUser,
			true,
		},
	}

	for _, c := range cases {
		err := c.user.Authenticate(c.newUser.Password)

		if !c.expectedError && err != nil {
			t.Errorf("%s: expected: no error, got: %v\n", c.caseMessage, err)
		}
	}
}
