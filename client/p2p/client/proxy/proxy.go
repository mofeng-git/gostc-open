// Copyright 2017 fatedier, fatedier@gmail.com
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

package proxy

import (
	"context"
	"io"
	"net"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	libio "github.com/fatedier/golib/io"
	libnet "github.com/fatedier/golib/net"
	pp "github.com/pires/go-proxyproto"
	"golang.org/x/time/rate"

	"gostc-sub/p2p/pkg/config/types"
	v1 "gostc-sub/p2p/pkg/config/v1"
	"gostc-sub/p2p/pkg/msg"
	"gostc-sub/p2p/pkg/transport"
	"gostc-sub/p2p/pkg/util/limit"
	"gostc-sub/p2p/pkg/util/xlog"
)

var proxyFactoryRegistry = map[reflect.Type]func(*BaseProxy, v1.ProxyConfigurer) Proxy{}

func RegisterProxyFactory(proxyConfType reflect.Type, factory func(*BaseProxy, v1.ProxyConfigurer) Proxy) {
	proxyFactoryRegistry[proxyConfType] = factory
}

// Proxy defines how to handle work connections for different proxy type.
type Proxy interface {
	Run() error
	// InWorkConn accept work connections registered to server.
	InWorkConn(net.Conn, *msg.StartWorkConn)
	SetInWorkConnCallback(func(*v1.ProxyBaseConfig, net.Conn, *msg.StartWorkConn) /* continue */ bool)
	Close()
}

func NewProxy(
	ctx context.Context,
	pxyConf v1.ProxyConfigurer,
	clientCfg *v1.ClientCommonConfig,
	msgTransporter transport.MessageTransporter,
) (pxy Proxy) {
	var limiter *rate.Limiter
	limitBytes := pxyConf.GetBaseConfig().Transport.BandwidthLimit.Bytes()
	if limitBytes > 0 && pxyConf.GetBaseConfig().Transport.BandwidthLimitMode == types.BandwidthLimitModeClient {
		limiter = rate.NewLimiter(rate.Limit(float64(limitBytes)), int(limitBytes))
	}

	baseProxy := BaseProxy{
		baseCfg:        pxyConf.GetBaseConfig(),
		clientCfg:      clientCfg,
		limiter:        limiter,
		msgTransporter: msgTransporter,
		xl:             xlog.FromContextSafe(ctx),
		ctx:            ctx,
	}

	factory := proxyFactoryRegistry[reflect.TypeOf(pxyConf)]
	if factory == nil {
		return nil
	}
	return factory(&baseProxy, pxyConf)
}

type BaseProxy struct {
	baseCfg            *v1.ProxyBaseConfig
	clientCfg          *v1.ClientCommonConfig
	msgTransporter     transport.MessageTransporter
	limiter            *rate.Limiter
	inWorkConnCallback func(*v1.ProxyBaseConfig, net.Conn, *msg.StartWorkConn) /* continue */ bool

	mu  sync.RWMutex
	xl  *xlog.Logger
	ctx context.Context
}

func (pxy *BaseProxy) Run() error {
	return nil
}

func (pxy *BaseProxy) Close() {
}

func (pxy *BaseProxy) SetInWorkConnCallback(cb func(*v1.ProxyBaseConfig, net.Conn, *msg.StartWorkConn) bool) {
	pxy.inWorkConnCallback = cb
}

func (pxy *BaseProxy) InWorkConn(conn net.Conn, m *msg.StartWorkConn) {
	if pxy.inWorkConnCallback != nil {
		if !pxy.inWorkConnCallback(pxy.baseCfg, conn, m) {
			return
		}
	}
	pxy.HandleTCPWorkConnection(conn, m, []byte(pxy.clientCfg.Auth.Token))
}

// Common handler for tcp work connections.
func (pxy *BaseProxy) HandleTCPWorkConnection(workConn net.Conn, m *msg.StartWorkConn, encKey []byte) {
	xl := pxy.xl
	baseCfg := pxy.baseCfg
	var (
		remote io.ReadWriteCloser
		err    error
	)
	remote = workConn
	if pxy.limiter != nil {
		remote = libio.WrapReadWriteCloser(limit.NewReader(workConn, pxy.limiter), limit.NewWriter(workConn, pxy.limiter), func() error {
			return workConn.Close()
		})
	}

	xl.Tracef("handle tcp work connection, useEncryption: %t, useCompression: %t",
		baseCfg.Transport.UseEncryption, baseCfg.Transport.UseCompression)
	if baseCfg.Transport.UseEncryption {
		remote, err = libio.WithEncryption(remote, encKey)
		if err != nil {
			workConn.Close()
			xl.Errorf("create encryption stream error: %v", err)
			return
		}
	}
	var compressionResourceRecycleFn func()
	if baseCfg.Transport.UseCompression {
		remote, compressionResourceRecycleFn = libio.WithCompressionFromPool(remote)
	}

	// check if we need to send proxy protocol info
	if m.SrcAddr != "" && m.SrcPort != 0 {
		if m.DstAddr == "" {
			m.DstAddr = "127.0.0.1"
		}
	}

	if baseCfg.Transport.ProxyProtocolVersion != "" && m.SrcAddr != "" && m.SrcPort != 0 {
		h := &pp.Header{
			Command: pp.PROXY,
		}

		if strings.Contains(m.SrcAddr, ".") {
			h.TransportProtocol = pp.TCPv4
		} else {
			h.TransportProtocol = pp.TCPv6
		}

		if baseCfg.Transport.ProxyProtocolVersion == "v1" {
			h.Version = 1
		} else if baseCfg.Transport.ProxyProtocolVersion == "v2" {
			h.Version = 2
		}
	}

	localConn, err := libnet.Dial(
		net.JoinHostPort(baseCfg.LocalIP, strconv.Itoa(baseCfg.LocalPort)),
		libnet.WithTimeout(10*time.Second),
	)
	if err != nil {
		workConn.Close()
		xl.Errorf("connect to local service [%s:%d] error: %v", baseCfg.LocalIP, baseCfg.LocalPort, err)
		return
	}

	xl.Debugf("join connections, localConn(l[%s] r[%s]) workConn(l[%s] r[%s])", localConn.LocalAddr().String(),
		localConn.RemoteAddr().String(), workConn.LocalAddr().String(), workConn.RemoteAddr().String())

	_, _, errs := libio.Join(localConn, remote)
	xl.Debugf("join connections closed")
	if len(errs) > 0 {
		xl.Tracef("join connections errors: %v", errs)
	}
	if compressionResourceRecycleFn != nil {
		compressionResourceRecycleFn()
	}
}
