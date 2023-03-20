package mongo

import (
	"context"
	"fmt"
	"gitlab.xiaoduoai.com/golib/xd_sdk/logger"
	"gitlab.xiaoduoai.com/golib/xd_sdk/mongoc_official"
	"myscript/config"
	"strings"
)

func Init() {
	mc := config.GetConfig().Mongos
	if err := InitMongoMap(mc); err != nil {
		panic(err)
	}
}

func InitMongoMap(mongoMap map[string]config.MongoConfig) error {
	var err error
	ctx := context.Background()
	for _, v := range mongoMap {
		if err = initMongo(ctx, &v); err != nil {
			return err
		} else {
			logger.Infof(ctx, "init mongo:%v", v)
		}
	}
	return nil
}

func initMongo(ctx context.Context, mc *config.MongoConfig) error {
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
