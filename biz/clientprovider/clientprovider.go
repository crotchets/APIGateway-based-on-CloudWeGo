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
	//todo
}

func NewClientProvider() *ClientProvider {
	return &ClientProvider{}
}

func (provider *ClientProvider) GetClient(rpcName string) (*genericclient.Client, error) {
	r, err := etcd.NewEtcdResolver([]string{"127.0.0.1:2379"})
	if err != nil {
		return nil, err
	}
	var opts []client.Option
	opts = append(opts, client.WithResolver(r))
	opts = append(opts, client.WithLoadBalancer(loadbalance.NewWeightedRandomBalancer()))

	// 设置长连接配置
	cfg := connpool.IdleConfig{
		MaxIdlePerAddress: 10,
		MaxIdleGlobal:     10,
		MaxIdleTimeout:    60 * time.Second,
	}
	opts = append(opts, client.WithLongConnection(cfg))

	//opts = append(opts, client.WithHostPorts("localhost:9999"))
	path, err := idlprovider.NewIdlProvider().GetIdl(rpcName)
	if err != nil {
		return nil, err
	}
	p, err := generic.NewThriftFileProvider(path)
	if err != nil {
		return nil, err
	}
	g, err := generic.JSONThriftGeneric(p)
	if err != nil {
		return nil, err
	}
	cli, _ := genericclient.NewClient(rpcName, g, opts...)
	return &cli, nil
}
