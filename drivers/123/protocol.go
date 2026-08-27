package _123

import (
	cryptorand "crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/OpenListTeam/OpenList/v4/drivers/base"
	"github.com/OpenListTeam/OpenList/v4/pkg/utils"
)

const (
	androidAppVersion  = "61"
	androidXAppVersion = "2.4.0"
	androidDeviceName  = "Xiaomi"
)

var (
	androidDeviceTypes = []string{
		"24075RP89G", "24076RP19G", "24076RP19I", "M1805E10A", "M2004J11G",
		"M2012K11AG", "M2104K10I", "22021211RG", "22021211RI", "21121210G",
		"23049PCD8G", "23049PCD8I", "23013PC75G", "24069PC21G", "24069PC21I",
		"23113RKC6G", "M1912G7BI", "M2007J20CI", "M2007J20CG", "M2007J20CT",
		"M2102J20SG", "M2102J20SI", "21061110AG", "2201116PG", "2201116PI",
		"22041216G", "22041216UG", "22111317PG", "22111317PI", "22101320G",
		"22101320I", "23122PCD1G", "23122PCD1I", "2311DRK48G", "2311DRK48I",
		"2312FRAFDI", "M2004J19PI",
	}
	androidOSVersions = []string{
		"Android_7.1.2", "Android_8.0.0", "Android_8.1.0", "Android_9.0",
		"Android_10", "Android_11", "Android_12", "Android_13",
		"Android_6.0.1", "Android_5.1.1", "Android_4.4.4", "Android_4.3",
		"Android_4.2.2", "Android_4.1.2",
	}
)

type pan123Endpoints struct {
	SignIn            string
	Logout            string
	UserInfo          string
	FileList          string
	DownloadInfo      string
	Mkdir             string
	Move              string
	Rename            string
	Trash             string
	UploadRequest     string
	UploadComplete    string
	S3PreSignedURLs   string
	S3Auth            string
	UploadCompleteV2  string
	S3Complete        string
	OfflineResolve    string
	OfflineSubmit     string
	OfflineTaskList   string
	OfflineTaskDelete string
	Origin            string
	Referer           string
}

var webEndpoints = pan123Endpoints{
	SignIn:            "https://login.123pan.com/api/user/sign_in",
	Logout:            "https://yun.123pan.com/b/api/user/logout",
	UserInfo:          "https://yun.123pan.com/b/api/user/info",
	FileList:          "https://yun.123pan.com/b/api/file/list/new",
	DownloadInfo:      "https://yun.123pan.com/b/api/file/download_info",
	Mkdir:             "https://yun.123pan.com/b/api/file/upload_request",
	Move:              "https://yun.123pan.com/b/api/file/mod_pid",
	Rename:            "https://yun.123pan.com/b/api/file/rename",
	Trash:             "https://yun.123pan.com/b/api/file/trash",
	UploadRequest:     "https://yun.123pan.com/b/api/file/upload_request",
	UploadComplete:    "https://yun.123pan.com/b/api/file/upload_complete",
	S3PreSignedURLs:   "https://yun.123pan.com/b/api/file/s3_repare_upload_parts_batch",
	S3Auth:            "https://yun.123pan.com/b/api/file/s3_upload_object/auth",
	UploadCompleteV2:  "https://yun.123pan.com/b/api/file/upload_complete/v2",
	S3Complete:        "https://yun.123pan.com/b/api/file/s3_complete_multipart_upload",
	OfflineResolve:    "https://yun.123pan.com/b/api/v2/offline_download/task/resolve",
	OfflineSubmit:     "https://yun.123pan.com/b/api/v2/offline_download/task/submit",
	OfflineTaskList:   "https://yun.123pan.com/b/api/offline_download/task/list",
	OfflineTaskDelete: "https://yun.123pan.com/b/api/offline_download/task/delete",
	Origin:            "https://yun.123pan.com",
	Referer:           "https://yun.123pan.com/",
}

var androidEndpoints = pan123Endpoints{
	SignIn:            "https://www.123pan.com/b/api/user/sign_in",
	Logout:            "https://www.123pan.com/b/api/user/logout",
	UserInfo:          "https://www.123pan.com/b/api/user/info",
	FileList:          "https://www.123pan.com/api/file/list/new",
	DownloadInfo:      "https://www.123pan.com/a/api/file/download_info",
	Mkdir:             "https://www.123pan.com/a/api/file/upload_request",
	Move:              "https://www.123pan.com/b/api/file/mod_pid",
	Rename:            "https://www.123pan.com/b/api/file/rename",
	Trash:             "https://www.123pan.com/a/api/file/trash",
	UploadRequest:     "https://www.123pan.com/b/api/file/upload_request",
	UploadComplete:    "https://www.123pan.com/b/api/file/upload_complete",
	S3PreSignedURLs:   "https://www.123pan.com/b/api/file/s3_repare_upload_parts_batch",
	S3Auth:            "https://www.123pan.com/b/api/file/s3_upload_object/auth",
	UploadCompleteV2:  "https://www.123pan.com/b/api/file/upload_complete/v2",
	S3Complete:        "https://www.123pan.com/b/api/file/s3_complete_multipart_upload",
	OfflineResolve:    "https://www.123pan.com/b/api/v2/offline_download/task/resolve",
	OfflineSubmit:     "https://www.123pan.com/b/api/v2/offline_download/task/submit",
	OfflineTaskList:   "https://www.123pan.com/b/api/offline_download/task/list",
	OfflineTaskDelete: "https://www.123pan.com/b/api/offline_download/task/delete",
	Origin:            "https://www.123pan.com",
	Referer:           "https://www.123pan.com/",
}

type androidProfile struct {
	LoginUUID  string
	OSVersion  string
	DeviceType string
	DeviceName string
}

func (d *Pan123) isAndroid() bool {
	return strings.EqualFold(strings.TrimSpace(d.Protocol), "android")
}

func (d *Pan123) endpoints() *pan123Endpoints {
	if d.isAndroid() {
		return &androidEndpoints
	}
	return &webEndpoints
}

func (d *Pan123) getAndroidProfile() androidProfile {
	d.androidProfileOnce.Do(func() {
		randomBytes := make([]byte, 18)
		if _, err := cryptorand.Read(randomBytes); err != nil {
			fallback := fmt.Sprintf("%032x", uint64(time.Now().UnixNano()))
			d.androidProfile = androidProfile{
				LoginUUID:  fallback,
				OSVersion:  androidOSVersions[0],
				DeviceType: androidDeviceTypes[0],
				DeviceName: androidDeviceName,
			}
			return
		}
		d.androidProfile = androidProfile{
			LoginUUID:  hex.EncodeToString(randomBytes[:16]),
			OSVersion:  androidOSVersions[int(randomBytes[16])%len(androidOSVersions)],
			DeviceType: androidDeviceTypes[int(randomBytes[17])%len(androidDeviceTypes)],
			DeviceName: androidDeviceName,
		}
	})
	return d.androidProfile
}

func webLoginHeaders(ep *pan123Endpoints) map[string]string {
	return map[string]string{
		"origin":      ep.Origin,
		"referer":     ep.Referer,
		"user-agent":  "Dart/2.19(dart:io)-openlist",
		"platform":    "web",
		"app-version": "3",
	}
}

func (d *Pan123) webRequestHeaders(ep *pan123Endpoints, token string) map[string]string {
	return map[string]string{
		"origin":        ep.Origin,
		"referer":       ep.Referer,
		"authorization": "Bearer " + token,
		"user-agent":    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) openlist-client",
		"platform":      d.Platform,
		"app-version":   "3",
	}
}

func androidHeaders(profile androidProfile, token string) map[string]string {
	headers := map[string]string{
		"content-type":  "application/json",
		"LoginUuid":     profile.LoginUUID,
		"user-agent":    fmt.Sprintf("123pan/v%s(%s;%s)", androidXAppVersion, profile.OSVersion, profile.DeviceName),
		"osversion":     profile.OSVersion,
		"platform":      "android",
		"devicetype":    profile.DeviceType,
		"devicename":    profile.DeviceName,
		"app-version":   androidAppVersion,
		"x-app-version": androidXAppVersion,
	}
	if token != "" {
		headers["authorization"] = "Bearer " + token
	}
	return headers
}

func (d *Pan123) loginSpec() (string, map[string]string, base.Json) {
	ep := d.endpoints()
	var body base.Json
	if utils.IsEmailFormat(d.Username) {
		body = base.Json{
			"mail":     d.Username,
			"password": d.Password,
			"type":     2,
		}
	} else if d.isAndroid() {
		body = base.Json{
			"passport": d.Username,
			"password": d.Password,
			"type":     1,
		}
	} else {
		body = base.Json{
			"passport": d.Username,
			"password": d.Password,
			"remember": true,
		}
	}
	if d.isAndroid() {
		return ep.SignIn, androidHeaders(d.getAndroidProfile(), ""), body
	}
	return ep.SignIn, webLoginHeaders(ep), body
}

func (d *Pan123) requestHeaders() map[string]string {
	ep := d.endpoints()
	if d.isAndroid() {
		return androidHeaders(d.getAndroidProfile(), d.AccessToken)
	}
	return d.webRequestHeaders(ep, d.AccessToken)
}

func (d *Pan123) requestURL(raw string) string {
	if d.isAndroid() {
		return raw
	}
	return GetApi(raw)
}
