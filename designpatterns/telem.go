package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build xxxx.go
func main() {
	// 读取文件内容
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("无法打开文件")
		return
	}
	defer file.Close()

	// 创建一个Scanner用于读取文件内容
	scanner := bufio.NewScanner(file)
	appName := ""
	// 逐行读取文件内容
	mapValues := make(map[string]string, 0)
	for scanner.Scan() {
		line := scanner.Text()
		// 将多个空白字符替换为一个空白字符
		if strings.Contains(line, "{") && strings.Contains(line, "}") {
			// 获取 { } 中的内容
			left := strings.Index(line, "{")
			right := strings.Index(line, "}")
			content := line[left : right+1]
			// 将 { } 中的内容替换回去
			line = strings.Replace(line, content, strings.Replace(content, " ", "$$space$$", -1), 1)
			//替换掉
			line = strings.Replace(strings.Replace(line, "{", "", -1), "}", "", -1)
		}
		// 将多个空白字符替换为一个空白字符
		reg := regexp.MustCompile(`\s+`)
		line = reg.ReplaceAllString(line, " ")
		line = strings.TrimSpace(line)
		strList := strings.Split(line, " ")
		fmt.Println("正在处理=>", line)
		// 将 $$space$$ 替换回空格
		mapValues[strList[0]] = strings.Replace(strList[1], "$$space$$", " ", -1)
		if strList[0] == "AppName" {
			appName = mapValues[strList[0]]
		}
	}

	mbData, fErr := os.ReadFile("mb.txt")
	if fErr != nil {
		fmt.Println("无法读取文件")
		return
	}

	mbContent := string(mbData)
	//替换内容
	for key, value := range mapValues {
		mbContent = strings.Replace(mbContent, fmt.Sprintf("{{%s}}", key), value, -1)
	}
	fmt.Println(mbContent)

	os.WriteFile(fmt.Sprintf("%s.txt", appName), []byte(mbContent), 0777)
}
