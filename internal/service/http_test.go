package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/nevercase/lunara-k8s/configs"
)

type FakeResponse struct {
	Success       httpResponse
	NoAuth        httpResponse
	SignedIn      httpResponse
	LoginFailed   httpResponse
	UnknownFailed httpResponse
}

var fakeResponse = FakeResponse{
	Success: httpResponse{
		ErrorCode: errNone,
		Message:   msgSuccess,
	},
	NoAuth: httpResponse{
		ErrorCode: errNoAuth,
		Message:   msgNoAuth,
	},
	SignedIn: httpResponse{
		ErrorCode: errSignedIn,
		Message:   msgSuccess,
	},
	LoginFailed: httpResponse{
		ErrorCode: errLoginFailed,
		Message:   msgFailed,
	},
	UnknownFailed: httpResponse{
		ErrorCode: errUnknown,
		Message:   msgFailed,
	},
}

var fakeConf = &configs.Config{
	Http: configs.HttpConfig{
		IP:   "0.0.0.0",
		Port: 9081,
	},
}

var domain = fmt.Sprintf("http://127.0.0.1:%d", fakeConf.Http.Port)
var sigInUrl = fmt.Sprintf("%s%s", domain, _routerSignIn)
var sigOutUrl = fmt.Sprintf("%s%s", domain, _routerSignOut)
var cookieTmp []*http.Cookie

func TestService_InitHttpServer(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	type fields struct {
		c           *configs.Config
		output      *os.File
		httpService *httpService
		ctx         context.Context
		cancel      context.CancelFunc
	}
	tests := []struct {
		name   string
		fields fields
		want   *httpService
	}{
		{
			name: "case1",
			fields: fields{
				c:      fakeConf,
				ctx:    ctx,
				cancel: cancel,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				c:           tt.fields.c,
				output:      tt.fields.output,
				httpService: tt.fields.httpService,
				ctx:         tt.fields.ctx,
				cancel:      tt.fields.cancel,
			}
			server := s.InitHttpServer()
			time.Sleep(time.Second * 1)
			// test sign in failed
			params := url.Values{}
			params.Add("account", "admin")
			params.Add("password", "12345611")
			if got, err := httpPost(sigInUrl, params); err != nil ||
				!reflect.DeepEqual(got, fakeResponse.LoginFailed) {
				t.Errorf("http request %s = %v, want %v, err (%v)", _routerSignIn, got, fakeResponse.LoginFailed, err)
			}
			// test sign in success
			params = url.Values{}
			params.Add("account", "admin")
			params.Add("password", "123456")
			if got, err := httpPost(sigInUrl, params); err != nil ||
				!reflect.DeepEqual(got, fakeResponse.Success) {
				t.Errorf("http request %s = %v, want %v, err (%v)", _routerSignIn, got, fakeResponse.Success, err)
			}
			// test signed in
			params = url.Values{}
			params.Add("account", "admin")
			params.Add("password", "123456")
			if got, err := httpPost(sigInUrl, params); err != nil ||
				!reflect.DeepEqual(got, fakeResponse.SignedIn) {
				t.Errorf("http request %s = %v, want %v, err (%v)", _routerSignIn, got, fakeResponse.SignedIn, err)
			}
			// test sign out success
			if got, err := httpGet(sigOutUrl); err != nil ||
				!reflect.DeepEqual(got, fakeResponse.Success) {
				t.Errorf("http request %s = %v, want %v, err (%v)", _routerSignOut, got, fakeResponse.Success, err)
			}
			cookieTmp = []*http.Cookie{}
			//test sign out failed
			if got, err := httpGet(sigOutUrl); err != nil ||
				!reflect.DeepEqual(got, fakeResponse.UnknownFailed) {
				t.Errorf("http request %s = %v, want %v, err (%v)", _routerSignOut, got, fakeResponse.UnknownFailed, err)
			}
			//server.ShutDown()
			_ = server
		})
	}
}

func httpGet(targetUrl string) (ret httpResponse, err error) {
	var (
		req     *http.Request
		res     *http.Response
		content []byte
	)
	client := &http.Client{
		Timeout: time.Second * 30,
	}
	req, _ = http.NewRequest("GET", targetUrl, strings.NewReader(""))
	if len(cookieTmp) > 0 {
		for _, v := range cookieTmp {
			req.AddCookie(v)
		}
	}
	fmt.Printf("req targetUrl:%s cookies:%v\n", targetUrl, req.Cookies())
	req.Header.Set("Content-Type", "application/json")
	if res, err = client.Do(req); err != nil {
		return ret, err
	}
	fmt.Printf("res targetUrl:%s cookies:%v\n", targetUrl, res.Cookies())
	if content, err = ioutil.ReadAll(res.Body); err != nil {
		return ret, err
	}
	if err = res.Body.Close(); err != nil {
		return ret, err
	}
	if err = json.Unmarshal(content, &ret); err != nil {
		return ret, err
	}
	return ret, nil
}

func httpPost(targetUrl string, params url.Values) (ret httpResponse, err error) {
	var (
		req     *http.Request
		res     *http.Response
		content []byte
	)
	client := &http.Client{
		Timeout: time.Second * 30,
	}
	req, _ = http.NewRequest("POST", targetUrl, strings.NewReader(params.Encode()))
	if len(cookieTmp) > 0 {
		for _, v := range cookieTmp {
			req.AddCookie(v)
		}
	}
	fmt.Printf("req targetUrl:%s cookies:%v\n", targetUrl, req.Cookies())
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")
	if res, err = client.Do(req); err != nil {
		return ret, err
	}
	fmt.Printf("res targetUrl:%s cookies:%v\n", targetUrl, res.Cookies())
	if content, err = ioutil.ReadAll(res.Body); err != nil {
		return ret, err
	}
	if len(res.Cookies()) > 0 {
		cookieTmp = res.Cookies()
	}
	fmt.Println("content:", string(content))
	if err = res.Body.Close(); err != nil {
		return ret, err
	}
	if err = json.Unmarshal(content, &ret); err != nil {
		return ret, err
	}
	return ret, nil
}

//func Test_httpService_ShutDown(t *testing.T) {
//	type fields struct {
//		server *http.Server
//	}
//	tests := []struct {
//		name   string
//		fields fields
//	}{
//		{
//			name: "case2",
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			h := &httpService{
//				server: tt.fields.server,
//			}
//			h.ShutDown()
//		})
//	}
//}
