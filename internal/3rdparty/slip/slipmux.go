package slip

import (
	"bytes"
	"io"
)

type SlipMuxReader struct {
	r *Reader
}

func NewSlipMuxReader(reader io.Reader) *SlipMuxReader {
	return &SlipMuxReader{
		r: NewReader(reader),
	}
}

type SlipMuxWriter struct {
	w *Writer
}

func NewSlipMuxWriter(writer io.Writer) *SlipMuxWriter {
	return &SlipMuxWriter{
		w: NewWriter(writer),
	}
}

const (
	FRAME_IPV4_START = 0x45
	FRAME_IPV4_END   = 0x4f
	FRAME_IPV6_START = 0x60
	FRAME_IPV6_END   = 0x6f
)

const (
	FRAME_UNKNOWN    = 0x00
	FRAME_DIAGNOSTIC = 0x0a // New Line Character
	FRAME_COAP       = 0xa9
)

func IsIpFrame(frame byte) bool {
	return IsIpv4Frame(frame) || IsIpv6Frame(frame)
}

func IsIpv6Frame(frame byte) bool {
	return frame >= FRAME_IPV6_START && frame <= FRAME_IPV6_END
}

func IsIpv4Frame(frame byte) bool {
	return frame >= FRAME_IPV4_START && frame <= FRAME_IPV4_END
}

// WritePacket writes a SlipMux packet with the given Frame prefix
// For registered frames see FRAME_* constants. IPV4 and IPV6 frames
// are not prepended as specified.
// To not use any Frame identifiers use SlipWriter directly.
func (s *SlipMuxWriter) WritePacket(frame byte, p []byte) error {
	if !IsIpFrame(frame) {
		p = append([]byte{frame}, p...)
	}

	if frame == FRAME_COAP {
		p = AppendFcs16(p, CalcFcs16(p))
	}

	return s.w.WritePacket(p)
}

func isInvalidFrame(frameType byte) bool {
	// ESC_ESC and ESC_END are not reserved
	return frameType == END ||
		frameType == ESC ||
		frameType == FRAME_UNKNOWN
}

// ReadPacket reads a packet from the stream.
// Blocks till packet is complete.
// Returns the packet, the frameType byte and an optional error
//
// The he packet frame is stripped from the returning byte slice.
// IPv4 and IPv6 frame identifiers are not stripped to keep
// backwards compatibility with SLIP
func (s *SlipMuxReader) ReadPacket() ([]byte, byte, error) {
	buf := bytes.Buffer{}

	for {
		p, isPrefix, err := s.r.ReadPacket()
		if err != nil && err != io.EOF {
			// EOF does not return here and must be handled
			// via Timeout in the application this is because
			// some streams might return EOF even if there
			// will be more data in future
			return nil, 0, err
		}
		buf.Write(p)
		if !isPrefix && len(p) > 0 {
			break
		}
	}

	res := buf.Bytes()

	frameType := res[0]

	// Ignore packets with invalid frame types
	if isInvalidFrame(frameType) {
		return s.ReadPacket()
	}

	if frameType == FRAME_COAP {
		// smallest CoAP message is frameType + 4 byte + 2 byte CRC = 7 bytes
		if len(res) < 7 {
			return s.ReadPacket()
		}

		// Ignore packets with bad Checksum
		if !CheckFsc16(res) {
			return s.ReadPacket()
		} else {
			res = RemoveFcs16(res)
		}
	}

	// Strip first byte of non IP packets
	if !IsIpFrame(frameType) {
		res = res[1:]
	}

	return res, frameType, nil

}
