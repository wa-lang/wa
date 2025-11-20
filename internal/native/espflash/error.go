// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package espflash

import "fmt"

type ErrorCode uint8

const (
	ErrROMUndefined         ErrorCode = 0x00
	ErrROMInvalidParam      ErrorCode = 0x01
	ErrROMMallocFailed      ErrorCode = 0x02
	ErrROMSendFailed        ErrorCode = 0x03
	ErrROMRecvFailed        ErrorCode = 0x04
	ErrROMInvalidFormat     ErrorCode = 0x05
	ErrROMResultError       ErrorCode = 0x06
	ErrROMChecksumError     ErrorCode = 0x07
	ErrROMFlashWriteError   ErrorCode = 0x08
	ErrROMFlashReadError    ErrorCode = 0x09
	ErrROMFlashReadLenError ErrorCode = 0x0a
	ErrROMDeflateFailed     ErrorCode = 0x0b
	ErrROMDeflateAdlerError ErrorCode = 0x0c
	ErrROMDeflateParamError ErrorCode = 0x0d
	ErrROMInvalidRAMSize    ErrorCode = 0x0e
	ErrROMInvalidRAMAddr    ErrorCode = 0x0f

	// Extended errors (0x64+)
	ErrROMInvalidParameter    ErrorCode = 0x64
	ErrROMInvalidFmt          ErrorCode = 0x65
	ErrROMDescriptionTooLong  ErrorCode = 0x66
	ErrROMBadEncoding         ErrorCode = 0x67
	ErrROMInsufficientStorage ErrorCode = 0x69
)

// description map
var romErrorDescriptions = map[ErrorCode]string{
	ErrROMUndefined:         "Undefined errors",
	ErrROMInvalidParam:      "The input parameter is invalid",
	ErrROMMallocFailed:      "Failed to malloc memory from system",
	ErrROMSendFailed:        "Failed to send out message",
	ErrROMRecvFailed:        "Failed to receive message",
	ErrROMInvalidFormat:     "The format of the received message is invalid",
	ErrROMResultError:       "Message is ok, but the running result is wrong",
	ErrROMChecksumError:     "Checksum error",
	ErrROMFlashWriteError:   "Flash write error",
	ErrROMFlashReadError:    "Flash read error",
	ErrROMFlashReadLenError: "Flash read length error",
	ErrROMDeflateFailed:     "Deflate failed error",
	ErrROMDeflateAdlerError: "Deflate Adler32 error",
	ErrROMDeflateParamError: "Deflate parameter error",
	ErrROMInvalidRAMSize:    "Invalid RAM binary size",
	ErrROMInvalidRAMAddr:    "Invalid RAM binary address",

	ErrROMInvalidParameter:    "Invalid parameter",
	ErrROMInvalidFmt:          "Invalid format",
	ErrROMDescriptionTooLong:  "Description too long",
	ErrROMBadEncoding:         "Bad encoding description",
	ErrROMInsufficientStorage: "Insufficient storage",
}

func (e ErrorCode) String() string {
	if s, ok := romErrorDescriptions[e]; ok {
		return s
	}
	return fmt.Sprintf("Unknown ROM error 0x%02X", uint8(e))
}

func (e ErrorCode) Error() string {
	return e.String()
}
