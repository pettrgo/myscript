package model

const (
	ShoppingCartAdded      = "shopping_cart_added"    //加购未下单
	Asked                  = "asked"                  //咨询未下单
	Created                = "created"                //下单未付款
	CheckAddr              = "checkaddr"              //核对地址
	Deposited              = "deposited"              //已付定金未付尾款
	Ungroup                = "ungroup"                //待成团
	Group                  = "group"                  //成团
	Paid                   = "paid"                   //已付款
	Pause                  = "paused"                 //订单暂停
	Locked                 = "locked"                 // 订单锁定
	ShippedDelay           = "shipped_delay"          //发货超时
	Shipped                = "shipped"                //已发货
	GotDelay               = "got_delay"              //揽件超时
	LogisticsUpdateDelay   = "logistics_update_delay" //物流更新超时
	LogisticsStagnateDelay = "logistics_stagnate_delay"
	SendScan               = "send_scan"       //验收提醒
	SendScanDelay          = "send_scan_delay" //验收提醒超时
	Signed                 = "signed"          //已签收
	Succeeded              = "succeeded"       //交易成功
	Finished               = "finished"        //京东交易成功
	TradeCanceled          = "trade_canceled"  // 京东订单取消
	Canceled               = "canceled"        // 京东自营订单取消
	Closed                 = "closed"          //订单取消
	Trigger                = "trigger"
	Refund                 = "refund"
	Rated                  = "rated"
	RefundAgree            = "refund_agree"
	RefundSuccess          = "refund_success"
	RefundReject           = "refund_reject"
)
