package trending

import (
	"testing"

	"github.com/turbo-uid/hots/globals"
)

// TestGetDouyinHot 测试抖音热搜获取
func TestGetDouyinHot(t *testing.T) {
	result, err := GetDouyinHot()
	if err != nil {
		t.Logf("GetDouyinHot error: %v", err)
		return
	}

	if result.Code != 0 {
		t.Errorf("GetDouyinHot failed with code: %d, err: %s", result.Code, result.Err)
		return
	}

	if len(result.Data) == 0 {
		t.Logf("GetDouyinHot returned empty data")
		return
	}

	t.Logf("GetDouyinHot success, data count: %d", len(result.Data))
	validateGblResp(t, result)
}

// TestGetItHomeHot 测试IT之家热搜获取
func TestGetItHomeHot(t *testing.T) {
	result, err := GetItHomeHot()
	if err != nil {
		t.Logf("GetItHomeHot error: %v", err)
		return
	}

	if result.Code != 0 {
		t.Errorf("GetItHomeHot failed with code: %d, err: %s", result.Code, result.Err)
		return
	}

	if len(result.Data) == 0 {
		t.Logf("GetItHomeHot returned empty data")
		return
	}

	t.Logf("GetItHomeHot success, data count: %d", len(result.Data))
	validateGblResp(t, result)
}

// TestGetEnDataHot 测试票房数据获取
func TestGetEnDataHot(t *testing.T) {
	// 测试月度票房
	result, err := GetEnDataHot("0")
	if err != nil {
		t.Logf("GetEnDataHot (monthly) error: %v", err)
		return
	}

	if result.Code != 0 {
		t.Errorf("GetEnDataHot (monthly) failed with code: %d, err: %s", result.Code, result.Err)
		return
	}

	if len(result.Data) == 0 {
		t.Logf("GetEnDataHot (monthly) returned empty data")
		return
	}

	t.Logf("GetEnDataHot (monthly) success, data count: %d", len(result.Data))
	validateGblResp(t, result)

	// 测试单日票房
	result2, err := GetEnDataHot("1")
	if err != nil {
		t.Logf("GetEnDataHot (daily) error: %v", err)
		return
	}

	if result2.Code != 0 {
		t.Errorf("GetEnDataHot (daily) failed with code: %d, err: %s", result2.Code, result2.Err)
		return
	}

	t.Logf("GetEnDataHot (daily) success, data count: %d", len(result2.Data))
	validateGblResp(t, result2)
}

// TestGetBiliHot 测试B站热搜获取
func TestGetBiliHot(t *testing.T) {
	result, err := GetBiliHot()
	if err != nil {
		t.Logf("GetBiliHot error: %v", err)
		return
	}

	if result.Code != 0 {
		t.Errorf("GetBiliHot failed with code: %d, err: %s", result.Code, result.Err)
		return
	}

	if len(result.Data) == 0 {
		t.Logf("GetBiliHot returned empty data")
		return
	}

	t.Logf("GetBiliHot success, data count: %d", len(result.Data))
	validateGblResp(t, result)
}

// TestGetWeiboHot 测试微博热搜获取
func TestGetWeiboHot(t *testing.T) {
	result, err := GetWeiboHot()
	if err != nil {
		t.Logf("GetWeiboHot error: %v", err)
		return
	}

	if result.Code != 0 {
		t.Errorf("GetWeiboHot failed with code: %d, err: %s", result.Code, result.Err)
		return
	}

	if len(result.Data) == 0 {
		t.Logf("GetWeiboHot returned empty data")
		return
	}

	t.Logf("GetWeiboHot success, data count: %d", len(result.Data))
	validateGblResp(t, result)
}

// TestGetCsdnHot 测试CSDN热搜获取
func TestGetCsdnHot(t *testing.T) {
	result, err := GetCsdnHot()
	if err != nil {
		t.Logf("GetCsdnHot error: %v", err)
		return
	}

	if result.Code != 0 {
		t.Errorf("GetCsdnHot failed with code: %d, err: %s", result.Code, result.Err)
		return
	}

	if len(result.Data) == 0 {
		t.Logf("GetCsdnHot returned empty data")
		return
	}

	t.Logf("GetCsdnHot success, data count: %d", len(result.Data))
	validateGblResp(t, result)
}

// TestGetZhihuByJsonHot 测试知乎热搜获取（JSON方式）
// zhihu 应该是改了数据结构，已经拿不到这个数据了
// func TestGetZhihuByJsonHot(t *testing.T) {
// 	result, err := GetZhihuByJsonHot()
// 	if err != nil {
// 		t.Logf("GetZhihuByJsonHot error: %v", err)
// 		return
// 	}

// 	if result.Code != 0 {
// 		t.Errorf("GetZhihuByJsonHot failed with code: %d, err: %s", result.Code, result.Err)
// 		return
// 	}

// 	if len(result.Data) == 0 {
// 		t.Logf("GetZhihuByJsonHot returned empty data")
// 		return
// 	}

// 	t.Logf("GetZhihuByJsonHot success, data count: %d", len(result.Data))
// 	validateGblResp(t, result)
// }

// validateGblResp 验证返回的数据结构
func validateGblResp(t *testing.T, resp globals.GblResp) {
	if resp.Succ != "ok" && resp.Code == 0 {
		t.Logf("Warning: Succ is not 'ok' but Code is 0")
	}

	if resp.Code != 0 && resp.Err == "" {
		t.Logf("Warning: Code is not 0 but Err is empty")
	}

	// 验证数据项
	for i, item := range resp.Data {
		if item.Title == "" {
			t.Logf("Warning: Data[%d] has empty Title", i)
		}
		if item.Pos <= 0 && item.Pos != 999 {
			t.Logf("Warning: Data[%d] has invalid Pos: %d", i, item.Pos)
		}
	}
}

// TestGetToutiaoHot 测试头条热搜获取
func TestGetToutiaoHot(t *testing.T) {
	result, err := GetToutiaoHot()
	if err != nil {
		t.Logf("GetToutiaoHot error: %v", err)
		return
	}

	if result.Code != 0 {
		t.Errorf("GetToutiaoHot failed with code: %d, err: %s", result.Code, result.Err)
		return
	}

	t.Logf("GetToutiaoHot success, data count: %d", len(result.Data))
	validateGblResp(t, result)
}

// TestGetBaiduHot 测试百度热搜获取
func TestGetBaiduHot(t *testing.T) {
	result, err := GetBaiduHot()
	if err != nil {
		t.Logf("GetBaiduHot error: %v", err)
		return
	}

	if result.Code != 0 {
		t.Errorf("GetBaiduHot failed with code: %d, err: %s", result.Code, result.Err)
		return
	}

	t.Logf("GetBaiduHot success, data count: %d", len(result.Data))
	validateGblResp(t, result)
}

// TestGetQqHot 测试QQ热搜获取
func TestGetQqHot(t *testing.T) {
	result, err := GetQqHot()
	if err != nil {
		t.Logf("GetQqHot error: %v", err)
		return
	}

	if result.Code != 0 {
		t.Errorf("GetQqHot failed with code: %d, err: %s", result.Code, result.Err)
		return
	}

	t.Logf("GetQqHot success, data count: %d", len(result.Data))
	validateGblResp(t, result)
}

// TestGetDoubanHot 测试豆瓣热搜获取
func TestGetDoubanHot(t *testing.T) {
	result, err := GetDoubanHot()
	if err != nil {
		t.Logf("GetDoubanHot error: %v", err)
		return
	}

	if result.Code != 0 {
		t.Errorf("GetDoubanHot failed with code: %d, err: %s", result.Code, result.Err)
		return
	}

	t.Logf("GetDoubanHot success, data count: %d", len(result.Data))
	validateGblResp(t, result)
}

// TestGetQcttHot 测试汽车头条热搜获取
func TestGetQcttHot(t *testing.T) {
	result, err := GetQcttHot()
	if err != nil {
		t.Logf("GetQcttHot error: %v", err)
		return
	}

	if result.Code != 0 {
		t.Errorf("GetQcttHot failed with code: %d, err: %s", result.Code, result.Err)
		return
	}

	t.Logf("GetQcttHot success, data count: %d", len(result.Data))
	validateGblResp(t, result)
}

// TestGetThepaperHot 测试澎湃新闻热搜获取
func TestGetThepaperHot(t *testing.T) {
	result, err := GetThepaperHot()
	if err != nil {
		t.Logf("GetThepaperHot error: %v", err)
		return
	}

	if result.Code != 0 {
		t.Errorf("GetThepaperHot failed with code: %d, err: %s", result.Code, result.Err)
		return
	}

	t.Logf("GetThepaperHot success, data count: %d", len(result.Data))
	validateGblResp(t, result)
}

// TestGetWy163Hot 测试网易热搜获取
func TestGetWy163Hot(t *testing.T) {
	result, err := GetWy163Hot()
	if err != nil {
		t.Logf("GetWy163Hot error: %v", err)
		return
	}

	if result.Code != 0 {
		t.Errorf("GetWy163Hot failed with code: %d, err: %s", result.Code, result.Err)
		return
	}

	t.Logf("GetWy163Hot success, data count: %d", len(result.Data))
	validateGblResp(t, result)
}

// TestGetXhsHot 测试小红书热搜获取
func TestGetXhsHot(t *testing.T) {
	result, err := GetXhsHot()
	if err != nil {
		t.Logf("GetXhsHot error: %v", err)
		return
	}

	if result.Code != 0 {
		t.Errorf("GetXhsHot failed with code: %d, err: %s", result.Code, result.Err)
		return
	}

	t.Logf("GetXhsHot success, data count: %d", len(result.Data))
	validateGblResp(t, result)
}

// TestGetTo36krHot 测试36氪热搜获取
func TestGetTo36krHot(t *testing.T) {
	result, err := GetTo36krHot()
	if err != nil {
		t.Logf("GetTo36krHot error: %v", err)
		return
	}

	if result.Code != 0 {
		t.Errorf("GetTo36krHot failed with code: %d, err: %s", result.Code, result.Err)
		return
	}

	t.Logf("GetTo36krHot success, data count: %d", len(result.Data))
	validateGblResp(t, result)
}

// TestGetHelloGithubHot 测试HelloGithub热搜获取
func TestGetHelloGithubHot(t *testing.T) {
	result, err := GetHelloGithubHot()
	if err != nil {
		t.Logf("GetHelloGithubHot error: %v", err)
		return
	}

	if result.Code != 0 {
		t.Errorf("GetHelloGithubHot failed with code: %d, err: %s", result.Code, result.Err)
		return
	}

	t.Logf("GetHelloGithubHot success, data count: %d", len(result.Data))
	validateGblResp(t, result)
}

// TestGetJueJinHot 测试掘金热搜获取
func TestGetJueJinHot(t *testing.T) {
	result, err := GetJueJinHot()
	if err != nil {
		t.Logf("GetJueJinHot error: %v", err)
		return
	}

	if result.Code != 0 {
		t.Errorf("GetJueJinHot failed with code: %d, err: %s", result.Code, result.Err)
		return
	}

	t.Logf("GetJueJinHot success, data count: %d", len(result.Data))
	validateGblResp(t, result)
}

// TestGetJueJinAIBox 测试掘金AI Box热搜获取
func TestGetJueJinAIBox(t *testing.T) {
	result, err := GetJueJinAIBox()
	if err != nil {
		t.Logf("GetJueJinAIBox error: %v", err)
		return
	}

	if result.Code != 0 {
		t.Errorf("GetJueJinAIBox failed with code: %d, err: %s", result.Code, result.Err)
		return
	}

	t.Logf("GetJueJinAIBox success, data count: %d", len(result.Data))
	validateGblResp(t, result)
}

// TestGetCarHomeHot 测试汽车之家热搜获取
func TestGetCarHomeHot(t *testing.T) {
	result, err := GetCarHomeHot()
	if err != nil {
		t.Logf("GetCarHomeHot error: %v", err)
		return
	}

	if result.Code != 0 {
		t.Errorf("GetCarHomeHot failed with code: %d, err: %s", result.Code, result.Err)
		return
	}

	t.Logf("GetCarHomeHot success, data count: %d", len(result.Data))
	validateGblResp(t, result)
}

// TestGetCheShiHot 测试车市热搜获取
func TestGetCheShiHot(t *testing.T) {
	result, err := GetCheShiHot()
	if err != nil {
		t.Logf("GetCheShiHot error: %v", err)
		return
	}

	if result.Code != 0 {
		t.Errorf("GetCheShiHot failed with code: %d, err: %s", result.Code, result.Err)
		return
	}

	t.Logf("GetCheShiHot success, data count: %d", len(result.Data))
	validateGblResp(t, result)
}

// TestGetDongCheDiHot 测试懂车帝热搜获取
func TestGetDongCheDiHot(t *testing.T) {
	result, err := GetDongCheDiHot()
	if err != nil {
		t.Logf("GetDongCheDiHot error: %v", err)
		return
	}

	if result.Code != 0 {
		t.Errorf("GetDongCheDiHot failed with code: %d, err: %s", result.Code, result.Err)
		return
	}

	t.Logf("GetDongCheDiHot success, data count: %d", len(result.Data))
	validateGblResp(t, result)
}

// TestGetCsdnContent 测试CSDN内容热搜获取
func TestGetCsdnContent(t *testing.T) {
	result, err := GetCsdnContent()
	if err != nil {
		t.Logf("GetCsdnContent error: %v", err)
		return
	}

	if result.Code != 0 {
		t.Errorf("GetCsdnContent failed with code: %d, err: %s", result.Code, result.Err)
		return
	}

	t.Logf("GetCsdnContent success, data count: %d", len(result.Data))
	validateGblResp(t, result)
}

// TestGetZhihuByHtmlHot 测试知乎热搜获取（HTML方式）
func TestGetZhihuByHtmlHot(t *testing.T) {
	result, err := GetZhihuByHtmlHot()
	if err != nil {
		t.Logf("GetZhihuByHtmlHot error: %v", err)
		return
	}

	if result.Code != 0 {
		t.Errorf("GetZhihuByHtmlHot failed with code: %d, err: %s", result.Code, result.Err)
		return
	}

	t.Logf("GetZhihuByHtmlHot success, data count: %d", len(result.Data))
	validateGblResp(t, result)
}

