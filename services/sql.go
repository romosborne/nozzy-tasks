package services

import (
	"database/sql"
	"fmt"
	"log"

	// SQLite3 import
	_ "github.com/mattn/go-sqlite3"
	"github.com/romosborne/nozzy-tasks/models"
)

type dbTask struct {
	ID        sql.NullInt64
	Title     sql.NullString
	Completed sql.NullBool
	ProjectID sql.NullInt64
}

// Datastore interface methods
type Datastore interface {
	AllTasks(userID int64) ([]*models.Project, error)
	SingleTask(taskID int64, userID int64) (*models.Task, error)
	CreateTask(task *models.Task) error
	DeleteTask(taskID int64, userID int64) error
	SetTaskCompletion(userID int64, taskID int64, complete bool) error
	CreateProject(project *models.Project) error
	DeleteProject(projectID int64, userID int64) error
	AddUser(user *models.User)
	CreateOrSetAuthToken(sub string, email string, authToken string)
	GetUserFromAuthToken(authToken string) (*models.User, error)
}

type SQL struct {
	DB *sql.DB
}

func NewSQL(source string) (*SQL, error) {
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
		email text,
		authToken text
	);
	create table if not exists projects (
		id integer primary key,
		name text,
		userId int,
		foreign key(userId) references users(id)
	);
	create table if not exists tasks (
		id integer primary key, 
		title text not null, 
		project int,
		completed bool,
		foreign key(project) references projects(id)
	);`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, err
	}

	return &SQL{DB: db}, nil
}

// AllTasks returns all the projects and tasks inside those projects
func (sql *SQL) AllTasks(userID int64) ([]*models.Project, error) {
	rows, err := sql.DB.Query(fmt.Sprintf("select t.id, t.title, t.completed, p.id, p.name from projects p left join tasks t on t.project = p.id where p.userId = '%d'", userID))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	projects := make([]*models.Project, 0)
	for rows.Next() {
		dbTask := new(dbTask)
		project := new(models.Project)
		err = rows.Scan(&dbTask.ID, &dbTask.Title, &dbTask.Completed, &project.ID, &project.Name)
		if err != nil {
			log.Fatal(err)
		}

		if found, index := contains(projects, project.ID); found {
			if dbTask.ID.Valid {
				task := &models.Task{
					ID:        dbTask.ID.Int64,
					Title:     dbTask.Title.String,
					Completed: dbTask.Completed.Bool,
					ProjectID: project.ID,
				}
				projects[index].Tasks = append(projects[index].Tasks, task)
			}
		} else {
			if dbTask.ID.Valid {
				task := &models.Task{
					ID:        dbTask.ID.Int64,
					Title:     dbTask.Title.String,
					Completed: dbTask.Completed.Bool,
					ProjectID: project.ID,
				}
				project.Tasks = []*models.Task{task}
			} else {
				project.Tasks = make([]*models.Task, 0)
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
func (sql *SQL) SingleTask(taskID int64, userID int64) (*models.Task, error) {
	row := sql.DB.QueryRow("select t.id, t.title, t.completed, t.project from tasks t left join projects p on t.project = p.id where t.id = ? and p.userId = ?", taskID, userID)

	task := new(models.Task)
	err := row.Scan(&task.ID, &task.Title, &task.Completed, &task.ProjectID)

	if err != nil {
		return nil, err
	}
	return task, nil
}

// DeleteTask deletes a task
func (sql *SQL) DeleteTask(taskID int64, userID int64) error {
	stmt := `delete from tasks 
	where exists (
		select t.id from tasks t
		join projects p on p.id = t.project
		where t.id = ?
		and p.userId = ?)
	and id = ?`

	sql.DB.Prepare(stmt)
	_, err := sql.DB.Exec(stmt, taskID, userID, taskID)

	return err
}

func (sql *SQL) SetTaskCompletion(userID int64, taskID int64, complete bool) error {
	stmt := `update tasks set completed = ? 
	where exists (
		select t.id from tasks t
		join projects p on p.id = t.project
		where t.id = ?
		and p.userId = ?)
	and id = ?`

	sql.DB.Prepare(stmt)
	_, err := sql.DB.Exec(stmt, complete, taskID, userID, taskID)

	return err
}

// CreateTask creates a task
func (sql *SQL) CreateTask(task *models.Task) error {
	tx, err := sql.DB.Begin()
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

// CreateProject create a project
func (sql *SQL) CreateProject(project *models.Project) error {
	tx, err := sql.DB.Begin()
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

// DeleteProject deletes a project
func (sql *SQL) DeleteProject(projectID int64, userID int64) error {
	stmt := `delete from tasks 
	where t.project = ?
	and id = ?;
	delete from projects
		where id = ?`

	sql.DB.Prepare(stmt)
	_, err := sql.DB.Exec(stmt, projectIdID, userID, projectID)

	return err
}

// AddUser adds a user to the database
func (sql *SQL) AddUser(user *models.User) {
	tx, err := sql.DB.Begin()
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

// CreateOrSetAuthToken sets the auth token of a user or creates a new user if none is found
func (sql *SQL) CreateOrSetAuthToken(sub string, email string, authToken string) {
	var count int
	_ = sql.DB.QueryRow("select count(*) from users where sub = ?", sub).Scan(&count)

	if count == 0 {
		// Create new user
		sql.AddUser(&models.User{
			Sub:       sub,
			Email:     email,
			AuthToken: authToken})
		return
	}

	tx, err := sql.DB.Begin()
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
		user := models.User{
			Sub:       sub,
			Email:     email,
			AuthToken: authToken,
		}
		sql.AddUser(&user)
	}
}

// GetUserFromAuthToken returns the user from the authtoken, or an error
func (sql *SQL) GetUserFromAuthToken(authToken string) (*models.User, error) {
	row := sql.DB.QueryRow(fmt.Sprintf("select id, sub, email from users where authToken = '%s'", authToken))

	var user models.User
	err := row.Scan(&user.ID, &user.Sub, &user.Email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func contains(slice []*models.Project, projectID int64) (bool, int) {
	for index, value := range slice {
		if value.ID == projectID {
			return true, index
		}
	}

	return false, 0
}
