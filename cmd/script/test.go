package script

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"myscript/esmodel/trade_orders"
	"myscript/model"
)

var platform string
var orderID string
var buyerID string

func init() {
	TestCmd.PersistentFlags().StringVarP(&platform, "平台", "p", "", "")
	TestCmd.PersistentFlags().StringVarP(&orderID, "订单号", "o", "", "")
	TestCmd.PersistentFlags().StringVarP(&buyerID, "买家id", "b", "", "")
}

var TestCmd = &cobra.Command{
	Use:   "test",
	Short: "测试用指令",
	Long:  "测试用指令",
	Run:   testMain,
}

func testMain(command *cobra.Command, args []string) {
	ctx := context.Background()
	o := testGetOrder(ctx, platform, orderID)
	o.OrderInfo.RealBuyerNick = buyerID
	testUpdateOrder(ctx, o)
	testUpdateOrder(ctx, o)
}

func testGetOrder(ctx context.Context, platform string, orderID string) *model.EsOrder {
	docID := fmt.Sprintf("%s_%s", platform, orderID)
	order := model.EsOrder{}
	tradeModel := trade_orders.Get()
	getSvc := tradeModel.GetService().Index(tradeModel.IndexName(ctx)).Id(docID)

	re, err := getSvc.Do(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("get order info succ, version: %d, seqNo: %d, primaryTerm: %d \n", re.Version, re.SeqNo, re.PrimaryTerm)
	fmt.Println(string(re.Source))
	order.SeqNo = re.SeqNo
	order.PrimaryTerm = re.PrimaryTerm
	order.OrderInfo = &model.Order{}
	if err = json.Unmarshal(re.Source, order.OrderInfo); err != nil {
		panic(err)
	}
	return &order
}

func testUpdateOrder(ctx context.Context, order *model.EsOrder) {
	tradeModel := trade_orders.Get()
	docID := fmt.Sprintf("%s_%s", order.OrderInfo.Platform, order.OrderInfo.OrderID)
	if _, err := tradeModel.UpdateService().Id(docID).IfSeqNo(*order.SeqNo).IfPrimaryTerm(*order.PrimaryTerm).Doc(order.OrderInfo).Do(ctx); err != nil {
		panic(err)
	}
	fmt.Println("update success")
}
