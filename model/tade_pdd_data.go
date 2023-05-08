package model

import (
	"context"
	"gitlab.xiaoduoai.com/golib/xd_sdk/mongoc_official"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type PddData struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	Platform   string             `bson:"platform,omitempty"`
	Tid        string             `bson:"order_id,omitempty" json:"tid,omitempty"`            //订单号
	Status     string             `bson:"status,omitempty" json:"status,omitempty"`           //消息类型
	SellerMemo string             `bson:"seller_memo,omitempty" json:"seller_memo,omitempty"` //卖家留言
	MallID     int64              `bson:"shop_id,omitempty" json:"mall_id,omitempty"`         //店铺id
	BuyerID    string             `bson:"buyer_id,omitempty" json:"buyer_id"`
	ShopName   string             `bson:"shop_name,omitempty" json:"-"`
	CreatedAt  int64              `bson:"created_at,omitempty" json:"-"`
	UpdatedAt  int64              `bson:"updated_at,omitempty" json:"-"`
	Address    string             `bson:"address,omitempty"`
	Phone      string             `bson:"phone,omitempty"`
	ReadFlag   bool               `bson:"read_flag" json:"-"` //已读标志  默认False
	UpdateTime time.Time          `bson:"update_time"`
}

type PddDataDB struct {
	*mongoc_official.Base
}

var pddDataDB = &PddDataDB{
	Base: mongoc_official.NewBaseModel("xdmp", "xdmp", "trade_pdd_data"),
}

func PddDataModel() *PddDataDB {
	return pddDataDB
}

func (p *PddDataDB) Find(ctx context.Context, query primitive.M) (*mongo.Cursor, error) {
	return p.C(ctx).Find(ctx, query)
}

func (p *PddDataDB) DeleteMany(ctx context.Context, filter primitive.M) error {
	_, err := p.C(ctx).DeleteMany(ctx, filter)
	return err
}
