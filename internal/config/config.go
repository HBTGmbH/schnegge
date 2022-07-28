package config

type ConfigKey int

const (
	TokenID ConfigKey = iota
	TokenSecret
	Date
	Hours
	Minutes
	Order
	Training
	Comment
	Verbose
	Server
	SortByProject
	NoSplash
)

type Config interface {
	GetValue(key ConfigKey) (value string, ok bool)
	AddValue(key ConfigKey, value string)
}

type SimpleConfig struct {
	values map[ConfigKey]string
}

func NewSimpleConfig() *SimpleConfig {
	var c SimpleConfig
	c.values = make(map[ConfigKey]string)
	return &c
}

func (c *SimpleConfig) GetValue(key ConfigKey) (value string, ok bool) {
	value, ok = c.values[key]
	return value, ok
}

func (c *SimpleConfig) AddValue(key ConfigKey, value string) {
	c.values[key] = value
}

type HierarchyConfig struct {
	values     map[ConfigKey]string
	BaseConfig Config
}

func NewHierarchyConfig(cfg Config) *HierarchyConfig {
	var c HierarchyConfig
	c.values = make(map[ConfigKey]string)
	c.BaseConfig = cfg
	return &c
}

func (c *HierarchyConfig) GetValue(key ConfigKey) (value string, ok bool) {
	value, ok = c.values[key]
	if ok {
		return value, ok
	}
	value, ok = c.BaseConfig.GetValue(key)
	return value, ok
}

func (c *HierarchyConfig) AddValue(key ConfigKey, value string) {
	c.values[key] = value
}
