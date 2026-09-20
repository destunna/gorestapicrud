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
