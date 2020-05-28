package models

import (
	"FoodBackend/pkg/api/dto"
	"FoodBackend/pkg/e"
	"FoodBackend/pkg/util"
	"github.com/jinzhu/gorm"
	"time"
)

type Book struct {
	Model
	//TypeId     int    `json:"type_id" gorm:"index"` //声明索引，如果你使用了自动迁移功能则会有所影响，在不使用则无影响
	Name                    string       `gorm:"type:varchar(50);unique_index";json:"name"`
	Content                 string       `json:"content"`
	AllowComments           int          `json:"allow_comments"`
	CreateUserId            int          `json:"create_user_id"`
	IsShareWeChatFriendZone int          `json:"is_share_wechat_friend_zone" gorm:"column:is_share_wechat_friend_zone"`
	MoreJson                BookMoreJson `json:"more_json";gorm:"type:json"`
	Status                  string       `json:"status"`
	Tag                     []Tag        `gorm:"many2many:book_tags" json:"tag"`
	//多对多关系
}

type BookMoreJson struct {
	Love string
}

func (book *Book) checkUnique(name string) bool {
	db.Select("id").Where(&Book{Name: name, CreateUserId: 1}).First(&book)
	if book.Id > 0 {
		return true
	}
	return false
}

func (book *Book) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreatedAt", time.Now().Unix())
	return nil
}

func (book *Book) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("UpdatedAt", time.Now().Unix())
	return nil
}

func (Book) GetByCon(maps interface{}) Book {
	var book Book
	db.Model(&book).Where(maps).First(&book)
	return book
}

func (Book) List(pageNum int, pageSize int, maps interface{}) (books []Book) {
	db.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&books) //预加载器，分别执行了2条SQL
	//db.Debug().Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&books) //预加载器，分别执行了2条SQL
	return
}

func (Book) Get(dto dto.GeneralGetDto) Book {
	var book Book
	//if me > 0 {
	//	db.Where("create_user_id=?", me)
	//}
	db.Where("id=?", dto.Id).First(&book)
	db.Model(&book).Related(&book.Tag)
	return book
}

func (Book) GetBooksByTagId(id int) []Book {
	var books []Book
	db.Where("tag_id", id).Find(&books)
	return books
}

func (Book) ChangeStatus(dto dto.BookChangeDto) int64 {
	return db.Model(&Book{Model: Model{Id: dto.Id}}).Update(map[string]interface{}{"status": dto.Status}).RowsAffected
}

func (Book) Update(dto dto.BookEditDto) int64 {
	ups := Book{
		Name:                    dto.Name,
		Content:                 dto.Content,
		AllowComments:           dto.AllowComments,
		IsShareWeChatFriendZone: dto.IsShareWeChatFriendZone,
		CreateUserId:            dto.CreateUserId,
		//MoreJson:                data["more_json"].(BookMoreJson),
		Status: dto.Status,
	}
	util.Log.Notice("bookModel:", ups)
	return db.Model(&Book{Model: Model{Id: dto.Id}}).Update(&ups).RowsAffected
}

func (Book) Create(dto dto.BookCreateDto) (Book, int) {
	var existOne Book
	db.Where("name = ? ", dto.Name).First(&existOne)
	if existOne.Id == 0 {
		book := Book{
			Name:                    dto.Name,
			Content:                 dto.Content,
			AllowComments:           dto.AllowComments,
			IsShareWeChatFriendZone: dto.IsShareWeChatFriendZone,
			CreateUserId:            dto.CreateUserId,
			//MoreJson:                data["more_json"].(BookMoreJson),
			Status: dto.Status,
		}
		util.Log.Notice("bookModel:", book)
		result := db.Create(&book)
		if result.Error == nil {
			return book, 0
		} else {
			util.Log.Error(result.Error.Error())
			return Book{}, e.ERROR
		}
	}
	return Book{}, e.BOOK_EXISTS
}

func (Book) Delete(book *Book) bool {
	//db.Delete(&Book{}, "id = ?", id)
	if db.Delete(book).GetErrors() == nil {
		return true
	}
	return false
}
