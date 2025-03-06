(module $__walang__
  (import "syscall_js" "print_bool" (func $syscall$js.__import__print_bool (param i32)))
  (import "syscall_js" "print_f32" (func $syscall$js.__import__print_f32 (param f32)))
  (import "syscall_js" "print_f64" (func $syscall$js.__import__print_f64 (param f64)))
  (import "syscall_js" "print_i32" (func $syscall$js.__import__print_i32 (param i32)))
  (import "syscall_js" "print_i64" (func $syscall$js.__import__print_i64 (param i64)))
  (import "syscall_js" "print_ptr" (func $syscall$js.__import__print_ptr (param i32)))
  (import "syscall_js" "print_rune" (func $syscall$js.__import__print_rune (param i32)))
  (import "syscall_js" "print_str" (func $syscall$js.__import__print_str (param i32) (param i32)))
  (import "syscall_js" "print_u32" (func $syscall$js.__import__print_u32 (param i32)))
  (import "syscall_js" "print_u64" (func $syscall$js.__import__print_u64 (param i64)))
  (import "syscall_js" "proc_exit" (func $syscall$js.__import__proc_exit (param i32)))
;; Copyright 2023 The Wa Authors. All rights reserved.

(memory $memory 1024)

(export "memory" (memory $memory))

;; +-----------------+---------------------+--------------+
;; | 0 <-- stack --> | <-- static-data --> | <-- heap --> |
;; +-----------------+---------------------+--------------+

(global $__stack_ptr (mut i32) (i32.const 1024))     ;; index=0
;;(global $__heap_base i32 (i32.const 1048576))     ;; index=1

(global $__heap_lfixed_cap i32 (i32.const 64)) ;; 固定尺寸空闲链表最大长度, 满时回收; 也可用于关闭 fixed 策略

;; ---------------------------------------------------------
;; package: runtime
;; ---------------------------------------------------------

;; file: heap.wat.ws

;; Copyright 2023 The Wa Authors. All rights reserved.

(func $runtime.throw
	unreachable
)

(func $runtime.getStackPtr (result i32)
	global.get $__stack_ptr
)

(func $runtime.setStackPtr (param $sp i32)
	local.get $sp
	global.set $__stack_ptr
)

(func $runtime.stackAlloc (param $size i32) (result i32)
	;; $__stack_ptr -= $size
	global.get $__stack_ptr
	local.get  $size
	i32.sub 
	global.set $__stack_ptr 

	;; return $__stack_ptr
	global.get $__stack_ptr
	return
)

(func $runtime.stackFree (param $size i32)
	;; $__stack_ptr += $size
	global.get $__stack_ptr
	local.get $size
	i32.add
	global.set $__stack_ptr 
)

(func $runtime.heapBase(result i32)
	global.get $__heap_base
)

(func $runtime.HeapAlloc (export "runtime.HeapAlloc") (param $nbytes i32) (result i32) ;;result = ptr
	(local $ptr i32)

	local.get $nbytes
	i32.eqz
	if
		i32.const 0
		return
	end

	local.get $nbytes
	i32.const 7
	i32.add
	i32.const 8
	i32.div_u
	i32.const 8
	i32.mul
	local.set $nbytes

	local.get $nbytes
	call $runtime.malloc
	local.set $ptr

	loop $zero
		local.get $nbytes
		i32.const 8
		i32.sub
		local.tee $nbytes
		local.get $ptr
		i32.add

		i64.const 0
		i64.store

		local.get $nbytes
		if
			br $zero
		end
	end ;;loop $zero

	local.get $ptr
)

(func $runtime.HeapFree (export "runtime.HeapFree") (param $ptr i32)
	local.get $ptr
	call $runtime.free
)

(func $runtime.Block.Init (param $ptr i32) (param $item_count i32) (param $release_func i32) (param $item_size i32) (result i32) ;;result = ptr
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

(func $runtime.Block.SetFinalizer (param $ptr i32) (param $release_func i32)
	local.get $ptr
	if
		local.get $ptr
		local.get $release_func
		i32.store offset=8 align=1
	end
)

(func $runtime.Block.HeapAlloc (export "runtime.Block.HeapAlloc") (param $item_count i32) (param $release_func i32) (param $item_size i32) (result i32) ;;result = ptr_block
  local.get $item_count
  local.get $item_size
  i32.mul
  i32.const 16
  i32.add
  call $runtime.HeapAlloc

  local.get $item_count
  local.get $release_func
  local.get $item_size
  call $runtime.Block.Init
)

(func $runtime.DupI32 (param $a i32) (result i32 i32) ;;r0 = r1 = p0
	local.get $a
	local.get $a
)

(func $runtime.SwapI32 (param $a i32) (param $b i32) (result i32 i32) ;;r0 = p1, r1 = p0
	local.get $b
	local.get $a
)

(func $runtime.Block.Retain (export "runtime.Block.Retain") (param $ptr i32) (result i32) ;;result = ptr
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

(func $runtime.Block.Release (export "runtime.Block.Release") (param $ptr i32)
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
					;; OnFree(data_ptr)
					local.get $data_ptr
					local.get $free_func
					call_indirect (type $$OnFree)

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
		call $runtime.HeapFree
	end  ;;ref_count == 0
)

(func $$wa.runtime.i32_ref_to_ptr (param $b i32) (param $d i32) (result i32) ;;result = ptr
	local.get $d
)

(func $$wa.runtime.i64_ref_to_ptr (param $b i32) (param $d i32) (result i32) ;;result = ptr
	local.get $d
)

(func $$wa.runtime.slice_to_ptr (param $b i32) (param $d i32) (param $l i32) (param $c i32) (result i32) ;;result = ptr
	local.get $d
)
;; file: heap_malloc.wat.ws

;; Copyright 2025 The Wa Authors. All rights reserved.

	;; 栈/静态数据/堆的内存布局
	;;
	;; 0               __stack_ptr           __heap_base    __heap_ptr  __heap_top
	;; |                 |                     |              |           |
	;; v                 v                     v              v           v
	;; +-----------------+---------------------+--------------+-----------+...........+
	;; | 0 <<<< stack ###| <-- static-data --> |### heap >>>> |...........|grow >>>>> |
	;; +-----------------+---------------------+--------------+-----------+...........+
	;;                                         |              |
	;; +---------------------------------------+              +------------------+
	;; |                                                                         |
	;; v                                                                         v
	;; +------------------------------+------------------------------------------+
	;; |  header: [5]heap_block_t     |### node:*heap_block_t......>>>....>>>>...|
	;; +-+----------------------------+----------+----------+--------------------+
	;; |                              |          ^          ^                    ^
	;; v                              v          |          |                    |
	;; +-----+-----+-----+-----+------+-----+    |          |                 __heap_ptr
	;; | l24 | l32 | l48 | l80 | l128 | nil | ---|----------|---> 空闲链表头, 头节点没有数据
	;; +--+--+--+--+--+--+--+--+--+---+-----+    |          |
	;; ^  |     |     |     |     |       +---------+     +---------+
	;; |  +-----|-----|-----|-----|-----> | 15 Byte | --> | 24 Byte | --> nil
	;; |        |     |     |     |       +---------+     +---------+
	;; |        |     |     |     |       +---------+
	;; |        +-----|-----|-----|-----> | 32 Byte | --> nil
	;; |              |     |     |       +---------+
	;; |              |     |     |       +---------+     +---------+
	;; |              +-----|-----|-----> | 40 Byte | --> | 48 Byte | --> nil
	;; |                    |     |       +---------+     +---------+
	;; |                    +-----|-----> nil
	;; |                          |       +---------+     +-----------+
	;; __heap_base          +-----+-----> | 81 Byte | --> | 1024 Byte | --> Head
	;;                      |     |       +---------+     +-----------+
	;;                      |     +-----> l128 负责大于80字节大小的变长内存分配
	;; __heap_l128_freep ---+     +-----> l128 是 **循环链表** 需要记录上次检索的位置
	;;                            +-----> l128 后面的 nil 避免头和分配的块数据相邻
	;;
	;; heap_fixed_header_t
	;; +----------+----------+
	;; | len:i32  | next:i32 |
	;; +----------+----------+
	;;
	;; heap_L128_header_t
	;; +------------+----------+
	;; | size=0:i32 | next:i32 |
	;; +------------+----------+
	;;
	;; heap_block_t
	;; +----------+----------+---------+
	;; | size:i32 | next:i32 | data[0] |
	;; +----------+----------+---------+
	;;
	(global $__heap_ptr  (mut i32) (i32.const 0))         ;; heap 当前位置指针
	(global $__heap_top  (mut i32) (i32.const 0))         ;; heap 最大位置指针(超过时要 grow 内存)
	(global $__heap_l128_freep (mut i32) (i32.const 0))   ;; l128 是循环链表, 记录当前的迭代位置
	(global $__heap_init_flag (mut i32) (i32.const 0)) ;; 是否已经初始化

	;; 合法的指针
	(func $heap_assert_valid_ptr (param $ptr i32)
		;; ptr > 0
		local.get $ptr
		i32.const 0
		i32.gt_s
		if else unreachable end

		;; 4 的倍数
		local.get $ptr
		i32.const 4
		i32.rem_s
		i32.eqz
		if else unreachable end
	)

	;; 是否允许固定大小的分配策略
	(func $heap_is_fixed_list_enabled (result i32)
		global.get $__heap_lfixed_cap
		i32.eqz
		if (result i32)
			i32.const 0
		else
			i32.const 1
		end
	)
	(func $heap_assert_fixed_list_enabled
		global.get $__heap_lfixed_cap
		i32.eqz
		if unreachable end
	)

	;; 判断是否为固定大小内存
	(func $heap_is_fixed_size (param $size i32) (result i32)
		;; 禁止了 fixed 策略?
		call $heap_is_fixed_list_enabled
		if else
			i32.const 0
			return
		end

		;; return ($size <= 80)? 1: 0
		local.get $size
		i32.const 80
		i32.le_s
		if (result i32)
			i32.const 1
		else
			i32.const 0
		end
	)

	;; size 对齐到8字节
	;; 块节点头部也是8字节对齐
	(func $heap_alignment8 (param $size i32) (result i32)
		;; (($size+7)/8)*8
		local.get $size
		i32.const 7
		i32.add
		i32.const 8
		i32.div_s
		i32.const 8
		i32.mul
	)
	(func $heap_assert_align8 (param $size i32)
		local.get $size
		i32.const 8
		i32.rem_s
		i32.eqz
		if else
			unreachable
		end
	)

	;; 初始化 heap_block_t
	(func $heap_block.init (param $ptr i32) (param $size i32) (param $next i32)
		;; $ptr[offset=0] = $size
		local.get $ptr
		local.get $size
		i32.store offset=0
		;; $ptr[offset=4] = $prev
		local.get $ptr
		local.get $next
		i32.store offset=4
	)
	(func $heap_block.set_size (param $ptr i32) (param $size i32)
		local.get $ptr
		local.get $size
		i32.store offset=0
	)
	(func $heap_block.set_next (param $ptr i32) (param $next i32)
		local.get $ptr
		local.get $next
		i32.store offset=4
	)

	;; 读取 heap_block_t 结构体成员
	(func $heap_block.size (param $ptr i32) (result i32)
		local.get $ptr
		i32.load offset=0
	)
	(func $heap_block.next (param $ptr i32) (result i32)
		local.get $ptr
		i32.load offset=4
	)
	(func $heap_block.data (param $ptr i32) (result i32)
		local.get $ptr
		i32.const 8
		i32.add
	)
	(func $heap_block.end (param $ptr i32) (result i32)
		;; $ptr.size + $ptr + 8
		local.get $ptr
		i32.load offset=0
		local.get $ptr
		i32.add
		i32.const 8
		i32.add
	)

	;; 根据要申请内存的大小返回合适的空闲链表, 和对齐后的大小
	;; func $heap_free_list.ptr_and_fixed_size(size: i32) => (ptr, fixed_size: i32)
	(func $heap_free_list.ptr_and_fixed_size (param $size i32) (result i32 i32)
		;; l24 | l32 | l48 | l80 | l128

		;; 禁止了 fixed 策略?
		call $heap_is_fixed_list_enabled
		if else
			;; return l128
			global.get $__heap_base
			i32.const 32
			i32.add
			local.get $size
			call $heap_alignment8
			return
		end

		;; if $size > 80: return (l128, max(128, $fixed_size))
		local.get $size
		i32.const 80
		i32.gt_s
		if
			;; $ptr = $__heap_base + 4*sizeof(heap_block_t)
			global.get $__heap_base
			i32.const 32
			i32.add
			;; $fixed_size = max($size, 128)
			block (result i32)
				local.get $size
				i32.const 128
				i32.le_s
				if (result i32)
					i32.const 128
				else
					local.get $size
					call $heap_alignment8
				end
			end
			;; return ($ptr, $fixed_size)
			return
		end

		;; if $size > 48: return (l80, 80)
		local.get $size
		i32.const 48
		i32.gt_s
		if
			global.get $__heap_base
			i32.const 24 ;; 3*sizeof(heap_block_t)
			i32.add
			i32.const 80
			return
		end

		;; if $size > 32: return (l48, 48)
		local.get $size
		i32.const 32
		i32.gt_s
		if
			global.get $__heap_base
			i32.const 16 ;; 2*sizeof(heap_block_t)
			i32.add
			i32.const 48
			return
		end

		;; if $size > 24: return (l32, 32)
		local.get $size
		i32.const 24
		i32.gt_s
		if
			global.get $__heap_base
			i32.const 8 ;; 1*sizeof(heap_block_t)
			i32.add
			i32.const 32
			return
		end

		;; return (l24, 24)
		global.get $__heap_base
		i32.const 24
		return
	)

	;; func wa_malloc_init_once()
	;; 基于 $__heap_base 全局常量值初始化堆
	;; 堆当前的大小由 $__heap_top 全局变量记录
	;; 当堆空间不足时, 将进行内存扩容操作, 扩容可能受最大页限制失败
	(func $wa_malloc_init_once
		;; 已经完成初始化则退出
		global.get $__heap_init_flag
		if return end

		;; 设置初始化标志
		i32.const 1
		global.set $__heap_init_flag

		;; assert $__stack_ptr > 0
		global.get $__stack_ptr
		i32.const 0
		i32.gt_s
		if else
			unreachable
		end

		;; assert $__stack_ptr < $__heap_base
		global.get $__stack_ptr
		global.get $__heap_base
		i32.lt_s
		if else
			unreachable
		end

		;; assert $__heap_base % 8 == 0
		global.get $__heap_base
		call $heap_assert_align8

		;; $__heap_ptr = $__heap_base + 6*sizeof(heap_block_t)
		global.get $__heap_base
		i32.const 48 ;; 6*sizeof(heap_block_t)
		i32.add
		global.set $__heap_ptr

		;; $__heap_top = memory.size * page_size
		memory.size
		i32.const 65536 ;; page size
		i32.mul
		global.set $__heap_top

		;; assert $__heap_top > $__heap_ptr
		global.get $__heap_top
		global.get $__heap_ptr
		i32.gt_s
		if else
			unreachable
		end

		;; 分别初始化空闲链表头(6*8=48个字节), l128 再单独初始化
		;; 固定大小的空闲链表头的 size 对应链表的节点数
		;; l24 | l32 | l48 | l80 | nil
		;; memset(__heap_base, 0, 48)

		global.get $__heap_base
		i32.const 0
		i32.const 48 ;; 6*sizeof(heap_block_t)
		memory.fill

		;; $__heap_l128_freep = $__heap_base + 4*sizeof(heap_block_t)
		;; 该字段改用 global 表示, 不占用内存空间
		global.get $__heap_base
		i32.const 32 ;; 4*sizeof(heap_block_t)
		i32.add
		global.set $__heap_l128_freep

		;; L128 是一个环
		global.get $__heap_l128_freep
		i32.const 0 ;; size
		global.get $__heap_l128_freep ;; next = l128
		call $heap_block.init
	)

	;; func malloc(size: i32) => i32
	;; 在堆上分配内存并返回地址, 内存不超过 2GB, 失败返回 0
	(func $runtime.malloc (param $size i32) (result i32)
		(local $free_list i32) ;; *heap_block_t, 空闲链表头
		(local $b i32) ;; *heap_block_t

		;; 仅初始化1次堆
		call $wa_malloc_init_once

		;; 输入参数对齐到8字节
		local.get $size
		call $heap_alignment8
		local.set $size

		;; 根据大小返回对应空闲链表的地址
		;; 并返回对齐到8字节的大小
		;; $free_list, $size = $heap_free_list_header.ptr_and_fixed_size(size)
		local.get $size
		call $heap_free_list.ptr_and_fixed_size
		local.set $size
		local.set $free_list

		;; 是否禁止了 fixed 策略
		call $heap_is_fixed_list_enabled
		if
			;; 定长和变长有不同的分配策略
			local.get $size
			call $heap_is_fixed_size
			if
				local.get $free_list
				call $wa_malloc_reuse_fixed
				local.tee $b
				i32.eqz
				if else
					;; if ($b != nil) return b->data;
					local.get $b
					call $heap_block.data
					return
				end
			end
		end

		;; 变长内存直接从变长空闲链表分配
		;; 定长分配失败再尝试从变长空闲链表分配
		local.get $size
		call $heap_reuse_varying

		;; 从空闲链表分配成功, 则返回
		local.tee $b
		i32.eqz
		if else
			;; if ($b != nil) return b->data;
			local.get $b
			call $heap_block.data
			return
		end

		;; 空闲链表分配失败, 则开辟新的内存
		local.get $size
		call $heap_new_allocation
		local.tee $b
		i32.eqz
		if else
			;; if ($b != nil) return b->data;
			local.get $b
			call $heap_block.data
			return
		end

		;; 可能的改进: 可尝试回收固定空闲链表, 因为可能有合并的动作产生更大的块但是机会很小
		i32.const 0
		return
	)

	;; func wa_malloc_reuse_fixed(free_list: *heap_block_t) => *heap_block_t
	(func $wa_malloc_reuse_fixed (param $free_list i32) (result i32)
		(local $p i32) ;; *heap_block_t

		;; 禁止了 fixed 策略
		call $heap_assert_fixed_list_enabled

		;; $free_list.size 对应节点的数量
		;; if $free_list.size == 0 { return nil }
		local.get $free_list
		call $heap_block.size
		i32.eqz
		if
			i32.const 0
			return
		end

		;; $free_list.size--
		local.get $free_list
		block (result i32)
			local.get $free_list
			call $heap_block.size
			i32.const 1
			i32.sub
		end
		call $heap_block.set_size

		;; $p = $free_list.next
		local.get $free_list
		call $heap_block.next
		local.set $p

		;; $free_list.next = $p.next
		local.get $free_list
		block (result i32)
			local.get $p
			call $heap_block.next
		end
		call $heap_block.set_next

		;; $p.next = nil
		local.get $p
		i32.const 0
		call $heap_block.set_next

		;; return $p
		local.get $p
		return
	)

	;; func heap_reuse_varying(size: i32) => *heap_block_t
	(func $heap_reuse_varying (param $nbytes i32) (result i32)
		(local $prevp i32) ;; 前趋节点指针
		(local $remaining i32) ;; 分裂后剩余的部分
		(local $p i32) ;; 当前节点指针
		
		;; $p = $prevp.next
		global.get $__heap_l128_freep
		local.set $prevp

		;; for $p = $prevp.next; true; $prevp, $p = $p, $p.next {}
		local.get $prevp
		call $heap_block.next
		local.set $p
		loop $continue
			;; if $p.size >= $nbytes + 8 { 分裂为2块 }
			local.get $p
			call $heap_block.size
			local.get $nbytes
			i32.const 8 ;; sizeof(heap_block_t)
			i32.add
			i32.ge_s
			if
				;; 剩余的部分裂为新块
				;; $remaining = $p.data + $nbytes
				local.get $p
				call $heap_block.data
				local.get $nbytes
				i32.add
				local.set $remaining

				;; $remaining.next = $p.next
				local.get $remaining
				local.get $p
				call $heap_block.next
				call $heap_block.set_next

				;; $remaining.size = $p.size - $nbytes - 8
				local.get $remaining
				local.get $p
				call $heap_block.size
				local.get $nbytes
				i32.sub
				i32.const 8
				i32.sub
				call $heap_block.set_size

				;; 取下 $p 块
				;; $prevp.next = $remaining
				local.get $prevp
				local.get $remaining
				call $heap_block.set_next

				;; $__heap_l128_freep = $prevp
				local.get $prevp
				global.set $__heap_l128_freep

				;; $p.size = $nbytes
				;; $p.data = nil
				local.get $p
				local.get $nbytes
				i32.const 0
				call $heap_block.init

				;; return $p
				local.get $p
				return
			end

			;; if $p.size >= $nbytes { 分配完整一块 }
			local.get $p
			call $heap_block.size
			local.get $nbytes
			i32.ge_s
			if
				;; $prevp.next = $p.next
				local.get $prevp
				block (result i32)
					local.get $p
					call $heap_block.next
				end
				call $heap_block.set_next

				;; $__heap_l128_freep = $prevp
				local.get $prevp
				global.set $__heap_l128_freep

				;; $p.size 不变
				;; $p.data = nil
				local.get $p
				i32.const 0
				call $heap_block.set_next

				;; return $p
				local.get $p
				return
			end

			;; if $p == $__heap_l128_freep { return nil }
			local.get $p
			global.get $__heap_l128_freep
			i32.eq
			if
				i32.const 0
				return
			end

			;; $prevp = $p
			local.get $p
			local.set $prevp

			;; $p = $p.next
			local.get $p
			call $heap_block.next
			local.set $p

			br $continue
		end
		unreachable
	)

	;; func heap_new_allocation(size: i32)=> *void
	(func $heap_new_allocation (param $size i32) (result i32)
		(local $ptr i32)
		(local $block_size i32)
		(local $pages i32)

		;; 先从堆后面未登记的空间分配

		;; $ptr = $__heap_ptr
		global.get $__heap_ptr
		local.set $ptr

		;; $block_size = sizeof(heap_block_t) + $size
		i32.const 8
		local.get $size
		i32.add
		local.set $block_size

		;; 如果已经超出内存最大空间, 则先扩容
		;; if heap_ptr+block_size >= heap_top { grow }
		global.get $__heap_ptr
		local.get $block_size
		i32.add
		global.get $__heap_top
		i32.ge_s
		if
			;; $pages = ($block_size+WASM_PAGE_SIZE-1) / WASM_PAGE_SIZE)
			local.get $block_size
			i32.const 65535 ;; WASM_PAGE_SIZE-1
			i32.add
			i32.const 65536 ;; WASM_PAGE_SIZE
			i32.div_s
			local.set $pages

			;; if memory.grow(pages) < 0 { return nil }
			local.get $pages
			memory.grow
			i32.const 0
			i32.lt_s
			if
				;; return nil
				i32.const 0
				return
			end

			;; 更新已经分配内存的大小
			;; heap_top += (pages * WASM_PAGE_SIZE);
			global.get $__heap_top
			block (result i32)
				local.get $pages
				i32.const 65536
				i32.mul
			end
			i32.add
			global.set $__heap_top

			;; 验证 memory.size
		end

		;; 成功, 更新堆当前指针
		;; $__heap_ptr += block_size;
		global.get $__heap_ptr
		local.get $block_size
		i32.add
		global.set $__heap_ptr

		;; 初始化返回的内存块
		;; 注意: 这里并未标注是否定长数据, 将在释放时根据实际大小确定
		;; $ptr.size = $size
		;; $ptr.next = nil
		local.get $ptr
		local.get $size
		i32.const 0
		call $heap_block.init

		;; return $ptr
		local.get $ptr
		return
	)

	;; func free(ptr: i32)
	(func $runtime.free (param $ptr i32)
		(local $size i32) ;; *heap_block_t
		(local $block i32) ;; *heap_block_t
		(local $freep i32) ;; 空闲链表指针

		;; 仅初始化1次堆
		call $wa_malloc_init_once

		;; 必须是合法的指针
		local.get $ptr
		call $heap_assert_valid_ptr

		;; 指针必须8字节对齐
		local.get $ptr
		call $heap_assert_align8

		;; heap_block_t *block = ptr - sizeof(heap_block_t);
		local.get $ptr
		i32.const 8
		i32.sub
		local.set $block

		;; $size = $block.size
		local.get $block
		call $heap_block.size
		local.set $size

		;; 允许fixed策略?
		call $heap_is_fixed_list_enabled
		if
			;; 如果是固定大小内存
			local.get $size
			call $heap_is_fixed_size
			if
				;; 根据 size 查询对应的空闲链表
				local.get $size
				call $heap_free_list.ptr_and_fixed_size
				drop
				local.get $block
				call $wa_lfixed_free_block
				return
			end
		end

		;; 变长大小内存
		local.get $block
		call $wa_l128_free
		return
	)

	;; 固定尺寸空闲链表释放
	(func $wa_lfixed_free_block (param $freep i32) (param $block i32)
		call $heap_assert_fixed_list_enabled

		;; 如果是空闲链表太长, 则回收到变长空闲链表
		;; if cap > 0 && $freep->size == cap { wa_lfixed_free_all() }
		local.get $freep
		call $heap_block.size
		global.get $__heap_lfixed_cap
		i32.eq
		if
			local.get $freep
			call $wa_lfixed_free_all
		end

		;; $block.next = $freep.next
		local.get $block
		block (result i32)
			local.get $freep
			call $heap_block.next
		end
		call $heap_block.set_next

		;; $freep.next = $block
		local.get $freep
		local.get $block
		call $heap_block.set_next

		;; $freep.size++
		local.get $freep
		block (result i32)
			local.get $freep
			call $heap_block.size
			i32.const 1
			i32.add
		end
		call $heap_block.set_size
	)

	;; 固定尺寸空闲链表释放到 l128
	(func $wa_lfixed_free_all (param $freep i32)
		(local $p i32)
		(local $temp i32)

		call $heap_assert_fixed_list_enabled

		;; p = $freep.next
		local.get $freep
		call $heap_block.next
		local.set $p

		;; while $p != nil
		block $break
			loop $continue
				;; if $p == nil { break }
				local.get $p
				i32.eqz
				if
					br $break
				end

				;; $temp = $p.next
				local.get $p
				call $heap_block.next
				local.set $temp

				;; $wa_l128_free($p)
				local.get $p
				call $wa_l128_free

				;; $p = $temp
				local.get $temp
				local.set $p

				br $continue
			end
		end

		;; 清空容量计数
		local.get $freep
		i32.const 0
		i32.const 0
		call $heap_block.init
	)

	;; kr_free 算法释放
	(func $wa_l128_free (param $bp i32)
		(local $p i32)

		;; for(p = freep; !(bp > p && bp < p->next); p = p->next) {
		;;   if (p >= p->s.ptr && (bp > p || bp < p->s.ptr)) 
		;;     break;
		;; }

		;; p = freep
		global.get $__heap_l128_freep
		local.set $p

		;; 查找到 bp 属于 [p, p->next] 区间结束循环
		;; while !(bp > p && bp < p->next) { ... }
		block $break
			loop $continue
				;; if bp > p && bp < p->next: break
				local.get $bp
				local.get $p
				i32.gt_s
				if
					;; if bp < p->next: break
					local.get $bp
					local.get $p
					call $heap_block.next
					i32.lt_s
					if
						br $break
					end
				end

				;; if (p >= p->next && (bp > p || bp < p->next)) break
				;; if p >= p->next && ...
				local.get $p
				local.get $p
				call $heap_block.next
				i32.ge_s
				if
					;; if bp > p || bp < p->next: break
					local.get $bp
					local.get $p
					i32.gt_s
					if
						br $break
					end

					;; if bp < p->next: break
					local.get $bp
					local.get $p
					call $heap_block.next
					i32.lt_s
					if
						br $break
					end
				end

				;; p = p.next
				local.get $p
				call $heap_block.next
				local.set $p

				br $continue
			end
		end

		;; if (bp + bp->s.size + 8 == p->s.ptr) { ... }
		local.get $bp
		local.get $bp
		call $heap_block.size		
		i32.add ;; 需要确保对齐, size 是字节数, bp 是 *heap_block_t
		i32.const 8 ;; sizeof(heap_block_t)
		i32.add
		local.get $p
		call $heap_block.next
		i32.eq
		if ;; join to upper nbr
			;; bp->s.size += p->s.ptr->s.size + 8;
			local.get $bp
			block (result i32)
				local.get $bp
				call $heap_block.size

				local.get $p
				call $heap_block.next
				call $heap_block.size

				i32.const 8

				i32.add
				i32.add
			end
			call $heap_block.set_size

			;; bp->s.ptr = p->s.ptr->s.ptr;
			local.get $bp
			block (result i32)
				local.get $p
				call $heap_block.next
				call $heap_block.next
			end
			call $heap_block.set_next

		else
			;; bp->s.ptr = p->s.ptr;
			local.get $bp
			local.get $p
			call $heap_block.next
			call $heap_block.set_next
		end

		;; if (p + p->s.size + 8 == bp) { ... }
		local.get $p
		local.get $p
		call $heap_block.size
		i32.add
		i32.const 8
		i32.add
		local.get $bp
		i32.eq
		if ;; join to lower nbr
			;; p->s.size += bp->s.size + 8;
			local.get $p
			block (result i32)
				local.get $p
				call $heap_block.size

				local.get $bp
				call $heap_block.size

				i32.const 8

				i32.add
				i32.add
			end
			call $heap_block.set_size

			;; p->s.ptr = bp->s.ptr;
			local.get $p
			local.get $bp
			call $heap_block.next
			call $heap_block.set_next
		else
			;; p->s.ptr = bp;
			local.get $p
			local.get $bp
			call $heap_block.set_next
		end

		;; freep = p;
		local.get $p
		global.set $__heap_l128_freep
	)
;; file: interface.wat.ws

(func $$wa.runtime.queryIface (param $d.b i32) (param $d.d i32) (param $itab i32) (param $eq i32) (param $ihash i32)
  (result i32 i32 i32 i32)
	(local $t i32)
	local.get $itab
	if (result i32 i32 i32 i32)
	  local.get $itab
	  i32.load offset=0 align=4
	  local.get $ihash
	  i32.const 0
	  call $runtime.getItab
	  local.set $t
	  local.get $t
	  if (result i32 i32 i32 i32)
	    local.get $d.b
		call $runtime.Block.Retain
	    local.get $d.d
	    local.get $t
	    local.get $eq
	  else
	    i32.const 0
	    i32.const 0
	    i32.const 0
	    i32.const 0
	    unreachable
	  end
	else
	  i32.const 0
	  i32.const 0
	  i32.const 0
	  i32.const 0
	  unreachable
	end
)

(func $$wa.runtime.queryIface_CommaOk (param $d.b i32) (param $d.d i32) (param $itab i32) (param $eq i32) (param $ihash i32)
  (result i32 i32 i32 i32 i32)
	(local $t i32)
	local.get $itab
	if (result i32 i32 i32 i32 i32)
	  local.get $itab
	  i32.load offset=0 align=4
	  local.get $ihash
	  i32.const 1
	  call $runtime.getItab
	  local.set $t
	  local.get $t
	  if (result i32 i32 i32 i32 i32)
	    local.get $d.b
		call $runtime.Block.Retain
	    local.get $d.d
	    local.get $t
	    local.get $eq
	    i32.const 1
	  else
	    i32.const 0
	    i32.const 0
	    i32.const 0
	    i32.const 0
	    i32.const 0
	  end
	else
	  i32.const 0
	  i32.const 0
	  i32.const 0
	  i32.const 0
	  i32.const 0
	end
)

(func $runtime.Compare (param $l.d.b i32) (param $l.d.d i32) (param $l.itab i32) (param $l.comp i32) (param $r.d.b i32) (param $r.d.d i32) (param $r.itab i32) (param $r.comp i32)
  (result i32)
  local.get $l.comp
  local.get $r.comp
  i32.lt_s
  if (result i32) ;;if l.comp < r.comp
    i32.const -1
  else
    local.get $l.comp
	local.get $r.comp
	i32.gt_s
	if (result i32) ;;if l.comp > r.comp
	  i32.const 1
	else ;;if l.comp == r.comp:
	  local.get $l.comp
	  if (result i32) ;;if comp != 0, compare by type.comp:
	    local.get $l.d.d
	    local.get $r.d.d
	    local.get $l.comp
		call_indirect (type $$wa.runtime.comp)
	  else ;;if comp == 0, compare as ref:
	    local.get $l.d.d
		local.get $r.d.d
		i32.lt_u
		if (result i32)
		  i32.const -1
		else
		  local.get $l.d.d
		  local.get $r.d.d
		  i32.gt_u
		end
	  end
	end
  end
)
;; file: map.wat.ws


;; file: string.wat.ws

(func $$wa.runtime.string_to_ptr (param $b i32) (param $d i32) (param $l i32) (result i32) ;;result = ptr
	local.get $d
)

(func $$wa.runtime.string_to_iter (param $b i32) (param $d i32) (param $l i32) (result i32 i32 i32)
	local.get $d
	local.get $l
	i32.const 0
)
;; ---------------------------------------------------------
;; package: syscall/js
;; ---------------------------------------------------------

;; file: z_abi.wat.ws

;; Copyright 2024 The Wa Authors. All rights reserved.

(func $$syscall/js.__linkname__string_to_ptr (param $b i32) (param $d i32) (param $l i32) (result i32) ;;result = ptr
	local.get $d
)
(data (i32.const 8224) "\24\24\77\61\64\73\24\24\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\30\61\73\73\65\72\74\20\66\61\69\6c\65\64\20\28\61\73\73\65\72\74\20\66\61\69\6c\65\64\3a\20\2b\69\29\6e\69\6c\20\6d\61\70\2e\6d\61\70\2e\77\61\3a\36\38\3a\38\70\61\6e\69\63\3a\20\74\72\75\65\66\61\6c\73\65\4e\61\4e\2b\49\6e\66\2d\49\6e\66\30\31\32\33\34\35\36\37\38\39\61\62\63\64\65\66\0a\5b\2f\5d\2b\2b\2b\2b\2b\2b\2b\2b\2b\2b\5b\3e\2b\2b\2b\2b\2b\2b\2b\2b\2b\2b\3c\2d\5d\3e\2b\2b\2b\2b\2e\2b\2e\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\01\00\00\00\ff\ff\ff\ff\08\96\00\00")
(table 24 funcref)
(elem (i32.const 1) $$string.$$compAddr)
(elem (i32.const 2) $$u8.$$block.$$OnFree)
(elem (i32.const 3) $$string.underlying.$$OnFree)
(elem (i32.const 4) $$runtime.mapImp.$$block.$$OnFree)
(elem (i32.const 5) $$runtime.mapImp.$ref.underlying.$$OnFree)
(elem (i32.const 6) $$runtime.mapIter.$$OnFree)
(elem (i32.const 7) $$runtime.mapNode.$$block.$$OnFree)
(elem (i32.const 8) $$runtime.mapNode.$ref.underlying.$$OnFree)
(elem (i32.const 9) $$void.$$block.$$OnFree)
(elem (i32.const 10) $$void.$ref.underlying.$$OnFree)
(elem (i32.const 11) $$i`0`.underlying.$$OnFree)
(elem (i32.const 12) $$runtime.mapNode.$$OnFree)
(elem (i32.const 13) $$runtime.mapNode.$ref.$$block.$$OnFree)
(elem (i32.const 14) $$runtime.mapNode.$ref.$slice.underlying.$$OnFree)
(elem (i32.const 15) $$runtime.mapImp.$$OnFree)
(elem (i32.const 16) $$runtime.mapNode.$ref.$array1.underlying.$$OnFree)
(elem (i32.const 17) $$$$$$.underlying.$$OnFree)
(elem (i32.const 18) $$$$$$.$array1.underlying.$$OnFree)
(elem (i32.const 19) $$$$$$.$$block.$$OnFree)
(elem (i32.const 20) $$$$$$.$slice.underlying.$$OnFree)
(elem (i32.const 21) $$runtime.defers.$$OnFree)
(elem (i32.const 22) $$runtime.defers.$array1.underlying.$$OnFree)
(elem (i32.const 23) $$brainfuck$bfpkg.BrainFuck.$$OnFree)
(type $$OnFree (func (param i32)))
(type $$wa.runtime.comp (func (param i32) (param i32) (result i32)))
(type $$$fnSig1 (func))
(global $$wa.runtime.closure_data (mut i32) (i32.const 0))
(global $$wa.runtime._concretTypeCount (mut i32) (i32.const 1))
(global $$wa.runtime._interfaceCount (mut i32) (i32.const 1))
(global $$wa.runtime._itabsPtr (mut i32) (i32.const 38416))
(global $runtime.defersStack.0 i32 (i32.const 0))
(global $runtime.defersStack.1 i32 (i32.const 8232))
(global $runtime.init$guard (mut i32) (i32.const 0))
(global $$knr_basep (mut i32) (i32.const 0))
(global $$knr_freep (mut i32) (i32.const 0))
(global $brainfuck$bfpkg.init$guard (mut i32) (i32.const 0))
(global $brainfuck.init$guard (mut i32) (i32.const 0))
(global $syscall$js.init$guard (mut i32) (i32.const 0))
(global $runtime.zptr (mut i32) (i32.const 8384))
(global $__heap_base i32 (i32.const 38448))


(func $$string.appendstr (param $x.0 i32) (param $x.1 i32) (param $x.2 i32) (param $y.0 i32) (param $y.1 i32) (param $y.2 i32) (result i32 i32 i32)
  (local $x_len i32)
  (local $y_len i32)
  (local $new_len i32)
  (local $item i32)
  (local $src i32)
  (local $dest i32)
  local.get $x.2
  local.set $x_len
  local.get $y.2
  local.set $y_len
  local.get $x_len
  local.get $y_len
  i32.add
  local.set $new_len
  local.get $new_len
  i32.const 1
  i32.mul
  i32.const 16
  i32.add
  call $runtime.HeapAlloc
  local.get $new_len
  i32.const 0
  i32.const 1
  call $runtime.Block.Init
  call $runtime.DupI32
  i32.const 16
  i32.add
  call $runtime.DupI32
  local.set $dest
  local.get $new_len
  local.get $x.1
  local.set $src
  block $block2
    loop $loop2
      local.get $x_len
      i32.eqz
      if
        br $block2
      else
      end
      local.get $src
      i32.load8_u offset=0 align=1
      local.set $item
      local.get $dest
      local.get $item
      i32.store8 offset=0 align=1
      local.get $src
      i32.const 1
      i32.add
      local.set $src
      local.get $dest
      i32.const 1
      i32.add
      local.set $dest
      local.get $x_len
      i32.const 1
      i32.sub
      local.set $x_len
      br $loop2
    end ;;loop2
  end ;;block2
  local.get $y.1
  local.set $src
  block $block3
    loop $loop3
      local.get $y_len
      i32.eqz
      if
        br $block3
      else
      end
      local.get $src
      i32.load8_u offset=0 align=1
      local.set $item
      local.get $dest
      local.get $item
      i32.store8 offset=0 align=1
      local.get $src
      i32.const 1
      i32.add
      local.set $src
      local.get $dest
      i32.const 1
      i32.add
      local.set $dest
      local.get $y_len
      i32.const 1
      i32.sub
      local.set $y_len
      br $loop3
    end ;;loop3
  end ;;block3
) ;;$string.appendstr

(func $$string.equal (param $x.0 i32) (param $x.1 i32) (param $x.2 i32) (param $y.0 i32) (param $y.1 i32) (param $y.2 i32) (result i32)
  (local $ret i32)
  i32.const 1
  local.set $ret
  local.get $x.2
  local.get $y.2
  i32.ne
  if
    i32.const 0
    local.set $ret
  else
    loop $loop1
      local.get $x.2
      if
        local.get $x.1
        local.get $x.2
        i32.add
        i32.const 1
        i32.sub
        i32.load8_u offset=0 align=1
        local.get $y.1
        local.get $x.2
        i32.add
        i32.const 1
        i32.sub
        i32.load8_u offset=0 align=1
        i32.eq
        if
          local.get $x.2
          i32.const 1
          i32.sub
          local.set $x.2
          br $loop1
        else
          i32.const 0
          local.set $ret
        end
      else
      end
    end ;;loop1
  end
  local.get $ret
) ;;$string.equal

(func $$string.$$compAddr (param $p0 i32) (param $p1 i32) (result i32)
  (local $v0.0 i32)
  (local $v0.1 i32)
  (local $v0.2 i32)
  (local $v1.0 i32)
  (local $v1.1 i32)
  (local $v1.2 i32)
  local.get $p0
  if
    local.get $p0
    i32.load offset=0 align=4
    call $runtime.Block.Retain
    local.get $p0
    i32.load offset=4 align=4
    local.get $p0
    i32.load offset=8 align=4
    local.set $v0.2
    local.set $v0.1
    local.get $v0.0
    call $runtime.Block.Release
    local.set $v0.0
  else
  end
  local.get $p1
  if
    local.get $p1
    i32.load offset=0 align=4
    call $runtime.Block.Retain
    local.get $p1
    i32.load offset=4 align=4
    local.get $p1
    i32.load offset=8 align=4
    local.set $v1.2
    local.set $v1.1
    local.get $v1.0
    call $runtime.Block.Release
    local.set $v1.0
  else
  end
  local.get $v0.0
  local.get $v0.1
  local.get $v0.2
  local.get $v1.0
  local.get $v1.1
  local.get $v1.2
  call $$wa.runtime.string_Comp
  local.get $v0.0
  call $runtime.Block.Release
  local.get $v1.0
  call $runtime.Block.Release
) ;;$string.$$compAddr

(func $$u8.$$block.$$OnFree (param $ptr i32)
  local.get $ptr
  i32.load offset=0 align=1
  call $runtime.Block.Release
  local.get $ptr
  i32.const 0
  i32.store offset=0 align=1
) ;;$u8.$$block.$$OnFree

(func $$string.underlying.$$OnFree (param $$ptr i32)
  local.get $$ptr
  i32.const 2
  call_indirect (type $$OnFree)
) ;;$string.underlying.$$OnFree

(func $runtime.ActivateEmptyInterface
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$t0.0.0 i32)
  (local $$t0.0.1 i32)
  (local $$t0.1 i32)
  (local $$t0.2 i32)
  (local $$t1.0 i32)
  (local $$t1.1 i32)
  (local $$t1.2 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;make interface{} <- string ("":string)
        i32.const 28
        call $runtime.HeapAlloc
        i32.const 1
        i32.const 3
        i32.const 12
        call $runtime.Block.Init
        call $runtime.DupI32
        i32.const 16
        i32.add
        call $runtime.DupI32
        call $runtime.DupI32
        i32.const 0
        call $runtime.SwapI32
        i32.load offset=0 align=1
        call $runtime.Block.Release
        i32.store offset=0 align=1
        call $runtime.DupI32
        i32.const 8224
        i32.store offset=4 align=4
        call $runtime.DupI32
        i32.const 0
        i32.store offset=8 align=4
        i32.const 1
        i32.const -1
        i32.const 0
        call $runtime.getItab
        i32.const 1
        local.set $$t0.2
        local.set $$t0.1
        local.set $$t0.0.1
        local.get $$t0.0.0
        call $runtime.Block.Release
        local.set $$t0.0.0

        ;;typeassert t0.(string)
        local.get $$t0.1
        i32.load offset=0 align=4
        i32.const 1
        i32.eq
        if (result i32 i32 i32)
          local.get $$t0.0.1
          i32.load offset=0 align=4
          call $runtime.Block.Retain
          local.get $$t0.0.1
          i32.load offset=4 align=4
          local.get $$t0.0.1
          i32.load offset=8 align=4
        else
          i32.const 0
          i32.const 8248
          i32.const 1
          unreachable
        end
        local.set $$t1.2
        local.set $$t1.1
        local.get $$t1.0
        call $runtime.Block.Release
        local.set $$t1.0

        ;;println(t1)
        local.get $$t1.1
        local.get $$t1.2
        call $$runtime.waPuts
        i32.const 10
        call $$runtime.waPrintChar

        ;;return
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$t0.0.0
  call $runtime.Block.Release
  local.get $$t1.0
  call $runtime.Block.Release
) ;;runtime.ActivateEmptyInterface

(func $$runtime.argsGet (param $result_argv i32) (param $result_argv_buf i32) (result i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$ret_0 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;return 0:i32
        i32.const 0
        local.set $$ret_0
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$ret_0
) ;;$runtime.argsGet

(func $$runtime.argsSizesGet (param $result_argc i32) (param $result_argv_len i32) (result i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$ret_0 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;return 0:i32
        i32.const 0
        local.set $$ret_0
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$ret_0
) ;;$runtime.argsSizesGet

(func $$runtime.assert (param $ok i32) (param $pos_msg_ptr i32) (param $pos_msg_len i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$t0 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_2
        block $$Block_1
          block $$Block_0
            block $$BlockSel
              local.get $$block_selector
              br_table 0 1 2 0
            end ;;$BlockSel
            i32.const 0
            local.set $$current_block

            ;;ok == 0:i32
            local.get $ok
            i32.const 0
            i32.eq
            local.set $$t0

            ;;if t0 goto 1 else 2
            local.get $$t0
            if
              br $$Block_0
            else
              br $$Block_1
            end

          end ;;$Block_0
          i32.const 1
          local.set $$current_block

          ;;printString("assert failed (":string)
          i32.const 0
          i32.const 8249
          i32.const 15
          call $runtime.printString

          ;;waPuts(pos_msg_ptr, pos_msg_len)
          local.get $pos_msg_ptr
          local.get $pos_msg_len
          call $$runtime.waPuts

          ;;waPrintRune(41:i32)
          i32.const 41
          call $$runtime.waPrintRune

          ;;waPrintRune(10:i32)
          i32.const 10
          call $$runtime.waPrintRune

          ;;procExit(1:i32)
          i32.const 1
          call $$runtime.procExit

          ;;jump 2
          br $$Block_1

        end ;;$Block_1
        i32.const 2
        local.set $$current_block

        ;;return
        br $$BlockFnBody

      end ;;$Block_2
    end ;;$BlockDisp
  end ;;$BlockFnBody
) ;;$runtime.assert

(func $$runtime.assertWithMessage (param $ok i32) (param $msg_ptr i32) (param $msg_len i32) (param $pos_msg_ptr i32) (param $pos_msg_len i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$t0 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_2
        block $$Block_1
          block $$Block_0
            block $$BlockSel
              local.get $$block_selector
              br_table 0 1 2 0
            end ;;$BlockSel
            i32.const 0
            local.set $$current_block

            ;;ok == 0:i32
            local.get $ok
            i32.const 0
            i32.eq
            local.set $$t0

            ;;if t0 goto 1 else 2
            local.get $$t0
            if
              br $$Block_0
            else
              br $$Block_1
            end

          end ;;$Block_0
          i32.const 1
          local.set $$current_block

          ;;printString("assert failed: ":string)
          i32.const 0
          i32.const 8264
          i32.const 15
          call $runtime.printString

          ;;waPuts(msg_ptr, msg_len)
          local.get $msg_ptr
          local.get $msg_len
          call $$runtime.waPuts

          ;;printString(" (":string)
          i32.const 0
          i32.const 8262
          i32.const 2
          call $runtime.printString

          ;;waPuts(pos_msg_ptr, pos_msg_len)
          local.get $pos_msg_ptr
          local.get $pos_msg_len
          call $$runtime.waPuts

          ;;waPrintRune(41:i32)
          i32.const 41
          call $$runtime.waPrintRune

          ;;waPrintRune(10:i32)
          i32.const 10
          call $$runtime.waPrintRune

          ;;procExit(1:i32)
          i32.const 1
          call $$runtime.procExit

          ;;jump 2
          br $$Block_1

        end ;;$Block_1
        i32.const 2
        local.set $$current_block

        ;;return
        br $$BlockFnBody

      end ;;$Block_2
    end ;;$BlockDisp
  end ;;$BlockFnBody
) ;;$runtime.assertWithMessage

(func $$wa.runtime.complex128_Add (param $a f64) (param $ai f64) (param $b f64) (param $bi f64) (result f64 f64)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$ret_0 f64)
  (local $$ret_1 f64)
  (local $$t0 f64)
  (local $$t1 f64)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;a + b
        local.get $a
        local.get $b
        f64.add
        local.set $$t0

        ;;ai + bi
        local.get $ai
        local.get $bi
        f64.add
        local.set $$t1

        ;;return t0, t1
        local.get $$t0
        local.set $$ret_0
        local.get $$t1
        local.set $$ret_1
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$ret_0
  local.get $$ret_1
) ;;$wa.runtime.complex128_Add

(func $$wa.runtime.complex128_Div (param $a f64) (param $ai f64) (param $b f64) (param $bi f64) (result f64 f64)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$ret_0 f64)
  (local $$ret_1 f64)
  (local $$t0 f64)
  (local $$t1 f64)
  (local $$t2 f64)
  (local $$t3 f64)
  (local $$t4 f64)
  (local $$t5 f64)
  (local $$t6 f64)
  (local $$t7 f64)
  (local $$t8 f64)
  (local $$t9 f64)
  (local $$t10 f64)
  (local $$t11 f64)
  (local $$t12 f64)
  (local $$t13 f64)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;a * b
        local.get $a
        local.get $b
        f64.mul
        local.set $$t0

        ;;ai * bi
        local.get $ai
        local.get $bi
        f64.mul
        local.set $$t1

        ;;t0 + t1
        local.get $$t0
        local.get $$t1
        f64.add
        local.set $$t2

        ;;b * b
        local.get $b
        local.get $b
        f64.mul
        local.set $$t3

        ;;bi * bi
        local.get $bi
        local.get $bi
        f64.mul
        local.set $$t4

        ;;t3 + t4
        local.get $$t3
        local.get $$t4
        f64.add
        local.set $$t5

        ;;t2 / t5
        local.get $$t2
        local.get $$t5
        f64.div
        local.set $$t6

        ;;ai * b
        local.get $ai
        local.get $b
        f64.mul
        local.set $$t7

        ;;a * bi
        local.get $a
        local.get $bi
        f64.mul
        local.set $$t8

        ;;t7 - t8
        local.get $$t7
        local.get $$t8
        f64.sub
        local.set $$t9

        ;;b * b
        local.get $b
        local.get $b
        f64.mul
        local.set $$t10

        ;;bi * bi
        local.get $bi
        local.get $bi
        f64.mul
        local.set $$t11

        ;;t10 + t11
        local.get $$t10
        local.get $$t11
        f64.add
        local.set $$t12

        ;;t9 / t12
        local.get $$t9
        local.get $$t12
        f64.div
        local.set $$t13

        ;;return t6, t13
        local.get $$t6
        local.set $$ret_0
        local.get $$t13
        local.set $$ret_1
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$ret_0
  local.get $$ret_1
) ;;$wa.runtime.complex128_Div

(func $$wa.runtime.complex128_Mul (param $a f64) (param $ai f64) (param $b f64) (param $bi f64) (result f64 f64)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$ret_0 f64)
  (local $$ret_1 f64)
  (local $$t0 f64)
  (local $$t1 f64)
  (local $$t2 f64)
  (local $$t3 f64)
  (local $$t4 f64)
  (local $$t5 f64)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;a * b
        local.get $a
        local.get $b
        f64.mul
        local.set $$t0

        ;;ai * bi
        local.get $ai
        local.get $bi
        f64.mul
        local.set $$t1

        ;;t0 - t1
        local.get $$t0
        local.get $$t1
        f64.sub
        local.set $$t2

        ;;ai * b
        local.get $ai
        local.get $b
        f64.mul
        local.set $$t3

        ;;a * bi
        local.get $a
        local.get $bi
        f64.mul
        local.set $$t4

        ;;t3 + t4
        local.get $$t3
        local.get $$t4
        f64.add
        local.set $$t5

        ;;return t2, t5
        local.get $$t2
        local.set $$ret_0
        local.get $$t5
        local.set $$ret_1
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$ret_0
  local.get $$ret_1
) ;;$wa.runtime.complex128_Mul

(func $$wa.runtime.complex128_Print (param $a f64) (param $ai f64)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$t0 i64)
  (local $$t1 f64)
  (local $$t2 i32)
  (local $$t3 i32)
  (local $$t4 i64)
  (local $$t5 f64)
  (local $$t6 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_8
        block $$Block_7
          block $$Block_6
            block $$Block_5
              block $$Block_4
                block $$Block_3
                  block $$Block_2
                    block $$Block_1
                      block $$Block_0
                        block $$BlockSel
                          local.get $$block_selector
                          br_table 0 1 2 3 4 5 6 7 8 0
                        end ;;$BlockSel
                        i32.const 0
                        local.set $$current_block

                        ;;printString("(":string)
                        i32.const 0
                        i32.const 8263
                        i32.const 1
                        call $runtime.printString

                        ;;convert i64 <- f64 (a)
                        local.get $a
                        i64.trunc_f64_s
                        local.set $$t0

                        ;;convert f64 <- i64 (t1)
                        local.get $$t0
                        f64.convert_i64_s
                        local.set $$t1

                        ;;t2 == a
                        local.get $$t1
                        local.get $a
                        f64.eq
                        local.set $$t2

                        ;;if t3 goto 1 else 3
                        local.get $$t2
                        if
                          br $$Block_0
                        else
                          br $$Block_2
                        end

                      end ;;$Block_0
                      i32.const 1
                      local.set $$current_block

                      ;;printI64(t1)
                      local.get $$t0
                      call $runtime.printI64

                      ;;jump 2
                      br $$Block_1

                    end ;;$Block_1
                    i32.const 2
                    local.set $$current_block

                    ;;ai >= 0:f64
                    local.get $ai
                    f64.const 0
                    f64.ge
                    local.set $$t3

                    ;;if t5 goto 4 else 5
                    local.get $$t3
                    if
                      br $$Block_3
                    else
                      br $$Block_4
                    end

                  end ;;$Block_2
                  i32.const 3
                  local.set $$current_block

                  ;;printF64(a)
                  local.get $a
                  call $runtime.printF64

                  ;;jump 2
                  i32.const 2
                  local.set $$block_selector
                  br $$BlockDisp

                end ;;$Block_3
                i32.const 4
                local.set $$current_block

                ;;printString("+":string)
                i32.const 0
                i32.const 8279
                i32.const 1
                call $runtime.printString

                ;;jump 5
                br $$Block_4

              end ;;$Block_4
              i32.const 5
              local.set $$current_block

              ;;convert i64 <- f64 (ai)
              local.get $ai
              i64.trunc_f64_s
              local.set $$t4

              ;;convert f64 <- i64 (t8)
              local.get $$t4
              f64.convert_i64_s
              local.set $$t5

              ;;t9 == ai
              local.get $$t5
              local.get $ai
              f64.eq
              local.set $$t6

              ;;if t10 goto 6 else 8
              local.get $$t6
              if
                br $$Block_5
              else
                br $$Block_7
              end

            end ;;$Block_5
            i32.const 6
            local.set $$current_block

            ;;printI64(t8)
            local.get $$t4
            call $runtime.printI64

            ;;jump 7
            br $$Block_6

          end ;;$Block_6
          i32.const 7
          local.set $$current_block

          ;;printString("i)":string)
          i32.const 0
          i32.const 8280
          i32.const 2
          call $runtime.printString

          ;;return
          br $$BlockFnBody

        end ;;$Block_7
        i32.const 8
        local.set $$current_block

        ;;printF64(ai)
        local.get $ai
        call $runtime.printF64

        ;;jump 7
        i32.const 7
        local.set $$block_selector
        br $$BlockDisp

      end ;;$Block_8
    end ;;$BlockDisp
  end ;;$BlockFnBody
) ;;$wa.runtime.complex128_Print

(func $$wa.runtime.complex128_Sub (param $a f64) (param $ai f64) (param $b f64) (param $bi f64) (result f64 f64)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$ret_0 f64)
  (local $$ret_1 f64)
  (local $$t0 f64)
  (local $$t1 f64)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;a - b
        local.get $a
        local.get $b
        f64.sub
        local.set $$t0

        ;;ai - bi
        local.get $ai
        local.get $bi
        f64.sub
        local.set $$t1

        ;;return t0, t1
        local.get $$t0
        local.set $$ret_0
        local.get $$t1
        local.set $$ret_1
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$ret_0
  local.get $$ret_1
) ;;$wa.runtime.complex128_Sub

(func $$wa.runtime.complex64_Add (param $a f32) (param $ai f32) (param $b f32) (param $bi f32) (result f32 f32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$ret_0 f32)
  (local $$ret_1 f32)
  (local $$t0 f32)
  (local $$t1 f32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;a + b
        local.get $a
        local.get $b
        f32.add
        local.set $$t0

        ;;ai + bi
        local.get $ai
        local.get $bi
        f32.add
        local.set $$t1

        ;;return t0, t1
        local.get $$t0
        local.set $$ret_0
        local.get $$t1
        local.set $$ret_1
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$ret_0
  local.get $$ret_1
) ;;$wa.runtime.complex64_Add

(func $$wa.runtime.complex64_Div (param $a f32) (param $ai f32) (param $b f32) (param $bi f32) (result f32 f32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$ret_0 f32)
  (local $$ret_1 f32)
  (local $$t0 f32)
  (local $$t1 f32)
  (local $$t2 f32)
  (local $$t3 f32)
  (local $$t4 f32)
  (local $$t5 f32)
  (local $$t6 f32)
  (local $$t7 f32)
  (local $$t8 f32)
  (local $$t9 f32)
  (local $$t10 f32)
  (local $$t11 f32)
  (local $$t12 f32)
  (local $$t13 f32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;a * b
        local.get $a
        local.get $b
        f32.mul
        local.set $$t0

        ;;ai * bi
        local.get $ai
        local.get $bi
        f32.mul
        local.set $$t1

        ;;t0 + t1
        local.get $$t0
        local.get $$t1
        f32.add
        local.set $$t2

        ;;b * b
        local.get $b
        local.get $b
        f32.mul
        local.set $$t3

        ;;bi * bi
        local.get $bi
        local.get $bi
        f32.mul
        local.set $$t4

        ;;t3 + t4
        local.get $$t3
        local.get $$t4
        f32.add
        local.set $$t5

        ;;t2 / t5
        local.get $$t2
        local.get $$t5
        f32.div
        local.set $$t6

        ;;ai * b
        local.get $ai
        local.get $b
        f32.mul
        local.set $$t7

        ;;a * bi
        local.get $a
        local.get $bi
        f32.mul
        local.set $$t8

        ;;t7 - t8
        local.get $$t7
        local.get $$t8
        f32.sub
        local.set $$t9

        ;;b * b
        local.get $b
        local.get $b
        f32.mul
        local.set $$t10

        ;;bi * bi
        local.get $bi
        local.get $bi
        f32.mul
        local.set $$t11

        ;;t10 + t11
        local.get $$t10
        local.get $$t11
        f32.add
        local.set $$t12

        ;;t9 / t12
        local.get $$t9
        local.get $$t12
        f32.div
        local.set $$t13

        ;;return t6, t13
        local.get $$t6
        local.set $$ret_0
        local.get $$t13
        local.set $$ret_1
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$ret_0
  local.get $$ret_1
) ;;$wa.runtime.complex64_Div

(func $$wa.runtime.complex64_Mul (param $a f32) (param $ai f32) (param $b f32) (param $bi f32) (result f32 f32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$ret_0 f32)
  (local $$ret_1 f32)
  (local $$t0 f32)
  (local $$t1 f32)
  (local $$t2 f32)
  (local $$t3 f32)
  (local $$t4 f32)
  (local $$t5 f32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;a * b
        local.get $a
        local.get $b
        f32.mul
        local.set $$t0

        ;;ai * bi
        local.get $ai
        local.get $bi
        f32.mul
        local.set $$t1

        ;;t0 - t1
        local.get $$t0
        local.get $$t1
        f32.sub
        local.set $$t2

        ;;ai * b
        local.get $ai
        local.get $b
        f32.mul
        local.set $$t3

        ;;a * bi
        local.get $a
        local.get $bi
        f32.mul
        local.set $$t4

        ;;t3 + t4
        local.get $$t3
        local.get $$t4
        f32.add
        local.set $$t5

        ;;return t2, t5
        local.get $$t2
        local.set $$ret_0
        local.get $$t5
        local.set $$ret_1
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$ret_0
  local.get $$ret_1
) ;;$wa.runtime.complex64_Mul

(func $$wa.runtime.complex64_Print (param $a f32) (param $ai f32)
  (local $$block_selector i32)
  (local $$current_block i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;print("(":string, a, ai, "i)":string)
        i32.const 8263
        i32.const 1
        call $$runtime.waPuts
        i32.const 32
        call $$runtime.waPrintRune
        local.get $a
        call $$runtime.waPrintF32
        i32.const 32
        call $$runtime.waPrintRune
        local.get $ai
        call $$runtime.waPrintF32
        i32.const 32
        call $$runtime.waPrintRune
        i32.const 8280
        i32.const 2
        call $$runtime.waPuts

        ;;return
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
) ;;$wa.runtime.complex64_Print

(func $$wa.runtime.complex64_Sub (param $a f32) (param $ai f32) (param $b f32) (param $bi f32) (result f32 f32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$ret_0 f32)
  (local $$ret_1 f32)
  (local $$t0 f32)
  (local $$t1 f32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;a - b
        local.get $a
        local.get $b
        f32.sub
        local.set $$t0

        ;;ai - bi
        local.get $ai
        local.get $bi
        f32.sub
        local.set $$t1

        ;;return t0, t1
        local.get $$t0
        local.set $$ret_0
        local.get $$t1
        local.set $$ret_1
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$ret_0
  local.get $$ret_1
) ;;$wa.runtime.complex64_Sub

(func $$runtime.environGet (param $result_environv i32) (param $result_environv_buf i32) (result i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$ret_0 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;return 0:i32
        i32.const 0
        local.set $$ret_0
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$ret_0
) ;;$runtime.environGet

(func $$runtime.environSizesGet (param $result_environc i32) (param $result_environv_len i32) (result i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$ret_0 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;return 0:i32
        i32.const 0
        local.set $$ret_0
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$ret_0
) ;;$runtime.environSizesGet

(func $$runtime.fdWrite (param $fd i32) (param $io i32) (param $iovs_len i32) (param $nwritten i32) (result i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$ret_0 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;return 0:i32
        i32.const 0
        local.set $$ret_0
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$ret_0
) ;;$runtime.fdWrite

(func $runtime.getItab (param $dhash i32) (param $ihash i32) (param $commanok i32) (result i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$ret_0 i32)
  (local $$t0 i32)
  (local $$t1 i32)
  (local $$t2 i32)
  (local $$t3 i32)
  (local $$t4 i32)
  (local $$t5 i32)
  (local $$t6 i32)
  (local $$t7 i32)
  (local $$t8 i32)
  (local $$t9 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;*_itabsPtr
        global.get $$wa.runtime._itabsPtr
        local.set $$t0

        ;;dhash - 1:i32
        local.get $dhash
        i32.const 1
        i32.sub
        local.set $$t1

        ;;*_interfaceCount
        global.get $$wa.runtime._interfaceCount
        local.set $$t2

        ;;t1 * t2
        local.get $$t1
        local.get $$t2
        i32.mul
        local.set $$t3

        ;;t3 - ihash
        local.get $$t3
        local.get $ihash
        i32.sub
        local.set $$t4

        ;;t4 - 1:i32
        local.get $$t4
        i32.const 1
        i32.sub
        local.set $$t5

        ;;t5 * 4:i32
        local.get $$t5
        i32.const 4
        i32.mul
        local.set $$t6

        ;;t0 + t6
        local.get $$t0
        local.get $$t6
        i32.add
        local.set $$t7

        ;;convert u32 <- i32 (t7)
        local.get $$t7
        local.set $$t8

        ;;getU32(t8)
        local.get $$t8
        call $runtime.getU32
        local.set $$t9

        ;;return t9
        local.get $$t9
        local.set $$ret_0
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$ret_0
) ;;runtime.getItab

(func $$wa.runtime.getTypePtr (param $hash i32) (result i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$ret_0 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;return 0:uintptr
        i32.const 0
        local.set $$ret_0
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$ret_0
) ;;$wa.runtime.getTypePtr

(func $runtime.getU32 (param $addr i32) (result i32)
  local.get $addr
  i32.load offset=0 align=4
) ;;runtime.getU32

(func $runtime.get_u8 (param $addr i32) (result i32)
  local.get $addr
  i32.load8_u offset=0 align=1
) ;;runtime.get_u8

(func $runtime.init
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$t0 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_2
        block $$Block_1
          block $$Block_0
            block $$BlockSel
              local.get $$block_selector
              br_table 0 1 2 0
            end ;;$BlockSel
            i32.const 0
            local.set $$current_block

            ;;*init$guard
            global.get $runtime.init$guard
            local.set $$t0

            ;;if t0 goto 2 else 1
            local.get $$t0
            if
              br $$Block_1
            else
              br $$Block_0
            end

          end ;;$Block_0
          i32.const 1
          local.set $$current_block

          ;;*init$guard = true:bool
          i32.const 1
          global.set $runtime.init$guard

          ;;syscall/js.init()
          call $syscall$js.init

          ;;jump 2
          br $$Block_1

        end ;;$Block_1
        i32.const 2
        local.set $$current_block

        ;;return
        br $$BlockFnBody

      end ;;$Block_2
    end ;;$BlockDisp
  end ;;$BlockFnBody
) ;;runtime.init

(func $runtime.knr_getBlockHeader (param $addr i32) (result i32 i32)
  local.get $addr
  i32.load offset=0 align=4
  local.get $addr
  i32.load offset=4 align=4
) ;;runtime.knr_getBlockHeader

(func $runtime.knr_setBlockHeader (param $addr i32) (param $data.0 i32) (param $data.1 i32)
  local.get $addr
  local.get $data.0
  i32.store offset=0 align=4
  local.get $addr
  local.get $data.1
  i32.store offset=4 align=4
) ;;runtime.knr_setBlockHeader

(func $$runtime.mapImp.$$block.$$OnFree (param $ptr i32)
  local.get $ptr
  i32.load offset=0 align=1
  call $runtime.Block.Release
  local.get $ptr
  i32.const 0
  i32.store offset=0 align=1
) ;;$runtime.mapImp.$$block.$$OnFree

(func $$runtime.mapImp.$ref.underlying.$$OnFree (param $$ptr i32)
  local.get $$ptr
  i32.const 4
  call_indirect (type $$OnFree)
) ;;$runtime.mapImp.$ref.underlying.$$OnFree

(func $$runtime.mapIter.$$OnFree (param $$ptr i32)
  local.get $$ptr
  i32.const 5
  call_indirect (type $$OnFree)
) ;;$runtime.mapIter.$$OnFree

(func $runtime.makeMapIter (param $m.0 i32) (param $m.1 i32) (result i32 i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$ret_0.0 i32)
  (local $$ret_0.1 i32)
  (local $$t0.0 i32)
  (local $$t0.1 i32)
  (local $$t1.0 i32)
  (local $$t1.1 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;new mapIter (complit)
        i32.const 28
        call $runtime.HeapAlloc
        i32.const 1
        i32.const 6
        i32.const 12
        call $runtime.Block.Init
        call $runtime.DupI32
        i32.const 16
        i32.add
        local.set $$t0.1
        local.get $$t0.0
        call $runtime.Block.Release
        local.set $$t0.0

        ;;&t0.m [#0]
        local.get $$t0.0
        call $runtime.Block.Retain
        local.get $$t0.1
        i32.const 0
        i32.add
        local.set $$t1.1
        local.get $$t1.0
        call $runtime.Block.Release
        local.set $$t1.0

        ;;*t1 = m
        local.get $$t1.1
        local.get $m.0
        call $runtime.Block.Retain
        local.get $$t1.1
        i32.load offset=0 align=1
        call $runtime.Block.Release
        i32.store offset=0 align=1
        local.get $$t1.1
        local.get $m.1
        i32.store offset=4 align=4

        ;;return t0
        local.get $$t0.0
        call $runtime.Block.Retain
        local.get $$t0.1
        local.set $$ret_0.1
        local.get $$ret_0.0
        call $runtime.Block.Release
        local.set $$ret_0.0
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$ret_0.0
  call $runtime.Block.Retain
  local.get $$ret_0.1
  local.get $$ret_0.0
  call $runtime.Block.Release
  local.get $$t0.0
  call $runtime.Block.Release
  local.get $$t1.0
  call $runtime.Block.Release
) ;;runtime.makeMapIter

(func $runtime.mapDelete (param $m.0 i32) (param $m.1 i32) (param $k.0.0 i32) (param $k.0.1 i32) (param $k.1 i32) (param $k.2 i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$t0 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_2
        block $$Block_1
          block $$Block_0
            block $$BlockSel
              local.get $$block_selector
              br_table 0 1 2 0
            end ;;$BlockSel
            i32.const 0
            local.set $$current_block

            ;;m == nil:*mapImp
            local.get $m.1
            i32.const 0
            i32.eq
            local.set $$t0

            ;;if t0 goto 1 else 2
            local.get $$t0
            if
              br $$Block_0
            else
              br $$Block_1
            end

          end ;;$Block_0
          i32.const 1
          local.set $$current_block

          ;;return
          br $$BlockFnBody

        end ;;$Block_1
        i32.const 2
        local.set $$current_block

        ;;(*mapImp).Delete(m, k)
        local.get $m.0
        local.get $m.1
        local.get $k.0.0
        local.get $k.0.1
        local.get $k.1
        local.get $k.2
        call $runtime.mapImp.Delete

        ;;return
        br $$BlockFnBody

      end ;;$Block_2
    end ;;$BlockDisp
  end ;;$BlockFnBody
) ;;runtime.mapDelete

(func $runtime.mapLen (param $m.0 i32) (param $m.1 i32) (result i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$ret_0 i32)
  (local $$t0 i32)
  (local $$t1 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_2
        block $$Block_1
          block $$Block_0
            block $$BlockSel
              local.get $$block_selector
              br_table 0 1 2 0
            end ;;$BlockSel
            i32.const 0
            local.set $$current_block

            ;;m == nil:*mapImp
            local.get $m.1
            i32.const 0
            i32.eq
            local.set $$t0

            ;;if t0 goto 1 else 2
            local.get $$t0
            if
              br $$Block_0
            else
              br $$Block_1
            end

          end ;;$Block_0
          i32.const 1
          local.set $$current_block

          ;;return 0:int
          i32.const 0
          local.set $$ret_0
          br $$BlockFnBody

        end ;;$Block_1
        i32.const 2
        local.set $$current_block

        ;;(*mapImp).Len(m)
        local.get $m.0
        local.get $m.1
        call $runtime.mapImp.Len
        local.set $$t1

        ;;return t1
        local.get $$t1
        local.set $$ret_0
        br $$BlockFnBody

      end ;;$Block_2
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$ret_0
) ;;runtime.mapLen

(func $runtime.mapLookup (param $m.0 i32) (param $m.1 i32) (param $k.0.0 i32) (param $k.0.1 i32) (param $k.1 i32) (param $k.2 i32) (result i32 i32 i32 i32 i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$ret_0.0.0 i32)
  (local $$ret_0.0.1 i32)
  (local $$ret_0.1 i32)
  (local $$ret_0.2 i32)
  (local $$ret_1 i32)
  (local $$t0 i32)
  (local $$t1.0.0.0 i32)
  (local $$t1.0.0.1 i32)
  (local $$t1.0.1 i32)
  (local $$t1.0.2 i32)
  (local $$t1.1 i32)
  (local $$t2.0.0 i32)
  (local $$t2.0.1 i32)
  (local $$t2.1 i32)
  (local $$t2.2 i32)
  (local $$t3 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_2
        block $$Block_1
          block $$Block_0
            block $$BlockSel
              local.get $$block_selector
              br_table 0 1 2 0
            end ;;$BlockSel
            i32.const 0
            local.set $$current_block

            ;;m == nil:*mapImp
            local.get $m.1
            i32.const 0
            i32.eq
            local.set $$t0

            ;;if t0 goto 1 else 2
            local.get $$t0
            if
              br $$Block_0
            else
              br $$Block_1
            end

          end ;;$Block_0
          i32.const 1
          local.set $$current_block

          ;;return nil:interface{}, false:bool
          i32.const 0
          i32.const 0
          i32.const 0
          i32.const 0
          local.set $$ret_0.2
          local.set $$ret_0.1
          local.set $$ret_0.0.1
          local.get $$ret_0.0.0
          call $runtime.Block.Release
          local.set $$ret_0.0.0
          i32.const 0
          local.set $$ret_1
          br $$BlockFnBody

        end ;;$Block_1
        i32.const 2
        local.set $$current_block

        ;;(*mapImp).Lookup(m, k)
        local.get $m.0
        local.get $m.1
        local.get $k.0.0
        local.get $k.0.1
        local.get $k.1
        local.get $k.2
        call $runtime.mapImp.Lookup
        local.set $$t1.1
        local.set $$t1.0.2
        local.set $$t1.0.1
        local.set $$t1.0.0.1
        local.get $$t1.0.0.0
        call $runtime.Block.Release
        local.set $$t1.0.0.0

        ;;extract t1 #0
        local.get $$t1.0.0.0
        call $runtime.Block.Retain
        local.get $$t1.0.0.1
        local.get $$t1.0.1
        local.get $$t1.0.2
        local.set $$t2.2
        local.set $$t2.1
        local.set $$t2.0.1
        local.get $$t2.0.0
        call $runtime.Block.Release
        local.set $$t2.0.0

        ;;extract t1 #1
        local.get $$t1.1
        local.set $$t3

        ;;return t2, t3
        local.get $$t2.0.0
        call $runtime.Block.Retain
        local.get $$t2.0.1
        local.get $$t2.1
        local.get $$t2.2
        local.set $$ret_0.2
        local.set $$ret_0.1
        local.set $$ret_0.0.1
        local.get $$ret_0.0.0
        call $runtime.Block.Release
        local.set $$ret_0.0.0
        local.get $$t3
        local.set $$ret_1
        br $$BlockFnBody

      end ;;$Block_2
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$ret_0.0.0
  call $runtime.Block.Retain
  local.get $$ret_0.0.1
  local.get $$ret_0.1
  local.get $$ret_0.2
  local.get $$ret_1
  local.get $$ret_0.0.0
  call $runtime.Block.Release
  local.get $$t1.0.0.0
  call $runtime.Block.Release
  local.get $$t2.0.0
  call $runtime.Block.Release
) ;;runtime.mapLookup

(func $$runtime.mapNode.$$block.$$OnFree (param $ptr i32)
  local.get $ptr
  i32.load offset=0 align=1
  call $runtime.Block.Release
  local.get $ptr
  i32.const 0
  i32.store offset=0 align=1
) ;;$runtime.mapNode.$$block.$$OnFree

(func $$runtime.mapNode.$ref.underlying.$$OnFree (param $$ptr i32)
  local.get $$ptr
  i32.const 7
  call_indirect (type $$OnFree)
) ;;$runtime.mapNode.$ref.underlying.$$OnFree

(func $$void.$$block.$$OnFree (param $ptr i32)
  local.get $ptr
  i32.load offset=0 align=1
  call $runtime.Block.Release
  local.get $ptr
  i32.const 0
  i32.store offset=0 align=1
) ;;$void.$$block.$$OnFree

(func $$void.$ref.underlying.$$OnFree (param $$ptr i32)
  local.get $$ptr
  i32.const 9
  call_indirect (type $$OnFree)
) ;;$void.$ref.underlying.$$OnFree

(func $$i`0`.underlying.$$OnFree (param $$ptr i32)
  local.get $$ptr
  i32.const 10
  call_indirect (type $$OnFree)
) ;;$i`0`.underlying.$$OnFree

(func $$runtime.mapNode.$$OnFree (param $$ptr i32)
  local.get $$ptr
  i32.const 8
  i32.add
  i32.const 8
  call_indirect (type $$OnFree)
  local.get $$ptr
  i32.const 16
  i32.add
  i32.const 8
  call_indirect (type $$OnFree)
  local.get $$ptr
  i32.const 28
  i32.add
  i32.const 11
  call_indirect (type $$OnFree)
  local.get $$ptr
  i32.const 44
  i32.add
  i32.const 11
  call_indirect (type $$OnFree)
) ;;$runtime.mapNode.$$OnFree

(func $$runtime.mapNode.$ref.$$block.$$OnFree (param $ptr i32)
  local.get $ptr
  i32.load offset=0 align=1
  call $runtime.Block.Release
  local.get $ptr
  i32.const 0
  i32.store offset=0 align=1
) ;;$runtime.mapNode.$ref.$$block.$$OnFree

(func $$runtime.mapNode.$ref.$slice.underlying.$$OnFree (param $$ptr i32)
  local.get $$ptr
  i32.const 13
  call_indirect (type $$OnFree)
) ;;$runtime.mapNode.$ref.$slice.underlying.$$OnFree

(func $$runtime.mapImp.$$OnFree (param $$ptr i32)
  local.get $$ptr
  i32.const 8
  call_indirect (type $$OnFree)
  local.get $$ptr
  i32.const 8
  i32.add
  i32.const 8
  call_indirect (type $$OnFree)
  local.get $$ptr
  i32.const 16
  i32.add
  i32.const 14
  call_indirect (type $$OnFree)
) ;;$runtime.mapImp.$$OnFree

(func $$runtime.mapNode.$ref.$array1.underlying.$$OnFree (param $$ptr i32)
  local.get $$ptr
  i32.const 8
  call_indirect (type $$OnFree)
) ;;$runtime.mapNode.$ref.$array1.underlying.$$OnFree

(func $runtime.mapMake (result i32 i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$ret_0.0 i32)
  (local $$ret_0.1 i32)
  (local $$t0.0 i32)
  (local $$t0.1 i32)
  (local $$t1.0 i32)
  (local $$t1.1 i32)
  (local $$t2.0 i32)
  (local $$t2.1 i32)
  (local $$t3.0 i32)
  (local $$t3.1 i32)
  (local $$t4.0 i32)
  (local $$t4.1 i32)
  (local $$t5.0 i32)
  (local $$t5.1 i32)
  (local $$t6.0 i32)
  (local $$t6.1 i32)
  (local $$t7.0 i32)
  (local $$t7.1 i32)
  (local $$t8.0 i32)
  (local $$t8.1 i32)
  (local $$t8.2 i32)
  (local $$t8.3 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;new mapNode (complit)
        i32.const 76
        call $runtime.HeapAlloc
        i32.const 1
        i32.const 12
        i32.const 60
        call $runtime.Block.Init
        call $runtime.DupI32
        i32.const 16
        i32.add
        local.set $$t0.1
        local.get $$t0.0
        call $runtime.Block.Release
        local.set $$t0.0

        ;;&t0.Color [#4]
        local.get $$t0.0
        call $runtime.Block.Retain
        local.get $$t0.1
        i32.const 24
        i32.add
        local.set $$t1.1
        local.get $$t1.0
        call $runtime.Block.Release
        local.set $$t1.0

        ;;*t1 = 1:mapColor
        local.get $$t1.1
        i32.const 1
        i32.store offset=0 align=4

        ;;new mapImp (complit)
        i32.const 48
        call $runtime.HeapAlloc
        i32.const 1
        i32.const 15
        i32.const 32
        call $runtime.Block.Init
        call $runtime.DupI32
        i32.const 16
        i32.add
        local.set $$t2.1
        local.get $$t2.0
        call $runtime.Block.Release
        local.set $$t2.0

        ;;&t2.NIL [#0]
        local.get $$t2.0
        call $runtime.Block.Retain
        local.get $$t2.1
        i32.const 0
        i32.add
        local.set $$t3.1
        local.get $$t3.0
        call $runtime.Block.Release
        local.set $$t3.0

        ;;&t2.root [#1]
        local.get $$t2.0
        call $runtime.Block.Retain
        local.get $$t2.1
        i32.const 8
        i32.add
        local.set $$t4.1
        local.get $$t4.0
        call $runtime.Block.Release
        local.set $$t4.0

        ;;&t2.nodes [#2]
        local.get $$t2.0
        call $runtime.Block.Retain
        local.get $$t2.1
        i32.const 16
        i32.add
        local.set $$t5.1
        local.get $$t5.0
        call $runtime.Block.Release
        local.set $$t5.0

        ;;new [1]*mapNode (slicelit)
        i32.const 24
        call $runtime.HeapAlloc
        i32.const 1
        i32.const 16
        i32.const 8
        call $runtime.Block.Init
        call $runtime.DupI32
        i32.const 16
        i32.add
        local.set $$t6.1
        local.get $$t6.0
        call $runtime.Block.Release
        local.set $$t6.0

        ;;&t6[0:int]
        local.get $$t6.0
        call $runtime.Block.Retain
        local.get $$t6.1
        i32.const 8
        i32.const 0
        i32.mul
        i32.add
        local.set $$t7.1
        local.get $$t7.0
        call $runtime.Block.Release
        local.set $$t7.0

        ;;*t7 = t0
        local.get $$t7.1
        local.get $$t0.0
        call $runtime.Block.Retain
        local.get $$t7.1
        i32.load offset=0 align=1
        call $runtime.Block.Release
        i32.store offset=0 align=1
        local.get $$t7.1
        local.get $$t0.1
        i32.store offset=4 align=4

        ;;slice t6[:]
        local.get $$t6.0
        call $runtime.Block.Retain
        local.get $$t6.1
        i32.const 8
        i32.const 0
        i32.mul
        i32.add
        i32.const 1
        i32.const 0
        i32.sub
        i32.const 1
        i32.const 0
        i32.sub
        local.set $$t8.3
        local.set $$t8.2
        local.set $$t8.1
        local.get $$t8.0
        call $runtime.Block.Release
        local.set $$t8.0

        ;;*t3 = t0
        local.get $$t3.1
        local.get $$t0.0
        call $runtime.Block.Retain
        local.get $$t3.1
        i32.load offset=0 align=1
        call $runtime.Block.Release
        i32.store offset=0 align=1
        local.get $$t3.1
        local.get $$t0.1
        i32.store offset=4 align=4

        ;;*t4 = t0
        local.get $$t4.1
        local.get $$t0.0
        call $runtime.Block.Retain
        local.get $$t4.1
        i32.load offset=0 align=1
        call $runtime.Block.Release
        i32.store offset=0 align=1
        local.get $$t4.1
        local.get $$t0.1
        i32.store offset=4 align=4

        ;;*t5 = t8
        local.get $$t5.1
        local.get $$t8.0
        call $runtime.Block.Retain
        local.get $$t5.1
        i32.load offset=0 align=1
        call $runtime.Block.Release
        i32.store offset=0 align=1
        local.get $$t5.1
        local.get $$t8.1
        i32.store offset=4 align=4
        local.get $$t5.1
        local.get $$t8.2
        i32.store offset=8 align=4
        local.get $$t5.1
        local.get $$t8.3
        i32.store offset=12 align=4

        ;;return t2
        local.get $$t2.0
        call $runtime.Block.Retain
        local.get $$t2.1
        local.set $$ret_0.1
        local.get $$ret_0.0
        call $runtime.Block.Release
        local.set $$ret_0.0
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$ret_0.0
  call $runtime.Block.Retain
  local.get $$ret_0.1
  local.get $$ret_0.0
  call $runtime.Block.Release
  local.get $$t0.0
  call $runtime.Block.Release
  local.get $$t1.0
  call $runtime.Block.Release
  local.get $$t2.0
  call $runtime.Block.Release
  local.get $$t3.0
  call $runtime.Block.Release
  local.get $$t4.0
  call $runtime.Block.Release
  local.get $$t5.0
  call $runtime.Block.Release
  local.get $$t6.0
  call $runtime.Block.Release
  local.get $$t7.0
  call $runtime.Block.Release
  local.get $$t8.0
  call $runtime.Block.Release
) ;;runtime.mapMake

(func $runtime.mapNext (param $iter.0.0 i32) (param $iter.0.1 i32) (param $iter.1 i32) (result i32 i32 i32 i32 i32 i32 i32 i32 i32 i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$ret_0 i32)
  (local $$ret_1.0.0 i32)
  (local $$ret_1.0.1 i32)
  (local $$ret_1.1 i32)
  (local $$ret_1.2 i32)
  (local $$ret_2.0.0 i32)
  (local $$ret_2.0.1 i32)
  (local $$ret_2.1 i32)
  (local $$ret_2.2 i32)
  (local $$ret_3 i32)
  (local $$t0.0 i32)
  (local $$t0.1 i32)
  (local $$t1.0 i32)
  (local $$t1.1 i32)
  (local $$t2.0 i32)
  (local $$t2.1 i32)
  (local $$t3 i32)
  (local $$t4.0 i32)
  (local $$t4.1 i32)
  (local $$t5 i32)
  (local $$t6.0 i32)
  (local $$t6.1.0.0 i32)
  (local $$t6.1.0.1 i32)
  (local $$t6.1.1 i32)
  (local $$t6.1.2 i32)
  (local $$t6.2.0.0 i32)
  (local $$t6.2.0.1 i32)
  (local $$t6.2.1 i32)
  (local $$t6.2.2 i32)
  (local $$t7 i32)
  (local $$t8.0.0 i32)
  (local $$t8.0.1 i32)
  (local $$t8.1 i32)
  (local $$t8.2 i32)
  (local $$t9.0.0 i32)
  (local $$t9.0.1 i32)
  (local $$t9.1 i32)
  (local $$t9.2 i32)
  (local $$t10.0 i32)
  (local $$t10.1 i32)
  (local $$t11 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_2
        block $$Block_1
          block $$Block_0
            block $$BlockSel
              local.get $$block_selector
              br_table 0 1 2 0
            end ;;$BlockSel
            i32.const 0
            local.set $$current_block

            ;;new mapIter (iter)
            i32.const 28
            call $runtime.HeapAlloc
            i32.const 1
            i32.const 6
            i32.const 12
            call $runtime.Block.Init
            call $runtime.DupI32
            i32.const 16
            i32.add
            local.set $$t0.1
            local.get $$t0.0
            call $runtime.Block.Release
            local.set $$t0.0

            ;;*t0 = iter
            local.get $$t0.1
            local.get $iter.0.0
            call $runtime.Block.Retain
            local.get $$t0.1
            i32.load offset=0 align=1
            call $runtime.Block.Release
            i32.store offset=0 align=1
            local.get $$t0.1
            local.get $iter.0.1
            i32.store offset=4 align=4
            local.get $$t0.1
            local.get $iter.1
            i32.store offset=8 align=4

            ;;&t0.m [#0]
            local.get $$t0.0
            call $runtime.Block.Retain
            local.get $$t0.1
            i32.const 0
            i32.add
            local.set $$t1.1
            local.get $$t1.0
            call $runtime.Block.Release
            local.set $$t1.0

            ;;*t1
            local.get $$t1.1
            i32.load offset=0 align=4
            call $runtime.Block.Retain
            local.get $$t1.1
            i32.load offset=4 align=4
            local.set $$t2.1
            local.get $$t2.0
            call $runtime.Block.Release
            local.set $$t2.0

            ;;t2 == nil:*mapImp
            local.get $$t2.1
            i32.const 0
            i32.eq
            local.set $$t3

            ;;if t3 goto 1 else 2
            local.get $$t3
            if
              br $$Block_0
            else
              br $$Block_1
            end

          end ;;$Block_0
          i32.const 1
          local.set $$current_block

          ;;&t0.pos [#1]
          local.get $$t0.0
          call $runtime.Block.Retain
          local.get $$t0.1
          i32.const 8
          i32.add
          local.set $$t4.1
          local.get $$t4.0
          call $runtime.Block.Release
          local.set $$t4.0

          ;;*t4
          local.get $$t4.1
          i32.load offset=0 align=4
          local.set $$t5

          ;;return false:bool, nil:interface{}, nil:interface{}, t5
          i32.const 0
          local.set $$ret_0
          i32.const 0
          i32.const 0
          i32.const 0
          i32.const 0
          local.set $$ret_1.2
          local.set $$ret_1.1
          local.set $$ret_1.0.1
          local.get $$ret_1.0.0
          call $runtime.Block.Release
          local.set $$ret_1.0.0
          i32.const 0
          i32.const 0
          i32.const 0
          i32.const 0
          local.set $$ret_2.2
          local.set $$ret_2.1
          local.set $$ret_2.0.1
          local.get $$ret_2.0.0
          call $runtime.Block.Release
          local.set $$ret_2.0.0
          local.get $$t5
          local.set $$ret_3
          br $$BlockFnBody

        end ;;$Block_1
        i32.const 2
        local.set $$current_block

        ;;(*mapIter).Next(t0)
        local.get $$t0.0
        local.get $$t0.1
        call $runtime.mapIter.Next
        local.set $$t6.2.2
        local.set $$t6.2.1
        local.set $$t6.2.0.1
        local.get $$t6.2.0.0
        call $runtime.Block.Release
        local.set $$t6.2.0.0
        local.set $$t6.1.2
        local.set $$t6.1.1
        local.set $$t6.1.0.1
        local.get $$t6.1.0.0
        call $runtime.Block.Release
        local.set $$t6.1.0.0
        local.set $$t6.0

        ;;extract t6 #0
        local.get $$t6.0
        local.set $$t7

        ;;extract t6 #1
        local.get $$t6.1.0.0
        call $runtime.Block.Retain
        local.get $$t6.1.0.1
        local.get $$t6.1.1
        local.get $$t6.1.2
        local.set $$t8.2
        local.set $$t8.1
        local.set $$t8.0.1
        local.get $$t8.0.0
        call $runtime.Block.Release
        local.set $$t8.0.0

        ;;extract t6 #2
        local.get $$t6.2.0.0
        call $runtime.Block.Retain
        local.get $$t6.2.0.1
        local.get $$t6.2.1
        local.get $$t6.2.2
        local.set $$t9.2
        local.set $$t9.1
        local.set $$t9.0.1
        local.get $$t9.0.0
        call $runtime.Block.Release
        local.set $$t9.0.0

        ;;&t0.pos [#1]
        local.get $$t0.0
        call $runtime.Block.Retain
        local.get $$t0.1
        i32.const 8
        i32.add
        local.set $$t10.1
        local.get $$t10.0
        call $runtime.Block.Release
        local.set $$t10.0

        ;;*t10
        local.get $$t10.1
        i32.load offset=0 align=4
        local.set $$t11

        ;;return t7, t8, t9, t11
        local.get $$t7
        local.set $$ret_0
        local.get $$t8.0.0
        call $runtime.Block.Retain
        local.get $$t8.0.1
        local.get $$t8.1
        local.get $$t8.2
        local.set $$ret_1.2
        local.set $$ret_1.1
        local.set $$ret_1.0.1
        local.get $$ret_1.0.0
        call $runtime.Block.Release
        local.set $$ret_1.0.0
        local.get $$t9.0.0
        call $runtime.Block.Retain
        local.get $$t9.0.1
        local.get $$t9.1
        local.get $$t9.2
        local.set $$ret_2.2
        local.set $$ret_2.1
        local.set $$ret_2.0.1
        local.get $$ret_2.0.0
        call $runtime.Block.Release
        local.set $$ret_2.0.0
        local.get $$t11
        local.set $$ret_3
        br $$BlockFnBody

      end ;;$Block_2
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$ret_0
  local.get $$ret_1.0.0
  call $runtime.Block.Retain
  local.get $$ret_1.0.1
  local.get $$ret_1.1
  local.get $$ret_1.2
  local.get $$ret_2.0.0
  call $runtime.Block.Retain
  local.get $$ret_2.0.1
  local.get $$ret_2.1
  local.get $$ret_2.2
  local.get $$ret_3
  local.get $$ret_1.0.0
  call $runtime.Block.Release
  local.get $$ret_2.0.0
  call $runtime.Block.Release
  local.get $$t0.0
  call $runtime.Block.Release
  local.get $$t1.0
  call $runtime.Block.Release
  local.get $$t2.0
  call $runtime.Block.Release
  local.get $$t4.0
  call $runtime.Block.Release
  local.get $$t6.2.0.0
  call $runtime.Block.Release
  local.get $$t6.1.0.0
  call $runtime.Block.Release
  local.get $$t8.0.0
  call $runtime.Block.Release
  local.get $$t9.0.0
  call $runtime.Block.Release
  local.get $$t10.0
  call $runtime.Block.Release
) ;;runtime.mapNext

(func $runtime.mapUpdate (param $m.0 i32) (param $m.1 i32) (param $k.0.0 i32) (param $k.0.1 i32) (param $k.1 i32) (param $k.2 i32) (param $v.0.0 i32) (param $v.0.1 i32) (param $v.1 i32) (param $v.2 i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$t0 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_2
        block $$Block_1
          block $$Block_0
            block $$BlockSel
              local.get $$block_selector
              br_table 0 1 2 0
            end ;;$BlockSel
            i32.const 0
            local.set $$current_block

            ;;m == nil:*mapImp
            local.get $m.1
            i32.const 0
            i32.eq
            local.set $$t0

            ;;if t0 goto 1 else 2
            local.get $$t0
            if
              br $$Block_0
            else
              br $$Block_1
            end

          end ;;$Block_0
          i32.const 1
          local.set $$current_block

          ;;panic "nil map.":string
          i32.const 8282
          i32.const 8
          i32.const 8290
          i32.const 11
          call $$runtime.panic_

        end ;;$Block_1
        i32.const 2
        local.set $$current_block

        ;;(*mapImp).Update(m, k, v)
        local.get $m.0
        local.get $m.1
        local.get $k.0.0
        local.get $k.0.1
        local.get $k.1
        local.get $k.2
        local.get $v.0.0
        local.get $v.0.1
        local.get $v.1
        local.get $v.2
        call $runtime.mapImp.Update

        ;;return
        br $$BlockFnBody

      end ;;$Block_2
    end ;;$BlockDisp
  end ;;$BlockFnBody
) ;;runtime.mapUpdate

(func $runtime.next_rune (param $iter.0 i32) (param $iter.1 i32) (param $iter.2 i32) (result i32 i32 i32 i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$ret_0 i32)
  (local $$ret_1 i32)
  (local $$ret_2 i32)
  (local $$ret_3 i32)
  (local $$t0.0 i32)
  (local $$t0.1 i32)
  (local $$t1.0 i32)
  (local $$t1.1 i32)
  (local $$t2 i32)
  (local $$t3.0 i32)
  (local $$t3.1 i32)
  (local $$t4 i32)
  (local $$t5 i32)
  (local $$t6.0 i32)
  (local $$t6.1 i32)
  (local $$t7 i32)
  (local $$t8.0 i32)
  (local $$t8.1 i32)
  (local $$t9 i32)
  (local $$t10.0 i32)
  (local $$t10.1 i32)
  (local $$t11 i32)
  (local $$t12 i32)
  (local $$t13.0 i32)
  (local $$t13.1 i32)
  (local $$t14 i32)
  (local $$t15 i32)
  (local $$t16 i32)
  (local $$t17 i32)
  (local $$t18 i32)
  (local $$t19 i32)
  (local $$t20 i32)
  (local $$t21.0 i32)
  (local $$t21.1 i32)
  (local $$t22 i32)
  (local $$t23.0 i32)
  (local $$t23.1 i32)
  (local $$t24 i32)
  (local $$t25 i32)
  (local $$t26 i32)
  (local $$t27 i32)
  (local $$t28 i32)
  (local $$t29 i32)
  (local $$t30.0 i32)
  (local $$t30.1 i32)
  (local $$t31 i32)
  (local $$t32 i32)
  (local $$t33.0 i32)
  (local $$t33.1 i32)
  (local $$t34 i32)
  (local $$t35 i32)
  (local $$t36 i32)
  (local $$t37 i32)
  (local $$t38 i32)
  (local $$t39 i32)
  (local $$t40 i32)
  (local $$t41.0 i32)
  (local $$t41.1 i32)
  (local $$t42 i32)
  (local $$t43 i32)
  (local $$t44.0 i32)
  (local $$t44.1 i32)
  (local $$t45 i32)
  (local $$t46 i32)
  (local $$t47 i32)
  (local $$t48 i32)
  (local $$t49 i32)
  (local $$t50 i32)
  (local $$t51.0 i32)
  (local $$t51.1 i32)
  (local $$t52 i32)
  (local $$t53 i32)
  (local $$t54.0 i32)
  (local $$t54.1 i32)
  (local $$t55 i32)
  (local $$t56 i32)
  (local $$t57 i32)
  (local $$t58 i32)
  (local $$t59 i32)
  (local $$t60 i32)
  (local $$t61 i32)
  (local $$t62 i32)
  (local $$t63.0 i32)
  (local $$t63.1 i32)
  (local $$t64 i32)
  (local $$t65 i32)
  (local $$t66.0 i32)
  (local $$t66.1 i32)
  (local $$t67 i32)
  (local $$t68 i32)
  (local $$t69 i32)
  (local $$t70 i32)
  (local $$t71 i32)
  (local $$t72 i32)
  (local $$t73 i32)
  (local $$t74.0 i32)
  (local $$t74.1 i32)
  (local $$t75 i32)
  (local $$t76 i32)
  (local $$t77 i32)
  (local $$t78.0 i32)
  (local $$t78.1 i32)
  (local $$t79 i32)
  (local $$t80 i32)
  (local $$t81 i32)
  (local $$t82 i32)
  (local $$t83 i32)
  (local $$t84 i32)
  (local $$t85.0 i32)
  (local $$t85.1 i32)
  (local $$t86 i32)
  (local $$t87 i32)
  (local $$t88.0 i32)
  (local $$t88.1 i32)
  (local $$t89 i32)
  (local $$t90 i32)
  (local $$t91 i32)
  (local $$t92 i32)
  (local $$t93 i32)
  (local $$t94 i32)
  (local $$t95 i32)
  (local $$t96 i32)
  (local $$t97.0 i32)
  (local $$t97.1 i32)
  (local $$t98 i32)
  (local $$t99 i32)
  (local $$t100.0 i32)
  (local $$t100.1 i32)
  (local $$t101 i32)
  (local $$t102 i32)
  (local $$t103 i32)
  (local $$t104 i32)
  (local $$t105 i32)
  (local $$t106 i32)
  (local $$t107 i32)
  (local $$t108 i32)
  (local $$t109.0 i32)
  (local $$t109.1 i32)
  (local $$t110 i32)
  (local $$t111 i32)
  (local $$t112.0 i32)
  (local $$t112.1 i32)
  (local $$t113 i32)
  (local $$t114 i32)
  (local $$t115 i32)
  (local $$t116 i32)
  (local $$t117 i32)
  (local $$t118 i32)
  (local $$t119 i32)
  (local $$t120.0 i32)
  (local $$t120.1 i32)
  (local $$t121 i32)
  (local $$t122 i32)
  (local $$t123 i32)
  (local $$t124 i32)
  (local $$t125.0 i32)
  (local $$t125.1 i32)
  (local $$t126 i32)
  (local $$t127 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_10
        block $$Block_9
          block $$Block_8
            block $$Block_7
              block $$Block_6
                block $$Block_5
                  block $$Block_4
                    block $$Block_3
                      block $$Block_2
                        block $$Block_1
                          block $$Block_0
                            block $$BlockSel
                              local.get $$block_selector
                              br_table 0 1 2 3 4 5 6 7 8 9 10 0
                            end ;;$BlockSel
                            i32.const 0
                            local.set $$current_block

                            ;;local stringIter (iter)
                            i32.const 28
                            call $runtime.HeapAlloc
                            i32.const 1
                            i32.const 0
                            i32.const 12
                            call $runtime.Block.Init
                            call $runtime.DupI32
                            i32.const 16
                            i32.add
                            local.set $$t0.1
                            local.get $$t0.0
                            call $runtime.Block.Release
                            local.set $$t0.0

                            ;;*t0 = iter
                            local.get $$t0.1
                            local.get $iter.0
                            i32.store offset=0 align=4
                            local.get $$t0.1
                            local.get $iter.1
                            i32.store offset=4 align=4
                            local.get $$t0.1
                            local.get $iter.2
                            i32.store offset=8 align=4

                            ;;&t0.pos [#2]
                            local.get $$t0.0
                            call $runtime.Block.Retain
                            local.get $$t0.1
                            i32.const 8
                            i32.add
                            local.set $$t1.1
                            local.get $$t1.0
                            call $runtime.Block.Release
                            local.set $$t1.0

                            ;;*t1
                            local.get $$t1.1
                            i32.load offset=0 align=4
                            local.set $$t2

                            ;;&t0.len [#1]
                            local.get $$t0.0
                            call $runtime.Block.Retain
                            local.get $$t0.1
                            i32.const 4
                            i32.add
                            local.set $$t3.1
                            local.get $$t3.0
                            call $runtime.Block.Release
                            local.set $$t3.0

                            ;;*t3
                            local.get $$t3.1
                            i32.load offset=0 align=4
                            local.set $$t4

                            ;;t2 >= t4
                            local.get $$t2
                            local.get $$t4
                            i32.ge_s
                            local.set $$t5

                            ;;if t5 goto 1 else 2
                            local.get $$t5
                            if
                              br $$Block_0
                            else
                              br $$Block_1
                            end

                          end ;;$Block_0
                          i32.const 1
                          local.set $$current_block

                          ;;&t0.pos [#2]
                          local.get $$t0.0
                          call $runtime.Block.Retain
                          local.get $$t0.1
                          i32.const 8
                          i32.add
                          local.set $$t6.1
                          local.get $$t6.0
                          call $runtime.Block.Release
                          local.set $$t6.0

                          ;;*t6
                          local.get $$t6.1
                          i32.load offset=0 align=4
                          local.set $$t7

                          ;;&t0.pos [#2]
                          local.get $$t0.0
                          call $runtime.Block.Retain
                          local.get $$t0.1
                          i32.const 8
                          i32.add
                          local.set $$t8.1
                          local.get $$t8.0
                          call $runtime.Block.Release
                          local.set $$t8.0

                          ;;*t8
                          local.get $$t8.1
                          i32.load offset=0 align=4
                          local.set $$t9

                          ;;return false:bool, t7, 0:rune, t9
                          i32.const 0
                          local.set $$ret_0
                          local.get $$t7
                          local.set $$ret_1
                          i32.const 0
                          local.set $$ret_2
                          local.get $$t9
                          local.set $$ret_3
                          br $$BlockFnBody

                        end ;;$Block_1
                        i32.const 2
                        local.set $$current_block

                        ;;&t0.ptr [#0]
                        local.get $$t0.0
                        call $runtime.Block.Retain
                        local.get $$t0.1
                        i32.const 0
                        i32.add
                        local.set $$t10.1
                        local.get $$t10.0
                        call $runtime.Block.Release
                        local.set $$t10.0

                        ;;*t10
                        local.get $$t10.1
                        i32.load offset=0 align=4
                        local.set $$t11

                        ;;convert u32 <- uint (t11)
                        local.get $$t11
                        local.set $$t12

                        ;;&t0.pos [#2]
                        local.get $$t0.0
                        call $runtime.Block.Retain
                        local.get $$t0.1
                        i32.const 8
                        i32.add
                        local.set $$t13.1
                        local.get $$t13.0
                        call $runtime.Block.Release
                        local.set $$t13.0

                        ;;*t13
                        local.get $$t13.1
                        i32.load offset=0 align=4
                        local.set $$t14

                        ;;convert u32 <- int (t14)
                        local.get $$t14
                        local.set $$t15

                        ;;t12 + t15
                        local.get $$t12
                        local.get $$t15
                        i32.add
                        local.set $$t16

                        ;;get_u8(t16)
                        local.get $$t16
                        call $runtime.get_u8
                        local.set $$t17

                        ;;convert i32 <- u8 (t17)
                        local.get $$t17
                        local.set $$t18

                        ;;t18 & 128:i32
                        local.get $$t18
                        i32.const 128
                        i32.and
                        local.set $$t19

                        ;;t19 == 0:i32
                        local.get $$t19
                        i32.const 0
                        i32.eq
                        local.set $$t20

                        ;;if t20 goto 3 else 4
                        local.get $$t20
                        if
                          br $$Block_2
                        else
                          br $$Block_3
                        end

                      end ;;$Block_2
                      i32.const 3
                      local.set $$current_block

                      ;;&t0.pos [#2]
                      local.get $$t0.0
                      call $runtime.Block.Retain
                      local.get $$t0.1
                      i32.const 8
                      i32.add
                      local.set $$t21.1
                      local.get $$t21.0
                      call $runtime.Block.Release
                      local.set $$t21.0

                      ;;*t21
                      local.get $$t21.1
                      i32.load offset=0 align=4
                      local.set $$t22

                      ;;&t0.pos [#2]
                      local.get $$t0.0
                      call $runtime.Block.Retain
                      local.get $$t0.1
                      i32.const 8
                      i32.add
                      local.set $$t23.1
                      local.get $$t23.0
                      call $runtime.Block.Release
                      local.set $$t23.0

                      ;;*t23
                      local.get $$t23.1
                      i32.load offset=0 align=4
                      local.set $$t24

                      ;;t24 + 1:int
                      local.get $$t24
                      i32.const 1
                      i32.add
                      local.set $$t25

                      ;;return true:bool, t22, t18, t25
                      i32.const 1
                      local.set $$ret_0
                      local.get $$t22
                      local.set $$ret_1
                      local.get $$t18
                      local.set $$ret_2
                      local.get $$t25
                      local.set $$ret_3
                      br $$BlockFnBody

                    end ;;$Block_3
                    i32.const 4
                    local.set $$current_block

                    ;;t18 & 224:i32
                    local.get $$t18
                    i32.const 224
                    i32.and
                    local.set $$t26

                    ;;t26 == 192:i32
                    local.get $$t26
                    i32.const 192
                    i32.eq
                    local.set $$t27

                    ;;if t27 goto 5 else 6
                    local.get $$t27
                    if
                      br $$Block_4
                    else
                      br $$Block_5
                    end

                  end ;;$Block_4
                  i32.const 5
                  local.set $$current_block

                  ;;t18 & 31:i32
                  local.get $$t18
                  i32.const 31
                  i32.and
                  local.set $$t28

                  ;;t28 << 6:uint64
                  local.get $$t28
                  i64.const 6
                  i32.wrap_i64
                  i32.shl
                  local.set $$t29

                  ;;&t0.ptr [#0]
                  local.get $$t0.0
                  call $runtime.Block.Retain
                  local.get $$t0.1
                  i32.const 0
                  i32.add
                  local.set $$t30.1
                  local.get $$t30.0
                  call $runtime.Block.Release
                  local.set $$t30.0

                  ;;*t30
                  local.get $$t30.1
                  i32.load offset=0 align=4
                  local.set $$t31

                  ;;convert u32 <- uint (t31)
                  local.get $$t31
                  local.set $$t32

                  ;;&t0.pos [#2]
                  local.get $$t0.0
                  call $runtime.Block.Retain
                  local.get $$t0.1
                  i32.const 8
                  i32.add
                  local.set $$t33.1
                  local.get $$t33.0
                  call $runtime.Block.Release
                  local.set $$t33.0

                  ;;*t33
                  local.get $$t33.1
                  i32.load offset=0 align=4
                  local.set $$t34

                  ;;convert u32 <- int (t34)
                  local.get $$t34
                  local.set $$t35

                  ;;t32 + t35
                  local.get $$t32
                  local.get $$t35
                  i32.add
                  local.set $$t36

                  ;;t36 + 1:u32
                  local.get $$t36
                  i32.const 1
                  i32.add
                  local.set $$t37

                  ;;get_u8(t37)
                  local.get $$t37
                  call $runtime.get_u8
                  local.set $$t38

                  ;;convert i32 <- u8 (t38)
                  local.get $$t38
                  local.set $$t39

                  ;;t39 & 63:i32
                  local.get $$t39
                  i32.const 63
                  i32.and
                  local.set $$t40

                  ;;&t0.pos [#2]
                  local.get $$t0.0
                  call $runtime.Block.Retain
                  local.get $$t0.1
                  i32.const 8
                  i32.add
                  local.set $$t41.1
                  local.get $$t41.0
                  call $runtime.Block.Release
                  local.set $$t41.0

                  ;;*t41
                  local.get $$t41.1
                  i32.load offset=0 align=4
                  local.set $$t42

                  ;;t29 | t40
                  local.get $$t29
                  local.get $$t40
                  i32.or
                  local.set $$t43

                  ;;&t0.pos [#2]
                  local.get $$t0.0
                  call $runtime.Block.Retain
                  local.get $$t0.1
                  i32.const 8
                  i32.add
                  local.set $$t44.1
                  local.get $$t44.0
                  call $runtime.Block.Release
                  local.set $$t44.0

                  ;;*t44
                  local.get $$t44.1
                  i32.load offset=0 align=4
                  local.set $$t45

                  ;;t45 + 2:int
                  local.get $$t45
                  i32.const 2
                  i32.add
                  local.set $$t46

                  ;;return true:bool, t42, t43, t46
                  i32.const 1
                  local.set $$ret_0
                  local.get $$t42
                  local.set $$ret_1
                  local.get $$t43
                  local.set $$ret_2
                  local.get $$t46
                  local.set $$ret_3
                  br $$BlockFnBody

                end ;;$Block_5
                i32.const 6
                local.set $$current_block

                ;;t18 & 240:i32
                local.get $$t18
                i32.const 240
                i32.and
                local.set $$t47

                ;;t47 == 224:i32
                local.get $$t47
                i32.const 224
                i32.eq
                local.set $$t48

                ;;if t48 goto 7 else 8
                local.get $$t48
                if
                  br $$Block_6
                else
                  br $$Block_7
                end

              end ;;$Block_6
              i32.const 7
              local.set $$current_block

              ;;t18 & 15:i32
              local.get $$t18
              i32.const 15
              i32.and
              local.set $$t49

              ;;t49 << 12:uint64
              local.get $$t49
              i64.const 12
              i32.wrap_i64
              i32.shl
              local.set $$t50

              ;;&t0.ptr [#0]
              local.get $$t0.0
              call $runtime.Block.Retain
              local.get $$t0.1
              i32.const 0
              i32.add
              local.set $$t51.1
              local.get $$t51.0
              call $runtime.Block.Release
              local.set $$t51.0

              ;;*t51
              local.get $$t51.1
              i32.load offset=0 align=4
              local.set $$t52

              ;;convert u32 <- uint (t52)
              local.get $$t52
              local.set $$t53

              ;;&t0.pos [#2]
              local.get $$t0.0
              call $runtime.Block.Retain
              local.get $$t0.1
              i32.const 8
              i32.add
              local.set $$t54.1
              local.get $$t54.0
              call $runtime.Block.Release
              local.set $$t54.0

              ;;*t54
              local.get $$t54.1
              i32.load offset=0 align=4
              local.set $$t55

              ;;convert u32 <- int (t55)
              local.get $$t55
              local.set $$t56

              ;;t53 + t56
              local.get $$t53
              local.get $$t56
              i32.add
              local.set $$t57

              ;;t57 + 1:u32
              local.get $$t57
              i32.const 1
              i32.add
              local.set $$t58

              ;;get_u8(t58)
              local.get $$t58
              call $runtime.get_u8
              local.set $$t59

              ;;convert i32 <- u8 (t59)
              local.get $$t59
              local.set $$t60

              ;;t60 & 63:i32
              local.get $$t60
              i32.const 63
              i32.and
              local.set $$t61

              ;;t61 << 6:uint64
              local.get $$t61
              i64.const 6
              i32.wrap_i64
              i32.shl
              local.set $$t62

              ;;&t0.ptr [#0]
              local.get $$t0.0
              call $runtime.Block.Retain
              local.get $$t0.1
              i32.const 0
              i32.add
              local.set $$t63.1
              local.get $$t63.0
              call $runtime.Block.Release
              local.set $$t63.0

              ;;*t63
              local.get $$t63.1
              i32.load offset=0 align=4
              local.set $$t64

              ;;convert u32 <- uint (t64)
              local.get $$t64
              local.set $$t65

              ;;&t0.pos [#2]
              local.get $$t0.0
              call $runtime.Block.Retain
              local.get $$t0.1
              i32.const 8
              i32.add
              local.set $$t66.1
              local.get $$t66.0
              call $runtime.Block.Release
              local.set $$t66.0

              ;;*t66
              local.get $$t66.1
              i32.load offset=0 align=4
              local.set $$t67

              ;;convert u32 <- int (t67)
              local.get $$t67
              local.set $$t68

              ;;t65 + t68
              local.get $$t65
              local.get $$t68
              i32.add
              local.set $$t69

              ;;t69 + 2:u32
              local.get $$t69
              i32.const 2
              i32.add
              local.set $$t70

              ;;get_u8(t70)
              local.get $$t70
              call $runtime.get_u8
              local.set $$t71

              ;;convert i32 <- u8 (t71)
              local.get $$t71
              local.set $$t72

              ;;t72 & 63:i32
              local.get $$t72
              i32.const 63
              i32.and
              local.set $$t73

              ;;&t0.pos [#2]
              local.get $$t0.0
              call $runtime.Block.Retain
              local.get $$t0.1
              i32.const 8
              i32.add
              local.set $$t74.1
              local.get $$t74.0
              call $runtime.Block.Release
              local.set $$t74.0

              ;;*t74
              local.get $$t74.1
              i32.load offset=0 align=4
              local.set $$t75

              ;;t50 | t62
              local.get $$t50
              local.get $$t62
              i32.or
              local.set $$t76

              ;;t76 | t73
              local.get $$t76
              local.get $$t73
              i32.or
              local.set $$t77

              ;;&t0.pos [#2]
              local.get $$t0.0
              call $runtime.Block.Retain
              local.get $$t0.1
              i32.const 8
              i32.add
              local.set $$t78.1
              local.get $$t78.0
              call $runtime.Block.Release
              local.set $$t78.0

              ;;*t78
              local.get $$t78.1
              i32.load offset=0 align=4
              local.set $$t79

              ;;t79 + 3:int
              local.get $$t79
              i32.const 3
              i32.add
              local.set $$t80

              ;;return true:bool, t75, t77, t80
              i32.const 1
              local.set $$ret_0
              local.get $$t75
              local.set $$ret_1
              local.get $$t77
              local.set $$ret_2
              local.get $$t80
              local.set $$ret_3
              br $$BlockFnBody

            end ;;$Block_7
            i32.const 8
            local.set $$current_block

            ;;t18 & 248:i32
            local.get $$t18
            i32.const 248
            i32.and
            local.set $$t81

            ;;t81 == 240:i32
            local.get $$t81
            i32.const 240
            i32.eq
            local.set $$t82

            ;;if t82 goto 9 else 10
            local.get $$t82
            if
              br $$Block_8
            else
              br $$Block_9
            end

          end ;;$Block_8
          i32.const 9
          local.set $$current_block

          ;;t18 & 7:i32
          local.get $$t18
          i32.const 7
          i32.and
          local.set $$t83

          ;;t83 << 18:uint64
          local.get $$t83
          i64.const 18
          i32.wrap_i64
          i32.shl
          local.set $$t84

          ;;&t0.ptr [#0]
          local.get $$t0.0
          call $runtime.Block.Retain
          local.get $$t0.1
          i32.const 0
          i32.add
          local.set $$t85.1
          local.get $$t85.0
          call $runtime.Block.Release
          local.set $$t85.0

          ;;*t85
          local.get $$t85.1
          i32.load offset=0 align=4
          local.set $$t86

          ;;convert u32 <- uint (t86)
          local.get $$t86
          local.set $$t87

          ;;&t0.pos [#2]
          local.get $$t0.0
          call $runtime.Block.Retain
          local.get $$t0.1
          i32.const 8
          i32.add
          local.set $$t88.1
          local.get $$t88.0
          call $runtime.Block.Release
          local.set $$t88.0

          ;;*t88
          local.get $$t88.1
          i32.load offset=0 align=4
          local.set $$t89

          ;;convert u32 <- int (t89)
          local.get $$t89
          local.set $$t90

          ;;t87 + t90
          local.get $$t87
          local.get $$t90
          i32.add
          local.set $$t91

          ;;t91 + 1:u32
          local.get $$t91
          i32.const 1
          i32.add
          local.set $$t92

          ;;get_u8(t92)
          local.get $$t92
          call $runtime.get_u8
          local.set $$t93

          ;;convert i32 <- u8 (t93)
          local.get $$t93
          local.set $$t94

          ;;t94 & 63:i32
          local.get $$t94
          i32.const 63
          i32.and
          local.set $$t95

          ;;t95 << 12:uint64
          local.get $$t95
          i64.const 12
          i32.wrap_i64
          i32.shl
          local.set $$t96

          ;;&t0.ptr [#0]
          local.get $$t0.0
          call $runtime.Block.Retain
          local.get $$t0.1
          i32.const 0
          i32.add
          local.set $$t97.1
          local.get $$t97.0
          call $runtime.Block.Release
          local.set $$t97.0

          ;;*t97
          local.get $$t97.1
          i32.load offset=0 align=4
          local.set $$t98

          ;;convert u32 <- uint (t98)
          local.get $$t98
          local.set $$t99

          ;;&t0.pos [#2]
          local.get $$t0.0
          call $runtime.Block.Retain
          local.get $$t0.1
          i32.const 8
          i32.add
          local.set $$t100.1
          local.get $$t100.0
          call $runtime.Block.Release
          local.set $$t100.0

          ;;*t100
          local.get $$t100.1
          i32.load offset=0 align=4
          local.set $$t101

          ;;convert u32 <- int (t101)
          local.get $$t101
          local.set $$t102

          ;;t99 + t102
          local.get $$t99
          local.get $$t102
          i32.add
          local.set $$t103

          ;;t103 + 2:u32
          local.get $$t103
          i32.const 2
          i32.add
          local.set $$t104

          ;;get_u8(t104)
          local.get $$t104
          call $runtime.get_u8
          local.set $$t105

          ;;convert i32 <- u8 (t105)
          local.get $$t105
          local.set $$t106

          ;;t106 & 63:i32
          local.get $$t106
          i32.const 63
          i32.and
          local.set $$t107

          ;;t107 << 6:uint64
          local.get $$t107
          i64.const 6
          i32.wrap_i64
          i32.shl
          local.set $$t108

          ;;&t0.ptr [#0]
          local.get $$t0.0
          call $runtime.Block.Retain
          local.get $$t0.1
          i32.const 0
          i32.add
          local.set $$t109.1
          local.get $$t109.0
          call $runtime.Block.Release
          local.set $$t109.0

          ;;*t109
          local.get $$t109.1
          i32.load offset=0 align=4
          local.set $$t110

          ;;convert u32 <- uint (t110)
          local.get $$t110
          local.set $$t111

          ;;&t0.pos [#2]
          local.get $$t0.0
          call $runtime.Block.Retain
          local.get $$t0.1
          i32.const 8
          i32.add
          local.set $$t112.1
          local.get $$t112.0
          call $runtime.Block.Release
          local.set $$t112.0

          ;;*t112
          local.get $$t112.1
          i32.load offset=0 align=4
          local.set $$t113

          ;;convert u32 <- int (t113)
          local.get $$t113
          local.set $$t114

          ;;t111 + t114
          local.get $$t111
          local.get $$t114
          i32.add
          local.set $$t115

          ;;t115 + 3:u32
          local.get $$t115
          i32.const 3
          i32.add
          local.set $$t116

          ;;get_u8(t116)
          local.get $$t116
          call $runtime.get_u8
          local.set $$t117

          ;;convert i32 <- u8 (t117)
          local.get $$t117
          local.set $$t118

          ;;t118 & 63:i32
          local.get $$t118
          i32.const 63
          i32.and
          local.set $$t119

          ;;&t0.pos [#2]
          local.get $$t0.0
          call $runtime.Block.Retain
          local.get $$t0.1
          i32.const 8
          i32.add
          local.set $$t120.1
          local.get $$t120.0
          call $runtime.Block.Release
          local.set $$t120.0

          ;;*t120
          local.get $$t120.1
          i32.load offset=0 align=4
          local.set $$t121

          ;;t84 | t96
          local.get $$t84
          local.get $$t96
          i32.or
          local.set $$t122

          ;;t122 | t108
          local.get $$t122
          local.get $$t108
          i32.or
          local.set $$t123

          ;;t123 | t119
          local.get $$t123
          local.get $$t119
          i32.or
          local.set $$t124

          ;;&t0.pos [#2]
          local.get $$t0.0
          call $runtime.Block.Retain
          local.get $$t0.1
          i32.const 8
          i32.add
          local.set $$t125.1
          local.get $$t125.0
          call $runtime.Block.Release
          local.set $$t125.0

          ;;*t125
          local.get $$t125.1
          i32.load offset=0 align=4
          local.set $$t126

          ;;t126 + 4:int
          local.get $$t126
          i32.const 4
          i32.add
          local.set $$t127

          ;;return true:bool, t121, t124, t127
          i32.const 1
          local.set $$ret_0
          local.get $$t121
          local.set $$ret_1
          local.get $$t124
          local.set $$ret_2
          local.get $$t127
          local.set $$ret_3
          br $$BlockFnBody

        end ;;$Block_9
        i32.const 10
        local.set $$current_block

        ;;return false:bool, 0:int, 0:rune, 0:int
        i32.const 0
        local.set $$ret_0
        i32.const 0
        local.set $$ret_1
        i32.const 0
        local.set $$ret_2
        i32.const 0
        local.set $$ret_3
        br $$BlockFnBody

      end ;;$Block_10
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$ret_0
  local.get $$ret_1
  local.get $$ret_2
  local.get $$ret_3
  local.get $$t0.0
  call $runtime.Block.Release
  local.get $$t1.0
  call $runtime.Block.Release
  local.get $$t3.0
  call $runtime.Block.Release
  local.get $$t6.0
  call $runtime.Block.Release
  local.get $$t8.0
  call $runtime.Block.Release
  local.get $$t10.0
  call $runtime.Block.Release
  local.get $$t13.0
  call $runtime.Block.Release
  local.get $$t21.0
  call $runtime.Block.Release
  local.get $$t23.0
  call $runtime.Block.Release
  local.get $$t30.0
  call $runtime.Block.Release
  local.get $$t33.0
  call $runtime.Block.Release
  local.get $$t41.0
  call $runtime.Block.Release
  local.get $$t44.0
  call $runtime.Block.Release
  local.get $$t51.0
  call $runtime.Block.Release
  local.get $$t54.0
  call $runtime.Block.Release
  local.get $$t63.0
  call $runtime.Block.Release
  local.get $$t66.0
  call $runtime.Block.Release
  local.get $$t74.0
  call $runtime.Block.Release
  local.get $$t78.0
  call $runtime.Block.Release
  local.get $$t85.0
  call $runtime.Block.Release
  local.get $$t88.0
  call $runtime.Block.Release
  local.get $$t97.0
  call $runtime.Block.Release
  local.get $$t100.0
  call $runtime.Block.Release
  local.get $$t109.0
  call $runtime.Block.Release
  local.get $$t112.0
  call $runtime.Block.Release
  local.get $$t120.0
  call $runtime.Block.Release
  local.get $$t125.0
  call $runtime.Block.Release
) ;;runtime.next_rune

(func $$runtime.panic_ (param $msg_ptr i32) (param $msg_len i32) (param $pos_msg_ptr i32) (param $pos_msg_len i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;waPrintString("panic: ":string)
        i32.const 0
        i32.const 8301
        i32.const 7
        call $$runtime.waPrintString

        ;;waPuts(msg_ptr, msg_len)
        local.get $msg_ptr
        local.get $msg_len
        call $$runtime.waPuts

        ;;waPrintString(" (":string)
        i32.const 0
        i32.const 8262
        i32.const 2
        call $$runtime.waPrintString

        ;;waPuts(pos_msg_ptr, pos_msg_len)
        local.get $pos_msg_ptr
        local.get $pos_msg_len
        call $$runtime.waPuts

        ;;waPrintRune(41:i32)
        i32.const 41
        call $$runtime.waPrintRune

        ;;waPrintRune(10:i32)
        i32.const 10
        call $$runtime.waPrintRune

        ;;procExit(1:i32)
        i32.const 1
        call $$runtime.procExit

        ;;return
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
) ;;$runtime.panic_

(func $runtime.popRunDeferStack
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$t0.0 i32)
  (local $$t0.1 i32)
  (local $$t0.2 i32)
  (local $$t0.3 i32)
  (local $$t1 i32)
  (local $$t2 i32)
  (local $$t3.0 i32)
  (local $$t3.1 i32)
  (local $$t3.2 i32)
  (local $$t3.3 i32)
  (local $$t4.0 i32)
  (local $$t4.1 i32)
  (local $$t5.0 i32)
  (local $$t5.1 i32)
  (local $$t6.0 i32)
  (local $$t6.1 i32)
  (local $$t6.2 i32)
  (local $$t6.3 i32)
  (local $$t7.0 i32)
  (local $$t7.1 i32)
  (local $$t7.2 i32)
  (local $$t7.3 i32)
  (local $$t8 i32)
  (local $$t9 i32)
  (local $$t10 i32)
  (local $$t11 i32)
  (local $$t12 i32)
  (local $$t13 i32)
  (local $$t14 i32)
  (local $$t15.0 i32)
  (local $$t15.1 i32)
  (local $$t16.0 i32)
  (local $$t16.1.0 i32)
  (local $$t16.1.1 i32)
  (local $$t17 i32)
  (local $$t18.0 i32)
  (local $$t18.1 i32)
  (local $$t19.0 i32)
  (local $$t19.1 i32)
  (local $$t19.2 i32)
  (local $$t19.3 i32)
  (local $$t20.0 i32)
  (local $$t20.1 i32)
  (local $$t21.0 i32)
  (local $$t21.1 i32)
  (local $$t22.0 i32)
  (local $$t22.1 i32)
  (local $$t22.2 i32)
  (local $$t22.3 i32)
  (local $$t23.0 i32)
  (local $$t23.1 i32)
  (local $$t23.2 i32)
  (local $$t23.3 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_3
        block $$Block_2
          block $$Block_1
            block $$Block_0
              block $$BlockSel
                local.get $$block_selector
                br_table 0 1 2 3 0
              end ;;$BlockSel
              i32.const 0
              local.set $$current_block

              ;;*defersStack
              i32.const 8232
              i32.load offset=0 align=4
              call $runtime.Block.Retain
              i32.const 8232
              i32.load offset=4 align=4
              i32.const 8232
              i32.load offset=8 align=4
              i32.const 8232
              i32.load offset=12 align=4
              local.set $$t0.3
              local.set $$t0.2
              local.set $$t0.1
              local.get $$t0.0
              call $runtime.Block.Release
              local.set $$t0.0

              ;;len(t0)
              local.get $$t0.2
              local.set $$t1

              ;;t1 - 1:int
              local.get $$t1
              i32.const 1
              i32.sub
              local.set $$t2

              ;;*defersStack
              i32.const 8232
              i32.load offset=0 align=4
              call $runtime.Block.Retain
              i32.const 8232
              i32.load offset=4 align=4
              i32.const 8232
              i32.load offset=8 align=4
              i32.const 8232
              i32.load offset=12 align=4
              local.set $$t3.3
              local.set $$t3.2
              local.set $$t3.1
              local.get $$t3.0
              call $runtime.Block.Release
              local.set $$t3.0

              ;;&t3[t2]
              local.get $$t3.0
              call $runtime.Block.Retain
              local.get $$t3.1
              i32.const 16
              local.get $$t2
              i32.mul
              i32.add
              local.set $$t4.1
              local.get $$t4.0
              call $runtime.Block.Release
              local.set $$t4.0

              ;;&t4.fns [#0]
              local.get $$t4.0
              call $runtime.Block.Retain
              local.get $$t4.1
              i32.const 0
              i32.add
              local.set $$t5.1
              local.get $$t5.0
              call $runtime.Block.Release
              local.set $$t5.0

              ;;*t5
              local.get $$t5.1
              i32.load offset=0 align=4
              call $runtime.Block.Retain
              local.get $$t5.1
              i32.load offset=4 align=4
              local.get $$t5.1
              i32.load offset=8 align=4
              local.get $$t5.1
              i32.load offset=12 align=4
              local.set $$t6.3
              local.set $$t6.2
              local.set $$t6.1
              local.get $$t6.0
              call $runtime.Block.Release
              local.set $$t6.0

              ;;slice t6[:]
              local.get $$t6.0
              call $runtime.Block.Retain
              local.get $$t6.1
              i32.const 12
              i32.const 0
              i32.mul
              i32.add
              local.get $$t6.2
              i32.const 0
              i32.sub
              local.get $$t6.3
              i32.const 0
              i32.sub
              local.set $$t7.3
              local.set $$t7.2
              local.set $$t7.1
              local.get $$t7.0
              call $runtime.Block.Release
              local.set $$t7.0

              ;;len(t7)
              local.get $$t7.2
              local.set $$t8

              ;;t8 - 1:int
              local.get $$t8
              i32.const 1
              i32.sub
              local.set $$t9

              ;;len(t7)
              local.get $$t7.2
              local.set $$t10

              ;;jump 1
              br $$Block_0

            end ;;$Block_0
            local.get $$current_block
            i32.const 0
            i32.eq
            if (result i32)
              i32.const -1
            else
              local.get $$t11
            end
            local.set $$t12
            i32.const 1
            local.set $$current_block

            ;;t11 + 1:int
            local.get $$t12
            i32.const 1
            i32.add
            local.set $$t11

            ;;t12 < t10
            local.get $$t11
            local.get $$t10
            i32.lt_s
            local.set $$t13

            ;;if t13 goto 2 else 3
            local.get $$t13
            if
              br $$Block_1
            else
              br $$Block_2
            end

          end ;;$Block_1
          i32.const 2
          local.set $$current_block

          ;;t9 - t12
          local.get $$t9
          local.get $$t11
          i32.sub
          local.set $$t14

          ;;&t7[t14]
          local.get $$t7.0
          call $runtime.Block.Retain
          local.get $$t7.1
          i32.const 12
          local.get $$t14
          i32.mul
          i32.add
          local.set $$t15.1
          local.get $$t15.0
          call $runtime.Block.Release
          local.set $$t15.0

          ;;*t15
          local.get $$t15.1
          i32.load offset=0 align=4
          local.get $$t15.1
          i32.load offset=4 align=4
          call $runtime.Block.Retain
          local.get $$t15.1
          i32.load offset=8 align=4
          local.set $$t16.1.1
          local.get $$t16.1.0
          call $runtime.Block.Release
          local.set $$t16.1.0
          local.set $$t16.0

          ;;t16()
          local.get $$t16.0
          local.get $$t16.1.1
          global.set $$wa.runtime.closure_data
          call_indirect (type $$$fnSig1)

          ;;t9 - t12
          local.get $$t9
          local.get $$t11
          i32.sub
          local.set $$t17

          ;;&t7[t18]
          local.get $$t7.0
          call $runtime.Block.Retain
          local.get $$t7.1
          i32.const 12
          local.get $$t17
          i32.mul
          i32.add
          local.set $$t18.1
          local.get $$t18.0
          call $runtime.Block.Release
          local.set $$t18.0

          ;;*t19 = nil:deferFn
          local.get $$t18.1
          i32.const 0
          i32.store offset=0 align=4
          local.get $$t18.1
          i32.const 0
          local.get $$t18.1
          i32.load offset=4 align=1
          call $runtime.Block.Release
          i32.store offset=4 align=1
          local.get $$t18.1
          i32.const 0
          i32.store offset=8 align=4

          ;;jump 1
          i32.const 1
          local.set $$block_selector
          br $$BlockDisp

        end ;;$Block_2
        i32.const 3
        local.set $$current_block

        ;;*defersStack
        i32.const 8232
        i32.load offset=0 align=4
        call $runtime.Block.Retain
        i32.const 8232
        i32.load offset=4 align=4
        i32.const 8232
        i32.load offset=8 align=4
        i32.const 8232
        i32.load offset=12 align=4
        local.set $$t19.3
        local.set $$t19.2
        local.set $$t19.1
        local.get $$t19.0
        call $runtime.Block.Release
        local.set $$t19.0

        ;;&t20[t2]
        local.get $$t19.0
        call $runtime.Block.Retain
        local.get $$t19.1
        i32.const 16
        local.get $$t2
        i32.mul
        i32.add
        local.set $$t20.1
        local.get $$t20.0
        call $runtime.Block.Release
        local.set $$t20.0

        ;;&t21.fns [#0]
        local.get $$t20.0
        call $runtime.Block.Retain
        local.get $$t20.1
        i32.const 0
        i32.add
        local.set $$t21.1
        local.get $$t21.0
        call $runtime.Block.Release
        local.set $$t21.0

        ;;*t22 = nil:[]deferFn
        local.get $$t21.1
        i32.const 0
        local.get $$t21.1
        i32.load offset=0 align=1
        call $runtime.Block.Release
        i32.store offset=0 align=1
        local.get $$t21.1
        i32.const 0
        i32.store offset=4 align=4
        local.get $$t21.1
        i32.const 0
        i32.store offset=8 align=4
        local.get $$t21.1
        i32.const 0
        i32.store offset=12 align=4

        ;;*defersStack
        i32.const 8232
        i32.load offset=0 align=4
        call $runtime.Block.Retain
        i32.const 8232
        i32.load offset=4 align=4
        i32.const 8232
        i32.load offset=8 align=4
        i32.const 8232
        i32.load offset=12 align=4
        local.set $$t22.3
        local.set $$t22.2
        local.set $$t22.1
        local.get $$t22.0
        call $runtime.Block.Release
        local.set $$t22.0

        ;;slice t23[:t2]
        local.get $$t22.0
        call $runtime.Block.Retain
        local.get $$t22.1
        i32.const 16
        i32.const 0
        i32.mul
        i32.add
        local.get $$t2
        i32.const 0
        i32.sub
        local.get $$t22.3
        i32.const 0
        i32.sub
        local.set $$t23.3
        local.set $$t23.2
        local.set $$t23.1
        local.get $$t23.0
        call $runtime.Block.Release
        local.set $$t23.0

        ;;*defersStack = t24
        i32.const 8232
        local.get $$t23.0
        call $runtime.Block.Retain
        i32.const 8232
        i32.load offset=0 align=1
        call $runtime.Block.Release
        i32.store offset=0 align=1
        i32.const 8232
        local.get $$t23.1
        i32.store offset=4 align=4
        i32.const 8232
        local.get $$t23.2
        i32.store offset=8 align=4
        i32.const 8232
        local.get $$t23.3
        i32.store offset=12 align=4

        ;;return
        br $$BlockFnBody

      end ;;$Block_3
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$t0.0
  call $runtime.Block.Release
  local.get $$t3.0
  call $runtime.Block.Release
  local.get $$t4.0
  call $runtime.Block.Release
  local.get $$t5.0
  call $runtime.Block.Release
  local.get $$t6.0
  call $runtime.Block.Release
  local.get $$t7.0
  call $runtime.Block.Release
  local.get $$t15.0
  call $runtime.Block.Release
  local.get $$t16.1.0
  call $runtime.Block.Release
  local.get $$t18.0
  call $runtime.Block.Release
  local.get $$t19.0
  call $runtime.Block.Release
  local.get $$t20.0
  call $runtime.Block.Release
  local.get $$t21.0
  call $runtime.Block.Release
  local.get $$t22.0
  call $runtime.Block.Release
  local.get $$t23.0
  call $runtime.Block.Release
) ;;runtime.popRunDeferStack

(func $runtime.printBool (param $v i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_3
        block $$Block_2
          block $$Block_1
            block $$Block_0
              block $$BlockSel
                local.get $$block_selector
                br_table 0 1 2 3 0
              end ;;$BlockSel
              i32.const 0
              local.set $$current_block

              ;;if v goto 1 else 3
              local.get $v
              if
                br $$Block_0
              else
                br $$Block_2
              end

            end ;;$Block_0
            i32.const 1
            local.set $$current_block

            ;;printString("true":string)
            i32.const 0
            i32.const 8308
            i32.const 4
            call $runtime.printString

            ;;jump 2
            br $$Block_1

          end ;;$Block_1
          i32.const 2
          local.set $$current_block

          ;;return
          br $$BlockFnBody

        end ;;$Block_2
        i32.const 3
        local.set $$current_block

        ;;printString("false":string)
        i32.const 0
        i32.const 8312
        i32.const 5
        call $runtime.printString

        ;;jump 2
        i32.const 2
        local.set $$block_selector
        br $$BlockDisp

      end ;;$Block_3
    end ;;$BlockDisp
  end ;;$BlockFnBody
) ;;runtime.printBool

(func $runtime.printBytes (param $b.0 i32) (param $b.1 i32) (param $b.2 i32) (param $b.3 i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$t0 i32)
  (local $$t1 i32)
  (local $$t2 i32)
  (local $$t3 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_2
        block $$Block_1
          block $$Block_0
            block $$BlockSel
              local.get $$block_selector
              br_table 0 1 2 0
            end ;;$BlockSel
            i32.const 0
            local.set $$current_block

            ;;len(b)
            local.get $b.2
            local.set $$t0

            ;;t0 > 0:int
            local.get $$t0
            i32.const 0
            i32.gt_s
            local.set $$t1

            ;;if t1 goto 1 else 2
            local.get $$t1
            if
              br $$Block_0
            else
              br $$Block_1
            end

          end ;;$Block_0
          i32.const 1
          local.set $$current_block

          ;;refToPtr_byteSlice(b)
          local.get $b.0
          local.get $b.1
          local.get $b.2
          local.get $b.3
          call $runtime.refToPtr_byteSlice
          local.set $$t2

          ;;convert i32 <- int (t0)
          local.get $$t0
          local.set $$t3

          ;;waPuts(t2, t3)
          local.get $$t2
          local.get $$t3
          call $$runtime.waPuts

          ;;return
          br $$BlockFnBody

        end ;;$Block_1
        i32.const 2
        local.set $$current_block

        ;;return
        br $$BlockFnBody

      end ;;$Block_2
    end ;;$BlockDisp
  end ;;$BlockFnBody
) ;;runtime.printBytes

(func $runtime.printF64 (param $v f64)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$t0 i32)
  (local $$t1 f64)
  (local $$t2 i32)
  (local $$t3 f64)
  (local $$t4 i32)
  (local $$t5 i32)
  (local $$t6 i32)
  (local $$t7.0 i32)
  (local $$t7.1 i32)
  (local $$t8.0 i32)
  (local $$t8.1 i32)
  (local $$t9 i32)
  (local $$t10 i32)
  (local $$t11 i32)
  (local $$t12 f64)
  (local $$t13 i32)
  (local $$t14 f64)
  (local $$t15 f64)
  (local $$t16 f64)
  (local $$t17 i32)
  (local $$t18 i32)
  (local $$t19 i32)
  (local $$t20 i32)
  (local $$t21.0 i32)
  (local $$t21.1 i32)
  (local $$t22 f64)
  (local $$t23.0 i32)
  (local $$t23.1 i32)
  (local $$t24 i32)
  (local $$t25 i32)
  (local $$t26 f64)
  (local $$t27 f64)
  (local $$t28 i32)
  (local $$t29 i32)
  (local $$t30 f64)
  (local $$t31 f64)
  (local $$t32 i32)
  (local $$t33 f64)
  (local $$t34 f64)
  (local $$t35 i32)
  (local $$t36 i32)
  (local $$t37 i32)
  (local $$t38 i32)
  (local $$t39 f64)
  (local $$t40 i32)
  (local $$t41 i32)
  (local $$t42 i32)
  (local $$t43.0 i32)
  (local $$t43.1 i32)
  (local $$t44 i32)
  (local $$t45 i32)
  (local $$t46 f64)
  (local $$t47 f64)
  (local $$t48 f64)
  (local $$t49 i32)
  (local $$t50.0 i32)
  (local $$t50.1 i32)
  (local $$t51.0 i32)
  (local $$t51.1 i32)
  (local $$t52 i32)
  (local $$t53.0 i32)
  (local $$t53.1 i32)
  (local $$t54.0 i32)
  (local $$t54.1 i32)
  (local $$t55.0 i32)
  (local $$t55.1 i32)
  (local $$t56 i32)
  (local $$t57 i32)
  (local $$t58 i32)
  (local $$t59.0 i32)
  (local $$t59.1 i32)
  (local $$t60 i32)
  (local $$t61.0 i32)
  (local $$t61.1 i32)
  (local $$t62 i32)
  (local $$t63 i32)
  (local $$t64 i32)
  (local $$t65.0 i32)
  (local $$t65.1 i32)
  (local $$t66 i32)
  (local $$t67 i32)
  (local $$t68 i32)
  (local $$t69 i32)
  (local $$t70.0 i32)
  (local $$t70.1 i32)
  (local $$t71 i32)
  (local $$t72 i32)
  (local $$t73 i32)
  (local $$t74.0 i32)
  (local $$t74.1 i32)
  (local $$t74.2 i32)
  (local $$t74.3 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_29
        block $$Block_28
          block $$Block_27
            block $$Block_26
              block $$Block_25
                block $$Block_24
                  block $$Block_23
                    block $$Block_22
                      block $$Block_21
                        block $$Block_20
                          block $$Block_19
                            block $$Block_18
                              block $$Block_17
                                block $$Block_16
                                  block $$Block_15
                                    block $$Block_14
                                      block $$Block_13
                                        block $$Block_12
                                          block $$Block_11
                                            block $$Block_10
                                              block $$Block_9
                                                block $$Block_8
                                                  block $$Block_7
                                                    block $$Block_6
                                                      block $$Block_5
                                                        block $$Block_4
                                                          block $$Block_3
                                                            block $$Block_2
                                                              block $$Block_1
                                                                block $$Block_0
                                                                  block $$BlockSel
                                                                    local.get $$block_selector
                                                                    br_table 0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 0
                                                                  end ;;$BlockSel
                                                                  i32.const 0
                                                                  local.set $$current_block

                                                                  ;;v != v
                                                                  local.get $v
                                                                  local.get $v
                                                                  f64.eq
                                                                  i32.eqz
                                                                  local.set $$t0

                                                                  ;;if t0 goto 1 else 3
                                                                  local.get $$t0
                                                                  if
                                                                    br $$Block_0
                                                                  else
                                                                    br $$Block_2
                                                                  end

                                                                end ;;$Block_0
                                                                i32.const 1
                                                                local.set $$current_block

                                                                ;;printString("NaN":string)
                                                                i32.const 0
                                                                i32.const 8317
                                                                i32.const 3
                                                                call $runtime.printString

                                                                ;;return
                                                                br $$BlockFnBody

                                                              end ;;$Block_1
                                                              i32.const 2
                                                              local.set $$current_block

                                                              ;;printString("+Inf":string)
                                                              i32.const 0
                                                              i32.const 8320
                                                              i32.const 4
                                                              call $runtime.printString

                                                              ;;return
                                                              br $$BlockFnBody

                                                            end ;;$Block_2
                                                            i32.const 3
                                                            local.set $$current_block

                                                            ;;v + v
                                                            local.get $v
                                                            local.get $v
                                                            f64.add
                                                            local.set $$t1

                                                            ;;t3 == v
                                                            local.get $$t1
                                                            local.get $v
                                                            f64.eq
                                                            local.set $$t2

                                                            ;;if t4 goto 6 else 7
                                                            local.get $$t2
                                                            if
                                                              br $$Block_5
                                                            else
                                                              br $$Block_6
                                                            end

                                                          end ;;$Block_3
                                                          i32.const 4
                                                          local.set $$current_block

                                                          ;;printString("-Inf":string)
                                                          i32.const 0
                                                          i32.const 8324
                                                          i32.const 4
                                                          call $runtime.printString

                                                          ;;return
                                                          br $$BlockFnBody

                                                        end ;;$Block_4
                                                        i32.const 5
                                                        local.set $$current_block

                                                        ;;v + v
                                                        local.get $v
                                                        local.get $v
                                                        f64.add
                                                        local.set $$t3

                                                        ;;t6 == v
                                                        local.get $$t3
                                                        local.get $v
                                                        f64.eq
                                                        local.set $$t4

                                                        ;;if t7 goto 9 else 10
                                                        local.get $$t4
                                                        if
                                                          br $$Block_8
                                                        else
                                                          br $$Block_9
                                                        end

                                                      end ;;$Block_5
                                                      i32.const 6
                                                      local.set $$current_block

                                                      ;;v > 0:f64
                                                      local.get $v
                                                      f64.const 0
                                                      f64.gt
                                                      local.set $$t5

                                                      ;;jump 7
                                                      br $$Block_6

                                                    end ;;$Block_6
                                                    local.get $$current_block
                                                    i32.const 3
                                                    i32.eq
                                                    if (result i32)
                                                      i32.const 0
                                                    else
                                                      local.get $$t5
                                                    end
                                                    local.set $$t6
                                                    i32.const 7
                                                    local.set $$current_block

                                                    ;;if t9 goto 2 else 5
                                                    local.get $$t6
                                                    if
                                                      i32.const 2
                                                      local.set $$block_selector
                                                      br $$BlockDisp
                                                    else
                                                      i32.const 5
                                                      local.set $$block_selector
                                                      br $$BlockDisp
                                                    end

                                                  end ;;$Block_7
                                                  i32.const 8
                                                  local.set $$current_block

                                                  ;;new [14]byte (buf)
                                                  i32.const 30
                                                  call $runtime.HeapAlloc
                                                  i32.const 1
                                                  i32.const 0
                                                  i32.const 14
                                                  call $runtime.Block.Init
                                                  call $runtime.DupI32
                                                  i32.const 16
                                                  i32.add
                                                  local.set $$t7.1
                                                  local.get $$t7.0
                                                  call $runtime.Block.Release
                                                  local.set $$t7.0

                                                  ;;&t10[0:int]
                                                  local.get $$t7.0
                                                  call $runtime.Block.Retain
                                                  local.get $$t7.1
                                                  i32.const 1
                                                  i32.const 0
                                                  i32.mul
                                                  i32.add
                                                  local.set $$t8.1
                                                  local.get $$t8.0
                                                  call $runtime.Block.Release
                                                  local.set $$t8.0

                                                  ;;*t11 = 43:byte
                                                  local.get $$t8.1
                                                  i32.const 43
                                                  i32.store8 offset=0 align=1

                                                  ;;v == 0:f64
                                                  local.get $v
                                                  f64.const 0
                                                  f64.eq
                                                  local.set $$t9

                                                  ;;if t12 goto 11 else 13
                                                  local.get $$t9
                                                  if
                                                    br $$Block_10
                                                  else
                                                    br $$Block_12
                                                  end

                                                end ;;$Block_8
                                                i32.const 9
                                                local.set $$current_block

                                                ;;v < 0:f64
                                                local.get $v
                                                f64.const 0
                                                f64.lt
                                                local.set $$t10

                                                ;;jump 10
                                                br $$Block_9

                                              end ;;$Block_9
                                              local.get $$current_block
                                              i32.const 5
                                              i32.eq
                                              if (result i32)
                                                i32.const 0
                                              else
                                                local.get $$t10
                                              end
                                              local.set $$t11
                                              i32.const 10
                                              local.set $$current_block

                                              ;;if t14 goto 4 else 8
                                              local.get $$t11
                                              if
                                                i32.const 4
                                                local.set $$block_selector
                                                br $$BlockDisp
                                              else
                                                i32.const 8
                                                local.set $$block_selector
                                                br $$BlockDisp
                                              end

                                            end ;;$Block_10
                                            i32.const 11
                                            local.set $$current_block

                                            ;;1:f64 / v
                                            f64.const 1
                                            local.get $v
                                            f64.div
                                            local.set $$t12

                                            ;;t15 < 0:f64
                                            local.get $$t12
                                            f64.const 0
                                            f64.lt
                                            local.set $$t13

                                            ;;if t16 goto 14 else 12
                                            local.get $$t13
                                            if
                                              br $$Block_13
                                            else
                                              br $$Block_11
                                            end

                                          end ;;$Block_11
                                          local.get $$current_block
                                          i32.const 11
                                          i32.eq
                                          if (result f64)
                                            local.get $v
                                          else
                                            local.get $$current_block
                                            i32.const 22
                                            i32.eq
                                            if (result f64)
                                              local.get $$t14
                                            else
                                              local.get $$current_block
                                              i32.const 14
                                              i32.eq
                                              if (result f64)
                                                local.get $v
                                              else
                                                local.get $$t15
                                              end
                                            end
                                          end
                                          local.get $$current_block
                                          i32.const 11
                                          i32.eq
                                          if (result i32)
                                            i32.const 0
                                          else
                                            local.get $$current_block
                                            i32.const 22
                                            i32.eq
                                            if (result i32)
                                              local.get $$t17
                                            else
                                              local.get $$current_block
                                              i32.const 14
                                              i32.eq
                                              if (result i32)
                                                i32.const 0
                                              else
                                                local.get $$t18
                                              end
                                            end
                                          end
                                          local.set $$t19
                                          local.set $$t16
                                          i32.const 12
                                          local.set $$current_block

                                          ;;jump 27
                                          br $$Block_26

                                        end ;;$Block_12
                                        i32.const 13
                                        local.set $$current_block

                                        ;;v < 0:f64
                                        local.get $v
                                        f64.const 0
                                        f64.lt
                                        local.set $$t20

                                        ;;if t19 goto 15 else 17
                                        local.get $$t20
                                        if
                                          br $$Block_14
                                        else
                                          br $$Block_16
                                        end

                                      end ;;$Block_13
                                      i32.const 14
                                      local.set $$current_block

                                      ;;&t10[0:int]
                                      local.get $$t7.0
                                      call $runtime.Block.Retain
                                      local.get $$t7.1
                                      i32.const 1
                                      i32.const 0
                                      i32.mul
                                      i32.add
                                      local.set $$t21.1
                                      local.get $$t21.0
                                      call $runtime.Block.Release
                                      local.set $$t21.0

                                      ;;*t20 = 45:byte
                                      local.get $$t21.1
                                      i32.const 45
                                      i32.store8 offset=0 align=1

                                      ;;jump 12
                                      i32.const 12
                                      local.set $$block_selector
                                      br $$BlockDisp

                                    end ;;$Block_14
                                    i32.const 15
                                    local.set $$current_block

                                    ;;-v
                                    f64.const 0
                                    local.get $v
                                    f64.sub
                                    local.set $$t22

                                    ;;&t10[0:int]
                                    local.get $$t7.0
                                    call $runtime.Block.Retain
                                    local.get $$t7.1
                                    i32.const 1
                                    i32.const 0
                                    i32.mul
                                    i32.add
                                    local.set $$t23.1
                                    local.get $$t23.0
                                    call $runtime.Block.Release
                                    local.set $$t23.0

                                    ;;*t22 = 45:byte
                                    local.get $$t23.1
                                    i32.const 45
                                    i32.store8 offset=0 align=1

                                    ;;jump 17
                                    br $$Block_16

                                  end ;;$Block_15
                                  i32.const 16
                                  local.set $$current_block

                                  ;;t26 + 1:int
                                  local.get $$t24
                                  i32.const 1
                                  i32.add
                                  local.set $$t25

                                  ;;t25 / 10:f64
                                  local.get $$t26
                                  f64.const 10
                                  f64.div
                                  local.set $$t27

                                  ;;jump 17
                                  br $$Block_16

                                end ;;$Block_16
                                local.get $$current_block
                                i32.const 13
                                i32.eq
                                if (result f64)
                                  local.get $v
                                else
                                  local.get $$current_block
                                  i32.const 16
                                  i32.eq
                                  if (result f64)
                                    local.get $$t27
                                  else
                                    local.get $$t22
                                  end
                                end
                                local.get $$current_block
                                i32.const 13
                                i32.eq
                                if (result i32)
                                  i32.const 0
                                else
                                  local.get $$current_block
                                  i32.const 16
                                  i32.eq
                                  if (result i32)
                                    local.get $$t25
                                  else
                                    i32.const 0
                                  end
                                end
                                local.set $$t24
                                local.set $$t26
                                i32.const 17
                                local.set $$current_block

                                ;;t25 >= 10:f64
                                local.get $$t26
                                f64.const 10
                                f64.ge
                                local.set $$t28

                                ;;if t27 goto 16 else 20
                                local.get $$t28
                                if
                                  i32.const 16
                                  local.set $$block_selector
                                  br $$BlockDisp
                                else
                                  br $$Block_19
                                end

                              end ;;$Block_17
                              i32.const 18
                              local.set $$current_block

                              ;;t31 - 1:int
                              local.get $$t17
                              i32.const 1
                              i32.sub
                              local.set $$t29

                              ;;t30 * 10:f64
                              local.get $$t30
                              f64.const 10
                              f64.mul
                              local.set $$t31

                              ;;jump 20
                              br $$Block_19

                            end ;;$Block_18
                            i32.const 19
                            local.set $$current_block

                            ;;jump 23
                            br $$Block_22

                          end ;;$Block_19
                          local.get $$current_block
                          i32.const 17
                          i32.eq
                          if (result f64)
                            local.get $$t26
                          else
                            local.get $$t31
                          end
                          local.get $$current_block
                          i32.const 17
                          i32.eq
                          if (result i32)
                            local.get $$t24
                          else
                            local.get $$t29
                          end
                          local.set $$t17
                          local.set $$t30
                          i32.const 20
                          local.set $$current_block

                          ;;t30 < 1:f64
                          local.get $$t30
                          f64.const 1
                          f64.lt
                          local.set $$t32

                          ;;if t32 goto 18 else 19
                          local.get $$t32
                          if
                            i32.const 18
                            local.set $$block_selector
                            br $$BlockDisp
                          else
                            i32.const 19
                            local.set $$block_selector
                            br $$BlockDisp
                          end

                        end ;;$Block_20
                        i32.const 21
                        local.set $$current_block

                        ;;t37 / 10:float64
                        local.get $$t33
                        f64.const 10
                        f64.div
                        local.set $$t34

                        ;;t38 + 1:int
                        local.get $$t35
                        i32.const 1
                        i32.add
                        local.set $$t36

                        ;;jump 23
                        br $$Block_22

                      end ;;$Block_21
                      i32.const 22
                      local.set $$current_block

                      ;;t30 + t37
                      local.get $$t30
                      local.get $$t33
                      f64.add
                      local.set $$t14

                      ;;t35 >= 10:f64
                      local.get $$t14
                      f64.const 10
                      f64.ge
                      local.set $$t37

                      ;;if t36 goto 24 else 12
                      local.get $$t37
                      if
                        br $$Block_23
                      else
                        i32.const 12
                        local.set $$block_selector
                        br $$BlockDisp
                      end

                    end ;;$Block_22
                    local.get $$current_block
                    i32.const 19
                    i32.eq
                    if (result f64)
                      f64.const 5
                    else
                      local.get $$t34
                    end
                    local.get $$current_block
                    i32.const 19
                    i32.eq
                    if (result i32)
                      i32.const 0
                    else
                      local.get $$t36
                    end
                    local.set $$t35
                    local.set $$t33
                    i32.const 23
                    local.set $$current_block

                    ;;t38 < 7:int
                    local.get $$t35
                    i32.const 7
                    i32.lt_s
                    local.set $$t38

                    ;;if t39 goto 21 else 22
                    local.get $$t38
                    if
                      i32.const 21
                      local.set $$block_selector
                      br $$BlockDisp
                    else
                      i32.const 22
                      local.set $$block_selector
                      br $$BlockDisp
                    end

                  end ;;$Block_23
                  i32.const 24
                  local.set $$current_block

                  ;;t31 + 1:int
                  local.get $$t17
                  i32.const 1
                  i32.add
                  local.set $$t18

                  ;;t35 / 10:f64
                  local.get $$t14
                  f64.const 10
                  f64.div
                  local.set $$t15

                  ;;jump 12
                  i32.const 12
                  local.set $$block_selector
                  br $$BlockDisp

                end ;;$Block_24
                i32.const 25
                local.set $$current_block

                ;;convert int <- f64 (t58)
                local.get $$t39
                i32.trunc_f64_s
                local.set $$t40

                ;;t59 + 2:int
                local.get $$t41
                i32.const 2
                i32.add
                local.set $$t42

                ;;&t10[t43]
                local.get $$t7.0
                call $runtime.Block.Retain
                local.get $$t7.1
                i32.const 1
                local.get $$t42
                i32.mul
                i32.add
                local.set $$t43.1
                local.get $$t43.0
                call $runtime.Block.Release
                local.set $$t43.0

                ;;t42 + 48:int
                local.get $$t40
                i32.const 48
                i32.add
                local.set $$t44

                ;;convert byte <- int (t45)
                local.get $$t44
                i32.const 255
                i32.and
                local.set $$t45

                ;;*t44 = t46
                local.get $$t43.1
                local.get $$t45
                i32.store8 offset=0 align=1

                ;;convert float64 <- int (t42)
                local.get $$t40
                f64.convert_i32_s
                local.set $$t46

                ;;t58 - t47
                local.get $$t39
                local.get $$t46
                f64.sub
                local.set $$t47

                ;;t48 * 10:f64
                local.get $$t47
                f64.const 10
                f64.mul
                local.set $$t48

                ;;t59 + 1:int
                local.get $$t41
                i32.const 1
                i32.add
                local.set $$t49

                ;;jump 27
                br $$Block_26

              end ;;$Block_25
              i32.const 26
              local.set $$current_block

              ;;&t10[1:int]
              local.get $$t7.0
              call $runtime.Block.Retain
              local.get $$t7.1
              i32.const 1
              i32.const 1
              i32.mul
              i32.add
              local.set $$t50.1
              local.get $$t50.0
              call $runtime.Block.Release
              local.set $$t50.0

              ;;&t10[2:int]
              local.get $$t7.0
              call $runtime.Block.Retain
              local.get $$t7.1
              i32.const 1
              i32.const 2
              i32.mul
              i32.add
              local.set $$t51.1
              local.get $$t51.0
              call $runtime.Block.Release
              local.set $$t51.0

              ;;*t52
              local.get $$t51.1
              i32.load8_u offset=0 align=1
              local.set $$t52

              ;;*t51 = t53
              local.get $$t50.1
              local.get $$t52
              i32.store8 offset=0 align=1

              ;;&t10[2:int]
              local.get $$t7.0
              call $runtime.Block.Retain
              local.get $$t7.1
              i32.const 1
              i32.const 2
              i32.mul
              i32.add
              local.set $$t53.1
              local.get $$t53.0
              call $runtime.Block.Release
              local.set $$t53.0

              ;;*t54 = 46:byte
              local.get $$t53.1
              i32.const 46
              i32.store8 offset=0 align=1

              ;;&t10[9:int]
              local.get $$t7.0
              call $runtime.Block.Retain
              local.get $$t7.1
              i32.const 1
              i32.const 9
              i32.mul
              i32.add
              local.set $$t54.1
              local.get $$t54.0
              call $runtime.Block.Release
              local.set $$t54.0

              ;;*t55 = 101:byte
              local.get $$t54.1
              i32.const 101
              i32.store8 offset=0 align=1

              ;;&t10[10:int]
              local.get $$t7.0
              call $runtime.Block.Retain
              local.get $$t7.1
              i32.const 1
              i32.const 10
              i32.mul
              i32.add
              local.set $$t55.1
              local.get $$t55.0
              call $runtime.Block.Release
              local.set $$t55.0

              ;;*t56 = 43:byte
              local.get $$t55.1
              i32.const 43
              i32.store8 offset=0 align=1

              ;;t18 < 0:int
              local.get $$t19
              i32.const 0
              i32.lt_s
              local.set $$t56

              ;;if t57 goto 28 else 29
              local.get $$t56
              if
                br $$Block_27
              else
                br $$Block_28
              end

            end ;;$Block_26
            local.get $$current_block
            i32.const 12
            i32.eq
            if (result f64)
              local.get $$t16
            else
              local.get $$t48
            end
            local.get $$current_block
            i32.const 12
            i32.eq
            if (result i32)
              i32.const 0
            else
              local.get $$t49
            end
            local.set $$t41
            local.set $$t39
            i32.const 27
            local.set $$current_block

            ;;t59 < 7:int
            local.get $$t41
            i32.const 7
            i32.lt_s
            local.set $$t57

            ;;if t60 goto 25 else 26
            local.get $$t57
            if
              i32.const 25
              local.set $$block_selector
              br $$BlockDisp
            else
              i32.const 26
              local.set $$block_selector
              br $$BlockDisp
            end

          end ;;$Block_27
          i32.const 28
          local.set $$current_block

          ;;-t18
          i32.const 0
          local.get $$t19
          i32.sub
          local.set $$t58

          ;;&t10[10:int]
          local.get $$t7.0
          call $runtime.Block.Retain
          local.get $$t7.1
          i32.const 1
          i32.const 10
          i32.mul
          i32.add
          local.set $$t59.1
          local.get $$t59.0
          call $runtime.Block.Release
          local.set $$t59.0

          ;;*t62 = 45:byte
          local.get $$t59.1
          i32.const 45
          i32.store8 offset=0 align=1

          ;;jump 29
          br $$Block_28

        end ;;$Block_28
        local.get $$current_block
        i32.const 26
        i32.eq
        if (result i32)
          local.get $$t19
        else
          local.get $$t58
        end
        local.set $$t60
        i32.const 29
        local.set $$current_block

        ;;&t10[11:int]
        local.get $$t7.0
        call $runtime.Block.Retain
        local.get $$t7.1
        i32.const 1
        i32.const 11
        i32.mul
        i32.add
        local.set $$t61.1
        local.get $$t61.0
        call $runtime.Block.Release
        local.set $$t61.0

        ;;t63 / 100:int
        local.get $$t60
        i32.const 100
        i32.div_s
        local.set $$t62

        ;;convert byte <- int (t65)
        local.get $$t62
        i32.const 255
        i32.and
        local.set $$t63

        ;;t66 + 48:byte
        local.get $$t63
        i32.const 48
        i32.add
        i32.const 255
        i32.and
        local.set $$t64

        ;;*t64 = t67
        local.get $$t61.1
        local.get $$t64
        i32.store8 offset=0 align=1

        ;;&t10[12:int]
        local.get $$t7.0
        call $runtime.Block.Retain
        local.get $$t7.1
        i32.const 1
        i32.const 12
        i32.mul
        i32.add
        local.set $$t65.1
        local.get $$t65.0
        call $runtime.Block.Release
        local.set $$t65.0

        ;;t63 / 10:int
        local.get $$t60
        i32.const 10
        i32.div_s
        local.set $$t66

        ;;convert byte <- int (t69)
        local.get $$t66
        i32.const 255
        i32.and
        local.set $$t67

        ;;t70 % 10:byte
        local.get $$t67
        i32.const 10
        i32.rem_u
        local.set $$t68

        ;;t71 + 48:byte
        local.get $$t68
        i32.const 48
        i32.add
        i32.const 255
        i32.and
        local.set $$t69

        ;;*t68 = t72
        local.get $$t65.1
        local.get $$t69
        i32.store8 offset=0 align=1

        ;;&t10[13:int]
        local.get $$t7.0
        call $runtime.Block.Retain
        local.get $$t7.1
        i32.const 1
        i32.const 13
        i32.mul
        i32.add
        local.set $$t70.1
        local.get $$t70.0
        call $runtime.Block.Release
        local.set $$t70.0

        ;;t63 % 10:int
        local.get $$t60
        i32.const 10
        i32.rem_s
        local.set $$t71

        ;;convert byte <- int (t74)
        local.get $$t71
        i32.const 255
        i32.and
        local.set $$t72

        ;;t75 + 48:byte
        local.get $$t72
        i32.const 48
        i32.add
        i32.const 255
        i32.and
        local.set $$t73

        ;;*t73 = t76
        local.get $$t70.1
        local.get $$t73
        i32.store8 offset=0 align=1

        ;;slice t10[:]
        local.get $$t7.0
        call $runtime.Block.Retain
        local.get $$t7.1
        i32.const 1
        i32.const 0
        i32.mul
        i32.add
        i32.const 14
        i32.const 0
        i32.sub
        i32.const 14
        i32.const 0
        i32.sub
        local.set $$t74.3
        local.set $$t74.2
        local.set $$t74.1
        local.get $$t74.0
        call $runtime.Block.Release
        local.set $$t74.0

        ;;printBytes(t77)
        local.get $$t74.0
        local.get $$t74.1
        local.get $$t74.2
        local.get $$t74.3
        call $runtime.printBytes

        ;;return
        br $$BlockFnBody

      end ;;$Block_29
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$t7.0
  call $runtime.Block.Release
  local.get $$t8.0
  call $runtime.Block.Release
  local.get $$t21.0
  call $runtime.Block.Release
  local.get $$t23.0
  call $runtime.Block.Release
  local.get $$t43.0
  call $runtime.Block.Release
  local.get $$t50.0
  call $runtime.Block.Release
  local.get $$t51.0
  call $runtime.Block.Release
  local.get $$t53.0
  call $runtime.Block.Release
  local.get $$t54.0
  call $runtime.Block.Release
  local.get $$t55.0
  call $runtime.Block.Release
  local.get $$t59.0
  call $runtime.Block.Release
  local.get $$t61.0
  call $runtime.Block.Release
  local.get $$t65.0
  call $runtime.Block.Release
  local.get $$t70.0
  call $runtime.Block.Release
  local.get $$t74.0
  call $runtime.Block.Release
) ;;runtime.printF64

(func $runtime.printHex (param $v i64)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$t0.0 i32)
  (local $$t0.1 i32)
  (local $$t1 i32)
  (local $$t2 i32)
  (local $$t3.0 i32)
  (local $$t3.1 i32)
  (local $$t4 i64)
  (local $$t5 i64)
  (local $$t6 i32)
  (local $$t7 i32)
  (local $$t8 i32)
  (local $$t9.0 i32)
  (local $$t9.1 i32)
  (local $$t10 i32)
  (local $$t11.0 i32)
  (local $$t11.1 i32)
  (local $$t12.0 i32)
  (local $$t12.1 i32)
  (local $$t12.2 i32)
  (local $$t12.3 i32)
  (local $$t13 i64)
  (local $$t14 i32)
  (local $$t15 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_4
        block $$Block_3
          block $$Block_2
            block $$Block_1
              block $$Block_0
                block $$BlockSel
                  local.get $$block_selector
                  br_table 0 1 2 3 4 0
                end ;;$BlockSel
                i32.const 0
                local.set $$current_block

                ;;new [64]byte (buf)
                i32.const 80
                call $runtime.HeapAlloc
                i32.const 1
                i32.const 0
                i32.const 64
                call $runtime.Block.Init
                call $runtime.DupI32
                i32.const 16
                i32.add
                local.set $$t0.1
                local.get $$t0.0
                call $runtime.Block.Release
                local.set $$t0.0

                ;;64:int - 1:int
                i32.const 64
                i32.const 1
                i32.sub
                local.set $$t1

                ;;jump 3
                br $$Block_2

              end ;;$Block_0
              i32.const 1
              local.set $$current_block

              ;;&t0[t13]
              local.get $$t0.0
              call $runtime.Block.Retain
              local.get $$t0.1
              i32.const 1
              local.get $$t2
              i32.mul
              i32.add
              local.set $$t3.1
              local.get $$t3.0
              call $runtime.Block.Release
              local.set $$t3.0

              ;;t12 % 16:u64
              local.get $$t4
              i64.const 16
              i64.rem_u
              local.set $$t5

              ;;"0123456789abcdef":untyped string[t3]
              i32.const 8328
              local.get $$t5
              i32.wrap_i64
              i32.add
              i32.load8_u offset=0 align=1
              local.set $$t6

              ;;*t2 = t4
              local.get $$t3.1
              local.get $$t6
              i32.store8 offset=0 align=1

              ;;t12 < 16:u64
              local.get $$t4
              i64.const 16
              i64.lt_u
              local.set $$t7

              ;;if t5 goto 2 else 4
              local.get $$t7
              if
                br $$Block_1
              else
                br $$Block_3
              end

            end ;;$Block_1
            i32.const 2
            local.set $$current_block

            ;;t13 - 1:int
            local.get $$t2
            i32.const 1
            i32.sub
            local.set $$t8

            ;;&t0[t6]
            local.get $$t0.0
            call $runtime.Block.Retain
            local.get $$t0.1
            i32.const 1
            local.get $$t8
            i32.mul
            i32.add
            local.set $$t9.1
            local.get $$t9.0
            call $runtime.Block.Release
            local.set $$t9.0

            ;;*t7 = 120:byte
            local.get $$t9.1
            i32.const 120
            i32.store8 offset=0 align=1

            ;;t6 - 1:int
            local.get $$t8
            i32.const 1
            i32.sub
            local.set $$t10

            ;;&t0[t8]
            local.get $$t0.0
            call $runtime.Block.Retain
            local.get $$t0.1
            i32.const 1
            local.get $$t10
            i32.mul
            i32.add
            local.set $$t11.1
            local.get $$t11.0
            call $runtime.Block.Release
            local.set $$t11.0

            ;;*t9 = 48:byte
            local.get $$t11.1
            i32.const 48
            i32.store8 offset=0 align=1

            ;;slice t0[t8:]
            local.get $$t0.0
            call $runtime.Block.Retain
            local.get $$t0.1
            i32.const 1
            local.get $$t10
            i32.mul
            i32.add
            i32.const 64
            local.get $$t10
            i32.sub
            i32.const 64
            local.get $$t10
            i32.sub
            local.set $$t12.3
            local.set $$t12.2
            local.set $$t12.1
            local.get $$t12.0
            call $runtime.Block.Release
            local.set $$t12.0

            ;;printBytes(t10)
            local.get $$t12.0
            local.get $$t12.1
            local.get $$t12.2
            local.get $$t12.3
            call $runtime.printBytes

            ;;return
            br $$BlockFnBody

          end ;;$Block_2
          local.get $$current_block
          i32.const 0
          i32.eq
          if (result i64)
            local.get $v
          else
            local.get $$t13
          end
          local.get $$current_block
          i32.const 0
          i32.eq
          if (result i32)
            local.get $$t1
          else
            local.get $$t14
          end
          local.set $$t2
          local.set $$t4
          i32.const 3
          local.set $$current_block

          ;;t13 > 0:int
          local.get $$t2
          i32.const 0
          i32.gt_s
          local.set $$t15

          ;;if t14 goto 1 else 2
          local.get $$t15
          if
            i32.const 1
            local.set $$block_selector
            br $$BlockDisp
          else
            i32.const 2
            local.set $$block_selector
            br $$BlockDisp
          end

        end ;;$Block_3
        i32.const 4
        local.set $$current_block

        ;;t12 / 16:u64
        local.get $$t4
        i64.const 16
        i64.div_u
        local.set $$t13

        ;;t13 - 1:int
        local.get $$t2
        i32.const 1
        i32.sub
        local.set $$t14

        ;;jump 3
        i32.const 3
        local.set $$block_selector
        br $$BlockDisp

      end ;;$Block_4
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$t0.0
  call $runtime.Block.Release
  local.get $$t3.0
  call $runtime.Block.Release
  local.get $$t9.0
  call $runtime.Block.Release
  local.get $$t11.0
  call $runtime.Block.Release
  local.get $$t12.0
  call $runtime.Block.Release
) ;;runtime.printHex

(func $runtime.printI64 (param $v i64)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$t0 i32)
  (local $$t1 i64)
  (local $$t2 i64)
  (local $$t3 i64)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_2
        block $$Block_1
          block $$Block_0
            block $$BlockSel
              local.get $$block_selector
              br_table 0 1 2 0
            end ;;$BlockSel
            i32.const 0
            local.set $$current_block

            ;;v < 0:i64
            local.get $v
            i64.const 0
            i64.lt_s
            local.set $$t0

            ;;if t0 goto 1 else 2
            local.get $$t0
            if
              br $$Block_0
            else
              br $$Block_1
            end

          end ;;$Block_0
          i32.const 1
          local.set $$current_block

          ;;printString("-":string)
          i32.const 0
          i32.const 8324
          i32.const 1
          call $runtime.printString

          ;;-v
          i64.const 0
          local.get $v
          i64.sub
          local.set $$t1

          ;;jump 2
          br $$Block_1

        end ;;$Block_1
        local.get $$current_block
        i32.const 0
        i32.eq
        if (result i64)
          local.get $v
        else
          local.get $$t1
        end
        local.set $$t2
        i32.const 2
        local.set $$current_block

        ;;convert u64 <- i64 (t3)
        local.get $$t2
        local.set $$t3

        ;;printU64(t4)
        local.get $$t3
        call $runtime.printU64

        ;;return
        br $$BlockFnBody

      end ;;$Block_2
    end ;;$BlockDisp
  end ;;$BlockFnBody
) ;;runtime.printI64

(func $runtime.printNewline
  (local $$block_selector i32)
  (local $$current_block i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;printString("\n":string)
        i32.const 0
        i32.const 8344
        i32.const 1
        call $runtime.printString

        ;;return
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
) ;;runtime.printNewline

(func $runtime.printSlice (param $s.0 i32) (param $s.1 i32) (param $s.2 i32) (param $s.3 i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$t0 i32)
  (local $$t1 i32)
  (local $$t2 i32)
  (local $$t3 i64)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;len(s)
        local.get $s.2
        local.set $$t0

        ;;cap(s)
        local.get $s.3
        local.set $$t1

        ;;print("[":string, t0, "/":string, t1, "]":string)
        i32.const 8345
        i32.const 1
        call $$runtime.waPuts
        i32.const 32
        call $$runtime.waPrintRune
        local.get $$t0
        call $$runtime.waPrintI32
        i32.const 32
        call $$runtime.waPrintRune
        i32.const 8346
        i32.const 1
        call $$runtime.waPuts
        i32.const 32
        call $$runtime.waPrintRune
        local.get $$t1
        call $$runtime.waPrintI32
        i32.const 32
        call $$runtime.waPrintRune
        i32.const 8347
        i32.const 1
        call $$runtime.waPuts

        ;;refToPtr_byteSlice(s)
        local.get $s.0
        local.get $s.1
        local.get $s.2
        local.get $s.3
        call $runtime.refToPtr_byteSlice
        local.set $$t2

        ;;convert u64 <- i32 (t3)
        local.get $$t2
        i64.extend_i32_u
        local.set $$t3

        ;;printHex(t4)
        local.get $$t3
        call $runtime.printHex

        ;;return
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
) ;;runtime.printSlice

(func $runtime.printSpace
  (local $$block_selector i32)
  (local $$current_block i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;printString(" ":string)
        i32.const 0
        i32.const 8255
        i32.const 1
        call $runtime.printString

        ;;return
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
) ;;runtime.printSpace

(func $runtime.printString (param $s.0 i32) (param $s.1 i32) (param $s.2 i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$t0 i32)
  (local $$t1 i32)
  (local $$t2 i32)
  (local $$t3 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_2
        block $$Block_1
          block $$Block_0
            block $$BlockSel
              local.get $$block_selector
              br_table 0 1 2 0
            end ;;$BlockSel
            i32.const 0
            local.set $$current_block

            ;;len(s)
            local.get $s.2
            local.set $$t0

            ;;t0 > 0:int
            local.get $$t0
            i32.const 0
            i32.gt_s
            local.set $$t1

            ;;if t1 goto 1 else 2
            local.get $$t1
            if
              br $$Block_0
            else
              br $$Block_1
            end

          end ;;$Block_0
          i32.const 1
          local.set $$current_block

          ;;refToPtr_string(s)
          local.get $s.0
          local.get $s.1
          local.get $s.2
          call $runtime.refToPtr_string
          local.set $$t2

          ;;convert i32 <- int (t0)
          local.get $$t0
          local.set $$t3

          ;;waPuts(t2, t3)
          local.get $$t2
          local.get $$t3
          call $$runtime.waPuts

          ;;return
          br $$BlockFnBody

        end ;;$Block_1
        i32.const 2
        local.set $$current_block

        ;;return
        br $$BlockFnBody

      end ;;$Block_2
    end ;;$BlockDisp
  end ;;$BlockFnBody
) ;;runtime.printString

(func $runtime.printU64 (param $v i64)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$t0.0 i32)
  (local $$t0.1 i32)
  (local $$t1 i32)
  (local $$t2 i32)
  (local $$t3.0 i32)
  (local $$t3.1 i32)
  (local $$t4 i64)
  (local $$t5 i64)
  (local $$t6 i64)
  (local $$t7 i32)
  (local $$t8 i32)
  (local $$t9.0 i32)
  (local $$t9.1 i32)
  (local $$t9.2 i32)
  (local $$t9.3 i32)
  (local $$t10 i64)
  (local $$t11 i32)
  (local $$t12 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_4
        block $$Block_3
          block $$Block_2
            block $$Block_1
              block $$Block_0
                block $$BlockSel
                  local.get $$block_selector
                  br_table 0 1 2 3 4 0
                end ;;$BlockSel
                i32.const 0
                local.set $$current_block

                ;;new [64]byte (buf)
                i32.const 80
                call $runtime.HeapAlloc
                i32.const 1
                i32.const 0
                i32.const 64
                call $runtime.Block.Init
                call $runtime.DupI32
                i32.const 16
                i32.add
                local.set $$t0.1
                local.get $$t0.0
                call $runtime.Block.Release
                local.set $$t0.0

                ;;64:int - 1:int
                i32.const 64
                i32.const 1
                i32.sub
                local.set $$t1

                ;;jump 3
                br $$Block_2

              end ;;$Block_0
              i32.const 1
              local.set $$current_block

              ;;&t0[t10]
              local.get $$t0.0
              call $runtime.Block.Retain
              local.get $$t0.1
              i32.const 1
              local.get $$t2
              i32.mul
              i32.add
              local.set $$t3.1
              local.get $$t3.0
              call $runtime.Block.Release
              local.set $$t3.0

              ;;t9 % 10:u64
              local.get $$t4
              i64.const 10
              i64.rem_u
              local.set $$t5

              ;;t3 + 48:u64
              local.get $$t5
              i64.const 48
              i64.add
              local.set $$t6

              ;;convert byte <- u64 (t4)
              local.get $$t6
              i32.wrap_i64
              i32.const 255
              i32.and
              local.set $$t7

              ;;*t2 = t5
              local.get $$t3.1
              local.get $$t7
              i32.store8 offset=0 align=1

              ;;t9 < 10:u64
              local.get $$t4
              i64.const 10
              i64.lt_u
              local.set $$t8

              ;;if t6 goto 2 else 4
              local.get $$t8
              if
                br $$Block_1
              else
                br $$Block_3
              end

            end ;;$Block_1
            i32.const 2
            local.set $$current_block

            ;;slice t0[t10:]
            local.get $$t0.0
            call $runtime.Block.Retain
            local.get $$t0.1
            i32.const 1
            local.get $$t2
            i32.mul
            i32.add
            i32.const 64
            local.get $$t2
            i32.sub
            i32.const 64
            local.get $$t2
            i32.sub
            local.set $$t9.3
            local.set $$t9.2
            local.set $$t9.1
            local.get $$t9.0
            call $runtime.Block.Release
            local.set $$t9.0

            ;;printBytes(t7)
            local.get $$t9.0
            local.get $$t9.1
            local.get $$t9.2
            local.get $$t9.3
            call $runtime.printBytes

            ;;return
            br $$BlockFnBody

          end ;;$Block_2
          local.get $$current_block
          i32.const 0
          i32.eq
          if (result i64)
            local.get $v
          else
            local.get $$t10
          end
          local.get $$current_block
          i32.const 0
          i32.eq
          if (result i32)
            local.get $$t1
          else
            local.get $$t11
          end
          local.set $$t2
          local.set $$t4
          i32.const 3
          local.set $$current_block

          ;;t10 > 0:int
          local.get $$t2
          i32.const 0
          i32.gt_s
          local.set $$t12

          ;;if t11 goto 1 else 2
          local.get $$t12
          if
            i32.const 1
            local.set $$block_selector
            br $$BlockDisp
          else
            i32.const 2
            local.set $$block_selector
            br $$BlockDisp
          end

        end ;;$Block_3
        i32.const 4
        local.set $$current_block

        ;;t9 / 10:u64
        local.get $$t4
        i64.const 10
        i64.div_u
        local.set $$t10

        ;;t10 - 1:int
        local.get $$t2
        i32.const 1
        i32.sub
        local.set $$t11

        ;;jump 3
        i32.const 3
        local.set $$block_selector
        br $$BlockDisp

      end ;;$Block_4
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$t0.0
  call $runtime.Block.Release
  local.get $$t3.0
  call $runtime.Block.Release
  local.get $$t9.0
  call $runtime.Block.Release
) ;;runtime.printU64

(func $$runtime.procExit (param $code i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;syscall/js.ProcExit(code)
        local.get $code
        call $syscall$js.ProcExit

        ;;return
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
) ;;$runtime.procExit

(func $$$$$$.underlying.$$OnFree (param $$ptr i32)
  local.get $$ptr
  i32.const 4
  i32.add
  i32.const 10
  call_indirect (type $$OnFree)
) ;;$$$$$.underlying.$$OnFree

(func $$$$$$.$array1.underlying.$$OnFree (param $$ptr i32)
  local.get $$ptr
  i32.const 17
  call_indirect (type $$OnFree)
) ;;$$$$$.$array1.underlying.$$OnFree

(func $$$$$$.$slice.append (param $x.0 i32) (param $x.1 i32) (param $x.2 i32) (param $x.3 i32) (param $y.0 i32) (param $y.1 i32) (param $y.2 i32) (param $y.3 i32) (result i32 i32 i32 i32)
  (local $item.0 i32)
  (local $item.1.0 i32)
  (local $item.1.1 i32)
  (local $x_len i32)
  (local $y_len i32)
  (local $new_len i32)
  (local $src i32)
  (local $dest i32)
  (local $new_cap i32)
  local.get $x.2
  local.set $x_len
  local.get $y.2
  local.set $y_len
  local.get $x_len
  local.get $y_len
  i32.add
  local.set $new_len
  local.get $new_len
  local.get $x.3
  i32.le_u
  if (result i32 i32 i32 i32)
    local.get $x.0
    call $runtime.Block.Retain
    local.get $x.1
    local.get $new_len
    local.get $x.3
    local.get $y.1
    local.set $src
    local.get $x.1
    i32.const 12
    local.get $x_len
    i32.mul
    i32.add
    local.set $dest
    block $block1
      loop $loop1
        local.get $y_len
        i32.eqz
        if
          br $block1
        else
        end
        local.get $src
        i32.load offset=0 align=4
        local.get $src
        i32.load offset=4 align=4
        call $runtime.Block.Retain
        local.get $src
        i32.load offset=8 align=4
        local.set $item.1.1
        local.get $item.1.0
        call $runtime.Block.Release
        local.set $item.1.0
        local.set $item.0
        local.get $dest
        local.get $item.0
        i32.store offset=0 align=4
        local.get $dest
        local.get $item.1.0
        call $runtime.Block.Retain
        local.get $dest
        i32.load offset=4 align=1
        call $runtime.Block.Release
        i32.store offset=4 align=1
        local.get $dest
        local.get $item.1.1
        i32.store offset=8 align=4
        local.get $src
        i32.const 12
        i32.add
        local.set $src
        local.get $dest
        i32.const 12
        i32.add
        local.set $dest
        local.get $y_len
        i32.const 1
        i32.sub
        local.set $y_len
        br $loop1
      end ;;loop1
    end ;;block1
  else
    local.get $new_len
    i32.const 2
    i32.mul
    local.set $new_cap
    local.get $new_cap
    i32.const 12
    i32.mul
    i32.const 16
    i32.add
    call $runtime.HeapAlloc
    local.get $new_cap
    i32.const 17
    i32.const 12
    call $runtime.Block.Init
    call $runtime.DupI32
    i32.const 16
    i32.add
    call $runtime.DupI32
    local.set $dest
    local.get $new_len
    local.get $new_cap
    local.get $x.1
    local.set $src
    block $block2
      loop $loop2
        local.get $x_len
        i32.eqz
        if
          br $block2
        else
        end
        local.get $src
        i32.load offset=0 align=4
        local.get $src
        i32.load offset=4 align=4
        call $runtime.Block.Retain
        local.get $src
        i32.load offset=8 align=4
        local.set $item.1.1
        local.get $item.1.0
        call $runtime.Block.Release
        local.set $item.1.0
        local.set $item.0
        local.get $dest
        local.get $item.0
        i32.store offset=0 align=4
        local.get $dest
        local.get $item.1.0
        call $runtime.Block.Retain
        local.get $dest
        i32.load offset=4 align=1
        call $runtime.Block.Release
        i32.store offset=4 align=1
        local.get $dest
        local.get $item.1.1
        i32.store offset=8 align=4
        local.get $src
        i32.const 12
        i32.add
        local.set $src
        local.get $dest
        i32.const 12
        i32.add
        local.set $dest
        local.get $x_len
        i32.const 1
        i32.sub
        local.set $x_len
        br $loop2
      end ;;loop2
    end ;;block2
    local.get $y.1
    local.set $src
    block $block3
      loop $loop3
        local.get $y_len
        i32.eqz
        if
          br $block3
        else
        end
        local.get $src
        i32.load offset=0 align=4
        local.get $src
        i32.load offset=4 align=4
        call $runtime.Block.Retain
        local.get $src
        i32.load offset=8 align=4
        local.set $item.1.1
        local.get $item.1.0
        call $runtime.Block.Release
        local.set $item.1.0
        local.set $item.0
        local.get $dest
        local.get $item.0
        i32.store offset=0 align=4
        local.get $dest
        local.get $item.1.0
        call $runtime.Block.Retain
        local.get $dest
        i32.load offset=4 align=1
        call $runtime.Block.Release
        i32.store offset=4 align=1
        local.get $dest
        local.get $item.1.1
        i32.store offset=8 align=4
        local.get $src
        i32.const 12
        i32.add
        local.set $src
        local.get $dest
        i32.const 12
        i32.add
        local.set $dest
        local.get $y_len
        i32.const 1
        i32.sub
        local.set $y_len
        br $loop3
      end ;;loop3
    end ;;block3
  end
  local.get $item.1.0
  call $runtime.Block.Release
) ;;$$$$$.$slice.append

(func $runtime.pushDeferFunc (param $f.0 i32) (param $f.1.0 i32) (param $f.1.1 i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$t0.0 i32)
  (local $$t0.1 i32)
  (local $$t0.2 i32)
  (local $$t0.3 i32)
  (local $$t1 i32)
  (local $$t2 i32)
  (local $$t3.0 i32)
  (local $$t3.1 i32)
  (local $$t3.2 i32)
  (local $$t3.3 i32)
  (local $$t4.0 i32)
  (local $$t4.1 i32)
  (local $$t5.0 i32)
  (local $$t5.1 i32)
  (local $$t6.0 i32)
  (local $$t6.1 i32)
  (local $$t6.2 i32)
  (local $$t6.3 i32)
  (local $$t7.0 i32)
  (local $$t7.1 i32)
  (local $$t8.0 i32)
  (local $$t8.1 i32)
  (local $$t9.0 i32)
  (local $$t9.1 i32)
  (local $$t9.2 i32)
  (local $$t9.3 i32)
  (local $$t10.0 i32)
  (local $$t10.1 i32)
  (local $$t11.0 i32)
  (local $$t11.1 i32)
  (local $$t12.0 i32)
  (local $$t12.1 i32)
  (local $$t12.2 i32)
  (local $$t12.3 i32)
  (local $$t13.0 i32)
  (local $$t13.1 i32)
  (local $$t13.2 i32)
  (local $$t13.3 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;*defersStack
        i32.const 8232
        i32.load offset=0 align=4
        call $runtime.Block.Retain
        i32.const 8232
        i32.load offset=4 align=4
        i32.const 8232
        i32.load offset=8 align=4
        i32.const 8232
        i32.load offset=12 align=4
        local.set $$t0.3
        local.set $$t0.2
        local.set $$t0.1
        local.get $$t0.0
        call $runtime.Block.Release
        local.set $$t0.0

        ;;len(t0)
        local.get $$t0.2
        local.set $$t1

        ;;t1 - 1:int
        local.get $$t1
        i32.const 1
        i32.sub
        local.set $$t2

        ;;*defersStack
        i32.const 8232
        i32.load offset=0 align=4
        call $runtime.Block.Retain
        i32.const 8232
        i32.load offset=4 align=4
        i32.const 8232
        i32.load offset=8 align=4
        i32.const 8232
        i32.load offset=12 align=4
        local.set $$t3.3
        local.set $$t3.2
        local.set $$t3.1
        local.get $$t3.0
        call $runtime.Block.Release
        local.set $$t3.0

        ;;&t3[t2]
        local.get $$t3.0
        call $runtime.Block.Retain
        local.get $$t3.1
        i32.const 16
        local.get $$t2
        i32.mul
        i32.add
        local.set $$t4.1
        local.get $$t4.0
        call $runtime.Block.Release
        local.set $$t4.0

        ;;&t4.fns [#0]
        local.get $$t4.0
        call $runtime.Block.Retain
        local.get $$t4.1
        i32.const 0
        i32.add
        local.set $$t5.1
        local.get $$t5.0
        call $runtime.Block.Release
        local.set $$t5.0

        ;;*defersStack
        i32.const 8232
        i32.load offset=0 align=4
        call $runtime.Block.Retain
        i32.const 8232
        i32.load offset=4 align=4
        i32.const 8232
        i32.load offset=8 align=4
        i32.const 8232
        i32.load offset=12 align=4
        local.set $$t6.3
        local.set $$t6.2
        local.set $$t6.1
        local.get $$t6.0
        call $runtime.Block.Release
        local.set $$t6.0

        ;;&t6[t2]
        local.get $$t6.0
        call $runtime.Block.Retain
        local.get $$t6.1
        i32.const 16
        local.get $$t2
        i32.mul
        i32.add
        local.set $$t7.1
        local.get $$t7.0
        call $runtime.Block.Release
        local.set $$t7.0

        ;;&t7.fns [#0]
        local.get $$t7.0
        call $runtime.Block.Retain
        local.get $$t7.1
        i32.const 0
        i32.add
        local.set $$t8.1
        local.get $$t8.0
        call $runtime.Block.Release
        local.set $$t8.0

        ;;*t8
        local.get $$t8.1
        i32.load offset=0 align=4
        call $runtime.Block.Retain
        local.get $$t8.1
        i32.load offset=4 align=4
        local.get $$t8.1
        i32.load offset=8 align=4
        local.get $$t8.1
        i32.load offset=12 align=4
        local.set $$t9.3
        local.set $$t9.2
        local.set $$t9.1
        local.get $$t9.0
        call $runtime.Block.Release
        local.set $$t9.0

        ;;new [1]deferFn (varargs)
        i32.const 28
        call $runtime.HeapAlloc
        i32.const 1
        i32.const 18
        i32.const 12
        call $runtime.Block.Init
        call $runtime.DupI32
        i32.const 16
        i32.add
        local.set $$t10.1
        local.get $$t10.0
        call $runtime.Block.Release
        local.set $$t10.0

        ;;&t10[0:int]
        local.get $$t10.0
        call $runtime.Block.Retain
        local.get $$t10.1
        i32.const 12
        i32.const 0
        i32.mul
        i32.add
        local.set $$t11.1
        local.get $$t11.0
        call $runtime.Block.Release
        local.set $$t11.0

        ;;*t11 = f
        local.get $$t11.1
        local.get $f.0
        i32.store offset=0 align=4
        local.get $$t11.1
        local.get $f.1.0
        call $runtime.Block.Retain
        local.get $$t11.1
        i32.load offset=4 align=1
        call $runtime.Block.Release
        i32.store offset=4 align=1
        local.get $$t11.1
        local.get $f.1.1
        i32.store offset=8 align=4

        ;;slice t10[:]
        local.get $$t10.0
        call $runtime.Block.Retain
        local.get $$t10.1
        i32.const 12
        i32.const 0
        i32.mul
        i32.add
        i32.const 1
        i32.const 0
        i32.sub
        i32.const 1
        i32.const 0
        i32.sub
        local.set $$t12.3
        local.set $$t12.2
        local.set $$t12.1
        local.get $$t12.0
        call $runtime.Block.Release
        local.set $$t12.0

        ;;append(t9, t12...)
        local.get $$t9.0
        local.get $$t9.1
        local.get $$t9.2
        local.get $$t9.3
        local.get $$t12.0
        local.get $$t12.1
        local.get $$t12.2
        local.get $$t12.3
        call $$$$$$.$slice.append
        local.set $$t13.3
        local.set $$t13.2
        local.set $$t13.1
        local.get $$t13.0
        call $runtime.Block.Release
        local.set $$t13.0

        ;;*t5 = t13
        local.get $$t5.1
        local.get $$t13.0
        call $runtime.Block.Retain
        local.get $$t5.1
        i32.load offset=0 align=1
        call $runtime.Block.Release
        i32.store offset=0 align=1
        local.get $$t5.1
        local.get $$t13.1
        i32.store offset=4 align=4
        local.get $$t5.1
        local.get $$t13.2
        i32.store offset=8 align=4
        local.get $$t5.1
        local.get $$t13.3
        i32.store offset=12 align=4

        ;;return
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$t0.0
  call $runtime.Block.Release
  local.get $$t3.0
  call $runtime.Block.Release
  local.get $$t4.0
  call $runtime.Block.Release
  local.get $$t5.0
  call $runtime.Block.Release
  local.get $$t6.0
  call $runtime.Block.Release
  local.get $$t7.0
  call $runtime.Block.Release
  local.get $$t8.0
  call $runtime.Block.Release
  local.get $$t9.0
  call $runtime.Block.Release
  local.get $$t10.0
  call $runtime.Block.Release
  local.get $$t11.0
  call $runtime.Block.Release
  local.get $$t12.0
  call $runtime.Block.Release
  local.get $$t13.0
  call $runtime.Block.Release
) ;;runtime.pushDeferFunc

(func $$$$$$.$$block.$$OnFree (param $ptr i32)
  local.get $ptr
  i32.load offset=0 align=1
  call $runtime.Block.Release
  local.get $ptr
  i32.const 0
  i32.store offset=0 align=1
) ;;$$$$$.$$block.$$OnFree

(func $$$$$$.$slice.underlying.$$OnFree (param $$ptr i32)
  local.get $$ptr
  i32.const 19
  call_indirect (type $$OnFree)
) ;;$$$$$.$slice.underlying.$$OnFree

(func $$runtime.defers.$$OnFree (param $$ptr i32)
  local.get $$ptr
  i32.const 20
  call_indirect (type $$OnFree)
) ;;$runtime.defers.$$OnFree

(func $$runtime.defers.$array1.underlying.$$OnFree (param $$ptr i32)
  local.get $$ptr
  i32.const 20
  call_indirect (type $$OnFree)
) ;;$runtime.defers.$array1.underlying.$$OnFree

(func $$runtime.defers.$slice.append (param $x.0 i32) (param $x.1 i32) (param $x.2 i32) (param $x.3 i32) (param $y.0 i32) (param $y.1 i32) (param $y.2 i32) (param $y.3 i32) (result i32 i32 i32 i32)
  (local $item.0.0 i32)
  (local $item.0.1 i32)
  (local $item.0.2 i32)
  (local $item.0.3 i32)
  (local $x_len i32)
  (local $y_len i32)
  (local $new_len i32)
  (local $src i32)
  (local $dest i32)
  (local $new_cap i32)
  local.get $x.2
  local.set $x_len
  local.get $y.2
  local.set $y_len
  local.get $x_len
  local.get $y_len
  i32.add
  local.set $new_len
  local.get $new_len
  local.get $x.3
  i32.le_u
  if (result i32 i32 i32 i32)
    local.get $x.0
    call $runtime.Block.Retain
    local.get $x.1
    local.get $new_len
    local.get $x.3
    local.get $y.1
    local.set $src
    local.get $x.1
    i32.const 16
    local.get $x_len
    i32.mul
    i32.add
    local.set $dest
    block $block1
      loop $loop1
        local.get $y_len
        i32.eqz
        if
          br $block1
        else
        end
        local.get $src
        i32.load offset=0 align=4
        call $runtime.Block.Retain
        local.get $src
        i32.load offset=4 align=4
        local.get $src
        i32.load offset=8 align=4
        local.get $src
        i32.load offset=12 align=4
        local.set $item.0.3
        local.set $item.0.2
        local.set $item.0.1
        local.get $item.0.0
        call $runtime.Block.Release
        local.set $item.0.0
        local.get $dest
        local.get $item.0.0
        call $runtime.Block.Retain
        local.get $dest
        i32.load offset=0 align=1
        call $runtime.Block.Release
        i32.store offset=0 align=1
        local.get $dest
        local.get $item.0.1
        i32.store offset=4 align=4
        local.get $dest
        local.get $item.0.2
        i32.store offset=8 align=4
        local.get $dest
        local.get $item.0.3
        i32.store offset=12 align=4
        local.get $src
        i32.const 16
        i32.add
        local.set $src
        local.get $dest
        i32.const 16
        i32.add
        local.set $dest
        local.get $y_len
        i32.const 1
        i32.sub
        local.set $y_len
        br $loop1
      end ;;loop1
    end ;;block1
  else
    local.get $new_len
    i32.const 2
    i32.mul
    local.set $new_cap
    local.get $new_cap
    i32.const 16
    i32.mul
    i32.const 16
    i32.add
    call $runtime.HeapAlloc
    local.get $new_cap
    i32.const 21
    i32.const 16
    call $runtime.Block.Init
    call $runtime.DupI32
    i32.const 16
    i32.add
    call $runtime.DupI32
    local.set $dest
    local.get $new_len
    local.get $new_cap
    local.get $x.1
    local.set $src
    block $block2
      loop $loop2
        local.get $x_len
        i32.eqz
        if
          br $block2
        else
        end
        local.get $src
        i32.load offset=0 align=4
        call $runtime.Block.Retain
        local.get $src
        i32.load offset=4 align=4
        local.get $src
        i32.load offset=8 align=4
        local.get $src
        i32.load offset=12 align=4
        local.set $item.0.3
        local.set $item.0.2
        local.set $item.0.1
        local.get $item.0.0
        call $runtime.Block.Release
        local.set $item.0.0
        local.get $dest
        local.get $item.0.0
        call $runtime.Block.Retain
        local.get $dest
        i32.load offset=0 align=1
        call $runtime.Block.Release
        i32.store offset=0 align=1
        local.get $dest
        local.get $item.0.1
        i32.store offset=4 align=4
        local.get $dest
        local.get $item.0.2
        i32.store offset=8 align=4
        local.get $dest
        local.get $item.0.3
        i32.store offset=12 align=4
        local.get $src
        i32.const 16
        i32.add
        local.set $src
        local.get $dest
        i32.const 16
        i32.add
        local.set $dest
        local.get $x_len
        i32.const 1
        i32.sub
        local.set $x_len
        br $loop2
      end ;;loop2
    end ;;block2
    local.get $y.1
    local.set $src
    block $block3
      loop $loop3
        local.get $y_len
        i32.eqz
        if
          br $block3
        else
        end
        local.get $src
        i32.load offset=0 align=4
        call $runtime.Block.Retain
        local.get $src
        i32.load offset=4 align=4
        local.get $src
        i32.load offset=8 align=4
        local.get $src
        i32.load offset=12 align=4
        local.set $item.0.3
        local.set $item.0.2
        local.set $item.0.1
        local.get $item.0.0
        call $runtime.Block.Release
        local.set $item.0.0
        local.get $dest
        local.get $item.0.0
        call $runtime.Block.Retain
        local.get $dest
        i32.load offset=0 align=1
        call $runtime.Block.Release
        i32.store offset=0 align=1
        local.get $dest
        local.get $item.0.1
        i32.store offset=4 align=4
        local.get $dest
        local.get $item.0.2
        i32.store offset=8 align=4
        local.get $dest
        local.get $item.0.3
        i32.store offset=12 align=4
        local.get $src
        i32.const 16
        i32.add
        local.set $src
        local.get $dest
        i32.const 16
        i32.add
        local.set $dest
        local.get $y_len
        i32.const 1
        i32.sub
        local.set $y_len
        br $loop3
      end ;;loop3
    end ;;block3
  end
  local.get $item.0.0
  call $runtime.Block.Release
) ;;$runtime.defers.$slice.append

(func $runtime.pushDeferStack
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$t0.0 i32)
  (local $$t0.1 i32)
  (local $$t0.2 i32)
  (local $$t0.3 i32)
  (local $$t1.0 i32)
  (local $$t1.1 i32)
  (local $$t2.0.0 i32)
  (local $$t2.0.1 i32)
  (local $$t2.0.2 i32)
  (local $$t2.0.3 i32)
  (local $$t3.0 i32)
  (local $$t3.1 i32)
  (local $$t4.0 i32)
  (local $$t4.1 i32)
  (local $$t5.0 i32)
  (local $$t5.1 i32)
  (local $$t5.2 i32)
  (local $$t5.3 i32)
  (local $$t6.0 i32)
  (local $$t6.1 i32)
  (local $$t6.2 i32)
  (local $$t6.3 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;*defersStack
        i32.const 8232
        i32.load offset=0 align=4
        call $runtime.Block.Retain
        i32.const 8232
        i32.load offset=4 align=4
        i32.const 8232
        i32.load offset=8 align=4
        i32.const 8232
        i32.load offset=12 align=4
        local.set $$t0.3
        local.set $$t0.2
        local.set $$t0.1
        local.get $$t0.0
        call $runtime.Block.Release
        local.set $$t0.0

        ;;local defers (complit)
        i32.const 32
        call $runtime.HeapAlloc
        i32.const 1
        i32.const 21
        i32.const 16
        call $runtime.Block.Init
        call $runtime.DupI32
        i32.const 16
        i32.add
        local.set $$t1.1
        local.get $$t1.0
        call $runtime.Block.Release
        local.set $$t1.0

        ;;*t1
        local.get $$t1.1
        i32.load offset=0 align=4
        call $runtime.Block.Retain
        local.get $$t1.1
        i32.load offset=4 align=4
        local.get $$t1.1
        i32.load offset=8 align=4
        local.get $$t1.1
        i32.load offset=12 align=4
        local.set $$t2.0.3
        local.set $$t2.0.2
        local.set $$t2.0.1
        local.get $$t2.0.0
        call $runtime.Block.Release
        local.set $$t2.0.0

        ;;new [1]defers (varargs)
        i32.const 32
        call $runtime.HeapAlloc
        i32.const 1
        i32.const 22
        i32.const 16
        call $runtime.Block.Init
        call $runtime.DupI32
        i32.const 16
        i32.add
        local.set $$t3.1
        local.get $$t3.0
        call $runtime.Block.Release
        local.set $$t3.0

        ;;&t3[0:int]
        local.get $$t3.0
        call $runtime.Block.Retain
        local.get $$t3.1
        i32.const 16
        i32.const 0
        i32.mul
        i32.add
        local.set $$t4.1
        local.get $$t4.0
        call $runtime.Block.Release
        local.set $$t4.0

        ;;*t4 = t2
        local.get $$t4.1
        local.get $$t2.0.0
        call $runtime.Block.Retain
        local.get $$t4.1
        i32.load offset=0 align=1
        call $runtime.Block.Release
        i32.store offset=0 align=1
        local.get $$t4.1
        local.get $$t2.0.1
        i32.store offset=4 align=4
        local.get $$t4.1
        local.get $$t2.0.2
        i32.store offset=8 align=4
        local.get $$t4.1
        local.get $$t2.0.3
        i32.store offset=12 align=4

        ;;slice t3[:]
        local.get $$t3.0
        call $runtime.Block.Retain
        local.get $$t3.1
        i32.const 16
        i32.const 0
        i32.mul
        i32.add
        i32.const 1
        i32.const 0
        i32.sub
        i32.const 1
        i32.const 0
        i32.sub
        local.set $$t5.3
        local.set $$t5.2
        local.set $$t5.1
        local.get $$t5.0
        call $runtime.Block.Release
        local.set $$t5.0

        ;;append(t0, t5...)
        local.get $$t0.0
        local.get $$t0.1
        local.get $$t0.2
        local.get $$t0.3
        local.get $$t5.0
        local.get $$t5.1
        local.get $$t5.2
        local.get $$t5.3
        call $$runtime.defers.$slice.append
        local.set $$t6.3
        local.set $$t6.2
        local.set $$t6.1
        local.get $$t6.0
        call $runtime.Block.Release
        local.set $$t6.0

        ;;*defersStack = t6
        i32.const 8232
        local.get $$t6.0
        call $runtime.Block.Retain
        i32.const 8232
        i32.load offset=0 align=1
        call $runtime.Block.Release
        i32.store offset=0 align=1
        i32.const 8232
        local.get $$t6.1
        i32.store offset=4 align=4
        i32.const 8232
        local.get $$t6.2
        i32.store offset=8 align=4
        i32.const 8232
        local.get $$t6.3
        i32.store offset=12 align=4

        ;;return
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$t0.0
  call $runtime.Block.Release
  local.get $$t1.0
  call $runtime.Block.Release
  local.get $$t2.0.0
  call $runtime.Block.Release
  local.get $$t3.0
  call $runtime.Block.Release
  local.get $$t4.0
  call $runtime.Block.Release
  local.get $$t5.0
  call $runtime.Block.Release
  local.get $$t6.0
  call $runtime.Block.Release
) ;;runtime.pushDeferStack

(func $runtime.refToPtr_byteSlice (param $t.0 i32) (param $t.1 i32) (param $t.2 i32) (param $t.3 i32) (result i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$ret_0 i32)
  (local $$t0 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;U8_slice_to_ptr(t)
        local.get $t.0
        local.get $t.1
        local.get $t.2
        local.get $t.3
        call $$wa.runtime.slice_to_ptr
        local.set $$t0

        ;;return t0
        local.get $$t0
        local.set $$ret_0
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$ret_0
) ;;runtime.refToPtr_byteSlice

(func $runtime.refToPtr_i32 (param $p.0 i32) (param $p.1 i32) (result i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$ret_0 i32)
  (local $$t0 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;I32_ref_to_ptr(p)
        local.get $p.0
        local.get $p.1
        call $$wa.runtime.i32_ref_to_ptr
        local.set $$t0

        ;;return t0
        local.get $$t0
        local.set $$ret_0
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$ret_0
) ;;runtime.refToPtr_i32

(func $runtime.refToPtr_string (param $s.0 i32) (param $s.1 i32) (param $s.2 i32) (result i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$ret_0 i32)
  (local $$t0 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;U8_string_to_ptr(s)
        local.get $s.0
        local.get $s.1
        local.get $s.2
        call $$wa.runtime.string_to_ptr
        local.set $$t0

        ;;return t0
        local.get $$t0
        local.set $$ret_0
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$ret_0
) ;;runtime.refToPtr_string

(func $$rune.$slice.append (param $x.0 i32) (param $x.1 i32) (param $x.2 i32) (param $x.3 i32) (param $y.0 i32) (param $y.1 i32) (param $y.2 i32) (param $y.3 i32) (result i32 i32 i32 i32)
  (local $item i32)
  (local $x_len i32)
  (local $y_len i32)
  (local $new_len i32)
  (local $src i32)
  (local $dest i32)
  (local $new_cap i32)
  local.get $x.2
  local.set $x_len
  local.get $y.2
  local.set $y_len
  local.get $x_len
  local.get $y_len
  i32.add
  local.set $new_len
  local.get $new_len
  local.get $x.3
  i32.le_u
  if (result i32 i32 i32 i32)
    local.get $x.0
    call $runtime.Block.Retain
    local.get $x.1
    local.get $new_len
    local.get $x.3
    local.get $y.1
    local.set $src
    local.get $x.1
    i32.const 4
    local.get $x_len
    i32.mul
    i32.add
    local.set $dest
    block $block1
      loop $loop1
        local.get $y_len
        i32.eqz
        if
          br $block1
        else
        end
        local.get $src
        i32.load offset=0 align=4
        local.set $item
        local.get $dest
        local.get $item
        i32.store offset=0 align=4
        local.get $src
        i32.const 4
        i32.add
        local.set $src
        local.get $dest
        i32.const 4
        i32.add
        local.set $dest
        local.get $y_len
        i32.const 1
        i32.sub
        local.set $y_len
        br $loop1
      end ;;loop1
    end ;;block1
  else
    local.get $new_len
    i32.const 2
    i32.mul
    local.set $new_cap
    local.get $new_cap
    i32.const 4
    i32.mul
    i32.const 16
    i32.add
    call $runtime.HeapAlloc
    local.get $new_cap
    i32.const 0
    i32.const 4
    call $runtime.Block.Init
    call $runtime.DupI32
    i32.const 16
    i32.add
    call $runtime.DupI32
    local.set $dest
    local.get $new_len
    local.get $new_cap
    local.get $x.1
    local.set $src
    block $block2
      loop $loop2
        local.get $x_len
        i32.eqz
        if
          br $block2
        else
        end
        local.get $src
        i32.load offset=0 align=4
        local.set $item
        local.get $dest
        local.get $item
        i32.store offset=0 align=4
        local.get $src
        i32.const 4
        i32.add
        local.set $src
        local.get $dest
        i32.const 4
        i32.add
        local.set $dest
        local.get $x_len
        i32.const 1
        i32.sub
        local.set $x_len
        br $loop2
      end ;;loop2
    end ;;block2
    local.get $y.1
    local.set $src
    block $block3
      loop $loop3
        local.get $y_len
        i32.eqz
        if
          br $block3
        else
        end
        local.get $src
        i32.load offset=0 align=4
        local.set $item
        local.get $dest
        local.get $item
        i32.store offset=0 align=4
        local.get $src
        i32.const 4
        i32.add
        local.set $src
        local.get $dest
        i32.const 4
        i32.add
        local.set $dest
        local.get $y_len
        i32.const 1
        i32.sub
        local.set $y_len
        br $loop3
      end ;;loop3
    end ;;block3
  end
) ;;$rune.$slice.append

(func $runtime.runeSliceFromString (param $s.0 i32) (param $s.1 i32) (param $s.2 i32) (result i32 i32 i32 i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$ret_0.0 i32)
  (local $$ret_0.1 i32)
  (local $$ret_0.2 i32)
  (local $$ret_0.3 i32)
  (local $$t0.0 i32)
  (local $$t0.1 i32)
  (local $$t0.2 i32)
  (local $$t1.0 i32)
  (local $$t1.1 i32)
  (local $$t1.2 i32)
  (local $$t1.3 i32)
  (local $$t2.0 i32)
  (local $$t2.1 i32)
  (local $$t2.2 i32)
  (local $$t2.3 i32)
  (local $$t3.0 i32)
  (local $$t3.1 i32)
  (local $$t3.2 i32)
  (local $$t4 i32)
  (local $$t5 i32)
  (local $$t6.0 i32)
  (local $$t6.1 i32)
  (local $$t7.0 i32)
  (local $$t7.1 i32)
  (local $$t8.0 i32)
  (local $$t8.1 i32)
  (local $$t8.2 i32)
  (local $$t8.3 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_3
        block $$Block_2
          block $$Block_1
            block $$Block_0
              block $$BlockSel
                local.get $$block_selector
                br_table 0 1 2 3 0
              end ;;$BlockSel
              i32.const 0
              local.set $$current_block

              ;;range s
              local.get $s.1
              local.get $s.2
              i32.const 0
              local.set $$t0.2
              local.set $$t0.1
              local.set $$t0.0

              ;;jump 1
              br $$Block_0

            end ;;$Block_0
            local.get $$current_block
            i32.const 0
            i32.eq
            if (result i32 i32 i32 i32)
              i32.const 0
              i32.const 0
              i32.const 0
              i32.const 0
            else
              local.get $$t1.0
              call $runtime.Block.Retain
              local.get $$t1.1
              local.get $$t1.2
              local.get $$t1.3
            end
            local.set $$t2.3
            local.set $$t2.2
            local.set $$t2.1
            local.get $$t2.0
            call $runtime.Block.Release
            local.set $$t2.0
            i32.const 1
            local.set $$current_block

            ;;next t0
            local.get $$t0.0
            local.get $$t0.1
            local.get $$t0.2
            call $runtime.next_rune
            local.set $$t0.2
            local.set $$t3.2
            local.set $$t3.1
            local.set $$t3.0

            ;;extract t2 #0
            local.get $$t3.0
            local.set $$t4

            ;;if t3 goto 2 else 3
            local.get $$t4
            if
              br $$Block_1
            else
              br $$Block_2
            end

          end ;;$Block_1
          i32.const 2
          local.set $$current_block

          ;;extract t2 #2
          local.get $$t3.2
          local.set $$t5

          ;;new [1]rune (varargs)
          i32.const 20
          call $runtime.HeapAlloc
          i32.const 1
          i32.const 0
          i32.const 4
          call $runtime.Block.Init
          call $runtime.DupI32
          i32.const 16
          i32.add
          local.set $$t6.1
          local.get $$t6.0
          call $runtime.Block.Release
          local.set $$t6.0

          ;;&t5[0:int]
          local.get $$t6.0
          call $runtime.Block.Retain
          local.get $$t6.1
          i32.const 4
          i32.const 0
          i32.mul
          i32.add
          local.set $$t7.1
          local.get $$t7.0
          call $runtime.Block.Release
          local.set $$t7.0

          ;;*t6 = t4
          local.get $$t7.1
          local.get $$t5
          i32.store offset=0 align=4

          ;;slice t5[:]
          local.get $$t6.0
          call $runtime.Block.Retain
          local.get $$t6.1
          i32.const 4
          i32.const 0
          i32.mul
          i32.add
          i32.const 1
          i32.const 0
          i32.sub
          i32.const 1
          i32.const 0
          i32.sub
          local.set $$t8.3
          local.set $$t8.2
          local.set $$t8.1
          local.get $$t8.0
          call $runtime.Block.Release
          local.set $$t8.0

          ;;append(t1, t7...)
          local.get $$t2.0
          local.get $$t2.1
          local.get $$t2.2
          local.get $$t2.3
          local.get $$t8.0
          local.get $$t8.1
          local.get $$t8.2
          local.get $$t8.3
          call $$rune.$slice.append
          local.set $$t1.3
          local.set $$t1.2
          local.set $$t1.1
          local.get $$t1.0
          call $runtime.Block.Release
          local.set $$t1.0

          ;;jump 1
          i32.const 1
          local.set $$block_selector
          br $$BlockDisp

        end ;;$Block_2
        i32.const 3
        local.set $$current_block

        ;;return t1
        local.get $$t2.0
        call $runtime.Block.Retain
        local.get $$t2.1
        local.get $$t2.2
        local.get $$t2.3
        local.set $$ret_0.3
        local.set $$ret_0.2
        local.set $$ret_0.1
        local.get $$ret_0.0
        call $runtime.Block.Release
        local.set $$ret_0.0
        br $$BlockFnBody

      end ;;$Block_3
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$ret_0.0
  call $runtime.Block.Retain
  local.get $$ret_0.1
  local.get $$ret_0.2
  local.get $$ret_0.3
  local.get $$ret_0.0
  call $runtime.Block.Release
  local.get $$t1.0
  call $runtime.Block.Release
  local.get $$t2.0
  call $runtime.Block.Release
  local.get $$t6.0
  call $runtime.Block.Release
  local.get $$t7.0
  call $runtime.Block.Release
  local.get $$t8.0
  call $runtime.Block.Release
) ;;runtime.runeSliceFromString

(func $runtime.set_u8 (param $addr i32) (param $data i32)
  local.get $addr
  local.get $data
  i32.store8 offset=0 align=1
) ;;runtime.set_u8

(func $runtime.stringFromRune (param $r i32) (result i32 i32 i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$ret_0.0 i32)
  (local $$ret_0.1 i32)
  (local $$ret_0.2 i32)
  (local $$t0.0 i32)
  (local $$t0.1 i32)
  (local $$t1.0 i32)
  (local $$t1.1 i32)
  (local $$t1.2 i32)
  (local $$t1.3 i32)
  (local $$t2 i32)
  (local $$t3 i32)
  (local $$t4 i32)
  (local $$t5.0 i32)
  (local $$t5.1 i32)
  (local $$t5.2 i32)
  (local $$t5.3 i32)
  (local $$t6.0 i32)
  (local $$t6.1 i32)
  (local $$t6.2 i32)
  (local $$t7.0 i32)
  (local $$t7.1 i32)
  (local $$t8 i32)
  (local $$t9.0 i32)
  (local $$t9.1 i32)
  (local $$t10 i32)
  (local $$t11.0 i32)
  (local $$t11.1 i32)
  (local $$t12 i32)
  (local $$t13 i32)
  (local $$t14 i32)
  (local $$t15.0 i32)
  (local $$t15.1 i32)
  (local $$t16 i32)
  (local $$t17 i32)
  (local $$t18 i32)
  (local $$t19 i32)
  (local $$t20.0 i32)
  (local $$t20.1 i32)
  (local $$t21 i32)
  (local $$t22.0 i32)
  (local $$t22.1 i32)
  (local $$t23 i32)
  (local $$t24 i32)
  (local $$t25 i32)
  (local $$t26.0 i32)
  (local $$t26.1 i32)
  (local $$t27 i32)
  (local $$t28 i32)
  (local $$t29 i32)
  (local $$t30 i32)
  (local $$t31.0 i32)
  (local $$t31.1 i32)
  (local $$t32 i32)
  (local $$t33 i32)
  (local $$t34 i32)
  (local $$t35 i32)
  (local $$t36.0 i32)
  (local $$t36.1 i32)
  (local $$t37 i32)
  (local $$t38.0 i32)
  (local $$t38.1 i32)
  (local $$t39 i32)
  (local $$t40 i32)
  (local $$t41 i32)
  (local $$t42.0 i32)
  (local $$t42.1 i32)
  (local $$t43 i32)
  (local $$t44 i32)
  (local $$t45 i32)
  (local $$t46 i32)
  (local $$t47.0 i32)
  (local $$t47.1 i32)
  (local $$t48 i32)
  (local $$t49 i32)
  (local $$t50 i32)
  (local $$t51 i32)
  (local $$t52 i32)
  (local $$t53 i32)
  (local $$t54 i32)
  (local $$t55.0 i32)
  (local $$t55.1 i32)
  (local $$t56 i32)
  (local $$t57.0 i32)
  (local $$t57.1 i32)
  (local $$t58 i32)
  (local $$t59 i32)
  (local $$t60 i32)
  (local $$t61.0 i32)
  (local $$t61.1 i32)
  (local $$t62 i32)
  (local $$t63 i32)
  (local $$t64 i32)
  (local $$t65 i32)
  (local $$t66.0 i32)
  (local $$t66.1 i32)
  (local $$t67 i32)
  (local $$t68 i32)
  (local $$t69 i32)
  (local $$t70 i32)
  (local $$t71.0 i32)
  (local $$t71.1 i32)
  (local $$t72 i32)
  (local $$t73 i32)
  (local $$t74 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_12
        block $$Block_11
          block $$Block_10
            block $$Block_9
              block $$Block_8
                block $$Block_7
                  block $$Block_6
                    block $$Block_5
                      block $$Block_4
                        block $$Block_3
                          block $$Block_2
                            block $$Block_1
                              block $$Block_0
                                block $$BlockSel
                                  local.get $$block_selector
                                  br_table 0 1 2 3 4 5 6 7 8 9 10 11 12 0
                                end ;;$BlockSel
                                i32.const 0
                                local.set $$current_block

                                ;;new [4]byte (makeslice)
                                i32.const 20
                                call $runtime.HeapAlloc
                                i32.const 1
                                i32.const 0
                                i32.const 4
                                call $runtime.Block.Init
                                call $runtime.DupI32
                                i32.const 16
                                i32.add
                                local.set $$t0.1
                                local.get $$t0.0
                                call $runtime.Block.Release
                                local.set $$t0.0

                                ;;slice t0[:0:int]
                                local.get $$t0.0
                                call $runtime.Block.Retain
                                local.get $$t0.1
                                i32.const 1
                                i32.const 0
                                i32.mul
                                i32.add
                                i32.const 0
                                i32.const 0
                                i32.sub
                                i32.const 4
                                i32.const 0
                                i32.sub
                                local.set $$t1.3
                                local.set $$t1.2
                                local.set $$t1.1
                                local.get $$t1.0
                                call $runtime.Block.Release
                                local.set $$t1.0

                                ;;convert uint32 <- rune (r)
                                local.get $r
                                local.set $$t2

                                ;;t2 <= 127:uint32
                                local.get $$t2
                                i32.const 127
                                i32.le_u
                                local.set $$t3

                                ;;if t3 goto 2 else 4
                                local.get $$t3
                                if
                                  br $$Block_1
                                else
                                  br $$Block_3
                                end

                              end ;;$Block_0
                              local.get $$current_block
                              i32.const 2
                              i32.eq
                              if (result i32)
                                i32.const 1
                              else
                                local.get $$current_block
                                i32.const 3
                                i32.eq
                                if (result i32)
                                  i32.const 2
                                else
                                  local.get $$current_block
                                  i32.const 5
                                  i32.eq
                                  if (result i32)
                                    i32.const 3
                                  else
                                    local.get $$current_block
                                    i32.const 7
                                    i32.eq
                                    if (result i32)
                                      i32.const 3
                                    else
                                      i32.const 4
                                    end
                                  end
                                end
                              end
                              local.set $$t4
                              i32.const 1
                              local.set $$current_block

                              ;;slice t1[:t4]
                              local.get $$t1.0
                              call $runtime.Block.Retain
                              local.get $$t1.1
                              i32.const 1
                              i32.const 0
                              i32.mul
                              i32.add
                              local.get $$t4
                              i32.const 0
                              i32.sub
                              local.get $$t1.3
                              i32.const 0
                              i32.sub
                              local.set $$t5.3
                              local.set $$t5.2
                              local.set $$t5.1
                              local.get $$t5.0
                              call $runtime.Block.Release
                              local.set $$t5.0

                              ;;convert string <- []byte (t5)
                              i32.const 0
                              i32.const 8224
                              i32.const 0
                              local.get $$t5.0
                              local.get $$t5.1
                              local.get $$t5.2
                              call $$string.appendstr
                              local.set $$t6.2
                              local.set $$t6.1
                              local.get $$t6.0
                              call $runtime.Block.Release
                              local.set $$t6.0

                              ;;return t6
                              local.get $$t6.0
                              call $runtime.Block.Retain
                              local.get $$t6.1
                              local.get $$t6.2
                              local.set $$ret_0.2
                              local.set $$ret_0.1
                              local.get $$ret_0.0
                              call $runtime.Block.Release
                              local.set $$ret_0.0
                              br $$BlockFnBody

                            end ;;$Block_1
                            i32.const 2
                            local.set $$current_block

                            ;;&t1[0:int]
                            local.get $$t1.0
                            call $runtime.Block.Retain
                            local.get $$t1.1
                            i32.const 1
                            i32.const 0
                            i32.mul
                            i32.add
                            local.set $$t7.1
                            local.get $$t7.0
                            call $runtime.Block.Release
                            local.set $$t7.0

                            ;;convert byte <- rune (r)
                            local.get $r
                            i32.const 255
                            i32.and
                            local.set $$t8

                            ;;*t7 = t8
                            local.get $$t7.1
                            local.get $$t8
                            i32.store8 offset=0 align=1

                            ;;jump 1
                            i32.const 1
                            local.set $$block_selector
                            br $$BlockDisp

                          end ;;$Block_2
                          i32.const 3
                          local.set $$current_block

                          ;;&t1[1:int]
                          local.get $$t1.0
                          call $runtime.Block.Retain
                          local.get $$t1.1
                          i32.const 1
                          i32.const 1
                          i32.mul
                          i32.add
                          local.set $$t9.1
                          local.get $$t9.0
                          call $runtime.Block.Release
                          local.set $$t9.0

                          ;;*t9
                          local.get $$t9.1
                          i32.load8_u offset=0 align=1
                          local.set $$t10

                          ;;&t1[0:int]
                          local.get $$t1.0
                          call $runtime.Block.Retain
                          local.get $$t1.1
                          i32.const 1
                          i32.const 0
                          i32.mul
                          i32.add
                          local.set $$t11.1
                          local.get $$t11.0
                          call $runtime.Block.Release
                          local.set $$t11.0

                          ;;r >> 6:uint64
                          local.get $r
                          i64.const 6
                          i32.wrap_i64
                          i32.shr_s
                          local.set $$t12

                          ;;convert byte <- rune (t12)
                          local.get $$t12
                          i32.const 255
                          i32.and
                          local.set $$t13

                          ;;192:byte | t13
                          i32.const 192
                          local.get $$t13
                          i32.or
                          local.set $$t14

                          ;;*t11 = t14
                          local.get $$t11.1
                          local.get $$t14
                          i32.store8 offset=0 align=1

                          ;;&t1[1:int]
                          local.get $$t1.0
                          call $runtime.Block.Retain
                          local.get $$t1.1
                          i32.const 1
                          i32.const 1
                          i32.mul
                          i32.add
                          local.set $$t15.1
                          local.get $$t15.0
                          call $runtime.Block.Release
                          local.set $$t15.0

                          ;;convert byte <- rune (r)
                          local.get $r
                          i32.const 255
                          i32.and
                          local.set $$t16

                          ;;t16 & 63:byte
                          local.get $$t16
                          i32.const 63
                          i32.and
                          local.set $$t17

                          ;;128:byte | t17
                          i32.const 128
                          local.get $$t17
                          i32.or
                          local.set $$t18

                          ;;*t15 = t18
                          local.get $$t15.1
                          local.get $$t18
                          i32.store8 offset=0 align=1

                          ;;jump 1
                          i32.const 1
                          local.set $$block_selector
                          br $$BlockDisp

                        end ;;$Block_3
                        i32.const 4
                        local.set $$current_block

                        ;;t2 <= 2047:uint32
                        local.get $$t2
                        i32.const 2047
                        i32.le_u
                        local.set $$t19

                        ;;if t19 goto 3 else 6
                        local.get $$t19
                        if
                          i32.const 3
                          local.set $$block_selector
                          br $$BlockDisp
                        else
                          br $$Block_5
                        end

                      end ;;$Block_4
                      i32.const 5
                      local.set $$current_block

                      ;;&t1[2:int]
                      local.get $$t1.0
                      call $runtime.Block.Retain
                      local.get $$t1.1
                      i32.const 1
                      i32.const 2
                      i32.mul
                      i32.add
                      local.set $$t20.1
                      local.get $$t20.0
                      call $runtime.Block.Release
                      local.set $$t20.0

                      ;;*t20
                      local.get $$t20.1
                      i32.load8_u offset=0 align=1
                      local.set $$t21

                      ;;&t1[0:int]
                      local.get $$t1.0
                      call $runtime.Block.Retain
                      local.get $$t1.1
                      i32.const 1
                      i32.const 0
                      i32.mul
                      i32.add
                      local.set $$t22.1
                      local.get $$t22.0
                      call $runtime.Block.Release
                      local.set $$t22.0

                      ;;65533:rune >> 12:uint64
                      i32.const 65533
                      i64.const 12
                      i32.wrap_i64
                      i32.shr_s
                      local.set $$t23

                      ;;convert byte <- rune (t23)
                      local.get $$t23
                      i32.const 255
                      i32.and
                      local.set $$t24

                      ;;224:byte | t24
                      i32.const 224
                      local.get $$t24
                      i32.or
                      local.set $$t25

                      ;;*t22 = t25
                      local.get $$t22.1
                      local.get $$t25
                      i32.store8 offset=0 align=1

                      ;;&t1[1:int]
                      local.get $$t1.0
                      call $runtime.Block.Retain
                      local.get $$t1.1
                      i32.const 1
                      i32.const 1
                      i32.mul
                      i32.add
                      local.set $$t26.1
                      local.get $$t26.0
                      call $runtime.Block.Release
                      local.set $$t26.0

                      ;;65533:rune >> 6:uint64
                      i32.const 65533
                      i64.const 6
                      i32.wrap_i64
                      i32.shr_s
                      local.set $$t27

                      ;;convert byte <- rune (t27)
                      local.get $$t27
                      i32.const 255
                      i32.and
                      local.set $$t28

                      ;;t28 & 63:byte
                      local.get $$t28
                      i32.const 63
                      i32.and
                      local.set $$t29

                      ;;128:byte | t29
                      i32.const 128
                      local.get $$t29
                      i32.or
                      local.set $$t30

                      ;;*t26 = t30
                      local.get $$t26.1
                      local.get $$t30
                      i32.store8 offset=0 align=1

                      ;;&t1[2:int]
                      local.get $$t1.0
                      call $runtime.Block.Retain
                      local.get $$t1.1
                      i32.const 1
                      i32.const 2
                      i32.mul
                      i32.add
                      local.set $$t31.1
                      local.get $$t31.0
                      call $runtime.Block.Release
                      local.set $$t31.0

                      ;;convert byte <- rune (65533:rune)
                      i32.const 65533
                      i32.const 255
                      i32.and
                      local.set $$t32

                      ;;t32 & 63:byte
                      local.get $$t32
                      i32.const 63
                      i32.and
                      local.set $$t33

                      ;;128:byte | t33
                      i32.const 128
                      local.get $$t33
                      i32.or
                      local.set $$t34

                      ;;*t31 = t34
                      local.get $$t31.1
                      local.get $$t34
                      i32.store8 offset=0 align=1

                      ;;jump 1
                      i32.const 1
                      local.set $$block_selector
                      br $$BlockDisp

                    end ;;$Block_5
                    i32.const 6
                    local.set $$current_block

                    ;;t2 > 1114111:uint32
                    local.get $$t2
                    i32.const 1114111
                    i32.gt_u
                    local.set $$t35

                    ;;if t35 goto 5 else 8
                    local.get $$t35
                    if
                      i32.const 5
                      local.set $$block_selector
                      br $$BlockDisp
                    else
                      br $$Block_7
                    end

                  end ;;$Block_6
                  i32.const 7
                  local.set $$current_block

                  ;;&t1[2:int]
                  local.get $$t1.0
                  call $runtime.Block.Retain
                  local.get $$t1.1
                  i32.const 1
                  i32.const 2
                  i32.mul
                  i32.add
                  local.set $$t36.1
                  local.get $$t36.0
                  call $runtime.Block.Release
                  local.set $$t36.0

                  ;;*t36
                  local.get $$t36.1
                  i32.load8_u offset=0 align=1
                  local.set $$t37

                  ;;&t1[0:int]
                  local.get $$t1.0
                  call $runtime.Block.Retain
                  local.get $$t1.1
                  i32.const 1
                  i32.const 0
                  i32.mul
                  i32.add
                  local.set $$t38.1
                  local.get $$t38.0
                  call $runtime.Block.Release
                  local.set $$t38.0

                  ;;r >> 12:uint64
                  local.get $r
                  i64.const 12
                  i32.wrap_i64
                  i32.shr_s
                  local.set $$t39

                  ;;convert byte <- rune (t39)
                  local.get $$t39
                  i32.const 255
                  i32.and
                  local.set $$t40

                  ;;224:byte | t40
                  i32.const 224
                  local.get $$t40
                  i32.or
                  local.set $$t41

                  ;;*t38 = t41
                  local.get $$t38.1
                  local.get $$t41
                  i32.store8 offset=0 align=1

                  ;;&t1[1:int]
                  local.get $$t1.0
                  call $runtime.Block.Retain
                  local.get $$t1.1
                  i32.const 1
                  i32.const 1
                  i32.mul
                  i32.add
                  local.set $$t42.1
                  local.get $$t42.0
                  call $runtime.Block.Release
                  local.set $$t42.0

                  ;;r >> 6:uint64
                  local.get $r
                  i64.const 6
                  i32.wrap_i64
                  i32.shr_s
                  local.set $$t43

                  ;;convert byte <- rune (t43)
                  local.get $$t43
                  i32.const 255
                  i32.and
                  local.set $$t44

                  ;;t44 & 63:byte
                  local.get $$t44
                  i32.const 63
                  i32.and
                  local.set $$t45

                  ;;128:byte | t45
                  i32.const 128
                  local.get $$t45
                  i32.or
                  local.set $$t46

                  ;;*t42 = t46
                  local.get $$t42.1
                  local.get $$t46
                  i32.store8 offset=0 align=1

                  ;;&t1[2:int]
                  local.get $$t1.0
                  call $runtime.Block.Retain
                  local.get $$t1.1
                  i32.const 1
                  i32.const 2
                  i32.mul
                  i32.add
                  local.set $$t47.1
                  local.get $$t47.0
                  call $runtime.Block.Release
                  local.set $$t47.0

                  ;;convert byte <- rune (r)
                  local.get $r
                  i32.const 255
                  i32.and
                  local.set $$t48

                  ;;t48 & 63:byte
                  local.get $$t48
                  i32.const 63
                  i32.and
                  local.set $$t49

                  ;;128:byte | t49
                  i32.const 128
                  local.get $$t49
                  i32.or
                  local.set $$t50

                  ;;*t47 = t50
                  local.get $$t47.1
                  local.get $$t50
                  i32.store8 offset=0 align=1

                  ;;jump 1
                  i32.const 1
                  local.set $$block_selector
                  br $$BlockDisp

                end ;;$Block_7
                i32.const 8
                local.set $$current_block

                ;;55296:uint32 <= t2
                i32.const 55296
                local.get $$t2
                i32.le_u
                local.set $$t51

                ;;if t51 goto 10 else 11
                local.get $$t51
                if
                  br $$Block_9
                else
                  br $$Block_10
                end

              end ;;$Block_8
              i32.const 9
              local.set $$current_block

              ;;t2 <= 65535:uint32
              local.get $$t2
              i32.const 65535
              i32.le_u
              local.set $$t52

              ;;if t52 goto 7 else 12
              local.get $$t52
              if
                i32.const 7
                local.set $$block_selector
                br $$BlockDisp
              else
                br $$Block_11
              end

            end ;;$Block_9
            i32.const 10
            local.set $$current_block

            ;;t2 <= 57343:uint32
            local.get $$t2
            i32.const 57343
            i32.le_u
            local.set $$t53

            ;;jump 11
            br $$Block_10

          end ;;$Block_10
          local.get $$current_block
          i32.const 8
          i32.eq
          if (result i32)
            i32.const 0
          else
            local.get $$t53
          end
          local.set $$t54
          i32.const 11
          local.set $$current_block

          ;;if t54 goto 5 else 9
          local.get $$t54
          if
            i32.const 5
            local.set $$block_selector
            br $$BlockDisp
          else
            i32.const 9
            local.set $$block_selector
            br $$BlockDisp
          end

        end ;;$Block_11
        i32.const 12
        local.set $$current_block

        ;;&t1[3:int]
        local.get $$t1.0
        call $runtime.Block.Retain
        local.get $$t1.1
        i32.const 1
        i32.const 3
        i32.mul
        i32.add
        local.set $$t55.1
        local.get $$t55.0
        call $runtime.Block.Release
        local.set $$t55.0

        ;;*t55
        local.get $$t55.1
        i32.load8_u offset=0 align=1
        local.set $$t56

        ;;&t1[0:int]
        local.get $$t1.0
        call $runtime.Block.Retain
        local.get $$t1.1
        i32.const 1
        i32.const 0
        i32.mul
        i32.add
        local.set $$t57.1
        local.get $$t57.0
        call $runtime.Block.Release
        local.set $$t57.0

        ;;r >> 18:uint64
        local.get $r
        i64.const 18
        i32.wrap_i64
        i32.shr_s
        local.set $$t58

        ;;convert byte <- rune (t58)
        local.get $$t58
        i32.const 255
        i32.and
        local.set $$t59

        ;;240:byte | t59
        i32.const 240
        local.get $$t59
        i32.or
        local.set $$t60

        ;;*t57 = t60
        local.get $$t57.1
        local.get $$t60
        i32.store8 offset=0 align=1

        ;;&t1[1:int]
        local.get $$t1.0
        call $runtime.Block.Retain
        local.get $$t1.1
        i32.const 1
        i32.const 1
        i32.mul
        i32.add
        local.set $$t61.1
        local.get $$t61.0
        call $runtime.Block.Release
        local.set $$t61.0

        ;;r >> 12:uint64
        local.get $r
        i64.const 12
        i32.wrap_i64
        i32.shr_s
        local.set $$t62

        ;;convert byte <- rune (t62)
        local.get $$t62
        i32.const 255
        i32.and
        local.set $$t63

        ;;t63 & 63:byte
        local.get $$t63
        i32.const 63
        i32.and
        local.set $$t64

        ;;128:byte | t64
        i32.const 128
        local.get $$t64
        i32.or
        local.set $$t65

        ;;*t61 = t65
        local.get $$t61.1
        local.get $$t65
        i32.store8 offset=0 align=1

        ;;&t1[2:int]
        local.get $$t1.0
        call $runtime.Block.Retain
        local.get $$t1.1
        i32.const 1
        i32.const 2
        i32.mul
        i32.add
        local.set $$t66.1
        local.get $$t66.0
        call $runtime.Block.Release
        local.set $$t66.0

        ;;r >> 6:uint64
        local.get $r
        i64.const 6
        i32.wrap_i64
        i32.shr_s
        local.set $$t67

        ;;convert byte <- rune (t67)
        local.get $$t67
        i32.const 255
        i32.and
        local.set $$t68

        ;;t68 & 63:byte
        local.get $$t68
        i32.const 63
        i32.and
        local.set $$t69

        ;;128:byte | t69
        i32.const 128
        local.get $$t69
        i32.or
        local.set $$t70

        ;;*t66 = t70
        local.get $$t66.1
        local.get $$t70
        i32.store8 offset=0 align=1

        ;;&t1[3:int]
        local.get $$t1.0
        call $runtime.Block.Retain
        local.get $$t1.1
        i32.const 1
        i32.const 3
        i32.mul
        i32.add
        local.set $$t71.1
        local.get $$t71.0
        call $runtime.Block.Release
        local.set $$t71.0

        ;;convert byte <- rune (r)
        local.get $r
        i32.const 255
        i32.and
        local.set $$t72

        ;;t72 & 63:byte
        local.get $$t72
        i32.const 63
        i32.and
        local.set $$t73

        ;;128:byte | t73
        i32.const 128
        local.get $$t73
        i32.or
        local.set $$t74

        ;;*t71 = t74
        local.get $$t71.1
        local.get $$t74
        i32.store8 offset=0 align=1

        ;;jump 1
        i32.const 1
        local.set $$block_selector
        br $$BlockDisp

      end ;;$Block_12
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$ret_0.0
  call $runtime.Block.Retain
  local.get $$ret_0.1
  local.get $$ret_0.2
  local.get $$ret_0.0
  call $runtime.Block.Release
  local.get $$t0.0
  call $runtime.Block.Release
  local.get $$t1.0
  call $runtime.Block.Release
  local.get $$t5.0
  call $runtime.Block.Release
  local.get $$t6.0
  call $runtime.Block.Release
  local.get $$t7.0
  call $runtime.Block.Release
  local.get $$t9.0
  call $runtime.Block.Release
  local.get $$t11.0
  call $runtime.Block.Release
  local.get $$t15.0
  call $runtime.Block.Release
  local.get $$t20.0
  call $runtime.Block.Release
  local.get $$t22.0
  call $runtime.Block.Release
  local.get $$t26.0
  call $runtime.Block.Release
  local.get $$t31.0
  call $runtime.Block.Release
  local.get $$t36.0
  call $runtime.Block.Release
  local.get $$t38.0
  call $runtime.Block.Release
  local.get $$t42.0
  call $runtime.Block.Release
  local.get $$t47.0
  call $runtime.Block.Release
  local.get $$t55.0
  call $runtime.Block.Release
  local.get $$t57.0
  call $runtime.Block.Release
  local.get $$t61.0
  call $runtime.Block.Release
  local.get $$t66.0
  call $runtime.Block.Release
  local.get $$t71.0
  call $runtime.Block.Release
) ;;runtime.stringFromRune

(func $runtime.stringFromRuneSlice (param $rs.0 i32) (param $rs.1 i32) (param $rs.2 i32) (param $rs.3 i32) (result i32 i32 i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$ret_0.0 i32)
  (local $$ret_0.1 i32)
  (local $$ret_0.2 i32)
  (local $$t0 i32)
  (local $$t1.0 i32)
  (local $$t1.1 i32)
  (local $$t1.2 i32)
  (local $$t2.0 i32)
  (local $$t2.1 i32)
  (local $$t2.2 i32)
  (local $$t3 i32)
  (local $$t4 i32)
  (local $$t5 i32)
  (local $$t6.0 i32)
  (local $$t6.1 i32)
  (local $$t7 i32)
  (local $$t8.0 i32)
  (local $$t8.1 i32)
  (local $$t8.2 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_3
        block $$Block_2
          block $$Block_1
            block $$Block_0
              block $$BlockSel
                local.get $$block_selector
                br_table 0 1 2 3 0
              end ;;$BlockSel
              i32.const 0
              local.set $$current_block

              ;;len(rs)
              local.get $rs.2
              local.set $$t0

              ;;jump 1
              br $$Block_0

            end ;;$Block_0
            local.get $$current_block
            i32.const 0
            i32.eq
            if (result i32 i32 i32)
              i32.const 0
              i32.const 8224
              i32.const 0
            else
              local.get $$t1.0
              call $runtime.Block.Retain
              local.get $$t1.1
              local.get $$t1.2
            end
            local.get $$current_block
            i32.const 0
            i32.eq
            if (result i32)
              i32.const -1
            else
              local.get $$t3
            end
            local.set $$t4
            local.set $$t2.2
            local.set $$t2.1
            local.get $$t2.0
            call $runtime.Block.Release
            local.set $$t2.0
            i32.const 1
            local.set $$current_block

            ;;t2 + 1:int
            local.get $$t4
            i32.const 1
            i32.add
            local.set $$t3

            ;;t3 < t0
            local.get $$t3
            local.get $$t0
            i32.lt_s
            local.set $$t5

            ;;if t4 goto 2 else 3
            local.get $$t5
            if
              br $$Block_1
            else
              br $$Block_2
            end

          end ;;$Block_1
          i32.const 2
          local.set $$current_block

          ;;&rs[t3]
          local.get $rs.0
          call $runtime.Block.Retain
          local.get $rs.1
          i32.const 4
          local.get $$t3
          i32.mul
          i32.add
          local.set $$t6.1
          local.get $$t6.0
          call $runtime.Block.Release
          local.set $$t6.0

          ;;*t5
          local.get $$t6.1
          i32.load offset=0 align=4
          local.set $$t7

          ;;stringFromRune(t6)
          local.get $$t7
          call $runtime.stringFromRune
          local.set $$t8.2
          local.set $$t8.1
          local.get $$t8.0
          call $runtime.Block.Release
          local.set $$t8.0

          ;;t1 + t7
          local.get $$t2.0
          local.get $$t2.1
          local.get $$t2.2
          local.get $$t8.0
          local.get $$t8.1
          local.get $$t8.2
          call $$string.appendstr
          local.set $$t1.2
          local.set $$t1.1
          local.get $$t1.0
          call $runtime.Block.Release
          local.set $$t1.0

          ;;jump 1
          i32.const 1
          local.set $$block_selector
          br $$BlockDisp

        end ;;$Block_2
        i32.const 3
        local.set $$current_block

        ;;return t1
        local.get $$t2.0
        call $runtime.Block.Retain
        local.get $$t2.1
        local.get $$t2.2
        local.set $$ret_0.2
        local.set $$ret_0.1
        local.get $$ret_0.0
        call $runtime.Block.Release
        local.set $$ret_0.0
        br $$BlockFnBody

      end ;;$Block_3
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$ret_0.0
  call $runtime.Block.Retain
  local.get $$ret_0.1
  local.get $$ret_0.2
  local.get $$ret_0.0
  call $runtime.Block.Release
  local.get $$t1.0
  call $runtime.Block.Release
  local.get $$t2.0
  call $runtime.Block.Release
  local.get $$t6.0
  call $runtime.Block.Release
  local.get $$t8.0
  call $runtime.Block.Release
) ;;runtime.stringFromRuneSlice

(func $$wa.runtime.string_Comp (param $x.0 i32) (param $x.1 i32) (param $x.2 i32) (param $y.0 i32) (param $y.1 i32) (param $y.2 i32) (result i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$ret_0 i32)
  (local $$t0.0 i32)
  (local $$t0.1 i32)
  (local $$t1.0 i32)
  (local $$t1.1 i32)
  (local $$t1.2 i32)
  (local $$t2.0 i32)
  (local $$t2.1 i32)
  (local $$t3.0 i32)
  (local $$t3.1 i32)
  (local $$t3.2 i32)
  (local $$t4.0 i32)
  (local $$t4.1 i32)
  (local $$t4.2 i32)
  (local $$t5.0 i32)
  (local $$t5.1 i32)
  (local $$t5.2 i32)
  (local $$t5.3 i32)
  (local $$t6 i32)
  (local $$t7 i32)
  (local $$t8 i32)
  (local $$t9 i32)
  (local $$t10.0 i32)
  (local $$t10.1 i32)
  (local $$t11.0 i32)
  (local $$t11.1 i32)
  (local $$t11.2 i32)
  (local $$t12.0 i32)
  (local $$t12.1 i32)
  (local $$t12.2 i32)
  (local $$t12.3 i32)
  (local $$t13 i32)
  (local $$t14 i32)
  (local $$t15 i32)
  (local $$t16 i32)
  (local $$t17.0 i32)
  (local $$t17.1 i32)
  (local $$t18 i32)
  (local $$t19 i32)
  (local $$t20 i32)
  (local $$t21 i32)
  (local $$t22 i32)
  (local $$t23 i32)
  (local $$t24 i32)
  (local $$t25 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_12
        block $$Block_11
          block $$Block_10
            block $$Block_9
              block $$Block_8
                block $$Block_7
                  block $$Block_6
                    block $$Block_5
                      block $$Block_4
                        block $$Block_3
                          block $$Block_2
                            block $$Block_1
                              block $$Block_0
                                block $$BlockSel
                                  local.get $$block_selector
                                  br_table 0 1 2 3 4 5 6 7 8 9 10 11 12 0
                                end ;;$BlockSel
                                i32.const 0
                                local.set $$current_block

                                ;;local stringIter (i1)
                                i32.const 28
                                call $runtime.HeapAlloc
                                i32.const 1
                                i32.const 0
                                i32.const 12
                                call $runtime.Block.Init
                                call $runtime.DupI32
                                i32.const 16
                                i32.add
                                local.set $$t0.1
                                local.get $$t0.0
                                call $runtime.Block.Release
                                local.set $$t0.0

                                ;;stringToIter(x)
                                local.get $x.0
                                local.get $x.1
                                local.get $x.2
                                call $$wa.runtime.string_to_iter
                                local.set $$t1.2
                                local.set $$t1.1
                                local.set $$t1.0

                                ;;*t0 = t1
                                local.get $$t0.1
                                local.get $$t1.0
                                i32.store offset=0 align=4
                                local.get $$t0.1
                                local.get $$t1.1
                                i32.store offset=4 align=4
                                local.get $$t0.1
                                local.get $$t1.2
                                i32.store offset=8 align=4

                                ;;local stringIter (i2)
                                i32.const 28
                                call $runtime.HeapAlloc
                                i32.const 1
                                i32.const 0
                                i32.const 12
                                call $runtime.Block.Init
                                call $runtime.DupI32
                                i32.const 16
                                i32.add
                                local.set $$t2.1
                                local.get $$t2.0
                                call $runtime.Block.Release
                                local.set $$t2.0

                                ;;stringToIter(y)
                                local.get $y.0
                                local.get $y.1
                                local.get $y.2
                                call $$wa.runtime.string_to_iter
                                local.set $$t3.2
                                local.set $$t3.1
                                local.set $$t3.0

                                ;;*t2 = t3
                                local.get $$t2.1
                                local.get $$t3.0
                                i32.store offset=0 align=4
                                local.get $$t2.1
                                local.get $$t3.1
                                i32.store offset=4 align=4
                                local.get $$t2.1
                                local.get $$t3.2
                                i32.store offset=8 align=4

                                ;;jump 1
                                br $$Block_0

                              end ;;$Block_0
                              i32.const 1
                              local.set $$current_block

                              ;;*t0
                              local.get $$t0.1
                              i32.load offset=0 align=4
                              local.get $$t0.1
                              i32.load offset=4 align=4
                              local.get $$t0.1
                              i32.load offset=8 align=4
                              local.set $$t4.2
                              local.set $$t4.1
                              local.set $$t4.0

                              ;;next_rune(t4)
                              local.get $$t4.0
                              local.get $$t4.1
                              local.get $$t4.2
                              call $runtime.next_rune
                              local.set $$t5.3
                              local.set $$t5.2
                              local.set $$t5.1
                              local.set $$t5.0

                              ;;extract t5 #0
                              local.get $$t5.0
                              local.set $$t6

                              ;;extract t5 #1
                              local.get $$t5.1
                              local.set $$t7

                              ;;extract t5 #2
                              local.get $$t5.2
                              local.set $$t8

                              ;;extract t5 #3
                              local.get $$t5.3
                              local.set $$t9

                              ;;&t0.pos [#2]
                              local.get $$t0.0
                              call $runtime.Block.Retain
                              local.get $$t0.1
                              i32.const 8
                              i32.add
                              local.set $$t10.1
                              local.get $$t10.0
                              call $runtime.Block.Release
                              local.set $$t10.0

                              ;;*t10 = t9
                              local.get $$t10.1
                              local.get $$t9
                              i32.store offset=0 align=4

                              ;;*t2
                              local.get $$t2.1
                              i32.load offset=0 align=4
                              local.get $$t2.1
                              i32.load offset=4 align=4
                              local.get $$t2.1
                              i32.load offset=8 align=4
                              local.set $$t11.2
                              local.set $$t11.1
                              local.set $$t11.0

                              ;;next_rune(t11)
                              local.get $$t11.0
                              local.get $$t11.1
                              local.get $$t11.2
                              call $runtime.next_rune
                              local.set $$t12.3
                              local.set $$t12.2
                              local.set $$t12.1
                              local.set $$t12.0

                              ;;extract t12 #0
                              local.get $$t12.0
                              local.set $$t13

                              ;;extract t12 #1
                              local.get $$t12.1
                              local.set $$t14

                              ;;extract t12 #2
                              local.get $$t12.2
                              local.set $$t15

                              ;;extract t12 #3
                              local.get $$t12.3
                              local.set $$t16

                              ;;&t2.pos [#2]
                              local.get $$t2.0
                              call $runtime.Block.Retain
                              local.get $$t2.1
                              i32.const 8
                              i32.add
                              local.set $$t17.1
                              local.get $$t17.0
                              call $runtime.Block.Release
                              local.set $$t17.0

                              ;;*t17 = t16
                              local.get $$t17.1
                              local.get $$t16
                              i32.store offset=0 align=4

                              ;;if t6 goto 4 else 5
                              local.get $$t6
                              if
                                br $$Block_3
                              else
                                br $$Block_4
                              end

                            end ;;$Block_1
                            i32.const 2
                            local.set $$current_block

                            ;;len(x)
                            local.get $x.2
                            local.set $$t18

                            ;;len(y)
                            local.get $y.2
                            local.set $$t19

                            ;;t18 < t19
                            local.get $$t18
                            local.get $$t19
                            i32.lt_s
                            local.set $$t20

                            ;;if t20 goto 9 else 10
                            local.get $$t20
                            if
                              br $$Block_8
                            else
                              br $$Block_9
                            end

                          end ;;$Block_2
                          i32.const 3
                          local.set $$current_block

                          ;;t8 < t15
                          local.get $$t8
                          local.get $$t15
                          i32.lt_s
                          local.set $$t21

                          ;;if t21 goto 6 else 7
                          local.get $$t21
                          if
                            br $$Block_5
                          else
                            br $$Block_6
                          end

                        end ;;$Block_3
                        i32.const 4
                        local.set $$current_block

                        ;;jump 5
                        br $$Block_4

                      end ;;$Block_4
                      local.get $$current_block
                      i32.const 1
                      i32.eq
                      if (result i32)
                        i32.const 0
                      else
                        local.get $$t13
                      end
                      local.set $$t22
                      i32.const 5
                      local.set $$current_block

                      ;;t22 != true:bool
                      local.get $$t22
                      i32.const 1
                      i32.eq
                      i32.eqz
                      local.set $$t23

                      ;;if t23 goto 2 else 3
                      local.get $$t23
                      if
                        i32.const 2
                        local.set $$block_selector
                        br $$BlockDisp
                      else
                        i32.const 3
                        local.set $$block_selector
                        br $$BlockDisp
                      end

                    end ;;$Block_5
                    i32.const 6
                    local.set $$current_block

                    ;;return -1:i32
                    i32.const -1
                    local.set $$ret_0
                    br $$BlockFnBody

                  end ;;$Block_6
                  i32.const 7
                  local.set $$current_block

                  ;;t8 > t15
                  local.get $$t8
                  local.get $$t15
                  i32.gt_s
                  local.set $$t24

                  ;;if t24 goto 8 else 1
                  local.get $$t24
                  if
                    br $$Block_7
                  else
                    i32.const 1
                    local.set $$block_selector
                    br $$BlockDisp
                  end

                end ;;$Block_7
                i32.const 8
                local.set $$current_block

                ;;return 1:i32
                i32.const 1
                local.set $$ret_0
                br $$BlockFnBody

              end ;;$Block_8
              i32.const 9
              local.set $$current_block

              ;;return -1:i32
              i32.const -1
              local.set $$ret_0
              br $$BlockFnBody

            end ;;$Block_9
            i32.const 10
            local.set $$current_block

            ;;t18 > t19
            local.get $$t18
            local.get $$t19
            i32.gt_s
            local.set $$t25

            ;;if t25 goto 11 else 12
            local.get $$t25
            if
              br $$Block_10
            else
              br $$Block_11
            end

          end ;;$Block_10
          i32.const 11
          local.set $$current_block

          ;;return 1:i32
          i32.const 1
          local.set $$ret_0
          br $$BlockFnBody

        end ;;$Block_11
        i32.const 12
        local.set $$current_block

        ;;return 0:i32
        i32.const 0
        local.set $$ret_0
        br $$BlockFnBody

      end ;;$Block_12
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$ret_0
  local.get $$t0.0
  call $runtime.Block.Release
  local.get $$t2.0
  call $runtime.Block.Release
  local.get $$t10.0
  call $runtime.Block.Release
  local.get $$t17.0
  call $runtime.Block.Release
) ;;$wa.runtime.string_Comp

(func $$wa.runtime.string_GEQ (param $x.0 i32) (param $x.1 i32) (param $x.2 i32) (param $y.0 i32) (param $y.1 i32) (param $y.2 i32) (result i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$ret_0 i32)
  (local $$t0 i32)
  (local $$t1 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;string_Comp(x, y)
        local.get $x.0
        local.get $x.1
        local.get $x.2
        local.get $y.0
        local.get $y.1
        local.get $y.2
        call $$wa.runtime.string_Comp
        local.set $$t0

        ;;t0 != -1:i32
        local.get $$t0
        i32.const -1
        i32.eq
        i32.eqz
        local.set $$t1

        ;;return t1
        local.get $$t1
        local.set $$ret_0
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$ret_0
) ;;$wa.runtime.string_GEQ

(func $$wa.runtime.string_GTR (param $x.0 i32) (param $x.1 i32) (param $x.2 i32) (param $y.0 i32) (param $y.1 i32) (param $y.2 i32) (result i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$ret_0 i32)
  (local $$t0 i32)
  (local $$t1 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;string_Comp(x, y)
        local.get $x.0
        local.get $x.1
        local.get $x.2
        local.get $y.0
        local.get $y.1
        local.get $y.2
        call $$wa.runtime.string_Comp
        local.set $$t0

        ;;t0 == 1:i32
        local.get $$t0
        i32.const 1
        i32.eq
        local.set $$t1

        ;;return t1
        local.get $$t1
        local.set $$ret_0
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$ret_0
) ;;$wa.runtime.string_GTR

(func $$wa.runtime.string_LEQ (param $x.0 i32) (param $x.1 i32) (param $x.2 i32) (param $y.0 i32) (param $y.1 i32) (param $y.2 i32) (result i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$ret_0 i32)
  (local $$t0 i32)
  (local $$t1 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;string_Comp(x, y)
        local.get $x.0
        local.get $x.1
        local.get $x.2
        local.get $y.0
        local.get $y.1
        local.get $y.2
        call $$wa.runtime.string_Comp
        local.set $$t0

        ;;t0 != 1:i32
        local.get $$t0
        i32.const 1
        i32.eq
        i32.eqz
        local.set $$t1

        ;;return t1
        local.get $$t1
        local.set $$ret_0
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$ret_0
) ;;$wa.runtime.string_LEQ

(func $$wa.runtime.string_LSS (param $x.0 i32) (param $x.1 i32) (param $x.2 i32) (param $y.0 i32) (param $y.1 i32) (param $y.2 i32) (result i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$ret_0 i32)
  (local $$t0 i32)
  (local $$t1 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;string_Comp(x, y)
        local.get $x.0
        local.get $x.1
        local.get $x.2
        local.get $y.0
        local.get $y.1
        local.get $y.2
        call $$wa.runtime.string_Comp
        local.set $$t0

        ;;t0 == -1:i32
        local.get $$t0
        i32.const -1
        i32.eq
        local.set $$t1

        ;;return t1
        local.get $$t1
        local.set $$ret_0
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$ret_0
) ;;$wa.runtime.string_LSS

(func $$runtime.waPrintBool (param $i i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;syscall/js.PrintBool(i)
        local.get $i
        call $syscall$js.PrintBool

        ;;return
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
) ;;$runtime.waPrintBool

(func $$runtime.waPrintChar (param $ch i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;waPrintRune(ch)
        local.get $ch
        call $$runtime.waPrintRune

        ;;return
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
) ;;$runtime.waPrintChar

(func $$runtime.waPrintF32 (param $i f32)
  (local $$block_selector i32)
  (local $$current_block i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;syscall/js.PrintF32(i)
        local.get $i
        call $syscall$js.PrintF32

        ;;return
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
) ;;$runtime.waPrintF32

(func $$runtime.waPrintF64 (param $i f64)
  (local $$block_selector i32)
  (local $$current_block i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;syscall/js.PrintF64(i)
        local.get $i
        call $syscall$js.PrintF64

        ;;return
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
) ;;$runtime.waPrintF64

(func $$runtime.waPrintI32 (param $i i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;syscall/js.PrintI32(i)
        local.get $i
        call $syscall$js.PrintI32

        ;;return
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
) ;;$runtime.waPrintI32

(func $$runtime.waPrintI64 (param $i i64)
  (local $$block_selector i32)
  (local $$current_block i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;syscall/js.PrintI64(i)
        local.get $i
        call $syscall$js.PrintI64

        ;;return
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
) ;;$runtime.waPrintI64

(func $$runtime.waPrintRune (param $ch i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;syscall/js.PrintRune(ch)
        local.get $ch
        call $syscall$js.PrintRune

        ;;return
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
) ;;$runtime.waPrintRune

(func $$runtime.waPrintString (param $s.0 i32) (param $s.1 i32) (param $s.2 i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;printString(s)
        local.get $s.0
        local.get $s.1
        local.get $s.2
        call $runtime.printString

        ;;return
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
) ;;$runtime.waPrintString

(func $$runtime.waPrintU32 (param $i i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;syscall/js.PrintU32(i)
        local.get $i
        call $syscall$js.PrintU32

        ;;return
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
) ;;$runtime.waPrintU32

(func $$runtime.waPrintU32Ptr (param $i i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;syscall/js.PrintU32Ptr(i)
        local.get $i
        call $syscall$js.PrintU32Ptr

        ;;return
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
) ;;$runtime.waPrintU32Ptr

(func $$runtime.waPrintU64 (param $i i64)
  (local $$block_selector i32)
  (local $$current_block i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;syscall/js.PrintU64(i)
        local.get $i
        call $syscall$js.PrintU64

        ;;return
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
) ;;$runtime.waPrintU64

(func $$runtime.waPuts (param $ptr i32) (param $len i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;syscall/js.Puts(ptr, len)
        local.get $ptr
        local.get $len
        call $syscall$js.Puts

        ;;return
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
) ;;$runtime.waPuts

(func $$brainfuck$bfpkg.BrainFuck.$$OnFree (param $$ptr i32)
  local.get $$ptr
  i32.const 30000
  i32.add
  i32.const 3
  call_indirect (type $$OnFree)
) ;;$brainfuck$bfpkg.BrainFuck.$$OnFree

(func $brainfuck$bfpkg.NewBrainFuck (param $code.0 i32) (param $code.1 i32) (param $code.2 i32) (result i32 i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$ret_0.0 i32)
  (local $$ret_0.1 i32)
  (local $$t0.0 i32)
  (local $$t0.1 i32)
  (local $$t1.0 i32)
  (local $$t1.1 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;new BrainFuck (complit)
        i32.const 30036
        call $runtime.HeapAlloc
        i32.const 1
        i32.const 23
        i32.const 30020
        call $runtime.Block.Init
        call $runtime.DupI32
        i32.const 16
        i32.add
        local.set $$t0.1
        local.get $$t0.0
        call $runtime.Block.Release
        local.set $$t0.0

        ;;&t0.code [#1]
        local.get $$t0.0
        call $runtime.Block.Retain
        local.get $$t0.1
        i32.const 30000
        i32.add
        local.set $$t1.1
        local.get $$t1.0
        call $runtime.Block.Release
        local.set $$t1.0

        ;;*t1 = code
        local.get $$t1.1
        local.get $code.0
        call $runtime.Block.Retain
        local.get $$t1.1
        i32.load offset=0 align=1
        call $runtime.Block.Release
        i32.store offset=0 align=1
        local.get $$t1.1
        local.get $code.1
        i32.store offset=4 align=4
        local.get $$t1.1
        local.get $code.2
        i32.store offset=8 align=4

        ;;return t0
        local.get $$t0.0
        call $runtime.Block.Retain
        local.get $$t0.1
        local.set $$ret_0.1
        local.get $$ret_0.0
        call $runtime.Block.Release
        local.set $$ret_0.0
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$ret_0.0
  call $runtime.Block.Retain
  local.get $$ret_0.1
  local.get $$ret_0.0
  call $runtime.Block.Release
  local.get $$t0.0
  call $runtime.Block.Release
  local.get $$t1.0
  call $runtime.Block.Release
) ;;brainfuck$bfpkg.NewBrainFuck

(func $brainfuck$bfpkg.init
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$t0 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_2
        block $$Block_1
          block $$Block_0
            block $$BlockSel
              local.get $$block_selector
              br_table 0 1 2 0
            end ;;$BlockSel
            i32.const 0
            local.set $$current_block

            ;;*init$guard
            global.get $brainfuck$bfpkg.init$guard
            local.set $$t0

            ;;if t0 goto 2 else 1
            local.get $$t0
            if
              br $$Block_1
            else
              br $$Block_0
            end

          end ;;$Block_0
          i32.const 1
          local.set $$current_block

          ;;*init$guard = true:bool
          i32.const 1
          global.set $brainfuck$bfpkg.init$guard

          ;;jump 2
          br $$Block_1

        end ;;$Block_1
        i32.const 2
        local.set $$current_block

        ;;return
        br $$BlockFnBody

      end ;;$Block_2
    end ;;$BlockDisp
  end ;;$BlockFnBody
) ;;brainfuck$bfpkg.init

(func $brainfuck.Run (export "brainfuck.Run") (param $code.0 i32) (param $code.1 i32) (param $code.2 i32) (result i32 i32 i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$ret_0.0 i32)
  (local $$ret_0.1 i32)
  (local $$ret_0.2 i32)
  (local $$t0.0 i32)
  (local $$t0.1 i32)
  (local $$t1.0 i32)
  (local $$t1.1 i32)
  (local $$t1.2 i32)
  (local $$t1.3 i32)
  (local $$t2.0 i32)
  (local $$t2.1 i32)
  (local $$t2.2 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;brainfuck/bfpkg.NewBrainFuck(code)
        local.get $code.0
        local.get $code.1
        local.get $code.2
        call $brainfuck$bfpkg.NewBrainFuck
        local.set $$t0.1
        local.get $$t0.0
        call $runtime.Block.Release
        local.set $$t0.0

        ;;(*brainfuck/bfpkg.BrainFuck).Run(t0)
        local.get $$t0.0
        local.get $$t0.1
        call $brainfuck$bfpkg.BrainFuck.Run
        local.set $$t1.3
        local.set $$t1.2
        local.set $$t1.1
        local.get $$t1.0
        call $runtime.Block.Release
        local.set $$t1.0

        ;;convert string <- []byte (t1)
        i32.const 0
        i32.const 8224
        i32.const 0
        local.get $$t1.0
        local.get $$t1.1
        local.get $$t1.2
        call $$string.appendstr
        local.set $$t2.2
        local.set $$t2.1
        local.get $$t2.0
        call $runtime.Block.Release
        local.set $$t2.0

        ;;return t2
        local.get $$t2.0
        call $runtime.Block.Retain
        local.get $$t2.1
        local.get $$t2.2
        local.set $$ret_0.2
        local.set $$ret_0.1
        local.get $$ret_0.0
        call $runtime.Block.Release
        local.set $$ret_0.0
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$ret_0.0
  call $runtime.Block.Retain
  local.get $$ret_0.1
  local.get $$ret_0.2
  local.get $$ret_0.0
  call $runtime.Block.Release
  local.get $$t0.0
  call $runtime.Block.Release
  local.get $$t1.0
  call $runtime.Block.Release
  local.get $$t2.0
  call $runtime.Block.Release
) ;;brainfuck.Run

(func $brainfuck.init
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$t0 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_2
        block $$Block_1
          block $$Block_0
            block $$BlockSel
              local.get $$block_selector
              br_table 0 1 2 0
            end ;;$BlockSel
            i32.const 0
            local.set $$current_block

            ;;*init$guard
            global.get $brainfuck.init$guard
            local.set $$t0

            ;;if t0 goto 2 else 1
            local.get $$t0
            if
              br $$Block_1
            else
              br $$Block_0
            end

          end ;;$Block_0
          i32.const 1
          local.set $$current_block

          ;;*init$guard = true:bool
          i32.const 1
          global.set $brainfuck.init$guard

          ;;runtime.init()
          call $runtime.init

          ;;brainfuck/bfpkg.init()
          call $brainfuck$bfpkg.init

          ;;jump 2
          br $$Block_1

        end ;;$Block_1
        i32.const 2
        local.set $$current_block

        ;;return
        br $$BlockFnBody

      end ;;$Block_2
    end ;;$BlockDisp
  end ;;$BlockFnBody
) ;;brainfuck.init

(func $brainfuck.main (export "brainfuck.main")
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$t0.0 i32)
  (local $$t0.1 i32)
  (local $$t1.0 i32)
  (local $$t1.1 i32)
  (local $$t1.2 i32)
  (local $$t1.3 i32)
  (local $$t2.0 i32)
  (local $$t2.1 i32)
  (local $$t2.2 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;brainfuck/bfpkg.NewBrainFuck("++++++++++[>+++++...":string)
        i32.const 0
        i32.const 8348
        i32.const 33
        call $brainfuck$bfpkg.NewBrainFuck
        local.set $$t0.1
        local.get $$t0.0
        call $runtime.Block.Release
        local.set $$t0.0

        ;;(*brainfuck/bfpkg.BrainFuck).Run(t0)
        local.get $$t0.0
        local.get $$t0.1
        call $brainfuck$bfpkg.BrainFuck.Run
        local.set $$t1.3
        local.set $$t1.2
        local.set $$t1.1
        local.get $$t1.0
        call $runtime.Block.Release
        local.set $$t1.0

        ;;convert string <- []byte (t1)
        i32.const 0
        i32.const 8224
        i32.const 0
        local.get $$t1.0
        local.get $$t1.1
        local.get $$t1.2
        call $$string.appendstr
        local.set $$t2.2
        local.set $$t2.1
        local.get $$t2.0
        call $runtime.Block.Release
        local.set $$t2.0

        ;;println(t2)
        local.get $$t2.1
        local.get $$t2.2
        call $$runtime.waPuts
        i32.const 10
        call $$runtime.waPrintChar

        ;;return
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$t0.0
  call $runtime.Block.Release
  local.get $$t1.0
  call $runtime.Block.Release
  local.get $$t2.0
  call $runtime.Block.Release
) ;;brainfuck.main

(func $syscall$js.PrintBool (param $v i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;__import__print_bool(v)
        local.get $v
        call $syscall$js.__import__print_bool

        ;;return
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
) ;;syscall$js.PrintBool

(func $syscall$js.PrintF32 (param $v f32)
  (local $$block_selector i32)
  (local $$current_block i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;__import__print_f32(v)
        local.get $v
        call $syscall$js.__import__print_f32

        ;;return
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
) ;;syscall$js.PrintF32

(func $syscall$js.PrintF64 (param $v f64)
  (local $$block_selector i32)
  (local $$current_block i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;__import__print_f64(v)
        local.get $v
        call $syscall$js.__import__print_f64

        ;;return
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
) ;;syscall$js.PrintF64

(func $syscall$js.PrintI32 (param $v i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;__import__print_i32(v)
        local.get $v
        call $syscall$js.__import__print_i32

        ;;return
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
) ;;syscall$js.PrintI32

(func $syscall$js.PrintI64 (param $v i64)
  (local $$block_selector i32)
  (local $$current_block i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;__import__print_i64(v)
        local.get $v
        call $syscall$js.__import__print_i64

        ;;return
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
) ;;syscall$js.PrintI64

(func $syscall$js.PrintRune (param $v i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;__import__print_rune(v)
        local.get $v
        call $syscall$js.__import__print_rune

        ;;return
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
) ;;syscall$js.PrintRune

(func $syscall$js.PrintString (param $s.0 i32) (param $s.1 i32) (param $s.2 i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$t0 i32)
  (local $$t1 i32)
  (local $$t2 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;__linkname__string_to_ptr(s)
        local.get $s.0
        local.get $s.1
        local.get $s.2
        call $$syscall/js.__linkname__string_to_ptr
        local.set $$t0

        ;;len(s)
        local.get $s.2
        local.set $$t1

        ;;convert i32 <- int (t1)
        local.get $$t1
        local.set $$t2

        ;;__import__print_str(t0, t2)
        local.get $$t0
        local.get $$t2
        call $syscall$js.__import__print_str

        ;;return
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
) ;;syscall$js.PrintString

(func $syscall$js.PrintU32 (param $v i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;__import__print_u32(v)
        local.get $v
        call $syscall$js.__import__print_u32

        ;;return
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
) ;;syscall$js.PrintU32

(func $syscall$js.PrintU32Ptr (param $v i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;__import__print_ptr(v)
        local.get $v
        call $syscall$js.__import__print_ptr

        ;;return
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
) ;;syscall$js.PrintU32Ptr

(func $syscall$js.PrintU64 (param $v i64)
  (local $$block_selector i32)
  (local $$current_block i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;__import__print_u64(v)
        local.get $v
        call $syscall$js.__import__print_u64

        ;;return
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
) ;;syscall$js.PrintU64

(func $syscall$js.ProcExit (param $v i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;__import__proc_exit(v)
        local.get $v
        call $syscall$js.__import__proc_exit

        ;;return
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
) ;;syscall$js.ProcExit

(func $syscall$js.Puts (param $ptr i32) (param $len i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;__import__print_str(ptr, len)
        local.get $ptr
        local.get $len
        call $syscall$js.__import__print_str

        ;;return
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
) ;;syscall$js.Puts

(func $syscall$js.init
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$t0 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_2
        block $$Block_1
          block $$Block_0
            block $$BlockSel
              local.get $$block_selector
              br_table 0 1 2 0
            end ;;$BlockSel
            i32.const 0
            local.set $$current_block

            ;;*init$guard
            global.get $syscall$js.init$guard
            local.set $$t0

            ;;if t0 goto 2 else 1
            local.get $$t0
            if
              br $$Block_1
            else
              br $$Block_0
            end

          end ;;$Block_0
          i32.const 1
          local.set $$current_block

          ;;*init$guard = true:bool
          i32.const 1
          global.set $syscall$js.init$guard

          ;;jump 2
          br $$Block_1

        end ;;$Block_1
        i32.const 2
        local.set $$current_block

        ;;return
        br $$BlockFnBody

      end ;;$Block_2
    end ;;$BlockDisp
  end ;;$BlockFnBody
) ;;syscall$js.init

(func $runtime.mapNode.Parent (param $this.0 i32) (param $this.1 i32) (param $m.0 i32) (param $m.1 i32) (result i32 i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$ret_0.0 i32)
  (local $$ret_0.1 i32)
  (local $$t0.0 i32)
  (local $$t0.1 i32)
  (local $$t1.0 i32)
  (local $$t1.1 i32)
  (local $$t1.2 i32)
  (local $$t1.3 i32)
  (local $$t2.0 i32)
  (local $$t2.1 i32)
  (local $$t3 i32)
  (local $$t4 i32)
  (local $$t5.0 i32)
  (local $$t5.1 i32)
  (local $$t6.0 i32)
  (local $$t6.1 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;&m.nodes [#2]
        local.get $m.0
        call $runtime.Block.Retain
        local.get $m.1
        i32.const 16
        i32.add
        local.set $$t0.1
        local.get $$t0.0
        call $runtime.Block.Release
        local.set $$t0.0

        ;;*t0
        local.get $$t0.1
        i32.load offset=0 align=4
        call $runtime.Block.Retain
        local.get $$t0.1
        i32.load offset=4 align=4
        local.get $$t0.1
        i32.load offset=8 align=4
        local.get $$t0.1
        i32.load offset=12 align=4
        local.set $$t1.3
        local.set $$t1.2
        local.set $$t1.1
        local.get $$t1.0
        call $runtime.Block.Release
        local.set $$t1.0

        ;;&this.parentIdx [#0]
        local.get $this.0
        call $runtime.Block.Retain
        local.get $this.1
        i32.const 0
        i32.add
        local.set $$t2.1
        local.get $$t2.0
        call $runtime.Block.Release
        local.set $$t2.0

        ;;*t2
        local.get $$t2.1
        i32.load offset=0 align=4
        local.set $$t3

        ;;changetype int <- mapNodeIdx (t3)
        local.get $$t3
        local.set $$t4

        ;;&t1[t4]
        local.get $$t1.0
        call $runtime.Block.Retain
        local.get $$t1.1
        i32.const 8
        local.get $$t4
        i32.mul
        i32.add
        local.set $$t5.1
        local.get $$t5.0
        call $runtime.Block.Release
        local.set $$t5.0

        ;;*t5
        local.get $$t5.1
        i32.load offset=0 align=4
        call $runtime.Block.Retain
        local.get $$t5.1
        i32.load offset=4 align=4
        local.set $$t6.1
        local.get $$t6.0
        call $runtime.Block.Release
        local.set $$t6.0

        ;;return t6
        local.get $$t6.0
        call $runtime.Block.Retain
        local.get $$t6.1
        local.set $$ret_0.1
        local.get $$ret_0.0
        call $runtime.Block.Release
        local.set $$ret_0.0
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$ret_0.0
  call $runtime.Block.Retain
  local.get $$ret_0.1
  local.get $$ret_0.0
  call $runtime.Block.Release
  local.get $$t0.0
  call $runtime.Block.Release
  local.get $$t1.0
  call $runtime.Block.Release
  local.get $$t2.0
  call $runtime.Block.Release
  local.get $$t5.0
  call $runtime.Block.Release
  local.get $$t6.0
  call $runtime.Block.Release
) ;;runtime.mapNode.Parent

(func $runtime.mapImp.Delete (param $this.0 i32) (param $this.1 i32) (param $k.0.0 i32) (param $k.0.1 i32) (param $k.1 i32) (param $k.2 i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$t0.0 i32)
  (local $$t0.1 i32)
  (local $$t1.0 i32)
  (local $$t1.1 i32)
  (local $$t2.0 i32)
  (local $$t2.1 i32)
  (local $$t3 i32)
  (local $$t4 i32)
  (local $$t5.0 i32)
  (local $$t5.1 i32)
  (local $$t6 i32)
  (local $$t7.0 i32)
  (local $$t7.1 i32)
  (local $$t8.0 i32)
  (local $$t8.1 i32)
  (local $$t8.2 i32)
  (local $$t8.3 i32)
  (local $$t9 i32)
  (local $$t10 i32)
  (local $$t11 i32)
  (local $$t12 i32)
  (local $$t13.0 i32)
  (local $$t13.1 i32)
  (local $$t14.0 i32)
  (local $$t14.1 i32)
  (local $$t14.2 i32)
  (local $$t14.3 i32)
  (local $$t15.0 i32)
  (local $$t15.1 i32)
  (local $$t16.0 i32)
  (local $$t16.1 i32)
  (local $$t16.2 i32)
  (local $$t16.3 i32)
  (local $$t17 i32)
  (local $$t18 i32)
  (local $$t19.0 i32)
  (local $$t19.1 i32)
  (local $$t20.0 i32)
  (local $$t20.1 i32)
  (local $$t21.0 i32)
  (local $$t21.1 i32)
  (local $$t22.0 i32)
  (local $$t22.1 i32)
  (local $$t23 i32)
  (local $$t24.0 i32)
  (local $$t24.1 i32)
  (local $$t25.0 i32)
  (local $$t25.1 i32)
  (local $$t25.2 i32)
  (local $$t25.3 i32)
  (local $$t26.0 i32)
  (local $$t26.1 i32)
  (local $$t27 i32)
  (local $$t28 i32)
  (local $$t29.0 i32)
  (local $$t29.1 i32)
  (local $$t30.0 i32)
  (local $$t30.1 i32)
  (local $$t31.0 i32)
  (local $$t31.1 i32)
  (local $$t32.0 i32)
  (local $$t32.1 i32)
  (local $$t33.0 i32)
  (local $$t33.1 i32)
  (local $$t34 i32)
  (local $$t35.0 i32)
  (local $$t35.1 i32)
  (local $$t36.0 i32)
  (local $$t36.1 i32)
  (local $$t37.0 i32)
  (local $$t37.1 i32)
  (local $$t37.2 i32)
  (local $$t37.3 i32)
  (local $$t38.0 i32)
  (local $$t38.1 i32)
  (local $$t39.0 i32)
  (local $$t39.1 i32)
  (local $$t39.2 i32)
  (local $$t39.3 i32)
  (local $$t40 i32)
  (local $$t41 i32)
  (local $$t42.0 i32)
  (local $$t42.1 i32)
  (local $$t42.2 i32)
  (local $$t42.3 i32)
  (local $$t43.0 i32)
  (local $$t43.1 i32)
  (local $$t44.0 i32)
  (local $$t44.1 i32)
  (local $$t45.0 i32)
  (local $$t45.1 i32)
  (local $$t46.0 i32)
  (local $$t46.1 i32)
  (local $$t47.0 i32)
  (local $$t47.1 i32)
  (local $$t48.0 i32)
  (local $$t48.1 i32)
  (local $$t49 i32)
  (local $$t50.0 i32)
  (local $$t50.1 i32)
  (local $$t51.0 i32)
  (local $$t51.1 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_7
        block $$Block_6
          block $$Block_5
            block $$Block_4
              block $$Block_3
                block $$Block_2
                  block $$Block_1
                    block $$Block_0
                      block $$BlockSel
                        local.get $$block_selector
                        br_table 0 1 2 3 4 5 6 7 0
                      end ;;$BlockSel
                      i32.const 0
                      local.set $$current_block

                      ;;(*mapImp).search(this, k)
                      local.get $this.0
                      local.get $this.1
                      local.get $k.0.0
                      local.get $k.0.1
                      local.get $k.1
                      local.get $k.2
                      call $runtime.mapImp.search
                      local.set $$t0.1
                      local.get $$t0.0
                      call $runtime.Block.Release
                      local.set $$t0.0

                      ;;&this.NIL [#0]
                      local.get $this.0
                      call $runtime.Block.Retain
                      local.get $this.1
                      i32.const 0
                      i32.add
                      local.set $$t1.1
                      local.get $$t1.0
                      call $runtime.Block.Release
                      local.set $$t1.0

                      ;;*t1
                      local.get $$t1.1
                      i32.load offset=0 align=4
                      call $runtime.Block.Retain
                      local.get $$t1.1
                      i32.load offset=4 align=4
                      local.set $$t2.1
                      local.get $$t2.0
                      call $runtime.Block.Release
                      local.set $$t2.0

                      ;;t0 == t2
                      local.get $$t0.1
                      local.get $$t2.1
                      i32.eq
                      local.set $$t3

                      ;;if t3 goto 1 else 2
                      local.get $$t3
                      if
                        br $$Block_0
                      else
                        br $$Block_1
                      end

                    end ;;$Block_0
                    i32.const 1
                    local.set $$current_block

                    ;;return
                    br $$BlockFnBody

                  end ;;$Block_1
                  i32.const 2
                  local.set $$current_block

                  ;;(*mapImp).delete(this, t0)
                  local.get $this.0
                  local.get $this.1
                  local.get $$t0.0
                  local.get $$t0.1
                  call $runtime.mapImp.delete
                  local.set $$t4

                  ;;&t0.NodeIdx [#1]
                  local.get $$t0.0
                  call $runtime.Block.Retain
                  local.get $$t0.1
                  i32.const 4
                  i32.add
                  local.set $$t5.1
                  local.get $$t5.0
                  call $runtime.Block.Release
                  local.set $$t5.0

                  ;;*t5
                  local.get $$t5.1
                  i32.load offset=0 align=4
                  local.set $$t6

                  ;;&this.nodes [#2]
                  local.get $this.0
                  call $runtime.Block.Retain
                  local.get $this.1
                  i32.const 16
                  i32.add
                  local.set $$t7.1
                  local.get $$t7.0
                  call $runtime.Block.Release
                  local.set $$t7.0

                  ;;*t7
                  local.get $$t7.1
                  i32.load offset=0 align=4
                  call $runtime.Block.Retain
                  local.get $$t7.1
                  i32.load offset=4 align=4
                  local.get $$t7.1
                  i32.load offset=8 align=4
                  local.get $$t7.1
                  i32.load offset=12 align=4
                  local.set $$t8.3
                  local.set $$t8.2
                  local.set $$t8.1
                  local.get $$t8.0
                  call $runtime.Block.Release
                  local.set $$t8.0

                  ;;len(t8)
                  local.get $$t8.2
                  local.set $$t9

                  ;;t9 - 1:int
                  local.get $$t9
                  i32.const 1
                  i32.sub
                  local.set $$t10

                  ;;changetype mapNodeIdx <- int (t10)
                  local.get $$t10
                  local.set $$t11

                  ;;t6 < t11
                  local.get $$t6
                  local.get $$t11
                  i32.lt_s
                  local.set $$t12

                  ;;if t12 goto 3 else 4
                  local.get $$t12
                  if
                    br $$Block_2
                  else
                    br $$Block_3
                  end

                end ;;$Block_2
                i32.const 3
                local.set $$current_block

                ;;&this.nodes [#2]
                local.get $this.0
                call $runtime.Block.Retain
                local.get $this.1
                i32.const 16
                i32.add
                local.set $$t13.1
                local.get $$t13.0
                call $runtime.Block.Release
                local.set $$t13.0

                ;;*t13
                local.get $$t13.1
                i32.load offset=0 align=4
                call $runtime.Block.Retain
                local.get $$t13.1
                i32.load offset=4 align=4
                local.get $$t13.1
                i32.load offset=8 align=4
                local.get $$t13.1
                i32.load offset=12 align=4
                local.set $$t14.3
                local.set $$t14.2
                local.set $$t14.1
                local.get $$t14.0
                call $runtime.Block.Release
                local.set $$t14.0

                ;;&this.nodes [#2]
                local.get $this.0
                call $runtime.Block.Retain
                local.get $this.1
                i32.const 16
                i32.add
                local.set $$t15.1
                local.get $$t15.0
                call $runtime.Block.Release
                local.set $$t15.0

                ;;*t15
                local.get $$t15.1
                i32.load offset=0 align=4
                call $runtime.Block.Retain
                local.get $$t15.1
                i32.load offset=4 align=4
                local.get $$t15.1
                i32.load offset=8 align=4
                local.get $$t15.1
                i32.load offset=12 align=4
                local.set $$t16.3
                local.set $$t16.2
                local.set $$t16.1
                local.get $$t16.0
                call $runtime.Block.Release
                local.set $$t16.0

                ;;len(t16)
                local.get $$t16.2
                local.set $$t17

                ;;t17 - 1:int
                local.get $$t17
                i32.const 1
                i32.sub
                local.set $$t18

                ;;&t14[t18]
                local.get $$t14.0
                call $runtime.Block.Retain
                local.get $$t14.1
                i32.const 8
                local.get $$t18
                i32.mul
                i32.add
                local.set $$t19.1
                local.get $$t19.0
                call $runtime.Block.Release
                local.set $$t19.0

                ;;*t19
                local.get $$t19.1
                i32.load offset=0 align=4
                call $runtime.Block.Retain
                local.get $$t19.1
                i32.load offset=4 align=4
                local.set $$t20.1
                local.get $$t20.0
                call $runtime.Block.Release
                local.set $$t20.0

                ;;&t20.NodeIdx [#1]
                local.get $$t20.0
                call $runtime.Block.Retain
                local.get $$t20.1
                i32.const 4
                i32.add
                local.set $$t21.1
                local.get $$t21.0
                call $runtime.Block.Release
                local.set $$t21.0

                ;;&t0.NodeIdx [#1]
                local.get $$t0.0
                call $runtime.Block.Retain
                local.get $$t0.1
                i32.const 4
                i32.add
                local.set $$t22.1
                local.get $$t22.0
                call $runtime.Block.Release
                local.set $$t22.0

                ;;*t22
                local.get $$t22.1
                i32.load offset=0 align=4
                local.set $$t23

                ;;*t21 = t23
                local.get $$t21.1
                local.get $$t23
                i32.store offset=0 align=4

                ;;&this.nodes [#2]
                local.get $this.0
                call $runtime.Block.Retain
                local.get $this.1
                i32.const 16
                i32.add
                local.set $$t24.1
                local.get $$t24.0
                call $runtime.Block.Release
                local.set $$t24.0

                ;;*t24
                local.get $$t24.1
                i32.load offset=0 align=4
                call $runtime.Block.Retain
                local.get $$t24.1
                i32.load offset=4 align=4
                local.get $$t24.1
                i32.load offset=8 align=4
                local.get $$t24.1
                i32.load offset=12 align=4
                local.set $$t25.3
                local.set $$t25.2
                local.set $$t25.1
                local.get $$t25.0
                call $runtime.Block.Release
                local.set $$t25.0

                ;;&t0.NodeIdx [#1]
                local.get $$t0.0
                call $runtime.Block.Retain
                local.get $$t0.1
                i32.const 4
                i32.add
                local.set $$t26.1
                local.get $$t26.0
                call $runtime.Block.Release
                local.set $$t26.0

                ;;*t26
                local.get $$t26.1
                i32.load offset=0 align=4
                local.set $$t27

                ;;changetype int <- mapNodeIdx (t27)
                local.get $$t27
                local.set $$t28

                ;;&t25[t28]
                local.get $$t25.0
                call $runtime.Block.Retain
                local.get $$t25.1
                i32.const 8
                local.get $$t28
                i32.mul
                i32.add
                local.set $$t29.1
                local.get $$t29.0
                call $runtime.Block.Release
                local.set $$t29.0

                ;;*t29 = t20
                local.get $$t29.1
                local.get $$t20.0
                call $runtime.Block.Retain
                local.get $$t29.1
                i32.load offset=0 align=1
                call $runtime.Block.Release
                i32.store offset=0 align=1
                local.get $$t29.1
                local.get $$t20.1
                i32.store offset=4 align=4

                ;;&t20.Left [#2]
                local.get $$t20.0
                call $runtime.Block.Retain
                local.get $$t20.1
                i32.const 8
                i32.add
                local.set $$t30.1
                local.get $$t30.0
                call $runtime.Block.Release
                local.set $$t30.0

                ;;*t30
                local.get $$t30.1
                i32.load offset=0 align=4
                call $runtime.Block.Retain
                local.get $$t30.1
                i32.load offset=4 align=4
                local.set $$t31.1
                local.get $$t31.0
                call $runtime.Block.Release
                local.set $$t31.0

                ;;&this.NIL [#0]
                local.get $this.0
                call $runtime.Block.Retain
                local.get $this.1
                i32.const 0
                i32.add
                local.set $$t32.1
                local.get $$t32.0
                call $runtime.Block.Release
                local.set $$t32.0

                ;;*t32
                local.get $$t32.1
                i32.load offset=0 align=4
                call $runtime.Block.Retain
                local.get $$t32.1
                i32.load offset=4 align=4
                local.set $$t33.1
                local.get $$t33.0
                call $runtime.Block.Release
                local.set $$t33.0

                ;;t31 != t33
                local.get $$t31.1
                local.get $$t33.1
                i32.eq
                i32.eqz
                local.set $$t34

                ;;if t34 goto 5 else 6
                local.get $$t34
                if
                  br $$Block_4
                else
                  br $$Block_5
                end

              end ;;$Block_3
              i32.const 4
              local.set $$current_block

              ;;&this.nodes [#2]
              local.get $this.0
              call $runtime.Block.Retain
              local.get $this.1
              i32.const 16
              i32.add
              local.set $$t35.1
              local.get $$t35.0
              call $runtime.Block.Release
              local.set $$t35.0

              ;;&this.nodes [#2]
              local.get $this.0
              call $runtime.Block.Retain
              local.get $this.1
              i32.const 16
              i32.add
              local.set $$t36.1
              local.get $$t36.0
              call $runtime.Block.Release
              local.set $$t36.0

              ;;*t36
              local.get $$t36.1
              i32.load offset=0 align=4
              call $runtime.Block.Retain
              local.get $$t36.1
              i32.load offset=4 align=4
              local.get $$t36.1
              i32.load offset=8 align=4
              local.get $$t36.1
              i32.load offset=12 align=4
              local.set $$t37.3
              local.set $$t37.2
              local.set $$t37.1
              local.get $$t37.0
              call $runtime.Block.Release
              local.set $$t37.0

              ;;&this.nodes [#2]
              local.get $this.0
              call $runtime.Block.Retain
              local.get $this.1
              i32.const 16
              i32.add
              local.set $$t38.1
              local.get $$t38.0
              call $runtime.Block.Release
              local.set $$t38.0

              ;;*t38
              local.get $$t38.1
              i32.load offset=0 align=4
              call $runtime.Block.Retain
              local.get $$t38.1
              i32.load offset=4 align=4
              local.get $$t38.1
              i32.load offset=8 align=4
              local.get $$t38.1
              i32.load offset=12 align=4
              local.set $$t39.3
              local.set $$t39.2
              local.set $$t39.1
              local.get $$t39.0
              call $runtime.Block.Release
              local.set $$t39.0

              ;;len(t39)
              local.get $$t39.2
              local.set $$t40

              ;;t40 - 1:int
              local.get $$t40
              i32.const 1
              i32.sub
              local.set $$t41

              ;;slice t37[:t41]
              local.get $$t37.0
              call $runtime.Block.Retain
              local.get $$t37.1
              i32.const 8
              i32.const 0
              i32.mul
              i32.add
              local.get $$t41
              i32.const 0
              i32.sub
              local.get $$t37.3
              i32.const 0
              i32.sub
              local.set $$t42.3
              local.set $$t42.2
              local.set $$t42.1
              local.get $$t42.0
              call $runtime.Block.Release
              local.set $$t42.0

              ;;*t35 = t42
              local.get $$t35.1
              local.get $$t42.0
              call $runtime.Block.Retain
              local.get $$t35.1
              i32.load offset=0 align=1
              call $runtime.Block.Release
              i32.store offset=0 align=1
              local.get $$t35.1
              local.get $$t42.1
              i32.store offset=4 align=4
              local.get $$t35.1
              local.get $$t42.2
              i32.store offset=8 align=4
              local.get $$t35.1
              local.get $$t42.3
              i32.store offset=12 align=4

              ;;return
              br $$BlockFnBody

            end ;;$Block_4
            i32.const 5
            local.set $$current_block

            ;;&t20.Left [#2]
            local.get $$t20.0
            call $runtime.Block.Retain
            local.get $$t20.1
            i32.const 8
            i32.add
            local.set $$t43.1
            local.get $$t43.0
            call $runtime.Block.Release
            local.set $$t43.0

            ;;*t43
            local.get $$t43.1
            i32.load offset=0 align=4
            call $runtime.Block.Retain
            local.get $$t43.1
            i32.load offset=4 align=4
            local.set $$t44.1
            local.get $$t44.0
            call $runtime.Block.Release
            local.set $$t44.0

            ;;(*mapNode).SetParent(t44, t20)
            local.get $$t44.0
            local.get $$t44.1
            local.get $$t20.0
            local.get $$t20.1
            call $runtime.mapNode.SetParent

            ;;jump 6
            br $$Block_5

          end ;;$Block_5
          i32.const 6
          local.set $$current_block

          ;;&t20.Right [#3]
          local.get $$t20.0
          call $runtime.Block.Retain
          local.get $$t20.1
          i32.const 16
          i32.add
          local.set $$t45.1
          local.get $$t45.0
          call $runtime.Block.Release
          local.set $$t45.0

          ;;*t46
          local.get $$t45.1
          i32.load offset=0 align=4
          call $runtime.Block.Retain
          local.get $$t45.1
          i32.load offset=4 align=4
          local.set $$t46.1
          local.get $$t46.0
          call $runtime.Block.Release
          local.set $$t46.0

          ;;&this.NIL [#0]
          local.get $this.0
          call $runtime.Block.Retain
          local.get $this.1
          i32.const 0
          i32.add
          local.set $$t47.1
          local.get $$t47.0
          call $runtime.Block.Release
          local.set $$t47.0

          ;;*t48
          local.get $$t47.1
          i32.load offset=0 align=4
          call $runtime.Block.Retain
          local.get $$t47.1
          i32.load offset=4 align=4
          local.set $$t48.1
          local.get $$t48.0
          call $runtime.Block.Release
          local.set $$t48.0

          ;;t47 != t49
          local.get $$t46.1
          local.get $$t48.1
          i32.eq
          i32.eqz
          local.set $$t49

          ;;if t50 goto 7 else 4
          local.get $$t49
          if
            br $$Block_6
          else
            i32.const 4
            local.set $$block_selector
            br $$BlockDisp
          end

        end ;;$Block_6
        i32.const 7
        local.set $$current_block

        ;;&t20.Right [#3]
        local.get $$t20.0
        call $runtime.Block.Retain
        local.get $$t20.1
        i32.const 16
        i32.add
        local.set $$t50.1
        local.get $$t50.0
        call $runtime.Block.Release
        local.set $$t50.0

        ;;*t51
        local.get $$t50.1
        i32.load offset=0 align=4
        call $runtime.Block.Retain
        local.get $$t50.1
        i32.load offset=4 align=4
        local.set $$t51.1
        local.get $$t51.0
        call $runtime.Block.Release
        local.set $$t51.0

        ;;(*mapNode).SetParent(t52, t20)
        local.get $$t51.0
        local.get $$t51.1
        local.get $$t20.0
        local.get $$t20.1
        call $runtime.mapNode.SetParent

        ;;jump 4
        i32.const 4
        local.set $$block_selector
        br $$BlockDisp

      end ;;$Block_7
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$t0.0
  call $runtime.Block.Release
  local.get $$t1.0
  call $runtime.Block.Release
  local.get $$t2.0
  call $runtime.Block.Release
  local.get $$t5.0
  call $runtime.Block.Release
  local.get $$t7.0
  call $runtime.Block.Release
  local.get $$t8.0
  call $runtime.Block.Release
  local.get $$t13.0
  call $runtime.Block.Release
  local.get $$t14.0
  call $runtime.Block.Release
  local.get $$t15.0
  call $runtime.Block.Release
  local.get $$t16.0
  call $runtime.Block.Release
  local.get $$t19.0
  call $runtime.Block.Release
  local.get $$t20.0
  call $runtime.Block.Release
  local.get $$t21.0
  call $runtime.Block.Release
  local.get $$t22.0
  call $runtime.Block.Release
  local.get $$t24.0
  call $runtime.Block.Release
  local.get $$t25.0
  call $runtime.Block.Release
  local.get $$t26.0
  call $runtime.Block.Release
  local.get $$t29.0
  call $runtime.Block.Release
  local.get $$t30.0
  call $runtime.Block.Release
  local.get $$t31.0
  call $runtime.Block.Release
  local.get $$t32.0
  call $runtime.Block.Release
  local.get $$t33.0
  call $runtime.Block.Release
  local.get $$t35.0
  call $runtime.Block.Release
  local.get $$t36.0
  call $runtime.Block.Release
  local.get $$t37.0
  call $runtime.Block.Release
  local.get $$t38.0
  call $runtime.Block.Release
  local.get $$t39.0
  call $runtime.Block.Release
  local.get $$t42.0
  call $runtime.Block.Release
  local.get $$t43.0
  call $runtime.Block.Release
  local.get $$t44.0
  call $runtime.Block.Release
  local.get $$t45.0
  call $runtime.Block.Release
  local.get $$t46.0
  call $runtime.Block.Release
  local.get $$t47.0
  call $runtime.Block.Release
  local.get $$t48.0
  call $runtime.Block.Release
  local.get $$t50.0
  call $runtime.Block.Release
  local.get $$t51.0
  call $runtime.Block.Release
) ;;runtime.mapImp.Delete

(func $runtime.mapImp.Len (param $this.0 i32) (param $this.1 i32) (result i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$ret_0 i32)
  (local $$t0.0 i32)
  (local $$t0.1 i32)
  (local $$t1.0 i32)
  (local $$t1.1 i32)
  (local $$t1.2 i32)
  (local $$t1.3 i32)
  (local $$t2 i32)
  (local $$t3 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;&this.nodes [#2]
        local.get $this.0
        call $runtime.Block.Retain
        local.get $this.1
        i32.const 16
        i32.add
        local.set $$t0.1
        local.get $$t0.0
        call $runtime.Block.Release
        local.set $$t0.0

        ;;*t0
        local.get $$t0.1
        i32.load offset=0 align=4
        call $runtime.Block.Retain
        local.get $$t0.1
        i32.load offset=4 align=4
        local.get $$t0.1
        i32.load offset=8 align=4
        local.get $$t0.1
        i32.load offset=12 align=4
        local.set $$t1.3
        local.set $$t1.2
        local.set $$t1.1
        local.get $$t1.0
        call $runtime.Block.Release
        local.set $$t1.0

        ;;len(t1)
        local.get $$t1.2
        local.set $$t2

        ;;t2 - 1:int
        local.get $$t2
        i32.const 1
        i32.sub
        local.set $$t3

        ;;return t3
        local.get $$t3
        local.set $$ret_0
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$ret_0
  local.get $$t0.0
  call $runtime.Block.Release
  local.get $$t1.0
  call $runtime.Block.Release
) ;;runtime.mapImp.Len

(func $runtime.mapImp.Lookup (param $this.0 i32) (param $this.1 i32) (param $k.0.0 i32) (param $k.0.1 i32) (param $k.1 i32) (param $k.2 i32) (result i32 i32 i32 i32 i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$ret_0.0.0 i32)
  (local $$ret_0.0.1 i32)
  (local $$ret_0.1 i32)
  (local $$ret_0.2 i32)
  (local $$ret_1 i32)
  (local $$t0.0 i32)
  (local $$t0.1 i32)
  (local $$t1.0 i32)
  (local $$t1.1 i32)
  (local $$t2.0 i32)
  (local $$t2.1 i32)
  (local $$t3 i32)
  (local $$t4.0 i32)
  (local $$t4.1 i32)
  (local $$t5.0.0 i32)
  (local $$t5.0.1 i32)
  (local $$t5.1 i32)
  (local $$t5.2 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_2
        block $$Block_1
          block $$Block_0
            block $$BlockSel
              local.get $$block_selector
              br_table 0 1 2 0
            end ;;$BlockSel
            i32.const 0
            local.set $$current_block

            ;;(*mapImp).search(this, k)
            local.get $this.0
            local.get $this.1
            local.get $k.0.0
            local.get $k.0.1
            local.get $k.1
            local.get $k.2
            call $runtime.mapImp.search
            local.set $$t0.1
            local.get $$t0.0
            call $runtime.Block.Release
            local.set $$t0.0

            ;;&this.NIL [#0]
            local.get $this.0
            call $runtime.Block.Retain
            local.get $this.1
            i32.const 0
            i32.add
            local.set $$t1.1
            local.get $$t1.0
            call $runtime.Block.Release
            local.set $$t1.0

            ;;*t1
            local.get $$t1.1
            i32.load offset=0 align=4
            call $runtime.Block.Retain
            local.get $$t1.1
            i32.load offset=4 align=4
            local.set $$t2.1
            local.get $$t2.0
            call $runtime.Block.Release
            local.set $$t2.0

            ;;t0 != t2
            local.get $$t0.1
            local.get $$t2.1
            i32.eq
            i32.eqz
            local.set $$t3

            ;;if t3 goto 1 else 2
            local.get $$t3
            if
              br $$Block_0
            else
              br $$Block_1
            end

          end ;;$Block_0
          i32.const 1
          local.set $$current_block

          ;;&t0.Val [#6]
          local.get $$t0.0
          call $runtime.Block.Retain
          local.get $$t0.1
          i32.const 44
          i32.add
          local.set $$t4.1
          local.get $$t4.0
          call $runtime.Block.Release
          local.set $$t4.0

          ;;*t4
          local.get $$t4.1
          i32.load offset=0 align=4
          call $runtime.Block.Retain
          local.get $$t4.1
          i32.load offset=4 align=4
          local.get $$t4.1
          i32.load offset=8 align=4
          local.get $$t4.1
          i32.load offset=12 align=4
          local.set $$t5.2
          local.set $$t5.1
          local.set $$t5.0.1
          local.get $$t5.0.0
          call $runtime.Block.Release
          local.set $$t5.0.0

          ;;return t5, true:bool
          local.get $$t5.0.0
          call $runtime.Block.Retain
          local.get $$t5.0.1
          local.get $$t5.1
          local.get $$t5.2
          local.set $$ret_0.2
          local.set $$ret_0.1
          local.set $$ret_0.0.1
          local.get $$ret_0.0.0
          call $runtime.Block.Release
          local.set $$ret_0.0.0
          i32.const 1
          local.set $$ret_1
          br $$BlockFnBody

        end ;;$Block_1
        i32.const 2
        local.set $$current_block

        ;;return nil:interface{}, false:bool
        i32.const 0
        i32.const 0
        i32.const 0
        i32.const 0
        local.set $$ret_0.2
        local.set $$ret_0.1
        local.set $$ret_0.0.1
        local.get $$ret_0.0.0
        call $runtime.Block.Release
        local.set $$ret_0.0.0
        i32.const 0
        local.set $$ret_1
        br $$BlockFnBody

      end ;;$Block_2
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$ret_0.0.0
  call $runtime.Block.Retain
  local.get $$ret_0.0.1
  local.get $$ret_0.1
  local.get $$ret_0.2
  local.get $$ret_1
  local.get $$ret_0.0.0
  call $runtime.Block.Release
  local.get $$t0.0
  call $runtime.Block.Release
  local.get $$t1.0
  call $runtime.Block.Release
  local.get $$t2.0
  call $runtime.Block.Release
  local.get $$t4.0
  call $runtime.Block.Release
  local.get $$t5.0.0
  call $runtime.Block.Release
) ;;runtime.mapImp.Lookup

(func $$runtime.mapNode.$ref.$slice.append (param $x.0 i32) (param $x.1 i32) (param $x.2 i32) (param $x.3 i32) (param $y.0 i32) (param $y.1 i32) (param $y.2 i32) (param $y.3 i32) (result i32 i32 i32 i32)
  (local $item.0 i32)
  (local $item.1 i32)
  (local $x_len i32)
  (local $y_len i32)
  (local $new_len i32)
  (local $src i32)
  (local $dest i32)
  (local $new_cap i32)
  local.get $x.2
  local.set $x_len
  local.get $y.2
  local.set $y_len
  local.get $x_len
  local.get $y_len
  i32.add
  local.set $new_len
  local.get $new_len
  local.get $x.3
  i32.le_u
  if (result i32 i32 i32 i32)
    local.get $x.0
    call $runtime.Block.Retain
    local.get $x.1
    local.get $new_len
    local.get $x.3
    local.get $y.1
    local.set $src
    local.get $x.1
    i32.const 8
    local.get $x_len
    i32.mul
    i32.add
    local.set $dest
    block $block1
      loop $loop1
        local.get $y_len
        i32.eqz
        if
          br $block1
        else
        end
        local.get $src
        i32.load offset=0 align=4
        call $runtime.Block.Retain
        local.get $src
        i32.load offset=4 align=4
        local.set $item.1
        local.get $item.0
        call $runtime.Block.Release
        local.set $item.0
        local.get $dest
        local.get $item.0
        call $runtime.Block.Retain
        local.get $dest
        i32.load offset=0 align=1
        call $runtime.Block.Release
        i32.store offset=0 align=1
        local.get $dest
        local.get $item.1
        i32.store offset=4 align=4
        local.get $src
        i32.const 8
        i32.add
        local.set $src
        local.get $dest
        i32.const 8
        i32.add
        local.set $dest
        local.get $y_len
        i32.const 1
        i32.sub
        local.set $y_len
        br $loop1
      end ;;loop1
    end ;;block1
  else
    local.get $new_len
    i32.const 2
    i32.mul
    local.set $new_cap
    local.get $new_cap
    i32.const 8
    i32.mul
    i32.const 16
    i32.add
    call $runtime.HeapAlloc
    local.get $new_cap
    i32.const 8
    i32.const 8
    call $runtime.Block.Init
    call $runtime.DupI32
    i32.const 16
    i32.add
    call $runtime.DupI32
    local.set $dest
    local.get $new_len
    local.get $new_cap
    local.get $x.1
    local.set $src
    block $block2
      loop $loop2
        local.get $x_len
        i32.eqz
        if
          br $block2
        else
        end
        local.get $src
        i32.load offset=0 align=4
        call $runtime.Block.Retain
        local.get $src
        i32.load offset=4 align=4
        local.set $item.1
        local.get $item.0
        call $runtime.Block.Release
        local.set $item.0
        local.get $dest
        local.get $item.0
        call $runtime.Block.Retain
        local.get $dest
        i32.load offset=0 align=1
        call $runtime.Block.Release
        i32.store offset=0 align=1
        local.get $dest
        local.get $item.1
        i32.store offset=4 align=4
        local.get $src
        i32.const 8
        i32.add
        local.set $src
        local.get $dest
        i32.const 8
        i32.add
        local.set $dest
        local.get $x_len
        i32.const 1
        i32.sub
        local.set $x_len
        br $loop2
      end ;;loop2
    end ;;block2
    local.get $y.1
    local.set $src
    block $block3
      loop $loop3
        local.get $y_len
        i32.eqz
        if
          br $block3
        else
        end
        local.get $src
        i32.load offset=0 align=4
        call $runtime.Block.Retain
        local.get $src
        i32.load offset=4 align=4
        local.set $item.1
        local.get $item.0
        call $runtime.Block.Release
        local.set $item.0
        local.get $dest
        local.get $item.0
        call $runtime.Block.Retain
        local.get $dest
        i32.load offset=0 align=1
        call $runtime.Block.Release
        i32.store offset=0 align=1
        local.get $dest
        local.get $item.1
        i32.store offset=4 align=4
        local.get $src
        i32.const 8
        i32.add
        local.set $src
        local.get $dest
        i32.const 8
        i32.add
        local.set $dest
        local.get $y_len
        i32.const 1
        i32.sub
        local.set $y_len
        br $loop3
      end ;;loop3
    end ;;block3
  end
  local.get $item.0
  call $runtime.Block.Release
) ;;$runtime.mapNode.$ref.$slice.append

(func $runtime.mapImp.Update (param $this.0 i32) (param $this.1 i32) (param $k.0.0 i32) (param $k.0.1 i32) (param $k.1 i32) (param $k.2 i32) (param $v.0.0 i32) (param $v.0.1 i32) (param $v.1 i32) (param $v.2 i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$t0.0 i32)
  (local $$t0.1 i32)
  (local $$t1.0 i32)
  (local $$t1.1 i32)
  (local $$t2.0 i32)
  (local $$t2.1 i32)
  (local $$t3 i32)
  (local $$t4.0 i32)
  (local $$t4.1 i32)
  (local $$t5.0 i32)
  (local $$t5.1 i32)
  (local $$t6.0 i32)
  (local $$t6.1 i32)
  (local $$t7.0 i32)
  (local $$t7.1 i32)
  (local $$t7.2 i32)
  (local $$t7.3 i32)
  (local $$t8 i32)
  (local $$t9 i32)
  (local $$t10.0 i32)
  (local $$t10.1 i32)
  (local $$t11.0 i32)
  (local $$t11.1 i32)
  (local $$t12.0 i32)
  (local $$t12.1 i32)
  (local $$t13.0 i32)
  (local $$t13.1 i32)
  (local $$t14.0 i32)
  (local $$t14.1 i32)
  (local $$t15.0 i32)
  (local $$t15.1 i32)
  (local $$t16.0 i32)
  (local $$t16.1 i32)
  (local $$t17.0 i32)
  (local $$t17.1 i32)
  (local $$t18.0 i32)
  (local $$t18.1 i32)
  (local $$t19.0 i32)
  (local $$t19.1 i32)
  (local $$t20.0 i32)
  (local $$t20.1 i32)
  (local $$t21.0 i32)
  (local $$t21.1 i32)
  (local $$t21.2 i32)
  (local $$t21.3 i32)
  (local $$t22.0 i32)
  (local $$t22.1 i32)
  (local $$t23.0 i32)
  (local $$t23.1 i32)
  (local $$t24.0 i32)
  (local $$t24.1 i32)
  (local $$t24.2 i32)
  (local $$t24.3 i32)
  (local $$t25.0 i32)
  (local $$t25.1 i32)
  (local $$t25.2 i32)
  (local $$t25.3 i32)
  (local $$t26.0 i32)
  (local $$t26.1 i32)
  (local $$t27.0 i32)
  (local $$t27.1 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_3
        block $$Block_2
          block $$Block_1
            block $$Block_0
              block $$BlockSel
                local.get $$block_selector
                br_table 0 1 2 3 0
              end ;;$BlockSel
              i32.const 0
              local.set $$current_block

              ;;(*mapImp).search(this, k)
              local.get $this.0
              local.get $this.1
              local.get $k.0.0
              local.get $k.0.1
              local.get $k.1
              local.get $k.2
              call $runtime.mapImp.search
              local.set $$t0.1
              local.get $$t0.0
              call $runtime.Block.Release
              local.set $$t0.0

              ;;&this.NIL [#0]
              local.get $this.0
              call $runtime.Block.Retain
              local.get $this.1
              i32.const 0
              i32.add
              local.set $$t1.1
              local.get $$t1.0
              call $runtime.Block.Release
              local.set $$t1.0

              ;;*t1
              local.get $$t1.1
              i32.load offset=0 align=4
              call $runtime.Block.Retain
              local.get $$t1.1
              i32.load offset=4 align=4
              local.set $$t2.1
              local.get $$t2.0
              call $runtime.Block.Release
              local.set $$t2.0

              ;;t0 == t2
              local.get $$t0.1
              local.get $$t2.1
              i32.eq
              local.set $$t3

              ;;if t3 goto 1 else 3
              local.get $$t3
              if
                br $$Block_0
              else
                br $$Block_2
              end

            end ;;$Block_0
            i32.const 1
            local.set $$current_block

            ;;new mapNode (complit)
            i32.const 76
            call $runtime.HeapAlloc
            i32.const 1
            i32.const 12
            i32.const 60
            call $runtime.Block.Init
            call $runtime.DupI32
            i32.const 16
            i32.add
            local.set $$t4.1
            local.get $$t4.0
            call $runtime.Block.Release
            local.set $$t4.0

            ;;&t4.NodeIdx [#1]
            local.get $$t4.0
            call $runtime.Block.Retain
            local.get $$t4.1
            i32.const 4
            i32.add
            local.set $$t5.1
            local.get $$t5.0
            call $runtime.Block.Release
            local.set $$t5.0

            ;;&this.nodes [#2]
            local.get $this.0
            call $runtime.Block.Retain
            local.get $this.1
            i32.const 16
            i32.add
            local.set $$t6.1
            local.get $$t6.0
            call $runtime.Block.Release
            local.set $$t6.0

            ;;*t6
            local.get $$t6.1
            i32.load offset=0 align=4
            call $runtime.Block.Retain
            local.get $$t6.1
            i32.load offset=4 align=4
            local.get $$t6.1
            i32.load offset=8 align=4
            local.get $$t6.1
            i32.load offset=12 align=4
            local.set $$t7.3
            local.set $$t7.2
            local.set $$t7.1
            local.get $$t7.0
            call $runtime.Block.Release
            local.set $$t7.0

            ;;len(t7)
            local.get $$t7.2
            local.set $$t8

            ;;changetype mapNodeIdx <- int (t8)
            local.get $$t8
            local.set $$t9

            ;;&t4.Left [#2]
            local.get $$t4.0
            call $runtime.Block.Retain
            local.get $$t4.1
            i32.const 8
            i32.add
            local.set $$t10.1
            local.get $$t10.0
            call $runtime.Block.Release
            local.set $$t10.0

            ;;&this.NIL [#0]
            local.get $this.0
            call $runtime.Block.Retain
            local.get $this.1
            i32.const 0
            i32.add
            local.set $$t11.1
            local.get $$t11.0
            call $runtime.Block.Release
            local.set $$t11.0

            ;;*t11
            local.get $$t11.1
            i32.load offset=0 align=4
            call $runtime.Block.Retain
            local.get $$t11.1
            i32.load offset=4 align=4
            local.set $$t12.1
            local.get $$t12.0
            call $runtime.Block.Release
            local.set $$t12.0

            ;;&t4.Right [#3]
            local.get $$t4.0
            call $runtime.Block.Retain
            local.get $$t4.1
            i32.const 16
            i32.add
            local.set $$t13.1
            local.get $$t13.0
            call $runtime.Block.Release
            local.set $$t13.0

            ;;&this.NIL [#0]
            local.get $this.0
            call $runtime.Block.Retain
            local.get $this.1
            i32.const 0
            i32.add
            local.set $$t14.1
            local.get $$t14.0
            call $runtime.Block.Release
            local.set $$t14.0

            ;;*t14
            local.get $$t14.1
            i32.load offset=0 align=4
            call $runtime.Block.Retain
            local.get $$t14.1
            i32.load offset=4 align=4
            local.set $$t15.1
            local.get $$t15.0
            call $runtime.Block.Release
            local.set $$t15.0

            ;;&t4.Color [#4]
            local.get $$t4.0
            call $runtime.Block.Retain
            local.get $$t4.1
            i32.const 24
            i32.add
            local.set $$t16.1
            local.get $$t16.0
            call $runtime.Block.Release
            local.set $$t16.0

            ;;&t4.Key [#5]
            local.get $$t4.0
            call $runtime.Block.Retain
            local.get $$t4.1
            i32.const 28
            i32.add
            local.set $$t17.1
            local.get $$t17.0
            call $runtime.Block.Release
            local.set $$t17.0

            ;;&t4.Val [#6]
            local.get $$t4.0
            call $runtime.Block.Retain
            local.get $$t4.1
            i32.const 44
            i32.add
            local.set $$t18.1
            local.get $$t18.0
            call $runtime.Block.Release
            local.set $$t18.0

            ;;*t5 = t9
            local.get $$t5.1
            local.get $$t9
            i32.store offset=0 align=4

            ;;*t10 = t12
            local.get $$t10.1
            local.get $$t12.0
            call $runtime.Block.Retain
            local.get $$t10.1
            i32.load offset=0 align=1
            call $runtime.Block.Release
            i32.store offset=0 align=1
            local.get $$t10.1
            local.get $$t12.1
            i32.store offset=4 align=4

            ;;*t13 = t15
            local.get $$t13.1
            local.get $$t15.0
            call $runtime.Block.Retain
            local.get $$t13.1
            i32.load offset=0 align=1
            call $runtime.Block.Release
            i32.store offset=0 align=1
            local.get $$t13.1
            local.get $$t15.1
            i32.store offset=4 align=4

            ;;*t16 = 0:mapColor
            local.get $$t16.1
            i32.const 0
            i32.store offset=0 align=4

            ;;*t17 = k
            local.get $$t17.1
            local.get $k.0.0
            call $runtime.Block.Retain
            local.get $$t17.1
            i32.load offset=0 align=1
            call $runtime.Block.Release
            i32.store offset=0 align=1
            local.get $$t17.1
            local.get $k.0.1
            i32.store offset=4 align=4
            local.get $$t17.1
            local.get $k.1
            i32.store offset=8 align=4
            local.get $$t17.1
            local.get $k.2
            i32.store offset=12 align=4

            ;;*t18 = v
            local.get $$t18.1
            local.get $v.0.0
            call $runtime.Block.Retain
            local.get $$t18.1
            i32.load offset=0 align=1
            call $runtime.Block.Release
            i32.store offset=0 align=1
            local.get $$t18.1
            local.get $v.0.1
            i32.store offset=4 align=4
            local.get $$t18.1
            local.get $v.1
            i32.store offset=8 align=4
            local.get $$t18.1
            local.get $v.2
            i32.store offset=12 align=4

            ;;&this.nodes [#2]
            local.get $this.0
            call $runtime.Block.Retain
            local.get $this.1
            i32.const 16
            i32.add
            local.set $$t19.1
            local.get $$t19.0
            call $runtime.Block.Release
            local.set $$t19.0

            ;;&this.nodes [#2]
            local.get $this.0
            call $runtime.Block.Retain
            local.get $this.1
            i32.const 16
            i32.add
            local.set $$t20.1
            local.get $$t20.0
            call $runtime.Block.Release
            local.set $$t20.0

            ;;*t20
            local.get $$t20.1
            i32.load offset=0 align=4
            call $runtime.Block.Retain
            local.get $$t20.1
            i32.load offset=4 align=4
            local.get $$t20.1
            i32.load offset=8 align=4
            local.get $$t20.1
            i32.load offset=12 align=4
            local.set $$t21.3
            local.set $$t21.2
            local.set $$t21.1
            local.get $$t21.0
            call $runtime.Block.Release
            local.set $$t21.0

            ;;new [1]*mapNode (varargs)
            i32.const 24
            call $runtime.HeapAlloc
            i32.const 1
            i32.const 16
            i32.const 8
            call $runtime.Block.Init
            call $runtime.DupI32
            i32.const 16
            i32.add
            local.set $$t22.1
            local.get $$t22.0
            call $runtime.Block.Release
            local.set $$t22.0

            ;;&t22[0:int]
            local.get $$t22.0
            call $runtime.Block.Retain
            local.get $$t22.1
            i32.const 8
            i32.const 0
            i32.mul
            i32.add
            local.set $$t23.1
            local.get $$t23.0
            call $runtime.Block.Release
            local.set $$t23.0

            ;;*t23 = t4
            local.get $$t23.1
            local.get $$t4.0
            call $runtime.Block.Retain
            local.get $$t23.1
            i32.load offset=0 align=1
            call $runtime.Block.Release
            i32.store offset=0 align=1
            local.get $$t23.1
            local.get $$t4.1
            i32.store offset=4 align=4

            ;;slice t22[:]
            local.get $$t22.0
            call $runtime.Block.Retain
            local.get $$t22.1
            i32.const 8
            i32.const 0
            i32.mul
            i32.add
            i32.const 1
            i32.const 0
            i32.sub
            i32.const 1
            i32.const 0
            i32.sub
            local.set $$t24.3
            local.set $$t24.2
            local.set $$t24.1
            local.get $$t24.0
            call $runtime.Block.Release
            local.set $$t24.0

            ;;append(t21, t24...)
            local.get $$t21.0
            local.get $$t21.1
            local.get $$t21.2
            local.get $$t21.3
            local.get $$t24.0
            local.get $$t24.1
            local.get $$t24.2
            local.get $$t24.3
            call $$runtime.mapNode.$ref.$slice.append
            local.set $$t25.3
            local.set $$t25.2
            local.set $$t25.1
            local.get $$t25.0
            call $runtime.Block.Release
            local.set $$t25.0

            ;;*t19 = t25
            local.get $$t19.1
            local.get $$t25.0
            call $runtime.Block.Retain
            local.get $$t19.1
            i32.load offset=0 align=1
            call $runtime.Block.Release
            i32.store offset=0 align=1
            local.get $$t19.1
            local.get $$t25.1
            i32.store offset=4 align=4
            local.get $$t19.1
            local.get $$t25.2
            i32.store offset=8 align=4
            local.get $$t19.1
            local.get $$t25.3
            i32.store offset=12 align=4

            ;;(*mapImp).insert(this, t4)
            local.get $this.0
            local.get $this.1
            local.get $$t4.0
            local.get $$t4.1
            call $runtime.mapImp.insert
            local.set $$t26.1
            local.get $$t26.0
            call $runtime.Block.Release
            local.set $$t26.0

            ;;jump 2
            br $$Block_1

          end ;;$Block_1
          i32.const 2
          local.set $$current_block

          ;;return
          br $$BlockFnBody

        end ;;$Block_2
        i32.const 3
        local.set $$current_block

        ;;&t0.Val [#6]
        local.get $$t0.0
        call $runtime.Block.Retain
        local.get $$t0.1
        i32.const 44
        i32.add
        local.set $$t27.1
        local.get $$t27.0
        call $runtime.Block.Release
        local.set $$t27.0

        ;;*t27 = v
        local.get $$t27.1
        local.get $v.0.0
        call $runtime.Block.Retain
        local.get $$t27.1
        i32.load offset=0 align=1
        call $runtime.Block.Release
        i32.store offset=0 align=1
        local.get $$t27.1
        local.get $v.0.1
        i32.store offset=4 align=4
        local.get $$t27.1
        local.get $v.1
        i32.store offset=8 align=4
        local.get $$t27.1
        local.get $v.2
        i32.store offset=12 align=4

        ;;jump 2
        i32.const 2
        local.set $$block_selector
        br $$BlockDisp

      end ;;$Block_3
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$t0.0
  call $runtime.Block.Release
  local.get $$t1.0
  call $runtime.Block.Release
  local.get $$t2.0
  call $runtime.Block.Release
  local.get $$t4.0
  call $runtime.Block.Release
  local.get $$t5.0
  call $runtime.Block.Release
  local.get $$t6.0
  call $runtime.Block.Release
  local.get $$t7.0
  call $runtime.Block.Release
  local.get $$t10.0
  call $runtime.Block.Release
  local.get $$t11.0
  call $runtime.Block.Release
  local.get $$t12.0
  call $runtime.Block.Release
  local.get $$t13.0
  call $runtime.Block.Release
  local.get $$t14.0
  call $runtime.Block.Release
  local.get $$t15.0
  call $runtime.Block.Release
  local.get $$t16.0
  call $runtime.Block.Release
  local.get $$t17.0
  call $runtime.Block.Release
  local.get $$t18.0
  call $runtime.Block.Release
  local.get $$t19.0
  call $runtime.Block.Release
  local.get $$t20.0
  call $runtime.Block.Release
  local.get $$t21.0
  call $runtime.Block.Release
  local.get $$t22.0
  call $runtime.Block.Release
  local.get $$t23.0
  call $runtime.Block.Release
  local.get $$t24.0
  call $runtime.Block.Release
  local.get $$t25.0
  call $runtime.Block.Release
  local.get $$t26.0
  call $runtime.Block.Release
  local.get $$t27.0
  call $runtime.Block.Release
) ;;runtime.mapImp.Update

(func $runtime.mapImp.delete (param $this.0 i32) (param $this.1 i32) (param $z.0 i32) (param $z.1 i32) (result i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$ret_0 i32)
  (local $$t0.0 i32)
  (local $$t0.1 i32)
  (local $$t1.0 i32)
  (local $$t1.1 i32)
  (local $$t2.0 i32)
  (local $$t2.1 i32)
  (local $$t3.0 i32)
  (local $$t3.1 i32)
  (local $$t4 i32)
  (local $$t5.0 i32)
  (local $$t5.1 i32)
  (local $$t6.0 i32)
  (local $$t6.1 i32)
  (local $$t7.0 i32)
  (local $$t7.1 i32)
  (local $$t8.0 i32)
  (local $$t8.1 i32)
  (local $$t9.0 i32)
  (local $$t9.1 i32)
  (local $$t10.0 i32)
  (local $$t10.1 i32)
  (local $$t11 i32)
  (local $$t12.0 i32)
  (local $$t12.1 i32)
  (local $$t13.0 i32)
  (local $$t13.1 i32)
  (local $$t14.0 i32)
  (local $$t14.1 i32)
  (local $$t15.0 i32)
  (local $$t15.1 i32)
  (local $$t16 i32)
  (local $$t17.0 i32)
  (local $$t17.1 i32)
  (local $$t18.0 i32)
  (local $$t18.1 i32)
  (local $$t19.0 i32)
  (local $$t19.1 i32)
  (local $$t20.0 i32)
  (local $$t20.1 i32)
  (local $$t21.0 i32)
  (local $$t21.1 i32)
  (local $$t22.0 i32)
  (local $$t22.1 i32)
  (local $$t23.0 i32)
  (local $$t23.1 i32)
  (local $$t24.0 i32)
  (local $$t24.1 i32)
  (local $$t25 i32)
  (local $$t26.0 i32)
  (local $$t26.1 i32)
  (local $$t27.0 i32)
  (local $$t27.1 i32)
  (local $$t28 i32)
  (local $$t29.0 i32)
  (local $$t29.1 i32)
  (local $$t30.0 i32)
  (local $$t30.1 i32)
  (local $$t31.0 i32)
  (local $$t31.1 i32)
  (local $$t32 i32)
  (local $$t33.0 i32)
  (local $$t33.1 i32)
  (local $$t34.0 i32)
  (local $$t34.1 i32)
  (local $$t35.0 i32)
  (local $$t35.1 i32)
  (local $$t36.0 i32)
  (local $$t36.1 i32)
  (local $$t37.0 i32)
  (local $$t37.1 i32)
  (local $$t38 i32)
  (local $$t39 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_16
        block $$Block_15
          block $$Block_14
            block $$Block_13
              block $$Block_12
                block $$Block_11
                  block $$Block_10
                    block $$Block_9
                      block $$Block_8
                        block $$Block_7
                          block $$Block_6
                            block $$Block_5
                              block $$Block_4
                                block $$Block_3
                                  block $$Block_2
                                    block $$Block_1
                                      block $$Block_0
                                        block $$BlockSel
                                          local.get $$block_selector
                                          br_table 0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 0
                                        end ;;$BlockSel
                                        i32.const 0
                                        local.set $$current_block

                                        ;;&z.Left [#2]
                                        local.get $z.0
                                        call $runtime.Block.Retain
                                        local.get $z.1
                                        i32.const 8
                                        i32.add
                                        local.set $$t0.1
                                        local.get $$t0.0
                                        call $runtime.Block.Release
                                        local.set $$t0.0

                                        ;;*t0
                                        local.get $$t0.1
                                        i32.load offset=0 align=4
                                        call $runtime.Block.Retain
                                        local.get $$t0.1
                                        i32.load offset=4 align=4
                                        local.set $$t1.1
                                        local.get $$t1.0
                                        call $runtime.Block.Release
                                        local.set $$t1.0

                                        ;;&this.NIL [#0]
                                        local.get $this.0
                                        call $runtime.Block.Retain
                                        local.get $this.1
                                        i32.const 0
                                        i32.add
                                        local.set $$t2.1
                                        local.get $$t2.0
                                        call $runtime.Block.Release
                                        local.set $$t2.0

                                        ;;*t2
                                        local.get $$t2.1
                                        i32.load offset=0 align=4
                                        call $runtime.Block.Retain
                                        local.get $$t2.1
                                        i32.load offset=4 align=4
                                        local.set $$t3.1
                                        local.get $$t3.0
                                        call $runtime.Block.Release
                                        local.set $$t3.0

                                        ;;t1 == t3
                                        local.get $$t1.1
                                        local.get $$t3.1
                                        i32.eq
                                        local.set $$t4

                                        ;;if t4 goto 1 else 4
                                        local.get $$t4
                                        if
                                          br $$Block_0
                                        else
                                          br $$Block_3
                                        end

                                      end ;;$Block_0
                                      i32.const 1
                                      local.set $$current_block

                                      ;;jump 2
                                      br $$Block_1

                                    end ;;$Block_1
                                    local.get $$current_block
                                    i32.const 1
                                    i32.eq
                                    if (result i32 i32)
                                      local.get $z.0
                                      call $runtime.Block.Retain
                                      local.get $z.1
                                    else
                                      local.get $$t5.0
                                      call $runtime.Block.Retain
                                      local.get $$t5.1
                                    end
                                    local.set $$t6.1
                                    local.get $$t6.0
                                    call $runtime.Block.Release
                                    local.set $$t6.0
                                    i32.const 2
                                    local.set $$current_block

                                    ;;&t5.Left [#2]
                                    local.get $$t6.0
                                    call $runtime.Block.Retain
                                    local.get $$t6.1
                                    i32.const 8
                                    i32.add
                                    local.set $$t7.1
                                    local.get $$t7.0
                                    call $runtime.Block.Release
                                    local.set $$t7.0

                                    ;;*t6
                                    local.get $$t7.1
                                    i32.load offset=0 align=4
                                    call $runtime.Block.Retain
                                    local.get $$t7.1
                                    i32.load offset=4 align=4
                                    local.set $$t8.1
                                    local.get $$t8.0
                                    call $runtime.Block.Release
                                    local.set $$t8.0

                                    ;;&this.NIL [#0]
                                    local.get $this.0
                                    call $runtime.Block.Retain
                                    local.get $this.1
                                    i32.const 0
                                    i32.add
                                    local.set $$t9.1
                                    local.get $$t9.0
                                    call $runtime.Block.Release
                                    local.set $$t9.0

                                    ;;*t8
                                    local.get $$t9.1
                                    i32.load offset=0 align=4
                                    call $runtime.Block.Retain
                                    local.get $$t9.1
                                    i32.load offset=4 align=4
                                    local.set $$t10.1
                                    local.get $$t10.0
                                    call $runtime.Block.Release
                                    local.set $$t10.0

                                    ;;t7 != t9
                                    local.get $$t8.1
                                    local.get $$t10.1
                                    i32.eq
                                    i32.eqz
                                    local.set $$t11

                                    ;;if t10 goto 5 else 7
                                    local.get $$t11
                                    if
                                      br $$Block_4
                                    else
                                      br $$Block_6
                                    end

                                  end ;;$Block_2
                                  i32.const 3
                                  local.set $$current_block

                                  ;;(*mapImp).successor(this, z)
                                  local.get $this.0
                                  local.get $this.1
                                  local.get $z.0
                                  local.get $z.1
                                  call $runtime.mapImp.successor
                                  local.set $$t5.1
                                  local.get $$t5.0
                                  call $runtime.Block.Release
                                  local.set $$t5.0

                                  ;;jump 2
                                  i32.const 2
                                  local.set $$block_selector
                                  br $$BlockDisp

                                end ;;$Block_3
                                i32.const 4
                                local.set $$current_block

                                ;;&z.Right [#3]
                                local.get $z.0
                                call $runtime.Block.Retain
                                local.get $z.1
                                i32.const 16
                                i32.add
                                local.set $$t12.1
                                local.get $$t12.0
                                call $runtime.Block.Release
                                local.set $$t12.0

                                ;;*t12
                                local.get $$t12.1
                                i32.load offset=0 align=4
                                call $runtime.Block.Retain
                                local.get $$t12.1
                                i32.load offset=4 align=4
                                local.set $$t13.1
                                local.get $$t13.0
                                call $runtime.Block.Release
                                local.set $$t13.0

                                ;;&this.NIL [#0]
                                local.get $this.0
                                call $runtime.Block.Retain
                                local.get $this.1
                                i32.const 0
                                i32.add
                                local.set $$t14.1
                                local.get $$t14.0
                                call $runtime.Block.Release
                                local.set $$t14.0

                                ;;*t14
                                local.get $$t14.1
                                i32.load offset=0 align=4
                                call $runtime.Block.Retain
                                local.get $$t14.1
                                i32.load offset=4 align=4
                                local.set $$t15.1
                                local.get $$t15.0
                                call $runtime.Block.Release
                                local.set $$t15.0

                                ;;t13 == t15
                                local.get $$t13.1
                                local.get $$t15.1
                                i32.eq
                                local.set $$t16

                                ;;if t16 goto 1 else 3
                                local.get $$t16
                                if
                                  i32.const 1
                                  local.set $$block_selector
                                  br $$BlockDisp
                                else
                                  i32.const 3
                                  local.set $$block_selector
                                  br $$BlockDisp
                                end

                              end ;;$Block_4
                              i32.const 5
                              local.set $$current_block

                              ;;&t5.Left [#2]
                              local.get $$t6.0
                              call $runtime.Block.Retain
                              local.get $$t6.1
                              i32.const 8
                              i32.add
                              local.set $$t17.1
                              local.get $$t17.0
                              call $runtime.Block.Release
                              local.set $$t17.0

                              ;;*t17
                              local.get $$t17.1
                              i32.load offset=0 align=4
                              call $runtime.Block.Retain
                              local.get $$t17.1
                              i32.load offset=4 align=4
                              local.set $$t18.1
                              local.get $$t18.0
                              call $runtime.Block.Release
                              local.set $$t18.0

                              ;;jump 6
                              br $$Block_5

                            end ;;$Block_5
                            local.get $$current_block
                            i32.const 5
                            i32.eq
                            if (result i32 i32)
                              local.get $$t18.0
                              call $runtime.Block.Retain
                              local.get $$t18.1
                            else
                              local.get $$t19.0
                              call $runtime.Block.Retain
                              local.get $$t19.1
                            end
                            local.set $$t20.1
                            local.get $$t20.0
                            call $runtime.Block.Release
                            local.set $$t20.0
                            i32.const 6
                            local.set $$current_block

                            ;;(*mapNode).Parent(t5, this)
                            local.get $$t6.0
                            local.get $$t6.1
                            local.get $this.0
                            local.get $this.1
                            call $runtime.mapNode.Parent
                            local.set $$t21.1
                            local.get $$t21.0
                            call $runtime.Block.Release
                            local.set $$t21.0

                            ;;(*mapNode).SetParent(t19, t20)
                            local.get $$t20.0
                            local.get $$t20.1
                            local.get $$t21.0
                            local.get $$t21.1
                            call $runtime.mapNode.SetParent

                            ;;(*mapNode).Parent(t5, this)
                            local.get $$t6.0
                            local.get $$t6.1
                            local.get $this.0
                            local.get $this.1
                            call $runtime.mapNode.Parent
                            local.set $$t22.1
                            local.get $$t22.0
                            call $runtime.Block.Release
                            local.set $$t22.0

                            ;;&this.NIL [#0]
                            local.get $this.0
                            call $runtime.Block.Retain
                            local.get $this.1
                            i32.const 0
                            i32.add
                            local.set $$t23.1
                            local.get $$t23.0
                            call $runtime.Block.Release
                            local.set $$t23.0

                            ;;*t23
                            local.get $$t23.1
                            i32.load offset=0 align=4
                            call $runtime.Block.Retain
                            local.get $$t23.1
                            i32.load offset=4 align=4
                            local.set $$t24.1
                            local.get $$t24.0
                            call $runtime.Block.Release
                            local.set $$t24.0

                            ;;t22 == t24
                            local.get $$t22.1
                            local.get $$t24.1
                            i32.eq
                            local.set $$t25

                            ;;if t25 goto 8 else 10
                            local.get $$t25
                            if
                              br $$Block_7
                            else
                              br $$Block_9
                            end

                          end ;;$Block_6
                          i32.const 7
                          local.set $$current_block

                          ;;&t5.Right [#3]
                          local.get $$t6.0
                          call $runtime.Block.Retain
                          local.get $$t6.1
                          i32.const 16
                          i32.add
                          local.set $$t26.1
                          local.get $$t26.0
                          call $runtime.Block.Release
                          local.set $$t26.0

                          ;;*t26
                          local.get $$t26.1
                          i32.load offset=0 align=4
                          call $runtime.Block.Retain
                          local.get $$t26.1
                          i32.load offset=4 align=4
                          local.set $$t19.1
                          local.get $$t19.0
                          call $runtime.Block.Release
                          local.set $$t19.0

                          ;;jump 6
                          i32.const 6
                          local.set $$block_selector
                          br $$BlockDisp

                        end ;;$Block_7
                        i32.const 8
                        local.set $$current_block

                        ;;&this.root [#1]
                        local.get $this.0
                        call $runtime.Block.Retain
                        local.get $this.1
                        i32.const 8
                        i32.add
                        local.set $$t27.1
                        local.get $$t27.0
                        call $runtime.Block.Release
                        local.set $$t27.0

                        ;;*t28 = t19
                        local.get $$t27.1
                        local.get $$t20.0
                        call $runtime.Block.Retain
                        local.get $$t27.1
                        i32.load offset=0 align=1
                        call $runtime.Block.Release
                        i32.store offset=0 align=1
                        local.get $$t27.1
                        local.get $$t20.1
                        i32.store offset=4 align=4

                        ;;jump 9
                        br $$Block_8

                      end ;;$Block_8
                      i32.const 9
                      local.set $$current_block

                      ;;t5 != z
                      local.get $$t6.1
                      local.get $z.1
                      i32.eq
                      i32.eqz
                      local.set $$t28

                      ;;if t29 goto 13 else 14
                      local.get $$t28
                      if
                        br $$Block_12
                      else
                        br $$Block_13
                      end

                    end ;;$Block_9
                    i32.const 10
                    local.set $$current_block

                    ;;(*mapNode).Parent(t5, this)
                    local.get $$t6.0
                    local.get $$t6.1
                    local.get $this.0
                    local.get $this.1
                    call $runtime.mapNode.Parent
                    local.set $$t29.1
                    local.get $$t29.0
                    call $runtime.Block.Release
                    local.set $$t29.0

                    ;;&t30.Left [#2]
                    local.get $$t29.0
                    call $runtime.Block.Retain
                    local.get $$t29.1
                    i32.const 8
                    i32.add
                    local.set $$t30.1
                    local.get $$t30.0
                    call $runtime.Block.Release
                    local.set $$t30.0

                    ;;*t31
                    local.get $$t30.1
                    i32.load offset=0 align=4
                    call $runtime.Block.Retain
                    local.get $$t30.1
                    i32.load offset=4 align=4
                    local.set $$t31.1
                    local.get $$t31.0
                    call $runtime.Block.Release
                    local.set $$t31.0

                    ;;t5 == t32
                    local.get $$t6.1
                    local.get $$t31.1
                    i32.eq
                    local.set $$t32

                    ;;if t33 goto 11 else 12
                    local.get $$t32
                    if
                      br $$Block_10
                    else
                      br $$Block_11
                    end

                  end ;;$Block_10
                  i32.const 11
                  local.set $$current_block

                  ;;(*mapNode).Parent(t5, this)
                  local.get $$t6.0
                  local.get $$t6.1
                  local.get $this.0
                  local.get $this.1
                  call $runtime.mapNode.Parent
                  local.set $$t33.1
                  local.get $$t33.0
                  call $runtime.Block.Release
                  local.set $$t33.0

                  ;;&t34.Left [#2]
                  local.get $$t33.0
                  call $runtime.Block.Retain
                  local.get $$t33.1
                  i32.const 8
                  i32.add
                  local.set $$t34.1
                  local.get $$t34.0
                  call $runtime.Block.Release
                  local.set $$t34.0

                  ;;*t35 = t19
                  local.get $$t34.1
                  local.get $$t20.0
                  call $runtime.Block.Retain
                  local.get $$t34.1
                  i32.load offset=0 align=1
                  call $runtime.Block.Release
                  i32.store offset=0 align=1
                  local.get $$t34.1
                  local.get $$t20.1
                  i32.store offset=4 align=4

                  ;;jump 9
                  i32.const 9
                  local.set $$block_selector
                  br $$BlockDisp

                end ;;$Block_11
                i32.const 12
                local.set $$current_block

                ;;(*mapNode).Parent(t5, this)
                local.get $$t6.0
                local.get $$t6.1
                local.get $this.0
                local.get $this.1
                call $runtime.mapNode.Parent
                local.set $$t35.1
                local.get $$t35.0
                call $runtime.Block.Release
                local.set $$t35.0

                ;;&t36.Right [#3]
                local.get $$t35.0
                call $runtime.Block.Retain
                local.get $$t35.1
                i32.const 16
                i32.add
                local.set $$t36.1
                local.get $$t36.0
                call $runtime.Block.Release
                local.set $$t36.0

                ;;*t37 = t19
                local.get $$t36.1
                local.get $$t20.0
                call $runtime.Block.Retain
                local.get $$t36.1
                i32.load offset=0 align=1
                call $runtime.Block.Release
                i32.store offset=0 align=1
                local.get $$t36.1
                local.get $$t20.1
                i32.store offset=4 align=4

                ;;jump 9
                i32.const 9
                local.set $$block_selector
                br $$BlockDisp

              end ;;$Block_12
              i32.const 13
              local.set $$current_block

              ;;jump 14
              br $$Block_13

            end ;;$Block_13
            i32.const 14
            local.set $$current_block

            ;;&t5.Color [#4]
            local.get $$t6.0
            call $runtime.Block.Retain
            local.get $$t6.1
            i32.const 24
            i32.add
            local.set $$t37.1
            local.get $$t37.0
            call $runtime.Block.Release
            local.set $$t37.0

            ;;*t38
            local.get $$t37.1
            i32.load offset=0 align=4
            local.set $$t38

            ;;t39 == 1:mapColor
            local.get $$t38
            i32.const 1
            i32.eq
            local.set $$t39

            ;;if t40 goto 15 else 16
            local.get $$t39
            if
              br $$Block_14
            else
              br $$Block_15
            end

          end ;;$Block_14
          i32.const 15
          local.set $$current_block

          ;;(*mapImp).deleteFixup(this, t19)
          local.get $this.0
          local.get $this.1
          local.get $$t20.0
          local.get $$t20.1
          call $runtime.mapImp.deleteFixup

          ;;jump 16
          br $$Block_15

        end ;;$Block_15
        i32.const 16
        local.set $$current_block

        ;;return 0:int
        i32.const 0
        local.set $$ret_0
        br $$BlockFnBody

      end ;;$Block_16
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$ret_0
  local.get $$t0.0
  call $runtime.Block.Release
  local.get $$t1.0
  call $runtime.Block.Release
  local.get $$t2.0
  call $runtime.Block.Release
  local.get $$t3.0
  call $runtime.Block.Release
  local.get $$t5.0
  call $runtime.Block.Release
  local.get $$t6.0
  call $runtime.Block.Release
  local.get $$t7.0
  call $runtime.Block.Release
  local.get $$t8.0
  call $runtime.Block.Release
  local.get $$t9.0
  call $runtime.Block.Release
  local.get $$t10.0
  call $runtime.Block.Release
  local.get $$t12.0
  call $runtime.Block.Release
  local.get $$t13.0
  call $runtime.Block.Release
  local.get $$t14.0
  call $runtime.Block.Release
  local.get $$t15.0
  call $runtime.Block.Release
  local.get $$t17.0
  call $runtime.Block.Release
  local.get $$t18.0
  call $runtime.Block.Release
  local.get $$t19.0
  call $runtime.Block.Release
  local.get $$t20.0
  call $runtime.Block.Release
  local.get $$t21.0
  call $runtime.Block.Release
  local.get $$t22.0
  call $runtime.Block.Release
  local.get $$t23.0
  call $runtime.Block.Release
  local.get $$t24.0
  call $runtime.Block.Release
  local.get $$t26.0
  call $runtime.Block.Release
  local.get $$t27.0
  call $runtime.Block.Release
  local.get $$t29.0
  call $runtime.Block.Release
  local.get $$t30.0
  call $runtime.Block.Release
  local.get $$t31.0
  call $runtime.Block.Release
  local.get $$t33.0
  call $runtime.Block.Release
  local.get $$t34.0
  call $runtime.Block.Release
  local.get $$t35.0
  call $runtime.Block.Release
  local.get $$t36.0
  call $runtime.Block.Release
  local.get $$t37.0
  call $runtime.Block.Release
) ;;runtime.mapImp.delete

(func $runtime.mapImp.deleteFixup (param $this.0 i32) (param $this.1 i32) (param $x.0 i32) (param $x.1 i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$t0.0 i32)
  (local $$t0.1 i32)
  (local $$t1.0 i32)
  (local $$t1.1 i32)
  (local $$t2.0 i32)
  (local $$t2.1 i32)
  (local $$t3.0 i32)
  (local $$t3.1 i32)
  (local $$t4 i32)
  (local $$t5.0 i32)
  (local $$t5.1 i32)
  (local $$t6.0 i32)
  (local $$t6.1 i32)
  (local $$t7.0 i32)
  (local $$t7.1 i32)
  (local $$t8.0 i32)
  (local $$t8.1 i32)
  (local $$t9.0 i32)
  (local $$t9.1 i32)
  (local $$t10.0 i32)
  (local $$t10.1 i32)
  (local $$t11.0 i32)
  (local $$t11.1 i32)
  (local $$t12 i32)
  (local $$t13.0 i32)
  (local $$t13.1 i32)
  (local $$t14 i32)
  (local $$t15 i32)
  (local $$t16.0 i32)
  (local $$t16.1 i32)
  (local $$t17.0 i32)
  (local $$t17.1 i32)
  (local $$t18.0 i32)
  (local $$t18.1 i32)
  (local $$t19.0 i32)
  (local $$t19.1 i32)
  (local $$t20 i32)
  (local $$t21 i32)
  (local $$t22.0 i32)
  (local $$t22.1 i32)
  (local $$t23.0 i32)
  (local $$t23.1 i32)
  (local $$t24.0 i32)
  (local $$t24.1 i32)
  (local $$t25.0 i32)
  (local $$t25.1 i32)
  (local $$t26 i32)
  (local $$t27 i32)
  (local $$t28.0 i32)
  (local $$t28.1 i32)
  (local $$t29.0 i32)
  (local $$t29.1 i32)
  (local $$t30.0 i32)
  (local $$t30.1 i32)
  (local $$t31.0 i32)
  (local $$t31.1 i32)
  (local $$t32.0 i32)
  (local $$t32.1 i32)
  (local $$t33.0 i32)
  (local $$t33.1 i32)
  (local $$t34.0 i32)
  (local $$t34.1 i32)
  (local $$t35.0 i32)
  (local $$t35.1 i32)
  (local $$t36.0 i32)
  (local $$t36.1 i32)
  (local $$t37.0 i32)
  (local $$t37.1 i32)
  (local $$t38.0 i32)
  (local $$t38.1 i32)
  (local $$t39 i32)
  (local $$t40 i32)
  (local $$t41.0 i32)
  (local $$t41.1 i32)
  (local $$t42.0 i32)
  (local $$t42.1 i32)
  (local $$t43.0 i32)
  (local $$t43.1 i32)
  (local $$t44.0 i32)
  (local $$t44.1 i32)
  (local $$t45 i32)
  (local $$t46 i32)
  (local $$t47.0 i32)
  (local $$t47.1 i32)
  (local $$t48.0 i32)
  (local $$t48.1 i32)
  (local $$t49.0 i32)
  (local $$t49.1 i32)
  (local $$t50 i32)
  (local $$t51 i32)
  (local $$t52.0 i32)
  (local $$t52.1 i32)
  (local $$t53.0 i32)
  (local $$t53.1 i32)
  (local $$t54.0 i32)
  (local $$t54.1 i32)
  (local $$t55.0 i32)
  (local $$t55.1 i32)
  (local $$t56.0 i32)
  (local $$t56.1 i32)
  (local $$t57.0 i32)
  (local $$t57.1 i32)
  (local $$t58.0 i32)
  (local $$t58.1 i32)
  (local $$t59.0 i32)
  (local $$t59.1 i32)
  (local $$t60.0 i32)
  (local $$t60.1 i32)
  (local $$t61.0 i32)
  (local $$t61.1 i32)
  (local $$t62.0 i32)
  (local $$t62.1 i32)
  (local $$t63 i32)
  (local $$t64.0 i32)
  (local $$t64.1 i32)
  (local $$t65.0 i32)
  (local $$t65.1 i32)
  (local $$t66.0 i32)
  (local $$t66.1 i32)
  (local $$t67.0 i32)
  (local $$t67.1 i32)
  (local $$t68.0 i32)
  (local $$t68.1 i32)
  (local $$t69.0 i32)
  (local $$t69.1 i32)
  (local $$t70.0 i32)
  (local $$t70.1 i32)
  (local $$t71.0 i32)
  (local $$t71.1 i32)
  (local $$t72.0 i32)
  (local $$t72.1 i32)
  (local $$t73.0 i32)
  (local $$t73.1 i32)
  (local $$t74.0 i32)
  (local $$t74.1 i32)
  (local $$t75.0 i32)
  (local $$t75.1 i32)
  (local $$t76.0 i32)
  (local $$t76.1 i32)
  (local $$t77.0 i32)
  (local $$t77.1 i32)
  (local $$t78.0 i32)
  (local $$t78.1 i32)
  (local $$t79.0 i32)
  (local $$t79.1 i32)
  (local $$t80.0 i32)
  (local $$t80.1 i32)
  (local $$t81.0 i32)
  (local $$t81.1 i32)
  (local $$t82 i32)
  (local $$t83 i32)
  (local $$t84.0 i32)
  (local $$t84.1 i32)
  (local $$t85.0 i32)
  (local $$t85.1 i32)
  (local $$t86.0 i32)
  (local $$t86.1 i32)
  (local $$t87.0 i32)
  (local $$t87.1 i32)
  (local $$t88 i32)
  (local $$t89 i32)
  (local $$t90.0 i32)
  (local $$t90.1 i32)
  (local $$t91.0 i32)
  (local $$t91.1 i32)
  (local $$t92.0 i32)
  (local $$t92.1 i32)
  (local $$t93 i32)
  (local $$t94 i32)
  (local $$t95.0 i32)
  (local $$t95.1 i32)
  (local $$t96.0 i32)
  (local $$t96.1 i32)
  (local $$t97.0 i32)
  (local $$t97.1 i32)
  (local $$t98.0 i32)
  (local $$t98.1 i32)
  (local $$t99.0 i32)
  (local $$t99.1 i32)
  (local $$t100.0 i32)
  (local $$t100.1 i32)
  (local $$t101.0 i32)
  (local $$t101.1 i32)
  (local $$t102.0 i32)
  (local $$t102.1 i32)
  (local $$t103.0 i32)
  (local $$t103.1 i32)
  (local $$t104.0 i32)
  (local $$t104.1 i32)
  (local $$t105.0 i32)
  (local $$t105.1 i32)
  (local $$t106 i32)
  (local $$t107.0 i32)
  (local $$t107.1 i32)
  (local $$t108.0 i32)
  (local $$t108.1 i32)
  (local $$t109.0 i32)
  (local $$t109.1 i32)
  (local $$t110.0 i32)
  (local $$t110.1 i32)
  (local $$t111.0 i32)
  (local $$t111.1 i32)
  (local $$t112.0 i32)
  (local $$t112.1 i32)
  (local $$t113.0 i32)
  (local $$t113.1 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_20
        block $$Block_19
          block $$Block_18
            block $$Block_17
              block $$Block_16
                block $$Block_15
                  block $$Block_14
                    block $$Block_13
                      block $$Block_12
                        block $$Block_11
                          block $$Block_10
                            block $$Block_9
                              block $$Block_8
                                block $$Block_7
                                  block $$Block_6
                                    block $$Block_5
                                      block $$Block_4
                                        block $$Block_3
                                          block $$Block_2
                                            block $$Block_1
                                              block $$Block_0
                                                block $$BlockSel
                                                  local.get $$block_selector
                                                  br_table 0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 0
                                                end ;;$BlockSel
                                                i32.const 0
                                                local.set $$current_block

                                                ;;jump 3
                                                br $$Block_2

                                              end ;;$Block_0
                                              i32.const 1
                                              local.set $$current_block

                                              ;;(*mapNode).Parent(t5, this)
                                              local.get $$t0.0
                                              local.get $$t0.1
                                              local.get $this.0
                                              local.get $this.1
                                              call $runtime.mapNode.Parent
                                              local.set $$t1.1
                                              local.get $$t1.0
                                              call $runtime.Block.Release
                                              local.set $$t1.0

                                              ;;&t0.Left [#2]
                                              local.get $$t1.0
                                              call $runtime.Block.Retain
                                              local.get $$t1.1
                                              i32.const 8
                                              i32.add
                                              local.set $$t2.1
                                              local.get $$t2.0
                                              call $runtime.Block.Release
                                              local.set $$t2.0

                                              ;;*t1
                                              local.get $$t2.1
                                              i32.load offset=0 align=4
                                              call $runtime.Block.Retain
                                              local.get $$t2.1
                                              i32.load offset=4 align=4
                                              local.set $$t3.1
                                              local.get $$t3.0
                                              call $runtime.Block.Release
                                              local.set $$t3.0

                                              ;;t5 == t2
                                              local.get $$t0.1
                                              local.get $$t3.1
                                              i32.eq
                                              local.set $$t4

                                              ;;if t3 goto 5 else 6
                                              local.get $$t4
                                              if
                                                br $$Block_4
                                              else
                                                br $$Block_5
                                              end

                                            end ;;$Block_1
                                            i32.const 2
                                            local.set $$current_block

                                            ;;&t5.Color [#4]
                                            local.get $$t0.0
                                            call $runtime.Block.Retain
                                            local.get $$t0.1
                                            i32.const 24
                                            i32.add
                                            local.set $$t5.1
                                            local.get $$t5.0
                                            call $runtime.Block.Release
                                            local.set $$t5.0

                                            ;;*t4 = 1:mapColor
                                            local.get $$t5.1
                                            i32.const 1
                                            i32.store offset=0 align=4

                                            ;;return
                                            br $$BlockFnBody

                                          end ;;$Block_2
                                          local.get $$current_block
                                          i32.const 0
                                          i32.eq
                                          if (result i32 i32)
                                            local.get $x.0
                                            call $runtime.Block.Retain
                                            local.get $x.1
                                          else
                                            local.get $$current_block
                                            i32.const 9
                                            i32.eq
                                            if (result i32 i32)
                                              local.get $$t6.0
                                              call $runtime.Block.Retain
                                              local.get $$t6.1
                                            else
                                              local.get $$current_block
                                              i32.const 16
                                              i32.eq
                                              if (result i32 i32)
                                                local.get $$t7.0
                                                call $runtime.Block.Retain
                                                local.get $$t7.1
                                              else
                                                local.get $$current_block
                                                i32.const 13
                                                i32.eq
                                                if (result i32 i32)
                                                  local.get $$t8.0
                                                  call $runtime.Block.Retain
                                                  local.get $$t8.1
                                                else
                                                  local.get $$t9.0
                                                  call $runtime.Block.Retain
                                                  local.get $$t9.1
                                                end
                                              end
                                            end
                                          end
                                          local.set $$t0.1
                                          local.get $$t0.0
                                          call $runtime.Block.Release
                                          local.set $$t0.0
                                          i32.const 3
                                          local.set $$current_block

                                          ;;&this.root [#1]
                                          local.get $this.0
                                          call $runtime.Block.Retain
                                          local.get $this.1
                                          i32.const 8
                                          i32.add
                                          local.set $$t10.1
                                          local.get $$t10.0
                                          call $runtime.Block.Release
                                          local.set $$t10.0

                                          ;;*t6
                                          local.get $$t10.1
                                          i32.load offset=0 align=4
                                          call $runtime.Block.Retain
                                          local.get $$t10.1
                                          i32.load offset=4 align=4
                                          local.set $$t11.1
                                          local.get $$t11.0
                                          call $runtime.Block.Release
                                          local.set $$t11.0

                                          ;;t5 != t7
                                          local.get $$t0.1
                                          local.get $$t11.1
                                          i32.eq
                                          i32.eqz
                                          local.set $$t12

                                          ;;if t8 goto 4 else 2
                                          local.get $$t12
                                          if
                                            br $$Block_3
                                          else
                                            i32.const 2
                                            local.set $$block_selector
                                            br $$BlockDisp
                                          end

                                        end ;;$Block_3
                                        i32.const 4
                                        local.set $$current_block

                                        ;;&t5.Color [#4]
                                        local.get $$t0.0
                                        call $runtime.Block.Retain
                                        local.get $$t0.1
                                        i32.const 24
                                        i32.add
                                        local.set $$t13.1
                                        local.get $$t13.0
                                        call $runtime.Block.Release
                                        local.set $$t13.0

                                        ;;*t9
                                        local.get $$t13.1
                                        i32.load offset=0 align=4
                                        local.set $$t14

                                        ;;t10 == 1:mapColor
                                        local.get $$t14
                                        i32.const 1
                                        i32.eq
                                        local.set $$t15

                                        ;;if t11 goto 1 else 2
                                        local.get $$t15
                                        if
                                          i32.const 1
                                          local.set $$block_selector
                                          br $$BlockDisp
                                        else
                                          i32.const 2
                                          local.set $$block_selector
                                          br $$BlockDisp
                                        end

                                      end ;;$Block_4
                                      i32.const 5
                                      local.set $$current_block

                                      ;;(*mapNode).Parent(t5, this)
                                      local.get $$t0.0
                                      local.get $$t0.1
                                      local.get $this.0
                                      local.get $this.1
                                      call $runtime.mapNode.Parent
                                      local.set $$t16.1
                                      local.get $$t16.0
                                      call $runtime.Block.Release
                                      local.set $$t16.0

                                      ;;&t12.Right [#3]
                                      local.get $$t16.0
                                      call $runtime.Block.Retain
                                      local.get $$t16.1
                                      i32.const 16
                                      i32.add
                                      local.set $$t17.1
                                      local.get $$t17.0
                                      call $runtime.Block.Release
                                      local.set $$t17.0

                                      ;;*t13
                                      local.get $$t17.1
                                      i32.load offset=0 align=4
                                      call $runtime.Block.Retain
                                      local.get $$t17.1
                                      i32.load offset=4 align=4
                                      local.set $$t18.1
                                      local.get $$t18.0
                                      call $runtime.Block.Release
                                      local.set $$t18.0

                                      ;;&t14.Color [#4]
                                      local.get $$t18.0
                                      call $runtime.Block.Retain
                                      local.get $$t18.1
                                      i32.const 24
                                      i32.add
                                      local.set $$t19.1
                                      local.get $$t19.0
                                      call $runtime.Block.Release
                                      local.set $$t19.0

                                      ;;*t15
                                      local.get $$t19.1
                                      i32.load offset=0 align=4
                                      local.set $$t20

                                      ;;t16 == 0:mapColor
                                      local.get $$t20
                                      i32.const 0
                                      i32.eq
                                      local.set $$t21

                                      ;;if t17 goto 7 else 8
                                      local.get $$t21
                                      if
                                        br $$Block_6
                                      else
                                        br $$Block_7
                                      end

                                    end ;;$Block_5
                                    i32.const 6
                                    local.set $$current_block

                                    ;;(*mapNode).Parent(t5, this)
                                    local.get $$t0.0
                                    local.get $$t0.1
                                    local.get $this.0
                                    local.get $this.1
                                    call $runtime.mapNode.Parent
                                    local.set $$t22.1
                                    local.get $$t22.0
                                    call $runtime.Block.Release
                                    local.set $$t22.0

                                    ;;&t18.Left [#2]
                                    local.get $$t22.0
                                    call $runtime.Block.Retain
                                    local.get $$t22.1
                                    i32.const 8
                                    i32.add
                                    local.set $$t23.1
                                    local.get $$t23.0
                                    call $runtime.Block.Release
                                    local.set $$t23.0

                                    ;;*t19
                                    local.get $$t23.1
                                    i32.load offset=0 align=4
                                    call $runtime.Block.Retain
                                    local.get $$t23.1
                                    i32.load offset=4 align=4
                                    local.set $$t24.1
                                    local.get $$t24.0
                                    call $runtime.Block.Release
                                    local.set $$t24.0

                                    ;;&t20.Color [#4]
                                    local.get $$t24.0
                                    call $runtime.Block.Retain
                                    local.get $$t24.1
                                    i32.const 24
                                    i32.add
                                    local.set $$t25.1
                                    local.get $$t25.0
                                    call $runtime.Block.Release
                                    local.set $$t25.0

                                    ;;*t21
                                    local.get $$t25.1
                                    i32.load offset=0 align=4
                                    local.set $$t26

                                    ;;t22 == 0:mapColor
                                    local.get $$t26
                                    i32.const 0
                                    i32.eq
                                    local.set $$t27

                                    ;;if t23 goto 14 else 15
                                    local.get $$t27
                                    if
                                      br $$Block_13
                                    else
                                      br $$Block_14
                                    end

                                  end ;;$Block_6
                                  i32.const 7
                                  local.set $$current_block

                                  ;;&t14.Color [#4]
                                  local.get $$t18.0
                                  call $runtime.Block.Retain
                                  local.get $$t18.1
                                  i32.const 24
                                  i32.add
                                  local.set $$t28.1
                                  local.get $$t28.0
                                  call $runtime.Block.Release
                                  local.set $$t28.0

                                  ;;*t24 = 1:mapColor
                                  local.get $$t28.1
                                  i32.const 1
                                  i32.store offset=0 align=4

                                  ;;(*mapNode).Parent(t5, this)
                                  local.get $$t0.0
                                  local.get $$t0.1
                                  local.get $this.0
                                  local.get $this.1
                                  call $runtime.mapNode.Parent
                                  local.set $$t29.1
                                  local.get $$t29.0
                                  call $runtime.Block.Release
                                  local.set $$t29.0

                                  ;;&t25.Color [#4]
                                  local.get $$t29.0
                                  call $runtime.Block.Retain
                                  local.get $$t29.1
                                  i32.const 24
                                  i32.add
                                  local.set $$t30.1
                                  local.get $$t30.0
                                  call $runtime.Block.Release
                                  local.set $$t30.0

                                  ;;*t26 = 0:mapColor
                                  local.get $$t30.1
                                  i32.const 0
                                  i32.store offset=0 align=4

                                  ;;(*mapNode).Parent(t5, this)
                                  local.get $$t0.0
                                  local.get $$t0.1
                                  local.get $this.0
                                  local.get $this.1
                                  call $runtime.mapNode.Parent
                                  local.set $$t31.1
                                  local.get $$t31.0
                                  call $runtime.Block.Release
                                  local.set $$t31.0

                                  ;;(*mapImp).leftRotate(this, t27)
                                  local.get $this.0
                                  local.get $this.1
                                  local.get $$t31.0
                                  local.get $$t31.1
                                  call $runtime.mapImp.leftRotate

                                  ;;(*mapNode).Parent(t5, this)
                                  local.get $$t0.0
                                  local.get $$t0.1
                                  local.get $this.0
                                  local.get $this.1
                                  call $runtime.mapNode.Parent
                                  local.set $$t32.1
                                  local.get $$t32.0
                                  call $runtime.Block.Release
                                  local.set $$t32.0

                                  ;;&t29.Right [#3]
                                  local.get $$t32.0
                                  call $runtime.Block.Retain
                                  local.get $$t32.1
                                  i32.const 16
                                  i32.add
                                  local.set $$t33.1
                                  local.get $$t33.0
                                  call $runtime.Block.Release
                                  local.set $$t33.0

                                  ;;*t30
                                  local.get $$t33.1
                                  i32.load offset=0 align=4
                                  call $runtime.Block.Retain
                                  local.get $$t33.1
                                  i32.load offset=4 align=4
                                  local.set $$t34.1
                                  local.get $$t34.0
                                  call $runtime.Block.Release
                                  local.set $$t34.0

                                  ;;jump 8
                                  br $$Block_7

                                end ;;$Block_7
                                local.get $$current_block
                                i32.const 5
                                i32.eq
                                if (result i32 i32)
                                  local.get $$t18.0
                                  call $runtime.Block.Retain
                                  local.get $$t18.1
                                else
                                  local.get $$t34.0
                                  call $runtime.Block.Retain
                                  local.get $$t34.1
                                end
                                local.set $$t35.1
                                local.get $$t35.0
                                call $runtime.Block.Release
                                local.set $$t35.0
                                i32.const 8
                                local.set $$current_block

                                ;;&t32.Left [#2]
                                local.get $$t35.0
                                call $runtime.Block.Retain
                                local.get $$t35.1
                                i32.const 8
                                i32.add
                                local.set $$t36.1
                                local.get $$t36.0
                                call $runtime.Block.Release
                                local.set $$t36.0

                                ;;*t33
                                local.get $$t36.1
                                i32.load offset=0 align=4
                                call $runtime.Block.Retain
                                local.get $$t36.1
                                i32.load offset=4 align=4
                                local.set $$t37.1
                                local.get $$t37.0
                                call $runtime.Block.Release
                                local.set $$t37.0

                                ;;&t34.Color [#4]
                                local.get $$t37.0
                                call $runtime.Block.Retain
                                local.get $$t37.1
                                i32.const 24
                                i32.add
                                local.set $$t38.1
                                local.get $$t38.0
                                call $runtime.Block.Release
                                local.set $$t38.0

                                ;;*t35
                                local.get $$t38.1
                                i32.load offset=0 align=4
                                local.set $$t39

                                ;;t36 == 1:mapColor
                                local.get $$t39
                                i32.const 1
                                i32.eq
                                local.set $$t40

                                ;;if t37 goto 11 else 10
                                local.get $$t40
                                if
                                  br $$Block_10
                                else
                                  br $$Block_9
                                end

                              end ;;$Block_8
                              i32.const 9
                              local.set $$current_block

                              ;;&t32.Color [#4]
                              local.get $$t35.0
                              call $runtime.Block.Retain
                              local.get $$t35.1
                              i32.const 24
                              i32.add
                              local.set $$t41.1
                              local.get $$t41.0
                              call $runtime.Block.Release
                              local.set $$t41.0

                              ;;*t38 = 0:mapColor
                              local.get $$t41.1
                              i32.const 0
                              i32.store offset=0 align=4

                              ;;(*mapNode).Parent(t5, this)
                              local.get $$t0.0
                              local.get $$t0.1
                              local.get $this.0
                              local.get $this.1
                              call $runtime.mapNode.Parent
                              local.set $$t6.1
                              local.get $$t6.0
                              call $runtime.Block.Release
                              local.set $$t6.0

                              ;;jump 3
                              i32.const 3
                              local.set $$block_selector
                              br $$BlockDisp

                            end ;;$Block_9
                            i32.const 10
                            local.set $$current_block

                            ;;&t32.Right [#3]
                            local.get $$t35.0
                            call $runtime.Block.Retain
                            local.get $$t35.1
                            i32.const 16
                            i32.add
                            local.set $$t42.1
                            local.get $$t42.0
                            call $runtime.Block.Release
                            local.set $$t42.0

                            ;;*t40
                            local.get $$t42.1
                            i32.load offset=0 align=4
                            call $runtime.Block.Retain
                            local.get $$t42.1
                            i32.load offset=4 align=4
                            local.set $$t43.1
                            local.get $$t43.0
                            call $runtime.Block.Release
                            local.set $$t43.0

                            ;;&t41.Color [#4]
                            local.get $$t43.0
                            call $runtime.Block.Retain
                            local.get $$t43.1
                            i32.const 24
                            i32.add
                            local.set $$t44.1
                            local.get $$t44.0
                            call $runtime.Block.Release
                            local.set $$t44.0

                            ;;*t42
                            local.get $$t44.1
                            i32.load offset=0 align=4
                            local.set $$t45

                            ;;t43 == 1:mapColor
                            local.get $$t45
                            i32.const 1
                            i32.eq
                            local.set $$t46

                            ;;if t44 goto 12 else 13
                            local.get $$t46
                            if
                              br $$Block_11
                            else
                              br $$Block_12
                            end

                          end ;;$Block_10
                          i32.const 11
                          local.set $$current_block

                          ;;&t32.Right [#3]
                          local.get $$t35.0
                          call $runtime.Block.Retain
                          local.get $$t35.1
                          i32.const 16
                          i32.add
                          local.set $$t47.1
                          local.get $$t47.0
                          call $runtime.Block.Release
                          local.set $$t47.0

                          ;;*t45
                          local.get $$t47.1
                          i32.load offset=0 align=4
                          call $runtime.Block.Retain
                          local.get $$t47.1
                          i32.load offset=4 align=4
                          local.set $$t48.1
                          local.get $$t48.0
                          call $runtime.Block.Release
                          local.set $$t48.0

                          ;;&t46.Color [#4]
                          local.get $$t48.0
                          call $runtime.Block.Retain
                          local.get $$t48.1
                          i32.const 24
                          i32.add
                          local.set $$t49.1
                          local.get $$t49.0
                          call $runtime.Block.Release
                          local.set $$t49.0

                          ;;*t47
                          local.get $$t49.1
                          i32.load offset=0 align=4
                          local.set $$t50

                          ;;t48 == 1:mapColor
                          local.get $$t50
                          i32.const 1
                          i32.eq
                          local.set $$t51

                          ;;if t49 goto 9 else 10
                          local.get $$t51
                          if
                            i32.const 9
                            local.set $$block_selector
                            br $$BlockDisp
                          else
                            i32.const 10
                            local.set $$block_selector
                            br $$BlockDisp
                          end

                        end ;;$Block_11
                        i32.const 12
                        local.set $$current_block

                        ;;&t32.Left [#2]
                        local.get $$t35.0
                        call $runtime.Block.Retain
                        local.get $$t35.1
                        i32.const 8
                        i32.add
                        local.set $$t52.1
                        local.get $$t52.0
                        call $runtime.Block.Release
                        local.set $$t52.0

                        ;;*t50
                        local.get $$t52.1
                        i32.load offset=0 align=4
                        call $runtime.Block.Retain
                        local.get $$t52.1
                        i32.load offset=4 align=4
                        local.set $$t53.1
                        local.get $$t53.0
                        call $runtime.Block.Release
                        local.set $$t53.0

                        ;;&t51.Color [#4]
                        local.get $$t53.0
                        call $runtime.Block.Retain
                        local.get $$t53.1
                        i32.const 24
                        i32.add
                        local.set $$t54.1
                        local.get $$t54.0
                        call $runtime.Block.Release
                        local.set $$t54.0

                        ;;*t52 = 1:mapColor
                        local.get $$t54.1
                        i32.const 1
                        i32.store offset=0 align=4

                        ;;&t32.Color [#4]
                        local.get $$t35.0
                        call $runtime.Block.Retain
                        local.get $$t35.1
                        i32.const 24
                        i32.add
                        local.set $$t55.1
                        local.get $$t55.0
                        call $runtime.Block.Release
                        local.set $$t55.0

                        ;;*t53 = 0:mapColor
                        local.get $$t55.1
                        i32.const 0
                        i32.store offset=0 align=4

                        ;;(*mapImp).rightRotate(this, t32)
                        local.get $this.0
                        local.get $this.1
                        local.get $$t35.0
                        local.get $$t35.1
                        call $runtime.mapImp.rightRotate

                        ;;(*mapNode).Parent(t5, this)
                        local.get $$t0.0
                        local.get $$t0.1
                        local.get $this.0
                        local.get $this.1
                        call $runtime.mapNode.Parent
                        local.set $$t56.1
                        local.get $$t56.0
                        call $runtime.Block.Release
                        local.set $$t56.0

                        ;;&t55.Right [#3]
                        local.get $$t56.0
                        call $runtime.Block.Retain
                        local.get $$t56.1
                        i32.const 16
                        i32.add
                        local.set $$t57.1
                        local.get $$t57.0
                        call $runtime.Block.Release
                        local.set $$t57.0

                        ;;*t56
                        local.get $$t57.1
                        i32.load offset=0 align=4
                        call $runtime.Block.Retain
                        local.get $$t57.1
                        i32.load offset=4 align=4
                        local.set $$t58.1
                        local.get $$t58.0
                        call $runtime.Block.Release
                        local.set $$t58.0

                        ;;jump 13
                        br $$Block_12

                      end ;;$Block_12
                      local.get $$current_block
                      i32.const 10
                      i32.eq
                      if (result i32 i32)
                        local.get $$t35.0
                        call $runtime.Block.Retain
                        local.get $$t35.1
                      else
                        local.get $$t58.0
                        call $runtime.Block.Retain
                        local.get $$t58.1
                      end
                      local.set $$t59.1
                      local.get $$t59.0
                      call $runtime.Block.Release
                      local.set $$t59.0
                      i32.const 13
                      local.set $$current_block

                      ;;&t58.Color [#4]
                      local.get $$t59.0
                      call $runtime.Block.Retain
                      local.get $$t59.1
                      i32.const 24
                      i32.add
                      local.set $$t60.1
                      local.get $$t60.0
                      call $runtime.Block.Release
                      local.set $$t60.0

                      ;;(*mapNode).Parent(t5, this)
                      local.get $$t0.0
                      local.get $$t0.1
                      local.get $this.0
                      local.get $this.1
                      call $runtime.mapNode.Parent
                      local.set $$t61.1
                      local.get $$t61.0
                      call $runtime.Block.Release
                      local.set $$t61.0

                      ;;&t60.Color [#4]
                      local.get $$t61.0
                      call $runtime.Block.Retain
                      local.get $$t61.1
                      i32.const 24
                      i32.add
                      local.set $$t62.1
                      local.get $$t62.0
                      call $runtime.Block.Release
                      local.set $$t62.0

                      ;;*t61
                      local.get $$t62.1
                      i32.load offset=0 align=4
                      local.set $$t63

                      ;;*t59 = t62
                      local.get $$t60.1
                      local.get $$t63
                      i32.store offset=0 align=4

                      ;;(*mapNode).Parent(t5, this)
                      local.get $$t0.0
                      local.get $$t0.1
                      local.get $this.0
                      local.get $this.1
                      call $runtime.mapNode.Parent
                      local.set $$t64.1
                      local.get $$t64.0
                      call $runtime.Block.Release
                      local.set $$t64.0

                      ;;&t63.Color [#4]
                      local.get $$t64.0
                      call $runtime.Block.Retain
                      local.get $$t64.1
                      i32.const 24
                      i32.add
                      local.set $$t65.1
                      local.get $$t65.0
                      call $runtime.Block.Release
                      local.set $$t65.0

                      ;;*t64 = 1:mapColor
                      local.get $$t65.1
                      i32.const 1
                      i32.store offset=0 align=4

                      ;;&t58.Right [#3]
                      local.get $$t59.0
                      call $runtime.Block.Retain
                      local.get $$t59.1
                      i32.const 16
                      i32.add
                      local.set $$t66.1
                      local.get $$t66.0
                      call $runtime.Block.Release
                      local.set $$t66.0

                      ;;*t65
                      local.get $$t66.1
                      i32.load offset=0 align=4
                      call $runtime.Block.Retain
                      local.get $$t66.1
                      i32.load offset=4 align=4
                      local.set $$t67.1
                      local.get $$t67.0
                      call $runtime.Block.Release
                      local.set $$t67.0

                      ;;&t66.Color [#4]
                      local.get $$t67.0
                      call $runtime.Block.Retain
                      local.get $$t67.1
                      i32.const 24
                      i32.add
                      local.set $$t68.1
                      local.get $$t68.0
                      call $runtime.Block.Release
                      local.set $$t68.0

                      ;;*t67 = 1:mapColor
                      local.get $$t68.1
                      i32.const 1
                      i32.store offset=0 align=4

                      ;;(*mapNode).Parent(t5, this)
                      local.get $$t0.0
                      local.get $$t0.1
                      local.get $this.0
                      local.get $this.1
                      call $runtime.mapNode.Parent
                      local.set $$t69.1
                      local.get $$t69.0
                      call $runtime.Block.Release
                      local.set $$t69.0

                      ;;(*mapImp).leftRotate(this, t68)
                      local.get $this.0
                      local.get $this.1
                      local.get $$t69.0
                      local.get $$t69.1
                      call $runtime.mapImp.leftRotate

                      ;;&this.root [#1]
                      local.get $this.0
                      call $runtime.Block.Retain
                      local.get $this.1
                      i32.const 8
                      i32.add
                      local.set $$t70.1
                      local.get $$t70.0
                      call $runtime.Block.Release
                      local.set $$t70.0

                      ;;*t70
                      local.get $$t70.1
                      i32.load offset=0 align=4
                      call $runtime.Block.Retain
                      local.get $$t70.1
                      i32.load offset=4 align=4
                      local.set $$t8.1
                      local.get $$t8.0
                      call $runtime.Block.Release
                      local.set $$t8.0

                      ;;jump 3
                      i32.const 3
                      local.set $$block_selector
                      br $$BlockDisp

                    end ;;$Block_13
                    i32.const 14
                    local.set $$current_block

                    ;;&t20.Color [#4]
                    local.get $$t24.0
                    call $runtime.Block.Retain
                    local.get $$t24.1
                    i32.const 24
                    i32.add
                    local.set $$t71.1
                    local.get $$t71.0
                    call $runtime.Block.Release
                    local.set $$t71.0

                    ;;*t72 = 1:mapColor
                    local.get $$t71.1
                    i32.const 1
                    i32.store offset=0 align=4

                    ;;(*mapNode).Parent(t5, this)
                    local.get $$t0.0
                    local.get $$t0.1
                    local.get $this.0
                    local.get $this.1
                    call $runtime.mapNode.Parent
                    local.set $$t72.1
                    local.get $$t72.0
                    call $runtime.Block.Release
                    local.set $$t72.0

                    ;;&t73.Color [#4]
                    local.get $$t72.0
                    call $runtime.Block.Retain
                    local.get $$t72.1
                    i32.const 24
                    i32.add
                    local.set $$t73.1
                    local.get $$t73.0
                    call $runtime.Block.Release
                    local.set $$t73.0

                    ;;*t74 = 0:mapColor
                    local.get $$t73.1
                    i32.const 0
                    i32.store offset=0 align=4

                    ;;(*mapNode).Parent(t5, this)
                    local.get $$t0.0
                    local.get $$t0.1
                    local.get $this.0
                    local.get $this.1
                    call $runtime.mapNode.Parent
                    local.set $$t74.1
                    local.get $$t74.0
                    call $runtime.Block.Release
                    local.set $$t74.0

                    ;;(*mapImp).rightRotate(this, t75)
                    local.get $this.0
                    local.get $this.1
                    local.get $$t74.0
                    local.get $$t74.1
                    call $runtime.mapImp.rightRotate

                    ;;(*mapNode).Parent(t5, this)
                    local.get $$t0.0
                    local.get $$t0.1
                    local.get $this.0
                    local.get $this.1
                    call $runtime.mapNode.Parent
                    local.set $$t75.1
                    local.get $$t75.0
                    call $runtime.Block.Release
                    local.set $$t75.0

                    ;;&t77.Left [#2]
                    local.get $$t75.0
                    call $runtime.Block.Retain
                    local.get $$t75.1
                    i32.const 8
                    i32.add
                    local.set $$t76.1
                    local.get $$t76.0
                    call $runtime.Block.Release
                    local.set $$t76.0

                    ;;*t78
                    local.get $$t76.1
                    i32.load offset=0 align=4
                    call $runtime.Block.Retain
                    local.get $$t76.1
                    i32.load offset=4 align=4
                    local.set $$t77.1
                    local.get $$t77.0
                    call $runtime.Block.Release
                    local.set $$t77.0

                    ;;jump 15
                    br $$Block_14

                  end ;;$Block_14
                  local.get $$current_block
                  i32.const 6
                  i32.eq
                  if (result i32 i32)
                    local.get $$t24.0
                    call $runtime.Block.Retain
                    local.get $$t24.1
                  else
                    local.get $$t77.0
                    call $runtime.Block.Retain
                    local.get $$t77.1
                  end
                  local.set $$t78.1
                  local.get $$t78.0
                  call $runtime.Block.Release
                  local.set $$t78.0
                  i32.const 15
                  local.set $$current_block

                  ;;&t80.Left [#2]
                  local.get $$t78.0
                  call $runtime.Block.Retain
                  local.get $$t78.1
                  i32.const 8
                  i32.add
                  local.set $$t79.1
                  local.get $$t79.0
                  call $runtime.Block.Release
                  local.set $$t79.0

                  ;;*t81
                  local.get $$t79.1
                  i32.load offset=0 align=4
                  call $runtime.Block.Retain
                  local.get $$t79.1
                  i32.load offset=4 align=4
                  local.set $$t80.1
                  local.get $$t80.0
                  call $runtime.Block.Release
                  local.set $$t80.0

                  ;;&t82.Color [#4]
                  local.get $$t80.0
                  call $runtime.Block.Retain
                  local.get $$t80.1
                  i32.const 24
                  i32.add
                  local.set $$t81.1
                  local.get $$t81.0
                  call $runtime.Block.Release
                  local.set $$t81.0

                  ;;*t83
                  local.get $$t81.1
                  i32.load offset=0 align=4
                  local.set $$t82

                  ;;t84 == 1:mapColor
                  local.get $$t82
                  i32.const 1
                  i32.eq
                  local.set $$t83

                  ;;if t85 goto 18 else 17
                  local.get $$t83
                  if
                    br $$Block_17
                  else
                    br $$Block_16
                  end

                end ;;$Block_15
                i32.const 16
                local.set $$current_block

                ;;&t80.Color [#4]
                local.get $$t78.0
                call $runtime.Block.Retain
                local.get $$t78.1
                i32.const 24
                i32.add
                local.set $$t84.1
                local.get $$t84.0
                call $runtime.Block.Release
                local.set $$t84.0

                ;;*t86 = 0:mapColor
                local.get $$t84.1
                i32.const 0
                i32.store offset=0 align=4

                ;;(*mapNode).Parent(t5, this)
                local.get $$t0.0
                local.get $$t0.1
                local.get $this.0
                local.get $this.1
                call $runtime.mapNode.Parent
                local.set $$t7.1
                local.get $$t7.0
                call $runtime.Block.Release
                local.set $$t7.0

                ;;jump 3
                i32.const 3
                local.set $$block_selector
                br $$BlockDisp

              end ;;$Block_16
              i32.const 17
              local.set $$current_block

              ;;&t80.Left [#2]
              local.get $$t78.0
              call $runtime.Block.Retain
              local.get $$t78.1
              i32.const 8
              i32.add
              local.set $$t85.1
              local.get $$t85.0
              call $runtime.Block.Release
              local.set $$t85.0

              ;;*t88
              local.get $$t85.1
              i32.load offset=0 align=4
              call $runtime.Block.Retain
              local.get $$t85.1
              i32.load offset=4 align=4
              local.set $$t86.1
              local.get $$t86.0
              call $runtime.Block.Release
              local.set $$t86.0

              ;;&t89.Color [#4]
              local.get $$t86.0
              call $runtime.Block.Retain
              local.get $$t86.1
              i32.const 24
              i32.add
              local.set $$t87.1
              local.get $$t87.0
              call $runtime.Block.Release
              local.set $$t87.0

              ;;*t90
              local.get $$t87.1
              i32.load offset=0 align=4
              local.set $$t88

              ;;t91 == 1:mapColor
              local.get $$t88
              i32.const 1
              i32.eq
              local.set $$t89

              ;;if t92 goto 19 else 20
              local.get $$t89
              if
                br $$Block_18
              else
                br $$Block_19
              end

            end ;;$Block_17
            i32.const 18
            local.set $$current_block

            ;;&t80.Right [#3]
            local.get $$t78.0
            call $runtime.Block.Retain
            local.get $$t78.1
            i32.const 16
            i32.add
            local.set $$t90.1
            local.get $$t90.0
            call $runtime.Block.Release
            local.set $$t90.0

            ;;*t93
            local.get $$t90.1
            i32.load offset=0 align=4
            call $runtime.Block.Retain
            local.get $$t90.1
            i32.load offset=4 align=4
            local.set $$t91.1
            local.get $$t91.0
            call $runtime.Block.Release
            local.set $$t91.0

            ;;&t94.Color [#4]
            local.get $$t91.0
            call $runtime.Block.Retain
            local.get $$t91.1
            i32.const 24
            i32.add
            local.set $$t92.1
            local.get $$t92.0
            call $runtime.Block.Release
            local.set $$t92.0

            ;;*t95
            local.get $$t92.1
            i32.load offset=0 align=4
            local.set $$t93

            ;;t96 == 1:mapColor
            local.get $$t93
            i32.const 1
            i32.eq
            local.set $$t94

            ;;if t97 goto 16 else 17
            local.get $$t94
            if
              i32.const 16
              local.set $$block_selector
              br $$BlockDisp
            else
              i32.const 17
              local.set $$block_selector
              br $$BlockDisp
            end

          end ;;$Block_18
          i32.const 19
          local.set $$current_block

          ;;&t80.Right [#3]
          local.get $$t78.0
          call $runtime.Block.Retain
          local.get $$t78.1
          i32.const 16
          i32.add
          local.set $$t95.1
          local.get $$t95.0
          call $runtime.Block.Release
          local.set $$t95.0

          ;;*t98
          local.get $$t95.1
          i32.load offset=0 align=4
          call $runtime.Block.Retain
          local.get $$t95.1
          i32.load offset=4 align=4
          local.set $$t96.1
          local.get $$t96.0
          call $runtime.Block.Release
          local.set $$t96.0

          ;;&t99.Color [#4]
          local.get $$t96.0
          call $runtime.Block.Retain
          local.get $$t96.1
          i32.const 24
          i32.add
          local.set $$t97.1
          local.get $$t97.0
          call $runtime.Block.Release
          local.set $$t97.0

          ;;*t100 = 1:mapColor
          local.get $$t97.1
          i32.const 1
          i32.store offset=0 align=4

          ;;&t80.Color [#4]
          local.get $$t78.0
          call $runtime.Block.Retain
          local.get $$t78.1
          i32.const 24
          i32.add
          local.set $$t98.1
          local.get $$t98.0
          call $runtime.Block.Release
          local.set $$t98.0

          ;;*t101 = 0:mapColor
          local.get $$t98.1
          i32.const 0
          i32.store offset=0 align=4

          ;;(*mapImp).leftRotate(this, t80)
          local.get $this.0
          local.get $this.1
          local.get $$t78.0
          local.get $$t78.1
          call $runtime.mapImp.leftRotate

          ;;(*mapNode).Parent(t5, this)
          local.get $$t0.0
          local.get $$t0.1
          local.get $this.0
          local.get $this.1
          call $runtime.mapNode.Parent
          local.set $$t99.1
          local.get $$t99.0
          call $runtime.Block.Release
          local.set $$t99.0

          ;;&t103.Left [#2]
          local.get $$t99.0
          call $runtime.Block.Retain
          local.get $$t99.1
          i32.const 8
          i32.add
          local.set $$t100.1
          local.get $$t100.0
          call $runtime.Block.Release
          local.set $$t100.0

          ;;*t104
          local.get $$t100.1
          i32.load offset=0 align=4
          call $runtime.Block.Retain
          local.get $$t100.1
          i32.load offset=4 align=4
          local.set $$t101.1
          local.get $$t101.0
          call $runtime.Block.Release
          local.set $$t101.0

          ;;jump 20
          br $$Block_19

        end ;;$Block_19
        local.get $$current_block
        i32.const 17
        i32.eq
        if (result i32 i32)
          local.get $$t78.0
          call $runtime.Block.Retain
          local.get $$t78.1
        else
          local.get $$t101.0
          call $runtime.Block.Retain
          local.get $$t101.1
        end
        local.set $$t102.1
        local.get $$t102.0
        call $runtime.Block.Release
        local.set $$t102.0
        i32.const 20
        local.set $$current_block

        ;;&t106.Color [#4]
        local.get $$t102.0
        call $runtime.Block.Retain
        local.get $$t102.1
        i32.const 24
        i32.add
        local.set $$t103.1
        local.get $$t103.0
        call $runtime.Block.Release
        local.set $$t103.0

        ;;(*mapNode).Parent(t5, this)
        local.get $$t0.0
        local.get $$t0.1
        local.get $this.0
        local.get $this.1
        call $runtime.mapNode.Parent
        local.set $$t104.1
        local.get $$t104.0
        call $runtime.Block.Release
        local.set $$t104.0

        ;;&t108.Color [#4]
        local.get $$t104.0
        call $runtime.Block.Retain
        local.get $$t104.1
        i32.const 24
        i32.add
        local.set $$t105.1
        local.get $$t105.0
        call $runtime.Block.Release
        local.set $$t105.0

        ;;*t109
        local.get $$t105.1
        i32.load offset=0 align=4
        local.set $$t106

        ;;*t107 = t110
        local.get $$t103.1
        local.get $$t106
        i32.store offset=0 align=4

        ;;(*mapNode).Parent(t5, this)
        local.get $$t0.0
        local.get $$t0.1
        local.get $this.0
        local.get $this.1
        call $runtime.mapNode.Parent
        local.set $$t107.1
        local.get $$t107.0
        call $runtime.Block.Release
        local.set $$t107.0

        ;;&t111.Color [#4]
        local.get $$t107.0
        call $runtime.Block.Retain
        local.get $$t107.1
        i32.const 24
        i32.add
        local.set $$t108.1
        local.get $$t108.0
        call $runtime.Block.Release
        local.set $$t108.0

        ;;*t112 = 1:mapColor
        local.get $$t108.1
        i32.const 1
        i32.store offset=0 align=4

        ;;&t106.Left [#2]
        local.get $$t102.0
        call $runtime.Block.Retain
        local.get $$t102.1
        i32.const 8
        i32.add
        local.set $$t109.1
        local.get $$t109.0
        call $runtime.Block.Release
        local.set $$t109.0

        ;;*t113
        local.get $$t109.1
        i32.load offset=0 align=4
        call $runtime.Block.Retain
        local.get $$t109.1
        i32.load offset=4 align=4
        local.set $$t110.1
        local.get $$t110.0
        call $runtime.Block.Release
        local.set $$t110.0

        ;;&t114.Color [#4]
        local.get $$t110.0
        call $runtime.Block.Retain
        local.get $$t110.1
        i32.const 24
        i32.add
        local.set $$t111.1
        local.get $$t111.0
        call $runtime.Block.Release
        local.set $$t111.0

        ;;*t115 = 1:mapColor
        local.get $$t111.1
        i32.const 1
        i32.store offset=0 align=4

        ;;(*mapNode).Parent(t5, this)
        local.get $$t0.0
        local.get $$t0.1
        local.get $this.0
        local.get $this.1
        call $runtime.mapNode.Parent
        local.set $$t112.1
        local.get $$t112.0
        call $runtime.Block.Release
        local.set $$t112.0

        ;;(*mapImp).rightRotate(this, t116)
        local.get $this.0
        local.get $this.1
        local.get $$t112.0
        local.get $$t112.1
        call $runtime.mapImp.rightRotate

        ;;&this.root [#1]
        local.get $this.0
        call $runtime.Block.Retain
        local.get $this.1
        i32.const 8
        i32.add
        local.set $$t113.1
        local.get $$t113.0
        call $runtime.Block.Release
        local.set $$t113.0

        ;;*t118
        local.get $$t113.1
        i32.load offset=0 align=4
        call $runtime.Block.Retain
        local.get $$t113.1
        i32.load offset=4 align=4
        local.set $$t9.1
        local.get $$t9.0
        call $runtime.Block.Release
        local.set $$t9.0

        ;;jump 3
        i32.const 3
        local.set $$block_selector
        br $$BlockDisp

      end ;;$Block_20
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$t0.0
  call $runtime.Block.Release
  local.get $$t1.0
  call $runtime.Block.Release
  local.get $$t2.0
  call $runtime.Block.Release
  local.get $$t3.0
  call $runtime.Block.Release
  local.get $$t5.0
  call $runtime.Block.Release
  local.get $$t6.0
  call $runtime.Block.Release
  local.get $$t7.0
  call $runtime.Block.Release
  local.get $$t8.0
  call $runtime.Block.Release
  local.get $$t9.0
  call $runtime.Block.Release
  local.get $$t10.0
  call $runtime.Block.Release
  local.get $$t11.0
  call $runtime.Block.Release
  local.get $$t13.0
  call $runtime.Block.Release
  local.get $$t16.0
  call $runtime.Block.Release
  local.get $$t17.0
  call $runtime.Block.Release
  local.get $$t18.0
  call $runtime.Block.Release
  local.get $$t19.0
  call $runtime.Block.Release
  local.get $$t22.0
  call $runtime.Block.Release
  local.get $$t23.0
  call $runtime.Block.Release
  local.get $$t24.0
  call $runtime.Block.Release
  local.get $$t25.0
  call $runtime.Block.Release
  local.get $$t28.0
  call $runtime.Block.Release
  local.get $$t29.0
  call $runtime.Block.Release
  local.get $$t30.0
  call $runtime.Block.Release
  local.get $$t31.0
  call $runtime.Block.Release
  local.get $$t32.0
  call $runtime.Block.Release
  local.get $$t33.0
  call $runtime.Block.Release
  local.get $$t34.0
  call $runtime.Block.Release
  local.get $$t35.0
  call $runtime.Block.Release
  local.get $$t36.0
  call $runtime.Block.Release
  local.get $$t37.0
  call $runtime.Block.Release
  local.get $$t38.0
  call $runtime.Block.Release
  local.get $$t41.0
  call $runtime.Block.Release
  local.get $$t42.0
  call $runtime.Block.Release
  local.get $$t43.0
  call $runtime.Block.Release
  local.get $$t44.0
  call $runtime.Block.Release
  local.get $$t47.0
  call $runtime.Block.Release
  local.get $$t48.0
  call $runtime.Block.Release
  local.get $$t49.0
  call $runtime.Block.Release
  local.get $$t52.0
  call $runtime.Block.Release
  local.get $$t53.0
  call $runtime.Block.Release
  local.get $$t54.0
  call $runtime.Block.Release
  local.get $$t55.0
  call $runtime.Block.Release
  local.get $$t56.0
  call $runtime.Block.Release
  local.get $$t57.0
  call $runtime.Block.Release
  local.get $$t58.0
  call $runtime.Block.Release
  local.get $$t59.0
  call $runtime.Block.Release
  local.get $$t60.0
  call $runtime.Block.Release
  local.get $$t61.0
  call $runtime.Block.Release
  local.get $$t62.0
  call $runtime.Block.Release
  local.get $$t64.0
  call $runtime.Block.Release
  local.get $$t65.0
  call $runtime.Block.Release
  local.get $$t66.0
  call $runtime.Block.Release
  local.get $$t67.0
  call $runtime.Block.Release
  local.get $$t68.0
  call $runtime.Block.Release
  local.get $$t69.0
  call $runtime.Block.Release
  local.get $$t70.0
  call $runtime.Block.Release
  local.get $$t71.0
  call $runtime.Block.Release
  local.get $$t72.0
  call $runtime.Block.Release
  local.get $$t73.0
  call $runtime.Block.Release
  local.get $$t74.0
  call $runtime.Block.Release
  local.get $$t75.0
  call $runtime.Block.Release
  local.get $$t76.0
  call $runtime.Block.Release
  local.get $$t77.0
  call $runtime.Block.Release
  local.get $$t78.0
  call $runtime.Block.Release
  local.get $$t79.0
  call $runtime.Block.Release
  local.get $$t80.0
  call $runtime.Block.Release
  local.get $$t81.0
  call $runtime.Block.Release
  local.get $$t84.0
  call $runtime.Block.Release
  local.get $$t85.0
  call $runtime.Block.Release
  local.get $$t86.0
  call $runtime.Block.Release
  local.get $$t87.0
  call $runtime.Block.Release
  local.get $$t90.0
  call $runtime.Block.Release
  local.get $$t91.0
  call $runtime.Block.Release
  local.get $$t92.0
  call $runtime.Block.Release
  local.get $$t95.0
  call $runtime.Block.Release
  local.get $$t96.0
  call $runtime.Block.Release
  local.get $$t97.0
  call $runtime.Block.Release
  local.get $$t98.0
  call $runtime.Block.Release
  local.get $$t99.0
  call $runtime.Block.Release
  local.get $$t100.0
  call $runtime.Block.Release
  local.get $$t101.0
  call $runtime.Block.Release
  local.get $$t102.0
  call $runtime.Block.Release
  local.get $$t103.0
  call $runtime.Block.Release
  local.get $$t104.0
  call $runtime.Block.Release
  local.get $$t105.0
  call $runtime.Block.Release
  local.get $$t107.0
  call $runtime.Block.Release
  local.get $$t108.0
  call $runtime.Block.Release
  local.get $$t109.0
  call $runtime.Block.Release
  local.get $$t110.0
  call $runtime.Block.Release
  local.get $$t111.0
  call $runtime.Block.Release
  local.get $$t112.0
  call $runtime.Block.Release
  local.get $$t113.0
  call $runtime.Block.Release
) ;;runtime.mapImp.deleteFixup

(func $runtime.mapImp.insert (param $this.0 i32) (param $this.1 i32) (param $z.0 i32) (param $z.1 i32) (result i32 i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$ret_0.0 i32)
  (local $$ret_0.1 i32)
  (local $$t0.0 i32)
  (local $$t0.1 i32)
  (local $$t1.0 i32)
  (local $$t1.1 i32)
  (local $$t2.0 i32)
  (local $$t2.1 i32)
  (local $$t3.0 i32)
  (local $$t3.1 i32)
  (local $$t4.0 i32)
  (local $$t4.1 i32)
  (local $$t5.0.0 i32)
  (local $$t5.0.1 i32)
  (local $$t5.1 i32)
  (local $$t5.2 i32)
  (local $$t6.0 i32)
  (local $$t6.1 i32)
  (local $$t7.0 i32)
  (local $$t7.1 i32)
  (local $$t8.0.0 i32)
  (local $$t8.0.1 i32)
  (local $$t8.1 i32)
  (local $$t8.2 i32)
  (local $$t9 i32)
  (local $$t10 i32)
  (local $$t11.0 i32)
  (local $$t11.1 i32)
  (local $$t12.0 i32)
  (local $$t12.1 i32)
  (local $$t13.0 i32)
  (local $$t13.1 i32)
  (local $$t14 i32)
  (local $$t15.0 i32)
  (local $$t15.1 i32)
  (local $$t16.0 i32)
  (local $$t16.1 i32)
  (local $$t17.0 i32)
  (local $$t17.1 i32)
  (local $$t18.0 i32)
  (local $$t18.1 i32)
  (local $$t19 i32)
  (local $$t20.0 i32)
  (local $$t20.1 i32)
  (local $$t21.0 i32)
  (local $$t21.1 i32)
  (local $$t22.0.0 i32)
  (local $$t22.0.1 i32)
  (local $$t22.1 i32)
  (local $$t22.2 i32)
  (local $$t23.0 i32)
  (local $$t23.1 i32)
  (local $$t24.0.0 i32)
  (local $$t24.0.1 i32)
  (local $$t24.1 i32)
  (local $$t24.2 i32)
  (local $$t25 i32)
  (local $$t26 i32)
  (local $$t27.0 i32)
  (local $$t27.1 i32)
  (local $$t28.0 i32)
  (local $$t28.1 i32)
  (local $$t29.0 i32)
  (local $$t29.1 i32)
  (local $$t30.0.0 i32)
  (local $$t30.0.1 i32)
  (local $$t30.1 i32)
  (local $$t30.2 i32)
  (local $$t31.0 i32)
  (local $$t31.1 i32)
  (local $$t32.0.0 i32)
  (local $$t32.0.1 i32)
  (local $$t32.1 i32)
  (local $$t32.2 i32)
  (local $$t33 i32)
  (local $$t34 i32)
  (local $$t35.0 i32)
  (local $$t35.1 i32)
  (local $$t36.0 i32)
  (local $$t36.1 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_12
        block $$Block_11
          block $$Block_10
            block $$Block_9
              block $$Block_8
                block $$Block_7
                  block $$Block_6
                    block $$Block_5
                      block $$Block_4
                        block $$Block_3
                          block $$Block_2
                            block $$Block_1
                              block $$Block_0
                                block $$BlockSel
                                  local.get $$block_selector
                                  br_table 0 1 2 3 4 5 6 7 8 9 10 11 12 0
                                end ;;$BlockSel
                                i32.const 0
                                local.set $$current_block

                                ;;&this.root [#1]
                                local.get $this.0
                                call $runtime.Block.Retain
                                local.get $this.1
                                i32.const 8
                                i32.add
                                local.set $$t0.1
                                local.get $$t0.0
                                call $runtime.Block.Release
                                local.set $$t0.0

                                ;;*t0
                                local.get $$t0.1
                                i32.load offset=0 align=4
                                call $runtime.Block.Retain
                                local.get $$t0.1
                                i32.load offset=4 align=4
                                local.set $$t1.1
                                local.get $$t1.0
                                call $runtime.Block.Release
                                local.set $$t1.0

                                ;;&this.NIL [#0]
                                local.get $this.0
                                call $runtime.Block.Retain
                                local.get $this.1
                                i32.const 0
                                i32.add
                                local.set $$t2.1
                                local.get $$t2.0
                                call $runtime.Block.Release
                                local.set $$t2.0

                                ;;*t2
                                local.get $$t2.1
                                i32.load offset=0 align=4
                                call $runtime.Block.Retain
                                local.get $$t2.1
                                i32.load offset=4 align=4
                                local.set $$t3.1
                                local.get $$t3.0
                                call $runtime.Block.Release
                                local.set $$t3.0

                                ;;jump 3
                                br $$Block_2

                              end ;;$Block_0
                              i32.const 1
                              local.set $$current_block

                              ;;&z.Key [#5]
                              local.get $z.0
                              call $runtime.Block.Retain
                              local.get $z.1
                              i32.const 28
                              i32.add
                              local.set $$t4.1
                              local.get $$t4.0
                              call $runtime.Block.Release
                              local.set $$t4.0

                              ;;*t4
                              local.get $$t4.1
                              i32.load offset=0 align=4
                              call $runtime.Block.Retain
                              local.get $$t4.1
                              i32.load offset=4 align=4
                              local.get $$t4.1
                              i32.load offset=8 align=4
                              local.get $$t4.1
                              i32.load offset=12 align=4
                              local.set $$t5.2
                              local.set $$t5.1
                              local.set $$t5.0.1
                              local.get $$t5.0.0
                              call $runtime.Block.Release
                              local.set $$t5.0.0

                              ;;&t14.Key [#5]
                              local.get $$t6.0
                              call $runtime.Block.Retain
                              local.get $$t6.1
                              i32.const 28
                              i32.add
                              local.set $$t7.1
                              local.get $$t7.0
                              call $runtime.Block.Release
                              local.set $$t7.0

                              ;;*t6
                              local.get $$t7.1
                              i32.load offset=0 align=4
                              call $runtime.Block.Retain
                              local.get $$t7.1
                              i32.load offset=4 align=4
                              local.get $$t7.1
                              i32.load offset=8 align=4
                              local.get $$t7.1
                              i32.load offset=12 align=4
                              local.set $$t8.2
                              local.set $$t8.1
                              local.set $$t8.0.1
                              local.get $$t8.0.0
                              call $runtime.Block.Release
                              local.set $$t8.0.0

                              ;;Compare(t5, t7)
                              local.get $$t5.0.0
                              local.get $$t5.0.1
                              local.get $$t5.1
                              local.get $$t5.2
                              local.get $$t8.0.0
                              local.get $$t8.0.1
                              local.get $$t8.1
                              local.get $$t8.2
                              call $runtime.Compare
                              local.set $$t9

                              ;;t8 < 0:i32
                              local.get $$t9
                              i32.const 0
                              i32.lt_s
                              local.set $$t10

                              ;;if t9 goto 4 else 5
                              local.get $$t10
                              if
                                br $$Block_3
                              else
                                br $$Block_4
                              end

                            end ;;$Block_1
                            i32.const 2
                            local.set $$current_block

                            ;;(*mapNode).SetParent(z, t15)
                            local.get $z.0
                            local.get $z.1
                            local.get $$t11.0
                            local.get $$t11.1
                            call $runtime.mapNode.SetParent

                            ;;&this.NIL [#0]
                            local.get $this.0
                            call $runtime.Block.Retain
                            local.get $this.1
                            i32.const 0
                            i32.add
                            local.set $$t12.1
                            local.get $$t12.0
                            call $runtime.Block.Release
                            local.set $$t12.0

                            ;;*t11
                            local.get $$t12.1
                            i32.load offset=0 align=4
                            call $runtime.Block.Retain
                            local.get $$t12.1
                            i32.load offset=4 align=4
                            local.set $$t13.1
                            local.get $$t13.0
                            call $runtime.Block.Release
                            local.set $$t13.0

                            ;;t15 == t12
                            local.get $$t11.1
                            local.get $$t13.1
                            i32.eq
                            local.set $$t14

                            ;;if t13 goto 8 else 10
                            local.get $$t14
                            if
                              br $$Block_7
                            else
                              br $$Block_9
                            end

                          end ;;$Block_2
                          local.get $$current_block
                          i32.const 0
                          i32.eq
                          if (result i32 i32)
                            local.get $$t1.0
                            call $runtime.Block.Retain
                            local.get $$t1.1
                          else
                            local.get $$current_block
                            i32.const 4
                            i32.eq
                            if (result i32 i32)
                              local.get $$t15.0
                              call $runtime.Block.Retain
                              local.get $$t15.1
                            else
                              local.get $$t16.0
                              call $runtime.Block.Retain
                              local.get $$t16.1
                            end
                          end
                          local.get $$current_block
                          i32.const 0
                          i32.eq
                          if (result i32 i32)
                            local.get $$t3.0
                            call $runtime.Block.Retain
                            local.get $$t3.1
                          else
                            local.get $$current_block
                            i32.const 4
                            i32.eq
                            if (result i32 i32)
                              local.get $$t6.0
                              call $runtime.Block.Retain
                              local.get $$t6.1
                            else
                              local.get $$t6.0
                              call $runtime.Block.Retain
                              local.get $$t6.1
                            end
                          end
                          local.set $$t11.1
                          local.get $$t11.0
                          call $runtime.Block.Release
                          local.set $$t11.0
                          local.set $$t6.1
                          local.get $$t6.0
                          call $runtime.Block.Release
                          local.set $$t6.0
                          i32.const 3
                          local.set $$current_block

                          ;;&this.NIL [#0]
                          local.get $this.0
                          call $runtime.Block.Retain
                          local.get $this.1
                          i32.const 0
                          i32.add
                          local.set $$t17.1
                          local.get $$t17.0
                          call $runtime.Block.Release
                          local.set $$t17.0

                          ;;*t16
                          local.get $$t17.1
                          i32.load offset=0 align=4
                          call $runtime.Block.Retain
                          local.get $$t17.1
                          i32.load offset=4 align=4
                          local.set $$t18.1
                          local.get $$t18.0
                          call $runtime.Block.Release
                          local.set $$t18.0

                          ;;t14 != t17
                          local.get $$t6.1
                          local.get $$t18.1
                          i32.eq
                          i32.eqz
                          local.set $$t19

                          ;;if t18 goto 1 else 2
                          local.get $$t19
                          if
                            i32.const 1
                            local.set $$block_selector
                            br $$BlockDisp
                          else
                            i32.const 2
                            local.set $$block_selector
                            br $$BlockDisp
                          end

                        end ;;$Block_3
                        i32.const 4
                        local.set $$current_block

                        ;;&t14.Left [#2]
                        local.get $$t6.0
                        call $runtime.Block.Retain
                        local.get $$t6.1
                        i32.const 8
                        i32.add
                        local.set $$t20.1
                        local.get $$t20.0
                        call $runtime.Block.Release
                        local.set $$t20.0

                        ;;*t19
                        local.get $$t20.1
                        i32.load offset=0 align=4
                        call $runtime.Block.Retain
                        local.get $$t20.1
                        i32.load offset=4 align=4
                        local.set $$t15.1
                        local.get $$t15.0
                        call $runtime.Block.Release
                        local.set $$t15.0

                        ;;jump 3
                        i32.const 3
                        local.set $$block_selector
                        br $$BlockDisp

                      end ;;$Block_4
                      i32.const 5
                      local.set $$current_block

                      ;;&t14.Key [#5]
                      local.get $$t6.0
                      call $runtime.Block.Retain
                      local.get $$t6.1
                      i32.const 28
                      i32.add
                      local.set $$t21.1
                      local.get $$t21.0
                      call $runtime.Block.Release
                      local.set $$t21.0

                      ;;*t21
                      local.get $$t21.1
                      i32.load offset=0 align=4
                      call $runtime.Block.Retain
                      local.get $$t21.1
                      i32.load offset=4 align=4
                      local.get $$t21.1
                      i32.load offset=8 align=4
                      local.get $$t21.1
                      i32.load offset=12 align=4
                      local.set $$t22.2
                      local.set $$t22.1
                      local.set $$t22.0.1
                      local.get $$t22.0.0
                      call $runtime.Block.Release
                      local.set $$t22.0.0

                      ;;&z.Key [#5]
                      local.get $z.0
                      call $runtime.Block.Retain
                      local.get $z.1
                      i32.const 28
                      i32.add
                      local.set $$t23.1
                      local.get $$t23.0
                      call $runtime.Block.Release
                      local.set $$t23.0

                      ;;*t23
                      local.get $$t23.1
                      i32.load offset=0 align=4
                      call $runtime.Block.Retain
                      local.get $$t23.1
                      i32.load offset=4 align=4
                      local.get $$t23.1
                      i32.load offset=8 align=4
                      local.get $$t23.1
                      i32.load offset=12 align=4
                      local.set $$t24.2
                      local.set $$t24.1
                      local.set $$t24.0.1
                      local.get $$t24.0.0
                      call $runtime.Block.Release
                      local.set $$t24.0.0

                      ;;Compare(t22, t24)
                      local.get $$t22.0.0
                      local.get $$t22.0.1
                      local.get $$t22.1
                      local.get $$t22.2
                      local.get $$t24.0.0
                      local.get $$t24.0.1
                      local.get $$t24.1
                      local.get $$t24.2
                      call $runtime.Compare
                      local.set $$t25

                      ;;t25 < 0:i32
                      local.get $$t25
                      i32.const 0
                      i32.lt_s
                      local.set $$t26

                      ;;if t26 goto 6 else 7
                      local.get $$t26
                      if
                        br $$Block_5
                      else
                        br $$Block_6
                      end

                    end ;;$Block_5
                    i32.const 6
                    local.set $$current_block

                    ;;&t14.Right [#3]
                    local.get $$t6.0
                    call $runtime.Block.Retain
                    local.get $$t6.1
                    i32.const 16
                    i32.add
                    local.set $$t27.1
                    local.get $$t27.0
                    call $runtime.Block.Release
                    local.set $$t27.0

                    ;;*t27
                    local.get $$t27.1
                    i32.load offset=0 align=4
                    call $runtime.Block.Retain
                    local.get $$t27.1
                    i32.load offset=4 align=4
                    local.set $$t16.1
                    local.get $$t16.0
                    call $runtime.Block.Release
                    local.set $$t16.0

                    ;;jump 3
                    i32.const 3
                    local.set $$block_selector
                    br $$BlockDisp

                  end ;;$Block_6
                  i32.const 7
                  local.set $$current_block

                  ;;return t14
                  local.get $$t6.0
                  call $runtime.Block.Retain
                  local.get $$t6.1
                  local.set $$ret_0.1
                  local.get $$ret_0.0
                  call $runtime.Block.Release
                  local.set $$ret_0.0
                  br $$BlockFnBody

                end ;;$Block_7
                i32.const 8
                local.set $$current_block

                ;;&this.root [#1]
                local.get $this.0
                call $runtime.Block.Retain
                local.get $this.1
                i32.const 8
                i32.add
                local.set $$t28.1
                local.get $$t28.0
                call $runtime.Block.Release
                local.set $$t28.0

                ;;*t29 = z
                local.get $$t28.1
                local.get $z.0
                call $runtime.Block.Retain
                local.get $$t28.1
                i32.load offset=0 align=1
                call $runtime.Block.Release
                i32.store offset=0 align=1
                local.get $$t28.1
                local.get $z.1
                i32.store offset=4 align=4

                ;;jump 9
                br $$Block_8

              end ;;$Block_8
              i32.const 9
              local.set $$current_block

              ;;(*mapImp).insertFixup(this, z)
              local.get $this.0
              local.get $this.1
              local.get $z.0
              local.get $z.1
              call $runtime.mapImp.insertFixup

              ;;return z
              local.get $z.0
              call $runtime.Block.Retain
              local.get $z.1
              local.set $$ret_0.1
              local.get $$ret_0.0
              call $runtime.Block.Release
              local.set $$ret_0.0
              br $$BlockFnBody

            end ;;$Block_9
            i32.const 10
            local.set $$current_block

            ;;&z.Key [#5]
            local.get $z.0
            call $runtime.Block.Retain
            local.get $z.1
            i32.const 28
            i32.add
            local.set $$t29.1
            local.get $$t29.0
            call $runtime.Block.Release
            local.set $$t29.0

            ;;*t31
            local.get $$t29.1
            i32.load offset=0 align=4
            call $runtime.Block.Retain
            local.get $$t29.1
            i32.load offset=4 align=4
            local.get $$t29.1
            i32.load offset=8 align=4
            local.get $$t29.1
            i32.load offset=12 align=4
            local.set $$t30.2
            local.set $$t30.1
            local.set $$t30.0.1
            local.get $$t30.0.0
            call $runtime.Block.Release
            local.set $$t30.0.0

            ;;&t15.Key [#5]
            local.get $$t11.0
            call $runtime.Block.Retain
            local.get $$t11.1
            i32.const 28
            i32.add
            local.set $$t31.1
            local.get $$t31.0
            call $runtime.Block.Release
            local.set $$t31.0

            ;;*t33
            local.get $$t31.1
            i32.load offset=0 align=4
            call $runtime.Block.Retain
            local.get $$t31.1
            i32.load offset=4 align=4
            local.get $$t31.1
            i32.load offset=8 align=4
            local.get $$t31.1
            i32.load offset=12 align=4
            local.set $$t32.2
            local.set $$t32.1
            local.set $$t32.0.1
            local.get $$t32.0.0
            call $runtime.Block.Release
            local.set $$t32.0.0

            ;;Compare(t32, t34)
            local.get $$t30.0.0
            local.get $$t30.0.1
            local.get $$t30.1
            local.get $$t30.2
            local.get $$t32.0.0
            local.get $$t32.0.1
            local.get $$t32.1
            local.get $$t32.2
            call $runtime.Compare
            local.set $$t33

            ;;t35 < 0:i32
            local.get $$t33
            i32.const 0
            i32.lt_s
            local.set $$t34

            ;;if t36 goto 11 else 12
            local.get $$t34
            if
              br $$Block_10
            else
              br $$Block_11
            end

          end ;;$Block_10
          i32.const 11
          local.set $$current_block

          ;;&t15.Left [#2]
          local.get $$t11.0
          call $runtime.Block.Retain
          local.get $$t11.1
          i32.const 8
          i32.add
          local.set $$t35.1
          local.get $$t35.0
          call $runtime.Block.Release
          local.set $$t35.0

          ;;*t37 = z
          local.get $$t35.1
          local.get $z.0
          call $runtime.Block.Retain
          local.get $$t35.1
          i32.load offset=0 align=1
          call $runtime.Block.Release
          i32.store offset=0 align=1
          local.get $$t35.1
          local.get $z.1
          i32.store offset=4 align=4

          ;;jump 9
          i32.const 9
          local.set $$block_selector
          br $$BlockDisp

        end ;;$Block_11
        i32.const 12
        local.set $$current_block

        ;;&t15.Right [#3]
        local.get $$t11.0
        call $runtime.Block.Retain
        local.get $$t11.1
        i32.const 16
        i32.add
        local.set $$t36.1
        local.get $$t36.0
        call $runtime.Block.Release
        local.set $$t36.0

        ;;*t38 = z
        local.get $$t36.1
        local.get $z.0
        call $runtime.Block.Retain
        local.get $$t36.1
        i32.load offset=0 align=1
        call $runtime.Block.Release
        i32.store offset=0 align=1
        local.get $$t36.1
        local.get $z.1
        i32.store offset=4 align=4

        ;;jump 9
        i32.const 9
        local.set $$block_selector
        br $$BlockDisp

      end ;;$Block_12
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$ret_0.0
  call $runtime.Block.Retain
  local.get $$ret_0.1
  local.get $$ret_0.0
  call $runtime.Block.Release
  local.get $$t0.0
  call $runtime.Block.Release
  local.get $$t1.0
  call $runtime.Block.Release
  local.get $$t2.0
  call $runtime.Block.Release
  local.get $$t3.0
  call $runtime.Block.Release
  local.get $$t4.0
  call $runtime.Block.Release
  local.get $$t5.0.0
  call $runtime.Block.Release
  local.get $$t6.0
  call $runtime.Block.Release
  local.get $$t7.0
  call $runtime.Block.Release
  local.get $$t8.0.0
  call $runtime.Block.Release
  local.get $$t11.0
  call $runtime.Block.Release
  local.get $$t12.0
  call $runtime.Block.Release
  local.get $$t13.0
  call $runtime.Block.Release
  local.get $$t15.0
  call $runtime.Block.Release
  local.get $$t16.0
  call $runtime.Block.Release
  local.get $$t17.0
  call $runtime.Block.Release
  local.get $$t18.0
  call $runtime.Block.Release
  local.get $$t20.0
  call $runtime.Block.Release
  local.get $$t21.0
  call $runtime.Block.Release
  local.get $$t22.0.0
  call $runtime.Block.Release
  local.get $$t23.0
  call $runtime.Block.Release
  local.get $$t24.0.0
  call $runtime.Block.Release
  local.get $$t27.0
  call $runtime.Block.Release
  local.get $$t28.0
  call $runtime.Block.Release
  local.get $$t29.0
  call $runtime.Block.Release
  local.get $$t30.0.0
  call $runtime.Block.Release
  local.get $$t31.0
  call $runtime.Block.Release
  local.get $$t32.0.0
  call $runtime.Block.Release
  local.get $$t35.0
  call $runtime.Block.Release
  local.get $$t36.0
  call $runtime.Block.Release
) ;;runtime.mapImp.insert

(func $runtime.mapImp.insertFixup (param $this.0 i32) (param $this.1 i32) (param $z.0 i32) (param $z.1 i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$t0.0 i32)
  (local $$t0.1 i32)
  (local $$t1.0 i32)
  (local $$t1.1 i32)
  (local $$t2.0 i32)
  (local $$t2.1 i32)
  (local $$t3.0 i32)
  (local $$t3.1 i32)
  (local $$t4.0 i32)
  (local $$t4.1 i32)
  (local $$t5.0 i32)
  (local $$t5.1 i32)
  (local $$t6 i32)
  (local $$t7.0 i32)
  (local $$t7.1 i32)
  (local $$t8.0 i32)
  (local $$t8.1 i32)
  (local $$t9.0 i32)
  (local $$t9.1 i32)
  (local $$t10.0 i32)
  (local $$t10.1 i32)
  (local $$t11.0 i32)
  (local $$t11.1 i32)
  (local $$t12.0 i32)
  (local $$t12.1 i32)
  (local $$t13.0 i32)
  (local $$t13.1 i32)
  (local $$t14.0 i32)
  (local $$t14.1 i32)
  (local $$t15.0 i32)
  (local $$t15.1 i32)
  (local $$t16 i32)
  (local $$t17 i32)
  (local $$t18.0 i32)
  (local $$t18.1 i32)
  (local $$t19.0 i32)
  (local $$t19.1 i32)
  (local $$t20.0 i32)
  (local $$t20.1 i32)
  (local $$t21.0 i32)
  (local $$t21.1 i32)
  (local $$t22.0 i32)
  (local $$t22.1 i32)
  (local $$t23 i32)
  (local $$t24 i32)
  (local $$t25.0 i32)
  (local $$t25.1 i32)
  (local $$t26.0 i32)
  (local $$t26.1 i32)
  (local $$t27.0 i32)
  (local $$t27.1 i32)
  (local $$t28.0 i32)
  (local $$t28.1 i32)
  (local $$t29.0 i32)
  (local $$t29.1 i32)
  (local $$t30 i32)
  (local $$t31 i32)
  (local $$t32.0 i32)
  (local $$t32.1 i32)
  (local $$t33.0 i32)
  (local $$t33.1 i32)
  (local $$t34.0 i32)
  (local $$t34.1 i32)
  (local $$t35.0 i32)
  (local $$t35.1 i32)
  (local $$t36.0 i32)
  (local $$t36.1 i32)
  (local $$t37.0 i32)
  (local $$t37.1 i32)
  (local $$t38.0 i32)
  (local $$t38.1 i32)
  (local $$t39.0 i32)
  (local $$t39.1 i32)
  (local $$t40.0 i32)
  (local $$t40.1 i32)
  (local $$t41.0 i32)
  (local $$t41.1 i32)
  (local $$t42 i32)
  (local $$t43.0 i32)
  (local $$t43.1 i32)
  (local $$t44.0 i32)
  (local $$t44.1 i32)
  (local $$t45.0 i32)
  (local $$t45.1 i32)
  (local $$t46.0 i32)
  (local $$t46.1 i32)
  (local $$t47.0 i32)
  (local $$t47.1 i32)
  (local $$t48.0 i32)
  (local $$t48.1 i32)
  (local $$t49.0 i32)
  (local $$t49.1 i32)
  (local $$t50.0 i32)
  (local $$t50.1 i32)
  (local $$t51.0 i32)
  (local $$t51.1 i32)
  (local $$t52.0 i32)
  (local $$t52.1 i32)
  (local $$t53.0 i32)
  (local $$t53.1 i32)
  (local $$t54.0 i32)
  (local $$t54.1 i32)
  (local $$t55.0 i32)
  (local $$t55.1 i32)
  (local $$t56.0 i32)
  (local $$t56.1 i32)
  (local $$t57.0 i32)
  (local $$t57.1 i32)
  (local $$t58.0 i32)
  (local $$t58.1 i32)
  (local $$t59.0 i32)
  (local $$t59.1 i32)
  (local $$t60.0 i32)
  (local $$t60.1 i32)
  (local $$t61 i32)
  (local $$t62.0 i32)
  (local $$t62.1 i32)
  (local $$t63.0 i32)
  (local $$t63.1 i32)
  (local $$t64.0 i32)
  (local $$t64.1 i32)
  (local $$t65.0 i32)
  (local $$t65.1 i32)
  (local $$t66.0 i32)
  (local $$t66.1 i32)
  (local $$t67.0 i32)
  (local $$t67.1 i32)
  (local $$t68.0 i32)
  (local $$t68.1 i32)
  (local $$t69.0 i32)
  (local $$t69.1 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_13
        block $$Block_12
          block $$Block_11
            block $$Block_10
              block $$Block_9
                block $$Block_8
                  block $$Block_7
                    block $$Block_6
                      block $$Block_5
                        block $$Block_4
                          block $$Block_3
                            block $$Block_2
                              block $$Block_1
                                block $$Block_0
                                  block $$BlockSel
                                    local.get $$block_selector
                                    br_table 0 1 2 3 4 5 6 7 8 9 10 11 12 13 0
                                  end ;;$BlockSel
                                  i32.const 0
                                  local.set $$current_block

                                  ;;jump 3
                                  br $$Block_2

                                end ;;$Block_0
                                i32.const 1
                                local.set $$current_block

                                ;;(*mapNode).Parent(t9, this)
                                local.get $$t0.0
                                local.get $$t0.1
                                local.get $this.0
                                local.get $this.1
                                call $runtime.mapNode.Parent
                                local.set $$t1.1
                                local.get $$t1.0
                                call $runtime.Block.Release
                                local.set $$t1.0

                                ;;(*mapNode).Parent(t9, this)
                                local.get $$t0.0
                                local.get $$t0.1
                                local.get $this.0
                                local.get $this.1
                                call $runtime.mapNode.Parent
                                local.set $$t2.1
                                local.get $$t2.0
                                call $runtime.Block.Release
                                local.set $$t2.0

                                ;;(*mapNode).Parent(t1, this)
                                local.get $$t2.0
                                local.get $$t2.1
                                local.get $this.0
                                local.get $this.1
                                call $runtime.mapNode.Parent
                                local.set $$t3.1
                                local.get $$t3.0
                                call $runtime.Block.Release
                                local.set $$t3.0

                                ;;&t2.Left [#2]
                                local.get $$t3.0
                                call $runtime.Block.Retain
                                local.get $$t3.1
                                i32.const 8
                                i32.add
                                local.set $$t4.1
                                local.get $$t4.0
                                call $runtime.Block.Release
                                local.set $$t4.0

                                ;;*t3
                                local.get $$t4.1
                                i32.load offset=0 align=4
                                call $runtime.Block.Retain
                                local.get $$t4.1
                                i32.load offset=4 align=4
                                local.set $$t5.1
                                local.get $$t5.0
                                call $runtime.Block.Release
                                local.set $$t5.0

                                ;;t0 == t4
                                local.get $$t1.1
                                local.get $$t5.1
                                i32.eq
                                local.set $$t6

                                ;;if t5 goto 4 else 5
                                local.get $$t6
                                if
                                  br $$Block_3
                                else
                                  br $$Block_4
                                end

                              end ;;$Block_1
                              i32.const 2
                              local.set $$current_block

                              ;;&this.root [#1]
                              local.get $this.0
                              call $runtime.Block.Retain
                              local.get $this.1
                              i32.const 8
                              i32.add
                              local.set $$t7.1
                              local.get $$t7.0
                              call $runtime.Block.Release
                              local.set $$t7.0

                              ;;*t6
                              local.get $$t7.1
                              i32.load offset=0 align=4
                              call $runtime.Block.Retain
                              local.get $$t7.1
                              i32.load offset=4 align=4
                              local.set $$t8.1
                              local.get $$t8.0
                              call $runtime.Block.Release
                              local.set $$t8.0

                              ;;&t7.Color [#4]
                              local.get $$t8.0
                              call $runtime.Block.Retain
                              local.get $$t8.1
                              i32.const 24
                              i32.add
                              local.set $$t9.1
                              local.get $$t9.0
                              call $runtime.Block.Release
                              local.set $$t9.0

                              ;;*t8 = 1:mapColor
                              local.get $$t9.1
                              i32.const 1
                              i32.store offset=0 align=4

                              ;;return
                              br $$BlockFnBody

                            end ;;$Block_2
                            local.get $$current_block
                            i32.const 0
                            i32.eq
                            if (result i32 i32)
                              local.get $z.0
                              call $runtime.Block.Retain
                              local.get $z.1
                            else
                              local.get $$current_block
                              i32.const 6
                              i32.eq
                              if (result i32 i32)
                                local.get $$t10.0
                                call $runtime.Block.Retain
                                local.get $$t10.1
                              else
                                local.get $$current_block
                                i32.const 10
                                i32.eq
                                if (result i32 i32)
                                  local.get $$t11.0
                                  call $runtime.Block.Retain
                                  local.get $$t11.1
                                else
                                  local.get $$current_block
                                  i32.const 9
                                  i32.eq
                                  if (result i32 i32)
                                    local.get $$t12.0
                                    call $runtime.Block.Retain
                                    local.get $$t12.1
                                  else
                                    local.get $$t13.0
                                    call $runtime.Block.Retain
                                    local.get $$t13.1
                                  end
                                end
                              end
                            end
                            local.set $$t0.1
                            local.get $$t0.0
                            call $runtime.Block.Release
                            local.set $$t0.0
                            i32.const 3
                            local.set $$current_block

                            ;;(*mapNode).Parent(t9, this)
                            local.get $$t0.0
                            local.get $$t0.1
                            local.get $this.0
                            local.get $this.1
                            call $runtime.mapNode.Parent
                            local.set $$t14.1
                            local.get $$t14.0
                            call $runtime.Block.Release
                            local.set $$t14.0

                            ;;&t10.Color [#4]
                            local.get $$t14.0
                            call $runtime.Block.Retain
                            local.get $$t14.1
                            i32.const 24
                            i32.add
                            local.set $$t15.1
                            local.get $$t15.0
                            call $runtime.Block.Release
                            local.set $$t15.0

                            ;;*t11
                            local.get $$t15.1
                            i32.load offset=0 align=4
                            local.set $$t16

                            ;;t12 == 0:mapColor
                            local.get $$t16
                            i32.const 0
                            i32.eq
                            local.set $$t17

                            ;;if t13 goto 1 else 2
                            local.get $$t17
                            if
                              i32.const 1
                              local.set $$block_selector
                              br $$BlockDisp
                            else
                              i32.const 2
                              local.set $$block_selector
                              br $$BlockDisp
                            end

                          end ;;$Block_3
                          i32.const 4
                          local.set $$current_block

                          ;;(*mapNode).Parent(t9, this)
                          local.get $$t0.0
                          local.get $$t0.1
                          local.get $this.0
                          local.get $this.1
                          call $runtime.mapNode.Parent
                          local.set $$t18.1
                          local.get $$t18.0
                          call $runtime.Block.Release
                          local.set $$t18.0

                          ;;(*mapNode).Parent(t14, this)
                          local.get $$t18.0
                          local.get $$t18.1
                          local.get $this.0
                          local.get $this.1
                          call $runtime.mapNode.Parent
                          local.set $$t19.1
                          local.get $$t19.0
                          call $runtime.Block.Release
                          local.set $$t19.0

                          ;;&t15.Right [#3]
                          local.get $$t19.0
                          call $runtime.Block.Retain
                          local.get $$t19.1
                          i32.const 16
                          i32.add
                          local.set $$t20.1
                          local.get $$t20.0
                          call $runtime.Block.Release
                          local.set $$t20.0

                          ;;*t16
                          local.get $$t20.1
                          i32.load offset=0 align=4
                          call $runtime.Block.Retain
                          local.get $$t20.1
                          i32.load offset=4 align=4
                          local.set $$t21.1
                          local.get $$t21.0
                          call $runtime.Block.Release
                          local.set $$t21.0

                          ;;&t17.Color [#4]
                          local.get $$t21.0
                          call $runtime.Block.Retain
                          local.get $$t21.1
                          i32.const 24
                          i32.add
                          local.set $$t22.1
                          local.get $$t22.0
                          call $runtime.Block.Release
                          local.set $$t22.0

                          ;;*t18
                          local.get $$t22.1
                          i32.load offset=0 align=4
                          local.set $$t23

                          ;;t19 == 0:mapColor
                          local.get $$t23
                          i32.const 0
                          i32.eq
                          local.set $$t24

                          ;;if t20 goto 6 else 7
                          local.get $$t24
                          if
                            br $$Block_5
                          else
                            br $$Block_6
                          end

                        end ;;$Block_4
                        i32.const 5
                        local.set $$current_block

                        ;;(*mapNode).Parent(t9, this)
                        local.get $$t0.0
                        local.get $$t0.1
                        local.get $this.0
                        local.get $this.1
                        call $runtime.mapNode.Parent
                        local.set $$t25.1
                        local.get $$t25.0
                        call $runtime.Block.Release
                        local.set $$t25.0

                        ;;(*mapNode).Parent(t21, this)
                        local.get $$t25.0
                        local.get $$t25.1
                        local.get $this.0
                        local.get $this.1
                        call $runtime.mapNode.Parent
                        local.set $$t26.1
                        local.get $$t26.0
                        call $runtime.Block.Release
                        local.set $$t26.0

                        ;;&t22.Left [#2]
                        local.get $$t26.0
                        call $runtime.Block.Retain
                        local.get $$t26.1
                        i32.const 8
                        i32.add
                        local.set $$t27.1
                        local.get $$t27.0
                        call $runtime.Block.Release
                        local.set $$t27.0

                        ;;*t23
                        local.get $$t27.1
                        i32.load offset=0 align=4
                        call $runtime.Block.Retain
                        local.get $$t27.1
                        i32.load offset=4 align=4
                        local.set $$t28.1
                        local.get $$t28.0
                        call $runtime.Block.Release
                        local.set $$t28.0

                        ;;&t24.Color [#4]
                        local.get $$t28.0
                        call $runtime.Block.Retain
                        local.get $$t28.1
                        i32.const 24
                        i32.add
                        local.set $$t29.1
                        local.get $$t29.0
                        call $runtime.Block.Release
                        local.set $$t29.0

                        ;;*t25
                        local.get $$t29.1
                        i32.load offset=0 align=4
                        local.set $$t30

                        ;;t26 == 0:mapColor
                        local.get $$t30
                        i32.const 0
                        i32.eq
                        local.set $$t31

                        ;;if t27 goto 10 else 11
                        local.get $$t31
                        if
                          br $$Block_9
                        else
                          br $$Block_10
                        end

                      end ;;$Block_5
                      i32.const 6
                      local.set $$current_block

                      ;;(*mapNode).Parent(t9, this)
                      local.get $$t0.0
                      local.get $$t0.1
                      local.get $this.0
                      local.get $this.1
                      call $runtime.mapNode.Parent
                      local.set $$t32.1
                      local.get $$t32.0
                      call $runtime.Block.Release
                      local.set $$t32.0

                      ;;&t28.Color [#4]
                      local.get $$t32.0
                      call $runtime.Block.Retain
                      local.get $$t32.1
                      i32.const 24
                      i32.add
                      local.set $$t33.1
                      local.get $$t33.0
                      call $runtime.Block.Release
                      local.set $$t33.0

                      ;;*t29 = 1:mapColor
                      local.get $$t33.1
                      i32.const 1
                      i32.store offset=0 align=4

                      ;;&t17.Color [#4]
                      local.get $$t21.0
                      call $runtime.Block.Retain
                      local.get $$t21.1
                      i32.const 24
                      i32.add
                      local.set $$t34.1
                      local.get $$t34.0
                      call $runtime.Block.Release
                      local.set $$t34.0

                      ;;*t30 = 1:mapColor
                      local.get $$t34.1
                      i32.const 1
                      i32.store offset=0 align=4

                      ;;(*mapNode).Parent(t9, this)
                      local.get $$t0.0
                      local.get $$t0.1
                      local.get $this.0
                      local.get $this.1
                      call $runtime.mapNode.Parent
                      local.set $$t35.1
                      local.get $$t35.0
                      call $runtime.Block.Release
                      local.set $$t35.0

                      ;;(*mapNode).Parent(t31, this)
                      local.get $$t35.0
                      local.get $$t35.1
                      local.get $this.0
                      local.get $this.1
                      call $runtime.mapNode.Parent
                      local.set $$t36.1
                      local.get $$t36.0
                      call $runtime.Block.Release
                      local.set $$t36.0

                      ;;&t32.Color [#4]
                      local.get $$t36.0
                      call $runtime.Block.Retain
                      local.get $$t36.1
                      i32.const 24
                      i32.add
                      local.set $$t37.1
                      local.get $$t37.0
                      call $runtime.Block.Release
                      local.set $$t37.0

                      ;;*t33 = 0:mapColor
                      local.get $$t37.1
                      i32.const 0
                      i32.store offset=0 align=4

                      ;;(*mapNode).Parent(t9, this)
                      local.get $$t0.0
                      local.get $$t0.1
                      local.get $this.0
                      local.get $this.1
                      call $runtime.mapNode.Parent
                      local.set $$t38.1
                      local.get $$t38.0
                      call $runtime.Block.Release
                      local.set $$t38.0

                      ;;(*mapNode).Parent(t34, this)
                      local.get $$t38.0
                      local.get $$t38.1
                      local.get $this.0
                      local.get $this.1
                      call $runtime.mapNode.Parent
                      local.set $$t10.1
                      local.get $$t10.0
                      call $runtime.Block.Release
                      local.set $$t10.0

                      ;;jump 3
                      i32.const 3
                      local.set $$block_selector
                      br $$BlockDisp

                    end ;;$Block_6
                    i32.const 7
                    local.set $$current_block

                    ;;(*mapNode).Parent(t9, this)
                    local.get $$t0.0
                    local.get $$t0.1
                    local.get $this.0
                    local.get $this.1
                    call $runtime.mapNode.Parent
                    local.set $$t39.1
                    local.get $$t39.0
                    call $runtime.Block.Release
                    local.set $$t39.0

                    ;;&t36.Right [#3]
                    local.get $$t39.0
                    call $runtime.Block.Retain
                    local.get $$t39.1
                    i32.const 16
                    i32.add
                    local.set $$t40.1
                    local.get $$t40.0
                    call $runtime.Block.Release
                    local.set $$t40.0

                    ;;*t37
                    local.get $$t40.1
                    i32.load offset=0 align=4
                    call $runtime.Block.Retain
                    local.get $$t40.1
                    i32.load offset=4 align=4
                    local.set $$t41.1
                    local.get $$t41.0
                    call $runtime.Block.Release
                    local.set $$t41.0

                    ;;t9 == t38
                    local.get $$t0.1
                    local.get $$t41.1
                    i32.eq
                    local.set $$t42

                    ;;if t39 goto 8 else 9
                    local.get $$t42
                    if
                      br $$Block_7
                    else
                      br $$Block_8
                    end

                  end ;;$Block_7
                  i32.const 8
                  local.set $$current_block

                  ;;(*mapNode).Parent(t9, this)
                  local.get $$t0.0
                  local.get $$t0.1
                  local.get $this.0
                  local.get $this.1
                  call $runtime.mapNode.Parent
                  local.set $$t43.1
                  local.get $$t43.0
                  call $runtime.Block.Release
                  local.set $$t43.0

                  ;;(*mapImp).leftRotate(this, t40)
                  local.get $this.0
                  local.get $this.1
                  local.get $$t43.0
                  local.get $$t43.1
                  call $runtime.mapImp.leftRotate

                  ;;jump 9
                  br $$Block_8

                end ;;$Block_8
                local.get $$current_block
                i32.const 7
                i32.eq
                if (result i32 i32)
                  local.get $$t0.0
                  call $runtime.Block.Retain
                  local.get $$t0.1
                else
                  local.get $$t43.0
                  call $runtime.Block.Retain
                  local.get $$t43.1
                end
                local.set $$t12.1
                local.get $$t12.0
                call $runtime.Block.Release
                local.set $$t12.0
                i32.const 9
                local.set $$current_block

                ;;(*mapNode).Parent(t42, this)
                local.get $$t12.0
                local.get $$t12.1
                local.get $this.0
                local.get $this.1
                call $runtime.mapNode.Parent
                local.set $$t44.1
                local.get $$t44.0
                call $runtime.Block.Release
                local.set $$t44.0

                ;;&t43.Color [#4]
                local.get $$t44.0
                call $runtime.Block.Retain
                local.get $$t44.1
                i32.const 24
                i32.add
                local.set $$t45.1
                local.get $$t45.0
                call $runtime.Block.Release
                local.set $$t45.0

                ;;*t44 = 1:mapColor
                local.get $$t45.1
                i32.const 1
                i32.store offset=0 align=4

                ;;(*mapNode).Parent(t42, this)
                local.get $$t12.0
                local.get $$t12.1
                local.get $this.0
                local.get $this.1
                call $runtime.mapNode.Parent
                local.set $$t46.1
                local.get $$t46.0
                call $runtime.Block.Release
                local.set $$t46.0

                ;;(*mapNode).Parent(t45, this)
                local.get $$t46.0
                local.get $$t46.1
                local.get $this.0
                local.get $this.1
                call $runtime.mapNode.Parent
                local.set $$t47.1
                local.get $$t47.0
                call $runtime.Block.Release
                local.set $$t47.0

                ;;&t46.Color [#4]
                local.get $$t47.0
                call $runtime.Block.Retain
                local.get $$t47.1
                i32.const 24
                i32.add
                local.set $$t48.1
                local.get $$t48.0
                call $runtime.Block.Release
                local.set $$t48.0

                ;;*t47 = 0:mapColor
                local.get $$t48.1
                i32.const 0
                i32.store offset=0 align=4

                ;;(*mapNode).Parent(t42, this)
                local.get $$t12.0
                local.get $$t12.1
                local.get $this.0
                local.get $this.1
                call $runtime.mapNode.Parent
                local.set $$t49.1
                local.get $$t49.0
                call $runtime.Block.Release
                local.set $$t49.0

                ;;(*mapNode).Parent(t48, this)
                local.get $$t49.0
                local.get $$t49.1
                local.get $this.0
                local.get $this.1
                call $runtime.mapNode.Parent
                local.set $$t50.1
                local.get $$t50.0
                call $runtime.Block.Release
                local.set $$t50.0

                ;;(*mapImp).rightRotate(this, t49)
                local.get $this.0
                local.get $this.1
                local.get $$t50.0
                local.get $$t50.1
                call $runtime.mapImp.rightRotate

                ;;jump 3
                i32.const 3
                local.set $$block_selector
                br $$BlockDisp

              end ;;$Block_9
              i32.const 10
              local.set $$current_block

              ;;(*mapNode).Parent(t9, this)
              local.get $$t0.0
              local.get $$t0.1
              local.get $this.0
              local.get $this.1
              call $runtime.mapNode.Parent
              local.set $$t51.1
              local.get $$t51.0
              call $runtime.Block.Release
              local.set $$t51.0

              ;;&t51.Color [#4]
              local.get $$t51.0
              call $runtime.Block.Retain
              local.get $$t51.1
              i32.const 24
              i32.add
              local.set $$t52.1
              local.get $$t52.0
              call $runtime.Block.Release
              local.set $$t52.0

              ;;*t52 = 1:mapColor
              local.get $$t52.1
              i32.const 1
              i32.store offset=0 align=4

              ;;&t24.Color [#4]
              local.get $$t28.0
              call $runtime.Block.Retain
              local.get $$t28.1
              i32.const 24
              i32.add
              local.set $$t53.1
              local.get $$t53.0
              call $runtime.Block.Release
              local.set $$t53.0

              ;;*t53 = 1:mapColor
              local.get $$t53.1
              i32.const 1
              i32.store offset=0 align=4

              ;;(*mapNode).Parent(t9, this)
              local.get $$t0.0
              local.get $$t0.1
              local.get $this.0
              local.get $this.1
              call $runtime.mapNode.Parent
              local.set $$t54.1
              local.get $$t54.0
              call $runtime.Block.Release
              local.set $$t54.0

              ;;(*mapNode).Parent(t54, this)
              local.get $$t54.0
              local.get $$t54.1
              local.get $this.0
              local.get $this.1
              call $runtime.mapNode.Parent
              local.set $$t55.1
              local.get $$t55.0
              call $runtime.Block.Release
              local.set $$t55.0

              ;;&t55.Color [#4]
              local.get $$t55.0
              call $runtime.Block.Retain
              local.get $$t55.1
              i32.const 24
              i32.add
              local.set $$t56.1
              local.get $$t56.0
              call $runtime.Block.Release
              local.set $$t56.0

              ;;*t56 = 0:mapColor
              local.get $$t56.1
              i32.const 0
              i32.store offset=0 align=4

              ;;(*mapNode).Parent(t9, this)
              local.get $$t0.0
              local.get $$t0.1
              local.get $this.0
              local.get $this.1
              call $runtime.mapNode.Parent
              local.set $$t57.1
              local.get $$t57.0
              call $runtime.Block.Release
              local.set $$t57.0

              ;;(*mapNode).Parent(t57, this)
              local.get $$t57.0
              local.get $$t57.1
              local.get $this.0
              local.get $this.1
              call $runtime.mapNode.Parent
              local.set $$t11.1
              local.get $$t11.0
              call $runtime.Block.Release
              local.set $$t11.0

              ;;jump 3
              i32.const 3
              local.set $$block_selector
              br $$BlockDisp

            end ;;$Block_10
            i32.const 11
            local.set $$current_block

            ;;(*mapNode).Parent(t9, this)
            local.get $$t0.0
            local.get $$t0.1
            local.get $this.0
            local.get $this.1
            call $runtime.mapNode.Parent
            local.set $$t58.1
            local.get $$t58.0
            call $runtime.Block.Release
            local.set $$t58.0

            ;;&t59.Left [#2]
            local.get $$t58.0
            call $runtime.Block.Retain
            local.get $$t58.1
            i32.const 8
            i32.add
            local.set $$t59.1
            local.get $$t59.0
            call $runtime.Block.Release
            local.set $$t59.0

            ;;*t60
            local.get $$t59.1
            i32.load offset=0 align=4
            call $runtime.Block.Retain
            local.get $$t59.1
            i32.load offset=4 align=4
            local.set $$t60.1
            local.get $$t60.0
            call $runtime.Block.Release
            local.set $$t60.0

            ;;t9 == t61
            local.get $$t0.1
            local.get $$t60.1
            i32.eq
            local.set $$t61

            ;;if t62 goto 12 else 13
            local.get $$t61
            if
              br $$Block_11
            else
              br $$Block_12
            end

          end ;;$Block_11
          i32.const 12
          local.set $$current_block

          ;;(*mapNode).Parent(t9, this)
          local.get $$t0.0
          local.get $$t0.1
          local.get $this.0
          local.get $this.1
          call $runtime.mapNode.Parent
          local.set $$t62.1
          local.get $$t62.0
          call $runtime.Block.Release
          local.set $$t62.0

          ;;(*mapImp).rightRotate(this, t63)
          local.get $this.0
          local.get $this.1
          local.get $$t62.0
          local.get $$t62.1
          call $runtime.mapImp.rightRotate

          ;;jump 13
          br $$Block_12

        end ;;$Block_12
        local.get $$current_block
        i32.const 11
        i32.eq
        if (result i32 i32)
          local.get $$t0.0
          call $runtime.Block.Retain
          local.get $$t0.1
        else
          local.get $$t62.0
          call $runtime.Block.Retain
          local.get $$t62.1
        end
        local.set $$t13.1
        local.get $$t13.0
        call $runtime.Block.Release
        local.set $$t13.0
        i32.const 13
        local.set $$current_block

        ;;(*mapNode).Parent(t65, this)
        local.get $$t13.0
        local.get $$t13.1
        local.get $this.0
        local.get $this.1
        call $runtime.mapNode.Parent
        local.set $$t63.1
        local.get $$t63.0
        call $runtime.Block.Release
        local.set $$t63.0

        ;;&t66.Color [#4]
        local.get $$t63.0
        call $runtime.Block.Retain
        local.get $$t63.1
        i32.const 24
        i32.add
        local.set $$t64.1
        local.get $$t64.0
        call $runtime.Block.Release
        local.set $$t64.0

        ;;*t67 = 1:mapColor
        local.get $$t64.1
        i32.const 1
        i32.store offset=0 align=4

        ;;(*mapNode).Parent(t65, this)
        local.get $$t13.0
        local.get $$t13.1
        local.get $this.0
        local.get $this.1
        call $runtime.mapNode.Parent
        local.set $$t65.1
        local.get $$t65.0
        call $runtime.Block.Release
        local.set $$t65.0

        ;;(*mapNode).Parent(t68, this)
        local.get $$t65.0
        local.get $$t65.1
        local.get $this.0
        local.get $this.1
        call $runtime.mapNode.Parent
        local.set $$t66.1
        local.get $$t66.0
        call $runtime.Block.Release
        local.set $$t66.0

        ;;&t69.Color [#4]
        local.get $$t66.0
        call $runtime.Block.Retain
        local.get $$t66.1
        i32.const 24
        i32.add
        local.set $$t67.1
        local.get $$t67.0
        call $runtime.Block.Release
        local.set $$t67.0

        ;;*t70 = 0:mapColor
        local.get $$t67.1
        i32.const 0
        i32.store offset=0 align=4

        ;;(*mapNode).Parent(t65, this)
        local.get $$t13.0
        local.get $$t13.1
        local.get $this.0
        local.get $this.1
        call $runtime.mapNode.Parent
        local.set $$t68.1
        local.get $$t68.0
        call $runtime.Block.Release
        local.set $$t68.0

        ;;(*mapNode).Parent(t71, this)
        local.get $$t68.0
        local.get $$t68.1
        local.get $this.0
        local.get $this.1
        call $runtime.mapNode.Parent
        local.set $$t69.1
        local.get $$t69.0
        call $runtime.Block.Release
        local.set $$t69.0

        ;;(*mapImp).leftRotate(this, t72)
        local.get $this.0
        local.get $this.1
        local.get $$t69.0
        local.get $$t69.1
        call $runtime.mapImp.leftRotate

        ;;jump 3
        i32.const 3
        local.set $$block_selector
        br $$BlockDisp

      end ;;$Block_13
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$t0.0
  call $runtime.Block.Release
  local.get $$t1.0
  call $runtime.Block.Release
  local.get $$t2.0
  call $runtime.Block.Release
  local.get $$t3.0
  call $runtime.Block.Release
  local.get $$t4.0
  call $runtime.Block.Release
  local.get $$t5.0
  call $runtime.Block.Release
  local.get $$t7.0
  call $runtime.Block.Release
  local.get $$t8.0
  call $runtime.Block.Release
  local.get $$t9.0
  call $runtime.Block.Release
  local.get $$t10.0
  call $runtime.Block.Release
  local.get $$t11.0
  call $runtime.Block.Release
  local.get $$t12.0
  call $runtime.Block.Release
  local.get $$t13.0
  call $runtime.Block.Release
  local.get $$t14.0
  call $runtime.Block.Release
  local.get $$t15.0
  call $runtime.Block.Release
  local.get $$t18.0
  call $runtime.Block.Release
  local.get $$t19.0
  call $runtime.Block.Release
  local.get $$t20.0
  call $runtime.Block.Release
  local.get $$t21.0
  call $runtime.Block.Release
  local.get $$t22.0
  call $runtime.Block.Release
  local.get $$t25.0
  call $runtime.Block.Release
  local.get $$t26.0
  call $runtime.Block.Release
  local.get $$t27.0
  call $runtime.Block.Release
  local.get $$t28.0
  call $runtime.Block.Release
  local.get $$t29.0
  call $runtime.Block.Release
  local.get $$t32.0
  call $runtime.Block.Release
  local.get $$t33.0
  call $runtime.Block.Release
  local.get $$t34.0
  call $runtime.Block.Release
  local.get $$t35.0
  call $runtime.Block.Release
  local.get $$t36.0
  call $runtime.Block.Release
  local.get $$t37.0
  call $runtime.Block.Release
  local.get $$t38.0
  call $runtime.Block.Release
  local.get $$t39.0
  call $runtime.Block.Release
  local.get $$t40.0
  call $runtime.Block.Release
  local.get $$t41.0
  call $runtime.Block.Release
  local.get $$t43.0
  call $runtime.Block.Release
  local.get $$t44.0
  call $runtime.Block.Release
  local.get $$t45.0
  call $runtime.Block.Release
  local.get $$t46.0
  call $runtime.Block.Release
  local.get $$t47.0
  call $runtime.Block.Release
  local.get $$t48.0
  call $runtime.Block.Release
  local.get $$t49.0
  call $runtime.Block.Release
  local.get $$t50.0
  call $runtime.Block.Release
  local.get $$t51.0
  call $runtime.Block.Release
  local.get $$t52.0
  call $runtime.Block.Release
  local.get $$t53.0
  call $runtime.Block.Release
  local.get $$t54.0
  call $runtime.Block.Release
  local.get $$t55.0
  call $runtime.Block.Release
  local.get $$t56.0
  call $runtime.Block.Release
  local.get $$t57.0
  call $runtime.Block.Release
  local.get $$t58.0
  call $runtime.Block.Release
  local.get $$t59.0
  call $runtime.Block.Release
  local.get $$t60.0
  call $runtime.Block.Release
  local.get $$t62.0
  call $runtime.Block.Release
  local.get $$t63.0
  call $runtime.Block.Release
  local.get $$t64.0
  call $runtime.Block.Release
  local.get $$t65.0
  call $runtime.Block.Release
  local.get $$t66.0
  call $runtime.Block.Release
  local.get $$t67.0
  call $runtime.Block.Release
  local.get $$t68.0
  call $runtime.Block.Release
  local.get $$t69.0
  call $runtime.Block.Release
) ;;runtime.mapImp.insertFixup

(func $runtime.mapImp.leftRotate (param $this.0 i32) (param $this.1 i32) (param $x.0 i32) (param $x.1 i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$t0.0 i32)
  (local $$t0.1 i32)
  (local $$t1.0 i32)
  (local $$t1.1 i32)
  (local $$t2.0 i32)
  (local $$t2.1 i32)
  (local $$t3.0 i32)
  (local $$t3.1 i32)
  (local $$t4 i32)
  (local $$t5.0 i32)
  (local $$t5.1 i32)
  (local $$t6.0 i32)
  (local $$t6.1 i32)
  (local $$t7.0 i32)
  (local $$t7.1 i32)
  (local $$t8.0 i32)
  (local $$t8.1 i32)
  (local $$t9.0 i32)
  (local $$t9.1 i32)
  (local $$t10.0 i32)
  (local $$t10.1 i32)
  (local $$t11.0 i32)
  (local $$t11.1 i32)
  (local $$t12.0 i32)
  (local $$t12.1 i32)
  (local $$t13.0 i32)
  (local $$t13.1 i32)
  (local $$t14 i32)
  (local $$t15.0 i32)
  (local $$t15.1 i32)
  (local $$t16.0 i32)
  (local $$t16.1 i32)
  (local $$t17.0 i32)
  (local $$t17.1 i32)
  (local $$t18.0 i32)
  (local $$t18.1 i32)
  (local $$t19.0 i32)
  (local $$t19.1 i32)
  (local $$t20.0 i32)
  (local $$t20.1 i32)
  (local $$t21 i32)
  (local $$t22.0 i32)
  (local $$t22.1 i32)
  (local $$t23.0 i32)
  (local $$t23.1 i32)
  (local $$t24.0 i32)
  (local $$t24.1 i32)
  (local $$t25.0 i32)
  (local $$t25.1 i32)
  (local $$t26.0 i32)
  (local $$t26.1 i32)
  (local $$t27 i32)
  (local $$t28.0 i32)
  (local $$t28.1 i32)
  (local $$t29.0 i32)
  (local $$t29.1 i32)
  (local $$t30.0 i32)
  (local $$t30.1 i32)
  (local $$t31.0 i32)
  (local $$t31.1 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_9
        block $$Block_8
          block $$Block_7
            block $$Block_6
              block $$Block_5
                block $$Block_4
                  block $$Block_3
                    block $$Block_2
                      block $$Block_1
                        block $$Block_0
                          block $$BlockSel
                            local.get $$block_selector
                            br_table 0 1 2 3 4 5 6 7 8 9 0
                          end ;;$BlockSel
                          i32.const 0
                          local.set $$current_block

                          ;;&x.Right [#3]
                          local.get $x.0
                          call $runtime.Block.Retain
                          local.get $x.1
                          i32.const 16
                          i32.add
                          local.set $$t0.1
                          local.get $$t0.0
                          call $runtime.Block.Release
                          local.set $$t0.0

                          ;;*t0
                          local.get $$t0.1
                          i32.load offset=0 align=4
                          call $runtime.Block.Retain
                          local.get $$t0.1
                          i32.load offset=4 align=4
                          local.set $$t1.1
                          local.get $$t1.0
                          call $runtime.Block.Release
                          local.set $$t1.0

                          ;;&this.NIL [#0]
                          local.get $this.0
                          call $runtime.Block.Retain
                          local.get $this.1
                          i32.const 0
                          i32.add
                          local.set $$t2.1
                          local.get $$t2.0
                          call $runtime.Block.Release
                          local.set $$t2.0

                          ;;*t2
                          local.get $$t2.1
                          i32.load offset=0 align=4
                          call $runtime.Block.Retain
                          local.get $$t2.1
                          i32.load offset=4 align=4
                          local.set $$t3.1
                          local.get $$t3.0
                          call $runtime.Block.Release
                          local.set $$t3.0

                          ;;t1 == t3
                          local.get $$t1.1
                          local.get $$t3.1
                          i32.eq
                          local.set $$t4

                          ;;if t4 goto 1 else 2
                          local.get $$t4
                          if
                            br $$Block_0
                          else
                            br $$Block_1
                          end

                        end ;;$Block_0
                        i32.const 1
                        local.set $$current_block

                        ;;return
                        br $$BlockFnBody

                      end ;;$Block_1
                      i32.const 2
                      local.set $$current_block

                      ;;&x.Right [#3]
                      local.get $x.0
                      call $runtime.Block.Retain
                      local.get $x.1
                      i32.const 16
                      i32.add
                      local.set $$t5.1
                      local.get $$t5.0
                      call $runtime.Block.Release
                      local.set $$t5.0

                      ;;*t5
                      local.get $$t5.1
                      i32.load offset=0 align=4
                      call $runtime.Block.Retain
                      local.get $$t5.1
                      i32.load offset=4 align=4
                      local.set $$t6.1
                      local.get $$t6.0
                      call $runtime.Block.Release
                      local.set $$t6.0

                      ;;&x.Right [#3]
                      local.get $x.0
                      call $runtime.Block.Retain
                      local.get $x.1
                      i32.const 16
                      i32.add
                      local.set $$t7.1
                      local.get $$t7.0
                      call $runtime.Block.Release
                      local.set $$t7.0

                      ;;&t6.Left [#2]
                      local.get $$t6.0
                      call $runtime.Block.Retain
                      local.get $$t6.1
                      i32.const 8
                      i32.add
                      local.set $$t8.1
                      local.get $$t8.0
                      call $runtime.Block.Release
                      local.set $$t8.0

                      ;;*t8
                      local.get $$t8.1
                      i32.load offset=0 align=4
                      call $runtime.Block.Retain
                      local.get $$t8.1
                      i32.load offset=4 align=4
                      local.set $$t9.1
                      local.get $$t9.0
                      call $runtime.Block.Release
                      local.set $$t9.0

                      ;;*t7 = t9
                      local.get $$t7.1
                      local.get $$t9.0
                      call $runtime.Block.Retain
                      local.get $$t7.1
                      i32.load offset=0 align=1
                      call $runtime.Block.Release
                      i32.store offset=0 align=1
                      local.get $$t7.1
                      local.get $$t9.1
                      i32.store offset=4 align=4

                      ;;&t6.Left [#2]
                      local.get $$t6.0
                      call $runtime.Block.Retain
                      local.get $$t6.1
                      i32.const 8
                      i32.add
                      local.set $$t10.1
                      local.get $$t10.0
                      call $runtime.Block.Release
                      local.set $$t10.0

                      ;;*t10
                      local.get $$t10.1
                      i32.load offset=0 align=4
                      call $runtime.Block.Retain
                      local.get $$t10.1
                      i32.load offset=4 align=4
                      local.set $$t11.1
                      local.get $$t11.0
                      call $runtime.Block.Release
                      local.set $$t11.0

                      ;;&this.NIL [#0]
                      local.get $this.0
                      call $runtime.Block.Retain
                      local.get $this.1
                      i32.const 0
                      i32.add
                      local.set $$t12.1
                      local.get $$t12.0
                      call $runtime.Block.Release
                      local.set $$t12.0

                      ;;*t12
                      local.get $$t12.1
                      i32.load offset=0 align=4
                      call $runtime.Block.Retain
                      local.get $$t12.1
                      i32.load offset=4 align=4
                      local.set $$t13.1
                      local.get $$t13.0
                      call $runtime.Block.Release
                      local.set $$t13.0

                      ;;t11 != t13
                      local.get $$t11.1
                      local.get $$t13.1
                      i32.eq
                      i32.eqz
                      local.set $$t14

                      ;;if t14 goto 3 else 4
                      local.get $$t14
                      if
                        br $$Block_2
                      else
                        br $$Block_3
                      end

                    end ;;$Block_2
                    i32.const 3
                    local.set $$current_block

                    ;;&t6.Left [#2]
                    local.get $$t6.0
                    call $runtime.Block.Retain
                    local.get $$t6.1
                    i32.const 8
                    i32.add
                    local.set $$t15.1
                    local.get $$t15.0
                    call $runtime.Block.Release
                    local.set $$t15.0

                    ;;*t15
                    local.get $$t15.1
                    i32.load offset=0 align=4
                    call $runtime.Block.Retain
                    local.get $$t15.1
                    i32.load offset=4 align=4
                    local.set $$t16.1
                    local.get $$t16.0
                    call $runtime.Block.Release
                    local.set $$t16.0

                    ;;(*mapNode).SetParent(t16, x)
                    local.get $$t16.0
                    local.get $$t16.1
                    local.get $x.0
                    local.get $x.1
                    call $runtime.mapNode.SetParent

                    ;;jump 4
                    br $$Block_3

                  end ;;$Block_3
                  i32.const 4
                  local.set $$current_block

                  ;;(*mapNode).Parent(x, this)
                  local.get $x.0
                  local.get $x.1
                  local.get $this.0
                  local.get $this.1
                  call $runtime.mapNode.Parent
                  local.set $$t17.1
                  local.get $$t17.0
                  call $runtime.Block.Release
                  local.set $$t17.0

                  ;;(*mapNode).SetParent(t6, t18)
                  local.get $$t6.0
                  local.get $$t6.1
                  local.get $$t17.0
                  local.get $$t17.1
                  call $runtime.mapNode.SetParent

                  ;;(*mapNode).Parent(x, this)
                  local.get $x.0
                  local.get $x.1
                  local.get $this.0
                  local.get $this.1
                  call $runtime.mapNode.Parent
                  local.set $$t18.1
                  local.get $$t18.0
                  call $runtime.Block.Release
                  local.set $$t18.0

                  ;;&this.NIL [#0]
                  local.get $this.0
                  call $runtime.Block.Retain
                  local.get $this.1
                  i32.const 0
                  i32.add
                  local.set $$t19.1
                  local.get $$t19.0
                  call $runtime.Block.Release
                  local.set $$t19.0

                  ;;*t21
                  local.get $$t19.1
                  i32.load offset=0 align=4
                  call $runtime.Block.Retain
                  local.get $$t19.1
                  i32.load offset=4 align=4
                  local.set $$t20.1
                  local.get $$t20.0
                  call $runtime.Block.Release
                  local.set $$t20.0

                  ;;t20 == t22
                  local.get $$t18.1
                  local.get $$t20.1
                  i32.eq
                  local.set $$t21

                  ;;if t23 goto 5 else 7
                  local.get $$t21
                  if
                    br $$Block_4
                  else
                    br $$Block_6
                  end

                end ;;$Block_4
                i32.const 5
                local.set $$current_block

                ;;&this.root [#1]
                local.get $this.0
                call $runtime.Block.Retain
                local.get $this.1
                i32.const 8
                i32.add
                local.set $$t22.1
                local.get $$t22.0
                call $runtime.Block.Release
                local.set $$t22.0

                ;;*t24 = t6
                local.get $$t22.1
                local.get $$t6.0
                call $runtime.Block.Retain
                local.get $$t22.1
                i32.load offset=0 align=1
                call $runtime.Block.Release
                i32.store offset=0 align=1
                local.get $$t22.1
                local.get $$t6.1
                i32.store offset=4 align=4

                ;;jump 6
                br $$Block_5

              end ;;$Block_5
              i32.const 6
              local.set $$current_block

              ;;&t6.Left [#2]
              local.get $$t6.0
              call $runtime.Block.Retain
              local.get $$t6.1
              i32.const 8
              i32.add
              local.set $$t23.1
              local.get $$t23.0
              call $runtime.Block.Release
              local.set $$t23.0

              ;;*t25 = x
              local.get $$t23.1
              local.get $x.0
              call $runtime.Block.Retain
              local.get $$t23.1
              i32.load offset=0 align=1
              call $runtime.Block.Release
              i32.store offset=0 align=1
              local.get $$t23.1
              local.get $x.1
              i32.store offset=4 align=4

              ;;(*mapNode).SetParent(x, t6)
              local.get $x.0
              local.get $x.1
              local.get $$t6.0
              local.get $$t6.1
              call $runtime.mapNode.SetParent

              ;;return
              br $$BlockFnBody

            end ;;$Block_6
            i32.const 7
            local.set $$current_block

            ;;(*mapNode).Parent(x, this)
            local.get $x.0
            local.get $x.1
            local.get $this.0
            local.get $this.1
            call $runtime.mapNode.Parent
            local.set $$t24.1
            local.get $$t24.0
            call $runtime.Block.Release
            local.set $$t24.0

            ;;&t27.Left [#2]
            local.get $$t24.0
            call $runtime.Block.Retain
            local.get $$t24.1
            i32.const 8
            i32.add
            local.set $$t25.1
            local.get $$t25.0
            call $runtime.Block.Release
            local.set $$t25.0

            ;;*t28
            local.get $$t25.1
            i32.load offset=0 align=4
            call $runtime.Block.Retain
            local.get $$t25.1
            i32.load offset=4 align=4
            local.set $$t26.1
            local.get $$t26.0
            call $runtime.Block.Release
            local.set $$t26.0

            ;;x == t29
            local.get $x.1
            local.get $$t26.1
            i32.eq
            local.set $$t27

            ;;if t30 goto 8 else 9
            local.get $$t27
            if
              br $$Block_7
            else
              br $$Block_8
            end

          end ;;$Block_7
          i32.const 8
          local.set $$current_block

          ;;(*mapNode).Parent(x, this)
          local.get $x.0
          local.get $x.1
          local.get $this.0
          local.get $this.1
          call $runtime.mapNode.Parent
          local.set $$t28.1
          local.get $$t28.0
          call $runtime.Block.Release
          local.set $$t28.0

          ;;&t31.Left [#2]
          local.get $$t28.0
          call $runtime.Block.Retain
          local.get $$t28.1
          i32.const 8
          i32.add
          local.set $$t29.1
          local.get $$t29.0
          call $runtime.Block.Release
          local.set $$t29.0

          ;;*t32 = t6
          local.get $$t29.1
          local.get $$t6.0
          call $runtime.Block.Retain
          local.get $$t29.1
          i32.load offset=0 align=1
          call $runtime.Block.Release
          i32.store offset=0 align=1
          local.get $$t29.1
          local.get $$t6.1
          i32.store offset=4 align=4

          ;;jump 6
          i32.const 6
          local.set $$block_selector
          br $$BlockDisp

        end ;;$Block_8
        i32.const 9
        local.set $$current_block

        ;;(*mapNode).Parent(x, this)
        local.get $x.0
        local.get $x.1
        local.get $this.0
        local.get $this.1
        call $runtime.mapNode.Parent
        local.set $$t30.1
        local.get $$t30.0
        call $runtime.Block.Release
        local.set $$t30.0

        ;;&t33.Right [#3]
        local.get $$t30.0
        call $runtime.Block.Retain
        local.get $$t30.1
        i32.const 16
        i32.add
        local.set $$t31.1
        local.get $$t31.0
        call $runtime.Block.Release
        local.set $$t31.0

        ;;*t34 = t6
        local.get $$t31.1
        local.get $$t6.0
        call $runtime.Block.Retain
        local.get $$t31.1
        i32.load offset=0 align=1
        call $runtime.Block.Release
        i32.store offset=0 align=1
        local.get $$t31.1
        local.get $$t6.1
        i32.store offset=4 align=4

        ;;jump 6
        i32.const 6
        local.set $$block_selector
        br $$BlockDisp

      end ;;$Block_9
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$t0.0
  call $runtime.Block.Release
  local.get $$t1.0
  call $runtime.Block.Release
  local.get $$t2.0
  call $runtime.Block.Release
  local.get $$t3.0
  call $runtime.Block.Release
  local.get $$t5.0
  call $runtime.Block.Release
  local.get $$t6.0
  call $runtime.Block.Release
  local.get $$t7.0
  call $runtime.Block.Release
  local.get $$t8.0
  call $runtime.Block.Release
  local.get $$t9.0
  call $runtime.Block.Release
  local.get $$t10.0
  call $runtime.Block.Release
  local.get $$t11.0
  call $runtime.Block.Release
  local.get $$t12.0
  call $runtime.Block.Release
  local.get $$t13.0
  call $runtime.Block.Release
  local.get $$t15.0
  call $runtime.Block.Release
  local.get $$t16.0
  call $runtime.Block.Release
  local.get $$t17.0
  call $runtime.Block.Release
  local.get $$t18.0
  call $runtime.Block.Release
  local.get $$t19.0
  call $runtime.Block.Release
  local.get $$t20.0
  call $runtime.Block.Release
  local.get $$t22.0
  call $runtime.Block.Release
  local.get $$t23.0
  call $runtime.Block.Release
  local.get $$t24.0
  call $runtime.Block.Release
  local.get $$t25.0
  call $runtime.Block.Release
  local.get $$t26.0
  call $runtime.Block.Release
  local.get $$t28.0
  call $runtime.Block.Release
  local.get $$t29.0
  call $runtime.Block.Release
  local.get $$t30.0
  call $runtime.Block.Release
  local.get $$t31.0
  call $runtime.Block.Release
) ;;runtime.mapImp.leftRotate

(func $runtime.mapImp.min (param $this.0 i32) (param $this.1 i32) (param $x.0 i32) (param $x.1 i32) (result i32 i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$ret_0.0 i32)
  (local $$ret_0.1 i32)
  (local $$t0.0 i32)
  (local $$t0.1 i32)
  (local $$t1.0 i32)
  (local $$t1.1 i32)
  (local $$t2 i32)
  (local $$t3.0 i32)
  (local $$t3.1 i32)
  (local $$t4.0 i32)
  (local $$t4.1 i32)
  (local $$t5.0 i32)
  (local $$t5.1 i32)
  (local $$t6.0 i32)
  (local $$t6.1 i32)
  (local $$t7.0 i32)
  (local $$t7.1 i32)
  (local $$t8.0 i32)
  (local $$t8.1 i32)
  (local $$t9.0 i32)
  (local $$t9.1 i32)
  (local $$t10.0 i32)
  (local $$t10.1 i32)
  (local $$t11.0 i32)
  (local $$t11.1 i32)
  (local $$t12 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_4
        block $$Block_3
          block $$Block_2
            block $$Block_1
              block $$Block_0
                block $$BlockSel
                  local.get $$block_selector
                  br_table 0 1 2 3 4 0
                end ;;$BlockSel
                i32.const 0
                local.set $$current_block

                ;;&this.NIL [#0]
                local.get $this.0
                call $runtime.Block.Retain
                local.get $this.1
                i32.const 0
                i32.add
                local.set $$t0.1
                local.get $$t0.0
                call $runtime.Block.Release
                local.set $$t0.0

                ;;*t0
                local.get $$t0.1
                i32.load offset=0 align=4
                call $runtime.Block.Retain
                local.get $$t0.1
                i32.load offset=4 align=4
                local.set $$t1.1
                local.get $$t1.0
                call $runtime.Block.Release
                local.set $$t1.0

                ;;x == t1
                local.get $x.1
                local.get $$t1.1
                i32.eq
                local.set $$t2

                ;;if t2 goto 1 else 4
                local.get $$t2
                if
                  br $$Block_0
                else
                  br $$Block_3
                end

              end ;;$Block_0
              i32.const 1
              local.set $$current_block

              ;;&this.NIL [#0]
              local.get $this.0
              call $runtime.Block.Retain
              local.get $this.1
              i32.const 0
              i32.add
              local.set $$t3.1
              local.get $$t3.0
              call $runtime.Block.Release
              local.set $$t3.0

              ;;*t3
              local.get $$t3.1
              i32.load offset=0 align=4
              call $runtime.Block.Retain
              local.get $$t3.1
              i32.load offset=4 align=4
              local.set $$t4.1
              local.get $$t4.0
              call $runtime.Block.Release
              local.set $$t4.0

              ;;return t4
              local.get $$t4.0
              call $runtime.Block.Retain
              local.get $$t4.1
              local.set $$ret_0.1
              local.get $$ret_0.0
              call $runtime.Block.Release
              local.set $$ret_0.0
              br $$BlockFnBody

            end ;;$Block_1
            i32.const 2
            local.set $$current_block

            ;;&t7.Left [#2]
            local.get $$t5.0
            call $runtime.Block.Retain
            local.get $$t5.1
            i32.const 8
            i32.add
            local.set $$t6.1
            local.get $$t6.0
            call $runtime.Block.Release
            local.set $$t6.0

            ;;*t5
            local.get $$t6.1
            i32.load offset=0 align=4
            call $runtime.Block.Retain
            local.get $$t6.1
            i32.load offset=4 align=4
            local.set $$t7.1
            local.get $$t7.0
            call $runtime.Block.Release
            local.set $$t7.0

            ;;jump 4
            br $$Block_3

          end ;;$Block_2
          i32.const 3
          local.set $$current_block

          ;;return t7
          local.get $$t5.0
          call $runtime.Block.Retain
          local.get $$t5.1
          local.set $$ret_0.1
          local.get $$ret_0.0
          call $runtime.Block.Release
          local.set $$ret_0.0
          br $$BlockFnBody

        end ;;$Block_3
        local.get $$current_block
        i32.const 0
        i32.eq
        if (result i32 i32)
          local.get $x.0
          call $runtime.Block.Retain
          local.get $x.1
        else
          local.get $$t7.0
          call $runtime.Block.Retain
          local.get $$t7.1
        end
        local.set $$t5.1
        local.get $$t5.0
        call $runtime.Block.Release
        local.set $$t5.0
        i32.const 4
        local.set $$current_block

        ;;&t7.Left [#2]
        local.get $$t5.0
        call $runtime.Block.Retain
        local.get $$t5.1
        i32.const 8
        i32.add
        local.set $$t8.1
        local.get $$t8.0
        call $runtime.Block.Release
        local.set $$t8.0

        ;;*t8
        local.get $$t8.1
        i32.load offset=0 align=4
        call $runtime.Block.Retain
        local.get $$t8.1
        i32.load offset=4 align=4
        local.set $$t9.1
        local.get $$t9.0
        call $runtime.Block.Release
        local.set $$t9.0

        ;;&this.NIL [#0]
        local.get $this.0
        call $runtime.Block.Retain
        local.get $this.1
        i32.const 0
        i32.add
        local.set $$t10.1
        local.get $$t10.0
        call $runtime.Block.Release
        local.set $$t10.0

        ;;*t10
        local.get $$t10.1
        i32.load offset=0 align=4
        call $runtime.Block.Retain
        local.get $$t10.1
        i32.load offset=4 align=4
        local.set $$t11.1
        local.get $$t11.0
        call $runtime.Block.Release
        local.set $$t11.0

        ;;t9 != t11
        local.get $$t9.1
        local.get $$t11.1
        i32.eq
        i32.eqz
        local.set $$t12

        ;;if t12 goto 2 else 3
        local.get $$t12
        if
          i32.const 2
          local.set $$block_selector
          br $$BlockDisp
        else
          i32.const 3
          local.set $$block_selector
          br $$BlockDisp
        end

      end ;;$Block_4
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$ret_0.0
  call $runtime.Block.Retain
  local.get $$ret_0.1
  local.get $$ret_0.0
  call $runtime.Block.Release
  local.get $$t0.0
  call $runtime.Block.Release
  local.get $$t1.0
  call $runtime.Block.Release
  local.get $$t3.0
  call $runtime.Block.Release
  local.get $$t4.0
  call $runtime.Block.Release
  local.get $$t5.0
  call $runtime.Block.Release
  local.get $$t6.0
  call $runtime.Block.Release
  local.get $$t7.0
  call $runtime.Block.Release
  local.get $$t8.0
  call $runtime.Block.Release
  local.get $$t9.0
  call $runtime.Block.Release
  local.get $$t10.0
  call $runtime.Block.Release
  local.get $$t11.0
  call $runtime.Block.Release
) ;;runtime.mapImp.min

(func $runtime.mapImp.rightRotate (param $this.0 i32) (param $this.1 i32) (param $x.0 i32) (param $x.1 i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$t0.0 i32)
  (local $$t0.1 i32)
  (local $$t1.0 i32)
  (local $$t1.1 i32)
  (local $$t2.0 i32)
  (local $$t2.1 i32)
  (local $$t3.0 i32)
  (local $$t3.1 i32)
  (local $$t4 i32)
  (local $$t5.0 i32)
  (local $$t5.1 i32)
  (local $$t6.0 i32)
  (local $$t6.1 i32)
  (local $$t7.0 i32)
  (local $$t7.1 i32)
  (local $$t8.0 i32)
  (local $$t8.1 i32)
  (local $$t9.0 i32)
  (local $$t9.1 i32)
  (local $$t10.0 i32)
  (local $$t10.1 i32)
  (local $$t11.0 i32)
  (local $$t11.1 i32)
  (local $$t12.0 i32)
  (local $$t12.1 i32)
  (local $$t13.0 i32)
  (local $$t13.1 i32)
  (local $$t14 i32)
  (local $$t15.0 i32)
  (local $$t15.1 i32)
  (local $$t16.0 i32)
  (local $$t16.1 i32)
  (local $$t17.0 i32)
  (local $$t17.1 i32)
  (local $$t18.0 i32)
  (local $$t18.1 i32)
  (local $$t19.0 i32)
  (local $$t19.1 i32)
  (local $$t20.0 i32)
  (local $$t20.1 i32)
  (local $$t21 i32)
  (local $$t22.0 i32)
  (local $$t22.1 i32)
  (local $$t23.0 i32)
  (local $$t23.1 i32)
  (local $$t24.0 i32)
  (local $$t24.1 i32)
  (local $$t25.0 i32)
  (local $$t25.1 i32)
  (local $$t26.0 i32)
  (local $$t26.1 i32)
  (local $$t27 i32)
  (local $$t28.0 i32)
  (local $$t28.1 i32)
  (local $$t29.0 i32)
  (local $$t29.1 i32)
  (local $$t30.0 i32)
  (local $$t30.1 i32)
  (local $$t31.0 i32)
  (local $$t31.1 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_9
        block $$Block_8
          block $$Block_7
            block $$Block_6
              block $$Block_5
                block $$Block_4
                  block $$Block_3
                    block $$Block_2
                      block $$Block_1
                        block $$Block_0
                          block $$BlockSel
                            local.get $$block_selector
                            br_table 0 1 2 3 4 5 6 7 8 9 0
                          end ;;$BlockSel
                          i32.const 0
                          local.set $$current_block

                          ;;&x.Left [#2]
                          local.get $x.0
                          call $runtime.Block.Retain
                          local.get $x.1
                          i32.const 8
                          i32.add
                          local.set $$t0.1
                          local.get $$t0.0
                          call $runtime.Block.Release
                          local.set $$t0.0

                          ;;*t0
                          local.get $$t0.1
                          i32.load offset=0 align=4
                          call $runtime.Block.Retain
                          local.get $$t0.1
                          i32.load offset=4 align=4
                          local.set $$t1.1
                          local.get $$t1.0
                          call $runtime.Block.Release
                          local.set $$t1.0

                          ;;&this.NIL [#0]
                          local.get $this.0
                          call $runtime.Block.Retain
                          local.get $this.1
                          i32.const 0
                          i32.add
                          local.set $$t2.1
                          local.get $$t2.0
                          call $runtime.Block.Release
                          local.set $$t2.0

                          ;;*t2
                          local.get $$t2.1
                          i32.load offset=0 align=4
                          call $runtime.Block.Retain
                          local.get $$t2.1
                          i32.load offset=4 align=4
                          local.set $$t3.1
                          local.get $$t3.0
                          call $runtime.Block.Release
                          local.set $$t3.0

                          ;;t1 == t3
                          local.get $$t1.1
                          local.get $$t3.1
                          i32.eq
                          local.set $$t4

                          ;;if t4 goto 1 else 2
                          local.get $$t4
                          if
                            br $$Block_0
                          else
                            br $$Block_1
                          end

                        end ;;$Block_0
                        i32.const 1
                        local.set $$current_block

                        ;;return
                        br $$BlockFnBody

                      end ;;$Block_1
                      i32.const 2
                      local.set $$current_block

                      ;;&x.Left [#2]
                      local.get $x.0
                      call $runtime.Block.Retain
                      local.get $x.1
                      i32.const 8
                      i32.add
                      local.set $$t5.1
                      local.get $$t5.0
                      call $runtime.Block.Release
                      local.set $$t5.0

                      ;;*t5
                      local.get $$t5.1
                      i32.load offset=0 align=4
                      call $runtime.Block.Retain
                      local.get $$t5.1
                      i32.load offset=4 align=4
                      local.set $$t6.1
                      local.get $$t6.0
                      call $runtime.Block.Release
                      local.set $$t6.0

                      ;;&x.Left [#2]
                      local.get $x.0
                      call $runtime.Block.Retain
                      local.get $x.1
                      i32.const 8
                      i32.add
                      local.set $$t7.1
                      local.get $$t7.0
                      call $runtime.Block.Release
                      local.set $$t7.0

                      ;;&t6.Right [#3]
                      local.get $$t6.0
                      call $runtime.Block.Retain
                      local.get $$t6.1
                      i32.const 16
                      i32.add
                      local.set $$t8.1
                      local.get $$t8.0
                      call $runtime.Block.Release
                      local.set $$t8.0

                      ;;*t8
                      local.get $$t8.1
                      i32.load offset=0 align=4
                      call $runtime.Block.Retain
                      local.get $$t8.1
                      i32.load offset=4 align=4
                      local.set $$t9.1
                      local.get $$t9.0
                      call $runtime.Block.Release
                      local.set $$t9.0

                      ;;*t7 = t9
                      local.get $$t7.1
                      local.get $$t9.0
                      call $runtime.Block.Retain
                      local.get $$t7.1
                      i32.load offset=0 align=1
                      call $runtime.Block.Release
                      i32.store offset=0 align=1
                      local.get $$t7.1
                      local.get $$t9.1
                      i32.store offset=4 align=4

                      ;;&t6.Right [#3]
                      local.get $$t6.0
                      call $runtime.Block.Retain
                      local.get $$t6.1
                      i32.const 16
                      i32.add
                      local.set $$t10.1
                      local.get $$t10.0
                      call $runtime.Block.Release
                      local.set $$t10.0

                      ;;*t10
                      local.get $$t10.1
                      i32.load offset=0 align=4
                      call $runtime.Block.Retain
                      local.get $$t10.1
                      i32.load offset=4 align=4
                      local.set $$t11.1
                      local.get $$t11.0
                      call $runtime.Block.Release
                      local.set $$t11.0

                      ;;&this.NIL [#0]
                      local.get $this.0
                      call $runtime.Block.Retain
                      local.get $this.1
                      i32.const 0
                      i32.add
                      local.set $$t12.1
                      local.get $$t12.0
                      call $runtime.Block.Release
                      local.set $$t12.0

                      ;;*t12
                      local.get $$t12.1
                      i32.load offset=0 align=4
                      call $runtime.Block.Retain
                      local.get $$t12.1
                      i32.load offset=4 align=4
                      local.set $$t13.1
                      local.get $$t13.0
                      call $runtime.Block.Release
                      local.set $$t13.0

                      ;;t11 != t13
                      local.get $$t11.1
                      local.get $$t13.1
                      i32.eq
                      i32.eqz
                      local.set $$t14

                      ;;if t14 goto 3 else 4
                      local.get $$t14
                      if
                        br $$Block_2
                      else
                        br $$Block_3
                      end

                    end ;;$Block_2
                    i32.const 3
                    local.set $$current_block

                    ;;&t6.Right [#3]
                    local.get $$t6.0
                    call $runtime.Block.Retain
                    local.get $$t6.1
                    i32.const 16
                    i32.add
                    local.set $$t15.1
                    local.get $$t15.0
                    call $runtime.Block.Release
                    local.set $$t15.0

                    ;;*t15
                    local.get $$t15.1
                    i32.load offset=0 align=4
                    call $runtime.Block.Retain
                    local.get $$t15.1
                    i32.load offset=4 align=4
                    local.set $$t16.1
                    local.get $$t16.0
                    call $runtime.Block.Release
                    local.set $$t16.0

                    ;;(*mapNode).SetParent(t16, x)
                    local.get $$t16.0
                    local.get $$t16.1
                    local.get $x.0
                    local.get $x.1
                    call $runtime.mapNode.SetParent

                    ;;jump 4
                    br $$Block_3

                  end ;;$Block_3
                  i32.const 4
                  local.set $$current_block

                  ;;(*mapNode).Parent(x, this)
                  local.get $x.0
                  local.get $x.1
                  local.get $this.0
                  local.get $this.1
                  call $runtime.mapNode.Parent
                  local.set $$t17.1
                  local.get $$t17.0
                  call $runtime.Block.Release
                  local.set $$t17.0

                  ;;(*mapNode).SetParent(t6, t18)
                  local.get $$t6.0
                  local.get $$t6.1
                  local.get $$t17.0
                  local.get $$t17.1
                  call $runtime.mapNode.SetParent

                  ;;(*mapNode).Parent(x, this)
                  local.get $x.0
                  local.get $x.1
                  local.get $this.0
                  local.get $this.1
                  call $runtime.mapNode.Parent
                  local.set $$t18.1
                  local.get $$t18.0
                  call $runtime.Block.Release
                  local.set $$t18.0

                  ;;&this.NIL [#0]
                  local.get $this.0
                  call $runtime.Block.Retain
                  local.get $this.1
                  i32.const 0
                  i32.add
                  local.set $$t19.1
                  local.get $$t19.0
                  call $runtime.Block.Release
                  local.set $$t19.0

                  ;;*t21
                  local.get $$t19.1
                  i32.load offset=0 align=4
                  call $runtime.Block.Retain
                  local.get $$t19.1
                  i32.load offset=4 align=4
                  local.set $$t20.1
                  local.get $$t20.0
                  call $runtime.Block.Release
                  local.set $$t20.0

                  ;;t20 == t22
                  local.get $$t18.1
                  local.get $$t20.1
                  i32.eq
                  local.set $$t21

                  ;;if t23 goto 5 else 7
                  local.get $$t21
                  if
                    br $$Block_4
                  else
                    br $$Block_6
                  end

                end ;;$Block_4
                i32.const 5
                local.set $$current_block

                ;;&this.root [#1]
                local.get $this.0
                call $runtime.Block.Retain
                local.get $this.1
                i32.const 8
                i32.add
                local.set $$t22.1
                local.get $$t22.0
                call $runtime.Block.Release
                local.set $$t22.0

                ;;*t24 = t6
                local.get $$t22.1
                local.get $$t6.0
                call $runtime.Block.Retain
                local.get $$t22.1
                i32.load offset=0 align=1
                call $runtime.Block.Release
                i32.store offset=0 align=1
                local.get $$t22.1
                local.get $$t6.1
                i32.store offset=4 align=4

                ;;jump 6
                br $$Block_5

              end ;;$Block_5
              i32.const 6
              local.set $$current_block

              ;;&t6.Right [#3]
              local.get $$t6.0
              call $runtime.Block.Retain
              local.get $$t6.1
              i32.const 16
              i32.add
              local.set $$t23.1
              local.get $$t23.0
              call $runtime.Block.Release
              local.set $$t23.0

              ;;*t25 = x
              local.get $$t23.1
              local.get $x.0
              call $runtime.Block.Retain
              local.get $$t23.1
              i32.load offset=0 align=1
              call $runtime.Block.Release
              i32.store offset=0 align=1
              local.get $$t23.1
              local.get $x.1
              i32.store offset=4 align=4

              ;;(*mapNode).SetParent(x, t6)
              local.get $x.0
              local.get $x.1
              local.get $$t6.0
              local.get $$t6.1
              call $runtime.mapNode.SetParent

              ;;return
              br $$BlockFnBody

            end ;;$Block_6
            i32.const 7
            local.set $$current_block

            ;;(*mapNode).Parent(x, this)
            local.get $x.0
            local.get $x.1
            local.get $this.0
            local.get $this.1
            call $runtime.mapNode.Parent
            local.set $$t24.1
            local.get $$t24.0
            call $runtime.Block.Release
            local.set $$t24.0

            ;;&t27.Left [#2]
            local.get $$t24.0
            call $runtime.Block.Retain
            local.get $$t24.1
            i32.const 8
            i32.add
            local.set $$t25.1
            local.get $$t25.0
            call $runtime.Block.Release
            local.set $$t25.0

            ;;*t28
            local.get $$t25.1
            i32.load offset=0 align=4
            call $runtime.Block.Retain
            local.get $$t25.1
            i32.load offset=4 align=4
            local.set $$t26.1
            local.get $$t26.0
            call $runtime.Block.Release
            local.set $$t26.0

            ;;x == t29
            local.get $x.1
            local.get $$t26.1
            i32.eq
            local.set $$t27

            ;;if t30 goto 8 else 9
            local.get $$t27
            if
              br $$Block_7
            else
              br $$Block_8
            end

          end ;;$Block_7
          i32.const 8
          local.set $$current_block

          ;;(*mapNode).Parent(x, this)
          local.get $x.0
          local.get $x.1
          local.get $this.0
          local.get $this.1
          call $runtime.mapNode.Parent
          local.set $$t28.1
          local.get $$t28.0
          call $runtime.Block.Release
          local.set $$t28.0

          ;;&t31.Left [#2]
          local.get $$t28.0
          call $runtime.Block.Retain
          local.get $$t28.1
          i32.const 8
          i32.add
          local.set $$t29.1
          local.get $$t29.0
          call $runtime.Block.Release
          local.set $$t29.0

          ;;*t32 = t6
          local.get $$t29.1
          local.get $$t6.0
          call $runtime.Block.Retain
          local.get $$t29.1
          i32.load offset=0 align=1
          call $runtime.Block.Release
          i32.store offset=0 align=1
          local.get $$t29.1
          local.get $$t6.1
          i32.store offset=4 align=4

          ;;jump 6
          i32.const 6
          local.set $$block_selector
          br $$BlockDisp

        end ;;$Block_8
        i32.const 9
        local.set $$current_block

        ;;(*mapNode).Parent(x, this)
        local.get $x.0
        local.get $x.1
        local.get $this.0
        local.get $this.1
        call $runtime.mapNode.Parent
        local.set $$t30.1
        local.get $$t30.0
        call $runtime.Block.Release
        local.set $$t30.0

        ;;&t33.Right [#3]
        local.get $$t30.0
        call $runtime.Block.Retain
        local.get $$t30.1
        i32.const 16
        i32.add
        local.set $$t31.1
        local.get $$t31.0
        call $runtime.Block.Release
        local.set $$t31.0

        ;;*t34 = t6
        local.get $$t31.1
        local.get $$t6.0
        call $runtime.Block.Retain
        local.get $$t31.1
        i32.load offset=0 align=1
        call $runtime.Block.Release
        i32.store offset=0 align=1
        local.get $$t31.1
        local.get $$t6.1
        i32.store offset=4 align=4

        ;;jump 6
        i32.const 6
        local.set $$block_selector
        br $$BlockDisp

      end ;;$Block_9
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$t0.0
  call $runtime.Block.Release
  local.get $$t1.0
  call $runtime.Block.Release
  local.get $$t2.0
  call $runtime.Block.Release
  local.get $$t3.0
  call $runtime.Block.Release
  local.get $$t5.0
  call $runtime.Block.Release
  local.get $$t6.0
  call $runtime.Block.Release
  local.get $$t7.0
  call $runtime.Block.Release
  local.get $$t8.0
  call $runtime.Block.Release
  local.get $$t9.0
  call $runtime.Block.Release
  local.get $$t10.0
  call $runtime.Block.Release
  local.get $$t11.0
  call $runtime.Block.Release
  local.get $$t12.0
  call $runtime.Block.Release
  local.get $$t13.0
  call $runtime.Block.Release
  local.get $$t15.0
  call $runtime.Block.Release
  local.get $$t16.0
  call $runtime.Block.Release
  local.get $$t17.0
  call $runtime.Block.Release
  local.get $$t18.0
  call $runtime.Block.Release
  local.get $$t19.0
  call $runtime.Block.Release
  local.get $$t20.0
  call $runtime.Block.Release
  local.get $$t22.0
  call $runtime.Block.Release
  local.get $$t23.0
  call $runtime.Block.Release
  local.get $$t24.0
  call $runtime.Block.Release
  local.get $$t25.0
  call $runtime.Block.Release
  local.get $$t26.0
  call $runtime.Block.Release
  local.get $$t28.0
  call $runtime.Block.Release
  local.get $$t29.0
  call $runtime.Block.Release
  local.get $$t30.0
  call $runtime.Block.Release
  local.get $$t31.0
  call $runtime.Block.Release
) ;;runtime.mapImp.rightRotate

(func $runtime.mapImp.search (param $this.0 i32) (param $this.1 i32) (param $key.0.0 i32) (param $key.0.1 i32) (param $key.1 i32) (param $key.2 i32) (result i32 i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$ret_0.0 i32)
  (local $$ret_0.1 i32)
  (local $$t0.0 i32)
  (local $$t0.1 i32)
  (local $$t1.0 i32)
  (local $$t1.1 i32)
  (local $$t2.0 i32)
  (local $$t2.1 i32)
  (local $$t3.0 i32)
  (local $$t3.1 i32)
  (local $$t4.0.0 i32)
  (local $$t4.0.1 i32)
  (local $$t4.1 i32)
  (local $$t4.2 i32)
  (local $$t5 i32)
  (local $$t6 i32)
  (local $$t7.0 i32)
  (local $$t7.1 i32)
  (local $$t8.0 i32)
  (local $$t8.1 i32)
  (local $$t9.0 i32)
  (local $$t9.1 i32)
  (local $$t10.0 i32)
  (local $$t10.1 i32)
  (local $$t11 i32)
  (local $$t12.0 i32)
  (local $$t12.1 i32)
  (local $$t13 i32)
  (local $$t14.0 i32)
  (local $$t14.1 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_7
        block $$Block_6
          block $$Block_5
            block $$Block_4
              block $$Block_3
                block $$Block_2
                  block $$Block_1
                    block $$Block_0
                      block $$BlockSel
                        local.get $$block_selector
                        br_table 0 1 2 3 4 5 6 7 0
                      end ;;$BlockSel
                      i32.const 0
                      local.set $$current_block

                      ;;&this.root [#1]
                      local.get $this.0
                      call $runtime.Block.Retain
                      local.get $this.1
                      i32.const 8
                      i32.add
                      local.set $$t0.1
                      local.get $$t0.0
                      call $runtime.Block.Release
                      local.set $$t0.0

                      ;;*t0
                      local.get $$t0.1
                      i32.load offset=0 align=4
                      call $runtime.Block.Retain
                      local.get $$t0.1
                      i32.load offset=4 align=4
                      local.set $$t1.1
                      local.get $$t1.0
                      call $runtime.Block.Release
                      local.set $$t1.0

                      ;;jump 3
                      br $$Block_2

                    end ;;$Block_0
                    i32.const 1
                    local.set $$current_block

                    ;;&t6.Key [#5]
                    local.get $$t2.0
                    call $runtime.Block.Retain
                    local.get $$t2.1
                    i32.const 28
                    i32.add
                    local.set $$t3.1
                    local.get $$t3.0
                    call $runtime.Block.Release
                    local.set $$t3.0

                    ;;*t2
                    local.get $$t3.1
                    i32.load offset=0 align=4
                    call $runtime.Block.Retain
                    local.get $$t3.1
                    i32.load offset=4 align=4
                    local.get $$t3.1
                    i32.load offset=8 align=4
                    local.get $$t3.1
                    i32.load offset=12 align=4
                    local.set $$t4.2
                    local.set $$t4.1
                    local.set $$t4.0.1
                    local.get $$t4.0.0
                    call $runtime.Block.Release
                    local.set $$t4.0.0

                    ;;Compare(t3, key)
                    local.get $$t4.0.0
                    local.get $$t4.0.1
                    local.get $$t4.1
                    local.get $$t4.2
                    local.get $key.0.0
                    local.get $key.0.1
                    local.get $key.1
                    local.get $key.2
                    call $runtime.Compare
                    local.set $$t5

                    ;;t4 < 0:i32
                    local.get $$t5
                    i32.const 0
                    i32.lt_s
                    local.set $$t6

                    ;;if t5 goto 4 else 5
                    local.get $$t6
                    if
                      br $$Block_3
                    else
                      br $$Block_4
                    end

                  end ;;$Block_1
                  i32.const 2
                  local.set $$current_block

                  ;;return t6
                  local.get $$t2.0
                  call $runtime.Block.Retain
                  local.get $$t2.1
                  local.set $$ret_0.1
                  local.get $$ret_0.0
                  call $runtime.Block.Release
                  local.set $$ret_0.0
                  br $$BlockFnBody

                end ;;$Block_2
                local.get $$current_block
                i32.const 0
                i32.eq
                if (result i32 i32)
                  local.get $$t1.0
                  call $runtime.Block.Retain
                  local.get $$t1.1
                else
                  local.get $$current_block
                  i32.const 4
                  i32.eq
                  if (result i32 i32)
                    local.get $$t7.0
                    call $runtime.Block.Retain
                    local.get $$t7.1
                  else
                    local.get $$t8.0
                    call $runtime.Block.Retain
                    local.get $$t8.1
                  end
                end
                local.set $$t2.1
                local.get $$t2.0
                call $runtime.Block.Release
                local.set $$t2.0
                i32.const 3
                local.set $$current_block

                ;;&this.NIL [#0]
                local.get $this.0
                call $runtime.Block.Retain
                local.get $this.1
                i32.const 0
                i32.add
                local.set $$t9.1
                local.get $$t9.0
                call $runtime.Block.Release
                local.set $$t9.0

                ;;*t7
                local.get $$t9.1
                i32.load offset=0 align=4
                call $runtime.Block.Retain
                local.get $$t9.1
                i32.load offset=4 align=4
                local.set $$t10.1
                local.get $$t10.0
                call $runtime.Block.Release
                local.set $$t10.0

                ;;t6 != t8
                local.get $$t2.1
                local.get $$t10.1
                i32.eq
                i32.eqz
                local.set $$t11

                ;;if t9 goto 1 else 2
                local.get $$t11
                if
                  i32.const 1
                  local.set $$block_selector
                  br $$BlockDisp
                else
                  i32.const 2
                  local.set $$block_selector
                  br $$BlockDisp
                end

              end ;;$Block_3
              i32.const 4
              local.set $$current_block

              ;;&t6.Right [#3]
              local.get $$t2.0
              call $runtime.Block.Retain
              local.get $$t2.1
              i32.const 16
              i32.add
              local.set $$t12.1
              local.get $$t12.0
              call $runtime.Block.Release
              local.set $$t12.0

              ;;*t10
              local.get $$t12.1
              i32.load offset=0 align=4
              call $runtime.Block.Retain
              local.get $$t12.1
              i32.load offset=4 align=4
              local.set $$t7.1
              local.get $$t7.0
              call $runtime.Block.Release
              local.set $$t7.0

              ;;jump 3
              i32.const 3
              local.set $$block_selector
              br $$BlockDisp

            end ;;$Block_4
            i32.const 5
            local.set $$current_block

            ;;t4 > 0:i32
            local.get $$t5
            i32.const 0
            i32.gt_s
            local.set $$t13

            ;;if t12 goto 6 else 7
            local.get $$t13
            if
              br $$Block_5
            else
              br $$Block_6
            end

          end ;;$Block_5
          i32.const 6
          local.set $$current_block

          ;;&t6.Left [#2]
          local.get $$t2.0
          call $runtime.Block.Retain
          local.get $$t2.1
          i32.const 8
          i32.add
          local.set $$t14.1
          local.get $$t14.0
          call $runtime.Block.Release
          local.set $$t14.0

          ;;*t13
          local.get $$t14.1
          i32.load offset=0 align=4
          call $runtime.Block.Retain
          local.get $$t14.1
          i32.load offset=4 align=4
          local.set $$t8.1
          local.get $$t8.0
          call $runtime.Block.Release
          local.set $$t8.0

          ;;jump 3
          i32.const 3
          local.set $$block_selector
          br $$BlockDisp

        end ;;$Block_6
        i32.const 7
        local.set $$current_block

        ;;return t6
        local.get $$t2.0
        call $runtime.Block.Retain
        local.get $$t2.1
        local.set $$ret_0.1
        local.get $$ret_0.0
        call $runtime.Block.Release
        local.set $$ret_0.0
        br $$BlockFnBody

      end ;;$Block_7
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$ret_0.0
  call $runtime.Block.Retain
  local.get $$ret_0.1
  local.get $$ret_0.0
  call $runtime.Block.Release
  local.get $$t0.0
  call $runtime.Block.Release
  local.get $$t1.0
  call $runtime.Block.Release
  local.get $$t2.0
  call $runtime.Block.Release
  local.get $$t3.0
  call $runtime.Block.Release
  local.get $$t4.0.0
  call $runtime.Block.Release
  local.get $$t7.0
  call $runtime.Block.Release
  local.get $$t8.0
  call $runtime.Block.Release
  local.get $$t9.0
  call $runtime.Block.Release
  local.get $$t10.0
  call $runtime.Block.Release
  local.get $$t12.0
  call $runtime.Block.Release
  local.get $$t14.0
  call $runtime.Block.Release
) ;;runtime.mapImp.search

(func $runtime.mapImp.successor (param $this.0 i32) (param $this.1 i32) (param $x.0 i32) (param $x.1 i32) (result i32 i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$ret_0.0 i32)
  (local $$ret_0.1 i32)
  (local $$t0.0 i32)
  (local $$t0.1 i32)
  (local $$t1.0 i32)
  (local $$t1.1 i32)
  (local $$t2 i32)
  (local $$t3.0 i32)
  (local $$t3.1 i32)
  (local $$t4.0 i32)
  (local $$t4.1 i32)
  (local $$t5.0 i32)
  (local $$t5.1 i32)
  (local $$t6.0 i32)
  (local $$t6.1 i32)
  (local $$t7.0 i32)
  (local $$t7.1 i32)
  (local $$t8.0 i32)
  (local $$t8.1 i32)
  (local $$t9 i32)
  (local $$t10.0 i32)
  (local $$t10.1 i32)
  (local $$t11.0 i32)
  (local $$t11.1 i32)
  (local $$t12.0 i32)
  (local $$t12.1 i32)
  (local $$t13.0 i32)
  (local $$t13.1 i32)
  (local $$t14.0 i32)
  (local $$t14.1 i32)
  (local $$t15.0 i32)
  (local $$t15.1 i32)
  (local $$t16.0 i32)
  (local $$t16.1 i32)
  (local $$t17.0 i32)
  (local $$t17.1 i32)
  (local $$t18.0 i32)
  (local $$t18.1 i32)
  (local $$t19 i32)
  (local $$t20.0 i32)
  (local $$t20.1 i32)
  (local $$t21.0 i32)
  (local $$t21.1 i32)
  (local $$t22 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_8
        block $$Block_7
          block $$Block_6
            block $$Block_5
              block $$Block_4
                block $$Block_3
                  block $$Block_2
                    block $$Block_1
                      block $$Block_0
                        block $$BlockSel
                          local.get $$block_selector
                          br_table 0 1 2 3 4 5 6 7 8 0
                        end ;;$BlockSel
                        i32.const 0
                        local.set $$current_block

                        ;;&this.NIL [#0]
                        local.get $this.0
                        call $runtime.Block.Retain
                        local.get $this.1
                        i32.const 0
                        i32.add
                        local.set $$t0.1
                        local.get $$t0.0
                        call $runtime.Block.Release
                        local.set $$t0.0

                        ;;*t0
                        local.get $$t0.1
                        i32.load offset=0 align=4
                        call $runtime.Block.Retain
                        local.get $$t0.1
                        i32.load offset=4 align=4
                        local.set $$t1.1
                        local.get $$t1.0
                        call $runtime.Block.Release
                        local.set $$t1.0

                        ;;x == t1
                        local.get $x.1
                        local.get $$t1.1
                        i32.eq
                        local.set $$t2

                        ;;if t2 goto 1 else 2
                        local.get $$t2
                        if
                          br $$Block_0
                        else
                          br $$Block_1
                        end

                      end ;;$Block_0
                      i32.const 1
                      local.set $$current_block

                      ;;&this.NIL [#0]
                      local.get $this.0
                      call $runtime.Block.Retain
                      local.get $this.1
                      i32.const 0
                      i32.add
                      local.set $$t3.1
                      local.get $$t3.0
                      call $runtime.Block.Release
                      local.set $$t3.0

                      ;;*t3
                      local.get $$t3.1
                      i32.load offset=0 align=4
                      call $runtime.Block.Retain
                      local.get $$t3.1
                      i32.load offset=4 align=4
                      local.set $$t4.1
                      local.get $$t4.0
                      call $runtime.Block.Release
                      local.set $$t4.0

                      ;;return t4
                      local.get $$t4.0
                      call $runtime.Block.Retain
                      local.get $$t4.1
                      local.set $$ret_0.1
                      local.get $$ret_0.0
                      call $runtime.Block.Release
                      local.set $$ret_0.0
                      br $$BlockFnBody

                    end ;;$Block_1
                    i32.const 2
                    local.set $$current_block

                    ;;&x.Right [#3]
                    local.get $x.0
                    call $runtime.Block.Retain
                    local.get $x.1
                    i32.const 16
                    i32.add
                    local.set $$t5.1
                    local.get $$t5.0
                    call $runtime.Block.Release
                    local.set $$t5.0

                    ;;*t5
                    local.get $$t5.1
                    i32.load offset=0 align=4
                    call $runtime.Block.Retain
                    local.get $$t5.1
                    i32.load offset=4 align=4
                    local.set $$t6.1
                    local.get $$t6.0
                    call $runtime.Block.Release
                    local.set $$t6.0

                    ;;&this.NIL [#0]
                    local.get $this.0
                    call $runtime.Block.Retain
                    local.get $this.1
                    i32.const 0
                    i32.add
                    local.set $$t7.1
                    local.get $$t7.0
                    call $runtime.Block.Release
                    local.set $$t7.0

                    ;;*t7
                    local.get $$t7.1
                    i32.load offset=0 align=4
                    call $runtime.Block.Retain
                    local.get $$t7.1
                    i32.load offset=4 align=4
                    local.set $$t8.1
                    local.get $$t8.0
                    call $runtime.Block.Release
                    local.set $$t8.0

                    ;;t6 != t8
                    local.get $$t6.1
                    local.get $$t8.1
                    i32.eq
                    i32.eqz
                    local.set $$t9

                    ;;if t9 goto 3 else 4
                    local.get $$t9
                    if
                      br $$Block_2
                    else
                      br $$Block_3
                    end

                  end ;;$Block_2
                  i32.const 3
                  local.set $$current_block

                  ;;&x.Right [#3]
                  local.get $x.0
                  call $runtime.Block.Retain
                  local.get $x.1
                  i32.const 16
                  i32.add
                  local.set $$t10.1
                  local.get $$t10.0
                  call $runtime.Block.Release
                  local.set $$t10.0

                  ;;*t10
                  local.get $$t10.1
                  i32.load offset=0 align=4
                  call $runtime.Block.Retain
                  local.get $$t10.1
                  i32.load offset=4 align=4
                  local.set $$t11.1
                  local.get $$t11.0
                  call $runtime.Block.Release
                  local.set $$t11.0

                  ;;(*mapImp).min(this, t11)
                  local.get $this.0
                  local.get $this.1
                  local.get $$t11.0
                  local.get $$t11.1
                  call $runtime.mapImp.min
                  local.set $$t12.1
                  local.get $$t12.0
                  call $runtime.Block.Release
                  local.set $$t12.0

                  ;;return t12
                  local.get $$t12.0
                  call $runtime.Block.Retain
                  local.get $$t12.1
                  local.set $$ret_0.1
                  local.get $$ret_0.0
                  call $runtime.Block.Release
                  local.set $$ret_0.0
                  br $$BlockFnBody

                end ;;$Block_3
                i32.const 4
                local.set $$current_block

                ;;(*mapNode).Parent(x, this)
                local.get $x.0
                local.get $x.1
                local.get $this.0
                local.get $this.1
                call $runtime.mapNode.Parent
                local.set $$t13.1
                local.get $$t13.0
                call $runtime.Block.Release
                local.set $$t13.0

                ;;jump 7
                br $$Block_6

              end ;;$Block_4
              i32.const 5
              local.set $$current_block

              ;;(*mapNode).Parent(t16, this)
              local.get $$t14.0
              local.get $$t14.1
              local.get $this.0
              local.get $this.1
              call $runtime.mapNode.Parent
              local.set $$t15.1
              local.get $$t15.0
              call $runtime.Block.Release
              local.set $$t15.0

              ;;jump 7
              br $$Block_6

            end ;;$Block_5
            i32.const 6
            local.set $$current_block

            ;;return t16
            local.get $$t14.0
            call $runtime.Block.Retain
            local.get $$t14.1
            local.set $$ret_0.1
            local.get $$ret_0.0
            call $runtime.Block.Release
            local.set $$ret_0.0
            br $$BlockFnBody

          end ;;$Block_6
          local.get $$current_block
          i32.const 4
          i32.eq
          if (result i32 i32)
            local.get $x.0
            call $runtime.Block.Retain
            local.get $x.1
          else
            local.get $$t14.0
            call $runtime.Block.Retain
            local.get $$t14.1
          end
          local.get $$current_block
          i32.const 4
          i32.eq
          if (result i32 i32)
            local.get $$t13.0
            call $runtime.Block.Retain
            local.get $$t13.1
          else
            local.get $$t15.0
            call $runtime.Block.Retain
            local.get $$t15.1
          end
          local.set $$t14.1
          local.get $$t14.0
          call $runtime.Block.Release
          local.set $$t14.0
          local.set $$t16.1
          local.get $$t16.0
          call $runtime.Block.Release
          local.set $$t16.0
          i32.const 7
          local.set $$current_block

          ;;&this.NIL [#0]
          local.get $this.0
          call $runtime.Block.Retain
          local.get $this.1
          i32.const 0
          i32.add
          local.set $$t17.1
          local.get $$t17.0
          call $runtime.Block.Release
          local.set $$t17.0

          ;;*t17
          local.get $$t17.1
          i32.load offset=0 align=4
          call $runtime.Block.Retain
          local.get $$t17.1
          i32.load offset=4 align=4
          local.set $$t18.1
          local.get $$t18.0
          call $runtime.Block.Release
          local.set $$t18.0

          ;;t16 != t18
          local.get $$t14.1
          local.get $$t18.1
          i32.eq
          i32.eqz
          local.set $$t19

          ;;if t19 goto 8 else 6
          local.get $$t19
          if
            br $$Block_7
          else
            i32.const 6
            local.set $$block_selector
            br $$BlockDisp
          end

        end ;;$Block_7
        i32.const 8
        local.set $$current_block

        ;;&t16.Right [#3]
        local.get $$t14.0
        call $runtime.Block.Retain
        local.get $$t14.1
        i32.const 16
        i32.add
        local.set $$t20.1
        local.get $$t20.0
        call $runtime.Block.Release
        local.set $$t20.0

        ;;*t20
        local.get $$t20.1
        i32.load offset=0 align=4
        call $runtime.Block.Retain
        local.get $$t20.1
        i32.load offset=4 align=4
        local.set $$t21.1
        local.get $$t21.0
        call $runtime.Block.Release
        local.set $$t21.0

        ;;t15 == t21
        local.get $$t16.1
        local.get $$t21.1
        i32.eq
        local.set $$t22

        ;;if t22 goto 5 else 6
        local.get $$t22
        if
          i32.const 5
          local.set $$block_selector
          br $$BlockDisp
        else
          i32.const 6
          local.set $$block_selector
          br $$BlockDisp
        end

      end ;;$Block_8
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$ret_0.0
  call $runtime.Block.Retain
  local.get $$ret_0.1
  local.get $$ret_0.0
  call $runtime.Block.Release
  local.get $$t0.0
  call $runtime.Block.Release
  local.get $$t1.0
  call $runtime.Block.Release
  local.get $$t3.0
  call $runtime.Block.Release
  local.get $$t4.0
  call $runtime.Block.Release
  local.get $$t5.0
  call $runtime.Block.Release
  local.get $$t6.0
  call $runtime.Block.Release
  local.get $$t7.0
  call $runtime.Block.Release
  local.get $$t8.0
  call $runtime.Block.Release
  local.get $$t10.0
  call $runtime.Block.Release
  local.get $$t11.0
  call $runtime.Block.Release
  local.get $$t12.0
  call $runtime.Block.Release
  local.get $$t13.0
  call $runtime.Block.Release
  local.get $$t14.0
  call $runtime.Block.Release
  local.get $$t15.0
  call $runtime.Block.Release
  local.get $$t16.0
  call $runtime.Block.Release
  local.get $$t17.0
  call $runtime.Block.Release
  local.get $$t18.0
  call $runtime.Block.Release
  local.get $$t20.0
  call $runtime.Block.Release
  local.get $$t21.0
  call $runtime.Block.Release
) ;;runtime.mapImp.successor

(func $runtime.mapNode.SetParent (param $this.0 i32) (param $this.1 i32) (param $x.0 i32) (param $x.1 i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$t0.0 i32)
  (local $$t0.1 i32)
  (local $$t1.0 i32)
  (local $$t1.1 i32)
  (local $$t2 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;&this.parentIdx [#0]
        local.get $this.0
        call $runtime.Block.Retain
        local.get $this.1
        i32.const 0
        i32.add
        local.set $$t0.1
        local.get $$t0.0
        call $runtime.Block.Release
        local.set $$t0.0

        ;;&x.NodeIdx [#1]
        local.get $x.0
        call $runtime.Block.Retain
        local.get $x.1
        i32.const 4
        i32.add
        local.set $$t1.1
        local.get $$t1.0
        call $runtime.Block.Release
        local.set $$t1.0

        ;;*t1
        local.get $$t1.1
        i32.load offset=0 align=4
        local.set $$t2

        ;;*t0 = t2
        local.get $$t0.1
        local.get $$t2
        i32.store offset=0 align=4

        ;;return
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$t0.0
  call $runtime.Block.Release
  local.get $$t1.0
  call $runtime.Block.Release
) ;;runtime.mapNode.SetParent

(func $runtime.mapIter.Next (param $this.0 i32) (param $this.1 i32) (result i32 i32 i32 i32 i32 i32 i32 i32 i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$ret_0 i32)
  (local $$ret_1.0.0 i32)
  (local $$ret_1.0.1 i32)
  (local $$ret_1.1 i32)
  (local $$ret_1.2 i32)
  (local $$ret_2.0.0 i32)
  (local $$ret_2.0.1 i32)
  (local $$ret_2.1 i32)
  (local $$ret_2.2 i32)
  (local $$t0.0 i32)
  (local $$t0.1 i32)
  (local $$t1 i32)
  (local $$t2.0 i32)
  (local $$t2.1 i32)
  (local $$t3.0 i32)
  (local $$t3.1 i32)
  (local $$t4 i32)
  (local $$t5 i32)
  (local $$t6.0 i32)
  (local $$t6.1 i32)
  (local $$t7 i32)
  (local $$t8 i32)
  (local $$t9.0 i32)
  (local $$t9.1 i32)
  (local $$t10.0 i32)
  (local $$t10.1 i32)
  (local $$t11.0 i32)
  (local $$t11.1 i32)
  (local $$t12.0 i32)
  (local $$t12.1 i32)
  (local $$t12.2 i32)
  (local $$t12.3 i32)
  (local $$t13.0 i32)
  (local $$t13.1 i32)
  (local $$t14 i32)
  (local $$t15.0 i32)
  (local $$t15.1 i32)
  (local $$t16.0 i32)
  (local $$t16.1 i32)
  (local $$t17.0 i32)
  (local $$t17.1 i32)
  (local $$t18.0.0 i32)
  (local $$t18.0.1 i32)
  (local $$t18.1 i32)
  (local $$t18.2 i32)
  (local $$t19.0 i32)
  (local $$t19.1 i32)
  (local $$t20.0.0 i32)
  (local $$t20.0.1 i32)
  (local $$t20.1 i32)
  (local $$t20.2 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_2
        block $$Block_1
          block $$Block_0
            block $$BlockSel
              local.get $$block_selector
              br_table 0 1 2 0
            end ;;$BlockSel
            i32.const 0
            local.set $$current_block

            ;;&this.pos [#1]
            local.get $this.0
            call $runtime.Block.Retain
            local.get $this.1
            i32.const 8
            i32.add
            local.set $$t0.1
            local.get $$t0.0
            call $runtime.Block.Release
            local.set $$t0.0

            ;;*t0
            local.get $$t0.1
            i32.load offset=0 align=4
            local.set $$t1

            ;;&this.m [#0]
            local.get $this.0
            call $runtime.Block.Retain
            local.get $this.1
            i32.const 0
            i32.add
            local.set $$t2.1
            local.get $$t2.0
            call $runtime.Block.Release
            local.set $$t2.0

            ;;*t2
            local.get $$t2.1
            i32.load offset=0 align=4
            call $runtime.Block.Retain
            local.get $$t2.1
            i32.load offset=4 align=4
            local.set $$t3.1
            local.get $$t3.0
            call $runtime.Block.Release
            local.set $$t3.0

            ;;(*mapImp).Len(t3)
            local.get $$t3.0
            local.get $$t3.1
            call $runtime.mapImp.Len
            local.set $$t4

            ;;t1 >= t4
            local.get $$t1
            local.get $$t4
            i32.ge_s
            local.set $$t5

            ;;if t5 goto 1 else 2
            local.get $$t5
            if
              br $$Block_0
            else
              br $$Block_1
            end

          end ;;$Block_0
          i32.const 1
          local.set $$current_block

          ;;return false:bool, nil:interface{}, nil:interface{}
          i32.const 0
          local.set $$ret_0
          i32.const 0
          i32.const 0
          i32.const 0
          i32.const 0
          local.set $$ret_1.2
          local.set $$ret_1.1
          local.set $$ret_1.0.1
          local.get $$ret_1.0.0
          call $runtime.Block.Release
          local.set $$ret_1.0.0
          i32.const 0
          i32.const 0
          i32.const 0
          i32.const 0
          local.set $$ret_2.2
          local.set $$ret_2.1
          local.set $$ret_2.0.1
          local.get $$ret_2.0.0
          call $runtime.Block.Release
          local.set $$ret_2.0.0
          br $$BlockFnBody

        end ;;$Block_1
        i32.const 2
        local.set $$current_block

        ;;&this.pos [#1]
        local.get $this.0
        call $runtime.Block.Retain
        local.get $this.1
        i32.const 8
        i32.add
        local.set $$t6.1
        local.get $$t6.0
        call $runtime.Block.Release
        local.set $$t6.0

        ;;*t6
        local.get $$t6.1
        i32.load offset=0 align=4
        local.set $$t7

        ;;t7 + 1:int
        local.get $$t7
        i32.const 1
        i32.add
        local.set $$t8

        ;;*t6 = t8
        local.get $$t6.1
        local.get $$t8
        i32.store offset=0 align=4

        ;;&this.m [#0]
        local.get $this.0
        call $runtime.Block.Retain
        local.get $this.1
        i32.const 0
        i32.add
        local.set $$t9.1
        local.get $$t9.0
        call $runtime.Block.Release
        local.set $$t9.0

        ;;*t9
        local.get $$t9.1
        i32.load offset=0 align=4
        call $runtime.Block.Retain
        local.get $$t9.1
        i32.load offset=4 align=4
        local.set $$t10.1
        local.get $$t10.0
        call $runtime.Block.Release
        local.set $$t10.0

        ;;&t10.nodes [#2]
        local.get $$t10.0
        call $runtime.Block.Retain
        local.get $$t10.1
        i32.const 16
        i32.add
        local.set $$t11.1
        local.get $$t11.0
        call $runtime.Block.Release
        local.set $$t11.0

        ;;*t11
        local.get $$t11.1
        i32.load offset=0 align=4
        call $runtime.Block.Retain
        local.get $$t11.1
        i32.load offset=4 align=4
        local.get $$t11.1
        i32.load offset=8 align=4
        local.get $$t11.1
        i32.load offset=12 align=4
        local.set $$t12.3
        local.set $$t12.2
        local.set $$t12.1
        local.get $$t12.0
        call $runtime.Block.Release
        local.set $$t12.0

        ;;&this.pos [#1]
        local.get $this.0
        call $runtime.Block.Retain
        local.get $this.1
        i32.const 8
        i32.add
        local.set $$t13.1
        local.get $$t13.0
        call $runtime.Block.Release
        local.set $$t13.0

        ;;*t13
        local.get $$t13.1
        i32.load offset=0 align=4
        local.set $$t14

        ;;&t12[t14]
        local.get $$t12.0
        call $runtime.Block.Retain
        local.get $$t12.1
        i32.const 8
        local.get $$t14
        i32.mul
        i32.add
        local.set $$t15.1
        local.get $$t15.0
        call $runtime.Block.Release
        local.set $$t15.0

        ;;*t15
        local.get $$t15.1
        i32.load offset=0 align=4
        call $runtime.Block.Retain
        local.get $$t15.1
        i32.load offset=4 align=4
        local.set $$t16.1
        local.get $$t16.0
        call $runtime.Block.Release
        local.set $$t16.0

        ;;&t16.Key [#5]
        local.get $$t16.0
        call $runtime.Block.Retain
        local.get $$t16.1
        i32.const 28
        i32.add
        local.set $$t17.1
        local.get $$t17.0
        call $runtime.Block.Release
        local.set $$t17.0

        ;;*t17
        local.get $$t17.1
        i32.load offset=0 align=4
        call $runtime.Block.Retain
        local.get $$t17.1
        i32.load offset=4 align=4
        local.get $$t17.1
        i32.load offset=8 align=4
        local.get $$t17.1
        i32.load offset=12 align=4
        local.set $$t18.2
        local.set $$t18.1
        local.set $$t18.0.1
        local.get $$t18.0.0
        call $runtime.Block.Release
        local.set $$t18.0.0

        ;;&t16.Val [#6]
        local.get $$t16.0
        call $runtime.Block.Retain
        local.get $$t16.1
        i32.const 44
        i32.add
        local.set $$t19.1
        local.get $$t19.0
        call $runtime.Block.Release
        local.set $$t19.0

        ;;*t19
        local.get $$t19.1
        i32.load offset=0 align=4
        call $runtime.Block.Retain
        local.get $$t19.1
        i32.load offset=4 align=4
        local.get $$t19.1
        i32.load offset=8 align=4
        local.get $$t19.1
        i32.load offset=12 align=4
        local.set $$t20.2
        local.set $$t20.1
        local.set $$t20.0.1
        local.get $$t20.0.0
        call $runtime.Block.Release
        local.set $$t20.0.0

        ;;return true:bool, t18, t20
        i32.const 1
        local.set $$ret_0
        local.get $$t18.0.0
        call $runtime.Block.Retain
        local.get $$t18.0.1
        local.get $$t18.1
        local.get $$t18.2
        local.set $$ret_1.2
        local.set $$ret_1.1
        local.set $$ret_1.0.1
        local.get $$ret_1.0.0
        call $runtime.Block.Release
        local.set $$ret_1.0.0
        local.get $$t20.0.0
        call $runtime.Block.Retain
        local.get $$t20.0.1
        local.get $$t20.1
        local.get $$t20.2
        local.set $$ret_2.2
        local.set $$ret_2.1
        local.set $$ret_2.0.1
        local.get $$ret_2.0.0
        call $runtime.Block.Release
        local.set $$ret_2.0.0
        br $$BlockFnBody

      end ;;$Block_2
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$ret_0
  local.get $$ret_1.0.0
  call $runtime.Block.Retain
  local.get $$ret_1.0.1
  local.get $$ret_1.1
  local.get $$ret_1.2
  local.get $$ret_2.0.0
  call $runtime.Block.Retain
  local.get $$ret_2.0.1
  local.get $$ret_2.1
  local.get $$ret_2.2
  local.get $$ret_1.0.0
  call $runtime.Block.Release
  local.get $$ret_2.0.0
  call $runtime.Block.Release
  local.get $$t0.0
  call $runtime.Block.Release
  local.get $$t2.0
  call $runtime.Block.Release
  local.get $$t3.0
  call $runtime.Block.Release
  local.get $$t6.0
  call $runtime.Block.Release
  local.get $$t9.0
  call $runtime.Block.Release
  local.get $$t10.0
  call $runtime.Block.Release
  local.get $$t11.0
  call $runtime.Block.Release
  local.get $$t12.0
  call $runtime.Block.Release
  local.get $$t13.0
  call $runtime.Block.Release
  local.get $$t15.0
  call $runtime.Block.Release
  local.get $$t16.0
  call $runtime.Block.Release
  local.get $$t17.0
  call $runtime.Block.Release
  local.get $$t18.0.0
  call $runtime.Block.Release
  local.get $$t19.0
  call $runtime.Block.Release
  local.get $$t20.0.0
  call $runtime.Block.Release
) ;;runtime.mapIter.Next

(func $$u8.$slice.append (param $x.0 i32) (param $x.1 i32) (param $x.2 i32) (param $x.3 i32) (param $y.0 i32) (param $y.1 i32) (param $y.2 i32) (param $y.3 i32) (result i32 i32 i32 i32)
  (local $item i32)
  (local $x_len i32)
  (local $y_len i32)
  (local $new_len i32)
  (local $src i32)
  (local $dest i32)
  (local $new_cap i32)
  local.get $x.2
  local.set $x_len
  local.get $y.2
  local.set $y_len
  local.get $x_len
  local.get $y_len
  i32.add
  local.set $new_len
  local.get $new_len
  local.get $x.3
  i32.le_u
  if (result i32 i32 i32 i32)
    local.get $x.0
    call $runtime.Block.Retain
    local.get $x.1
    local.get $new_len
    local.get $x.3
    local.get $y.1
    local.set $src
    local.get $x.1
    i32.const 1
    local.get $x_len
    i32.mul
    i32.add
    local.set $dest
    block $block1
      loop $loop1
        local.get $y_len
        i32.eqz
        if
          br $block1
        else
        end
        local.get $src
        i32.load8_u offset=0 align=1
        local.set $item
        local.get $dest
        local.get $item
        i32.store8 offset=0 align=1
        local.get $src
        i32.const 1
        i32.add
        local.set $src
        local.get $dest
        i32.const 1
        i32.add
        local.set $dest
        local.get $y_len
        i32.const 1
        i32.sub
        local.set $y_len
        br $loop1
      end ;;loop1
    end ;;block1
  else
    local.get $new_len
    i32.const 2
    i32.mul
    local.set $new_cap
    local.get $new_cap
    i32.const 1
    i32.mul
    i32.const 16
    i32.add
    call $runtime.HeapAlloc
    local.get $new_cap
    i32.const 0
    i32.const 1
    call $runtime.Block.Init
    call $runtime.DupI32
    i32.const 16
    i32.add
    call $runtime.DupI32
    local.set $dest
    local.get $new_len
    local.get $new_cap
    local.get $x.1
    local.set $src
    block $block2
      loop $loop2
        local.get $x_len
        i32.eqz
        if
          br $block2
        else
        end
        local.get $src
        i32.load8_u offset=0 align=1
        local.set $item
        local.get $dest
        local.get $item
        i32.store8 offset=0 align=1
        local.get $src
        i32.const 1
        i32.add
        local.set $src
        local.get $dest
        i32.const 1
        i32.add
        local.set $dest
        local.get $x_len
        i32.const 1
        i32.sub
        local.set $x_len
        br $loop2
      end ;;loop2
    end ;;block2
    local.get $y.1
    local.set $src
    block $block3
      loop $loop3
        local.get $y_len
        i32.eqz
        if
          br $block3
        else
        end
        local.get $src
        i32.load8_u offset=0 align=1
        local.set $item
        local.get $dest
        local.get $item
        i32.store8 offset=0 align=1
        local.get $src
        i32.const 1
        i32.add
        local.set $src
        local.get $dest
        i32.const 1
        i32.add
        local.set $dest
        local.get $y_len
        i32.const 1
        i32.sub
        local.set $y_len
        br $loop3
      end ;;loop3
    end ;;block3
  end
) ;;$u8.$slice.append

(func $brainfuck$bfpkg.BrainFuck.Run (param $this.0 i32) (param $this.1 i32) (result i32 i32 i32 i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$ret_0.0 i32)
  (local $$ret_0.1 i32)
  (local $$ret_0.2 i32)
  (local $$ret_0.3 i32)
  (local $$t0.0 i32)
  (local $$t0.1 i32)
  (local $$t1.0 i32)
  (local $$t1.1 i32)
  (local $$t1.2 i32)
  (local $$t1.3 i32)
  (local $$t2.0 i32)
  (local $$t2.1 i32)
  (local $$t3.0 i32)
  (local $$t3.1 i32)
  (local $$t3.2 i32)
  (local $$t4.0 i32)
  (local $$t4.1 i32)
  (local $$t5 i32)
  (local $$t6 i32)
  (local $$t7 i32)
  (local $$t8.0 i32)
  (local $$t8.1 i32)
  (local $$t8.2 i32)
  (local $$t8.3 i32)
  (local $$t9.0 i32)
  (local $$t9.1 i32)
  (local $$t9.2 i32)
  (local $$t9.3 i32)
  (local $$t10.0 i32)
  (local $$t10.1 i32)
  (local $$t11 i32)
  (local $$t12.0 i32)
  (local $$t12.1 i32)
  (local $$t13.0 i32)
  (local $$t13.1 i32)
  (local $$t13.2 i32)
  (local $$t14 i32)
  (local $$t15 i32)
  (local $$t16.0 i32)
  (local $$t16.1 i32)
  (local $$t16.2 i32)
  (local $$t16.3 i32)
  (local $$t17.0 i32)
  (local $$t17.1 i32)
  (local $$t18 i32)
  (local $$t19 i32)
  (local $$t20.0 i32)
  (local $$t20.1 i32)
  (local $$t21 i32)
  (local $$t22 i32)
  (local $$t23.0 i32)
  (local $$t23.1 i32)
  (local $$t24 i32)
  (local $$t25 i32)
  (local $$t26 i32)
  (local $$t27.0 i32)
  (local $$t27.1 i32)
  (local $$t28.0 i32)
  (local $$t28.1 i32)
  (local $$t29 i32)
  (local $$t30.0 i32)
  (local $$t30.1 i32)
  (local $$t31 i32)
  (local $$t32 i32)
  (local $$t33 i32)
  (local $$t34.0 i32)
  (local $$t34.1 i32)
  (local $$t35.0 i32)
  (local $$t35.1 i32)
  (local $$t36 i32)
  (local $$t37.0 i32)
  (local $$t37.1 i32)
  (local $$t38 i32)
  (local $$t39 i32)
  (local $$t40 i32)
  (local $$t41.0 i32)
  (local $$t41.1 i32)
  (local $$t42.0 i32)
  (local $$t42.1 i32)
  (local $$t43 i32)
  (local $$t44.0 i32)
  (local $$t44.1 i32)
  (local $$t45 i32)
  (local $$t46 i32)
  (local $$t47 i32)
  (local $$t48.0 i32)
  (local $$t48.1 i32)
  (local $$t49.0 i32)
  (local $$t49.1 i32)
  (local $$t50 i32)
  (local $$t51.0 i32)
  (local $$t51.1 i32)
  (local $$t52 i32)
  (local $$t53 i32)
  (local $$t54 i32)
  (local $$t55.0 i32)
  (local $$t55.1 i32)
  (local $$t56.0 i32)
  (local $$t56.1 i32)
  (local $$t57 i32)
  (local $$t58.0 i32)
  (local $$t58.1 i32)
  (local $$t59 i32)
  (local $$t60.0 i32)
  (local $$t60.1 i32)
  (local $$t61.0 i32)
  (local $$t61.1 i32)
  (local $$t62.0 i32)
  (local $$t62.1 i32)
  (local $$t62.2 i32)
  (local $$t62.3 i32)
  (local $$t63 i32)
  (local $$t64 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_20
        block $$Block_19
          block $$Block_18
            block $$Block_17
              block $$Block_16
                block $$Block_15
                  block $$Block_14
                    block $$Block_13
                      block $$Block_12
                        block $$Block_11
                          block $$Block_10
                            block $$Block_9
                              block $$Block_8
                                block $$Block_7
                                  block $$Block_6
                                    block $$Block_5
                                      block $$Block_4
                                        block $$Block_3
                                          block $$Block_2
                                            block $$Block_1
                                              block $$Block_0
                                                block $$BlockSel
                                                  local.get $$block_selector
                                                  br_table 0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 0
                                                end ;;$BlockSel
                                                i32.const 0
                                                local.set $$current_block

                                                ;;new [64]byte (makeslice)
                                                i32.const 80
                                                call $runtime.HeapAlloc
                                                i32.const 1
                                                i32.const 0
                                                i32.const 64
                                                call $runtime.Block.Init
                                                call $runtime.DupI32
                                                i32.const 16
                                                i32.add
                                                local.set $$t0.1
                                                local.get $$t0.0
                                                call $runtime.Block.Release
                                                local.set $$t0.0

                                                ;;slice t0[:0:int]
                                                local.get $$t0.0
                                                call $runtime.Block.Retain
                                                local.get $$t0.1
                                                i32.const 1
                                                i32.const 0
                                                i32.mul
                                                i32.add
                                                i32.const 0
                                                i32.const 0
                                                i32.sub
                                                i32.const 64
                                                i32.const 0
                                                i32.sub
                                                local.set $$t1.3
                                                local.set $$t1.2
                                                local.set $$t1.1
                                                local.get $$t1.0
                                                call $runtime.Block.Release
                                                local.set $$t1.0

                                                ;;jump 3
                                                br $$Block_2

                                              end ;;$Block_0
                                              i32.const 1
                                              local.set $$current_block

                                              ;;&this.code [#1]
                                              local.get $this.0
                                              call $runtime.Block.Retain
                                              local.get $this.1
                                              i32.const 30000
                                              i32.add
                                              local.set $$t2.1
                                              local.get $$t2.0
                                              call $runtime.Block.Release
                                              local.set $$t2.0

                                              ;;*t2
                                              local.get $$t2.1
                                              i32.load offset=0 align=4
                                              call $runtime.Block.Retain
                                              local.get $$t2.1
                                              i32.load offset=4 align=4
                                              local.get $$t2.1
                                              i32.load offset=8 align=4
                                              local.set $$t3.2
                                              local.set $$t3.1
                                              local.get $$t3.0
                                              call $runtime.Block.Release
                                              local.set $$t3.0

                                              ;;&this.pc [#3]
                                              local.get $this.0
                                              call $runtime.Block.Retain
                                              local.get $this.1
                                              i32.const 30016
                                              i32.add
                                              local.set $$t4.1
                                              local.get $$t4.0
                                              call $runtime.Block.Release
                                              local.set $$t4.0

                                              ;;*t4
                                              local.get $$t4.1
                                              i32.load offset=0 align=4
                                              local.set $$t5

                                              ;;t3[t5]
                                              local.get $$t3.1
                                              local.get $$t5
                                              i32.add
                                              i32.load8_u offset=0 align=1
                                              local.set $$t6

                                              ;;t6 == 62:byte
                                              local.get $$t6
                                              i32.const 62
                                              i32.eq
                                              local.set $$t7

                                              ;;if t7 goto 5 else 7
                                              local.get $$t7
                                              if
                                                br $$Block_4
                                              else
                                                br $$Block_6
                                              end

                                            end ;;$Block_1
                                            i32.const 2
                                            local.set $$current_block

                                            ;;return t8
                                            local.get $$t8.0
                                            call $runtime.Block.Retain
                                            local.get $$t8.1
                                            local.get $$t8.2
                                            local.get $$t8.3
                                            local.set $$ret_0.3
                                            local.set $$ret_0.2
                                            local.set $$ret_0.1
                                            local.get $$ret_0.0
                                            call $runtime.Block.Release
                                            local.set $$ret_0.0
                                            br $$BlockFnBody

                                          end ;;$Block_2
                                          local.get $$current_block
                                          i32.const 0
                                          i32.eq
                                          if (result i32 i32 i32 i32)
                                            local.get $$t1.0
                                            call $runtime.Block.Retain
                                            local.get $$t1.1
                                            local.get $$t1.2
                                            local.get $$t1.3
                                          else
                                            local.get $$t9.0
                                            call $runtime.Block.Retain
                                            local.get $$t9.1
                                            local.get $$t9.2
                                            local.get $$t9.3
                                          end
                                          local.set $$t8.3
                                          local.set $$t8.2
                                          local.set $$t8.1
                                          local.get $$t8.0
                                          call $runtime.Block.Release
                                          local.set $$t8.0
                                          i32.const 3
                                          local.set $$current_block

                                          ;;&this.pc [#3]
                                          local.get $this.0
                                          call $runtime.Block.Retain
                                          local.get $this.1
                                          i32.const 30016
                                          i32.add
                                          local.set $$t10.1
                                          local.get $$t10.0
                                          call $runtime.Block.Release
                                          local.set $$t10.0

                                          ;;*t9
                                          local.get $$t10.1
                                          i32.load offset=0 align=4
                                          local.set $$t11

                                          ;;&this.code [#1]
                                          local.get $this.0
                                          call $runtime.Block.Retain
                                          local.get $this.1
                                          i32.const 30000
                                          i32.add
                                          local.set $$t12.1
                                          local.get $$t12.0
                                          call $runtime.Block.Release
                                          local.set $$t12.0

                                          ;;*t11
                                          local.get $$t12.1
                                          i32.load offset=0 align=4
                                          call $runtime.Block.Retain
                                          local.get $$t12.1
                                          i32.load offset=4 align=4
                                          local.get $$t12.1
                                          i32.load offset=8 align=4
                                          local.set $$t13.2
                                          local.set $$t13.1
                                          local.get $$t13.0
                                          call $runtime.Block.Release
                                          local.set $$t13.0

                                          ;;len(t12)
                                          local.get $$t13.2
                                          local.set $$t14

                                          ;;t10 != t13
                                          local.get $$t11
                                          local.get $$t14
                                          i32.eq
                                          i32.eqz
                                          local.set $$t15

                                          ;;if t14 goto 1 else 2
                                          local.get $$t15
                                          if
                                            i32.const 1
                                            local.set $$block_selector
                                            br $$BlockDisp
                                          else
                                            i32.const 2
                                            local.set $$block_selector
                                            br $$BlockDisp
                                          end

                                        end ;;$Block_3
                                        local.get $$current_block
                                        i32.const 5
                                        i32.eq
                                        if (result i32 i32 i32 i32)
                                          local.get $$t8.0
                                          call $runtime.Block.Retain
                                          local.get $$t8.1
                                          local.get $$t8.2
                                          local.get $$t8.3
                                        else
                                          local.get $$current_block
                                          i32.const 6
                                          i32.eq
                                          if (result i32 i32 i32 i32)
                                            local.get $$t8.0
                                            call $runtime.Block.Retain
                                            local.get $$t8.1
                                            local.get $$t8.2
                                            local.get $$t8.3
                                          else
                                            local.get $$current_block
                                            i32.const 8
                                            i32.eq
                                            if (result i32 i32 i32 i32)
                                              local.get $$t8.0
                                              call $runtime.Block.Retain
                                              local.get $$t8.1
                                              local.get $$t8.2
                                              local.get $$t8.3
                                            else
                                              local.get $$current_block
                                              i32.const 10
                                              i32.eq
                                              if (result i32 i32 i32 i32)
                                                local.get $$t8.0
                                                call $runtime.Block.Retain
                                                local.get $$t8.1
                                                local.get $$t8.2
                                                local.get $$t8.3
                                              else
                                                local.get $$current_block
                                                i32.const 12
                                                i32.eq
                                                if (result i32 i32 i32 i32)
                                                  local.get $$t8.0
                                                  call $runtime.Block.Retain
                                                  local.get $$t8.1
                                                  local.get $$t8.2
                                                  local.get $$t8.3
                                                else
                                                  local.get $$current_block
                                                  i32.const 14
                                                  i32.eq
                                                  if (result i32 i32 i32 i32)
                                                    local.get $$t8.0
                                                    call $runtime.Block.Retain
                                                    local.get $$t8.1
                                                    local.get $$t8.2
                                                    local.get $$t8.3
                                                  else
                                                    local.get $$current_block
                                                    i32.const 17
                                                    i32.eq
                                                    if (result i32 i32 i32 i32)
                                                      local.get $$t16.0
                                                      call $runtime.Block.Retain
                                                      local.get $$t16.1
                                                      local.get $$t16.2
                                                      local.get $$t16.3
                                                    else
                                                      local.get $$current_block
                                                      i32.const 20
                                                      i32.eq
                                                      if (result i32 i32 i32 i32)
                                                        local.get $$t8.0
                                                        call $runtime.Block.Retain
                                                        local.get $$t8.1
                                                        local.get $$t8.2
                                                        local.get $$t8.3
                                                      else
                                                        local.get $$current_block
                                                        i32.const 16
                                                        i32.eq
                                                        if (result i32 i32 i32 i32)
                                                          local.get $$t8.0
                                                          call $runtime.Block.Retain
                                                          local.get $$t8.1
                                                          local.get $$t8.2
                                                          local.get $$t8.3
                                                        else
                                                          local.get $$t8.0
                                                          call $runtime.Block.Retain
                                                          local.get $$t8.1
                                                          local.get $$t8.2
                                                          local.get $$t8.3
                                                        end
                                                      end
                                                    end
                                                  end
                                                end
                                              end
                                            end
                                          end
                                        end
                                        local.set $$t9.3
                                        local.set $$t9.2
                                        local.set $$t9.1
                                        local.get $$t9.0
                                        call $runtime.Block.Release
                                        local.set $$t9.0
                                        i32.const 4
                                        local.set $$current_block

                                        ;;&this.pc [#3]
                                        local.get $this.0
                                        call $runtime.Block.Retain
                                        local.get $this.1
                                        i32.const 30016
                                        i32.add
                                        local.set $$t17.1
                                        local.get $$t17.0
                                        call $runtime.Block.Release
                                        local.set $$t17.0

                                        ;;*t16
                                        local.get $$t17.1
                                        i32.load offset=0 align=4
                                        local.set $$t18

                                        ;;t17 + 1:int
                                        local.get $$t18
                                        i32.const 1
                                        i32.add
                                        local.set $$t19

                                        ;;*t16 = t18
                                        local.get $$t17.1
                                        local.get $$t19
                                        i32.store offset=0 align=4

                                        ;;jump 3
                                        i32.const 3
                                        local.set $$block_selector
                                        br $$BlockDisp

                                      end ;;$Block_4
                                      i32.const 5
                                      local.set $$current_block

                                      ;;&this.pos [#2]
                                      local.get $this.0
                                      call $runtime.Block.Retain
                                      local.get $this.1
                                      i32.const 30012
                                      i32.add
                                      local.set $$t20.1
                                      local.get $$t20.0
                                      call $runtime.Block.Release
                                      local.set $$t20.0

                                      ;;*t19
                                      local.get $$t20.1
                                      i32.load offset=0 align=4
                                      local.set $$t21

                                      ;;t20 + 1:int
                                      local.get $$t21
                                      i32.const 1
                                      i32.add
                                      local.set $$t22

                                      ;;*t19 = t21
                                      local.get $$t20.1
                                      local.get $$t22
                                      i32.store offset=0 align=4

                                      ;;jump 4
                                      i32.const 4
                                      local.set $$block_selector
                                      br $$BlockDisp

                                    end ;;$Block_5
                                    i32.const 6
                                    local.set $$current_block

                                    ;;&this.pos [#2]
                                    local.get $this.0
                                    call $runtime.Block.Retain
                                    local.get $this.1
                                    i32.const 30012
                                    i32.add
                                    local.set $$t23.1
                                    local.get $$t23.0
                                    call $runtime.Block.Release
                                    local.set $$t23.0

                                    ;;*t22
                                    local.get $$t23.1
                                    i32.load offset=0 align=4
                                    local.set $$t24

                                    ;;t23 - 1:int
                                    local.get $$t24
                                    i32.const 1
                                    i32.sub
                                    local.set $$t25

                                    ;;*t22 = t24
                                    local.get $$t23.1
                                    local.get $$t25
                                    i32.store offset=0 align=4

                                    ;;jump 4
                                    i32.const 4
                                    local.set $$block_selector
                                    br $$BlockDisp

                                  end ;;$Block_6
                                  i32.const 7
                                  local.set $$current_block

                                  ;;t6 == 60:byte
                                  local.get $$t6
                                  i32.const 60
                                  i32.eq
                                  local.set $$t26

                                  ;;if t25 goto 6 else 9
                                  local.get $$t26
                                  if
                                    i32.const 6
                                    local.set $$block_selector
                                    br $$BlockDisp
                                  else
                                    br $$Block_8
                                  end

                                end ;;$Block_7
                                i32.const 8
                                local.set $$current_block

                                ;;&this.mem [#0]
                                local.get $this.0
                                call $runtime.Block.Retain
                                local.get $this.1
                                i32.const 0
                                i32.add
                                local.set $$t27.1
                                local.get $$t27.0
                                call $runtime.Block.Release
                                local.set $$t27.0

                                ;;&this.pos [#2]
                                local.get $this.0
                                call $runtime.Block.Retain
                                local.get $this.1
                                i32.const 30012
                                i32.add
                                local.set $$t28.1
                                local.get $$t28.0
                                call $runtime.Block.Release
                                local.set $$t28.0

                                ;;*t27
                                local.get $$t28.1
                                i32.load offset=0 align=4
                                local.set $$t29

                                ;;&t26[t28]
                                local.get $$t27.0
                                call $runtime.Block.Retain
                                local.get $$t27.1
                                i32.const 1
                                local.get $$t29
                                i32.mul
                                i32.add
                                local.set $$t30.1
                                local.get $$t30.0
                                call $runtime.Block.Release
                                local.set $$t30.0

                                ;;*t29
                                local.get $$t30.1
                                i32.load8_u offset=0 align=1
                                local.set $$t31

                                ;;t30 + 1:byte
                                local.get $$t31
                                i32.const 1
                                i32.add
                                i32.const 255
                                i32.and
                                local.set $$t32

                                ;;*t29 = t31
                                local.get $$t30.1
                                local.get $$t32
                                i32.store8 offset=0 align=1

                                ;;jump 4
                                i32.const 4
                                local.set $$block_selector
                                br $$BlockDisp

                              end ;;$Block_8
                              i32.const 9
                              local.set $$current_block

                              ;;t6 == 43:byte
                              local.get $$t6
                              i32.const 43
                              i32.eq
                              local.set $$t33

                              ;;if t32 goto 8 else 11
                              local.get $$t33
                              if
                                i32.const 8
                                local.set $$block_selector
                                br $$BlockDisp
                              else
                                br $$Block_10
                              end

                            end ;;$Block_9
                            i32.const 10
                            local.set $$current_block

                            ;;&this.mem [#0]
                            local.get $this.0
                            call $runtime.Block.Retain
                            local.get $this.1
                            i32.const 0
                            i32.add
                            local.set $$t34.1
                            local.get $$t34.0
                            call $runtime.Block.Release
                            local.set $$t34.0

                            ;;&this.pos [#2]
                            local.get $this.0
                            call $runtime.Block.Retain
                            local.get $this.1
                            i32.const 30012
                            i32.add
                            local.set $$t35.1
                            local.get $$t35.0
                            call $runtime.Block.Release
                            local.set $$t35.0

                            ;;*t34
                            local.get $$t35.1
                            i32.load offset=0 align=4
                            local.set $$t36

                            ;;&t33[t35]
                            local.get $$t34.0
                            call $runtime.Block.Retain
                            local.get $$t34.1
                            i32.const 1
                            local.get $$t36
                            i32.mul
                            i32.add
                            local.set $$t37.1
                            local.get $$t37.0
                            call $runtime.Block.Release
                            local.set $$t37.0

                            ;;*t36
                            local.get $$t37.1
                            i32.load8_u offset=0 align=1
                            local.set $$t38

                            ;;t37 - 1:byte
                            local.get $$t38
                            i32.const 1
                            i32.sub
                            i32.const 255
                            i32.and
                            local.set $$t39

                            ;;*t36 = t38
                            local.get $$t37.1
                            local.get $$t39
                            i32.store8 offset=0 align=1

                            ;;jump 4
                            i32.const 4
                            local.set $$block_selector
                            br $$BlockDisp

                          end ;;$Block_10
                          i32.const 11
                          local.set $$current_block

                          ;;t6 == 45:byte
                          local.get $$t6
                          i32.const 45
                          i32.eq
                          local.set $$t40

                          ;;if t39 goto 10 else 13
                          local.get $$t40
                          if
                            i32.const 10
                            local.set $$block_selector
                            br $$BlockDisp
                          else
                            br $$Block_12
                          end

                        end ;;$Block_11
                        i32.const 12
                        local.set $$current_block

                        ;;&this.mem [#0]
                        local.get $this.0
                        call $runtime.Block.Retain
                        local.get $this.1
                        i32.const 0
                        i32.add
                        local.set $$t41.1
                        local.get $$t41.0
                        call $runtime.Block.Release
                        local.set $$t41.0

                        ;;&this.pos [#2]
                        local.get $this.0
                        call $runtime.Block.Retain
                        local.get $this.1
                        i32.const 30012
                        i32.add
                        local.set $$t42.1
                        local.get $$t42.0
                        call $runtime.Block.Release
                        local.set $$t42.0

                        ;;*t41
                        local.get $$t42.1
                        i32.load offset=0 align=4
                        local.set $$t43

                        ;;&t40[t42]
                        local.get $$t41.0
                        call $runtime.Block.Retain
                        local.get $$t41.1
                        i32.const 1
                        local.get $$t43
                        i32.mul
                        i32.add
                        local.set $$t44.1
                        local.get $$t44.0
                        call $runtime.Block.Release
                        local.set $$t44.0

                        ;;*t43
                        local.get $$t44.1
                        i32.load8_u offset=0 align=1
                        local.set $$t45

                        ;;t44 == 0:byte
                        local.get $$t45
                        i32.const 0
                        i32.eq
                        local.set $$t46

                        ;;if t45 goto 16 else 4
                        local.get $$t46
                        if
                          br $$Block_15
                        else
                          i32.const 4
                          local.set $$block_selector
                          br $$BlockDisp
                        end

                      end ;;$Block_12
                      i32.const 13
                      local.set $$current_block

                      ;;t6 == 91:byte
                      local.get $$t6
                      i32.const 91
                      i32.eq
                      local.set $$t47

                      ;;if t46 goto 12 else 15
                      local.get $$t47
                      if
                        i32.const 12
                        local.set $$block_selector
                        br $$BlockDisp
                      else
                        br $$Block_14
                      end

                    end ;;$Block_13
                    i32.const 14
                    local.set $$current_block

                    ;;&this.mem [#0]
                    local.get $this.0
                    call $runtime.Block.Retain
                    local.get $this.1
                    i32.const 0
                    i32.add
                    local.set $$t48.1
                    local.get $$t48.0
                    call $runtime.Block.Release
                    local.set $$t48.0

                    ;;&this.pos [#2]
                    local.get $this.0
                    call $runtime.Block.Retain
                    local.get $this.1
                    i32.const 30012
                    i32.add
                    local.set $$t49.1
                    local.get $$t49.0
                    call $runtime.Block.Release
                    local.set $$t49.0

                    ;;*t48
                    local.get $$t49.1
                    i32.load offset=0 align=4
                    local.set $$t50

                    ;;&t47[t49]
                    local.get $$t48.0
                    call $runtime.Block.Retain
                    local.get $$t48.1
                    i32.const 1
                    local.get $$t50
                    i32.mul
                    i32.add
                    local.set $$t51.1
                    local.get $$t51.0
                    call $runtime.Block.Release
                    local.set $$t51.0

                    ;;*t50
                    local.get $$t51.1
                    i32.load8_u offset=0 align=1
                    local.set $$t52

                    ;;t51 != 0:byte
                    local.get $$t52
                    i32.const 0
                    i32.eq
                    i32.eqz
                    local.set $$t53

                    ;;if t52 goto 19 else 4
                    local.get $$t53
                    if
                      br $$Block_18
                    else
                      i32.const 4
                      local.set $$block_selector
                      br $$BlockDisp
                    end

                  end ;;$Block_14
                  i32.const 15
                  local.set $$current_block

                  ;;t6 == 93:byte
                  local.get $$t6
                  i32.const 93
                  i32.eq
                  local.set $$t54

                  ;;if t53 goto 14 else 18
                  local.get $$t54
                  if
                    i32.const 14
                    local.set $$block_selector
                    br $$BlockDisp
                  else
                    br $$Block_17
                  end

                end ;;$Block_15
                i32.const 16
                local.set $$current_block

                ;;(*BrainFuck).loop(this, 1:int)
                local.get $this.0
                local.get $this.1
                i32.const 1
                call $brainfuck$bfpkg.BrainFuck.loop

                ;;jump 4
                i32.const 4
                local.set $$block_selector
                br $$BlockDisp

              end ;;$Block_16
              i32.const 17
              local.set $$current_block

              ;;&this.mem [#0]
              local.get $this.0
              call $runtime.Block.Retain
              local.get $this.1
              i32.const 0
              i32.add
              local.set $$t55.1
              local.get $$t55.0
              call $runtime.Block.Release
              local.set $$t55.0

              ;;&this.pos [#2]
              local.get $this.0
              call $runtime.Block.Retain
              local.get $this.1
              i32.const 30012
              i32.add
              local.set $$t56.1
              local.get $$t56.0
              call $runtime.Block.Release
              local.set $$t56.0

              ;;*t56
              local.get $$t56.1
              i32.load offset=0 align=4
              local.set $$t57

              ;;&t55[t57]
              local.get $$t55.0
              call $runtime.Block.Retain
              local.get $$t55.1
              i32.const 1
              local.get $$t57
              i32.mul
              i32.add
              local.set $$t58.1
              local.get $$t58.0
              call $runtime.Block.Release
              local.set $$t58.0

              ;;*t58
              local.get $$t58.1
              i32.load8_u offset=0 align=1
              local.set $$t59

              ;;new [1]byte (varargs)
              i32.const 17
              call $runtime.HeapAlloc
              i32.const 1
              i32.const 0
              i32.const 1
              call $runtime.Block.Init
              call $runtime.DupI32
              i32.const 16
              i32.add
              local.set $$t60.1
              local.get $$t60.0
              call $runtime.Block.Release
              local.set $$t60.0

              ;;&t60[0:int]
              local.get $$t60.0
              call $runtime.Block.Retain
              local.get $$t60.1
              i32.const 1
              i32.const 0
              i32.mul
              i32.add
              local.set $$t61.1
              local.get $$t61.0
              call $runtime.Block.Release
              local.set $$t61.0

              ;;*t61 = t59
              local.get $$t61.1
              local.get $$t59
              i32.store8 offset=0 align=1

              ;;slice t60[:]
              local.get $$t60.0
              call $runtime.Block.Retain
              local.get $$t60.1
              i32.const 1
              i32.const 0
              i32.mul
              i32.add
              i32.const 1
              i32.const 0
              i32.sub
              i32.const 1
              i32.const 0
              i32.sub
              local.set $$t62.3
              local.set $$t62.2
              local.set $$t62.1
              local.get $$t62.0
              call $runtime.Block.Release
              local.set $$t62.0

              ;;append(t8, t62...)
              local.get $$t8.0
              local.get $$t8.1
              local.get $$t8.2
              local.get $$t8.3
              local.get $$t62.0
              local.get $$t62.1
              local.get $$t62.2
              local.get $$t62.3
              call $$u8.$slice.append
              local.set $$t16.3
              local.set $$t16.2
              local.set $$t16.1
              local.get $$t16.0
              call $runtime.Block.Release
              local.set $$t16.0

              ;;jump 4
              i32.const 4
              local.set $$block_selector
              br $$BlockDisp

            end ;;$Block_17
            i32.const 18
            local.set $$current_block

            ;;t6 == 46:byte
            local.get $$t6
            i32.const 46
            i32.eq
            local.set $$t63

            ;;if t64 goto 17 else 20
            local.get $$t63
            if
              i32.const 17
              local.set $$block_selector
              br $$BlockDisp
            else
              br $$Block_19
            end

          end ;;$Block_18
          i32.const 19
          local.set $$current_block

          ;;(*BrainFuck).loop(this, -1:int)
          local.get $this.0
          local.get $this.1
          i32.const -1
          call $brainfuck$bfpkg.BrainFuck.loop

          ;;jump 4
          i32.const 4
          local.set $$block_selector
          br $$BlockDisp

        end ;;$Block_19
        i32.const 20
        local.set $$current_block

        ;;t6 == 44:byte
        local.get $$t6
        i32.const 44
        i32.eq
        local.set $$t64

        ;;jump 4
        i32.const 4
        local.set $$block_selector
        br $$BlockDisp

      end ;;$Block_20
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$ret_0.0
  call $runtime.Block.Retain
  local.get $$ret_0.1
  local.get $$ret_0.2
  local.get $$ret_0.3
  local.get $$ret_0.0
  call $runtime.Block.Release
  local.get $$t0.0
  call $runtime.Block.Release
  local.get $$t1.0
  call $runtime.Block.Release
  local.get $$t2.0
  call $runtime.Block.Release
  local.get $$t3.0
  call $runtime.Block.Release
  local.get $$t4.0
  call $runtime.Block.Release
  local.get $$t8.0
  call $runtime.Block.Release
  local.get $$t9.0
  call $runtime.Block.Release
  local.get $$t10.0
  call $runtime.Block.Release
  local.get $$t12.0
  call $runtime.Block.Release
  local.get $$t13.0
  call $runtime.Block.Release
  local.get $$t16.0
  call $runtime.Block.Release
  local.get $$t17.0
  call $runtime.Block.Release
  local.get $$t20.0
  call $runtime.Block.Release
  local.get $$t23.0
  call $runtime.Block.Release
  local.get $$t27.0
  call $runtime.Block.Release
  local.get $$t28.0
  call $runtime.Block.Release
  local.get $$t30.0
  call $runtime.Block.Release
  local.get $$t34.0
  call $runtime.Block.Release
  local.get $$t35.0
  call $runtime.Block.Release
  local.get $$t37.0
  call $runtime.Block.Release
  local.get $$t41.0
  call $runtime.Block.Release
  local.get $$t42.0
  call $runtime.Block.Release
  local.get $$t44.0
  call $runtime.Block.Release
  local.get $$t48.0
  call $runtime.Block.Release
  local.get $$t49.0
  call $runtime.Block.Release
  local.get $$t51.0
  call $runtime.Block.Release
  local.get $$t55.0
  call $runtime.Block.Release
  local.get $$t56.0
  call $runtime.Block.Release
  local.get $$t58.0
  call $runtime.Block.Release
  local.get $$t60.0
  call $runtime.Block.Release
  local.get $$t61.0
  call $runtime.Block.Release
  local.get $$t62.0
  call $runtime.Block.Release
) ;;brainfuck$bfpkg.BrainFuck.Run

(func $brainfuck$bfpkg.BrainFuck.loop (param $this.0 i32) (param $this.1 i32) (param $inc i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$t0.0 i32)
  (local $$t0.1 i32)
  (local $$t1.0 i32)
  (local $$t1.1 i32)
  (local $$t1.2 i32)
  (local $$t2.0 i32)
  (local $$t2.1 i32)
  (local $$t3 i32)
  (local $$t4 i32)
  (local $$t5 i32)
  (local $$t6 i32)
  (local $$t7 i32)
  (local $$t8 i32)
  (local $$t9 i32)
  (local $$t10 i32)
  (local $$t11 i32)
  (local $$t12.0 i32)
  (local $$t12.1 i32)
  (local $$t13 i32)
  (local $$t14 i32)
  (local $$t15 i32)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_7
        block $$Block_6
          block $$Block_5
            block $$Block_4
              block $$Block_3
                block $$Block_2
                  block $$Block_1
                    block $$Block_0
                      block $$BlockSel
                        local.get $$block_selector
                        br_table 0 1 2 3 4 5 6 7 0
                      end ;;$BlockSel
                      i32.const 0
                      local.set $$current_block

                      ;;jump 3
                      br $$Block_2

                    end ;;$Block_0
                    i32.const 1
                    local.set $$current_block

                    ;;&this.code [#1]
                    local.get $this.0
                    call $runtime.Block.Retain
                    local.get $this.1
                    i32.const 30000
                    i32.add
                    local.set $$t0.1
                    local.get $$t0.0
                    call $runtime.Block.Release
                    local.set $$t0.0

                    ;;*t0
                    local.get $$t0.1
                    i32.load offset=0 align=4
                    call $runtime.Block.Retain
                    local.get $$t0.1
                    i32.load offset=4 align=4
                    local.get $$t0.1
                    i32.load offset=8 align=4
                    local.set $$t1.2
                    local.set $$t1.1
                    local.get $$t1.0
                    call $runtime.Block.Release
                    local.set $$t1.0

                    ;;&this.pc [#3]
                    local.get $this.0
                    call $runtime.Block.Retain
                    local.get $this.1
                    i32.const 30016
                    i32.add
                    local.set $$t2.1
                    local.get $$t2.0
                    call $runtime.Block.Release
                    local.set $$t2.0

                    ;;*t2
                    local.get $$t2.1
                    i32.load offset=0 align=4
                    local.set $$t3

                    ;;t3 + inc
                    local.get $$t3
                    local.get $inc
                    i32.add
                    local.set $$t4

                    ;;t1[t4]
                    local.get $$t1.1
                    local.get $$t4
                    i32.add
                    i32.load8_u offset=0 align=1
                    local.set $$t5

                    ;;t5 == 91:byte
                    local.get $$t5
                    i32.const 91
                    i32.eq
                    local.set $$t6

                    ;;if t6 goto 5 else 7
                    local.get $$t6
                    if
                      br $$Block_4
                    else
                      br $$Block_6
                    end

                  end ;;$Block_1
                  i32.const 2
                  local.set $$current_block

                  ;;return
                  br $$BlockFnBody

                end ;;$Block_2
                local.get $$current_block
                i32.const 0
                i32.eq
                if (result i32)
                  local.get $inc
                else
                  local.get $$t7
                end
                local.set $$t8
                i32.const 3
                local.set $$current_block

                ;;t7 != 0:int
                local.get $$t8
                i32.const 0
                i32.eq
                i32.eqz
                local.set $$t9

                ;;if t8 goto 1 else 2
                local.get $$t9
                if
                  i32.const 1
                  local.set $$block_selector
                  br $$BlockDisp
                else
                  i32.const 2
                  local.set $$block_selector
                  br $$BlockDisp
                end

              end ;;$Block_3
              local.get $$current_block
              i32.const 5
              i32.eq
              if (result i32)
                local.get $$t10
              else
                local.get $$current_block
                i32.const 6
                i32.eq
                if (result i32)
                  local.get $$t11
                else
                  local.get $$t8
                end
              end
              local.set $$t7
              i32.const 4
              local.set $$current_block

              ;;&this.pc [#3]
              local.get $this.0
              call $runtime.Block.Retain
              local.get $this.1
              i32.const 30016
              i32.add
              local.set $$t12.1
              local.get $$t12.0
              call $runtime.Block.Release
              local.set $$t12.0

              ;;*t10
              local.get $$t12.1
              i32.load offset=0 align=4
              local.set $$t13

              ;;t11 + inc
              local.get $$t13
              local.get $inc
              i32.add
              local.set $$t14

              ;;*t10 = t12
              local.get $$t12.1
              local.get $$t14
              i32.store offset=0 align=4

              ;;jump 3
              i32.const 3
              local.set $$block_selector
              br $$BlockDisp

            end ;;$Block_4
            i32.const 5
            local.set $$current_block

            ;;t7 + 1:int
            local.get $$t8
            i32.const 1
            i32.add
            local.set $$t10

            ;;jump 4
            i32.const 4
            local.set $$block_selector
            br $$BlockDisp

          end ;;$Block_5
          i32.const 6
          local.set $$current_block

          ;;t7 - 1:int
          local.get $$t8
          i32.const 1
          i32.sub
          local.set $$t11

          ;;jump 4
          i32.const 4
          local.set $$block_selector
          br $$BlockDisp

        end ;;$Block_6
        i32.const 7
        local.set $$current_block

        ;;t5 == 93:byte
        local.get $$t5
        i32.const 93
        i32.eq
        local.set $$t15

        ;;if t15 goto 6 else 4
        local.get $$t15
        if
          i32.const 6
          local.set $$block_selector
          br $$BlockDisp
        else
          i32.const 4
          local.set $$block_selector
          br $$BlockDisp
        end

      end ;;$Block_7
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$t0.0
  call $runtime.Block.Release
  local.get $$t1.0
  call $runtime.Block.Release
  local.get $$t2.0
  call $runtime.Block.Release
  local.get $$t12.0
  call $runtime.Block.Release
) ;;brainfuck$bfpkg.BrainFuck.loop

(func $_start (export "_start")
  call $brainfuck.init
) ;;_start

(func $_main (export "_main")
  call $brainfuck.main
) ;;_main
) ;;module