package esmodel

import (
	"context"
	"github.com/pkg/errors"
	"gitlab.xiaoduoai.com/golib/xd_sdk/es_official"
	"myscript/config"
)

type EsModel interface {
	Name() string
}

func InitEs() {
	ctx := context.Background()
	esConfigs := config.GetConfig().EsConfigs
	for _, c := range esConfigs {
		if err := es_official.ClientsMgr().NewClient(ctx, c.IndexName, c.Addrs...); err != nil {
			panic(errors.Wrap(err, "init es client failed"))
		}
	}
}
