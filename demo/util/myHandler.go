package util

import (
	"log"

	"github.com/siddontang/go-mysql/client"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/server"
)

// MyHandler MySQL proxy
type MyHandler struct {
	server.EmptyHandler
	conn *client.Conn
}

// NewMyhandler 连接数据库
func NewMyhandler() MyHandler {
	conn, err := client.Connect("127.0.0.1:3306", "root", "linyifan", "fyouku")
	if err != nil {
		log.Fatal(err)
	}
	return MyHandler{conn: conn}
}

// HandleQuery 解析SQL语句
func (h MyHandler) HandleQuery(query string) (*mysql.Result, error) {
	result, err := h.conn.Execute(query)
	if err != nil {
		return nil, err
	}
	return result, nil
}
