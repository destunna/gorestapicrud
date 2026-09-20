package main

import (
	"github.com/labstack/echo/v4"

	"gorestapicrud/cmd/handlers"
	"gorestapicrud/cmd/storage"
)

func main() {
	e := echo.New()
	e.GET("/", handlers.Home)

	storage.InitDB()

	e.POST("/users", handlers.HandleCreateUser)
	e.PUT("/users/:id", handlers.HandleUpdateUser)
	e.GET("/users/:id", handlers.HandleGetUser)
	e.GET("/users", handlers.HandleGetAllUsers)
	e.DELETE("/users/:id", handlers.HandleDeleteUser)

	e.Logger.Fatal(e.Start(":8080"))
}
