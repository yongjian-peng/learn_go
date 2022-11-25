package main

import "fmt"

type BodyMap map[string]interface{}

func (bm *BodyMap) Set(key string, value interface{}) *BodyMap {
	b := *bm
	b[key] = value
	return bm
}

func (bm BodyMap) SetA(key string, value interface{}) BodyMap {
	bm[key] = value
	return bm
}

func main() {
	bm := make(BodyMap, 0)
	bma := bm.Set("A", "VA")
	fmt.Println(fmt.Sprintf("bma:%p", bma))
	bmb := bma.Set("B", "BA")
	fmt.Println(fmt.Sprintf("bmb:%p", bmb))
	bmc := bmb.Set("C", "CA")
	fmt.Println(fmt.Sprintf("bmc:%p", bmc))
	fmt.Println(bm)
	bmm := make(BodyMap, 0)
	bmma := bmm.SetA("A", "VA")
	fmt.Println(fmt.Sprintf("bmma:%p", &bmma))
	bmmb := bmma.SetA("B", "VB")
	fmt.Println(fmt.Sprintf("bmmb:%p", &bmmb))
	bmmc := bmmb.SetA("C", "VC")
	fmt.Println(fmt.Sprintf("bmmc:%p", &bmmc))
}
