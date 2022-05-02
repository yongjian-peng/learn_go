package handingerror

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func ReadFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "open failed") // 包装根因
	}
	defer f.Close()
}

// demo
func main() {
	_, err := ReadConfig()
	if err != nil {
		fmt.Printf("original error: %T %v\n", errors.Cause(err), errors.Cause(err)) // 获取根因
		fmt.Printf("stack trace:\n%+v\n", err)
	}
}

func ReadConfig() ([]byte, error) {
	home := os.Getenv("HOME")
	config, err := ReadFile(filepath.Join(home, ".settngs.xml"))

	return config, errors.WithMessage(err, "could not read config")
}

func main() {
	_, err := ReadConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Write(w io.Write, buf []byte) error {
	_, err := w.Write(buf)
	return errors.Wrap(err, "write failede")
}

func Write2(w io.Write, buf []byte) error {
	_, err := w.Write(buf)
	if err != nil {
		log.Println("unable to write:", err)
		return err
	}
	return nil
}
