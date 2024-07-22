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

package proxywasm

import (
	"errors"
	"fmt"
	"math"

	"wa-lang.org/wa/internal/3rdparty/wazero/imports/proxywasm/internal"
	"wa-lang.org/wa/internal/3rdparty/wazero/imports/proxywasm/types"
)

// GetVMConfiguration is used for retrieving configurations given in the "vm_config.configuration" field.
// This hostcall is only available during types.PluginContext.OnVMStart call.
func GetVMConfiguration() ([]byte, error) {
	return getBuffer(internal.BufferTypeVMConfiguration, 0, math.MaxInt32)
}

// GetPluginConfiguration is used for retrieving configurations given in the "config.configuration" field.
// This hostcall is only available during types.PluginContext.OnPluginStart call.
func GetPluginConfiguration() ([]byte, error) {
	return getBuffer(internal.BufferTypePluginConfiguration, 0, math.MaxInt32)
}

// SetTickPeriodMilliSeconds sets the tick interval of types.PluginContext.OnTick calls.
// Only available for types.PluginContext.
func SetTickPeriodMilliSeconds(millSec uint32) error {
	return internal.StatusToError(internal.ProxySetTickPeriodMilliseconds(millSec))
}

// SetEffectiveContext sets the effective context to "context_id".
// This hostcall is usually used to change the context after receiving
// types.PluginContext.OnQueueReady or types.PluginContext.OnTick
func SetEffectiveContext(contextID uint32) error {
	return internal.StatusToError(internal.ProxySetEffectiveContext(contextID))
}

// RegisterSharedQueue registers the shared queue on this plugin context.
// "Register" means that OnQueueReady is called for this plugin context whenever a new item is enqueued on that queueID.
// Only available for types.PluginContext. The returned queueID can be used for Enqueue/DequeueSharedQueue.
// Note that "name" must be unique across all Wasm VMs which share the same "vm_id".
// That means you can use "vm_id" to separate shared queue namespace.
func RegisterSharedQueue(name string) (queueID uint32, err error) {
	ptr := internal.StringBytePtr(name)
	st := internal.ProxyRegisterSharedQueue(ptr, len(name), &queueID)
	return queueID, internal.StatusToError(st)
}

// ResolveSharedQueue acquires the queueID for the given vmID and queueName.
// The returned queueID can be used for Enqueue/DequeueSharedQueue.
func ResolveSharedQueue(vmID, queueName string) (queueID uint32, err error) {
	var ret uint32
	st := internal.ProxyResolveSharedQueue(internal.StringBytePtr(vmID),
		len(vmID), internal.StringBytePtr(queueName), len(queueName), &ret)
	return ret, internal.StatusToError(st)
}

// EnqueueSharedQueue enqueues data to the shared queue of the given queueID.
// In order to get queue id for a target queue, use "ResolveSharedQueue" first.
func EnqueueSharedQueue(queueID uint32, data []byte) error {
	return internal.StatusToError(internal.ProxyEnqueueSharedQueue(queueID, &data[0], len(data)))
}

// DequeueSharedQueue dequeues data from the shared queue of the given queueID.
// In order to get queue id for a target queue, use "ResolveSharedQueue" first.
func DequeueSharedQueue(queueID uint32) ([]byte, error) {
	var raw *byte
	var size int
	st := internal.ProxyDequeueSharedQueue(queueID, &raw, &size)
	if st != internal.StatusOK {
		return nil, internal.StatusToError(st)
	}
	return internal.RawBytePtrToByteSlice(raw, size), nil
}

// PluginDone must be called when OnPluginDone returns false indicating that the plugin is in pending state
// right before deletion by the hosts. Only available for types.PluginContext.
func PluginDone() {
	internal.ProxyDone()
}

// DispatchHttpCall is for dispatching HTTP calls to a remote cluster. This can be used by all contexts
// including Tcp and Root contexts. "cluster" arg specifies the remote cluster the host will send
// the request against with "headers", "body", and "trailers" arguments. "callBack" function is called if the host successfully
// made the request and received the response from the remote cluster.
// When the callBack function is called, the "GetHttpCallResponseHeaders", "GetHttpCallResponseBody", "GetHttpCallResponseTrailers"
// calls are available for accessing the response information.
func DispatchHttpCall(
	cluster string,
	headers [][2]string,
	body []byte,
	trailers [][2]string,
	timeoutMillisecond uint32,
	callBack func(numHeaders, bodySize, numTrailers int),
) (calloutID uint32, err error) {
	shs := internal.SerializeMap(headers)
	hp := &shs[0]
	hl := len(shs)

	sts := internal.SerializeMap(trailers)
	tp := &sts[0]
	tl := len(sts)

	var bodyPtr *byte
	if len(body) > 0 {
		bodyPtr = &body[0]
	}

	u := internal.StringBytePtr(cluster)
	switch st := internal.ProxyHttpCall(u, len(cluster),
		hp, hl, bodyPtr, len(body), tp, tl, timeoutMillisecond, &calloutID); st {
	case internal.StatusOK:
		internal.RegisterHttpCallout(calloutID, callBack)
		return calloutID, nil
	default:
		return 0, internal.StatusToError(st)
	}
}

// GetHttpCallResponseHeaders is used for retrieving HTTP response headers
// returned by a remote cluster in response to the DispatchHttpCall.
// Only available during "callback" function passed to the DispatchHttpCall.
func GetHttpCallResponseHeaders() ([][2]string, error) {
	return getMap(internal.MapTypeHttpCallResponseHeaders)
}

// GetHttpCallResponseBody is used for retrieving HTTP response body
// returned by a remote cluster in response to the DispatchHttpCall.
// Only available during "callback" function passed to the DispatchHttpCall.
func GetHttpCallResponseBody(start, maxSize int) ([]byte, error) {
	return getBuffer(internal.BufferTypeHttpCallResponseBody, start, maxSize)
}

// GetHttpCallResponseTrailers is used for retrieving HTTP response trailers
// returned by a remote cluster in response to the DispatchHttpCall.
// Only available during "callback" function passed to DispatchHttpCall.
func GetHttpCallResponseTrailers() ([][2]string, error) {
	return getMap(internal.MapTypeHttpCallResponseTrailers)
}

// GetDownstreamData can be used for retrieving TCP downstream data buffered in the host.
// Returned bytes beginning from "start" to "start" + "maxSize" in the buffer.
// Only available during types.TcpContext.OnDownstreamData.
func GetDownstreamData(start, maxSize int) ([]byte, error) {
	return getBuffer(internal.BufferTypeDownstreamData, start, maxSize)
}

// AppendDownstreamData appends the given bytes to the downstream TCP data buffered in the host.
// Only available during types.TcpContext.OnDownstreamData.
func AppendDownstreamData(data []byte) error {
	return appendToBuffer(internal.BufferTypeDownstreamData, data)
}

// PrependDownstreamData prepends the given bytes to the downstream TCP data buffered in the host.
// Only available during types.TcpContext.OnDownstreamData.
func PrependDownstreamData(data []byte) error {
	return prependToBuffer(internal.BufferTypeDownstreamData, data)
}

// ReplaceDownstreamData replaces the downstream TCP data buffered in the host
// with the given bytes. Only available during types.TcpContext.OnDownstreamData.
func ReplaceDownstreamData(data []byte) error {
	return replaceBuffer(internal.BufferTypeDownstreamData, data)
}

// GetUpstreamData can be used for retrieving upstream TCP data buffered in the host.
// Returned bytes beginning from "start" to "start" + "maxSize" in the buffer.
// Only available during types.TcpContext.OnUpstreamData.
func GetUpstreamData(start, maxSize int) ([]byte, error) {
	return getBuffer(internal.BufferTypeUpstreamData, start, maxSize)
}

// AppendUpstreamData appends the given bytes to the upstream TCP data buffered in the host.
// Only available during types.TcpContext.OnUpstreamData.
func AppendUpstreamData(data []byte) error {
	return appendToBuffer(internal.BufferTypeUpstreamData, data)
}

// PrependUpstreamData prepends the given bytes to the upstream TCP data buffered in the host.
// Only available during types.TcpContext.OnUpstreamData.
func PrependUpstreamData(data []byte) error {
	return prependToBuffer(internal.BufferTypeUpstreamData, data)
}

// ReplaceUpstreamData replaces the upstream TCP data buffered in the host
// with the given bytes. Only available during types.TcpContext.OnUpstreamData.
func ReplaceUpstreamData(data []byte) error {
	return replaceBuffer(internal.BufferTypeUpstreamData, data)
}

// ContinueTcpStream continues interating on the TCP connection
// after types.Action.Pause was returned by types.TcpContext.
// Only available for types.TcpContext.
func ContinueTcpStream() error {
	// Note that internal.ProxyContinueStream is not implemented in Envoy,
	// so we intentionally choose to pass StreamTypeDownstream here while
	// the name itself is not indicating "continue downstream".
	return internal.StatusToError(internal.ProxyContinueStream(internal.StreamTypeDownstream))
}

// CloseDownstream closes the downstream TCP connection for this Tcp context.
// Only available for types.TcpContext.
func CloseDownstream() error {
	return internal.StatusToError(internal.ProxyCloseStream(internal.StreamTypeDownstream))
}

// CloseUpstream closes the upstream TCP connection for this Tcp context.
// Only available for types.TcpContext.
func CloseUpstream() error {
	return internal.StatusToError(internal.ProxyCloseStream(internal.StreamTypeUpstream))
}

// GetHttpRequestHeaders is used for retrieving HTTP request headers.
// Only available during types.HttpContext.OnHttpRequestHeaders and
// types.HttpContext.OnHttpStreamDone.
func GetHttpRequestHeaders() ([][2]string, error) {
	return getMap(internal.MapTypeHttpRequestHeaders)
}

// ReplaceHttpRequestHeaders is used for replacing HTTP request headers
// with given headers. Only available during
// types.HttpContext.OnHttpRequestHeaders.
func ReplaceHttpRequestHeaders(headers [][2]string) error {
	return setMap(internal.MapTypeHttpRequestHeaders, headers)
}

// GetHttpRequestHeader is used for retrieving an HTTP request header value
// for given "key". Only available during types.HttpContext.OnHttpRequestHeaders and
// types.HttpContext.OnHttpStreamDone.
// If multiple values are present for the key, the "first" value found in the host is returned.
// See https://github.com/envoyproxy/envoy/blob/72bf41fb0ecc039f196be02f534bfc2c9c69f348/source/extensions/common/wasm/context.cc#L762-L763
// for detail.
func GetHttpRequestHeader(key string) (string, error) {
	return getMapValue(internal.MapTypeHttpRequestHeaders, key)
}

// RemoveHttpRequestHeader is used for removing an HTTP request header with
// a given "key". Only available during types.HttpContext.OnHttpRequestHeaders.
func RemoveHttpRequestHeader(key string) error {
	return removeMapValue(internal.MapTypeHttpRequestHeaders, key)
}

// ReplaceHttpRequestHeader replaces a value for given "key" from request headers.
// Only available during types.HttpContext.OnHttpRequestHeaders.
// If multiple values are present for the key, only the "first" value in the host is replaced.
// See https://github.com/envoyproxy/envoy/blob/72bf41fb0ecc039f196be02f534bfc2c9c69f348/envoy/http/header_map.h#L547-L549
// for detail.
func ReplaceHttpRequestHeader(key, value string) error {
	return replaceMapValue(internal.MapTypeHttpRequestHeaders, key, value)
}

// AddHttpRequestHeader adds a value for given "key" to the request headers.
// Only available during types.HttpContext.OnHttpRequestHeaders.
func AddHttpRequestHeader(key, value string) error {
	return addMapValue(internal.MapTypeHttpRequestHeaders, key, value)
}

// GetHttpRequestBody is used for retrieving the entire HTTP request body.
// Only available during types.HttpContext.OnHttpRequestBody.
func GetHttpRequestBody(start, maxSize int) ([]byte, error) {
	return getBuffer(internal.BufferTypeHttpRequestBody, start, maxSize)
}

// AppendHttpRequestBody appends the given bytes to the HTTP request body buffer.
// Only available during types.HttpContext.OnHttpRequestBody.
// Please note that you must remove the "content-length" header during OnHttpRequestHeaders.
// Otherwise, the wrong content-length is sent to the upstream and that might result in a client crash.
func AppendHttpRequestBody(data []byte) error {
	return appendToBuffer(internal.BufferTypeHttpRequestBody, data)
}

// PrependHttpRequestBody prepends the given bytes to the HTTP request body buffer.
// Only available during types.HttpContext.OnHttpRequestBody.
// Please note that you must remove the "content-length" header during OnHttpRequestHeaders.
// Otherwise, the wrong content-length is sent to the upstream and that might result in client crash.
func PrependHttpRequestBody(data []byte) error {
	return prependToBuffer(internal.BufferTypeHttpRequestBody, data)
}

// ReplaceHttpRequestBody replaces the HTTP request body buffer with the given bytes.
// Only available during types.HttpContext.OnHttpRequestBody.
// Please note that you must remove the "content-length" header during OnHttpRequestHeaders.
// Otherwise, the wrong content-length is sent to the upstream and that might result in client crash,
// if the size of the data differs from the original size.
func ReplaceHttpRequestBody(data []byte) error {
	return replaceBuffer(internal.BufferTypeHttpRequestBody, data)
}

// GetHttpRequestTrailers is used for retrieving HTTP request trailers.
// Only available during types.HttpContext.OnHttpRequestTrailers and
// types.HttpContext.OnHttpStreamDone.
func GetHttpRequestTrailers() ([][2]string, error) {
	return getMap(internal.MapTypeHttpRequestTrailers)
}

// ReplaceHttpRequestTrailers is used for replacing HTTP request trailers
// with given trailers. Only available during
// types.HttpContext.OnHttpRequestTrailers.
func ReplaceHttpRequestTrailers(trailers [][2]string) error {
	return setMap(internal.MapTypeHttpRequestTrailers, trailers)
}

// GetHttpRequestTrailer is used for retrieving HTTP request trailer value
// for given "key". Only available during types.HttpContext.OnHttpRequestTrailers and
// types.HttpContext.OnHttpStreamDone.
// If multiple values are present for the key, the "first" value found in the host is returned.
// See https://github.com/envoyproxy/envoy/blob/72bf41fb0ecc039f196be02f534bfc2c9c69f348/source/extensions/common/wasm/context.cc#L762-L763
// for detail.
func GetHttpRequestTrailer(key string) (string, error) {
	return getMapValue(internal.MapTypeHttpRequestTrailers, key)
}

// RemoveHttpRequestTrailer removes all values for given "key" from the request trailers.
// Only available during types.HttpContext.OnHttpRequestTrailers.
func RemoveHttpRequestTrailer(key string) error {
	return removeMapValue(internal.MapTypeHttpRequestTrailers, key)
}

// ReplaceHttpRequestTrailer replaces a value for given "key" from the request trailers.
// Only available during types.HttpContext.OnHttpRequestTrailers.
// If multiple values are present for the key, only the "first" value in the host is replaced.
// See https://github.com/envoyproxy/envoy/blob/72bf41fb0ecc039f196be02f534bfc2c9c69f348/envoy/http/header_map.h#L547-L549
// for detail.
func ReplaceHttpRequestTrailer(key, value string) error {
	return replaceMapValue(internal.MapTypeHttpRequestTrailers, key, value)
}

// AddHttpRequestTrailer adds a value for given "key" to the request trailers.
// Only available during types.HttpContext.OnHttpRequestTrailers.
func AddHttpRequestTrailer(key, value string) error {
	return addMapValue(internal.MapTypeHttpRequestTrailers, key, value)
}

// ResumeHttpRequest can be used for resuming HTTP request processing that was stopped
// after returning the types.Action.Pause. Only available during types.HttpContext.
func ResumeHttpRequest() error {
	return internal.StatusToError(internal.ProxyContinueStream(internal.StreamTypeRequest))
}

// GetHttpResponseHeaders is used for retrieving HTTP response headers.
// Only available during types.HttpContext.OnHttpResponseHeaders and
// types.HttpContext.OnHttpStreamDone.
func GetHttpResponseHeaders() ([][2]string, error) {
	return getMap(internal.MapTypeHttpResponseHeaders)
}

// ReplaceHttpResponseHeaders is used for replacing HTTP response headers
// with given headers. Only available during
// types.HttpContext.OnHttpResponseHeaders.
func ReplaceHttpResponseHeaders(headers [][2]string) error {
	return setMap(internal.MapTypeHttpResponseHeaders, headers)
}

// GetHttpResponseHeader is used for retrieving an HTTP response header value
// for a given "key". Only available during types.HttpContext.OnHttpResponseHeaders and
// types.HttpContext.OnHttpStreamDone.
// If multiple values are present for the key, the "first" value found in the host is returned.
// See https://github.com/envoyproxy/envoy/blob/72bf41fb0ecc039f196be02f534bfc2c9c69f348/source/extensions/common/wasm/context.cc#L762-L763
// for detail.
func GetHttpResponseHeader(key string) (string, error) {
	return getMapValue(internal.MapTypeHttpResponseHeaders, key)
}

// RemoveHttpResponseHeader removes all values for given "key" from the response headers.
// Only available during types.HttpContext.OnHttpResponseHeaders.
func RemoveHttpResponseHeader(key string) error {
	return removeMapValue(internal.MapTypeHttpResponseHeaders, key)
}

// ReplaceHttpResponseHeader replaces the value for given "key" from the response headers.
// Only available during types.HttpContext.OnHttpResponseHeaders.
// If multiple values are present for the key, only the "first" value in the host is replaced.
// See https://github.com/envoyproxy/envoy/blob/72bf41fb0ecc039f196be02f534bfc2c9c69f348/envoy/http/header_map.h#L547-L549
// for detail.
func ReplaceHttpResponseHeader(key, value string) error {
	return replaceMapValue(internal.MapTypeHttpResponseHeaders, key, value)
}

// AddHttpResponseHeader adds a value for a given "key" to the response headers.
// Only available during types.HttpContext.OnHttpResponseHeaders.
func AddHttpResponseHeader(key, value string) error {
	return addMapValue(internal.MapTypeHttpResponseHeaders, key, value)
}

// GetHttpResponseBody is used for retrieving the entire HTTP response body.
// Only available during types.HttpContext.OnHttpResponseBody.
func GetHttpResponseBody(start, maxSize int) ([]byte, error) {
	return getBuffer(internal.BufferTypeHttpResponseBody, start, maxSize)
}

// AppendHttpResponseBody appends the given bytes to the HTTP response body buffer.
// Only available during types.HttpContext.OnHttpResponseBody.
// Please note that you must remove the "content-length" header during OnHttpResponseHeaders.
// Otherwise, the wrong content-length is sent to the upstream and that might result in client crash.
func AppendHttpResponseBody(data []byte) error {
	return appendToBuffer(internal.BufferTypeHttpResponseBody, data)
}

// PrependHttpResponseBody prepends the given bytes to the HTTP response body buffer.
// Only available during types.HttpContext.OnHttpResponseBody.
// Please note that you must remove the "content-length" header during OnHttpResponseHeaders.
// Otherwise, the wrong content-length is sent to the upstream and that might result in client crash.
func PrependHttpResponseBody(data []byte) error {
	return prependToBuffer(internal.BufferTypeHttpResponseBody, data)
}

// ReplaceHttpResponseBody replaces the http response body buffer with the given bytes.
// Only available during types.HttpContext.OnHttpResponseBody.
// Please note that you must remove "content-length" header during OnHttpResponseHeaders.
// Otherwise, the wrong content-length is sent to the upstream and that might result in client crash
// if the size of the data differs from the original one.
func ReplaceHttpResponseBody(data []byte) error {
	return replaceBuffer(internal.BufferTypeHttpResponseBody, data)
}

// GetHttpResponseTrailers is used for retrieving HTTP response trailers.
// Only available during types.HttpContext.OnHttpResponseTrailers and
// types.HttpContext.OnHttpStreamDone.
func GetHttpResponseTrailers() ([][2]string, error) {
	return getMap(internal.MapTypeHttpResponseTrailers)
}

// ReplaceHttpResponseTrailers is used for replacing HTTP response trailers
// with given trailers. Only available during
// types.HttpContext.OnHttpResponseTrailers.
func ReplaceHttpResponseTrailers(trailers [][2]string) error {
	return setMap(internal.MapTypeHttpResponseTrailers, trailers)
}

// GetHttpResponseTrailer is used for retrieving an HTTP response trailer value
// for a given "key". Only available during types.HttpContext.OnHttpResponseTrailers and
// types.HttpContext.OnHttpStreamDone.
// If multiple values are present for the key, the "first" value found in the host is returned.
// See https://github.com/envoyproxy/envoy/blob/72bf41fb0ecc039f196be02f534bfc2c9c69f348/source/extensions/common/wasm/context.cc#L762-L763
// for detail.
func GetHttpResponseTrailer(key string) (string, error) {
	return getMapValue(internal.MapTypeHttpResponseTrailers, key)
}

// RemoveHttpResponseTrailer removes all values for given "key" from the response trailers.
// Only available during types.HttpContext.OnHttpResponseTrailers.
func RemoveHttpResponseTrailer(key string) error {
	return removeMapValue(internal.MapTypeHttpResponseTrailers, key)
}

// ReplaceHttpResponseTrailer replaces a value for given "key" from the response trailers.
// Only available during types.HttpContext.OnHttpResponseHeaders.
// If multiple values are present for the key, only the "first" value in the host is replaced.
// See https://github.com/envoyproxy/envoy/blob/72bf41fb0ecc039f196be02f534bfc2c9c69f348/envoy/http/header_map.h#L547-L549
// for detail.
func ReplaceHttpResponseTrailer(key, value string) error {
	return replaceMapValue(internal.MapTypeHttpResponseTrailers, key, value)
}

// AddHttpResponseTrailer adds a value for given "key" to the response trailers.
// Only available during types.HttpContext.OnHttpResponseHeaders.
func AddHttpResponseTrailer(key, value string) error {
	return addMapValue(internal.MapTypeHttpResponseTrailers, key, value)
}

// ResumeHttpResponse can be used to resume the HTTP response processing that was stopped
// after returning types.Action.Pause. Only available during types.HttpContext.
func ResumeHttpResponse() error {
	return internal.StatusToError(internal.ProxyContinueStream(internal.StreamTypeResponse))
}

// SendHttpResponse sends an HTTP response to the downstream with given information (headers, statusCode, body).
// The function returns an error if used outside the types.HttpContext.
// You cannot use this call after types.HttpContext.OnHttpResponseHeaders returns
// Continue because the response headers may have already arrived downstream, and
// there is no way to override the headers that were already sent.
// After you've invoked this function, you *must* return types.Action.Pause to
// stop further processing of the initial HTTP request/response.
// Note that the gRPCStatus can be set to -1 if this is a not gRPC stream.
func SendHttpResponse(statusCode uint32, headers [][2]string, body []byte, gRPCStatus int32) error {
	shs := internal.SerializeMap(headers)
	var bp *byte
	if len(body) > 0 {
		bp = &body[0]
	}
	hp := &shs[0]
	hl := len(shs)
	return internal.StatusToError(
		internal.ProxySendLocalResponse(
			statusCode, nil, 0,
			bp, len(body), hp, hl, gRPCStatus,
		),
	)
}

// GetSharedData is used for retrieving the value for given "key".
// For thread-safe updates you must use the returned "cas" value
// when calling SetSharedData for the same key.
func GetSharedData(key string) (value []byte, cas uint32, err error) {
	var raw *byte
	var size int

	st := internal.ProxyGetSharedData(internal.StringBytePtr(key), len(key), &raw, &size, &cas)
	if st != internal.StatusOK {
		return nil, 0, internal.StatusToError(st)
	}
	return internal.RawBytePtrToByteSlice(raw, size), cas, nil
}

// SetSharedData is used for setting key-value pairs in the shared data storage
// which is defined per "vm_config.vm_id" in the hosts.
// The function returns ErrorStatusCasMismatch when the given CAS value
// doesn't match the current value. The error indicates other Wasm VMs
// have already set a value on the same key, and the current CAS for the key gets incremented.
// Implementing a retry logic to handle this error is recommended.
// Setting the cas value to 0 will never return ErrorStatusCasMismatch,
// and the call will always succeed. However, it is not thread-safe.
// Another VM may have already incremented the value, and the value you
// see is already different from the one stored when you call this function.
func SetSharedData(key string, data []byte, cas uint32) error {
	var dataPtr *byte
	if len(data) > 0 { // Empty data is allowed to set, so we need this check.
		dataPtr = &data[0]
	}
	st := internal.ProxySetSharedData(internal.StringBytePtr(key), len(key), dataPtr, len(data), cas)
	return internal.StatusToError(st)
}

// GetProperty is used for retrieving property/metadata in the host
// for a given path.
// Available path and properties depend on the host implementation.
// For Envoy, please refer to https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/advanced/attributes
//
// Note: if the target property is map-type, use GetPropertyMap instead. This GetProperty returns
// raw un-serialized bytes for such properties. For example, if you have the metadata as
//
//	clusters:
//	- name: web_service
//	  metadata:
//	   filter_metadata:
//	     foo:
//	       my_value: '1234'
//	       my_map:
//	         k1: v1
//	         k2: v2
//
// Then,
//   - use GetPropertyMap for {"cluster_metadata", "filter_metadata", "foo", "my_map"}.
//   - use GetProperty    for {"cluster_metadata", "filter_metadata", "foo", "my_value"}.
//
// Note: you cannot get the raw bytes of protobuf. For example, accessing {"cluster_metadata", "filter_data", "foo"}) doesn't
// return the protobuf bytes, but instead this returns the serialized map of "foo". Therefore, we recommend to access individual
// "leaf" fields (not the middle or top field of metadata) to avoid the need to figure out the (host-dependent) encoding of properties.
func GetProperty(path []string) ([]byte, error) {
	if len(path) == 0 {
		return nil, errors.New("path must not be empty")
	}
	var ret *byte
	var retSize int
	raw := internal.SerializePropertyPath(path)

	err := internal.StatusToError(internal.ProxyGetProperty(&raw[0], len(raw), &ret, &retSize))
	if err != nil {
		return nil, err
	}
	return internal.RawBytePtrToByteSlice(ret, retSize), nil
}

// GetPropertyMap is the same as GetProperty but can be used to decode map-typed properties.
// See GetProperty for detail.
//
// Note: this operates under the assumption that the path is encoded as map, therefore this might
// cause panic if it is used on the non-map types.
func GetPropertyMap(path []string) ([][2]string, error) {
	b, err := GetProperty(path)
	if err != nil {
		return nil, err
	}

	return internal.DeserializeMap(b), nil
}

// SetProperty is used for setting property/metadata in the host
// for a given path.
// Available path and properties depend on the host implementation.
// For Envoy, please refer to https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/advanced/attributes
func SetProperty(path []string, data []byte) error {
	if len(path) == 0 {
		return errors.New("path must not be empty")
	} else if len(data) == 0 {
		return errors.New("data must not be empty")
	}
	raw := internal.SerializePropertyPath(path)
	return internal.StatusToError(internal.ProxySetProperty(
		&raw[0], len(raw), &data[0], len(data),
	))
}

// CallForeignFunction calls a foreign function of given funcName defined by host implementations.
// Foreign functions are host-specific functions, so please refer to the doc of your host implementation for detail.
func CallForeignFunction(funcName string, param []byte) (ret []byte, err error) {
	f := internal.StringBytePtr(funcName)

	var returnData *byte
	var returnSize int

	var paramPtr *byte
	if len(param) != 0 {
		paramPtr = &param[0]
	}

	switch st := internal.ProxyCallForeignFunction(f, len(funcName), paramPtr, len(param), &returnData, &returnSize); st {
	case internal.StatusOK:
		return internal.RawBytePtrToByteSlice(returnData, returnSize), nil
	default:
		return nil, internal.StatusToError(st)
	}
}

// LogTrace emits a message as a log with Trace log level.
func LogTrace(msg string) {
	internal.ProxyLog(internal.LogLevelTrace, internal.StringBytePtr(msg), len(msg))
}

// LogTracef formats according to a format specifier and emits as a log with Trace log level.
//
// Note that not all combinations of format and args are supported by tinygo.
// For example, %v with a map will cause a panic. See
// https://tinygo.org/docs/reference/lang-support/stdlib/#fmt for more
// information.
func LogTracef(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	internal.ProxyLog(internal.LogLevelTrace, internal.StringBytePtr(msg), len(msg))
}

// LogDebug emits a message as a log with Debug log level.
func LogDebug(msg string) {
	internal.ProxyLog(internal.LogLevelDebug, internal.StringBytePtr(msg), len(msg))
}

// LogDebugf formats according to a format specifier and emits as a log with Debug log level.
//
// Note that not all combinations of format and args are supported by tinygo.
// For example, %v with a map will cause a panic. See
// https://tinygo.org/docs/reference/lang-support/stdlib/#fmt for more
// information.
func LogDebugf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	internal.ProxyLog(internal.LogLevelDebug, internal.StringBytePtr(msg), len(msg))
}

// LogInfo emits a message as a log with Info log level.
func LogInfo(msg string) {
	internal.ProxyLog(internal.LogLevelInfo, internal.StringBytePtr(msg), len(msg))
}

// LogInfof formats according to a format specifier and emits as a log with Info log level.
//
// Note that not all combinations of format and args are supported by tinygo.
// For example, %v with a map will cause a panic. See
// https://tinygo.org/docs/reference/lang-support/stdlib/#fmt for more
// information.
func LogInfof(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	internal.ProxyLog(internal.LogLevelInfo, internal.StringBytePtr(msg), len(msg))
}

// LogWarn emits a message as a log with Warn log level.
func LogWarn(msg string) {
	internal.ProxyLog(internal.LogLevelWarn, internal.StringBytePtr(msg), len(msg))
}

// LogWarnf formats according to a format specifier and emits as a log with Warn log level.
//
// Note that not all combinations of format and args are supported by tinygo.
// For example, %v with a map will cause a panic. See
// https://tinygo.org/docs/reference/lang-support/stdlib/#fmt for more
// information.
func LogWarnf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	internal.ProxyLog(internal.LogLevelWarn, internal.StringBytePtr(msg), len(msg))
}

// LogError emits a message as a log with Error log level.
func LogError(msg string) {
	internal.ProxyLog(internal.LogLevelError, internal.StringBytePtr(msg), len(msg))
}

// LogErrorf formats according to a format specifier and emits as a log with Error log level.
//
// Note that not all combinations of format and args are supported by tinygo.
// For example, %v with a map will cause a panic. See
// https://tinygo.org/docs/reference/lang-support/stdlib/#fmt for more
// information.
func LogErrorf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	internal.ProxyLog(internal.LogLevelError, internal.StringBytePtr(msg), len(msg))
}

// LogCritical emits a message as a log with Critical log level.
func LogCritical(msg string) {
	internal.ProxyLog(internal.LogLevelCritical, internal.StringBytePtr(msg), len(msg))
}

// LogCriticalf formats according to a format specifier and emits as a log with Critical log level.
//
// Note that not all combinations of format and args are supported by tinygo.
// For example, %v with a map will cause a panic. See
// https://tinygo.org/docs/reference/lang-support/stdlib/#fmt for more
// information.
func LogCriticalf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	internal.ProxyLog(internal.LogLevelCritical, internal.StringBytePtr(msg), len(msg))
}

type (
	// MetricCounter represents a counter metric.
	// Use DefineCounterMetric for initialization.
	MetricCounter uint32
	// MetricGauge represents a gauge metric.
	// Use DefineGaugeMetric for initialization.
	MetricGauge uint32
	// MetricHistogram represents a histogram metric.
	// Use DefineHistogramMetric for initialization.
	MetricHistogram uint32
)

// DefineCounterMetric returns MetricCounter for a name.
func DefineCounterMetric(name string) MetricCounter {
	var id uint32
	ptr := internal.StringBytePtr(name)
	st := internal.ProxyDefineMetric(internal.MetricTypeCounter, ptr, len(name), &id)
	if err := internal.StatusToError(st); err != nil {
		panic(fmt.Sprintf("define metric of name %s: %v", name, internal.StatusToError(st)))
	}
	return MetricCounter(id)
}

// Value returns the current value for this counter.
func (m MetricCounter) Value() uint64 {
	var val uint64
	st := internal.ProxyGetMetric(uint32(m), &val)
	if err := internal.StatusToError(st); err != nil {
		panic(fmt.Sprintf("get metric of  %d: %v", uint32(m), internal.StatusToError(st)))
	}
	return val
}

// Increment increments the current value by an offset for this counter.
func (m MetricCounter) Increment(offset uint64) {
	if err := internal.StatusToError(internal.ProxyIncrementMetric(uint32(m), int64(offset))); err != nil {
		panic(fmt.Sprintf("increment %d by %d: %v", uint32(m), offset, err))
	}
}

// DefineCounterMetric returns MetricGauge for a name.
func DefineGaugeMetric(name string) MetricGauge {
	var id uint32
	ptr := internal.StringBytePtr(name)
	st := internal.ProxyDefineMetric(internal.MetricTypeGauge, ptr, len(name), &id)
	if err := internal.StatusToError(st); err != nil {
		panic(fmt.Sprintf("error define metric of name %s: %v", name, internal.StatusToError(st)))
	}
	return MetricGauge(id)
}

// Value returns the current value for this gauge.
func (m MetricGauge) Value() int64 {
	var val uint64
	if err := internal.StatusToError(internal.ProxyGetMetric(uint32(m), &val)); err != nil {
		panic(fmt.Sprintf("get metric of  %d: %v", uint32(m), err))
	}
	return int64(val)
}

// Add adds an offset to the current value for this gauge.
func (m MetricGauge) Add(offset int64) {
	if err := internal.StatusToError(internal.ProxyIncrementMetric(uint32(m), offset)); err != nil {
		panic(fmt.Sprintf("error adding %d by %d: %v", uint32(m), offset, err))
	}
}

// DefineHistogramMetric returns MetricHistogram for a name.
func DefineHistogramMetric(name string) MetricHistogram {
	var id uint32
	ptr := internal.StringBytePtr(name)
	st := internal.ProxyDefineMetric(internal.MetricTypeHistogram, ptr, len(name), &id)
	if err := internal.StatusToError(st); err != nil {
		panic(fmt.Sprintf("error define metric of name %s: %v", name, internal.StatusToError(st)))
	}
	return MetricHistogram(id)
}

// Value returns the current value for this histogram.
func (m MetricHistogram) Value() uint64 {
	var val uint64
	st := internal.ProxyGetMetric(uint32(m), &val)
	if err := internal.StatusToError(st); err != nil {
		panic(fmt.Sprintf("get metric of  %d: %v", uint32(m), internal.StatusToError(st)))
	}
	return val
}

// Record records a value for this histogram.
func (m MetricHistogram) Record(value uint64) {
	if err := internal.StatusToError(internal.ProxyRecordMetric(uint32(m), value)); err != nil {
		panic(fmt.Sprintf("error adding %d: %v", uint32(m), err))
	}
}

func setMap(mapType internal.MapType, headers [][2]string) error {
	shs := internal.SerializeMap(headers)
	hp := &shs[0]
	hl := len(shs)
	return internal.StatusToError(internal.ProxySetHeaderMapPairs(mapType, hp, hl))
}

func getMapValue(mapType internal.MapType, key string) (string, error) {
	var rvs int
	var raw *byte
	if st := internal.ProxyGetHeaderMapValue(
		mapType, internal.StringBytePtr(key), len(key), &raw, &rvs,
	); st != internal.StatusOK {
		return "", internal.StatusToError(st)
	}

	ret := internal.RawBytePtrToString(raw, rvs)
	return ret, nil
}

func removeMapValue(mapType internal.MapType, key string) error {
	return internal.StatusToError(
		internal.ProxyRemoveHeaderMapValue(mapType, internal.StringBytePtr(key), len(key)),
	)
}

func replaceMapValue(mapType internal.MapType, key, value string) error {
	return internal.StatusToError(
		internal.ProxyReplaceHeaderMapValue(
			mapType, internal.StringBytePtr(key), len(key), internal.StringBytePtr(value), len(value),
		),
	)
}

func addMapValue(mapType internal.MapType, key, value string) error {
	return internal.StatusToError(
		internal.ProxyAddHeaderMapValue(
			mapType, internal.StringBytePtr(key), len(key), internal.StringBytePtr(value), len(value),
		),
	)
}

func getMap(mapType internal.MapType) ([][2]string, error) {
	var rvs int
	var raw *byte

	st := internal.ProxyGetHeaderMapPairs(mapType, &raw, &rvs)
	if st != internal.StatusOK {
		return nil, internal.StatusToError(st)
	} else if raw == nil {
		return nil, types.ErrorStatusNotFound
	}

	bs := internal.RawBytePtrToByteSlice(raw, rvs)
	return internal.DeserializeMap(bs), nil
}

func getBuffer(bufType internal.BufferType, start, maxSize int) ([]byte, error) {
	var retData *byte
	var retSize int
	switch st := internal.ProxyGetBufferBytes(bufType, start, maxSize, &retData, &retSize); st {
	case internal.StatusOK:
		if retData == nil {
			return nil, types.ErrorStatusNotFound
		}
		return internal.RawBytePtrToByteSlice(retData, retSize), nil
	default:
		return nil, internal.StatusToError(st)
	}
}

func appendToBuffer(bufType internal.BufferType, buffer []byte) error {
	var bufferData *byte
	if len(buffer) != 0 {
		bufferData = &buffer[0]
	}
	return internal.StatusToError(internal.ProxySetBufferBytes(bufType, math.MaxInt32, 0, bufferData, len(buffer)))
}

func prependToBuffer(bufType internal.BufferType, buffer []byte) error {
	var bufferData *byte
	if len(buffer) != 0 {
		bufferData = &buffer[0]
	}
	return internal.StatusToError(internal.ProxySetBufferBytes(bufType, 0, 0, bufferData, len(buffer)))
}

func replaceBuffer(bufType internal.BufferType, buffer []byte) error {
	var bufferData *byte
	if len(buffer) != 0 {
		bufferData = &buffer[0]
	}
	return internal.StatusToError(
		internal.ProxySetBufferBytes(bufType, 0, math.MaxInt32, bufferData, len(buffer)),
	)
}
