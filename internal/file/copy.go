package file

import (
	"io"
	"os"
)

func Copy(srcPath string, dstPath string) error {
	// 元ファイルを開く
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	// コピー先を作成
	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	buffer := make([]byte, 1024)
	for {
		n, err := src.Read(buffer)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		_, err = dst.Write(buffer[:n])
		if err != nil {
			return err
		}
	}

	return nil
}
