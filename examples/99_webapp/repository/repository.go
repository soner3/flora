// Copyright © 2026 Soner Astan astansoner@gmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	flora.Component `flora:"inject(readDB=replicaDB, writeDB=masterDB)"`
	master          *sql.DB
	replica         *sql.DB
}

// NewBookRepository receives the Master and Replica databases seamlessly thanks to Flora.
func NewBookRepository(writeDB *sql.DB, readDB *sql.DB) *BookRepository {
	fmt.Println("-> [Repo] BookRepository initialized with Master & Replica DB")
	return &BookRepository{
		master:  writeDB,
		replica: readDB,
	}
}

// GetAll intentionally uses the REPLICA for reading!
func (r *BookRepository) GetAll() ([]Book, error) {
	rows, err := r.replica.Query("SELECT id, title, author FROM books")
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
