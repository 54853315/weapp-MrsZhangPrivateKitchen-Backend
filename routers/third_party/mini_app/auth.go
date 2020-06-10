package mini_app

import (
	"FoodBackend/models"
	api "FoodBackend/pkg/api/controllers"
	"FoodBackend/pkg/e"
	"FoodBackend/pkg/util"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/medivhzhan/weapp/v2"
	"log"
	"net/http"
)

var userModel = models.User{}

type auth struct {
	JsCode string `json:"js_code" binding:"required"`
	Thumb  string `json:"thumb" binding:"required"`
	Name   string `json:"name" binding:"required"`
}

func GetAuth(c *gin.Context) {
	var request auth
	code := e.ERROR
	data := make(map[string]interface{})
	var controller = api.BaseController{}
	if controller.BindAndValidate(c, &request) {
		appCode = request.JsCode
		code = e.SUCCESS
		res, err := weapp.Login(appId, appSecret, appCode)
		log.Printf("weapp.Login() returns: %#v", res)
		if err != nil {
			// 处理一般错误信息
			return
		}

		if err := res.GetResponseError(); err != nil {
			// 处理微信返回错误信息
			return
		}

		openId := res.OpenID
		sessionKey := res.SessionKey

		user := userModel.GetUserByCondition(map[string]interface{}{"wx_open_id": openId})

		moreJsonMap := map[string]string{"session_key": sessionKey}

		moreJson, _ := json.Marshal(moreJsonMap)

		token, err := util.GenerateToken(openId)

		if err != nil {
			code = e.ERROR_AUTH_TOKEN
		} else {
			data["token"] = token
		}

		if user.Id == 0 {
			insert := make(map[string]interface{})
			insert["open_id"] = openId
			insert["name"] = request.Name
			insert["thumb"] = request.Thumb
			insert["api_token"] = token
			insert["more_json"] = moreJson
			user = userModel.Create(insert)
		} else {
			// 更新用户信息
			update := make(map[string]interface{})
			update["name"] = request.Name
			update["thumb"] = request.Thumb
			update["more_json"] = moreJson
			update["api_token"] = token
			userModel.Update(user.Id, update)
		}

		data["uid"] = user.Id

	} else {
		util.Log.Debug("Jwt auth fail : ", err.Error())
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})

}
