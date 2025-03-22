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

type ServerConfig struct {
	Auth                            AuthServerConfig      `json:"auth,omitempty"`
	BindAddr                        string                `json:"bindAddr,omitempty"`
	BindPort                        int                   `json:"bindPort,omitempty"`
	Transport                       ServerTransportConfig `json:"transport,omitempty"`
	DetailedErrorsToClient          *bool                 `json:"detailedErrorsToClient,omitempty"`
	MaxPortsPerClient               int64                 `json:"maxPortsPerClient,omitempty"`
	UserConnTimeout                 int64                 `json:"userConnTimeout,omitempty"`
	NatHoleAnalysisDataReserveHours int64                 `json:"natholeAnalysisDataReserveHours,omitempty"`
	HTTPPlugins                     []HTTPPluginOptions   `json:"httpPlugins,omitempty"`
}

type AuthServerConfig struct {
	Method           AuthMethod           `json:"method,omitempty"`
	AdditionalScopes []AuthScope          `json:"additionalScopes,omitempty"`
	Token            string               `json:"token,omitempty"`
	OIDC             AuthOIDCServerConfig `json:"oidc,omitempty"`
}

type AuthOIDCServerConfig struct {
	Issuer          string `json:"issuer,omitempty"`
	Audience        string `json:"audience,omitempty"`
	SkipExpiryCheck bool   `json:"skipExpiryCheck,omitempty"`
	SkipIssuerCheck bool   `json:"skipIssuerCheck,omitempty"`
}

type ServerTransportConfig struct {
	TCPMux                  *bool           `json:"tcpMux,omitempty"`
	TCPMuxKeepaliveInterval int64           `json:"tcpMuxKeepaliveInterval,omitempty"`
	TCPKeepAlive            int64           `json:"tcpKeepalive,omitempty"`
	MaxPoolCount            int64           `json:"maxPoolCount,omitempty"`
	HeartbeatTimeout        int64           `json:"heartbeatTimeout,omitempty"`
	QUIC                    QUICOptions     `json:"quic,omitempty"`
	TLS                     TLSServerConfig `json:"tls,omitempty"`
}

type TLSServerConfig struct {
	Force bool `json:"force,omitempty"`
	TLSConfig
}
