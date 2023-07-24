package idl_manager

import (
	"crypto/sha256"
	"errors"
	"io"
	"os"
)

type IdlInfo struct {
	name string
	hash string
}
type IdlManager struct {
	m map[string]IdlInfo
}

var manager *IdlManager

func GetManager() *IdlManager {
	if manager == nil {
		manager = &IdlManager{make(map[string]IdlInfo)}
	}
	return manager
}
func AddIDL(name string, idl interface{}) error {
	if _, exist := (GetManager().m)[name]; exist {
		return errors.New("IDL exists")
	}
	newFile, err := os.Create("./idls/" + name + ".thrift")
	if err != nil {
		return err
	}
	defer newFile.Close()
	if _, err = newFile.WriteString(idl.(string)); err != nil {
		return err
	}
	hash := sha256.New()
	if _, err = io.Copy(hash, newFile); err != nil {
		return err
	}
	GetManager().m[name] = IdlInfo{name, string(hash.Sum(nil))}
	return nil
}

func DelIDL(name string) error {
	if _, exist := (GetManager().m)[name]; !exist {
		return errors.New("no such IDL")
	}
	if err := os.Remove("./idls/" + name + ".thrift"); err != nil {
		return err
	}
	delete(GetManager().m, name)
	return nil
}

func getIDL() (interface{}, error) {
	//todo 和idl_provider需要职责划分
	return nil, nil
}
