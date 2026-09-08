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

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	return nil
}
