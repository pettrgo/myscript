package script

import (
	"context"
	"encoding/csv"
	"github.com/spf13/cobra"
	"gitlab.xiaoduoai.com/golib/xd_sdk/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"myscript/config"
	"myscript/model"
	"myscript/storage/mongo"
	"os"
)

func init() {
	TaskSearchCmd.PersistentFlags().StringVarP(&platform, "platform", "p", "", "平台")
}

var TaskSearchCmd = &cobra.Command{
	Use:   "task_search",
	Short: "搜索任务",
	Long:  "搜索任务",
	PreRun: func(cmd *cobra.Command, args []string) {
		conf := config.GetConfig()
		mongo.Init(conf.Mongos)
	},
	Run: taskSearchMain,
}

type taskHandler struct {
	//taskCh  chan *model.Task
	Lines [][]string
}

func taskSearchMain(command *cobra.Command, args []string) {
	ctx := context.Background()

	h := &taskHandler{
		//taskCh: make(chan *model.Task),
		Lines: make([][]string, 0),
	}
	header := []string{"店铺id"}
	h.searchTask(ctx)
	if err := h.export(ctx, header); err != nil {
		panic(err)
	}
}

func (h *taskHandler) searchTask(ctx context.Context) {
	query := primitive.M{
		"platform":  platform,
		"enable":    true,
		"node_type": model.Created,
		"is_deleted": primitive.M{
			"$ne": true,
		},
	}
	cur, err := model.TaskModel().Find(ctx, query)
	if err != nil {
		panic(err)
	}
	shopDupCheck := make(map[string]struct{})
	for cur.Next(ctx) {
		task := model.Task{}
		if err := cur.Decode(&task); err != nil {
			logger.Errorf(ctx, "task decode failed, err: %v", err)
			continue
		}
		ok, err := task.IsEnableCard(ctx)
		if err != nil {
			logger.Errorf(ctx, "check task enable card failed, task_id: %v err: %v", task.OID, err)
			continue
		}
		if ok {
			if _, ok := shopDupCheck[task.ShopID]; ok {
				continue
			}
			h.Lines = append(h.Lines, []string{task.ShopID})
			shopDupCheck[task.ShopID] = struct{}{}
		}
	}
}

//func (h *taskHandler) handlerTask(ctx context.Context) {
//	for task := range h.taskCh {
//		ok, err := task.IsEnableCard(ctx)
//		if err != nil {
//			logger.Errorf(ctx, "check task enable card failed, task_id: %v err: %v", task.OID, err)
//			continue
//		}
//		if ok {
//			h.shopIDs = append(h.shopIDs, task.ShopID)
//		}
//	}
//}

func (h *taskHandler) export(ctx context.Context, header []string) error {
	file, err := os.Create("./created_card.csv")
	if err != nil {
		return err
	}
	defer file.Close()
	file.WriteString("\xEF\xBB\xBF")
	writer := csv.NewWriter(file)
	defer writer.Flush()
	if err := writer.Write(header); err != nil {
		return err
	}
	for _, line := range h.Lines {
		if err := writer.Write(line); err != nil {
			logger.WithError(err).Errorf(ctx, "write csv failed, line: %v", line)
		}
	}
	return nil
}
