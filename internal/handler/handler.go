package handler

import (
	"context"
	"minitrello/internal/handler/middleware"
	"minitrello/pkg/models"
	"minitrello/pkg/models/projects"
	"minitrello/pkg/models/tasks"

	"github.com/gin-gonic/gin"
)

type IService interface {
	//User
	RegisterUser(ctx context.Context, user *models.UserInput) error
	GetByID(ctx context.Context, id int) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	LoginUser(ctx context.Context, email, password string) (string, error)
	UpdateUserName(ctx context.Context, userID int, newName string) error
	DeleteUser(ctx context.Context, id int) error
	//Project
	CreateProject(ctx context.Context, project *projects.ProjectInput, userID int) (int, error)
	GetUserProjectRole(ctx context.Context, projectID int, userID int) (string, error)
	AddMember(ctx context.Context, projectID int, userID int, role string) error
	DeleteMember(ctx context.Context, projectID, userID int) error
	GetProjectByID(ctx context.Context, id int) (*projects.Project, error)
	DeleteProject(ctx context.Context, ProjectID int) error
	GetUserProjects(ctx context.Context, userID int) ([]*projects.Project, error)
	//Tasks
	CreateTask(ctx context.Context, input *tasks.TaskInput) (int, error)
	GetTaskByID(ctx context.Context, taskID int) (*tasks.Task, error)
	UpdateTask(ctx context.Context, task *tasks.Task) error
	DeleteTask(ctx context.Context, taskID int) error
	ListTasksByUser(ctx context.Context, userID int) ([]*tasks.Task, error)
}
type Handler struct {
	service IService
}

func NewHandler(srv IService) *Handler {
	return &Handler{
		service: srv,
	}
}
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()
	//Пользователи
	auth := router.Group("/auth")
	{
		auth.POST("/register", h.RegisterUser)
		auth.POST("/login", h.Login)
	}

	user := router.Group("/user", middleware.JWTauth())
	{
		user.GET("/email/:email", middleware.RequireRole("admin"), h.GetByEmail)
		user.GET("/id/:id", middleware.RequireRole("admin"), h.GetByID)
		user.PATCH("/updatename", h.UpdateUserName)
		user.DELETE("/delete/:id", middleware.RequireRole("admin"), h.Delete)

	}
	//Проекты
	project := router.Group("/project", middleware.JWTauth())
	{
		project.POST("/create", h.CreateProject)
		project.GET("/id/:id", h.GetProjectByID)
		project.POST("/addmember/:project_id", middleware.RequireProjectRole(h.service, "owner"), h.AddMember)
		project.DELETE("/deletemember/:project_id/:user_id", middleware.RequireProjectRole(h.service, "owner"), h.DeleteMember)
		project.DELETE("/deleteproject/:project_id", middleware.RequireProjectRole(h.service, "owner"), h.DeleteProject)
		project.GET("/myprojects", h.GetUserProjects)
	}
	//Задачи
	tasks := router.Group("/tasks", middleware.JWTauth())
	{
		tasks.POST("/create/:project_id", middleware.RequireProjectRole(h.service, "owner"), h.CreateTask)
		tasks.GET("/id/:id", h.GetTaskByID)
		tasks.PUT("/update/:id", h.UpdateTask)
		tasks.DELETE("/delete/:id", h.DeleteTask)
		tasks.GET("/mytasks", h.ListTasksByUser)
	}

	return router
}
