// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package appplay

import (
	"encoding/json"
	"fmt"
	"net/http"

	"wa-lang.org/wa/api"
	"wa-lang.org/wa/internal/token"
	"wa-lang.org/wa/internal/xlang"
)

func (p *WebServer) fmtHandler(w http.ResponseWriter, r *http.Request) {
	var (
		in  = []byte(r.FormValue("body"))
		err error
	)

	resp, err := p.fmtCode(in)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, fmt.Sprintf("error encoding response: %v", err), http.StatusInternalServerError)
		return
	}
}

func (p *WebServer) fmtCode(code []byte) (*fmtResponse, error) {
	var filename string
	switch xlang.DetectLang("", code) {
	case token.LangType_Wa:
		filename = "prog.wa"
	case token.LangType_Wz:
		filename = "prog.wz"
	default:
		return &fmtResponse{
			Body: string(code),
		}, nil
	}

	output, err := api.FormatCode(filename, string(code))
	if err != nil {
		resp := &fmtResponse{
			Error: err.Error(),
		}
		return resp, nil
	}

	resp := &fmtResponse{
		Body: string(output),
	}

	return resp, nil
}
