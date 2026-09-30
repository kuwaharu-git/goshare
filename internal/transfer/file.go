package transfer

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"

	"github.com/kuwaharu-git/goshare/internal/progress"
)

func SendFile(address string, path string) error {
	src, err := os.Open(path)
	if err != nil {
		return err
	}
	defer src.Close()
	fileInfo, err := src.Stat()
	if err != nil {
		return err
	}
	filename := filepath.Base(path)
	filenameBytes := []byte(filename)
	filenameLength := uint16(len(filenameBytes))
	size := fileInfo.Size()

	progressReader := progress.NewReader(src, size)

	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	defer conn.Close()
	if err := binary.Write(conn, binary.BigEndian, filenameLength); err != nil {
		return err
	}

	_, err = conn.Write(filenameBytes)
	if err != nil {
		return err
	}
	_, err = io.Copy(conn, progressReader)
	if err != nil {
		return err
	}

	return nil
}

func handleConnection(conn net.Conn, dir string) error {
	defer conn.Close()
	var filenameLength uint16
	binary.Read(
		conn,
		binary.BigEndian,
		&filenameLength,
	)
	filenameBuffer := make([]byte, filenameLength)
	_, err := io.ReadFull(conn, filenameBuffer)
	filename := string(filenameBuffer)
	dstPath := filepath.Join(dir, filename)
	dst, err := os.Create(dstPath)
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
