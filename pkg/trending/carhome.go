package trending

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/turbo-uid/hots/globals"
)

type CarHomeShellResponse struct {
	Code    int           `json:"returncode"`
	Message string        `json:"message"`
	Data    []CarHomeData `json:"result"`
}

type CarHomeData struct {
	Title  string `json:"title"`
	Desc   string `json:"subtitle"`
	HotVal int    `json:"order"`
	ToUrl  string `json:"url"`
	BizId  int    `json:"bizId"`
}

var CarHomeUrl string = "https://content.api.autohome.com.cn/pc/rank/list?ranktype=1&count=20"

func GetCarHomeHot() (globals.GblResp, error) {
	var resultResp globals.GblResp

	resp, err := http.Get(CarHomeUrl)
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

	var shellResp CarHomeShellResponse
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
		newData.Desc = v.Desc
		newData.HotVal = strconv.Itoa(v.HotVal)
		newData.Pos = k + 1
		newData.ToUrl = v.ToUrl

		resultResp.Data = append(resultResp.Data, newData)
	}

	return resultResp, nil
}

