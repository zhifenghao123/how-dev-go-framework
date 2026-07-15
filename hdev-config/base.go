package hdevconfig

import (
	"os"
)

type FileReader struct{}

// readFile 读取文件内容
func (FileReader) readFile(file string) ([]byte, error) {
	return os.ReadFile(file)
}
