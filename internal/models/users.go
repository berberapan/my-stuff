package models

import "time"

type User struct {
	ID             int
	Email          string
	HashedPassword []byte
	CreatedAt      time.Time
}
