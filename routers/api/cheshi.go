package api

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/turbo-uid/hots/globals"
	"github.com/turbo-uid/hots/utils"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

func CheShiHot(c *gin.Context) {

	if cacheResult, found := globals.GoCache.Get(utils.GetHotCacheKey(globals.CheShiFlag)); found {
		globals.GoLogger.Infof("API GET GCACHE %s", utils.GetHotCacheKey(globals.CheShiFlag))

		c.JSON(http.StatusOK, cacheResult)
		return
	}

	var resultResp globals.GblResp

	resp, err := http.Get("https://news.cheshi.com/djbd/")
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

	doc.Find(".fall_list").Each(func(i int, s *goquery.Selection) {

		if i < 30 {
			titleTxt := s.Find(".list_txt h3 a").Text()
			trim_title := strings.TrimSpace(titleTxt)

			desc := s.Find(".list_txt .txt").Text()
			rex := regexp.MustCompile(`[\s\r\n]+`) // 匹配空格、回车换行符
			trim_desc := rex.ReplaceAllString(desc, "")

			img_src, _ := s.Find(".list_img a img").Attr("data-original")
			href_src, _ := s.Find(".list_img a").Attr("href")

			var newData globals.GblRespData

			newData.Title = trim_title
			newData.HotVal = "0"
			newData.Desc = trim_desc
			newData.ToUrl = href_src
			newData.Pos = i + 1
			newData.Lab = ""
			newData.Icon = img_src

			resultResp.Data = append(resultResp.Data, newData)

		}
	})

	if len(resultResp.Data) > 0 {
		resultResp.Succ = "ok"
		resultResp.Code = 0
	}

	globals.GoCache.Set(utils.GetHotCacheKey(globals.CheShiFlag), resultResp, globals.HotCacheExpired)

	globals.GoLogger.Infof("API SET GCACHE %s DATA LEN %d", utils.GetHotCacheKey(globals.CheShiFlag), len(resultResp.Data))

	c.JSON(http.StatusOK, resultResp)
}
