package trending

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/turbo-uid/hots/globals"
)

type CsdnShellResponse struct {
	Code    int        `json:"code"`
	TraceId string     `json:"traceId"`
	Data    []CsdnData `json:"data"`
}

type CsdnData struct {
	Title          string   `json:"articleTitle"`
	Url            string   `json:"articleDetailUrl"`
	PcHotRankScore string   `json:"pcHotRankScore"`
	HotRankScore   string   `json:"hotRankScore"`
	Author         string   `json:"nickName"`
	PicList        []string `json:"picList"`
}

func GetCsdnHot() (globals.GblResp, error) {
	var resultResp globals.GblResp

	resp, err := http.Get("https://blog.csdn.net/phoenix/web/blog/hot-rank?page=0&pageSize=30")
	if err != nil {
		resultResp.Code = 1
		resultResp.Err = "Failed to fetch data"
		return resultResp, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		resultResp.Code = 1
		resultResp.Err = "Failed to read response body"
		return resultResp, err
	}

	var shellResp CsdnShellResponse
	err = json.Unmarshal(body, &shellResp)
	if err != nil {
		resultResp.Code = 1
		resultResp.Err = "Failed to parse JSON"
		return resultResp, err
	}

	resultResp.Succ = "ok"
	resultResp.Code = 0

	for k, v := range shellResp.Data {
		var newData globals.GblRespData

		newData.Title = v.Title
		newData.Desc = ""
		newData.HotVal = v.PcHotRankScore
		newData.Pos = k + 1
		newData.ToUrl = v.Url

		if len(v.PicList) > 0 {
			newData.Icon = v.PicList[0]
		}
		resultResp.Data = append(resultResp.Data, newData)
	}

	return resultResp, nil
}

func GetCsdnContent() (globals.GblResp, error) {
	var resultResp globals.GblResp

	resp, err := http.Get("https://blog.csdn.net/phoenix/web/blog/hot-rank?page=0&pageSize=50&child_channel=%E4%BA%BA%E5%B7%A5%E6%99%BA%E8%83%BD&type=")
	if err != nil {
		resultResp.Code = 1
		resultResp.Err = "Failed to fetch data"
		return resultResp, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		resultResp.Code = 1
		resultResp.Err = "Failed to read response body"
		return resultResp, err
	}

	var shellResp CsdnShellResponse
	err = json.Unmarshal(body, &shellResp)
	if err != nil {
		resultResp.Code = 1
		resultResp.Err = "Failed to parse JSON"
		return resultResp, err
	}

	resultResp.Succ = "ok"
	resultResp.Code = 0

	for k, v := range shellResp.Data {
		var newData globals.GblRespData

		newData.Title = v.Title
		newData.Desc = ""
		newData.HotVal = v.PcHotRankScore
		newData.Pos = k + 1
		newData.ToUrl = v.Url

		if len(v.PicList) > 0 {
			newData.Icon = v.PicList[0]
		}
		resultResp.Data = append(resultResp.Data, newData)
	}

	return resultResp, nil
}

