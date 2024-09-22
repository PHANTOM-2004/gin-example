package v1

import (
	"gin-example/models"
	"gin-example/pkg/app"
	"gin-example/pkg/e"
	"gin-example/pkg/setting"
	"gin-example/pkg/util"
	"gin-example/service/article_service"
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/unknwon/com"
)

// 获取单个文章
func GetArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	// 如果有输入error首先返回
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		app.Response(c, http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	articleService := article_service.Article{ID: id}
	exist, err := articleService.ExistByID()
	// TODO: give right error code
	if err != nil {
		app.Response(c, http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	if !exist {
		app.Response(c, http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	article, err := articleService.Get()
	if err != nil {
		app.Response(c, http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	app.Response(c, http.StatusOK, e.SUCCESS, article)
}

// 获取多个文章
func GetArticles(c *gin.Context) {
	data := make(map[string]interface{})
	maps := make(map[string]interface{})
	valid := validation.Validation{}

	if arg := c.Query("state"); arg != "" {
		state := com.StrTo(arg).MustInt()
		maps["state"] = state
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	if arg := c.Query("tag_id"); arg != "" {
		tagId := com.StrTo(arg).MustInt()
		maps["tag_id"] = tagId
		valid.Min(tagId, 1, "tag_id").Message("标签ID必须大于0")
	}

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		app.Response(c, http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	data["lists"] = models.GetArticles(util.GetPage(c), setting.AppSetting.PageSize, maps)
	data["total"] = models.GetArticleTotal(maps)

	app.Response(c, http.StatusOK, e.SUCCESS, data)
}

// 新增文章
func AddArticle(c *gin.Context) {
	tagID := com.StrTo(c.Query("tag_id")).MustInt()
	title := c.Query("title")
	desc := c.Query("desc")
	content := c.Query("content")
	createdBy := c.Query("created_by")
	state := com.StrTo(c.DefaultQuery("state", "0")).MustInt()
	cover_image_url := c.Query("cover_image_url")

	// do some validation
	valid := validation.Validation{}
	valid.Min(tagID, 1, "tag_id").Message("标签ID必须大于0")
	valid.Required(title, "title").Message("标题不能为空")
	valid.Required(desc, "desc").Message("简述不能为空")
	valid.Required(content, "content").Message("内容不能为空")
	valid.Required(createdBy, "created_by").Message("创建人不能为空")
	valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	// 暂且可以允许是空的
	valid.Range(len(cover_image_url), 0, 256, "cover_image_url").Message("文章封面最长路径为256")

	// check http arguments
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		app.Response(c, http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	if !models.ExistTagByID(tagID) {
		app.Response(c, http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	as := article_service.Article{
		TagID:         tagID,
		Title:         title,
		Desc:          desc,
		Content:       content,
		CreatedBy:     createdBy,
		State:         state,
		CoverImageUrl: cover_image_url,
	}

	err := as.Add()
	if err != nil {
		log.Fatal(err)
	}

	app.Response(c, http.StatusOK, e.SUCCESS, nil)
}

// 修改文章
func EditArticle(c *gin.Context) {
	// TODO: don't mix valid here
	valid := validation.Validation{}

	id := com.StrTo(c.Param("id")).MustInt()
	tagID := com.StrTo(c.Query("tag_id")).MustInt()
	title := c.Query("title")
	desc := c.Query("desc")
	content := c.Query("content")
	modifiedBy := c.Query("modified_by")
	cover_image_url := c.Query("cover_image_url")

	if arg := c.Query("state"); arg != "" {
		state := com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	valid.Min(id, 1, "id").Message("ID必须大于0")
	valid.MaxSize(title, 100, "title").Message("标题最长为100字符")
	valid.MaxSize(desc, 255, "desc").Message("简述最长为255字符")
	valid.MaxSize(content, 65535, "content").Message("内容最长为65535字符")
	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长为100字符")
	// 修改的时候认为必须是非空
	valid.Range(len(cover_image_url), 1, 256, "cover_image_url").Message("文章封面路径非空且最长路径为256")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		app.Response(c, http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	if existArticle, err := models.ExistArticleByID(id); !existArticle {
		if err != nil {
			log.Fatal(err)
			return
		}
		app.Response(c, http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	if existTag := models.ExistTagByID(tagID); !existTag {
		app.Response(c, http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	as := article_service.Article{
		ID:            id,
		TagID:         tagID,
		Desc:          desc,
		Content:       content,
		ModifiedBy:    modifiedBy,
		CoverImageUrl: cover_image_url,
	}
	err := as.Edit()
	if err != nil {
		log.Fatal(err)
	}

	app.Response(c, http.StatusOK, e.SUCCESS, nil)
}

// 删除文章
func DeleteArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		app.Response(c, http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	models.DeleteArticle(id)

	app.Response(c, http.StatusOK, e.SUCCESS, nil)
}
