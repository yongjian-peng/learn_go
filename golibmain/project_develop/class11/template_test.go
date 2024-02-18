package main

import (
	"fmt"
	"testing"
)

func TestTemplate(t *testing.T) {
	// 做西红柿
	xihongshi := &XiHongShi{}
	doCook(xihongshi)

	fmt.Println("\n=====> 做另外一道菜")

	// 做炒鸡蛋
	chaojidan := &ChaoJiDan{}
	doCook(chaojidan)
}
