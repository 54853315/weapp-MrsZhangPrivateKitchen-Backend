package util

import (
	"FoodBackend/pkg/setting"
	"github.com/satori/go.uuid"
	"os"
	"reflect"
	"strings"
)

func IsContain(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

func StructToMapViaReflect(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

func PrintStruct(obj interface{}) {
	t := reflect.TypeOf(obj)
	k := t.Kind()
	if k != reflect.Struct {
		Log.Fatal("This is not a struct.")
	}
	v := reflect.ValueOf(obj)
	Log.Notice("Fields:")
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		val := v.Field(i).Interface()
		Log.Noticef("%6s:%v = %v\n", f.Name, f.Type, val)
	}

	Log.Notice("Methods:")
	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		Log.Noticef("%6s %v\n", m.Name, m.Type)
	}
}

func Empty(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String:
		return value.Len() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return value.IsNil()
	}
	return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
}

func GetUniqueId() string {
	u4 := uuid.NewV4()
	return u4.String()
}

func GetFileBuff(filepath string) []byte {
	file, err := os.Open(filepath)
	if err != nil {
		Log.Fatal(err)
	}
	buff := make([]byte, 512)
	if _, err := file.Read(buff); err != nil {
		Log.Fatal(err)
	}
	return buff
}

func GetUrl(uri string) string {
	//为URI追加本站域名，适合操作图片URL
	//Log.Debugf("网址%s，检测http的结果值为%d", uri, strings.Index(uri, "http"))
	if strings.Index(uri, "http") == -1 {
		return setting.AppUrl + "/" + uri
	}
	return uri
}

func RemoveDomain(url string) string {
	//移除网址中的本站域名，适合操作图片URL
	return strings.Replace(url, setting.AppUrl+"/", "", 1)
}
