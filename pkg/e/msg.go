package e

var MsgFlags = map[int]string{
	SUCCESS:                        "ok",
	ERROR:                          "fail",
	INVALID_PARAMS:                 "请求参数错误",
	UNAUTHORIZED:                   "未授权",
	NOT_FOUND:                      "资源无法找到",
	ERROR_EXIST_TAG:                "已存在该标签名称",
	ERROR_NOT_EXIST_TAG:            "该标签不存在",
	ERROR_NOT_EXIST_ARTICLE:        "该文章不存在",
	ERROR_NOT_EXIST:                "该内容不存在",
	ERROR_AUTH_CHECK_TOKEN_FAIL:    "Token鉴权失败",
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT: "Token已超时",
	ERROR_AUTH_TOKEN:               "Token生成失败",
	ERROR_AUTH:                     "Token错误",

	ERROR_JS_CODE: "请传递正确的js_code",
	ID_MISSING:    "ID必须",

	BOOK_STATUS_FAIL: "status错误",
	BOOK_EXISTS:      "食物已经创建过",
	BOOK_NOT_EXISTS:  "食物不存在",
}

var StatusFlags = map[string]string{
	"已发布": "publish",
	"未发布": "private",
}

func GetStatusMsg(str string) string {
	msg, ok := StatusFlags[str]
	if ok {
		return msg
	}
	return ""
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
