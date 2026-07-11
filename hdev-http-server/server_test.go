package ihttp

import (
	"context"
	"fmt"
	"github.com/zhifenghao123/how-dev-go-framework/hdev-http-server/pb"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

var ctx context.Context
var stop context.CancelFunc

func init() {
	ctx, stop = context.WithCancel(context.Background())
}

// test server init
func Test_Init(t *testing.T) {

	s := HttpServer{
		Route: "config/test_router.yml",
		Ioc:   i,
		Ctx:   ctx,
		Port:  8080,
	}
	s.SetRouter(nil)
	s.SetHandler(nil)
	go func() {
		time.Sleep(1 * time.Second)
		s, _ := proto.Marshal(&pb.UserInfo{UserList: []*pb.UserInfo_User{
			{
				Username: "djy",
				Age:      10,
				Graduate: "123",
			},
		}})
		doRequest(t, "GET", "http://127.0.0.1:8080/test/pre1/pre2/test1",
			200, "hello test1") //TestController@Test1
		doRequest(t, "POST", "http://127.0.0.1:8080/test/test2/test10?a=1&b=aaa&c=0&d=1.23456",
			200, "error test2") //TestController@Test2
		doRequest(t, "GET", "http://127.0.0.1:8080/test/test2/test10",
			404, "404 page not found") // Page not found
		doRequest(t, "GET", "http://127.0.0.1:8080/test/test2/pre3/pre4/test8",
			200, string(s)) //TestController@Test3  protobuf
		doRequest(t, "GET", "http://127.0.0.1:8080/test/test2/test9",
			200, "null") //TestController@Test4
		stop()
	}()

	s.Execute()
}

func doRequest(t *testing.T, method string, uri string, code int, msg string) {
	request, err := http.NewRequest(method, uri, nil)
	if err != nil {
		fmt.Println(err)
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Errorf("Test %s %s failed %v", method, uri, err)
	}

	if response.StatusCode != code || string(body) != msg {
		t.Errorf("Test %s %s failed", method, uri)
		fmt.Println("response : code[", response.StatusCode, "], body : [", string(body), "]")
	}
}
