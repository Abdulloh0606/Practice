package projects

import "time"

type ProjectMembers struct {
	ID        int       `db:"id"`
	ProjectID int       `db:"project_id"`
	UserID    int       `db:"user_id"`
	Role      string    `db:"role"`
	AddedAt   time.Time `db:"added_at"`
}
