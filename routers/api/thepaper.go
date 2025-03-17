package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/turbo-uid/hots/globals"
	"github.com/turbo-uid/hots/utils"
)

type ThepaperShellResponse struct {
	ResultCode int          `json:"resultCode"`
	ResultMsg  string       `json:"resultMsg"`
	Data       ThepaperData `json:"data"`
}

type ThepaperData struct {
	HotNews []ThepaperList `json:"hotNews"`
}

type ThepaperList struct {
	Title          string `json:"name"`
	Icon           string `json:"sharePic"` // https://imgpai.thepaper.cn/newpai/image/1736252107447_Eef8na_1736252107733.png
	ContId         string `json:"contId"`   // https://www.thepaper.cn/newsDetail_forward_
	PubTimeNew     string `json:"pubTimeNew"`
	PraiseTimes    string `json:"praiseTimes"`    // 点赞
	InteractionNum string `json:"interactionNum"` // 评论
}

func ThepaperHot(c *gin.Context) {

	if cacheResult, found := globals.GoCache.Get(utils.GetHotCacheKey(globals.ThepaperFlag)); found {
		globals.GoLogger.Infof("API GET GCACHE %s", utils.GetHotCacheKey(globals.ThepaperFlag))

		c.JSON(http.StatusOK, cacheResult)
		return
	}

	var resultResp globals.GblResp

	resp, err := http.Get("https://cache.thepaper.cn/contentapi/wwwIndex/rightSidebar")
	if err != nil {
		resultResp.Code = 1
		resultResp.Err = "Failed to fetch data"
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

	var shellResp ThepaperShellResponse
	err = json.Unmarshal(body, &shellResp)
	if err != nil {
		resultResp.Code = 1
		resultResp.Err = "Failed to parse JSON"
		c.JSON(http.StatusOK, resultResp)
		return
	}

	resultResp.Succ = "ok"
	resultResp.Code = 0

	for k, v := range shellResp.Data.HotNews {
		var newData globals.GblRespData

		newData.Title = v.Title
		newData.Desc = fmt.Sprintf("评论数: %s 点赞数: %s 更新时间: %s", v.InteractionNum, v.PraiseTimes, v.PubTimeNew)
		newData.HotVal = ""
		newData.Pos = k + 1
		newData.ToUrl = fmt.Sprintf("https://www.thepaper.cn/newsDetail_forward_%s", v.ContId)
		newData.Lab = ""
		newData.Icon = v.Icon
		resultResp.Data = append(resultResp.Data, newData)
	}

	globals.GoCache.Set(utils.GetHotCacheKey(globals.ThepaperFlag), resultResp, globals.HotCacheExpired)

	globals.GoLogger.Infof("API SET GCACHE %s DATA LEN %d", utils.GetHotCacheKey(globals.ThepaperFlag), len(resultResp.Data))
	// 将解析后的数据作为响应返回给客户端
	c.JSON(http.StatusOK, resultResp)
}
