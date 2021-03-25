package characters

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jtejido/go-marvel/cache"
	"github.com/jtejido/go-marvel/config"
	"github.com/jtejido/go-marvel/core"
	"github.com/jtejido/go-marvel/middleware"
	"github.com/jtejido/go-marvel/service"
	"net/http"
	"strconv"
)

var etagCache *cache.Cache = cache.New()

func charactersHandler(conf *config.Config) func(c *gin.Context) {
	return func(c *gin.Context) {
		marvelService := service.GetService(conf.Key, conf.Secret, conf.Timeout)
		var etagStr string
		etag, ok := c.Get("etag")
		if ok {
			etagStr = etag.(string)
		}

		res, err := marvelService.All(etagStr, c.Query("limit"), c.Query("offset"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, core.NewErrorResponseWithCode(err.Error(), 500))
			return
		}

		if etagStr != "" {
			if res.Etag != etagStr {
				etagCache.Set(c.Request.RequestURI, etagStr)
			}
		}

		c.JSON(http.StatusOK, res.Result)
	}
}

func characterHandler(conf *config.Config) func(c *gin.Context) {
	return func(c *gin.Context) {
		if c.Param("id") == "" {
			c.JSON(http.StatusBadRequest, core.NewErrorResponseWithCode("id required", 400))
		}

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, core.NewErrorResponseWithCode("invalid id", 400))
		}

		marvelService := service.GetService(conf.Key, conf.Secret, conf.Timeout)
		var etagStr string
		etag, ok := c.Get("etag")
		if ok {
			etagStr = etag.(string)
		}

		res, err := marvelService.ByID(etagStr, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, core.NewErrorResponseWithCode(err.Error(), 500))
			return
		}

		if etagStr != "" {
			if res.Etag != etagStr {
				etagCache.Set(c.Request.RequestURI, etagStr)
			}
		}

		c.JSON(http.StatusOK, res.Result)
	}
}

func Setup(ctx context.Context, conf *config.Config, g *gin.RouterGroup) {
	g.GET("", middleware.ETagCacheCheck(etagCache), charactersHandler(conf))
	g.GET("/:id", middleware.ETagCacheCheck(etagCache), characterHandler(conf))
}
