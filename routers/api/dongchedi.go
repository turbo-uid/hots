package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/turbo-uid/hots/globals"
	"github.com/turbo-uid/hots/utils"

	"github.com/gin-gonic/gin"
)

type DongCheDiShellResponse struct {
	Status  int           `json:"status"`
	Message string        `json:"message"`
	Data    DongCheDiData `json:"data"`
}

type DongCheDiData struct {
	List []DongCheDiList `json:"list"`
}

type DongCheDiList struct {
	Title   string `json:"title"`
	HotVal  int    `json:"count"`
	GroupId string `json:"group_id"`
}

// var DongCheDiUrl string = "https://www.dongchedi.com/motor/pc/content/pgc_content_rank?aid=1839&app_name=auto_web_pc&rank_type=pgc_video_total_rank" // 视频榜
var DongCheDiUrl string = "https://www.dongchedi.com/motor/pc/content/pgc_content_rank?aid=1839&app_name=auto_web_pc&rank_type=pgc_article_total_rank" // 文章榜

func DongCheDiHot(c *gin.Context) {

	if cacheResult, found := globals.GoCache.Get(utils.GetHotCacheKey(globals.DongCheDiFlag)); found {
		globals.GoLogger.Infof("API GET GCACHE %s", utils.GetHotCacheKey(globals.DongCheDiFlag))

		c.JSON(http.StatusOK, cacheResult)
		return
	}

	// 统一输出结果
	var resultResp globals.GblResp

	resp, err := http.Get(DongCheDiUrl)

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

	var shellResp DongCheDiShellResponse
	err = json.Unmarshal(body, &shellResp)
	if err != nil {
		resultResp.Code = 1
		resultResp.Err = "Failed to parse JSON"
		c.JSON(http.StatusOK, resultResp)
		return
	}

	resultResp.Succ = "ok"
	resultResp.Code = 0

	for k, v := range shellResp.Data.List {
		var newData globals.GblRespData
		type DongCheDiList struct {
			Title   string `json:"title"`
			HotVal  int    `json:"count"`
			GroupId string `json:"group_id"`
		}
		newData.Title = v.Title
		newData.Desc = ""
		newData.HotVal = strconv.Itoa(v.HotVal)
		newData.Pos = k + 1
		// https://www.dongchedi.com/video/7457505435944223273
		// https://www.dongchedi.com/article/7459679563439473203
		newData.ToUrl = fmt.Sprintf("https://www.dongchedi.com/article/%s", v.GroupId)

		resultResp.Data = append(resultResp.Data, newData)
	}

	globals.GoCache.Set(utils.GetHotCacheKey(globals.DongCheDiFlag), resultResp, globals.HotCacheExpired)

	globals.GoLogger.Infof("API SET GCACHE %s DATA LEN %d", utils.GetHotCacheKey(globals.DongCheDiFlag), len(resultResp.Data))

	c.JSON(http.StatusOK, resultResp)
}
