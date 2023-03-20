package model

import (
	"context"
	"github.com/globalsign/mgo/bson"
	"gitlab.xiaoduoai.com/golib/xd_sdk/mongoc_official"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type UserServiceInfo struct {
	ID               primitive.ObjectID     `bson:"_id,omitempty" json:"id,omitempty"`
	CreateTime       time.Time              `bson:"create_time,omitempty" json:"create_time,omitempty"`
	UpdateTime       time.Time              `bson:"update_time,omitempty" json:"update_time,omitempty"`
	UserID           primitive.ObjectID     `bson:"user_id,omitempty" json:"user_id,omitempty"`
	Nick             string                 `bson:"nick,omitempty" json:"nick,omitempty"`
	ServiceName      string                 `bson:"service_name,omitempty" json:"service_name,omitempty"`
	ServiceLevel     string                 `bson:"service_level,omitempty" json:"service_level,omitempty"`
	ExpireTime       time.Time              `bson:"expire_time,omitempty" json:"expire_time,omitempty"`
	SopExpireTime    time.Time              `bson:"sop_expire_time,omitempty" json:"sop_expire_time,omitempty"`
	InviteBy         string                 `bson:"invite_by,omitempty" json:"invite_by,omitempty"`
	Settings         map[string]interface{} `bson:"settings,omitempty" json:"settings,omitempty"`
	SendMsgSettings  map[string]interface{} `bson:"send_msg_settings,omitempty" json:"send_msg_settings,omitempty"`
	GuideStatus      map[string]bool        `bson:"guide_status,omitempty" json:"guide_status,omitempty"`
	LastChannel      string                 `bson:"last_channel,omitempty" json:"last_channel,omitempty"`
	GoodsSizeConfig  map[string]interface{} `bson:"goods_size_config,omitempty" json:"goods_size_config,omitempty"`
	IsCVD            bool                   `bson:"is_cvd,omitempty" json:"is_cvd,omitempty"`
	BusyRespSettings map[string]interface{} `bson:"busy_resp_settings,omitempty" json:"busy_resp_settings,omitempty"`

	PluginExpireTime   time.Time `bson:"plugin_expire_time,omitempty" json:"plugin_expire_time,omitempty"`
	PluginServiceLevel string    `bson:"plugin_service_level,omitempty" json:"plugin_service_level,omitempty"`
	KlkExpireTime      time.Time `bson:"klk_expire_time,omitempty" json:"klk_expire_time,omitempty"`
	KlkServiceLevel    string    `bson:"klk_service_level,omitempty" json:"klk_service_level,omitempty"`
	YuchiExpireTime    time.Time `bson:"yuchi_expire_time,omitempty" json:"yuchi_expire_time,omitempty"`
	YuchiServiceLevel  string    `bson:"yuchi_service_level,omitempty" json:"yuchi_service_level,omitempty"`
	OriServiceLevel    string    `bson:"ori_service_level,omitempty" json:"ori_service_level,omitempty"`
	Channel            string    `bson:"channel,omitempty" json:"channel,omitempty"`

	IsExpire bool     `bson:"-,omitempty" json:"is_expire,omitempty"`
	Features []string `bson:"features,omitempty" json:"features,omitempty"`

	AfterSalesVersion int `bson:"after_sales_version,omitempty" json:"after_sales_version,omitempty"` // 售后机器人版本

	IsZhongji          int       `bson:"is_zhongji,omitempty" json:"is_zhongji,omitempty"`   // 快手判断是否是中级版
	MarketExpireTime   time.Time `bson:"market_expire_time" json:"market_expire_time"`       //服务市场过期时间
	TrueAuthExpireTime time.Time `bson:"true_auth_expire_time" json:"true_auth_expire_time"` // 通过淘宝诊断工具查到的真实授权过期时间
}

func (u *UserServiceInfo) GetExpireTime() time.Time {
	expireTime := u.ExpireTime
	if u.KlkExpireTime.After(expireTime) {
		expireTime = u.KlkExpireTime
	}
	if u.PluginExpireTime.After(expireTime) {
		expireTime = u.PluginExpireTime
	}
	if u.YuchiExpireTime.After(expireTime) {
		expireTime = u.YuchiExpireTime
	}
	return expireTime
}

type UserServiceDB struct {
	*mongoc_official.Base
}

var userServiceDB = &UserServiceDB{
	Base: mongoc_official.NewBaseModel("xdmp", "xdmp", "user_service"),
}

func UserServiceModel() *UserServiceDB {
	return userServiceDB
}

func (u *UserServiceDB) FindByUserID(ctx context.Context, userID primitive.ObjectID) (*UserServiceInfo, error) {
	info := UserServiceInfo{}
	err := u.FindOne(ctx, bson.M{"user_id": userID}, &info)
	return &info, err
}
