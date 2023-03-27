package script

import (
	"context"
	"encoding/json"
	"github.com/spf13/cobra"
	"gitlab.xiaoduoai.com/golib/xd_sdk/pubsub"
	"myscript/config"
	"myscript/queue"
	"time"
)

var PubTransferMsgCmd = &cobra.Command{
	Use:   "pub_transfer_msg",
	Short: "发布转接消息",
	Long:  "发布转接消息",
	PreRun: func(cmd *cobra.Command, args []string) {
		conf := config.GetConfig()
		queue.Init(conf.QueueConfig)
	},
	Run: pubTransferMsgMain,
}

// TransferResult 转接结果存储字段文档：https://xiaoduoai.feishu.cn/wiki/wikcnNNm1pGhty2dFgLtpuKfRTd

type TransferResult struct {
	ID            string `json:"id" mapstructure:"id"`
	ShopID        string `json:"shop_id" mapstructure:"shop_id"`
	Status        int    `json:"status" mapstructure:"status"`
	Type          int    `json:"type" mapstructure:"type"`
	Platform      string `json:"platform" mapstructure:"platform"`
	Source        int    `json:"source" mapstructure:"source"`
	Buyer         string `json:"buyer" mapstructure:"buyer"`
	RealBuyerNick string `json:"real_buyer_nick" mapstructure:"real_buyer_nick"`
	BuyerOneID    string `json:"buyer_one_id" mapstructure:"buyer_one_id"`
	FromSeller    string `json:"from_seller" mapstructure:"from_seller"`
	ToSeller      string `json:"to_seller" mapstructure:"to_seller"`
	Date          int64  `json:"date" mapstructure:"date"`
	FailureReason int    `json:"failure_reason" mapstructure:"failure_reason"`
	Extend        string `json:"extend,omitempty" mapstructure:"extend"`
}

func pubTransferMsgMain(command *cobra.Command, args []string) {
	ctx := context.Background()
	tr := TransferResult{
		Status:     1,
		Date:       time.Now().Unix(),
		Platform:   "tb",
		ShopID:     "test_shop_id",
		Buyer:      "one_id_test",
		BuyerOneID: "one_id_test",
		FromSeller: "seller01",
		ToSeller:   "seller02",
	}
	pub := queue.GetPubClient(ctx, "tb", "transfer_result")
	msg, err := json.Marshal(tr)
	if err != nil {
		panic(err)
	}
	if err := pub.Send(ctx, pubsub.ProducerMessage{
		Payload: msg,
	}); err != nil {
		panic(err)
	}
}
