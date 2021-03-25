package service

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/jtejido/go-marvel/cache"
	"github.com/jtejido/go-marvel/model"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	marvelService *Service
	bodyCache     *cache.Cache
	once          sync.Once
)

const (
	charactersUrl string = "http://gateway.marvel.com/v1/public/characters"
)

type Service struct {
	key     string
	secret  string
	timeout int
}

func GetService(key, secret string, timeout int) *Service {
	once.Do(func() {
		marvelService = &Service{key, secret, timeout}
		bodyCache = cache.New()
	})

	return marvelService
}

func (service *Service) All(etag, limit, offset string) (*model.Response, error) {
	newEtag, data, err := service.all(etag, limit, offset)
	if err != nil {
		return nil, err
	}

	res := new(model.Response)
	res.Etag = newEtag
	err = json.Unmarshal(data, &res.Result)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (service *Service) ByID(etag string, id int) (*model.Response, error) {
	newEtag, data, err := service.getCharacter(etag, id)
	if err != nil {
		return nil, err
	}

	res := new(model.Response)
	res.Etag = newEtag
	err = json.Unmarshal(data, &res.Result)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (service *Service) getCharacter(etag string, id int) (newEtag string, data []byte, err error) {
	request, err := http.NewRequest("GET", charactersUrl+"/"+fmt.Sprintf("%d", id), nil)
	if err != nil {
		return etag, nil, err
	}
	if etag != "" {
		request.Header.Set("If-None-Match", etag)
	}
	q := request.URL.Query()
	h := md5.New()
	now := time.Now()
	io.WriteString(h, now.String())
	io.WriteString(h, service.secret)
	io.WriteString(h, service.key)

	q.Add("apikey", service.key)
	q.Add("hash", fmt.Sprintf("%x", h.Sum(nil)))
	q.Add("ts", now.String())

	request.URL.RawQuery = q.Encode()
	client := http.Client{}

	if service.timeout > 0 {
		client.Timeout = time.Duration(service.timeout) * time.Second
	}

	response, err := client.Do(request)
	if err != nil {
		log.Printf("Error sending request: %s", err.Error())
		return etag, nil, err
	}

	respEtags, ok := response.Header["Etag"]
	if !ok || len(respEtags) == 0 {
		log.Println("malformed response. etag expected. Aborting.")
		return etag, nil, fmt.Errorf("Invalid response.")
	}

	if response.StatusCode == http.StatusNotModified {
		log.Println("Looking-up in cache...")
		data, err = bodyCache.Get(etag)
		if err != nil {
			log.Println("Not found in cache. Aborting")
			return etag, nil, err
		}
		return etag, data, err
	}

	data, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("Error reading body: %s", err.Error())
		return etag, nil, err
	}

	if respEtags[0] != etag {
		etag = respEtags[0]
	}

	bodyCache.Set(etag, data)

	return etag, data, nil
}

func (service *Service) all(etag, limit, offset string) (newEtag string, data []byte, err error) {
	request, err := http.NewRequest("GET", charactersUrl, nil)
	if err != nil {
		return etag, nil, err
	}
	if etag != "" {
		request.Header.Set("If-None-Match", etag)
	}

	q := request.URL.Query()
	h := md5.New()
	now := time.Now()
	io.WriteString(h, now.String())
	io.WriteString(h, service.secret)
	io.WriteString(h, service.key)

	if limit != "" {
		q.Add("limit", limit)
	}

	if offset != "" {
		q.Add("offset", offset)
	}

	q.Add("apikey", service.key)
	q.Add("hash", fmt.Sprintf("%x", h.Sum(nil)))
	q.Add("ts", now.String())
	q.Add("orderBy", "modified")

	request.URL.RawQuery = q.Encode()
	client := http.Client{}

	if service.timeout > 0 {
		client.Timeout = time.Duration(service.timeout) * time.Second
	}

	response, err := client.Do(request)
	if err != nil {
		log.Printf("Error getting response from API: %s", err.Error())
		return etag, nil, err
	}

	respEtags, ok := response.Header["Etag"]
	if !ok || len(respEtags) == 0 {
		log.Println("Failed to get etag response. Aborting.")
		return etag, nil, fmt.Errorf("Invalid response.")
	}

	if response.StatusCode == http.StatusNotModified {
		log.Println("Looking-up in cache...")
		data, err = bodyCache.Get(etag)
		if err != nil {
			log.Println("Not found in cache. Aborting")
			return etag, nil, err
		}
		return etag, data, err
	}

	data, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("Error reading body: %s", err.Error())
		return etag, nil, err
	}

	if respEtags[0] != etag {
		etag = respEtags[0]
	}

	bodyCache.Set(etag, data)
	return etag, data, nil
}
