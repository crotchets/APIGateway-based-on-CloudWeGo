package clientprovider

import (
	"APIGateway/biz/idlprovider"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/connpool"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/loadbalance"
	etcd "github.com/kitex-contrib/registry-etcd"
	"time"
)

type ClientProvider struct {
	cli genericclient.Client
}

func NewClientProvider() *ClientProvider {
	return &ClientProvider{}
}

func (provider *ClientProvider) GetClient(rpcName string, version string) (*genericclient.Client, error) {
	if provider.cli == nil {
		r, err := etcd.NewEtcdResolver([]string{"127.0.0.1:2379"}) // ETCD服务发现
		if err != nil {
			return nil, err
		}
		var opts []client.Option
		opts = append(opts, client.WithResolver(r))                                           // 解析器
		opts = append(opts, client.WithLoadBalancer(loadbalance.NewWeightedRandomBalancer())) // 负载均衡

		// 设置长连接配置
		cfg := connpool.IdleConfig{
			MaxIdlePerAddress: 10,
			MaxIdleGlobal:     10,
			MaxIdleTimeout:    60 * time.Second,
		}
		opts = append(opts, client.WithLongConnection(cfg))

		content, err := idlprovider.GetIdlProvider().GetIdl(rpcName, version) // 获取IDL文件内容，现在provider只是一个无情的转发机器，等待加入内容缓存
		if err != nil {
			return nil, err
		}
		p, err := generic.NewThriftContentProvider(content, nil)
		if err != nil {
			return nil, err
		}
		g, err := generic.JSONThriftGeneric(p)
		if err != nil {
			return nil, err
		}
		provider.cli, _ = genericclient.NewClient(rpcName, g, opts...)
	}
	return &provider.cli, nil
}
