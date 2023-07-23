package client_provider

import (
	"APIGateway/biz/idl_provider"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/loadbalance"
	etcd "github.com/kitex-contrib/registry-etcd"
)

type MyProvider struct {
	//todo
}

func NewMyProvider() *MyProvider {
	return &MyProvider{}
}

func (provider *MyProvider) GetClient(rpcName string) (*genericclient.Client, error) {
	r, err := etcd.NewEtcdResolver([]string{"127.0.0.1:2379"})
	if err != nil {
		return nil, err
	}
	var opts []client.Option
	opts = append(opts, client.WithResolver(r))
	opts = append(opts, client.WithLoadBalancer(loadbalance.NewWeightedRandomBalancer()))

	//opts = append(opts, client.WithHostPorts("localhost:9999"))
	path, err := idl_provider.NewMyIdlProvider().GetIdl(rpcName)
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
