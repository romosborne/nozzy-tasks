package models

import (
	"database/sql"

	// SQLite3 import
	_ "github.com/mattn/go-sqlite3"
)

// Datastore interface methods
type Datastore interface {
	AllTasks(userID string) ([]*Project, error)
	SingleTask(taskID int64, userID string) (*Task, error)
	CreateTask(task *Task) error
	CreateProject(project *Project) error
	AddUser(user *User)
}

// DB is a Custom DB for Datastore interface
type DB struct {
	*sql.DB
}

// NewDB initializes the DB
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
		sub text,
		email text
	);
	create table if not exists projects (
		id integer primary key,
		name text,
		userId int,
		foreign key(userId) references users(id));
	create table if not exists tasks (
		id integer primary key, 
		title text not null, 
		project int,
		completed bool,
		foreign key(project) references projects(id));`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, err
	}

	myDb := DB{db}

	return &myDb, nil
}
