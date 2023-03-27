package trade_orders

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"gitlab.xiaoduoai.com/golib/xd_sdk/es_official"
	"gitlab.xiaoduoai.com/golib/xd_sdk/logger"
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

func (t *TradeOrdersEsModel) UpdateOrdersByVersion(ctx context.Context, orders []*model.EsOrder) error {

	if len(orders) == 0 {
		return nil
	}

	bulkRequest := t.Client().Bulk()

	orderIDs := make([]string, 0, len(orders))
	for _, o := range orders {
		docID := fmt.Sprintf("%s_%s", o.OrderInfo.Platform, o.OrderInfo.OrderID)
		if o.SeqNo == nil || o.PrimaryTerm == nil {
			logger.Errorf(ctx, "seqNo or primaryTerm is nil, order_id: %v", o.OrderInfo.OrderID)
			continue
		}
		orderIDs = append(orderIDs, o.OrderInfo.OrderID)
		indexReq := elastic.NewBulkIndexRequest().Index(t.IndexName(ctx)).Id(docID).IfSeqNo(*o.SeqNo).IfPrimaryTerm(*o.PrimaryTerm).Doc(o.OrderInfo)
		bulkRequest = bulkRequest.Add(indexReq)
	}
	logger.Infof(ctx, "current handler order ids: %v", orderIDs)

	_, err := bulkRequest.Do(ctx)
	//for _, b := range br.Failed() {
	//	logger.WithField("update_by", "batch update").Errorf(ctx, utils.UnsafeMarshal(ctx, b.Error))
	//}

	return err
}
