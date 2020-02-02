package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nevercase/lunara-k8s/configs"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"testing"
	"time"
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
		Port: 8081,
	},
}

var domain = fmt.Sprintf("http://127.0.0.1:%d", fakeConf.Http.Port)
var sigInUrl = fmt.Sprintf("%s%s", domain, _routerSignIn)
var sigOutUrl = fmt.Sprintf("%s%s", domain, _routerSignOut)

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
			var client = &http.Client{
				Timeout: 30 * time.Second,
			}
			// test sign in failed
			params := url.Values{}
			params.Add("account", "admin")
			params.Add("password", "12345611")
			if got, err := httpPost(client, sigInUrl, params); err != nil ||
				!reflect.DeepEqual(got, fakeResponse.LoginFailed) {
				t.Errorf("http request %s = %v, want %v, err (%v)", _routerSignIn, got, fakeResponse.LoginFailed, err)
			}
			// test sign in success
			params = url.Values{}
			params.Add("account", "admin")
			params.Add("password", "123456")
			if got, err := httpPost(client, sigInUrl, params); err != nil ||
				!reflect.DeepEqual(got, fakeResponse.Success) {
				t.Errorf("http request %s = %v, want %v, err (%v)", _routerSignIn, got, fakeResponse.Success, err)
			}
			// test sign out success
			if got, err := httpGet(client, sigOutUrl); err != nil ||
				!reflect.DeepEqual(got, fakeResponse.Success) {
				t.Errorf("http request %s = %v, want %v, err (%v)", _routerSignOut, got, fakeResponse.Success, err)
			}
			server.ShutDown()
		})
	}
}

func httpGet(client *http.Client, targetUrl string) (ret httpResponse, err error) {
	var (
		res     *http.Response
		content []byte
	)
	if res, err = client.Get(targetUrl); err != nil {
		return ret, err
	}
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

func httpPost(client *http.Client, targetUrl string, params url.Values) (ret httpResponse, err error) {
	var (
		res     *http.Response
		content []byte
	)
	if res, err = client.PostForm(targetUrl, params); err != nil {
		return ret, err
	}
	if content, err = ioutil.ReadAll(res.Body); err != nil {
		return ret, err
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

func signalHandler() {
	var (
		ch = make(chan os.Signal, 1)
	)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		sig := <-ch
		fmt.Printf("get a signal %s, stop the lunara-k8s service\n", sig.String())
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
