package util

import (
	"reflect"
)

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
