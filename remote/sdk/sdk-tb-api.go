package sdk

import (
	"context"
	"gitlab.xiaoduoai.com/golib/xd_sdk/httpclient"
	"myscript/config"
	"sync"
	"time"
)

type TbApiClient struct {
	SdkClient httpclient.Client
	Once      sync.Once
}

var (
	client *TbApiClient
	once   sync.Once
)

func TBClient() *TbApiClient {
	once.Do(func() {
		conf := config.GetConfig().RemoteService.SdkTbApi
		sdkClient, err := httpclient.NewClient(
			httpclient.WithAddress(conf.Addr),
			httpclient.WithTimeout(time.Duration(conf.Timeout)*time.Second),
		)
		if err != nil {
			panic(err)
		}
		client = &TbApiClient{
			SdkClient: sdkClient,
		}
	})
	return client
}

func (client *TbApiClient) BatchDeleteJdp(ctx context.Context, shopIDs []string) error {
	for _, shopID := range shopIDs {
		_, _ = client.SdkClient.NewRequest(ctx).SetQueryParams(map[string]string{
			"shop_id": shopID,
		}).Get("/sdk/jushita/jdp/user/delete")
	}
	return nil
}
