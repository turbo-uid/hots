package trending

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/turbo-uid/hots/globals"
)

type EnDataShellResponse struct {
	Status  int         `json:"status"`
	Des     string      `json:"des"`
	Version int         `json:"version"`
	Data    EnDataTable `json:"data"`
}

type EnDataTable struct {
	Table0 []EnDataTableList `json:"table0"`
}

type EnDataTableList struct {
	MovieName   string  `json:"MovieName"`
	ReleaseTime string  `json:"ReleaseTime"`
	BoxOffice   float64 `json:"BoxOffice"`
	Irank       int     `json:"Irank"`
}

var EnDataMUrl string = "https://ys.endata.cn/enlib-api/api/home/getrank_mainland.do"
var EnDataSUrl string = "https://ys.endata.cn/enlib-api/api/home/getrank_singleday.do"

func GetEnDataHot(dataType string) (globals.GblResp, error) {
	var EnDataUrl string

	if dataType == "1" {
		EnDataUrl = EnDataSUrl
	} else {
		EnDataUrl = EnDataMUrl
	}

	// 统一输出结果
	var resultResp globals.GblResp

	formData := url.Values{
		"r":    {fmtRandomNum()},
		"top":  {"50"},
		"type": {dataType},
	}

	// 发送 POST 请求
	resp, err := http.PostForm(EnDataUrl, formData)
	if err != nil {
		resultResp.Code = 1
		resultResp.Err = "Failed to send POST request"
		return resultResp, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		resultResp.Code = 1
		resultResp.Err = "Failed to read response body"
		return resultResp, err
	}

	var shellResp EnDataShellResponse
	err = json.Unmarshal(body, &shellResp)
	if err != nil {
		resultResp.Code = 1
		resultResp.Err = "Failed to parse JSON"
		return resultResp, err
	}

	resultResp.Succ = "ok"
	resultResp.Code = 0

	for _, v := range shellResp.Data.Table0 {
		var newData globals.GblRespData

		newData.Title = v.MovieName
		newData.Desc = v.ReleaseTime
		newData.HotVal = fmtBoxOffice(v.BoxOffice)
		newData.Pos = v.Irank
		newData.ToUrl = ""

		resultResp.Data = append(resultResp.Data, newData)
	}

	return resultResp, nil
}

func fmtBoxOffice(num float64) string {
	if num >= 1e8 {
		return fmt.Sprintf("%.2f亿", num/1e8)
	} else if num >= 1e4 {
		return fmt.Sprintf("%.2f万", num/1e4)
	} else {
		return fmt.Sprintf("%.2f", num)
	}
}

func fmtRandomNum() string {
	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Float64() * (0.1)
	return fmt.Sprintf("%.17f", randomNum)
}

