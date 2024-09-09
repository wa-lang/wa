;; Copyright 2024 The Wa Authors. All rights reserved.

(global $PALETTE0 i32 (i32.const 0x04))
(global $PALETTE1 i32 (i32.const 0x08))
(global $PALETTE2 i32 (i32.const 0x0c))
(global $PALETTE3 i32 (i32.const 0x10))
(global $DRAW_COLORS i32 (i32.const 0x14))
(global $GAMEPAD1 i32 (i32.const 0x16))
(global $GAMEPAD2 i32 (i32.const 0x17))
(global $GAMEPAD3 i32 (i32.const 0x18))
(global $GAMEPAD4 i32 (i32.const 0x19))
(global $MOUSE_X i32 (i32.const 0x1a))
(global $MOUSE_Y i32 (i32.const 0x1c))
(global $MOUSE_BUTTONS i32 (i32.const 0x1e))
(global $SYSTEM_FLAGS i32 (i32.const 0x1f))
(global $NETPLAY i32 (i32.const 0x20))

(global $FRAMEBUFFER_BLOCK i32 (i32.const 0x98))
(global $FRAMEBUFFER i32 (i32.const 0xa0))

(global $BUTTON_1 i32 (i32.const 1))
(global $BUTTON_2 i32 (i32.const 2))
(global $BUTTON_LEFT i32 (i32.const 16))
(global $BUTTON_RIGHT i32 (i32.const 32))
(global $BUTTON_UP i32 (i32.const 64))
(global $BUTTON_DOWN i32 (i32.const 128))

(global $MOUSE_LEFT i32 (i32.const 1))
(global $MOUSE_RIGHT i32 (i32.const 2))
(global $MOUSE_MIDDLE i32 (i32.const 4))

(global $SYSTEM_PRESERVE_FRAMEBUFFER i32 (i32.const 1))
(global $SYSTEM_HIDE_GAMEPAD_OVERLAY i32 (i32.const 2))

(global $BLIT_2BPP i32 (i32.const 1))
(global $BLIT_1BPP i32 (i32.const 0))
(global $BLIT_FLIP_X i32 (i32.const 2))
(global $BLIT_FLIP_Y i32 (i32.const 4))
(global $BLIT_ROTATE i32 (i32.const 8))

(global $TONE_PULSE1 i32 (i32.const 0))
(global $TONE_PULSE2 i32 (i32.const 1))
(global $TONE_TRIANGLE i32 (i32.const 2))
(global $TONE_NOISE i32 (i32.const 3))
(global $TONE_MODE1 i32 (i32.const 0))
(global $TONE_MODE2 i32 (i32.const 4))
(global $TONE_MODE3 i32 (i32.const 8))
(global $TONE_MODE4 i32 (i32.const 12))
(global $TONE_PAN_LEFT i32 (i32.const 16))
(global $TONE_PAN_RIGHT i32 (i32.const 32))
(global $TONE_NOTE_MODE i32 (i32.const 64))

(func $$syscall/wasm4.getMemory (param $blk i32) (param $ptr i32) (param $len i32) (param $cap i32) (result i32 i32 i32 i32)
	local.get $blk
	local.get $ptr
	local.get $len
	local.get $cap
	return
)

(func $$syscall/wasm4.getFramebuffer (result i32 i32 i32 i32)
	global.get $FRAMEBUFFER_BLOCK
	global.get $FRAMEBUFFER
	i32.const 6400
	i32.const 6400
	return
)

(func $$syscall/wasm4.getPalette0 (result i32)
	global.get $PALETTE0
	i32.load
)
(func $$syscall/wasm4.getPalette1 (result i32)
	global.get $PALETTE1
	i32.load
)

(func $$syscall/wasm4.getPalette2 (result i32)
	global.get $PALETTE2
	i32.load
)

(func $$syscall/wasm4.getPalette3 (result i32)
	global.get $PALETTE3
	i32.load
)

(func $$syscall/wasm4.GetPalette (result i32 i32 i32 i32)
	global.get $PALETTE0
	i32.load
	global.get $PALETTE1
	i32.load
	global.get $PALETTE2
	i32.load
	global.get $PALETTE3
	i32.load
	return
)

(func $$syscall/wasm4.setPalette0 (param $a i32)
	global.get $PALETTE0
	local.get $a
	i32.store
)
(func $$syscall/wasm4.setPalette1 (param $a i32)
	global.get $PALETTE1
	local.get $a
	i32.store
)
(func $$syscall/wasm4.setPalette2 (param $a i32)
	global.get $PALETTE2
	local.get $a
	i32.store
)
(func $$syscall/wasm4.setPalette3 (param $a i32)
	global.get $PALETTE3
	local.get $a
	i32.store
)

(func $$syscall/wasm4.SetPalette (param $a0 i32) (param $a1 i32) (param $a2 i32) (param $a3 i32)
	global.get $PALETTE0
	local.get $a0
	i32.store

	global.get $PALETTE1
	local.get $a1
	i32.store

	global.get $PALETTE2
	local.get $a2
	i32.store

	global.get $PALETTE3
	local.get $a3
	i32.store

	return
)

(func $$syscall/wasm4.GetDrawColors (result i32)
	global.get $DRAW_COLORS
	i32.load16_u
	return
)

(func $$syscall/wasm4.SetDrawColors (param $a i32)
	global.get $DRAW_COLORS
	local.get $a
	i32.store16
	return
)

(func $$syscall/wasm4.GetGamePad1 (result i32)
	global.get $GAMEPAD1
	i32.load8_u
	return
)
(func $$syscall/wasm4.GetGamePad2 (result i32)
	global.get $GAMEPAD2
	i32.load8_u
	return
)
(func $$syscall/wasm4.GetGamePad3 (result i32)
	global.get $GAMEPAD3
	i32.load8_u
	return
)
(func $$syscall/wasm4.GetGamePad4 (result i32)
	global.get $GAMEPAD4
	i32.load8_u
	return
)

(func $$syscall/wasm4.GetMouseX (result i32)
	global.get $MOUSE_X
	i32.load16_u
	return
)
(func $$syscall/wasm4.GetMouseY (result i32)
	global.get $MOUSE_Y
	i32.load16_u
	return
)
(func $$syscall/wasm4.GetMouseButtons (result i32)
	global.get $MOUSE_BUTTONS
	i32.load8_u
	return
)
