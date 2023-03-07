package config

import "github.com/spf13/viper"

type confInfo struct {
	Name string
	Type string
	Path string
}
type config struct {
	viper *viper.Viper
}

var (
	Conf   *config
	Secret *config
)

func init() {
	ci := confInfo{
		Name: "confs",
		Type: "yaml",
		Path: "config/conf",
	}
	si := confInfo{
		Name: "secrets",
		Type: "yaml",
		Path: "config/secret",
	}
	Conf = &config{getConf(ci)}
	Secret = &config{getConf(si)}
}

// viper用于解析yaml
func getConf(ci confInfo) *viper.Viper {
	v := viper.New()
	v.SetConfigName(ci.Name) // 与yaml文件名一致
	v.SetConfigType(ci.Type)
	v.AddConfigPath(ci.Path)
	v.ReadInConfig()
	return v
}

func (c *config) GetString(key string) string {
	return c.viper.GetString(key)
}
