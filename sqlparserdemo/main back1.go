package main

import (
	"fmt"

	"github.com/siddontang/go-log/log"
	"github.com/xwb1989/sqlparser"
)

func mainback1() {
	sql := "select * from user as a"
	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		log.Fatal(err)
	}
	switch stmt := stmt.(type) {
	case *sqlparser.Select:
		// 获取表明
		for _, node := range stmt.From {
			fmt.Printf("%T\n", node)
			getTable := node.(*sqlparser.AliasedTableExpr)
			// 打印别名
			fmt.Println(getTable.As.String())
			fmt.Printf("%T\n", getTable.Expr)
			// 打印表名
			fmt.Println(getTable.Expr.(sqlparser.TableName).Name)
		}

	case *sqlparser.Insert:
		fmt.Println("Insert")
	default:
		fmt.Println("default")
	}
}
