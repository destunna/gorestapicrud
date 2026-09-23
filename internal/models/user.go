package models

import "time"

type User struct {
	Id        int       `json:"id"`
	FullName  string    `json:"full_name"`
	Age       string    `json:"age"`
	Habits    string    `json:"habits"`
	Alive     bool      `json:"alive"`
	CreatedAt time.Time `json:"created_at"`
}

type UserRepository interface {
	CreateUser(user User) (User, error)
	UpdateUser(user User, id int) (User, error)
	GetUser(id int) (User, error)
	GetAllUsers(page, limit int) ([]User, error)
	DeleteUser(id int) error
}
