package handler

import (
	"errors"
	"fmt"
	"minitrello/pkg/errs"
	"minitrello/pkg/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// регистрация
func (h *Handler) RegisterUser(c *gin.Context) {
	var request models.UserInput
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(400, gin.H{"error: ": err.Error()})
		return
	}
	err = h.service.RegisterUser(c.Request.Context(), &models.UserInput{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		if errors.Is(err, errs.ErrEmailAlreadyExists) {
			c.JSON(409, gin.H{
				"error": "Email already exists",
			})
			return
		}
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"massage: ": "Account created!"})
}

// вход
func (h *Handler) Login(c *gin.Context) {
	var request models.UserInput
	if err := c.BindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	token, err := h.service.LoginUser(c.Request.Context(), request.Email, request.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": "invalid credentials"})
		return
	}
	fmt.Println(token)
	c.JSON(200, gin.H{"token": token})
}

// получение пользователя по почте
func (h *Handler) GetByEmail(c *gin.Context) {
	email := c.Param("email")

	user, err := h.service.GetByEmail(c.Request.Context(), email)
	if err != nil {
		c.JSON(404, gin.H{"error": "user not found"})
		return
	}

	c.JSON(200, gin.H{
		"id":         user.ID,
		"name":       user.Name,
		"email":      user.Email,
		"role":       user.Role,
		"created_at": user.CreatedAt,
	})
}

// получение пользователя по ID
func (h *Handler) GetByID(c *gin.Context) {
	SID := c.Param("id")
	id, _ := strconv.Atoi(SID)

	user, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(404, gin.H{"error": "user not found"})
		return
	}

	c.JSON(200, gin.H{
		"id":         user.ID,
		"name":       user.Name,
		"email":      user.Email,
		"role":       user.Role,
		"created_at": user.CreatedAt,
	})
}

func (h *Handler) UpdateUserName(c *gin.Context) {
	var input struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "uncorrect name input"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	err := h.service.UpdateUserName(c.Request.Context(), userID.(int), input.Name)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to update name"})
		return
	}

	c.JSON(200, gin.H{
		"message":  "name updated successfully",
		"new_name": input.Name,
	})
}

// удаление пользователя
func (h *Handler) Delete(c *gin.Context) {
	SID := c.Param("id")
	id, err := strconv.Atoi(SID)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid user ID",
		})
		return
	}
	err = h.service.DeleteUser(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			c.JSON(404, gin.H{
				"error": "user not found",
			})
			return
		}
		c.JSON(500, gin.H{
			"error": "internal server error",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "User deleted succesfully!",
	})

}
