package service

import (
	"context"
	"log"
	"os"

	"github.com/nevercase/lunara-k8s/configs"
	"github.com/nevercase/lunara-k8s/internal/kubernetes"
)

type Service struct {
	c             *configs.Config
	output        *os.File
	httpService   *httpService
	k8sController *kubernetes.K8SController
	ctx           context.Context
	cancel        context.CancelFunc
}

func NewService(c *configs.Config) *Service {
	file, err := os.OpenFile(c.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Panic("OpenFile err: ", err)
	}
	log.SetOutput(file)
	ctx, cancel := context.WithCancel(context.Background())
	s := &Service{
		c:             c,
		output:        file,
		k8sController: kubernetes.NewK8SController(c.Kubernetes),
		ctx:           ctx,
		cancel:        cancel,
	}
	s.httpService = s.InitHttpServer()
	return s
}

func (s *Service) Close() {
	s.httpService.ShutDown()
	s.cancel()
	if err := s.output.Close(); err != nil {
		log.Println("log file close")
	}
}
