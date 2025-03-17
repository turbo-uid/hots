package routers

import (
	"github.com/turbo-uid/hots/routers/api"
	"github.com/turbo-uid/hots/routers/middlewares"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Recovery())

	r.Use(middlewares.CorsReq())

	r.Use(middlewares.Logger())

	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/hot/bili", api.BiliHot)
		apiGroup.GET("/hot/weibo", api.WeiboHot)
		apiGroup.GET("/hot/douyin", api.DouyinHot)
		apiGroup.GET("/hot/toutiao", api.ToutiaoHot)
		apiGroup.GET("/hot/douban", api.DoubanHot)
		apiGroup.GET("/hot/thepaper", api.ThepaperHot)
		apiGroup.GET("/hot/xhs", api.XhsHot)
		apiGroup.GET("/hot/wy163", api.Wy163Hot)
		apiGroup.GET("/hot/qq", api.QqHot)
		apiGroup.GET("/hot/baidu", api.BaiduHot)
		apiGroup.GET("/hot/zhihu/v1", api.ZhihuByHtmlHot)
		apiGroup.GET("/hot/zhihu/v2", api.ZhihuByJsonHot)
		apiGroup.GET("/hot/36kr", api.To36krHot)
		apiGroup.GET("/hot/csdn", api.CsdnHot)
		apiGroup.GET("/hot/csdn-content", api.CsdnContent)
		apiGroup.GET("/hot/hellogithub", api.HelloGithubHot)
		apiGroup.GET("/hot/ithome", api.ItHomeHot)
		apiGroup.GET("/hot/juejin", api.JueJinHot)
		apiGroup.GET("/hot/juejin-aibox", api.JueJinAIBox)
		apiGroup.GET("/hot/carhome", api.CarHomeHot)
		apiGroup.GET("/hot/dongchedi", api.DongCheDiHot)
		apiGroup.GET("/hot/cheshi", api.CheShiHot)
		apiGroup.GET("/hot/qctt", api.QcttHot)
		apiGroup.GET("/hot/endata", api.EnDataHot)
		apiGroup.GET("/hot/toolify", api.ToolifyHot)
	}

	return r
}
