package trending

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/turbo-uid/hots/globals"
)

func GetItHomeHot() (globals.GblResp, error) {
	var resultResp globals.GblResp

	resp, err := http.Get("https://m.ithome.com/rankm/")
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

	doc.Find(".rank-box .placeholder").Each(func(i int, s *goquery.Selection) {
		// 这是因为it之家把日榜、周榜、月榜放在了一起
		if i < 10 {
			title := s.Find(".plc-title").Text()
			new_title := strings.TrimSpace(title)
			hot := s.Find(".review-num").Text()
			new_hot := strings.TrimSpace(hot)
			pos := s.Find(".rank-num").Text()
			new_pos := strings.TrimSpace(pos)
			href, _ := s.Find("a").Attr("href")
			cover, _ := s.Find("img").Attr("data-original")

			var newData globals.GblRespData

			newData.Title = new_title
			newData.HotVal = new_hot
			newData.ToUrl = href
			newData.Pos, _ = strconv.Atoi(new_pos)
			newData.Icon = cover

			resultResp.Data = append(resultResp.Data, newData)
		}
	})

	if len(resultResp.Data) > 0 {
		resultResp.Succ = "ok"
		resultResp.Code = 0
	}

	return resultResp, nil
}

