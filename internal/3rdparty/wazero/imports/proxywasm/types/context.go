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

package types

// There are four types of these interfaces which you are supposed to implement in order to extend your network proxies.
// They are VMContext, PluginContext, TcpContext and HttpContext, and their relationship can be described as the following diagram:
//
//	                        Wasm Virtual Machine(VM)
//	                   (corresponds to VM configuration)
//	┌────────────────────────────────────────────────────────────────────────────┐
//	│                                                      TcpContext            │
//	│                                                  ╱ (Each Tcp stream)       │
//	│                                                 ╱                          │
//	│                      1: N                      ╱ 1: N                      │
//	│       VMContext  ──────────  PluginContext                                 │
//	│  (VM configuration)     (Plugin configuration) ╲ 1: N                      │
//	│                                                 ╲                          │
//	│                                                  ╲   HttpContext           │
//	│                                                   (Each Http stream)       │
//	└────────────────────────────────────────────────────────────────────────────┘
//
// To summarize,
//
// 1) VMContext corresponds to each Wasm Virtual Machine, and only one VMContext exists in each VM.
// Note that in Envoy, Wasm VMs are created per "vm_config" field in envoy.yaml. For example having different "vm_config.configuration" fields
// results in multiple VMs being created and each of them corresponds to each "vm_config.configuration".
//
// 2) VMContext is parent of PluginContext, and is responsible for creating arbitrary number of PluginContexts.
//
// 3) PluginContext corresponds to each plugin configuration in the host. In Envoy, each plugin configuration is given at HttpFilter or NetworkFilter
// on listeners. That said, a plugin context corresponds to an Http or Network filter on a listener and is in charge of creating the "filter instances" for
// each Http or Tcp stream. These "filter instances" are represented by HttpContexts or TcpContexts.
//
// 4) PluginContext is a parent of TcpContext and HttpContext, and is responsible for creating arbitrary number of these contexts.
//
// 5) TcpContext is responsible for handling individual Tcp stream events.
//
// 6) HttpContext is responsible for handling individual Http stream events.

// VMContext corresponds to a Wasm VM machine and its configuration.
// It's the entrypoint for extending the network proxy.
// Its lifetime matches the Wasm Virtual Machines on the host.
type VMContext interface {
	// OnVMStart is called after the VM is created and main function is called.
	// During this call, GetVMConfiguration hostcall is available and can be used to
	// retrieve the configuration set at vm_config.configuration in the host configuration.
	// This is mainly used for doing Wasm VM-wide initialization.
	OnVMStart(vmConfigurationSize int) OnVMStartStatus

	// NewPluginContext is used for creating PluginContext for each plugin configuration.
	NewPluginContext(contextID uint32) PluginContext
}

// PluginContext corresponds to different plugin configurations (config.configuration).
// Each configuration is typically given at the HTTP/TCP filter in a listener in the hosts.
// PluginContext is responsible for creating the "filter instances" for each TCP/HTTP stream on the listener.
type PluginContext interface {
	// OnPluginStart is called for all plugin contexts (after OnVmStart if this is the VM context).
	// During this call, GetPluginConfiguration is available and can be used to
	// retrieve the configuration set at config.configuration in the host configuration.
	OnPluginStart(pluginConfigurationSize int) OnPluginStartStatus

	// OnPluginDone is called right before the plugin contexts are deleted by hosts.
	// Return false to indicate plugin is in a pending state and there's more work left.
	// In that case you must call PluginDone() function once the work is completed to indicate that
	// hosts can kill this context.
	OnPluginDone() bool

	// OnQueueReady is called when the queue is ready after calling the RegisterQueue hostcall.
	// Note that the queue might be dequeued by another VM running in another thread, so it's
	// possible queue will be empty during the OnQueueReady even if it is not dequeued by this VM.
	OnQueueReady(queueID uint32)

	// OnTick is called when SetTickPeriodMilliSeconds hostcall is called by this plugin context.
	// This can be used to do asynchronous tasks in parallel to the stream processing.
	OnTick()

	// The following functions are used for creating contexts on streams,
	// and developers *must* implement either of them corresponding to
	// extension points. For example, if you configure this plugin context is running
	// at Http filters, then NewHttpContext must be implemented. Same goes for
	// Tcp filters.

	// NewTcpContext is used for creating TcpContext for each Tcp stream.
	// Return nil to indicate this PluginContext is not for TcpContext.
	NewTcpContext(contextID uint32) TcpContext

	// NewHttpContext is used for creating HttpContext for each Http stream.
	// Return nil to indicate this PluginContext is not for HttpContext.
	NewHttpContext(contextID uint32) HttpContext
}

// TcpContext corresponds to each Tcp stream and is created by PluginContext via NewTcpContext.
type TcpContext interface {
	// OnNewConnection is called when the Tcp connection is established between downstream and upstream.
	OnNewConnection() Action

	// OnDownstreamData is called when a data frame arrives from the downstream connection.
	OnDownstreamData(dataSize int, endOfStream bool) Action

	// OnDownstreamClose is called when the downstream connection is closed.
	OnDownstreamClose(peerType PeerType)

	// OnUpstreamData is called when a data frame arrives from the upstream connection.
	OnUpstreamData(dataSize int, endOfStream bool) Action

	// OnUpstreamClose is called when the upstream connection is closed.
	OnUpstreamClose(peerType PeerType)

	// OnStreamDone is called before the host deletes this context.
	// You can retrieve the stream information (such as remote addresses, etc.) during this call.
	// This can be used to implement logging features.
	OnStreamDone()
}

// HttpContext corresponds to each Http stream and is created by PluginContext via NewHttpContext.
type HttpContext interface {
	// OnHttpRequestHeaders is called when request headers arrive.
	// Return types.ActionPause if you want to stop sending headers to the upstream.
	OnHttpRequestHeaders(numHeaders int, endOfStream bool) Action

	// OnHttpRequestBody is called when a request body *frame* arrives.
	// Note that this is potentially called multiple times until we see end_of_stream = true.
	// Return types.ActionPause if you want to buffer the body and stop sending body to the upstream.
	// Even after returning types.ActionPause, this will be called when an unseen frame arrives.
	OnHttpRequestBody(bodySize int, endOfStream bool) Action

	// OnHttpRequestTrailers is called when request trailers arrive.
	// Return types.ActionPause if you want to stop sending trailers to the upstream.
	OnHttpRequestTrailers(numTrailers int) Action

	// OnHttpResponseHeaders is called when response headers arrive.
	// Return types.ActionPause if you want to stop sending headers to downstream.
	OnHttpResponseHeaders(numHeaders int, endOfStream bool) Action

	// OnHttpResponseBody is called when a response body *frame* arrives.
	// Note that this is potentially called multiple times until we see end_of_stream = true.
	// Return types.ActionPause if you want to buffer the body and stop sending body to the downtream.
	// Even after returning types.ActionPause, this will be called when an unseen frame arrives.
	OnHttpResponseBody(bodySize int, endOfStream bool) Action

	// OnHttpResponseTrailers is called when response trailers arrive.
	// Return types.ActionPause if you want to stop sending trailers to the downstream.
	OnHttpResponseTrailers(numTrailers int) Action

	// OnHttpStreamDone is called before the host deletes this context.
	// You can retrieve the HTTP request/response information (such as headers, etc.) during this call.
	// This can be used to implement logging features.
	OnHttpStreamDone()
}

// DefaultContexts are a no-op implementation of contexts.
// Users can embed them into their custom contexts, so that
// they only have to implement methods they want.
type (
	// DefaultVMContext provides the no-op implementation of the VMContext interface.
	DefaultVMContext struct{}

	// DefaultPluginContext provides the no-op implementation of the PluginContext interface.
	DefaultPluginContext struct{}

	// DefaultTcpContext provides the no-op implementation of the TcpContext interface.
	DefaultTcpContext struct{}

	// DefaultHttpContext provides the no-op implementation of the HttpContext interface.
	DefaultHttpContext struct{}
)

// impl VMContext

func (*DefaultVMContext) OnVMStart(vmConfigurationSize int) OnVMStartStatus { return OnVMStartStatusOK }
func (*DefaultVMContext) NewPluginContext(contextID uint32) PluginContext {
	return &DefaultPluginContext{}
}

// impl PluginContext

func (*DefaultPluginContext) OnQueueReady(uint32) {}
func (*DefaultPluginContext) OnTick()             {}
func (*DefaultPluginContext) OnPluginStart(int) OnPluginStartStatus {
	return OnPluginStartStatusOK
}
func (*DefaultPluginContext) OnPluginDone() bool                { return true }
func (*DefaultPluginContext) NewTcpContext(uint32) TcpContext   { return nil }
func (*DefaultPluginContext) NewHttpContext(uint32) HttpContext { return nil }

// impl TcpContext

func (*DefaultTcpContext) OnDownstreamData(int, bool) Action { return ActionContinue }
func (*DefaultTcpContext) OnDownstreamClose(PeerType)        {}
func (*DefaultTcpContext) OnNewConnection() Action           { return ActionContinue }
func (*DefaultTcpContext) OnUpstreamData(int, bool) Action   { return ActionContinue }
func (*DefaultTcpContext) OnUpstreamClose(PeerType)          {}
func (*DefaultTcpContext) OnStreamDone()                     {}

// impl HttpContext

func (*DefaultHttpContext) OnHttpRequestHeaders(int, bool) Action  { return ActionContinue }
func (*DefaultHttpContext) OnHttpRequestBody(int, bool) Action     { return ActionContinue }
func (*DefaultHttpContext) OnHttpRequestTrailers(int) Action       { return ActionContinue }
func (*DefaultHttpContext) OnHttpResponseHeaders(int, bool) Action { return ActionContinue }
func (*DefaultHttpContext) OnHttpResponseBody(int, bool) Action    { return ActionContinue }
func (*DefaultHttpContext) OnHttpResponseTrailers(int) Action      { return ActionContinue }
func (*DefaultHttpContext) OnHttpStreamDone()                      {}

var (
	_ VMContext     = &DefaultVMContext{}
	_ PluginContext = &DefaultPluginContext{}
	_ TcpContext    = &DefaultTcpContext{}
	_ HttpContext   = &DefaultHttpContext{}
)
