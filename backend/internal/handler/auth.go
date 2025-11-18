package handler  //该包负责HTTP请求处理
//前后端交互流程：
//HTTP请求->handler（HTTP请求处理层）->service（业务逻辑层）->repo（数据访问层）->数据库（前端发起请求流程）
//数据库->repo->service->handler->HTTP响应（后端响应前端流程）
//代码实现时，从底层逐步向上实现，即编码顺序：数据库->repo->service->handler->HTTP响应

//数据库层：负责与数据库交互，执行SQL语句。
//数据访问层（repo）：负责定义数据访问接口，封装数据库操作。
//业务逻辑层（service）：负责实现业务逻辑，调用数据访问层进行数据操作。
//HTTP请求处理层（handler）：负责处理HTTP请求，调用业务逻辑层进行业务处理。
//HTTP响应：负责将业务处理结果返回给前端。


import (
    "net/http"  //Go标准库的HTTP功能

    "se_practice/backend/internal/util"  //项目内部的工具包（统一响应格式）
)

//处理用户注册
func Register(w http.ResponseWriter, r *http.Request) {
    //------------检查是否为Post请求---------------//
    if r.Method != http.MethodPost {
        util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
        return
    }
    //------------检查是否为Post请求---------------//
    //Post请求在传输用户账号、密码时更安全

    util.OK(w, map[string]string{"result": "register stub"})
}

//处理用户登录
func Login(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
        return
    }
    util.OK(w, map[string]string{"token": "stub.jwt.token"})
}


