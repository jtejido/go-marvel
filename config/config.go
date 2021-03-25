package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"os"
)

// Our choice of persistent GraphDB is Neo4j, you can create your own config loader for other preferences.
type Config struct {
	*viper.Viper
	Key     string `mapstructure:"key"`
	Secret  string `mapstructure:"secret"`
	Listen  string `mapstructure:"listen"`
	Timeout int    `mapstructure:"timeout"`
	Token   string `mapstructure:"token"`
	Debug   bool   `mapstructure:"debug"`
}

func LoadConfig(filename string) (*Config, error) {
	conf := &Config{Viper: viper.New()}
	var configFile string

	if filename != "" {
		configFile = filename
	} else if os.Getenv("MARVEL_API_CONFIG") != "" {
		configFile = os.Getenv("MARVEL_API_CONFIG")
	} else {
		return nil, fmt.Errorf("Unable to find config file")
	}

	conf.Set("env", os.Getenv("env"))
	conf.SetConfigFile(configFile)
	if err := conf.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := conf.Unmarshal(conf); err != nil {
		return nil, err
	}

	go func() {
		conf.WatchConfig()
		// https://github.com/gohugoio/hugo/blob/master/watcher/batcher.go
		// https://github.com/spf13/viper/issues/609
		// for some reason this fires twice on a Win machine, and the way some editors save files.
		conf.OnConfigChange(func(e fsnotify.Event) {
			log.Println("Configuration has been changed...")
			// only re-read if file has been modified
			if err := conf.ReadInConfig(); err != nil {
				if err == nil {
					log.Println("Reading failed after configuration update: no data was read")
				} else {
					log.Fatalf("Reading failed after configuration update: %s \n", err.Error())
				}

				return
			} else {
				log.Println("Successfully re-read config file...")
			}

		})
	}()
	return conf, nil
}
