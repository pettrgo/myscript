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
	"time"
)

var OrderRefreshCmd = &cobra.Command{
	Use:   "order_refresh",
	Short: "刷新trade_orders订单数据",
	Long:  "刷新trade_orders订单数据",
	Run:   process,
}

type tradeOrderHandler struct {
	Hits              chan json.RawMessage
	TradeOrdersClient *trade_orders.TradeOrdersEsModel
}

func process(command *cobra.Command, args []string) {
	fmt.Println("start order refresh")
	g, ctx := errgroup.WithContext(context.Background())
	handler := tradeOrderHandler{
		Hits:              make(chan json.RawMessage),
		TradeOrdersClient: trade_orders.Get(),
	}

	g.Go(func() error {
		return handler.search(ctx)
	})
	for i := 0; i < 10; i++ {
		g.Go(func() error {
			return handler.consume(ctx)
		})
	}
	g.Wait()
	fmt.Println("end order refresh")
}

func (h *tradeOrderHandler) search(ctx context.Context) error {
	defer close(h.Hits)
	scroll := h.TradeOrdersClient.ScrollService().Size(1000)
	cnt := 0
	count := 0
	for {
		results, err := scroll.Do(ctx)
		if err != nil {
			if err != io.EOF {
				return err
			}
			logger.Infof(ctx, "scroll orders finished, count: %v", count)
			return nil
		}
		count += len(results.Hits.Hits)
		for _, hit := range results.Hits.Hits {
			select {
			case h.Hits <- hit.Source:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		cnt++
		if cnt%1000 == 0 {
			time.Sleep(time.Millisecond * 100)
		}
	}
}

func (h *tradeOrderHandler) consume(ctx context.Context) error {
	orders := make([]*model.Order, 0, 100)
	for hit := range h.Hits {
		order := &model.Order{}
		if err := json.Unmarshal(hit, order); err != nil {
			logger.WithError(err).Errorf(ctx, "unmarshal order failed, order:%s", string(hit))
			continue
		}
		if !order.NeedUpdate() {
			//fmt.Println("skip update order")
			logger.Infof(ctx, "skip current order, order_id: %v", order.OrderID)
			continue
		}
		//fmt.Println(json2.UnsafeMarshalString(order))
		orders = append(orders, order.FieldClipping())
		if len(orders) == 100 {
			//fmt.Printf("upsert orders, orders len: %d, orders: %s \n", len(orders), gjson.UnsafeMarshalString(orders))
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
		logger.Infof(ctx, "finally update orders, len:", len(orders))
	}
	return nil
}
