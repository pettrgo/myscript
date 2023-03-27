package queue

import (
	"context"
	"gitlab.xiaoduoai.com/golib/xd_sdk/logger"
	"gitlab.xiaoduoai.com/golib/xd_sdk/pubsub"
	"gitlab.xiaoduoai.com/golib/xd_sdk/pubsub/pulsar"
)

const (
	TopicTransferResult = "transfer_result"
)

type PlatformSubOptions map[string]pulsar.SubOptions

type PlatformPubOptions map[string]pulsar.PubOptions

type Config struct {
	Subscribers map[string]PlatformSubOptions `mapstructure:"subscribers"`
	Publishers  map[string]PlatformPubOptions `mapstructure:"publishers"`
}

type QClient struct {
	Pub map[string]map[string]pubsub.Publisher
	Sub map[string]map[string]pubsub.Subscriber
}

var pubSubCli = QClient{
	Pub: map[string]map[string]pubsub.Publisher{},
	Sub: map[string]map[string]pubsub.Subscriber{},
}

func Init(conf Config) {
	err := initPubCli(conf)
	if err != nil {
		panic(err)
	}
	err = initSubCli(conf)
	if err != nil {
		panic(err)
	}
}

func initPubCli(conf Config) error {
	for k, value := range conf.Publishers {
		pubs := map[string]pubsub.Publisher{}
		pubSubCli.Pub[k] = pubs

		for pf, v := range value {
			cli, err := pulsar.NewPublisher(v)
			if err != nil {
				return err
			}
			pubs[pf] = cli
		}
	}
	logger.Debugf(context.Background(), "queueCli: %+v", pubSubCli)
	return nil
}

func initSubCli(conf Config) error {
	for k, value := range conf.Subscribers {
		subs := map[string]pubsub.Subscriber{}
		pubSubCli.Sub[k] = subs

		for pf, v := range value {
			cli, err := pulsar.NewSubscriber(v)
			if err != nil {
				return err
			}
			subs[pf] = cli
		}
	}
	logger.Debugf(context.Background(), "queueCli: %+v", pubSubCli)
	return nil
}

func GetPubClient(ctx context.Context, platform, key string) pubsub.Publisher {
	pubcli := pubSubCli.Pub[key][platform]
	if pubcli == nil {
		logger.Errorf(ctx, "platform:%v key:%v pubcli is nil", platform, key)
	}

	return pubcli
}

func GetSubClient(platform, key string) pubsub.Subscriber {
	return pubSubCli.Sub[key][platform]
}
