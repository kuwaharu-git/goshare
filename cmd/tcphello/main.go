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
			fmt.Println("送信先アドレスと送る文字列を指定してください。")
			return
		}
		address := os.Args[2]
		message := os.Args[3]
		err := transfer.SendMessage(address, message)
		if err != nil {
			fmt.Println(err)
		}
	}

	if mode == "receive" {
		if len(os.Args) < 3 {
			fmt.Println("送信元アドレス指定してください。")
			return
		}
		address := os.Args[2]
		err := transfer.ReceiveMessage(address)
		if err != nil {
			fmt.Println(err)
		}
	}
}
