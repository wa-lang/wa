// 版权 @2023 凹语言 作者。保留所有权利。

package lsp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type lspRequest struct {
	ID      int           `json:"id"`
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

type lspResponse struct {
	ID      int             `json:"id"`
	JSONRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	Error   *lspError       `json:"error"`
}

// JSONRPC错误
type lspError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (err *lspError) Error() string {
	return fmt.Sprintf("lsp.Error{code:%d, message: %q}", err.Code, err.Message)
}

// 原生调用方法
func lspCall(host, method string, args ...interface{}) (result json.RawMessage, err error) {
	return jsonrpcCall(host, method, args...)
}

// jsonrpc 调用
func jsonrpcCall(host, method string, args ...interface{}) (result json.RawMessage, err error) {
	reqBody, err := json.Marshal(&lspRequest{
		JSONRPC: "2.0",
		Method:  method,
		Params:  args,
		ID:      1,
	})
	if err != nil {
		return nil, err
	}

	r, err := http.Post(host, "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	respBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var resp = new(lspResponse)
	if err := json.Unmarshal(respBody, resp); err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	return resp.Result, nil
}

// JSONRPC调用, 需要通过response提供返回类型
func lsprpcCall(host, method string, response interface{}, args ...interface{}) (err error) {
	resp_Result, err := jsonrpcCall(host, method, args...)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(resp_Result, response); err != nil {
		return err
	}
	return nil
}
