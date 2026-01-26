;; Copyright (C) 2026 武汉凹语言科技有限公司
;; SPDX-License-Identifier: AGPL-3.0-or-later

(module $fib_01
    (import "env" "print_i64" (func $env.print_i64 (param i64)))

    (memory 1)

    (func $main (export "_start")
        i64.const 123
        call $env.print_i64
    )

    (global $__heap_l128_freep (mut i32) (i32.const 0))   ;; l128 是循环链表, 记录当前的迭代位置

	(func $heap_block.size (param $ptr i32) (result i32)
		local.get $ptr
		i32.load offset=0
	)
	(func $heap_block.set_size (param $ptr i32) (param $size i32)
		local.get $ptr
		local.get $size
		i32.store offset=0
	)
	(func $heap_block.next (param $ptr i32) (result i32)
		local.get $ptr
		i32.load offset=4
	)
	(func $heap_block.set_next (param $ptr i32) (param $next i32)
		local.get $ptr
		local.get $next
		i32.store offset=4
	)
	(func $heap_block.data (param $ptr i32) (result i32)
		local.get $ptr
		i32.const 8
		i32.add
	)
    (func $heap_block.init (param $ptr i32) (param $size i32) (param $next i32)
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
)
