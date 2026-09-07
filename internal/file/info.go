package file

import (
	"os"
)

type Info struct {
	Name string
	Size int64
}

func GetInfo(path string) (Info, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return Info{}, err
	}
	info := Info{
		Name: fileInfo.Name(),
		Size: fileInfo.Size(),
	}
	return info, nil
}
