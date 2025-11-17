package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/turbo-uid/hots/globals"
	"github.com/turbo-uid/hots/pkg/trending"
	"github.com/turbo-uid/hots/utils"
)

func EnDataHot(c *gin.Context) {
	var EnDataFlag, EnDataType string

	s := c.DefaultQuery("t", "m")
	if s == "s" {
		EnDataFlag = globals.EnDataSFlag
		EnDataType = "1"
	} else {
		EnDataFlag = globals.EnDataMFlag
		EnDataType = "0"
	}

	if cacheResult, found := globals.GoCache.Get(utils.GetHotCacheKey(EnDataFlag)); found {
		globals.GoLogger.Infof("API GET GCACHE %s", utils.GetHotCacheKey(EnDataFlag))
		c.JSON(http.StatusOK, cacheResult)
		return
	}

	resultResp, err := trending.GetEnDataHot(EnDataType)
	if err != nil {
		c.JSON(http.StatusOK, resultResp)
		return
	}

	globals.GoCache.Set(utils.GetHotCacheKey(EnDataFlag), resultResp, globals.HotCacheExpired)
	globals.GoLogger.Infof("API SET GCACHE %s DATA LEN %d", utils.GetHotCacheKey(EnDataFlag), len(resultResp.Data))
	c.JSON(http.StatusOK, resultResp)
}

