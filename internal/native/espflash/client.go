// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package espflash

import (
	"crypto/md5"
	"encoding/binary"
	"errors"
	"fmt"
	"os"

	"wa-lang.org/wa/internal/3rdparty/serial"
	"wa-lang.org/wa/internal/3rdparty/slip"
)

type Client struct {
	port *serial.Port
	r    *slip.Reader
	w    *slip.Writer
}

func Open(name string, baudrate int) (*Client, error) {
	p := new(Client)
	port, err := serial.OpenPort(
		&serial.Config{
			Name: name,
			Baud: baudrate,
		},
	)
	if err != nil {
		return nil, err
	}

	p.port = port
	p.r = slip.NewReader(port)
	p.w = slip.NewWriter(port)
	return p, nil
}

func (p *Client) Close() error {
	return p.port.Close()
}

func (c *Client) SendCommand(cmd CmdOpcode, payload []byte, checksum uint32) ([]byte, error) {
	if err := c.sendRequest(cmd, payload, checksum); err != nil {
		return nil, err
	}
	return c.readResponse()
}

func (c *Client) sendFrame(data []byte) error {
	return c.w.WritePacket(data)
}

func (c *Client) readFrame() (p []byte, isPrefix bool, err error) {
	return c.r.ReadPacket()
}

func (c *Client) sendRequest(cmd CmdOpcode, payload []byte, checksum uint32) error {
	// Assemble request frame
	frame := make([]byte, 0, 8+len(payload))
	frame = append(frame, 0x00) // direction: 0x00 = host -> target
	frame = append(frame, byte(cmd))

	// payload length
	lenb := make([]byte, 2)
	binary.LittleEndian.PutUint16(lenb, uint16(len(payload)))
	frame = append(frame, lenb...)

	// payload
	frame = append(frame, payload...)

	// checksum
	cksum := make([]byte, 4)
	binary.LittleEndian.PutUint32(cksum, checksum)
	frame = append(frame, cksum...)

	return c.sendFrame(frame)
}

// readResponse reads a response frame:
//
//	Status(1) + Length(2) + Payload
func (c *Client) readResponse() ([]byte, error) {
	resp, _, err := c.readFrame()
	if err != nil {
		return nil, err
	}
	if len(resp) < 3 {
		return nil, errors.New("invalid response length")
	}

	status := resp[0]
	if status != 0 {
		return nil, ErrorCode(status)
	}

	respLen := binary.LittleEndian.Uint16(resp[1:3])

	if int(respLen)+3 != len(resp) {
		return nil, errors.New("invalid response length field")
	}

	if status != 0 {
		return nil, fmt.Errorf("rom loader error: 0x%02x", status)
	}

	// ok
	return resp[3:], nil
}

func (f *Client) CmdSync() error {
	payload := make([]byte, 36)
	payload[0] = 0x07
	payload[1] = 0x07
	payload[2] = 0x12
	payload[3] = 0x20
	for i := 4; i < 36; i++ {
		payload[i] = 0x55
	}

	_, err := f.SendCommand(CmdSync, payload, 0)
	return err
}
func (f *Client) CmdEraseFlash() error {
	_, err := f.SendCommand(CmdEraseFlash, nil, 0)
	return err
}
func (f *Client) CmdBeginWrite(addr, totalSize, blockSize uint32) error {
	if blockSize == 0 {
		return errors.New("blockSize cannot be 0")
	}
	packets := (totalSize + blockSize - 1) / blockSize

	eraseSize := packets * blockSize

	payload := make([]byte, 5*4)
	binary.LittleEndian.PutUint32(payload[0:], eraseSize)
	binary.LittleEndian.PutUint32(payload[4:], packets)
	binary.LittleEndian.PutUint32(payload[8:], blockSize)
	binary.LittleEndian.PutUint32(payload[12:], addr)
	binary.LittleEndian.PutUint32(payload[16:], 0)

	_, err := f.SendCommand(CmdFlashBegin, payload, 0)
	return err
}
func (f *Client) CmdWriteBlock(blockIndex uint32, data []byte) error {
	payload := make([]byte, 16+len(data))
	binary.LittleEndian.PutUint32(payload[0:], uint32(len(data)))
	binary.LittleEndian.PutUint32(payload[4:], blockIndex)
	binary.LittleEndian.PutUint32(payload[8:], 0)
	binary.LittleEndian.PutUint32(payload[12:], 0)
	copy(payload[16:], data)

	// TODO: 计算 checksum, 按 ROM 协议
	var sum uint32
	for _, b := range data {
		sum += uint32(b)
	}

	_, err := f.SendCommand(CmdFlashData, payload, sum)
	return err
}
func (f *Client) CmdEndWrite(reboot bool) error {
	payload := make([]byte, 4)

	if reboot {
		binary.LittleEndian.PutUint32(payload, 0)
	} else {
		binary.LittleEndian.PutUint32(payload, 1)
	}

	_, err := f.SendCommand(CmdFlashEnd, payload, 0)
	return err
}

func (f *Client) CmdFlashBinary(addr uint32, bin []byte) error {
	const blockSize = 0x400 // 1024 bytes

	if err := f.CmdBeginWrite(addr, uint32(len(bin)), blockSize); err != nil {
		return err
	}

	var blockIndex uint32
	for offset := 0; offset < len(bin); offset += blockSize {
		end := offset + blockSize
		if end > len(bin) {
			end = len(bin)
		}
		chunk := bin[offset:end]

		if err := f.CmdWriteBlock(blockIndex, chunk); err != nil {
			return err
		}
		blockIndex++
	}

	if err := f.CmdEndWrite(false); err != nil {
		return err
	}

	return nil
}

func (c *Client) FlashFile(addr uint32, filePath string) error {
	bin, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	if err := c.CmdSync(); err != nil {
		fmt.Println("Sync failed, retrying...")
		if err := c.CmdSync(); err != nil {
			return fmt.Errorf("sync failed: %v", err)
		}
	}

	if err := c.CmdFlashBinary(addr, bin); err != nil {
		return fmt.Errorf("flash failed: %v", err)
	}

	ok, err := c.verifyFlashMD5(addr, bin)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("flash verify failed: MD5 mismatch")
	}

	return nil
}

func (c *Client) verifyFlashMD5(addr uint32, bin []byte) (bool, error) {

	const CmdSPIFlashMD5 = 0x13

	payload := make([]byte, 16)
	binary.LittleEndian.PutUint32(payload[0:], addr)
	binary.LittleEndian.PutUint32(payload[4:], uint32(len(bin)))
	binary.LittleEndian.PutUint32(payload[8:], 0)
	binary.LittleEndian.PutUint32(payload[12:], 0)

	resp, err := c.SendCommand(CmdSPIFlashMD5, payload, 0)
	if err != nil {
		return false, fmt.Errorf("SPI_FLASH_MD5 command failed: %v", err)
	}

	if len(resp) < 16 {
		return false, errors.New("MD5 response too short")
	}

	deviceMD5 := resp[:16]
	localMD5 := md5.Sum(bin)

	return string(deviceMD5) == string(localMD5[:]), nil
}
