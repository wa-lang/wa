// 版权 @2023 凹语言 作者。保留所有权利。

package lsp

import (
	"encoding/json"
	"fmt"
	"io"
)

type lspChannel struct {
	r   io.Reader
	w   io.Writer
	dec *json.Decoder
}

type lspRequest struct {
	ID      int               `json:"id"`
	JSONRPC string            `json:"jsonrpc"`
	Method  string            `json:"method"`
	Params  []json.RawMessage `json:"params"`
}

type lspResponse struct {
	ID      int         `json:"id"`
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
	Error   *lspError   `json:"error"`
}

// JSONRPC错误
type lspError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (err *lspError) Error() string {
	return fmt.Sprintf("lsp.Error{code:%d, message: %q}", err.Code, err.Message)
}

func (p *lspChannel) RecvRequest() (req *lspRequest, err error) {
	if p.dec == nil {
		p.dec = json.NewDecoder(p.r)
	}
	req = new(lspRequest)
	err = p.dec.Decode(&req)
	return
}

func (p *lspChannel) SendRespose(resp *lspResponse) error {
	reqBody, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	_, err = p.w.Write(reqBody)
	if err != nil {
		return err
	}
	return nil
}
