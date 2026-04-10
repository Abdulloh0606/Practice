package repository

import (
	"context"
	"errors"
	"minitrello/pkg/errs"
	"minitrello/pkg/models/tasks"

	"github.com/jackc/pgx/v5"
)

func (r *Repository) CreateTask(ctx context.Context, task *tasks.Task) (int, error) {
	query := `
		insert into tasks(name, description, status, project_id, assigned_to, deadline)
		values($1, $2, $3, $4, $5, $6)
		returning id;
	`
	var id int
	err := r.postgres.QueryRow(ctx, query,
		task.Name,
		task.Description,
		task.Status,
		task.ProjectID,
		task.AssignedTo,
		task.Deadline,
	).Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *Repository) GetTaskByID(ctx context.Context, taskID int) (*tasks.Task, error) {
	query := `
        select id, name, description, comment, status, project_id, assigned_to, deadline, created_at
        from tasks
        where id=$1
    `
	var task tasks.Task
	err := r.postgres.QueryRow(ctx, query, taskID).Scan(
		&task.ID,
		&task.Name,
		&task.Description,
		&task.Comment,
		&task.Status,
		&task.ProjectID,
		&task.AssignedTo,
		&task.Deadline,
		&task.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrTaskNotFound
		}
		return nil, err
	}
	return &task, nil
}

func (r *Repository) UpdateTask(ctx context.Context, task *tasks.Task) error {
    query := `
        update tasks
        set name=$1,
            description=$2,
            status=$3,
            assigned_to=$4,
            deadline=$5
        where id=$6
    `
    _, err := r.postgres.Exec(ctx, query,
        task.Name,
        task.Description,
        task.Status,
        task.AssignedTo,
        task.Deadline,
        task.ID,
    )
    return err
}

func (r *Repository) DeleteTask(ctx context.Context, taskID int) error {
  
    query := `select 1 FROM tasks where id=$1`
    var exist int
    err := r.postgres.QueryRow(ctx, query, taskID).Scan(&exist)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return errs.ErrTaskNotFound
        }
        return err
    }

    query = `delete from tasks where id=$1`
    _, err = r.postgres.Exec(ctx, query, taskID)
    if err != nil {
        return err
    }

    return nil
}

func (r *Repository) ListTasksByUser(ctx context.Context, userID int) ([]*tasks.Task, error) {
    query := `
        select id, name, description, comment, status, project_id, assigned_to, deadline, created_at
        from tasks
        where assigned_to = $1
        order by created_at desc
    `
    rows, err := r.postgres.Query(ctx, query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var tasksList []*tasks.Task

    for rows.Next() {
        var t tasks.Task
        err := rows.Scan(
            &t.ID,
            &t.Name,
            &t.Description,
            &t.Comment,
            &t.Status,
            &t.ProjectID,
            &t.AssignedTo,
            &t.Deadline,
            &t.CreatedAt,
        )
        if err != nil {
            return nil, err
        }
        tasksList = append(tasksList, &t)
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return tasksList, nil
}
