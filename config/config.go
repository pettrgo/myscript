package config

import (
	"encoding/json"
	"fmt"
	"gitlab.xiaoduoai.com/golib/xd_sdk/pubsub/pulsar"
	"os"
)

var config = &Config{}

type Config struct {
	ServiceName   string            `json:"service_name"`
	OrderPub      pulsar.PubOptions `json:"order_pub"`
	EsConfigs     []EsConfig        `json:"es_configs"`
	Logger        Logger            `json:"logger"`
	RemoteService RemoteService     `json:"remote_service"`
	Mongos        MongoConfigs      `json:"mongos"`
}

type EsConfig struct {
	IndexName      string   `json:"index_name"`
	Addrs          []string `json:"addrs" `
	RequestTimeout int      `json:"timeout"`
	AutoInit       bool     `json:"auto_init"`
}

type MongoConfigs map[string]MongoConfig

type MongoConfig struct {
	Addrs          []string `json:"addrs" mapstructure:"addrs" example:"127.0.0.1:27017"`
	Source         string   `json:"source" mapstructure:"source" example:""`
	ReplicaSetName string   `json:"replica_set_name" mapstructure:"replica_set_name" example:""`
	Timeout        int      `json:"timeout" mapstructure:"timeout" example:"2"`
	Username       string   `json:"username" mapstructure:"username" example:""`
	Password       string   `json:"password" mapstructure:"password" example:""`
	Mode           *int     `json:"mode,omitempty" mapstructure:"mode,omitempty" example:"3"`
	Alias          string   `json:"alias" mapstructure:"alias" example:"default"`
	AppName        string   `mapstructure:"app_name"`
}

type Logger struct {
	Level string `json:"level"`
	Path  string `json:"path"`
}

type RemoteService struct {
	SdkTbApi RemoteServiceConfig `json:"sdk_tb_api"`
}

type RemoteServiceConfig struct {
	Addr    string `json:"addr"`
	Timeout int    `json:"timeout"`
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
