package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/turbo-uid/hots/globals"
	"github.com/turbo-uid/hots/utils"
)

type ZhihuShellResponse struct {
	SubAppName   string             `json:"subAppName"`
	SpanName     string             `json:"spanName"`
	InitialState ZhInitialStateData `json:"initialState"`
}

type ZhInitialStateData struct {
	Topstory ZhTopstoryData `json:"topstory"`
}

type ZhTopstoryData struct {
	HotList []ZhHotData `json:"hotList"`
}

type ZhHotData struct {
	ZhId      string       `json:"id"`
	ZhType    string       `json:"type"`
	StyleType string       `json:"styleType"`
	CardId    string       `json:"cardId"`
	Target    ZhTargetData `json:"target"`
}

type ZhTargetData struct {
	TitleArea   ZhHotTitle   `json:"titleArea"`
	ExcerptArea ZhHotExcerpt `json:"excerptArea"`
	ImageArea   ZhHotImage   `json:"imageArea"`
	MetricsArea ZhHotMetrics `json:"metricsArea"`
	LabelArea   ZhHotLabel   `json:"labelArea"`
	Link        ZhHotLink    `json:"link"`
}

type ZhHotTitle struct {
	Text string `json:"text"`
}
type ZhHotExcerpt struct {
	Text string `json:"text"`
}
type ZhHotImage struct {
	Text string `json:"url"`
}
type ZhHotMetrics struct {
	Text string `json:"text"`
}
type ZhHotLabel struct {
	Trend int `json:"trend"`
}
type ZhHotLink struct {
	Text string `json:"url"`
}

func ZhihuByHtmlHot(c *gin.Context) {
	// 统一输出结果
	var resultResp globals.GblResp

	resp, err := http.Get("https://www.zhihu.com/billboard")
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

	doc.Find(".Card .HotList-item").Each(func(i int, s *goquery.Selection) {

		hotScore := s.Find(".HotList-itemMetrics").Text()
		trim_hot := strings.TrimSpace(hotScore)
		titleTxt := s.Find(".HotList-itemTitle").Text()

		trim_title := strings.TrimSpace(titleTxt)
		pos := s.Find(".HotList-itemIndex").Text()
		trim_pos := strings.TrimSpace(pos)

		src, _ := s.Find(".HotList-itemImgContainer img").Attr("src")

		var num_pos int
		if len(trim_pos) == 0 {
			num_pos = 0
		} else {
			num_pos, _ = strconv.Atoi(trim_pos)
		}

		var newData globals.GblRespData

		newData.Title = trim_title
		newData.HotVal = trim_hot
		newData.Desc = ""
		newData.ToUrl = ""
		newData.Pos = num_pos
		// newData.Lab = trim_lab
		newData.Icon = src

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

	// 将解析后的数据作为响应返回给客户端
	c.JSON(http.StatusOK, resultResp)
}

func ZhihuByJsonHot(c *gin.Context) {

	if cacheResult, found := globals.GoCache.Get(utils.GetHotCacheKey(globals.ZhihuFlag)); found {
		globals.GoLogger.Infof("API GET GCACHE %s", utils.GetHotCacheKey(globals.ZhihuFlag))

		c.JSON(http.StatusOK, cacheResult)
		return
	}

	var resultResp globals.GblResp

	resp, err := http.Get("https://www.zhihu.com/billboard")
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
		resultResp.Err = "Error reading the response body: " + err.Error()
		c.JSON(http.StatusOK, resultResp)
		return
	}

	rex := regexp.MustCompile(`<script id="js-initialData" type="text\/json">(.*?)<\/script>`)

	match := rex.FindStringSubmatch(string(body))

	if len(match) > 1 {
		// 解析JSON响应
		var shellResp ZhihuShellResponse

		jsonData := []byte(match[1])
		err = json.Unmarshal(jsonData, &shellResp)
		if err != nil {
			resultResp.Code = 1
			resultResp.Err = "Failed to parse JSON: " + err.Error()
			c.JSON(http.StatusOK, resultResp)
			return
		}

		for k, v := range shellResp.InitialState.Topstory.HotList {
			var newData globals.GblRespData

			newData.Title = v.Target.TitleArea.Text
			newData.Desc = v.Target.ExcerptArea.Text // 此数据过长，缓存的话可以比较消耗内存，可以考虑不要
			newData.HotVal = v.Target.MetricsArea.Text
			newData.Pos = k + 1
			newData.ToUrl = v.Target.Link.Text
			// newData.IsTop = 0
			newData.Icon = v.Target.ImageArea.Text
			resultResp.Data = append(resultResp.Data, newData)
		}

		if len(resultResp.Data) > 0 {
			resultResp.Succ = "ok"
			resultResp.Code = 0
		}
	} else {
		resultResp.Code = 1
		resultResp.Err = "No match found in the fetched HTML content."
	}

	globals.GoCache.Set(utils.GetHotCacheKey(globals.ZhihuFlag), resultResp, globals.HotCacheExpired)

	globals.GoLogger.Infof("API SET GCACHE %s DATA LEN %d", utils.GetHotCacheKey(globals.ZhihuFlag), len(resultResp.Data))

	c.JSON(http.StatusOK, resultResp)
}
