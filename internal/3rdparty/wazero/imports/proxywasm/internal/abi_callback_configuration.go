// Copyright 2020-2021 Tetrate
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package internal

import (
	"time"

	"wa-lang.org/wa/internal/3rdparty/wazero/imports/proxywasm/types"
)

//export proxy_on_vm_start
func proxyOnVMStart(_ uint32, vmConfigurationSize int) types.OnVMStartStatus {
	if recordTiming {
		defer logTiming("proxyOnVMStart", time.Now())
	}
	return currentState.vmContext.OnVMStart(vmConfigurationSize)
}

//export proxy_on_configure
func proxyOnConfigure(pluginContextID uint32, pluginConfigurationSize int) types.OnPluginStartStatus {
	if recordTiming {
		defer logTiming("proxyOnConfigure", time.Now())
	}
	ctx, ok := currentState.pluginContexts[pluginContextID]
	if !ok {
		panic("invalid context on proxy_on_configure")
	}
	currentState.setActiveContextID(pluginContextID)
	return ctx.context.OnPluginStart(pluginConfigurationSize)
}
