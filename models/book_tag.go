package models

type BookTag struct {
	Model
	BookId int `json:"book_id"`
	TagId  int `json:"tag_id"`
}

var BookTagTableName = "book_tags"
