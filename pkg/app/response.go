package app

import (
	"gin-example/pkg/e"

	"github.com/gin-gonic/gin"
)

func Response(c *gin.Context, httpCode int, errCode int, data any) {
	c.JSON(httpCode, gin.H{
		"code": httpCode,
		"msg":  e.String(errCode),
		"data": data,
	})
}
