package models

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ItemModelInterface interface{}

type Item struct {
	ID              int
	Name            string
	Description     string
	Accessories     string
	Place           string
	AdditionalNotes string
	UserID          int
	CreatedAt       time.Time
}

type ItemModel struct {
	DB *pgxpool.Pool
}
