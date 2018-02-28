package models

// Project is a collection of tasks
type Project struct {
	ID     int64   `json:"id"`
	Name   string  `json:"name"`
	UserID int64   `json:"-"`
	Tasks  []*Task `json:"tasks"`
}
