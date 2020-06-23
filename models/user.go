package models

import (
	"database/sql/driver"
	"encoding/json"
)

type User struct {
	Model
	ApiToken string      `json:"-"`
	MoreJson interface{} `json:"-" sql:"TYPE:json"`
	Thumb    string      `json:"thumb"`
	IsAdmin  bool        `json:"-"`
	IsEnable bool        `json:"-"`
	WxOpenId string      `json:"-"`
	Name     string      `json:"name"`
	Password string      `json:"-"`
	Book     []Book      `gorm:"foreignkey=CreateUserId" json:"books"`
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

func (User) Get(id int) (user User) {
	db.Where("id=?", id).First(&user)
	return
}

func (User) Update(id int, data interface{}) bool {
	db.Debug().Model(&User{}).Where("id = ?", id).Updates(data)
	return true
}

func (User) Create(data map[string]interface{}) (user User) {

	user = User{
		Name:     data["name"].(string),
		WxOpenId: data["open_id"].(string),
		ApiToken: data["api_token"].(string),
		Thumb:    data["thumb"].(string),
		IsAdmin:  false,
		MoreJson: data["more_json"],
		IsEnable: false,
	}

	db.Debug().Create(&user)
	return
}

func GetUserIdByToken(token string) int {
	var user User
	db.Debug().Select("id").First(&user, "api_token = ?", token)
	return user.Id
}
