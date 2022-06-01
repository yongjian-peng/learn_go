package system

import "fmt"

type ClearFileService struct{}

func (ClearFileService *ClearFileService) ClearFile() {
	fmt.Print("clearFile")
}
