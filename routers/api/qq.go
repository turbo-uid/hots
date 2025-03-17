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

type QqShellResponse struct {
	Ret    int            `json:"ret"`
	Idlist []QqIdlistData `json:"idlist"`
}

type QqIdlistData struct {
	IdsHash  string         `json:"ids_hash"`
	Newslist []QqActualData `json:"newslist"`
}

type QqActualData struct {
	Desc      string     `json:"abstract"`
	Longtitle string     `json:"longtitle"`
	ShareUrl  string     `json:"shareUrl"`
	MiniImage string     `json:"miniProShareImage"`
	HotEvent  QqHotEvent `json:"hotEvent"`
}

type QqHotEvent struct {
	Title  string `json:"title"`
	HotVal int    `json:"hotScore"`
	Pos    int    `json:"ranking"`
	IsTop  int    `json:"is_top"`
}

func QqHot(c *gin.Context) {

	if cacheResult, found := globals.GoCache.Get(utils.GetHotCacheKey(globals.QqFlag)); found {
		globals.GoLogger.Infof("API GET GCACHE %s", utils.GetHotCacheKey(globals.QqFlag))

		c.JSON(http.StatusOK, cacheResult)
		return
	}

	var resultResp globals.GblResp

	resp, err := http.Get("https://r.inews.qq.com/gw/event/hot_ranking_list?page_size=51")
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

	// 解析JSON响应
	var shellResp QqShellResponse
	err = json.Unmarshal(body, &shellResp)
	if err != nil {
		resultResp.Code = 1
		resultResp.Err = "Failed to parse JSON"
		c.JSON(http.StatusOK, resultResp)
		return
	}

	resultResp.Succ = "ok"
	resultResp.Code = 0

	// type QqActualData struct {
	// 	Desc string    `json:"abstract"`
	// 	Longtitle string    `json:"longtitle"`
	// 	ShareUrl string `json:"shareUrl"`
	// 	HotEvent QqHotEvent `json:"hotEvent"`
	// }

	// type QqHotEvent struct {
	// 	Title string `json:"title"`
	// 	HotVal int `json:"hotScore"`
	// 	Icon string    `json:"rec_night_icon"`
	// 	Pos int  `json:"ranking"`
	// 	IsTop int `json:"is_top"`
	// }

	listData := shellResp.Idlist[0]
	for _, v := range listData.Newslist {
		if len(v.ShareUrl) > 0 && len(v.Longtitle) > 0 {
			var newData globals.GblRespData

			newData.Title = v.HotEvent.Title
			newData.Desc = v.Longtitle
			newData.HotVal = strconv.Itoa(v.HotEvent.HotVal)
			newData.Pos = v.HotEvent.Pos - 1
			newData.ToUrl = v.ShareUrl
			newData.Icon = v.MiniImage
			newData.IsTop = 0
			// 兼容其他平台
			if v.HotEvent.IsTop == 1 && v.HotEvent.Pos == 1 {
				newData.IsTop = 1
				newData.Pos = 999
				newData.HotVal = "0"
			}
			// newData.Lab = ""
			resultResp.Data = append(resultResp.Data, newData)
		}
	}

	globals.GoCache.Set(utils.GetHotCacheKey(globals.QqFlag), resultResp, globals.HotCacheExpired)

	globals.GoLogger.Infof("API SET GCACHE %s DATA LEN %d", utils.GetHotCacheKey(globals.QqFlag), len(resultResp.Data))
	// 将解析后的数据作为响应返回给客户端
	c.JSON(http.StatusOK, resultResp)
}
