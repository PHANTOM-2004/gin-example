package e

import log "github.com/sirupsen/logrus"

const (
	SUCCESS        = 200
	ERROR          = 500
	INVALID_PARAMS = 400

	ERROR_EXIST_TAG         = 10001
	ERROR_NOT_EXIST_TAG     = 10002
	ERROR_NOT_EXIST_ARTICLE = 10003

	ERROR_AUTH_CHECK_TOKEN_FAIL    = 20001
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT = 20002
	ERROR_AUTH_TOKEN               = 20003
	ERROR_AUTH                     = 20004
)

type Etype struct {
	code int
}

func NewEtype(code int) (res Etype) {
	res.code = code
	return
}

func (e *Etype) String() (res string) {
	switch e.code {
	case SUCCESS:
		res = "ok"
	case ERROR:
		res = "fail"
	case INVALID_PARAMS:
		res = "请求参数错误"
	case ERROR_EXIST_TAG:
		res = "已存在该标签"
	case ERROR_NOT_EXIST_TAG:
		res = "该标签不存在"
	case ERROR_NOT_EXIST_ARTICLE:
		res = "该文章不存在"
	case ERROR_AUTH_CHECK_TOKEN_FAIL:
		res = "Token鉴权失败"
	case ERROR_AUTH_CHECK_TOKEN_TIMEOUT:
		res = "Token已超时"
	case ERROR_AUTH_TOKEN:
		res = "Token生成失败"
	case ERROR_AUTH:
		res = "Token错误"
	default:
		log.Fatal("未知错误码")
	}
	return
}
