package idl_mangaer

type IdlInfo struct {
	name string
	hash uint
}
type IdlManager struct {
	m map[string]IdlInfo
}

func AddIDL(name string, idl interface{}) error {
	//todo
	return nil
}

func DelIDL(name string) error {
	//todo
	return nil
}

func getIDL() (interface{}, error) {
	//todo 和idl_provider需要职责划分
	return nil, nil
}
