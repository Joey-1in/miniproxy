package main

import (
	"fmt"

	"github.com/siddontang/go-log/log"
	"github.com/xwb1989/sqlparser"
)

func main() {
	sql := "select * from user as a"
	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		log.Fatal(err)
	}
	switch stmt := stmt.(type) {
	case *sqlparser.Select:
		// 获取要拿的字段
		buf := sqlparser.NewTrackedBuffer(nil)
		stmt.SelectExprs.Format(buf)
		fmt.Println(buf.String())
	case *sqlparser.Insert:
		fmt.Println("Insert")
	default:
		fmt.Println("default")
	}
}
