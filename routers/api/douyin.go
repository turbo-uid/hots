package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/turbo-uid/hots/globals"
	"github.com/turbo-uid/hots/pkg/trending"
	"github.com/turbo-uid/hots/utils"
)

func DouyinHot(c *gin.Context) {
	if cacheResult, found := globals.GoCache.Get(utils.GetHotCacheKey(globals.DouyinFlag)); found {
		globals.GoLogger.Infof("API GET GCACHE %s", utils.GetHotCacheKey(globals.DouyinFlag))
		c.JSON(http.StatusOK, cacheResult)
		return
	}

	resultResp, err := trending.GetDouyinHot()
	if err != nil {
		c.JSON(http.StatusOK, resultResp)
		return
	}

	globals.GoCache.Set(utils.GetHotCacheKey(globals.DouyinFlag), resultResp, globals.HotCacheExpired)
	globals.GoLogger.Infof("API SET GCACHE %s DATA LEN %d", utils.GetHotCacheKey(globals.DouyinFlag), len(resultResp.Data))
	c.JSON(http.StatusOK, resultResp)
}
