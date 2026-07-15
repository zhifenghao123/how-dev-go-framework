package controller

import (
	"context"
	"how-dev-backend-test2/model"
	"how-dev-backend-test2/service"
	"how-dev-backend-test2/tools"
	"net/http"
)

type UserController struct {
	UserService *service.UserService `wired:"true"`
}

func (u *UserController) AddUser(ctx context.Context, req *http.Request) (rsp any, err error) {
	var params model.AddUserReq
	if err = tools.RequestParam(&params, req); err != nil {
		println("parse params error:", err.Error())
		return
	}
	result, err := u.UserService.AddUser(ctx, params)
	if err != nil {
		return nil, err
	}
	rsp = result
	return rsp, nil
}

// Login 用户登录接口
func (u *UserController) Login(ctx context.Context, req *http.Request) (rsp any, err error) {
	var params model.LoginReq
	if err = tools.RequestParam(&params, req); err != nil {
		println("parse login params error:", err.Error())
		return nil, err
	}

	result, err := u.UserService.Login(ctx, params)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetUserDetail 获取用户详情接口
func (u *UserController) GetUserDetail(ctx context.Context, req *http.Request) (rsp any, err error) {
	var params struct {
		UserId int64 `json:"userId"`
	}
	if err = tools.RequestParam(&params, req); err != nil {
		println("parse user detail params error:", err.Error())
		return nil, err
	}

	result, err := u.UserService.GetUserDetail(ctx, params.UserId)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *UserController) GetUserById(ctx context.Context, req *http.Request) (rsp any, err error) {
	var params model.GetUserReq
	if err = tools.RequestParam(&params, req); err != nil {
		println("parse params error:", err.Error())
		return
	}
	result, err := u.UserService.GetUserByUserId(ctx, params)
	if err != nil {
		return nil, err
	}
	rsp = result
	return rsp, nil
}
