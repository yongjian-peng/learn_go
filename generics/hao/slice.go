package hao

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"unsafe"
)

type slice struct {
	array unsafe.Pointer // 指向存放数据的数组指针
	len   int            // 长度有多大
	cap   int            // 容量有多大
}

// type Person struct {
// 	Name   string // 指向存放数据的数组指针
// 	Sexual string // 长度有多大
// 	Age    int    // 容量有多大
// }

func RunSlice() {
	fmt.Println("bar")
	a := make([]int, 32)
	a[5] = 5
	b := a[1:16]
	a = append(a, 1)
	a[2] = 42
	fmt.Println(a)
	fmt.Println(b)
}

func RunSlice2() {
	path := []byte("AAAA/BBBBBBBBB")
	sepIndex := bytes.IndexByte(path, '/')

	// dir1 := path[:sepIndex]
	dir1 := path[:sepIndex:sepIndex]
	dir2 := path[sepIndex+1:]
	var foo []int
	foo = make([]int, 5)
	foo[3] = 42
	foo[4] = 100

	fmt.Println("dir1 =>", string(dir1))
	fmt.Println("dir2 =>", string(dir2))

	dir1 = append(dir1, "suffix"...)

	fmt.Println("dir1 =>", string(dir1)) // dir1 => AAAAsuffix
	fmt.Println("dir2 =>", string(dir2)) // dir2 => uffixBBBB
}

func RunSlice3() {
	// var data []int{}
	// v1 := data{}
	// v2 := data{}

	// fmt.Println("v1 == v2:", reflect.DeepEqual(v1, v2))

	m1 := map[string]string{"one": "a", "two": "b"}
	m2 := map[string]string{"two": "b", "one": "a"}
	fmt.Println("m1 == m2:", reflect.DeepEqual(m1, m2))

	s1 := []int{1, 2, 3}
	s2 := []int{1, 2, 3}
	fmt.Println("s1 == s2:", reflect.DeepEqual(s1, s2))
}

// func PrintPerson(p *Person) {
// 	fmt.Printf("Name=%s, Sexual=%s, Age=%d\n",
// 		p.Name, p.Sexual, p.Age)
// }
// func (p *Person) Print() {
// 	fmt.Printf("Name=%s, Sexual=%s, Age=%d\n",
// 		p.Name, p.Sexual, p.Age)
// }

// func RunSlice4() {
// 	var p = Person{
// 		Name:   "Hao Chen",
// 		Sexual: "Male",
// 		Age:    44,
// 	}
// 	// PrintPerson(&p)
// 	p.Print()
// }

type WithName struct {
	Name string
}

type Country struct {
	Name string
}
type City struct {
	Name string
}
type Stringable interface {
	ToString() string
}
type Printable interface {
	PrintStr()
}

func (c Country) ToString() string {
	return "Country = " + c.Name
}
func (c City) ToString() string {
	return "City = " + c.Name
}

func PrintStr(p Stringable) {
	fmt.Println(p.ToString())
}

func RunSlice5() {
	// c1 := Country{"China"}
	// c2 := City{"Beijing"}
	// c1 := Country{WithName{"Chain"}}
	// c2 := City{WithName{"Beijing"}}
	c1 := Country{"USA"}
	c2 := City{"Los Angeles"}
	PrintStr(c1)
	PrintStr(c2)
}

type Shape interface {
	Sides() int
	Area() int
}
type Square struct {
	len int
}
type Square2 struct {
	len2 int
}

func (s *Square) Sides() int {
	return 4
}
func (n *Square2) Area() int {
	return 6
}
func RunSlice6() {
	// var _ Shape = (*Square)(*Square2)
	// s := Square{len: 5}
	// n := Square2{len2: 6}
	// fmt.Printf("%d%d\n", s.Sides(), n.Area())
}

func Close(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func RunSlice7() {
	r, err := os.Open("a")
	if err != nil {
		log.Fatalf("error opening 'a'\n")
	}

	defer Close(r) // 使用 defer 关键字在函数退出时关闭文件
	r, err = os.Open("b")
	if err != nil {
		log.Fatalf("error opening 'b'\n")
	}
	defer Close(r) // 使用 defer 关键字在函数退出时关闭文件
}

type Point struct {
	Longitude     string
	Latitrude     string
	Distance      string
	ElevationGain string
	ElevationLoss string
}

type Reader struct {
	r   io.Reader
	err error
}

func (r *Reader) read(data interface{}) {
	if r.err == nil {
		r.err = binary.Read(r.r, binary.BigEndian, data)
	}
}

// func parse(r io.Reader) (*Point, error) {
func parse(input io.Reader) (*Point, error) {
	// var p Point
	// var err error
	// if err := binary.Read(r, binary.BigEndian, &p.Longitude); err != nil {
	// 	return nil, err
	// }
	// if err := binary.Read(r, binary.BigEndian, &p.Latitrude); err != nil {
	// 	return nil, err
	// }
	// if err := binary.Read(r, binary.BigEndian, &p.Distance); err != nil {
	// 	return nil, err
	// }
	// if err := binary.Read(r, binary.BigEndian, &p.ElevationGain); err != nil {
	// 	return nil, err
	// }
	// if err := binary.Read(r, binary.BigEndian, &p.ElevationLoss); err != nil {
	// 	return nil, err
	// }
	// return &p, err

	// var p Point
	// var err error
	// read := func(data interface{}) {
	// 	if err != nil {
	// 		return
	// 	}
	// 	err = binary.Read(r, binary.BigEndian, data)
	// }
	// read(&p.Longitude)
	// read(&p.Latitrude)
	// read(&p.Distance)
	// read(&p.ElevationGain)
	// read(&p.ElevationLoss)

	// if err != nil {
	// 	return &p, err
	// }

	// return &p, nil

	// scanner := bufio.NewScanner(input)
	// for scanner.Scan() {
	// 	token := scanner.Text()
	// 	// process token
	// }

	// if err := scanner.Err(); err != nil {
	// 	// process the error
	// }

	var p Point
	r := Reader{r: input}

	r.read(&p.Longitude)
	r.read(&p.Latitrude)
	r.read(&p.Distance)
	r.read(&p.ElevationGain)
	r.read(&p.ElevationLoss)

	if r.err != nil {
		return nil, r.err
	}
	return &p, nil
}

// 长度不够，少一个Weight
var b = []byte{0x48, 0x61, 0x6f, 0x20, 0x43, 0x68, 0x65, 0x6e, 0x00, 0x00, 0x2c, 0x2c}
var r = bytes.NewReader(b)

type Person struct {
	Name   [10]byte
	Age    uint8
	Weight uint8
	err    error
}

func (p *Person) read(data interface{}) {
	if p.err == nil {
		p.err = binary.Read(r, binary.BigEndian, data)
	}
}

func (p *Person) ReadName() *Person {
	p.read(&p.Name)
	return p
}

func (p *Person) ReadAge() *Person {
	p.read(&p.Age)
	return p
}

func (p *Person) ReadWeight() *Person {
	p.read(&p.Weight)
	return p
}

func (p *Person) Print() *Person {
	if p.err == nil {
		fmt.Printf("Name=%s, Age=%d, Weight=%d\n", p.Name, p.Age, p.Weight)
	}
	return p
}

func RunSlice8() {
	// var data [10]byte
	// data[0] = 'A'
	// data[1] = 'E'
	p := Person{}
	p.ReadName().ReadAge().ReadWeight().Print()
	fmt.Println(p.err) // EOF 错误
}
