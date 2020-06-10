package api

import (
	"FoodBackend/pkg/api/dto"
	"FoodBackend/pkg/e"
	"FoodBackend/pkg/setting"
	"FoodBackend/pkg/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

type BaseController struct {
}

func (bc *BaseController) BindAndValidate(c *gin.Context, obj interface{}) bool {
	if err := dto.Bind(c, obj); err != nil {
		message(c, e.ERROR, err.Error())
		return false
	}
	return true
}

// 返回200，带数据 —— 用于操作正常，带数据
func resp(c *gin.Context, data map[string]interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": data,
	})
}

//返回200，带状态 —— 用于操作成功，带提示
func ok(c *gin.Context, code int) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  e.GetMsg(code),
	})
}

//返回自定义code，带信息 —— 用于错误提示
func message(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  message,
	})
}

//返回错误，带信息 —— 用于操作失败，带提示
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

func getUploadPath() string {
	sec, err := setting.Cfg.GetSection("upload")
	if err != nil {
		util.Log.Fatal("Fail to get section 'upload': %v", err)
	}
	return sec.Key("UPLOAD_SAVE_DIR").String() + "/"
}

func getUrl(uri string) string {
	return setting.AppUrl + "/" + uri
}
