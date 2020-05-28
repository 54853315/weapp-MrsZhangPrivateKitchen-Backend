package api

import (
	"FoodBackend/pkg/api/dto"
	"FoodBackend/pkg/e"
	"github.com/gin-gonic/gin"
	"net/http"
)

type BaseController struct {
}

func (bc *BaseController) BindAndValidate(c *gin.Context, obj interface{}) bool {
	if err := dto.Bind(c, obj); err != nil {
		rawOk(c, err.Error())
		return false
	}
	return true
}

func resp(c *gin.Context, data map[string]interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": data,
	})
}
func ok(c *gin.Context, code int) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  e.GetMsg(code),
	})
}

func rawOk(c *gin.Context, message string) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  message,
	})
}

func fail(c *gin.Context, code int) {
	//currentLang,_ := c.Cookie("lang")
	//currentLang := GetLang
	//errs.Message = i18n.Tr(middleware.GetLang(), errs.Langkey)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		//"code":     errs.Code,
		//"msg":      errs.Message,
		//"moreinfo": errs.Moreinfo,
	})
}
