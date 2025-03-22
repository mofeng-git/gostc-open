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

type ProxyTransport struct {
	UseEncryption        bool   `json:"useEncryption,omitempty"`
	UseCompression       bool   `json:"useCompression,omitempty"`
	BandwidthLimit       string `json:"bandwidthLimit,omitempty"`
	BandwidthLimitMode   string `json:"bandwidthLimitMode,omitempty"`
	ProxyProtocolVersion string `json:"proxyProtocolVersion,omitempty"`
}

type ProxyBackend struct {
	LocalIP   string `json:"localIP,omitempty"`
	LocalPort int    `json:"localPort,omitempty"`
}

type ProxyBaseConfig struct {
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Transport   ProxyTransport    `json:"transport,omitempty"`
	Metadatas   map[string]string `json:"metadatas,omitempty"`
	ProxyBackend
}

type STCPProxyConfig struct {
	ProxyBaseConfig

	Secretkey  string   `json:"secretKey,omitempty"`
	AllowUsers []string `json:"allowUsers,omitempty"`
}

type XTCPProxyConfig struct {
	ProxyBaseConfig

	Secretkey  string   `json:"secretKey,omitempty"`
	AllowUsers []string `json:"allowUsers,omitempty"`
}
