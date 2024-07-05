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

//go:build tinygo

package internal

//export proxy_log
func ProxyLog(logLevel LogLevel, messageData *byte, messageSize int) Status

//export proxy_send_local_response
func ProxySendLocalResponse(statusCode uint32, statusCodeDetailData *byte, statusCodeDetailsSize int,
	bodyData *byte, bodySize int, headersData *byte, headersSize int, grpcStatus int32) Status

//export proxy_get_shared_data
func ProxyGetSharedData(keyData *byte, keySize int, returnValueData **byte, returnValueSize *int, returnCas *uint32) Status

//export proxy_set_shared_data
func ProxySetSharedData(keyData *byte, keySize int, valueData *byte, valueSize int, cas uint32) Status

//export proxy_register_shared_queue
func ProxyRegisterSharedQueue(nameData *byte, nameSize int, returnID *uint32) Status

//export proxy_resolve_shared_queue
func ProxyResolveSharedQueue(vmIDData *byte, vmIDSize int, nameData *byte, nameSize int, returnID *uint32) Status

//export proxy_dequeue_shared_queue
func ProxyDequeueSharedQueue(queueID uint32, returnValueData **byte, returnValueSize *int) Status

//export proxy_enqueue_shared_queue
func ProxyEnqueueSharedQueue(queueID uint32, valueData *byte, valueSize int) Status

//export proxy_get_header_map_value
func ProxyGetHeaderMapValue(mapType MapType, keyData *byte, keySize int, returnValueData **byte, returnValueSize *int) Status

//export proxy_add_header_map_value
func ProxyAddHeaderMapValue(mapType MapType, keyData *byte, keySize int, valueData *byte, valueSize int) Status

//export proxy_replace_header_map_value
func ProxyReplaceHeaderMapValue(mapType MapType, keyData *byte, keySize int, valueData *byte, valueSize int) Status

//export proxy_remove_header_map_value
func ProxyRemoveHeaderMapValue(mapType MapType, keyData *byte, keySize int) Status

//export proxy_get_header_map_pairs
func ProxyGetHeaderMapPairs(mapType MapType, returnValueData **byte, returnValueSize *int) Status

//export proxy_set_header_map_pairs
func ProxySetHeaderMapPairs(mapType MapType, mapData *byte, mapSize int) Status

//export proxy_get_buffer_bytes
func ProxyGetBufferBytes(bufferType BufferType, start int, maxSize int, returnBufferData **byte, returnBufferSize *int) Status

//export proxy_set_buffer_bytes
func ProxySetBufferBytes(bufferType BufferType, start int, maxSize int, bufferData *byte, bufferSize int) Status

//export proxy_continue_stream
func ProxyContinueStream(streamType StreamType) Status

//export proxy_close_stream
func ProxyCloseStream(streamType StreamType) Status

//export proxy_http_call
func ProxyHttpCall(upstreamData *byte, upstreamSize int, headerData *byte, headerSize int,
	bodyData *byte, bodySize int, trailersData *byte, trailersSize int, timeout uint32, calloutIDPtr *uint32,
) Status

//export proxy_call_foreign_function
func ProxyCallForeignFunction(funcNamePtr *byte, funcNameSize int, paramPtr *byte, paramSize int, returnData **byte, returnSize *int) Status

//export proxy_set_tick_period_milliseconds
func ProxySetTickPeriodMilliseconds(period uint32) Status

//export proxy_set_effective_context
func ProxySetEffectiveContext(contextID uint32) Status

//export proxy_done
func ProxyDone() Status

//export proxy_define_metric
func ProxyDefineMetric(metricType MetricType, metricNameData *byte, metricNameSize int, returnMetricIDPtr *uint32) Status

//export proxy_increment_metric
func ProxyIncrementMetric(metricID uint32, offset int64) Status

//export proxy_record_metric
func ProxyRecordMetric(metricID uint32, value uint64) Status

//export proxy_get_metric
func ProxyGetMetric(metricID uint32, returnMetricValue *uint64) Status

//export proxy_get_property
func ProxyGetProperty(pathData *byte, pathSize int, returnValueData **byte, returnValueSize *int) Status

//export proxy_set_property
func ProxySetProperty(pathData *byte, pathSize int, valueData *byte, valueSize int) Status
