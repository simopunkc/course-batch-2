package domain

import "time"

type User struct {
	ID        int
	Name      string
	Email     string
	Password  string
	NoHP      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
