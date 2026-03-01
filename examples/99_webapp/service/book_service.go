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
