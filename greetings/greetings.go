package greetings

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// Hello returns a greeting for the named preson
func Hello(name string) (string, error) {
	if name == "" {
		return "", errors.New("empty name")
	}
	message := fmt.Sprintf(randomFormat(), name)
	return message, nil
}

// Hellos 返回一个映射，将每个命名的人关联起来
func Hellos(names []string) (map[string]string, error) {
	// 将名称与消息关联的映射
	messages := make(map[string]string)
	// 遍历接收到的名称切片，调用
	// Hello 函数为每个名字获取一条消息
	for _, name := range names {
		message, err := Hello(name)
		if err != nil {
			return nil, err
		}
		messages[name] = message
	}
	return messages, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randomFormat() string {
	formats := []string{
		"Hi, %v, Welcome!",
		"Great to see you %v",
		"Hail, %v Well met!",
	}
	return formats[rand.Intn(len(formats))]
}
