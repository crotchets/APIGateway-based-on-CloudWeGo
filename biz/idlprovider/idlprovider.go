package idlprovider

import "APIGateway/biz/idlmanager"

type IdlProvider struct {
	//todo: add idl content cache
}

// 切换为单例模式
var idlProvider *IdlProvider

func GetIdlProvider() *IdlProvider {
	if idlProvider == nil {
		idlProvider = &IdlProvider{}
	}
	return idlProvider
}

func (provider *IdlProvider) GetIdl(rpcName string, version string) (content string, err error) {
	return idlmanager.GetManager().GetIDL(rpcName, version)
}
