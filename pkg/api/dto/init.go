package dto

import (
	"FoodBackend/pkg/matchers"
	"FoodBackend/pkg/util"
	"bytes"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"io"
	"net/http"
	"strings"
)

//NOTE 记得一定要用"github.com/go-playground/validator/v10"

func init() {
	// Register custom validate methods
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		//_ = v.RegisterValidation("imageValidate", imageValidate)
		_ = v.RegisterValidation("pwdValidate", pwdValidate)
		//_ = v.RegisterValidation("permsValidate", permsValidate)
	} else {
		util.Log.Fatal("Gin fail to registered custom validator(v10)")
	}
}

//var validate *validator.Validate
// Bind : bind request dto and auto verify parameters
func Bind(c *gin.Context, obj interface{}) error {
	_ = c.ShouldBindUri(obj)
	var tagErrorMsg []string
	if err := c.ShouldBind(obj); err != nil {
		if fieldErr, ok := err.(validator.ValidationErrors); ok {
			for _, v := range fieldErr {
				if _, has := ValidateErrorMessage[v.Tag()]; has {
					value := v.Value()
					field := v.Field()
					tagErrorMsg = append(tagErrorMsg, fmt.Sprintf(ValidateErrorMessage[v.Tag()], field, value))
				} else {
					tagErrorMsg = append(tagErrorMsg, err.Error())
				}
			}
		} else if err.Error() == "unexpected end of JSON input" {
			tagErrorMsg = append(tagErrorMsg, "incomplete data!")
		}
	} else if strings.Index(c.GetHeader("content-type"), "multipart/form-data") != -1 {
		file, _, err := c.Request.FormFile("file")
		defer file.Close()
		if err != nil {
			tagErrorMsg = append(tagErrorMsg, err.Error())
		}
		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, file); err != nil {
			util.Log.Fatal(err)
			tagErrorMsg = append(tagErrorMsg, err.Error())
		}

		fileType := http.DetectContentType(buf.Bytes())
		util.Log.Notice("上传的真实文件类型为：", fileType)

		if !matchers.Image(fileType) {
			tagErrorMsg = append(tagErrorMsg, "file incorrect.")
		}
	}
	if len(tagErrorMsg) > 0 {
		return errors.New(strings.Join(tagErrorMsg, ","))
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
