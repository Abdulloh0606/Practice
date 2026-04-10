package tasks

import "time"

type TaskInput struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	ProjectID   int       `json:"project_id"`
	AssignedTo  int       `json:"assigned_to,omitempty"`
	Deadline    time.Time `json:"deadline,omitempty"`
}
