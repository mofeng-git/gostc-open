// Copyright 2023 The frp Authors
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

package v1

import (
	"github.com/samber/lo"

	"gostc-sub/p2p/pkg/config/types"
	"gostc-sub/p2p/pkg/util/util"
)

type ServerConfig struct {
	Auth AuthServerConfig `json:"auth,omitempty"`
	// BindAddr specifies the address that the server binds to. By default,
	// this value is "0.0.0.0".
	BindAddr string `json:"bindAddr,omitempty"`
	// BindPort specifies the port that the server listens on. By default, this
	// value is 7000.
	BindPort int `json:"bindPort,omitempty"`

	Transport ServerTransportConfig `json:"transport,omitempty"`

	// DetailedErrorsToClient defines whether to send the specific error (with
	// debug info) to frpc. By default, this value is true.
	DetailedErrorsToClient *bool `json:"detailedErrorsToClient,omitempty"`
	// MaxPortsPerClient specifies the maximum number of ports a single client
	// may proxy to. If this value is 0, no limit will be applied.
	MaxPortsPerClient int64 `json:"maxPortsPerClient,omitempty"`
	// UserConnTimeout specifies the maximum time to wait for a work
	// connection. By default, this value is 10.
	UserConnTimeout int64 `json:"userConnTimeout,omitempty"`
	// UDPPacketSize specifies the UDP packet size
	// By default, this value is 1500
	//UDPPacketSize int64 `json:"udpPacketSize,omitempty"`
	// NatHoleAnalysisDataReserveHours specifies the hours to reserve nat hole analysis data.
	NatHoleAnalysisDataReserveHours int64 `json:"natholeAnalysisDataReserveHours,omitempty"`

	AllowPorts []types.PortsRange `json:"allowPorts,omitempty"`

	HTTPPlugins []HTTPPluginOptions `json:"httpPlugins,omitempty"`
}

func (c *ServerConfig) Complete() {
	c.Auth.Complete()
	c.Transport.Complete()

	c.BindAddr = util.EmptyOr(c.BindAddr, "0.0.0.0")
	c.BindPort = util.EmptyOr(c.BindPort, 7000)

	c.DetailedErrorsToClient = util.EmptyOr(c.DetailedErrorsToClient, lo.ToPtr(true))
	c.UserConnTimeout = util.EmptyOr(c.UserConnTimeout, 10)
	c.NatHoleAnalysisDataReserveHours = util.EmptyOr(c.NatHoleAnalysisDataReserveHours, 7*24)
}

type AuthServerConfig struct {
	Method           AuthMethod           `json:"method,omitempty"`
	AdditionalScopes []AuthScope          `json:"additionalScopes,omitempty"`
	Token            string               `json:"token,omitempty"`
	OIDC             AuthOIDCServerConfig `json:"oidc,omitempty"`
}

func (c *AuthServerConfig) Complete() {
	c.Method = util.EmptyOr(c.Method, "token")
}

type AuthOIDCServerConfig struct {
	// Issuer specifies the issuer to verify OIDC tokens with. This issuer
	// will be used to load public keys to verify signature and will be compared
	// with the issuer claim in the OIDC token.
	Issuer string `json:"issuer,omitempty"`
	// Audience specifies the audience OIDC tokens should contain when validated.
	// If this value is empty, audience ("client ID") verification will be skipped.
	Audience string `json:"audience,omitempty"`
	// SkipExpiryCheck specifies whether to skip checking if the OIDC token is
	// expired.
	SkipExpiryCheck bool `json:"skipExpiryCheck,omitempty"`
	// SkipIssuerCheck specifies whether to skip checking if the OIDC token's
	// issuer claim matches the issuer specified in OidcIssuer.
	SkipIssuerCheck bool `json:"skipIssuerCheck,omitempty"`
}

type ServerTransportConfig struct {
	// TCPMux toggles TCP stream multiplexing. This allows multiple requests
	// from a client to share a single TCP connection. By default, this value
	// is true.
	// $HideFromDoc
	TCPMux *bool `json:"tcpMux,omitempty"`
	// TCPMuxKeepaliveInterval specifies the keep alive interval for TCP stream multiplier.
	// If TCPMux is true, heartbeat of application layer is unnecessary because it can only rely on heartbeat in TCPMux.
	TCPMuxKeepaliveInterval int64 `json:"tcpMuxKeepaliveInterval,omitempty"`
	// TCPKeepAlive specifies the interval between keep-alive probes for an active network connection between frpc and frps.
	// If negative, keep-alive probes are disabled.
	TCPKeepAlive int64 `json:"tcpKeepalive,omitempty"`
	// MaxPoolCount specifies the maximum pool size for each proxy. By default,
	// this value is 5.
	MaxPoolCount int64 `json:"maxPoolCount,omitempty"`
	// HeartBeatTimeout specifies the maximum time to wait for a heartbeat
	// before terminating the connection. It is not recommended to change this
	// value. By default, this value is 90. Set negative value to disable it.
	HeartbeatTimeout int64 `json:"heartbeatTimeout,omitempty"`
	// QUIC options.
	QUIC QUICOptions `json:"quic,omitempty"`
	// TLS specifies TLS settings for the connection from the client.
	TLS TLSServerConfig `json:"tls,omitempty"`
}

func (c *ServerTransportConfig) Complete() {
	c.TCPMux = util.EmptyOr(c.TCPMux, lo.ToPtr(true))
	c.TCPMuxKeepaliveInterval = util.EmptyOr(c.TCPMuxKeepaliveInterval, 30)
	c.TCPKeepAlive = util.EmptyOr(c.TCPKeepAlive, 7200)
	c.MaxPoolCount = util.EmptyOr(c.MaxPoolCount, 5)
	if lo.FromPtr(c.TCPMux) {
		// If TCPMux is enabled, heartbeat of application layer is unnecessary because we can rely on heartbeat in tcpmux.
		c.HeartbeatTimeout = util.EmptyOr(c.HeartbeatTimeout, -1)
	} else {
		c.HeartbeatTimeout = util.EmptyOr(c.HeartbeatTimeout, 90)
	}
	c.QUIC.Complete()
	if c.TLS.TrustedCaFile != "" {
		c.TLS.Force = true
	}
}

type TLSServerConfig struct {
	// Force specifies whether to only accept TLS-encrypted connections.
	Force bool `json:"force,omitempty"`

	TLSConfig
}
