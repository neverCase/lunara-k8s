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

type auth struct {
	status    int
	loginTime time.Time
}

func (s *Service) InitHttpServer() *httpService {
	router := gin.New()
	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{Output: s.output}), gin.RecoveryWithWriter(s.output))
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("sessionStore", store))
	router.GET("/login", func(c *gin.Context) {
		session := sessions.Default(c)
		if v := session.Get("auth"); v != nil {
			auth := v.(auth)
			if auth.status == authActive {
				c.JSON(http.StatusOK, s.getResponse(errSignedIn, msgSuccess))
				return
			}
		}
		// verify authentication
		login := &loginRequest{
			Account:  c.Query("account"),
			Password: c.Query("password"),
		}
		if check := s.checkAuth(login); check == false {
			c.JSON(http.StatusOK, s.getResponse(errLoginFailed, msgLoginFailed))
			return
		}
		auth := auth{
			status:    authInActive,
			loginTime: time.Now(),
		}
		session.Set("auth", auth)
		session.Save()
		c.JSON(http.StatusOK, s.getResponse(errNone, msgSuccess))
	})
	router.GET("/incr", func(c *gin.Context) {
		session := sessions.Default(c)
		if v := session.Get("auth"); v != nil {
			auth := v.(auth)
			if auth.status == authActive {
				// todo
				c.JSON(200, s.getResponse(errNone, msgSuccess))
				return
			}
		}
		c.JSON(http.StatusOK, s.getResponse(errNoAuth, msgNoAuth))
		return
	})
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.c.Http.IP, s.c.Http.Port),
		Handler: router,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				log.Println("The server closed under request")
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
