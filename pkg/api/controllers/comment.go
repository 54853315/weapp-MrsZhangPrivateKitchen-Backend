package api

import (
	jwt "FoodBackend/middleware"
	"FoodBackend/models"
	"FoodBackend/pkg/api/dto"
	"FoodBackend/pkg/e"
	"FoodBackend/pkg/util"
	"fmt"
	"github.com/gin-gonic/gin"
)

type CommentController struct {
	BaseController
}

var commentModel = models.Comment{}

func (self *CommentController) List(c *gin.Context) {

	var listDto dto.CommentListDto
	if self.BindAndValidate(c, &listDto) {

		var comments []models.Comment
		var total int64
		db := models.GetDB()
		//@TODO 需要增加条件：必须是自己发布的
		//db.Where("create_user_id")
		for sk, sv := range dto.CommentListSearchMapping {
			db = db.Where(fmt.Sprintf("%s = ?", sk), c.Query(sv))
		}

		db.Preload("Comment").Offset(listDto.Skip).Limit(listDto.Limit).Find(&comments)
		db.Model(&models.Comment{}).Count(&total)
		resp(c, map[string]interface{}{
			"result": comments,
			"total":  total,
		})
	}
}

func (self *CommentController) Create(c *gin.Context) {
	var dto dto.CommentCreateDto
	dto.CreateUserId = jwt.UserId
	if self.BindAndValidate(c, &dto) {
		util.Log.Notice(dto)
		//@TODO 准备MoreJson
		newBook, err := commentModel.Create(dto)
		if err > 0 {
			fail(c, err)
			return
		}
		resp(c, map[string]interface{}{
			"result": newBook,
		})
	}
}

func (self *CommentController) Delete(c *gin.Context) {
	var dto dto.GeneralDelDto
	if self.BindAndValidate(c, &dto) {
		if commentModel.Delete(&models.Comment{Model: models.Model{Id: dto.Id}}) {
			fail(c, e.ERROR_NOT_EXIST)
			return
		}
		ok(c, e.SUCCESS)
	}
}
