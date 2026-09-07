package main

import (
	"fmt"
	"os"

	"github.com/kuwaharu-git/goshare/internal/file"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("ファイルパスを指定してください")
		return
	}

	path := os.Args[1]

	info, err := file.GetInfo(path)
	if err != nil {
		fmt.Println("エラー:", err)
		return
	}

	fmt.Printf("Name: %s\n", info.Name)
	fmt.Printf("Size: %d bytes\n", info.Size)
}
