package idlprovider

type IdlProvider struct {
	//todo
}

func NewIdlProvider() *IdlProvider {
	return &IdlProvider{}
}

func (provider *IdlProvider) GetIdl(rpcName string) (path string, err error) {
	//todo 这个函数的参数，返回目前都比较随意
	return "./idls/student.thrift", nil
}
