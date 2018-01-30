package file

import (
	"os"
)

func WriteData(fileName string, data []byte) (int, error) {
	f, e := os.Create(fileName) //创建文件
	if e != nil {
		return 0, e
	}
	defer f.Close()

	return f.Write(data)
}
