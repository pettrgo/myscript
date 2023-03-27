package script

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/spf13/cobra"
	"gitlab.xiaoduoai.com/golib/xd_sdk/logger"
	"io"
	"myscript/config"
	"myscript/esmodel"
	"myscript/esmodel/trade_orders"
	"myscript/model"
	"myscript/storage/mongo"
	"myscript/utils"
	"os"
	"time"
)

var DayOrderShopSearchCmd = &cobra.Command{
	Use:   "day_order_shop_search",
	Short: "查询某天订单店铺数据统计",
	Long:  "查询某天订单店铺数据统计",
	PreRun: func(cmd *cobra.Command, args []string) {
		conf := config.GetConfig()
		esmodel.Init()
		mongo.Init(conf.Mongos)
	},
	Run: dayOrderShopSearchMain,
}

type dayOrderSearchHandler struct {
	OrderCh chan *model.EsOrder
	//ParallelCh        chan *model.EsOrder
	TradeOrdersClient *trade_orders.TradeOrdersEsModel
	shopMap           map[string]int64
	Lines             [][]string
}

func init() {
	DayOrderShopSearchCmd.Flags().StringVarP(&platform, "platform", "p", "", "平台")
}

var now = time.Now()

func dayOrderShopSearchMain(command *cobra.Command, args []string) {
	fmt.Println("start day order shop search")
	ctx := context.Background()
	handler := dayOrderSearchHandler{
		TradeOrdersClient: trade_orders.Get(),
		shopMap:           make(map[string]int64),
	}
	esmodel.Init()
	header := []string{"店铺id", "店铺名", "主账号", "过期时间", "当天订单数"}
	handler.Lines = append(handler.Lines, header)
	if err := handler.search(ctx); err != nil {
		logger.WithError(err).Errorf(ctx, "handler search failed")
		return
	}
	handler.shopHandler(ctx)
	handler.writeCSV(ctx)
	logger.Infof(ctx, "get shop map: %s", utils.UnsafeMarshal(handler.shopMap))
	fmt.Println("end day order shop search")
}

func (h *dayOrderSearchHandler) search(ctx context.Context) error {
	// 2023-02-24之前的订单
	rangeQuery := elastic.NewRangeQuery("CreatedAt").Gte(1678982400).Lte(1679068800)
	//rangeQuery := elastic.NewRangeQuery("CreatedAt").Gte(0).Lte(time.Now().Unix())
	termQuery := elastic.NewTermQuery("Platform", platform)
	query := elastic.NewBoolQuery().Must(rangeQuery, termQuery)
	DSL, _ := query.Source()
	fmt.Println(utils.UnsafeMarshal(DSL))
	scroll := h.TradeOrdersClient.ScrollService().SearchSource(elastic.NewSearchSource().SeqNoAndPrimaryTerm(true)).Query(query).Size(1000)
	for {
		results, err := scroll.Do(ctx)
		if err != nil {
			if err != io.EOF {
				return err
			}
			return nil
		}
		for _, hit := range results.Hits.Hits {
			//fmt.Println(hit.Source)
			order := model.Order{}
			if err := json.Unmarshal(hit.Source, &order); err != nil {
				continue
			}
			h.shopMap[order.ShopID] += 1
		}
	}
}

func (h *dayOrderSearchHandler) shopHandler(ctx context.Context) {
	count := 0
	for shopID, cnt := range h.shopMap {
		if shopID == "" {
			continue
		}
		shop, err := model.ShopModel().FindByShopID(ctx, shopID)
		if err != nil {
			logger.WithError(err).Errorf(ctx, "find shop failed, shop_id: %v", shopID)
			continue
		}
		us, err := model.UserServiceModel().FindByUserID(ctx, shop.UserID)
		if err != nil {
			logger.WithError(err).Errorf(ctx, "get user service info failed, user_id: %v", shop.UserID)
			continue
		}
		expireTime := us.GetExpireTime()
		if expireTime.After(now) {
			continue
		}
		h.Lines = append(h.Lines, []string{shopID, shop.PlatShopName, shop.PlatUserID, expireTime.Format("2006-01-02 15:04:05"), fmt.Sprintf("%d'", cnt)})
		if count%100 == 0 {
			time.Sleep(100 * time.Millisecond)
		}
		count++
	}
}

//func (h *tradeOrderHandler) parallelTestWorker(ctx context.Context) error {
//	rand.Seed(time.Now().UnixNano())
//	for o := range h.ParallelCh {
//		// 50%概率
//		if rand.Int()%100 < 5 {
//			testUpdateOrder(ctx, o)
//		}
//	}
//	return nil
//}

func (h *dayOrderSearchHandler) writeCSV(ctx context.Context) error {
	file, err := os.Create("./day_order.csv")
	if err != nil {
		return err
	}
	defer file.Close()
	file.WriteString("\xEF\xBB\xBF")
	writer := csv.NewWriter(file)
	defer writer.Flush()
	for _, line := range h.Lines {
		if err := writer.Write(line); err != nil {
			logger.WithError(err).Errorf(ctx, "write csv failed, line: %v", line)
		}
	}
	return nil
}
