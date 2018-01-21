package models

import (
	"fmt"
	"log"
)

type Task struct {
	ID        int      `json:"id"`
	Title     string   `json:"title"`
	Completed bool     `json:"completed"`
	Project   *Project `json:"project"`
}

func (db *DB) AddStuff() {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into tasks(title) values(?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	for i := 0; i < 10; i++ {
		_, err = stmt.Exec(fmt.Sprintf("こんにちわ世界%03d", i))
		if err != nil {
			log.Fatal(err)
		}
	}
	tx.Commit()
}

func (db *DB) AllTasks() ([]*Task, error) {
	rows, err := db.Query("select id, title from tasks")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	tasks := make([]*Task, 0)
	for rows.Next() {
		task := new(Task)
		err = rows.Scan(&task.ID, &task.Title)
		if err != nil {
			log.Fatal(err)
		}
		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (db *DB) SingleTask(taskID int) (*Task, error) {
	row := db.QueryRow(fmt.Sprintf("select id, title from tasks where id = %d", taskID))

	task := new(Task)

	err := row.Scan(&task.ID, &task.Title)

	if err != nil {
		return nil, err
	}

	return task, nil
}
