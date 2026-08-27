package _123

import (
	"net/url"
	"reflect"
	"testing"
	"time"
)

func TestEffectiveProtocol(t *testing.T) {
	for _, tc := range []struct {
		protocol string
		android  bool
	}{
		{"", false},
		{"web", false},
		{"invalid", false},
		{"android", true},
		{" Android ", true},
		{"ANDROID", true},
	} {
		d := &Pan123{Addition: Addition{Protocol: tc.protocol}}
		if got := d.isAndroid(); got != tc.android {
			t.Fatalf("protocol %q: got android=%v, want %v", tc.protocol, got, tc.android)
		}
	}
}

func TestProtocolEndpoints(t *testing.T) {
	wantWeb := pan123Endpoints{
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
	wantAndroid := pan123Endpoints{
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
	if !reflect.DeepEqual(webEndpoints, wantWeb) {
		t.Fatalf("web endpoints changed:\n got: %#v\nwant: %#v", webEndpoints, wantWeb)
	}
	if !reflect.DeepEqual(androidEndpoints, wantAndroid) {
		t.Fatalf("android endpoints changed:\n got: %#v\nwant: %#v", androidEndpoints, wantAndroid)
	}
}

func TestProtocolHeaders(t *testing.T) {
	web := (&Pan123{Addition: Addition{Platform: "web"}}).webRequestHeaders(&webEndpoints, "token")
	wantWeb := map[string]string{
		"origin":        "https://yun.123pan.com",
		"referer":       "https://yun.123pan.com/",
		"authorization": "Bearer token",
		"user-agent":    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) openlist-client",
		"platform":      "web",
		"app-version":   "3",
	}
	if !reflect.DeepEqual(web, wantWeb) {
		t.Fatalf("web headers changed: got %#v, want %#v", web, wantWeb)
	}

	profile := androidProfile{
		LoginUUID:  "0123456789abcdef0123456789abcdef",
		OSVersion:  "Android_13",
		DeviceType: "M2102J20SG",
		DeviceName: "Xiaomi",
	}
	loginHeaders := androidHeaders(profile, "")
	if _, ok := loginHeaders["authorization"]; ok {
		t.Fatal("Android login must not send an empty Authorization header")
	}
	requestHeaders := androidHeaders(profile, "token")
	wantAndroid := map[string]string{
		"content-type":  "application/json",
		"authorization": "Bearer token",
		"LoginUuid":     profile.LoginUUID,
		"user-agent":    "123pan/v2.4.0(Android_13;Xiaomi)",
		"osversion":     "Android_13",
		"platform":      "android",
		"devicetype":    "M2102J20SG",
		"devicename":    "Xiaomi",
		"app-version":   "61",
		"x-app-version": "2.4.0",
	}
	if !reflect.DeepEqual(requestHeaders, wantAndroid) {
		t.Fatalf("Android headers changed: got %#v, want %#v", requestHeaders, wantAndroid)
	}
}

func TestAndroidProfileLifetime(t *testing.T) {
	firstDriver := &Pan123{}
	first := firstDriver.getAndroidProfile()
	again := firstDriver.getAndroidProfile()
	if first != again {
		t.Fatalf("profile changed within driver lifetime: %#v != %#v", first, again)
	}
	if len(first.LoginUUID) != 32 {
		t.Fatalf("LoginUuid length = %d, want 32", len(first.LoginUUID))
	}
	second := (&Pan123{}).getAndroidProfile()
	if first.LoginUUID == second.LoginUUID {
		t.Fatal("different driver instances unexpectedly share LoginUuid")
	}
}

func TestLoginSpec(t *testing.T) {
	tests := []struct {
		name     string
		driver   *Pan123
		wantURL  string
		wantBody map[string]interface{}
	}{
		{
			name:     "web passport",
			driver:   &Pan123{Addition: Addition{Username: "13800000000", Password: "secret", Protocol: "web"}},
			wantURL:  webEndpoints.SignIn,
			wantBody: map[string]interface{}{"passport": "13800000000", "password": "secret", "remember": true},
		},
		{
			name:     "android passport",
			driver:   &Pan123{Addition: Addition{Username: "13800000000", Password: "secret", Protocol: "android"}},
			wantURL:  androidEndpoints.SignIn,
			wantBody: map[string]interface{}{"passport": "13800000000", "password": "secret", "type": 1},
		},
		{
			name:     "android email",
			driver:   &Pan123{Addition: Addition{Username: "user@example.com", Password: "secret", Protocol: "android"}},
			wantURL:  androidEndpoints.SignIn,
			wantBody: map[string]interface{}{"mail": "user@example.com", "password": "secret", "type": 2},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotURL, _, gotBody := tc.driver.loginSpec()
			if gotURL != tc.wantURL {
				t.Fatalf("URL = %q, want %q", gotURL, tc.wantURL)
			}
			for key, want := range tc.wantBody {
				if got := gotBody[key]; !reflect.DeepEqual(got, want) {
					t.Fatalf("body[%q] = %#v, want %#v", key, got, want)
				}
			}
			if len(gotBody) != len(tc.wantBody) {
				t.Fatalf("body has unexpected fields: %#v", gotBody)
			}
		})
	}
}

func TestSignPathFixedVector(t *testing.T) {
	now := time.Date(2026, 1, 2, 3, 4, 0, 0, time.FixedZone("UTC+8", 8*60*60))
	key, value := signPathAt("/b/api/file/list/new", "web", "3", now, "1234567")
	if key != "1188383987" || value != "1767294240-1234567-2261425176" {
		t.Fatalf("signature = %s=%s", key, value)
	}
}

func TestRequestURLProtocol(t *testing.T) {
	raw := "https://www.123pan.com/api/file/list/new?existing=value"
	android := (&Pan123{Addition: Addition{Protocol: "android"}}).requestURL(raw)
	if android != raw {
		t.Fatalf("Android URL was modified: %q", android)
	}

	webRaw := "https://yun.123pan.com/b/api/file/list/new?existing=value"
	web := (&Pan123{Addition: Addition{Protocol: "web"}}).requestURL(webRaw)
	parsed, err := url.Parse(web)
	if err != nil {
		t.Fatal(err)
	}
	query := parsed.Query()
	if query.Get("existing") != "value" {
		t.Fatalf("existing query was lost: %q", web)
	}
	if len(query) != 2 {
		t.Fatalf("signed Web URL should contain the original and signature query: %q", web)
	}
}
