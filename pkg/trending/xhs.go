package trending

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/turbo-uid/hots/globals"
)

type XhsShellResponse struct {
	Succ bool    `json:"success"`
	Msg  string  `json:"msg"`
	Code int     `json:"code"`
	Data XhsData `json:"data"`
}

type XhsData struct {
	HotListId       string         `json:"hot_list_id"`
	IsNewHotListExp bool           `json:"is_new_hot_list_exp"`
	Items           []XhsItemsList `json:"items"`
}

type XhsItemsList struct {
	Title  string `json:"title"`
	HotVal string `json:"score"`
	Label  string `json:"word_type"`
}

func GetXhsHot() (globals.GblResp, error) {
	var resultResp globals.GblResp

	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://edith.xiaohongshu.com/api/sns/v1/search/hot_list", nil)
	if err != nil {
		resultResp.Code = 1
		resultResp.Err = "Error creating request"
		return resultResp, err
	}

	req.Header.Set("xy-direction", "22")
	req.Header.Set("accept-language", "zh-Hans-CN;q=1")
	req.Header.Set("shield", "XYAAAAAQAAAAEAAABTAAAAUzUWEe4xG1IYD9/c+qCLOlKGmTtFa+lG434Oe+FTRagxxoaz6rUWSZ3+juJYz8RZqct+oNMyZQxLEBaBEL+H3i0RhOBVGrauzVSARchIWFYwbwkV")
	req.Header.Set("xy-platform-info", "platform=iOS&version=8.7&build=8070515&deviceId=C323D3A5-6A27-4CE6-AA0E-51C9D4C26A24&bundle=com.xingin.discover")
	req.Header.Set("xy-common-params", "app_id=ECFAAF02&build=8070515&channel=AppStore&deviceId=C323D3A5-6A27-4CE6-AA0E-51C9D4C26A24&device_fingerprint=20230920120211bd7b71a80778509cf4211099ea911000010d2f20f6050264&device_fingerprint1=20230920120211bd7b71a80778509cf4211099ea911000010d2f20f6050264&device_model=phone&fid=1695182528-0-0-63b29d709954a1bb8c8733eb2fb58f29&gid=7dc4f3d168c355f1a886c54a898c6ef21fe7b9a847359afc77fc24ad&identifier_flag=0&lang=zh-Hans&launch_id=716882697&platform=iOS&project_id=ECFAAF&sid=session.1695189743787849952190&t=1695190591&teenager=0&tz=Asia/Shanghai&uis=light&version=8.7")
	req.Header.Set("referer", "https://app.xhs.cn/")

	resp, err := client.Do(req)
	if err != nil {
		resultResp.Code = 1
		resultResp.Err = "Error making GET request"
		return resultResp, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		resultResp.Code = 1
		resultResp.Err = "Failed to read response body"
		return resultResp, err
	}

	var shellResp XhsShellResponse
	err = json.Unmarshal(body, &shellResp)
	if err != nil {
		resultResp.Code = 1
		resultResp.Err = "Failed to parse JSON"
		return resultResp, err
	}

	resultResp.Succ = "ok"
	resultResp.Code = 0

	for k, v := range shellResp.Data.Items {
		var newData globals.GblRespData

		newData.Title = v.Title
		newData.Desc = ""
		newData.HotVal = v.HotVal
		newData.Pos = k + 1
		newData.ToUrl = ""
		if v.Label == "无" {
			newData.Lab = ""
		} else {
			newData.Lab = v.Label
		}

		resultResp.Data = append(resultResp.Data, newData)
	}

	return resultResp, nil
}

