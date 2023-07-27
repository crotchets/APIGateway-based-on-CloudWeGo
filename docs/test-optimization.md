# APIGateway-based-on-CloudWeGo性能测试和优化报告
## 小组成员
- [x] 陈皓鑫🌶️ [@crotchets](https://github.com/211250236)
- [x] 张哲恺🌴 [@Corax](https://github.com/KYCoraxxx)
- [x] 张铭铭🍵 [@TTHA](https://github.com/T-THA)
## 测试方案说明

本小组实现的API网关一共实现了两个功能，如下所示：

- 通过HTTP实现IDL的多版本控制与管理
- 接收HTTP请求并转发给RPC服务器完成对学生信息管理的处理

### IDL多版本控制与管理

IDL多版本控制与管理由HTTP服务器自行管理，不需要转发请求，因此主要进行对该模块增、删、查的功能及面对高并发压力的能力进行测试

### API网关学生信息管理

RPC服务端一共实现了两种数据管理方式，分别是暂存于运行内存和存储于公网服务器postgresql数据库。对于两者而言，仍然主要对其增、查及面对高并发压力的能力进行测试，其中对于存储于公网服务器数据库的RPC服务器，还要额外考虑在性能资源有限的情况下如何面对高压力访问

## IDL多版本控制于管理

### 高压力连续查询

#### 性能测试数据

由于ab不支持PATCH请求，因此我们写了一个简易的发压脚本，如下所示
```go
package main

import (
	"fmt"
	"os/exec"
)

const REQUEST_TIMES int = 100000
const COWORK_NUMS int = 10

var c = make(chan int, REQUEST_TIMES)
var args []string

func test(id int) {
	for len(c) < REQUEST_TIMES {
		c <- id
		cmd := exec.Command("curl", args...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(out[:]))
	}
}

func main() {
	args = append(args, "-H")
	args = append(args, "Method: get")
	args = append(args, "-X")
	args = append(args, "PATCH")
	args = append(args, "http://127.0.0.1:8888/idl/student/1.0")
	for i := 1; i <= COWORK_NUMS; i++ {
		go test(i)
	}
	for len(c) <= REQUEST_TIMES {

	}
}
```
1. 返回数据
```text
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100   874  100   874    0     0  1203k      0 --:--:-- --:--:-- --:--:--  853k
{"msg":"namespace go demo\r\n\r\n//--------------------request \u0026 response--------------\r\nstruct College {\r\n    1: required string name(go.tag = 'json:\"name\"'),\r\n    2: string address(go.ta
g = 'json:\"address\"'),\r\n}\r\n\r\nstruct Student {\r\n    1: required i32 id,\r\n    2: required string name,\r\n    3: required College college,\r\n    4: optional list\u003cstring\u003e email,\r\n
    5: optional string sex,\r\n}\r\n\r\nstruct RegisterResp {\r\n    1: bool success,\r\n    2: string message,\r\n}\r\n\r\nstruct QueryReq {\r\n    1: required i32 id,\r\n}\r\nstruct GetPortReq{\r\n}\
r\nstruct GetPortResp{\r\n    1: string port\r\n}\r\n//----------------------service-------------------\r\nservice StudentService {\r\n    RegisterResp Register(1: Student student)\r\n    Student Query(1: QueryReq req)\r\n    GetPortResp GetPort(1: GetPortReq req)\r\n}"}

......
```
2. pprof监测结果
![](img/optimization/idl.png)
![](img/optimization/top.png)

#### 性能优化方案

从pprof的监测结果来看，os.ReadFile接口占用了相当的时间，考察此处代码的写法:

```go
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
```

因此优化方案为为读文件添加缓存，并定时更新文件内容:

```go
package idlmanager

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

type IdlInfo struct {
	name    string // 对应目录文件名
	content string // 文件内容
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

var flag bool = false

func GetManager() *IdlManager {
	if manager == nil {
		manager = &IdlManager{make(map[string]IdlInfo)}
		go manager.update()
		for !flag {
		} // wait for the first updating
	}
	return manager
}

func (man *IdlManager) update() {
	for {
		files, err := readIDLFileFromPath(idlRootDirectory)
		if err != nil {
			log.Fatal(err)
		} else {
			for _, file := range files {
				ct, err := ioutil.ReadFile(idlRootDirectory + file)
				if err != nil {
					log.Fatal(err)
				}
				manager.m[getFileName(file)] = IdlInfo{name: file, content: string(ct[:])}
			}
		}
		flag = true
		time.Sleep(30 * time.Second)
	}
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

	man.m[targetFile] = IdlInfo{filename, idl.(string)}
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
	} else {
		return man.m[targetFile].content, nil
	}
}
```

#### 优化后性能数据

![](img/optimization/idl-opt.png)
![](img/optimization/top-opt.png)

从pprof监测的结果来看，IDLManage模块原本在读取文件上就要花费的0.69s，在缓存的帮助下，整个模块耗时降低到了0.24s，性能得到了提升

## API网关学生信息管理

### 较高压力连续查询

#### 性能测试数据

```bash
ab -n 1000 -c 10 -H 'IDLVersion: 1.0' -T 'application/json' -p data.json http://127.0.0.1:8888/agw/student/Query
```
查询中使用的`data.json`文件如下所示
```json
{
  "id" : 1
}
```
测试结果如下：
1. 返回结果
```text
This is ApacheBench, Version 2.3 <$Revision: 1879490 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking 127.0.0.1 (be patient)
Completed 100 requests
Completed 200 requests
Completed 300 requests
Completed 400 requests
Completed 500 requests
Completed 600 requests
Completed 700 requests
Completed 800 requests
Completed 900 requests
Completed 1000 requests
Finished 1000 requests


Server Software:        hertz
Server Hostname:        127.0.0.1
Server Port:            8888

Document Path:          /agw/student/Query
Document Length:        0 bytes

Concurrency Level:      10
Time taken for tests:   0.413 seconds
Complete requests:      1000
Failed requests:        0
Total transferred:      109000 bytes
Total body sent:        187000
HTML transferred:       0 bytes
Requests per second:    2423.98 [#/sec] (mean)
Time per request:       4.125 [ms] (mean)
Time per request:       0.413 [ms] (mean, across all concurrent requests)
Transfer rate:          258.02 [Kbytes/sec] received
                        442.66 kb/s sent
                        700.68 kb/s total

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.1      0       0
Processing:     0    3   4.2      2     126
Waiting:        0    3   4.2      2     126
Total:          0    3   4.2      3     126

Percentage of the requests served within a certain time (ms)
  50%      3
  66%      3
  75%      4
  80%      4
  90%      5
  95%      6
  98%      7
  99%     10
 100%    126 (longest request)
```
2. RPC端情况
```text
2023/07/26 11:34:22 dial tcp 110.42.252.167:5432: connect: connection refused
exit status 1
```
在面对较高压力的情况下，性能一般的公网服务器数据库进程被打挂了❌

#### 性能优化方案

为RPC服务端添加数据查询缓存，查询一次之后就将查询的数据存入缓存，之后的相同查询就直接从缓存中取数据，不再访问数据库

```go
func (s *StudentServiceImpl) Query(ctx context.Context, req *demo.QueryReq) (resp *demo.Student, err error) {
	resp = demo.NewStudent()
	var oldStudent demo.Student
	if value, exist := studentMap[req.Id]; exist {
		fmt.Println("Use Cache")
		resp = value
		return
	} else {
		fmt.Println("Query Database")
		err = QueryFromDatabase(req.Id, &oldStudent)
		if err != nil {
			return
		}
		if oldStudent.Id == -1 {
			var student = demo.Student{
				Id:      -1,
				Name:    "Student Not Exist",
				College: &demo.College{Name: "Unknown", Address: "Unknown"},
				Email:   nil,
			}
			resp = &student
		} else {
			resp = &oldStudent
			studentMap[req.Id] = &oldStudent
		}
		return
	}
}
```

同样也为注册方法添加缓存

```go
func (s *StudentServiceImpl) Register(ctx context.Context, student *demo.Student) (resp *demo.RegisterResp, err error) {
	resp = demo.NewRegisterResp()
	var newStudent demo.Student
	if _, exist := studentMap[student.Id]; exist {
		resp.Success = false
		resp.Message = "Register Failed: Student Information Already Exists"
	} else {
		err = QueryFromDatabase(student.Id, &newStudent)
		if err != nil {
			resp.Success = false
			resp.Message = "Internal Exception"
		}
		if newStudent.Id > 0 {
			studentMap[student.Id] = &newStudent
			resp.Success = false
			resp.Message = "Register Failed: Student Information Already Exists"
		} else {
			err = InsertIntoDatabase(student)
			if err != nil {
				resp.Success = false
				resp.Message = "Internal Exception"
			}
			resp.Success = true
			resp.Message = "Register Success"
		}
		fmt.Println(resp)
	}
	return
}
```

#### 优化后性能数据

1. 返回数据

```text
This is ApacheBench, Version 2.3 <$Revision: 1879490 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking 127.0.0.1 (be patient)
Completed 100 requests
Completed 200 requests
Completed 300 requests
Completed 400 requests
Completed 500 requests
Completed 600 requests
Completed 700 requests
Completed 800 requests
Completed 900 requests
Completed 1000 requests
Finished 1000 requests


Server Software:        hertz
Server Hostname:        127.0.0.1
Server Port:            8888

Document Path:          /agw/student/Query
Document Length:        169 bytes

Concurrency Level:      10
Time taken for tests:   0.303 seconds
Complete requests:      1000
Failed requests:        0
Total transferred:      327000 bytes
Total body sent:        187000
HTML transferred:       169000 bytes
Requests per second:    3299.42 [#/sec] (mean)
Time per request:       3.031 [ms] (mean)
Time per request:       0.303 [ms] (mean, across all concurrent requests)
Transfer rate:          1053.62 [Kbytes/sec] received
                        602.53 kb/s sent
                        1656.15 kb/s total

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.1      0       1
Processing:     0    1   7.9      0     250
Waiting:        0    1   7.9      0     250
Total:          0    1   7.9      0     250

Percentage of the requests served within a certain time (ms)
  50%      0
  66%      1
  75%      1
  80%      1
  90%      1
  95%      1
  98%      1
  99%      2
 100%    250 (longest request)
```

2. RPC端情况

```text
Query Database
Use Cache
Use Cache
......
```

可见性能得到明显提升，而且公网服务器存活了下来😊

### 高压力连续查询

#### 性能测试数据

将测试命令更改为
```bash
ab -n 100000 -c 10 -H 'IDLVersion: 1.0' -T 'application/json' -p data.json http://127.0.0.1:8888/agw/student/Query
```
1. 返回数据
```text
This is ApacheBench, Version 2.3 <$Revision: 1879490 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking 127.0.0.1 (be patient)
Completed 10000 requests
Completed 20000 requests
Completed 30000 requests
Completed 40000 requests
Completed 50000 requests
Completed 60000 requests
Completed 70000 requests
Completed 80000 requests
Completed 90000 requests
Completed 100000 requests
Finished 100000 requests


Server Software:        hertz
Server Hostname:        127.0.0.1
Server Port:            8888

Document Path:          /agw/student/Query
Document Length:        169 bytes

Concurrency Level:      10
Time taken for tests:   5.994 seconds
Complete requests:      100000
Failed requests:        0
Total transferred:      32700000 bytes
Total body sent:        18700000
HTML transferred:       16900000 bytes
Requests per second:    16684.15 [#/sec] (mean)
Time per request:       0.599 [ms] (mean)
Time per request:       0.060 [ms] (mean, across all concurrent requests)
Transfer rate:          5327.85 [Kbytes/sec] received
                        3046.81 kb/s sent
                        8374.66 kb/s total

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.1      0       1
Processing:     0    1   0.2      0       6
Waiting:        0    0   0.2      0       6
Total:          0    1   0.2      1       6
ERROR: The median and mean for the processing time are more than twice the standard
       deviation apart. These results are NOT reliable.

Percentage of the requests served within a certain time (ms)
  50%      1
  66%      1
  75%      1
  80%      1
  90%      1
  95%      1
  98%      1
  99%      1
 100%      6 (longest request)
```
2. pprof监测结果

![](img/optimization/query.png)

#### 性能评估

根据性能测试数据来看，在经过数据缓存优化之后，系统的性能已经基本满足期望要求

### 高压力连续注册

#### 性能测试数据

测试命令为:
```bash
ab -n 100000 -c 10 -H 'IDLVersion: 1.0' -T 'application/json' -p data.json http://127.0.0.1:8888/agw/student/Register
```
使用的`data.json`文件如下所示:
```json
{
    "id": 4,
    "name" : "KFC",
    "college" : {"name": "KFC", "address": "Thursday"},
    "email" : ["2631197015@qq.com", "211250245@smail.nju.edu.cn"],
    "sex" : "female"
}
```

1. 返回数据
```text
This is ApacheBench, Version 2.3 <$Revision: 1879490 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking 127.0.0.1 (be patient)
Completed 10000 requests
Completed 20000 requests
Completed 30000 requests
Completed 40000 requests
Completed 50000 requests
Completed 60000 requests
Completed 70000 requests
Completed 80000 requests
Completed 90000 requests
Completed 100000 requests
Finished 100000 requests


Server Software:        hertz
Server Hostname:        127.0.0.1
Server Port:            8888

Document Path:          /agw/student/Register
Document Length:        53 bytes

Concurrency Level:      10
Time taken for tests:   5.934 seconds
Complete requests:      100000
Failed requests:        99999
   (Connect: 0, Receive: 0, Length: 99999, Exceptions: 0)
Total transferred:      24599964 bytes
Total body sent:        35900000
HTML transferred:       8899964 bytes
Requests per second:    16851.88 [#/sec] (mean)
Time per request:       0.593 [ms] (mean)
Time per request:       0.059 [ms] (mean, across all concurrent requests)
Transfer rate:          4048.39 [Kbytes/sec] received
                        5908.03 kb/s sent
                        9956.43 kb/s total

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.1      0       1
Processing:     0    0   3.8      0    1016
Waiting:        0    0   3.8      0    1016
Total:          0    0   3.8      0    1016

Percentage of the requests served within a certain time (ms)
  50%      0
  66%      0
  75%      1
  80%      1
  90%      1
  95%      1
  98%      1
  99%      1
 100%   1016 (longest request)
```
2. pprof监测结果

![](img/optimization/insert.png)

#### 性能评估

根据性能测试数据来看，在经过数据缓存优化之后，系统的性能已经基本满足期望要求