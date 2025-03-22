// Copyright 2018 fatedier, fatedier@gmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package frps

import (
	"context"
	v1 "gostc-sub/p2p/pkg/config/v1"
	"gostc-sub/p2p/server"
)

type Service struct {
	cfg      v1.ServerConfig
	svc      *server.Service
	stopChan chan struct{}
}

func NewService(cfg v1.ServerConfig) *Service {
	return &Service{
		cfg:      cfg,
		stopChan: make(chan struct{}),
	}
}

func (s *Service) Start() (err error) {
	s.cfg.Complete()
	s.svc, err = server.NewService(&s.cfg)
	if err != nil {
		return err
	}
	go s.svc.Run(context.Background())
	return nil
}

func (s *Service) Stop() {
	_ = s.svc.Close()
	close(s.stopChan)
}

func (s *Service) Wait() {
	select {
	case <-s.stopChan:
	}
}
