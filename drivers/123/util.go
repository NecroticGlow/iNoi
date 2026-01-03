package _123

import (
	"context"
	"errors"
	"fmt"
	"hash/crc32"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/OpenListTeam/OpenList/v4/drivers/base"
	"github.com/OpenListTeam/OpenList/v4/pkg/utils"
	"github.com/go-resty/resty/v2"
	jsoniter "github.com/json-iterator/go"
	log "github.com/sirupsen/logrus"
	"github.com/google/uuid"
)

/*
	======== 常量 =========
*/

const (
	LoginApi     = "https://login.123pan.com/api"
	MainApi      = "https://www.123pan.com/b/api"

	SignIn       = LoginApi + "/user/sign_in"
	UserInfo     = MainApi + "/user/info"
	FileList     = MainApi + "/file/list/new"
	DownloadInfo = MainApi + "/file/download_info"
)

/*
	======== Android 签名 =========
*/

func signPath(path, os, version string) (string, string) {
	table := []byte{'a','d','e','f','g','h','l','m','y','i','j','n','o','p','k','q','r','s','t','u','b','c','v','w','s','z'}
	random := fmt.Sprintf("%.f", math.Round(1e7*rand.Float64()))
	now := time.Now().In(time.FixedZone("CST", 8*3600))

	timestamp := fmt.Sprint(now.Unix())
	nowStr := []byte(now.Format("200601021504"))
	for i := range nowStr {
		nowStr[i] = table[nowStr[i]-48]
	}

	timeSign := fmt.Sprint(crc32.ChecksumIEEE(nowStr))
	data := strings.Join([]string{
		timestamp,
		random,
		path,
		os,
		version,
		timeSign,
	}, "|")

	dataSign := fmt.Sprint(crc32.ChecksumIEEE([]byte(data)))
	return timeSign, timestamp + "-" + random + "-" + dataSign
}

func GetApi(raw string) string {
	u, _ := url.Parse(raw)
	q := u.Query()
	q.Add(signPath(u.Path, "android", "2.5.4"))
	u.RawQuery = q.Encode()
	return u.String()
}

/*
	======== 登录 =========
*/

func (d *Pan123) login() error {
	var body base.Json
	if utils.IsEmailFormat(d.Username) {
		body = base.Json{
			"mail": d.Username,
			"password": d.Password,
			"type": 2,
		}
	} else {
		body = base.Json{
			"passport": d.Username,
			"password": d.Password,
			"remember": true,
		}
	}

	res, err := base.RestyClient.R().
		SetHeaders(androidHeaders("")).
		SetBody(body).
		Post(SignIn)
	if err != nil {
		return err
	}

	if utils.Json.Get(res.Body(), "code").ToInt() != 200 {
		return errors.New(utils.Json.Get(res.Body(), "message").ToString())
	}

	d.AccessToken = utils.Json.Get(res.Body(), "data", "token").ToString()
	return nil
}

/*
	======== 请求核心 =========
*/

func androidHeaders(token string) map[string]string {
	h := map[string]string{
		"User-Agent":      "123pan/v2.5.4(Android_14.0.0;Xiaomi)",
		"Accept-Encoding":"gzip",
		"Content-Type":   "application/json",
		"osversion":      "Android_14.0.0",
		"platform":       "android",
		"devicetype":     "M2104K10I",
		"devicename":     "Xiaomi",
		"loginuuid":      uuid.New().String(),
		"App-Version":    "77",
		"X-App-Version":  "2.5.4",
		"Origin":         "https://www.123pan.com",
		"Referer":        "https://www.123pan.com/",
	}
	if token != "" {
		h["Authorization"] = "Bearer " + token
	}
	return h
}

func (d *Pan123) Request(
	url string,
	method string,
	cb base.ReqCallback,
	resp interface{},
) ([]byte, error) {

	isRetry := false

retry:
	req := base.RestyClient.R().
		SetHeaders(androidHeaders(d.AccessToken))

	if cb != nil {
		cb(req)
	}
	if resp != nil {
		req.SetResult(resp)
	}

	res, err := req.Execute(method, GetApi(url))
	if err != nil {
		return nil, err
	}

	body := res.Body()
	code := utils.Json.Get(body, "code").ToInt()

	if code != 0 {
		if code == 401 && !isRetry {
			if err := d.login(); err != nil {
				return nil, err
			}
			isRetry = true
			goto retry
		}
		return nil, errors.New(utils.Json.Get(body, "message").ToString())
	}

	return body, nil
}

/*
	======== 文件列表 =========
*/

func (d *Pan123) getFiles(ctx context.Context, parentId, name string) ([]File, error) {
	page := 1
	var res []File

	for {
		var resp Files
		_, err := d.Request(FileList, http.MethodGet, func(r *resty.Request) {
			r.SetContext(ctx)
			r.SetQueryParams(map[string]string{
				"driveId": "0",
				"limit": "100",
				"Page": strconv.Itoa(page),
				"parentFileId": parentId,
			})
		}, &resp)

		if err != nil {
			return nil, err
		}

		res = append(res, resp.Data.InfoList...)
		if resp.Data.Next == "-1" || len(resp.Data.InfoList) == 0 {
			break
		}
		page++
	}
	return res, nil
}

/*
	======== 用户信息 =========
*/

func (d *Pan123) getUserInfo(ctx context.Context) (*UserInfoResp, error) {
	var resp UserInfoResp
	_, err := d.Request(UserInfo, http.MethodGet, func(r *resty.Request) {
		r.SetContext(ctx)
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
