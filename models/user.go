package models

// A User is a user of NozzyTasks
type User struct {
	ID        int64
	Sub       string
	Email     string
	AuthToken string
}
