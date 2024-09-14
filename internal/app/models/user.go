package models

import "time"

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}
