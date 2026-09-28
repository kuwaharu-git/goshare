package transfer

import (
	"fmt"
	"net"
)

func SendMessage(address string, message string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Write([]byte(message))
	if err != nil {
		return err
	}
	return nil
}

func ReceiveMessage(address string) error {
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
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return err
	}
	message := string(buffer[:n])
	fmt.Println((message))
	return nil

}
