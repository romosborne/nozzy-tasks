package models

import (
	"fmt"
	"log"
)

// A User is a user of NozzyTasks
type User struct {
	ID        int64
	Sub       string
	Email     string
	AuthToken string
}

// AddUser adds a user to the database
func (db *DB) AddUser(user *User) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into users(sub, email) values(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(user.Sub, user.Email)

	if err != nil {
		log.Fatal(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	tx.Commit()

	user.ID = id
}

// SetAuthToken sets the auth token of a user
func (db *DB) SetAuthToken(sub string, email string, authToken string) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare(fmt.Sprintf("update users set authToken = '%s' where sub = '%s'", authToken, sub))
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	result, err := stmt.Exec()

	if err != nil {
		log.Fatal(err)
	}

	tx.Commit()

	if rowsAffected, _ := result.RowsAffected(); rowsAffected != 1 {
		user := User{
			Sub:       sub,
			Email:     email,
			AuthToken: authToken,
		}
		db.AddUser(&user)
	}
}

// GetUserFromAuthToken returns the user from the authtoken, or an error
func (db *DB) GetUserFromAuthToken(authToken string) (*User, error) {
	row := db.QueryRow(fmt.Sprintf("select id, sub, email from users where authToken = '%s'", authToken))

	var user User
	err := row.Scan(&user.ID, &user.Sub, &user.Email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
