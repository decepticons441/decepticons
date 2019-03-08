package users

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"
	"database/sql"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestNewSqlStore(t *testing.T) {
	defer func() {
        if r := recover(); r == nil {
            t.Errorf("The code did not panic")
        }
    }()

    // The following is the code under test
    NewSqlStore(nil)
}
func TestInsert(t *testing.T) {
	// Good Case Scenario
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("error creating mock database, err: %s", err)
	}

	// ensure that mock database is closed after the test ends
	defer db.Close()
	store := NewSqlStore(db)

	newUser := &NewUser{
		Email:        "hello@example.com",
		Password:     "password",
		PasswordConf: "password",
		UserName:     "nehay100",
		FirstName:    "Neha",
		LastName:     "Yadav",
	}
	user, err := newUser.ToUser()
	if err != nil {
		t.Errorf("error when converting newUser to user: %v", err)
	}
	
	mock.ExpectExec(regexp.QuoteMeta("insert into users (email, userName, passHash, firstName, lastName, photoURL)")).
		WithArgs(user.Email, user.UserName, user.PassHash, user.FirstName, user.LastName, user.PhotoURL).
		WillReturnResult(sqlmock.NewResult(2, 1))

	addedInsert, err := store.Insert(user)
	if err != nil {
		t.Errorf("error when inserting new user to mock db: %v", err)
	}

	if !reflect.DeepEqual(addedInsert, user) {
		t.Errorf("addedUser fields don't match information inserted: prevUser(%v) and addedUser(%v), err: %v", user, addedInsert, err)
	}

	user.ID = -1

	mock.ExpectExec(regexp.QuoteMeta("insert into users (email, userName, passHash,  firstName, lastName, photoURL)")).
		WithArgs(user.Email, user.UserName, user.PassHash, user.FirstName, user.LastName, user.PhotoURL).
		WillReturnError(fmt.Errorf("error when executing insert method"))
	_, err = store.Insert(user)
	if err == nil {
		t.Errorf("expecting an error when inserting badUser(%v): %v", user, err)
	}
	err2 := mock.ExpectationsWereMet()
	if err2 != nil {
		t.Errorf("unmet mock sql expectations: %v", err)
	}
}

func TestUpdate(t *testing.T) {
	// Good Case Scenario
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("error creating mock database, err: %s", err)
	}

	// ensure that mock database is closed after the test ends
	defer db.Close()
	store := NewSqlStore(db)

	newUser := &NewUser{
		Email:        "hello@example.com",
		Password:     "password",
		PasswordConf: "password",
		UserName:     "nehay100",
		FirstName:    "Neha",
		LastName:     "Yadav",
	}
	user, err := newUser.ToUser()
	if err != nil {
		t.Errorf("error when converting newUser to user: %v", err)
	}
	update2 := &Updates{
		FirstName: "Lorelai",
		LastName:  "Gilmore",
	}

	col := []string{"id", "email", "userName", "passHash",  "firstName", "lastName", "photoURL"}

	rows := sqlmock.NewRows(col)
	rows.AddRow(user.ID, user.Email, user.UserName, user.PassHash, 
		user.FirstName, user.LastName, user.PhotoURL)
	mock.ExpectExec(regexp.QuoteMeta("update users set firstName=?, lastName=? where id=?")).
		WithArgs(update2.FirstName, update2.LastName, user.ID).
		WillReturnResult(sqlmock.NewResult(2, 1))

	mock.ExpectQuery(sqlInitSelectID).WithArgs(user.ID).WillReturnRows(rows)
	updatedUser, err := store.Update(user.ID, update2)
	if err != nil {
		t.Errorf("error when updating new user in mock db: %v", err)
	}

	if !reflect.DeepEqual(updatedUser, user) {
		t.Errorf("updatedUser fields don't match information updated: prevUser(%v) and updatedUser(%v), err: %v", user, updatedUser, err)
	}

	// Bad Case Scenario #1: Updates are the same
	badUser := &User{
		ID:        -1,
		Email:     "hello@example.com",
		UserName:  "nehay100",
		PassHash:  []byte{},
		FirstName: "Neha",
		LastName:  "Yadav",
		PhotoURL:  "",
	}

	update2.FirstName = "Neha"
	update2.LastName = "Yadav"

	mock.ExpectExec(regexp.QuoteMeta("update users set firstName=?, lastName=? where id=?")).
		WithArgs(update2.FirstName, update2.LastName, user.ID).
		WillReturnResult(sqlmock.NewResult(3, 0))
	_, err = store.Update(user.ID, update2)
	if err != nil {
		t.Errorf("expecting an error when updating badUser(%v): %v", badUser, err)
	}

	// Bad Case Scenario #2: Updates are Invalid
	update2.FirstName = ""
	update2.LastName = ""

	mock.ExpectExec(regexp.QuoteMeta("update users set firstName=?, lastName=? where id=?")).
		WithArgs(update2.FirstName, update2.LastName, user.ID).
		WillReturnResult(sqlmock.NewResult(3, 0))
	_, err = store.Update(user.ID, update2)
	if err != nil {
		t.Errorf("expecting an error when updating badUser(%v): %v", badUser, err)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unmet mock sql expectations: %v", err)
	}
}

func TestDelete(t *testing.T) {
	// Good Case Scenario
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("error creating mock database, err: %s", err)
	}

	// ensure that mock database is closed after the test ends
	defer db.Close()
	store := NewSqlStore(db)

	newUser := &NewUser{
		Email:        "hello@example.com",
		Password:     "password",
		PasswordConf: "password",
		UserName:     "nehay100",
		FirstName:    "Neha",
		LastName:     "Yadav",
	}
	user, err := newUser.ToUser()
	if err != nil {
		t.Errorf("error when converting newUser to user: %v", err)
	}

	col := []string{"id", "email", "userName", "passHash",  "firstName", "lastName", "photoURL"}
	
	rows := sqlmock.NewRows(col)
	rows.AddRow(user.ID, user.Email, user.UserName, user.PassHash, 
		user.FirstName, user.LastName, user.PhotoURL)
	
	mock.ExpectExec(regexp.QuoteMeta("delete from users where id=?")).
		WithArgs(user.ID).
		WillReturnResult(sqlmock.NewResult(2, 1))

	err = store.Delete(user.ID)
	if err != nil {
		t.Errorf("error when deleting user in mock db: %v", err)
	}

	// Bad Case Scenario #1: Errors were thrown

	mock.ExpectExec(regexp.QuoteMeta("delete from users where id=?")).
		WithArgs(-1).
		WillReturnError(fmt.Errorf("error when deleting a user in mock db"))
	err = store.Delete(-1)
	if err == nil {
		t.Errorf("expecting an error when deleting non-exisitant user: %v", err)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unmet mock sql expectations: %v", err)
	}
}

func TestGetByID(t *testing.T) {
	// Good Case Scenario
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("error creating mock database, err: %s", err)
	}

	// ensure that mock database is closed after the test ends
	defer db.Close()
	store := NewSqlStore(db)

	newUser := &NewUser{
		Email:        "hello@example.com",
		Password:     "password",
		PasswordConf: "password",
		UserName:     "nehay100",
		FirstName:    "Neha",
		LastName:     "Yadav",
	}
	user, err := newUser.ToUser()
	if err != nil {
		t.Errorf("error when converting newUser to user: %v", err)
	}

	col := []string{"id", "email", "userName", "passHash",  "firstName", "lastName", "photoURL"}

	rows := sqlmock.NewRows(col)
	rows.AddRow(user.ID, user.Email, user.UserName, user.PassHash, 
		user.FirstName, user.LastName, user.PhotoURL)

	mock.ExpectQuery(regexp.QuoteMeta(sqlInitSelectID)).WithArgs(user.ID).WillReturnRows(rows)
	getUser, err := store.GetByID(user.ID)
	if err != nil {
		t.Errorf("error when getting user by id in mock db: %v", err)
	}

	if !reflect.DeepEqual(getUser, user) && err == nil {
		t.Errorf("getUser fields don't match information received: prevUser(%v) and getUser(%v), err: %v", user, getUser, err)
	}

	// Bad Case Scenario #1: Getting a wrong ID
	mock.ExpectQuery(sqlInitSelectID).WithArgs(-1).
		WillReturnError(sql.ErrNoRows)
	_, err = store.GetByID(-1)
	if err == nil {
		t.Errorf("expecting an error (%v) when getting user by id in mock db when an error wasn't received", err)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unmet mock sql expectations: %v", err)
	}
}

func TestGetByEmail(t *testing.T) {
	// Good Case Scenario
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("error creating mock database, err: %s", err)
	}

	// ensure that mock database is closed after the test ends
	defer db.Close()
	store := NewSqlStore(db)

	newUser := &NewUser{
		Email:        "hello@example.com",
		Password:     "password",
		PasswordConf: "password",
		UserName:     "nehay100",
		FirstName:    "Neha",
		LastName:     "Yadav",
	}
	user, err := newUser.ToUser()
	if err != nil {
		t.Errorf("error when converting newUser to user: %v", err)
	}

	col := []string{"id", "email", "userName", "passHash", "firstName", "lastName", "photoURL"}

	rows := sqlmock.NewRows(col)
	rows.AddRow(user.ID, user.Email, user.UserName, user.PassHash, 
		user.FirstName, user.LastName, user.PhotoURL)

	mock.ExpectQuery(regexp.QuoteMeta(sqlInitSelectEmail)).WithArgs(user.Email).WillReturnRows(rows)
	getUser, err := store.GetByEmail(user.Email)
	if err != nil {
		t.Errorf("error when getting user by email in mock db: %v", err)
	}

	if !reflect.DeepEqual(getUser, user) && err == nil {
		t.Errorf("getUser fields don't match information received: prevUser(%v) and getUser(%v), err: %v", user, getUser, err)
	}

	// Bad Case Scenario #1: Getting a wrong email
	mock.ExpectQuery(sqlInitSelectEmail).WithArgs("").
		WillReturnError(sql.ErrNoRows)
	_, err = store.GetByEmail("")
	if err == nil {
		t.Errorf("expecting an error (%v) when getting user by email in mock db when an error wasn't received", err)
	}
	
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unmet mock sql expectations: %v", err)
	}
}

func TestGetByUserName(t *testing.T) {
	// Good Case Scenario
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("error creating mock database, err: %s", err)
	}

	// ensure that mock database is closed after the test ends
	defer db.Close()
	store := NewSqlStore(db)

	newUser := &NewUser{
		Email:        "hello@example.com",
		Password:     "password",
		PasswordConf: "password",
		UserName:     "nehay100",
		FirstName:    "Neha",
		LastName:     "Yadav",
	}
	user, err := newUser.ToUser()
	if err != nil {
		t.Errorf("error when converting newUser to user: %v", err)
	}

	col := []string{"id", "email", "userName", "passHash",  "firstName", "lastName", "photoURL"}

	rows := sqlmock.NewRows(col)
	rows.AddRow(user.ID, user.Email, user.UserName, user.PassHash, 
		user.FirstName, user.LastName, user.PhotoURL)

	mock.ExpectQuery(regexp.QuoteMeta(sqlInitSelectUserName)).WithArgs(user.UserName).WillReturnRows(rows)
	getUser, err := store.GetByUserName(user.UserName)
	if err != nil {
		t.Errorf("error when getting user by username in mock db: %v", err)
	}

	if !reflect.DeepEqual(getUser, user) && err == nil {
		t.Errorf("getUser fields don't match information received: prevUser(%v) and getUser(%v), err: %v", user, getUser, err)
	}

	// Bad Case Scenario #1: Getting a wrong update
	mock.ExpectQuery(sqlInitSelectUserName).WithArgs("").
		WillReturnError(sql.ErrNoRows)
	_, err = store.GetByUserName("")
	if err == nil {
		t.Errorf("expecting an error (%v) when getting user by username in mock db when an error wasn't received", err)
	}
	
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unmet mock sql expectations: %v", err)
	}
}
