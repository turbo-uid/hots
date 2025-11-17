package trending

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/turbo-uid/hots/globals"
)

type JueJinShellResponse struct {
	ErrNo  int          `json:"err_no"`
	ErrMsg string       `json:"err_msg"`
	Data   []JueJinData `json:"data"`
}

type JueJinData struct {
	Content        JueJinDataContent        `json:"content"`
	ContentCounter JueJinDataContentCounter `json:"content_counter"`
	Author         JueJinDataAuthor         `json:"author"`
}

type JueJinDataContent struct {
	Title     string `json:"title"`
	ContentId string `json:"content_id"`
}

type JueJinDataContentCounter struct {
	HotRank int `json:"hot_rank"`
	View    int `json:"view"`
}

type JueJinDataAuthor struct {
	Name string `json:"name"`
}

func GetJueJinHot() (globals.GblResp, error) {
	var resultResp globals.GblResp

	resp, err := http.Get("https://api.juejin.cn/content_api/v1/content/article_rank?category_id=1&type=hot")
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

	var shellResp JueJinShellResponse
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

		newData.Title = v.Content.Title
		newData.Desc = ""
		newData.HotVal = strconv.Itoa(v.ContentCounter.HotRank)
		newData.Pos = k + 1
		newData.ToUrl = fmt.Sprintf("https://juejin.cn/post/%s", v.Content.ContentId)

		resultResp.Data = append(resultResp.Data, newData)
	}

	return resultResp, nil
}

func GetJueJinAIBox() (globals.GblResp, error) {
	var resultResp globals.GblResp

	resp, err := http.Get("https://api.juejin.cn/content_api/v1/content/article_rank?category_id=6809637773935378440&type=hot")
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

	var shellResp JueJinShellResponse
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

		newData.Title = v.Content.Title
		newData.Desc = ""
		newData.HotVal = strconv.Itoa(v.ContentCounter.HotRank)
		newData.Pos = k + 1
		newData.ToUrl = fmt.Sprintf("https://juejin.cn/post/%s", v.Content.ContentId)

		resultResp.Data = append(resultResp.Data, newData)
	}

	return resultResp, nil
}

