package middleware

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jtejido/go-marvel/cache"
	"net/http"
)

type Response struct {
	Etag string `json:"etag"`
}

// this holds etag from Marvel API per URI, etags will be checked from service when it's modified or not.
type etagCacheWriter struct {
	gin.ResponseWriter
	cache      *cache.Cache
	requestURI string
}

func (w etagCacheWriter) Write(b []byte) (int, error) {
	status := w.Status()
	var res Response
	json.Unmarshal(b, &res)
	if status == http.StatusOK {
		w.cache.Set(w.requestURI, []byte(res.Etag))
	}
	return w.ResponseWriter.Write(b)
}

func ETagCacheCheck(cache *cache.Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		etag, err := cache.Get(c.Request.RequestURI)
		if err == nil {
			c.Set("etag", string(etag))
			c.Next()
		} else {
			c.Writer = etagCacheWriter{cache: cache, requestURI: c.Request.RequestURI, ResponseWriter: c.Writer}
			c.Next()
		}
	}
}
