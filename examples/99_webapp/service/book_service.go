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

package service

import (
	"fmt"

	"github.com/soner3/flora"
	"github.com/soner3/flora/examples/99_full_tutorial/repository"
)

// ---------------------------------------------------------
// 1. The Interface Definition
// ---------------------------------------------------------
// By defining the interface here, the service doesn't care HOW
// the books are stored (SQL, NoSQL, Memory). It just cares
// that it can call "GetAll()".
type BookRepository interface {
	GetAll() ([]repository.Book, error)
}

// ---------------------------------------------------------
// 2. The Service Struct
// ---------------------------------------------------------
type BookService struct {
	flora.Component

	// We now store the interface, not the pointer to the concrete struct!
	repo BookRepository
}

// NewBookService constructor.
// Flora scans the project, finds 'repository.BookRepository' (the struct),
// sees that it implements the 'BookRepository' interface, and
// injects it here automatically.
func NewBookService(repo BookRepository) *BookService {
	fmt.Println("-> [Service] BookService initialized with Repository Interface")
	return &BookService{
		repo: repo,
	}
}

// GetAllBooks calls the interface method.
func (s *BookService) GetAllBooks() ([]repository.Book, error) {
	return s.repo.GetAll()
}
