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

type WeiboShellResponse struct {
	Ok   int       `json:"ok"`
	Data WeiboData `json:"data"`
}

type WeiboData struct {
	Wlist    []WeiboList    `json:"realtime"`
	WTopList []WeiboTopList `json:"hotgovs"`
}

type WeiboList struct {
	Title    string `json:"word"`
	Desc     string `json:"note"`
	HotVal   int    `json:"num"`
	Icon     string `json:"icon"`
	Pos      int    `json:"realpos"`
	Label    string `json:"label_name"`
	FlagDesc string `json:"flag_desc"`
}

type WeiboTopList struct {
	Title string `json:"word"`
	Desc  string `json:"name"`
	// HotVal int `json:"num"`
	Icon  string `json:"icon"`
	Pos   int    `json:"pos"`
	Label string `json:"icon_desc"`
}

func WeiboHot(c *gin.Context) {

	if cacheResult, found := globals.GoCache.Get(utils.GetHotCacheKey(globals.WeiboFlag)); found {
		globals.GoLogger.Infof("API GET GCACHE %s", utils.GetHotCacheKey(globals.WeiboFlag))

		c.JSON(http.StatusOK, cacheResult)
		return
	}

	// 统一输出结果
	var resultResp globals.GblResp

	resp, err := http.Get("https://weibo.com/ajax/side/hotSearch")
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

	var shellResp WeiboShellResponse
	err = json.Unmarshal(body, &shellResp)
	if err != nil {
		resultResp.Code = 1
		resultResp.Err = "Failed to parse JSON"
		c.JSON(http.StatusOK, resultResp)
		return
	}

	resultResp.Succ = "ok"
	resultResp.Code = 0

	for _, v := range shellResp.Data.WTopList {
		var newData globals.GblRespData

		newData.Title = utils.RemoveChar(v.Title, "#")
		newData.Desc = ""
		newData.HotVal = "0" // 兼容其他平台
		newData.Pos = 999
		newData.ToUrl = fmt.Sprintf("https://s.weibo.com/weibo?q=%%23%s%%23&t=31", utils.RemoveChar(v.Title, "#"))
		newData.IsTop = 1
		newData.Lab = v.Label
		resultResp.Data = append(resultResp.Data, newData)
	}

	for _, v := range shellResp.Data.Wlist {
		var newData globals.GblRespData

		newData.Title = v.Title
		newData.Desc = ""
		newData.HotVal = fmt.Sprintf("%s%s", v.FlagDesc, strconv.Itoa(v.HotVal)) // 兼容其他平台
		newData.Pos = v.Pos
		newData.ToUrl = fmt.Sprintf("https://s.weibo.com/weibo?q=%%23%s%%23&t=31", v.Title)
		newData.Lab = v.Label
		resultResp.Data = append(resultResp.Data, newData)
	}

	globals.GoCache.Set(utils.GetHotCacheKey(globals.WeiboFlag), resultResp, globals.HotCacheExpired)

	globals.GoLogger.Infof("API SET GCACHE %s DATA LEN %d", utils.GetHotCacheKey(globals.WeiboFlag), len(resultResp.Data))
	// 将解析后的数据作为响应返回给客户端
	c.JSON(http.StatusOK, resultResp)
}
