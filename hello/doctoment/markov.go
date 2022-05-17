// 版权所有 2011 Go 作者。版权所有。
// 此源代码的使用由 BSD 样式管理
// 可以在 LICENSE 文件中找到的许可证。

/*
生成随机文本：马尔可夫链算法

基于“设计和实现”一章中介绍的程序
《编程实践》（Kernighan 和 Pike，Addison-Wesley 1999）。
另见 Computer Recreations, Scientific American 260, 122 - 125 (1989)。

马尔可夫链算法通过创建统计模型来生成文本
给定前缀的潜在文本后缀。考虑这段文字：

我不是数字！我是自由人！

我们的马尔可夫链算法会将此文本排列到这组前缀中
和后缀，或“链”：（此表假定前缀长度为两个单词。）

前缀 后缀

““ ““        一世
““ 我是
我是一个
我不是
一个自由的人！
我是免费的
我不是
一个号码！一世
数字！我是
不是数字！

要使用此表生成文本，我们选择一个初始前缀（“I am”，例如
例如），随机选择与该前缀关联的后缀之一
概率由输入统计信息（“a”）确定，
然后通过从前缀中删除第一个单词来创建一个新前缀
并附加后缀（使新前缀为“am a”）。重复这个过程
直到我们找不到当前前缀的任何后缀或超出单词
限制。 （字数限制是必要的，因为链表可能包含循环。）

我们的这个程序版本从标准输入读取文本，将其解析为
马尔可夫链，并将生成的文本写入标准输出。
可以使用 -prefix 和 -words 指定前缀和输出长度
命令行上的标志。
*/
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Prefix 是一个或多个单词的马尔可夫链前缀。
type Prefix []string

// String 以字符串形式返回前缀（用作映射键）
func (p Prefix) String() string {
	return strings.Join(p, " ")
}

// Shift 从前缀中删除第一个单词并附加给定的单词
func (p Prefix) Shift(word string) {
	copy(p, p[1:])
	p[len(p)-1] = word
}

// 链包含前缀到后缀列表的映射（“链”）。
// 前缀是由空格连接的 prefixLen 字串。
// 后缀是一个单词。 一个前缀可以有多个后缀。
type Chain struct {
	chain     map[string][]string
	prefixLen int
}

// NewChain 返回一个带有 prefixLen 单词前缀的新链。
func NewChain(prefixLen int) *Chain {
	return &Chain{make(map[string][]string), prefixLen}
}

// Build 从提供的 Reader 中读取文本，并
// 将其解析为存储在 Chain 中的前缀和后缀。
func (c *Chain) Build(r io.Reader) {
	br := bufio.NewReader(r)
	p := make(Prefix, c.prefixLen)
	for {
		var s string
		if _, err := fmt.Fscan(br, &s); err != nil {
			break
		}
		key := p.String()
		c.chain[key] = append(c.chain[key], s)
		p.Shift(s)
	}
}

// Generate 返回从 Chain 生成的最多 n 个单词的字符串。
func (c *Chain) Generate(n int) string {
	p := make(Prefix, c.prefixLen)
	var words []string
	for i := 0; i < n; i++ {
		choices := c.chain[p.String()]
		if len(choices) == 0 {
			break
		}
		next := choices[rand.Intn(len(choices))]
		words = append(words, next)
		p.Shift(next)
	}
	return strings.Join(words, " ")
}

func main() {
	// 注册命令行标志
	numWords := flag.Int("words", 100, "maximum number of words to print")
	prefixLen := flag.Int("prefix", 2, "prefix length in words")

	flag.Parse()                     // Parse command-line flags.
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator.

	c := NewChain(*prefixLen)     // Initialize a new Chain.
	c.Build(os.Stdin)             // Build chains from standard input.
	text := c.Generate(*numWords) // Generate text.
	fmt.Println(text)             // Write text to standard output.
}
