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
	"github.com/soner3/flora/examples/99_webapp/repository"
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
// 2. Clean Architecture: Type Aliasing
// ---------------------------------------------------------
// In Hexagonal or Clean Architecture, your core domain (like this service)
// should not have hard dependencies on external frameworks ("Framework Pollution").
// To keep your domain structs clean, Flora fully supports Go Type Aliases.
// You can define this alias once in a shared domain package and use it everywhere.
type DIComponent = flora.Component

// ---------------------------------------------------------
// 3. The Service Struct
// ---------------------------------------------------------
type BookService struct {
	// We embed the neutral alias instead of the framework type directly!
	// Flora resolves this back to flora.Component during the AST scan.
	DIComponent `flora:"constructor=ProvideBookService"`

	// We now store the interface, not the pointer to the concrete struct!
	repo BookRepository
}

// NewBookService constructor.
// Flora scans the project, finds the concrete 'repository.BookRepository',
// sees that it implements the 'BookRepository' interface, and
// injects it here automatically.
func ProvideBookService(repo BookRepository) *BookService {
	fmt.Println("-> [Service] BookService initialized with Repository Interface")
	return &BookService{
		repo: repo,
	}
}

// GetAllBooks calls the interface method.
func (s *BookService) GetAllBooks() ([]repository.Book, error) {
	return s.repo.GetAll()
}
