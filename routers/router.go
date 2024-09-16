package routers

import (
	"gin-example/docs"
	"gin-example/middleware/jwt"
	"gin-example/pkg/setting"
	"gin-example/pkg/upload"
	"gin-example/routers/api"
	"net/http"

	"github.com/gin-gonic/gin"

	v1 "gin-example/routers/api/v1"

	log "github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter() *gin.Engine {
	defer log.Info("router initialized")

	r := gin.New()

	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	gin.SetMode(setting.ServerSetting.RunMode)

	{
		r.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "test",
			})
		})

		r.GET("/auth", api.GetAuth)

		r.POST("/upload", api.UploadImage)

		r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))

		// add swagger support
		// 我们可以使用gin-swagger产生swagg而.json, 利用swagger进行serve.
		// 也就是首先swag init, 得到`json:`, 然后利用go-swagger进行serve
		// go install github.com/go-swagger/go-swagger/cmd/swagger@latest
		docs.SwaggerInfo.BasePath = "/api/v1"
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	apiv1 := r.Group("/api/v1")
	apiv1.Use(jwt.JWT())

	{
		apiv1.GET("/tags", v1.GetTags)
		apiv1.POST("/tags", v1.AddTag)
		apiv1.PUT("/tags/:id", v1.EditTag)
		apiv1.DELETE("/tags/:id", v1.DeleteTag)
	}

	{
		// 获取文章列表
		apiv1.GET("/articles", v1.GetArticles)
		// 获取指定文章
		apiv1.GET("/articles/:id", v1.GetArticle)
		// 新建文章
		apiv1.POST("/articles", v1.AddArticle)
		// 更新指定文章
		apiv1.PUT("/articles/:id", v1.EditArticle)
		// 删除指定文章
		apiv1.DELETE("/articles/:id", v1.DeleteArticle)
	}

	return r
}
