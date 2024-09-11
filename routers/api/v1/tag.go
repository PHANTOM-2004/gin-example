package v1

import (
	"gin-example/models"
	"gin-example/pkg/e"
	"gin-example/pkg/setting"
	"gin-example/pkg/util"
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

	valid := validation.Validation{}
	valid.Required(name, "name").Message("名称不能为空")
	valid.MaxSize(name, 100, "name").Message("名称最长100字符")
	valid.Required(createdBy, "created_by").Message("创建人最长100字符")
	valid.Range(state, 0, 1, "state").Message("状态只能为0/1")

	verror := valid.HasErrors()

	code := e.INVALID_PARAMS

	if !verror {
		// 没错误发生, 并且这个tag还不存在, 意味着成功插入
		tag_exist := models.ExistTag(name)
		if !tag_exist {
			code = e.SUCCESS
			models.AddTag(name, state, createdBy)
		} else {
			// tag存在, 不能重复插入
			log.WithFields(log.Fields{
				"tag_name": name,
			}).Warn(e.String(code))
			code = e.ERROR_EXIST_TAG

		}
	}

	// 存在数据错误
	//
	for _, err := range valid.Errors {
		log.Warnf("validation error: [%s]: %s", err.Key, err.Message)
	}

	log.WithFields(log.Fields{
		"name":       name,
		"state":      state,
		"created_by": createdBy,
	}).Debug("add tag:", e.String(code))

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.String(code),
		"data": make(map[string]string),
	})
}

// 修改文章标签
// apiv1.PUT("/tags/:id", v1.EditTag)

func EditTag(c *gin.Context) {
	id := com.StrTo(c.Query("id")).MustInt() // 修改的时候使用的是id进行修改
	name := c.Query("name")
	modifiedBy := c.Query("modified_by")

	valid := validation.Validation{}

	state := -1

	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只能是0或者1")
	}

	valid.Required(id, "id").Message("id不能为空") // 注意这里id制定的是字段名称
	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人名称最长100字符")
	valid.MaxSize(name, 100, "name").Message("名称最长100字符")

	code := e.INVALID_PARAMS
	verror := valid.HasErrors()

	if !verror {
		tag_exist := models.ExistTag(name)

		if tag_exist {
			code = e.SUCCESS
			// 如果标签存在才能edit, 这是显而易见的
			data := make(map[string]any)
			data["modified_by"] = modifiedBy // 这里传递修改人的姓名
			if name != "" {
				// NOTE: 这里为什么要记录name?我还没有搞明白
				data["name"] = name // 这里传递
			}
			if state != -1 {
				data["state"] = state
			}

			models.EditTag(id, data)
		} else {
			// 不存在, 直接抛出错误
			code = e.ERROR_NOT_EXIST_TAG
			// 编辑不存在的标签
			log.Warnf("try edit non-existed tag[%s], aborted\n", name)
		}
	}

	// 存在数据错误
	for _, err := range valid.Errors {
		log.Warnf("validation error: [%s]: %s", err.Key, err.Message)
	}
	// debug日志记录
	log.WithFields(log.Fields{
		"tag_name":    name,
		"modified_by": modifiedBy,
		"state":       state,
	}).Debug("edit_tag:", e.String(code))

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.String(code),
		"data": make(map[string]string),
	})
}

// 删除文章标签
func DeleteTag(c *gin.Context) {
	id := com.StrTo(c.Query("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	code := e.INVALID_PARAMS
	verror := valid.HasErrors()

	if !verror {
		code = e.SUCCESS
		tag_exist := models.ExistTagByID(id)
		if tag_exist {
			models.DeleteTag(id)
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.String(code),
		"data": make(map[string]string),
	})
}
