package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"schnegge/internal/base"
	"schnegge/internal/config"
	"time"
)

func doGetRequest(cfg config.Config, url string) (*http.Response, error) {
	client := &http.Client{Timeout: 2 * time.Second}
	req, _ := http.NewRequest("GET", getServer(cfg)+url, nil)
	tokenId, _ := cfg.GetValue(config.TokenID)
	tokenSecret, _ := cfg.GetValue(config.TokenSecret)
	req.Header.Set("x-api-key", tokenId+":"+tokenSecret)
	base.Log.Println("Request", req)
	return client.Do(req)
}

func doPostRequest(cfg config.Config, url string, obj interface{}) (*http.Response, error) {
	j, _ := json.Marshal(obj)
	base.Log.Println("json:", j)
	client := &http.Client{Timeout: 2 * time.Second}
	req, _ := http.NewRequest("POST", getServer(cfg)+url, bytes.NewBuffer(j))
	tokenId, _ := cfg.GetValue(config.TokenID)
	tokenSecret, _ := cfg.GetValue(config.TokenSecret)
	req.Header.Set("x-api-key", tokenId+":"+tokenSecret)
	req.Header.Set("content-type", "application/json")
	req.Header.Set("accept", "*/*")
	base.Log.Println("Request", req)
	return client.Do(req)
}

func getServer(cfg config.Config) string {
	server, ok := cfg.GetValue(config.Server)
	if !ok {
		// default value
		server = "https://salat.hbt.de"
	}
	return server
}

func checkConfig(cfg config.Config) bool {
	result := true
	_, ok := cfg.GetValue(config.TokenID)
	if !ok {
		fmt.Println("Bitte einen Token in Salat erzeugen und die TokenID setzen")
		result = false
	}
	_, ok = cfg.GetValue(config.TokenSecret)
	if !ok {
		fmt.Println("Bitte einen Token in Salat erzeugen und das TokenSecret setzen")
		result = false
	}
	return result
}
