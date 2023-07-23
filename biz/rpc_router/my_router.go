package rpc_router

import (
	"APIGateway/biz/client_provider"
	"context"
)

type MyRouter struct {
	//todo
}

func NewMyRouter() *MyRouter {
	return &MyRouter{}
}
func (router MyRouter) Forward(ctx context.Context, req interface{}, rpcName string, methodName string) (resp interface{}, err error) {
	//todo 将请求转发给rpc服务器，并将回复返回
	client, err := client_provider.NewMyProvider().GetClient(rpcName)
	if err != nil {
		return "", err
	}
	resp, err = (*client).GenericCall(ctx, methodName, req)
	if err != nil {
		return "", err
	}
	return
}
