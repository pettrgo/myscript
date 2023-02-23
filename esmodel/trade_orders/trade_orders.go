package trade_orders

import (
	"context"
	"fmt"
	"gitlab.xiaoduoai.com/golib/xd_sdk/es_official"
	"myscript/model"
	"sync"
)

var initHandler sync.Once

const tradeOrdersClientName = "trade_orders"

const tradeOrdersIndex = "trade_orders"

var tradeOrderEsModel *TradeOrdersEsModel

type TradeOrdersEsModel struct {
	*es_official.Base
}

func Get() *TradeOrdersEsModel {
	initHandler.Do(func() {
		tradeOrderEsModel = &TradeOrdersEsModel{
			es_official.NewBaseModelV7(tradeOrdersClientName, tradeOrdersIndex),
		}
	})
	return tradeOrderEsModel
}

func (t *TradeOrdersEsModel) Name() string {
	return "trade_orders"
}

func (t *TradeOrdersEsModel) UpdateOrders(ctx context.Context, orders []*model.Order) error {
	values := make([]interface{}, len(orders))
	for idx, order := range orders {
		values[idx] = order
	}
	if err := t.BatchUpsert(ctx, values, func(i int) string {
		order := orders[i]
		return fmt.Sprintf("%s_%s", order.Platform, order.OrderID)
	}); err != nil {
		return err
	}
	return nil
}
