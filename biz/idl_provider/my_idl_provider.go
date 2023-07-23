package idl_provider

type MyIdlProvider struct {
	//todo
}

func NewMyIdlProvider() *MyIdlProvider {
	return &MyIdlProvider{}
}

func (provider *MyIdlProvider) GetIdl(rpcName string) (path string, err error) {
	//todo 这个函数的参数，返回目前都比较随意
	return "./idls/student.thrift", nil
}
