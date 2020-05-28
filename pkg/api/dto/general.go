package dto

import (
	"github.com/go-playground/validator"
	"strings"
)

// GeneralListDto - General list request params
type GeneralListDto struct {
	Skip  int    `form:"skip,default=0" json:"skip"`
	Limit int    `form:"limit,default=20" json:"limit" binding:"max=10000"`
	Order string `form:"order" json:"order"`
	Q     string `form:"q" json:"q"`
}

//type GeneralTreeDto struct {
//	Q string `form:"q" json:"q"`
//}

type GeneralAuthDto struct {
	CreateUserId int
	//Token string `form:"token" binding:"required"`
}

type GeneralDelDto struct {
	GeneralAuthDto
	Id int `uri:"id" json:"id" binding:"required"`
}
type GeneralGetDto struct {
	GeneralAuthDto
	Id     int    `uri:"id" json:"id" binding:"required"`
	Status string `uri:"status" form:"status" json:"status" `
}

// TransformSearch - transform search query
func TransformSearch(qs string, mapping map[string]string) (ss map[string]string) {
	ss = make(map[string]string)
	for _, v := range strings.Split(qs, ",") {
		vs := strings.Split(v, "=")
		if _, ok := mapping[vs[0]]; ok {
			ss[mapping[vs[0]]] = vs[1]
		}
	}
	return
}

func OwnedValidate(fl validator.FieldLevel) bool {
	//检测当前模型的字段内是否有create_user_id，并且是否为我自己
	//用dto的CreateUserId或者去User模型里读取userid
	return true
}
