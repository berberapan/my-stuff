package models

import "time"

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
