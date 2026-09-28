package transfer

import (
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

func ReceiveFile(address string, path string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer listener.Close()
	conn, err := listener.Accept()
	if err != nil {
		return err
	}
	defer conn.Close()
	dst, err := os.Create(path)
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
