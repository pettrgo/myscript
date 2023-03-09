package script

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.xiaoduoai.com/golib/xd_sdk/logger"
	"golang.org/x/sync/errgroup"
	"io"
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
	OrderCh           chan *model.EsOrder
	TradeOrdersClient *trade_orders.TradeOrdersEsModel
}

func process(command *cobra.Command, args []string) {
	fmt.Println("start order refresh")
	g, ctx := errgroup.WithContext(context.Background())
	handler := tradeOrderHandler{
		OrderCh:           make(chan *model.EsOrder),
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

func (h *tradeOrderHandler) search(ctx context.Context) error {
	defer close(h.OrderCh)
	scroll := h.TradeOrdersClient.ScrollService().Size(1000)
	for {
		results, err := scroll.Do(ctx)
		if err != nil {
			if err != io.EOF {
				return err
			}
			return nil
		}
		esOrders := make([]*model.EsOrder, 0, len(results.Hits.Hits))
		for _, hit := range results.Hits.Hits {
			order := model.Order{}
			if err := json.Unmarshal(hit.Source, &order); err != nil {
				continue
			}
			esOrders = append(esOrders, &model.EsOrder{
				OrderInfo:   &order,
				SeqNo:       hit.SeqNo,
				PrimaryTerm: hit.PrimaryTerm,
			})
			//select {
			//case h.OrderCh <- hit.Source:
			//case <-ctx.Done():
			//	return ctx.Err()
			//}
		}
		for _, o := range esOrders {
			select {
			case h.OrderCh <- o:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
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
			//logger.Infof(ctx, "skip current order, order_id: %v", order.OrderID)
			continue
		}
		o.OrderInfo = o.OrderInfo.FieldClipping()
		orders = append(orders, o)
		if len(orders) == 100 {
			if err := h.TradeOrdersClient.UpdateOrdersByVersion(ctx, orders); err != nil {
				logger.WithError(err).Error(ctx, "update orders failed")
			}
			logger.Infof(ctx, "update per 100 orders, len: %v", len(orders))
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
