package models

import (
	"database/sql/driver"
	"encoding/json"
)

type User struct {
	Model
	ApiToken string
	MoreJson interface{} `sql:"TYPE:json"`
	Thumb    string      `json:"thumb"`
	IsAdmin  bool        `json:"is_admin"`
	IsEnable bool        `json:"is_enable"`
	WxOpenId string      `json:"wx_open_id"`
	Name     string      `json:"name"`
	Password string      `json:"password"`
	Book     []Book      `gorm:"foreignkey=CreateUserId";json:"books"`
}

type UserMoreJson struct {
	Love string
}

//func CheckUser(username, password string) bool {
//	var user User
//	db.Select("id").Where(User{Username: username, Password: password}).First(&user)
//	if user.ID > 0 {
//		return true
//	}
//	return false
//}

func (model User) Value() (driver.Value, error) {
	b, err := json.Marshal(model)
	return string(b), err
}

func (model *User) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), model)
}

func (User) CheckWeChatAuth(openId string) bool {
	var user User
	db.Select("id").Where(User{WxOpenId: openId}).First(&user)
	if user.Id > 0 {
		return true
	}
	return false
}

func (User) GetUserByCondition(maps interface{}) User {
	var user User
	db.Where(maps).First(&user)
	return user
}

func (User) GetUserById(id int) (user User) {
	//db.Where("id=?", id).First(&user)
	db.First(&user, id)
	return
}

func (User) UpdateUser(id int, data interface{}) bool {
	db.Model(&User{}).Where("id = ?", id).Updates(data)
	return true
}

func (User) CreateUser(data map[string]interface{}) bool {
	db.Create(&User{
		Name:     data["name"].(string),
		WxOpenId: data["open_id"].(string),
		Thumb:    data["thumb"].(string),
		IsAdmin:  false,
		MoreJson: data["more_json"],
		//CreatedAt: data["created_at"].(string),
		IsEnable: false,
	})
	return true
}

func GetUserIdByToken(token string) int {
	var user User
	db.Debug().Select("id").First(&user, "api_token = ?", token)
	return user.Id
}
