(module $__walang__
	(import "env" "traceUtf8" (func $runtime.traceUtf8 (param i32) (param i32)))
	(import "env" "blit" (func $syscall$wasm4.__import__blit (param i32) (param i32) (param i32) (param i32) (param i32) (param i32)))
	(import "env" "rect" (func $syscall$wasm4.__import__rect (param i32) (param i32) (param i32) (param i32)))
	(import "env" "tone" (func $syscall$wasm4.__import__tone (param i32) (param i32) (param i32) (param i32)))
	(import "env" "memory" (memory 1 2))
	(table 34 funcref)
	(type $$OnFree (func (param i32)))
	(type $$wa.runtime.comp (func (param i32) (param i32) (result i32)))
	(type $$$fnSig1 (func))
	(type $$$fnSig2 (func (param i32) (param i32) (result i64)))
	(type $$$fnSig3 (func (param i32) (param i32) (param i64)))
	(type $$$fnSig4 (func (param i32) (param i32) (result i64)))
	(type $$$fnSig5 (func (param i32) (param i32) (result i32 i32 i32)))
	(type $$$fnSig6 (func (param i32) (param i32)))
	(type $$$fnSig7 (func (param i32) (result i32)))
	(global $__stack_ptr (mut i32) (i32.const 14656))
	(global $__heap_lfixed_cap i32 (i32.const 0))
	(global $__heap_ptr (mut i32) (i32.const 0))
	(global $__heap_top (mut i32) (i32.const 0))
	(global $__heap_l128_freep (mut i32) (i32.const 0))
	(global $__heap_init_flag (mut i32) (i32.const 0))
	(global $$wa.runtime.closure_data (mut i32) (i32.const 0))
	(global $$wa.runtime._concretTypeCount (mut i32) (i32.const 3))
	(global $$wa.runtime._interfaceCount (mut i32) (i32.const 4))
	(global $$wa.runtime._itabsPtr (mut i32) (i32.const 29704))
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
	(global $math$rand.fe.0 i32 (i32.const 0))
	(global $math$rand.fe.1 i32 (i32.const 14904))
	(global $math$rand.fn.0 i32 (i32.const 0))
	(global $math$rand.fn.1 i32 (i32.const 15928))
	(global $math$rand.globalRand.0 i32 (i32.const 0))
	(global $math$rand.globalRand.1 i32 (i32.const 16440))
	(global $math$rand.init$guard (mut i32) (i32.const 0))
	(global $math$rand.ke.0 i32 (i32.const 0))
	(global $math$rand.ke.1 i32 (i32.const 16448))
	(global $math$rand.kn.0 i32 (i32.const 0))
	(global $math$rand.kn.1 i32 (i32.const 17472))
	(global $math$rand.rngCooked.0 i32 (i32.const 0))
	(global $math$rand.rngCooked.1 i32 (i32.const 17984))
	(global $math$rand.we.0 i32 (i32.const 0))
	(global $math$rand.we.1 i32 (i32.const 22840))
	(global $math$rand.wn.0 i32 (i32.const 0))
	(global $math$rand.wn.1 i32 (i32.const 23864))
	(global $errors.init$guard (mut i32) (i32.const 0))
	(global $syscall$wasm4.init$guard (mut i32) (i32.const 0))
	(global $w4snake.frameCount.0 i32 (i32.const 0))
	(global $w4snake.frameCount.1 i32 (i32.const 24376))
	(global $w4snake.fruit.0 i32 (i32.const 0))
	(global $w4snake.fruit.1 i32 (i32.const 24380))
	(global $w4snake.fruitSprite.0 i32 (i32.const 0))
	(global $w4snake.fruitSprite.1 i32 (i32.const 24384))
	(global $w4snake.init$guard (mut i32) (i32.const 0))
	(global $w4snake.randInt.0 i32 (i32.const 0))
	(global $w4snake.randInt.1 i32 (i32.const 24400))
	(global $w4snake.sfxEat.0 i32 (i32.const 0))
	(global $w4snake.sfxEat.1 i32 (i32.const 24412))
	(global $w4snake.snake.0 i32 (i32.const 0))
	(global $w4snake.snake.1 i32 (i32.const 24448))
	(global $runtime.zptr (mut i32) (i32.const 24760))
	(global $__heap_base i32 (i32.const 29776))
	(func $runtime.throw
		unreachable
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
	(func $$wa.runtime.queryIface_CommaOk (param $d.b i32) (param $d.d i32) (param $itab i32) (param $eq i32) (param $ihash i32) (result i32 i32 i32 i32 i32)
		(local $t i32)
		local.get $itab
		if(result i32 i32 i32 i32 i32)
			local.get $itab
			i32.load
			local.get $ihash
			i32.const 1
			call $runtime.getItab
			local.set $t
			local.get $t
			if(result i32 i32 i32 i32 i32)
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
	(func $$wa.runtime.string_to_ptr (param $b i32) (param $d i32) (param $l i32) (result i32)
		local.get $d
	)
	(func $$wa.runtime.string_to_iter (param $b i32) (param $d i32) (param $l i32) (result i32 i32 i32)
		local.get $d
		local.get $l
		i32.const 0
	)
	(func $$syscall/wasm4.__linkname__slice_data_ptr (param $b i32) (param $d i32) (param $l i32) (param $c i32) (result i32)
		local.get $d
	)
	(func $$syscall/wasm4.__linkname__make_slice (param $blk i32) (param $ptr i32) (param $len i32) (param $cap i32) (result i32 i32 i32 i32)
		local.get $blk
		local.get $ptr
		local.get $len
		local.get $cap
		return
	)
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
				end
				local.get $src
				i32.load8_u align=1
				local.set $item
				local.get $dest
				local.get $item
				i32.store8 align=1
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
			end
		end
		local.get $y.1
		local.set $src
		block $block3
			loop $loop3
				local.get $y_len
				i32.eqz
				if
					br $block3
				end
				local.get $src
				i32.load8_u align=1
				local.set $item
				local.get $dest
				local.get $item
				i32.store8 align=1
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
			end
		end
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
						br_table  0 0
					end
					i32.const 0
					local.set $$current_block
					global.get $$wa.runtime._itabsPtr
					local.set $$t0
					local.get $dhash
					i32.const 1
					i32.sub
					local.set $$t1
					global.get $$wa.runtime._interfaceCount
					local.set $$t2
					local.get $$t1
					local.get $$t2
					i32.mul
					local.set $$t3
					local.get $$t3
					local.get $ihash
					i32.sub
					local.set $$t4
					local.get $$t4
					i32.const 1
					i32.sub
					local.set $$t5
					local.get $$t5
					i32.const 4
					i32.mul
					local.set $$t6
					local.get $$t0
					local.get $$t6
					i32.add
					local.set $$t7
					local.get $$t7
					local.set $$t8
					local.get $$t8
					call $runtime.getU32
					local.set $$t9
					local.get $$t9
					local.set $$ret_0
					br $$BlockFnBody
				end
			end
		end
		local.get $$ret_0
	)
	(func $runtime.getU32 (param $addr i32) (result i32)
		local.get $addr
		i32.load
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
	(func $$runtime.panic_ (param $msg_ptr i32) (param $msg_len i32) (param $pos_msg_ptr i32) (param $pos_msg_len i32)
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
					i32.const 0
					i32.const 24529
					i32.const 7
					call $$runtime.waPrintString
					local.get $msg_ptr
					local.get $msg_len
					call $$runtime.waPuts
					i32.const 0
					i32.const 24490
					i32.const 2
					call $$runtime.waPrintString
					local.get $pos_msg_ptr
					local.get $pos_msg_len
					call $$runtime.waPuts
					i32.const 41
					call $$runtime.waPrintRune
					i32.const 10
					call $$runtime.waPrintRune
					i32.const 1
					call $$runtime.procExit
					br $$BlockFnBody
				end
			end
		end
	)
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
								br_table  0 1 2 0
							end
							i32.const 0
							local.set $$current_block
							local.get $s.2
							local.set $$t0
							local.get $$t0
							i32.const 0
							i32.gt_s
							local.set $$t1
							local.get $$t1
							if
								br $$Block_0
							else
								br $$Block_1
							end
						end
						i32.const 1
						local.set $$current_block
						local.get $s.0
						local.get $s.1
						local.get $s.2
						call $runtime.refToPtr_string
						local.set $$t2
						local.get $$t0
						local.set $$t3
						local.get $$t2
						local.get $$t3
						call $$runtime.waPuts
						br $$BlockFnBody
					end
					i32.const 2
					local.set $$current_block
					br $$BlockFnBody
				end
			end
		end
	)
	(func $$runtime.procExit (param $code i32)
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
					call $runtime.throw
					br $$BlockFnBody
				end
			end
		end
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
						br_table  0 0
					end
					i32.const 0
					local.set $$current_block
					local.get $s.0
					local.get $s.1
					local.get $s.2
					call $$wa.runtime.string_to_ptr
					local.set $$t0
					local.get $$t0
					local.set $$ret_0
					br $$BlockFnBody
				end
			end
		end
		local.get $$ret_0
	)
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
																		br_table  0 1 2 3 4 5 6 7 8 9 10 11 12 0
																	end
																	i32.const 0
																	local.set $$current_block
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
																	local.get $r
																	local.set $$t2
																	local.get $$t2
																	i32.const 127
																	i32.le_u
																	local.set $$t3
																	local.get $$t3
																	if
																		br $$Block_1
																	else
																		br $$Block_3
																	end
																end
																local.get $$current_block
																i32.const 2
																i32.eq
																if(result i32)
																	i32.const 1
																else
																	local.get $$current_block
																	i32.const 3
																	i32.eq
																	if(result i32)
																		i32.const 2
																	else
																		local.get $$current_block
																		i32.const 5
																		i32.eq
																		if(result i32)
																			i32.const 3
																		else
																			local.get $$current_block
																			i32.const 7
																			i32.eq
																			if(result i32)
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
																i32.const 0
																i32.const 14784
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
															end
															i32.const 2
															local.set $$current_block
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
															local.get $r
															i32.const 255
															i32.and
															local.set $$t8
															local.get $$t7.1
															local.get $$t8
															i32.store8 align=1
															i32.const 1
															local.set $$block_selector
															br $$BlockDisp
														end
														i32.const 3
														local.set $$current_block
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
														local.get $$t9.1
														i32.load8_u align=1
														local.set $$t10
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
														local.get $r
														i64.const 6
														i32.wrap_i64
														i32.shr_s
														local.set $$t12
														local.get $$t12
														i32.const 255
														i32.and
														local.set $$t13
														i32.const 192
														local.get $$t13
														i32.or
														local.set $$t14
														local.get $$t11.1
														local.get $$t14
														i32.store8 align=1
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
														local.get $r
														i32.const 255
														i32.and
														local.set $$t16
														local.get $$t16
														i32.const 63
														i32.and
														local.set $$t17
														i32.const 128
														local.get $$t17
														i32.or
														local.set $$t18
														local.get $$t15.1
														local.get $$t18
														i32.store8 align=1
														i32.const 1
														local.set $$block_selector
														br $$BlockDisp
													end
													i32.const 4
													local.set $$current_block
													local.get $$t2
													i32.const 2047
													i32.le_u
													local.set $$t19
													local.get $$t19
													if
														i32.const 3
														local.set $$block_selector
														br $$BlockDisp
													else
														br $$Block_5
													end
												end
												i32.const 5
												local.set $$current_block
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
												local.get $$t20.1
												i32.load8_u align=1
												local.set $$t21
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
												i32.const 65533
												i64.const 12
												i32.wrap_i64
												i32.shr_s
												local.set $$t23
												local.get $$t23
												i32.const 255
												i32.and
												local.set $$t24
												i32.const 224
												local.get $$t24
												i32.or
												local.set $$t25
												local.get $$t22.1
												local.get $$t25
												i32.store8 align=1
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
												i32.const 65533
												i64.const 6
												i32.wrap_i64
												i32.shr_s
												local.set $$t27
												local.get $$t27
												i32.const 255
												i32.and
												local.set $$t28
												local.get $$t28
												i32.const 63
												i32.and
												local.set $$t29
												i32.const 128
												local.get $$t29
												i32.or
												local.set $$t30
												local.get $$t26.1
												local.get $$t30
												i32.store8 align=1
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
												i32.const 65533
												i32.const 255
												i32.and
												local.set $$t32
												local.get $$t32
												i32.const 63
												i32.and
												local.set $$t33
												i32.const 128
												local.get $$t33
												i32.or
												local.set $$t34
												local.get $$t31.1
												local.get $$t34
												i32.store8 align=1
												i32.const 1
												local.set $$block_selector
												br $$BlockDisp
											end
											i32.const 6
											local.set $$current_block
											local.get $$t2
											i32.const 1114111
											i32.gt_u
											local.set $$t35
											local.get $$t35
											if
												i32.const 5
												local.set $$block_selector
												br $$BlockDisp
											else
												br $$Block_7
											end
										end
										i32.const 7
										local.set $$current_block
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
										local.get $$t36.1
										i32.load8_u align=1
										local.set $$t37
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
										local.get $r
										i64.const 12
										i32.wrap_i64
										i32.shr_s
										local.set $$t39
										local.get $$t39
										i32.const 255
										i32.and
										local.set $$t40
										i32.const 224
										local.get $$t40
										i32.or
										local.set $$t41
										local.get $$t38.1
										local.get $$t41
										i32.store8 align=1
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
										local.get $r
										i64.const 6
										i32.wrap_i64
										i32.shr_s
										local.set $$t43
										local.get $$t43
										i32.const 255
										i32.and
										local.set $$t44
										local.get $$t44
										i32.const 63
										i32.and
										local.set $$t45
										i32.const 128
										local.get $$t45
										i32.or
										local.set $$t46
										local.get $$t42.1
										local.get $$t46
										i32.store8 align=1
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
										local.get $r
										i32.const 255
										i32.and
										local.set $$t48
										local.get $$t48
										i32.const 63
										i32.and
										local.set $$t49
										i32.const 128
										local.get $$t49
										i32.or
										local.set $$t50
										local.get $$t47.1
										local.get $$t50
										i32.store8 align=1
										i32.const 1
										local.set $$block_selector
										br $$BlockDisp
									end
									i32.const 8
									local.set $$current_block
									i32.const 55296
									local.get $$t2
									i32.le_u
									local.set $$t51
									local.get $$t51
									if
										br $$Block_9
									else
										br $$Block_10
									end
								end
								i32.const 9
								local.set $$current_block
								local.get $$t2
								i32.const 65535
								i32.le_u
								local.set $$t52
								local.get $$t52
								if
									i32.const 7
									local.set $$block_selector
									br $$BlockDisp
								else
									br $$Block_11
								end
							end
							i32.const 10
							local.set $$current_block
							local.get $$t2
							i32.const 57343
							i32.le_u
							local.set $$t53
							br $$Block_10
						end
						local.get $$current_block
						i32.const 8
						i32.eq
						if(result i32)
							i32.const 0
						else
							local.get $$t53
						end
						local.set $$t54
						i32.const 11
						local.set $$current_block
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
					end
					i32.const 12
					local.set $$current_block
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
					local.get $$t55.1
					i32.load8_u align=1
					local.set $$t56
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
					local.get $r
					i64.const 18
					i32.wrap_i64
					i32.shr_s
					local.set $$t58
					local.get $$t58
					i32.const 255
					i32.and
					local.set $$t59
					i32.const 240
					local.get $$t59
					i32.or
					local.set $$t60
					local.get $$t57.1
					local.get $$t60
					i32.store8 align=1
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
					local.get $r
					i64.const 12
					i32.wrap_i64
					i32.shr_s
					local.set $$t62
					local.get $$t62
					i32.const 255
					i32.and
					local.set $$t63
					local.get $$t63
					i32.const 63
					i32.and
					local.set $$t64
					i32.const 128
					local.get $$t64
					i32.or
					local.set $$t65
					local.get $$t61.1
					local.get $$t65
					i32.store8 align=1
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
					local.get $r
					i64.const 6
					i32.wrap_i64
					i32.shr_s
					local.set $$t67
					local.get $$t67
					i32.const 255
					i32.and
					local.set $$t68
					local.get $$t68
					i32.const 63
					i32.and
					local.set $$t69
					i32.const 128
					local.get $$t69
					i32.or
					local.set $$t70
					local.get $$t66.1
					local.get $$t70
					i32.store8 align=1
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
					local.get $r
					i32.const 255
					i32.and
					local.set $$t72
					local.get $$t72
					i32.const 63
					i32.and
					local.set $$t73
					i32.const 128
					local.get $$t73
					i32.or
					local.set $$t74
					local.get $$t71.1
					local.get $$t74
					i32.store8 align=1
					i32.const 1
					local.set $$block_selector
					br $$BlockDisp
				end
			end
		end
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
	(func $$runtime.waPrintRune (param $ch i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$t0.0 i32)
		(local $$t0.1 i32)
		(local $$t0.2 i32)
		block $$BlockFnBody
			loop $$BlockDisp
				block $$Block_0
					block $$BlockSel
						local.get $$block_selector
						br_table  0 0
					end
					i32.const 0
					local.set $$current_block
					local.get $ch
					call $runtime.stringFromRune
					local.set $$t0.2
					local.set $$t0.1
					local.get $$t0.0
					call $runtime.Block.Release
					local.set $$t0.0
					local.get $$t0.0
					local.get $$t0.1
					local.get $$t0.2
					call $runtime.printString
					br $$BlockFnBody
				end
			end
		end
		local.get $$t0.0
		call $runtime.Block.Release
	)
	(func $$runtime.waPrintString (param $s.0 i32) (param $s.1 i32) (param $s.2 i32)
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
					local.get $s.0
					local.get $s.1
					local.get $s.2
					call $runtime.printString
					br $$BlockFnBody
				end
			end
		end
	)
	(func $$runtime.waPuts (param $ptr i32) (param $len i32)
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
					local.get $ptr
					local.get $len
					call $runtime.traceUtf8
					br $$BlockFnBody
				end
			end
		end
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
	(func $$math$rand.Source.underlying.$$OnFree (param $$ptr i32)
		local.get $$ptr
		i32.const 10
		call_indirect 0 (type $$OnFree)
	)
	(func $$math$rand.Source64.underlying.$$OnFree (param $$ptr i32)
		local.get $$ptr
		i32.const 10
		call_indirect 0 (type $$OnFree)
	)
	(func $$math$rand.Rand.$$OnFree (param $$ptr i32)
		local.get $$ptr
		i32.const 23
		call_indirect 0 (type $$OnFree)
		local.get $$ptr
		i32.const 16
		i32.add
		i32.const 24
		call_indirect 0 (type $$OnFree)
	)
	(func $math$rand.New (param $src.0.0 i32) (param $src.0.1 i32) (param $src.1 i32) (param $src.2 i32) (result i32 i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$ret_0.0 i32)
		(local $$ret_0.1 i32)
		(local $$t0.0.0.0 i32)
		(local $$t0.0.0.1 i32)
		(local $$t0.0.1 i32)
		(local $$t0.0.2 i32)
		(local $$t0.1 i32)
		(local $$t1.0.0 i32)
		(local $$t1.0.1 i32)
		(local $$t1.1 i32)
		(local $$t1.2 i32)
		(local $$t2 i32)
		(local $$t3.0 i32)
		(local $$t3.1 i32)
		(local $$t4.0 i32)
		(local $$t4.1 i32)
		(local $$t5.0 i32)
		(local $$t5.1 i32)
		block $$BlockFnBody
			loop $$BlockDisp
				block $$Block_0
					block $$BlockSel
						local.get $$block_selector
						br_table  0 0
					end
					i32.const 0
					local.set $$current_block
					local.get $src.0.0
					local.get $src.0.1
					local.get $src.1
					local.get $src.2
					i32.const -3
					call $$wa.runtime.queryIface_CommaOk
					local.set $$t0.1
					local.set $$t0.0.2
					local.set $$t0.0.1
					local.set $$t0.0.0.1
					local.get $$t0.0.0.0
					call $runtime.Block.Release
					local.set $$t0.0.0.0
					local.get $$t0.0.0.0
					call $runtime.Block.Retain
					local.get $$t0.0.0.1
					local.get $$t0.0.1
					local.get $$t0.0.2
					local.set $$t1.2
					local.set $$t1.1
					local.set $$t1.0.1
					local.get $$t1.0.0
					call $runtime.Block.Release
					local.set $$t1.0.0
					local.get $$t0.1
					local.set $$t2
					i32.const 64
					call $runtime.HeapAlloc
					i32.const 1
					i32.const 25
					i32.const 48
					call $runtime.Block.Init
					call $runtime.DupI32
					i32.const 16
					i32.add
					local.set $$t3.1
					local.get $$t3.0
					call $runtime.Block.Release
					local.set $$t3.0
					local.get $$t3.0
					call $runtime.Block.Retain
					local.get $$t3.1
					i32.const 0
					i32.add
					local.set $$t4.1
					local.get $$t4.0
					call $runtime.Block.Release
					local.set $$t4.0
					local.get $$t3.0
					call $runtime.Block.Retain
					local.get $$t3.1
					i32.const 16
					i32.add
					local.set $$t5.1
					local.get $$t5.0
					call $runtime.Block.Release
					local.set $$t5.0
					local.get $$t4.1
					local.get $src.0.0
					call $runtime.Block.Retain
					local.get $$t4.1
					i32.load align=1
					call $runtime.Block.Release
					i32.store align=1
					local.get $$t4.1
					local.get $src.0.1
					i32.store offset=4
					local.get $$t4.1
					local.get $src.1
					i32.store offset=8
					local.get $$t4.1
					local.get $src.2
					i32.store offset=12
					local.get $$t5.1
					local.get $$t1.0.0
					call $runtime.Block.Retain
					local.get $$t5.1
					i32.load align=1
					call $runtime.Block.Release
					i32.store align=1
					local.get $$t5.1
					local.get $$t1.0.1
					i32.store offset=4
					local.get $$t5.1
					local.get $$t1.1
					i32.store offset=8
					local.get $$t5.1
					local.get $$t1.2
					i32.store offset=12
					local.get $$t3.0
					call $runtime.Block.Retain
					local.get $$t3.1
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
		local.get $$t0.0.0.0
		call $runtime.Block.Release
		local.get $$t1.0.0
		call $runtime.Block.Release
		local.get $$t3.0
		call $runtime.Block.Release
		local.get $$t4.0
		call $runtime.Block.Release
		local.get $$t5.0
		call $runtime.Block.Release
	)
	(func $math$rand.NewSource (param $seed i64) (result i32 i32 i32 i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$ret_0.0.0 i32)
		(local $$ret_0.0.1 i32)
		(local $$ret_0.1 i32)
		(local $$ret_0.2 i32)
		(local $$t0.0 i32)
		(local $$t0.1 i32)
		(local $$t1.0.0 i32)
		(local $$t1.0.1 i32)
		(local $$t1.1 i32)
		(local $$t1.2 i32)
		block $$BlockFnBody
			loop $$BlockDisp
				block $$Block_0
					block $$BlockSel
						local.get $$block_selector
						br_table  0 0
					end
					i32.const 0
					local.set $$current_block
					local.get $seed
					call $math$rand.newSource
					local.set $$t0.1
					local.get $$t0.0
					call $runtime.Block.Release
					local.set $$t0.0
					local.get $$t0.0
					call $runtime.Block.Retain
					local.get $$t0.1
					i32.const 2
					i32.const -2
					i32.const 0
					call $runtime.getItab
					i32.const 0
					local.set $$t1.2
					local.set $$t1.1
					local.set $$t1.0.1
					local.get $$t1.0.0
					call $runtime.Block.Release
					local.set $$t1.0.0
					local.get $$t1.0.0
					call $runtime.Block.Retain
					local.get $$t1.0.1
					local.get $$t1.1
					local.get $$t1.2
					local.set $$ret_0.2
					local.set $$ret_0.1
					local.set $$ret_0.0.1
					local.get $$ret_0.0.0
					call $runtime.Block.Release
					local.set $$ret_0.0.0
					br $$BlockFnBody
				end
			end
		end
		local.get $$ret_0.0.0
		call $runtime.Block.Retain
		local.get $$ret_0.0.1
		local.get $$ret_0.1
		local.get $$ret_0.2
		local.get $$ret_0.0.0
		call $runtime.Block.Release
		local.get $$t0.0
		call $runtime.Block.Release
		local.get $$t1.0.0
		call $runtime.Block.Release
	)
	(func $math$rand.init
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$t0 i32)
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
		(local $$t6.0.0 i32)
		(local $$t6.0.1 i32)
		(local $$t6.1 i32)
		(local $$t6.2 i32)
		(local $$t7.0.0 i32)
		(local $$t7.0.1 i32)
		(local $$t7.1 i32)
		(local $$t7.2 i32)
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
							global.get $math$rand.init$guard
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
						global.set $math$rand.init$guard
						call $math.init
						call $errors.init
						i32.const 64
						call $runtime.HeapAlloc
						i32.const 1
						i32.const 25
						i32.const 48
						call $runtime.Block.Init
						call $runtime.DupI32
						i32.const 16
						i32.add
						local.set $$t1.1
						local.get $$t1.0
						call $runtime.Block.Release
						local.set $$t1.0
						local.get $$t1.0
						call $runtime.Block.Retain
						local.get $$t1.1
						i32.const 0
						i32.add
						local.set $$t2.1
						local.get $$t2.0
						call $runtime.Block.Release
						local.set $$t2.0
						i32.const 4880
						call $runtime.HeapAlloc
						i32.const 1
						i32.const 0
						i32.const 4864
						call $runtime.Block.Init
						call $runtime.DupI32
						i32.const 16
						i32.add
						local.set $$t3.1
						local.get $$t3.0
						call $runtime.Block.Release
						local.set $$t3.0
						local.get $$t1.0
						call $runtime.Block.Retain
						local.get $$t1.1
						i32.const 16
						i32.add
						local.set $$t4.1
						local.get $$t4.0
						call $runtime.Block.Release
						local.set $$t4.0
						i32.const 4880
						call $runtime.HeapAlloc
						i32.const 1
						i32.const 0
						i32.const 4864
						call $runtime.Block.Init
						call $runtime.DupI32
						i32.const 16
						i32.add
						local.set $$t5.1
						local.get $$t5.0
						call $runtime.Block.Release
						local.set $$t5.0
						local.get $$t3.0
						call $runtime.Block.Retain
						local.get $$t3.1
						i32.const 2
						i32.const -2
						i32.const 0
						call $runtime.getItab
						i32.const 0
						local.set $$t6.2
						local.set $$t6.1
						local.set $$t6.0.1
						local.get $$t6.0.0
						call $runtime.Block.Release
						local.set $$t6.0.0
						local.get $$t2.1
						local.get $$t6.0.0
						call $runtime.Block.Retain
						local.get $$t2.1
						i32.load align=1
						call $runtime.Block.Release
						i32.store align=1
						local.get $$t2.1
						local.get $$t6.0.1
						i32.store offset=4
						local.get $$t2.1
						local.get $$t6.1
						i32.store offset=8
						local.get $$t2.1
						local.get $$t6.2
						i32.store offset=12
						local.get $$t5.0
						call $runtime.Block.Retain
						local.get $$t5.1
						i32.const 2
						i32.const -3
						i32.const 0
						call $runtime.getItab
						i32.const 0
						local.set $$t7.2
						local.set $$t7.1
						local.set $$t7.0.1
						local.get $$t7.0.0
						call $runtime.Block.Release
						local.set $$t7.0.0
						local.get $$t4.1
						local.get $$t7.0.0
						call $runtime.Block.Retain
						local.get $$t4.1
						i32.load align=1
						call $runtime.Block.Release
						i32.store align=1
						local.get $$t4.1
						local.get $$t7.0.1
						i32.store offset=4
						local.get $$t4.1
						local.get $$t7.1
						i32.store offset=8
						local.get $$t4.1
						local.get $$t7.2
						i32.store offset=12
						i32.const 16440
						local.get $$t1.0
						call $runtime.Block.Retain
						i32.const 16440
						i32.load align=1
						call $runtime.Block.Release
						i32.store align=1
						i32.const 16440
						local.get $$t1.1
						i32.store offset=4
						br $$Block_1
					end
					i32.const 2
					local.set $$current_block
					br $$BlockFnBody
				end
			end
		end
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
		local.get $$t6.0.0
		call $runtime.Block.Release
		local.get $$t7.0.0
		call $runtime.Block.Release
	)
	(func $math$rand.newSource (param $seed i64) (result i32 i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$ret_0.0 i32)
		(local $$ret_0.1 i32)
		(local $$t0.0 i32)
		(local $$t0.1 i32)
		block $$BlockFnBody
			loop $$BlockDisp
				block $$Block_0
					block $$BlockSel
						local.get $$block_selector
						br_table  0 0
					end
					i32.const 0
					local.set $$current_block
					i32.const 4880
					call $runtime.HeapAlloc
					i32.const 1
					i32.const 0
					i32.const 4864
					call $runtime.Block.Init
					call $runtime.DupI32
					i32.const 16
					i32.add
					local.set $$t0.1
					local.get $$t0.0
					call $runtime.Block.Release
					local.set $$t0.0
					local.get $$t0.0
					local.get $$t0.1
					local.get $seed
					call $math$rand.rngSource.Seed
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
	)
	(func $math$rand.seedrand (param $x i32) (result i32)
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
							local.get $x
							i32.const 44488
							i32.div_s
							local.set $$t0
							local.get $x
							i32.const 44488
							i32.rem_s
							local.set $$t1
							i32.const 48271
							local.get $$t1
							i32.mul
							local.set $$t2
							i32.const 3399
							local.get $$t0
							i32.mul
							local.set $$t3
							local.get $$t2
							local.get $$t3
							i32.sub
							local.set $$t4
							local.get $$t4
							i32.const 0
							i32.lt_s
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
						local.get $$t4
						i32.const 2147483647
						i32.add
						local.set $$t6
						br $$Block_1
					end
					local.get $$current_block
					i32.const 0
					i32.eq
					if(result i32)
						local.get $$t4
					else
						local.get $$t6
					end
					local.set $$t7
					i32.const 2
					local.set $$current_block
					local.get $$t7
					local.set $$ret_0
					br $$BlockFnBody
				end
			end
		end
		local.get $$ret_0
	)
	(func $$errors.errorString.$$OnFree (param $$ptr i32)
		local.get $$ptr
		i32.const 3
		call_indirect 0 (type $$OnFree)
	)
	(func $errors.init
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
							global.get $errors.init$guard
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
						global.set $errors.init$guard
						br $$Block_1
					end
					i32.const 2
					local.set $$current_block
					br $$BlockFnBody
				end
			end
		end
	)
	(func $syscall$wasm4.BlitI32 (param $sprite.0 i32) (param $sprite.1 i32) (param $sprite.2 i32) (param $sprite.3 i32) (param $x i32) (param $y i32) (param $width i32) (param $height i32) (param $flags i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$t0 i32)
		block $$BlockFnBody
			loop $$BlockDisp
				block $$Block_0
					block $$BlockSel
						local.get $$block_selector
						br_table  0 0
					end
					i32.const 0
					local.set $$current_block
					local.get $sprite.0
					local.get $sprite.1
					local.get $sprite.2
					local.get $sprite.3
					call $$syscall/wasm4.__linkname__slice_data_ptr
					local.set $$t0
					local.get $$t0
					local.get $x
					local.get $y
					local.get $width
					local.get $height
					local.get $flags
					call $syscall$wasm4.__import__blit
					br $$BlockFnBody
				end
			end
		end
	)
	(func $syscall$wasm4.GetGamePad1 (result i32)
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
					i32.const 22
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
					i32.const 1
					i32.const 0
					i32.mul
					i32.add
					local.set $$t1.1
					local.get $$t1.0
					call $runtime.Block.Release
					local.set $$t1.0
					local.get $$t1.1
					i32.load8_u align=1
					local.set $$t2
					local.get $$t2
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
	(func $syscall$wasm4.RectI32 (param $x i32) (param $y i32) (param $width i32) (param $height i32)
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
					local.get $x
					local.get $y
					local.get $width
					local.get $height
					call $syscall$wasm4.__import__rect
					br $$BlockFnBody
				end
			end
		end
	)
	(func $syscall$wasm4.SetDrawColors (param $a i32) (param $b i32) (param $c i32) (param $d i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$t0.0 i32)
		(local $$t0.1 i32)
		(local $$t0.2 i32)
		(local $$t0.3 i32)
		(local $$t1 i32)
		(local $$t2 i32)
		(local $$t3 i32)
		(local $$t4 i32)
		(local $$t5 i32)
		(local $$t6 i32)
		(local $$t7 i32)
		(local $$t8 i32)
		(local $$t9 i32)
		(local $$t10 i32)
		(local $$t11.0 i32)
		(local $$t11.1 i32)
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
					i32.const 0
					i32.const 20
					i32.const 1
					i32.const 1
					call $$syscall/wasm4.__linkname__make_slice
					local.set $$t0.3
					local.set $$t0.2
					local.set $$t0.1
					local.get $$t0.0
					call $runtime.Block.Release
					local.set $$t0.0
					local.get $a
					local.set $$t1
					local.get $b
					local.set $$t2
					local.get $$t2
					i64.const 4
					i32.wrap_i64
					i32.shl
					local.set $$t3
					local.get $$t1
					local.get $$t3
					i32.add
					local.set $$t4
					local.get $c
					local.set $$t5
					local.get $$t5
					i64.const 8
					i32.wrap_i64
					i32.shl
					local.set $$t6
					local.get $$t4
					local.get $$t6
					i32.add
					local.set $$t7
					local.get $d
					local.set $$t8
					local.get $$t8
					i64.const 12
					i32.wrap_i64
					i32.shl
					local.set $$t9
					local.get $$t7
					local.get $$t9
					i32.add
					local.set $$t10
					local.get $$t0.0
					call $runtime.Block.Retain
					local.get $$t0.1
					i32.const 2
					i32.const 0
					i32.mul
					i32.add
					local.set $$t11.1
					local.get $$t11.0
					call $runtime.Block.Release
					local.set $$t11.0
					local.get $$t10
					i32.const 65535
					i32.and
					local.set $$t12
					local.get $$t11.1
					local.get $$t12
					i32.store16
					br $$BlockFnBody
				end
			end
		end
		local.get $$t0.0
		call $runtime.Block.Release
		local.get $$t11.0
		call $runtime.Block.Release
	)
	(func $syscall$wasm4.SetPalette (param $a0 i32) (param $a1 i32) (param $a2 i32) (param $a3 i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$t0.0 i32)
		(local $$t0.1 i32)
		(local $$t0.2 i32)
		(local $$t0.3 i32)
		(local $$t1.0 i32)
		(local $$t1.1 i32)
		(local $$t2.0 i32)
		(local $$t2.1 i32)
		(local $$t3.0 i32)
		(local $$t3.1 i32)
		(local $$t4.0 i32)
		(local $$t4.1 i32)
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
					i32.const 4
					i32.const 4
					i32.const 4
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
					i32.const 4
					i32.const 0
					i32.mul
					i32.add
					local.set $$t1.1
					local.get $$t1.0
					call $runtime.Block.Release
					local.set $$t1.0
					local.get $$t0.0
					call $runtime.Block.Retain
					local.get $$t0.1
					i32.const 4
					i32.const 1
					i32.mul
					i32.add
					local.set $$t2.1
					local.get $$t2.0
					call $runtime.Block.Release
					local.set $$t2.0
					local.get $$t0.0
					call $runtime.Block.Retain
					local.get $$t0.1
					i32.const 4
					i32.const 2
					i32.mul
					i32.add
					local.set $$t3.1
					local.get $$t3.0
					call $runtime.Block.Release
					local.set $$t3.0
					local.get $$t0.0
					call $runtime.Block.Retain
					local.get $$t0.1
					i32.const 4
					i32.const 3
					i32.mul
					i32.add
					local.set $$t4.1
					local.get $$t4.0
					call $runtime.Block.Release
					local.set $$t4.0
					local.get $$t1.1
					local.get $a0
					i32.store
					local.get $$t2.1
					local.get $a1
					i32.store
					local.get $$t3.1
					local.get $a2
					i32.store
					local.get $$t4.1
					local.get $a3
					i32.store
					br $$BlockFnBody
				end
			end
		end
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
	)
	(func $syscall$wasm4.Tone (param $frequency i32) (param $duration i32) (param $volume i32) (param $flags i32)
		(local $$block_selector i32)
		(local $$current_block i32)
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
					local.get $frequency
					local.set $$t0
					local.get $duration
					local.set $$t1
					local.get $volume
					local.set $$t2
					local.get $flags
					local.set $$t3
					local.get $$t0
					local.get $$t1
					local.get $$t2
					local.get $$t3
					call $syscall$wasm4.__import__tone
					br $$BlockFnBody
				end
			end
		end
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
	(func $math$rand.Rand.Intn.$bound (param $n i32) (result i32)
		global.get $$wa.runtime.closure_data
		i32.load
		global.get $$wa.runtime.closure_data
		i32.load offset=4
		i32.const 0
		global.set $$wa.runtime.closure_data
		local.get $n
		call $math$rand.Rand.Intn
	)
	(func $$math$rand.Rand.$$block.$$OnFree (param $ptr i32)
		local.get $ptr
		i32.load align=1
		call $runtime.Block.Release
		local.get $ptr
		i32.const 0
		i32.store align=1
	)
	(func $$math$rand.Rand.$ref.underlying.$$OnFree (param $$ptr i32)
		local.get $$ptr
		i32.const 28
		call_indirect 0 (type $$OnFree)
	)
	(func $w4snake.init
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$t0 i32)
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
		(local $$t16.0 i32)
		(local $$t16.1 i32)
		(local $$t17.0 i32)
		(local $$t17.1 i32)
		(local $$t18.0 i32)
		(local $$t18.1 i32)
		(local $$t18.2 i32)
		(local $$t18.3 i32)
		(local $$t19.0.0 i32)
		(local $$t19.0.1 i32)
		(local $$t19.1 i32)
		(local $$t19.2 i32)
		(local $$t20.0 i32)
		(local $$t20.1 i32)
		(local $$t21.0 i32)
		(local $$t21.1.0 i32)
		(local $$t21.1.1 i32)
		(local $$t22.0 i32)
		(local $$t22.1.0 i32)
		(local $$t22.1.1 i32)
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
							global.get $w4snake.init$guard
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
						global.set $w4snake.init$guard
						call $runtime.init
						call $math$rand.init
						call $syscall$wasm4.init
						i32.const 32
						call $runtime.HeapAlloc
						i32.const 1
						i32.const 0
						i32.const 16
						call $runtime.Block.Init
						call $runtime.DupI32
						i32.const 16
						i32.add
						local.set $$t1.1
						local.get $$t1.0
						call $runtime.Block.Release
						local.set $$t1.0
						local.get $$t1.0
						call $runtime.Block.Retain
						local.get $$t1.1
						i32.const 1
						i32.const 0
						i32.mul
						i32.add
						local.set $$t2.1
						local.get $$t2.0
						call $runtime.Block.Release
						local.set $$t2.0
						local.get $$t2.1
						i32.const 0
						i32.store8 align=1
						local.get $$t1.0
						call $runtime.Block.Retain
						local.get $$t1.1
						i32.const 1
						i32.const 1
						i32.mul
						i32.add
						local.set $$t3.1
						local.get $$t3.0
						call $runtime.Block.Release
						local.set $$t3.0
						local.get $$t3.1
						i32.const 160
						i32.store8 align=1
						local.get $$t1.0
						call $runtime.Block.Retain
						local.get $$t1.1
						i32.const 1
						i32.const 2
						i32.mul
						i32.add
						local.set $$t4.1
						local.get $$t4.0
						call $runtime.Block.Release
						local.set $$t4.0
						local.get $$t4.1
						i32.const 2
						i32.store8 align=1
						local.get $$t1.0
						call $runtime.Block.Retain
						local.get $$t1.1
						i32.const 1
						i32.const 3
						i32.mul
						i32.add
						local.set $$t5.1
						local.get $$t5.0
						call $runtime.Block.Release
						local.set $$t5.0
						local.get $$t5.1
						i32.const 0
						i32.store8 align=1
						local.get $$t1.0
						call $runtime.Block.Retain
						local.get $$t1.1
						i32.const 1
						i32.const 4
						i32.mul
						i32.add
						local.set $$t6.1
						local.get $$t6.0
						call $runtime.Block.Release
						local.set $$t6.0
						local.get $$t6.1
						i32.const 14
						i32.store8 align=1
						local.get $$t1.0
						call $runtime.Block.Retain
						local.get $$t1.1
						i32.const 1
						i32.const 5
						i32.mul
						i32.add
						local.set $$t7.1
						local.get $$t7.0
						call $runtime.Block.Release
						local.set $$t7.0
						local.get $$t7.1
						i32.const 240
						i32.store8 align=1
						local.get $$t1.0
						call $runtime.Block.Retain
						local.get $$t1.1
						i32.const 1
						i32.const 6
						i32.mul
						i32.add
						local.set $$t8.1
						local.get $$t8.0
						call $runtime.Block.Release
						local.set $$t8.0
						local.get $$t8.1
						i32.const 54
						i32.store8 align=1
						local.get $$t1.0
						call $runtime.Block.Retain
						local.get $$t1.1
						i32.const 1
						i32.const 7
						i32.mul
						i32.add
						local.set $$t9.1
						local.get $$t9.0
						call $runtime.Block.Release
						local.set $$t9.0
						local.get $$t9.1
						i32.const 92
						i32.store8 align=1
						local.get $$t1.0
						call $runtime.Block.Retain
						local.get $$t1.1
						i32.const 1
						i32.const 8
						i32.mul
						i32.add
						local.set $$t10.1
						local.get $$t10.0
						call $runtime.Block.Release
						local.set $$t10.0
						local.get $$t10.1
						i32.const 214
						i32.store8 align=1
						local.get $$t1.0
						call $runtime.Block.Retain
						local.get $$t1.1
						i32.const 1
						i32.const 9
						i32.mul
						i32.add
						local.set $$t11.1
						local.get $$t11.0
						call $runtime.Block.Release
						local.set $$t11.0
						local.get $$t11.1
						i32.const 87
						i32.store8 align=1
						local.get $$t1.0
						call $runtime.Block.Retain
						local.get $$t1.1
						i32.const 1
						i32.const 10
						i32.mul
						i32.add
						local.set $$t12.1
						local.get $$t12.0
						call $runtime.Block.Release
						local.set $$t12.0
						local.get $$t12.1
						i32.const 213
						i32.store8 align=1
						local.get $$t1.0
						call $runtime.Block.Retain
						local.get $$t1.1
						i32.const 1
						i32.const 11
						i32.mul
						i32.add
						local.set $$t13.1
						local.get $$t13.0
						call $runtime.Block.Release
						local.set $$t13.0
						local.get $$t13.1
						i32.const 87
						i32.store8 align=1
						local.get $$t1.0
						call $runtime.Block.Retain
						local.get $$t1.1
						i32.const 1
						i32.const 12
						i32.mul
						i32.add
						local.set $$t14.1
						local.get $$t14.0
						call $runtime.Block.Release
						local.set $$t14.0
						local.get $$t14.1
						i32.const 53
						i32.store8 align=1
						local.get $$t1.0
						call $runtime.Block.Retain
						local.get $$t1.1
						i32.const 1
						i32.const 13
						i32.mul
						i32.add
						local.set $$t15.1
						local.get $$t15.0
						call $runtime.Block.Release
						local.set $$t15.0
						local.get $$t15.1
						i32.const 92
						i32.store8 align=1
						local.get $$t1.0
						call $runtime.Block.Retain
						local.get $$t1.1
						i32.const 1
						i32.const 14
						i32.mul
						i32.add
						local.set $$t16.1
						local.get $$t16.0
						call $runtime.Block.Release
						local.set $$t16.0
						local.get $$t16.1
						i32.const 15
						i32.store8 align=1
						local.get $$t1.0
						call $runtime.Block.Retain
						local.get $$t1.1
						i32.const 1
						i32.const 15
						i32.mul
						i32.add
						local.set $$t17.1
						local.get $$t17.0
						call $runtime.Block.Release
						local.set $$t17.0
						local.get $$t17.1
						i32.const 240
						i32.store8 align=1
						local.get $$t1.0
						call $runtime.Block.Retain
						local.get $$t1.1
						i32.const 1
						i32.const 0
						i32.mul
						i32.add
						i32.const 16
						i32.const 0
						i32.sub
						i32.const 16
						i32.const 0
						i32.sub
						local.set $$t18.3
						local.set $$t18.2
						local.set $$t18.1
						local.get $$t18.0
						call $runtime.Block.Release
						local.set $$t18.0
						i32.const 24384
						local.get $$t18.0
						call $runtime.Block.Retain
						i32.const 24384
						i32.load align=1
						call $runtime.Block.Release
						i32.store align=1
						i32.const 24384
						local.get $$t18.1
						i32.store offset=4
						i32.const 24384
						local.get $$t18.2
						i32.store offset=8
						i32.const 24384
						local.get $$t18.3
						i32.store offset=12
						i64.const 1
						call $math$rand.NewSource
						local.set $$t19.2
						local.set $$t19.1
						local.set $$t19.0.1
						local.get $$t19.0.0
						call $runtime.Block.Release
						local.set $$t19.0.0
						local.get $$t19.0.0
						local.get $$t19.0.1
						local.get $$t19.1
						local.get $$t19.2
						call $math$rand.New
						local.set $$t20.1
						local.get $$t20.0
						call $runtime.Block.Release
						local.set $$t20.0
						i32.const 27
						local.set $$t21.0
						i32.const 24
						call $runtime.HeapAlloc
						i32.const 1
						i32.const 29
						i32.const 8
						call $runtime.Block.Init
						call $runtime.DupI32
						i32.const 16
						i32.add
						local.set $$t21.1.1
						local.get $$t21.1.0
						call $runtime.Block.Release
						local.set $$t21.1.0
						local.get $$t21.1.1
						local.get $$t20.0
						call $runtime.Block.Retain
						local.get $$t21.1.1
						i32.load align=1
						call $runtime.Block.Release
						i32.store align=1
						local.get $$t21.1.1
						local.get $$t20.1
						i32.store offset=4
						local.get $$t21.0
						local.get $$t21.1.0
						call $runtime.Block.Retain
						local.get $$t21.1.1
						local.set $$t22.1.1
						local.get $$t22.1.0
						call $runtime.Block.Release
						local.set $$t22.1.0
						local.set $$t22.0
						i32.const 24400
						local.get $$t22.0
						i32.store
						i32.const 24400
						local.get $$t22.1.0
						call $runtime.Block.Retain
						i32.const 24400
						i32.load offset=4 align=1
						call $runtime.Block.Release
						i32.store offset=4 align=1
						i32.const 24400
						local.get $$t22.1.1
						i32.store offset=8
						br $$Block_1
					end
					i32.const 2
					local.set $$current_block
					br $$BlockFnBody
				end
			end
		end
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
		local.get $$t16.0
		call $runtime.Block.Release
		local.get $$t17.0
		call $runtime.Block.Release
		local.get $$t18.0
		call $runtime.Block.Release
		local.get $$t19.0.0
		call $runtime.Block.Release
		local.get $$t20.0
		call $runtime.Block.Release
		local.get $$t21.1.0
		call $runtime.Block.Release
		local.get $$t22.1.0
		call $runtime.Block.Release
	)
	(func $w4snake.input
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
		(local $$t10 i32)
		(local $$t11 i32)
		(local $$t12 i32)
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
																	call $syscall$wasm4.GetGamePad1
																	local.set $$t0
																	local.get $$t0
																	i32.const 1
																	i32.and
																	local.set $$t1
																	local.get $$t1
																	i32.const 0
																	i32.eq
																	i32.eqz
																	local.set $$t2
																	local.get $$t2
																	if
																		br $$Block_0
																	else
																		br $$Block_1
																	end
																end
																i32.const 1
																local.set $$current_block
																i32.const 24472
																i32.const 1
																i32.store8 align=1
																br $$Block_1
															end
															i32.const 2
															local.set $$current_block
															local.get $$t0
															i32.const 2
															i32.and
															local.set $$t3
															local.get $$t3
															i32.const 0
															i32.eq
															i32.eqz
															local.set $$t4
															local.get $$t4
															if
																br $$Block_2
															else
																br $$Block_3
															end
														end
														i32.const 3
														local.set $$current_block
														i32.const 24472
														i32.const 0
														i32.store8 align=1
														br $$Block_3
													end
													i32.const 4
													local.set $$current_block
													local.get $$t0
													i32.const 64
													i32.and
													local.set $$t5
													local.get $$t5
													i32.const 0
													i32.eq
													i32.eqz
													local.set $$t6
													local.get $$t6
													if
														br $$Block_4
													else
														br $$Block_5
													end
												end
												i32.const 5
												local.set $$current_block
												i32.const 0
												i32.const 24448
												call $w4snake.Snake.Up
												br $$Block_5
											end
											i32.const 6
											local.set $$current_block
											local.get $$t0
											i32.const 128
											i32.and
											local.set $$t7
											local.get $$t7
											i32.const 0
											i32.eq
											i32.eqz
											local.set $$t8
											local.get $$t8
											if
												br $$Block_6
											else
												br $$Block_7
											end
										end
										i32.const 7
										local.set $$current_block
										i32.const 0
										i32.const 24448
										call $w4snake.Snake.Down
										br $$Block_7
									end
									i32.const 8
									local.set $$current_block
									local.get $$t0
									i32.const 16
									i32.and
									local.set $$t9
									local.get $$t9
									i32.const 0
									i32.eq
									i32.eqz
									local.set $$t10
									local.get $$t10
									if
										br $$Block_8
									else
										br $$Block_9
									end
								end
								i32.const 9
								local.set $$current_block
								i32.const 0
								i32.const 24448
								call $w4snake.Snake.Left
								br $$Block_9
							end
							i32.const 10
							local.set $$current_block
							local.get $$t0
							i32.const 32
							i32.and
							local.set $$t11
							local.get $$t11
							i32.const 0
							i32.eq
							i32.eqz
							local.set $$t12
							local.get $$t12
							if
								br $$Block_10
							else
								br $$Block_11
							end
						end
						i32.const 11
						local.set $$current_block
						i32.const 0
						i32.const 24448
						call $w4snake.Snake.Right
						br $$Block_11
					end
					i32.const 12
					local.set $$current_block
					br $$BlockFnBody
				end
			end
		end
	)
	(func $w4snake.start (export "start")
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
					i32.const 16513011
					i32.const 15052931
					i32.const 4353629
					i32.const 2107453
					call $syscall$wasm4.SetPalette
					i32.const 0
					i32.const 24448
					call $w4snake.Snake.Reset
					br $$BlockFnBody
				end
			end
		end
	)
	(func $$w4snake.Point.$slice.append (param $x.0 i32) (param $x.1 i32) (param $x.2 i32) (param $x.3 i32) (param $y.0 i32) (param $y.1 i32) (param $y.2 i32) (param $y.3 i32) (result i32 i32 i32 i32)
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
		if(result i32 i32 i32 i32)
			local.get $x.0
			call $runtime.Block.Retain
			local.get $x.1
			local.get $new_len
			local.get $x.3
			local.get $y.1
			local.set $src
			local.get $x.1
			i32.const 2
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
					end
					local.get $src
					i32.load8_u align=1
					local.get $src
					i32.load8_u offset=1 align=1
					local.set $item.1
					local.set $item.0
					local.get $dest
					local.get $item.0
					i32.store8 align=1
					local.get $dest
					local.get $item.1
					i32.store8 offset=1 align=1
					local.get $src
					i32.const 2
					i32.add
					local.set $src
					local.get $dest
					i32.const 2
					i32.add
					local.set $dest
					local.get $y_len
					i32.const 1
					i32.sub
					local.set $y_len
					br $loop1
				end
			end
		else
			local.get $new_len
			i32.const 2
			i32.mul
			local.set $new_cap
			local.get $new_cap
			i32.const 2
			i32.mul
			i32.const 16
			i32.add
			call $runtime.HeapAlloc
			local.get $new_cap
			i32.const 0
			i32.const 2
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
					end
					local.get $src
					i32.load8_u align=1
					local.get $src
					i32.load8_u offset=1 align=1
					local.set $item.1
					local.set $item.0
					local.get $dest
					local.get $item.0
					i32.store8 align=1
					local.get $dest
					local.get $item.1
					i32.store8 offset=1 align=1
					local.get $src
					i32.const 2
					i32.add
					local.set $src
					local.get $dest
					i32.const 2
					i32.add
					local.set $dest
					local.get $x_len
					i32.const 1
					i32.sub
					local.set $x_len
					br $loop2
				end
			end
			local.get $y.1
			local.set $src
			block $block3
				loop $loop3
					local.get $y_len
					i32.eqz
					if
						br $block3
					end
					local.get $src
					i32.load8_u align=1
					local.get $src
					i32.load8_u offset=1 align=1
					local.set $item.1
					local.set $item.0
					local.get $dest
					local.get $item.0
					i32.store8 align=1
					local.get $dest
					local.get $item.1
					i32.store8 offset=1 align=1
					local.get $src
					i32.const 2
					i32.add
					local.set $src
					local.get $dest
					i32.const 2
					i32.add
					local.set $dest
					local.get $y_len
					i32.const 1
					i32.sub
					local.set $y_len
					br $loop3
				end
			end
		end
	)
	(func $w4snake.update (export "update")
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$t0 i32)
		(local $$t1 i32)
		(local $$t2 i32)
		(local $$t3 i32)
		(local $$t4 i32)
		(local $$t5 i32)
		(local $$t6.0 i32)
		(local $$t6.1 i32)
		(local $$t6.2 i32)
		(local $$t6.3 i32)
		(local $$t7 i32)
		(local $$t8 i32)
		(local $$t9 i32)
		(local $$t10 i32)
		(local $$t11.0 i32)
		(local $$t11.1 i32)
		(local $$t11.2 i32)
		(local $$t11.3 i32)
		(local $$t12.0 i32)
		(local $$t12.1 i32)
		(local $$t13.0 i32)
		(local $$t13.1 i32)
		(local $$t14.0 i32)
		(local $$t14.1 i32)
		(local $$t15 i32)
		(local $$t16.0 i32)
		(local $$t16.1 i32)
		(local $$t16.2 i32)
		(local $$t16.3 i32)
		(local $$t17.0 i32)
		(local $$t17.1 i32)
		(local $$t17.2 i32)
		(local $$t17.3 i32)
		(local $$t18.0 i32)
		(local $$t18.1 i32)
		(local $$t18.2 i32)
		(local $$t18.3 i32)
		(local $$t19 i32)
		(local $$t20 i32)
		(local $$t21.0 i32)
		(local $$t21.1 i32)
		(local $$t22.0 i32)
		(local $$t22.1 i32)
		(local $$t23.0 i32)
		(local $$t23.1 i32)
		(local $$t24.0 i32)
		(local $$t24.1 i32)
		(local $$t25.0 i32)
		(local $$t25.1 i32)
		(local $$t25.2 i32)
		(local $$t25.3 i32)
		(local $$t26.0 i32)
		(local $$t26.1 i32)
		(local $$t26.2 i32)
		(local $$t26.3 i32)
		(local $$t27.0 i32)
		(local $$t27.1.0 i32)
		(local $$t27.1.1 i32)
		(local $$t28 i32)
		(local $$t29 i32)
		(local $$t30 i32)
		(local $$t31.0 i32)
		(local $$t31.1.0 i32)
		(local $$t31.1.1 i32)
		(local $$t32 i32)
		(local $$t33 i32)
		(local $$t34 i32)
		block $$BlockFnBody
			loop $$BlockDisp
				block $$Block_5
					block $$Block_4
						block $$Block_3
							block $$Block_2
								block $$Block_1
									block $$Block_0
										block $$BlockSel
											local.get $$block_selector
											br_table  0 1 2 3 4 5 0
										end
										i32.const 0
										local.set $$current_block
										call $w4snake.input
										i32.const 24376
										i32.load
										local.set $$t0
										local.get $$t0
										i32.const 1
										i32.add
										local.set $$t1
										i32.const 24376
										local.get $$t1
										i32.store
										i32.const 24376
										i32.load
										local.set $$t2
										local.get $$t2
										i32.const 10
										i32.rem_s
										local.set $$t3
										local.get $$t3
										i32.const 0
										i32.eq
										local.set $$t4
										local.get $$t4
										if
											br $$Block_0
										else
											br $$Block_1
										end
									end
									i32.const 1
									local.set $$current_block
									i32.const 0
									i32.const 24448
									call $w4snake.Snake.Update
									i32.const 0
									i32.const 24448
									call $w4snake.Snake.IsDead
									local.set $$t5
									local.get $$t5
									if
										br $$Block_2
									else
										br $$Block_3
									end
								end
								i32.const 2
								local.set $$current_block
								i32.const 0
								i32.const 24448
								call $w4snake.Snake.Draw
								i32.const 1
								i32.const 2
								i32.const 3
								i32.const 4
								call $syscall$wasm4.SetDrawColors
								i32.const 24384
								i32.load
								call $runtime.Block.Retain
								i32.const 24384
								i32.load offset=4
								i32.const 24384
								i32.load offset=8
								i32.const 24384
								i32.load offset=12
								local.set $$t6.3
								local.set $$t6.2
								local.set $$t6.1
								local.get $$t6.0
								call $runtime.Block.Release
								local.set $$t6.0
								i32.const 24380
								i32.load8_u align=1
								local.set $$t7
								local.get $$t7
								local.set $$t8
								i32.const 24381
								i32.load8_u align=1
								local.set $$t9
								local.get $$t9
								local.set $$t10
								local.get $$t6.0
								local.get $$t6.1
								local.get $$t6.2
								local.get $$t6.3
								local.get $$t8
								local.get $$t10
								i32.const 8
								i32.const 8
								i32.const 1
								call $syscall$wasm4.BlitI32
								br $$BlockFnBody
							end
							i32.const 3
							local.set $$current_block
							i32.const 0
							i32.const 24448
							call $w4snake.Snake.Reset
							br $$Block_3
						end
						i32.const 4
						local.set $$current_block
						i32.const 24448
						i32.load
						call $runtime.Block.Retain
						i32.const 24448
						i32.load offset=4
						i32.const 24448
						i32.load offset=8
						i32.const 24448
						i32.load offset=12
						local.set $$t11.3
						local.set $$t11.2
						local.set $$t11.1
						local.get $$t11.0
						call $runtime.Block.Release
						local.set $$t11.0
						local.get $$t11.0
						call $runtime.Block.Retain
						local.get $$t11.1
						i32.const 2
						i32.const 0
						i32.mul
						i32.add
						local.set $$t12.1
						local.get $$t12.0
						call $runtime.Block.Release
						local.set $$t12.0
						local.get $$t12.1
						i32.load8_u align=1
						local.get $$t12.1
						i32.load8_u offset=1 align=1
						local.set $$t13.1
						local.set $$t13.0
						i32.const 24380
						i32.load8_u align=1
						i32.const 24380
						i32.load8_u offset=1 align=1
						local.set $$t14.1
						local.set $$t14.0
						local.get $$t13.0
						local.get $$t14.0
						i32.eq
						local.get $$t13.1
						local.get $$t14.1
						i32.eq
						i32.and
						local.set $$t15
						local.get $$t15
						if
							br $$Block_4
						else
							i32.const 2
							local.set $$block_selector
							br $$BlockDisp
						end
					end
					i32.const 5
					local.set $$current_block
					i32.const 24448
					i32.load
					call $runtime.Block.Retain
					i32.const 24448
					i32.load offset=4
					i32.const 24448
					i32.load offset=8
					i32.const 24448
					i32.load offset=12
					local.set $$t16.3
					local.set $$t16.2
					local.set $$t16.1
					local.get $$t16.0
					call $runtime.Block.Release
					local.set $$t16.0
					i32.const 24448
					i32.load
					call $runtime.Block.Retain
					i32.const 24448
					i32.load offset=4
					i32.const 24448
					i32.load offset=8
					i32.const 24448
					i32.load offset=12
					local.set $$t17.3
					local.set $$t17.2
					local.set $$t17.1
					local.get $$t17.0
					call $runtime.Block.Release
					local.set $$t17.0
					i32.const 24448
					i32.load
					call $runtime.Block.Retain
					i32.const 24448
					i32.load offset=4
					i32.const 24448
					i32.load offset=8
					i32.const 24448
					i32.load offset=12
					local.set $$t18.3
					local.set $$t18.2
					local.set $$t18.1
					local.get $$t18.0
					call $runtime.Block.Release
					local.set $$t18.0
					local.get $$t18.2
					local.set $$t19
					local.get $$t19
					i32.const 1
					i32.sub
					local.set $$t20
					local.get $$t17.0
					call $runtime.Block.Retain
					local.get $$t17.1
					i32.const 2
					local.get $$t20
					i32.mul
					i32.add
					local.set $$t21.1
					local.get $$t21.0
					call $runtime.Block.Release
					local.set $$t21.0
					local.get $$t21.1
					i32.load8_u align=1
					local.get $$t21.1
					i32.load8_u offset=1 align=1
					local.set $$t22.1
					local.set $$t22.0
					i32.const 18
					call $runtime.HeapAlloc
					i32.const 1
					i32.const 0
					i32.const 2
					call $runtime.Block.Init
					call $runtime.DupI32
					i32.const 16
					i32.add
					local.set $$t23.1
					local.get $$t23.0
					call $runtime.Block.Release
					local.set $$t23.0
					local.get $$t23.0
					call $runtime.Block.Retain
					local.get $$t23.1
					i32.const 2
					i32.const 0
					i32.mul
					i32.add
					local.set $$t24.1
					local.get $$t24.0
					call $runtime.Block.Release
					local.set $$t24.0
					local.get $$t24.1
					local.get $$t22.0
					i32.store8 align=1
					local.get $$t24.1
					local.get $$t22.1
					i32.store8 offset=1 align=1
					local.get $$t23.0
					call $runtime.Block.Retain
					local.get $$t23.1
					i32.const 2
					i32.const 0
					i32.mul
					i32.add
					i32.const 1
					i32.const 0
					i32.sub
					i32.const 1
					i32.const 0
					i32.sub
					local.set $$t25.3
					local.set $$t25.2
					local.set $$t25.1
					local.get $$t25.0
					call $runtime.Block.Release
					local.set $$t25.0
					local.get $$t16.0
					local.get $$t16.1
					local.get $$t16.2
					local.get $$t16.3
					local.get $$t25.0
					local.get $$t25.1
					local.get $$t25.2
					local.get $$t25.3
					call $$w4snake.Point.$slice.append
					local.set $$t26.3
					local.set $$t26.2
					local.set $$t26.1
					local.get $$t26.0
					call $runtime.Block.Release
					local.set $$t26.0
					i32.const 24448
					local.get $$t26.0
					call $runtime.Block.Retain
					i32.const 24448
					i32.load align=1
					call $runtime.Block.Release
					i32.store align=1
					i32.const 24448
					local.get $$t26.1
					i32.store offset=4
					i32.const 24448
					local.get $$t26.2
					i32.store offset=8
					i32.const 24448
					local.get $$t26.3
					i32.store offset=12
					i32.const 24400
					i32.load
					i32.const 24400
					i32.load offset=4
					call $runtime.Block.Retain
					i32.const 24400
					i32.load offset=8
					local.set $$t27.1.1
					local.get $$t27.1.0
					call $runtime.Block.Release
					local.set $$t27.1.0
					local.set $$t27.0
					i32.const 20
					local.get $$t27.0
					local.get $$t27.1.1
					global.set $$wa.runtime.closure_data
					call_indirect 0 (type $$$fnSig7)
					local.set $$t28
					local.get $$t28
					i32.const 8
					i32.mul
					local.set $$t29
					local.get $$t29
					i32.const 255
					i32.and
					local.set $$t30
					i32.const 24380
					local.get $$t30
					i32.store8 align=1
					i32.const 24400
					i32.load
					i32.const 24400
					i32.load offset=4
					call $runtime.Block.Retain
					i32.const 24400
					i32.load offset=8
					local.set $$t31.1.1
					local.get $$t31.1.0
					call $runtime.Block.Release
					local.set $$t31.1.0
					local.set $$t31.0
					i32.const 20
					local.get $$t31.0
					local.get $$t31.1.1
					global.set $$wa.runtime.closure_data
					call_indirect 0 (type $$$fnSig7)
					local.set $$t32
					local.get $$t32
					i32.const 8
					i32.mul
					local.set $$t33
					local.get $$t33
					i32.const 255
					i32.and
					local.set $$t34
					i32.const 24381
					local.get $$t34
					i32.store8 align=1
					i32.const 0
					i32.const 24412
					call $w4snake.Sound.play
					i32.const 2
					local.set $$block_selector
					br $$BlockDisp
				end
			end
		end
		local.get $$t6.0
		call $runtime.Block.Release
		local.get $$t11.0
		call $runtime.Block.Release
		local.get $$t12.0
		call $runtime.Block.Release
		local.get $$t16.0
		call $runtime.Block.Release
		local.get $$t17.0
		call $runtime.Block.Release
		local.get $$t18.0
		call $runtime.Block.Release
		local.get $$t21.0
		call $runtime.Block.Release
		local.get $$t23.0
		call $runtime.Block.Release
		local.get $$t24.0
		call $runtime.Block.Release
		local.get $$t25.0
		call $runtime.Block.Release
		local.get $$t26.0
		call $runtime.Block.Release
		local.get $$t27.1.0
		call $runtime.Block.Release
		local.get $$t31.1.0
		call $runtime.Block.Release
	)
	(func $math$rand.Rand.Int31 (param $this.0 i32) (param $this.1 i32) (result i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$ret_0 i32)
		(local $$t0 i64)
		(local $$t1 i64)
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
					local.get $this.0
					local.get $this.1
					call $math$rand.Rand.Int63
					local.set $$t0
					local.get $$t0
					i64.const 32
					i64.shr_s
					local.set $$t1
					local.get $$t1
					i32.wrap_i64
					local.set $$t2
					local.get $$t2
					local.set $$ret_0
					br $$BlockFnBody
				end
			end
		end
		local.get $$ret_0
	)
	(func $math$rand.Rand.Int31n (param $this.0 i32) (param $this.1 i32) (param $n i32) (result i32)
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
		(local $$t10 i32)
		(local $$t11 i32)
		(local $$t12 i32)
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
													br_table  0 1 2 3 4 5 6 7 0
												end
												i32.const 0
												local.set $$current_block
												local.get $n
												i32.const 0
												i32.le_s
												local.set $$t0
												local.get $$t0
												if
													br $$Block_0
												else
													br $$Block_1
												end
											end
											i32.const 1
											local.set $$current_block
											i32.const 24605
											i32.const 26
											i32.const 24631
											i32.const 13
											call $$runtime.panic_
										end
										i32.const 2
										local.set $$current_block
										local.get $n
										i32.const 1
										i32.sub
										local.set $$t1
										local.get $n
										local.get $$t1
										i32.and
										local.set $$t2
										local.get $$t2
										i32.const 0
										i32.eq
										local.set $$t3
										local.get $$t3
										if
											br $$Block_2
										else
											br $$Block_3
										end
									end
									i32.const 3
									local.set $$current_block
									local.get $this.0
									local.get $this.1
									call $math$rand.Rand.Int31
									local.set $$t4
									local.get $n
									i32.const 1
									i32.sub
									local.set $$t5
									local.get $$t4
									local.get $$t5
									i32.and
									local.set $$t6
									local.get $$t6
									local.set $$ret_0
									br $$BlockFnBody
								end
								i32.const 4
								local.set $$current_block
								local.get $n
								local.set $$t7
								i32.const -2147483648
								local.get $$t7
								i32.rem_u
								local.set $$t8
								i32.const 2147483647
								local.get $$t8
								i32.sub
								local.set $$t9
								local.get $$t9
								local.set $$t10
								local.get $this.0
								local.get $this.1
								call $math$rand.Rand.Int31
								local.set $$t11
								br $$Block_6
							end
							i32.const 5
							local.set $$current_block
							local.get $this.0
							local.get $this.1
							call $math$rand.Rand.Int31
							local.set $$t12
							br $$Block_6
						end
						i32.const 6
						local.set $$current_block
						local.get $$t13
						local.get $n
						i32.rem_s
						local.set $$t14
						local.get $$t14
						local.set $$ret_0
						br $$BlockFnBody
					end
					local.get $$current_block
					i32.const 4
					i32.eq
					if(result i32)
						local.get $$t11
					else
						local.get $$t12
					end
					local.set $$t13
					i32.const 7
					local.set $$current_block
					local.get $$t13
					local.get $$t10
					i32.gt_s
					local.set $$t15
					local.get $$t15
					if
						i32.const 5
						local.set $$block_selector
						br $$BlockDisp
					else
						i32.const 6
						local.set $$block_selector
						br $$BlockDisp
					end
				end
			end
		end
		local.get $$ret_0
	)
	(func $math$rand.Rand.Int63 (param $this.0 i32) (param $this.1 i32) (result i64)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$ret_0 i64)
		(local $$t0.0 i32)
		(local $$t0.1 i32)
		(local $$t1.0.0 i32)
		(local $$t1.0.1 i32)
		(local $$t1.1 i32)
		(local $$t1.2 i32)
		(local $$t2 i64)
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
					call $runtime.Block.Retain
					local.get $$t0.1
					i32.load offset=4
					local.get $$t0.1
					i32.load offset=8
					local.get $$t0.1
					i32.load offset=12
					local.set $$t1.2
					local.set $$t1.1
					local.set $$t1.0.1
					local.get $$t1.0.0
					call $runtime.Block.Release
					local.set $$t1.0.0
					local.get $$t1.0.0
					local.get $$t1.0.1
					local.get $$t1.1
					i32.load offset=8
					call_indirect 0 (type $$$fnSig2)
					local.set $$t2
					local.get $$t2
					local.set $$ret_0
					br $$BlockFnBody
				end
			end
		end
		local.get $$ret_0
		local.get $$t0.0
		call $runtime.Block.Release
		local.get $$t1.0.0
		call $runtime.Block.Release
	)
	(func $math$rand.Rand.Int63n (param $this.0 i32) (param $this.1 i32) (param $n i64) (result i64)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$ret_0 i64)
		(local $$t0 i32)
		(local $$t1 i64)
		(local $$t2 i64)
		(local $$t3 i32)
		(local $$t4 i64)
		(local $$t5 i64)
		(local $$t6 i64)
		(local $$t7 i64)
		(local $$t8 i64)
		(local $$t9 i64)
		(local $$t10 i64)
		(local $$t11 i64)
		(local $$t12 i64)
		(local $$t13 i64)
		(local $$t14 i64)
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
													br_table  0 1 2 3 4 5 6 7 0
												end
												i32.const 0
												local.set $$current_block
												local.get $n
												i64.const 0
												i64.le_s
												local.set $$t0
												local.get $$t0
												if
													br $$Block_0
												else
													br $$Block_1
												end
											end
											i32.const 1
											local.set $$current_block
											i32.const 24644
											i32.const 26
											i32.const 24670
											i32.const 12
											call $$runtime.panic_
										end
										i32.const 2
										local.set $$current_block
										local.get $n
										i64.const 1
										i64.sub
										local.set $$t1
										local.get $n
										local.get $$t1
										i64.and
										local.set $$t2
										local.get $$t2
										i64.const 0
										i64.eq
										local.set $$t3
										local.get $$t3
										if
											br $$Block_2
										else
											br $$Block_3
										end
									end
									i32.const 3
									local.set $$current_block
									local.get $this.0
									local.get $this.1
									call $math$rand.Rand.Int63
									local.set $$t4
									local.get $n
									i64.const 1
									i64.sub
									local.set $$t5
									local.get $$t4
									local.get $$t5
									i64.and
									local.set $$t6
									local.get $$t6
									local.set $$ret_0
									br $$BlockFnBody
								end
								i32.const 4
								local.set $$current_block
								local.get $n
								local.set $$t7
								i64.const -9223372036854775808
								local.get $$t7
								i64.rem_u
								local.set $$t8
								i64.const 9223372036854775807
								local.get $$t8
								i64.sub
								local.set $$t9
								local.get $$t9
								local.set $$t10
								local.get $this.0
								local.get $this.1
								call $math$rand.Rand.Int63
								local.set $$t11
								br $$Block_6
							end
							i32.const 5
							local.set $$current_block
							local.get $this.0
							local.get $this.1
							call $math$rand.Rand.Int63
							local.set $$t12
							br $$Block_6
						end
						i32.const 6
						local.set $$current_block
						local.get $$t13
						local.get $n
						i64.rem_s
						local.set $$t14
						local.get $$t14
						local.set $$ret_0
						br $$BlockFnBody
					end
					local.get $$current_block
					i32.const 4
					i32.eq
					if(result i64)
						local.get $$t11
					else
						local.get $$t12
					end
					local.set $$t13
					i32.const 7
					local.set $$current_block
					local.get $$t13
					local.get $$t10
					i64.gt_s
					local.set $$t15
					local.get $$t15
					if
						i32.const 5
						local.set $$block_selector
						br $$BlockDisp
					else
						i32.const 6
						local.set $$block_selector
						br $$BlockDisp
					end
				end
			end
		end
		local.get $$ret_0
	)
	(func $math$rand.Rand.Intn (param $this.0 i32) (param $this.1 i32) (param $n i32) (result i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$ret_0 i32)
		(local $$t0 i32)
		(local $$t1 i32)
		(local $$t2 i32)
		(local $$t3 i32)
		(local $$t4 i32)
		(local $$t5 i64)
		(local $$t6 i64)
		(local $$t7 i32)
		block $$BlockFnBody
			loop $$BlockDisp
				block $$Block_4
					block $$Block_3
						block $$Block_2
							block $$Block_1
								block $$Block_0
									block $$BlockSel
										local.get $$block_selector
										br_table  0 1 2 3 4 0
									end
									i32.const 0
									local.set $$current_block
									local.get $n
									i32.const 0
									i32.le_s
									local.set $$t0
									local.get $$t0
									if
										br $$Block_0
									else
										br $$Block_1
									end
								end
								i32.const 1
								local.set $$current_block
								i32.const 24682
								i32.const 24
								i32.const 24706
								i32.const 13
								call $$runtime.panic_
							end
							i32.const 2
							local.set $$current_block
							local.get $n
							i32.const 2147483647
							i32.le_s
							local.set $$t1
							local.get $$t1
							if
								br $$Block_2
							else
								br $$Block_3
							end
						end
						i32.const 3
						local.set $$current_block
						local.get $n
						local.set $$t2
						local.get $this.0
						local.get $this.1
						local.get $$t2
						call $math$rand.Rand.Int31n
						local.set $$t3
						local.get $$t3
						local.set $$t4
						local.get $$t4
						local.set $$ret_0
						br $$BlockFnBody
					end
					i32.const 4
					local.set $$current_block
					local.get $n
					i64.extend_i32_s
					local.set $$t5
					local.get $this.0
					local.get $this.1
					local.get $$t5
					call $math$rand.Rand.Int63n
					local.set $$t6
					local.get $$t6
					i32.wrap_i64
					local.set $$t7
					local.get $$t7
					local.set $$ret_0
					br $$BlockFnBody
				end
			end
		end
		local.get $$ret_0
	)
	(func $w4snake.Sound.play (param $this.0 i32) (param $this.1 i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$t0.0 i32)
		(local $$t0.1 i32)
		(local $$t1 i32)
		(local $$t2.0 i32)
		(local $$t2.1 i32)
		(local $$t3 i32)
		(local $$t4 i32)
		(local $$t5 i32)
		(local $$t6.0 i32)
		(local $$t6.1 i32)
		(local $$t7 i32)
		(local $$t8 i32)
		(local $$t9.0 i32)
		(local $$t9.1 i32)
		(local $$t10 i32)
		(local $$t11 i32)
		(local $$t12 i32)
		(local $$t13.0 i32)
		(local $$t13.1 i32)
		(local $$t14 i32)
		(local $$t15 i32)
		(local $$t16.0 i32)
		(local $$t16.1 i32)
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
		(local $$t26 i32)
		(local $$t27 i32)
		(local $$t28.0 i32)
		(local $$t28.1 i32)
		(local $$t29 i32)
		(local $$t30 i32)
		(local $$t31 i32)
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
					local.get $this.0
					call $runtime.Block.Retain
					local.get $this.1
					i32.const 4
					i32.add
					local.set $$t2.1
					local.get $$t2.0
					call $runtime.Block.Release
					local.set $$t2.0
					local.get $$t2.1
					i32.load
					local.set $$t3
					local.get $$t3
					i64.const 16
					i32.wrap_i64
					i32.shl
					local.set $$t4
					local.get $$t1
					local.get $$t4
					i32.or
					local.set $$t5
					local.get $this.0
					call $runtime.Block.Retain
					local.get $this.1
					i32.const 8
					i32.add
					local.set $$t6.1
					local.get $$t6.0
					call $runtime.Block.Release
					local.set $$t6.0
					local.get $$t6.1
					i32.load
					local.set $$t7
					local.get $$t7
					i64.const 24
					i32.wrap_i64
					i32.shl
					local.set $$t8
					local.get $this.0
					call $runtime.Block.Retain
					local.get $this.1
					i32.const 12
					i32.add
					local.set $$t9.1
					local.get $$t9.0
					call $runtime.Block.Release
					local.set $$t9.0
					local.get $$t9.1
					i32.load
					local.set $$t10
					local.get $$t10
					i64.const 16
					i32.wrap_i64
					i32.shl
					local.set $$t11
					local.get $$t8
					local.get $$t11
					i32.or
					local.set $$t12
					local.get $this.0
					call $runtime.Block.Retain
					local.get $this.1
					i32.const 16
					i32.add
					local.set $$t13.1
					local.get $$t13.0
					call $runtime.Block.Release
					local.set $$t13.0
					local.get $$t13.1
					i32.load
					local.set $$t14
					local.get $$t12
					local.get $$t14
					i32.or
					local.set $$t15
					local.get $this.0
					call $runtime.Block.Retain
					local.get $this.1
					i32.const 20
					i32.add
					local.set $$t16.1
					local.get $$t16.0
					call $runtime.Block.Release
					local.set $$t16.0
					local.get $$t16.1
					i32.load
					local.set $$t17
					local.get $$t17
					i64.const 8
					i32.wrap_i64
					i32.shl
					local.set $$t18
					local.get $$t15
					local.get $$t18
					i32.or
					local.set $$t19
					local.get $this.0
					call $runtime.Block.Retain
					local.get $this.1
					i32.const 28
					i32.add
					local.set $$t20.1
					local.get $$t20.0
					call $runtime.Block.Release
					local.set $$t20.0
					local.get $$t20.1
					i32.load
					local.set $$t21
					local.get $this.0
					call $runtime.Block.Retain
					local.get $this.1
					i32.const 32
					i32.add
					local.set $$t22.1
					local.get $$t22.0
					call $runtime.Block.Release
					local.set $$t22.0
					local.get $$t22.1
					i32.load
					local.set $$t23
					local.get $$t23
					i64.const 2
					i32.wrap_i64
					i32.shl
					local.set $$t24
					local.get $$t21
					local.get $$t24
					i32.or
					local.set $$t25
					local.get $$t5
					local.set $$t26
					local.get $$t19
					local.set $$t27
					local.get $this.0
					call $runtime.Block.Retain
					local.get $this.1
					i32.const 24
					i32.add
					local.set $$t28.1
					local.get $$t28.0
					call $runtime.Block.Release
					local.set $$t28.0
					local.get $$t28.1
					i32.load
					local.set $$t29
					local.get $$t29
					local.set $$t30
					local.get $$t25
					local.set $$t31
					local.get $$t26
					local.get $$t27
					local.get $$t30
					local.get $$t31
					call $syscall$wasm4.Tone
					br $$BlockFnBody
				end
			end
		end
		local.get $$t0.0
		call $runtime.Block.Release
		local.get $$t2.0
		call $runtime.Block.Release
		local.get $$t6.0
		call $runtime.Block.Release
		local.get $$t9.0
		call $runtime.Block.Release
		local.get $$t13.0
		call $runtime.Block.Release
		local.get $$t16.0
		call $runtime.Block.Release
		local.get $$t20.0
		call $runtime.Block.Release
		local.get $$t22.0
		call $runtime.Block.Release
		local.get $$t28.0
		call $runtime.Block.Release
	)
	(func $w4snake.Snake.Down (param $this.0 i32) (param $this.1 i32)
		(local $$block_selector i32)
		(local $$current_block i32)
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
		(local $$t7.0 i32)
		(local $$t7.1 i32)
		(local $$t8.0 i32)
		(local $$t8.1 i32)
		(local $$t9.0 i32)
		(local $$t9.1 i32)
		block $$BlockFnBody
			loop $$BlockDisp
				block $$Block_4
					block $$Block_3
						block $$Block_2
							block $$Block_1
								block $$Block_0
									block $$BlockSel
										local.get $$block_selector
										br_table  0 1 2 3 4 0
									end
									i32.const 0
									local.set $$current_block
									local.get $this.0
									call $runtime.Block.Retain
									local.get $this.1
									i32.const 24
									i32.add
									local.set $$t0.1
									local.get $$t0.0
									call $runtime.Block.Release
									local.set $$t0.0
									local.get $$t0.1
									i32.load8_u align=1
									local.set $$t1
									local.get $$t1
									if
										br $$Block_0
									else
										br $$Block_1
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
							i32.const 16
							i32.add
							local.set $$t2.1
							local.get $$t2.0
							call $runtime.Block.Release
							local.set $$t2.0
							local.get $$t2.0
							call $runtime.Block.Retain
							local.get $$t2.1
							i32.const 4
							i32.add
							local.set $$t3.1
							local.get $$t3.0
							call $runtime.Block.Release
							local.set $$t3.0
							local.get $$t3.1
							i32.load
							local.set $$t4
							local.get $$t4
							i32.const 0
							i32.eq
							local.set $$t5
							local.get $$t5
							if
								br $$Block_2
							else
								br $$Block_3
							end
						end
						i32.const 3
						local.set $$current_block
						local.get $this.0
						call $runtime.Block.Retain
						local.get $this.1
						i32.const 16
						i32.add
						local.set $$t6.1
						local.get $$t6.0
						call $runtime.Block.Release
						local.set $$t6.0
						local.get $$t6.0
						call $runtime.Block.Retain
						local.get $$t6.1
						i32.const 0
						i32.add
						local.set $$t7.1
						local.get $$t7.0
						call $runtime.Block.Release
						local.set $$t7.0
						local.get $this.0
						call $runtime.Block.Retain
						local.get $this.1
						i32.const 16
						i32.add
						local.set $$t8.1
						local.get $$t8.0
						call $runtime.Block.Release
						local.set $$t8.0
						local.get $$t8.0
						call $runtime.Block.Retain
						local.get $$t8.1
						i32.const 4
						i32.add
						local.set $$t9.1
						local.get $$t9.0
						call $runtime.Block.Release
						local.set $$t9.0
						local.get $$t7.1
						i32.const 0
						i32.store
						local.get $$t9.1
						i32.const 8
						i32.store
						br $$Block_3
					end
					i32.const 4
					local.set $$current_block
					br $$BlockFnBody
				end
			end
		end
		local.get $$t0.0
		call $runtime.Block.Release
		local.get $$t2.0
		call $runtime.Block.Release
		local.get $$t3.0
		call $runtime.Block.Release
		local.get $$t6.0
		call $runtime.Block.Release
		local.get $$t7.0
		call $runtime.Block.Release
		local.get $$t8.0
		call $runtime.Block.Release
		local.get $$t9.0
		call $runtime.Block.Release
	)
	(func $w4snake.Snake.Draw (param $this.0 i32) (param $this.1 i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$t0.0 i32)
		(local $$t0.1 i32)
		(local $$t1.0 i32)
		(local $$t1.1 i32)
		(local $$t2.0 i32)
		(local $$t2.1 i32)
		(local $$t2.2 i32)
		(local $$t2.3 i32)
		(local $$t3 i32)
		(local $$t4 i32)
		(local $$t5 i32)
		(local $$t6 i32)
		(local $$t7.0 i32)
		(local $$t7.1 i32)
		(local $$t8.0 i32)
		(local $$t8.1 i32)
		(local $$t9.0 i32)
		(local $$t9.1 i32)
		(local $$t10 i32)
		(local $$t11 i32)
		(local $$t12.0 i32)
		(local $$t12.1 i32)
		(local $$t13 i32)
		(local $$t14 i32)
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
		(local $$t19.0 i32)
		(local $$t19.1 i32)
		(local $$t20.0 i32)
		(local $$t20.1 i32)
		(local $$t21 i32)
		(local $$t22 i32)
		(local $$t23.0 i32)
		(local $$t23.1 i32)
		(local $$t24 i32)
		(local $$t25 i32)
		block $$BlockFnBody
			loop $$BlockDisp
				block $$Block_3
					block $$Block_2
						block $$Block_1
							block $$Block_0
								block $$BlockSel
									local.get $$block_selector
									br_table  0 1 2 3 0
								end
								i32.const 0
								local.set $$current_block
								i32.const 3
								i32.const 4
								i32.const 3
								i32.const 4
								call $syscall$wasm4.SetDrawColors
								i32.const 18
								call $runtime.HeapAlloc
								i32.const 1
								i32.const 0
								i32.const 2
								call $runtime.Block.Init
								call $runtime.DupI32
								i32.const 16
								i32.add
								local.set $$t0.1
								local.get $$t0.0
								call $runtime.Block.Release
								local.set $$t0.0
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
								call $runtime.Block.Retain
								local.get $$t1.1
								i32.load offset=4
								local.get $$t1.1
								i32.load offset=8
								local.get $$t1.1
								i32.load offset=12
								local.set $$t2.3
								local.set $$t2.2
								local.set $$t2.1
								local.get $$t2.0
								call $runtime.Block.Release
								local.set $$t2.0
								local.get $$t2.2
								local.set $$t3
								br $$Block_0
							end
							local.get $$current_block
							i32.const 0
							i32.eq
							if(result i32)
								i32.const -1
							else
								local.get $$t4
							end
							local.set $$t5
							i32.const 1
							local.set $$current_block
							local.get $$t5
							i32.const 1
							i32.add
							local.set $$t4
							local.get $$t4
							local.get $$t3
							i32.lt_s
							local.set $$t6
							local.get $$t6
							if
								br $$Block_1
							else
								br $$Block_2
							end
						end
						i32.const 2
						local.set $$current_block
						local.get $$t2.0
						call $runtime.Block.Retain
						local.get $$t2.1
						i32.const 2
						local.get $$t4
						i32.mul
						i32.add
						local.set $$t7.1
						local.get $$t7.0
						call $runtime.Block.Release
						local.set $$t7.0
						local.get $$t7.1
						i32.load8_u align=1
						local.get $$t7.1
						i32.load8_u offset=1 align=1
						local.set $$t8.1
						local.set $$t8.0
						local.get $$t0.1
						local.get $$t8.0
						i32.store8 align=1
						local.get $$t0.1
						local.get $$t8.1
						i32.store8 offset=1 align=1
						local.get $$t0.0
						call $runtime.Block.Retain
						local.get $$t0.1
						i32.const 0
						i32.add
						local.set $$t9.1
						local.get $$t9.0
						call $runtime.Block.Release
						local.set $$t9.0
						local.get $$t9.1
						i32.load8_u align=1
						local.set $$t10
						local.get $$t10
						local.set $$t11
						local.get $$t0.0
						call $runtime.Block.Retain
						local.get $$t0.1
						i32.const 1
						i32.add
						local.set $$t12.1
						local.get $$t12.0
						call $runtime.Block.Release
						local.set $$t12.0
						local.get $$t12.1
						i32.load8_u align=1
						local.set $$t13
						local.get $$t13
						local.set $$t14
						local.get $$t11
						local.get $$t14
						i32.const 8
						i32.const 8
						call $syscall$wasm4.RectI32
						i32.const 1
						local.set $$block_selector
						br $$BlockDisp
					end
					i32.const 3
					local.set $$current_block
					i32.const 4
					i32.const 0
					i32.const 3
					i32.const 4
					call $syscall$wasm4.SetDrawColors
					i32.const 18
					call $runtime.HeapAlloc
					i32.const 1
					i32.const 0
					i32.const 2
					call $runtime.Block.Init
					call $runtime.DupI32
					i32.const 16
					i32.add
					local.set $$t15.1
					local.get $$t15.0
					call $runtime.Block.Release
					local.set $$t15.0
					local.get $this.0
					call $runtime.Block.Retain
					local.get $this.1
					i32.const 0
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
					i32.const 2
					i32.const 0
					i32.mul
					i32.add
					local.set $$t18.1
					local.get $$t18.0
					call $runtime.Block.Release
					local.set $$t18.0
					local.get $$t18.1
					i32.load8_u align=1
					local.get $$t18.1
					i32.load8_u offset=1 align=1
					local.set $$t19.1
					local.set $$t19.0
					local.get $$t15.1
					local.get $$t19.0
					i32.store8 align=1
					local.get $$t15.1
					local.get $$t19.1
					i32.store8 offset=1 align=1
					local.get $$t15.0
					call $runtime.Block.Retain
					local.get $$t15.1
					i32.const 0
					i32.add
					local.set $$t20.1
					local.get $$t20.0
					call $runtime.Block.Release
					local.set $$t20.0
					local.get $$t20.1
					i32.load8_u align=1
					local.set $$t21
					local.get $$t21
					local.set $$t22
					local.get $$t15.0
					call $runtime.Block.Retain
					local.get $$t15.1
					i32.const 1
					i32.add
					local.set $$t23.1
					local.get $$t23.0
					call $runtime.Block.Release
					local.set $$t23.0
					local.get $$t23.1
					i32.load8_u align=1
					local.set $$t24
					local.get $$t24
					local.set $$t25
					local.get $$t22
					local.get $$t25
					i32.const 8
					i32.const 8
					call $syscall$wasm4.RectI32
					br $$BlockFnBody
				end
			end
		end
		local.get $$t0.0
		call $runtime.Block.Release
		local.get $$t1.0
		call $runtime.Block.Release
		local.get $$t2.0
		call $runtime.Block.Release
		local.get $$t7.0
		call $runtime.Block.Release
		local.get $$t9.0
		call $runtime.Block.Release
		local.get $$t12.0
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
		local.get $$t23.0
		call $runtime.Block.Release
	)
	(func $w4snake.Snake.IsDead (param $this.0 i32) (param $this.1 i32) (result i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$ret_0 i32)
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
		(local $$t4.0 i32)
		(local $$t4.1 i32)
		(local $$t5.0 i32)
		(local $$t5.1 i32)
		(local $$t5.2 i32)
		(local $$t5.3 i32)
		(local $$t6 i32)
		(local $$t7.0 i32)
		(local $$t7.1 i32)
		(local $$t8.0 i32)
		(local $$t8.1 i32)
		(local $$t9 i32)
		(local $$t10 i32)
		(local $$t11.0 i32)
		(local $$t11.1 i32)
		(local $$t12.0 i32)
		(local $$t12.1 i32)
		(local $$t12.2 i32)
		(local $$t12.3 i32)
		(local $$t13 i32)
		(local $$t14 i32)
		block $$BlockFnBody
			loop $$BlockDisp
				block $$Block_5
					block $$Block_4
						block $$Block_3
							block $$Block_2
								block $$Block_1
									block $$Block_0
										block $$BlockSel
											local.get $$block_selector
											br_table  0 1 2 3 4 5 0
										end
										i32.const 0
										local.set $$current_block
										br $$Block_2
									end
									i32.const 1
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
									call $runtime.Block.Retain
									local.get $$t0.1
									i32.load offset=4
									local.get $$t0.1
									i32.load offset=8
									local.get $$t0.1
									i32.load offset=12
									local.set $$t1.3
									local.set $$t1.2
									local.set $$t1.1
									local.get $$t1.0
									call $runtime.Block.Release
									local.set $$t1.0
									local.get $$t1.0
									call $runtime.Block.Retain
									local.get $$t1.1
									i32.const 2
									i32.const 0
									i32.mul
									i32.add
									local.set $$t2.1
									local.get $$t2.0
									call $runtime.Block.Release
									local.set $$t2.0
									local.get $$t2.1
									i32.load8_u align=1
									local.get $$t2.1
									i32.load8_u offset=1 align=1
									local.set $$t3.1
									local.set $$t3.0
									local.get $this.0
									call $runtime.Block.Retain
									local.get $this.1
									i32.const 0
									i32.add
									local.set $$t4.1
									local.get $$t4.0
									call $runtime.Block.Release
									local.set $$t4.0
									local.get $$t4.1
									i32.load
									call $runtime.Block.Retain
									local.get $$t4.1
									i32.load offset=4
									local.get $$t4.1
									i32.load offset=8
									local.get $$t4.1
									i32.load offset=12
									local.set $$t5.3
									local.set $$t5.2
									local.set $$t5.1
									local.get $$t5.0
									call $runtime.Block.Release
									local.set $$t5.0
									local.get $$t5.0
									call $runtime.Block.Retain
									local.get $$t5.1
									i32.const 2
									local.get $$t6
									i32.mul
									i32.add
									local.set $$t7.1
									local.get $$t7.0
									call $runtime.Block.Release
									local.set $$t7.0
									local.get $$t7.1
									i32.load8_u align=1
									local.get $$t7.1
									i32.load8_u offset=1 align=1
									local.set $$t8.1
									local.set $$t8.0
									local.get $$t3.0
									local.get $$t8.0
									i32.eq
									local.get $$t3.1
									local.get $$t8.1
									i32.eq
									i32.and
									local.set $$t9
									local.get $$t9
									if
										br $$Block_3
									else
										br $$Block_4
									end
								end
								i32.const 2
								local.set $$current_block
								i32.const 0
								local.set $$ret_0
								br $$BlockFnBody
							end
							local.get $$current_block
							i32.const 0
							i32.eq
							if(result i32)
								i32.const 1
							else
								local.get $$t10
							end
							local.set $$t6
							i32.const 3
							local.set $$current_block
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
							local.get $$t12.2
							local.set $$t13
							local.get $$t6
							local.get $$t13
							i32.lt_s
							local.set $$t14
							local.get $$t14
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
						i32.const 1
						local.set $$ret_0
						br $$BlockFnBody
					end
					i32.const 5
					local.set $$current_block
					local.get $$t6
					i32.const 1
					i32.add
					local.set $$t10
					i32.const 3
					local.set $$block_selector
					br $$BlockDisp
				end
			end
		end
		local.get $$ret_0
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
		local.get $$t7.0
		call $runtime.Block.Release
		local.get $$t11.0
		call $runtime.Block.Release
		local.get $$t12.0
		call $runtime.Block.Release
	)
	(func $w4snake.Snake.Left (param $this.0 i32) (param $this.1 i32)
		(local $$block_selector i32)
		(local $$current_block i32)
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
		(local $$t7.0 i32)
		(local $$t7.1 i32)
		(local $$t8.0 i32)
		(local $$t8.1 i32)
		(local $$t9.0 i32)
		(local $$t9.1 i32)
		block $$BlockFnBody
			loop $$BlockDisp
				block $$Block_4
					block $$Block_3
						block $$Block_2
							block $$Block_1
								block $$Block_0
									block $$BlockSel
										local.get $$block_selector
										br_table  0 1 2 3 4 0
									end
									i32.const 0
									local.set $$current_block
									local.get $this.0
									call $runtime.Block.Retain
									local.get $this.1
									i32.const 24
									i32.add
									local.set $$t0.1
									local.get $$t0.0
									call $runtime.Block.Release
									local.set $$t0.0
									local.get $$t0.1
									i32.load8_u align=1
									local.set $$t1
									local.get $$t1
									if
										br $$Block_0
									else
										br $$Block_1
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
							i32.const 16
							i32.add
							local.set $$t2.1
							local.get $$t2.0
							call $runtime.Block.Release
							local.set $$t2.0
							local.get $$t2.0
							call $runtime.Block.Retain
							local.get $$t2.1
							i32.const 0
							i32.add
							local.set $$t3.1
							local.get $$t3.0
							call $runtime.Block.Release
							local.set $$t3.0
							local.get $$t3.1
							i32.load
							local.set $$t4
							local.get $$t4
							i32.const 0
							i32.eq
							local.set $$t5
							local.get $$t5
							if
								br $$Block_2
							else
								br $$Block_3
							end
						end
						i32.const 3
						local.set $$current_block
						local.get $this.0
						call $runtime.Block.Retain
						local.get $this.1
						i32.const 16
						i32.add
						local.set $$t6.1
						local.get $$t6.0
						call $runtime.Block.Release
						local.set $$t6.0
						local.get $$t6.0
						call $runtime.Block.Retain
						local.get $$t6.1
						i32.const 0
						i32.add
						local.set $$t7.1
						local.get $$t7.0
						call $runtime.Block.Release
						local.set $$t7.0
						local.get $this.0
						call $runtime.Block.Retain
						local.get $this.1
						i32.const 16
						i32.add
						local.set $$t8.1
						local.get $$t8.0
						call $runtime.Block.Release
						local.set $$t8.0
						local.get $$t8.0
						call $runtime.Block.Retain
						local.get $$t8.1
						i32.const 4
						i32.add
						local.set $$t9.1
						local.get $$t9.0
						call $runtime.Block.Release
						local.set $$t9.0
						local.get $$t7.1
						i32.const -8
						i32.store
						local.get $$t9.1
						i32.const 0
						i32.store
						br $$Block_3
					end
					i32.const 4
					local.set $$current_block
					br $$BlockFnBody
				end
			end
		end
		local.get $$t0.0
		call $runtime.Block.Release
		local.get $$t2.0
		call $runtime.Block.Release
		local.get $$t3.0
		call $runtime.Block.Release
		local.get $$t6.0
		call $runtime.Block.Release
		local.get $$t7.0
		call $runtime.Block.Release
		local.get $$t8.0
		call $runtime.Block.Release
		local.get $$t9.0
		call $runtime.Block.Release
	)
	(func $w4snake.Snake.Reset (param $this.0 i32) (param $this.1 i32)
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
		(local $$t11.2 i32)
		(local $$t11.3 i32)
		(local $$t12.0 i32)
		(local $$t12.1 i32)
		(local $$t13.0 i32)
		(local $$t13.1 i32)
		(local $$t14.0 i32)
		(local $$t14.1 i32)
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
					i32.const 22
					call $runtime.HeapAlloc
					i32.const 1
					i32.const 0
					i32.const 6
					call $runtime.Block.Init
					call $runtime.DupI32
					i32.const 16
					i32.add
					local.set $$t1.1
					local.get $$t1.0
					call $runtime.Block.Release
					local.set $$t1.0
					local.get $$t1.0
					call $runtime.Block.Retain
					local.get $$t1.1
					i32.const 2
					i32.const 0
					i32.mul
					i32.add
					local.set $$t2.1
					local.get $$t2.0
					call $runtime.Block.Release
					local.set $$t2.0
					local.get $$t2.0
					call $runtime.Block.Retain
					local.get $$t2.1
					i32.const 0
					i32.add
					local.set $$t3.1
					local.get $$t3.0
					call $runtime.Block.Release
					local.set $$t3.0
					local.get $$t2.0
					call $runtime.Block.Retain
					local.get $$t2.1
					i32.const 1
					i32.add
					local.set $$t4.1
					local.get $$t4.0
					call $runtime.Block.Release
					local.set $$t4.0
					local.get $$t3.1
					i32.const 16
					i32.store8 align=1
					local.get $$t4.1
					i32.const 0
					i32.store8 align=1
					local.get $$t1.0
					call $runtime.Block.Retain
					local.get $$t1.1
					i32.const 2
					i32.const 1
					i32.mul
					i32.add
					local.set $$t5.1
					local.get $$t5.0
					call $runtime.Block.Release
					local.set $$t5.0
					local.get $$t5.0
					call $runtime.Block.Retain
					local.get $$t5.1
					i32.const 0
					i32.add
					local.set $$t6.1
					local.get $$t6.0
					call $runtime.Block.Release
					local.set $$t6.0
					local.get $$t5.0
					call $runtime.Block.Retain
					local.get $$t5.1
					i32.const 1
					i32.add
					local.set $$t7.1
					local.get $$t7.0
					call $runtime.Block.Release
					local.set $$t7.0
					local.get $$t6.1
					i32.const 8
					i32.store8 align=1
					local.get $$t7.1
					i32.const 0
					i32.store8 align=1
					local.get $$t1.0
					call $runtime.Block.Retain
					local.get $$t1.1
					i32.const 2
					i32.const 2
					i32.mul
					i32.add
					local.set $$t8.1
					local.get $$t8.0
					call $runtime.Block.Release
					local.set $$t8.0
					local.get $$t8.0
					call $runtime.Block.Retain
					local.get $$t8.1
					i32.const 0
					i32.add
					local.set $$t9.1
					local.get $$t9.0
					call $runtime.Block.Release
					local.set $$t9.0
					local.get $$t8.0
					call $runtime.Block.Retain
					local.get $$t8.1
					i32.const 1
					i32.add
					local.set $$t10.1
					local.get $$t10.0
					call $runtime.Block.Release
					local.set $$t10.0
					local.get $$t9.1
					i32.const 0
					i32.store8 align=1
					local.get $$t10.1
					i32.const 0
					i32.store8 align=1
					local.get $$t1.0
					call $runtime.Block.Retain
					local.get $$t1.1
					i32.const 2
					i32.const 0
					i32.mul
					i32.add
					i32.const 3
					i32.const 0
					i32.sub
					i32.const 3
					i32.const 0
					i32.sub
					local.set $$t11.3
					local.set $$t11.2
					local.set $$t11.1
					local.get $$t11.0
					call $runtime.Block.Release
					local.set $$t11.0
					local.get $$t0.1
					local.get $$t11.0
					call $runtime.Block.Retain
					local.get $$t0.1
					i32.load align=1
					call $runtime.Block.Release
					i32.store align=1
					local.get $$t0.1
					local.get $$t11.1
					i32.store offset=4
					local.get $$t0.1
					local.get $$t11.2
					i32.store offset=8
					local.get $$t0.1
					local.get $$t11.3
					i32.store offset=12
					local.get $this.0
					call $runtime.Block.Retain
					local.get $this.1
					i32.const 16
					i32.add
					local.set $$t12.1
					local.get $$t12.0
					call $runtime.Block.Release
					local.set $$t12.0
					local.get $$t12.0
					call $runtime.Block.Retain
					local.get $$t12.1
					i32.const 0
					i32.add
					local.set $$t13.1
					local.get $$t13.0
					call $runtime.Block.Release
					local.set $$t13.0
					local.get $$t12.0
					call $runtime.Block.Retain
					local.get $$t12.1
					i32.const 4
					i32.add
					local.set $$t14.1
					local.get $$t14.0
					call $runtime.Block.Release
					local.set $$t14.0
					local.get $$t13.1
					i32.const 8
					i32.store
					local.get $$t14.1
					i32.const 0
					i32.store
					br $$BlockFnBody
				end
			end
		end
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
	)
	(func $w4snake.Snake.Right (param $this.0 i32) (param $this.1 i32)
		(local $$block_selector i32)
		(local $$current_block i32)
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
		(local $$t7.0 i32)
		(local $$t7.1 i32)
		(local $$t8.0 i32)
		(local $$t8.1 i32)
		(local $$t9.0 i32)
		(local $$t9.1 i32)
		block $$BlockFnBody
			loop $$BlockDisp
				block $$Block_4
					block $$Block_3
						block $$Block_2
							block $$Block_1
								block $$Block_0
									block $$BlockSel
										local.get $$block_selector
										br_table  0 1 2 3 4 0
									end
									i32.const 0
									local.set $$current_block
									local.get $this.0
									call $runtime.Block.Retain
									local.get $this.1
									i32.const 24
									i32.add
									local.set $$t0.1
									local.get $$t0.0
									call $runtime.Block.Release
									local.set $$t0.0
									local.get $$t0.1
									i32.load8_u align=1
									local.set $$t1
									local.get $$t1
									if
										br $$Block_0
									else
										br $$Block_1
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
							i32.const 16
							i32.add
							local.set $$t2.1
							local.get $$t2.0
							call $runtime.Block.Release
							local.set $$t2.0
							local.get $$t2.0
							call $runtime.Block.Retain
							local.get $$t2.1
							i32.const 0
							i32.add
							local.set $$t3.1
							local.get $$t3.0
							call $runtime.Block.Release
							local.set $$t3.0
							local.get $$t3.1
							i32.load
							local.set $$t4
							local.get $$t4
							i32.const 0
							i32.eq
							local.set $$t5
							local.get $$t5
							if
								br $$Block_2
							else
								br $$Block_3
							end
						end
						i32.const 3
						local.set $$current_block
						local.get $this.0
						call $runtime.Block.Retain
						local.get $this.1
						i32.const 16
						i32.add
						local.set $$t6.1
						local.get $$t6.0
						call $runtime.Block.Release
						local.set $$t6.0
						local.get $$t6.0
						call $runtime.Block.Retain
						local.get $$t6.1
						i32.const 0
						i32.add
						local.set $$t7.1
						local.get $$t7.0
						call $runtime.Block.Release
						local.set $$t7.0
						local.get $this.0
						call $runtime.Block.Retain
						local.get $this.1
						i32.const 16
						i32.add
						local.set $$t8.1
						local.get $$t8.0
						call $runtime.Block.Release
						local.set $$t8.0
						local.get $$t8.0
						call $runtime.Block.Retain
						local.get $$t8.1
						i32.const 4
						i32.add
						local.set $$t9.1
						local.get $$t9.0
						call $runtime.Block.Release
						local.set $$t9.0
						local.get $$t7.1
						i32.const 8
						i32.store
						local.get $$t9.1
						i32.const 0
						i32.store
						br $$Block_3
					end
					i32.const 4
					local.set $$current_block
					br $$BlockFnBody
				end
			end
		end
		local.get $$t0.0
		call $runtime.Block.Release
		local.get $$t2.0
		call $runtime.Block.Release
		local.get $$t3.0
		call $runtime.Block.Release
		local.get $$t6.0
		call $runtime.Block.Release
		local.get $$t7.0
		call $runtime.Block.Release
		local.get $$t8.0
		call $runtime.Block.Release
		local.get $$t9.0
		call $runtime.Block.Release
	)
	(func $w4snake.Snake.Up (param $this.0 i32) (param $this.1 i32)
		(local $$block_selector i32)
		(local $$current_block i32)
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
		(local $$t7.0 i32)
		(local $$t7.1 i32)
		(local $$t8.0 i32)
		(local $$t8.1 i32)
		(local $$t9.0 i32)
		(local $$t9.1 i32)
		block $$BlockFnBody
			loop $$BlockDisp
				block $$Block_4
					block $$Block_3
						block $$Block_2
							block $$Block_1
								block $$Block_0
									block $$BlockSel
										local.get $$block_selector
										br_table  0 1 2 3 4 0
									end
									i32.const 0
									local.set $$current_block
									local.get $this.0
									call $runtime.Block.Retain
									local.get $this.1
									i32.const 24
									i32.add
									local.set $$t0.1
									local.get $$t0.0
									call $runtime.Block.Release
									local.set $$t0.0
									local.get $$t0.1
									i32.load8_u align=1
									local.set $$t1
									local.get $$t1
									if
										br $$Block_0
									else
										br $$Block_1
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
							i32.const 16
							i32.add
							local.set $$t2.1
							local.get $$t2.0
							call $runtime.Block.Release
							local.set $$t2.0
							local.get $$t2.0
							call $runtime.Block.Retain
							local.get $$t2.1
							i32.const 4
							i32.add
							local.set $$t3.1
							local.get $$t3.0
							call $runtime.Block.Release
							local.set $$t3.0
							local.get $$t3.1
							i32.load
							local.set $$t4
							local.get $$t4
							i32.const 0
							i32.eq
							local.set $$t5
							local.get $$t5
							if
								br $$Block_2
							else
								br $$Block_3
							end
						end
						i32.const 3
						local.set $$current_block
						local.get $this.0
						call $runtime.Block.Retain
						local.get $this.1
						i32.const 16
						i32.add
						local.set $$t6.1
						local.get $$t6.0
						call $runtime.Block.Release
						local.set $$t6.0
						local.get $$t6.0
						call $runtime.Block.Retain
						local.get $$t6.1
						i32.const 0
						i32.add
						local.set $$t7.1
						local.get $$t7.0
						call $runtime.Block.Release
						local.set $$t7.0
						local.get $this.0
						call $runtime.Block.Retain
						local.get $this.1
						i32.const 16
						i32.add
						local.set $$t8.1
						local.get $$t8.0
						call $runtime.Block.Release
						local.set $$t8.0
						local.get $$t8.0
						call $runtime.Block.Retain
						local.get $$t8.1
						i32.const 4
						i32.add
						local.set $$t9.1
						local.get $$t9.0
						call $runtime.Block.Release
						local.set $$t9.0
						local.get $$t7.1
						i32.const 0
						i32.store
						local.get $$t9.1
						i32.const -8
						i32.store
						br $$Block_3
					end
					i32.const 4
					local.set $$current_block
					br $$BlockFnBody
				end
			end
		end
		local.get $$t0.0
		call $runtime.Block.Release
		local.get $$t2.0
		call $runtime.Block.Release
		local.get $$t3.0
		call $runtime.Block.Release
		local.get $$t6.0
		call $runtime.Block.Release
		local.get $$t7.0
		call $runtime.Block.Release
		local.get $$t8.0
		call $runtime.Block.Release
		local.get $$t9.0
		call $runtime.Block.Release
	)
	(func $w4snake.Snake.Update (param $this.0 i32) (param $this.1 i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$t0.0 i32)
		(local $$t0.1 i32)
		(local $$t1 i32)
		(local $$t2.0 i32)
		(local $$t2.1 i32)
		(local $$t3.0 i32)
		(local $$t3.1 i32)
		(local $$t3.2 i32)
		(local $$t3.3 i32)
		(local $$t4 i32)
		(local $$t5 i32)
		(local $$t6.0 i32)
		(local $$t6.1 i32)
		(local $$t7.0 i32)
		(local $$t7.1 i32)
		(local $$t7.2 i32)
		(local $$t7.3 i32)
		(local $$t8 i32)
		(local $$t9.0 i32)
		(local $$t9.1 i32)
		(local $$t10.0 i32)
		(local $$t10.1 i32)
		(local $$t11.0 i32)
		(local $$t11.1 i32)
		(local $$t11.2 i32)
		(local $$t11.3 i32)
		(local $$t12 i32)
		(local $$t13.0 i32)
		(local $$t13.1 i32)
		(local $$t14.0 i32)
		(local $$t14.1 i32)
		(local $$t15 i32)
		(local $$t16.0 i32)
		(local $$t16.1 i32)
		(local $$t17.0 i32)
		(local $$t17.1 i32)
		(local $$t17.2 i32)
		(local $$t17.3 i32)
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
		(local $$t24 i32)
		(local $$t25 i32)
		(local $$t26.0 i32)
		(local $$t26.1 i32)
		(local $$t27.0 i32)
		(local $$t27.1 i32)
		(local $$t28 i32)
		(local $$t29 i32)
		(local $$t30 i32)
		(local $$t31 i32)
		(local $$t32.0 i32)
		(local $$t32.1 i32)
		(local $$t33.0 i32)
		(local $$t33.1 i32)
		(local $$t33.2 i32)
		(local $$t33.3 i32)
		(local $$t34.0 i32)
		(local $$t34.1 i32)
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
		(local $$t40 i32)
		(local $$t41 i32)
		(local $$t42.0 i32)
		(local $$t42.1 i32)
		(local $$t43.0 i32)
		(local $$t43.1 i32)
		(local $$t44 i32)
		(local $$t45 i32)
		(local $$t46 i32)
		(local $$t47 i32)
		(local $$t48.0 i32)
		(local $$t48.1 i32)
		(local $$t49.0 i32)
		(local $$t49.1 i32)
		(local $$t49.2 i32)
		(local $$t49.3 i32)
		(local $$t50.0 i32)
		(local $$t50.1 i32)
		(local $$t51.0 i32)
		(local $$t51.1 i32)
		(local $$t52 i32)
		(local $$t53 i32)
		(local $$t54 i32)
		(local $$t55.0 i32)
		(local $$t55.1 i32)
		(local $$t56.0 i32)
		(local $$t56.1 i32)
		(local $$t56.2 i32)
		(local $$t56.3 i32)
		(local $$t57.0 i32)
		(local $$t57.1 i32)
		(local $$t58.0 i32)
		(local $$t58.1 i32)
		(local $$t59.0 i32)
		(local $$t59.1 i32)
		(local $$t60.0 i32)
		(local $$t60.1 i32)
		(local $$t60.2 i32)
		(local $$t60.3 i32)
		(local $$t61.0 i32)
		(local $$t61.1 i32)
		(local $$t62.0 i32)
		(local $$t62.1 i32)
		(local $$t63 i32)
		(local $$t64 i32)
		(local $$t65.0 i32)
		(local $$t65.1 i32)
		(local $$t66.0 i32)
		(local $$t66.1 i32)
		(local $$t66.2 i32)
		(local $$t66.3 i32)
		(local $$t67.0 i32)
		(local $$t67.1 i32)
		(local $$t68.0 i32)
		(local $$t68.1 i32)
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
														local.get $this.0
														call $runtime.Block.Retain
														local.get $this.1
														i32.const 24
														i32.add
														local.set $$t0.1
														local.get $$t0.0
														call $runtime.Block.Release
														local.set $$t0.0
														local.get $$t0.1
														i32.load8_u align=1
														local.set $$t1
														local.get $$t1
														if
															br $$Block_0
														else
															br $$Block_1
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
												local.set $$t2.1
												local.get $$t2.0
												call $runtime.Block.Release
												local.set $$t2.0
												local.get $$t2.1
												i32.load
												call $runtime.Block.Retain
												local.get $$t2.1
												i32.load offset=4
												local.get $$t2.1
												i32.load offset=8
												local.get $$t2.1
												i32.load offset=12
												local.set $$t3.3
												local.set $$t3.2
												local.set $$t3.1
												local.get $$t3.0
												call $runtime.Block.Release
												local.set $$t3.0
												local.get $$t3.2
												local.set $$t4
												local.get $$t4
												i32.const 1
												i32.sub
												local.set $$t5
												br $$Block_4
											end
											i32.const 3
											local.set $$current_block
											local.get $this.0
											call $runtime.Block.Retain
											local.get $this.1
											i32.const 0
											i32.add
											local.set $$t6.1
											local.get $$t6.0
											call $runtime.Block.Release
											local.set $$t6.0
											local.get $$t6.1
											i32.load
											call $runtime.Block.Retain
											local.get $$t6.1
											i32.load offset=4
											local.get $$t6.1
											i32.load offset=8
											local.get $$t6.1
											i32.load offset=12
											local.set $$t7.3
											local.set $$t7.2
											local.set $$t7.1
											local.get $$t7.0
											call $runtime.Block.Release
											local.set $$t7.0
											local.get $$t7.0
											call $runtime.Block.Retain
											local.get $$t7.1
											i32.const 2
											local.get $$t8
											i32.mul
											i32.add
											local.set $$t9.1
											local.get $$t9.0
											call $runtime.Block.Release
											local.set $$t9.0
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
											call $runtime.Block.Retain
											local.get $$t10.1
											i32.load offset=4
											local.get $$t10.1
											i32.load offset=8
											local.get $$t10.1
											i32.load offset=12
											local.set $$t11.3
											local.set $$t11.2
											local.set $$t11.1
											local.get $$t11.0
											call $runtime.Block.Release
											local.set $$t11.0
											local.get $$t8
											i32.const 1
											i32.sub
											local.set $$t12
											local.get $$t11.0
											call $runtime.Block.Retain
											local.get $$t11.1
											i32.const 2
											local.get $$t12
											i32.mul
											i32.add
											local.set $$t13.1
											local.get $$t13.0
											call $runtime.Block.Release
											local.set $$t13.0
											local.get $$t13.1
											i32.load8_u align=1
											local.get $$t13.1
											i32.load8_u offset=1 align=1
											local.set $$t14.1
											local.set $$t14.0
											local.get $$t9.1
											local.get $$t14.0
											i32.store8 align=1
											local.get $$t9.1
											local.get $$t14.1
											i32.store8 offset=1 align=1
											local.get $$t8
											i32.const 1
											i32.sub
											local.set $$t15
											br $$Block_4
										end
										i32.const 4
										local.set $$current_block
										local.get $this.0
										call $runtime.Block.Retain
										local.get $this.1
										i32.const 0
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
										i32.const 2
										i32.const 0
										i32.mul
										i32.add
										local.set $$t18.1
										local.get $$t18.0
										call $runtime.Block.Release
										local.set $$t18.0
										local.get $$t18.0
										call $runtime.Block.Retain
										local.get $$t18.1
										i32.const 0
										i32.add
										local.set $$t19.1
										local.get $$t19.0
										call $runtime.Block.Release
										local.set $$t19.0
										local.get $this.0
										call $runtime.Block.Retain
										local.get $this.1
										i32.const 0
										i32.add
										local.set $$t20.1
										local.get $$t20.0
										call $runtime.Block.Release
										local.set $$t20.0
										local.get $$t20.1
										i32.load
										call $runtime.Block.Retain
										local.get $$t20.1
										i32.load offset=4
										local.get $$t20.1
										i32.load offset=8
										local.get $$t20.1
										i32.load offset=12
										local.set $$t21.3
										local.set $$t21.2
										local.set $$t21.1
										local.get $$t21.0
										call $runtime.Block.Release
										local.set $$t21.0
										local.get $$t21.0
										call $runtime.Block.Retain
										local.get $$t21.1
										i32.const 2
										i32.const 0
										i32.mul
										i32.add
										local.set $$t22.1
										local.get $$t22.0
										call $runtime.Block.Release
										local.set $$t22.0
										local.get $$t22.0
										call $runtime.Block.Retain
										local.get $$t22.1
										i32.const 0
										i32.add
										local.set $$t23.1
										local.get $$t23.0
										call $runtime.Block.Release
										local.set $$t23.0
										local.get $$t23.1
										i32.load8_u align=1
										local.set $$t24
										local.get $$t24
										local.set $$t25
										local.get $this.0
										call $runtime.Block.Retain
										local.get $this.1
										i32.const 16
										i32.add
										local.set $$t26.1
										local.get $$t26.0
										call $runtime.Block.Release
										local.set $$t26.0
										local.get $$t26.0
										call $runtime.Block.Retain
										local.get $$t26.1
										i32.const 0
										i32.add
										local.set $$t27.1
										local.get $$t27.0
										call $runtime.Block.Release
										local.set $$t27.0
										local.get $$t27.1
										i32.load
										local.set $$t28
										local.get $$t25
										local.get $$t28
										i32.add
										local.set $$t29
										local.get $$t29
										i32.const 160
										i32.rem_s
										local.set $$t30
										local.get $$t30
										i32.const 255
										i32.and
										local.set $$t31
										local.get $$t19.1
										local.get $$t31
										i32.store8 align=1
										local.get $this.0
										call $runtime.Block.Retain
										local.get $this.1
										i32.const 0
										i32.add
										local.set $$t32.1
										local.get $$t32.0
										call $runtime.Block.Release
										local.set $$t32.0
										local.get $$t32.1
										i32.load
										call $runtime.Block.Retain
										local.get $$t32.1
										i32.load offset=4
										local.get $$t32.1
										i32.load offset=8
										local.get $$t32.1
										i32.load offset=12
										local.set $$t33.3
										local.set $$t33.2
										local.set $$t33.1
										local.get $$t33.0
										call $runtime.Block.Release
										local.set $$t33.0
										local.get $$t33.0
										call $runtime.Block.Retain
										local.get $$t33.1
										i32.const 2
										i32.const 0
										i32.mul
										i32.add
										local.set $$t34.1
										local.get $$t34.0
										call $runtime.Block.Release
										local.set $$t34.0
										local.get $$t34.0
										call $runtime.Block.Retain
										local.get $$t34.1
										i32.const 1
										i32.add
										local.set $$t35.1
										local.get $$t35.0
										call $runtime.Block.Release
										local.set $$t35.0
										local.get $this.0
										call $runtime.Block.Retain
										local.get $this.1
										i32.const 0
										i32.add
										local.set $$t36.1
										local.get $$t36.0
										call $runtime.Block.Release
										local.set $$t36.0
										local.get $$t36.1
										i32.load
										call $runtime.Block.Retain
										local.get $$t36.1
										i32.load offset=4
										local.get $$t36.1
										i32.load offset=8
										local.get $$t36.1
										i32.load offset=12
										local.set $$t37.3
										local.set $$t37.2
										local.set $$t37.1
										local.get $$t37.0
										call $runtime.Block.Release
										local.set $$t37.0
										local.get $$t37.0
										call $runtime.Block.Retain
										local.get $$t37.1
										i32.const 2
										i32.const 0
										i32.mul
										i32.add
										local.set $$t38.1
										local.get $$t38.0
										call $runtime.Block.Release
										local.set $$t38.0
										local.get $$t38.0
										call $runtime.Block.Retain
										local.get $$t38.1
										i32.const 1
										i32.add
										local.set $$t39.1
										local.get $$t39.0
										call $runtime.Block.Release
										local.set $$t39.0
										local.get $$t39.1
										i32.load8_u align=1
										local.set $$t40
										local.get $$t40
										local.set $$t41
										local.get $this.0
										call $runtime.Block.Retain
										local.get $this.1
										i32.const 16
										i32.add
										local.set $$t42.1
										local.get $$t42.0
										call $runtime.Block.Release
										local.set $$t42.0
										local.get $$t42.0
										call $runtime.Block.Retain
										local.get $$t42.1
										i32.const 4
										i32.add
										local.set $$t43.1
										local.get $$t43.0
										call $runtime.Block.Release
										local.set $$t43.0
										local.get $$t43.1
										i32.load
										local.set $$t44
										local.get $$t41
										local.get $$t44
										i32.add
										local.set $$t45
										local.get $$t45
										i32.const 160
										i32.rem_s
										local.set $$t46
										local.get $$t46
										i32.const 255
										i32.and
										local.set $$t47
										local.get $$t35.1
										local.get $$t47
										i32.store8 align=1
										local.get $this.0
										call $runtime.Block.Retain
										local.get $this.1
										i32.const 0
										i32.add
										local.set $$t48.1
										local.get $$t48.0
										call $runtime.Block.Release
										local.set $$t48.0
										local.get $$t48.1
										i32.load
										call $runtime.Block.Retain
										local.get $$t48.1
										i32.load offset=4
										local.get $$t48.1
										i32.load offset=8
										local.get $$t48.1
										i32.load offset=12
										local.set $$t49.3
										local.set $$t49.2
										local.set $$t49.1
										local.get $$t49.0
										call $runtime.Block.Release
										local.set $$t49.0
										local.get $$t49.0
										call $runtime.Block.Retain
										local.get $$t49.1
										i32.const 2
										i32.const 0
										i32.mul
										i32.add
										local.set $$t50.1
										local.get $$t50.0
										call $runtime.Block.Release
										local.set $$t50.0
										local.get $$t50.0
										call $runtime.Block.Retain
										local.get $$t50.1
										i32.const 0
										i32.add
										local.set $$t51.1
										local.get $$t51.0
										call $runtime.Block.Release
										local.set $$t51.0
										local.get $$t51.1
										i32.load8_u align=1
										local.set $$t52
										local.get $$t52
										i32.const 160
										i32.gt_u
										local.set $$t53
										local.get $$t53
										if
											br $$Block_5
										else
											br $$Block_6
										end
									end
									local.get $$current_block
									i32.const 2
									i32.eq
									if(result i32)
										local.get $$t5
									else
										local.get $$t15
									end
									local.set $$t8
									i32.const 5
									local.set $$current_block
									local.get $$t8
									i32.const 0
									i32.gt_s
									local.set $$t54
									local.get $$t54
									if
										i32.const 3
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
								i32.const 0
								i32.add
								local.set $$t55.1
								local.get $$t55.0
								call $runtime.Block.Release
								local.set $$t55.0
								local.get $$t55.1
								i32.load
								call $runtime.Block.Retain
								local.get $$t55.1
								i32.load offset=4
								local.get $$t55.1
								i32.load offset=8
								local.get $$t55.1
								i32.load offset=12
								local.set $$t56.3
								local.set $$t56.2
								local.set $$t56.1
								local.get $$t56.0
								call $runtime.Block.Release
								local.set $$t56.0
								local.get $$t56.0
								call $runtime.Block.Retain
								local.get $$t56.1
								i32.const 2
								i32.const 0
								i32.mul
								i32.add
								local.set $$t57.1
								local.get $$t57.0
								call $runtime.Block.Release
								local.set $$t57.0
								local.get $$t57.0
								call $runtime.Block.Retain
								local.get $$t57.1
								i32.const 0
								i32.add
								local.set $$t58.1
								local.get $$t58.0
								call $runtime.Block.Release
								local.set $$t58.0
								local.get $$t58.1
								i32.const 152
								i32.store8 align=1
								br $$Block_6
							end
							i32.const 7
							local.set $$current_block
							local.get $this.0
							call $runtime.Block.Retain
							local.get $this.1
							i32.const 0
							i32.add
							local.set $$t59.1
							local.get $$t59.0
							call $runtime.Block.Release
							local.set $$t59.0
							local.get $$t59.1
							i32.load
							call $runtime.Block.Retain
							local.get $$t59.1
							i32.load offset=4
							local.get $$t59.1
							i32.load offset=8
							local.get $$t59.1
							i32.load offset=12
							local.set $$t60.3
							local.set $$t60.2
							local.set $$t60.1
							local.get $$t60.0
							call $runtime.Block.Release
							local.set $$t60.0
							local.get $$t60.0
							call $runtime.Block.Retain
							local.get $$t60.1
							i32.const 2
							i32.const 0
							i32.mul
							i32.add
							local.set $$t61.1
							local.get $$t61.0
							call $runtime.Block.Release
							local.set $$t61.0
							local.get $$t61.0
							call $runtime.Block.Retain
							local.get $$t61.1
							i32.const 1
							i32.add
							local.set $$t62.1
							local.get $$t62.0
							call $runtime.Block.Release
							local.set $$t62.0
							local.get $$t62.1
							i32.load8_u align=1
							local.set $$t63
							local.get $$t63
							i32.const 160
							i32.gt_u
							local.set $$t64
							local.get $$t64
							if
								br $$Block_7
							else
								br $$Block_8
							end
						end
						i32.const 8
						local.set $$current_block
						local.get $this.0
						call $runtime.Block.Retain
						local.get $this.1
						i32.const 0
						i32.add
						local.set $$t65.1
						local.get $$t65.0
						call $runtime.Block.Release
						local.set $$t65.0
						local.get $$t65.1
						i32.load
						call $runtime.Block.Retain
						local.get $$t65.1
						i32.load offset=4
						local.get $$t65.1
						i32.load offset=8
						local.get $$t65.1
						i32.load offset=12
						local.set $$t66.3
						local.set $$t66.2
						local.set $$t66.1
						local.get $$t66.0
						call $runtime.Block.Release
						local.set $$t66.0
						local.get $$t66.0
						call $runtime.Block.Retain
						local.get $$t66.1
						i32.const 2
						i32.const 0
						i32.mul
						i32.add
						local.set $$t67.1
						local.get $$t67.0
						call $runtime.Block.Release
						local.set $$t67.0
						local.get $$t67.0
						call $runtime.Block.Retain
						local.get $$t67.1
						i32.const 1
						i32.add
						local.set $$t68.1
						local.get $$t68.0
						call $runtime.Block.Release
						local.set $$t68.0
						local.get $$t68.1
						i32.const 152
						i32.store8 align=1
						br $$Block_8
					end
					i32.const 9
					local.set $$current_block
					br $$BlockFnBody
				end
			end
		end
		local.get $$t0.0
		call $runtime.Block.Release
		local.get $$t2.0
		call $runtime.Block.Release
		local.get $$t3.0
		call $runtime.Block.Release
		local.get $$t6.0
		call $runtime.Block.Release
		local.get $$t7.0
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
		local.get $$t20.0
		call $runtime.Block.Release
		local.get $$t21.0
		call $runtime.Block.Release
		local.get $$t22.0
		call $runtime.Block.Release
		local.get $$t23.0
		call $runtime.Block.Release
		local.get $$t26.0
		call $runtime.Block.Release
		local.get $$t27.0
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
		local.get $$t42.0
		call $runtime.Block.Release
		local.get $$t43.0
		call $runtime.Block.Release
		local.get $$t48.0
		call $runtime.Block.Release
		local.get $$t49.0
		call $runtime.Block.Release
		local.get $$t50.0
		call $runtime.Block.Release
		local.get $$t51.0
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
		local.get $$t65.0
		call $runtime.Block.Release
		local.get $$t66.0
		call $runtime.Block.Release
		local.get $$t67.0
		call $runtime.Block.Release
		local.get $$t68.0
		call $runtime.Block.Release
	)
	(func $math$rand.rngSource.Int63 (param $this.0 i32) (param $this.1 i32) (result i64)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$ret_0 i64)
		(local $$t0 i64)
		(local $$t1 i64)
		(local $$t2 i64)
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
					local.get $this.1
					call $math$rand.rngSource.Uint64
					local.set $$t0
					local.get $$t0
					i64.const 9223372036854775807
					i64.and
					local.set $$t1
					local.get $$t1
					local.set $$t2
					local.get $$t2
					local.set $$ret_0
					br $$BlockFnBody
				end
			end
		end
		local.get $$ret_0
	)
	(func $math$rand.rngSource.Seed (param $this.0 i32) (param $this.1 i32) (param $seed i64)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$t0.0 i32)
		(local $$t0.1 i32)
		(local $$t1.0 i32)
		(local $$t1.1 i32)
		(local $$t2 i64)
		(local $$t3 i32)
		(local $$t4 i64)
		(local $$t5 i64)
		(local $$t6 i32)
		(local $$t7 i64)
		(local $$t8 i32)
		(local $$t9 i32)
		(local $$t10 i32)
		(local $$t11 i32)
		(local $$t12 i32)
		(local $$t13 i32)
		(local $$t14 i32)
		(local $$t15 i32)
		(local $$t16 i64)
		(local $$t17 i64)
		(local $$t18 i32)
		(local $$t19 i64)
		(local $$t20 i64)
		(local $$t21 i64)
		(local $$t22 i32)
		(local $$t23 i64)
		(local $$t24 i64)
		(local $$t25.0 i32)
		(local $$t25.1 i32)
		(local $$t26 i64)
		(local $$t27 i64)
		(local $$t28.0 i32)
		(local $$t28.1 i32)
		(local $$t29.0 i32)
		(local $$t29.1 i32)
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
														i32.const 0
														i32.store
														local.get $this.0
														call $runtime.Block.Retain
														local.get $this.1
														i32.const 4
														i32.add
														local.set $$t1.1
														local.get $$t1.0
														call $runtime.Block.Release
														local.set $$t1.0
														local.get $$t1.1
														i32.const 334
														i32.store
														local.get $seed
														i64.const 2147483647
														i64.rem_s
														local.set $$t2
														local.get $$t2
														i64.const 0
														i64.lt_s
														local.set $$t3
														local.get $$t3
														if
															br $$Block_0
														else
															br $$Block_1
														end
													end
													i32.const 1
													local.set $$current_block
													local.get $$t2
													i64.const 2147483647
													i64.add
													local.set $$t4
													br $$Block_1
												end
												local.get $$current_block
												i32.const 0
												i32.eq
												if(result i64)
													local.get $$t2
												else
													local.get $$t4
												end
												local.set $$t5
												i32.const 2
												local.set $$current_block
												local.get $$t5
												i64.const 0
												i64.eq
												local.set $$t6
												local.get $$t6
												if
													br $$Block_2
												else
													br $$Block_3
												end
											end
											i32.const 3
											local.set $$current_block
											br $$Block_3
										end
										local.get $$current_block
										i32.const 2
										i32.eq
										if(result i64)
											local.get $$t5
										else
											i64.const 89482311
										end
										local.set $$t7
										i32.const 4
										local.set $$current_block
										local.get $$t7
										i32.wrap_i64
										local.set $$t8
										br $$Block_6
									end
									i32.const 5
									local.set $$current_block
									local.get $$t9
									call $math$rand.seedrand
									local.set $$t10
									local.get $$t11
									i32.const 0
									i32.ge_s
									local.set $$t12
									local.get $$t12
									if
										br $$Block_7
									else
										br $$Block_8
									end
								end
								i32.const 6
								local.set $$current_block
								br $$BlockFnBody
							end
							local.get $$current_block
							i32.const 4
							i32.eq
							if(result i32)
								local.get $$t8
							else
								local.get $$t13
							end
							local.get $$current_block
							i32.const 4
							i32.eq
							if(result i32)
								i32.const -20
							else
								local.get $$t14
							end
							local.set $$t11
							local.set $$t9
							i32.const 7
							local.set $$current_block
							local.get $$t11
							i32.const 607
							i32.lt_s
							local.set $$t15
							local.get $$t15
							if
								i32.const 5
								local.set $$block_selector
								br $$BlockDisp
							else
								i32.const 6
								local.set $$block_selector
								br $$BlockDisp
							end
						end
						i32.const 8
						local.set $$current_block
						local.get $$t10
						i64.extend_i32_s
						local.set $$t16
						local.get $$t16
						i64.const 40
						i64.shl
						local.set $$t17
						local.get $$t10
						call $math$rand.seedrand
						local.set $$t18
						local.get $$t18
						i64.extend_i32_s
						local.set $$t19
						local.get $$t19
						i64.const 20
						i64.shl
						local.set $$t20
						local.get $$t17
						local.get $$t20
						i64.xor
						local.set $$t21
						local.get $$t18
						call $math$rand.seedrand
						local.set $$t22
						local.get $$t22
						i64.extend_i32_s
						local.set $$t23
						local.get $$t21
						local.get $$t23
						i64.xor
						local.set $$t24
						i32.const 0
						i32.const 17984
						i32.const 8
						local.get $$t11
						i32.mul
						i32.add
						local.set $$t25.1
						local.get $$t25.0
						call $runtime.Block.Release
						local.set $$t25.0
						local.get $$t25.1
						i64.load
						local.set $$t26
						local.get $$t24
						local.get $$t26
						i64.xor
						local.set $$t27
						local.get $this.0
						call $runtime.Block.Retain
						local.get $this.1
						i32.const 8
						i32.add
						local.set $$t28.1
						local.get $$t28.0
						call $runtime.Block.Release
						local.set $$t28.0
						local.get $$t28.0
						call $runtime.Block.Retain
						local.get $$t28.1
						i32.const 8
						local.get $$t11
						i32.mul
						i32.add
						local.set $$t29.1
						local.get $$t29.0
						call $runtime.Block.Release
						local.set $$t29.0
						local.get $$t29.1
						local.get $$t27
						i64.store align=8
						br $$Block_8
					end
					local.get $$current_block
					i32.const 5
					i32.eq
					if(result i32)
						local.get $$t10
					else
						local.get $$t22
					end
					local.set $$t13
					i32.const 9
					local.set $$current_block
					local.get $$t11
					i32.const 1
					i32.add
					local.set $$t14
					i32.const 7
					local.set $$block_selector
					br $$BlockDisp
				end
			end
		end
		local.get $$t0.0
		call $runtime.Block.Release
		local.get $$t1.0
		call $runtime.Block.Release
		local.get $$t25.0
		call $runtime.Block.Release
		local.get $$t28.0
		call $runtime.Block.Release
		local.get $$t29.0
		call $runtime.Block.Release
	)
	(func $math$rand.rngSource.Uint64 (param $this.0 i32) (param $this.1 i32) (result i64)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$ret_0 i64)
		(local $$t0.0 i32)
		(local $$t0.1 i32)
		(local $$t1 i32)
		(local $$t2 i32)
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
		(local $$t10 i32)
		(local $$t11 i32)
		(local $$t12.0 i32)
		(local $$t12.1 i32)
		(local $$t13 i32)
		(local $$t14 i32)
		(local $$t15.0 i32)
		(local $$t15.1 i32)
		(local $$t16 i32)
		(local $$t17 i32)
		(local $$t18.0 i32)
		(local $$t18.1 i32)
		(local $$t19.0 i32)
		(local $$t19.1 i32)
		(local $$t20 i32)
		(local $$t21.0 i32)
		(local $$t21.1 i32)
		(local $$t22 i64)
		(local $$t23.0 i32)
		(local $$t23.1 i32)
		(local $$t24.0 i32)
		(local $$t24.1 i32)
		(local $$t25 i32)
		(local $$t26.0 i32)
		(local $$t26.1 i32)
		(local $$t27 i64)
		(local $$t28 i64)
		(local $$t29.0 i32)
		(local $$t29.1 i32)
		(local $$t30.0 i32)
		(local $$t30.1 i32)
		(local $$t31 i32)
		(local $$t32.0 i32)
		(local $$t32.1 i32)
		(local $$t33 i64)
		block $$BlockFnBody
			loop $$BlockDisp
				block $$Block_4
					block $$Block_3
						block $$Block_2
							block $$Block_1
								block $$Block_0
									block $$BlockSel
										local.get $$block_selector
										br_table  0 1 2 3 4 0
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
									local.get $$t1
									i32.const 1
									i32.sub
									local.set $$t2
									local.get $$t0.1
									local.get $$t2
									i32.store
									local.get $this.0
									call $runtime.Block.Retain
									local.get $this.1
									i32.const 0
									i32.add
									local.set $$t3.1
									local.get $$t3.0
									call $runtime.Block.Release
									local.set $$t3.0
									local.get $$t3.1
									i32.load
									local.set $$t4
									local.get $$t4
									i32.const 0
									i32.lt_s
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
								local.get $this.0
								call $runtime.Block.Retain
								local.get $this.1
								i32.const 0
								i32.add
								local.set $$t6.1
								local.get $$t6.0
								call $runtime.Block.Release
								local.set $$t6.0
								local.get $$t6.1
								i32.load
								local.set $$t7
								local.get $$t7
								i32.const 607
								i32.add
								local.set $$t8
								local.get $$t6.1
								local.get $$t8
								i32.store
								br $$Block_1
							end
							i32.const 2
							local.set $$current_block
							local.get $this.0
							call $runtime.Block.Retain
							local.get $this.1
							i32.const 4
							i32.add
							local.set $$t9.1
							local.get $$t9.0
							call $runtime.Block.Release
							local.set $$t9.0
							local.get $$t9.1
							i32.load
							local.set $$t10
							local.get $$t10
							i32.const 1
							i32.sub
							local.set $$t11
							local.get $$t9.1
							local.get $$t11
							i32.store
							local.get $this.0
							call $runtime.Block.Retain
							local.get $this.1
							i32.const 4
							i32.add
							local.set $$t12.1
							local.get $$t12.0
							call $runtime.Block.Release
							local.set $$t12.0
							local.get $$t12.1
							i32.load
							local.set $$t13
							local.get $$t13
							i32.const 0
							i32.lt_s
							local.set $$t14
							local.get $$t14
							if
								br $$Block_2
							else
								br $$Block_3
							end
						end
						i32.const 3
						local.set $$current_block
						local.get $this.0
						call $runtime.Block.Retain
						local.get $this.1
						i32.const 4
						i32.add
						local.set $$t15.1
						local.get $$t15.0
						call $runtime.Block.Release
						local.set $$t15.0
						local.get $$t15.1
						i32.load
						local.set $$t16
						local.get $$t16
						i32.const 607
						i32.add
						local.set $$t17
						local.get $$t15.1
						local.get $$t17
						i32.store
						br $$Block_3
					end
					i32.const 4
					local.set $$current_block
					local.get $this.0
					call $runtime.Block.Retain
					local.get $this.1
					i32.const 8
					i32.add
					local.set $$t18.1
					local.get $$t18.0
					call $runtime.Block.Release
					local.set $$t18.0
					local.get $this.0
					call $runtime.Block.Retain
					local.get $this.1
					i32.const 4
					i32.add
					local.set $$t19.1
					local.get $$t19.0
					call $runtime.Block.Release
					local.set $$t19.0
					local.get $$t19.1
					i32.load
					local.set $$t20
					local.get $$t18.0
					call $runtime.Block.Retain
					local.get $$t18.1
					i32.const 8
					local.get $$t20
					i32.mul
					i32.add
					local.set $$t21.1
					local.get $$t21.0
					call $runtime.Block.Release
					local.set $$t21.0
					local.get $$t21.1
					i64.load
					local.set $$t22
					local.get $this.0
					call $runtime.Block.Retain
					local.get $this.1
					i32.const 8
					i32.add
					local.set $$t23.1
					local.get $$t23.0
					call $runtime.Block.Release
					local.set $$t23.0
					local.get $this.0
					call $runtime.Block.Retain
					local.get $this.1
					i32.const 0
					i32.add
					local.set $$t24.1
					local.get $$t24.0
					call $runtime.Block.Release
					local.set $$t24.0
					local.get $$t24.1
					i32.load
					local.set $$t25
					local.get $$t23.0
					call $runtime.Block.Retain
					local.get $$t23.1
					i32.const 8
					local.get $$t25
					i32.mul
					i32.add
					local.set $$t26.1
					local.get $$t26.0
					call $runtime.Block.Release
					local.set $$t26.0
					local.get $$t26.1
					i64.load
					local.set $$t27
					local.get $$t22
					local.get $$t27
					i64.add
					local.set $$t28
					local.get $this.0
					call $runtime.Block.Retain
					local.get $this.1
					i32.const 8
					i32.add
					local.set $$t29.1
					local.get $$t29.0
					call $runtime.Block.Release
					local.set $$t29.0
					local.get $this.0
					call $runtime.Block.Retain
					local.get $this.1
					i32.const 4
					i32.add
					local.set $$t30.1
					local.get $$t30.0
					call $runtime.Block.Release
					local.set $$t30.0
					local.get $$t30.1
					i32.load
					local.set $$t31
					local.get $$t29.0
					call $runtime.Block.Retain
					local.get $$t29.1
					i32.const 8
					local.get $$t31
					i32.mul
					i32.add
					local.set $$t32.1
					local.get $$t32.0
					call $runtime.Block.Release
					local.set $$t32.0
					local.get $$t32.1
					local.get $$t28
					i64.store align=8
					local.get $$t28
					local.set $$t33
					local.get $$t33
					local.set $$ret_0
					br $$BlockFnBody
				end
			end
		end
		local.get $$ret_0
		local.get $$t0.0
		call $runtime.Block.Release
		local.get $$t3.0
		call $runtime.Block.Release
		local.get $$t6.0
		call $runtime.Block.Release
		local.get $$t9.0
		call $runtime.Block.Release
		local.get $$t12.0
		call $runtime.Block.Release
		local.get $$t15.0
		call $runtime.Block.Release
		local.get $$t18.0
		call $runtime.Block.Release
		local.get $$t19.0
		call $runtime.Block.Release
		local.get $$t21.0
		call $runtime.Block.Release
		local.get $$t23.0
		call $runtime.Block.Release
		local.get $$t24.0
		call $runtime.Block.Release
		local.get $$t26.0
		call $runtime.Block.Release
		local.get $$t29.0
		call $runtime.Block.Release
		local.get $$t30.0
		call $runtime.Block.Release
		local.get $$t32.0
		call $runtime.Block.Release
	)
	(func $errors.errorString.Error (param $this.0 i32) (param $this.1 i32) (result i32 i32 i32)
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
					call $runtime.Block.Retain
					local.get $$t0.1
					i32.load offset=4
					local.get $$t0.1
					i32.load offset=8
					local.set $$t1.2
					local.set $$t1.1
					local.get $$t1.0
					call $runtime.Block.Release
					local.set $$t1.0
					local.get $$t1.0
					call $runtime.Block.Retain
					local.get $$t1.1
					local.get $$t1.2
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
		local.get $$ret_0.0
		call $runtime.Block.Release
		local.get $$t0.0
		call $runtime.Block.Release
		local.get $$t1.0
		call $runtime.Block.Release
	)
	(func $_start (export "_start")
		call $w4snake.init
	)
	(func $_main (export "_main"))
	(data (i32.const 14784) "\24\24\77\61\64\73\24\24\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\9b\1a\86\a0\49\fa\a8\bd\05\3f\4e\7b\9d\ee\21\3e\c6\4b\ac\7e\4f\7e\92\be\f5\44\c8\19\a0\01\fa\3e\91\4f\c1\16\6c\c1\56\bf\4b\55\55\55\55\55\a5\3f\cd\9c\d1\1f\fd\d8\e5\3d\5d\1f\29\a9\e5\e5\5a\be\a1\48\7d\56\e3\1d\c7\3e\03\df\bf\19\a0\01\2a\bf\d0\f7\10\11\11\11\81\3f\48\55\55\55\55\55\c5\bf\00\00\80\3f\2f\2a\70\3f\33\85\66\3f\04\28\5f\3f\78\08\59\3f\d5\b5\53\3f\b9\f4\4e\3f\8e\a1\4a\3f\1c\a5\46\3f\44\ef\42\3f\51\74\3f\3f\75\2b\3c\3f\db\0d\39\3f\1a\16\36\3f\d3\3f\33\3f\6e\87\30\3f\eb\e9\2d\3f\c4\64\2b\3f\d1\f5\28\3f\36\9b\26\3f\58\53\24\3f\cd\1c\22\3f\59\f6\1f\3f\e2\de\1d\3f\6d\d5\1b\3f\1a\d9\19\3f\1f\e9\17\3f\c6\04\16\3f\69\2b\14\3f\71\5c\12\3f\56\97\10\3f\99\db\0e\3f\c6\28\0d\3f\73\7e\0b\3f\3e\dc\09\3f\ca\41\08\3f\c4\ae\06\3f\dc\22\05\3f\ca\9d\03\3f\47\1f\02\3f\16\a7\00\3f\f0\69\fe\3e\6c\91\fb\3e\37\c4\f8\3e\ea\01\f6\3e\2a\4a\f3\3e\9c\9c\f0\3e\ec\f8\ed\3e\cc\5e\eb\3e\ef\cd\e8\3e\0f\46\e6\3e\e7\c6\e3\3e\37\50\e1\3e\c1\e1\de\3e\4b\7b\dc\3e\9d\1c\da\3e\82\c5\d7\3e\c7\75\d5\3e\3b\2d\d3\3e\b1\eb\d0\3e\fb\b0\ce\3e\f0\7c\cc\3e\65\4f\ca\3e\34\28\c8\3e\38\07\c6\3e\4c\ec\c3\3e\4e\d7\c1\3e\1b\c8\bf\3e\95\be\bd\3e\9c\ba\bb\3e\12\bc\b9\3e\da\c2\b7\3e\d9\ce\b5\3e\f4\df\b3\3e\12\f6\b1\3e\19\11\b0\3e\f1\30\ae\3e\83\55\ac\3e\b9\7e\aa\3e\7c\ac\a8\3e\b8\de\a6\3e\59\15\a5\3e\49\50\a3\3e\77\8f\a1\3e\d0\d2\9f\3e\42\1a\9e\3e\ba\65\9c\3e\29\b5\9a\3e\7e\08\99\3e\a9\5f\97\3e\9a\ba\95\3e\43\19\94\3e\94\7b\92\3e\80\e1\90\3e\f8\4a\8f\3e\ef\b7\8d\3e\58\28\8c\3e\27\9c\8a\3e\4e\13\89\3e\c3\8d\87\3e\78\0b\86\3e\62\8c\84\3e\78\10\83\3e\ac\97\81\3e\f5\21\80\3e\92\5e\7d\3e\3b\7f\7a\3e\d0\a5\77\3e\40\d2\74\3e\77\04\72\3e\62\3c\6f\3e\f1\79\6c\3e\11\bd\69\3e\b2\05\67\3e\c2\53\64\3e\33\a7\61\3e\f3\ff\5e\3e\f4\5d\5c\3e\26\c1\59\3e\7a\29\57\3e\e2\96\54\3e\50\09\52\3e\b7\80\4f\3e\07\fd\4c\3e\35\7e\4a\3e\33\04\48\3e\f5\8e\45\3e\6e\1e\43\3e\93\b2\40\3e\56\4b\3e\3e\ae\e8\3b\3e\8e\8a\39\3e\eb\30\37\3e\bb\db\34\3e\f3\8a\32\3e\88\3e\30\3e\70\f6\2d\3e\a2\b2\2b\3e\13\73\29\3e\bb\37\27\3e\8f\00\25\3e\86\cd\22\3e\98\9e\20\3e\bc\73\1e\3e\e9\4c\1c\3e\17\2a\1a\3e\3d\0b\18\3e\54\f0\15\3e\54\d9\13\3e\34\c6\11\3e\ed\b6\0f\3e\79\ab\0d\3e\cf\a3\0b\3e\e9\9f\09\3e\bf\9f\07\3e\4c\a3\05\3e\87\aa\03\3e\6c\b5\01\3e\e5\87\ff\3d\2b\ac\fb\3d\9d\d7\f7\3d\30\0a\f4\3d\d8\43\f0\3d\89\84\ec\3d\38\cc\e8\3d\db\1a\e5\3d\68\70\e1\3d\d3\cc\dd\3d\13\30\da\3d\1e\9a\d6\3d\ea\0a\d3\3d\6f\82\cf\3d\a2\00\cc\3d\7c\85\c8\3d\f4\10\c5\3d\01\a3\c1\3d\9c\3b\be\3d\bc\da\ba\3d\5a\80\b7\3d\6f\2c\b4\3d\f3\de\b0\3d\df\97\ad\3d\2e\57\aa\3d\d8\1c\a7\3d\d7\e8\a3\3d\25\bb\a0\3d\bd\93\9d\3d\99\72\9a\3d\b4\57\97\3d\09\43\94\3d\93\34\91\3d\4d\2c\8e\3d\34\2a\8b\3d\44\2e\88\3d\79\38\85\3d\cf\48\82\3d\86\be\7e\3d\a5\f7\78\3d\f5\3c\73\3d\72\8e\6d\3d\18\ec\67\3d\e3\55\62\3d\d1\cb\5c\3d\de\4d\57\3d\0a\dc\51\3d\54\76\4c\3d\bb\1c\47\3d\41\cf\41\3d\e6\8d\3c\3d\ac\58\37\3d\96\2f\32\3d\a9\12\2d\3d\e8\01\28\3d\59\fd\22\3d\02\05\1e\3d\ec\18\19\3d\1e\39\14\3d\a3\65\0f\3d\85\9e\0a\3d\d0\e3\05\3d\93\35\01\3d\b6\27\f9\3c\74\fd\ef\3c\83\ec\e6\3c\0b\f5\dd\3c\37\17\d5\3c\38\53\cc\3c\43\a9\c3\3c\8f\19\bb\3c\5c\a4\b2\3c\ed\49\aa\3c\8e\0a\a2\3c\91\e6\99\3c\4f\de\91\3c\2b\f2\89\3c\90\22\82\3c\ef\df\74\3c\c9\b5\65\3c\d3\c7\56\3c\53\17\48\3c\b7\a5\39\3c\98\74\2b\3c\c6\85\1d\3c\4f\db\0f\3c\91\77\02\3c\90\ba\ea\3b\4f\1f\d1\3b\fa\24\b8\3b\be\d4\9f\3b\eb\39\88\3b\9c\c5\62\3b\48\c4\36\3b\5d\a3\0c\3b\ab\5d\c9\3a\58\90\7d\3a\e2\18\ee\39\00\00\80\3f\78\ae\76\3f\39\b0\6f\3f\3a\bd\69\3f\92\6c\64\3f\db\8c\5f\3f\16\02\5b\3f\86\ba\56\3f\0c\aa\52\3f\f0\c7\4e\3f\a6\0d\4b\3f\1e\76\47\3f\53\fd\43\3f\03\a0\40\3f\7d\5b\3d\3f\82\2d\3a\3f\2b\14\37\3f\d6\0d\34\3f\18\19\31\3f\b7\34\2e\3f\9d\5f\2b\3f\d4\98\28\3f\84\df\25\3f\e7\32\23\3f\4f\92\20\3f\1c\fd\1d\3f\bd\72\1b\3f\b1\f2\18\3f\7d\7c\16\3f\b5\0f\14\3f\f3\ab\11\3f\d9\50\0f\3f\12\fe\0c\3f\4c\b3\0a\3f\3e\70\08\3f\a1\34\06\3f\36\00\04\3f\c0\d2\01\3f\0d\58\ff\3e\a7\17\fb\3e\ec\e3\f6\3e\7d\bc\f2\3e\02\a1\ee\3e\28\91\ea\3e\a0\8c\e6\3e\20\93\e2\3e\60\a4\de\3e\1e\c0\da\3e\1c\e6\d6\3e\1b\16\d3\3e\e5\4f\cf\3e\42\93\cb\3e\fe\df\c7\3e\e9\35\c4\3e\d4\94\c0\3e\92\fc\bc\3e\f8\6c\b9\3e\de\e5\b5\3e\1e\67\b2\3e\92\f0\ae\3e\16\82\ab\3e\8a\1b\a8\3e\cc\bc\a4\3e\bf\65\a1\3e\44\16\9e\3e\41\ce\9a\3e\98\8d\97\3e\32\54\94\3e\f5\21\91\3e\cb\f6\8d\3e\9c\d2\8a\3e\54\b5\87\3e\de\9e\84\3e\27\8f\81\3e\3a\0c\7d\3e\5b\07\77\3e\91\0f\71\3e\bc\24\6b\3e\c0\46\65\3e\7f\75\5f\3e\df\b0\59\3e\c8\f8\53\3e\22\4d\4e\3e\d7\ad\48\3e\d5\1a\43\3e\07\94\3d\3e\5e\19\38\3e\ca\aa\32\3e\3e\48\2d\3e\ad\f1\27\3e\0d\a7\22\3e\55\68\1d\3e\7f\35\18\3e\85\0e\13\3e\64\f3\0d\3e\1b\e4\08\3e\aa\e0\03\3e\26\d2\fd\3d\b7\fa\f3\3d\16\3b\ea\3d\56\93\e0\3d\8f\03\d7\3d\df\8b\cd\3d\6b\2c\c4\3d\5f\e5\ba\3d\eb\b6\b1\3d\4c\a1\a8\3d\c4\a4\9f\3d\a1\c1\96\3d\3b\f8\8d\3d\f8\48\85\3d\9a\68\79\3d\7b\75\68\3d\c6\b9\57\3d\da\36\47\3d\4f\ee\36\3d\fe\e1\26\3d\16\14\17\3d\2a\87\07\3d\a9\7c\f0\3c\bd\7a\d2\3c\fd\11\b5\3c\77\4e\98\3c\42\80\78\3c\5f\fa\41\3c\b6\4d\0d\3c\58\d4\b5\3b\f2\f4\2e\3b\00\00\00\00\00\00\00\00\39\a1\90\e2\00\00\00\00\bc\de\ea\9b\71\ac\77\c3\90\b9\dd\d4\b8\3f\89\de\7c\e8\a8\e4\6a\f1\df\e8\ab\de\f2\eb\e8\a6\49\ee\fd\4e\20\f0\8e\db\9b\f1\bb\58\d4\f2\4b\10\da\f3\78\6d\b8\f4\8a\ad\77\f5\3d\e8\1d\f6\84\b7\af\f6\73\a5\30\f7\51\76\a3\f7\b6\5b\0a\f8\9d\18\67\f8\4f\1b\bb\f8\62\90\07\f9\ca\70\4d\f9\7d\8c\8d\f9\8a\92\c8\f9\5b\17\ff\f9\96\99\31\fa\f8\85\60\fa\62\3a\8c\fa\4e\08\b5\fa\c8\36\db\fa\10\04\ff\fa\ea\a6\20\fb\b4\4f\40\fb\51\29\5e\fb\e9\59\7a\fb\8c\03\95\fb\ba\44\ae\fb\d8\38\c6\fb\92\f8\dc\fb\30\9a\f2\fb\df\31\07\fc\ed\d1\1a\fc\02\8b\2d\fc\4d\6c\3f\fc\ac\83\50\fc\d1\dd\60\fc\62\86\70\fc\10\88\7f\fc\b4\ec\8d\fc\62\bd\9b\fc\7c\02\a9\fc\c3\c3\b5\fc\64\08\c2\fc\0a\d7\cd\fc\e3\35\d9\fc\b0\2a\e4\fc\ce\ba\ee\fc\3b\eb\f8\fc\a0\c0\02\fd\59\3f\0c\fd\7b\6b\15\fd\d6\48\1e\fd\ff\da\26\fd\52\25\2f\fd\f7\2a\37\fd\e5\ee\3e\fd\e7\73\46\fd\9e\bc\4d\fd\85\cb\54\fd\f2\a2\5b\fd\1b\45\62\fd\15\b4\68\fd\da\f1\6e\fd\47\00\75\fd\20\e1\7a\fd\12\96\80\fd\b4\20\86\fd\85\82\8b\fd\f5\bc\90\fd\5e\d1\95\fd\0b\c1\9a\fd\36\8d\9f\fd\08\37\a4\fd\9e\bf\a8\fd\06\28\ad\fd\41\71\b1\fd\46\9c\b5\fd\fd\a9\b9\fd\46\9b\bd\fd\f6\70\c1\fd\d8\2b\c5\fd\ac\cc\c8\fd\2d\54\cc\fd\0b\c3\cf\fd\ef\19\d3\fd\7a\59\d6\fd\45\82\d9\fd\e5\94\dc\fd\e6\91\df\fd\ce\79\e2\fd\1f\4d\e5\fd\52\0c\e8\fd\de\b7\ea\fd\34\50\ed\fd\be\d5\ef\fd\e3\48\f2\fd\06\aa\f4\fd\84\f9\f6\fd\b6\37\f9\fd\f4\64\fb\fd\8d\81\fd\fd\d0\8d\ff\fd\08\8a\01\fe\7a\76\03\fe\6c\53\05\fe\1c\21\07\fe\c9\df\08\fe\ab\8f\0a\fe\fb\30\0c\fe\ec\c3\0d\fe\b1\48\0f\fe\76\bf\10\fe\69\28\12\fe\b4\83\13\fe\7c\d1\14\fe\e7\11\16\fe\16\45\17\fe\2a\6b\18\fe\3e\84\19\fe\70\90\1a\fe\d6\8f\1b\fe\89\82\1c\fe\9b\68\1d\fe\20\42\1e\fe\26\0f\1f\fe\bc\cf\1f\fe\ed\83\20\fe\c3\2b\21\fe\45\c7\21\fe\78\56\22\fe\5f\d9\22\fe\fb\4f\23\fe\4a\ba\23\fe\49\18\24\fe\f2\69\24\fe\3c\af\24\fe\1e\e8\24\fe\8b\14\25\fe\74\34\25\fe\c7\47\25\fe\70\4e\25\fe\5a\48\25\fe\6a\35\25\fe\86\15\25\fe\8f\e8\24\fe\64\ae\24\fe\e1\66\24\fe\df\11\24\fe\34\af\23\fe\b4\3e\23\fe\2c\c0\22\fe\6b\33\22\fe\38\98\21\fe\58\ee\20\fe\8c\35\20\fe\92\6d\1f\fe\21\96\1e\fe\f0\ae\1d\fe\ac\b7\1c\fe\02\b0\1b\fe\98\97\1a\fe\0d\6e\19\fe\fd\32\18\fe\fe\e5\16\fe\9d\86\15\fe\64\14\14\fe\d3\8e\12\fe\65\f5\10\fe\8c\47\0f\fe\b1\84\0d\fe\36\ac\0b\fe\73\bd\09\fe\b5\b7\07\fe\40\9a\05\fe\4c\64\03\fe\04\15\01\fe\88\ab\fe\fd\e9\26\fc\fd\29\86\f9\fd\3b\c8\f6\fd\01\ec\f3\fd\4a\f0\f0\fd\d1\d3\ed\fd\3d\95\ea\fd\1e\33\e7\fd\e9\ab\e3\fd\fb\fd\df\fd\91\27\dc\fd\cd\26\d8\fd\a8\f9\d3\fd\fc\9d\cf\fd\76\11\cb\fd\98\51\c6\fd\b3\5b\c1\fd\e2\2c\bc\fd\06\c2\b6\fd\be\17\b1\fd\63\2a\ab\fd\fd\f5\a4\fd\40\76\9e\fd\7a\a6\97\fd\92\81\90\fd\f2\01\89\fd\82\21\81\fd\8e\d9\78\fd\bb\22\70\fd\ed\f4\66\fd\32\47\5d\fd\9c\0f\53\fd\2b\43\48\fd\9a\d5\3c\fd\36\b9\30\fd\a4\de\23\fd\9e\34\16\fd\a3\a7\07\fd\9b\21\f8\fc\5b\89\e7\fc\20\c2\d5\fc\db\aa\c2\fc\5e\1d\ae\fc\4e\ed\97\fc\d4\e6\7f\fc\f3\cc\65\fc\62\57\49\fc\c8\2f\2a\fc\19\ee\07\fc\c1\13\e2\fb\1a\05\b8\fb\78\00\89\fb\a5\11\54\fb\05\00\18\fb\82\34\d3\fa\76\92\83\fa\32\3b\26\fa\1c\2d\b7\f9\a2\a1\30\f9\23\f0\89\f8\d2\77\b5\f7\0c\65\9c\f6\f0\30\15\f5\3c\0e\cb\f2\5d\b1\ef\ee\cf\6e\da\e6\12\22\ad\76\00\00\00\00\53\1b\0f\60\a6\47\e4\6c\a2\46\5b\72\1d\05\60\75\eb\21\49\77\bd\25\9a\78\c3\45\90\79\5d\ce\4b\7a\9f\62\df\7a\a6\82\56\7b\c6\a8\b8\7b\22\e7\0a\7c\e7\cc\50\7c\5b\ec\8c\7c\d6\2c\c1\7c\d2\fe\ee\7c\0b\7e\17\7d\83\88\3b\7d\6c\ce\5b\7d\64\dd\78\7d\86\28\93\7d\57\0e\ab\7d\30\dd\c0\7d\88\d6\d4\7d\85\31\e7\7d\ea\1c\f8\7d\a3\c0\07\7e\fa\3e\16\7e\87\b5\23\7e\fd\3d\30\7e\c2\ee\3b\7e\77\db\46\7e\5d\15\51\7e\b3\ab\5a\7e\f7\ab\63\7e\2c\22\6c\7e\06\19\74\7e\18\9a\7b\7e\fa\ad\82\7e\63\5c\89\7e\4b\ac\8f\7e\fb\a3\95\7e\24\49\9b\7e\ef\a0\a0\7e\0d\b0\a5\7e\c3\7a\aa\7e\f3\04\af\7e\2a\52\b3\7e\a5\65\b7\7e\59\42\bb\7e\fd\ea\be\7e\0a\62\c2\7e\c4\a9\c5\7e\41\c4\c8\7e\65\b3\cb\7e\ed\78\ce\7e\71\16\d1\7e\62\8d\d3\7e\12\df\d5\7e\b4\0c\d8\7e\5c\17\da\7e\05\00\dc\7e\8e\c7\dd\7e\bf\6e\df\7e\47\f6\e0\7e\be\5e\e2\7e\a9\a8\e3\7e\73\d4\e4\7e\76\e2\e5\7e\f5\d2\e6\7e\20\a6\e7\7e\10\5c\e8\7e\cd\f4\e8\7e\47\70\e9\7e\59\ce\e9\7e\ca\0e\ea\7e\47\31\ea\7e\68\35\ea\7e\ab\1a\ea\7e\71\e0\e9\7e\02\86\e9\7e\88\0a\e9\7e\08\6d\e8\7e\6a\ac\e7\7e\69\c7\e6\7e\9c\bc\e5\7e\67\8a\e4\7e\fc\2e\e3\7e\57\a8\e1\7e\2f\f4\df\7e\fa\0f\de\7e\d9\f8\db\7e\94\ab\d9\7e\8d\24\d7\7e\ae\5f\d4\7e\5c\58\d1\7e\5f\09\ce\7e\cb\6c\ca\7e\e2\7b\c6\7e\ee\2e\c2\7e\1a\7d\bd\7e\35\5c\b8\7e\75\c0\b2\7e\20\9c\ac\7e\27\df\a5\7e\9f\76\9e\7e\16\4c\96\7e\ba\44\8d\7e\33\40\83\7e\28\17\78\7e\33\99\6b\7e\1a\8a\5d\7e\ed\9d\4d\7e\7a\73\3b\7e\2f\8c\26\7e\f5\3f\0e\7e\5d\aa\f1\7d\72\8c\cf\7d\1e\1a\a6\7d\fb\a0\72\7d\97\e0\30\7d\ab\b4\d9\7c\1a\0f\60\7c\dc\0b\a9\7b\76\21\72\7a\e5\64\d6\77\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\e0\ff\ff\ff\ff\ff\ff\ff\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\1f\00\00\00\00\00\00\00\65\27\0b\31\ec\c4\82\2d\92\b5\d6\2d\ad\99\0c\2e\37\17\29\2e\88\84\42\2e\8f\c6\59\2e\c6\66\6f\2e\82\df\81\2e\d9\86\8b\2e\15\c0\94\2e\48\9c\9d\2e\ae\28\a6\2e\c5\6f\ae\2e\08\7a\b6\2e\6f\4e\be\2e\cb\f2\c5\2e\03\6c\cd\2e\46\be\d4\2e\2f\ed\db\2e\df\fb\e2\2e\12\ed\e9\2e\34\c3\f0\2e\66\80\f7\2e\93\26\fe\2e\b7\5b\02\2f\42\9a\05\2f\9c\cf\08\2f\67\fc\0b\2f\37\21\0f\2f\93\3e\12\2f\f7\54\15\2f\d5\64\18\2f\97\6e\1b\2f\9f\72\1e\2f\46\71\21\2f\e3\6a\24\2f\c3\5f\27\2f\31\50\2a\2f\72\3c\2d\2f\c6\24\30\2f\6b\09\33\2f\9b\ea\35\2f\8c\c8\38\2f\71\a3\3b\2f\7c\7b\3e\2f\db\50\41\2f\b9\23\44\2f\43\f4\46\2f\9e\c2\49\2f\f2\8e\4c\2f\64\59\4f\2f\16\22\52\2f\2b\e9\54\2f\c2\ae\57\2f\fc\72\5a\2f\f6\35\5d\2f\cd\f7\5f\2f\9d\b8\62\2f\81\78\65\2f\94\37\68\2f\f0\f5\6a\2f\ab\b3\6d\2f\e0\70\70\2f\a4\2d\73\2f\0f\ea\75\2f\37\a6\78\2f\31\62\7b\2f\13\1e\7e\2f\f9\6c\80\2f\f0\ca\81\2f\f9\28\83\2f\1e\87\84\2f\68\e5\85\2f\e1\43\87\2f\92\a2\88\2f\83\01\8a\2f\bf\60\8b\2f\4d\c0\8c\2f\37\20\8e\2f\85\80\8f\2f\3f\e1\90\2f\6e\42\92\2f\1a\a4\93\2f\4c\06\95\2f\0b\69\96\2f\61\cc\97\2f\54\30\99\2f\ed\94\9a\2f\35\fa\9b\2f\32\60\9d\2f\ee\c6\9e\2f\70\2e\a0\2f\bf\96\a1\2f\e5\ff\a2\2f\e8\69\a4\2f\d1\d4\a5\2f\a8\40\a7\2f\74\ad\a8\2f\3e\1b\aa\2f\0e\8a\ab\2f\eb\f9\ac\2f\df\6a\ae\2f\f0\dc\af\2f\27\50\b1\2f\8d\c4\b2\2f\29\3a\b4\2f\04\b1\b5\2f\26\29\b7\2f\99\a2\b8\2f\63\1d\ba\2f\8f\99\bb\2f\24\17\bd\2f\2b\96\be\2f\ae\16\c0\2f\b6\98\c1\2f\4b\1c\c3\2f\76\a1\c4\2f\42\28\c6\2f\b8\b0\c7\2f\e0\3a\c9\2f\c6\c6\ca\2f\72\54\cc\2f\ef\e3\cd\2f\47\75\cf\2f\84\08\d1\2f\b2\9d\d2\2f\da\34\d4\2f\07\ce\d5\2f\45\69\d7\2f\9f\06\d9\2f\20\a6\da\2f\d4\47\dc\2f\c7\eb\dd\2f\05\92\df\2f\9a\3a\e1\2f\94\e5\e2\2f\ff\92\e4\2f\e8\42\e6\2f\5c\f5\e7\2f\6a\aa\e9\2f\1f\62\eb\2f\8b\1c\ed\2f\ba\d9\ee\2f\be\99\f0\2f\a4\5c\f2\2f\7d\22\f4\2f\59\eb\f5\2f\48\b7\f7\2f\5b\86\f9\2f\a5\58\fb\2f\36\2e\fd\2f\20\07\ff\2f\bc\71\00\30\a7\61\01\30\5d\53\02\30\e6\46\03\30\4e\3c\04\30\a0\33\05\30\e5\2c\06\30\2b\28\07\30\7b\25\08\30\e3\24\09\30\6f\26\0a\30\2c\2a\0b\30\27\30\0c\30\6d\38\0d\30\0d\43\0e\30\15\50\0f\30\95\5f\10\30\9b\71\11\30\37\86\12\30\7b\9d\13\30\77\b7\14\30\3e\d4\15\30\e0\f3\16\30\73\16\18\30\08\3c\19\30\b6\64\1a\30\90\90\1b\30\ad\bf\1c\30\24\f2\1d\30\0c\28\1f\30\7f\61\20\30\96\9e\21\30\6c\df\22\30\1d\24\24\30\c5\6c\25\30\84\b9\26\30\78\0a\28\30\c4\5f\29\30\88\b9\2a\30\ea\17\2c\30\10\7b\2d\30\20\e3\2e\30\45\50\30\30\a9\c2\31\30\7b\3a\33\30\ea\b7\34\30\29\3b\36\30\6e\c4\37\30\ee\53\39\30\e7\e9\3a\30\96\86\3c\30\3c\2a\3e\30\1f\d5\3f\30\89\87\41\30\c8\41\43\30\2e\04\45\30\14\cf\46\30\d7\a2\48\30\da\7f\4a\30\88\66\4c\30\52\57\4e\30\b2\52\50\30\2a\59\52\30\46\6b\54\30\9c\89\56\30\ce\b4\58\30\8b\ed\5a\30\8f\34\5d\30\a7\8a\5f\30\b2\f0\61\30\a2\67\64\30\7f\f0\66\30\6b\8c\69\30\a3\3c\6c\30\85\02\6f\30\93\df\71\30\79\d5\74\30\13\e6\77\30\75\13\7b\30\f2\5f\7e\30\16\e7\80\30\8c\b0\82\30\0e\8e\84\30\8c\81\86\30\40\8d\88\30\c3\b3\8a\30\1e\f8\8c\30\e5\5d\8f\30\5e\e9\91\30\ad\9f\94\30\1c\87\97\30\71\a7\9a\30\76\0a\9e\30\bb\bc\a1\30\be\ce\a5\30\c2\56\aa\30\d7\73\af\30\3b\53\b5\30\87\3a\bc\30\ff\9c\c4\30\e0\4e\cf\30\f3\1c\de\30\c9\4e\f6\30\34\a3\ed\30\a4\6d\0b\2f\49\ca\39\2f\7f\64\5a\2f\bb\72\74\2f\b7\49\85\2f\73\06\8f\2f\6b\cc\97\2f\f7\d5\9f\2f\60\4a\a7\2f\7f\45\ae\2f\d9\db\b4\2f\0b\1d\bb\2f\43\15\c1\2f\35\ce\c6\2f\c3\4f\cc\2f\69\a0\d1\2f\93\c5\d6\2f\d4\c3\db\2f\16\9f\e0\2f\b5\5a\e5\2f\a0\f9\e9\2f\66\7e\ee\2f\48\eb\f2\2f\47\42\f7\2f\2c\85\fb\2f\91\b5\ff\2f\74\ea\01\30\3f\f2\03\30\c2\f2\05\30\89\ec\07\30\12\e0\09\30\d4\cd\0b\30\3a\b6\0d\30\aa\99\0f\30\80\78\11\30\14\53\13\30\b8\29\15\30\b8\fc\16\30\5d\cc\18\30\ec\98\1a\30\a4\62\1c\30\c3\29\1e\30\85\ee\1f\30\20\b1\21\30\cc\71\23\30\bb\30\25\30\1f\ee\26\30\28\aa\28\30\06\65\2a\30\e4\1e\2c\30\f1\d7\2d\30\55\90\2f\30\3d\48\31\30\d1\ff\32\30\3b\b7\34\30\a2\6e\36\30\2f\26\38\30\0a\de\39\30\5b\96\3b\30\48\4f\3d\30\fa\08\3f\30\98\c3\40\30\4a\7f\42\30\39\3c\44\30\8e\fa\45\30\71\ba\47\30\0d\7c\49\30\8d\3f\4b\30\1c\05\4d\30\e7\cc\4e\30\1c\97\50\30\ea\63\52\30\82\33\54\30\16\06\56\30\da\db\57\30\04\b5\59\30\cb\91\5b\30\6a\72\5d\30\1e\57\5f\30\27\40\61\30\c6\2d\63\30\42\20\65\30\e6\17\67\30\fd\14\69\30\db\17\6b\30\d5\20\6d\30\4a\30\6f\30\9a\46\71\30\2d\64\73\30\75\89\75\30\e7\b6\77\30\05\ed\79\30\58\2c\7c\30\76\75\7e\30\7f\64\80\30\d1\93\81\30\0f\c9\82\30\a0\04\84\30\f7\46\85\30\8f\90\86\30\f2\e1\87\30\b9\3b\89\30\8b\9e\8a\30\26\0b\8c\30\5c\82\8d\30\1a\05\8f\30\6f\94\90\30\8c\31\92\30\d1\dd\93\30\d5\9a\95\30\70\6a\97\30\cc\4e\99\30\7a\4a\9b\30\8e\60\9d\30\c2\94\9f\30\ae\eb\a1\30\0c\6b\a4\30\2b\1a\a7\30\90\02\aa\30\f7\30\ad\30\00\b7\b0\30\16\ae\b4\30\ef\3c\b9\30\f6\a2\be\30\9f\53\c5\30\06\47\ce\30\e2\53\dc\30\00\00\00\00\50\50\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\8c\00\00\00\14\00\00\00\04\00\00\00\0a\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\30\61\73\73\65\72\74\20\66\61\69\6c\65\64\20\28\61\73\73\65\72\74\20\66\61\69\6c\65\64\3a\20\2b\69\29\6e\69\6c\20\6d\61\70\2e\6d\61\70\2e\77\61\3a\36\38\3a\38\70\61\6e\69\63\3a\20\74\72\75\65\66\61\6c\73\65\4e\61\4e\2b\49\6e\66\2d\49\6e\66\30\31\32\33\34\35\36\37\38\39\61\62\63\64\65\66\0a\5b\2f\5d\72\61\6e\64\2e\47\65\74\52\61\6e\64\6f\6d\44\61\74\61\3a\20\75\6e\73\75\70\70\6f\72\74\69\6e\76\61\6c\69\64\20\61\72\67\75\6d\65\6e\74\20\74\6f\20\49\6e\74\33\31\6e\72\61\6e\64\2e\77\61\3a\31\31\33\3a\38\69\6e\76\61\6c\69\64\20\61\72\67\75\6d\65\6e\74\20\74\6f\20\49\6e\74\36\33\6e\72\61\6e\64\2e\77\61\3a\39\35\3a\38\69\6e\76\61\6c\69\64\20\61\72\67\75\6d\65\6e\74\20\74\6f\20\49\6e\74\6e\72\61\6e\64\2e\77\61\3a\31\35\36\3a\38\69\6e\76\61\6c\69\64\20\61\72\67\75\6d\65\6e\74\20\74\6f\20\53\68\75\66\66\6c\65\72\61\6e\64\2e\77\61\3a\32\33\30\3a\38\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\01\00\00\00\ff\ff\ff\ff\02\00\00\00\ff\ff\ff\ff\02\00\00\00\fe\ff\ff\ff\1e\00\00\00\1f\00\00\00\02\00\00\00\fd\ff\ff\ff\1e\00\00\00\1f\00\00\00\20\00\00\00\00\00\00\00\03\00\00\00\ff\ff\ff\ff\03\00\00\00\fc\ff\ff\ff\21\00\00\00\00\00\00\00\b8\73\00\00\00\00\00\00\00\00\00\00\00\00\00\00\c0\73\00\00\c8\73\00\00\d8\73\00\00\00\00\00\00\f0\73\00\00\00\00\00\00\00\00\00\00\f8\73\00\00")
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
	(elem (i32.const 23) $$math$rand.Source.underlying.$$OnFree)
	(elem (i32.const 24) $$math$rand.Source64.underlying.$$OnFree)
	(elem (i32.const 25) $$math$rand.Rand.$$OnFree)
	(elem (i32.const 26) $$errors.errorString.$$OnFree)
	(elem (i32.const 27) $math$rand.Rand.Intn.$bound)
	(elem (i32.const 28) $$math$rand.Rand.$$block.$$OnFree)
	(elem (i32.const 29) $$math$rand.Rand.$ref.underlying.$$OnFree)
	(elem (i32.const 30) $math$rand.rngSource.Int63)
	(elem (i32.const 31) $math$rand.rngSource.Seed)
	(elem (i32.const 32) $math$rand.rngSource.Uint64)
	(elem (i32.const 33) $errors.errorString.Error)
)
