package models

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Item struct {
	ID              int
	Name            string
	Description     string
	Accessories     string
	Place           string
	Manual          string
	Receipt         string
	WarrantyExp     time.Time
	InsuranceExp    time.Time
	AdditionalNotes string
	UserID          int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type ItemModel struct {
	DB *pgxpool.Pool
}
