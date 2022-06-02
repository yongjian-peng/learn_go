package logic

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
	"unsafe"
	"upload_log/global"
	"upload_log/utils/filehash"
)

func BuildOssDataName(filename string) string {
	osspath := global.Config.UploadLogData.OssObjectPrefix
	// 构造 上传oss 路径+名称
	var build strings.Builder
	build.WriteString(osspath)
	build.WriteString(filename)

	return build.String()
}

func BuildLocalDataName(filename string) string {
	uploadPath := global.Config.UploadLogData.BizSearchDir
	// 构造 上传oss 路径+名称
	var build strings.Builder
	build.WriteString(uploadPath)
	build.WriteString(filename)

	return build.String()
}

func BuildFileHash(filepath string) string {
	f, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	start := time.Now()
	defer func() {
		fmt.Println("HashWriter 耗时：", time.Now().Sub(start))
	}()

	hw := filehash.NewHashWriter(false, true, false)
	defer hw.Close()

	b := make([]byte, 32*1024)
	for {
		n, err := f.Read(b)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("读文件出错：", err)
			return ""
		}
		hw.Write(b[:n])
	}

	md5Buf, sha1Buf, sha256Buf := hw.Sum(nil)

	//hex.EncodeToString(md5Buf)
	fmt.Printf("md5: %x\n", md5Buf)
	fmt.Printf("sha1: %x\n", sha1Buf)
	fmt.Printf("sha256: %x\n", sha256Buf)

	return hex.EncodeToString(sha1Buf)

	// return BytesToString(sha1Buf)

	// return BytesToString(sha1Buf)

}

// string 装 []byte
func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

// []byte 转 string
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
