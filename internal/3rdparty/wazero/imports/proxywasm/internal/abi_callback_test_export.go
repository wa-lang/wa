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

//go:build !tinygo

package internal

import "wa-lang.org/wa/internal/3rdparty/wazero/imports/proxywasm/types"

// this file exists only for proxytest package, therefore, these functions are not included in a resulting Wasm binary.

func ProxyOnVMStart(pluginContextID uint32, vmConfigurationSize int) types.OnVMStartStatus {
	return proxyOnVMStart(pluginContextID, vmConfigurationSize)
}

func ProxyOnConfigure(pluginContextID uint32, vmConfigurationSize int) types.OnPluginStartStatus {
	return proxyOnConfigure(pluginContextID, vmConfigurationSize)
}

func ProxyOnNewConnection(contextID uint32) types.Action {
	return proxyOnNewConnection(contextID)
}

func ProxyOnDownstreamData(contextID uint32, dataSize int, endOfStream bool) types.Action {
	return proxyOnDownstreamData(contextID, dataSize, endOfStream)
}

func ProxyOnDownstreamConnectionClose(contextID uint32, pType types.PeerType) {
	proxyOnDownstreamConnectionClose(contextID, pType)
}

func ProxyOnUpstreamData(contextID uint32, dataSize int, endOfStream bool) types.Action {
	return proxyOnUpstreamData(contextID, dataSize, endOfStream)
}

func ProxyOnUpstreamConnectionClose(contextID uint32, pType types.PeerType) {
	proxyOnUpstreamConnectionClose(contextID, pType)
}

func ProxyOnRequestHeaders(contextID uint32, numHeaders int, endOfStream bool) types.Action {
	return proxyOnRequestHeaders(contextID, numHeaders, endOfStream)
}

func ProxyOnRequestBody(contextID uint32, bodySize int, endOfStream bool) types.Action {
	return proxyOnRequestBody(contextID, bodySize, endOfStream)
}

func ProxyOnRequestTrailers(contextID uint32, numTrailers int) types.Action {
	return proxyOnRequestTrailers(contextID, numTrailers)
}

func ProxyOnResponseHeaders(contextID uint32, numHeaders int, endOfStream bool) types.Action {
	return proxyOnResponseHeaders(contextID, numHeaders, endOfStream)
}

func ProxyOnResponseBody(contextID uint32, bodySize int, endOfStream bool) types.Action {
	return proxyOnResponseBody(contextID, bodySize, endOfStream)
}

func ProxyOnResponseTrailers(contextID uint32, numTrailers int) types.Action {
	return proxyOnResponseTrailers(contextID, numTrailers)
}

func ProxyOnHttpCallResponse(pluginContextID, calloutID uint32, numHeaders, bodySize, numTrailers int) {
	proxyOnHttpCallResponse(pluginContextID, calloutID, numHeaders, bodySize, numTrailers)
}

func ProxyOnContextCreate(contextID uint32, pluginContextID uint32) {
	proxyOnContextCreate(contextID, pluginContextID)
}

func ProxyOnDone(contextID uint32) bool {
	return proxyOnDone(contextID)
}

func ProxyOnLog(pluginContextID uint32) {
	proxyOnLog(pluginContextID)
}

func ProxyOnQueueReady(contextID, queueID uint32) {
	proxyOnQueueReady(contextID, queueID)
}

func ProxyOnTick(pluginContextID uint32) {
	proxyOnTick(pluginContextID)
}

func ProxyOnDelete(contextID uint32) {
	proxyOnDelete(contextID)
}
