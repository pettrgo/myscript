package model

import (
	"context"
	"gitlab.xiaoduoai.com/golib/xd_sdk/mongoc_official"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ShopInfo struct {
	ID                     primitive.ObjectID     `bson:"_id,omitempty" json:"id,omitempty"`
	CreateTime             time.Time              `bson:"create_time,omitempty" json:"create_time,omitempty"`
	UpdateTime             time.Time              `bson:"update_time,omitempty" json:"update_time,omitempty"`
	Platform               string                 `bson:"platform,omitempty" json:"platform,omitempty"`
	PlatShopID             string                 `bson:"plat_shop_id,omitempty" json:"plat_shop_id,omitempty"`
	PlatShopName           string                 `bson:"plat_shop_name,omitempty" json:"plat_shop_name,omitempty"`
	PlatShopCid            string                 `bson:"plat_shop_cid,omitempty" json:"plat_shop_cid,omitempty"`
	PlatUserID             string                 `bson:"plat_user_id,omitempty" json:"plat_user_id,omitempty"`           // 店铺第一次使用晓多机器人的主账号ID, 不改变
	PlatUserRealID         string                 `bson:"plat_user_real_id,omitempty" json:"plat_user_real_id,omitempty"` // 表示店铺最新的主账号ID
	PlatUserName           string                 `bson:"plat_user_name,omitempty" json:"plat_user_name,omitempty"`
	PlatChildUserIds       []string               `bson:"plat_child_user_ids,omitempty" json:"plat_child_user_ids,omitempty"`
	PlatExtra              map[string]interface{} `bson:"plat_extra,omitempty" json:"plat_extra,omitempty"`
	UserID                 primitive.ObjectID     `bson:"user_id,omitempty" json:"user_id,omitempty"`
	PicPath                string                 `bson:"pic_path,omitempty" json:"pic_path"`
	IsCVD                  string                 `bson:"is_cvd,omitempty" json:"is_cvd,omitempty"`
	Flag                   string                 `bson:"flag,omitempty" json:"flag,omitempty"`
	IsDeleted              bool                   `bson:"is_deleted,omitempty" json:"is_deleted,omitempty"`
	HistoryNicks           []string               `bson:"history_nicks,omitempty" json:"history_nicks,omitempty"`
	AccountLimitV          int                    `bson:"account_limit_v,omitempty" json:"account_limit_v,omitempty"`
	AccountLimit           int                    `bson:"account_limit,omitempty" json:"account_limit"`
	RmSuccessChildUserIds  []string               `bson:"rm_success_child_user_ids,omitempty" json:"rm_success_child_user_ids,omitempty"`
	JdType                 string                 `bson:"jd_type,omitempty" json:"jd_type,omitempty"`
	IP                     string                 `bson:"ip,omitempty" json:"ip,omitempty"`
	MainUsername           string                 `bson:"main_username,omitempty" json:"main_username,omitempty"`
	CategoryId             primitive.ObjectID     `bson:"category_id,omitempty" json:"category_id,omitempty"`
	ReminderVersion        string                 `bson:"reminder_version,omitempty" json:"reminder_version,omitempty"`
	VenderID               int64                  `bson:"vender_id,omitempty" json:"vender_id,omitempty"`
	Version                *int                   `bson:"version,omitempty" json:"version,omitempty"`
	IsCategoryEdited       *bool                  `bson:"is_category_edited,omitempty" json:"is_category_edited,omitempty"`
	VipDuty                string                 `bson:"vip_duty,omitempty" json:"vip_duty,omitempty"`
	Duty                   string                 `bson:"duty,omitempty" json:"duty,omitempty"`
	CustomerArea           string                 `bson:"customer_area,omitempty" json:"customer_area,omitempty"`
	DisabledSubcategoryIds []primitive.ObjectID   `bson:"disabled_subcategory_ids" json:"disabled_subcategory_ids"`
	SupplierVenderIds      []string               `bson:"supplier_vender_ids" json:"supplier_vender_ids"`
	AddSupplierVenderIds   []string               `bson:"add_supplier_vender_ids" json:"add_supplier_vender_ids"`
	SNType                 string                 `bson:"sn_type,omitempty" json:"sn_type,omitempty"`
	CategoryAlias          string                 `bson:"category_alias" json:"category_alias"`
	OpenId                 string                 `bson:"open_id" json:"open_id"`
	RobotSpinWithPre       string                 `bson:"-" json:"robot_spin_with_pre"`
}

type ShopDB struct {
	*mongoc_official.Base
}

var shopDB = ShopDB{
	Base: mongoc_official.NewBaseModel("xdmp", "xdmp", "shop"),
}

func ShopModel() *ShopDB {
	return &shopDB
}

func (s *ShopDB) FindByShopID(ctx context.Context, shopID string) (*ShopInfo, error) {
	id, err := primitive.ObjectIDFromHex(shopID)
	if err != nil {
		return nil, err
	}
	info := ShopInfo{}
	if err := s.FindOne(ctx, primitive.M{"_id": id}, &info); err != nil {
		return nil, err
	}
	return &info, nil
}
