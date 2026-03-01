package repository

import (
	"database/sql"
	"fmt"

	"github.com/soner3/flora"
)

// Book represents our database model
type Book struct {
	ID     int
	Title  string
	Author string
}

// BookRepository handles database operations for books.
type BookRepository struct {
	flora.Component

	db *sql.DB
}

func NewBookRepository(db *sql.DB) *BookRepository {
	fmt.Println("-> [Repo] BookRepository initialized")
	return &BookRepository{
		db: db,
	}
}

// GetAll fetches all books from the database
func (r *BookRepository) GetAll() ([]Book, error) {
	rows, err := r.db.Query("SELECT id, title, author FROM books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var b Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Author); err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	return books, nil
}
