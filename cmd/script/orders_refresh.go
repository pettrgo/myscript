package script

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/spf13/cobra"
	"gitlab.xiaoduoai.com/golib/xd_sdk/logger"
	"golang.org/x/sync/errgroup"
	"io"
	"myscript/esmodel"
	"myscript/esmodel/trade_orders"
	"myscript/model"
	"time"
)

var OrderRefreshCmd = &cobra.Command{
	Use:   "order_refresh",
	Short: "刷新trade_orders订单数据",
	Long:  "刷新trade_orders订单数据",
	Run:   process,
}

type tradeOrderHandler struct {
	OrderCh chan *model.EsOrder
	//ParallelCh        chan *model.EsOrder
	TradeOrdersClient *trade_orders.TradeOrdersEsModel
}

func process(command *cobra.Command, args []string) {
	fmt.Println("start order refresh")
	g, ctx := errgroup.WithContext(context.Background())
	handler := tradeOrderHandler{
		OrderCh: make(chan *model.EsOrder),
		//ParallelCh:        make(chan *model.EsOrder),
		TradeOrdersClient: trade_orders.Get(),
	}
	esmodel.Init()

	g.Go(func() error {
		return handler.search(ctx)
	})
	//g.Go(func() error {
	//	return handler.parallelTestWorker(ctx)
	//})
	for i := 0; i < 50; i++ {
		g.Go(func() error {
			return handler.consume(ctx)
		})
	}
	g.Wait()
	fmt.Println("end order refresh")
}

func (h *tradeOrderHandler) search(ctx context.Context) error {
	defer func() {
		close(h.OrderCh)
		//close(h.ParallelCh)
	}()
	// 2023-02-24之前的订单
	query := elastic.NewRangeQuery("UpdatedAt").Lte(1677168000)
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
			order := model.Order{}
			if err := json.Unmarshal(hit.Source, &order); err != nil {
				continue
			}
			esOrder := &model.EsOrder{
				OrderInfo:   &order,
				SeqNo:       hit.SeqNo,
				PrimaryTerm: hit.PrimaryTerm,
			}
			select {
			case h.OrderCh <- esOrder:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		//for _, o := range esOrders {
		//	select {
		//	case h.OrderCh <- o:
		//	case <-ctx.Done():
		//		return ctx.Err()
		//	}
		//	//fmt.Println("put order into order chan")
		//	//h.OrderCh <- o
		//	//fmt.Println("put order into parallel chan")
		//	//h.ParallelCh <- o
		//}
	}
}

//func (h *tradeOrderHandler) search(ctx context.Context) error {
//	defer close(h.OrderCh)
//	//cnt := 0
//	//orderIDMap := make(map[string]interface{})
//	if err := h.TradeOrdersClient.DoWithScrollSession(ctx, func(scrollSession *es_official.ScrollSession) error {
//		for {
//			orders := make([]*model.Order, 0)
//			extra := map[string]interface{}{
//				"size": 1000,
//			}
//			_, scrollId, err := scrollSession.Scroll(ctx, "", nil, elastic.NewBoolQuery(), extra, &orders)
//			if err != nil {
//				return err
//			}
//			if scrollId == "" {
//				return nil
//			}
//			for _, o := range orders {
//				//orderIDMap[o.OrderID] = struct{}{}
//				//cnt++
//				h.OrderCh <- o
//			}
//		}
//	}); err != nil {
//		logger.WithError(err).Error(ctx, "scroll orders failed")
//		return err
//	}
//	//logger.Infof(ctx, "after scroll, scroll count: %v, count2: %v", cnt, len(orderIDMap))
//	return nil
//}

func (h *tradeOrderHandler) consume(ctx context.Context) error {
	orders := make([]*model.EsOrder, 0, 100)
	for o := range h.OrderCh {
		if !o.OrderInfo.NeedUpdate() {
			//logger.Infof(ctx, "skip current order, order_id: %v", o.OrderInfo.OrderID)
			continue
		}
		o.OrderInfo = o.OrderInfo.FieldClipping()
		orders = append(orders, o)
		if len(orders) == 100 {
			if err := h.TradeOrdersClient.UpdateOrdersByVersion(ctx, orders); err != nil {
				logger.WithError(err).Error(ctx, "update orders failed")
			}
			logger.Infof(ctx, "update per 100 orders, len: %v", len(orders))
			time.Sleep(10 * time.Millisecond)
			orders = make([]*model.EsOrder, 0, 100)
		}
	}
	if len(orders) > 0 {
		if err := h.TradeOrdersClient.UpdateOrdersByVersion(ctx, orders); err != nil {
			logger.WithError(err).Error(ctx, "update orders failed")
		}
		logger.Infof(ctx, "finally update orders, len:%v", len(orders))
	}
	return nil
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
