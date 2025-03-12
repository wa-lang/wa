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

(func $runtime.Block.HeapAlloc (export "runtime.Block.HeapAlloc") (param $item_count i32) (param $release_func i32) (param $item_size i32) (result i32 i32) ;;result = ptr_block, ptr_data
  (local $b i32)
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
  
  local.tee $b
  local.get $b
  i32.const 16
  i32.add
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
