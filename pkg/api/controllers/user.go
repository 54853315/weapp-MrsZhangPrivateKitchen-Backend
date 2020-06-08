package api

import (
	"FoodBackend/models"
	"FoodBackend/pkg/api/dto"
	"FoodBackend/pkg/e"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	BaseController
}

var userModel = models.User{}

func (self *UserController) Get(c *gin.Context) {
	var gDto dto.GeneralGetDto
	if self.BindAndValidate(c, &gDto) {
		data := userModel.Get(gDto.Id)
		//book not found
		if data.Id < 1 {
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

func (self *UserController) Update(c *gin.Context) {

}

func (self *UserController) buildMoreJson() models.UserMoreJson {
	structThing := models.UserMoreJson{Love: "HAHA"}
	return structThing
}
