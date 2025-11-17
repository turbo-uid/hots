package trending

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/turbo-uid/hots/globals"
)

type DongCheDiShellResponse struct {
	Status  int           `json:"status"`
	Message string        `json:"message"`
	Data    DongCheDiData `json:"data"`
}

type DongCheDiData struct {
	List []DongCheDiList `json:"list"`
}

type DongCheDiList struct {
	Title   string `json:"title"`
	HotVal  int    `json:"count"`
	GroupId string `json:"group_id"`
}

var DongCheDiUrl string = "https://www.dongchedi.com/motor/pc/content/pgc_content_rank?aid=1839&app_name=auto_web_pc&rank_type=pgc_article_total_rank"

func GetDongCheDiHot() (globals.GblResp, error) {
	var resultResp globals.GblResp

	resp, err := http.Get(DongCheDiUrl)
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

	var shellResp DongCheDiShellResponse
	err = json.Unmarshal(body, &shellResp)
	if err != nil {
		resultResp.Code = 1
		resultResp.Err = "Failed to parse JSON"
		return resultResp, err
	}

	resultResp.Succ = "ok"
	resultResp.Code = 0

	for k, v := range shellResp.Data.List {
		var newData globals.GblRespData

		newData.Title = v.Title
		newData.Desc = ""
		newData.HotVal = strconv.Itoa(v.HotVal)
		newData.Pos = k + 1
		newData.ToUrl = fmt.Sprintf("https://www.dongchedi.com/article/%s", v.GroupId)

		resultResp.Data = append(resultResp.Data, newData)
	}

	return resultResp, nil
}

