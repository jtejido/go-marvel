package middleware

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jtejido/go-marvel/cache"
	"log"
	"net/http"
)

// this holds etag from Marvel API per URI, etags will be checked from service when it's modified or not.
type etagCacheWriter struct {
	bodyCacheWriter
	requestURI string
}

func (w etagCacheWriter) Write(b []byte) (int, error) {
	status := w.Status()
	var res map[string]interface{}
	json.Unmarshal(b, &res)
	if status == http.StatusOK {
		if v, ok := res["etag"]; ok {
			log.Printf("Writing etag: %s in cache.", v.(string))
			w.cache.Set(w.requestURI, []byte(v.(string)))
		}
	}
	return w.bodyCacheWriter.Write(b)
}

func ETagCacheCheck(cache *cache.Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("Looking-up etag for URI: %s in cache.", c.Request.RequestURI)
		etag, err := cache.Get(c.Request.RequestURI)

		if err == nil {
			log.Println("ETag found.")
			c.Set("etag", string(etag))
			c.Next()
		} else {
			log.Println("No ETag found.")
			etagWriter := etagCacheWriter{}
			etagWriter.cache = cache
			etagWriter.requestURI = c.Request.RequestURI
			etagWriter.ResponseWriter = c.Writer
			c.Writer = etagWriter
			c.Next()
		}
	}
}

// this holds body content from Marvel API per URI.
type bodyCacheWriter struct {
	gin.ResponseWriter
	cache *cache.Cache
}

func (w bodyCacheWriter) Write(b []byte) (int, error) {
	status := w.Status()
	var res map[string]interface{}
	json.Unmarshal(b, &res)
	if status == http.StatusOK {
		if v, ok := res["etag"]; ok {
			log.Printf("Writing body for etag: %s in cache.", v.(string))
			w.cache.Set(v.(string), b)
		}
	}
	return w.ResponseWriter.Write(b)
}

func BodyCacheCheck(cache *cache.Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		etag, ok := c.Get("etag")
		if ok {
			log.Printf("Looking-up body for etag: %s in cache.", etag.(string))
			response, err := cache.Get(etag.(string))
			if err == nil {
				log.Println("Body found.")
				c.Data(http.StatusOK, "application/json", response)
				c.Abort()
			} else {
				log.Println("No Body found.")
				c.Writer = bodyCacheWriter{cache: cache, ResponseWriter: c.Writer}
			}
		}

		c.Next()
	}
}
