package models

import (
	"FoodBackend/pkg/api/dto"
	"FoodBackend/pkg/e"
	"FoodBackend/pkg/util"
	"github.com/jinzhu/gorm"
	"time"
)

type Comment struct {
	Model
	Content      string `json:"content"`
	CreateUserId int    `json:"create_user_id"`
	Status       string `json:"status" default="pending" `
	BookId       int
	Book         Book
}

type CommentMoreJson struct {
}

func (comment *Comment) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreatedAt", time.Now().Unix())
	return nil
}

func (comment *Comment) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("UpdatedAt", time.Now().Unix())
	return nil
}

func (comment *Comment) List(pageNum int, pageSize int, maps interface{}) (comments []Comment) {
	db.Model(&Comment{}).Related("Comment").Where(maps).Offset(pageNum).Limit(pageSize).Find(&comments) //预加载器，分别执行了2条SQL
	return
}

func (comment *Comment) Get(dto dto.GeneralGetDto) *Comment {
	//if me > 0 {
	//	db.Where("create_user_id=?", me)
	//}
	db.Where("id=?", dto.Id).First(&comment)
	db.Model(&comment).Related(&comment.Book)
	return comment
}

func (Comment) Update(dto dto.CommentEditDto) int64 {
	ups := Comment{
		Content: dto.Content,
		BookId:  dto.BookId,
		//MoreJson:                data["more_json"].(CommentMoreJson),
	}
	util.Log.Notice("bookModel:", ups)
	return db.Model(&Comment{Model: Model{Id: dto.Id}}).Update(&ups).RowsAffected
}

func (Comment) Create(data dto.CommentCreateDto) (Comment, int) {
	book := Book{}.GetByCon(map[string]interface{}{"id": data.BookId, "status": "publish"})
	if book.Id > 0 {
		comment := Comment{
			Content: data.Content,
			BookId:  data.BookId,
			Book:    book,
			//MoreJson:                data["more_json"].(CommentMoreJson),
		}
		util.Log.Notice("Model:", comment)
		result := db.Create(&comment)
		if result.Error == nil {
			//@TODO 返回的内容中，Book为空了
			//comment.Book = Book{}
			return comment, 0
		} else {
			util.Log.Error(result.Error.Error())
			return Comment{}, e.ERROR
		}
	}
	return Comment{}, e.BOOK_NOT_EXISTS
}

func (Comment) Delete(book *Comment) bool {
	//db.Delete(&Comment{}, "id = ?", id)
	if db.Delete(book).GetErrors() == nil {
		return true
	}
	return false
}
