package main

import (
	"fmt"
	"mySQLproxy/util"

	"github.com/xwb1989/sqlparser"
)

func main() {
	// sql := "select * from user a where id = 11 and name = '王五' order by id desc limit 0,10" // 命中，只会产生一条sql语句
	// sql := "select * from user a where name = '王五' order by id desc limit 0,10" // 没有命中，会产生多条sql语句（有几个分表就会有几条）

	// sql := "select * from user a where id > 1100 and id < 50 and name = '王五' order by id desc limit 0,10"
	sql := "select * from user"
	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		fmt.Println("err", err)
	}
	switch stmt := stmt.(type) {
	case *sqlparser.Select:
		// sqls := util.AliasedTableSQLTest(stmt.SelectExprs, stmt.From, stmt.Where)
		// for _, sql := range sqls {
		// 	fmt.Println(sql)
		// }
		fmt.Println(stmt)
		sqls2 := util.AliasedTableSQL(stmt)
		fmt.Println(sqls2)
		for _, sql := range sqls2 {
			fmt.Println(sql)
		}
	case *sqlparser.Insert:
		fmt.Println("Insert")
	default:
		fmt.Println("default")
	}

	// 取集合test
	// s1 := mapset.NewSetFromSlice([]interface{}{"user1", "user2"})
	// s2 := mapset.NewSetFromSlice([]interface{}{"user1"})

	// fmt.Println(s1.Intersect(s2).ToSlice())
}
