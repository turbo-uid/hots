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

type DouyinShellResponse struct {
	StatusCode int        `json:"status_code"`
	Data       DouyinData `json:"data"`
}

type DouyinData struct {
	DyList []DouyinList `json:"word_list"`
}

type DouyinList struct {
	Title  string `json:"word"`
	HotVal int    `json:"hot_value"`
	Pos    int    `json:"position"`
	Label  int    `json:"label"`
	IsN1   bool   `json:"is_n1"`
}

func DouyinHot(c *gin.Context) {

	if cacheResult, found := globals.GoCache.Get(utils.GetHotCacheKey(globals.DouyinFlag)); found {
		globals.GoLogger.Infof("API GET GCACHE %s", utils.GetHotCacheKey(globals.DouyinFlag))

		c.JSON(http.StatusOK, cacheResult)
		return
	}

	// 统一输出结果
	var resultResp globals.GblResp

	resp, err := http.Get("https://aweme-lq.snssdk.com/aweme/v1/hot/search/list/?device_platform=webapp&aid=6383&channel=channel_pc_web&detail_list=1&source=6&main_billboard_count=5&update_version_code=170400&pc_client_type=1&pc_libra_divert=Windows&version_code=170400&version_name=17.4.0&cookie_enabled=true&screen_width=1920&screen_height=1080&browser_language=zh-CN&browser_platform=Win32&browser_name=Chrome&browser_version=131.0.0.0&browser_online=true&engine_name=Blink&engine_version=131.0.0.0&os_name=Windows") // 替换为实际的第三方API URL

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

	var shellResp DouyinShellResponse
	err = json.Unmarshal(body, &shellResp)
	if err != nil {
		resultResp.Code = 1
		resultResp.Err = "Failed to parse JSON"
		c.JSON(http.StatusOK, resultResp)
		return
	}

	resultResp.Succ = "ok"
	resultResp.Code = 0

	for _, v := range shellResp.Data.DyList {
		var newData globals.GblRespData

		newData.Title = v.Title
		newData.Desc = ""
		newData.HotVal = strconv.Itoa(v.HotVal)
		newData.Pos = v.Pos
		newData.ToUrl = fmt.Sprintf("https://www.douyin.com/root/search/%s?aid=8f302f2a-b661-4a1b-a88a-1027f4475461&type=general", v.Title)
		if v.Label == 3 {
			newData.Lab = "热"
		} else if v.Label == 1 {
			newData.Lab = "新"
		} else if v.Label == 8 {
			newData.Lab = "独家"
		} else if v.Label == 5 {
			newData.Lab = "首发"
		} else {
			newData.Lab = ""
		}
		if v.IsN1 || (v.HotVal == 0 && v.Pos == 0) {
			newData.Pos = 999
			newData.IsTop = 1
		}

		resultResp.Data = append(resultResp.Data, newData)
	}

	globals.GoCache.Set(utils.GetHotCacheKey(globals.DouyinFlag), resultResp, globals.HotCacheExpired)

	globals.GoLogger.Infof("API SET GCACHE %s DATA LEN %d", utils.GetHotCacheKey(globals.DouyinFlag), len(resultResp.Data))
	// 将解析后的数据作为响应返回给客户端
	c.JSON(http.StatusOK, resultResp)
}
