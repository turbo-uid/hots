package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/turbo-uid/hots/globals"
	"github.com/turbo-uid/hots/utils"

	"github.com/gin-gonic/gin"
)

type HelloGithubShellResponse struct {
	Success bool              `json:"success"`
	Page    int               `json:"page"`
	Data    []HelloGithubData `json:"data"`
}

type HelloGithubData struct {
	Title       string `json:"title"`
	Author      string `json:"author"`
	Desc        string `json:"summary"`
	ClicksTotal int    `json:"clicks_total"`
	ItemId      string `json:"item_id"`
}

func HelloGithubHot(c *gin.Context) {

	if cacheResult, found := globals.GoCache.Get(utils.GetHotCacheKey(globals.HelloGithubFlag)); found {
		globals.GoLogger.Infof("API GET GCACHE %s", utils.GetHotCacheKey(globals.HelloGithubFlag))

		c.JSON(http.StatusOK, cacheResult)
		return
	}

	// 统一输出结果
	var resultResp globals.GblResp

	resp, err := http.Get("https://abroad.hellogithub.com/v1/?sort_by=all&tid=&page=1")
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

	var shellResp HelloGithubShellResponse
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
		newData.HotVal = fmt.Sprintf("%d", v.ClicksTotal)
		newData.Pos = k + 1
		newData.ToUrl = fmt.Sprintf("https://hellogithub.com/repository/%s", v.ItemId)

		resultResp.Data = append(resultResp.Data, newData)
	}

	globals.GoCache.Set(utils.GetHotCacheKey(globals.HelloGithubFlag), resultResp, globals.HotCacheExpired)

	globals.GoLogger.Infof("API SET GCACHE %s DATA LEN %d", utils.GetHotCacheKey(globals.HelloGithubFlag), len(resultResp.Data))

	c.JSON(http.StatusOK, resultResp)
}
