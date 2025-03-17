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

type Wy163ShellResponse struct {
	Code int       `json:"code"`
	Mag  string    `json:"message"`
	Data Wy163Data `json:"data"`
}

type Wy163Data struct {
	RequestId string      `json:"requestId"`
	HotRank   []Wy163List `json:"hotRank"`
}

type Wy163List struct {
	Title  string `json:"hotWord"`
	Desc   string `json:"searchWord"`
	HotVal string `json:"exp"`
	Pos    int    `json:"rank"`
}

func Wy163Hot(c *gin.Context) {

	if cacheResult, found := globals.GoCache.Get(utils.GetHotCacheKey(globals.Wy163Flag)); found {
		globals.GoLogger.Infof("API GET GCACHE %s", utils.GetHotCacheKey(globals.Wy163Flag))

		c.JSON(http.StatusOK, cacheResult)
		return
	}

	var resultResp globals.GblResp

	resp, err := http.Get("https://gw.m.163.com/search/api/v2/hot-search")
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

	var shellResp Wy163ShellResponse
	err = json.Unmarshal(body, &shellResp)
	if err != nil {
		resultResp.Code = 1
		resultResp.Err = "Failed to parse JSON"
		c.JSON(http.StatusOK, resultResp)
		return
	}

	resultResp.Succ = "ok"
	resultResp.Code = 0

	for _, v := range shellResp.Data.HotRank {
		var newData globals.GblRespData

		newData.Title = v.Title
		newData.Desc = ""
		newData.HotVal = v.HotVal
		newData.Pos = v.Pos
		newData.ToUrl = fmt.Sprintf("https://m.163.com/search?keyword=%s", v.Title)
		newData.Lab = ""
		resultResp.Data = append(resultResp.Data, newData)
	}

	globals.GoCache.Set(utils.GetHotCacheKey(globals.Wy163Flag), resultResp, globals.HotCacheExpired)

	globals.GoLogger.Infof("API SET GCACHE %s DATA LEN %d", utils.GetHotCacheKey(globals.Wy163Flag), len(resultResp.Data))

	c.JSON(http.StatusOK, resultResp)
}
