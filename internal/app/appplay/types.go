// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package appplay

import (
	"time"
)

type Snippet struct {
	Body []byte
}

type Request struct {
	Body string
}

type Response struct {
	Errors string
	Events []Event
}

type Event struct {
	Message string
	Kind    string        // "stdout" or "stderr"
	Delay   time.Duration // time to wait before printing Message
}

type editData struct {
	Snippet *Snippet
}

type fmtResponse struct {
	Body  string
	Error string
}
