package model

// AddUserReq
type AddUserReq struct {
	UserId     int64  `json:"UserId"`     // 用户id
	Username   string `json:"Username"`   // 用户名
	Password   string `json:"Password"`   // 密码
	Email      string `json:"Email"`      // 邮箱
	Phone      string `json:"Phone"`      // 手机号
	Status     int    `json:"Status"`     // 状态
	CreateTime string `json:"CreateTime"` // 创建时间
	UpdateTime string `json:"UpdateTime"` // 更新时间
}

// AddUserResp
type AddUserResp struct {
}

type GetUserReq struct {
	UserId int64 `json:"UserId"`
}

type GetUserResp struct {
	UserId     int64  `json:"UserId"`     // 用户id
	Username   string `json:"Username"`   // 用户名
	Password   string `json:"Password"`   // 密码
	Email      string `json:"Email"`      // 邮箱
	Phone      string `json:"Phone"`      // 手机号
	Status     int    `json:"Status"`     // 状态
	CreateTime string `json:"CreateTime"` // 创建时间
	UpdateTime string `json:"UpdateTime"` // 更新时间
}

// LoginReq 登录请求
type LoginReq struct {
	Username string `json:"username" binding:"required"` // 用户名
	Password string `json:"password" binding:"required"` // 密码
}

type LoginData struct {
	Token    string `json:"token"`    // JWT token
	UserId   int64  `json:"userId"`   // 用户ID
	Username string `json:"username"` // 用户名
	Email    string `json:"email"`    // 邮箱
}

// LoginResp 登录响应
type LoginResp struct {
	Success bool      `json:"success"` // 是否成功
	Data    LoginData `json:"data"`    // 登录数据
}

// UserDetailResp 用户详情响应
type UserDetailResp struct {
	UserId     int64  `json:"userId"`     // 用户id
	Username   string `json:"username"`   // 用户名
	Email      string `json:"email"`      // 邮箱
	Phone      string `json:"phone"`      // 手机号
	Status     int    `json:"status"`     // 状态
	CreateTime string `json:"createTime"` // 创建时间
	UpdateTime string `json:"updateTime"` // 更新时间
}
