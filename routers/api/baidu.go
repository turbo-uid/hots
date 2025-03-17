package api

import (
	"net/http"
	"strconv"
	"strings"

	// "unicode/utf8"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/turbo-uid/hots/globals"
	"github.com/turbo-uid/hots/utils"
)

func BaiduHot(c *gin.Context) {

	if cacheResult, found := globals.GoCache.Get(utils.GetHotCacheKey(globals.BaiduFlag)); found {
		globals.GoLogger.Infof("API GET GCACHE %s", utils.GetHotCacheKey(globals.BaiduFlag))

		c.JSON(http.StatusOK, cacheResult)
		return
	}

	var resultResp globals.GblResp

	resp, err := http.Get("https://top.baidu.com/board?tab=realtime")
	if err != nil {
		resultResp.Code = 1
		resultResp.Err = "Failed to fetch data"
		c.JSON(http.StatusOK, resultResp)
		return
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		resultResp.Code = 1
		resultResp.Err = "Failed to read response body"
		c.JSON(http.StatusOK, resultResp)
	}

	doc.Find(".category-wrap_iQLoo").Each(func(i int, s *goquery.Selection) {

		hotScore := s.Find(".hot-index_1Bl1a").Text()
		trim_hot := strings.TrimSpace(hotScore)
		titleTxt := s.Find(".c-single-text-ellipsis").Text()
		trim_title := strings.TrimSpace(titleTxt)
		pos := s.Find(".index_1Ew5p").Text()
		trim_pos := strings.TrimSpace(pos)
		href, _ := s.Find(".title_dIF3B").Attr("href")

		label := s.Find(".c-text").Text()
		trim_lab := strings.TrimSpace(label)

		desc := s.Find(".small_Uvkd3").Text()
		trim_desc := strings.TrimSpace(desc)
		str_length := len(trim_desc)
		if str_length > 0 {
			trim_desc = trim_desc[0 : str_length-13]
		}

		img_src, _ := s.Find("a img").Attr("src")

		var num_pos int
		if len(trim_pos) == 0 {
			num_pos = 0
		} else {
			num_pos, _ = strconv.Atoi(trim_pos)
		}

		var newData globals.GblRespData

		newData.Title = trim_title
		newData.HotVal = trim_hot
		newData.Desc = trim_desc
		newData.ToUrl = href
		newData.Pos = num_pos
		newData.Lab = trim_lab
		newData.Icon = img_src

		if i == 0 && num_pos == 0 && len(trim_pos) == 0 {
			newData.IsTop = 1
			newData.Pos = 999
			newData.HotVal = "0"
		}

		resultResp.Data = append(resultResp.Data, newData)
	})

	if len(resultResp.Data) > 0 {
		resultResp.Succ = "ok"
		resultResp.Code = 0
	}

	globals.GoCache.Set(utils.GetHotCacheKey(globals.BaiduFlag), resultResp, globals.HotCacheExpired)

	globals.GoLogger.Infof("API SET GCACHE %s DATA LEN %d", utils.GetHotCacheKey(globals.BaiduFlag), len(resultResp.Data))

	c.JSON(http.StatusOK, resultResp)
}
