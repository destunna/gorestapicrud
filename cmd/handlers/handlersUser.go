package handlers

import (
	"gorestapicrud/cmd/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	Repo models.UserRepository
}

func NewUserHandler(repo models.UserRepository) *UserHandler {
	return &UserHandler{Repo: repo}
}

func (h *UserHandler) HandleCreateUser(c echo.Context) error {
	user := models.User{}

	// c.Bind() - функция для привязки тела запроса к user переменной
	c.Bind(&user)

	newUser, err := h.Repo.CreateUser(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, newUser)
}

func (h *UserHandler) HandleUpdateUser(c echo.Context) error {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	user := models.User{}
	c.Bind(&user)

	updateUser, err := h.Repo.UpdateUser(user, idInt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, updateUser)
}

func (h *UserHandler) HandleGetUser(c echo.Context) error {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}

	user := models.User{}

	currentUser, err := h.Repo.GetUser(user, idInt)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	return c.JSON(http.StatusOK, currentUser)
}

func (h *UserHandler) HandleGetAllUsers(c echo.Context) error {
	page := c.QueryParam("page")
	pageInt, err := strconv.Atoi(page)

	limit := 2

	users, err := h.Repo.GetAllUsers(pageInt, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Если пользователей в базе нет, вернется пустой массив `[]`, а не `null`
	if users == nil {
		users = []models.User{}
	}

	return c.JSON(http.StatusOK, users)
}

func (h *UserHandler) HandleDeleteUser(c echo.Context) error {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}

	user := models.User{}

	deleteErr := h.Repo.DeleteUser(user, idInt)
	if deleteErr != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	return c.JSON(http.StatusOK, map[string]string{"success": "User deleted"})
}
