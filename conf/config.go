package conf

// Config struct
type Config struct {
	Models map[string][]string
	Rule   Rule
}

// NewConfig 工厂模式
func NewConfig() Config {
	mod := make(map[string][]string)
	mod["user"] = []string{"user1", "user2"}
	return Config{Models: mod, Rule: UseRangeRule()}
}

// UseRangeRule 范围分表
func UseRangeRule() *RangeRule {
	rangerule := NewRangeRule("id")
	// 0-1000取user1表
	rangerule.AddRange(1000, 0, "user1")
	// 1000-3000取user2表
	rangerule.AddRange(3000, 1000, "user2")
	return rangerule
}
