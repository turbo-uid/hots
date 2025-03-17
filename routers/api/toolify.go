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

type ToolifyShellResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    ToolifyData `json:"data"`
}

type ToolifyData struct {
	Title int               `json:"current_page"`
	Data  []ToolifyDataList `json:"data"`
}

type ToolifyDataList struct {
	Name              string   `json:"name"`
	MonthVisitedVount int      `json:"month_visited_count"`
	Growth            int      `json:"growth"`
	GrowthRate        float64  `json:"growth_rate"`
	Description       string   `json:"description"`
	Tags              []string `json:"tags"`
	Date              string   `json:"date"`
}

type AIResp struct {
	Succ string       `json:"succ"`
	Err  string       `json:"err"`
	Code int          `json:"code"`
	Data []AIRespData `json:"data"`
}

type AIRespData struct {
	Name          string   `json:"name"`
	MonthlyVisits string   `json:"monthlyVisits"`
	Growth        string   `json:"growth"`
	GrowthRate    string   `json:"growthRate"`
	Description   string   `json:"description"`
	Tags          []string `json:"tags"`
	Expanded      bool     `json:"expanded"`
}

var ToolifyUrl string = "https://www.toolify.ai/self-api/v1/top/month-top?page=1&per_page=50&direction=desc&order_by=growth"

func ToolifyHot(c *gin.Context) {

	if cacheResult, found := globals.GoCache.Get(utils.GetHotCacheKey(globals.ToolifyFlag)); found {
		globals.GoLogger.Infof("API GET GCACHE %s", utils.GetHotCacheKey(globals.ToolifyFlag))

		c.JSON(http.StatusOK, cacheResult)
		return
	}

	// 统一输出结果
	var resultResp AIResp

	client := &http.Client{}

	req, err := http.NewRequest("GET", ToolifyUrl, nil)
	if err != nil {
		resultResp.Code = 1
		resultResp.Err = "Error creating request"
		c.JSON(http.StatusOK, resultResp)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		resultResp.Code = 1
		resultResp.Err = "Error making GET request"
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

	var shellResp ToolifyShellResponse
	err = json.Unmarshal(body, &shellResp)
	if err != nil {
		resultResp.Code = 1
		resultResp.Err = "Failed to parse JSON"
		c.JSON(http.StatusOK, resultResp)
		return
	}

	resultResp.Succ = "ok"
	resultResp.Code = 0

	for _, v := range shellResp.Data.Data {
		var newData AIRespData

		newData.Name = v.Name
		newData.Description = v.Description
		newData.Tags = v.Tags
		newData.Expanded = false
		newData.MonthlyVisits = fmtVisitedCount(v.MonthVisitedVount)
		newData.Growth = fmt.Sprintf("+%s", fmtVisitedCount(v.Growth))
		newData.GrowthRate = fmt.Sprintf("%.2f%%", v.GrowthRate*100)

		resultResp.Data = append(resultResp.Data, newData)
	}

	globals.GoCache.Set(utils.GetHotCacheKey(globals.ToolifyFlag), resultResp, globals.HotCacheExpired)

	globals.GoLogger.Infof("API SET GCACHE %s DATA LEN %d", utils.GetHotCacheKey(globals.ToolifyFlag), len(resultResp.Data))

	c.JSON(http.StatusOK, resultResp)
}

func fmtVisitedCount(num int) string {

	if num >= 1e8 {
		// 数字大于等于1亿，转换成“亿”
		return fmt.Sprintf("%.2f亿", float64(num)/1e8)
	} else if num >= 1e4 {
		// 数字大于等于1万，小于1亿，转换成“万”
		return fmt.Sprintf("%.2f万", float64(num)/1e4)
	} else {
		// 数字小于1万，直接显示原值
		return fmt.Sprintf("%.2f", float64(num))
	}
}
