package model

import "encoding/json"

type RuleType string

// RuleSetting is setting of rule, Args is json raw message
type RuleSetting struct {
	Type RuleType        `bson:"type" json:"type"`
	Args json.RawMessage `bson:"args" json:"args"`
}

const SendMessageMultiTricksRule = RuleType("send_message_multi_tricks")

type SendMessageMultiTricksRuleSetting struct {
	Replies []*Reply `json:"replies"`
}

type Reply struct {
	Round      int    `json:"round"`
	Enable     bool   `json:"enable"`
	AgeingID   string `json:"ageing_id"`
	StateDelay *int   `json:"state_delay,omitempty"`
	PushByDay  *int   `json:"push_by_day,omitempty"`
	//Images     []*model.Image `json:"images"`
	Message    string `json:"message"`
	SendAction int    `json:"send_action"` // 发送模式：0默认值 话术+卡片，1-仅话术，2-仅卡片；

	AfterShoppingCartGoodsAdded bool   `json:"after_shopping_cart_goods_added"`
	ReceiveMessageFilter        bool   `json:"receive_message_filter"`
	ContentType                 string `json:"content_type"`
	DisablePaymentCardMsg       bool   `json:"disable_payment_card_msg"`

	//Extensions []model.ExtraInfo `json:"extensions"`
}
