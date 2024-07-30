package models

import "fmt"

type Book struct {
	ID     int    `json:"id"`
	Title string `json:"title"`
	Author string `json:"author"`
	Status string `json:"status"`
}

// string representation of Book, should output json format
func (b Book) String() string {
	return fmt.Sprintf("Book{ID: %d, Title: %s, Author: %s, Status: %s}", b.ID, b.Title, b.Author, b.Status)
}