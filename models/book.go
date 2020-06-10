package models

import (
	"FoodBackend/pkg/api/dto"
	"FoodBackend/pkg/e"
	"FoodBackend/pkg/util"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

type Book struct {
	Model
	//TypeId     int    `json:"type_id" gorm:"index"` //声明索引，如果你使用了自动迁移功能则会有所影响，在不使用则无影响
	Name                    string       `json:"name" gorm:"type:varchar(50);unique_index"`
	Content                 string       `json:"content"`
	AllowComments           int          `json:"allow_comments"`
	CreateUserId            int          `json:"create_user_id"`
	IsShareWeChatFriendZone int          `json:"is_share_wechat_friend_zone" gorm:"column:is_share_wechat_friend_zone"`
	MoreJson                BookMoreJson `json:"more_json";gorm:"type:json"`
	Status                  string       `json:"status"`
	Tag                     []Tag        `json:"tags";gorm:"many2many:book_tags"`
	User                    User         `gorm:"-";json:"user";gorm:"foreignkey=CreateUserId;association_foreignkey=CreateUserId"`
	FileUrlJson             interface{}  `json:"file_url_json";gorm:"type:json"`
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

func (model *Book) List(listDto dto.GeneralListDto) (books []Book, total int64) {

	var tag Tag
	allRelationTags := tag.GetAllRelatedTags()

	for sk, sv := range dto.TransformSearch(listDto.Q, dto.BookListSearchMapping) {
		if sk == "name" {
			db = db.Where(fmt.Sprintf("%s LIKE ?", sk), "%"+sv+"%")
		} else {
			db = db.Where(fmt.Sprintf("%s = ?", sk), sv)
		}
	}
	//db.Model(model).Related(&model.User,"CreateUserId").Offset(listDto.Skip).Limit(listDto.Limit).Find(&books)	//NOTE not working
	//@TODO 如果未登录，则强制只能看发布的books
	db.Model(model).Where("status = ?", "publish").Offset(listDto.Skip).Limit(listDto.Limit).Order("created_at DESC", true).Find(&books)

	for bookIndex, book := range books {
		if book.CreateUserId > 0 {
			db.Model(&books[bookIndex]).Related(&books[bookIndex].User, "CreateUserId")
		}
		for _, tag := range allRelationTags {
			if tag.BookId == book.Id {
				books[bookIndex].Tag = append(books[bookIndex].Tag, tag)
			}
		}
	}
	db.Model(&books).Count(&total)
	return books, total
}

func (Book) Get(dto dto.GeneralGetDto) (book Book) {
	//if me > 0 {
	//	db.Where("create_user_id=?", me)
	//}
	//db.Preload("Tag").Where("id=?", dto.Id).Find(&book)	//NOTE 因官方这个many2many + preload有返回行数的Bug，所以不用
	db.Where("id=?", dto.Id).Find(&book)
	db.Model(&book).Related(&book.User, "CreateUserId")
	var tags Tag
	tags.GetTagsByBookId(dto.Id)
	book.Tag = []Tag{tags}
	return
}

func (Book) GetBooksByTagId(id int) (books []Book) {
	db.Where("tag_id", id).Find(&books)
	return
}

func (Book) ChangeStatus(dto dto.BookChangeDto) int64 {
	return db.Model(&Book{Model: Model{Id: dto.Id}}).Update(map[string]interface{}{"status": dto.Status}).RowsAffected
}

func (Book) Update(dto dto.BookEditDto) int64 {
	fileJson, _ := json.Marshal(dto.Files)
	ups := Book{
		//Name:                    dto.Name,
		Content:                 dto.Content,
		AllowComments:           dto.AllowComments,
		IsShareWeChatFriendZone: dto.IsShareWeChatFriendZone,
		CreateUserId:            dto.CreateUserId,
		FileUrlJson:             fileJson,
		//MoreJson:                data["more_json"].(BookMoreJson),
		Status: dto.Status,
	}
	util.Log.Notice("bookModel:", ups)
	affected := db.Model(&Book{Model: Model{Id: dto.Id}}).Update(&ups).RowsAffected

	if affected > 0 {
		//@TODO 创建Tag关联，还要去除多余的关联
		tag := Tag{
			BookId: ups.Id,
			//Name:   dto.Name,
		}
		util.Log.Notice("tagModel:", tag)
		//db.Create(&tag)
	}
	return 0
}

func (Book) Create(dto dto.BookCreateDto) (Book, int) {
	var existOne Book
	//db.Where("name = ? ", dto.Name).First(&existOne)
	fileJson, _ := json.Marshal(dto.Files)
	if existOne.Id == 0 {
		book := Book{
			//Name:                    dto.Name,
			Content:                 dto.Content,
			AllowComments:           dto.AllowComments,
			IsShareWeChatFriendZone: dto.IsShareWeChatFriendZone,
			CreateUserId:            dto.CreateUserId,
			FileUrlJson:             fileJson,
			Status:                  dto.Status,
		}
		util.Log.Notice("bookModel:", book)
		result := db.Create(&book)
		if result.Error == nil {
			//@TODO 创建Tag关联
			tag := Tag{
				BookId: book.Id,
				//Name:   dto.Name,
			}
			util.Log.Notice("tagModel:", tag)
			db.Create(&tag)
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
