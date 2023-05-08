package script

import (
	"context"
	"github.com/spf13/cobra"
	"gitlab.xiaoduoai.com/golib/xd_sdk/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"myscript/config"
	"myscript/model"
	"myscript/storage/mongo"
)

var DeleteExpireTradePddData = &cobra.Command{
	Use:   "delete_expire_trade_pdd_data",
	Short: "删除过期的trade_pdd_data表数据",
	Long:  "删除过期的trade_pdd_data表数据",
	PreRun: func(cmd *cobra.Command, args []string) {
		conf := config.GetConfig().Mongos
		mongo.InitMongoMap(conf)
	},
	Run: process,
}

type DeleteExpireTradePddDataHandler struct {
	ch chan primitive.ObjectID
}

func DeleteExpireTradePddDataMain() {

}

func (h *DeleteExpireTradePddDataHandler) findExpireTradePddData(ctx context.Context) {
	defer func() {
		close(h.ch)
	}()
	cur, err := model.PddDataModel().Find(ctx, primitive.M{})
	if err != nil {
		return
	}
	for cur.Next(ctx) {
		var pddData model.PddData
		if err := cur.Decode(&pddData); err != nil {
			continue
		}
		h.ch <- pddData.ID
	}
}

func (h *DeleteExpireTradePddDataHandler) deleteExpireTradePddData(ctx context.Context) {
	ids := make([]primitive.ObjectID, 0, 100)
	for id := range h.ch {
		ids = append(ids, id)
		if len(ids) == 100 {
			filter := primitive.M{"_id": primitive.M{"$in": ids}}
			if err := model.PddDataModel().DeleteMany(ctx, filter); err != nil {
				logger.Errorf(ctx, "delete failed, ids: %v, err: %v", ids, err)
			}
		}
	}
}
