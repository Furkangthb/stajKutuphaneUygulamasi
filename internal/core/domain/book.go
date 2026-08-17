package domain

import "time"

type Book struct {
	ID          int       `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Author      string    `db:"author" json:"author"`
	Genre       string    `db:"genre" json:"genre"`
	PublishDate time.Time `db:"publish_date" json:"publish_date"`
	Description string    `db:"description" json:"description"`
	StockCount  int       `db:"stock_count" json:"stock_count"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}
