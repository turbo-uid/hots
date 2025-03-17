package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/turbo-uid/hots/globals"
	"github.com/turbo-uid/hots/utils"
)

type ToutiaoShellResponse struct {
	Status  string           `json:"status"`
	ImprId  string           `json:"impr_id"`
	Data    []ToutiaoList    `json:"data"`
	TopData []ToutiaoTopList `json:"fixed_top_data"`
}

type ToutiaoList struct {
	Title  string `json:"Title"`
	Desc   string `json:"QueryWord"`
	HotVal string `json:"HotValue"`
	Label  string `json:"Label"`
	ToUrl  string `json:"Url"`
}

type ToutiaoTopList struct {
	Title string `json:"Title"`
	ToUrl string `json:"Url"`
}

func ToutiaoHot(c *gin.Context) {

	if cacheResult, found := globals.GoCache.Get(utils.GetHotCacheKey(globals.ToutiaoFlag)); found {
		globals.GoLogger.Infof("API GET GCACHE %s", utils.GetHotCacheKey(globals.ToutiaoFlag))

		c.JSON(http.StatusOK, cacheResult)
		return
	}

	var resultResp globals.GblResp

	resp, err := http.Get("https://www.toutiao.com/hot-event/hot-board/?origin=toutiao_pc")
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

	var shellResp ToutiaoShellResponse
	err = json.Unmarshal(body, &shellResp)
	if err != nil {
		resultResp.Code = 1
		resultResp.Err = "Failed to parse JSON"
		c.JSON(http.StatusOK, resultResp)
		return
	}

	resultResp.Succ = "ok"
	resultResp.Code = 0

	for _, v := range shellResp.TopData {
		var newData globals.GblRespData

		newData.Title = v.Title
		newData.Desc = ""
		newData.HotVal = "0" // 兼容其他平台
		newData.Pos = 999
		newData.ToUrl = v.ToUrl
		newData.IsTop = 1
		newData.Lab = ""
		resultResp.Data = append(resultResp.Data, newData)
	}

	for k, v := range shellResp.Data {
		var newData globals.GblRespData

		newData.Title = v.Title
		newData.Desc = ""
		newData.HotVal = v.HotVal
		newData.Pos = k + 1
		newData.ToUrl = v.ToUrl
		newData.Lab = v.Label
		resultResp.Data = append(resultResp.Data, newData)
	}

	globals.GoCache.Set(utils.GetHotCacheKey(globals.ToutiaoFlag), resultResp, globals.HotCacheExpired)

	globals.GoLogger.Infof("API SET GCACHE %s DATA LEN %d", utils.GetHotCacheKey(globals.ToutiaoFlag), len(resultResp.Data))
	// 将解析后的数据作为响应返回给客户端
	c.JSON(http.StatusOK, resultResp)
}
