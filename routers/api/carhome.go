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

type CarHomeShellResponse struct {
	Code    int           `json:"returncode"`
	Message string        `json:"message"`
	Data    []CarHomeData `json:"result"`
}

type CarHomeData struct {
	Title  string `json:"title"`
	Desc   string `json:"subtitle"`
	HotVal int    `json:"order"`
	ToUrl  string `json:"url"`
	BizId  int    `json:"bizId"`
}

var CarHomeUrl string = "https://content.api.autohome.com.cn/pc/rank/list?ranktype=1&count=20"

func CarHomeHot(c *gin.Context) {

	if cacheResult, found := globals.GoCache.Get(utils.GetHotCacheKey(globals.CarHomeFlag)); found {
		globals.GoLogger.Infof("API GET GCACHE %s", utils.GetHotCacheKey(globals.CarHomeFlag))

		c.JSON(http.StatusOK, cacheResult)
		return
	}

	// 统一输出结果
	var resultResp globals.GblResp

	resp, err := http.Get(CarHomeUrl)

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

	var shellResp CarHomeShellResponse
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

		resultResp.Data = append(resultResp.Data, newData)
	}

	globals.GoCache.Set(utils.GetHotCacheKey(globals.CarHomeFlag), resultResp, globals.HotCacheExpired)

	globals.GoLogger.Infof("API SET GCACHE %s DATA LEN %d", utils.GetHotCacheKey(globals.CarHomeFlag), len(resultResp.Data))

	c.JSON(http.StatusOK, resultResp)
}
