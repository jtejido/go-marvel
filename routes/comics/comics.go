package comics

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

var (
	Cache *cache.Cache = cache.New()
)

func comicsHandler(conf *config.Config) func(c *gin.Context) {
	return func(c *gin.Context) {
		marvelService := service.GetService(conf.Key, conf.Secret, conf.Timeout)
		var etagStr string
		etag, ok := c.Get("etag")
		if ok {
			etagStr = etag.(string)
		}

		res, err := marvelService.Comics(etagStr, c.Query("limit"), c.Query("offset"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, core.NewErrorResponseWithCode(err.Error(), 500))
			return
		}

		c.JSON(http.StatusOK, res)
	}
}

func comicHandler(conf *config.Config) func(c *gin.Context) {
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

		res, err := marvelService.ComicByID(etagStr, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, core.NewErrorResponseWithCode(err.Error(), 500))
			return
		}

		c.JSON(http.StatusOK, res)
	}
}

func comicCharactersHandler(conf *config.Config) func(c *gin.Context) {
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

		res, err := marvelService.ComicCharactersByID(etagStr, id, c.Query("limit"), c.Query("offset"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, core.NewErrorResponseWithCode(err.Error(), 500))
			return
		}

		c.JSON(http.StatusOK, res)
	}
}

func Setup(ctx context.Context, conf *config.Config, g *gin.RouterGroup) {
	g.GET("", middleware.ETagCacheCheck(Cache), middleware.BodyCacheCheck(Cache), comicsHandler(conf))
	g.GET("/:id", middleware.ETagCacheCheck(Cache), middleware.BodyCacheCheck(Cache), comicHandler(conf))
	g.GET("/:id/characters", middleware.ETagCacheCheck(Cache), middleware.BodyCacheCheck(Cache), comicCharactersHandler(conf))
}
