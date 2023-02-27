package script

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/spf13/cobra"
	"gitlab.xiaoduoai.com/golib/xd_sdk/es_official"
	"gitlab.xiaoduoai.com/golib/xd_sdk/logger"
	"golang.org/x/sync/errgroup"
	"myscript/esmodel/trade_orders"
	"myscript/model"
)

var OrderRefreshCmd = &cobra.Command{
	Use:   "order_refresh",
	Short: "刷新trade_orders订单数据",
	Long:  "刷新trade_orders订单数据",
	Run:   process,
}

type tradeOrderHandler struct {
	OrderCh           chan *model.Order
	TradeOrdersClient *trade_orders.TradeOrdersEsModel
}

func process(command *cobra.Command, args []string) {
	fmt.Println("start order refresh")
	g, ctx := errgroup.WithContext(context.Background())
	handler := tradeOrderHandler{
		OrderCh:           make(chan *model.Order),
		TradeOrdersClient: trade_orders.Get(),
	}

	g.Go(func() error {
		return handler.search(ctx)
	})
	for i := 0; i < 50; i++ {
		g.Go(func() error {
			return handler.consume(ctx)
		})
	}
	g.Wait()
	fmt.Println("end order refresh")
}

//func (h *tradeOrderHandler) search(ctx context.Context) error {
//	defer close(h.Hits)
//	scroll := h.TradeOrdersClient.ScrollService().Size(1000)
//	cnt := 0
//	count := 0
//	for {
//		results, err := scroll.Do(ctx)
//		if err != nil {
//			if err != io.EOF {
//				return err
//			}
//			logger.Infof(ctx, "scroll orders finished, count: %v", count)
//			return nil
//		}
//		count += len(results.Hits.Hits)
//		for _, hit := range results.Hits.Hits {
//			select {
//			case h.Hits <- hit.Source:
//			case <-ctx.Done():
//				return ctx.Err()
//			}
//		}
//		cnt++
//		if cnt%1000 == 0 {
//			time.Sleep(time.Millisecond * 100)
//		}
//	}
//}

func (h *tradeOrderHandler) search(ctx context.Context) error {
	defer close(h.OrderCh)
	//cnt := 0
	//orderIDMap := make(map[string]interface{})
	if err := h.TradeOrdersClient.DoWithScrollSession(ctx, func(scrollSession *es_official.ScrollSession) error {
		for {
			orders := make([]*model.Order, 0)
			extra := map[string]interface{}{
				"size": 1000,
			}
			_, scrollId, err := scrollSession.Scroll(ctx, "", nil, elastic.NewBoolQuery(), extra, &orders)
			if err != nil {
				return err
			}
			if scrollId == "" {
				return nil
			}
			for _, o := range orders {
				//orderIDMap[o.OrderID] = struct{}{}
				//cnt++
				h.OrderCh <- o
			}
		}
	}); err != nil {
		logger.WithError(err).Error(ctx, "scroll orders failed")
		return err
	}
	//logger.Infof(ctx, "after scroll, scroll count: %v, count2: %v", cnt, len(orderIDMap))
	return nil
}

func (h *tradeOrderHandler) consume(ctx context.Context) error {
	orders := make([]*model.Order, 0, 100)
	for order := range h.OrderCh {
		if !order.NeedUpdate() {
			//logger.Infof(ctx, "skip current order, order_id: %v", order.OrderID)
			continue
		}
		orders = append(orders, order.FieldClipping())
		if len(orders) == 100 {
			if err := h.TradeOrdersClient.UpdateOrders(ctx, orders); err != nil {
				logger.WithError(err).Error(ctx, "update orders failed")
			}
			logger.Infof(ctx, "update per 100 orders, len: %v", len(orders))
			orders = make([]*model.Order, 0, 100)
		}
	}
	if len(orders) > 0 {
		if err := h.TradeOrdersClient.UpdateOrders(ctx, orders); err != nil {
			logger.WithError(err).Error(ctx, "update orders failed")
		}
		logger.Infof(ctx, "finally update orders, len:%v", len(orders))
	}
	return nil
}
