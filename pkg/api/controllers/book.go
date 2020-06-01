package api

import (
	jwt "FoodBackend/middleware"
	"FoodBackend/models"
	"FoodBackend/pkg/api/dto"
	"FoodBackend/pkg/e"
	"github.com/gin-gonic/gin"
)

type BookController struct {
	BaseController
}

var bookModel = models.Book{}

func (self *BookController) List(c *gin.Context) {
	var listDto dto.GeneralListDto
	if self.BindAndValidate(c, &listDto) {

		//var b models.Book
		//var u models.User
		//models.GetDB().Model(&b).Related(&u,"CreateUserId")
		//
		//resp(c, map[string]interface{}{
		//	"b": b,
		//	"u":  u,
		//})
		//return

		books, total := bookModel.List(listDto)
		resp(c, map[string]interface{}{
			"result": books,
			"total":  total,
		})
	}
}

func (self *BookController) Get(c *gin.Context) {
	var gDto dto.GeneralGetDto
	if self.BindAndValidate(c, &gDto) {
		data := bookModel.Get(gDto)
		//role not found

		if gDto.Id < 1 {
			fail(c, e.ERROR_NOT_EXIST)
			return
		}
		// todo: get feature permission list
		// data permission list
		resp(c, map[string]interface{}{
			"result": map[string]interface{}{
				"detail": data,
			},
		})
	}
}

func (self *BookController) buildMoreJson() models.BookMoreJson {
	structThing := models.BookMoreJson{Love: "HAHA"}
	return structThing
}

//func (b *book) AfterCreate(scope *gorm.Scope) (err error) {
//	util.Log.Notice("AfterCreate()")
//	util.Log.Notice(b)
//	//if b.ID == 1{
//	//scope.DB().Model()
//	//}
//}

func (self *BookController) Create(c *gin.Context) {
	var bookDto dto.BookCreateDto
	bookDto.CreateUserId = jwt.UserId
	if self.BindAndValidate(c, &bookDto) {
		//@TODO 准备MoreJson
		newBook, err := bookModel.Create(bookDto)
		if err > 0 {
			fail(c, err)
			return
		}
		resp(c, map[string]interface{}{
			"result": newBook,
		})
	}
}

func (self *BookController) ChangeStatus(c *gin.Context) {
	var bookDto dto.BookChangeDto
	if self.BindAndValidate(c, &bookDto) {
		if bookModel.ChangeStatus(bookDto) < 0 {
			fail(c, e.ERROR)
			return
		}
		ok(c, e.SUCCESS)
	}
}

func (self *BookController) Update(c *gin.Context) {
	var bookDto dto.BookEditDto
	bookDto.CreateUserId = jwt.UserId
	if self.BindAndValidate(c, &bookDto) {
		if bookModel.Update(bookDto) < 0 {
			fail(c, e.ERROR)
			return
		}
		ok(c, e.SUCCESS)
	}
}

func (self *BookController) Delete(c *gin.Context) {
	var dto dto.GeneralDelDto
	if self.BindAndValidate(c, &dto) {
		if bookModel.Delete(&models.Book{Model: models.Model{Id: dto.Id}}) {
			fail(c, e.ERROR_NOT_EXIST)
			return
		}
		ok(c, e.SUCCESS)
	}
}
