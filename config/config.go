package config

import (
	"context"
	"encoding/json"
	"fmt"
	"gitlab.xiaoduoai.com/golib/xd_sdk/logger"
	"gitlab.xiaoduoai.com/golib/xd_sdk/pubsub/pulsar"
	"io/ioutil"
)

var config = &Config{}

type Config struct {
	OrderPub pulsar.PubOptions `json:"order_pub"`
}

func Init(ctx context.Context) {
	data, err := ioutil.ReadFile("/mnt/d/work/workspace/GoWork/src/myscript/config/myscript.conf")
	if err != nil {
		logger.Fatalf(ctx, "read config failed, err: %v \n", err)
		return
	}
	if err := json.Unmarshal(data, config); err != nil {
		fmt.Println(string(data))
		logger.Fatalf(ctx, "unmarshal config failed, err: %v \n", err)
	}
	fmt.Println(config)
}

func GetConfig() *Config {
	return config
}
