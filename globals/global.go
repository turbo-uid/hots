package globals

import (
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

var GoCache *cache.Cache

var HotCacheExpired time.Duration = 2 * time.Minute

var GoLogger *logrus.Logger

// var GoLogFields = logrus.Fields{}

type GblRespData struct {
	Title  string `json:"title"`
	Desc   string `json:"desc"`
	HotVal string `json:"hot_val"`
	Icon   string `json:"icon"`
	Pos    int    `json:"pos"`
	ToUrl  string `json:"to_url"`
	Lab    string `json:"label"`
	IsTop  int    `json:"is_top"`
}

type GblResp struct {
	Succ string        `json:"succ"`
	Err  string        `json:"err"`
	Code int           `json:"code"`
	Data []GblRespData `json:"data"`
}

// 娱乐榜
var (
	BiliFlag     string = "bilibili"
	WeiboFlag    string = "weibo"
	DouyinFlag   string = "douyin"
	DoubanFlag   string = "douban"
	ThepaperFlag string = "thepaper"
	ToutiaoFlag  string = "toutiao"
	XhsFlag      string = "xiaohongshu"
	Wy163Flag    string = "wy163"
	QqFlag       string = "qq"
	BaiduFlag    string = "baidu"
	ZhihuFlag    string = "zhihu"
	To36krFlag   string = "36kr"
)

// 技术榜
var (
	CsdnFlag        string = "csdn"
	CsdnContentFlag string = "csdn-content"
	HelloGithubFlag string = "hellogithub"
	ItHomeFlag      string = "ithome"
	JueJinFlag      string = "juejin"
	JueJinAIBoxFlag string = "juejin-aibox"
)

// 汽车榜
var (
	CarHomeFlag   string = "carhome"
	DongCheDiFlag string = "dongchedi"
	CheShiFlag    string = "cheshi"
	QcttFlag      string = "qctt"
)

// 票房榜
var (
	EnDataMFlag string = "endata_m"
	EnDataSFlag string = "endata_s"
)

// ai榜
var (
	ToolifyFlag string = "toolify"
)
