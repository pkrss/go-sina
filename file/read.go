package file

import (
	"io/ioutil"
)

func ReadData(fileName string) ([]byte, error) {
	return ioutil.ReadFile(fileName)
}
