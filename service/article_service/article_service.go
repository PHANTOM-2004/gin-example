package article_service

import (
	"encoding/json"
	"gin-example/models"
	"gin-example/pkg/gredis"
	"gin-example/service/cache_service"

	log "github.com/sirupsen/logrus"
)

type Article struct {
	ID            int
	TagID         int
	Title         string
	Desc          string
	Content       string
	CoverImageUrl string
	State         int
	CreatedBy     string
	ModifiedBy    string

	PageNum  int
	PageSize int
}

func (a *Article) Get() (*models.Article, error) {
	cacheArticle := &models.Article{}

	cache := cache_service.Article{ID: a.ID}
	key := cache.GetArticleKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			log.Debug(err)
		}
		json.Unmarshal(data, cacheArticle)
		return cacheArticle, nil
	}

	article, err := models.GetArticle(a.ID)
	if err != nil {
		return nil, err
	}

	gredis.Set(key, article, 3600)
	return article, nil
}

func (a *Article) ExistByID() (exist bool, err error) {
	return models.ExistArticleByID(a.ID)
}

func (a *Article) Edit() error {
	data := map[string]any{
		"tag_id":          a.TagID,
		"title":           a.Title,
		"desc":            a.Desc,
		"content":         a.Content,
		"cover_image_url": a.CoverImageUrl,
		"state":           a.State,
		"modified_by":     a.ModifiedBy,
	}
	err := models.EditArticle(a.ID, data)
	log.Info(data)
	return err
}

func (a *Article) Add() error {
	data := map[string]any{
		"tag_id":          a.TagID,
		"title":           a.Title,
		"desc":            a.Desc,
		"content":         a.Content,
		"cover_image_url": a.CoverImageUrl,
		"state":           a.State,
		"modified_by":     a.ModifiedBy,
		"created_by":       a.CreatedBy,
	}
	err := models.AddArticle(data)
	log.Info(data)
	return err
}
