package models

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ItemModelInterface interface {
	Insert(name, description, accessories, places, addtionalNotes string, id int) error
	AllItems(id int) ([]Item, error)
}

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

func (im *ItemModel) Insert(name, description, accessories, place, additionalNotes string, id int) error {
	statement := `INSERT INTO items (name, description, accessories, place, additional_notes, user_id, created_at)
    VALUES($1, $2, $3, $4, $5, $6, NOW())`

	_, err := im.DB.Exec(context.Background(), statement, name, description, accessories, place, additionalNotes, id)
	if err != nil {
		return err
	}
	return nil
}

func (im *ItemModel) AllItems(id int) ([]Item, error) {
	statement := "SELECT * FROM items WHERE user_id = $1"

	rows, err := im.DB.Query(context.Background(), statement, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Item

	for rows.Next() {
		var i Item
		err = rows.Scan(&i.ID, &i.Name, &i.Description, &i.Accessories, &i.Place, &i.AdditionalNotes, &i.UserID, &i.CreatedAt)
		if err != nil {
			return nil, err
		}
		items = append(items, i)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
