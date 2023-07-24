namespace go apigatewayservice
struct Req{
    1: string ServiceName(api.path="serviceName")
    2: string IDLVersion(api.header="IDLVersion")
    3:string MethodName(api.path="methodName")
}
struct IDLManageReq{
    1: string IDLName(api.path="IDLName")
    2: string IDLVersion(api.path="IDLVersion")
    3: string method(api.header="Method")
}
struct Resp{
    1:string msg
}

service APIService{
    Resp APIPost(1: Req req)(api.post="agw/:serviceName/*methodName")
    Resp IDLManage(1: IDLManageReq req)(api.patch="idl/:IDLName/:IDLVersion")
}