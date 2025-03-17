package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/turbo-uid/hots/globals"
	"github.com/turbo-uid/hots/utils"

	"github.com/gin-gonic/gin"
)

type CsdnShellResponse struct {
	Code    int        `json:"code"`
	TraceId string     `json:"traceId"`
	Data    []CsdnData `json:"data"`
}

type CsdnData struct {
	Title          string   `json:"articleTitle"`
	Url            string   `json:"articleDetailUrl"`
	PcHotRankScore string   `json:"pcHotRankScore"`
	HotRankScore   string   `json:"hotRankScore"`
	Author         string   `json:"nickName"`
	PicList        []string `json:"picList"`
}

func CsdnHot(c *gin.Context) {

	if cacheResult, found := globals.GoCache.Get(utils.GetHotCacheKey(globals.CsdnFlag)); found {
		globals.GoLogger.Infof("API GET GCACHE %s", utils.GetHotCacheKey(globals.CsdnFlag))

		c.JSON(http.StatusOK, cacheResult)
		return
	}

	// 统一输出结果
	var resultResp globals.GblResp

	resp, err := http.Get("https://blog.csdn.net/phoenix/web/blog/hot-rank?page=0&pageSize=30")
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

	var shellResp CsdnShellResponse
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
		newData.Desc = "" // v.Desc
		newData.HotVal = v.PcHotRankScore
		newData.Pos = k + 1
		newData.ToUrl = v.Url

		if len(v.PicList) > 0 {
			newData.Icon = v.PicList[0]
		}
		resultResp.Data = append(resultResp.Data, newData)
	}

	globals.GoCache.Set(utils.GetHotCacheKey(globals.CsdnFlag), resultResp, globals.HotCacheExpired)

	globals.GoLogger.Infof("API SET GCACHE %s DATA LEN %d", utils.GetHotCacheKey(globals.CsdnFlag), len(resultResp.Data))

	c.JSON(http.StatusOK, resultResp)
}

func CsdnContent(c *gin.Context) {

	if cacheResult, found := globals.GoCache.Get(utils.GetHotCacheKey(globals.CsdnContentFlag)); found {
		globals.GoLogger.Infof("API GET GCACHE %s", utils.GetHotCacheKey(globals.CsdnContentFlag))

		c.JSON(http.StatusOK, cacheResult)
		return
	}

	// 统一输出结果
	var resultResp globals.GblResp

	resp, err := http.Get("https://blog.csdn.net/phoenix/web/blog/hot-rank?page=0&pageSize=50&child_channel=%E4%BA%BA%E5%B7%A5%E6%99%BA%E8%83%BD&type=")
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

	var shellResp CsdnShellResponse
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
		newData.Desc = "" // v.Desc
		newData.HotVal = v.PcHotRankScore
		newData.Pos = k + 1
		newData.ToUrl = v.Url

		if len(v.PicList) > 0 {
			newData.Icon = v.PicList[0]
		}
		resultResp.Data = append(resultResp.Data, newData)
	}

	globals.GoCache.Set(utils.GetHotCacheKey(globals.CsdnContentFlag), resultResp, globals.HotCacheExpired)

	globals.GoLogger.Infof("API SET GCACHE %s DATA LEN %d", utils.GetHotCacheKey(globals.CsdnContentFlag), len(resultResp.Data))

	c.JSON(http.StatusOK, resultResp)
}
