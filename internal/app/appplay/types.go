// 版权 @2023 凹语言 作者。保留所有权利。

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
