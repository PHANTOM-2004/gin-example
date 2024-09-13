package util

import (
	"gin-example/pkg/setting"

	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

func GetPage(c *gin.Context) int {
	result := 0
	p := c.Query("page") // query the key in url,
  // if there is no key, return ""
  //
	page, _ := com.StrTo(p).Int()
	if page > 0 {
		result = (page - 1) * setting.AppSetting.PageSize
	}

	return result
}
