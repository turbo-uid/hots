package trending

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/turbo-uid/hots/globals"
)

type QcttData struct {
	Title   string   `json:"title"`
	Author  string   `json:"authorName"`
	PicList []string `json:"picUrlList"`
}

const QcttUrl string = "https://www.qctt.cn/channelDataList?page=1&id=1"

func GetQcttHot() (globals.GblResp, error) {
	var resultResp globals.GblResp

	resp, err := http.Get(QcttUrl)
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

	var shellResp []QcttData
	err = json.Unmarshal(body, &shellResp)
	if err != nil {
		resultResp.Code = 1
		resultResp.Err = "Failed to parse JSON"
		return resultResp, err
	}

	resultResp.Succ = "ok"
	resultResp.Code = 0

	for k, v := range shellResp {
		var newData globals.GblRespData

		newData.Title = v.Title
		newData.Desc = ""
		newData.HotVal = ""
		newData.Pos = k + 1
		newData.ToUrl = ""

		if len(v.PicList) > 0 {
			newData.Icon = v.PicList[0]
		}
		resultResp.Data = append(resultResp.Data, newData)
	}

	return resultResp, nil
}

