// Copyright 2015 mint.zhao.chiu@gmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.
package sms

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/aiwuTech/devKit/convert"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	_BaseUrl_YTX = "app.cloopen.com:8883" // 服务http地址
	_Version_YTX = "2013-12-26"           // 服务版本号
)

type YuntongxunService struct {
	account string
	token   string
	appId   string
}

func NewYuntongxun(account, token, appId string) *YuntongxunService {
	return &YuntongxunService{
		account: account,
		token:   token,
		appId:   appId,
	}
}

func (this *YuntongxunService) getSigAuth() (string, string) {
	date := time.Now()
	sig := getMd5String([]byte(fmt.Sprintf("%s%s%s", this.account, this.token, date.Format("20060102150405"))))
	auth := getBase64String([]byte(fmt.Sprintf("%s:%s", this.account, date.Format("20060102150405"))))
	return sig, auth
}

func (this *YuntongxunService) GetUserInfo() (*SmsUser, error) {
	return nil, nil
}

func (this *YuntongxunService) SendSMS(string, []string) (*SmsResult, error) {
	return nil, nil
}

type SendSMSRequest struct {
	AppId      string   `json:"appId"`
	To         string   `json:"to"`
	TemplateId string   `json:"templateId"`
	Datas      []string `json:"datas"`
}

type SendSMSResponse struct {
	StatusCode string `json:"statusCode"`
	StatusMsg  string `json:"statusMsg"`
}

func (this *YuntongxunService) SendSMS_Tpl(templateId int64, tos []string, args []string) (*SmsResult, error) {
	request := SendSMSRequest{
		AppId:      this.appId,
		TemplateId: convert.Int642str(templateId),
		To:         strings.Join(tos, ","),
		Datas:      args,
	}

	sig, auth := this.getSigAuth()
	values := url.Values{}
	values.Add("sig", sig)

	u := &url.URL{
		Scheme:   "https",
		Host:     _BaseUrl_YTX,
		Path:     fmt.Sprintf("/%s/Accounts/%s/SMS/TemplateSMS", _Version_YTX, this.account),
		RawQuery: values.Encode(),
	}

	body, error := json.Marshal(request)
	if error != nil {
		return nil, error
	}

	req, error := http.NewRequest("POST", u.String(), bytes.NewReader(body))
	if error != nil {
		log.Errorf("yuntongxun request err: %v", error)
		return nil, error
	}
	req.Header.Add("Authorization", auth)
	req.Header.Add("Content-Type", "application/json;charset=utf-8;")
	req.Header.Add("Accept", "application/json;")

	httpClient := http.Client{}
	res, error := httpClient.Do(req)
	if error != nil {
		return nil, error
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(res.Status)
	}

	data, error := ioutil.ReadAll(res.Body)
	if error != nil {
		return nil, error
	}
	response := &SendSMSResponse{}
	error = json.Unmarshal(data, response)
	if error != nil {
		return nil, error
	}
	if response.StatusCode != "000000" {
		return nil, errors.New(response.StatusMsg)
	}

	return &SmsResult{}, nil
}

func getMd5String(data []byte) string {
	h := md5.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

func getBase64String(data []byte) string {
	h := base64.StdEncoding
	return h.EncodeToString(data)
}
