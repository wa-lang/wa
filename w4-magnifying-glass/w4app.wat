(module $__walang__
	(import "env" "memory" (memory 1 2))
	(table 26 funcref)
	(type $$OnFree (func (param i32)))
	(type $$wa.runtime.comp (func (param i32) (param i32) (result i32)))
	(type $$$fnSig1 (func))
	(global $__stack_ptr (mut i32) (i32.const 14656))
	(global $__heap_lfixed_cap i32 (i32.const 0))
	(global $__heap_ptr (mut i32) (i32.const 0))
	(global $__heap_top (mut i32) (i32.const 0))
	(global $__heap_l128_freep (mut i32) (i32.const 0))
	(global $__heap_init_flag (mut i32) (i32.const 0))
	(global $$wa.runtime.closure_data (mut i32) (i32.const 0))
	(global $$wa.runtime._concretTypeCount (mut i32) (i32.const 1))
	(global $$wa.runtime._interfaceCount (mut i32) (i32.const 1))
	(global $$wa.runtime._itabsPtr (mut i32) (i32.const 15112))
	(global $runtime.defersStack.0 i32 (i32.const 0))
	(global $runtime.defersStack.1 i32 (i32.const 14792))
	(global $runtime.init$guard (mut i32) (i32.const 0))
	(global $$knr_basep (mut i32) (i32.const 0))
	(global $$knr_freep (mut i32) (i32.const 0))
	(global $math._cos.0 i32 (i32.const 0))
	(global $math._cos.1 i32 (i32.const 14808))
	(global $math._sin.0 i32 (i32.const 0))
	(global $math._sin.1 i32 (i32.const 14856))
	(global $math.init$guard (mut i32) (i32.const 0))
	(global $syscall$wasm4.init$guard (mut i32) (i32.const 0))
	(global $w4app.cells.0 i32 (i32.const 0))
	(global $w4app.cells.1 i32 (i32.const 14904))
	(global $w4app.cellsFrame.0 i32 (i32.const 0))
	(global $w4app.cellsFrame.1 i32 (i32.const 14912))
	(global $w4app.center_x.0 i32 (i32.const 0))
	(global $w4app.center_x.1 i32 (i32.const 14920))
	(global $w4app.center_y.0 i32 (i32.const 0))
	(global $w4app.center_y.1 i32 (i32.const 14924))
	(global $w4app.height.0 i32 (i32.const 0))
	(global $w4app.height.1 i32 (i32.const 14928))
	(global $w4app.init$guard (mut i32) (i32.const 0))
	(global $w4app.width.0 i32 (i32.const 0))
	(global $w4app.width.1 i32 (i32.const 14932))
	(global $runtime.zptr (mut i32) (i32.const 15040))
	(global $__heap_base i32 (i32.const 15136))
	(func $$math.waSqrtF64 (param $x f64) (result f64)
		local.get $x
		f64.sqrt
	)
	(func $runtime.HeapAlloc (export "runtime.HeapAlloc") (param $nbytes i32) (result i32)
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
			i64.store align=8
			local.get $nbytes
			if
				br $zero
			end
		end
		local.get $ptr
	)
	(func $runtime.HeapFree (export "runtime.HeapFree") (param $ptr i32)
		local.get $ptr
		call $runtime.free
	)
	(func $runtime.Block.Init (param $ptr i32) (param $item_count i32) (param $release_func i32) (param $item_size i32) (result i32)
		local.get $ptr
		local.get $ptr
		if
			local.get $ptr
			i32.const 1
			i32.store align=1
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
	(func $runtime.Block.HeapAlloc (export "runtime.Block.HeapAlloc") (param $item_count i32) (param $release_func i32) (param $item_size i32) (result i32)
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
	(func $runtime.DupI32 (param $a i32) (result i32 i32)
		local.get $a
		local.get $a
	)
	(func $runtime.Block.Retain (export "runtime.Block.Retain") (param $ptr i32) (result i32)
		local.get $ptr
		local.get $ptr
		if
			local.get $ptr
			local.get $ptr
			i32.load align=1
			i32.const 1
			i32.add
			i32.store align=1
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
		i32.load align=1
		i32.const 1
		i32.sub
		local.set $ref_count
		local.get $ref_count
		if
			local.get $ptr
			local.get $ref_count
			i32.store align=1
		else
			local.get $ptr
			i32.load offset=8 align=1
			local.set $free_func
			local.get $free_func
			if
				local.get $ptr
				i32.load offset=4 align=1
				local.set $item_count
				local.get $item_count
				if
					local.get $ptr
					i32.load offset=12 align=1
					local.set $item_size
					local.get $ptr
					i32.const 16
					i32.add
					local.set $data_ptr
					loop $free_next
						local.get $data_ptr
						local.get $free_func
						call_indirect 0 (type $$OnFree)
						local.get $item_count
						i32.const 1
						i32.sub
						local.set $item_count
						local.get $item_count
						if
							local.get $data_ptr
							local.get $item_size
							i32.add
							local.set $data_ptr
							br $free_next
						end
					end
				end
			end
			local.get $ptr
			call $runtime.HeapFree
		end
	)
	(func $heap_assert_valid_ptr (param $ptr i32)
		local.get $ptr
		i32.const 0
		i32.gt_s
		if
		else
			unreachable
		end
		local.get $ptr
		i32.const 4
		i32.rem_s
		i32.eqz
		if
		else
			unreachable
		end
	)
	(func $heap_is_fixed_list_enabled (result i32)
		global.get $__heap_lfixed_cap
		i32.eqz
		if(result i32)
			i32.const 0
		else
			i32.const 1
		end
	)
	(func $heap_assert_fixed_list_enabled
		global.get $__heap_lfixed_cap
		i32.eqz
		if
			unreachable
		end
	)
	(func $heap_is_fixed_size (param $size i32) (result i32)
		call $heap_is_fixed_list_enabled
		if
		else
			i32.const 0
			return
		end
		local.get $size
		i32.const 80
		i32.le_s
		if(result i32)
			i32.const 1
		else
			i32.const 0
		end
	)
	(func $heap_alignment8 (param $size i32) (result i32)
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
		if
		else
			unreachable
		end
	)
	(func $heap_block.init (param $ptr i32) (param $size i32) (param $next i32)
		local.get $ptr
		local.get $size
		i32.store
		local.get $ptr
		local.get $next
		i32.store offset=4
	)
	(func $heap_block.set_size (param $ptr i32) (param $size i32)
		local.get $ptr
		local.get $size
		i32.store
	)
	(func $heap_block.set_next (param $ptr i32) (param $next i32)
		local.get $ptr
		local.get $next
		i32.store offset=4
	)
	(func $heap_block.size (param $ptr i32) (result i32)
		local.get $ptr
		i32.load
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
	(func $heap_free_list.ptr_and_fixed_size (param $size i32) (result i32 i32)
		call $heap_is_fixed_list_enabled
		if
		else
			global.get $__heap_base
			i32.const 32
			i32.add
			local.get $size
			call $heap_alignment8
			return
		end
		local.get $size
		i32.const 80
		i32.gt_s
		if
			global.get $__heap_base
			i32.const 32
			i32.add
			block(result i32)
				local.get $size
				i32.const 128
				i32.le_s
				if(result i32)
					i32.const 128
				else
					local.get $size
					call $heap_alignment8
				end
			end
			return
		end
		local.get $size
		i32.const 48
		i32.gt_s
		if
			global.get $__heap_base
			i32.const 24
			i32.add
			i32.const 80
			return
		end
		local.get $size
		i32.const 32
		i32.gt_s
		if
			global.get $__heap_base
			i32.const 16
			i32.add
			i32.const 48
			return
		end
		local.get $size
		i32.const 24
		i32.gt_s
		if
			global.get $__heap_base
			i32.const 8
			i32.add
			i32.const 32
			return
		end
		global.get $__heap_base
		i32.const 24
		return
	)
	(func $wa_malloc_init_once
		global.get $__heap_init_flag
		if
			return
		end
		i32.const 1
		global.set $__heap_init_flag
		global.get $__stack_ptr
		i32.const 0
		i32.gt_s
		if
		else
			unreachable
		end
		global.get $__stack_ptr
		global.get $__heap_base
		i32.lt_s
		if
		else
			unreachable
		end
		global.get $__heap_base
		call $heap_assert_align8
		global.get $__heap_base
		i32.const 48
		i32.add
		global.set $__heap_ptr
		memory.size
		i32.const 65536
		i32.mul
		global.set $__heap_top
		global.get $__heap_top
		global.get $__heap_ptr
		i32.gt_s
		if
		else
			unreachable
		end
		global.get $__heap_base
		i32.const 0
		i32.const 48
		memory.fill
		global.get $__heap_base
		i32.const 32
		i32.add
		global.set $__heap_l128_freep
		global.get $__heap_l128_freep
		i32.const 0
		global.get $__heap_l128_freep
		call $heap_block.init
	)
	(func $runtime.malloc (param $size i32) (result i32)
		(local $free_list i32)
		(local $b i32)
		call $wa_malloc_init_once
		local.get $size
		call $heap_alignment8
		local.set $size
		local.get $size
		call $heap_free_list.ptr_and_fixed_size
		local.set $size
		local.set $free_list
		call $heap_is_fixed_list_enabled
		if
			local.get $size
			call $heap_is_fixed_size
			if
				local.get $free_list
				call $wa_malloc_reuse_fixed
				local.tee $b
				i32.eqz
				if
				else
					local.get $b
					call $heap_block.data
					return
				end
			end
		end
		local.get $size
		call $heap_reuse_varying
		local.tee $b
		i32.eqz
		if
		else
			local.get $b
			call $heap_block.data
			return
		end
		local.get $size
		call $heap_new_allocation
		local.tee $b
		i32.eqz
		if
		else
			local.get $b
			call $heap_block.data
			return
		end
		i32.const 0
		return
	)
	(func $wa_malloc_reuse_fixed (param $free_list i32) (result i32)
		(local $p i32)
		call $heap_assert_fixed_list_enabled
		local.get $free_list
		call $heap_block.size
		i32.eqz
		if
			i32.const 0
			return
		end
		local.get $free_list
		block(result i32)
			local.get $free_list
			call $heap_block.size
			i32.const 1
			i32.sub
		end
		call $heap_block.set_size
		local.get $free_list
		call $heap_block.next
		local.set $p
		local.get $free_list
		block(result i32)
			local.get $p
			call $heap_block.next
		end
		call $heap_block.set_next
		local.get $p
		i32.const 0
		call $heap_block.set_next
		local.get $p
		return
	)
	(func $heap_reuse_varying (param $nbytes i32) (result i32)
		(local $prevp i32)
		(local $remaining i32)
		(local $p i32)
		global.get $__heap_l128_freep
		local.set $prevp
		local.get $prevp
		call $heap_block.next
		local.set $p
		loop $continue
			local.get $p
			call $heap_block.size
			local.get $nbytes
			i32.const 8
			i32.add
			i32.ge_s
			if
				local.get $p
				call $heap_block.data
				local.get $nbytes
				i32.add
				local.set $remaining
				local.get $remaining
				local.get $p
				call $heap_block.next
				call $heap_block.set_next
				local.get $remaining
				local.get $p
				call $heap_block.size
				local.get $nbytes
				i32.sub
				i32.const 8
				i32.sub
				call $heap_block.set_size
				local.get $prevp
				local.get $remaining
				call $heap_block.set_next
				local.get $prevp
				global.set $__heap_l128_freep
				local.get $p
				local.get $nbytes
				i32.const 0
				call $heap_block.init
				local.get $p
				return
			end
			local.get $p
			call $heap_block.size
			local.get $nbytes
			i32.ge_s
			if
				local.get $prevp
				block(result i32)
					local.get $p
					call $heap_block.next
				end
				call $heap_block.set_next
				local.get $prevp
				global.set $__heap_l128_freep
				local.get $p
				i32.const 0
				call $heap_block.set_next
				local.get $p
				return
			end
			local.get $p
			global.get $__heap_l128_freep
			i32.eq
			if
				i32.const 0
				return
			end
			local.get $p
			local.set $prevp
			local.get $p
			call $heap_block.next
			local.set $p
			br $continue
		end
		unreachable
	)
	(func $heap_new_allocation (param $size i32) (result i32)
		(local $ptr i32)
		(local $block_size i32)
		(local $pages i32)
		global.get $__heap_ptr
		local.set $ptr
		i32.const 8
		local.get $size
		i32.add
		local.set $block_size
		global.get $__heap_ptr
		local.get $block_size
		i32.add
		global.get $__heap_top
		i32.ge_s
		if
			local.get $block_size
			i32.const 65535
			i32.add
			i32.const 65536
			i32.div_s
			local.set $pages
			local.get $pages
			memory.grow
			i32.const 0
			i32.lt_s
			if
				i32.const 0
				return
			end
			global.get $__heap_top
			block(result i32)
				local.get $pages
				i32.const 65536
				i32.mul
			end
			i32.add
			global.set $__heap_top
		end
		global.get $__heap_ptr
		local.get $block_size
		i32.add
		global.set $__heap_ptr
		local.get $ptr
		local.get $size
		i32.const 0
		call $heap_block.init
		local.get $ptr
		return
	)
	(func $runtime.free (param $ptr i32)
		(local $size i32)
		(local $block i32)
		(local $freep i32)
		call $wa_malloc_init_once
		local.get $ptr
		call $heap_assert_valid_ptr
		local.get $ptr
		call $heap_assert_align8
		local.get $ptr
		i32.const 8
		i32.sub
		local.set $block
		local.get $block
		call $heap_block.size
		local.set $size
		call $heap_is_fixed_list_enabled
		if
			local.get $size
			call $heap_is_fixed_size
			if
				local.get $size
				call $heap_free_list.ptr_and_fixed_size
				drop
				local.get $block
				call $wa_lfixed_free_block
				return
			end
		end
		local.get $block
		call $wa_l128_free
		return
	)
	(func $wa_lfixed_free_block (param $freep i32) (param $block i32)
		call $heap_assert_fixed_list_enabled
		local.get $freep
		call $heap_block.size
		global.get $__heap_lfixed_cap
		i32.eq
		if
			local.get $freep
			call $wa_lfixed_free_all
		end
		local.get $block
		block(result i32)
			local.get $freep
			call $heap_block.next
		end
		call $heap_block.set_next
		local.get $freep
		local.get $block
		call $heap_block.set_next
		local.get $freep
		block(result i32)
			local.get $freep
			call $heap_block.size
			i32.const 1
			i32.add
		end
		call $heap_block.set_size
	)
	(func $wa_lfixed_free_all (param $freep i32)
		(local $p i32)
		(local $temp i32)
		call $heap_assert_fixed_list_enabled
		local.get $freep
		call $heap_block.next
		local.set $p
		block $break
			loop $continue
				local.get $p
				i32.eqz
				if
					br $break
				end
				local.get $p
				call $heap_block.next
				local.set $temp
				local.get $p
				call $wa_l128_free
				local.get $temp
				local.set $p
				br $continue
			end
		end
		local.get $freep
		i32.const 0
		i32.const 0
		call $heap_block.init
	)
	(func $wa_l128_free (param $bp i32)
		(local $p i32)
		global.get $__heap_l128_freep
		local.set $p
		block $break
			loop $continue
				local.get $bp
				local.get $p
				i32.gt_s
				if
					local.get $bp
					local.get $p
					call $heap_block.next
					i32.lt_s
					if
						br $break
					end
				end
				local.get $p
				local.get $p
				call $heap_block.next
				i32.ge_s
				if
					local.get $bp
					local.get $p
					i32.gt_s
					if
						br $break
					end
					local.get $bp
					local.get $p
					call $heap_block.next
					i32.lt_s
					if
						br $break
					end
				end
				local.get $p
				call $heap_block.next
				local.set $p
				br $continue
			end
		end
		local.get $bp
		local.get $bp
		call $heap_block.size
		i32.add
		i32.const 8
		i32.add
		local.get $p
		call $heap_block.next
		i32.eq
		if
			local.get $bp
			block(result i32)
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
			local.get $bp
			block(result i32)
				local.get $p
				call $heap_block.next
				call $heap_block.next
			end
			call $heap_block.set_next
		else
			local.get $bp
			local.get $p
			call $heap_block.next
			call $heap_block.set_next
		end
		local.get $p
		local.get $p
		call $heap_block.size
		i32.add
		i32.const 8
		i32.add
		local.get $bp
		i32.eq
		if
			local.get $p
			block(result i32)
				local.get $p
				call $heap_block.size
				local.get $bp
				call $heap_block.size
				i32.const 8
				i32.add
				i32.add
			end
			call $heap_block.set_size
			local.get $p
			local.get $bp
			call $heap_block.next
			call $heap_block.set_next
		else
			local.get $p
			local.get $bp
			call $heap_block.set_next
		end
		local.get $p
		global.set $__heap_l128_freep
	)
	(func $$wa.runtime.string_to_iter (param $b i32) (param $d i32) (param $l i32) (result i32 i32 i32)
		local.get $d
		local.get $l
		i32.const 0
	)
	(func $$syscall/wasm4.__linkname__make_slice (param $blk i32) (param $ptr i32) (param $len i32) (param $cap i32) (result i32 i32 i32 i32)
		local.get $blk
		local.get $ptr
		local.get $len
		local.get $cap
		return
	)
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
			i32.load
			call $runtime.Block.Retain
			local.get $p0
			i32.load offset=4
			local.get $p0
			i32.load offset=8
			local.set $v0.2
			local.set $v0.1
			local.get $v0.0
			call $runtime.Block.Release
			local.set $v0.0
		end
		local.get $p1
		if
			local.get $p1
			i32.load
			call $runtime.Block.Retain
			local.get $p1
			i32.load offset=4
			local.get $p1
			i32.load offset=8
			local.set $v1.2
			local.set $v1.1
			local.get $v1.0
			call $runtime.Block.Release
			local.set $v1.0
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
	)
	(func $$u8.$$block.$$OnFree (param $ptr i32)
		local.get $ptr
		i32.load align=1
		call $runtime.Block.Release
		local.get $ptr
		i32.const 0
		i32.store align=1
	)
	(func $$string.underlying.$$OnFree (param $$ptr i32)
		local.get $$ptr
		i32.const 2
		call_indirect 0 (type $$OnFree)
	)
	(func $runtime.get_u8 (param $addr i32) (result i32)
		local.get $addr
		i32.load8_u align=1
	)
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
								br_table  0 1 2 0
							end
							i32.const 0
							local.set $$current_block
							global.get $runtime.init$guard
							local.set $$t0
							local.get $$t0
							if
								br $$Block_1
							else
								br $$Block_0
							end
						end
						i32.const 1
						local.set $$current_block
						i32.const 1
						global.set $runtime.init$guard
						call $syscall$wasm4.init
						br $$Block_1
					end
					i32.const 2
					local.set $$current_block
					br $$BlockFnBody
				end
			end
		end
	)
	(func $$runtime.mapImp.$$block.$$OnFree (param $ptr i32)
		local.get $ptr
		i32.load align=1
		call $runtime.Block.Release
		local.get $ptr
		i32.const 0
		i32.store align=1
	)
	(func $$runtime.mapImp.$ref.underlying.$$OnFree (param $$ptr i32)
		local.get $$ptr
		i32.const 4
		call_indirect 0 (type $$OnFree)
	)
	(func $$runtime.mapIter.$$OnFree (param $$ptr i32)
		local.get $$ptr
		i32.const 5
		call_indirect 0 (type $$OnFree)
	)
	(func $$runtime.mapNode.$$block.$$OnFree (param $ptr i32)
		local.get $ptr
		i32.load align=1
		call $runtime.Block.Release
		local.get $ptr
		i32.const 0
		i32.store align=1
	)
	(func $$runtime.mapNode.$ref.underlying.$$OnFree (param $$ptr i32)
		local.get $$ptr
		i32.const 7
		call_indirect 0 (type $$OnFree)
	)
	(func $$void.$$block.$$OnFree (param $ptr i32)
		local.get $ptr
		i32.load align=1
		call $runtime.Block.Release
		local.get $ptr
		i32.const 0
		i32.store align=1
	)
	(func $$void.$ref.underlying.$$OnFree (param $$ptr i32)
		local.get $$ptr
		i32.const 9
		call_indirect 0 (type $$OnFree)
	)
	(func $$i`0`.underlying.$$OnFree (param $$ptr i32)
		local.get $$ptr
		i32.const 10
		call_indirect 0 (type $$OnFree)
	)
	(func $$runtime.mapNode.$$OnFree (param $$ptr i32)
		local.get $$ptr
		i32.const 8
		i32.add
		i32.const 8
		call_indirect 0 (type $$OnFree)
		local.get $$ptr
		i32.const 16
		i32.add
		i32.const 8
		call_indirect 0 (type $$OnFree)
		local.get $$ptr
		i32.const 28
		i32.add
		i32.const 11
		call_indirect 0 (type $$OnFree)
		local.get $$ptr
		i32.const 44
		i32.add
		i32.const 11
		call_indirect 0 (type $$OnFree)
	)
	(func $$runtime.mapNode.$ref.$$block.$$OnFree (param $ptr i32)
		local.get $ptr
		i32.load align=1
		call $runtime.Block.Release
		local.get $ptr
		i32.const 0
		i32.store align=1
	)
	(func $$runtime.mapNode.$ref.$slice.underlying.$$OnFree (param $$ptr i32)
		local.get $$ptr
		i32.const 13
		call_indirect 0 (type $$OnFree)
	)
	(func $$runtime.mapImp.$$OnFree (param $$ptr i32)
		local.get $$ptr
		i32.const 8
		call_indirect 0 (type $$OnFree)
		local.get $$ptr
		i32.const 8
		i32.add
		i32.const 8
		call_indirect 0 (type $$OnFree)
		local.get $$ptr
		i32.const 16
		i32.add
		i32.const 14
		call_indirect 0 (type $$OnFree)
	)
	(func $$runtime.mapNode.$ref.$array1.underlying.$$OnFree (param $$ptr i32)
		local.get $$ptr
		i32.const 8
		call_indirect 0 (type $$OnFree)
	)
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
																br_table  0 1 2 3 4 5 6 7 8 9 10 0
															end
															i32.const 0
															local.set $$current_block
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
															local.get $$t0.1
															local.get $iter.0
															i32.store
															local.get $$t0.1
															local.get $iter.1
															i32.store offset=4
															local.get $$t0.1
															local.get $iter.2
															i32.store offset=8
															local.get $$t0.0
															call $runtime.Block.Retain
															local.get $$t0.1
															i32.const 8
															i32.add
															local.set $$t1.1
															local.get $$t1.0
															call $runtime.Block.Release
															local.set $$t1.0
															local.get $$t1.1
															i32.load
															local.set $$t2
															local.get $$t0.0
															call $runtime.Block.Retain
															local.get $$t0.1
															i32.const 4
															i32.add
															local.set $$t3.1
															local.get $$t3.0
															call $runtime.Block.Release
															local.set $$t3.0
															local.get $$t3.1
															i32.load
															local.set $$t4
															local.get $$t2
															local.get $$t4
															i32.ge_s
															local.set $$t5
															local.get $$t5
															if
																br $$Block_0
															else
																br $$Block_1
															end
														end
														i32.const 1
														local.set $$current_block
														local.get $$t0.0
														call $runtime.Block.Retain
														local.get $$t0.1
														i32.const 8
														i32.add
														local.set $$t6.1
														local.get $$t6.0
														call $runtime.Block.Release
														local.set $$t6.0
														local.get $$t6.1
														i32.load
														local.set $$t7
														local.get $$t0.0
														call $runtime.Block.Retain
														local.get $$t0.1
														i32.const 8
														i32.add
														local.set $$t8.1
														local.get $$t8.0
														call $runtime.Block.Release
														local.set $$t8.0
														local.get $$t8.1
														i32.load
														local.set $$t9
														i32.const 0
														local.set $$ret_0
														local.get $$t7
														local.set $$ret_1
														i32.const 0
														local.set $$ret_2
														local.get $$t9
														local.set $$ret_3
														br $$BlockFnBody
													end
													i32.const 2
													local.set $$current_block
													local.get $$t0.0
													call $runtime.Block.Retain
													local.get $$t0.1
													i32.const 0
													i32.add
													local.set $$t10.1
													local.get $$t10.0
													call $runtime.Block.Release
													local.set $$t10.0
													local.get $$t10.1
													i32.load
													local.set $$t11
													local.get $$t11
													local.set $$t12
													local.get $$t0.0
													call $runtime.Block.Retain
													local.get $$t0.1
													i32.const 8
													i32.add
													local.set $$t13.1
													local.get $$t13.0
													call $runtime.Block.Release
													local.set $$t13.0
													local.get $$t13.1
													i32.load
													local.set $$t14
													local.get $$t14
													local.set $$t15
													local.get $$t12
													local.get $$t15
													i32.add
													local.set $$t16
													local.get $$t16
													call $runtime.get_u8
													local.set $$t17
													local.get $$t17
													local.set $$t18
													local.get $$t18
													i32.const 128
													i32.and
													local.set $$t19
													local.get $$t19
													i32.const 0
													i32.eq
													local.set $$t20
													local.get $$t20
													if
														br $$Block_2
													else
														br $$Block_3
													end
												end
												i32.const 3
												local.set $$current_block
												local.get $$t0.0
												call $runtime.Block.Retain
												local.get $$t0.1
												i32.const 8
												i32.add
												local.set $$t21.1
												local.get $$t21.0
												call $runtime.Block.Release
												local.set $$t21.0
												local.get $$t21.1
												i32.load
												local.set $$t22
												local.get $$t0.0
												call $runtime.Block.Retain
												local.get $$t0.1
												i32.const 8
												i32.add
												local.set $$t23.1
												local.get $$t23.0
												call $runtime.Block.Release
												local.set $$t23.0
												local.get $$t23.1
												i32.load
												local.set $$t24
												local.get $$t24
												i32.const 1
												i32.add
												local.set $$t25
												i32.const 1
												local.set $$ret_0
												local.get $$t22
												local.set $$ret_1
												local.get $$t18
												local.set $$ret_2
												local.get $$t25
												local.set $$ret_3
												br $$BlockFnBody
											end
											i32.const 4
											local.set $$current_block
											local.get $$t18
											i32.const 224
											i32.and
											local.set $$t26
											local.get $$t26
											i32.const 192
											i32.eq
											local.set $$t27
											local.get $$t27
											if
												br $$Block_4
											else
												br $$Block_5
											end
										end
										i32.const 5
										local.set $$current_block
										local.get $$t18
										i32.const 31
										i32.and
										local.set $$t28
										local.get $$t28
										i64.const 6
										i32.wrap_i64
										i32.shl
										local.set $$t29
										local.get $$t0.0
										call $runtime.Block.Retain
										local.get $$t0.1
										i32.const 0
										i32.add
										local.set $$t30.1
										local.get $$t30.0
										call $runtime.Block.Release
										local.set $$t30.0
										local.get $$t30.1
										i32.load
										local.set $$t31
										local.get $$t31
										local.set $$t32
										local.get $$t0.0
										call $runtime.Block.Retain
										local.get $$t0.1
										i32.const 8
										i32.add
										local.set $$t33.1
										local.get $$t33.0
										call $runtime.Block.Release
										local.set $$t33.0
										local.get $$t33.1
										i32.load
										local.set $$t34
										local.get $$t34
										local.set $$t35
										local.get $$t32
										local.get $$t35
										i32.add
										local.set $$t36
										local.get $$t36
										i32.const 1
										i32.add
										local.set $$t37
										local.get $$t37
										call $runtime.get_u8
										local.set $$t38
										local.get $$t38
										local.set $$t39
										local.get $$t39
										i32.const 63
										i32.and
										local.set $$t40
										local.get $$t0.0
										call $runtime.Block.Retain
										local.get $$t0.1
										i32.const 8
										i32.add
										local.set $$t41.1
										local.get $$t41.0
										call $runtime.Block.Release
										local.set $$t41.0
										local.get $$t41.1
										i32.load
										local.set $$t42
										local.get $$t29
										local.get $$t40
										i32.or
										local.set $$t43
										local.get $$t0.0
										call $runtime.Block.Retain
										local.get $$t0.1
										i32.const 8
										i32.add
										local.set $$t44.1
										local.get $$t44.0
										call $runtime.Block.Release
										local.set $$t44.0
										local.get $$t44.1
										i32.load
										local.set $$t45
										local.get $$t45
										i32.const 2
										i32.add
										local.set $$t46
										i32.const 1
										local.set $$ret_0
										local.get $$t42
										local.set $$ret_1
										local.get $$t43
										local.set $$ret_2
										local.get $$t46
										local.set $$ret_3
										br $$BlockFnBody
									end
									i32.const 6
									local.set $$current_block
									local.get $$t18
									i32.const 240
									i32.and
									local.set $$t47
									local.get $$t47
									i32.const 224
									i32.eq
									local.set $$t48
									local.get $$t48
									if
										br $$Block_6
									else
										br $$Block_7
									end
								end
								i32.const 7
								local.set $$current_block
								local.get $$t18
								i32.const 15
								i32.and
								local.set $$t49
								local.get $$t49
								i64.const 12
								i32.wrap_i64
								i32.shl
								local.set $$t50
								local.get $$t0.0
								call $runtime.Block.Retain
								local.get $$t0.1
								i32.const 0
								i32.add
								local.set $$t51.1
								local.get $$t51.0
								call $runtime.Block.Release
								local.set $$t51.0
								local.get $$t51.1
								i32.load
								local.set $$t52
								local.get $$t52
								local.set $$t53
								local.get $$t0.0
								call $runtime.Block.Retain
								local.get $$t0.1
								i32.const 8
								i32.add
								local.set $$t54.1
								local.get $$t54.0
								call $runtime.Block.Release
								local.set $$t54.0
								local.get $$t54.1
								i32.load
								local.set $$t55
								local.get $$t55
								local.set $$t56
								local.get $$t53
								local.get $$t56
								i32.add
								local.set $$t57
								local.get $$t57
								i32.const 1
								i32.add
								local.set $$t58
								local.get $$t58
								call $runtime.get_u8
								local.set $$t59
								local.get $$t59
								local.set $$t60
								local.get $$t60
								i32.const 63
								i32.and
								local.set $$t61
								local.get $$t61
								i64.const 6
								i32.wrap_i64
								i32.shl
								local.set $$t62
								local.get $$t0.0
								call $runtime.Block.Retain
								local.get $$t0.1
								i32.const 0
								i32.add
								local.set $$t63.1
								local.get $$t63.0
								call $runtime.Block.Release
								local.set $$t63.0
								local.get $$t63.1
								i32.load
								local.set $$t64
								local.get $$t64
								local.set $$t65
								local.get $$t0.0
								call $runtime.Block.Retain
								local.get $$t0.1
								i32.const 8
								i32.add
								local.set $$t66.1
								local.get $$t66.0
								call $runtime.Block.Release
								local.set $$t66.0
								local.get $$t66.1
								i32.load
								local.set $$t67
								local.get $$t67
								local.set $$t68
								local.get $$t65
								local.get $$t68
								i32.add
								local.set $$t69
								local.get $$t69
								i32.const 2
								i32.add
								local.set $$t70
								local.get $$t70
								call $runtime.get_u8
								local.set $$t71
								local.get $$t71
								local.set $$t72
								local.get $$t72
								i32.const 63
								i32.and
								local.set $$t73
								local.get $$t0.0
								call $runtime.Block.Retain
								local.get $$t0.1
								i32.const 8
								i32.add
								local.set $$t74.1
								local.get $$t74.0
								call $runtime.Block.Release
								local.set $$t74.0
								local.get $$t74.1
								i32.load
								local.set $$t75
								local.get $$t50
								local.get $$t62
								i32.or
								local.set $$t76
								local.get $$t76
								local.get $$t73
								i32.or
								local.set $$t77
								local.get $$t0.0
								call $runtime.Block.Retain
								local.get $$t0.1
								i32.const 8
								i32.add
								local.set $$t78.1
								local.get $$t78.0
								call $runtime.Block.Release
								local.set $$t78.0
								local.get $$t78.1
								i32.load
								local.set $$t79
								local.get $$t79
								i32.const 3
								i32.add
								local.set $$t80
								i32.const 1
								local.set $$ret_0
								local.get $$t75
								local.set $$ret_1
								local.get $$t77
								local.set $$ret_2
								local.get $$t80
								local.set $$ret_3
								br $$BlockFnBody
							end
							i32.const 8
							local.set $$current_block
							local.get $$t18
							i32.const 248
							i32.and
							local.set $$t81
							local.get $$t81
							i32.const 240
							i32.eq
							local.set $$t82
							local.get $$t82
							if
								br $$Block_8
							else
								br $$Block_9
							end
						end
						i32.const 9
						local.set $$current_block
						local.get $$t18
						i32.const 7
						i32.and
						local.set $$t83
						local.get $$t83
						i64.const 18
						i32.wrap_i64
						i32.shl
						local.set $$t84
						local.get $$t0.0
						call $runtime.Block.Retain
						local.get $$t0.1
						i32.const 0
						i32.add
						local.set $$t85.1
						local.get $$t85.0
						call $runtime.Block.Release
						local.set $$t85.0
						local.get $$t85.1
						i32.load
						local.set $$t86
						local.get $$t86
						local.set $$t87
						local.get $$t0.0
						call $runtime.Block.Retain
						local.get $$t0.1
						i32.const 8
						i32.add
						local.set $$t88.1
						local.get $$t88.0
						call $runtime.Block.Release
						local.set $$t88.0
						local.get $$t88.1
						i32.load
						local.set $$t89
						local.get $$t89
						local.set $$t90
						local.get $$t87
						local.get $$t90
						i32.add
						local.set $$t91
						local.get $$t91
						i32.const 1
						i32.add
						local.set $$t92
						local.get $$t92
						call $runtime.get_u8
						local.set $$t93
						local.get $$t93
						local.set $$t94
						local.get $$t94
						i32.const 63
						i32.and
						local.set $$t95
						local.get $$t95
						i64.const 12
						i32.wrap_i64
						i32.shl
						local.set $$t96
						local.get $$t0.0
						call $runtime.Block.Retain
						local.get $$t0.1
						i32.const 0
						i32.add
						local.set $$t97.1
						local.get $$t97.0
						call $runtime.Block.Release
						local.set $$t97.0
						local.get $$t97.1
						i32.load
						local.set $$t98
						local.get $$t98
						local.set $$t99
						local.get $$t0.0
						call $runtime.Block.Retain
						local.get $$t0.1
						i32.const 8
						i32.add
						local.set $$t100.1
						local.get $$t100.0
						call $runtime.Block.Release
						local.set $$t100.0
						local.get $$t100.1
						i32.load
						local.set $$t101
						local.get $$t101
						local.set $$t102
						local.get $$t99
						local.get $$t102
						i32.add
						local.set $$t103
						local.get $$t103
						i32.const 2
						i32.add
						local.set $$t104
						local.get $$t104
						call $runtime.get_u8
						local.set $$t105
						local.get $$t105
						local.set $$t106
						local.get $$t106
						i32.const 63
						i32.and
						local.set $$t107
						local.get $$t107
						i64.const 6
						i32.wrap_i64
						i32.shl
						local.set $$t108
						local.get $$t0.0
						call $runtime.Block.Retain
						local.get $$t0.1
						i32.const 0
						i32.add
						local.set $$t109.1
						local.get $$t109.0
						call $runtime.Block.Release
						local.set $$t109.0
						local.get $$t109.1
						i32.load
						local.set $$t110
						local.get $$t110
						local.set $$t111
						local.get $$t0.0
						call $runtime.Block.Retain
						local.get $$t0.1
						i32.const 8
						i32.add
						local.set $$t112.1
						local.get $$t112.0
						call $runtime.Block.Release
						local.set $$t112.0
						local.get $$t112.1
						i32.load
						local.set $$t113
						local.get $$t113
						local.set $$t114
						local.get $$t111
						local.get $$t114
						i32.add
						local.set $$t115
						local.get $$t115
						i32.const 3
						i32.add
						local.set $$t116
						local.get $$t116
						call $runtime.get_u8
						local.set $$t117
						local.get $$t117
						local.set $$t118
						local.get $$t118
						i32.const 63
						i32.and
						local.set $$t119
						local.get $$t0.0
						call $runtime.Block.Retain
						local.get $$t0.1
						i32.const 8
						i32.add
						local.set $$t120.1
						local.get $$t120.0
						call $runtime.Block.Release
						local.set $$t120.0
						local.get $$t120.1
						i32.load
						local.set $$t121
						local.get $$t84
						local.get $$t96
						i32.or
						local.set $$t122
						local.get $$t122
						local.get $$t108
						i32.or
						local.set $$t123
						local.get $$t123
						local.get $$t119
						i32.or
						local.set $$t124
						local.get $$t0.0
						call $runtime.Block.Retain
						local.get $$t0.1
						i32.const 8
						i32.add
						local.set $$t125.1
						local.get $$t125.0
						call $runtime.Block.Release
						local.set $$t125.0
						local.get $$t125.1
						i32.load
						local.set $$t126
						local.get $$t126
						i32.const 4
						i32.add
						local.set $$t127
						i32.const 1
						local.set $$ret_0
						local.get $$t121
						local.set $$ret_1
						local.get $$t124
						local.set $$ret_2
						local.get $$t127
						local.set $$ret_3
						br $$BlockFnBody
					end
					i32.const 10
					local.set $$current_block
					i32.const 0
					local.set $$ret_0
					i32.const 0
					local.set $$ret_1
					i32.const 0
					local.set $$ret_2
					i32.const 0
					local.set $$ret_3
					br $$BlockFnBody
				end
			end
		end
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
	)
	(func $$$$$$.underlying.$$OnFree (param $$ptr i32)
		local.get $$ptr
		i32.const 4
		i32.add
		i32.const 10
		call_indirect 0 (type $$OnFree)
	)
	(func $$$$$$.$array1.underlying.$$OnFree (param $$ptr i32)
		local.get $$ptr
		i32.const 17
		call_indirect 0 (type $$OnFree)
	)
	(func $$$$$$.$$block.$$OnFree (param $ptr i32)
		local.get $ptr
		i32.load align=1
		call $runtime.Block.Release
		local.get $ptr
		i32.const 0
		i32.store align=1
	)
	(func $$$$$$.$slice.underlying.$$OnFree (param $$ptr i32)
		local.get $$ptr
		i32.const 19
		call_indirect 0 (type $$OnFree)
	)
	(func $$runtime.defers.$$OnFree (param $$ptr i32)
		local.get $$ptr
		i32.const 20
		call_indirect 0 (type $$OnFree)
	)
	(func $$runtime.defers.$array1.underlying.$$OnFree (param $$ptr i32)
		local.get $$ptr
		i32.const 20
		call_indirect 0 (type $$OnFree)
	)
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
																		br_table  0 1 2 3 4 5 6 7 8 9 10 11 12 0
																	end
																	i32.const 0
																	local.set $$current_block
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
																	local.get $x.0
																	local.get $x.1
																	local.get $x.2
																	call $$wa.runtime.string_to_iter
																	local.set $$t1.2
																	local.set $$t1.1
																	local.set $$t1.0
																	local.get $$t0.1
																	local.get $$t1.0
																	i32.store
																	local.get $$t0.1
																	local.get $$t1.1
																	i32.store offset=4
																	local.get $$t0.1
																	local.get $$t1.2
																	i32.store offset=8
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
																	local.get $y.0
																	local.get $y.1
																	local.get $y.2
																	call $$wa.runtime.string_to_iter
																	local.set $$t3.2
																	local.set $$t3.1
																	local.set $$t3.0
																	local.get $$t2.1
																	local.get $$t3.0
																	i32.store
																	local.get $$t2.1
																	local.get $$t3.1
																	i32.store offset=4
																	local.get $$t2.1
																	local.get $$t3.2
																	i32.store offset=8
																	br $$Block_0
																end
																i32.const 1
																local.set $$current_block
																local.get $$t0.1
																i32.load
																local.get $$t0.1
																i32.load offset=4
																local.get $$t0.1
																i32.load offset=8
																local.set $$t4.2
																local.set $$t4.1
																local.set $$t4.0
																local.get $$t4.0
																local.get $$t4.1
																local.get $$t4.2
																call $runtime.next_rune
																local.set $$t5.3
																local.set $$t5.2
																local.set $$t5.1
																local.set $$t5.0
																local.get $$t5.0
																local.set $$t6
																local.get $$t5.1
																local.set $$t7
																local.get $$t5.2
																local.set $$t8
																local.get $$t5.3
																local.set $$t9
																local.get $$t0.0
																call $runtime.Block.Retain
																local.get $$t0.1
																i32.const 8
																i32.add
																local.set $$t10.1
																local.get $$t10.0
																call $runtime.Block.Release
																local.set $$t10.0
																local.get $$t10.1
																local.get $$t9
																i32.store
																local.get $$t2.1
																i32.load
																local.get $$t2.1
																i32.load offset=4
																local.get $$t2.1
																i32.load offset=8
																local.set $$t11.2
																local.set $$t11.1
																local.set $$t11.0
																local.get $$t11.0
																local.get $$t11.1
																local.get $$t11.2
																call $runtime.next_rune
																local.set $$t12.3
																local.set $$t12.2
																local.set $$t12.1
																local.set $$t12.0
																local.get $$t12.0
																local.set $$t13
																local.get $$t12.1
																local.set $$t14
																local.get $$t12.2
																local.set $$t15
																local.get $$t12.3
																local.set $$t16
																local.get $$t2.0
																call $runtime.Block.Retain
																local.get $$t2.1
																i32.const 8
																i32.add
																local.set $$t17.1
																local.get $$t17.0
																call $runtime.Block.Release
																local.set $$t17.0
																local.get $$t17.1
																local.get $$t16
																i32.store
																local.get $$t6
																if
																	br $$Block_3
																else
																	br $$Block_4
																end
															end
															i32.const 2
															local.set $$current_block
															local.get $x.2
															local.set $$t18
															local.get $y.2
															local.set $$t19
															local.get $$t18
															local.get $$t19
															i32.lt_s
															local.set $$t20
															local.get $$t20
															if
																br $$Block_8
															else
																br $$Block_9
															end
														end
														i32.const 3
														local.set $$current_block
														local.get $$t8
														local.get $$t15
														i32.lt_s
														local.set $$t21
														local.get $$t21
														if
															br $$Block_5
														else
															br $$Block_6
														end
													end
													i32.const 4
													local.set $$current_block
													br $$Block_4
												end
												local.get $$current_block
												i32.const 1
												i32.eq
												if(result i32)
													i32.const 0
												else
													local.get $$t13
												end
												local.set $$t22
												i32.const 5
												local.set $$current_block
												local.get $$t22
												i32.const 1
												i32.eq
												i32.eqz
												local.set $$t23
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
											end
											i32.const 6
											local.set $$current_block
											i32.const -1
											local.set $$ret_0
											br $$BlockFnBody
										end
										i32.const 7
										local.set $$current_block
										local.get $$t8
										local.get $$t15
										i32.gt_s
										local.set $$t24
										local.get $$t24
										if
											br $$Block_7
										else
											i32.const 1
											local.set $$block_selector
											br $$BlockDisp
										end
									end
									i32.const 8
									local.set $$current_block
									i32.const 1
									local.set $$ret_0
									br $$BlockFnBody
								end
								i32.const 9
								local.set $$current_block
								i32.const -1
								local.set $$ret_0
								br $$BlockFnBody
							end
							i32.const 10
							local.set $$current_block
							local.get $$t18
							local.get $$t19
							i32.gt_s
							local.set $$t25
							local.get $$t25
							if
								br $$Block_10
							else
								br $$Block_11
							end
						end
						i32.const 11
						local.set $$current_block
						i32.const 1
						local.set $$ret_0
						br $$BlockFnBody
					end
					i32.const 12
					local.set $$current_block
					i32.const 0
					local.set $$ret_0
					br $$BlockFnBody
				end
			end
		end
		local.get $$ret_0
		local.get $$t0.0
		call $runtime.Block.Release
		local.get $$t2.0
		call $runtime.Block.Release
		local.get $$t10.0
		call $runtime.Block.Release
		local.get $$t17.0
		call $runtime.Block.Release
	)
	(func $math.Sqrt (param $x f64) (result f64)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$ret_0 f64)
		(local $$t0 f64)
		block $$BlockFnBody
			loop $$BlockDisp
				block $$Block_0
					block $$BlockSel
						local.get $$block_selector
						br_table  0 0
					end
					i32.const 0
					local.set $$current_block
					local.get $x
					call $$math.waSqrtF64
					local.set $$t0
					local.get $$t0
					local.set $$ret_0
					br $$BlockFnBody
				end
			end
		end
		local.get $$ret_0
	)
	(func $math.init
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
								br_table  0 1 2 0
							end
							i32.const 0
							local.set $$current_block
							global.get $math.init$guard
							local.set $$t0
							local.get $$t0
							if
								br $$Block_1
							else
								br $$Block_0
							end
						end
						i32.const 1
						local.set $$current_block
						i32.const 1
						global.set $math.init$guard
						br $$Block_1
					end
					i32.const 2
					local.set $$current_block
					br $$BlockFnBody
				end
			end
		end
	)
	(func $syscall$wasm4.GetFramebuffer (result i32 i32 i32 i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$ret_0.0 i32)
		(local $$ret_0.1 i32)
		(local $$ret_0.2 i32)
		(local $$ret_0.3 i32)
		(local $$t0.0 i32)
		(local $$t0.1 i32)
		(local $$t0.2 i32)
		(local $$t0.3 i32)
		block $$BlockFnBody
			loop $$BlockDisp
				block $$Block_0
					block $$BlockSel
						local.get $$block_selector
						br_table  0 0
					end
					i32.const 0
					local.set $$current_block
					i32.const 0
					i32.const 160
					i32.const 6400
					i32.const 6400
					call $$syscall/wasm4.__linkname__make_slice
					local.set $$t0.3
					local.set $$t0.2
					local.set $$t0.1
					local.get $$t0.0
					call $runtime.Block.Release
					local.set $$t0.0
					local.get $$t0.0
					call $runtime.Block.Retain
					local.get $$t0.1
					local.get $$t0.2
					local.get $$t0.3
					local.set $$ret_0.3
					local.set $$ret_0.2
					local.set $$ret_0.1
					local.get $$ret_0.0
					call $runtime.Block.Release
					local.set $$ret_0.0
					br $$BlockFnBody
				end
			end
		end
		local.get $$ret_0.0
		call $runtime.Block.Retain
		local.get $$ret_0.1
		local.get $$ret_0.2
		local.get $$ret_0.3
		local.get $$ret_0.0
		call $runtime.Block.Release
		local.get $$t0.0
		call $runtime.Block.Release
	)
	(func $syscall$wasm4.GetMouseX (result i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$ret_0 i32)
		(local $$t0.0 i32)
		(local $$t0.1 i32)
		(local $$t0.2 i32)
		(local $$t0.3 i32)
		(local $$t1.0 i32)
		(local $$t1.1 i32)
		(local $$t2 i32)
		(local $$t3 i32)
		block $$BlockFnBody
			loop $$BlockDisp
				block $$Block_0
					block $$BlockSel
						local.get $$block_selector
						br_table  0 0
					end
					i32.const 0
					local.set $$current_block
					i32.const 0
					i32.const 26
					i32.const 1
					i32.const 1
					call $$syscall/wasm4.__linkname__make_slice
					local.set $$t0.3
					local.set $$t0.2
					local.set $$t0.1
					local.get $$t0.0
					call $runtime.Block.Release
					local.set $$t0.0
					local.get $$t0.0
					call $runtime.Block.Retain
					local.get $$t0.1
					i32.const 2
					i32.const 0
					i32.mul
					i32.add
					local.set $$t1.1
					local.get $$t1.0
					call $runtime.Block.Release
					local.set $$t1.0
					local.get $$t1.1
					i32.load16_u
					local.set $$t2
					local.get $$t2
					local.set $$t3
					local.get $$t3
					local.set $$ret_0
					br $$BlockFnBody
				end
			end
		end
		local.get $$ret_0
		local.get $$t0.0
		call $runtime.Block.Release
		local.get $$t1.0
		call $runtime.Block.Release
	)
	(func $syscall$wasm4.GetMouseY (result i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$ret_0 i32)
		(local $$t0.0 i32)
		(local $$t0.1 i32)
		(local $$t0.2 i32)
		(local $$t0.3 i32)
		(local $$t1.0 i32)
		(local $$t1.1 i32)
		(local $$t2 i32)
		(local $$t3 i32)
		block $$BlockFnBody
			loop $$BlockDisp
				block $$Block_0
					block $$BlockSel
						local.get $$block_selector
						br_table  0 0
					end
					i32.const 0
					local.set $$current_block
					i32.const 0
					i32.const 28
					i32.const 1
					i32.const 1
					call $$syscall/wasm4.__linkname__make_slice
					local.set $$t0.3
					local.set $$t0.2
					local.set $$t0.1
					local.get $$t0.0
					call $runtime.Block.Release
					local.set $$t0.0
					local.get $$t0.0
					call $runtime.Block.Retain
					local.get $$t0.1
					i32.const 2
					i32.const 0
					i32.mul
					i32.add
					local.set $$t1.1
					local.get $$t1.0
					call $runtime.Block.Release
					local.set $$t1.0
					local.get $$t1.1
					i32.load16_u
					local.set $$t2
					local.get $$t2
					local.set $$t3
					local.get $$t3
					local.set $$ret_0
					br $$BlockFnBody
				end
			end
		end
		local.get $$ret_0
		local.get $$t0.0
		call $runtime.Block.Release
		local.get $$t1.0
		call $runtime.Block.Release
	)
	(func $syscall$wasm4.init
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
								br_table  0 1 2 0
							end
							i32.const 0
							local.set $$current_block
							global.get $syscall$wasm4.init$guard
							local.set $$t0
							local.get $$t0
							if
								br $$Block_1
							else
								br $$Block_0
							end
						end
						i32.const 1
						local.set $$current_block
						i32.const 1
						global.set $syscall$wasm4.init$guard
						br $$Block_1
					end
					i32.const 2
					local.set $$current_block
					br $$BlockFnBody
				end
			end
		end
	)
	(func $$u8.$slice.underlying.$$OnFree (param $$ptr i32)
		local.get $$ptr
		i32.const 2
		call_indirect 0 (type $$OnFree)
	)
	(func $$w4app.Framebuffer.$$OnFree (param $$ptr i32)
		local.get $$ptr
		i32.const 23
		call_indirect 0 (type $$OnFree)
	)
	(func $w4app.FramebufferInstance (export "w4app.FramebufferInstance") (result i32 i32)
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
		(local $$t2.2 i32)
		(local $$t2.3 i32)
		block $$BlockFnBody
			loop $$BlockDisp
				block $$Block_0
					block $$BlockSel
						local.get $$block_selector
						br_table  0 0
					end
					i32.const 0
					local.set $$current_block
					i32.const 32
					call $runtime.HeapAlloc
					i32.const 1
					i32.const 24
					i32.const 16
					call $runtime.Block.Init
					call $runtime.DupI32
					i32.const 16
					i32.add
					local.set $$t0.1
					local.get $$t0.0
					call $runtime.Block.Release
					local.set $$t0.0
					local.get $$t0.0
					call $runtime.Block.Retain
					local.get $$t0.1
					i32.const 0
					i32.add
					local.set $$t1.1
					local.get $$t1.0
					call $runtime.Block.Release
					local.set $$t1.0
					call $syscall$wasm4.GetFramebuffer
					local.set $$t2.3
					local.set $$t2.2
					local.set $$t2.1
					local.get $$t2.0
					call $runtime.Block.Release
					local.set $$t2.0
					local.get $$t1.1
					local.get $$t2.0
					call $runtime.Block.Retain
					local.get $$t1.1
					i32.load align=1
					call $runtime.Block.Release
					i32.store align=1
					local.get $$t1.1
					local.get $$t2.1
					i32.store offset=4
					local.get $$t1.1
					local.get $$t2.2
					i32.store offset=8
					local.get $$t1.1
					local.get $$t2.3
					i32.store offset=12
					local.get $$t0.0
					call $runtime.Block.Retain
					local.get $$t0.1
					local.set $$ret_0.1
					local.get $$ret_0.0
					call $runtime.Block.Release
					local.set $$ret_0.0
					br $$BlockFnBody
				end
			end
		end
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
	)
	(func $$w4app.BitImage.$$OnFree (param $$ptr i32)
		local.get $$ptr
		i32.const 8
		i32.add
		i32.const 23
		call_indirect 0 (type $$OnFree)
	)
	(func $w4app.NewBitImage (export "w4app.NewBitImage") (param $w i32) (param $h i32) (result i32 i32)
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
		(local $$t4 i32)
		(local $$t5 i32)
		(local $$t6.0 i32)
		(local $$t6.1 i32)
		(local $$t6.2 i32)
		(local $$t6.3 i32)
		block $$BlockFnBody
			loop $$BlockDisp
				block $$Block_0
					block $$BlockSel
						local.get $$block_selector
						br_table  0 0
					end
					i32.const 0
					local.set $$current_block
					i32.const 40
					call $runtime.HeapAlloc
					i32.const 1
					i32.const 25
					i32.const 24
					call $runtime.Block.Init
					call $runtime.DupI32
					i32.const 16
					i32.add
					local.set $$t0.1
					local.get $$t0.0
					call $runtime.Block.Release
					local.set $$t0.0
					local.get $$t0.0
					call $runtime.Block.Retain
					local.get $$t0.1
					i32.const 0
					i32.add
					local.set $$t1.1
					local.get $$t1.0
					call $runtime.Block.Release
					local.set $$t1.0
					local.get $$t0.0
					call $runtime.Block.Retain
					local.get $$t0.1
					i32.const 4
					i32.add
					local.set $$t2.1
					local.get $$t2.0
					call $runtime.Block.Release
					local.set $$t2.0
					local.get $$t0.0
					call $runtime.Block.Retain
					local.get $$t0.1
					i32.const 8
					i32.add
					local.set $$t3.1
					local.get $$t3.0
					call $runtime.Block.Release
					local.set $$t3.0
					local.get $w
					local.get $h
					i32.mul
					local.set $$t4
					local.get $$t4
					i32.const 8
					i32.div_s
					local.set $$t5
					local.get $$t5
					i32.const 1
					i32.mul
					i32.const 16
					i32.add
					call $runtime.HeapAlloc
					local.get $$t5
					i32.const 0
					i32.const 1
					call $runtime.Block.Init
					call $runtime.DupI32
					i32.const 16
					i32.add
					local.get $$t5
					local.get $$t5
					local.set $$t6.3
					local.set $$t6.2
					local.set $$t6.1
					local.get $$t6.0
					call $runtime.Block.Release
					local.set $$t6.0
					local.get $$t1.1
					local.get $w
					i32.store
					local.get $$t2.1
					local.get $h
					i32.store
					local.get $$t3.1
					local.get $$t6.0
					call $runtime.Block.Retain
					local.get $$t3.1
					i32.load align=1
					call $runtime.Block.Release
					i32.store align=1
					local.get $$t3.1
					local.get $$t6.1
					i32.store offset=4
					local.get $$t3.1
					local.get $$t6.2
					i32.store offset=8
					local.get $$t3.1
					local.get $$t6.3
					i32.store offset=12
					local.get $$t0.0
					call $runtime.Block.Retain
					local.get $$t0.1
					local.set $$ret_0.1
					local.get $$ret_0.0
					call $runtime.Block.Release
					local.set $$ret_0.0
					br $$BlockFnBody
				end
			end
		end
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
		local.get $$t6.0
		call $runtime.Block.Release
	)
	(func $w4app.StepGame (export "w4app.StepGame")
		(local $$block_selector i32)
		(local $$current_block i32)
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
		(local $$t10.0 i32)
		(local $$t10.1 i32)
		(local $$t11 i32)
		(local $$t12 i32)
		(local $$t13 i32)
		(local $$t14 i32)
		(local $$t15 i32)
		(local $$t16 i32)
		(local $$t17 i32)
		(local $$t18 i32)
		(local $$t19 i32)
		(local $$t20 i32)
		(local $$t21 i32)
		(local $$t22 i32)
		(local $$t23.0 i32)
		(local $$t23.1 i32)
		(local $$t24.0 i32)
		(local $$t24.1 i32)
		block $$BlockFnBody
			loop $$BlockDisp
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
																					br_table  0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 0
																				end
																				i32.const 0
																				local.set $$current_block
																				call $syscall$wasm4.GetMouseX
																				local.set $$t0
																				local.get $$t0
																				i32.const 0
																				i32.gt_s
																				local.set $$t1
																				local.get $$t1
																				if
																					br $$Block_2
																				else
																					br $$Block_1
																				end
																			end
																			i32.const 1
																			local.set $$current_block
																			i32.const 14920
																			local.get $$t0
																			i32.store
																			br $$Block_1
																		end
																		i32.const 2
																		local.set $$current_block
																		call $syscall$wasm4.GetMouseY
																		local.set $$t2
																		local.get $$t2
																		i32.const 0
																		i32.gt_s
																		local.set $$t3
																		local.get $$t3
																		if
																			br $$Block_5
																		else
																			br $$Block_4
																		end
																	end
																	i32.const 3
																	local.set $$current_block
																	i32.const 14932
																	i32.load
																	local.set $$t4
																	local.get $$t4
																	i32.const 1
																	i32.sub
																	local.set $$t5
																	local.get $$t0
																	local.get $$t5
																	i32.lt_s
																	local.set $$t6
																	local.get $$t6
																	if
																		i32.const 1
																		local.set $$block_selector
																		br $$BlockDisp
																	else
																		i32.const 2
																		local.set $$block_selector
																		br $$BlockDisp
																	end
																end
																i32.const 4
																local.set $$current_block
																i32.const 14924
																local.get $$t2
																i32.store
																br $$Block_4
															end
															i32.const 5
															local.set $$current_block
															br $$Block_8
														end
														i32.const 6
														local.set $$current_block
														i32.const 14928
														i32.load
														local.set $$t7
														local.get $$t7
														i32.const 1
														i32.sub
														local.set $$t8
														local.get $$t2
														local.get $$t8
														i32.lt_s
														local.set $$t9
														local.get $$t9
														if
															i32.const 4
															local.set $$block_selector
															br $$BlockDisp
														else
															i32.const 5
															local.set $$block_selector
															br $$BlockDisp
														end
													end
													i32.const 7
													local.set $$current_block
													br $$Block_11
												end
												i32.const 8
												local.set $$current_block
												i32.const 14904
												i32.load
												call $runtime.Block.Retain
												i32.const 14904
												i32.load offset=4
												local.set $$t10.1
												local.get $$t10.0
												call $runtime.Block.Release
												local.set $$t10.0
												local.get $$t10.0
												local.get $$t10.1
												call $w4app.drawFrambuffer
												br $$BlockFnBody
											end
											local.get $$current_block
											i32.const 5
											i32.eq
											if(result i32)
												i32.const 0
											else
												local.get $$t11
											end
											local.set $$t12
											i32.const 9
											local.set $$current_block
											i32.const 14932
											i32.load
											local.set $$t13
											local.get $$t12
											local.get $$t13
											i32.lt_s
											local.set $$t14
											local.get $$t14
											if
												i32.const 7
												local.set $$block_selector
												br $$BlockDisp
											else
												i32.const 8
												local.set $$block_selector
												br $$BlockDisp
											end
										end
										i32.const 10
										local.set $$current_block
										i32.const 14920
										i32.load
										local.set $$t15
										i32.const 14924
										i32.load
										local.set $$t16
										local.get $$t12
										local.get $$t17
										local.get $$t15
										local.get $$t16
										call $w4app.genColorGray
										local.set $$t18
										local.get $$t18
										i32.const 0
										i32.eq
										i32.eqz
										local.set $$t19
										local.get $$t19
										if
											br $$Block_12
										else
											br $$Block_14
										end
									end
									i32.const 11
									local.set $$current_block
									local.get $$t12
									i32.const 1
									i32.add
									local.set $$t11
									i32.const 9
									local.set $$block_selector
									br $$BlockDisp
								end
								local.get $$current_block
								i32.const 7
								i32.eq
								if(result i32)
									i32.const 0
								else
									local.get $$t20
								end
								local.set $$t17
								i32.const 12
								local.set $$current_block
								i32.const 14928
								i32.load
								local.set $$t21
								local.get $$t17
								local.get $$t21
								i32.lt_s
								local.set $$t22
								local.get $$t22
								if
									i32.const 10
									local.set $$block_selector
									br $$BlockDisp
								else
									i32.const 11
									local.set $$block_selector
									br $$BlockDisp
								end
							end
							i32.const 13
							local.set $$current_block
							i32.const 14904
							i32.load
							call $runtime.Block.Retain
							i32.const 14904
							i32.load offset=4
							local.set $$t23.1
							local.get $$t23.0
							call $runtime.Block.Release
							local.set $$t23.0
							local.get $$t23.0
							local.get $$t23.1
							local.get $$t12
							local.get $$t17
							i32.const 1
							call $w4app.BitImage.Set
							br $$Block_13
						end
						i32.const 14
						local.set $$current_block
						local.get $$t17
						i32.const 1
						i32.add
						local.set $$t20
						i32.const 12
						local.set $$block_selector
						br $$BlockDisp
					end
					i32.const 15
					local.set $$current_block
					i32.const 14904
					i32.load
					call $runtime.Block.Retain
					i32.const 14904
					i32.load offset=4
					local.set $$t24.1
					local.get $$t24.0
					call $runtime.Block.Release
					local.set $$t24.0
					local.get $$t24.0
					local.get $$t24.1
					local.get $$t12
					local.get $$t17
					i32.const 0
					call $w4app.BitImage.Set
					i32.const 14
					local.set $$block_selector
					br $$BlockDisp
				end
			end
		end
		local.get $$t10.0
		call $runtime.Block.Release
		local.get $$t23.0
		call $runtime.Block.Release
		local.get $$t24.0
		call $runtime.Block.Release
	)
	(func $w4app.Update (export "update")
		(local $$block_selector i32)
		(local $$current_block i32)
		block $$BlockFnBody
			loop $$BlockDisp
				block $$Block_0
					block $$BlockSel
						local.get $$block_selector
						br_table  0 0
					end
					i32.const 0
					local.set $$current_block
					call $w4app.StepGame
					br $$BlockFnBody
				end
			end
		end
	)
	(func $w4app._sq (param $x f32) (result f32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$ret_0 f32)
		(local $$t0 f32)
		block $$BlockFnBody
			loop $$BlockDisp
				block $$Block_0
					block $$BlockSel
						local.get $$block_selector
						br_table  0 0
					end
					i32.const 0
					local.set $$current_block
					local.get $x
					local.get $x
					f32.mul
					local.set $$t0
					local.get $$t0
					local.set $$ret_0
					br $$BlockFnBody
				end
			end
		end
		local.get $$ret_0
	)
	(func $w4app._sqrt (param $x f32) (result f32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$ret_0 f32)
		(local $$t0 f64)
		(local $$t1 f64)
		(local $$t2 f32)
		block $$BlockFnBody
			loop $$BlockDisp
				block $$Block_0
					block $$BlockSel
						local.get $$block_selector
						br_table  0 0
					end
					i32.const 0
					local.set $$current_block
					local.get $x
					f64.promote_f32
					local.set $$t0
					local.get $$t0
					call $math.Sqrt
					local.set $$t1
					local.get $$t1
					f32.demote_f64
					local.set $$t2
					local.get $$t2
					local.set $$ret_0
					br $$BlockFnBody
				end
			end
		end
		local.get $$ret_0
	)
	(func $w4app.clearBit (param $n i32) (param $pos i32) (result i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$ret_0 i32)
		(local $$t0 i32)
		(local $$t1 i32)
		(local $$t2 i32)
		(local $$t3 i32)
		block $$BlockFnBody
			loop $$BlockDisp
				block $$Block_0
					block $$BlockSel
						local.get $$block_selector
						br_table  0 0
					end
					i32.const 0
					local.set $$current_block
					i32.const 1
					local.get $pos
					i32.shl
					local.set $$t0
					i32.const -1
					local.get $$t0
					i32.xor
					local.set $$t1
					local.get $$t1
					local.set $$t2
					local.get $n
					local.get $$t2
					i32.and
					local.set $$t3
					local.get $$t3
					local.set $$ret_0
					br $$BlockFnBody
				end
			end
		end
		local.get $$ret_0
	)
	(func $w4app.drawFrambuffer (param $m.0 i32) (param $m.1 i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$t0 i32)
		(local $$t1 i32)
		(local $$t2 i32)
		(local $$t3 i32)
		(local $$t4 i32)
		(local $$t5 i32)
		(local $$t6 i32)
		(local $$t7 i32)
		(local $$t8 i32)
		(local $$t9.0 i32)
		(local $$t9.1 i32)
		(local $$t10.0 i32)
		(local $$t10.1 i32)
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
															br_table  0 1 2 3 4 5 6 7 8 9 0
														end
														i32.const 0
														local.set $$current_block
														br $$Block_2
													end
													i32.const 1
													local.set $$current_block
													br $$Block_5
												end
												i32.const 2
												local.set $$current_block
												br $$BlockFnBody
											end
											local.get $$current_block
											i32.const 0
											i32.eq
											if(result i32)
												i32.const 0
											else
												local.get $$t0
											end
											local.set $$t1
											i32.const 3
											local.set $$current_block
											i32.const 14932
											i32.load
											local.set $$t2
											local.get $$t1
											local.get $$t2
											i32.lt_s
											local.set $$t3
											local.get $$t3
											if
												i32.const 1
												local.set $$block_selector
												br $$BlockDisp
											else
												i32.const 2
												local.set $$block_selector
												br $$BlockDisp
											end
										end
										i32.const 4
										local.set $$current_block
										local.get $m.0
										local.get $m.1
										local.get $$t1
										local.get $$t4
										call $w4app.BitImage.At
										local.set $$t5
										local.get $$t5
										if
											br $$Block_6
										else
											br $$Block_8
										end
									end
									i32.const 5
									local.set $$current_block
									local.get $$t1
									i32.const 1
									i32.add
									local.set $$t0
									i32.const 3
									local.set $$block_selector
									br $$BlockDisp
								end
								local.get $$current_block
								i32.const 1
								i32.eq
								if(result i32)
									i32.const 0
								else
									local.get $$t6
								end
								local.set $$t4
								i32.const 6
								local.set $$current_block
								i32.const 14928
								i32.load
								local.set $$t7
								local.get $$t4
								local.get $$t7
								i32.lt_s
								local.set $$t8
								local.get $$t8
								if
									i32.const 4
									local.set $$block_selector
									br $$BlockDisp
								else
									i32.const 5
									local.set $$block_selector
									br $$BlockDisp
								end
							end
							i32.const 7
							local.set $$current_block
							i32.const 14912
							i32.load
							call $runtime.Block.Retain
							i32.const 14912
							i32.load offset=4
							local.set $$t9.1
							local.get $$t9.0
							call $runtime.Block.Release
							local.set $$t9.0
							local.get $$t9.0
							local.get $$t9.1
							local.get $$t1
							local.get $$t4
							i32.const 4
							call $w4app.Framebuffer.Set
							br $$Block_7
						end
						i32.const 8
						local.set $$current_block
						local.get $$t4
						i32.const 1
						i32.add
						local.set $$t6
						i32.const 6
						local.set $$block_selector
						br $$BlockDisp
					end
					i32.const 9
					local.set $$current_block
					i32.const 14912
					i32.load
					call $runtime.Block.Retain
					i32.const 14912
					i32.load offset=4
					local.set $$t10.1
					local.get $$t10.0
					call $runtime.Block.Release
					local.set $$t10.0
					local.get $$t10.0
					local.get $$t10.1
					local.get $$t1
					local.get $$t4
					i32.const 2
					call $w4app.Framebuffer.Set
					i32.const 8
					local.set $$block_selector
					br $$BlockDisp
				end
			end
		end
		local.get $$t9.0
		call $runtime.Block.Release
		local.get $$t10.0
		call $runtime.Block.Release
	)
	(func $w4app.genColorGray (param $i i32) (param $j i32) (param $center_x i32) (param $center_y i32) (result i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$ret_0 i32)
		(local $$t0 i32)
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
		(local $$t11 i32)
		(local $$t12 i32)
		(local $$t13 f32)
		(local $$t14 f32)
		(local $$t15 i32)
		(local $$t16 i32)
		(local $$t17 i32)
		(local $$t18 f32)
		(local $$t19 f32)
		(local $$t20 f32)
		(local $$t21 f32)
		(local $$t22 f32)
		(local $$t23 f32)
		(local $$t24 f32)
		(local $$t25 f32)
		(local $$t26 f32)
		(local $$t27 f32)
		(local $$t28 f32)
		(local $$t29 f32)
		(local $$t30 f32)
		(local $$t31 f32)
		(local $$t32 f32)
		(local $$t33 f32)
		(local $$t34 f32)
		(local $$t35 f32)
		(local $$t36 f32)
		(local $$t37 i32)
		(local $$t38 f32)
		(local $$t39 f32)
		(local $$t40 f32)
		(local $$t41 f32)
		(local $$t42 i32)
		(local $$t43 i32)
		(local $$t44 i32)
		(local $$t45 i32)
		(local $$t46 i32)
		(local $$t47 i32)
		(local $$t48 i32)
		(local $$t49 i32)
		(local $$t50 i32)
		(local $$t51 i32)
		(local $$t52 i32)
		(local $$t53 i32)
		block $$BlockFnBody
			loop $$BlockDisp
				block $$Block_2
					block $$Block_1
						block $$Block_0
							block $$BlockSel
								local.get $$block_selector
								br_table  0 1 2 0
							end
							i32.const 0
							local.set $$current_block
							i32.const 14932
							i32.load
							local.set $$t0
							local.get $i
							f32.convert_i32_s
							local.set $$t1
							local.get $center_x
							f32.convert_i32_s
							local.set $$t2
							local.get $$t1
							local.get $$t2
							f32.sub
							local.set $$t3
							local.get $$t3
							call $w4app._sq
							local.set $$t4
							local.get $j
							f32.convert_i32_s
							local.set $$t5
							local.get $center_y
							f32.convert_i32_s
							local.set $$t6
							local.get $$t5
							local.get $$t6
							f32.sub
							local.set $$t7
							local.get $$t7
							call $w4app._sq
							local.set $$t8
							local.get $$t4
							local.get $$t8
							f32.add
							local.set $$t9
							local.get $$t9
							call $w4app._sqrt
							local.set $$t10
							local.get $$t0
							i32.const 2
							i32.div_s
							local.set $$t11
							local.get $$t11
							i32.const 2
							i32.div_s
							local.set $$t12
							local.get $$t12
							f32.convert_i32_s
							local.set $$t13
							local.get $$t10
							local.get $$t13
							f32.div
							local.set $$t14
							local.get $$t14
							f32.const 1
							f32.lt
							local.set $$t15
							local.get $$t15
							if
								br $$Block_0
							else
								br $$Block_1
							end
						end
						i32.const 1
						local.set $$current_block
						local.get $i
						local.get $center_x
						i32.sub
						local.set $$t16
						local.get $j
						local.get $center_y
						i32.sub
						local.set $$t17
						local.get $$t14
						call $w4app._sq
						local.set $$t18
						f32.const 1
						local.get $$t18
						f32.sub
						local.set $$t19
						local.get $$t19
						call $w4app._sqrt
						local.set $$t20
						i32.const 3
						f32.convert_i32_s
						local.set $$t21
						local.get $$t14
						local.get $$t21
						f32.div
						local.set $$t22
						local.get $$t22
						call $w4app._sq
						local.set $$t23
						f32.const 1
						local.get $$t23
						f32.sub
						local.set $$t24
						local.get $$t24
						call $w4app._sqrt
						local.set $$t25
						local.get $$t20
						local.get $$t25
						f32.mul
						local.set $$t26
						local.get $$t14
						call $w4app._sq
						local.set $$t27
						i32.const 3
						f32.convert_i32_s
						local.set $$t28
						local.get $$t27
						local.get $$t28
						f32.div
						local.set $$t29
						local.get $$t26
						local.get $$t29
						f32.add
						local.set $$t30
						i32.const 3
						f32.convert_i32_s
						local.set $$t31
						local.get $$t30
						local.get $$t31
						f32.mul
						local.set $$t32
						local.get $$t16
						f32.convert_i32_s
						local.set $$t33
						local.get $$t33
						local.get $$t32
						f32.div
						local.set $$t34
						local.get $center_x
						f32.convert_i32_s
						local.set $$t35
						local.get $$t34
						local.get $$t35
						f32.add
						local.set $$t36
						local.get $$t36
						i32.trunc_f32_s
						local.set $$t37
						local.get $$t17
						f32.convert_i32_s
						local.set $$t38
						local.get $$t38
						local.get $$t32
						f32.div
						local.set $$t39
						local.get $center_y
						f32.convert_i32_s
						local.set $$t40
						local.get $$t39
						local.get $$t40
						f32.add
						local.set $$t41
						local.get $$t41
						i32.trunc_f32_s
						local.set $$t42
						br $$Block_1
					end
					local.get $$current_block
					i32.const 0
					i32.eq
					if(result i32)
						local.get $i
					else
						local.get $$t37
					end
					local.get $$current_block
					i32.const 0
					i32.eq
					if(result i32)
						local.get $j
					else
						local.get $$t42
					end
					local.set $$t44
					local.set $$t43
					i32.const 2
					local.set $$current_block
					local.get $$t0
					i32.const 2
					i32.div_s
					local.set $$t45
					local.get $$t45
					i32.const 8
					i32.div_s
					local.set $$t46
					local.get $$t43
					local.get $$t46
					i32.div_s
					local.set $$t47
					local.get $$t0
					i32.const 2
					i32.div_s
					local.set $$t48
					local.get $$t48
					i32.const 8
					i32.div_s
					local.set $$t49
					local.get $$t44
					local.get $$t49
					i32.div_s
					local.set $$t50
					local.get $$t47
					local.get $$t50
					i32.add
					local.set $$t51
					local.get $$t51
					i32.const 2
					i32.rem_s
					local.set $$t52
					local.get $$t52
					i32.const 255
					i32.mul
					local.set $$t53
					local.get $$t53
					local.set $$ret_0
					br $$BlockFnBody
				end
			end
		end
		local.get $$ret_0
	)
	(func $w4app.hasBit (param $n i32) (param $pos i32) (result i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$ret_0 i32)
		(local $$t0 i32)
		(local $$t1 i32)
		(local $$t2 i32)
		block $$BlockFnBody
			loop $$BlockDisp
				block $$Block_0
					block $$BlockSel
						local.get $$block_selector
						br_table  0 0
					end
					i32.const 0
					local.set $$current_block
					i32.const 1
					local.get $pos
					i32.shl
					local.set $$t0
					local.get $n
					local.get $$t0
					i32.and
					local.set $$t1
					local.get $$t1
					i32.const 0
					i32.eq
					i32.eqz
					local.set $$t2
					local.get $$t2
					local.set $$ret_0
					br $$BlockFnBody
				end
			end
		end
		local.get $$ret_0
	)
	(func $w4app.init
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$t0 i32)
		(local $$t1 i32)
		(local $$t2 i32)
		(local $$t3.0 i32)
		(local $$t3.1 i32)
		(local $$t4.0 i32)
		(local $$t4.1 i32)
		block $$BlockFnBody
			loop $$BlockDisp
				block $$Block_2
					block $$Block_1
						block $$Block_0
							block $$BlockSel
								local.get $$block_selector
								br_table  0 1 2 0
							end
							i32.const 0
							local.set $$current_block
							global.get $w4app.init$guard
							local.set $$t0
							local.get $$t0
							if
								br $$Block_1
							else
								br $$Block_0
							end
						end
						i32.const 1
						local.set $$current_block
						i32.const 1
						global.set $w4app.init$guard
						call $runtime.init
						call $math.init
						call $syscall$wasm4.init
						i32.const 14932
						i32.load
						local.set $$t1
						i32.const 14928
						i32.load
						local.set $$t2
						local.get $$t1
						local.get $$t2
						call $w4app.NewBitImage
						local.set $$t3.1
						local.get $$t3.0
						call $runtime.Block.Release
						local.set $$t3.0
						i32.const 14904
						local.get $$t3.0
						call $runtime.Block.Retain
						i32.const 14904
						i32.load align=1
						call $runtime.Block.Release
						i32.store align=1
						i32.const 14904
						local.get $$t3.1
						i32.store offset=4
						call $w4app.FramebufferInstance
						local.set $$t4.1
						local.get $$t4.0
						call $runtime.Block.Release
						local.set $$t4.0
						i32.const 14912
						local.get $$t4.0
						call $runtime.Block.Retain
						i32.const 14912
						i32.load align=1
						call $runtime.Block.Release
						i32.store align=1
						i32.const 14912
						local.get $$t4.1
						i32.store offset=4
						br $$Block_1
					end
					i32.const 2
					local.set $$current_block
					br $$BlockFnBody
				end
			end
		end
		local.get $$t3.0
		call $runtime.Block.Release
		local.get $$t4.0
		call $runtime.Block.Release
	)
	(func $w4app.setBit (param $n i32) (param $pos i32) (result i32)
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
						br_table  0 0
					end
					i32.const 0
					local.set $$current_block
					i32.const 1
					local.get $pos
					i32.shl
					local.set $$t0
					local.get $n
					local.get $$t0
					i32.or
					local.set $$t1
					local.get $$t1
					local.set $$ret_0
					br $$BlockFnBody
				end
			end
		end
		local.get $$ret_0
	)
	(func $w4app.BitImage.At (param $this.0 i32) (param $this.1 i32) (param $x i32) (param $y i32) (result i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$ret_0 i32)
		(local $$t0.0 i32)
		(local $$t0.1 i32)
		(local $$t1 i32)
		(local $$t2 i32)
		(local $$t3 i32)
		(local $$t4 i32)
		(local $$t5.0 i32)
		(local $$t5.1 i32)
		(local $$t6.0 i32)
		(local $$t6.1 i32)
		(local $$t6.2 i32)
		(local $$t6.3 i32)
		(local $$t7.0 i32)
		(local $$t7.1 i32)
		(local $$t8 i32)
		(local $$t9 i32)
		(local $$t10 i32)
		(local $$t11 i32)
		(local $$t12 i32)
		block $$BlockFnBody
			loop $$BlockDisp
				block $$Block_0
					block $$BlockSel
						local.get $$block_selector
						br_table  0 0
					end
					i32.const 0
					local.set $$current_block
					local.get $this.0
					call $runtime.Block.Retain
					local.get $this.1
					i32.const 0
					i32.add
					local.set $$t0.1
					local.get $$t0.0
					call $runtime.Block.Release
					local.set $$t0.0
					local.get $$t0.1
					i32.load
					local.set $$t1
					local.get $y
					local.get $$t1
					i32.mul
					local.set $$t2
					local.get $$t2
					local.get $x
					i32.add
					local.set $$t3
					local.get $$t3
					i32.const 8
					i32.div_s
					local.set $$t4
					local.get $this.0
					call $runtime.Block.Retain
					local.get $this.1
					i32.const 8
					i32.add
					local.set $$t5.1
					local.get $$t5.0
					call $runtime.Block.Release
					local.set $$t5.0
					local.get $$t5.1
					i32.load
					call $runtime.Block.Retain
					local.get $$t5.1
					i32.load offset=4
					local.get $$t5.1
					i32.load offset=8
					local.get $$t5.1
					i32.load offset=12
					local.set $$t6.3
					local.set $$t6.2
					local.set $$t6.1
					local.get $$t6.0
					call $runtime.Block.Release
					local.set $$t6.0
					local.get $$t6.0
					call $runtime.Block.Retain
					local.get $$t6.1
					i32.const 1
					local.get $$t4
					i32.mul
					i32.add
					local.set $$t7.1
					local.get $$t7.0
					call $runtime.Block.Release
					local.set $$t7.0
					local.get $$t7.1
					i32.load8_u align=1
					local.set $$t8
					local.get $$t8
					local.set $$t9
					local.get $x
					i32.const 8
					i32.rem_s
					local.set $$t10
					local.get $$t10
					local.set $$t11
					local.get $$t9
					local.get $$t11
					call $w4app.hasBit
					local.set $$t12
					local.get $$t12
					local.set $$ret_0
					br $$BlockFnBody
				end
			end
		end
		local.get $$ret_0
		local.get $$t0.0
		call $runtime.Block.Release
		local.get $$t5.0
		call $runtime.Block.Release
		local.get $$t6.0
		call $runtime.Block.Release
		local.get $$t7.0
		call $runtime.Block.Release
	)
	(func $w4app.BitImage.Set (param $this.0 i32) (param $this.1 i32) (param $x i32) (param $y i32) (param $c i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$t0 i32)
		(local $$t1.0 i32)
		(local $$t1.1 i32)
		(local $$t2 i32)
		(local $$t3 i32)
		(local $$t4 i32)
		(local $$t5 i32)
		(local $$t6.0 i32)
		(local $$t6.1 i32)
		(local $$t7 i32)
		(local $$t8 i32)
		(local $$t9 i32)
		(local $$t10.0 i32)
		(local $$t10.1 i32)
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
		(local $$t17.0 i32)
		(local $$t17.1 i32)
		(local $$t17.2 i32)
		(local $$t17.3 i32)
		(local $$t18.0 i32)
		(local $$t18.1 i32)
		(local $$t19 i32)
		(local $$t20 i32)
		(local $$t21 i32)
		(local $$t22 i32)
		(local $$t23 i32)
		(local $$t24 i32)
		(local $$t25.0 i32)
		(local $$t25.1 i32)
		(local $$t26.0 i32)
		(local $$t26.1 i32)
		(local $$t26.2 i32)
		(local $$t26.3 i32)
		(local $$t27.0 i32)
		(local $$t27.1 i32)
		(local $$t28.0 i32)
		(local $$t28.1 i32)
		(local $$t29.0 i32)
		(local $$t29.1 i32)
		(local $$t29.2 i32)
		(local $$t29.3 i32)
		(local $$t30.0 i32)
		(local $$t30.1 i32)
		(local $$t31 i32)
		(local $$t32 i32)
		(local $$t33 i32)
		(local $$t34 i32)
		(local $$t35 i32)
		(local $$t36 i32)
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
														br_table  0 1 2 3 4 5 6 7 8 0
													end
													i32.const 0
													local.set $$current_block
													local.get $x
													i32.const 0
													i32.lt_s
													local.set $$t0
													local.get $$t0
													if
														br $$Block_0
													else
														br $$Block_4
													end
												end
												i32.const 1
												local.set $$current_block
												br $$BlockFnBody
											end
											i32.const 2
											local.set $$current_block
											local.get $this.0
											call $runtime.Block.Retain
											local.get $this.1
											i32.const 0
											i32.add
											local.set $$t1.1
											local.get $$t1.0
											call $runtime.Block.Release
											local.set $$t1.0
											local.get $$t1.1
											i32.load
											local.set $$t2
											local.get $y
											local.get $$t2
											i32.mul
											local.set $$t3
											local.get $$t3
											local.get $x
											i32.add
											local.set $$t4
											local.get $$t4
											i32.const 8
											i32.div_s
											local.set $$t5
											local.get $c
											if
												br $$Block_5
											else
												br $$Block_7
											end
										end
										i32.const 3
										local.set $$current_block
										local.get $this.0
										call $runtime.Block.Retain
										local.get $this.1
										i32.const 4
										i32.add
										local.set $$t6.1
										local.get $$t6.0
										call $runtime.Block.Release
										local.set $$t6.0
										local.get $$t6.1
										i32.load
										local.set $$t7
										local.get $y
										local.get $$t7
										i32.ge_s
										local.set $$t8
										local.get $$t8
										if
											i32.const 1
											local.set $$block_selector
											br $$BlockDisp
										else
											i32.const 2
											local.set $$block_selector
											br $$BlockDisp
										end
									end
									i32.const 4
									local.set $$current_block
									local.get $y
									i32.const 0
									i32.lt_s
									local.set $$t9
									local.get $$t9
									if
										i32.const 1
										local.set $$block_selector
										br $$BlockDisp
									else
										i32.const 3
										local.set $$block_selector
										br $$BlockDisp
									end
								end
								i32.const 5
								local.set $$current_block
								local.get $this.0
								call $runtime.Block.Retain
								local.get $this.1
								i32.const 0
								i32.add
								local.set $$t10.1
								local.get $$t10.0
								call $runtime.Block.Release
								local.set $$t10.0
								local.get $$t10.1
								i32.load
								local.set $$t11
								local.get $x
								local.get $$t11
								i32.ge_s
								local.set $$t12
								local.get $$t12
								if
									i32.const 1
									local.set $$block_selector
									br $$BlockDisp
								else
									i32.const 4
									local.set $$block_selector
									br $$BlockDisp
								end
							end
							i32.const 6
							local.set $$current_block
							local.get $this.0
							call $runtime.Block.Retain
							local.get $this.1
							i32.const 8
							i32.add
							local.set $$t13.1
							local.get $$t13.0
							call $runtime.Block.Release
							local.set $$t13.0
							local.get $$t13.1
							i32.load
							call $runtime.Block.Retain
							local.get $$t13.1
							i32.load offset=4
							local.get $$t13.1
							i32.load offset=8
							local.get $$t13.1
							i32.load offset=12
							local.set $$t14.3
							local.set $$t14.2
							local.set $$t14.1
							local.get $$t14.0
							call $runtime.Block.Release
							local.set $$t14.0
							local.get $$t14.0
							call $runtime.Block.Retain
							local.get $$t14.1
							i32.const 1
							local.get $$t5
							i32.mul
							i32.add
							local.set $$t15.1
							local.get $$t15.0
							call $runtime.Block.Release
							local.set $$t15.0
							local.get $this.0
							call $runtime.Block.Retain
							local.get $this.1
							i32.const 8
							i32.add
							local.set $$t16.1
							local.get $$t16.0
							call $runtime.Block.Release
							local.set $$t16.0
							local.get $$t16.1
							i32.load
							call $runtime.Block.Retain
							local.get $$t16.1
							i32.load offset=4
							local.get $$t16.1
							i32.load offset=8
							local.get $$t16.1
							i32.load offset=12
							local.set $$t17.3
							local.set $$t17.2
							local.set $$t17.1
							local.get $$t17.0
							call $runtime.Block.Release
							local.set $$t17.0
							local.get $$t17.0
							call $runtime.Block.Retain
							local.get $$t17.1
							i32.const 1
							local.get $$t5
							i32.mul
							i32.add
							local.set $$t18.1
							local.get $$t18.0
							call $runtime.Block.Release
							local.set $$t18.0
							local.get $$t18.1
							i32.load8_u align=1
							local.set $$t19
							local.get $$t19
							local.set $$t20
							local.get $x
							i32.const 8
							i32.rem_s
							local.set $$t21
							local.get $$t21
							local.set $$t22
							local.get $$t20
							local.get $$t22
							call $w4app.setBit
							local.set $$t23
							local.get $$t23
							i32.const 255
							i32.and
							local.set $$t24
							local.get $$t15.1
							local.get $$t24
							i32.store8 align=1
							br $$Block_6
						end
						i32.const 7
						local.set $$current_block
						br $$BlockFnBody
					end
					i32.const 8
					local.set $$current_block
					local.get $this.0
					call $runtime.Block.Retain
					local.get $this.1
					i32.const 8
					i32.add
					local.set $$t25.1
					local.get $$t25.0
					call $runtime.Block.Release
					local.set $$t25.0
					local.get $$t25.1
					i32.load
					call $runtime.Block.Retain
					local.get $$t25.1
					i32.load offset=4
					local.get $$t25.1
					i32.load offset=8
					local.get $$t25.1
					i32.load offset=12
					local.set $$t26.3
					local.set $$t26.2
					local.set $$t26.1
					local.get $$t26.0
					call $runtime.Block.Release
					local.set $$t26.0
					local.get $$t26.0
					call $runtime.Block.Retain
					local.get $$t26.1
					i32.const 1
					local.get $$t5
					i32.mul
					i32.add
					local.set $$t27.1
					local.get $$t27.0
					call $runtime.Block.Release
					local.set $$t27.0
					local.get $this.0
					call $runtime.Block.Retain
					local.get $this.1
					i32.const 8
					i32.add
					local.set $$t28.1
					local.get $$t28.0
					call $runtime.Block.Release
					local.set $$t28.0
					local.get $$t28.1
					i32.load
					call $runtime.Block.Retain
					local.get $$t28.1
					i32.load offset=4
					local.get $$t28.1
					i32.load offset=8
					local.get $$t28.1
					i32.load offset=12
					local.set $$t29.3
					local.set $$t29.2
					local.set $$t29.1
					local.get $$t29.0
					call $runtime.Block.Release
					local.set $$t29.0
					local.get $$t29.0
					call $runtime.Block.Retain
					local.get $$t29.1
					i32.const 1
					local.get $$t5
					i32.mul
					i32.add
					local.set $$t30.1
					local.get $$t30.0
					call $runtime.Block.Release
					local.set $$t30.0
					local.get $$t30.1
					i32.load8_u align=1
					local.set $$t31
					local.get $$t31
					local.set $$t32
					local.get $x
					i32.const 8
					i32.rem_s
					local.set $$t33
					local.get $$t33
					local.set $$t34
					local.get $$t32
					local.get $$t34
					call $w4app.clearBit
					local.set $$t35
					local.get $$t35
					i32.const 255
					i32.and
					local.set $$t36
					local.get $$t27.1
					local.get $$t36
					i32.store8 align=1
					i32.const 7
					local.set $$block_selector
					br $$BlockDisp
				end
			end
		end
		local.get $$t1.0
		call $runtime.Block.Release
		local.get $$t6.0
		call $runtime.Block.Release
		local.get $$t10.0
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
		local.get $$t30.0
		call $runtime.Block.Release
	)
	(func $w4app.Framebuffer.Set (param $this.0 i32) (param $this.1 i32) (param $x i32) (param $y i32) (param $v i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$t0 i32)
		(local $$t1 i32)
		(local $$t2 i32)
		(local $$t3 i32)
		(local $$t4 i32)
		(local $$t5 i32)
		(local $$t6 i32)
		(local $$t7.0 i32)
		(local $$t7.1 i32)
		(local $$t8.0 i32)
		(local $$t8.1 i32)
		(local $$t8.2 i32)
		(local $$t8.3 i32)
		(local $$t9.0 i32)
		(local $$t9.1 i32)
		(local $$t10 i32)
		(local $$t11.0 i32)
		(local $$t11.1 i32)
		(local $$t12.0 i32)
		(local $$t12.1 i32)
		(local $$t12.2 i32)
		(local $$t12.3 i32)
		(local $$t13.0 i32)
		(local $$t13.1 i32)
		(local $$t14 i32)
		(local $$t15 i32)
		(local $$t16 i32)
		(local $$t17 i32)
		block $$BlockFnBody
			loop $$BlockDisp
				block $$Block_0
					block $$BlockSel
						local.get $$block_selector
						br_table  0 0
					end
					i32.const 0
					local.set $$current_block
					local.get $y
					i32.const 160
					i32.mul
					local.set $$t0
					local.get $$t0
					local.get $x
					i32.add
					local.set $$t1
					local.get $$t1
					i32.const 4
					i32.div_s
					local.set $$t2
					local.get $x
					i32.const 4
					i32.rem_s
					local.set $$t3
					local.get $$t3
					i32.const 2
					i32.mul
					local.set $$t4
					local.get $$t4
					i32.const 255
					i32.and
					local.set $$t5
					i32.const 3
					local.get $$t5
					i32.shl
					i32.const 255
					i32.and
					local.set $$t6
					local.get $this.0
					call $runtime.Block.Retain
					local.get $this.1
					i32.const 0
					i32.add
					local.set $$t7.1
					local.get $$t7.0
					call $runtime.Block.Release
					local.set $$t7.0
					local.get $$t7.1
					i32.load
					call $runtime.Block.Retain
					local.get $$t7.1
					i32.load offset=4
					local.get $$t7.1
					i32.load offset=8
					local.get $$t7.1
					i32.load offset=12
					local.set $$t8.3
					local.set $$t8.2
					local.set $$t8.1
					local.get $$t8.0
					call $runtime.Block.Release
					local.set $$t8.0
					local.get $$t8.0
					call $runtime.Block.Retain
					local.get $$t8.1
					i32.const 1
					local.get $$t2
					i32.mul
					i32.add
					local.set $$t9.1
					local.get $$t9.0
					call $runtime.Block.Release
					local.set $$t9.0
					local.get $v
					local.get $$t5
					i32.shl
					i32.const 255
					i32.and
					local.set $$t10
					local.get $this.0
					call $runtime.Block.Retain
					local.get $this.1
					i32.const 0
					i32.add
					local.set $$t11.1
					local.get $$t11.0
					call $runtime.Block.Release
					local.set $$t11.0
					local.get $$t11.1
					i32.load
					call $runtime.Block.Retain
					local.get $$t11.1
					i32.load offset=4
					local.get $$t11.1
					i32.load offset=8
					local.get $$t11.1
					i32.load offset=12
					local.set $$t12.3
					local.set $$t12.2
					local.set $$t12.1
					local.get $$t12.0
					call $runtime.Block.Release
					local.set $$t12.0
					local.get $$t12.0
					call $runtime.Block.Retain
					local.get $$t12.1
					i32.const 1
					local.get $$t2
					i32.mul
					i32.add
					local.set $$t13.1
					local.get $$t13.0
					call $runtime.Block.Release
					local.set $$t13.0
					local.get $$t13.1
					i32.load8_u align=1
					local.set $$t14
					i32.const -1
					local.get $$t6
					i32.xor
					i32.const 255
					i32.and
					local.set $$t15
					local.get $$t14
					local.get $$t15
					i32.and
					local.set $$t16
					local.get $$t10
					local.get $$t16
					i32.or
					local.set $$t17
					local.get $$t9.1
					local.get $$t17
					i32.store8 align=1
					br $$BlockFnBody
				end
			end
		end
		local.get $$t7.0
		call $runtime.Block.Release
		local.get $$t8.0
		call $runtime.Block.Release
		local.get $$t9.0
		call $runtime.Block.Release
		local.get $$t11.0
		call $runtime.Block.Release
		local.get $$t12.0
		call $runtime.Block.Release
		local.get $$t13.0
		call $runtime.Block.Release
	)
	(func $_start (export "_start")
		call $w4app.init
	)
	(func $_main (export "_main"))
	(data (i32.const 14784) "\24\24\77\61\64\73\24\24\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\9b\1a\86\a0\49\fa\a8\bd\05\3f\4e\7b\9d\ee\21\3e\c6\4b\ac\7e\4f\7e\92\be\f5\44\c8\19\a0\01\fa\3e\91\4f\c1\16\6c\c1\56\bf\4b\55\55\55\55\55\a5\3f\cd\9c\d1\1f\fd\d8\e5\3d\5d\1f\29\a9\e5\e5\5a\be\a1\48\7d\56\e3\1d\c7\3e\03\df\bf\19\a0\01\2a\bf\d0\f7\10\11\11\11\81\3f\48\55\55\55\55\55\c5\bf\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\50\00\00\00\50\00\00\00\a0\00\00\00\a0\00\00\00\30\61\73\73\65\72\74\20\66\61\69\6c\65\64\20\28\61\73\73\65\72\74\20\66\61\69\6c\65\64\3a\20\2b\69\29\6e\69\6c\20\6d\61\70\2e\6d\61\70\2e\77\61\3a\36\38\3a\38\70\61\6e\69\63\3a\20\74\72\75\65\66\61\6c\73\65\4e\61\4e\2b\49\6e\66\2d\49\6e\66\30\31\32\33\34\35\36\37\38\39\61\62\63\64\65\66\0a\5b\2f\5d\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\01\00\00\00\ff\ff\ff\ff\00\3b\00\00")
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
	(elem (i32.const 23) $$u8.$slice.underlying.$$OnFree)
	(elem (i32.const 24) $$w4app.Framebuffer.$$OnFree)
	(elem (i32.const 25) $$w4app.BitImage.$$OnFree)
)
