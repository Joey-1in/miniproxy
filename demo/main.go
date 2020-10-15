package main

import (
	"mySQLproxy/demo/util"
	"net"

	_ "github.com/siddontang/go-mysql/driver"
	"github.com/siddontang/go-mysql/server"
)

func main() {
	l, _ := net.Listen("tcp", "0.0.0.0:3309")

	for {
		c, _ := l.Accept()
		go func() {
			conn, _ := server.NewConn(c, "root", "123", util.NewMyhandler())
			for {
				conn.HandleCommand()
			}
		}()
	}
}
