package config

import (
	"encoding/json"
	"gitlab.xiaoduoai.com/golib/xd_sdk/pubsub/pulsar"
	"os"
)

var config = &Config{}

type Config struct {
	OrderPub  pulsar.PubOptions `json:"order_pub"`
	EsConfigs []EsConfig        `json:"es_configs"`
}

type EsConfig struct {
	IndexName      string   `mapstructure:"index_name" json:"index_name" toml:"index_name"`
	Addrs          []string `mapstructure:"addrs" json:"addrs" toml:"addrs"`
	RequestTimeout int      `mapstructure:"timeout" json:"timeout" toml:"timeout"`
	AutoInit       bool     `mapstructure:"auto_init" json:"auto_init" toml:"auto_init"`
}

func Init(configFile string) {
	//data, err := ioutil.ReadFile("/mnt/d/work/workspace/GoWork/src/myscript/config/myscript.conf")
	//if err != nil {
	//	logger.Fatalf(ctx, "read config failed, err: %v \n", err)
	//	return
	//}
	data, err := os.ReadFile(configFile)
	if err != nil {
		panic(err)
	}
	//fmt.Println("configFile: ", string(data))
	if err := json.Unmarshal(data, config); err != nil {
		panic(err)
	}
}

func GetConfig() *Config {
	return config
}
