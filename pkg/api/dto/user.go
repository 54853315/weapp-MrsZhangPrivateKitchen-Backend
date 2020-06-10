package dto

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

// UserListSearchMapping - define search query keys in user list page
var UserListSearchMapping = map[string]string{
	"n": "name",
}

//UserCreateDto - binding user creation params
type UserCreateDto struct {
	Id       int    `json:"id"`
	ApiToken string `json:"api_token" binding:"required"`
	WxOpenId string `json:"wx_open_id" binding:"required"`
	Name     string `form:"name" json:"name" binding:"required"`
	Sex      int    `form:"sex" json:"sex"`
	//Password      string    `form:"password" json:"password" binding:"required,pwdValidate"`
	Thumb  string `json:"thumb"`
	Status int    `form:"status,default=1" json:"status"`
}

//UserCreateDto - binding user edition params
type UserEditDto struct {
	Id       int    `uri:"id" json:"id" binding:"required"`
	ApiToken string `json:"api_token" binding:"required"`
	WxOpenId string `json:"wx_open_id" binding:"required"`
	Name     string `form:"name" json:"name" binding:"required"`
	Sex      int    `form:"sex" json:"sex"`
	Thumb    string `json:"thumb"`
	Status   int    `form:"status,default=1" json:"status"`
}

type UserMoreJson struct {
}

// password validator
func pwdValidate(fl validator.FieldLevel) bool {
	reg := regexp.MustCompile(`^[a-zA-Z0-9!@#$%^&*]{6,}$`)
	if val, ok := fl.Field().Interface().(string); ok {
		if !reg.Match([]byte(val)) {
			return false
		}
	}
	return true
}
