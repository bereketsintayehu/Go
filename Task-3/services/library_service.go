package services

import (
	"fmt"
	"library/models"
	"sync"
)

type LibraryService interface  {
	AddBook (book models.Book) (models.Book, error)
	RemoveBook (bookID int) (models.Book, error)
	BorrowBook (bookID, memberID int) (models.Book, error)
	ReturnBook (bookID, memberID int) (models.Book, error)
	ListAvailableBooks () []models.Book
	ListBorrowedBooks () []models.Book
}

type libraryService struct {
	books map[int]models.Book
	members map[int]models.Member
	bookIDCounter int
	mu sync.Mutex
}

// TODO: Define error types

func NewLibraryService() LibraryService {
	return &libraryService{
		books: make(map[int]models.Book),
		members: make(map[int]models.Member),
	}
}

func (ls *libraryService) AddBook (book models.Book) (models.Book, error) {
	
	ls.mu.Lock()
	defer ls.mu.Unlock()
	book.ID = ls.bookIDCounter
	ls.bookIDCounter++
	ls.books[book.ID] = book

	return book, nil
}

func (ls *libraryService) RemoveBook (bookID int) (models.Book, error) {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	book, exists := ls.books[bookID]
	if !exists {
		fmt.Println("Book not found")
		// TODO: Return error
		return book, nil
	}
	delete(ls.books, bookID)
	return book, nil
}

func (ls *libraryService) BorrowBook (bookID, memberID int) (models.Book, error) {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	book, exists := ls.books[bookID]
	if !exists {
		fmt.Println("Book not found")
		// TODO: Return error
		return book, nil
	}
	if book.Status == "borrowed" {
		fmt.Println("Book is already borrowed")
		// TODO: Return error
		return book, nil
	}
	book.Status = "borrowed"
	ls.books[bookID] = book
	return book, nil
}

func (ls *libraryService) ReturnBook (bookID, memberID int) (models.Book, error) {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	book, exists := ls.books[bookID]
	if !exists {
		fmt.Println("Book not found")
		// TODO: Return error
		return book, nil
	}
	if book.Status == "available" {
		fmt.Println("Book is already available")
		return book, nil
	}

	book.Status = "available"
	ls.books[bookID] = book
	return book, nil
}

func (ls *libraryService) ListAvailableBooks () []models.Book {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	availableBooks := []models.Book{}
	for _, book := range ls.books {
		if book.Status == "available" {
			availableBooks = append(availableBooks, book)
		}
	}
	return availableBooks
}

func (ls *libraryService) ListBorrowedBooks () []models.Book {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	borrowedBooks := []models.Book{}
	for _, book := range ls.books {
		if book.Status == "borrowed" {
			borrowedBooks = append(borrowedBooks, book)
		}
	}
	return borrowedBooks
}
