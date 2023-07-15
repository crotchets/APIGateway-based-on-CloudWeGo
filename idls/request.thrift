namespace go apigatewayservice
struct Req{
    1: string ServiceName(api.path="serviceName")
    2: string IDLVersion(api.header="IDLVersion")
    3:string MethodName(api.path="methodName")
}
struct Resp{
    1:string msg
}

service APIService{
    Resp APIPost(1: Req req)(api.post="agw/:serviceName/*methodName")
}