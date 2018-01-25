package models

import (
	"log"
)

// A User is a user of NozzyTasks
type User struct {
	ID    int64  `json:"id"`
	Sub   string `json:"sub"`
	Email string `json:"email"`
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
