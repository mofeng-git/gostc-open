// Copyright 2019 fatedier, fatedier@gmail.com
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

package controller

import (
	"gostc-sub/p2p/pkg/nathole"
	plugin "gostc-sub/p2p/pkg/plugin/server"
	"gostc-sub/p2p/server/visitor"
)

// All resource managers and controllers
type ResourceController struct {
	// Manage all visitor listeners
	VisitorManager *visitor.Manager

	// Controller for nat hole connections
	NatHoleController *nathole.Controller

	// All server manager plugin
	PluginManager *plugin.Manager
}

func (rc *ResourceController) Close() error {
	return nil
}
