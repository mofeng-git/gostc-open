// Copyright 2016 fatedier, fatedier@gmail.com
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

package frpc

import (
	"context"
	"gostc-sub/pkg/p2p/client"
	v1 "gostc-sub/pkg/p2p/pkg/config/v1"
	"time"
)

type Service struct {
	common      v1.ClientCommonConfig
	proxyCfgs   []v1.ProxyConfigurer
	visitorCfgs []v1.VisitorConfigurer
	svc         *client.Service
	stopChan    chan struct{}
}

func NewService(common v1.ClientCommonConfig, proxyCfgs []v1.ProxyConfigurer, visitorCfgs []v1.VisitorConfigurer) *Service {
	return &Service{
		common:      common,
		proxyCfgs:   proxyCfgs,
		visitorCfgs: visitorCfgs,
		stopChan:    make(chan struct{}),
	}
}

func (s *Service) Start() (err error) {
	s.common.Complete()
	for i := 0; i < len(s.proxyCfgs); i++ {
		s.proxyCfgs[i].Complete("")
	}
	for i := 0; i < len(s.visitorCfgs); i++ {
		s.visitorCfgs[i].Complete(&s.common)
	}
	if s.svc, err = client.NewService(client.ServiceOptions{
		Common:      &s.common,
		ProxyCfgs:   s.proxyCfgs,
		VisitorCfgs: s.visitorCfgs,
	}); err != nil {
		return err
	}
	go func() {
		err = s.svc.Run(context.Background())
	}()
	time.Sleep(time.Second)
	return err
}

func (s *Service) Stop() {
	s.svc.Close()
	close(s.stopChan)
}

func (s *Service) Wait() {
	select {
	case <-s.stopChan:
	}
}
