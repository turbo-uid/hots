package trending

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/turbo-uid/hots/globals"
)

func GetCheShiHot() (globals.GblResp, error) {
	var resultResp globals.GblResp

	resp, err := http.Get("https://news.cheshi.com/djbd/")
	if err != nil {
		resultResp.Code = 1
		resultResp.Err = "Failed to fetch data"
		return resultResp, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		resultResp.Code = 1
		resultResp.Err = "Failed to read response body"
		return resultResp, err
	}

	doc.Find(".fall_list").Each(func(i int, s *goquery.Selection) {
		if i < 30 {
			titleTxt := s.Find(".list_txt h3 a").Text()
			trim_title := strings.TrimSpace(titleTxt)

			desc := s.Find(".list_txt .txt").Text()
			rex := regexp.MustCompile(`[\s\r\n]+`)
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

	return resultResp, nil
}

