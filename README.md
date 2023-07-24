# APIGateway-based-on-CloudWeGo
> Project of CloudWeGo in NJU SE
## Contributor
- [x] 陈皓鑫🌶️ [@crotchets](https://github.com/211250236)
- [x] 张哲恺🌴 [@Corax](https://github.com/KYCoraxxx)
- [x] 张铭铭🍵 [@TTHA](https://github.com/T-THA)
## User Documentation
🎉🎉🎉 请前往docs目录查看详细用户文档
## Developer Documentation
⚠️⚠️⚠️ 下面的文档仅为开发者指引，用户请勿参考
### Quick Start
#### 1. Clone the repository
```bash
git clone git@github.com:crotchets/APIGateway-based-on-CloudWeGo.git  #SSH用户
git clone https://github.com/crotchets/APIGateway-based-on-CloudWeGo.git  #HTTPS用户
```
#### 2. Checkout to developer branch
```bash
git checkout -b demo  #我也不知道为啥叫demo
git pull origin demo
```
#### 3. Run the project
```bash
cd APIGateway-based-on-CloudWeGo
go run .
```
#### 4. Test the project
##### IDL Management
```bash
curl -H "Content-Type: application/json" -H "Method: get" -X PATCH http://127.0.0.1:8888/idl/student/1.0
curl -H "Content-Type: application/json" -H "Method: get" -X PATCH http://127.0.0.1:8888/idl/student/1.1
curl -H "Content-Type: application/json" -H "Method: delete" -X PATCH http://127.0.0.1:8888/idl/student/1.1
curl -H "Content-Type: text/plain" -H "Method: add" -T ./idls/student/1.0.thrift -X PATCH http://127.0.0.1:8888/idl/student/1.1
curl -H "Content-Type: application/json" -H "Method: get" -X PATCH http://127.0.0.1:8888/idl/student/1.1
```
上面的五条指令分别表示:
1. 通过`get`方法获取`student`服务的`1.0`版本的`idl`
2. 通过`get`方法获取`student`服务的`1.1`版本的`idl`
3. 通过`delete`方法删除`student`服务的`1.1`版本的`idl`
4. 通过`add`方法添加`student`服务的`1.1`版本的`idl`
5. 通过`get`方法获取`student`服务的`1.1`版本的`idl`

执行上述指令后你应该看到的结果:
1. `student`服务的`1.0`版本的`idl`内容
2. `student`服务的`1.1`版本的`idl`内容，由于它之前就是从`1.0`版本的`idl`中复制过来的，所以内容应该相同
3. `student`服务的`1.1`版本的`idl`被删除，你可以在项目结构的/idls/student/目录下查看
4. `student`服务的`1.1`版本的`idl`被添加，你可以在项目结构的/idls/student/目录下查看
5. `student`服务的`1.1`版本的`idl`内容，由于它就是从`1.0`版本的`idl`中复制过来的，所以内容应该相同

##### Student Service
尚未进行接入测试，所以暂时无法提供测试指令

>TODO: 挖个坑在这

### Project Structure

>此部分并非严格的项目结构说明，仅是引导开发者快速掌握整体结构，请按照先后顺序阅读
 
**1. /biz/handler**

apiservice.go处理收到的http请求，并根据业务逻辑进行相应的调用，详细细节参见代码文件注释

#### RPC Server

**2. /biz/rpcrouter**

APIPost方法请求rpcrouter.go进行转发，详细细节参见代码文件注释

**3. /biz/clientprovider**

Forward方法请求clientprovider.go提供RPC调用的客户端，随后进行客户端的泛化调用，并将结果返回，详细细节参见代码文件注释

#### IDL Manager

**2. /biz/idlmanager**

IDLManage方法根据请求标头中的参数调用idlmanager.go中的对应方法，详细细节参见代码文件注释



