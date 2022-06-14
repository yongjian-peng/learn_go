package upload

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"upload_log/global"
	"upload_log/utils"

	"go.uber.org/zap"
)

type Local struct{}

//@author: [piexlmax](https://github.com/piexlmax)
//@author: [ccfish86](https://github.com/ccfish86)
//@author: [SliverHorn](https://github.com/SliverHorn)
//@object: *Local
//@function: UploadFile
//@description: 上传文件
//@param: file *multipart.FileHeader
//@return: string, string, error

func (*Local) UploadFile(filename string) (string, string, error) {

	dir, file := path.Split(filename)
	fmt.Printf("path.Split(%q) = dir: %q, file: %q\n", filename, dir, file) // dir: "static/", file: "myfile.css"

	// 读取文件后缀
	ext := path.Ext(filename)
	// 读取文件名并加密
	name := strings.TrimSuffix(filename, ext)
	name = utils.MD5V([]byte(name))
	// 拼接新文件名
	newfilename := name + "_" + time.Now().Format("20060102150405") + ext
	// 尝试创建此路径
	mkdirErr := os.MkdirAll(global.Config.Local.Path, os.ModePerm)
	if mkdirErr != nil {
		log.Println("function os.MkdirAll() Filed", zap.Any("err", mkdirErr.Error()))
		return "", "", errors.New("function os.MkdirAll() Filed, err:" + mkdirErr.Error())
	}
	// 拼接路径和文件名
	p := global.Config.Local.Path + "/" + newfilename

	f, openError := os.Open(filename) // 读取文件
	if openError != nil {
		log.Println("function file.Open() Filed", zap.Any("err", openError.Error()))
		return "", "", errors.New("function file.Open() Filed, err:" + openError.Error())
	}
	defer f.Close() // 创建文件 defer 关闭

	out, createErr := os.Create(p)
	if createErr != nil {
		log.Println("function os.Create() Filed", zap.Any("err", createErr.Error()))

		return "", "", errors.New("function os.Create() Filed, err:" + createErr.Error())
	}
	defer out.Close() // 创建文件 defer 关闭

	_, copyErr := io.Copy(out, f) // 传输（拷贝）文件
	if copyErr != nil {
		log.Println("function io.Copy() Filed", zap.Any("err", copyErr.Error()))
		return "", "", errors.New("function io.Copy() Filed, err:" + copyErr.Error())
	}
	return p, filename, nil
}

//@author: [piexlmax](https://github.com/piexlmax)
//@author: [ccfish86](https://github.com/ccfish86)
//@author: [SliverHorn](https://github.com/SliverHorn)
//@object: *Local
//@function: DeleteFile
//@description: 删除文件
//@param: key string
//@return: error

func (*Local) DeleteFile(key string) error {
	p := global.Config.Local.Path + "/" + key
	if strings.Contains(p, global.Config.Local.Path) {
		if err := os.Remove(p); err != nil {
			return errors.New("本地文件删除失败, err:" + err.Error())
		}
	}
	return nil
}
