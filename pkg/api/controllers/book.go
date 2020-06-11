package api

import (
	jwt "FoodBackend/middleware"
	"FoodBackend/models"
	"FoodBackend/pkg/api/dto"
	"FoodBackend/pkg/e"
	"FoodBackend/pkg/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"strings"
	"time"
)

type BookController struct {
	BaseController
}

var bookModel = models.Book{}
var imageModel = models.Image{}

type timeLine struct {
	Date  string        `json:"date"`
	Books []models.Book `json:"books"`
}

func (self *BookController) timeLine(books []models.Book) []timeLine { //NOTE 切片的timeline结构体来模拟map，避免GO1.12后输出时会自动排序的问题

	var timelines []timeLine

	for i := 0; i < len(books); i++ {
		day := books[i].CreatedAt.Day()
		month := int(books[i].CreatedAt.Month())
		dateString := fmt.Sprintf("%0d-%d", month, day)
		existsDateInStrut := false

		for timelineKey, timelineItem := range timelines {
			if timelineItem.Date == dateString {
				timelines[timelineKey].Books = append(timelines[timelineKey].Books, books[i])
				existsDateInStrut = true
			}
		}

		if !existsDateInStrut {
			timelines = append(timelines, timeLine{
				Date: dateString,
				Books: []models.Book{
					books[i],
				},
			})
		}
	}
	return timelines
}

func (self *BookController) List(c *gin.Context) {
	var listDto dto.GeneralListDto
	if self.BindAndValidate(c, &listDto) {
		result, total := bookModel.List(listDto)
		books := self.timeLine(result)
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
		//book not found
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

func (self *BookController) Create(c *gin.Context) {
	var bookDto dto.BookCreateDto
	bookDto.CreateUserId = jwt.UserId
	if self.BindAndValidate(c, &bookDto) {
		util.Log.Debug(len(bookDto.Files))
	}
	newFile := make([]string, len(bookDto.Files))
	for _, file := range bookDto.Files {
		newFile = append(newFile, util.RemoveDomain(file))
	}
	bookDto.Files = newFile
	newBook, err := bookModel.Create(bookDto)
	if err > 0 {
		fail(c, err)
		return
	}
	//查找标签
	tagModel.CreateTagsByBookStore(newBook.Id, bookDto.Content)
	resp(c, map[string]interface{}{
		"result": newBook,
	})
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
			//@TODO 删除关联的图片文件
			fail(c, e.ERROR_NOT_EXIST)
			return
		}
		ok(c, e.SUCCESS)
	}
}

func (self BookController) Upload(c *gin.Context) {
	var uploadDto dto.UploadDto
	uploadDto.CreateUserId = jwt.UserId

	if self.BindAndValidate(c, &uploadDto) {
		fileExt := uploadDto.File.Filename[strings.LastIndex(uploadDto.File.Filename, "."):]
		FileName := util.GetUniqueId() + fileExt
		pathName := time.Now().Format("2006/01/02")
		savePath := getUploadPath() + pathName

		if err := os.MkdirAll(savePath, os.ModePerm); err != nil {
			util.Log.Error(err)
		}
		savePath += "/" + FileName

		_ = c.SaveUploadedFile(uploadDto.File, savePath)

		uploadDto.Url = savePath
		_, err := imageModel.Create(uploadDto)
		if err > 0 {
			fail(c, err)
			return
		}

		resp(c, map[string]interface{}{
			"result": map[string]interface{}{
				"savePath":  util.GetUrl(savePath),
				"imageName": FileName,
			},
		})
	}

}
