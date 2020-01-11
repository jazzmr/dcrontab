package master

import (
	"encoding/json"
	"io/ioutil"
)

var (
	config *Config
)

type Config struct {
	ApiPort         int      `json:"apiPort"`
	ApiReadTimeout  int      `json:"apiReadTimeout"`
	ApiWriteTimeout int      `json:"apiWriteTimeout"`
	EtcdEndpoints   []string `json:"etcdEndpoints"`
	EtcdDialTimeout int      `json:"etcdDialTimeout"`
}

// 加载配置
func InitConfig(filename string) (err error) {
	var (
		content []byte
	)

	if content, err = ioutil.ReadFile(filename); err != nil {
		return
	}

	config = &Config{}

	if err = json.Unmarshal(content, config); err != nil {
		return
	}
	return
}

func GetConfig() *Config {
	return config
}
