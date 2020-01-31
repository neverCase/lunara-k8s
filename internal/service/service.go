package service

import (
	"context"
	"github.com/nevercase/lunara-k8s/configs"
)

type Service struct {
	c      *configs.Config
	ctx    context.Context
		cancel context.CancelFunc
}

func NewService(c *configs.Config) *Service {
	ctx, cancel := context.WithCancel(context.Background())
	s := &Service{
		c:      c,
		ctx:    ctx,
		cancel: cancel,
	}
	return s
}

func (s *Service) Close() {
	s.cancel()
}
