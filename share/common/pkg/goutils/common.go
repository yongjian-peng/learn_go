package goutils

import (
	"github.com/valyala/bytebufferpool"
	"go.uber.org/zap"
	"share/common/pkg/logger"
	"time"
)

const (
	TimeLayout   = "2006-01-02 15:04:05"
	TimeLayout_2 = "20060102150405"
	DateLayout   = "2006-01-02"
	NULL         = ""
)

type File struct {
	Name    string `json:"name"`
	Content []byte `json:"content"`
}

// request_id 请求的全局的id
// tag、detailed 表示超时发生位置的两个字符串参数。
// start 程序开始执行的时间
// timeLimit  函数执行超时阀值，单位是秒。
func ExecutionTime(logFileName, requestId, tag, detailed string, start time.Time, timeLimit float64) {
	totalTime := time.Now().Sub(start).Seconds()
	logger.ApiInfo(logFileName, requestId, tag+detailed, zap.Float64("totalTime", totalTime))
	if totalTime > timeLimit {
		// 可以执行的逻辑 如果操作了指定的时间 则触发机制
	}
}

// ParseParamSquareBrackets
func ParseParamSquareBrackets(k string) (string, error) {
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)

	kbytes := []byte(k)

	for i, b := range kbytes {

		if b == '[' && kbytes[i+1] != ']' {
			if err := bb.WriteByte('.'); err != nil {
				return "", err
			}
		}

		if b == '[' || b == ']' {
			continue
		}

		if err := bb.WriteByte(b); err != nil {
			return "", err
		}
	}

	return bb.String(), nil
}
