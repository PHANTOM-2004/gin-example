package v1

import (
	"gin-example/models"
	"gin-example/pkg/e"
	"gin-example/pkg/setting"
	"gin-example/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

// 获取多个文章标签
func GetTags(c *gin.Context) {
	/* GET /path?id=1234&name=Manu&value=
	   c.Query("id") == "1234"
	   c.Query("name") == "Manu"
	   c.Query("value") == ""
	   c.Query("wtf") == ""
	*/
	name := c.Query("name")

	maps := make(map[string]any)
	data := make(map[string]any)

	// 这里的maps是为了存储相应的信息,在这个地方就是
	// 存储着name的键值对, 以及state的键值对
	if name != "" {
		maps["name"] = name
	}
	state := -1
	if arg := c.Query("state"); arg != "" {
		// get the state arguments
		state = com.StrTo(arg).MustInt()
		maps["state"] = state
	}

	code := e.SUCCESS

	// 我们通过model进行查询,查询对象是maps
	data["list"] = models.GetTags(
		util.GetPage(c),
		setting.PageSize,
		maps,
	)

	// 利用model进行查询, 查询对象是maps, 这里只记录数量
	data["total"] = models.GetTagTotal(maps)

	c.JSON(
		http.StatusOK,
		gin.H{
			"code": code,
			"msg":  e.String(code),
			"data": data,
		},
	)
}

// 新增文章标签
func AddTag(c *gin.Context) {
}

// 修改文章标签
func EditTag(c *gin.Context) {
}

// 删除文章标签
func DeleteTag(c *gin.Context) {
}
