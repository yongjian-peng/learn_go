package handingerror

import (
	"encoding/json"
	"io"
	"log"
)

// 每个错误只处理一次 每个代码掉一层 一层 记录日志 返回上层错误信息
func WriteAll(w io.Writer, buf []byte) error {
	_, err := w.Write(buf)

	if err != nil {
		log.Println("unable to write:", err)
	}
	return nil
}

func WriteConfig(w io.Writer, conf *Config) error {
	buf, err := json.Marshal(conf)
	if err != nil {
		log.Printf("could not marshal config: %v", err)
		return err
	}
	if err := WriteAll(w, buf); err != nil {
		log.Println("could not write config:%v", err)
		return err
	}
	return nil
}
