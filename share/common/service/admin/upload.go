package admin

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"os"
	"path"
	"share/common/pkg/appError"
	"share/common/pkg/config"
	"share/common/pkg/goutils"
)

type UploadService struct {
	*Service
}

func Upload(c *fiber.Ctx) *UploadService {
	return &UploadService{Service: NewService(c, "uploadService")}
}

func (s *UploadService) UploadImg() error {
	// 从表单字段 "document"获取第一个文件:
	file, err := s.C.FormFile("file")
	if err != nil {
		return s.Error(appError.NewError(err.Error()))
	}

	contentType := file.Header.Get("Content-Type")
	if !goutils.InSlice[string](contentType, []string{"image/png", "image/jpeg", "image/gif"}) {
		return s.Error(appError.NewError("图片类型错误"))
	}

	fileExt := path.Ext(file.Filename)
	if !goutils.InSlice[string](fileExt, []string{".png", ".jpg", ".jpeg", ".gif"}) {
		return s.Error(appError.NewError("图片类型错误"))
	}

	filePath := fmt.Sprintf("/assets/upload/images/%s%s", goutils.Md5(fmt.Sprintf("%d:%d", s.AdminId, goutils.GetCurTimeUnixNano())), fileExt)
	saveFilePath := fmt.Sprintf("./views%s", filePath)
	err = s.C.SaveFile(file, saveFilePath)
	if err != nil {
		return s.Error(appError.NewError(err.Error()))
	}

	//检查文件内容是否是图片格式
	if !goutils.CheckIsImageFile(saveFilePath) {
		_ = os.Remove(saveFilePath)
		return s.Error(appError.NewError("图片类型错误"))
	}

	return s.Success(fiber.Map{
		"fileUrl": fmt.Sprintf("%s%s", config.AppConfig.Server.AssertUrl, filePath),
	})
}
