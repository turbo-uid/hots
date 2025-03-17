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

type JueJinShellResponse struct {
	ErrNo  int          `json:"err_no"`
	ErrMsg string       `json:"err_msg"`
	Data   []JueJinData `json:"data"`
}

type JueJinData struct {
	Content        JueJinDataContent        `json:"content"`
	ContentCounter JueJinDataContentCounter `json:"content_counter"`
	Author         JueJinDataAuthor         `json:"author"`
}

type JueJinDataContent struct {
	Title     string `json:"title"`
	ContentId string `json:"content_id"`
}

type JueJinDataContentCounter struct {
	HotRank int `json:"hot_rank"`
	View    int `json:"view"`
}

type JueJinDataAuthor struct {
	Name string `json:"name"`
}

func JueJinHot(c *gin.Context) {

	if cacheResult, found := globals.GoCache.Get(utils.GetHotCacheKey(globals.JueJinFlag)); found {
		globals.GoLogger.Infof("API GET GCACHE %s", utils.GetHotCacheKey(globals.JueJinFlag))

		c.JSON(http.StatusOK, cacheResult)
		return
	}

	// 统一输出结果
	var resultResp globals.GblResp

	resp, err := http.Get("https://api.juejin.cn/content_api/v1/content/article_rank?category_id=1&type=hot")
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

	var shellResp JueJinShellResponse
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

		newData.Title = v.Content.Title
		newData.Desc = ""
		newData.HotVal = strconv.Itoa(v.ContentCounter.HotRank)
		newData.Pos = k + 1
		newData.ToUrl = fmt.Sprintf("https://juejin.cn/post/%s", v.Content.ContentId)

		resultResp.Data = append(resultResp.Data, newData)
	}

	globals.GoCache.Set(utils.GetHotCacheKey(globals.JueJinFlag), resultResp, globals.HotCacheExpired)

	globals.GoLogger.Infof("API SET GCACHE %s DATA LEN %d", utils.GetHotCacheKey(globals.JueJinFlag), len(resultResp.Data))

	c.JSON(http.StatusOK, resultResp)
}

func JueJinAIBox(c *gin.Context) {

	if cacheResult, found := globals.GoCache.Get(utils.GetHotCacheKey(globals.JueJinAIBoxFlag)); found {
		globals.GoLogger.Infof("API GET GCACHE %s", utils.GetHotCacheKey(globals.JueJinAIBoxFlag))

		c.JSON(http.StatusOK, cacheResult)
		return
	}

	// 统一输出结果
	var resultResp globals.GblResp

	resp, err := http.Get("https://api.juejin.cn/content_api/v1/content/article_rank?category_id=6809637773935378440&type=hot")
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

	var shellResp JueJinShellResponse
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

		newData.Title = v.Content.Title
		newData.Desc = ""
		newData.HotVal = strconv.Itoa(v.ContentCounter.HotRank)
		newData.Pos = k + 1
		newData.ToUrl = fmt.Sprintf("https://juejin.cn/post/%s", v.Content.ContentId)

		resultResp.Data = append(resultResp.Data, newData)
	}

	globals.GoCache.Set(utils.GetHotCacheKey(globals.JueJinAIBoxFlag), resultResp, globals.HotCacheExpired)

	globals.GoLogger.Infof("API SET GCACHE %s DATA LEN %d", utils.GetHotCacheKey(globals.JueJinAIBoxFlag), len(resultResp.Data))

	c.JSON(http.StatusOK, resultResp)
}
