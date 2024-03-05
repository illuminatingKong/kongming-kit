package config

type Option interface {
	get() interface{}
}
type option struct{ value interface{} }

func (o *option) get() interface{} { return o.value }

func (c *Config) SetConfigxPath(in string) {
	if in != "" {
		ab := c.abPath(in)
		c.Logger.Info("adding path to search paths: " + ab)
		c.ConfigsDir = ConfigsDir(ab)
		c.AddConfigPath(ab)
	}
}
