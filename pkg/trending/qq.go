package trending

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/turbo-uid/hots/globals"
)

type QqShellResponse struct {
	Ret    int            `json:"ret"`
	Idlist []QqIdlistData `json:"idlist"`
}

type QqIdlistData struct {
	IdsHash  string         `json:"ids_hash"`
	Newslist []QqActualData `json:"newslist"`
}

type QqActualData struct {
	Desc      string     `json:"abstract"`
	Longtitle string     `json:"longtitle"`
	ShareUrl  string     `json:"shareUrl"`
	MiniImage string     `json:"miniProShareImage"`
	HotEvent  QqHotEvent `json:"hotEvent"`
}

type QqHotEvent struct {
	Title  string `json:"title"`
	HotVal int    `json:"hotScore"`
	Pos    int    `json:"ranking"`
	IsTop  int    `json:"is_top"`
}

func GetQqHot() (globals.GblResp, error) {
	var resultResp globals.GblResp

	resp, err := http.Get("https://r.inews.qq.com/gw/event/hot_ranking_list?page_size=51")
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

	var shellResp QqShellResponse
	err = json.Unmarshal(body, &shellResp)
	if err != nil {
		resultResp.Code = 1
		resultResp.Err = "Failed to parse JSON"
		return resultResp, err
	}

	resultResp.Succ = "ok"
	resultResp.Code = 0

	listData := shellResp.Idlist[0]
	for _, v := range listData.Newslist {
		if len(v.ShareUrl) > 0 && len(v.Longtitle) > 0 {
			var newData globals.GblRespData

			newData.Title = v.HotEvent.Title
			newData.Desc = v.Longtitle
			newData.HotVal = strconv.Itoa(v.HotEvent.HotVal)
			newData.Pos = v.HotEvent.Pos - 1
			newData.ToUrl = v.ShareUrl
			newData.Icon = v.MiniImage
			newData.IsTop = 0
			if v.HotEvent.IsTop == 1 && v.HotEvent.Pos == 1 {
				newData.IsTop = 1
				newData.Pos = 999
				newData.HotVal = "0"
			}
			resultResp.Data = append(resultResp.Data, newData)
		}
	}

	return resultResp, nil
}

