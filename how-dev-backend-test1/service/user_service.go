package service

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"how-dev-backend-test1/dao/user_dao"
	"how-dev-backend-test1/model"
	"log"
	"time"
)

type UserService struct {
	UserReadDao  *user_dao.UserReadDao  `wired:"true"`
	UserWriteDao *user_dao.UserWriteDao `wired:"true"`
}

func (u *UserService) AddUser(ctx context.Context, user model.AddUserReq) (result bool, err error) {
	createUser := &user_dao.TableUser{
		UserId:     user.UserId,
		Username:   user.Username,
		Password:   user.Password,
		Email:      user.Email,
		Phone:      user.Phone,
		Status:     user.Status,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	_, err = u.UserWriteDao.Create(ctx, createUser)
	if err != nil {
		log.Printf("add user error: %v", err)
		return false, err
	}
	return true, nil
}

func (u *UserService) GetUserByUserId(ctx context.Context, req model.GetUserReq) (model.GetUserResp, error) {
	var userRsp model.GetUserResp
	user, err := u.UserReadDao.GetUserByUserId(ctx, req.UserId)
	if err != nil {
		log.Printf("get user by id error: %v", err)
		return userRsp, err
	}
	userRsp.UserId = user.UserId
	userRsp.Username = user.Username
	userRsp.Email = user.Email
	userRsp.Phone = user.Phone
	userRsp.Status = user.Status
	userRsp.CreateTime = user.CreateTime.Format("2006-01-02 15:04:05")
	userRsp.UpdateTime = user.UpdateTime.Format("2006-01-02 15:04:05")
	return userRsp, nil
}

// Login 用户登录 - 使用mock数据
func (u *UserService) Login(ctx context.Context, req model.LoginReq) (model.LoginResp, error) {
	var loginResp model.LoginResp

	// Mock用户数据
	mockUsers := []struct {
		UserId   int64
		Username string
		Password string
		Email    string
	}{
		{1001, "admin", "21232f297a57a5a743894a0e4a801fc3", "admin@example.com"}, // admin的MD5
		{1002, "user1", "5d41402abc4b2a76b9719d911017c592", "user1@example.com"}, // hello的MD5
		{1003, "test", "098f6bcd4621d373cade4e832627b4f6", "test@example.com"},   // test的MD5
	}

	// 对输入密码进行MD5加密
	inputPasswordMD5 := fmt.Sprintf("%x", md5.Sum([]byte(req.Password)))

	// 查找匹配的用户
	var matchedUser *struct {
		UserId   int64
		Username string
		Password string
		Email    string
	}

	for _, user := range mockUsers {
		if user.Username == req.Username && user.Password == inputPasswordMD5 {
			matchedUser = &user
			break
		}
	}

	if matchedUser == nil {
		return loginResp, errors.New("用户名或密码错误")
	}

	// 生成简单的token (实际项目中应该使用JWT)
	token := fmt.Sprintf("token_%d_%d", matchedUser.UserId, time.Now().Unix())

	var loginData model.LoginData
	loginData.Token = token
	loginData.UserId = matchedUser.UserId
	loginData.Username = matchedUser.Username
	loginData.Email = matchedUser.Email

	loginResp.Success = true
	loginResp.Data = loginData

	log.Printf("用户登录成功: %s", req.Username)
	return loginResp, nil
}

// GetUserDetail 获取用户详情 - 使用mock数据
func (u *UserService) GetUserDetail(ctx context.Context, userId int64) (model.UserDetailResp, error) {
	var userDetail model.UserDetailResp

	// Mock用户详情数据
	mockUserDetails := map[int64]model.UserDetailResp{
		1001: {
			UserId:     1001,
			Username:   "admin",
			Email:      "admin@example.com",
			Phone:      "13800138001",
			Status:     1,
			CreateTime: "2024-01-01 10:00:00",
			UpdateTime: "2024-10-14 18:00:00",
		},
		1002: {
			UserId:     1002,
			Username:   "user1",
			Email:      "user1@example.com",
			Phone:      "13800138002",
			Status:     1,
			CreateTime: "2024-02-01 10:00:00",
			UpdateTime: "2024-10-14 17:30:00",
		},
		1003: {
			UserId:     1003,
			Username:   "test",
			Email:      "test@example.com",
			Phone:      "13800138003",
			Status:     1,
			CreateTime: "2024-03-01 10:00:00",
			UpdateTime: "2024-10-14 17:00:00",
		},
	}

	detail, exists := mockUserDetails[userId]
	if !exists {
		return userDetail, errors.New("用户不存在")
	}

	log.Printf("获取用户详情成功: userId=%d", userId)
	return detail, nil
}
