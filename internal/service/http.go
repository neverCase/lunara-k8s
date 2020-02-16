package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type httpService struct {
	server *http.Server
}

const (
	authInActive = iota
	authActive
)

const (
	_authKey       = "auth"
	_routerSignIn  = "/signin"
	_routerSignOut = "/signout"
)

const (
	_cookieTimeout = iota
	_maxAge        = 3600
)

type auth struct {
	Status    int `json:"status"`
	LoginTime int `json:"login_time"`
}

func header() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method               //请求方法
		origin := c.Request.Header.Get("Origin") //请求头部
		var headerKeys []string                  // 声明请求头keys
		for k, v := range c.Request.Header {
			_ = v
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Origin", origin)                                       // 这是允许访问所有域
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE") //服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
			//  header的类型
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			//              允许跨域设置                                                                                                      可以返回其他子段
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar") // 跨域关键设置 让浏览器可以解析
			c.Header("Access-Control-Max-Age", "172800")                                                                                                                                                           // 缓存请求信息 单位为秒
			c.Header("Access-Control-Allow-Credentials", "true")                                                                                                                                                  //  跨域请求是否需要带cookie信息 默认设置为true
			c.Set("content-type", "application/json")                                                                                                                                                              // 设置返回格式是json
		}

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		// 处理请求
		c.Next() //  处理请求
	}
}

func (s *Service) InitHttpServer() *httpService {
	router := gin.New()
	if s.output != nil {
		router.Use(gin.LoggerWithConfig(gin.LoggerConfig{Output: s.output}), gin.RecoveryWithWriter(s.output))
	}
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("sessionStore", store))
	router.Use(header())
	router.POST(_routerSignIn, func(c *gin.Context) {
		s.handleSignIn(c)
	})
	router.GET(_routerSignOut, func(c *gin.Context) {
		s.handleSignOut(c)
	})
	router.GET("/incr", func(c *gin.Context) {
		session := sessions.Default(c)
		var count int
		v := session.Get("count")
		if v == nil {
			count = 0
		} else {
			count = v.(int)
			count++
		}
		session.Set("count", count)
		session.Save()
		c.JSON(200, gin.H{"count": count})
	})
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.c.Http.IP, s.c.Http.Port),
		Handler: router,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				log.Println("The server closed under request err:", err)
			} else {
				log.Fatal("The server closed unexpected err:", err)
			}
		}
	}()
	return &httpService{
		server: server,
	}
}

func (h *httpService) ShutDown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := h.server.Shutdown(ctx); err != nil {
		log.Println("http.Server shutdown err:", err)
	}
}

func (s *Service) handleSignIn(c *gin.Context) {
	session := sessions.Default(c)
	if v := session.Get(_authKey); v != nil {
		var authCache auth
		if err := json.Unmarshal(v.([]byte), &authCache); err != nil {
			log.Printf("handleSignIn Unmarshal data:%v err:%v\n", v.(string), err)
			c.JSON(http.StatusOK, s.getResponse(errUnknown, msgFailed))
			return
		}
		if authCache.Status == authActive {
			c.JSON(http.StatusOK, s.getResponse(errSignedIn, msgSuccess))
			return
		} else {
			log.Println("signed out or auth timeout")
			c.JSON(http.StatusOK, s.getResponse(errUnknown, msgFailed))
			return
		}
	}
	// verify authentication
	login := &loginRequest{
		Account:  c.PostForm("account"),
		Password: c.PostForm("password"),
	}
	if check := s.checkAuth(login); check == false {
		c.JSON(http.StatusOK, s.getResponse(errLoginFailed, msgFailed))
		return
	}
	authData := auth{
		Status:    authActive,
		LoginTime: int(time.Now().Unix()),
	}
	var (
		str []byte
		err error
	)
	if str, err = json.Marshal(authData); err != nil {
		log.Printf("handleSignIn Marshal data:%v err:%v\n", string(str), err)
		c.JSON(http.StatusOK, s.getResponse(errUnknown, msgFailed))
		return
	}
	session.Set(_authKey, str)
	session.Options(sessions.Options{MaxAge: _maxAge})
	session.Save()
	c.JSON(http.StatusOK, s.getResponse(errNone, msgSuccess))
}

func (s *Service) handleSignOut(c *gin.Context) {
	t, err := c.Cookie("sessionStore")
	if err != nil {
		log.Println("get cookie err:", err)
	} else {
		log.Println("handleSignOut cookie:", t)
	}
	session := sessions.Default(c)
	if v := session.Get(_authKey); v != nil {
		var (
			authCache auth
			str       []byte
			err       error
		)
		if err := json.Unmarshal(v.([]byte), &authCache); err != nil {
			log.Printf("handleSignOut Unmarshal data:%v err:%v\n", v.(string), err)
			c.JSON(http.StatusOK, s.getResponse(errUnknown, msgFailed))
			return
		}
		log.Printf("handleSignOut -------- cookie data:%v \n", v)
		authCache.Status = authInActive
		if str, err = json.Marshal(authCache); err != nil {
			log.Printf("handleSignIn Marshal data:%v err:%v\n", string(str), err)
			c.JSON(http.StatusOK, s.getResponse(errUnknown, msgFailed))
			return
		}
		session.Delete(_authKey)
		session.Options(sessions.Options{MaxAge: _cookieTimeout})
		session.Save()
		c.JSON(http.StatusOK, s.getResponse(errNone, msgSuccess))
		return
	} else {
		log.Println("no session")
		c.JSON(http.StatusOK, s.getResponse(errUnknown, msgFailed))
	}
}
