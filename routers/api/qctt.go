package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/turbo-uid/hots/globals"
	"github.com/turbo-uid/hots/utils"

	"github.com/gin-gonic/gin"
)

type QcttData struct {
	Title   string   `json:"title"`
	Author  string   `json:"authorName"`
	PicList []string `json:"picUrlList"`
}

const QcttUrl string = "https://www.qctt.cn/channelDataList?page=1&id=1"

func QcttHot(c *gin.Context) {

	if cacheResult, found := globals.GoCache.Get(utils.GetHotCacheKey(globals.QcttFlag)); found {
		globals.GoLogger.Infof("API GET GCACHE %s", utils.GetHotCacheKey(globals.QcttFlag))

		c.JSON(http.StatusOK, cacheResult)
		return
	}

	// 统一输出结果
	var resultResp globals.GblResp

	resp, err := http.Get(QcttUrl)
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

	var shellResp []QcttData
	err = json.Unmarshal(body, &shellResp)
	if err != nil {
		resultResp.Code = 1
		resultResp.Err = "Failed to parse JSON"
		c.JSON(http.StatusOK, resultResp)
		return
	}

	resultResp.Succ = "ok"
	resultResp.Code = 0

	for k, v := range shellResp {
		var newData globals.GblRespData

		newData.Title = v.Title
		newData.Desc = "" // v.Desc
		newData.HotVal = ""
		newData.Pos = k + 1
		newData.ToUrl = "" // resourceLoc 暂不处理

		if len(v.PicList) > 0 {
			newData.Icon = v.PicList[0]
		}
		resultResp.Data = append(resultResp.Data, newData)
	}

	globals.GoCache.Set(utils.GetHotCacheKey(globals.QcttFlag), resultResp, globals.HotCacheExpired)

	globals.GoLogger.Infof("API SET GCACHE %s DATA LEN %d", utils.GetHotCacheKey(globals.QcttFlag), len(resultResp.Data))

	c.JSON(http.StatusOK, resultResp)
}
