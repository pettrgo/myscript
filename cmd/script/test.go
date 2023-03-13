package script

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/spf13/cobra"
	"gitlab.xiaoduoai.com/golib/xd_sdk/logger"
	"math/rand"
	"myscript/esmodel/trade_orders"
	"myscript/model"
	"time"
)

var platform string
var orderID string
var buyerID string

var orderTemplate = "{\"ID\":\"640aa8f81f46dc000173ee94\",\"Platform\":\"tb\",\"OrderID\":\"3249724142611744742\",\"OriginalOrderID\":\"3249724142611744742\",\"ShopID\":\"\",\"SellerID\":\"方太官方旗舰店\",\"BuyerID\":\"one_id_806744247\",\"OriBuyerID\":\"\",\"Payment\":8199,\"Address\":\"安徽省舒*县经济开发区古**路青**园西苑门面房米卡造型\",\"Province\":\"\",\"City\":\"\",\"Town\":\"\",\"Street\":\"\",\"DecryptAddress\":\"\",\"Status\":\"created\",\"OriginStatus\":\"WAIT_BUYER_PAY\",\"OrderType\":\"\",\"StepTradeStatus\":\"\",\"StepPaidFee\":0,\"PushStatus\":\"created\",\"StatusHistory\":[{\"Status\":\"created\",\"Time\":1678420216}],\"BuyerRemark\":\"\",\"SellerMemo\":\"\",\"BalanceUsed\":\"\",\"SellerDiscount\":\"\",\"OrderSellerPrice\":\"\",\"PayType\":\"\",\"FreightPrice\":\"\",\"CreatedAt\":1678420210,\"UpdatedAt\":1678420210,\"ItemIDs\":[\"567256065770\"],\"Tbext\":{\"BuyerNick\":\"one_id_806744247\",\"BuyerRemark\":\"\",\"Items\":[{\"ItemID\":\"567256065770\",\"ItemName\":\"方太EMQ5/T+TH29B变频抽吸油烟机燃气灶套餐烟机套装烟灶官方旗舰\",\"SkuID\":\"4862006477762\",\"SkuName\":\"颜色分类:【直流爆款】TH29B(天然气5.0kw*，液化气4.5kw*);燃料种类:天然气\",\"Count\":1,\"Price\":0}],\"Receiver\":{\"Name\":\"K**\",\"Phone\":\"\",\"Zip\":\"231300\",\"DecryptName\":\"\",\"DecryptPhone\":\"\",\"DecryptMobile\":\"\"},\"SellerFlag\":0,\"OriginalOrder\":\"{\\\"adjust_fee\\\":\\\"0.00\\\",\\\"alipay_no\\\":\\\"2023031022001185321423046139\\\",\\\"alipay_point\\\":0,\\\"available_confirm_fee\\\":\\\"8199.00\\\",\\\"buyer_alipay_no\\\":\\\"150********\\\",\\\"buyer_area\\\":\\\"安徽阜阳电信\\\",\\\"buyer_cod_fee\\\":\\\"0.00\\\",\\\"buyer_email\\\":\\\"\\\",\\\"buyer_nick\\\":\\\"许**\\\",\\\"buyer_obtain_point_fee\\\":0,\\\"buyer_open_uid\\\":\\\"AAG6ZosOAKdpZ1hrVDjspd_R\\\",\\\"buyer_rate\\\":false,\\\"cod_fee\\\":\\\"0.00\\\",\\\"cod_status\\\":\\\"NEW_CREATED\\\",\\\"commission_fee\\\":\\\"0.00\\\",\\\"coupon_fee\\\":0,\\\"created\\\":\\\"2023-03-10 11:50:10\\\",\\\"discount_fee\\\":\\\"0.00\\\",\\\"has_post_fee\\\":true,\\\"is_3D\\\":false,\\\"is_brand_sale\\\":false,\\\"is_daixiao\\\":false,\\\"is_force_wlb\\\":false,\\\"is_gift\\\":false,\\\"is_lgtype\\\":true,\\\"is_part_consign\\\":false,\\\"is_sh_ship\\\":true,\\\"is_wt\\\":false,\\\"modified\\\":\\\"2023-03-10 11:50:10\\\",\\\"new_presell\\\":false,\\\"nr_shop_guide_id\\\":\\\"\\\",\\\"nr_shop_guide_name\\\":\\\"\\\",\\\"num\\\":1,\\\"num_iid\\\":567256065770,\\\"oaid\\\":\\\"13n7ici5cK6d1yK74OiaepReCWVwVTO7jVaOCCFdc4F67Dy06m65Yq4hMYmgqMurUOvibZNmr7K\\\",\\\"orders\\\":{\\\"order\\\":[{\\\"adjust_fee\\\":\\\"0.00\\\",\\\"buyer_rate\\\":false,\\\"cid\\\":50018263,\\\"discount_fee\\\":\\\"1800.00\\\",\\\"is_daixiao\\\":false,\\\"is_oversold\\\":false,\\\"is_www\\\":true,\\\"num\\\":1,\\\"num_iid\\\":567256065770,\\\"oid\\\":\\\"3249724142611744742\\\",\\\"order_from\\\":\\\"WAP,WAP\\\",\\\"outer_sku_id\\\":\\\"TC017434\\\",\\\"payment\\\":\\\"8199.00\\\",\\\"price\\\":\\\"9999.00\\\",\\\"refund_status\\\":\\\"NO_REFUND\\\",\\\"seller_rate\\\":false,\\\"seller_type\\\":\\\"B\\\",\\\"sku_id\\\":\\\"4862006477762\\\",\\\"sku_properties_name\\\":\\\"颜色分类:【直流爆款】TH29B(天然气5.0kw*，液化气4.5kw*);燃料种类:天然气\\\",\\\"snapshot_url\\\":\\\"v:3249724142611744742_1\\\",\\\"status\\\":\\\"WAIT_BUYER_PAY\\\",\\\"store_code\\\":\\\"QDHEWL-0003\\\",\\\"timeout_action_time\\\":\\\"2023-03-11 11:50:10\\\",\\\"title\\\":\\\"方太EMQ5/T+TH29B变频抽吸油烟机燃气灶套餐烟机套装烟灶官方旗舰\\\",\\\"total_fee\\\":\\\"8199.00\\\"}]},\\\"payment\\\":\\\"8199.00\\\",\\\"pcc_af\\\":0,\\\"platform_subsidy_fee\\\":\\\"0.00\\\",\\\"point_fee\\\":0,\\\"post_fee\\\":\\\"0.00\\\",\\\"price\\\":\\\"9999.00\\\",\\\"real_point_fee\\\":0,\\\"received_payment\\\":\\\"0.00\\\",\\\"receiver_address\\\":\\\"舒*县经济开发区古**路青**园西苑门面房米卡造型\\\",\\\"receiver_city\\\":\\\"六安市\\\",\\\"receiver_country\\\":\\\"\\\",\\\"receiver_district\\\":\\\"舒城县\\\",\\\"receiver_mobile\\\":\\\"*******1125\\\",\\\"receiver_name\\\":\\\"K**\\\",\\\"receiver_state\\\":\\\"安徽省\\\",\\\"receiver_town\\\":\\\"舒城县经济开发区\\\",\\\"receiver_zip\\\":\\\"231300\\\",\\\"seller_alipay_no\\\":\\\"***otile@fotile.com\\\",\\\"seller_can_rate\\\":false,\\\"seller_cod_fee\\\":\\\"0.00\\\",\\\"seller_email\\\":\\\"liubta@fotile.com\\\",\\\"seller_flag\\\":0,\\\"seller_mobile\\\":\\\"13566612988\\\",\\\"seller_name\\\":\\\"宁波方太营销有限公司\\\",\\\"seller_nick\\\":\\\"方太官方旗舰店\\\",\\\"seller_phone\\\":\\\"4000-315888\\\",\\\"seller_rate\\\":false,\\\"service_type\\\":\\\"\\\",\\\"shipping_type\\\":\\\"express\\\",\\\"sid\\\":\\\"3249724142611744742\\\",\\\"snapshot_url\\\":\\\"v:3249724142611744742_1\\\",\\\"status\\\":\\\"WAIT_BUYER_PAY\\\",\\\"tid\\\":\\\"3249724142611744742\\\",\\\"tid_str\\\":\\\"3249724142611744742\\\",\\\"timeout_action_time\\\":\\\"2023-03-11 11:50:10\\\",\\\"title\\\":\\\"方太官方旗舰店\\\",\\\"total_fee\\\":\\\"9999.00\\\",\\\"trade_attr\\\":\\\"{\\\\\\\"erpHold\\\\\\\":\\\\\\\"0\\\\\\\",\\\\\\\"tmallDelivery\\\\\\\":\\\\\\\"true\\\\\\\"}\\\",\\\"trade_from\\\":\\\"WAP,WAP\\\",\\\"type\\\":\\\"fixed\\\"}\"},\"Jdext\":null,\"Pddext\":null,\"Lazy\":false,\"NotInWhitelist\":false,\"UpdateTime\":\"2023-03-10T03:50:16.380398254Z\",\"UpdateDesc\":\"trade worker insert\",\"BuyerOpenUID\":\"AAG6ZosOAKdpZ1hrVDjspd_R\",\"RealBuyerNick\":\"许**\",\"LogisticsInfo\":null}"

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
	//o := testGetOrder(ctx, "tb", "3249837903036436522")
	//testUpdateOrder(ctx, o)
	//parallelUpdateOrderTest(ctx)
	//BatchInsertNewOrders(ctx)
	//testRand()
	testQuery(ctx)
}

func parallelUpdateOrderTest(ctx context.Context) {
	orders := make([]*model.EsOrder, 0)
	o1 := testGetOrder(ctx, platform, "3249837903036436522")
	orders = append(orders, o1)
	o2 := testGetOrder(ctx, platform, "3249599834713436522")
	orders = append(orders, o2)
	testUpdateOrder(ctx, o2)
	testBatchUpdateOrder(ctx, orders)
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
	//fmt.Println(string(re.Source))
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
	if order.SeqNo == nil || order.PrimaryTerm == nil {
		fmt.Println("current seqNo or primaryTerm is nil, docID:", docID)
		return
	}
	if _, err := tradeModel.UpdateService().Id(docID).IfSeqNo(*order.SeqNo).IfPrimaryTerm(*order.PrimaryTerm).Doc(map[string]bool{
		"is_parallel_update": true,
	}).Do(ctx); err != nil {
		logger.WithField("update_by", "one_update").Errorf(ctx, "test update order failed, docID: %v, err: %v", docID, err)
	}

	//fmt.Println("update success")
}

func testBatchUpdateOrder(ctx context.Context, orders []*model.EsOrder) {
	tradeModel := trade_orders.Get()
	for idx, o := range orders {
		o.OrderInfo.Tbext.OriginalOrder = fmt.Sprintf("test-version:%d", idx)
	}
	if err := tradeModel.UpdateOrdersByVersion(ctx, orders); err != nil {
		panic(err)
	}
	fmt.Println("batch update success")
}

func BatchInsertNewOrders(ctx context.Context) {
	order := model.Order{}
	if err := json.Unmarshal([]byte(orderTemplate), &order); err != nil {
		panic(err)
	}
	tradeModel := trade_orders.Get()
	for i := 0; i < 50000; i++ {
		order.OrderID = fmt.Sprintf("test_%d", time.Now().UnixNano())
		order.BuyerID = "one_id_test_buyer"
		order.CreatedAt = time.Now().Unix()
		order.UpdatedAt = time.Now().Unix()
		order.UpdateTime = time.Now()
		docID := fmt.Sprintf("%s_%s", order.Platform, order.OrderID)
		if err := tradeModel.InsertBodyJSON(ctx, docID, "", order); err != nil {
			logger.Errorf(ctx, "error")
		}
	}
}

func testRand() {
	trueCnt := 0
	falseCnt := 0
	for i := 0; i < 50000; i++ {
		rand.Seed(time.Now().UnixNano())
		if rand.Int()%100 < 3 {
			trueCnt++
		} else {
			falseCnt++
		}

	}
	fmt.Println("trueCnt:", trueCnt)
	fmt.Println("falseCnt:", falseCnt)
}

func testQuery(ctx context.Context) {
	query := elastic.NewRangeQuery("UpdatedAt").Lte(1678675500)
	re, err := trade_orders.Get().SearchService().Query(query).Size(1000).Do(ctx)
	if err != nil {
		panic(err)
	}
	//s, _ := query.Source()
	//fmt.Println(utils.UnsafeMarshal(ctx, s))
	fmt.Println(re.TotalHits())
}
