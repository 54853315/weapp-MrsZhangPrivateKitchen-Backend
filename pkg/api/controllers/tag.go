package api

import (
	"FoodBackend/models"
	"FoodBackend/pkg/api/dto"
	"fmt"
	"github.com/gin-gonic/gin"
)

type TagController struct {
	BaseController
}

var tagModel = models.Tag{}

func (self *TagController) List(c *gin.Context) {

	var listDto dto.GeneralListDto
	if self.BindAndValidate(c, &listDto) {
		var tags []models.Tag
		var total int64
		db := models.GetDB()
		for sk, sv := range dto.TransformSearch(listDto.Q, dto.BookListSearchMapping) {
			if sk == "name" {
				db = db.Where(fmt.Sprintf("%s LIKE ?", sk), "%"+sv+"%")
			} else {
				db = db.Where(fmt.Sprintf("%s = ?", sk), sv)
			}
		}
		db.Offset(listDto.Skip).Find(&tags)
		db.Model(&models.Book{}).Count(&total)
		resp(c, map[string]interface{}{
			"result": tags,
			"total":  total,
		})
	}
}
