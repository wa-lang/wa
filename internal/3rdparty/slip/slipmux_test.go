package slip

import (
	"bytes"
	"strconv"
	"testing"
)

var readDataMux = []struct {
	data     []byte
	expected []byte
	frame    byte // Frame Type
	err      error
}{
	// Diagnostic message
	{[]byte{END, FRAME_DIAGNOSTIC, 'h', 'i', END}, []byte{'h', 'i'}, FRAME_DIAGNOSTIC, nil},

	// IP Frames are return including FrameType byte
	{[]byte{END, FRAME_IPV4_START, 1, 2, 3, 4, 5, 6, END}, []byte{FRAME_IPV4_START, 1, 2, 3, 4, 5, 6}, FRAME_IPV4_START, nil},
	{[]byte{END, FRAME_IPV4_END, 1, 2, 3, 4, 5, 6, END}, []byte{FRAME_IPV4_END, 1, 2, 3, 4, 5, 6}, FRAME_IPV4_END, nil},
	{[]byte{END, FRAME_IPV6_START, 1, 2, 3, 4, 5, 6, END}, []byte{FRAME_IPV6_START, 1, 2, 3, 4, 5, 6}, FRAME_IPV6_START, nil},
	{[]byte{END, FRAME_IPV6_END, 1, 2, 3, 4, 5, 6, END}, []byte{FRAME_IPV6_END, 1, 2, 3, 4, 5, 6}, FRAME_IPV6_END, nil},

	// CoAP Messages are CRC checked
	// Short (< 4 msg + 2 crc = 6 bytes) CoAP messages are ignored
	{[]byte{END, FRAME_COAP, 1, 2, END, FRAME_DIAGNOSTIC, 'x', END}, []byte{'x'}, FRAME_DIAGNOSTIC, nil},
	{[]byte{END, FRAME_COAP, 1, 2, 3, 4, 5, END, FRAME_DIAGNOSTIC, 'x', END}, []byte{'x'}, FRAME_DIAGNOSTIC, nil},
	// CoAP msg with valid checksum - [152 177] is the checksum for [FRAME_COAP, 1 2 3 4]
	{[]byte{END, FRAME_COAP, 1, 2, 3, 4, 152, 177, END, FRAME_DIAGNOSTIC, 0, END}, []byte{1, 2, 3, 4}, FRAME_COAP, nil},
	// CoAP msg with invalid checksum - [145 57] is the checksum for [1 2 3 4]
	{[]byte{END, FRAME_COAP, 1, 2, 3, 4, 145, 57, END, FRAME_DIAGNOSTIC, 0, END}, []byte{0}, FRAME_DIAGNOSTIC, nil},

	// Unknown frame is just accepted
	{[]byte{1, 2, 3, END}, []byte{2, 3}, 1, nil},

	// Reserved frames are ignored silently
	// ------------------------------------
	// data "FRAME_UNKNOWN" is ignored when used as FrameType
	{[]byte{FRAME_UNKNOWN, 2, 3, END, FRAME_DIAGNOSTIC, 'x', END}, []byte{'x'}, FRAME_DIAGNOSTIC, nil},
	// End gets naturally ignored when used as FrameType
	{[]byte{END, 2, 3, END, FRAME_DIAGNOSTIC, 'x', END}, []byte{3}, 2, nil},
	// data "ESC" is ignored when used as FrameType
	{[]byte{ESC, ESC_ESC, 2, 3, END, FRAME_DIAGNOSTIC, 'x', END}, []byte{'x'}, FRAME_DIAGNOSTIC, nil},
	// data "END" is ignored when used as FrameType
	{[]byte{ESC, ESC_END, 2, 3, END, FRAME_DIAGNOSTIC, 'x', END}, []byte{'x'}, FRAME_DIAGNOSTIC, nil},
	// data "ESC_ESC" is NOT ignored when used as FrameType
	{[]byte{ESC_ESC, 2, 3, END, FRAME_DIAGNOSTIC, 'x', END}, []byte{2, 3}, ESC_ESC, nil},
	// data "ESC_END" is NOT ignored when used as FrameType
	{[]byte{ESC_END, 2, 3, END, FRAME_DIAGNOSTIC, 'x', END}, []byte{2, 3}, ESC_END, nil},

	// Non terminated data would lead to blocking read!
}

var writeDataMux = []struct {
	frame    byte
	data     []byte
	expected []byte
	err      error
}{
	// Just data. Starts with END and ends with END
	{FRAME_DIAGNOSTIC, []byte{1, 2, 3}, []byte{END, FRAME_DIAGNOSTIC, 1, 2, 3, END}, nil},
	// Diagnostic messages and Escape sequences
	{FRAME_DIAGNOSTIC, []byte{'h', 'i'}, []byte{END, FRAME_DIAGNOSTIC, 'h', 'i', END}, nil},
	{FRAME_DIAGNOSTIC, []byte{END}, []byte{END, FRAME_DIAGNOSTIC, ESC, ESC_END, END}, nil},
	{FRAME_DIAGNOSTIC, []byte{ESC}, []byte{END, FRAME_DIAGNOSTIC, ESC, ESC_ESC, END}, nil},
	{FRAME_DIAGNOSTIC, []byte{ESC_END}, []byte{END, FRAME_DIAGNOSTIC, ESC_END, END}, nil},
	{FRAME_DIAGNOSTIC, []byte{ESC_ESC}, []byte{END, FRAME_DIAGNOSTIC, ESC_ESC, END}, nil},
	{FRAME_DIAGNOSTIC, []byte{END, ESC}, []byte{END, FRAME_DIAGNOSTIC, ESC, ESC_END, ESC, ESC_ESC, END}, nil},

	// Ip Packets don't get additional frame type byte
	{FRAME_IPV4_START, []byte{FRAME_IPV4_START, 1, 2, 3}, []byte{END, FRAME_IPV4_START, 1, 2, 3, END}, nil},
	{FRAME_IPV4_END, []byte{FRAME_IPV4_END, 1, 2, 3}, []byte{END, FRAME_IPV4_END, 1, 2, 3, END}, nil},
	{FRAME_IPV6_START, []byte{FRAME_IPV6_START, 1, 2, 3}, []byte{END, FRAME_IPV6_START, 1, 2, 3, END}, nil},
	{FRAME_IPV6_END, []byte{FRAME_IPV6_END, 1, 2, 3}, []byte{END, FRAME_IPV6_END, 1, 2, 3, END}, nil},

	// CoAP packets get checksum
	{FRAME_COAP, []byte{1, 2, 3, 4}, []byte{END, FRAME_COAP, 1, 2, 3, 4, 152, 177, END}, nil},

	// Unknown types are handled like all non ip packets
	{1, []byte{1, 2, 3, 4}, []byte{END, 1, 1, 2, 3, 4, END}, nil},
}

func TestReadMux(t *testing.T) {
	for i, d := range readDataMux {
		r := NewSlipMuxReader(bytes.NewReader(d.data))
		p, frame, err := r.ReadPacket()

		if err == nil && d.err != nil {
			t.Error(strconv.Itoa(i), "Expected error", d.err.Error(), "but got", err)
		} else if err != nil && d.err == nil {
			t.Error(strconv.Itoa(i), "Expected error", d.err, "but got", err)
		} else if err != d.err && err.Error() != d.err.Error() {
			t.Error(strconv.Itoa(i), "Expected error", d.err, "but got", err)
		}

		if frame != d.frame {
			t.Error(strconv.Itoa(i), "Expected frame", d.frame, "but got", frame)
		}
		if !eqBytes(p, d.expected) {
			t.Error(strconv.Itoa(i), "Expected data", d.expected, "but got", p)
		}
	}
}

func TestWriteMux(t *testing.T) {
	for i, d := range writeDataMux {
		buf := &bytes.Buffer{}
		w := NewSlipMuxWriter(buf)
		err := w.WritePacket(d.frame, d.data)

		if err == nil && d.err != nil {
			t.Error(strconv.Itoa(i), "Expected error", d.err.Error(), "but got", err)
		} else if err != nil && d.err == nil {
			t.Error(strconv.Itoa(i), "Expected error", d.err, "but got", err)
		} else if err != d.err && err.Error() != d.err.Error() {
			t.Error(strconv.Itoa(i), "Expected error", d.err, "but got", err)
		}

		if !eqBytes(buf.Bytes(), d.expected) {
			t.Error(strconv.Itoa(i), "Expected data", d.expected, "but got", buf.Bytes())
		}
	}
}

func TestWriteAndReadMux(t *testing.T) {
	for i, d := range writeDataMux {
		buf := &bytes.Buffer{}
		w := NewSlipMuxWriter(buf)
		err := w.WritePacket(FRAME_DIAGNOSTIC, d.data)

		if err != nil {
			t.Error("Unexpected error:", err)
		}

		r := NewSlipMuxReader(buf)
		p, frameType, err := r.ReadPacket()

		if err != nil {
			t.Error("Unexpected error:", err)
		}
		if frameType != FRAME_DIAGNOSTIC {
			t.Error("Expected frameType to be FRAME_DIAGNOSTIC but is", frameType)
		}
		if !eqBytes(p, d.data) {
			t.Error(strconv.Itoa(i), "Expected data", d.data, "but got", p)
		}
	}
}
