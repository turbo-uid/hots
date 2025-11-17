package trending

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/turbo-uid/hots/globals"
)

type ThepaperShellResponse struct {
	ResultCode int          `json:"resultCode"`
	ResultMsg  string       `json:"resultMsg"`
	Data       ThepaperData `json:"data"`
}

type ThepaperData struct {
	HotNews []ThepaperList `json:"hotNews"`
}

type ThepaperList struct {
	Title          string `json:"name"`
	Icon           string `json:"sharePic"`
	ContId         string `json:"contId"`
	PubTimeNew     string `json:"pubTimeNew"`
	PraiseTimes    string `json:"praiseTimes"`
	InteractionNum string `json:"interactionNum"`
}

func GetThepaperHot() (globals.GblResp, error) {
	var resultResp globals.GblResp

	resp, err := http.Get("https://cache.thepaper.cn/contentapi/wwwIndex/rightSidebar")
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

	var shellResp ThepaperShellResponse
	err = json.Unmarshal(body, &shellResp)
	if err != nil {
		resultResp.Code = 1
		resultResp.Err = "Failed to parse JSON"
		return resultResp, err
	}

	resultResp.Succ = "ok"
	resultResp.Code = 0

	for k, v := range shellResp.Data.HotNews {
		var newData globals.GblRespData

		newData.Title = v.Title
		newData.Desc = fmt.Sprintf("评论数: %s 点赞数: %s 更新时间: %s", v.InteractionNum, v.PraiseTimes, v.PubTimeNew)
		newData.HotVal = ""
		newData.Pos = k + 1
		newData.ToUrl = fmt.Sprintf("https://www.thepaper.cn/newsDetail_forward_%s", v.ContId)
		newData.Lab = ""
		newData.Icon = v.Icon
		resultResp.Data = append(resultResp.Data, newData)
	}

	return resultResp, nil
}

