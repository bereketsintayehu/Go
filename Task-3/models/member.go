package models

type Member struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	BorrowedBooks []int `json:"borrowed_books"`
}
