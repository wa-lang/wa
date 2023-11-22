// 版权 @2023 凹语言 作者。保留所有权利。

package appplay

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"wa-lang.org/wa/api"
	"wa-lang.org/wa/internal/wamime"
)

func (p *WebServer) runHandler(w http.ResponseWriter, r *http.Request) {
	var req Request
	version := r.PostFormValue("version")
	if version == "2" {
		req.Body = r.PostFormValue("body")
	} else {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, fmt.Sprintf("error decoding request: %v", err), http.StatusBadRequest)
			return
		}
	}
	resp, err := p.compileAndRun(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, fmt.Sprintf("error encoding response: %v", err), http.StatusInternalServerError)
		return
	}
}

func (p *WebServer) compileAndRun(req *Request) (*Response, error) {
	tmpDir, err := ioutil.TempDir("", "sandbox")
	if err != nil {
		return nil, fmt.Errorf("error creating temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	filename := "prog.wa"
	if wamime.GetCodeMime(filename, []byte(req.Body)) == "wz" {
		filename = "prog.wz"
	}

	result, err := api.RunCode(api.DefaultConfig(), filename, req.Body, "")
	if err != nil {
		resp := &Response{Errors: err.Error()}
		return resp, nil
	}

	resp := &Response{
		Events: []Event{
			{
				Message: string(result),
				Kind:    "stdout",
			},
		},
	}

	return resp, nil
}
