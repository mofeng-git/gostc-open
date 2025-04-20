package utils

import (
	"fmt"
	"net"
)

func IsUse(bind string, port int) error {
	address := fmt.Sprintf("%s:%d", bind, port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err // 端口被占用
	}
	defer listener.Close()
	packetConn, err := net.ListenPacket("udp", address)
	if err != nil {
		return err // 端口被占用
	}
	defer packetConn.Close()
	return nil
}
