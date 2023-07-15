namespace go demo

//--------------------request & response--------------
struct College {
    1: required string name(go.tag = 'json:"name"'),
    2: string address(go.tag = 'json:"address"'),
}

struct Student {
    1: required i32 id(api.body="id"),
    2: required string name(api.body="name"),
    3: required College college(api.body="college"),
    4: optional list<string> email(api.body="email"),

}

struct RegisterResp {
    1: bool success,
    2: string message,
}

struct QueryReq {
    1: required i32 id(api.query="id"),
}

struct GetPortReq{
}
struct GetPortResp{
    1: string port
}
//----------------------service-------------------
service StudentService {
    RegisterResp Register(1: Student student)(api.post="/add-student-info")
    Student Query(1: QueryReq req)(api.get="/query")
    GetPortResp GetPort(1: GetPortReq req)(api.get="/port")
}