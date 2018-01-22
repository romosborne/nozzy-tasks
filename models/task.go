package models

import (
	"fmt"
	"log"
)

type Task struct {
	ID        int64    `json:"id"`
	Title     string   `json:"title"`
	Completed bool     `json:"completed"`
	Project   *Project `json:"project"`
}

func (db *DB) AllTasks() ([]*Task, error) {
	rows, err := db.Query("select t.id, t.title, t.completed, p.id, p.name from tasks t left join projects p on t.project = p.id")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	tasks := make([]*Task, 0)
	for rows.Next() {
		task := new(Task)
		project := new(Project)
		err = rows.Scan(&task.ID, &task.Title, &task.Completed, &project.ID, &project.Name)
		if err != nil {
			log.Fatal(err)
		}

		task.Project = project
		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (db *DB) SingleTask(taskID int64) (*Task, error) {
	row := db.QueryRow(fmt.Sprintf("select t.id, t.title, t.completed, p.id, p.name from tasks t left join projects p on t.project = p.id where t.id = %d", taskID))

	task := new(Task)
	project := new(Project)
	err := row.Scan(&task.ID, &task.Title, &task.Completed, &project.ID, &project.Name)

	if err != nil {
		return nil, err
	}

	task.Project = project

	return task, nil
}

func (db *DB) CreateTask(task *Task) error {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into tasks(title, project, completed) values(?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var projectID string
	if task.Project == nil {
		projectID = "null"
	} else {
		projectID = string(task.Project.ID)
	}

	result, err := stmt.Exec(task.Title, projectID, task.Completed)
	if err != nil {
		log.Fatal(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	tx.Commit()

	task.ID = id

	return nil
}
