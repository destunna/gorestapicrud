package main

import (
	"log"

	"github.com/labstack/echo/v4"

	"gorestapicrud/internal/handlers"
	"gorestapicrud/internal/repositories"
	"gorestapicrud/internal/storage"
)

func main() {
	e := echo.New()
	e.GET("/", handlers.Home)

	db, err := storage.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := repositories.NewUserRepository(db)

	userHandler := handlers.NewUserHandler(repo)

	e.POST("/users", userHandler.HandleCreateUser)
	e.PUT("/users/:id", userHandler.HandleUpdateUser)
	e.GET("/users/:id", userHandler.HandleGetUser)
	e.GET("/users", userHandler.HandleGetAllUsers)
	e.DELETE("/users/:id", userHandler.HandleDeleteUser)

	e.Logger.Fatal(e.Start(":8080"))
}
