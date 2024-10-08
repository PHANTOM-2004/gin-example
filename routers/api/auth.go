package api

import (
	"gin-example/models"
	"gin-example/pkg/e"
	"gin-example/pkg/util"
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

func GetAuth(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	valid := validation.Validation{}
	a := auth{Username: username, Password: password}
	ok, _ := valid.Valid(&a)

	data := make(map[string]any)
	code := e.INVALID_PARAMS

	log.Debug("getting auth")

	if ok {
		isExist := models.CheckAuth(username, password)
		if isExist {
			token, err := util.GenerateToken(username, password)

			if err != nil {
				code = e.ERROR_AUTH_TOKEN
			} else {
				data["token"] = token

				code = e.SUCCESS
			}

		} else {
			code = e.ERROR_AUTH
		}

	} else {
		for _, err := range valid.Errors {
			log.Warn(err.Key, err.Message)
		}
	}

	log.Debug("response auth")
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.String(code),
		"data": data,
	})
}
