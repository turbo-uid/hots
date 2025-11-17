package trending

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/turbo-uid/hots/globals"
)

type HelloGithubShellResponse struct {
	Success bool              `json:"success"`
	Page    int               `json:"page"`
	Data    []HelloGithubData `json:"data"`
}

type HelloGithubData struct {
	Title       string `json:"title"`
	Author      string `json:"author"`
	Desc        string `json:"summary"`
	ClicksTotal int    `json:"clicks_total"`
	ItemId      string `json:"item_id"`
}

func GetHelloGithubHot() (globals.GblResp, error) {
	var resultResp globals.GblResp

	resp, err := http.Get("https://abroad.hellogithub.com/v1/?sort_by=all&tid=&page=1")
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

	var shellResp HelloGithubShellResponse
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
		newData.HotVal = fmt.Sprintf("%d", v.ClicksTotal)
		newData.Pos = k + 1
		newData.ToUrl = fmt.Sprintf("https://hellogithub.com/repository/%s", v.ItemId)

		resultResp.Data = append(resultResp.Data, newData)
	}

	return resultResp, nil
}

