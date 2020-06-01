package models

import (
	"FoodBackend/pkg/setting"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"time"
)

var db *gorm.DB

type Model struct {
	Id        int       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func init() {
	var (
		err                                               error
		dbType, dbName, user, password, host, tablePrefix string
	)

	sec, err := setting.Cfg.GetSection("database")
	if err != nil {
		log.Fatal(2, "Fail to get section 'database': %v", err)
	}

	dbType = sec.Key("TYPE").String()
	dbName = sec.Key("NAME").String()
	user = sec.Key("USER").String()
	password = sec.Key("PASSWORD").String()
	host = sec.Key("HOST").String()
	tablePrefix = sec.Key("TABLE_PREFIX").String()

	db, err = gorm.Open(dbType, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		user,
		password,
		host,
		dbName))

	if err != nil {
		log.Println(err)
	}

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return tablePrefix + defaultTableName
	}

	//db.SingularTable(true)
	// 启用Logger，显示详细日志
	db.LogMode(true)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
}

func CloseDB() {
	defer db.Close()
}

func GetDB() *gorm.DB {
	return db
}

func BeforeInsert() {

}

func Create(c *gin.Context, model interface{}, data map[string]interface{}) {

	_ = c.BindJSON(&model) //绑定一个请求主体到一个类型

	db.Model(model).Create(data)
	return
	//db.Create(&Comment{
	//	Content: data["content"].(string),
	//	//CreateUserId:
	//	//MoreJson:      data["name"].(string),
	//	//CreatedAt: data["created_at"].(string),
	//	Status: data["status"].(string),
	//})

	//values := make([]interface{}, 0)
	//sql := "INSERT INTO `" + tablename + "` (" //+strings.Join(allFields, ",")+") VALUES ("
	//var ks []string
	//var vs []string
	//for k, v := range data { //注意：golang中对象的遍历，字段的排列是随机的
	//	ks = append(ks, "`"+k+"`") //保存所有字段
	//	vs = append(vs, "?")       //提供相应的占位符
	//	values = append(values, v) //对应保存相应的值
	//}
	////生成正常的插入语句
	//sql += strings.Join(ks, ",") + ") VALUES (" + strings.Join(vs, ",") + ")"
	//db.Exec(sql, values)
}

func Update(tablename string, params map[string]interface{}, id string) {
	values := make([]interface{}, 0)
	sql := "UPDATE `" + tablename + "` set " //+strings.Join(allFields, ",")+") VALUES ("
	var ks string
	index := 0
	psLen := len(params)
	for k, v := range params { //遍历对象
		index++
		values = append(values, v) //参数
		ks += "`" + k + "` =  ?"   //修改一个key的语句
		if index < psLen {         //非最后一个key，加逗号
			ks += ","
		}
	}
	values = append(values, id) //主键ID是单独的
	sql += ks + " WHERE id = ? "
	db.Exec(sql, values)
}

func Delete() {}
