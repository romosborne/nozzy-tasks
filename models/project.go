package models

import "log"

type Project struct {
	ID     int64   `json:"id"`
	Name   string  `json:"name"`
	UserID string  `json:"-"`
	Tasks  []*Task `json:"tasks"`
}

func (db *DB) CreateProject(project *Project) error {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into projects(name, userId) values(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(project.Name, project.UserID)
	if err != nil {
		log.Fatal(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	tx.Commit()

	project.ID = id

	return nil
}
