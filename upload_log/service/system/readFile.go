package system

import "fmt"

type ReadFileService struct{}

func (ReadFileService *ReadFileService) ReadFile(file string) {
	fmt.Println("readlog")
}
