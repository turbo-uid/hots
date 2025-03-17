package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/turbo-uid/hots/globals"
	"github.com/turbo-uid/hots/utils"
)

type BiliShellResponse struct {
	Code int      `json:"code"`
	Data BiliData `json:"data"`
}

type BiliData struct {
	Blist    []BiliList    `json:"list"`
	BTopList []BiliTopList `json:"top_list"`
	Trackid  string        `json:"trackid"`
}

type BiliList struct {
	Title  string `json:"keyword"`
	Desc   string `json:"show_name"`
	HotVal int    `json:"hot_id"`
	Icon   string `json:"icon"`
	Pos    int    `json:"position"`
}

type BiliTopList struct {
	Title  string `json:"keyword"`
	Desc   string `json:"show_name"`
	HotVal int    `json:"hot_id"`
	Icon   string `json:"icon"`
	Pos    int    `json:"position"`
}

func BiliHot(c *gin.Context) {

	if cacheResult, found := globals.GoCache.Get(utils.GetHotCacheKey(globals.BiliFlag)); found {
		globals.GoLogger.Infof("API GET GCACHE %s", utils.GetHotCacheKey(globals.BiliFlag))

		c.JSON(http.StatusOK, cacheResult)
		return
	}

	// 统一输出结果
	var resultResp globals.GblResp

	resp, err := http.Get("https://app.bilibili.com/x/v2/search/trending/ranking?limit=30")
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

	var shellResp BiliShellResponse
	err = json.Unmarshal(body, &shellResp)
	if err != nil {
		resultResp.Code = 1
		resultResp.Err = "Failed to parse JSON"
		c.JSON(http.StatusOK, resultResp)
		return
	}

	resultResp.Succ = "ok"
	resultResp.Code = 0

	// https://search.bilibili.com/all?keyword= &from_source=webtop_search&spm_id_from=333.1007&search_source=4
	for _, v := range shellResp.Data.BTopList {
		var newData globals.GblRespData

		newData.Title = v.Title
		newData.Desc = "" // v.Desc
		newData.HotVal = strconv.Itoa(v.HotVal)
		newData.Pos = 999
		newData.ToUrl = fmt.Sprintf("https://search.bilibili.com/all?keyword=%s&from_source=webtop_search&spm_id_from=333.1007&search_source=4", v.Title)
		newData.IsTop = 1
		resultResp.Data = append(resultResp.Data, newData)
	}

	for _, v := range shellResp.Data.Blist {
		var newData globals.GblRespData

		newData.Title = v.Title
		newData.Desc = "" // v.Desc
		newData.HotVal = strconv.Itoa(v.HotVal)
		newData.Pos = v.Pos
		newData.ToUrl = fmt.Sprintf("https://search.bilibili.com/all?keyword=%s&from_source=webtop_search&spm_id_from=333.1007&search_source=4", v.Title)
		resultResp.Data = append(resultResp.Data, newData)
	}

	globals.GoCache.Set(utils.GetHotCacheKey(globals.BiliFlag), resultResp, globals.HotCacheExpired)

	globals.GoLogger.Infof("API SET GCACHE %s DATA LEN %d", utils.GetHotCacheKey(globals.BiliFlag), len(resultResp.Data))

	c.JSON(http.StatusOK, resultResp)
}
