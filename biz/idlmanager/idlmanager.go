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
	name string // 对应目录文件名
	hash string // 暂时无用
}
type IdlManager struct {
	m map[string]IdlInfo // 建立从名称+版本到对应目录文件的映射
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
				res = append(res, file.Name()+"/"+version)
			}
		} else {
			res = append(res, file.Name())
		}
	}
	return res, nil
}

func getFileName(rawName string) string {
	noSuffixName := strings.TrimSuffix(rawName, idlFileSuffix)
	parts := strings.Split(noSuffixName, "/")
	return strings.Join(parts, "")
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
				manager.m[getFileName(file)] = IdlInfo{name: file, hash: string(hash.Sum(nil))}
			}
		}
	}
	return manager
}

func (man *IdlManager) AddIDL(name string, version string, idl interface{}) error {
	targetFile := name + version
	if _, exist := man.m[targetFile]; exist {
		return errors.New("IDL already exists")
	}
	filename := name + "/" + version + idlFileSuffix
	newFile, err := os.Create(idlRootDirectory + filename)
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

	man.m[targetFile] = IdlInfo{filename, string(hash.Sum(nil))}
	return nil
}

func (man *IdlManager) DelIDL(name string, version string) error {
	targetFile := name + version
	if _, exist := man.m[targetFile]; !exist {
		return errors.New("no such IDL")
	}
	if err := os.Remove(idlRootDirectory + man.m[targetFile].name); err != nil {
		return err
	}

	delete(man.m, targetFile)
	return nil
}

func (man *IdlManager) GetIDL(name string, version string) (string, error) {
	targetFile := name + version
	if _, exist := man.m[targetFile]; !exist {
		return "", errors.New("no such IDL")
	}
	if file, err := ioutil.ReadFile(idlRootDirectory + man.m[targetFile].name); err != nil {
		return "", err
	} else {
		return string(file[:]), nil
	}
}
