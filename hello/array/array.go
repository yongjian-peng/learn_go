package main

import "fmt"

func main()  {
	//这里我们使用range来统计一个slice的元素个数
	//数组也可以采用这种方法
	nums := []int{2,3,4}
	sum := 0
	for _,num := range nums{
		sum += num
	}
	fmt.Println("sum:",sum)

	//range在数组和slice中同样提供每个项的索引和值
	//上面我们不需要索引，所以我们使用空值定义符`_`来忽略它
	//有时候我们实际上是需要这个索引的
	for i,num := range nums{
		if num == 3{
			fmt.Println("index:",i)
		}
	}

	//range在map中迭代键值树
	kvs := map[string]string{"a":"apple","b":"banana"}
	for k,v := range kvs{
		fmt.Println("%s->%s\n",k,v)
	}

	var payAdapterImpl = map[string][]string{
		"PAYTM":    {"wechat.wechat", "wechat.hk"},
		"CASHFREE": {"alipay.linx", "alipay.yedpay"},
	}

	for k,v := range payAdapterImpl{
		fmt.Println("%s->%s\n",k,v)
	}

	fmt.Println(kvs["b"])

	//range在字符串中迭代unicode编码。
	//第一个返回值是rune的起始字节位置，然后第二个是rune自己
	for i,c := range "go1"{
		fmt.Println(i,c)
	}
}