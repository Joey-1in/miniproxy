package util

import (
	"mySQLproxy/conf"
	"strconv"

	"github.com/xwb1989/sqlparser"
)

// AliasedTableSQLTest 处理SQL语句
func AliasedTableSQLTest(selectFields sqlparser.SelectExprs, from sqlparser.TableExprs, where *sqlparser.Where) []string {
	// 获取配置内容
	config := conf.NewConfig()
	SQLs := make([]string, 1)
	for _, node := range from {
		// 获取表名
		tableName := node.(*sqlparser.AliasedTableExpr).Expr.(sqlparser.TableName).Name
		// 别名
		as := node.(*sqlparser.AliasedTableExpr).As
		if mtables, ok := config.Models[tableName.String()]; ok {
			for _, mtable := range mtables {
				// 组装失去了语句
				newSQL := &sqlparser.Select{}     // select
				newSQL.SelectExprs = selectFields // select xx, xx
				newTe := &sqlparser.AliasedTableExpr{As: as, Expr: sqlparser.TableName{Name: sqlparser.NewTableIdent(mtable)}}
				newSQL.From = append(newSQL.From, newTe) // select xx, xx from xxx
				newSQL.Where = where
				buf := sqlparser.NewTrackedBuffer(nil)
				newSQL.Format(buf)
				SQLs = append(SQLs, buf.String())
			}
		}
	}
	return SQLs
}

// AliasedTableSQL 处理SQL语句
// 处理条件更多的语句
func AliasedTableSQL(stmt *sqlparser.Select) []string {
	SQLs := make([]string, 0)
	node := ParseMutiWhere(stmt.Where)
	for _, te := range stmt.From {
		// 获取表名
		tableName := te.(*sqlparser.AliasedTableExpr).Expr.(sqlparser.TableName).Name
		// 别名
		as := te.(*sqlparser.AliasedTableExpr).As
		if mtables, ok := config.Models[tableName.String()]; ok {
			for _, mtable := range mtables {
				if node != nil && len(node) > 0 && !Contains(node, mtable) {
					continue
				}
				sql := forkSQL(stmt, mtable, as)
				SQLs = append(SQLs, sql)
			}
		}
	}
	return SQLs
}

// 获取配置内容
var config = conf.NewConfig()

// 条件运算符
var operator = []interface{}{"=", "<", "<=", ">", ">="}

// GetString 辅助函数
func GetString(expr sqlparser.Expr) string {
	buf := sqlparser.NewTrackedBuffer(nil)
	expr.Format(buf)
	return buf.String()
}

// GetInt 辅助函数
func GetInt(expr sqlparser.Expr, defValue int) int {
	str := GetString(expr)
	istr, err := strconv.Atoi(str)
	if err != nil {
		return defValue
	}
	return istr
}

// forkSQL 处理
func forkSQL(stmt *sqlparser.Select, mtable string, as sqlparser.TableIdent) string {
	newSQL := &sqlparser.Select{}
	newSQL.SelectExprs = stmt.SelectExprs
	newTe := &sqlparser.AliasedTableExpr{
		As:   as,
		Expr: sqlparser.TableName{Name: sqlparser.NewTableIdent(mtable)},
	}
	newSQL.From = append(newSQL.From, newTe)
	newSQL.Where = stmt.Where
	newSQL.OrderBy = stmt.OrderBy
	newSQL.Limit = stmt.Limit // select xx, xx from xxx where xx order by xx limit xx
	buf := sqlparser.NewTrackedBuffer(nil)
	newSQL.Format(buf)
	return buf.String()
}

// ParseWhere test
func ParseWhere(expr sqlparser.Expr) []interface{} {
	if expr == nil {
		return nil
	}
	ce := expr.(*sqlparser.ComparisonExpr)
	rule := config.Rule.(*conf.RangeRule)
	column := rule.Column
	if GetString(ce.Left) == column {
		if Contains(operator, ce.Operator) {
			node := rule.GetNode(GetInt(ce.Right, 0), ce.Operator)
			return node
		}
	}
	return nil
}

// getNode 递归处理
func getNode(expr sqlparser.Expr, isLeft bool) []interface{} {
	if andExpr, ok := expr.(*sqlparser.AndExpr); ok {
		if isLeft {
			return getNode(andExpr.Left, isLeft)
		} else {
			return getNode(andExpr.Right, isLeft)
		}
	} else if cExpr, ok := expr.(*sqlparser.ComparisonExpr); ok {
		return ParseWhere(cExpr)
	} else {
		return nil
	}
}

// ParseMutiWhere 处理多条件问题
func ParseMutiWhere(where *sqlparser.Where) []interface{} {
	if where == nil {
		return nil
	}
	exps := ParseWhereToSlice(where.Expr)
	ret := make([]interface{}, 0)
	for _, exp := range exps {
		parsedNode := ParseWhere(exp)
		if parsedNode == nil || len(parsedNode) == 0 {
			continue
		}
		if len(ret) == 0 {
			ret = parsedNode
		}
		ret = IntersectSlice(ret, parsedNode)
	}
	return ret
}

// ParseWhereToSlice 处理多条件
func ParseWhereToSlice(expr sqlparser.Expr) []sqlparser.Expr {
	exprList := make([]sqlparser.Expr, 0)
	temp := expr
	for {
		//如果是and类型 则取右边
		if andExpr, ok := temp.(*sqlparser.AndExpr); ok {
			exprList = append(exprList, andExpr.Right)
			temp = andExpr.Left
		} else {
			exprList = append(exprList, temp)
			break
		}
	}
	return exprList
}
