package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
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

func (s *Service) InitHttpServer() *httpService {
	router := gin.New()
	if s.output != nil {
		router.Use(gin.LoggerWithConfig(gin.LoggerConfig{Output: s.output}), gin.RecoveryWithWriter(s.output))
	}
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("sessionStore", store))
	router.POST(_routerSignIn, func(c *gin.Context) {
		s.handleSignIn(c)
	})
	router.GET(_routerSignOut, func(c *gin.Context) {
		s.handleSignOut(c)
	})
	router.GET("/incr", func(c *gin.Context) {
		//session := sessions.Default(c)
		//if v := session.Get("auth"); v != nil {
		//	auth := v.(auth)
		//	log.Println("session auth:", session)
		//	if auth.status == authActive {
		//		// todo
		//		c.JSON(200, s.getResponse(errNone, msgSuccess))
		//		return
		//	}
		//} else {
		//	log.Println("no session")
		//}
		//c.JSON(http.StatusOK, s.getResponse(errNone, msgSuccess))
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
			str []byte
			err error
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
