package service

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
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
	once          sync.Once
)

const (
	charactersUrl string = "http://gateway.marvel.com/v1/public/characters"
	comicsUrl     string = "http://gateway.marvel.com/v1/public/comics"
)

type Service struct {
	key     string
	secret  string
	timeout int
}

func GetService(key, secret string, timeout int) *Service {
	once.Do(func() {
		marvelService = &Service{key, secret, timeout}
	})

	return marvelService
}

func (service *Service) Comics(etag, limit, offset string) (*model.ComicDataWrapper, error) {
	data, err := service.comics(etag, limit, offset)
	if err != nil {
		return nil, err
	}

	res := new(model.ComicDataWrapper)
	res.Etag = etag
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (service *Service) ComicByID(etag string, id int) (*model.ComicDataWrapper, error) {
	data, err := service.comicByID(etag, id)
	if err != nil {
		return nil, err
	}

	res := new(model.ComicDataWrapper)
	res.Etag = etag
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (service *Service) ComicCharactersByID(etag string, id int, limit, offset string) (*model.CharacterDataWrapper, error) {
	data, err := service.comicCharactersByID(etag, id, limit, offset)
	if err != nil {
		return nil, err
	}

	res := new(model.CharacterDataWrapper)
	res.Etag = etag
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (service *Service) Characters(etag, limit, offset string) (*model.CharacterDataWrapper, error) {
	data, err := service.characters(etag, limit, offset)
	if err != nil {
		return nil, err
	}

	res := new(model.CharacterDataWrapper)
	res.Etag = etag
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (service *Service) CharacterByID(etag string, id int) (*model.CharacterDataWrapper, error) {
	data, err := service.characterByID(etag, id)
	if err != nil {
		return nil, err
	}

	res := new(model.CharacterDataWrapper)
	res.Etag = etag
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (service *Service) characterByID(etag string, id int) (data []byte, err error) {
	request, err := http.NewRequest("GET", charactersUrl+"/"+fmt.Sprintf("%d", id), nil)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	respEtags, ok := response.Header["Etag"]
	if !ok || len(respEtags) == 0 {
		log.Println("ETag expected. Aborting.")
		return nil, fmt.Errorf("Invalid response.")
	}

	if response.StatusCode == http.StatusNotModified {
		return nil, err
	}

	data, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("Error reading body: %s", err.Error())
		return nil, err
	}

	return data, nil
}

func (service *Service) comicByID(etag string, id int) (data []byte, err error) {
	request, err := http.NewRequest("GET", comicsUrl+"/"+fmt.Sprintf("%d", id), nil)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	respEtags, ok := response.Header["Etag"]
	if !ok || len(respEtags) == 0 {
		log.Println("ETag expected. Aborting.")
		return nil, fmt.Errorf("Invalid response.")
	}

	if response.StatusCode == http.StatusNotModified {
		return nil, err
	}

	data, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("Error reading body: %s", err.Error())
		return nil, err
	}

	return data, nil
}

func (service *Service) comicCharactersByID(etag string, id int, limit, offset string) (data []byte, err error) {
	request, err := http.NewRequest("GET", comicsUrl+"/"+fmt.Sprintf("%d", id)+"/characters", nil)
	if err != nil {
		return nil, err
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

	request.URL.RawQuery = q.Encode()
	client := http.Client{}

	if service.timeout > 0 {
		client.Timeout = time.Duration(service.timeout) * time.Second
	}

	response, err := client.Do(request)
	if err != nil {
		log.Printf("Error sending request: %s", err.Error())
		return nil, err
	}

	respEtags, ok := response.Header["Etag"]
	if !ok || len(respEtags) == 0 {
		log.Println("ETag expected. Aborting.")
		return nil, fmt.Errorf("Invalid response.")
	}

	if response.StatusCode == http.StatusNotModified {
		return nil, err
	}

	data, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("Error reading body: %s", err.Error())
		return nil, err
	}

	return data, nil
}

func (service *Service) characters(etag, limit, offset string) (data []byte, err error) {
	request, err := http.NewRequest("GET", charactersUrl, nil)
	if err != nil {
		return nil, err
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

	request.URL.RawQuery = q.Encode()
	client := http.Client{}

	if service.timeout > 0 {
		client.Timeout = time.Duration(service.timeout) * time.Second
	}

	response, err := client.Do(request)
	if err != nil {
		log.Printf("Error getting response from API: %s", err.Error())
		return nil, err
	}

	respEtags, ok := response.Header["Etag"]
	if !ok || len(respEtags) == 0 {
		log.Println("ETag expected. Aborting.")
		return nil, fmt.Errorf("Invalid response.")
	}

	if response.StatusCode == http.StatusNotModified {
		return nil, err
	}

	data, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("Error reading body: %s", err.Error())
		return nil, err
	}

	return data, nil
}

func (service *Service) comics(etag, limit, offset string) (data []byte, err error) {
	request, err := http.NewRequest("GET", comicsUrl, nil)
	if err != nil {
		return nil, err
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

	request.URL.RawQuery = q.Encode()
	client := http.Client{}

	if service.timeout > 0 {
		client.Timeout = time.Duration(service.timeout) * time.Second
	}

	response, err := client.Do(request)
	if err != nil {
		log.Printf("Error getting response from API: %s", err.Error())
		return nil, err
	}

	respEtags, ok := response.Header["Etag"]
	if !ok || len(respEtags) == 0 {
		log.Println("ETag expected. Aborting.")
		return nil, fmt.Errorf("Invalid response.")
	}

	if response.StatusCode == http.StatusNotModified {
		return nil, err
	}

	data, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("Error reading body: %s", err.Error())
		return nil, err
	}

	return data, nil
}
