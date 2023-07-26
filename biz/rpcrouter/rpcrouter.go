package rpcrouter

import (
	"APIGateway/biz/clientprovider"
	"context"
)

// 接收请求，将请求转发给rpc服务器，并将回复返回

type RPCRouter struct{}

var rpcrouter *RPCRouter

func GetRPCRouter() *RPCRouter {
	if rpcrouter == nil {
		rpcrouter = &RPCRouter{}
	}
	return rpcrouter
}
func (router *RPCRouter) Forward(ctx context.Context, req interface{}, rpcName string, version string, methodName string) (resp interface{}, err error) {
	client, err := clientprovider.GetClientProvider().GetClient(rpcName, version) // 获取RPC客户端
	if err != nil {
		return "", err
	}
	resp, err = (*client).GenericCall(ctx, methodName, req) // 泛化调用并返回结果
	if err != nil {
		return "", err
	}
	return
}
