package dto

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"strings"
)

func init() {
	// Register custom validate methods
	//if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
	//	//_ = v.RegisterValidation("customValidate", customValidate)
	//	_ = v.RegisterValidation("pwdValidate", pwdValidate)
	//	//_ = v.RegisterValidation("permsValidate", permsValidate)
	//} else {
	//	util.Log.Fatal("Gin fail to registered custom validator(v10)")
	//}
}

//var validate *validator.Validate
// Bind : bind request dto and auto verify parameters
func Bind(c *gin.Context, obj interface{}) error {
	_ = c.ShouldBindUri(obj)
	if err := c.ShouldBind(obj); err != nil {
		if fieldErr, ok := err.(validator.ValidationErrors); ok {
			var tagErrorMsg []string
			//util.Log.Debug("fieldErr:", fieldErr)
			for _, v := range fieldErr {
				if _, has := ValidateErrorMessage[v.Tag()]; has {
					//v.Tag() like requests
					//v.Field() like `json:` name
					//util.Log.Noticef(reflect.FuncOf(v.Tag()))
					value := v.Value()
					field := v.Field()

					//util.Log.Noticef("%T,%T", v.Field, v.Value)
					//util.Log.Noticef("%v,%v", v.Field, v.Value)
					//util.Log.Noticef("%#v,%#v", v.Field, v.Value)
					//util.Log.Noticef("%+v,%+v", v.Field, v.Value)
					//util.Log.Noticef("%%,%%", v.Field, v.Value)
					//util.Log.Noticef("%x,%x", v.Field, v.Value)

					tagErrorMsg = append(tagErrorMsg, fmt.Sprintf(ValidateErrorMessage[v.Tag()], field, value))
				} else {
					tagErrorMsg = append(tagErrorMsg, err.Error())
				}
			}
			return errors.New(strings.Join(tagErrorMsg, ","))
		}
	}
	return nil
}

//ValidateErrorMessage : customize error messages
var ValidateErrorMessage = map[string]string{
	"customValidate": "%s can not be %s",
	"required":       "%v is required,got empty %v",
	"oneof":          "%v is not valid parameter,%v is not legal",
	"pwdValidate":    "%s is not a valid password",
	"permsValidate":  "%s contains comma",
}
