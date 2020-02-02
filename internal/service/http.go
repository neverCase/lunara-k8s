package service

import (
	"context"
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

type auth struct {
	status    int
	loginTime time.Time
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
		c.JSON(http.StatusOK, s.getResponse(errNone, msgSuccess))
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
		auth := v.(auth)
		if auth.status == authActive {
			c.JSON(http.StatusOK, s.getResponse(errSignedIn, msgSuccess))
			return
		}
	}
	// verify authentication
	login := &loginRequest{
		Account:  c.PostForm("account"),
		Password: c.PostForm("password"),
	}
	fmt.Println("login:", login)
	if check := s.checkAuth(login); check == false {
		c.JSON(http.StatusOK, s.getResponse(errLoginFailed, msgFailed))
		return
	}
	authData := auth{
		status:    authInActive,
		loginTime: time.Now(),
	}
	session.Set(_authKey, authData)
	session.Save()
	c.JSON(http.StatusOK, s.getResponse(errNone, msgSuccess))
	session = sessions.Default(c)
	if v := session.Get(_authKey); v != nil {
		auth := v.(auth)
		log.Println("auth:", auth)
	} else {
		log.Println("no session")
	}
}

func (s *Service) handleSignOut(c *gin.Context) {
	session := sessions.Default(c)
	if v := session.Get(_authKey); v != nil {
		auth := v.(auth)
		log.Println("handleSignOut auth:", auth)
		if auth.status == authActive {
			session.Clear()
			session.Save()
			c.JSON(http.StatusOK, s.getResponse(errNone, msgSuccess))
			return
		}
	} else {
		log.Println("no session")
	}
	c.JSON(http.StatusOK, s.getResponse(errUnknown, msgFailed))
}
