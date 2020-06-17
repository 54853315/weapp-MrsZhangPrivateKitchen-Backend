package models

import (
	"FoodBackend/pkg/api/dto"
	"FoodBackend/pkg/e"
	"FoodBackend/pkg/util"
	"encoding/json"
	"errors"
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
	User                    User         `json:"user" gorm:"-;foreignkey=CreateUserId;association_foreignkey=CreateUserId"`
	FileUrlJson             FileJson     `json:"file_url_json" gorm:"column:file_url_json;default:'[]';type:json"`

	// NOTE 切片Struct 在 Gorm内不识别，无法进行json数组存储
	// 具体见：https://github.com/go-gorm/gorm/issues/1879#issuecomment-643954492
	//FileUrlJson             []FileJson     `json:"file_url_json";gorm:"column:file_url_json;type:json"`
}

type FileJson NormalJson

func (f *FileJson) MarshalJSON() ([]byte, error) {
	var jsonArr []string
	for _, v := range *f {
		jsonArr = append(jsonArr, util.GetUrl(string(v)))
	}
	return json.Marshal(jsonArr)
}
func (f *FileJson) Scan(input interface{}) error {
	switch value := input.(type) {
	case []byte:
		return json.Unmarshal(value, &f)
	default:
		return errors.New("not supported")
	}
}

type BookMoreJson struct {
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

func (book Book) All(where map[string]interface{}) (books []Book) {
	db.Model(book).Where(where).Find(&books)
	return books
}

func (book *Book) List(listDto dto.GeneralListDto) (books []Book, total int64) {

	var tag Tag
	allRelationTags := tag.GetAllRelatedTags()

	for sk, sv := range dto.TransformSearch(listDto.Q, dto.BookListSearchMapping) {
		if sk == "name" {
			db = db.Where(fmt.Sprintf("%s LIKE ?", sk), "%"+sv+"%")
		} else {
			db = db.Where(fmt.Sprintf("%s = ?", sk), sv)
		}
	}
	//db.Model(book).Related(&book.User,"CreateUserId").Offset(listDto.Skip).Limit(listDto.Limit).Find(&books)	//NOTE not working
	//@TODO 如果未登录，则强制只能看发布的books
	db.Model(book).Where("status = ?", "publish").Offset(listDto.Skip).Limit(listDto.Limit).Order("created_at DESC", true).Find(&books)

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
	//db.Preload("Tag").Where("id=?", dto.Id).Find(&book)	//NOTE 因官方这个many2many + preload有返回行数的Bug，所以不用
	db.Where("id=?", dto.Id).Find(&book)
	db.Model(&book).Related(&book.User, "CreateUserId")
	var tags Tag
	if tags.GetTagsByBookId(dto.Id) > 0 {
		book.Tag = []Tag{tags}
	}
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
	var fileJson []string
	for _, file := range dto.Files {
		if file != "" {
			fileJson = append(fileJson, file)
		}
	}
	ups := Book{
		Content:                 dto.Content,
		AllowComments:           dto.AllowComments,
		IsShareWeChatFriendZone: dto.IsShareWeChatFriendZone,
		CreateUserId:            dto.CreateUserId,
		FileUrlJson:             fileJson,
		Status:                  dto.Status,
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
	var fileJson []string
	for _, file := range dto.Files {
		if file != "" {
			fileJson = append(fileJson, file)
		}
	}
	book := Book{
		Content:                 dto.Content,
		AllowComments:           dto.AllowComments,
		IsShareWeChatFriendZone: dto.IsShareWeChatFriendZone,
		CreateUserId:            dto.CreateUserId,
		FileUrlJson:             fileJson,
		Status:                  dto.Status,
	}
	util.Log.Notice("bookModel:", book)
	result := db.Debug().Create(&book)
	util.Log.Debug("报错情况：", result.GetErrors())
	if result.Error == nil {
		return book, 0
	} else {
		util.Log.Error(result.Error.Error())
		return Book{}, e.ERROR
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

// 在一个事务中更新数据
func (book *Book) AfterDelete(tx *gorm.DB) (err error) {
	if len(book.FileUrlJson) > 0 {

	}
	return
}
