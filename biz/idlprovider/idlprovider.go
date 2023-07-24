package idlprovider

type IdlProvider struct {
	//todo
}

// 切换为单例模式
var idlProvider *IdlProvider

func GetIdlProvider() *IdlProvider {
	if idlProvider == nil {
		idlProvider = &IdlProvider{}
	}
	return idlProvider
}

func (provider *IdlProvider) GetIdl(rpcName string) (path string, err error) {
	//todo 这个函数的参数，返回目前都比较随意
	target := "./idls/" + rpcName + ".thrift"
	return target, nil
}
