# APIGateway-based-on-CloudWeGo开发者记录文档
## 小组成员
- [x] 陈皓鑫🌶️ [@crotchets](https://github.com/211250236)
- [x] 张哲恺🌴 [@Corax](https://github.com/KYCoraxxx)
- [x] 张铭铭🍵 [@TTHA](https://github.com/T-THA)
## Quick Start
⚠️⚠️⚠️ 下面的文档仅为开发者指引，用户请勿参考 ⚠️⚠️⚠️
### 1. Clone the repository
```bash
git clone git@github.com:crotchets/APIGateway-based-on-CloudWeGo.git  #SSH用户
git clone https://github.com/crotchets/APIGateway-based-on-CloudWeGo.git  #HTTPS用户
```
### 2. Checkout to developer branch
```bash
cd APIGateway-based-on-CloudWeGo
git checkout -b dev
git pull origin dev
```
### 3. Run the project
```bash
go run .
```
### 4. Test the project
#### IDL Management
```bash
curl -H "Content-Type: application/json" -H "Method: get" -X PATCH http://127.0.0.1:8888/idl/student/1.0
curl -H "Content-Type: application/json" -H "Method: get" -X PATCH http://127.0.0.1:8888/idl/student/1.1
curl -H "Content-Type: application/json" -H "Method: delete" -X PATCH http://127.0.0.1:8888/idl/student/1.1
curl -H "Content-Type: text/plain" -H "Method: add" -T ./idls/student/1.0.thrift -X PATCH http://127.0.0.1:8888/idl/student/1.1
curl -H "Content-Type: application/json" -H "Method: get" -X PATCH http://127.0.0.1:8888/idl/student/1.1
```

#### Student Service
1. 运行etcd程序
```bash
etcd --log-level debug
```
2. 运行rpc server

下列同种类命令均为二选一执行即可
```bash
 git clone git@github.com:KYCoraxxx/rpc-server-for-cloudwego-project.git  #SSH用户
 git clone https://github.com/KYCoraxxx/rpc-server-for-cloudwego-project.git  #HTTPS用户
 
 cd rpc-server-for-cloudwego-project
 
 git checkout db-required #使用公网数据库存储数据
 git checkout local-storage #使用内存暂存数据
 
 go run .
```
3. 运行测试指令
```bash 
curl -H "Content-Type: application/json" -H "IDLVersion: 1.0" -X POST http://127.0.0.1:8888/agw/student/Register -d '{"id": 1, "name" : "Xinshen", "college" : {"name": "NJU", "address": "ikuan g"}, "email" : ["2631197015@qq.com", "211250245@smail.nju.edu.cn"], "sex" : "male"}' 
curl -H "Content-Type: application/vbjson" -H "IDLVersion: 1.0" -X POST http://127.0.0.1:8888/agw/student/Register -d '{"id": 2, "name" : "Corax", "college" : {"name": "NJU", "address": "ikuan g"}, "email" : ["2631197015@qq.com", "211250245@smail.nju.edu.cn"], "sex" : "male"}' 
curl -H "Content-Type: application/json" -H "IDLVersion: 1.0" -X POST http://127.0.0.1:8888/agw/student/Query -d '{"id" : 1}' 
curl -H "Content-Type: application/json" -H "IDLVersion: 1.0" -X POST http://127.0.0.1:8888/agw/student/Query -d '{"id" : 2}' 
```

## Project Structure

>此部分并非严格的项目结构说明，仅是引导开发者快速掌握整体结构，请按照先后顺序阅读

**1. /biz/handler**

apiservice.go处理收到的http请求，并根据业务逻辑进行相应的调用，详细细节参见代码文件注释

### RPC Server

**2. /biz/rpcrouter**

APIPost方法请求rpcrouter.go进行转发，详细细节参见代码文件注释

**3. /biz/clientprovider**

Forward方法请求clientprovider.go提供RPC调用的客户端，随后进行客户端的泛化调用，并将结果返回，详细细节参见代码文件注释

### IDL Manager

**2. /biz/idlmanager**

IDLManage方法根据请求标头中的参数调用idlmanager.go中的对应方法，详细细节参见代码文件注释
