package upload

import (
	"gin-example/pkg/setting"
	"gin-example/pkg/util"
	"os"
	"path/filepath"
	"strings"
)


func GetImageFullUrl(name string) string {
	return setting.AppSetting.ImagePrefixUrl + "/" + GetImagePath() + name
}

func GetImageMaxSize() int64{
  return setting.AppSetting.ImageMaxSize
}

func GetImageName(name string) string {
	ext := filepath.Ext(name)
	base := filepath.Base(name)
	f := util.EncodeMD5(strings.TrimSuffix(base, ext))

	return f + ext
}

func GetImagePath() string {
	return setting.AppSetting.ImageSavePath
}


func GetImageFullPath() string {
	return setting.AppSetting.RuntimeRootPath + GetImagePath()
}

func CheckImageExt(fileName string) bool {
	// No dots: ""
	// One dot: ".js"
	// Two dots: ".js"
	ext := filepath.Ext(fileName)
	for _, allowExt := range setting.AppSetting.ImageAllowExts {
		equal := strings.EqualFold(
			strings.ToUpper(allowExt),
			strings.ToUpper(ext),
		)
		if equal {
			return true
		}
	}

	return false
}

func CheckImage(src string) error {
	dir, err := os.Getwd()
	if err != nil {
		// log.Fatalf("os.Getwd err: %v", err)
		return err
	}

	_, err = os.Stat(src)

	perm := os.IsPermission(err)
	notExist := os.IsNotExist(err)

	if perm {
		return err
	}

	if notExist {
		path := dir + "/" + src
		// 如果路径不存在那么创建这个路径
		err = os.Mkdir(path, 0600)
		if err != nil {
			return err
		}
	}

	return err
}
