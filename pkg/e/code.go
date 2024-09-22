package e

import log "github.com/sirupsen/logrus"

type Ecode = int

const (
	CACHE_ARTICLE = "ARTICLE"
	CACHE_TAG     = "TAG"
)

const (
	SUCCESS        Ecode = 200
	ERROR          Ecode = 500
	INVALID_PARAMS Ecode = 400

	ERROR_EXIST_TAG         Ecode = 10001
	ERROR_NOT_EXIST_TAG     Ecode = 10002
	ERROR_NOT_EXIST_ARTICLE Ecode = 10003

	ERROR_AUTH_CHECK_TOKEN_FAIL    Ecode = 20001
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT Ecode = 20002
	ERROR_AUTH_TOKEN               Ecode = 20003
	ERROR_AUTH                     Ecode = 20004

	ERROR_UPLOAD_IMAGE_INVALID_FMT Ecode = 30001
	ERROR_UPLOAD_IMAGE_CHECK_FAIL  Ecode = 30002
	ERROR_UPLOAD_IMAGE_SAVE_FAIL   Ecode = 30003
)

func String(e Ecode) (res string) {
	switch e {
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
	case ERROR_UPLOAD_IMAGE_CHECK_FAIL:
		res = "图片检查失败"
	case ERROR_UPLOAD_IMAGE_INVALID_FMT:
		res = "图片格式不合法"
	case ERROR_UPLOAD_IMAGE_SAVE_FAIL:
		res = "图片保存失败"

	default:
		log.Fatalf("未知错误码[%d]", e)
	}
	return
}
