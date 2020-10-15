package main

import (
	"mySQLproxy/controller"
	"net"

	_ "github.com/siddontang/go-mysql/driver"
	"github.com/siddontang/go-mysql/server"
)

func main() {
	l, _ := net.Listen("tcp", "0.0.0.0:3309")

	for {
		c, _ := l.Accept()
		go func() {
			conn, _ := server.NewConn(c, "root", "123", controller.NewMyhandler())
			for {
				conn.HandleCommand()
			}
		}()
	}
}
