package main

import (
	"flag"
	"fmt"
	"github.com/nevercase/lunara-k8s/configs"
	"github.com/nevercase/lunara-k8s/internal/service"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
	if err := configs.Init(); err != nil {
		log.Fatal(err)
	}
	signalHandler(service.NewService(configs.GetConfig()))
}

func signalHandler(s *service.Service) {
	var (
		ch = make(chan os.Signal, 1)
	)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		sig := <-ch
		fmt.Printf("get a signal %s, stop the lunara-k8s service\n", sig.String())
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			s.Close()
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
