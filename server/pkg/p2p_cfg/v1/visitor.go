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

type VisitorTransport struct {
	UseEncryption  bool `json:"useEncryption,omitempty"`
	UseCompression bool `json:"useCompression,omitempty"`
}

type VisitorBaseConfig struct {
	Name      string           `json:"name"`
	Type      string           `json:"type"`
	Transport VisitorTransport `json:"transport,omitempty"`
	SecretKey string           `json:"secretKey,omitempty"`
	// if the server user is not set, it defaults to the current user
	ServerUser string `json:"serverUser,omitempty"`
	ServerName string `json:"serverName,omitempty"`
	BindAddr   string `json:"bindAddr,omitempty"`
	// BindPort is the port that visitor listens on.
	// It can be less than 0, it means don't bind to the port and only receive connections redirected from
	// other visitors. (This is not supported for SUDP now)
	BindPort int `json:"bindPort,omitempty"`
}

type VisitorConfigurer interface {
	Complete(*ClientCommonConfig)
	GetBaseConfig() *VisitorBaseConfig
}

type VisitorType string

type STCPVisitorConfig struct {
	VisitorBaseConfig
}

type XTCPVisitorConfig struct {
	VisitorBaseConfig

	Protocol          string `json:"protocol,omitempty"`
	KeepTunnelOpen    bool   `json:"keepTunnelOpen,omitempty"`
	MaxRetriesAnHour  int    `json:"maxRetriesAnHour,omitempty"`
	MinRetryInterval  int    `json:"minRetryInterval,omitempty"`
	FallbackTo        string `json:"fallbackTo,omitempty"`
	FallbackTimeoutMs int    `json:"fallbackTimeoutMs,omitempty"`
}
