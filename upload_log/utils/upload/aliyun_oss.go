package upload

import (
	"errors"
	"fmt"
	"log"
	"path"
	"time"

	"upload_log/global"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"go.uber.org/zap"
)

type AliyunOSS struct{}

func (*AliyunOSS) UploadFile(filename string) (string, string, error) {
	bucket, err := NewBucket()
	if err != nil {
		log.Println("function AliyunOSS.NewBucket() Failed", zap.Any("err", err.Error()))
		return "", "", errors.New("function AliyunOSS.NewBucket() Failed, err:" + err.Error())
	}

	dir, file := path.Split(filename)
	fmt.Printf("path.Split(%q) = dir: %q, file: %q\n", filename, dir, file) // dir: "static/", file: "myfile.css"

	// 上传阿里云路径 文件名格式 自己可以改 建议保证唯一性
	// yunFileTmpPath := filepath.Join("uploads", time.Now().Format("2006-01-02")) + "/" + file.Filename
	yunFileTmpPath := global.Config.AliyunOSS.BasePath + "/" + "uploads" + "/" + time.Now().Format("2006-01-02") + "/" + file

	// 上传文件(路径)。
	err = bucket.PutObjectFromFile(yunFileTmpPath, filename)
	if err != nil {
		log.Println("function formUploader.Put() Failed", zap.Any("err", err.Error()))
		return "", "", errors.New("function formUploader.Put() Failed, err:" + err.Error())
	}

	return global.Config.AliyunOSS.BucketUrl + "/" + yunFileTmpPath, yunFileTmpPath, nil
}

func (*AliyunOSS) DeleteFile(key string) error {
	bucket, err := NewBucket()
	if err != nil {
		log.Println("function AliyunOSS.NewBucket() Failed", zap.Any("err", err.Error()))
		return errors.New("function AliyunOSS.NewBucket() Failed, err:" + err.Error())
	}

	// 删除单个文件。objectName表示删除OSS文件时需要指定包含文件后缀在内的完整路径，例如abc/efg/123.jpg。
	// 如需删除文件夹，请将objectName设置为对应的文件夹名称。如果文件夹非空，则需要将文件夹下的所有object删除后才能删除该文件夹。
	err = bucket.DeleteObject(key)
	if err != nil {
		log.Println("function bucketManager.Delete() Filed", zap.Any("err", err.Error()))
		return errors.New("function bucketManager.Delete() Filed, err:" + err.Error())
	}

	return nil
}

func NewBucket() (*oss.Bucket, error) {
	// 创建OSSClient实例。
	client, err := oss.New(global.Config.AliyunOSS.Endpoint, global.Config.AliyunOSS.AccessKeyId, global.Config.AliyunOSS.AccessKeySecret)
	if err != nil {
		return nil, err
	}

	// 获取存储空间。
	bucket, err := client.Bucket(global.Config.AliyunOSS.BucketName)
	if err != nil {
		return nil, err
	}

	return bucket, nil
}
