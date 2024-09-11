package jwt

import (
	"gin-example/pkg/e"
	"gin-example/pkg/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func JWT() gin.HandlerFunc {
	return jwtFunc
}

func jwtFunc(c *gin.Context) {
	var code int
	var data any

	code = e.SUCCESS
	token := c.Query("token")

	if token == "" {
		code = e.INVALID_PARAMS
	} else {
		log.Infof("get jwt token[%s]", token)
		claims, err := util.ParseToken(token)
    //验证结果
		if err != nil {
			code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
		} else if time.Now().Unix() > claims.ExpiresAt {
			code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
		}
	}

	if code != e.SUCCESS {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": code,
			"msg":  e.String(code),
			"data": data,
		})

		c.Abort()
		return
	}

	c.Next()
}
