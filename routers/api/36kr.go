package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/turbo-uid/hots/globals"
	"github.com/turbo-uid/hots/utils"

	"github.com/gin-gonic/gin"
)

type To36krShellResponse struct {
	Code int        `json:"code"`
	Data To36krData `json:"data"`
}

type To36krData struct {
	HotRankList []To36krHotRank `json:"hotRankList"`
}

type To36krHotRank struct {
	ItemId   int64             `json:"itemId"`
	Route    string            `json:"route"`
	Material To36krHotMaterial `json:"templateMaterial"`
}
type To36krHotMaterial struct {
	Title  string `json:"widgetTitle"`
	HotVal int    `json:"statRead"`
	Icon   string `json:"widgetImage"`
	ItemId int64  `json:"itemId"`
	Author string `json:"authorName"`
}

type To36krReq struct {
	PartnerId string         `json:"partner_id"`
	Timestamp int64          `json:"timestamp"`
	Param     To36krReqParam `json:"param"`
}
type To36krReqParam struct {
	SiteId     int `json:"siteId"`
	PlatformId int `json:"platformId"`
}

var To36krUrl string = "https://gateway.36kr.com/api/mis/nav/home/nav/rank/hot"

func To36krHot(c *gin.Context) {

	if cacheResult, found := globals.GoCache.Get(utils.GetHotCacheKey(globals.To36krFlag)); found {
		globals.GoLogger.Infof("API GET GCACHE %s", utils.GetHotCacheKey(globals.To36krFlag))

		c.JSON(http.StatusOK, cacheResult)
		return
	}

	// 统一输出结果
	var resultResp globals.GblResp

	reqParam := To36krReqParam{SiteId: 1, PlatformId: 2}
	reqData := To36krReq{PartnerId: "wap", Timestamp: utils.GetTimestamp(), Param: reqParam}
	jsonData, err := json.Marshal(reqData)
	if err != nil {
		resultResp.Code = 1
		resultResp.Err = "Error marshaling JSON"
		c.JSON(http.StatusOK, resultResp)
		return
	}

	req, err := http.NewRequest("POST", To36krUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		resultResp.Code = 1
		resultResp.Err = "Error creating request"
		c.JSON(http.StatusOK, resultResp)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		resultResp.Code = 1
		resultResp.Err = "Error making POST request"
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

	var shellResp To36krShellResponse
	err = json.Unmarshal(body, &shellResp)
	if err != nil {
		resultResp.Code = 1
		resultResp.Err = "Failed to parse JSON"
		c.JSON(http.StatusOK, resultResp)
		return
	}

	resultResp.Succ = "ok"
	resultResp.Code = 0

	for k, v := range shellResp.Data.HotRankList {
		var newData globals.GblRespData

		newData.Title = v.Material.Title
		newData.Desc = ""
		newData.HotVal = strconv.Itoa(v.Material.HotVal) // 兼容其他平台
		newData.Pos = k + 1
		newData.ToUrl = fmt.Sprintf("https://m.36kr.com/p/%d", v.ItemId)
		newData.Icon = v.Material.Icon
		resultResp.Data = append(resultResp.Data, newData)
	}

	globals.GoCache.Set(utils.GetHotCacheKey(globals.To36krFlag), resultResp, globals.HotCacheExpired)

	globals.GoLogger.Infof("API SET GCACHE %s DATA LEN %d", utils.GetHotCacheKey(globals.To36krFlag), len(resultResp.Data))
	// 将解析后的数据作为响应返回给客户端
	c.JSON(http.StatusOK, resultResp)
}
