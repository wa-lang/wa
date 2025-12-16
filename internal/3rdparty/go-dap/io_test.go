// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dap

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"reflect"
	"strings"
	"testing"
	"time"
)

func Test_WriteBaseMessage(t *testing.T) {
	tests := []struct {
		input       string
		wantWritten string
		wantErr     error
	}{
		{``, "Content-Length: 0\r\n\r\n", nil},
		{`a`, "Content-Length: 1\r\n\r\na", nil},
		{`{}`, "Content-Length: 2\r\n\r\n{}", nil},
		{`{"a":0 "b":"blah"}`, "Content-Length: 18\r\n\r\n{\"a\":0 \"b\":\"blah\"}", nil},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			var buf bytes.Buffer
			gotErr := WriteBaseMessage(&buf, []byte(test.input))
			gotWritten := buf.String()
			if gotErr != test.wantErr {
				t.Errorf("got err=%#v, want %#v", gotErr, test.wantErr)
			}
			if gotErr == nil && gotWritten != test.wantWritten {
				t.Errorf("got written=%q, want %q", gotWritten, test.wantWritten)
			}
		})
	}
}

func Test_ReadBaseMessage(t *testing.T) {
	tests := []struct {
		input         string
		wantBytesRead []byte
		wantBytesLeft []byte
		wantErr       error
	}{
		{"", nil, []byte(""), io.EOF},
		{"random stuff\r\nabc", nil, []byte("c"), ErrHeaderDelimiterNotCrLfCrLf},
		{"Cache-Control: no-cache\r\n\r\n", nil, []byte(""), ErrHeaderNotContentLength},
		{"Content-Length 1\r\n\r\nabc", nil, []byte("abc"), ErrHeaderNotContentLength},
		{"Content-Length: 10\r\n\r\nabc", nil, []byte(""), io.ErrUnexpectedEOF},
		{"Content-Length: 3\r\n\r\nabc", []byte("abc"), []byte(""), nil},
		{"Content-Length: 4194305\r\n\r\nabc", nil, []byte("abc"), ErrHeaderContentTooLong},
		{"Content-Length: 6506440440440\r\n\r\nabc", nil, []byte("abc"), ErrHeaderContentTooLong},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(test.input))
			gotBytes, gotErr := ReadBaseMessage(reader)
			if gotErr != test.wantErr {
				t.Errorf("got err=%#v, want %#v", gotErr, test.wantErr)
			}
			if gotErr == nil && !bytes.Equal(gotBytes, test.wantBytesRead) {
				t.Errorf("got bytes=%q, want %q", gotBytes, test.wantBytesRead)
			}
			bytesLeft, _ := io.ReadAll(reader)
			if !bytes.Equal(bytesLeft, test.wantBytesLeft) {
				t.Errorf("got bytesLeft=%q, want %q", bytesLeft, test.wantBytesLeft)
			}
		})
	}
}

func Test_readContentLengthHeader(t *testing.T) {
	tests := []struct {
		input         string
		wantBytesLeft string // Bytes left in the reader after header reading
		wantLen       int64  // Extracted content length value
		wantErr       error
	}{
		{"", "", 0, io.EOF},
		{"Cache-Control: no-cache", "", 0, io.EOF},
		{"Cache-Control: no-cache\r", "", 0, io.EOF},
		{"Cache-Control: no-cache\rabc", "", 0, ErrHeaderDelimiterNotCrLfCrLf},
		{"Cache-Control: no-cache\r\n", "", 0, io.ErrUnexpectedEOF},
		{"Cache-Control: no-cache\r\n\r", "", 0, io.ErrUnexpectedEOF},
		{"Cache-Control: no-cache\r\n\r\n", "", 0, ErrHeaderNotContentLength},
		{"Cache-Control: no-cache\r\n\r\nabc", "abc", 0, ErrHeaderNotContentLength},
		{"Content-Length: 3 abc", "", 0, io.EOF},
		{"Content-Length: 3\nabc", "", 0, io.EOF},
		{"Content-Length: 3\rabc", "", 0, ErrHeaderDelimiterNotCrLfCrLf},
		{"Content-Length: 3\r\nabc", "c", 0, ErrHeaderDelimiterNotCrLfCrLf},
		{"Content-Length: 3\r\n\rabc", "bc", 0, ErrHeaderDelimiterNotCrLfCrLf},
		{"Content-Length: 3\r \n\r\nabc", "\nabc", 0, ErrHeaderDelimiterNotCrLfCrLf},
		{"Content-Length: 3\r\n \r\nabc", "\nabc", 0, ErrHeaderDelimiterNotCrLfCrLf},
		{"Content-Length: 3\r\n\r \nabc", "\nabc", 0, ErrHeaderDelimiterNotCrLfCrLf},
		{"Content-Length 3\r\n\r\nabc", "abc", 0, ErrHeaderNotContentLength},
		{"_Content-Length: 3\r\n\r\nabc", "abc", 0, ErrHeaderNotContentLength},
		{"Content-Length: 3_\r\n\r\nabc", "abc", 0, ErrHeaderNotContentLength},
		{"Content-Length: x\r\n\r\nabc", "abc", 0, ErrHeaderNotContentLength},
		{"Content-Length: 3.0\r\n\r\nabc", "abc", 0, ErrHeaderNotContentLength},
		{"Content-Length: -3\r\n\r\nabc", "abc", 0, ErrHeaderNotContentLength},
		{"Content-Length: 0\r\n\r\nabc", "abc", 0, nil},
		{"Content-Length: 3\r\n\r\nabc", "abc", 3, nil},
		{"Content-Length: 9223372036854775807\r\n\r\nabc", "abc", 9223372036854775807, nil},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(test.input))
			gotLen, gotErr := readContentLengthHeader(reader)
			if gotErr != test.wantErr {
				t.Errorf("got err=%#v, want %#v", gotErr, test.wantErr)
			}
			if gotErr == nil && gotLen != test.wantLen {
				t.Errorf("got len=%d, want %d", gotLen, test.wantLen)
			}
			bytesLeft, _ := io.ReadAll(reader)
			if string(bytesLeft) != test.wantBytesLeft {
				t.Errorf("got bytesLeft=%q, want %q", bytesLeft, test.wantBytesLeft)
			}
		})
	}
}

func TestWriteRead(t *testing.T) {
	writeContent := [][]byte{
		[]byte("this is"),
		[]byte("a read write"),
		[]byte("test"),
	}
	var buf bytes.Buffer
	for _, wc := range writeContent {
		if err := WriteBaseMessage(&buf, wc); err != nil {
			t.Fatal(err)
		}
	}
	reader := bufio.NewReader(&buf)
	for _, wc := range writeContent {
		rc, err := ReadBaseMessage(reader)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(rc, wc) {
			t.Fatalf("got %q, want %q", rc, wc)
		}
	}
}

// readMessagesAndNotify reads messages one by one until EOF.
// Notifies of a read via messagesRead channel.
func readMessagesAndNotify(t *testing.T, r io.Reader, messagesRead chan<- []byte) {
	reader := bufio.NewReader(r)
	for {
		msg, err := ReadBaseMessage(reader)
		if err == io.EOF {
			close(messagesRead)
			break
		}
		// On error, this will send "" as the content read
		messagesRead <- msg
	}
}

func writeOrFail(t *testing.T, w io.Writer, data string) {
	if n, err := w.Write([]byte(data)); err != nil || n < len(data) {
		t.Fatal(err)
	}
}

func TestReadMessageInParts(t *testing.T) {
	// This test uses separate goroutines to write and read messages
	// and relies on blocking channel operations between them to ensure that
	// the expected number of messages is read for what is written.
	// Otherwise, the test will time out.
	messagesRead := make(chan []byte)
	r, w := io.Pipe()
	header := "Content-Length: 11"
	baddelim := "\r\r\r\r"
	content1 := "message one"
	content2 := "message two"
	nocontent := ""

	// This will keep blocking to read a full message or EOF.
	go readMessagesAndNotify(t, r, messagesRead)

	checkNoMessageRead := func() {
		time.Sleep(100 * time.Millisecond) // Let read goroutine run
		select {
		case msg := <-messagesRead:
			t.Errorf("got %q, want none", msg)
		default:
		}
	}

	checkOneMessageRead := func(want []byte) {
		got := <-messagesRead
		if !bytes.Equal(got, want) {
			t.Errorf("got %q, want %q", got, want)
		}
	}

	// Good message written in full
	writeOrFail(t, w, header+crLfcrLf+content1)
	checkOneMessageRead([]byte(content1))

	// Good message written in parts
	writeOrFail(t, w, header)
	checkNoMessageRead()
	writeOrFail(t, w, crLfcrLf)
	checkNoMessageRead()
	writeOrFail(t, w, content2)
	checkOneMessageRead([]byte(content2))

	// Bad message written in full
	writeOrFail(t, w, header+baddelim)
	checkOneMessageRead([]byte(nocontent))

	// Bad meassage written in parts
	writeOrFail(t, w, header)
	checkNoMessageRead()
	writeOrFail(t, w, baddelim)
	checkOneMessageRead([]byte(nocontent))

	w.Close() // "sends" EOF
}

func TestReadWriteWithCodec(t *testing.T) {
	// Tests end-to-end write and read from a buffer using the DAP codec.
	req := InitializeRequest{
		Request: Request{
			ProtocolMessage: ProtocolMessage{
				Type: "request",
				Seq:  121,
			},
			Command: "initialize",
		},
		Arguments: InitializeRequestArguments{
			ClientID:                     "vscode",
			ClientName:                   "Visual Studio Code",
			AdapterID:                    "go",
			PathFormat:                   "path",
			LinesStartAt1:                true,
			ColumnsStartAt1:              false,
			SupportsVariableType:         true,
			SupportsVariablePaging:       true,
			SupportsRunInTerminalRequest: false,
			Locale:                       "en-us",
		},
	}

	baseReq, err := json.Marshal(req)
	if err != nil {
		t.Error(err)
	}

	buf := new(bytes.Buffer)
	err = WriteBaseMessage(buf, baseReq)
	if err != nil {
		t.Error(err)
	}

	reader := bufio.NewReader(buf)
	msg, err := ReadBaseMessage(reader)
	if err != nil {
		t.Error(err)
	}

	readReqPtr, err := DecodeProtocolMessage(msg)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(readReqPtr, &req) {
		t.Errorf("got req=%#v, want %#v", readReqPtr, req)
	}
}

const cancelReqString = "Content-Length: 75\r\n\r\n{\"seq\":25,\"type\":\"request\",\"command\":\"cancel\",\"arguments\":{\"requestId\":24}}"

var cancelReqStruct = CancelRequest{
	Request: Request{
		ProtocolMessage: ProtocolMessage{
			Type: "request",
			Seq:  25,
		},
		Command: "cancel",
	},
	Arguments: &CancelArguments{RequestId: 24},
}

func TestWriteProtocolMessage(t *testing.T) {
	buf := new(bytes.Buffer)
	err := WriteProtocolMessage(buf, &cancelReqStruct)
	if err != nil {
		t.Error(err)
	}

	if buf.String() != cancelReqString {
		t.Errorf("got %#v, want %#v", buf.String(), cancelReqString)
	}
}

func TestReadProtocolMessage(t *testing.T) {
	reader := bufio.NewReader(strings.NewReader(cancelReqString))

	msg, err := ReadProtocolMessage(reader)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(msg, &cancelReqStruct) {
		t.Errorf("got req=%#v, want %#v", msg, &cancelReqStruct)
	}
}
