package model

import (
	"encoding/json"
	"github.com/globalsign/mgo/bson"
	"myscript/utils"
	"time"
)

type OrderStatus struct {
	Status string `bson:"status,omitempty"`
	Time   int64  `bson:"time,omitempty"`
}

type LogisticsInfo struct {
	TrackingNo  string `bson:"tracking_no"`
	Company     string `bson:"company"`
	ShipTime    int64  `bson:"ship_time"`
	DeliveryId  string `bson:"delivery_id"`
	CompanyName string `bson:"company_name"`
}

// 订单结构定义。
type Order struct {
	ID               bson.ObjectId   `bson:"_id,omitempty"`
	Platform         string          `bson:"platform,omitempty"`
	OrderID          string          `bson:"order_id,omitempty"`
	OriginalOrderID  string          `bson:"original_order_id,omitempty"` // 为记录订单原始订单ID
	ShopID           string          `bson:"shop_id,omitempty"`
	SellerID         string          `bson:"seller_id,omitempty"`
	BuyerID          string          `bson:"buyer_id,omitempty"`
	OriBuyerID       string          `bson:"ori_buyer_id,omitempty"` // 为记录京东原始buyer_id，京东的BuyerID转为小写
	Payment          float64         `bson:"payment,omitempty"`
	Address          string          `bson:"address,omitempty"`
	Province         string          `bson:"province,omitempty"`
	City             string          `bson:"city,omitempty"`
	Town             string          `bson:"town,omitempty"`
	Street           string          `bson:"street,omitempty"`
	DecryptAddress   string          `bson:"decrypt_address,omitempty"`
	Status           string          `bson:"status,omitempty"`
	OriginStatus     string          `bson:"origin_status,omitempty"`
	OrderType        string          `bson:"order_type,omitempty"`
	StepTradeStatus  string          `bson:"step_trade_status,omitempty"`
	StepPaidFee      float64         `bson:"step_paid_fee,omitempty"`
	PushStatus       string          `bson:"push_status,omitempty"`
	StatusHistory    []OrderStatus   `bson:"status_history,omitempty"`
	BuyerRemark      string          `bson:"buyer_remark,omitempty"`
	SellerMemo       string          `bson:"seller_memo,omitempty"`        // 卖家留言
	BalanceUsed      string          `bson:"balance_used,omitempty"`       // 余额支付
	SellerDiscount   string          `bson:"seller_discount,omitempty"`    // 商家优惠金额
	OrderSellerPrice string          `bson:"order_seller_price,omitempty"` // 订单货款金额（结算金额）
	PayType          string          `bson:"pay_type,omitempty"`           // 支付方式
	FreightPrice     string          `bson:"freight_price,omitempty"`      // 商品的运费
	CreatedAt        int64           `bson:"created_at,omitempty"`
	UpdatedAt        int64           `bson:"updated_at,omitempty"`
	ItemIDs          []string        `bson:"item_ids,omitempty"`
	Tbext            *Tb             `bson:"tbext,omitempty"`
	Jdext            *Jd             `bson:"jdext,omitempty"`
	Pddext           *Pdd            `bson:"pddext,omitempty"`
	Lazy             bool            `bson:"lazy,omitempty"`
	NotInWhitelist   bool            `bson:"not_in_whitelist,omitempty"`
	UpdateTime       time.Time       `bson:"update_time,omitempty"`
	UpdateDesc       string          `bson:"update_desc,omitempty"`
	BuyerOpenUID     string          `bson:"buyer_open_uid,omitempty"`
	RealBuyerNick    string          `bson:"real_buyer_nick,omitempty"`
	LogisticsInfo    []LogisticsInfo `bson:"logistics_info"`
}

type Jd struct {
	BuyerNick   string         `bson:"buyer_nick,omitempty"`
	StateDesc   string         `bson:"state_desc,omitempty"`
	Items       []*JdItem      `bson:"items,omitempty"`
	Receiver    *JdReceiver    `bson:"receiver,omitempty"`
	ConsignInfo *JdConsignInfo `bson:"consign_info,omitempty"`
}

type JdItem struct {
	ItemID   string  `bson:"item_id,omitempty"`
	ItemName string  `bson:"item_name,omitempty"`
	SkuID    string  `bson:"sku_id,omitempty"`
	SkuName  string  `bson:"sku_name,omitempty"`
	Count    int64   `bson:"count,omitempty"`
	Price    float32 `bson:"price,omitempty"`
}

type JdReceiver struct {
	Name  string `bson:"receiver,omitempty"`
	Phone string `bson:"receiver_phone,omitempty"`
}

type JdConsignInfo struct {
	Province string `bson:"province,omitempty"`
	City     string `bson:"city,omitempty"`
	County   string `bson:"county,omitempty"`
	Town     string `bson:"town,omitempty"`
}

// 淘宝特有字段。
type Tb struct {
	BuyerNick     string      `bson:"buyer_nick,omitempty"`
	BuyerRemark   string      `bson:"buyer_remark,omitempty"`
	Items         []*TbItem   `bson:"items,omitempty"`
	Receiver      *TbReceiver `bson:"receiver,omitempty"`
	SellerFlag    int64       `bson:"seller_flag,omitempty"`
	OriginalOrder string      `bson:"original_order,omitempty"`
}

type OriginalOrder struct {
	SellerNick  *string `json:"seller_nick,omitempty"`
	BuyerNick   *string `json:"buyer_nick,omitempty"`
	Created     *string `json:"created,omitempty"`
	Modified    *string `json:"modified,omitempty"`
	ServiceTags *struct {
		LogisticsTag []struct {
			OrderID                *string `json:"order_id,omitempty"`
			LogisticServiceTagList *struct {
				LogisticServiceTag []struct {
					ServiceTag  *string `json:"service_tag,omitempty"`
					ServiceType *string `json:"service_type,omitempty"`
				} `json:"logistic_service_tag,omitempty"`
			} `json:"logistic_service_tag_list,omitempty"`
		} `json:"logistics_tag,omitempty"`
	} `json:"service_tags,omitempty"`
	Orders *struct {
		Order []struct {
			Num              *int64  `json:"num,omitempty"`
			NumIid           *int64  `json:"num_iid,omitempty"`
			OuterIid         *string `json:"outer_iid,omitempty"`
			Cid              *int64  `json:"cid,omitempty"`
			Payment          *string `json:"payment,omitempty"`
			PicPath          *string `json:"pic_path,omitempty"`
			Sku              *string `json:"sku_properties_name,omitempty"`
			SkuID            *string `json:"sku_id,omitempty"`
			OuterSkuID       *string `json:"outer_sku_id,omitempty"`
			Title            *string `json:"title,omitempty"`
			Price            *string `json:"price,omitempty"`
			Oid              *string `json:"oid,omitempty"`
			OrderAttr        *string `json:"order_attr,omitempty"`
			ShippingType     *string `json:"shipping_type,omitempty"`
			LogisticsCompany *string `json:"logistics_company,omitempty"`
			InvoiceNo        *string `json:"invoice_no,omitempty"`
			ConsignTime      *string `json:"consign_time,omitempty"`
			Status           *string `json:"status,omitempty"`
			RefundStatus     *string `json:"refund_status,omitempty"`
			TotalFee         *string `json:"total_fee,omitempty"`
		} `json:"order,omitempty"`
	} `json:"orders,omitempty"`
	Payment                 *string `json:"payment,omitempty"`
	ReceiverAddress         *string `json:"receiver_address,omitempty"`
	ReceiverCity            *string `json:"receiver_city,omitempty"`
	ReceiverCountry         *string `json:"receiver_country,omitempty"`
	ReceiverDistrict        *string `json:"receiver_district,omitempty"`
	ReceiverMobile          *string `json:"receiver_mobile,omitempty"`
	ReceiverName            *string `json:"receiver_name,omitempty"`
	ReceiverPhone           *string `json:"receiver_phone,omitempty"`
	ReceiverState           *string `json:"receiver_state,omitempty"`
	ReceiverTown            *string `json:"receiver_town,omitempty"`
	ReceiverZip             *string `json:"receiver_zip,omitempty"`
	BuyerMessage            *string `json:"buyer_message,omitempty"`
	BuyerMemo               *string `json:"buyer_memo,omitempty"`
	Tid                     *string `json:"tid,omitempty"`
	Status                  *string `json:"status,omitempty"`
	StepTradeStatus         *string `json:"step_trade_status,omitempty"`
	DecryptReceiverAddress  *string `json:"decrypt_receiver_address,omitempty"`
	DecryptReceiverCity     *string `json:"decrypt_receiver_city,omitempty"`
	DecryptReceiverDistrict *string `json:"decrypt_receiver_district,omitempty"`
	DecryptReceiverMobile   *string `json:"decrypt_receiver_mobile,omitempty"`
	DecryptReceiverName     *string `json:"decrypt_receiver_name,omitempty"`
	DecryptReceiverPhone    *string `json:"decrypt_receiver_phone,omitempty"`
	DecryptReceiverState    *string `json:"decrypt_receiver_state,omitempty"`
	DecryptReceiverTown     *string `json:"decrypt_receiver_town,omitempty"`
	BuyerAlipayNo           *string `json:"buyer_alipay_no,omitempty"`
	BuyerArea               *string `json:"buyer_area,omitempty"`
	BuyerEmail              *string `json:"buyer_email,omitempty"`
	BuyerIP                 *string `json:"buyer_ip,omitempty"`
	BuyerRate               *bool   `json:"buyer_rate,omitempty"`
	PostFee                 *string `json:"post_fee,omitempty"`
	SellerAlipayNo          *string `json:"seller_alipay_no,omitempty"`
	SellerEmail             *string `json:"seller_email,omitempty"`
	SellerMobile            *string `json:"seller_mobile,omitempty"`
	SellerName              *string `json:"seller_name,omitempty"`
	SellerFlag              *int    `json:"seller_flag,omitempty"`
	SellerMemo              *string `json:"seller_memo,omitempty"`
	Sid                     *string `json:"sid,omitempty"`
	LogisticServiceOs       *struct {
		OsDate  *string `json:"osDate,omitempty"`
		OsRange *string `json:"osRange,omitempty"`
	} `json:"logistic_service_os,omitempty"` // 预约派送信息
	RealBuyerNick *string `json:"real_buyer_nick,omitempty"`
}

func (m *Order) NeedUpdate() bool {
	if m.Tbext != nil && m.Tbext.OriginalOrder != "" {
		ori := utils.SortJsonStr(m.Tbext.OriginalOrder)
		newOriOrder := utils.SortJsonStr(m.Tbext.FieldClipping().OriginalOrder)
		if ori != newOriOrder {
			return true
		}
	}
	return false
}

func (m *Order) FieldClipping() *Order {
	ret := &Order{
		ID:               m.ID,
		Platform:         m.Platform,
		OrderID:          m.OrderID,
		OriginalOrderID:  m.OriginalOrderID,
		ShopID:           m.ShopID,
		SellerID:         m.SellerID,
		BuyerID:          m.BuyerID,
		OriBuyerID:       m.OriBuyerID,
		Payment:          m.Payment,
		Address:          m.Address,
		Province:         m.Province,
		City:             m.City,
		Town:             m.Town,
		Street:           m.Street,
		DecryptAddress:   m.DecryptAddress,
		Status:           m.Status,
		OriginStatus:     m.OriginStatus,
		OrderType:        m.OrderType,
		StepTradeStatus:  m.StepTradeStatus,
		StepPaidFee:      m.StepPaidFee,
		PushStatus:       m.PushStatus,
		StatusHistory:    m.StatusHistory,
		BuyerRemark:      m.BuyerRemark,
		SellerMemo:       m.SellerMemo,
		BalanceUsed:      m.BalanceUsed,
		SellerDiscount:   m.SellerDiscount,
		OrderSellerPrice: m.OrderSellerPrice,
		PayType:          m.PayType,
		FreightPrice:     m.FreightPrice,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
		ItemIDs:          m.ItemIDs,
		Tbext:            m.Tbext.FieldClipping(),
		Jdext:            m.Jdext,
		Pddext:           m.Pddext,
		Lazy:             m.Lazy,
		NotInWhitelist:   m.NotInWhitelist,
		UpdateTime:       m.UpdateTime,
		UpdateDesc:       m.UpdateDesc,
		BuyerOpenUID:     m.BuyerOpenUID,
		RealBuyerNick:    m.RealBuyerNick,
		LogisticsInfo:    m.LogisticsInfo,
	}
	return ret
}

func (t *Tb) FieldClipping() *Tb {
	if t == nil {
		return t
	}
	ori := t.OriginalOrder
	temp := struct {
		SellerNick  *string `json:"seller_nick,omitempty"`
		BuyerNick   *string `json:"buyer_nick,omitempty"`
		Created     *string `json:"created,omitempty"`
		Modified    *string `json:"modified,omitempty"`
		ServiceTags *struct {
			LogisticsTag []struct {
				OrderID                *string `json:"order_id,omitempty"`
				LogisticServiceTagList *struct {
					LogisticServiceTag []struct {
						ServiceTag  *string `json:"service_tag,omitempty"`
						ServiceType *string `json:"service_type,omitempty"`
					} `json:"logistic_service_tag,omitempty"`
				} `json:"logistic_service_tag_list,omitempty"`
			} `json:"logistics_tag,omitempty"`
		} `json:"service_tags,omitempty"`
		Orders *struct {
			Order []struct {
				Num              *int64  `json:"num,omitempty"`
				NumIid           *int64  `json:"num_iid,omitempty"`
				OuterIid         *string `json:"outer_iid,omitempty"`
				Cid              *int64  `json:"cid,omitempty"`
				Payment          *string `json:"payment,omitempty"`
				PicPath          *string `json:"pic_path,omitempty"`
				Sku              *string `json:"sku_properties_name,omitempty"`
				SkuID            *string `json:"sku_id,omitempty"`
				OuterSkuID       *string `json:"outer_sku_id,omitempty"`
				Title            *string `json:"title,omitempty"`
				Price            *string `json:"price,omitempty"`
				Oid              *string `json:"oid,omitempty"`
				OrderAttr        *string `json:"order_attr,omitempty"`
				ShippingType     *string `json:"shipping_type,omitempty"`
				LogisticsCompany *string `json:"logistics_company,omitempty"`
				InvoiceNo        *string `json:"invoice_no,omitempty"`
				ConsignTime      *string `json:"consign_time,omitempty"`
				Status           *string `json:"status,omitempty"`
				RefundStatus     *string `json:"refund_status,omitempty"`
				TotalFee         *string `json:"total_fee,omitempty"`
			} `json:"order,omitempty"`
		} `json:"orders,omitempty"`
		Payment                 *string `json:"payment,omitempty"`
		ReceiverAddress         *string `json:"receiver_address,omitempty"`
		ReceiverCity            *string `json:"receiver_city,omitempty"`
		ReceiverCountry         *string `json:"receiver_country,omitempty"`
		ReceiverDistrict        *string `json:"receiver_district,omitempty"`
		ReceiverMobile          *string `json:"receiver_mobile,omitempty"`
		ReceiverName            *string `json:"receiver_name,omitempty"`
		ReceiverPhone           *string `json:"receiver_phone,omitempty"`
		ReceiverState           *string `json:"receiver_state,omitempty"`
		ReceiverTown            *string `json:"receiver_town,omitempty"`
		ReceiverZip             *string `json:"receiver_zip,omitempty"`
		BuyerMessage            *string `json:"buyer_message,omitempty"`
		BuyerMemo               *string `json:"buyer_memo,omitempty"`
		Tid                     *string `json:"tid,omitempty"`
		Status                  *string `json:"status,omitempty"`
		StepTradeStatus         *string `json:"step_trade_status,omitempty"`
		DecryptReceiverAddress  *string `json:"decrypt_receiver_address,omitempty"`
		DecryptReceiverCity     *string `json:"decrypt_receiver_city,omitempty"`
		DecryptReceiverDistrict *string `json:"decrypt_receiver_district,omitempty"`
		DecryptReceiverMobile   *string `json:"decrypt_receiver_mobile,omitempty"`
		DecryptReceiverName     *string `json:"decrypt_receiver_name,omitempty"`
		DecryptReceiverPhone    *string `json:"decrypt_receiver_phone,omitempty"`
		DecryptReceiverState    *string `json:"decrypt_receiver_state,omitempty"`
		DecryptReceiverTown     *string `json:"decrypt_receiver_town,omitempty"`
		BuyerAlipayNo           *string `json:"buyer_alipay_no,omitempty"`
		BuyerArea               *string `json:"buyer_area,omitempty"`
		BuyerEmail              *string `json:"buyer_email,omitempty"`
		BuyerIP                 *string `json:"buyer_ip,omitempty"`
		BuyerRate               *bool   `json:"buyer_rate,omitempty"`
		PostFee                 *string `json:"post_fee,omitempty"`
		SellerAlipayNo          *string `json:"seller_alipay_no,omitempty"`
		SellerEmail             *string `json:"seller_email,omitempty"`
		SellerMobile            *string `json:"seller_mobile,omitempty"`
		SellerName              *string `json:"seller_name,omitempty"`
		SellerFlag              *int    `json:"seller_flag,omitempty"`
		SellerMemo              *string `json:"seller_memo,omitempty"`
		Sid                     *string `json:"sid,omitempty"`
		LogisticServiceOs       *struct {
			OsDate  *string `json:"osDate,omitempty"`
			OsRange *string `json:"osRange,omitempty"`
		} `json:"logistic_service_os,omitempty"` // 预约派送信息
		RealBuyerNick *string `json:"real_buyer_nick,omitempty"`
	}{}
	if err := json.Unmarshal([]byte(ori), &temp); err != nil {
		return t
	}
	data, err := json.Marshal(temp)
	if err != nil {
		return t
	}
	return &Tb{
		BuyerNick:     t.BuyerNick,
		BuyerRemark:   t.BuyerRemark,
		Items:         t.Items,
		Receiver:      t.Receiver,
		SellerFlag:    t.SellerFlag,
		OriginalOrder: string(data),
	}
}

// 商品结构定义。
type TbItem struct {
	ItemID   string  `bson:"item_id,omitempty"`
	ItemName string  `bson:"item_name,omitempty"`
	SkuID    string  `bson:"sku_id,omitempty"`
	SkuName  string  `bson:"sku_name,omitempty"`
	Count    int64   `bson:"count,omitempty"`
	Price    float32 `bson:"price,omitempty"`
}

// 收件人定义。
type TbReceiver struct {
	Name          string `bson:"name,omitempty"`
	Phone         string `bson:"phone,omitempty"`
	Zip           string `bson:"zip,omitempty"`
	DecryptName   string `bson:"decrypt_name,omitempty"`
	DecryptPhone  string `bson:"decrypt_phone,omitempty"`
	DecryptMobile string `bson:"decrypt_mobile,omitempty"`
}

type Pdd struct {
	BuyerNick    string       `bson:"buyer_nick,omitempty"`
	StateDesc    string       `bson:"state_desc,omitempty"`
	Items        []*PddItem   `bson:"items,omitempty"`
	Receiver     *PddReceiver `bson:"receiver,omitempty"`
	GroupOrderId string       `bson:"group_order_id,omitempty"`
}

type PddItem struct {
	ItemID   string  `bson:"item_id,omitempty"`
	ItemName string  `bson:"item_name,omitempty"`
	SkuID    string  `bson:"sku_id,omitempty"`
	SkuName  string  `bson:"sku_name,omitempty"`
	Count    int64   `bson:"count,omitempty"`
	Price    float32 `bson:"price,omitempty"`
}

type PddReceiver struct {
	Name  string `bson:"receiver,omitempty"`
	Phone string `bson:"receiver_phone,omitempty"`
}
