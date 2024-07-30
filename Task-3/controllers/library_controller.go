package controllers

import (
	"fmt"
	"library/models"
	"library/services"
)

type LibraryController struct {
	libraryService services.LibraryService
}

func NewLibraryController(libraryService services.LibraryService) *LibraryController {
	return &LibraryController{
		libraryService: libraryService,
	}
}

func (lc *LibraryController) AddBook() {
	var book models.Book
	fmt.Println("Enter book title:")
	fmt.Scanln(&book.Title)
	fmt.Println("Enter book author:")
	fmt.Scanln(&book.Author)
	book.Status = "available"
	book, err := lc.libraryService.AddBook(book)
	if err != nil {
		fmt.Println("Error adding book")
		return
	}
	fmt.Println("Book added successfully")
}

// RemoveBook removes a book from the library
func (lc *LibraryController) RemoveBook() {
	var bookID int
	fmt.Println("Enter book ID:")
	fmt.Scanln(&bookID)
	book, err := lc.libraryService.RemoveBook(bookID)
	if err != nil {
		fmt.Println("Error removing book")
		return
	}
	fmt.Println("Book removed successfully", book)
}

// BorrowBook borrows a book from the library
func (lc *LibraryController) BorrowBook() {
	var bookID, memberID int
	fmt.Println("Enter book ID:")
	fmt.Scanln(&bookID)
	fmt.Println("Enter member ID:")
	fmt.Scanln(&memberID)
	book, err := lc.libraryService.BorrowBook(bookID, memberID)
	if err != nil {
		fmt.Println("Error borrowing book")
		return
	}
	fmt.Println("Book borrowed successfully", book)
}

// ReturnBook returns a book to the library
func (lc *LibraryController) ReturnBook() {
	var bookID, memberID int
	fmt.Println("Enter book ID:")
	fmt.Scanln(&bookID)
	fmt.Println("Enter member ID:")
	fmt.Scanln(&memberID)
	book, err := lc.libraryService.ReturnBook(bookID, memberID)
	if err != nil {
		fmt.Println("Error returning book")
		return
	}
	fmt.Println("Book returned successfully", book)
}

func (lc *LibraryController) ListAvailableBooks() {
	books := lc.libraryService.ListAvailableBooks()
	fmt.Println("Available Books in library:")
	for _, book := range books {
		fmt.Println(book)
	}
}

func (lc *LibraryController) ListBorrowedBooks() {
	books := lc.libraryService.ListBorrowedBooks()
	fmt.Println("Borrowed Books in library:")
	for _, book := range books {
		fmt.Println(book)
	}
}

func (lc *LibraryController) Menu() {
	for {
		fmt.Println("1. Add Book")
		fmt.Println("2. Remove Book")
		fmt.Println("3. Borrow Book")
		fmt.Println("4. Return Book")
		fmt.Println("5. List Available Books")
		fmt.Println("6. List Borrowed Books")
		fmt.Println("7. Exit")
		var choice int
		fmt.Scanln(&choice)
		switch choice {
		case 1:
			lc.AddBook()
		case 2:
			lc.RemoveBook()
		case 3:
			lc.BorrowBook()
		case 4:
			lc.ReturnBook()
		case 5:
			lc.ListAvailableBooks()
		case 6:
			lc.ListBorrowedBooks()
		case 7:
			return
		default:
			fmt.Println("Invalid choice")
		}
	}
}