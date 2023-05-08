package script

//var TradeOrderCleanCmd = &cobra.Command{
//	Use:   "trade_order_clean",
//	Short: "清理trade_orders订单数据",
//	Long:  "清理trade_orders订单数据",
//	PreRun: func(cmd *cobra.Command, args []string) {
//		esmodel.Init()
//	},
//	Run: process,
//}
//
//type TradeOrderCleanHandler struct {
//	TradeOrderClient *trade_orders.TradeOrdersEsModel
//}
//
//func tradeOrderCleanMain() {
//	h := TradeOrderCleanHandler{
//		TradeOrderClient: trade_orders.Get(),
//	}
//	//h.TradeOrderClient.DeleteService().Do(ctx)
//}
