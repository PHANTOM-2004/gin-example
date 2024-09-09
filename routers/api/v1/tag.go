package v1

import (
	"gin-example/models"
	"gin-example/pkg/e"
	"gin-example/pkg/setting"
	"gin-example/util"
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
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

// 新增文章tag
func AddTag(c *gin.Context) {
	name := c.Query("name")
	state_str := c.DefaultQuery("state", "0")
	state := com.StrTo(state_str).MustInt()
	createdBy := c.Query("created_by")

	log.WithFields(log.Fields{
		"name":       name,
		"state":      state,
		"created_by": createdBy,
	}).Debug("query for add tag")

	valid := validation.Validation{}
	valid.Required(name, "name").Message("名称不能为空")
	valid.MaxSize(name, 100, "name").Message("名称最长100字符")
	valid.Required(createdBy, "created_by").Message("创建人最长100字符")
	valid.Range(state, 0, 1, "state").Message("状态只能为0/1")

	tag_exist := models.ExistTag(name)
	verror := valid.HasErrors()

	code := e.INVALID_PARAMS

	if !verror && !tag_exist {
		// 没错误发生, 并且这个tag还不存在, 意味着成功插入
		code = e.SUCCESS
		models.AddTag(name, state, createdBy)
	}
	// tag存在, 不能重复插入
	if tag_exist {
		log.WithFields(log.Fields{
			"tag_name": name,
		}).Warn(e.String(code))

		code = e.ERROR_EXIST_TAG
	}

	// 存在数据错误
  //
	for _, err := range valid.Errors {
		log.Warnf("validation error: [%s]: %s", err.Key, err.Message)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.String(code),
		"data": make(map[string]string),
	})
}

// 修改文章标签
func EditTag(c *gin.Context) {
}

// 删除文章标签
func DeleteTag(c *gin.Context) {
}
