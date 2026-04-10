package app

import (
	"log"
	"minitrello/internal/handler"
	"minitrello/internal/repository"
	"minitrello/internal/service"
	"minitrello/pkg/db/postgresql"
)

func Run() {
	conn, err := postgresql.NewPostgresPgxCon()
	if err != nil {
		panic(err)
	}
	repo := repository.NewRepository(conn)
	service := service.NewService(repo)
	handler := handler.NewHandler(service)
	router := handler.InitRoutes()

	if err := router.Run(":8080"); err != nil {
		log.Fatal("failed to run server: ", err)
	}
}
