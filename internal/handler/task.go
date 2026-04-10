package handler

import (
	"errors"
	"minitrello/pkg/errs"
	"minitrello/pkg/models/tasks"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateTask(c *gin.Context) {
	var input tasks.TaskInput
	if err := c.BindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	taskID, err := h.service.CreateTask(c.Request.Context(), &input)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to create task"})
		return
	}

	c.JSON(201, gin.H{
		"message":    "task created",
		"task_id":    taskID,
		"task_name":  input.Name,
		"project_id": input.ProjectID,
	})
}

func (h *Handler) GetTaskByID(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid task ID"})
		return
	}

	task, err := h.service.GetTaskByID(c.Request.Context(), taskID)
	if err != nil {
		if errors.Is(err, errs.ErrTaskNotFound) {
			c.JSON(404, gin.H{"error": "task not found"})
			return
		}
		c.JSON(500, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"id":          task.ID,
		"name":        task.Name,
		"description": task.Description,
		"comment":     task.Comment,
		"status":      task.Status,
		"project_id":  task.ProjectID,
		"assigned_to": task.AssignedTo,
		"deadline":    task.Deadline,
		"created_at":  task.CreatedAt,
	})
}

func (h *Handler) UpdateTask(c *gin.Context) {

	taskIDStr := c.Param("id")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid task ID"})
		return
	}

	var input tasks.TaskInput
	if err := c.BindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// создаём объект tasks.Task с обновлёнными полями
	task := &tasks.Task{
		ID:          taskID,
		Name:        input.Name,
		Description: &input.Description,
		Status:      &input.Status,
		ProjectID:   input.ProjectID,
		AssignedTo:  &input.AssignedTo,
		Deadline:    &input.Deadline,
	}

	if err := h.service.UpdateTask(c.Request.Context(), task); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "task updated successfully",
		"task_id": taskID,
	})
}

func (h *Handler) DeleteTask(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid task ID"})
		return
	}

	err = h.service.DeleteTask(c.Request.Context(), taskID)
	if err != nil {
		if errors.Is(err, errs.ErrTaskNotFound) {
			c.JSON(404, gin.H{"error": "task not found"})
			return
		}
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(200, gin.H{"message": "task deleted successfully"})
}

func (h *Handler) ListTasksByUser(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(401, gin.H{"error": "unauthorized"})
        return
    }

    tasksList, err := h.service.ListTasksByUser(c.Request.Context(), userID.(int))
    if err != nil {
        c.JSON(500, gin.H{"error": "internal server error"})
        return
    }

    var response []gin.H
    for _, t := range tasksList {
        response = append(response, gin.H{
            "id":          t.ID,
            "name":        t.Name,
            "description": t.Description,
            "comment":     t.Comment,
            "status":      t.Status,
            "project_id":  t.ProjectID,
            "assigned_to": t.AssignedTo,
            "deadline":    t.Deadline,
            "created_at":  t.CreatedAt,
        })
    }

    c.JSON(200, response)
}
