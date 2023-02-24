package config

import (
	"encoding/json"
	"fmt"
	"gitlab.xiaoduoai.com/golib/xd_sdk/pubsub/pulsar"
	"os"
)

var config = &Config{}

type Config struct {
	ServiceName string            `json:"service_name"`
	OrderPub    pulsar.PubOptions `json:"order_pub"`
	EsConfigs   []EsConfig        `json:"es_configs"`
	Logger      Logger            `json:"logger"`
}

type EsConfig struct {
	IndexName      string   `json:"index_name"`
	Addrs          []string `json:"addrs" `
	RequestTimeout int      `json:"timeout"`
	AutoInit       bool     `json:"auto_init"`
}

type Logger struct {
	Level string `json:"level"`
	Path  string `json:"path"`
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
	fmt.Println("configFile: ", string(data))
	if err := json.Unmarshal(data, config); err != nil {
		panic(err)
	}
}

func GetConfig() *Config {
	return config
}
