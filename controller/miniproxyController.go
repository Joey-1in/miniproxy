package controller

import (
	"fmt"
	"log"
	"mySQLproxy/util"

	"github.com/siddontang/go-mysql/client"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/server"
	"github.com/xwb1989/sqlparser"
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
	// result, err := h.conn.Execute(query)
	// if err != nil {
	// 	return nil, err
	// }
	// return result, nil

	// fmt.Println("query", query)
	stmt, err := sqlparser.Parse(query)
	if err != nil {
		return nil, err
	}
	switch stmt := stmt.(type) {
	case *sqlparser.Select:
		sqls := util.AliasedTableSQL(stmt)
		if len(sqls) > 0 {
			// 有拆分
			results := make([]*mysql.Result, 0)
			for index, sql := range sqls {
				fmt.Println("sql: ", index, sql)
				r, err := h.conn.Execute(sql)
				if err != nil {
					return nil, err
				}
				results = append(results, r)
			}
			if len(results) > 1 {
				for index, result := range results {
					if index == 0 {
						continue
					}
					results[0].RowDatas = append(results[0].RowDatas, result.RowDatas...)
					results[0].AffectedRows = result.AffectedRows
				}
				return results[0], nil
			} else {
				return results[0], nil
			}
		} else {
			return h.conn.Execute(query)
		}
	default:
		return h.conn.Execute(query)
	}
}
