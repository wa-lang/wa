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
	"wa-lang.org/wa/internal/3rdparty/wazero/imports/proxywasm/types"
)

type (
	pluginContextState struct {
		context       types.PluginContext
		httpCallbacks map[uint32]*httpCallbackAttribute
	}

	httpCallbackAttribute struct {
		callback        func(numHeaders, bodySize, numTrailers int)
		callerContextID uint32
	}
)

type state struct {
	vmContext      types.VMContext
	pluginContexts map[uint32]*pluginContextState
	httpContexts   map[uint32]types.HttpContext
	tcpContexts    map[uint32]types.TcpContext

	contextIDToRootID map[uint32]uint32
	activeContextID   uint32
}

var currentState = &state{
	pluginContexts:    make(map[uint32]*pluginContextState),
	httpContexts:      make(map[uint32]types.HttpContext),
	tcpContexts:       make(map[uint32]types.TcpContext),
	contextIDToRootID: make(map[uint32]uint32),
}

func SetVMContext(vmContext types.VMContext) {
	currentState.vmContext = vmContext
}

func RegisterHttpCallout(calloutID uint32, callback func(numHeaders, bodySize, numTrailers int)) {
	currentState.registerHttpCallOut(calloutID, callback)
}

func (s *state) createPluginContext(contextID uint32) {
	ctx := s.vmContext.NewPluginContext(contextID)
	s.pluginContexts[contextID] = &pluginContextState{
		context:       ctx,
		httpCallbacks: map[uint32]*httpCallbackAttribute{},
	}

	// NOTE: this is a temporary work around for avoiding nil pointer panic
	// when users make http dispatch(es) on PluginContext.
	// See https://github.com/tetratelabs/proxy-wasm-go-sdk/issues/110
	// TODO: refactor
	s.contextIDToRootID[contextID] = contextID
}

func (s *state) createTcpContext(contextID uint32, pluginContextID uint32) bool {
	root, ok := s.pluginContexts[pluginContextID]
	if !ok {
		panic("invalid plugin context id")
	}

	if _, ok := s.tcpContexts[contextID]; ok {
		panic("context id duplicated")
	}

	ctx := root.context.NewTcpContext(contextID)
	if ctx == nil {
		// NewTcpContext is not defined by the user
		return false
	}
	s.contextIDToRootID[contextID] = pluginContextID
	s.tcpContexts[contextID] = ctx
	return true
}

func (s *state) createHttpContext(contextID uint32, pluginContextID uint32) bool {
	root, ok := s.pluginContexts[pluginContextID]
	if !ok {
		panic("invalid plugin context id")
	}

	if _, ok := s.httpContexts[contextID]; ok {
		panic("context id duplicated")
	}

	ctx := root.context.NewHttpContext(contextID)
	if ctx == nil {
		// NewHttpContext is not defined by the user
		return false
	}
	s.contextIDToRootID[contextID] = pluginContextID
	s.httpContexts[contextID] = ctx
	return true
}

func (s *state) registerHttpCallOut(calloutID uint32, callback func(numHeaders, bodySize, numTrailers int)) {
	r := s.pluginContexts[s.contextIDToRootID[s.activeContextID]]
	r.httpCallbacks[calloutID] = &httpCallbackAttribute{callback: callback, callerContextID: s.activeContextID}
}

func (s *state) setActiveContextID(contextID uint32) {
	s.activeContextID = contextID
}
