package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/turbo-uid/hots/globals"
	"github.com/turbo-uid/hots/utils"

	"github.com/gin-gonic/gin"
)

type EnDataShellResponse struct {
	Status  int         `json:"status"`
	Des     string      `json:"des"`
	Version int         `json:"version"`
	Data    EnDataTable `json:"data"`
}

type EnDataTable struct {
	Table0 []EnDataTableList `json:"table0"`
}

type EnDataTableList struct {
	MovieName   string  `json:"MovieName"`
	ReleaseTime string  `json:"ReleaseTime"`
	BoxOffice   float64 `json:"BoxOffice"`
	Irank       int     `json:"Irank"`
}

var EnDataMUrl string = "https://ys.endata.cn/enlib-api/api/home/getrank_mainland.do"
var EnDataSUrl string = "https://ys.endata.cn/enlib-api/api/home/getrank_singleday.do"

func EnDataHot(c *gin.Context) {
	var EnDataUrl, EnDataFlag, EnDataType string

	s := c.DefaultQuery("t", "m")
	if s == "s" {
		EnDataUrl = EnDataSUrl
		EnDataFlag = globals.EnDataSFlag
		EnDataType = "1"
	} else {
		EnDataUrl = EnDataMUrl
		EnDataFlag = globals.EnDataMFlag
		EnDataType = "0"
	}

	if cacheResult, found := globals.GoCache.Get(utils.GetHotCacheKey(EnDataFlag)); found {
		globals.GoLogger.Infof("API GET GCACHE %s", utils.GetHotCacheKey(EnDataFlag))

		c.JSON(http.StatusOK, cacheResult)
		return
	}

	// 统一输出结果
	var resultResp globals.GblResp

	formData := url.Values{
		"r":    {fmtRandomNum()},
		"top":  {"50"},
		"type": {EnDataType},
	}

	// 发送 POST 请求
	resp, err := http.PostForm(EnDataUrl, formData)
	if err != nil {
		resultResp.Code = 1
		resultResp.Err = "Failed to send POST request"
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

	var shellResp EnDataShellResponse
	err = json.Unmarshal(body, &shellResp)
	if err != nil {
		resultResp.Code = 1
		resultResp.Err = "Failed to parse JSON"
		c.JSON(http.StatusOK, resultResp)
		return
	}

	resultResp.Succ = "ok"
	resultResp.Code = 0

	for _, v := range shellResp.Data.Table0 {
		var newData globals.GblRespData

		type EnDataTableList struct {
			MovieName   string  `json:"MovieName"`
			ReleaseTime string  `json:"ReleaseTime"`
			BoxOffice   float64 `json:"BoxOffice"`
			Irank       int     `json:"Irank"`
		}

		newData.Title = v.MovieName
		newData.Desc = v.ReleaseTime
		newData.HotVal = fmtBoxOffice(v.BoxOffice)
		newData.Pos = v.Irank
		newData.ToUrl = ""

		resultResp.Data = append(resultResp.Data, newData)
	}

	globals.GoCache.Set(utils.GetHotCacheKey(EnDataFlag), resultResp, globals.HotCacheExpired)

	globals.GoLogger.Infof("API SET GCACHE %s DATA LEN %d", utils.GetHotCacheKey(EnDataFlag), len(resultResp.Data))

	c.JSON(http.StatusOK, resultResp)
}

func fmtBoxOffice(num float64) string {

	if num >= 1e8 {

		return fmt.Sprintf("%.2f亿", num/1e8)
	} else if num >= 1e4 {

		return fmt.Sprintf("%.2f万", num/1e4)
	} else {

		return fmt.Sprintf("%.2f", num)
	}
}

func fmtRandomNum() string {

	rand.Seed(time.Now().UnixNano())

	randomNum := rand.Float64() * (0.1)

	return fmt.Sprintf("%.17f", randomNum)
}
