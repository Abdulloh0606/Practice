package handler

import (
	"errors"
	"minitrello/pkg/errs"
	"minitrello/pkg/models/projects"
	"strconv"

	"github.com/gin-gonic/gin"
)

// создание проекта
func (h *Handler) CreateProject(c *gin.Context) {
	var input projects.ProjectInput
	err := c.BindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "uncorrect project name",
		})
		return
	}
	id, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{
			"error": "Unauthorized",
		})
		return
	}
	UserID := id.(int)
	projectId, err := h.service.CreateProject(c.Request.Context(), &input, UserID)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "internal server error",
		})
		return
	}
	c.JSON(201, gin.H{
		"message":      "project created!",
		"project_id":   projectId,
		"project_name": input.Name,
	})
}

func (h *Handler) GetProjectByID(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid project ID",
		})
		return
	}
	project, err := h.service.GetProjectByID(c.Request.Context(), projectID)
	if err != nil {
		if errors.Is(err, errs.ErrProjectNotFound) {
			c.JSON(404, gin.H{"error": "project not found"})
			return
		}
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if project == nil {
		c.JSON(404, gin.H{"error": "project not found"})
		return
	}
	c.JSON(200, gin.H{
		"project_name": project.Name,
		"created_by":   project.CreatedBy,
		"created_at":   project.CreatedAt,
	})
}

// добавление участника в проект
func (h *Handler) AddMember(c *gin.Context) {
	projectIDStr := c.Param("project_id")
	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid project ID",
		})
		return
	}
	var input *projects.AddMemberInput

	err = c.BindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err,
		})
	}

	role := "member"
	err = h.service.AddMember(c.Request.Context(), projectID, input.UserID, role)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "internal server error",
		})
	}
	c.JSON(200, gin.H{
		"message":    "user added in project succesfully",
		"project_id": projectID,
		"user_id":    input.UserID,
	})

}

// удаление учатсника из проекта
func (h *Handler) DeleteMember(c *gin.Context) {
	// Получаем project_id
	projectIDStr := c.Param("project_id")
	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid project ID"})
		return
	}

	// Получаем user_id
	userIDStr := c.Param("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid user ID"})
		return
	}

	err = h.service.DeleteMember(c.Request.Context(), projectID, userID)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			c.JSON(404, gin.H{"error": "user not found in project"})
			return
		}
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(200, gin.H{
		"message":    "user removed from project successfully",
		"project_id": projectID,
		"user_id":    userID,
	})
}

//удаление проекта

func (h *Handler) DeleteProject(c *gin.Context) {
	projectIDStr := c.Param("project_id")
	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid project ID"})
		return
	}

	err = h.service.DeleteProject(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(200, gin.H{
		"message":    "project deleted successfully",
		"project_id": projectID,
	})
}

// получение проектов где участвует пользователь
func (h *Handler) GetUserProjects(c *gin.Context) {
	userID := c.GetInt("user_id")

	projects, err := h.service.GetUserProjects(c.Request.Context(), userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(200, gin.H{
		"projects": projects,
	})
}
