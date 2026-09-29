package transfer

import (
	"fmt"
	"io"
	"net"
	"os"
)

func SendFile(address string, path string) error {
	src, err := os.Open(path)
	if err != nil {
		return err
	}
	defer src.Close()

	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = io.Copy(conn, src)
	if err != nil {
		return err
	}

	return nil
}

func handleConnection(conn net.Conn, dir string) error {
	defer conn.Close()
	dst, err := os.CreateTemp(dir, "received-*")
	if err != nil {
		return err
	}
	defer dst.Close()
	_, err = io.Copy(dst, conn)
	if err != nil {
		return err
	}
	return nil
}

func ReceiveFile(address string, dir string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go func(conn net.Conn, dir string) {
			err := handleConnection(conn, dir)
			if err != nil {
				fmt.Println("receive error:", err)
			}
		}(conn, dir)

	}

}
