package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/turbo-uid/hots/globals"
	"github.com/turbo-uid/hots/pkg/trending"
	"github.com/turbo-uid/hots/utils"
)

func XhsHot(c *gin.Context) {
	if cacheResult, found := globals.GoCache.Get(utils.GetHotCacheKey(globals.XhsFlag)); found {
		globals.GoLogger.Infof("API GET GCACHE %s", utils.GetHotCacheKey(globals.XhsFlag))
		c.JSON(http.StatusOK, cacheResult)
		return
	}

	resultResp, err := trending.GetXhsHot()
	if err != nil {
		c.JSON(http.StatusOK, resultResp)
		return
	}

	globals.GoCache.Set(utils.GetHotCacheKey(globals.XhsFlag), resultResp, globals.HotCacheExpired)
	globals.GoLogger.Infof("API SET GCACHE %s DATA LEN %d", utils.GetHotCacheKey(globals.XhsFlag), len(resultResp.Data))
	c.JSON(http.StatusOK, resultResp)
}
