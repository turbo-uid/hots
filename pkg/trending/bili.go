package trending

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/turbo-uid/hots/globals"
)

type BiliShellResponse struct {
	Code int      `json:"code"`
	Data BiliData `json:"data"`
}

type BiliData struct {
	Blist    []BiliList    `json:"list"`
	BTopList []BiliTopList `json:"top_list"`
	Trackid  string        `json:"trackid"`
}

type BiliList struct {
	Title  string `json:"keyword"`
	Desc   string `json:"show_name"`
	HotVal int    `json:"hot_id"`
	Icon   string `json:"icon"`
	Pos    int    `json:"position"`
}

type BiliTopList struct {
	Title  string `json:"keyword"`
	Desc   string `json:"show_name"`
	HotVal int    `json:"hot_id"`
	Icon   string `json:"icon"`
	Pos    int    `json:"position"`
}

func GetBiliHot() (globals.GblResp, error) {
	var resultResp globals.GblResp

	resp, err := http.Get("https://app.bilibili.com/x/v2/search/trending/ranking?limit=30")
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

	var shellResp BiliShellResponse
	err = json.Unmarshal(body, &shellResp)
	if err != nil {
		resultResp.Code = 1
		resultResp.Err = "Failed to parse JSON"
		return resultResp, err
	}

	resultResp.Succ = "ok"
	resultResp.Code = 0

	for _, v := range shellResp.Data.BTopList {
		var newData globals.GblRespData

		newData.Title = v.Title
		newData.Desc = ""
		newData.HotVal = strconv.Itoa(v.HotVal)
		newData.Pos = 999
		newData.ToUrl = fmt.Sprintf("https://search.bilibili.com/all?keyword=%s&from_source=webtop_search&spm_id_from=333.1007&search_source=4", v.Title)
		newData.IsTop = 1
		resultResp.Data = append(resultResp.Data, newData)
	}

	for _, v := range shellResp.Data.Blist {
		var newData globals.GblRespData

		newData.Title = v.Title
		newData.Desc = ""
		newData.HotVal = strconv.Itoa(v.HotVal)
		newData.Pos = v.Pos
		newData.ToUrl = fmt.Sprintf("https://search.bilibili.com/all?keyword=%s&from_source=webtop_search&spm_id_from=333.1007&search_source=4", v.Title)
		resultResp.Data = append(resultResp.Data, newData)
	}

	return resultResp, nil
}

