package models

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (db *DB) AddUser(user *User) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into users(username, password) values(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(user.Username, hash)

	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()
}

func (db *DB) CheckPassword(username string, password string) (bool, error) {
	row := db.QueryRow(fmt.Sprintf("select password from users where username = %s", username))

	var dBHash []byte
	err := row.Scan(&dBHash)
	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword(dBHash, []byte(password))

	if err != nil {
		return false, nil
	}

	return true, nil
}
