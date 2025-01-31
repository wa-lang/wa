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
