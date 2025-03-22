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

type ClientCommonConfig struct {
	Auth              AuthClientConfig      `json:"auth,omitempty"`
	User              string                `json:"user,omitempty"`
	ServerAddr        string                `json:"serverAddr,omitempty"`
	ServerPort        int                   `json:"serverPort,omitempty"`
	NatHoleSTUNServer string                `json:"natHoleStunServer,omitempty"`
	DNSServer         string                `json:"dnsServer,omitempty"`
	LoginFailExit     *bool                 `json:"loginFailExit,omitempty"`
	Transport         ClientTransportConfig `json:"transport,omitempty"`
	Metadatas         map[string]string     `json:"metadatas,omitempty"`
}

type ClientTransportConfig struct {
	Protocol                string          `json:"protocol,omitempty"`
	DialServerTimeout       int64           `json:"dialServerTimeout,omitempty"`
	DialServerKeepAlive     int64           `json:"dialServerKeepalive,omitempty"`
	ConnectServerLocalIP    string          `json:"connectServerLocalIP,omitempty"`
	ProxyURL                string          `json:"proxyURL,omitempty"`
	PoolCount               int             `json:"poolCount,omitempty"`
	TCPMux                  *bool           `json:"tcpMux,omitempty"`
	TCPMuxKeepaliveInterval int64           `json:"tcpMuxKeepaliveInterval,omitempty"`
	QUIC                    QUICOptions     `json:"quic,omitempty"`
	HeartbeatInterval       int64           `json:"heartbeatInterval,omitempty"`
	HeartbeatTimeout        int64           `json:"heartbeatTimeout,omitempty"`
	TLS                     TLSClientConfig `json:"tls,omitempty"`
}

type TLSClientConfig struct {
	Enable                    *bool `json:"enable,omitempty"`
	DisableCustomTLSFirstByte *bool `json:"disableCustomTLSFirstByte,omitempty"`
	TLSConfig
}

type AuthClientConfig struct {
	Method           AuthMethod           `json:"method,omitempty"`
	AdditionalScopes []AuthScope          `json:"additionalScopes,omitempty"`
	Token            string               `json:"token,omitempty"`
	OIDC             AuthOIDCClientConfig `json:"oidc,omitempty"`
}

type AuthOIDCClientConfig struct {
	ClientID                 string            `json:"clientID,omitempty"`
	ClientSecret             string            `json:"clientSecret,omitempty"`
	Audience                 string            `json:"audience,omitempty"`
	Scope                    string            `json:"scope,omitempty"`
	TokenEndpointURL         string            `json:"tokenEndpointURL,omitempty"`
	AdditionalEndpointParams map[string]string `json:"additionalEndpointParams,omitempty"`
}
