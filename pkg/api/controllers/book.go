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
	"path/filepath"
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
			"skip":   listDto.Skip,
		})
	}
}

func (self *BookController) Get(c *gin.Context) {
	var gDto dto.GeneralGetDto
	if self.BindAndValidate(c, &gDto) {
		data := bookModel.Get(gDto)
		//book not found
		if data.Id < 1 {
			fail(c, e.ERROR_NOT_EXIST)
			return
		}

		resp(c, map[string]interface{}{
			"result": map[string]interface{}{
				"detail": &data, //如果此处不使用&data，则MarshalJSON()并不会执行
			},
		})
	}
}

//func (self *BookController) buildMoreJson() models.BookMoreJson {
//	structThing := models.BookMoreJson{Love: "HAHA"}
//	return structThing
//}

func (self *BookController) Create(c *gin.Context) {
	var bookDto dto.BookCreateDto
	bookDto.CreateUserId = jwt.UserId
	if self.BindAndValidate(c, &bookDto) {
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
		if bookModel.Update(bookDto) <= 0 {
			fail(c, e.ERROR)
			return
		}
		//查找标签
		tagModel.CreateTagsByBookStore(bookDto.Id, bookDto.Content)
		resp(c, map[string]interface{}{
			"result": bookDto,
		})
		ok(c, e.SUCCESS)
	}
}

func (self *BookController) Delete(c *gin.Context) {
	var dto dto.GeneralDelDto
	if self.BindAndValidate(c, &dto) {
		if bookModel.Delete(&models.Book{CreateUserId: dto.CreateUserId, Model: models.Model{Id: dto.Id}}) {
			fail(c, e.ERROR_NOT_EXIST)
			return
		}
		ok(c, e.SUCCESS)
	}
}

func (self BookController) ClearPictureBeforeDisplayCreate(c *gin.Context) {
	var allFiles []string
	var usedFiles []string
	uploadPath := getUploadPath()
	_ = filepath.Walk(uploadPath, func(path string, info os.FileInfo, err error) error {
		//获取当前目录下的所有文件或目录信息
		if !info.IsDir() {
			allFiles = append(allFiles, path)
		}
		return nil
	})

	books := bookModel.All(map[string]interface{}{
		"create_user_id": string(jwt.UserId),
	})
	for _, book := range books {
		for _, bookFiles := range book.FileUrlJson {
			usedFiles = append(usedFiles, bookFiles)
		}
	}

	//util.Log.Noticef("找到的总文件数：%d，已使用的总文件数：%d。", len(allFiles), len(usedFiles))

	breakFor := false
	for _, file := range allFiles {
		for _, useFile := range usedFiles {
			if useFile == file {
				breakFor = true
			}
		}

		if breakFor {
			breakFor = false
			break
		}

		if err := os.Remove(file); err != nil {
			util.Log.Errorf("删除文件失败：", err)
		}

	}

	CleanUploadEmptySubDir()

	resp(c, map[string]interface{}{
		"result": books,
	})
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
