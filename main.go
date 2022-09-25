package main

import (
	"context"
	"gitlab.xiaoduoai.com/golib/xd_sdk/logger"
	"gitlab.xiaoduoai.com/golib/xd_sdk/pubsub"
	"myscript/config"
	"myscript/queue/order_pulsar"
)

var sdkRsp = "{\"buyer_nick\":\"t_1479649095865_099\",\"buyer_rate\":false,\"created\":\"2022-07-26 14:21:54\",\"modified\":\"2022-07-26 14:21:54\",\"new_presell\":false,\"num\":1,\"orders\":{\"order\":[{\"adjust_fee\":\"0.00\",\"buyer_rate\":false,\"cid\":125436003,\"discount_fee\":\"0.00\",\"is_daixiao\":false,\"is_oversold\":false,\"num\":1,\"num_iid\":636371344190,\"oid\":\"2783398428720436522\",\"order_from\":\"WAP,WAP\",\"outer_iid\":\"bianma\",\"payment\":\"0.01\",\"pic_path\":\"https://img.alicdn.com/bao/uploaded/i4/2486708185/O1CN018T503U2AKmSJ6k6Ww_!!2486708185.jpg\",\"price\":\"0.01\",\"refund_status\":\"NO_REFUND\",\"seller_rate\":false,\"seller_type\":\"C\",\"sku_id\":\"4736663246902\",\"sku_properties_name\":\"颜色分类:翠绿色;套餐类型:标准套餐;成色:官翻新机\",\"snapshot_url\":\"v:2781353054620436522_1\",\"status\":\"WAIT_BUYER_PAY\",\"title\":\"测试-多sku分级\"}]},\"step_trade_status\":\"FRONT_NOPAID_FINAL_NOPAID\",\"step_paid_fee\":\"0.01\",\"payment\":\"0.01\",\"post_fee\":\"0.00\",\"receiver_city\":\"成都市\",\"receiver_state\":\"四川省\",\"receiver_town\":\"中和街道\",\"receiver_zip\":\"000000\",\"seller_nick\":\"木月瑞希尔\",\"status\":\"WAIT_BUYER_PAY\",\"tid\":\"2782667268547436522\",\"tid_str\":\"2782667268547436522\",\"total_fee\":\"0.01\",\"trade_from\":\"WAP,WAP\",\"type\":\"fixed\"}"

type TmcOrder struct {
	Status string `json:"status"`
	Trade  Trade  `json:"trade"`
}

type Trade struct {
	BuyerNick    string `json:"buyer_nick"`
	BuyerOpenUid string `json:"buyer_open_uid"`
	EndTime      int64  `json:"end_time"`
	Iid          int64  `json:"iid"`
	Modified     string `json:"modified"`
	Oid          int64  `json:"oid"`
	Payment      string `json:"payment"`
	PostFee      string `json:"post_fee"`
	SellerNick   string `json:"seller_nick"`
	Status       string `json:"status"`
	StoreCode    string `json:"store_code"`
	Tid          int64  `json:"tid"`
	Type         string `json:"type"`
}

func main() {
	ctx := context.Background()
	//ctx = metadata.AddTags(ctx, map[string]string{
	//	"flow_color":       "test",
	//	"GetTradeFullInfo": sdkRsp,
	//})
	defer ResourceClose(ctx)
	initFunction(ctx)
	data := generateTmcOrder(ctx)
	sendTestOrderMsg(ctx, data, 1)
}

func generateTmcOrder(ctx context.Context) []byte {
	//tmcOrder := TmcOrder{
	//	Status: "TradeBuyerStepPay",
	//	Trade: Trade{
	//		BuyerNick:    "t_1479649095865_099",
	//		BuyerOpenUid: "AAH5ZosOAKdpZ1hrVDhur_Tv",
	//		EndTime:      1658816514000,
	//		Iid:          636371344190,
	//		Modified:     "2022-07-26T14:21:54.139642599+08:00",
	//		Oid:          2783398428720436522,
	//		Payment:      "0.01",
	//		PostFee:      "0.00",
	//		SellerNick:   "木月瑞希尔",
	//		Status:       "TRADE_NO_CREATE_PAY",
	//		Tid:          2783398428720436522,
	//		Type:         "guarantee_trade",
	//	},
	//}
	//data, err := json.Marshal(tmcOrder)
	//if err != nil {
	//	logger.Fatalf(ctx, "marshal tmc order failed, err: %v", err)
	//}
	jsonStr := "{\"status\":\"create\",\"trade\":{\"buyer_nick\":\"t_1479649095865_099\",\"buyer_open_uid\":\"AAH5ZosOAKdpZ1hrVDhur_Tv\",\"end_time\":1660813549000,\"iid\":636371344190,\"modified\":\"2022-08-18T17:05:49.69876505+08:00\",\"oid\":2828425394223436522,\"payment\":\"0.01\",\"post_fee\":\"0.00\",\"seller_nick\":\"木月瑞希尔\",\"status\":\"TRADE_NO_CREATE_PAY\",\"tid\":2828425394223436522,\"type\":\"guarantee_trade\"}}"
	return []byte(jsonStr)
}

func initFunction(ctx context.Context) {
	config.Init(ctx)
	order_pulsar.InitPulsar(ctx)
}

func sendTestOrderMsg(ctx context.Context, data []byte, round int) {
	logger.Infoln(ctx, "start to push msg to pulsar")
	if err := order_pulsar.GetOrderPulsar().Send(ctx, pubsub.ProducerMessage{
		Payload: data,
	}); err != nil {
		logger.Errorf(ctx, "push msg to pulsar failed, err: %v \n", err)
		return
	}
	logger.Infoln(ctx, "push msg to pulsar success")
}

type metadataCtxKey struct{}

type Metadata map[string]string

func WithMetadata(ctx context.Context, metadata Metadata) context.Context {
	return context.WithValue(ctx, metadataCtxKey{}, metadata)
}

func ResourceClose(ctx context.Context) {
	order_pulsar.Close(ctx)
}
