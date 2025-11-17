package trending

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

func GetToolifyHot() (AIResp, error) {
	var resultResp AIResp

	client := &http.Client{}

	req, err := http.NewRequest("GET", ToolifyUrl, nil)
	if err != nil {
		resultResp.Code = 1
		resultResp.Err = "Error creating request"
		return resultResp, err
	}

	resp, err := client.Do(req)
	if err != nil {
		resultResp.Code = 1
		resultResp.Err = "Error making GET request"
		return resultResp, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		resultResp.Code = 1
		resultResp.Err = "Failed to read response body"
		return resultResp, err
	}

	var shellResp ToolifyShellResponse
	err = json.Unmarshal(body, &shellResp)
	if err != nil {
		resultResp.Code = 1
		resultResp.Err = "Failed to parse JSON"
		return resultResp, err
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

	return resultResp, nil
}

func fmtVisitedCount(num int) string {
	if num >= 1e8 {
		return fmt.Sprintf("%.2f亿", float64(num)/1e8)
	} else if num >= 1e4 {
		return fmt.Sprintf("%.2f万", float64(num)/1e4)
	} else {
		return fmt.Sprintf("%.2f", float64(num))
	}
}

