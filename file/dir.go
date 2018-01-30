package file

import (
	"os"
	"strings"
)

// 文件所在目录
func FileDir(fileName string) string {
	pos := strings.LastIndexAny(fileName, "/\\")
	if pos > 0 {
		fileName = fileName[:pos]
	}

	return fileName
}

// 创建目录
func CreateDir(fileName string) error {
	return os.MkdirAll(fileName, 0666)
}
