package order_pulsar

import (
	"context"
	"gitlab.xiaoduoai.com/golib/xd_sdk/logger"
	"gitlab.xiaoduoai.com/golib/xd_sdk/pubsub"
	"gitlab.xiaoduoai.com/golib/xd_sdk/pubsub/pulsar"
	"myscript/config"
)

var orderPulsar pubsub.Publisher

func InitPulsar(ctx context.Context) {
	var err error
	if orderPulsar, err = pulsar.NewPublisher(config.GetConfig().OrderPub); err != nil {
		logger.Fatalf(ctx, "init order pulsar failed, err: %v", err)
	}
}

func Close(ctx context.Context) {
	orderPulsar.Close()
}

func GetOrderPulsar() pubsub.Publisher {
	return orderPulsar
}
