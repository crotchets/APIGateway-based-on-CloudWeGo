package idlmanager

import (
	"crypto/sha256"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type IdlInfo struct {
	name string
	hash string
}
type IdlManager struct {
	m map[string]IdlInfo
}

var manager *IdlManager

const idlRootDirectory string = "./idls/"
const idlFileSuffix string = ".thrift"

func readIDLFileFromPath(path string) ([]string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var res []string

	for _, file := range files {
		if file.IsDir() {
			subFiles, err := readIDLFileFromPath(path + file.Name())
			if err != nil {
				return nil, err
			}

			for _, version := range subFiles {
				res = append(res, file.Name()+version)
			}
		} else {
			res = append(res, file.Name())
		}
	}
	return res, nil
}

func GetManager() *IdlManager {
	if manager == nil {
		manager = &IdlManager{make(map[string]IdlInfo)}

		files, err := readIDLFileFromPath(idlRootDirectory)
		if err != nil {
			log.Fatal(err)
		} else {
			hash := sha256.New()
			for _, file := range files {
				filename := strings.TrimSuffix(file, idlFileSuffix)
				manager.m[filename] = IdlInfo{name: file, hash: string(hash.Sum(nil))}
			}
		}
	}
	return manager
}
func AddIDL(name string, version string, idl interface{}) error {
	if _, exist := (GetManager().m)[name+version]; exist {
		return errors.New("IDL exists")
	}
	newFile, err := os.Create("./idls/" + name + "/" + version + ".thrift")
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
	GetManager().m[name+version] = IdlInfo{name, string(hash.Sum(nil))}
	return nil
}

func DelIDL(name string, version string) error {
	if _, exist := (GetManager().m)[name+version]; !exist {
		return errors.New("no such IDL")
	}
	if err := os.Remove("./idls/" + name + "/" + version + ".thrift"); err != nil {
		return err
	}
	delete(GetManager().m, name+version)
	return nil
}

func getIDL() (interface{}, error) {
	//todo 和idl_provider需要职责划分

	// idl provider 返回路径， getIDL返回文件内容
	return nil, nil
}
