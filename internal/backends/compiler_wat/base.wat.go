// 版权 @2022 凹语言 作者。保留所有权利。

package compiler_wat

import (
	_ "embed"
)

//go:embed base.wat
var __base_wat_data string

// wasm 内存 1 页大小
const _WASM_PAGE_SIZE = 65536

// WASM 约定栈和内存管理
// 相关全局变量地址必须和 base_wasm.go 保持一致
const (
	// ;; heap 和 stack 状态(__heap_base 只读)
	// ;; | 0 <-- stack --> | <-- static-data --> | <-- heap --> |
	// (global $$stack_prt (mut i32) (i32.const 1024)) ;; index=0
	// (global $$heap_base i32 (i32.const 2048))       ;; index=1
	__stack_ptr_index = 0
	__heap_base_index = 1
)

// 内置函数名字
const (
	// 栈函数
	_waStackPtr   = "$$StackPtr"
	_waStackAlloc = "$$StackAlloc"

	// 堆管理函数
	_waHeapPtr   = "$$HeapPtr"
	_waHeapAlloc = "$$HeapAlloc"
	_waHeapFree  = "$$HeapFree"

	// 输出函数
	_waPrintString = "$$waPuts"
	_waPrintRune   = "$$waPrintRune"
	_waPrintInt32  = "$$waPrintI32"

	// 开始函数
	_waStart = "_start"
)

// todo: PrintI32 => waPrintI32, 以 wa 开头

const modBaseWat_wa = `
(memory $memory 1)

(export "memory" (memory $memory))

;;Remove these functions if they've been implemented in .wa
(global $$heap_ptr (mut i32) (i32.const 2048))

(func $$waGetHeapPtr(result i32)
	global.get $$heap_ptr
)

(func $$waLoadI32(param $ptr i32) (result i32)
	local.get $ptr
	i32.load offset=0 align=1
)
(func $$waStoreI32(param $ptr i32) (param $value i32)
	local.get $ptr
	i32.const 1
	i32.store offset=0 align=1
)


(func $$waHeapAlloc (param $size i32) (result i32) ;;result = ptr
	;;Todo
	global.get $$heap_ptr

	;; global.get $$heap_ptr
	;; call $$runtime.waPrintI32
	;; i32.const 32
	;; call $$runtime.waPrintRune
	;; local.get $size
	;; call $$runtime.waPrintI32
	;; i32.const 10
	;; call $$runtime.waPrintRune

	global.get $$heap_ptr
	local.get $size
	i32.add
	global.set $$heap_ptr
)
(func $$waHeapFree (param $ptr i32)
	;;Todo
	;; i32.const 126
	;; call $$runtime.waPrintRune
	;; local.get $ptr
	;; call $$runtime.waPrintI32
	;; i32.const 10
	;; call $$runtime.waPrintRune
)

(func $$wa.RT.Block.Init (param $ptr i32) (param $item_count i32) (param $release_func i32) (param $item_size i32) (result i32) ;;result = ptr
  local.get $ptr

  local.get $ptr
  if
    local.get $ptr
    i32.const 1
    i32.store offset=0 align=1

    local.get $ptr
    local.get $item_count
    i32.store offset=4 align=1

    local.get $ptr
    local.get $release_func
    i32.store offset=8 align=1

	local.get $ptr
	local.get $item_size
    i32.store offset=12 align=1
  end
)

(func $$wa.RT.DupWatStack (param $p i32) (result i32) (result i32) ;;result0 = result1 = p
  local.get $p
  local.get $p
)

(func $$wa.RT.Block.Retain (param $ptr i32) (result i32) ;;result = ptr
  local.get $ptr

  local.get $ptr
  if
    local.get $ptr
    local.get $ptr
    i32.load offset=0 align=1
    i32.const 1
    i32.add
    i32.store offset=0 align=1
  end
)

(func $$wa.RT.Block.Release (param $ptr i32)
  ;;Todo
  (local $ref_count i32)
  (local $item_count i32)
  (local $free_func i32)
  (local $item_size i32)
  (local $data_ptr i32)

  local.get $ptr
  i32.const 0
  i32.eq
  if
    return
  end

  local.get $ptr
  i32.load offset=0 align=1
  i32.const 1
  i32.sub
  local.set $ref_count

  local.get $ref_count
  if
    local.get $ptr
    local.get $ref_count
	i32.store offset=0 align=1

  else  ;;ref_count == 0
    local.get $ptr
	i32.load offset=8 align=1
	local.set $free_func

	local.get $free_func
	if  ;;free_func != 0
	  local.get $ptr
	  i32.load offset=4 align=1
	  local.set $item_count

	  local.get $item_count
	  if  ;;item_count > 0
        local.get $ptr
        i32.load offset=12 align=1
        local.set $item_size

        local.get $ptr
        i32.const 16
        i32.add
        local.set $data_ptr

        loop $free_next
          ;; onFree(data_ptr)
          local.get $data_ptr
          local.get $free_func
          call_indirect (type $$onFree)

          ;; item_count--
          local.get $item_count
          i32.const 1
          i32.sub
          local.set $item_count

          local.get $item_count
          if  ;;while item_count>0
            ;; data_ptr += item_size
            local.get $data_ptr
            local.get $item_size
            i32.add
            local.set $data_ptr

            br $free_next  ;;continue
          end  ;;while item_count>0        
        end  ;;loop $free_next
	  end  ;;if item_count > 0
	end  ;;free_func != 0

	local.get $ptr
	call $$waHeapFree
  end  ;;ref_count == 0
)

`

const modBaseWat_wasi = `
(import "wasi_snapshot_preview1" "fd_write"
	(func $$FdWrite (param i32 i32 i32 i32) (result i32))
)

(memory $memory 1)

(export "memory" (memory $memory))
(export "_start" (func $_start))
;; (export "main.main" (func $main.main))

;; | 0 <-- stack --> | <-- static-data --> | <-- heap --> |
(global $$stack_prt (mut i32) (i32.const 1024)) ;; index=0
(global $$heap_base i32 (i32.const 2048))       ;; index=1

(func $$StackPtr (result i32)
	(global.get $$stack_prt)
)

(func $$StackAlloc (param $size i32) (result i32)
	;; $$stack_prt -= $size
	(global.set $$stack_prt (i32.sub (global.get $$stack_prt) (local.get  $size)))
	;; return $$stack_prt
	(return (global.get $$stack_prt))
)

(func $$HeapPtr (result i32)
	(global.get $$heap_base)
)

(func $$HeapAlloc (param $size i32) (result i32)
	;; {{$$HeapAlloc/body/begin}}
	unreachable
	;; {{$$HeapAlloc/body/end}}
)

(func $$HeapFree (param $ptr i32)
	;; {{$$HeapFree/body/begin}}
	unreachable
	;; {{$$HeapFree/body/end}}
)

(func $$Puts (param $str i32) (param $len i32)
	;; {{$$Puts/body/begin}}

	(local $sp i32)
	(local $p_iov i32)
	(local $p_nwritten i32)
	(local $stdout i32)

	(local.set $sp (global.get $$stack_prt))

	(local.set $p_iov (call $$StackAlloc (i32.const 8)))

	(local.set $p_nwritten (call $$StackAlloc (i32.const 4)))

	(i32.store offset=0 align=1 (local.get $p_iov) (local.get $str))
	(i32.store offset=4 align=1 (local.get $p_iov) (local.get $len))

	(local.set $stdout (i32.const 1))

	(call $$FdWrite
		(local.get $stdout)
		(local.get $p_iov) (i32.const 1)
		(local.get $p_nwritten)
	)

	(global.set $$stack_prt (local.get $sp))
	drop

	;; {{$$Puts/body/end}}
)

(func $$waPrintRune (param $ch i32)
	;; {{$$waPrintRune/body/begin}}

	(local $sp i32)
	(local $p_ch i32)

	(local.set $sp (global.get $$stack_prt))

	(local.set $p_ch (call $$StackAlloc (i32.const 4)))
	(i32.store offset=0 align=1 (local.get $p_ch) (local.get $ch))

	(call $$Puts (local.get $p_ch) (i32.const 1))

	(global.set $$stack_prt (local.get $sp))

	;; {{$$waPrintRune/body/begin}}
)

(func $$waPrintI32 (param $x i32)
	;; if $x == 0 { print '0'; return }
	(i32.eq (local.get $x) (i32.const 0))
	if
		(call $$waPrintRune (i32.const 48)) ;; '0'
		(return)
	end

	;; if $x < 0 { $x = 0-$x; print '-'; }
	(i32.lt_s (local.get $x) (i32.const 0))
	if 
		(local.set $x (i32.sub (i32.const 0) (local.get $x)))
		(call $$waPrintRune (i32.const 45)) ;; '-'
	end

	local.get $x
	call $$$print_i32
)

(func $$$print_i32 (param $x i32)
	;; {{$$$print_i32/body/begin}}

	(local $div i32)
	(local $rem i32)

	;; if $x == 0 { print '0'; return }
	(i32.eq (local.get $x) (i32.const 0))
	if
		(return)
	end

	;; print_i32($x / 10)
	;; puchar($x%10 + '0')
	(local.set $div (i32.div_s (local.get $x) (i32.const 10)))
	(local.set $rem (i32.rem_s (local.get $x) (i32.const 10)))
	(call $$$print_i32 (local.get $div))
	(call $$waPrintRune (i32.add (local.get $rem) (i32.const 48))) ;; '0'

	;; {{$$$print_i32/body/end}}
)

(func $_start
	;; {{$_start/body/begin}}
	;; (call $main.init)
	(call $main)
	;; {{$_start/body/end}}
)
`
