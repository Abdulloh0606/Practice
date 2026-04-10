package projects

import "time"

type Project struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	CreatedBy int       `db:"created_by"`
}

type ProjectInput struct {
	Name string `json:"name"`
}

type AddMemberInput struct {
	UserID int `json:"user_id"`
}
