package main

import (
	"github.com/labstack/echo/v4"

	"gorestapicrud/cmd/handlers"
	"gorestapicrud/cmd/models"
	"gorestapicrud/cmd/repositories"
	"gorestapicrud/cmd/storage"
)

type MainUserRepository struct{}

func (m *MainUserRepository) CreateUser(user models.User) (models.User, error) {
	return repositories.CreateUser(user)
}

func (m *MainUserRepository) UpdateUser(user models.User, id int) (models.User, error) {
	return repositories.UpdateUser(user, id)
}

func (m *MainUserRepository) GetUser(user models.User, id int) (models.User, error) {
	return repositories.GetUser(user, id)
}

func (m *MainUserRepository) GetAllUsers(page int, limit int) ([]models.User, error) {
	return repositories.GetAllUsers(page, limit)
}

func (m *MainUserRepository) DeleteUser(user models.User, id int) error {
	return repositories.DeleteUser(user, id)
}

func main() {
	e := echo.New()
	e.GET("/", handlers.Home)

	storage.InitDB()

	dbRepo := &MainUserRepository{}

	userHandler := handlers.NewUserHandler(dbRepo)

	e.POST("/users", userHandler.HandleCreateUser)
	e.PUT("/users/:id", userHandler.HandleUpdateUser)
	e.GET("/users/:id", userHandler.HandleGetUser)
	e.GET("/users", userHandler.HandleGetAllUsers)
	e.DELETE("/users/:id", userHandler.HandleDeleteUser)

	e.Logger.Fatal(e.Start(":8080"))
}
