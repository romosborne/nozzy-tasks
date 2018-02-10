package models

import (
	"database/sql"
	"fmt"
	"log"
)

// A Task is the smallest item of work to do
type Task struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
	ProjectID int64  `json:"project_id"`
}

type TaskRequest struct {
	Title          string `json:"title"`
	ProjectID      int64  `json:"project_id"`
	NewProjectName string `json:"new_project_name"`
}

type TaskCompletionRequest struct {
	TaskID    int64 `json:"task_id"`
	Completed bool  `json:"completed"`
}

type dbTask struct {
	ID        sql.NullInt64
	Title     sql.NullString
	Completed sql.NullBool
	ProjectID sql.NullInt64
}

func contains(slice []*Project, projectID int64) (bool, int) {
	for index, value := range slice {
		if value.ID == projectID {
			return true, index
		}
	}

	return false, 0
}

// AllTasks returns all the projects and tasks inside those projects
func (db *DB) AllTasks(userID int64) ([]*Project, error) {
	rows, err := db.Query(fmt.Sprintf("select t.id, t.title, t.completed, p.id, p.name from projects p left join tasks t on t.project = p.id where p.userId = '%d'", userID))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	projects := make([]*Project, 0)
	for rows.Next() {
		dbTask := new(dbTask)
		project := new(Project)
		err = rows.Scan(&dbTask.ID, &dbTask.Title, &dbTask.Completed, &project.ID, &project.Name)
		if err != nil {
			log.Fatal(err)
		}

		if found, index := contains(projects, project.ID); found {
			if dbTask.ID.Valid {
				task := &Task{
					ID:        dbTask.ID.Int64,
					Title:     dbTask.Title.String,
					Completed: dbTask.Completed.Bool,
					ProjectID: project.ID,
				}
				projects[index].Tasks = append(projects[index].Tasks, task)
			}
		} else {
			if dbTask.ID.Valid {
				task := &Task{
					ID:        dbTask.ID.Int64,
					Title:     dbTask.Title.String,
					Completed: dbTask.Completed.Bool,
					ProjectID: project.ID,
				}
				project.Tasks = []*Task{task}
			} else {
				project.Tasks = make([]*Task, 0)
			}
			projects = append(projects, project)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return projects, nil
}

// SingleTask returns a single task
func (db *DB) SingleTask(taskID int64, userID int64) (*Task, error) {
	row := db.QueryRow("select t.id, t.title, t.completed, t.project from tasks t left join projects p on t.project = p.id where t.id = ? and p.userId = ?", taskID, userID)

	task := new(Task)
	err := row.Scan(&task.ID, &task.Title, &task.Completed, &task.ProjectID)

	if err != nil {
		return nil, err
	}
	return task, nil
}

func (db *DB) SetTaskCompletion(userID int64, taskID int64, complete bool) error {

	stmt := `update tasks set completed = ? 
	where exists (
		select t.id from tasks t
		join projects p on p.id = t.project
		where t.id = ?
		and p.userId = ?)
	and id = ?`

	db.Prepare(stmt)
	_, err := db.Exec(stmt, complete, taskID, userID, taskID)

	return err
}

// CreateTask creates a task
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

	result, err := stmt.Exec(task.Title, task.ProjectID, task.Completed)
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
