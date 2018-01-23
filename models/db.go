package models

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Datastore interface {
	AllTasks() ([]*Task, error)
	SingleTask(taskID int64) (*Task, error)
	CreateTask(task *Task) error
	AddUser(user *User)
	CheckPassword(username string, password string) (bool, error)
}

type DB struct {
	*sql.DB
}

func NewDB(source string) (*DB, error) {
	db, err := sql.Open("sqlite3", source)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	sqlStmt := `
	create table if not exists users (
		id integer primary key,
		username text,
		password text
	);
	create table if not exists projects (
		id integer primary key,
		name text);
	create table if not exists tasks (
		id integer primary key, 
		title text not null, 
		project int,
		completed bool,
		foreign key(project) references projects(id))`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, err
	}

	myDb := DB{db}
	//myDb.AddUser(&User{Username: "username", Password: "password"})
	//myDb.AddStuff()

	return &myDb, nil
}
