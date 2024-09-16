package api

import (
	"gin-example/pkg/e"
	"gin-example/pkg/upload"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func UploadImage(c *gin.Context) {
	code := e.SUCCESS
	data := make(map[string]string)

	_, image, err := c.Request.FormFile("image")
	if err != nil {
		log.Warn(err)

		code = e.ERROR
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  e.String(code),
			"data": data,
		})
	}

	if image == nil {
		code = e.INVALID_PARAMS
	} else {
		imageName := upload.GetImageName(image.Filename)
		fullPath := upload.GetImageFullPath()
    log.Info(fullPath)
		savePath := upload.GetImagePath()

		src := fullPath + imageName

		sizeValid := image.Size <= upload.GetImageMaxSize()
		extValid := upload.CheckImageExt(imageName)

		if !extValid || !sizeValid {
			log.Warnf("ext: %v, size: %v", extValid, sizeValid)
			code = e.ERROR_UPLOAD_IMAGE_INVALID_FMT
		} else if err := upload.CheckImage(fullPath); err != nil {
			// valid to create
			log.Warn(err)
			code = e.ERROR_UPLOAD_IMAGE_CHECK_FAIL
		} else if err := c.SaveUploadedFile(image, src); err != nil {
			// unable to save
			log.Warn(err)
			code = e.ERROR_UPLOAD_IMAGE_SAVE_FAIL
		} else {
			// no error, check passed
			log.Info("image uploaded")
			data["image_url"] = upload.GetImageFullUrl(imageName)
			data["image_save_url"] = savePath + imageName
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.String(code),
		"data": data,
	})
}
