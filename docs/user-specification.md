# APIGateway-based-on-CloudWeGo用户说明文档
## 小组成员
- [x] 陈皓鑫🌶️ [@crotchets](https://github.com/211250236)
- [x] 张哲恺🌴 [@Corax](https://github.com/KYCoraxxx)
- [x] 张铭铭🍵 [@TTHA](https://github.com/T-THA)
## 项目结构
### 1. 项目目录
以下列出项目主要目录及文件，其中`...`表示省略的目录或文件：
```
├── biz
│    ├── clientprovider
│    │    └── clientprovider.go
│    ├── handler
│    │    └── apigatewayservice
│    │         └── apiservice.go
│    ├── idlmanager
│    │    └── idlmanager.go
│    ├── idlprovider
│    │    └── idlprovider.go
│    ├── model
│    │    └── apigatewayservice
│    │         └── request.go
│    ├── router
│    │    ├── apigatewayservice
│    │    │    ├── middleware.go
│    │    │    └── request.go
│    │    └── register.go
│    └── rpcrouter
│         └── rpcrouter.go
├── go.mod
├── idls
│    ├── request
│    │    └── 1.0.thrift
│    └── student
│         ├── 1.0.thrift
│         └── 1.1.thrift
├── main.go
├── router.go
├── router_gen.go
├── ...
```
### 2. 关键接口及方法描述
#### clientprovider
- 路径：**/biz/handler/clientprovider/clientprovider.go**
- 提供RPC调用的客户端对象，随后进行客户端的泛化调用
- 关键方法：
  - `func (provider *ClientProvider) GetClient(serviceName string, version string) (*genericclient.Client, error)`
    - 功能：获取指定服务的指定版本的`genericclient.Client`对象
    - 参数：
        - `serviceName`：服务名，`string`类型
        - `version`：版本号，`string`类型
    - 返回值：
        - `*genericclient.Client`对象
        - `error`对象
#### apiservice
- 路径：**/biz/handler/apigatewayservice/apiservice.go**
- 处理收到的http请求，并根据业务逻辑进行相应的调用
- 关键方法：
  - `func APIPost(ctx context.Context, c *app.RequestContext)`
    - 路由：`/agw/:serviceName/*methodName [POST]`
    - 功能：处理收到的请求，转化为http请求，并进行RPCRouter转发，返回RPCRouter的调用结果
    - 参数：
        - `ctx`：上下文，`context.Context`类型
        - `c`：请求上下文，`*app.RequestContext`类型
  - `func IDLManage(ctx context.Context, c *app.RequestContext)`
    - 路由：`idl/:IDLName/:IDLVersion [PATCH]`
    - 功能：处理收到的idl相关请求，包括获取、增、删，将请求转发给IDLManager进行处理，并返回处理结果
    - 参数：
        - `ctx`：上下文，`context.Context`类型
        - `c`：请求上下文，`*app.RequestContext`类型
#### idlmanager
- 路径：**/biz/handler/idlmanager/idlmanager.go**
- idl管理的业务核心模块，包括获取、增、删
- 关键方法：
  - `func readIDLFileFromPath(path string) ([]string, error)`
    - 功能：从指定路径读取idl文件内容
    - 参数：
        - `path`：文件路径，`string`类型
    - 返回值：
        - `[]string`类型的idl文件内容
        - `error`对象
  - `func (manager *IDLManager) GetIDL(IDLName string, IDLVersion string) (string, error)`
    - 功能：获取指定服务的指定版本的idl内容
    - 参数：
        - `IDLName`：idl名，`string`类型
        - `IDLVersion`：idl版本号，`string`类型
    - 返回值：
        - `string`类型的idl内容，若不存在指定的idl则返回提示信息
        - `error`对象
  - `func (manager *IDLManager) AddIDL(IDLName string, IDLVersion string, IDLContent string) error`
    - 功能：添加指定服务的指定版本的idl内容
    - 参数：
        - `IDLName`：idl名，`string`类型
        - `IDLVersion`：idl版本号，`string`类型
        - `IDLContent`：idl内容，`string`类型
    - 返回值：
        - `error`对象
  - `func (manager *IDLManager) DeleteIDL(IDLName string, IDLVersion string) error`
    - 功能：删除指定服务的指定版本的idl内容
    - 参数：
        - `IDLName`：idl名，`string`类型
        - `IDLVersion`：idl版本号，`string`类型
    - 返回值：
        - `error`对象
#### rpcrouter
- 路径：**/biz/rpcrouter/rpcrouter.go**
- RPC路由模块，负责将收到的请求转发给RPC Server
- 关键方法：
  - `func (router *RPCRouter) Forward(ctx context.Context, req interface{}, rpcName string, version string, methodName string) (resp interface{}, err error)`
    - 功能：将收到的请求转发给RPC Server
    - 参数：
        - `ctx`：上下文，`context.Context`类型
        - `req`：请求，`interface{}`类型
        - `rpcName`：RPC服务名，`string`类型
        - `version`：版本号，`string`类型
        - `methodName`：方法名，`string`类型
    - 返回值：
        - `resp`：响应，`interface{}`类型
        - `err`：错误，`error`类型
## 部署步骤
确保本地环境中已经安装了`go`和`etcd`，并且已经做好了环境变量配置。
### 1. 准备调用端
调用端即为本项目仓库，可以使用ssh方式或者https方式克隆，选择其中一种执行即可：
```bash
git clone git@github.com:crotchets/APIGateway-based-on-CloudWeGo.git  #SSH用户
git clone https://github.com/crotchets/APIGateway-based-on-CloudWeGo.git  #HTTPS用户
```
可以看到，目录下产生新的文件夹`APIGateway-based-on-CloudWeGo`，进入该文件夹：
```bash
cd APIGateway-based-on-CloudWeGo
```
`main`分支是没有代码的，需要切换分支：
```bash
git checkout -b dev
```
切换分支后，拉取最新代码：
```bash
git pull origin dev
```
更新mod依赖，在目录下执行：
```bash
go mod tidy
```
此外，我们还需要修改idls目录的权限，使其可以在后续测试时进行读写：
```bash
chmod 777 idls/
```
### 2. 准备rpc server
返回上一级目录：
```bash
cd ..
```
克隆rpc server仓库，下列同种类命令均为二选一执行即可：
```bash
git clone git@github.com:KYCoraxxx/rpc-server-for-cloudwego-project.git  #SSH用户
git clone https://github.com/KYCoraxxx/rpc-server-for-cloudwego-project.git  #HTTPS用户
```
可以看到，目录下产生新的文件夹`rpc-server-for-cloudwego-project`，进入该文件夹：
```bash
cd rpc-server-for-cloudwego-project
```
可以使用两种数据存储方式，选择其中一种执行即可，这里通过切换分支选择。
使用公网数据库：
```bash
git checkout -b db-required  #使用公网数据库存储数据
git pull origin db-required
```
使用内存暂存数据：
```bash
git checkout -b local-storage  #使用内存暂存数据
git pull origin local-storage
```
至此，rpc server就准备完毕了。
### 3. 项目运行
首先开启本地etcd服务：
```bash
etcd --log-level debug
```
然后在`APIGateway-based-on-CloudWeGo`目录下运行项目：
```bash
cd APIGateway-based-on-CloudWeGo
go run .
```
然后在`rpc-server-for-cloudwego-project`目录下运行rpc server：
```bash
cd rpc-server-for-cloudwego-project
go run .
```
### 4. 测试
#### IDL Management测试
- 通过`get`方法获取`student`服务的`1.0`版本的`idl`：
    ```bash
    curl -H "Content-Type: application/json" -H "Method: get" -X PATCH http://127.0.0.1:8888/idl/student/1.0
    ```
    执行后应当看到`student`服务的`1.0`版本的`idl`内容
- 通过`get`方法获取`student`服务的`1.1`版本的`idl`：
    ```bash
    curl -H "Content-Type: application/json" -H "Method: get" -X PATCH http://127.0.0.1:8888/idl/student/1.1
    ```
    执行后应当看到`student`服务的`1.1`版本的`idl`内容，由于它之前就是从`1.0`版本的`idl`中复制过来的，所以内容应该相同
- 通过`delete`方法删除`student`服务的`1.1`版本的`idl`：
    ```bash
    curl -H "Content-Type: application/json" -H "Method: delete" -X PATCH http://127.0.0.1:8888/idl/student/1.1
    ```
    执行后应当看到`student`服务的`1.1`版本的`idl`被删除，可以在`APIGateway-based-on-CloudWeGo`的`/idls/student/`目录下查看
- 通过`add`方法添加`student`服务的`1.1`版本的`idl`：
    ```bash
    curl -H "Content-Type: text/plain" -H "Method: add" -T ./idls/student/1.0.thrift -X PATCH http://127.0.0.1:8888/idl/student/1.1
    ```
    执行后应当看到`student`服务的`1.1`版本的`idl`被添加，可以在`APIGateway-based-on-CloudWeGo`的`/idls/student/`目录下查看
- 通过`get`方法获取`student`服务的`1.1`版本的`idl`：
    ```bash
    curl -H "Content-Type: application/json" -H "Method: get" -X PATCH http://127.0.0.1:8888/idl/student/1.1
    ```
    执行后应当看到`student`服务的`1.1`版本的`idl`内容，由于它就是从`1.0`版本的`idl`中复制过来的，所以内容应该相同

#### Student Service测试
- 通过发送`POST`请求，请求注册：
    ```bash
    curl -H "Content-Type: application/json" -H "IDLVersion: 1.0" -X POST http://127.0.0.1:8888/agw/student/Register -d '{"id": 1, "name" : "Xinshen", "college" : {"name": "NJU", "address": "ikuang"}, "email" : ["123456789@qq.com", "211250236@smail.nju.edu.cn"], "sex" : "male"}' 
    ```
    ```bash
    curl -H "Content-Type: application/json" -H "IDLVersion: 1.0" -X POST http://127.0.0.1:8888/agw/student/Register -d '{"id": 2, "name" : "Corax", "college" : {"name": "NJU", "address": "ikun"}, "email" : ["2631197015@qq.com", "211250245@smail.nju.edu.cn"], "sex" : "male"}' 
    ```
    ```bash
    curl -H "Content-Type: application/json" -H "IDLVersion: 1.0" -X POST http://127.0.0.1:8888/agw/student/Register -d '{"id": 3, "name" : "TTHA", "college" : {"name": "NJU", "address": "iming"}, "email" : ["1919810@qq.com", "211252112@smail.nju.edu.cn"], "sex" : "female"}' 
    ```
- 通过发送`POST`请求，请求查询：
    ```bash
    curl -H "Content-Type: application/json" -H "IDLVersion: 1.0" -X POST http://127.0.0.1:8888/agw/student/Query -d '{"id" : 1}'
    ```
    ```bash
    curl -H "Content-Type: application/json" -H "IDLVersion: 1.0" -X POST http://127.0.0.1:8888/agw/student/Query -d '{"id" : 2}'
    ```
    ```bash
    curl -H "Content-Type: application/json" -H "IDLVersion: 1.0" -X POST http://127.0.0.1:8888/agw/student/Query -d '{"id" : 3}'
    ```
    执行后应当可以看到注册的信息。
- 查询之前不存在的学生信息：
    ```bash
    curl -H "Content-Type: application/json" -H "IDLVersion: 1.0" -X POST http://127.0.0.1:8888/agw/student/Query -d '{"id" : 4}'
    ```
    执行后应当可以看到`id`为`4`的学生不存在。
- 重复添加学生信息：
    ```bash
    curl -H "Content-Type: application/json" -H "IDLVersion: 1.0" -X POST http://127.0.0.1:8888/agw/student/Register -d '{"id": 3, "name" : "TTHA", "college" : {"name": "NJU", "address": "iming"}, "email" : ["1919810@qq.com", "211252112@smail.nju.edu.cn"], "sex" : "female"}'
    ```
    执行后应当提示`id`为`3`的学生已经存在。
## 请求接口描述
### 1. IDL Management
- **接口描述**：IDL管理，包括获取IDL、添加IDL、删除IDL等
- **接口地址**：`/idl/{service}/{version}`
- **请求方法**：`PATCH`
- **请求参数**：
    - `service`：服务名，`string`类型，必填
    - `version`：版本号，`string`类型，必填
    - `Method`：请求方法，`string`类型，必填，取值范围为`get`、`add`、`delete`
    - `IDLVersion`：IDL版本号，`string`类型，必填
    - `IDL`：IDL内容，`string`类型，必填
### 2. Student Service
- **接口描述**：学生信息管理，包括注册学生信息、查询学生信息等
- **接口地址**：`/agw/student/{method}`
- **请求方法**：`POST`
- **请求参数**：
    - `method`：方法名，`string`类型，必填，取值范围为`Register`、`Query`
    - `Body`：请求体，`string`类型，必填
      - `id`：学生id，`int`类型
      - `name`：学生姓名，`string`类型
      - `college`：学生学院，`struct`类型
        - `name`：学院名称，`string`类型
        - `address`：学院地址，`string`类型
      - `email`：学生邮箱，`[]string`类型
      - `sex`：学生性别，`string`类型
    - `IDLVersion`：IDL版本号，`string`类型，必填
    - `IDL`：IDL内容，`string`类型，必填


