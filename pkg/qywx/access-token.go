package qywx

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type AccessToken struct {
	Code      int    `json:"errcode"`
	Message   string `json:"errmsg"`
	Token     string `json:"access_token"`
	ExpiresIn int    `json:"expires_in"`
}

const QYAPIEndpoint = "https://qyapi.weixin.qq.com"

var accessToken *AccessToken

func GetAccessToken() string {
	return accessToken.Token
}

func freshAccessToken(agentID, secret string) {
	resp, err := http.Get(fmt.Sprintf(QYAPIEndpoint+"/cgi-bin/gettoken?corpid=%s&corpsecret=%s", agentID, secret))
	if err == nil {
		if bs, err := ioutil.ReadAll(resp.Body); err == nil {
			token := new(AccessToken)
			if err := json.Unmarshal(bs, token); err == nil {
				accessToken = token
			}
		}
	}
}

func FreshTokenTask(agentID, secret string, interval time.Duration) error {

	ticker := time.NewTicker(interval)
	freshAccessToken(agentID, secret)
	for {
		<-ticker.C
		log.Println("Begin to refresh acess token")
		freshAccessToken(agentID, secret)
		log.Printf("New access token is %s", accessToken.Token)
	}
}
