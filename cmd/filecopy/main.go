package main

import (
	"fmt"
	"os"

	"github.com/kuwaharu-git/goshare/internal/file"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("コピー元と、コピー先のファイルパスを指定してください")
		return
	}

	srcPath := os.Args[1]
	dstPath := os.Args[2]

	err := file.Copy(srcPath, dstPath)
	if err != nil {
		fmt.Println("コピーに失敗しました:", err)
		return
	}

	fmt.Println("コピー完了")

}
