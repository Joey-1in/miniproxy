package conf

import (
	"github.com/xwb1989/sqlparser"
)

// Rule interface
type Rule interface {
}

// RangeRule 分表规则
type RangeRule struct {
	Column string
	Ranges []*Range
}

// NewRangeRule RangeRule工厂模式
func NewRangeRule(column string) *RangeRule {
	return &RangeRule{Column: column, Ranges: make([]*Range, 0)}
}

// AddRange 新增RangeRule
func (r *RangeRule) AddRange(max, min int, node string) {
	r.Ranges = append(r.Ranges, &Range{Max: max, Min: min, Node: node})
}

// GetNode 根据传值判断要取那个表的数据库
func (r *RangeRule) GetNode(value int, operator string) []interface{} {
	getIndex := -1
	for index, r := range r.Ranges {
		if value >= r.Min && value <= r.Max {
			getIndex = index
			break
		}
	}
	if getIndex < 0 {
		return nil
	}
	// 如果是 =
	if operator == sqlparser.EqualStr {
		return []interface{}{r.Ranges[getIndex].Node}
	}
	// 如果是 < 和 <=
	if operator == sqlparser.LessThanStr || operator == sqlparser.LessEqualStr {
		return r.GetNodes(getIndex, true)
	}
	// 如果是 > 和 >=
	if operator == sqlparser.GreaterEqualStr || operator == sqlparser.GreaterThanStr {
		return r.GetNodes(getIndex, false)
	}
	return nil
}

// GetNodes TODO
func (r *RangeRule) GetNodes(index int, less bool) []interface{} {
	result := make([]interface{}, 0)
	if less {
		for _, r := range r.Ranges[0 : index+1] {
			result = append(result, r.Node)
		}
	} else {
		for _, r := range r.Ranges[index:] {
			result = append(result, r.Node)
		}
	}
	return result
}

// Range where
type Range struct {
	Max  int
	Min  int
	Node string
}
