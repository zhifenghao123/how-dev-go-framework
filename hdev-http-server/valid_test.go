package ihttp

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"
)

// TestValidateGet test validate
func TestValidateGet(t *testing.T) {
	testQuery := `a=1&b=apple&c=true&d=12.23`
	queryReq, _ := http.NewRequest(`GET`, `https://iheat.qq.com?`+testQuery, nil)

	data, err := Validate(queryReq, []Rule{
		{Key: "a", Type: "int", Rule: "required|include", Enum: []interface{}{1, 2, 3}, Msg: "a校验失败"},
		{Key: "b", Type: "string", Rule: "required|between:1,10", Msg: "b不能为空，长度为1~3"},
		{Key: "c", Type: "bool", Rule: "required", Msg: "c不能为空"},
		{Key: "d", Type: "float64", Rule: "required", Msg: "d不能为空"},
	})

	assert.Equal(t, nil, err)
	assert.Equal(t, 1, data.Int(`a`))
	assert.Equal(t, int64(0), data.Int64(`a`))
	assert.Equal(t, "apple", data.String(`b`))
	assert.Equal(t, true, data.Bool(`c`))
	assert.Equal(t, 12.23, data.Float64(`d`))
}

// TestValidatePostForm test multipart params validate
func TestValidatePostForm(t *testing.T) {
	postData := map[string]string{
		`name`: `yaoan`,
		`num`:  `12.34`,
	}
	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)
	for k, v := range postData {
		_ = w.WriteField(k, v)
	}
	_ = w.Close()
	req, _ := http.NewRequest(`POST`, `https://iheat.qq.com`, body)
	req.Header.Set("Content-Type", w.FormDataContentType())

	data, err := Validate(req, []Rule{
		{Key: `num`, Type: "float64", Rule: "required|include", Enum: []interface{}{12.34, 10.24}, Msg: "num校验失败"},
		{Key: `name`, Type: "string", Rule: "required|between:1,5", Msg: "name不能为空，长度为1~3"},
	})

	fmt.Println(`err`, err)
	if err != nil {
		return
	}
	fmt.Println(`num:`, data.Float64(`num`))
	fmt.Println(`name:`, data.String(`name`))
}

// TestValidatePostXForm test post x-www-form-urlencoded params validate
func TestValidatePostXForm(t *testing.T) {
	testBody := strings.NewReader(`name=ann&num=123`)
	req, err := http.NewRequest(`POST`, `https://iheat.qq.com`, testBody)
	if err != nil {
		fmt.Println(`request format error:`, err)
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", "nameCook=anyone")

	data, err := Validate(req, []Rule{
		{Key: `name`, Type: `string`, Rule: "required|include", Enum: []interface{}{`ann`, `yao`}, Msg: "name校验失败"},
		{Key: `num`, Type: "string", Rule: "required|between:1,5", Msg: "num不能为空，长度为1~2"},
		{Key: `Cookie|header`, Type: "string", Rule: "required", Msg: "nameCook不能为空"},
		{Key: "Content-Type|header", Type: "string", Rule: "required|include",
			Enum: []interface{}{`application/x-www-form-urlencoded`, `text/html; charset=utf-8`},
			Msg:  "Content-Type校验失败"},
	})

	fmt.Println(`err`, err)
	if err != nil {
		return
	}
	fmt.Println(`name:`, data.String(`name`))
	fmt.Println(`num:`, data.String(`num`))
	fmt.Println(`Cookie:`, data.String(`Cookie`))
	fmt.Println(`Content-Type:`, data.String(`Content-Type`))
}
