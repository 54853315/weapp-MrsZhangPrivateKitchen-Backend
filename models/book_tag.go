package models

type BookTag struct {
	Id     int `gorm:"primary_key" json:"id"`
	BookId int `json:"book_id"`
	TagId  int `json:"tag_id"`
}

var BookTagTableName = "book_tags"
