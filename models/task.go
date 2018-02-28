package models

// A Task is the smallest item of work to do
type Task struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
	ProjectID int64  `json:"project_id"`
}

// A TaskRequest is a request to create a task
type TaskRequest struct {
	Title          string `json:"title"`
	ProjectID      int64  `json:"project_id"`
	NewProjectName string `json:"new_project_name"`
}

// A TaskCompletionRequest is a request to change the completedness of a task
type TaskCompletionRequest struct {
	TaskID    int64 `json:"task_id"`
	Completed bool  `json:"completed"`
}
