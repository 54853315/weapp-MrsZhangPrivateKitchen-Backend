package models

import (
	"FoodBackend/pkg/util"
	"github.com/jinzhu/gorm"
	"regexp"
	"time"
)

type Tag struct {
	Model
	Name   string `json:"name"`
	Book   []Book `gorm:"many2many:book_tags";json:"books"`
	BookId int    `gorm:"-"`
}

func (tag *Tag) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreatedAt", time.Now().Unix())
	return nil
}

func (tag *Tag) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("UpdatedAt", time.Now().Unix())
	return nil
}

func ExistsTagById(id int) bool {
	var tag Tag
	db.Select("id").Where("id = ?", id).First(&tag)
	if tag.Id > 0 {
		return true
	}
	return false
}

// NOTE 等gorm官方升级后解决many2many + preload 的rows bug后，可以不用这个方法拿数据了
func (tag *Tag) GetTagsByBookId(bookId int) (tags []Tag) {
	db.Where("id IN (?)", db.Table("book_tags").Select("tag_id").Where("book_id = ?", bookId).SubQuery()).Find(&tags)
	return
}

func (Tag) GetAllRelatedTags() (tags []Tag) {
	db.Table(BookTagTableName).Select(BookTagTableName + ".book_id,tags.*").Joins("left join tags on tags.id = book_tags.tag_id").Find(&tags)
	return
}

//创建标签关联
func (Tag) CreateTagsByBookStore(bookId int, content string) {
	const tagReg = "#([0-9a-zA-Z]+) "
	compile := regexp.MustCompile(tagReg)
	subMatch := compile.FindAllSubmatch([]byte(content), -1)
	for _, m := range subMatch {
		util.Log.Noticef("开始处理查找到的标签%s", string(m[1]))
		tag := Tag{
			Name: string(m[1]),
		}
		db.Where(tag).FirstOrCreate(&tag)
		bookTag := BookTag{
			BookId: bookId,
			TagId:  tag.Id,
		}
		db.Where(bookTag).FirstOrCreate(&bookTag)
	}
}
