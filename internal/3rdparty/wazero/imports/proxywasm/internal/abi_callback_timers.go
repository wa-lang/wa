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

import "time"

//export proxy_on_tick
func proxyOnTick(pluginContextID uint32) {
	if recordTiming {
		defer logTiming("proxyOnTick", time.Now())
	}
	ctx, ok := currentState.pluginContexts[pluginContextID]
	if !ok {
		panic("invalid root_context_id")
	}
	currentState.setActiveContextID(pluginContextID)
	ctx.context.OnTick()
}
