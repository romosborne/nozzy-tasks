package models

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Datastore interface {
	AllTasks() ([]*Task, error)
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

	myDb.AddStuff()

	return &myDb, nil
}