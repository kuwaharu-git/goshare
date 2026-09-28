package main

import (
	"fmt"
	"os"

	"github.com/kuwaharu-git/goshare/internal/transfer"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("モードを選択してください(send or receive)")
		return
	}

	mode := os.Args[1]

	if mode == "send" {
		if len(os.Args) < 4 {
			fmt.Println("送信先アドレスと送るファイルのパスを指定してください。")
			return
		}
		address := os.Args[2]
		path := os.Args[3]
		err := transfer.SendFile(address, path)
		if err != nil {
			fmt.Println(err)
		}
	}

	if mode == "receive" {
		if len(os.Args) < 4 {
			fmt.Println("送信元アドレスと保存先ファイルパスを指定してください。")
			return
		}
		address := os.Args[2]
		path := os.Args[3]
		err := transfer.ReceiveFile(address, path)
		if err != nil {
			fmt.Println(err)
		}
	}
}
