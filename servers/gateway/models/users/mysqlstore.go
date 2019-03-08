package users

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/nehay100/assignments-nehay100/servers/gateway/indexes"
)

// MySQLStore represents a users.Store backed by MYSQL
type SqlStore struct {
	db *sql.DB
}

const sqlSelectUser = "select id, email, userName, passHash, firstName, lastName, photoURL from users"
const sqlInitSelectID = sqlSelectUser + " where id = ?"
const sqlInitSelectEmail = sqlSelectUser + " where email = ?"
const sqlInitSelectUserName = sqlSelectUser + " where username = ?"
const update = "update users set firstName=?, lastName=? where id=?"
const delete = "delete from users where id=?"
const insertUser = "insert into users (email, userName, passHash, firstName, lastName, photoURL) values (?, ?, ?, ?, ?, ?)"
const insertSignIn = "insert into signin (id, signingTimeDate, ipAddress) values (?, ?, ?)"

// NewMySqlStore constructs a new MySQLStore.
// it will panic if the db pointer is nil.
func NewSqlStore(db *sql.DB) *SqlStore {
	if db == nil {
		panic("empty database pointer")
	}
	return &SqlStore{
		db: db,
	}
}
func (s *SqlStore) get(param string, value interface{}) (*User, error) {
	query := fmt.Sprintf(sqlSelectUser+" where %v = ?", param)
	row := s.db.QueryRow(query, value)

	u := &User{}

	if err := row.Scan(&u.ID, &u.Email, &u.UserName, &u.PassHash,
		&u.FirstName, &u.LastName, &u.PhotoURL); err != nil {
		return nil, err
	}

	return u, nil
}

// GetByID returns the User with the given ID
func (s *SqlStore) GetByID(id int64) (*User, error) {
	return s.get("id", id)
}

// GetByEmail returns the User with the given email
func (s *SqlStore) GetByEmail(email string) (*User, error) {
	return s.get("email", email)
}

// GetByUserName returns the User with the given Username
func (s *SqlStore) GetByUserName(username string) (*User, error) {
	return s.get("username", username)
}

// Insert inserts the user into the database, and returns
// the newly-inserted User, complete with the DBMS-assigned ID
func (s *SqlStore) Insert(user *User) (*User, error) {
	//insert a new row into the "contacts" table
	//use ? markers for the values to defeat SQL
	//injection attacks

	// insq := "insert into users (email, userName, passHash, firstName, lastName, photoURL) values (?, ?, ?, ?, ?, ?)"
	res, err := s.db.Exec(insertUser, user.Email, user.UserName, user.PassHash, user.FirstName, user.LastName, user.PhotoURL)
	if err != nil {
		return nil, fmt.Errorf("error inserting new row: %v", err)
	}
	//get the auto-assigned ID for the new row
	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting new ID: %v, err: %v", id, err)
	}

	user.ID = id

	return user, nil
}

// Insert inserts the user into the database, and returns
// the newly-inserted User, complete with the DBMS-assigned ID
func (s *SqlStore) SignInInsert(userSignIn *UserSignIn) (*UserSignIn, error) {
	//insert a new row into the "contacts" table
	//use ? markers for the values to defeat SQL
	//injection attacks
	// insq := "insert into signin (id, signingTimeDate, ipAddress) values (?, ?, ?)"
	res, err := s.db.Exec(insertSignIn, userSignIn.ID, userSignIn.DateTime, userSignIn.IPAddress)
	if err != nil {
		return nil, fmt.Errorf("error inserting new row: %v", err)
	} else {
		//get the auto-assigned ID for the new row
		id, err := res.LastInsertId()
		if err != nil {
			return nil, fmt.Errorf("error getting new ID: %v, err: %v", id, err)
		}

		userSignIn.ID = id
	}
	return userSignIn, nil
}

func (s *SqlStore) AddAllToTrie() (*indexes.Trie, error) {
	rows, err := s.db.Query("select * from users")
	if err != nil {
		return nil, fmt.Errorf("error when getting all users: %v", err)
	}
	// row := s.db.QueryRow(sqlInitSelectID, user.ID)
	defer rows.Close()
	// row.Next
	u := &User{}
	for rows.Next() {
		if err := rows.Scan(&u.ID, &u.UserName, &u.FirstName, &u.LastName); err != nil {
			return nil, err
		}
	}

	username := strings.ToLower(u.UserName)
	firstname := strings.ToLower(u.FirstName)
	lastname := strings.ToLower(u.LastName)

	var wordSplit []string
	if strings.Contains(username, " ") {
		wordSplit = strings.Split(username, " ")
	} else if strings.Contains(firstname, " ") {
		wordSplit = strings.Split(firstname, " ")
	} else if strings.Contains(lastname, " ") {
		wordSplit = strings.Split(lastname, " ")
	}

	trie := indexes.NewTrie()
	for _, elem := range wordSplit {
		trie.Add(elem, u.ID)
	}
	return trie, nil
}

// Update applies UserUpdates to the given user ID
// and returns the newly-updated user
func (s *SqlStore) Update(id int64, updates *Updates) (*User, error) {
	//insert a new row into the "contacts" table
	//use ? markers for the values to defeat SQL
	//injection attacks
	// updq := "update users set firstName=?, lastName=? where id=?"
	res, err := s.db.Exec(update, updates.FirstName, updates.LastName, id)
	if err != nil {
		return nil, fmt.Errorf("error updating row: %v", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("error when trying to get affected rows: %v", err)
	}

	// Update not made in table
	if rowsAffected == 0 {
		return nil, err
	}

	return s.GetByID(id)
}

// Delete deletes the user with the given ID
func (s *SqlStore) Delete(id int64) error {
	_, err := s.db.Exec(delete, id)
	if err != nil {
		return fmt.Errorf("error when deleting user with id (%d), err: %v", id, err)
	}
	return nil
}
