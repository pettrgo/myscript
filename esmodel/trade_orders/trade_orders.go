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

	for _, o := range orders {
		docID := fmt.Sprintf("%s_%s", o.OrderInfo.Platform, o.OrderInfo.OrderID)
		if o.SeqNo == nil || o.PrimaryTerm == nil {
			logger.Errorf(ctx, "seqNo or primaryTerm is nil, order_id: %v", o.OrderInfo.OrderID)
			continue
		}
		indexReq := elastic.NewBulkIndexRequest().Index(t.IndexName(ctx)).Id(docID).IfSeqNo(*o.SeqNo).IfPrimaryTerm(*o.PrimaryTerm).Doc(o.OrderInfo)

		bulkRequest = bulkRequest.Add(indexReq)
	}

	_, err := bulkRequest.Do(ctx)
	//for _, b := range br.Failed() {
	//	logger.Errorf(ctx, utils.UnsafeMarshal(ctx, b.Error))
	//}

	return err
}
