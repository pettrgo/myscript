package model

import (
	"context"
	"encoding/json"
	"gitlab.xiaoduoai.com/golib/xd_sdk/logger"
	"gitlab.xiaoduoai.com/golib/xd_sdk/mongoc_official"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"myscript/utils"
	"sort"
)

type Task struct {
	OID      primitive.ObjectID `bson:"_id" json:"id"`
	Platform PlatformName       `bson:"platform" json:"platform"` // 平台：tb(淘宝)、jd(京东)
	ShopID   string             `bson:"shop_id" json:"shop_id"`
	NodeType string             `bson:"node_type" json:"node_type"`

	//AtAny       []*NodeState   `bson:"at_any" json:"at_any"`
	Name        string         `bson:"name" json:"name"`
	Enable      bool           `bson:"enable" json:"enable"`
	Status      string         `json:"status,omitempty"`         // 列表显示的任务状态：未开始-unstart、进行中-starting、已结束-ended、已关闭-closed
	MsgType     string         `bson:"msg_type" json:"msg_type"` // 发送短信渠道：千牛平台-client， 短信-sms
	ActivityID  string         `bson:"activity_id,omitempty" json:"activity_id,omitempty"`
	TemplateKey TemplateKey    `bson:"template_key" json:"template_key"`
	Hidden      bool           `bson:"hidden" json:"hidden"` // 是否在任务列表显示
	Rules       []*RuleSetting `bson:"rules" json:"rules"`
	IsDeleted   bool           `bson:"is_deleted" json:"is_deleted"`
	LastEnable  bool           `bson:"last_enable" json:"last_enable"`
	SendCount   string         `bson:"-" json:"send_count" `
	Weight      int            `bson:"weight,omitempty" json:"weight,omitempty"`

	IsABTest bool   `bson:"is_abtest,omitempty" json:"is_abtest,omitempty"`
	PlanID   string `bson:"plan_id,omitempty" json:"plan_id,omitempty"`

	ExtraType map[string]bool `bson:"extra_type" json:"extra_type"`
}

type TemplateKey int

// TemplateKeys
const (
	TemplateDefault         TemplateKey = iota // 通常的发送模式
	TemplateTimeRange                          // 指定时间段发送
	TemplateExclusiveSeller                    // 专属客服
)

func (tk TemplateKey) String() string {
	switch tk {
	case TemplateDefault:
		return "<Template(default)>"
	case TemplateTimeRange:
		return "<Template(time range)>"
	default:
		return "<Template(unknown)>"
	}
}

func (t *Task) FindRule(ruleType RuleType) *RuleSetting {
	for _, rule := range t.Rules {
		if rule.Type == ruleType {
			return rule
		}
	}
	return nil
}

func (t *Task) IsEnableCard(ctx context.Context) (bool, error) {
	// 仅任务开启且为下单未付款场景任务需要判断
	if !t.Enable || t.NodeType != Created {
		return false, nil
	}
	rule := t.FindRule(SendMessageMultiTricksRule)
	if rule == nil {
		return false, nil
	}
	s := SendMessageMultiTricksRuleSetting{}
	if err := json.Unmarshal(rule.Args, &s); err != nil {
		return false, err
	}
	sort.Slice(s.Replies, func(i, j int) bool {
		return s.Replies[i].Round < s.Replies[j].Round
	})
	for _, reply := range s.Replies {
		if reply.Enable && (reply.SendAction == 0 || reply.SendAction == 2) {
			if t.ShopID == "5f8ff0c0a3967d00188dca48" {
				logger.Debugf(ctx, "task: %s", utils.UnsafeMarshal(t))
			}
			return true, nil
		}
	}
	return false, nil
}

type TaskDB struct {
	*mongoc_official.Base
}

var taskDB = &TaskDB{
	Base: mongoc_official.NewBaseModel("reminder", "reminder", "task"),
}

func TaskModel() *TaskDB {
	return taskDB
}

func (t *TaskDB) Find(ctx context.Context, query primitive.M) (*mongo.Cursor, error) {
	return t.C(ctx).Find(ctx, query)
}
