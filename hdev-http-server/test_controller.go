package ihttp

import (
	"context"
	"errors"
	"github.com/zhifenghao123/how-dev-go-framework/hdev-http-server/pb"
	"log"
	"net/http"
)

// TestController /*
type TestController struct {
}

// Test1 action for test string
func (TestController) Test1() string {
	return "hello test1"
}

// Test2 action for test error
func (TestController) Test2(_ context.Context, r *http.Request, res http.ResponseWriter) error {
	data, err := Validate(r, []Rule{
		{Key: "a", Type: "int", Rule: "required|include", Enum: []interface{}{1, 2, 3}, Msg: "a校验失败"},
		{Key: "b", Type: "string", Rule: "required|between:1,3", Msg: "b不能为空，长度为1~3"},
		{Key: "c", Type: "bool", Rule: "required", Msg: "c不能为空"},
		{Key: "d", Type: "float64", Rule: "required", Msg: "d不能为空"},
	})
	if err != nil {
		return err
	}
	log.Println(data.Int("a"), data.String("b"), data.Bool("c"), data.Float64("d"))
	res.Header().Set("test123", "123123123")
	return errors.New("error test2")
}

// Test3 action for test protobuf
func (TestController) Test3(_ context.Context, _ *http.Request, _ http.ResponseWriter) (*pb.UserInfo, error) {
	return &pb.UserInfo{UserList: []*pb.UserInfo_User{
		{
			Username: "djy",
			Age:      10,
			Graduate: "123",
		},
	}}, nil
}

// Test4 action for test nil
func (TestController) Test4(context.Context, *http.Request) {
	log.Println("lala test4")
	return
}
