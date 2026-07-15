package ihttp

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"
)

var testReq *http.Request

func init() {
	testJson := `{"admin_user":"yaoan,sheltonliu","first_reviewer":"yaoan","second_reviewer":"",` +
		`"third_reviewer":"","available_budget":8000,"mba_code":"MBA-10000010","budget_id":"100073",` +
		`"activity_type":"reward","game_id":"jmx","game_name":"FIFA足球世界","title":"安徒生开发环境联调",` +
		`"content_type":"video","content_theme":"视频分发","channel_type":"douyin,kuaishou",` +
		`"start_time":"2021-10-21 17:25:00","end_time":"2021-10-30 16:25:54","show_reward":"1000",` +
		`"activity_desc":"安徒生开发环境联调","reward_count":75,` +
		`"banner_image":"https://tgl-image-1259052923.cos.com/source/20211021/heat_1634804905_5184.jpeg",` +
		`"cover_image":"https://tgl-image-1259052923.cos.com/source/20211021/heat_1634804898_3368.jpeg",` +
		`"detail_image":null,"activity_intro":"<p>安徒生联调奖金池</p>"}`
	testReq, _ = http.NewRequest("POST", `https://iheat.qq.com`, bytes.NewBuffer([]byte(testJson)))
}

// ActivityVO 活动详情VO
type ActivityVO struct {
	Title         string `json:"title" note:"活动标题" maxLen:"200" required:"1"`
	MbaCode       string `json:"mba_code" note:"立项号" required:"1"`
	BudgetId      string `json:"budget_id" note:"预算id" required:"1"`
	GameId        string `json:"game_id" note:"游戏id"`
	GameName      string `json:"game_name" note:"游戏名称" required:"1" maxLen:"50"`
	ActivityType  string `json:"activity_type" note:"活动类型" required:"1" include:"reward"`
	ContentType   string `json:"content_type" note:"内容类型" required:"1" maxLen:"20" enum:"all,video,article,live"`
	ContentTheme  string `json:"content_theme" note:"征集题材" maxLen:"50"`
	ActivityTag   string `json:"activity_tag" note:"活动标签" maxLen:"100"`
	ShowReward    string `json:"show_reward" note:"展示预算" required:"1"`
	ActivityDesc  string `json:"activity_desc" note:"活动描述" maxLen:"25"`
	ChannelType   string `json:"channel_type" note:"活动渠道" enum:"kuaishou,bilibili,channels,redbook,weibo,douyin"`
	ActivityIntro string `json:"activity_intro" note:"活动介绍" maxLen:"5000"`
	CoverImage    string `json:"cover_image" note:"活动封面" required:"1" maxLen:"200" type:"image"`
	DetailImage   string `json:"detail_image" note:"自定义详情图" maxLen:"200" type:"image"`
	BannerImage   string `json:"banner_image" note:"banner图" maxLen:"200" type:"image"`
	StartTime     string `json:"start_time" note:"征稿开始时间" required:"1" type:"time"`
	EndTime       string `json:"end_time" note:"征稿结束时间" required:"1" type:"time"`

	// 审核相关
	FirstReviewer  string `json:"first_reviewer" note:"预算一级审批人" required:"1" maxLen:"50"`
	SecondReviewer string `json:"second_reviewer" note:"预算二级审批人" maxLen:"50"`
	ThirdReviewer  string `json:"third_reviewer" note:"预算三级审批人"  maxLen:"50"`
}

// Test_JsonCheck test json check
func Test_JsonCheck(t *testing.T) {
	data := ActivityVO{}
	checkErr := CheckJson(testReq, &data)

	fmt.Println(`checkErr`, checkErr)
	fmt.Println(`data`, data)
}

// BenchmarkCheckJson benchmark test json check
func BenchmarkCheckJson(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var data ActivityVO
		_ = CheckJson(testReq, &data)
	}
}
