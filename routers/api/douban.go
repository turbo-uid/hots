package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/turbo-uid/hots/globals"
	"github.com/turbo-uid/hots/utils"
)

type DoubanShellResponse struct {
	Data []DoubanList `json:"gallery_topics"`
}

type DoubanList struct {
	Title    string         `json:"title"`
	Desc     string         `json:"card_subtitle"`
	HotVal   int            `json:"read_count"`
	ToUrl    string         `json:"url"`
	TailIcon DoubanTailIcon `json:"tail_icon"`
}

type DoubanTailIcon struct {
	Text    string `json:"text"`
	BgColor string `json:"bg_color"`
}

func DoubanHot(c *gin.Context) {

	if cacheResult, found := globals.GoCache.Get(utils.GetHotCacheKey(globals.DoubanFlag)); found {
		globals.GoLogger.Infof("API GET GCACHE %s", utils.GetHotCacheKey(globals.DoubanFlag))

		c.JSON(http.StatusOK, cacheResult)
		return
	}

	// 统一输出结果
	var resultResp globals.GblResp

	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://m.douban.com/rexxar/api/v2/search/hots?ck=", nil)
	if err != nil {
		resultResp.Code = 1
		resultResp.Err = "Error creating request"
		c.JSON(http.StatusOK, resultResp)
		return
	}
	// 设置请求头部
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Origin", "https://www.douban.com")
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Referer", "https://www.douban.com/")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Mobile/15E148 Safari/604.1")

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

	var shellResp DoubanShellResponse
	err = json.Unmarshal(body, &shellResp)
	if err != nil {
		resultResp.Code = 1
		resultResp.Err = "Failed to parse JSON"
		c.JSON(http.StatusOK, resultResp)
		return
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
		newData.Lab = v.TailIcon.Text
		resultResp.Data = append(resultResp.Data, newData)
	}

	globals.GoCache.Set(utils.GetHotCacheKey(globals.DoubanFlag), resultResp, globals.HotCacheExpired)

	globals.GoLogger.Infof("API SET GCACHE %s DATA LEN %d", utils.GetHotCacheKey(globals.DoubanFlag), len(resultResp.Data))
	// 将解析后的数据作为响应返回给客户端
	c.JSON(http.StatusOK, resultResp)
}
