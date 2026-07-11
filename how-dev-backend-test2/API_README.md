# 用户管理API接口文档

## 概述
本项目为how-dev-go-app后端服务，提供用户管理相关的API接口，包括用户登录、用户详情查询等功能。当前版本使用mock数据进行测试，不依赖数据库。

## 接口列表

### 1. 用户登录 (Login)
**接口地址**: `POST /interface`
**Action**: `Login`

**请求参数**:
```json
{
  "RequestId": "test_Login",
  "Action": "Login",
  "username": "admin",
  "password": "admin"
}
```

**响应示例**:
```json
{
  "token": "token_1001_1697270400",
  "userId": 1001,
  "username": "admin",
  "email": "admin@example.com"
}
```

**Mock测试账号**:
- 用户名: `admin`, 密码: `admin`
- 用户名: `user1`, 密码: `hello`  
- 用户名: `test`, 密码: `test`

### 2. 获取用户详情 (GetUserDetail)
**接口地址**: `POST /interface`
**Action**: `GetUserDetail`

**请求参数**:
```json
{
  "RequestId": "test_GetUserDetail",
  "Action": "GetUserDetail",
  "userId": 1001
}
```

**响应示例**:
```json
{
  "userId": 1001,
  "username": "admin",
  "email": "admin@example.com",
  "phone": "13800138001",
  "status": 1,
  "createTime": "2024-01-01 10:00:00",
  "updateTime": "2024-10-14 18:00:00"
}
```

**可用的Mock用户ID**:
- `1001`: admin用户
- `1002`: user1用户
- `1003`: test用户

### 3. 添加用户 (AddUser)
**接口地址**: `POST /interface`
**Action**: `AddUser`

**请求参数**:
```json
{
  "RequestId": "test_AddUser",
  "Action": "AddUser",
  "UserId": 1004,
  "Username": "newuser",
  "Password": "password123",
  "Email": "newuser@example.com",
  "Phone": "13800138004",
  "Status": 1
}
```

### 4. 根据ID获取用户 (GetUserById)
**接口地址**: `POST /interface`
**Action**: `GetUserById`

**请求参数**:
```json
{
  "RequestId": "test_GetUserById",
  "Action": "GetUserById",
  "UserId": 1001
}
```

## 前后端对接说明

### 登录流程
1. 前端发送登录请求到 `/interface` 接口，Action为 `Login`
2. 后端验证用户名密码（当前使用MD5加密验证）
3. 验证成功返回token和用户基本信息
4. 前端保存token用于后续请求认证

### 用户详情查询
1. 前端携带userId发送请求到 `/interface` 接口，Action为 `GetUserDetail`
2. 后端根据userId返回用户详细信息
3. 前端展示用户详情页面

## 测试方法

### 使用HTTP文件测试
项目提供了 `httptest/user.http` 文件，包含所有接口的测试用例：

1. 启动服务: `go run main/main.go`
2. 使用IDE的HTTP Client或Postman导入测试文件
3. 执行各个测试用例验证接口功能

### Mock数据说明
当前所有接口都使用mock数据，不会真正操作数据库：

**登录验证**: 
- 密码使用MD5加密存储和验证
- 预设了3个测试账号供验证

**用户详情**:
- 预设了3个用户的完整信息
- 包含用户基本信息和时间戳

## 注意事项
1. 当前版本为开发测试版本，使用mock数据
2. 密码采用MD5加密（生产环境建议使用更安全的加密方式）
3. Token生成较为简单（生产环境建议使用JWT）
4. 所有接口都通过统一的 `/interface` 入口，通过Action字段区分具体操作