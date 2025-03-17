package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/turbo-uid/hots/globals"
	"github.com/turbo-uid/hots/utils"
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

func XhsHot(c *gin.Context) {

	if cacheResult, found := globals.GoCache.Get(utils.GetHotCacheKey(globals.XhsFlag)); found {
		globals.GoLogger.Infof("API GET GCACHE %s", utils.GetHotCacheKey(globals.XhsFlag))

		c.JSON(http.StatusOK, cacheResult)
		return
	}

	var resultResp globals.GblResp

	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://edith.xiaohongshu.com/api/sns/v1/search/hot_list", nil)
	if err != nil {
		resultResp.Code = 1
		resultResp.Err = "Error creating request"
		c.JSON(http.StatusOK, resultResp)
		return
	}
	// 设置请求头部
	// req.Header.Set("x-legacy-fid", " 1695182528-0-0-63b29d709954a1bb8c8733eb2fb58f29")
	req.Header.Set("xy-direction", "22")
	// req.Header.Set("x-xray-traceid", "c5589dd4090584f0d665967ecbd70541")
	// req.Header.Set("x-b3-traceid", "138c7d341c388db4")
	// req.Header.Set("x-mini-gid", " 7dc4f3d168c355f1a886c54a898c6ef21fe7b9a847359afc77fc24ad")
	req.Header.Set("accept-language", "zh-Hans-CN;q=1")
	// req.Header.Set("x-mini-sig", "04b34b5ab16c061892f155202ae7df67f303d96d8ebe6af24337a34b056b2526")
	// req.Header.Set("x-legacy-did", "C323D3A5-6A27-4CE6-AA0E-51C9D4C26A24")
	// req.Header.Set("x-net-core", "crn")
	req.Header.Set("shield", "XYAAAAAQAAAAEAAABTAAAAUzUWEe4xG1IYD9/c+qCLOlKGmTtFa+lG434Oe+FTRagxxoaz6rUWSZ3+juJYz8RZqct+oNMyZQxLEBaBEL+H3i0RhOBVGrauzVSARchIWFYwbwkV")
	// req.Header.Set("x-legacy-smid", "20230920120211bd7b71a80778509cf4211099ea911000010d2f20f6050264")
	req.Header.Set("xy-platform-info", "platform=iOS&version=8.7&build=8070515&deviceId=C323D3A5-6A27-4CE6-AA0E-51C9D4C26A24&bundle=com.xingin.discover")
	// req.Header.Set("mode", "gslb")
	req.Header.Set("xy-common-params", "app_id=ECFAAF02&build=8070515&channel=AppStore&deviceId=C323D3A5-6A27-4CE6-AA0E-51C9D4C26A24&device_fingerprint=20230920120211bd7b71a80778509cf4211099ea911000010d2f20f6050264&device_fingerprint1=20230920120211bd7b71a80778509cf4211099ea911000010d2f20f6050264&device_model=phone&fid=1695182528-0-0-63b29d709954a1bb8c8733eb2fb58f29&gid=7dc4f3d168c355f1a886c54a898c6ef21fe7b9a847359afc77fc24ad&identifier_flag=0&lang=zh-Hans&launch_id=716882697&platform=iOS&project_id=ECFAAF&sid=session.1695189743787849952190&t=1695190591&teenager=0&tz=Asia/Shanghai&uis=light&version=8.7")
	// req.Header.Set("x-legacy-sid", "session.1695189743787849952190")
	// req.Header.Set("x-raw-ptr", "0")
	req.Header.Set("referer", "https://app.xhs.cn/")
	// req.Header.Set("cookie", "acw_tc=2c0be1613d1a3c5a6d5cc9108c2172e9f4e0958c7ccf9908562a2dfb7f9014b8")

	resp, err := client.Do(req)
	if err != nil {
		resultResp.Code = 1
		resultResp.Err = "Error making GET request"
		c.JSON(http.StatusOK, resultResp)
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		resultResp.Code = 1
		resultResp.Err = "Failed to read response body"
		c.JSON(http.StatusOK, resultResp)
		return
	}

	var shellResp XhsShellResponse
	err = json.Unmarshal(body, &shellResp)
	if err != nil {
		resultResp.Code = 1
		resultResp.Err = "Failed to parse JSON"
		c.JSON(http.StatusOK, resultResp)
		return
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

	globals.GoCache.Set(utils.GetHotCacheKey(globals.XhsFlag), resultResp, globals.HotCacheExpired)

	globals.GoLogger.Infof("API SET GCACHE %s DATA LEN %d", utils.GetHotCacheKey(globals.XhsFlag), len(resultResp.Data))
	// 将解析后的数据作为响应返回给客户端
	c.JSON(http.StatusOK, resultResp)
}
