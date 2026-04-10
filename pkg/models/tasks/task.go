package tasks

import "time"

type Task struct {
	ID          int       `db:"id"`
	Name        string    `db:"name"`
	Description *string    `db:"description"`
	Comment     *string    `db:"comment"`
	Status      *string    `db:"status"`
	ProjectID   int       `db:"project_id"`
	CreatedAt   time.Time `db:"created_at"`
	AssignedTo  *int       `db:"assigned_to"`
	Deadline    *time.Time `db:"deadline"`
}
