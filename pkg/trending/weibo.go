package trending

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/turbo-uid/hots/globals"
	"github.com/turbo-uid/hots/utils"
)

type WeiboShellResponse struct {
	Ok   int       `json:"ok"`
	Data WeiboData `json:"data"`
}

type WeiboData struct {
	Wlist    []WeiboList    `json:"realtime"`
	WTopList []WeiboTopList `json:"hotgovs"`
}

type WeiboList struct {
	Title    string `json:"word"`
	Desc     string `json:"note"`
	HotVal   int    `json:"num"`
	Icon     string `json:"icon"`
	Pos      int    `json:"realpos"`
	Label    string `json:"label_name"`
	FlagDesc string `json:"flag_desc"`
}

type WeiboTopList struct {
	Title string `json:"word"`
	Desc  string `json:"name"`
	Icon  string `json:"icon"`
	Pos   int    `json:"pos"`
	Label string `json:"icon_desc"`
}

func GetWeiboHot() (globals.GblResp, error) {
	var resultResp globals.GblResp

	resp, err := http.Get("https://weibo.com/ajax/side/hotSearch")
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

	var shellResp WeiboShellResponse
	err = json.Unmarshal(body, &shellResp)
	if err != nil {
		resultResp.Code = 1
		resultResp.Err = "Failed to parse JSON"
		return resultResp, err
	}

	resultResp.Succ = "ok"
	resultResp.Code = 0

	for _, v := range shellResp.Data.WTopList {
		var newData globals.GblRespData

		newData.Title = utils.RemoveChar(v.Title, "#")
		newData.Desc = ""
		newData.HotVal = "0"
		newData.Pos = 999
		newData.ToUrl = fmt.Sprintf("https://s.weibo.com/weibo?q=%%23%s%%23&t=31", utils.RemoveChar(v.Title, "#"))
		newData.IsTop = 1
		newData.Lab = v.Label
		resultResp.Data = append(resultResp.Data, newData)
	}

	for _, v := range shellResp.Data.Wlist {
		var newData globals.GblRespData

		newData.Title = v.Title
		newData.Desc = ""
		newData.HotVal = fmt.Sprintf("%s%s", v.FlagDesc, strconv.Itoa(v.HotVal))
		newData.Pos = v.Pos
		newData.ToUrl = fmt.Sprintf("https://s.weibo.com/weibo?q=%%23%s%%23&t=31", v.Title)
		newData.Lab = v.Label
		resultResp.Data = append(resultResp.Data, newData)
	}

	return resultResp, nil
}

