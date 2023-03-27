package mongo

import (
	"context"
	"fmt"
	"gitlab.xiaoduoai.com/golib/xd_sdk/logger"
	"gitlab.xiaoduoai.com/golib/xd_sdk/mongoc_official"
	"myscript/utils"
	"strings"
)

type MongoMap map[string]MongoConfig

type MongoConfig struct {
	Addrs          []string `json:"addrs" mapstructure:"addrs" example:"127.0.0.1:27017"`
	Source         string   `json:"source" mapstructure:"source" example:""`
	ReplicaSetName string   `json:"replica_set_name" mapstructure:"replica_set_name" example:""`
	Timeout        int      `json:"timeout" mapstructure:"timeout" example:"2"`
	Username       string   `json:"username" mapstructure:"username" example:""`
	Password       string   `json:"password" mapstructure:"password" example:""`
	Mode           *int     `json:"mode,omitempty" mapstructure:"mode,omitempty" example:"3"`
	Alias          string   `json:"alias" mapstructure:"alias" example:"default"`
	AppName        string   `mapstructure:"app_name"`
}

func Init(mc MongoMap) {
	if err := InitMongoMap(mc); err != nil {
		panic(err)
	}
}

func InitMongoMap(mongoMap MongoMap) error {
	var err error
	ctx := context.Background()
	fmt.Println(utils.UnsafeMarshal(mongoMap))
	for _, v := range mongoMap {
		if err = initMongo(ctx, &v); err != nil {
			return err
		} else {
			logger.Infof(ctx, "init mongo:%v", v)
		}
	}
	return nil
}

func initMongo(ctx context.Context, mc *MongoConfig) error {
	mgoUri := fmt.Sprintf("mongodb://%s:%s@%s/?replicaSet=%s&auth=%s&w=majority&wtimeoutMS=%v&appName=%s",
		mc.Username, mc.Password, strings.Join(mc.Addrs, ","), mc.ReplicaSetName, mc.Source, mc.Timeout, mc.AppName)
	err := mongoc_official.AddClientAddress(ctx, mc.AppName, mgoUri)

	if err != nil {
		return err
	}

	//dialInfo := mc.getMongoDialInfo()
	//s, err := mgo.DialWithInfo(dialInfo)
	//if err != nil {
	//	return err
	//}
	//s.SetMode(mgo.Mode(*mc.Mode), true)
	//sessionLock.Lock()
	//defer sessionLock.Unlock()
	//if _, ok := sessions[mc.Alias]; ok {
	//	return errors.New("duplicate session:" + mc.Alias)
	//}
	//if len(sessions) == 0 {
	//	defaultSession = s
	//}
	//sessions[mc.Alias] = s
	return nil
}
