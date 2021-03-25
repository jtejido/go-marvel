package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/jtejido/go-marvel/config"
	"github.com/jtejido/go-marvel/core"
	"github.com/jtejido/go-marvel/routes/characters"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"time"
)

var (
	wg       *sync.WaitGroup
	signals  chan os.Signal
	confPath string
)

func init() {
	signals = make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	wg = &sync.WaitGroup{}
	flag.StringVar(&confPath, "config", "", "Location of the config file.")
}

func routes(conf *config.Config, router *gin.Engine, close func()) {
	store := cookie.NewStore([]byte(conf.Token))
	router.Use(sessions.Sessions("marvel", store))
	router.Use(core.CORS)

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, core.NewErrorResponse("Route not found"))
	})

	ctx, cancel := context.WithCancel(context.WithValue(context.TODO(), "wg", wg))

	characters.Setup(ctx, conf, router.Group("/characters", core.CORS))

	<-signals

	cancel()
	wg.Wait()
	close()

	fmt.Println("Server is shutdown")
}

// client main
func main() {
	flag.Parse()

	runtime.GOMAXPROCS(runtime.NumCPU())

	gin.SetMode(gin.ReleaseMode)
	conf, err := config.LoadConfig(confPath)
	if err != nil {
		panic(err)
	}
	router := gin.New()

	listen := conf.Listen

	if conf.Debug {
		file, err := os.OpenFile("/tmp/api.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal(err)
		}

		log.SetOutput(file)
	}

	srv := &http.Server{
		ReadTimeout: 15 * time.Second,
		IdleTimeout: 120 * time.Second,
		Handler:     router,
		Addr:        listen,
	}

	go routes(conf, router, func() {
		srv.Close()
	})

	fmt.Printf("Server is listening at %s\n", listen)

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
