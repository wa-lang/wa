(module $__walang__
	(import "env" "traceUtf8" (func $runtime.traceUtf8 (param i32) (param i32)))
	(import "env" "line" (func $syscall$wasm4.__import__line (param i32) (param i32) (param i32) (param i32)))
	(import "env" "rect" (func $syscall$wasm4.__import__rect (param i32) (param i32) (param i32) (param i32)))
	(import "env" "textUtf8" (func $syscall$wasm4.__import__textUtf8 (param i32) (param i32) (param i32) (param i32)))
	(import "env" "tone" (func $syscall$wasm4.__import__tone (param i32) (param i32) (param i32) (param i32)))
	(import "env" "traceUtf8" (func $syscall$wasm4.__import__traceUtf8 (param i32) (param i32)))
	(import "env" "memory" (memory 1))
	(table 49 funcref)
	(type $$onFree (func (param i32)))
	(type $$wa.runtime.comp (func (param i32) (param i32) (result i32)))
	(type $$$fnSig1 (func))
	(type $$$fnSig2 (func (param i32) (param i32) (result i32 i32 i32)))
	(type $$$fnSig3 (func (param i32) (param i32)))
	(type $$$fnSig4 (func (param i32) (param i32) (param i32) (param i32) (result i32)))
	(type $$$fnSig5 (func (param i32) (param i32) (param i32)))
	(type $$$fnSig6 (func (param i32) (param i32) (result i32)))
	(type $$$fnSig7 (func (param i32) (param i32) (result i32)))
	(type $$$fnSig8 (func (result i32)))
	(global $__stack_ptr (mut i32) (i32.const 14656))
	(global $__heap_max i32 (i32.const 65536))
	(global $$wa.runtime.closure_data (mut i32) (i32.const 0))
	(global $$wa.runtime._concretTypeCount (mut i32) (i32.const 4))
	(global $$wa.runtime._interfaceCount (mut i32) (i32.const 3))
	(global $$wa.runtime._itabsPtr (mut i32) (i32.const 46280))
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
	(global $math$bits.deBruijn32tab.0 i32 (i32.const 0))
	(global $math$bits.deBruijn32tab.1 i32 (i32.const 14904))
	(global $math$bits.deBruijn64tab.0 i32 (i32.const 0))
	(global $math$bits.deBruijn64tab.1 i32 (i32.const 14936))
	(global $math$bits.init$guard (mut i32) (i32.const 0))
	(global $errors.init$guard (mut i32) (i32.const 0))
	(global $strconv.ErrRange.0 i32 (i32.const 0))
	(global $strconv.ErrRange.1 i32 (i32.const 15000))
	(global $strconv.ErrSyntax.0 i32 (i32.const 0))
	(global $strconv.ErrSyntax.1 i32 (i32.const 15016))
	(global $strconv.detailedPowersOfTen.0 i32 (i32.const 0))
	(global $strconv.detailedPowersOfTen.1 i32 (i32.const 15032))
	(global $strconv.f64info.0 i32 (i32.const 0))
	(global $strconv.f64info.1 i32 (i32.const 26168))
	(global $strconv.float32info.0 i32 (i32.const 0))
	(global $strconv.float32info.1 i32 (i32.const 26180))
	(global $strconv.float32pow10.0 i32 (i32.const 0))
	(global $strconv.float32pow10.1 i32 (i32.const 26192))
	(global $strconv.float64pow10.0 i32 (i32.const 0))
	(global $strconv.float64pow10.1 i32 (i32.const 26240))
	(global $strconv.init$guard (mut i32) (i32.const 0))
	(global $strconv.isGraphic.0 i32 (i32.const 0))
	(global $strconv.isGraphic.1 i32 (i32.const 26424))
	(global $strconv.isNotPrint16.0 i32 (i32.const 0))
	(global $strconv.isNotPrint16.1 i32 (i32.const 26456))
	(global $strconv.isNotPrint32.0 i32 (i32.const 0))
	(global $strconv.isNotPrint32.1 i32 (i32.const 26720))
	(global $strconv.isPrint16.0 i32 (i32.const 0))
	(global $strconv.isPrint16.1 i32 (i32.const 26910))
	(global $strconv.isPrint32.0 i32 (i32.const 0))
	(global $strconv.isPrint32.1 i32 (i32.const 27780))
	(global $strconv.leftcheats.0 i32 (i32.const 0))
	(global $strconv.leftcheats.1 i32 (i32.const 29652))
	(global $strconv.optimize.0 i32 (i32.const 0))
	(global $strconv.optimize.1 i32 (i32.const 30628))
	(global $strconv.powtab.0 i32 (i32.const 0))
	(global $strconv.powtab.1 i32 (i32.const 30632))
	(global $strconv.u64pow10.0 i32 (i32.const 0))
	(global $strconv.u64pow10.1 i32 (i32.const 30672))
	(global $syscall$wasm4.init$guard (mut i32) (i32.const 0))
	(global $unicode$utf8.acceptRanges.0 i32 (i32.const 0))
	(global $unicode$utf8.acceptRanges.1 i32 (i32.const 30832))
	(global $unicode$utf8.first.0 i32 (i32.const 0))
	(global $unicode$utf8.first.1 i32 (i32.const 30864))
	(global $unicode$utf8.init$guard (mut i32) (i32.const 0))
	(global $w42048.init$guard (mut i32) (i32.const 0))
	(global $w42048.sfxDown.0 i32 (i32.const 0))
	(global $w42048.sfxDown.1 i32 (i32.const 31120))
	(global $w42048.sfxLeft.0 i32 (i32.const 0))
	(global $w42048.sfxLeft.1 i32 (i32.const 31156))
	(global $w42048.sfxRight.0 i32 (i32.const 0))
	(global $w42048.sfxRight.1 i32 (i32.const 31192))
	(global $w42048.sfxUp.0 i32 (i32.const 0))
	(global $w42048.sfxUp.1 i32 (i32.const 31228))
	(global $w42048.ui.0 i32 (i32.const 0))
	(global $w42048.ui.1 i32 (i32.const 31264))
	(global $w42048$game.init$guard (mut i32) (i32.const 0))
	(global $w42048$palettes.All.0 i32 (i32.const 0))
	(global $w42048$palettes.All.1 i32 (i32.const 31272))
	(global $w42048$palettes.BluesGB.0 i32 (i32.const 0))
	(global $w42048$palettes.BluesGB.1 i32 (i32.const 31288))
	(global $w42048$palettes.EN4.0 i32 (i32.const 0))
	(global $w42048$palettes.EN4.1 i32 (i32.const 31304))
	(global $w42048$palettes.GBChocolate.0 i32 (i32.const 0))
	(global $w42048$palettes.GBChocolate.1 i32 (i32.const 31320))
	(global $w42048$palettes.Grapefruit.0 i32 (i32.const 0))
	(global $w42048$palettes.Grapefruit.1 i32 (i32.const 31336))
	(global $w42048$palettes.GreyMist.0 i32 (i32.const 0))
	(global $w42048$palettes.GreyMist.1 i32 (i32.const 31352))
	(global $w42048$palettes.IceCreamGB.0 i32 (i32.const 0))
	(global $w42048$palettes.IceCreamGB.1 i32 (i32.const 31368))
	(global $w42048$palettes.Keeby.0 i32 (i32.const 0))
	(global $w42048$palettes.Keeby.1 i32 (i32.const 31384))
	(global $w42048$palettes.Lightgreen.0 i32 (i32.const 0))
	(global $w42048$palettes.Lightgreen.1 i32 (i32.const 31400))
	(global $w42048$palettes.Platinum.0 i32 (i32.const 0))
	(global $w42048$palettes.Platinum.1 i32 (i32.const 31416))
	(global $w42048$palettes.Warmlight.0 i32 (i32.const 0))
	(global $w42048$palettes.Warmlight.1 i32 (i32.const 31432))
	(global $w42048$palettes.init$guard (mut i32) (i32.const 0))
	(global $runtime.zptr (mut i32) (i32.const 35048))
	(global $__heap_base i32 (i32.const 46352))
	(func $runtime.throw
		unreachable
	)
	(func $runtime.heapBase (result i32)
		global.get $__heap_base
	)
	(func $runtime.heapMax (result i32)
		global.get $__heap_max
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
						call_indirect 0 (type $$onFree)
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
	(func $$wa.runtime.string_to_ptr (param $b i32) (param $d i32) (param $l i32) (result i32)
		local.get $d
	)
	(func $$wa.runtime.string_to_iter (param $b i32) (param $d i32) (param $l i32) (result i32 i32 i32)
		local.get $d
		local.get $l
		i32.const 0
	)
	(func $$syscall/wasm4.__linkname__string_data_ptr (param $b i32) (param $d i32) (param $l i32) (result i32)
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
	(func $$u8.$$block.$$onFree (param $ptr i32)
		local.get $ptr
		i32.load align=1
		call $runtime.Block.Release
		local.get $ptr
		i32.const 0
		i32.store align=1
	)
	(func $$string.underlying.$$onFree (param $$ptr i32)
		local.get $$ptr
		i32.const 2
		call_indirect 0 (type $$onFree)
	)
	(func $runtime.free (param $ap i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$t0 i32)
		(local $$t1 i32)
		(local $$t2.0 i32)
		(local $$t2.1 i32)
		(local $$t3.0 i32)
		(local $$t3.1 i32)
		(local $$t4 i32)
		(local $$t5 i32)
		(local $$t6 i32)
		(local $$t7 i32)
		(local $$t8 i32)
		(local $$t9.0 i32)
		(local $$t9.1 i32)
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
		(local $$t24 i32)
		(local $$t25 i32)
		(local $$t26 i32)
		(local $$t27.0 i32)
		(local $$t27.1 i32)
		(local $$t28 i32)
		(local $$t29.0 i32)
		(local $$t29.1 i32)
		(local $$t30 i32)
		(local $$t31 i32)
		(local $$t32 i32)
		(local $$t33 i32)
		(local $$t34.0 i32)
		(local $$t34.1 i32)
		(local $$t35 i32)
		(local $$t36 i32)
		(local $$t37 i32)
		(local $$t38 i32)
		(local $$t39 i32)
		(local $$t40.0 i32)
		(local $$t40.1 i32)
		(local $$t41 i32)
		(local $$t42 i32)
		(local $$t43 i32)
		(local $$t44 i32)
		(local $$t45.0 i32)
		(local $$t45.1 i32)
		(local $$t46.0 i32)
		(local $$t46.1 i32)
		block $$BlockFnBody
			loop $$BlockDisp
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
																								br_table  0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 0
																							end
																							i32.const 0
																							local.set $$current_block
																							local.get $ap
																							i32.const 0
																							i32.eq
																							local.set $$t0
																							local.get $$t0
																							if
																								br $$Block_0
																							else
																								br $$Block_2
																							end
																						end
																						i32.const 1
																						local.set $$current_block
																						br $$BlockFnBody
																					end
																					i32.const 2
																					local.set $$current_block
																					local.get $ap
																					i32.const 8
																					i32.sub
																					local.set $$t1
																					i32.const 0
																					local.set $$t2.0
																					i32.const 0
																					local.set $$t2.1
																					local.get $$t1
																					call $runtime.knr_getBlockHeader
																					local.set $$t3.1
																					local.set $$t3.0
																					local.get $$t3.0
																					local.get $$t3.1
																					local.set $$t2.1
																					local.set $$t2.0
																					local.get $$t2.1
																					local.set $$t4
																					local.get $$t4
																					i32.const 1
																					i32.le_u
																					local.set $$t5
																					local.get $$t5
																					if
																						br $$Block_3
																					else
																						br $$Block_4
																					end
																				end
																				i32.const 3
																				local.set $$current_block
																				local.get $ap
																				i32.const 8
																				i32.rem_u
																				local.set $$t6
																				local.get $$t6
																				i32.const 0
																				i32.eq
																				i32.eqz
																				local.set $$t7
																				local.get $$t7
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
																			br $$BlockFnBody
																		end
																		i32.const 5
																		local.set $$current_block
																		global.get $$knr_freep
																		local.set $$t8
																		i32.const 0
																		local.set $$t9.0
																		i32.const 0
																		local.set $$t9.1
																		local.get $$t8
																		call $runtime.knr_getBlockHeader
																		local.set $$t10.1
																		local.set $$t10.0
																		local.get $$t10.0
																		local.get $$t10.1
																		local.set $$t9.1
																		local.set $$t9.0
																		br $$Block_7
																	end
																	i32.const 6
																	local.set $$current_block
																	local.get $$t9.0
																	local.set $$t11
																	local.get $$t12
																	local.get $$t11
																	i32.ge_u
																	local.set $$t13
																	local.get $$t13
																	if
																		br $$Block_10
																	else
																		br $$Block_9
																	end
																end
																i32.const 7
																local.set $$current_block
																local.get $$t2.1
																local.set $$t14
																local.get $$t14
																i32.const 8
																i32.mul
																local.set $$t15
																local.get $$t1
																local.get $$t15
																i32.add
																local.set $$t16
																local.get $$t9.0
																local.set $$t17
																local.get $$t16
																local.get $$t17
																i32.eq
																local.set $$t18
																local.get $$t18
																if
																	br $$Block_12
																else
																	br $$Block_14
																end
															end
															local.get $$current_block
															i32.const 5
															i32.eq
															if(result i32)
																local.get $$t8
															else
																local.get $$t19
															end
															local.set $$t12
															i32.const 8
															local.set $$current_block
															local.get $$t1
															local.get $$t12
															i32.gt_u
															local.set $$t20
															local.get $$t20
															if
																br $$Block_8
															else
																i32.const 6
																local.set $$block_selector
																br $$BlockDisp
															end
														end
														i32.const 9
														local.set $$current_block
														local.get $$t9.0
														local.set $$t21
														local.get $$t1
														local.get $$t21
														i32.lt_u
														local.set $$t22
														local.get $$t22
														if
															i32.const 7
															local.set $$block_selector
															br $$BlockDisp
														else
															i32.const 6
															local.set $$block_selector
															br $$BlockDisp
														end
													end
													i32.const 10
													local.set $$current_block
													local.get $$t9.0
													local.set $$t19
													local.get $$t19
													call $runtime.knr_getBlockHeader
													local.set $$t23.1
													local.set $$t23.0
													local.get $$t23.0
													local.get $$t23.1
													local.set $$t9.1
													local.set $$t9.0
													i32.const 8
													local.set $$block_selector
													br $$BlockDisp
												end
												i32.const 11
												local.set $$current_block
												local.get $$t1
												local.get $$t12
												i32.gt_u
												local.set $$t24
												local.get $$t24
												if
													i32.const 7
													local.set $$block_selector
													br $$BlockDisp
												else
													br $$Block_11
												end
											end
											i32.const 12
											local.set $$current_block
											local.get $$t9.0
											local.set $$t25
											local.get $$t1
											local.get $$t25
											i32.lt_u
											local.set $$t26
											local.get $$t26
											if
												i32.const 7
												local.set $$block_selector
												br $$BlockDisp
											else
												i32.const 10
												local.set $$block_selector
												br $$BlockDisp
											end
										end
										i32.const 13
										local.set $$current_block
										i32.const 0
										local.set $$t27.0
										i32.const 0
										local.set $$t27.1
										local.get $$t9.0
										local.set $$t28
										local.get $$t28
										call $runtime.knr_getBlockHeader
										local.set $$t29.1
										local.set $$t29.0
										local.get $$t29.0
										local.get $$t29.1
										local.set $$t27.1
										local.set $$t27.0
										local.get $$t27.1
										local.set $$t30
										local.get $$t2.1
										local.set $$t31
										local.get $$t31
										local.get $$t30
										i32.add
										local.set $$t32
										local.get $$t32
										local.set $$t2.1
										local.get $$t27.0
										local.set $$t33
										local.get $$t33
										local.set $$t2.0
										local.get $$t2.0
										local.get $$t2.1
										local.set $$t34.1
										local.set $$t34.0
										local.get $$t1
										local.get $$t34.0
										local.get $$t34.1
										call $runtime.knr_setBlockHeader
										br $$Block_13
									end
									i32.const 14
									local.set $$current_block
									local.get $$t9.1
									local.set $$t35
									local.get $$t35
									i32.const 8
									i32.mul
									local.set $$t36
									local.get $$t12
									local.get $$t36
									i32.add
									local.set $$t37
									local.get $$t37
									local.get $$t1
									i32.eq
									local.set $$t38
									local.get $$t38
									if
										br $$Block_15
									else
										br $$Block_17
									end
								end
								i32.const 15
								local.set $$current_block
								local.get $$t9.0
								local.set $$t39
								local.get $$t39
								local.set $$t2.0
								local.get $$t2.0
								local.get $$t2.1
								local.set $$t40.1
								local.set $$t40.0
								local.get $$t1
								local.get $$t40.0
								local.get $$t40.1
								call $runtime.knr_setBlockHeader
								i32.const 14
								local.set $$block_selector
								br $$BlockDisp
							end
							i32.const 16
							local.set $$current_block
							local.get $$t2.1
							local.set $$t41
							local.get $$t9.1
							local.set $$t42
							local.get $$t42
							local.get $$t41
							i32.add
							local.set $$t43
							local.get $$t43
							local.set $$t9.1
							local.get $$t2.0
							local.set $$t44
							local.get $$t44
							local.set $$t9.0
							local.get $$t9.0
							local.get $$t9.1
							local.set $$t45.1
							local.set $$t45.0
							local.get $$t12
							local.get $$t45.0
							local.get $$t45.1
							call $runtime.knr_setBlockHeader
							br $$Block_16
						end
						i32.const 17
						local.set $$current_block
						local.get $$t12
						global.set $$knr_freep
						br $$BlockFnBody
					end
					i32.const 18
					local.set $$current_block
					local.get $$t1
					local.set $$t9.0
					local.get $$t9.0
					local.get $$t9.1
					local.set $$t46.1
					local.set $$t46.0
					local.get $$t12
					local.get $$t46.0
					local.get $$t46.1
					call $runtime.knr_setBlockHeader
					i32.const 17
					local.set $$block_selector
					br $$BlockDisp
				end
			end
		end
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
	(func $runtime.knr_getBlockHeader (param $addr i32) (result i32 i32)
		local.get $addr
		i32.load
		local.get $addr
		i32.load offset=4
	)
	(func $runtime.knr_setBlockHeader (param $addr i32) (param $data.0 i32) (param $data.1 i32)
		local.get $addr
		local.get $data.0
		i32.store
		local.get $addr
		local.get $data.1
		i32.store offset=4
	)
	(func $$runtime.mapImp.$$block.$$onFree (param $ptr i32)
		local.get $ptr
		i32.load align=1
		call $runtime.Block.Release
		local.get $ptr
		i32.const 0
		i32.store align=1
	)
	(func $$runtime.mapImp.$ref.underlying.$$onFree (param $$ptr i32)
		local.get $$ptr
		i32.const 4
		call_indirect 0 (type $$onFree)
	)
	(func $$runtime.mapIter.$$onFree (param $$ptr i32)
		local.get $$ptr
		i32.const 5
		call_indirect 0 (type $$onFree)
	)
	(func $runtime.malloc (param $nbytes i32) (result i32)
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
		(local $$t9.0 i32)
		(local $$t9.1 i32)
		(local $$t10 i32)
		(local $$t11 i32)
		(local $$t12 i32)
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
		(local $$t29 i32)
		(local $$t30 i32)
		(local $$t31 i32)
		(local $$t32 i32)
		(local $$t33 i32)
		(local $$t34 i32)
		(local $$t35 i32)
		(local $$t36 i32)
		(local $$t37 i32)
		(local $$t38.0 i32)
		(local $$t38.1 i32)
		(local $$t39 i32)
		(local $$t40 i32)
		(local $$t41 i32)
		(local $$t42 i32)
		(local $$t43 i32)
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
		(local $$t50.0 i32)
		(local $$t50.1 i32)
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
																			br_table  0 1 2 3 4 5 6 7 8 9 10 11 12 13 0
																		end
																		i32.const 0
																		local.set $$current_block
																		local.get $nbytes
																		i32.const 0
																		i32.eq
																		local.set $$t0
																		local.get $$t0
																		if
																			br $$Block_0
																		else
																			br $$Block_2
																		end
																	end
																	i32.const 1
																	local.set $$current_block
																	i32.const 0
																	local.set $$ret_0
																	br $$BlockFnBody
																end
																i32.const 2
																local.set $$current_block
																global.get $$knr_basep
																local.set $$t1
																local.get $$t1
																i32.const 0
																i32.eq
																local.set $$t2
																local.get $$t2
																if
																	br $$Block_3
																else
																	br $$Block_4
																end
															end
															i32.const 3
															local.set $$current_block
															call $runtime.heapMax
															local.set $$t3
															global.get $$knr_basep
															local.set $$t4
															local.get $$t3
															local.get $$t4
															i32.sub
															local.set $$t5
															local.get $nbytes
															local.get $$t5
															i32.ge_u
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
														call $runtime.heapBase
														local.set $$t7
														local.get $$t7
														global.set $$knr_basep
														global.get $$knr_basep
														local.set $$t8
														local.get $$t8
														global.set $$knr_freep
														i32.const 0
														local.set $$t9.0
														i32.const 0
														local.set $$t9.1
														global.get $$knr_basep
														local.set $$t10
														call $runtime.heapMax
														local.set $$t11
														global.get $$knr_basep
														local.set $$t12
														local.get $$t11
														local.get $$t12
														i32.sub
														local.set $$t13
														local.get $$t13
														i32.const 8
														i32.div_u
														local.set $$t14
														local.get $$t14
														i32.const 1
														i32.sub
														local.set $$t15
														local.get $$t10
														local.set $$t9.0
														local.get $$t15
														local.set $$t9.1
														global.get $$knr_basep
														local.set $$t16
														local.get $$t9.0
														local.get $$t9.1
														local.set $$t17.1
														local.set $$t17.0
														local.get $$t16
														local.get $$t17.0
														local.get $$t17.1
														call $runtime.knr_setBlockHeader
														br $$Block_4
													end
													i32.const 5
													local.set $$current_block
													local.get $nbytes
													i32.const 8
													i32.add
													local.set $$t18
													local.get $$t18
													i32.const 1
													i32.sub
													local.set $$t19
													local.get $$t19
													i32.const 8
													i32.div_u
													local.set $$t20
													local.get $$t20
													i32.const 1
													i32.add
													local.set $$t21
													global.get $$knr_freep
													local.set $$t22
													i32.const 0
													local.set $$t23.0
													i32.const 0
													local.set $$t23.1
													local.get $$t22
													call $runtime.knr_getBlockHeader
													local.set $$t24.1
													local.set $$t24.0
													local.get $$t24.0
													local.get $$t24.1
													local.set $$t23.1
													local.set $$t23.0
													local.get $$t23.0
													local.set $$t25
													i32.const 0
													local.set $$t26.0
													i32.const 0
													local.set $$t26.1
													local.get $$t25
													call $runtime.knr_getBlockHeader
													local.set $$t27.1
													local.set $$t27.0
													local.get $$t27.0
													local.get $$t27.1
													local.set $$t26.1
													local.set $$t26.0
													br $$Block_5
												end
												local.get $$current_block
												i32.const 5
												i32.eq
												if(result i32)
													local.get $$t22
												else
													local.get $$t28
												end
												local.get $$current_block
												i32.const 5
												i32.eq
												if(result i32)
													local.get $$t25
												else
													local.get $$t30
												end
												local.set $$t28
												local.set $$t29
												i32.const 6
												local.set $$current_block
												local.get $$t26.1
												local.set $$t31
												local.get $$t31
												local.get $$t21
												i32.ge_u
												local.set $$t32
												local.get $$t32
												if
													br $$Block_6
												else
													br $$Block_7
												end
											end
											i32.const 7
											local.set $$current_block
											local.get $$t26.1
											local.set $$t33
											local.get $$t33
											local.get $$t21
											i32.eq
											local.set $$t34
											local.get $$t34
											if
												br $$Block_8
											else
												br $$Block_10
											end
										end
										i32.const 8
										local.set $$current_block
										global.get $$knr_freep
										local.set $$t35
										local.get $$t28
										local.get $$t35
										i32.eq
										local.set $$t36
										local.get $$t36
										if
											br $$Block_11
										else
											br $$Block_12
										end
									end
									i32.const 9
									local.set $$current_block
									local.get $$t26.0
									local.set $$t37
									local.get $$t37
									local.set $$t23.0
									local.get $$t23.0
									local.get $$t23.1
									local.set $$t38.1
									local.set $$t38.0
									local.get $$t29
									local.get $$t38.0
									local.get $$t38.1
									call $runtime.knr_setBlockHeader
									br $$Block_9
								end
								local.get $$current_block
								i32.const 9
								i32.eq
								if(result i32)
									local.get $$t28
								else
									local.get $$t39
								end
								local.set $$t40
								i32.const 10
								local.set $$current_block
								local.get $$t29
								global.set $$knr_freep
								local.get $$t40
								i32.const 8
								i32.add
								local.set $$t41
								local.get $$t41
								local.set $$ret_0
								br $$BlockFnBody
							end
							i32.const 11
							local.set $$current_block
							local.get $$t26.1
							local.set $$t42
							local.get $$t42
							local.get $$t21
							i32.sub
							local.set $$t43
							local.get $$t43
							local.set $$t26.1
							local.get $$t26.0
							local.get $$t26.1
							local.set $$t44.1
							local.set $$t44.0
							local.get $$t28
							local.get $$t44.0
							local.get $$t44.1
							call $runtime.knr_setBlockHeader
							local.get $$t26.1
							local.set $$t45
							local.get $$t45
							i32.const 8
							i32.mul
							local.set $$t46
							local.get $$t28
							local.get $$t46
							i32.add
							local.set $$t39
							local.get $$t39
							call $runtime.knr_getBlockHeader
							local.set $$t47.1
							local.set $$t47.0
							local.get $$t47.0
							local.get $$t47.1
							local.set $$t26.1
							local.set $$t26.0
							local.get $$t21
							local.set $$t26.1
							local.get $$t26.0
							local.get $$t26.1
							local.set $$t48.1
							local.set $$t48.0
							local.get $$t39
							local.get $$t48.0
							local.get $$t48.1
							call $runtime.knr_setBlockHeader
							i32.const 10
							local.set $$block_selector
							br $$BlockDisp
						end
						i32.const 12
						local.set $$current_block
						i32.const 0
						local.set $$ret_0
						br $$BlockFnBody
					end
					i32.const 13
					local.set $$current_block
					local.get $$t28
					call $runtime.knr_getBlockHeader
					local.set $$t49.1
					local.set $$t49.0
					local.get $$t49.0
					local.get $$t49.1
					local.set $$t23.1
					local.set $$t23.0
					local.get $$t26.0
					local.set $$t30
					local.get $$t30
					call $runtime.knr_getBlockHeader
					local.set $$t50.1
					local.set $$t50.0
					local.get $$t50.0
					local.get $$t50.1
					local.set $$t26.1
					local.set $$t26.0
					i32.const 6
					local.set $$block_selector
					br $$BlockDisp
				end
			end
		end
		local.get $$ret_0
	)
	(func $$runtime.mapNode.$$block.$$onFree (param $ptr i32)
		local.get $ptr
		i32.load align=1
		call $runtime.Block.Release
		local.get $ptr
		i32.const 0
		i32.store align=1
	)
	(func $$runtime.mapNode.$ref.underlying.$$onFree (param $$ptr i32)
		local.get $$ptr
		i32.const 7
		call_indirect 0 (type $$onFree)
	)
	(func $$void.$$block.$$onFree (param $ptr i32)
		local.get $ptr
		i32.load align=1
		call $runtime.Block.Release
		local.get $ptr
		i32.const 0
		i32.store align=1
	)
	(func $$void.$ref.underlying.$$onFree (param $$ptr i32)
		local.get $$ptr
		i32.const 9
		call_indirect 0 (type $$onFree)
	)
	(func $$i`0`.underlying.$$onFree (param $$ptr i32)
		local.get $$ptr
		i32.const 10
		call_indirect 0 (type $$onFree)
	)
	(func $$runtime.mapNode.$$onFree (param $$ptr i32)
		local.get $$ptr
		i32.const 8
		i32.add
		i32.const 8
		call_indirect 0 (type $$onFree)
		local.get $$ptr
		i32.const 16
		i32.add
		i32.const 8
		call_indirect 0 (type $$onFree)
		local.get $$ptr
		i32.const 28
		i32.add
		i32.const 11
		call_indirect 0 (type $$onFree)
		local.get $$ptr
		i32.const 44
		i32.add
		i32.const 11
		call_indirect 0 (type $$onFree)
	)
	(func $$runtime.mapNode.$ref.$$block.$$onFree (param $ptr i32)
		local.get $ptr
		i32.load align=1
		call $runtime.Block.Release
		local.get $ptr
		i32.const 0
		i32.store align=1
	)
	(func $$runtime.mapNode.$ref.$slice.underlying.$$onFree (param $$ptr i32)
		local.get $$ptr
		i32.const 13
		call_indirect 0 (type $$onFree)
	)
	(func $$runtime.mapImp.$$onFree (param $$ptr i32)
		local.get $$ptr
		i32.const 8
		call_indirect 0 (type $$onFree)
		local.get $$ptr
		i32.const 8
		i32.add
		i32.const 8
		call_indirect 0 (type $$onFree)
		local.get $$ptr
		i32.const 16
		i32.add
		i32.const 14
		call_indirect 0 (type $$onFree)
	)
	(func $$runtime.mapNode.$ref.$array1.underlying.$$onFree (param $$ptr i32)
		local.get $$ptr
		i32.const 8
		call_indirect 0 (type $$onFree)
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
					i32.const 31501
					i32.const 7
					call $$runtime.waPrintString
					local.get $msg_ptr
					local.get $msg_len
					call $$runtime.waPuts
					i32.const 0
					i32.const 31462
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
	(func $$$$$$.underlying.$$onFree (param $$ptr i32)
		local.get $$ptr
		i32.const 4
		i32.add
		i32.const 10
		call_indirect 0 (type $$onFree)
	)
	(func $$$$$$.$array1.underlying.$$onFree (param $$ptr i32)
		local.get $$ptr
		i32.const 17
		call_indirect 0 (type $$onFree)
	)
	(func $$$$$$.$$block.$$onFree (param $ptr i32)
		local.get $ptr
		i32.load align=1
		call $runtime.Block.Release
		local.get $ptr
		i32.const 0
		i32.store align=1
	)
	(func $$$$$$.$slice.underlying.$$onFree (param $$ptr i32)
		local.get $$ptr
		i32.const 19
		call_indirect 0 (type $$onFree)
	)
	(func $$runtime.defers.$$onFree (param $$ptr i32)
		local.get $$ptr
		i32.const 20
		call_indirect 0 (type $$onFree)
	)
	(func $$runtime.defers.$array1.underlying.$$onFree (param $$ptr i32)
		local.get $$ptr
		i32.const 20
		call_indirect 0 (type $$onFree)
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
	(func $math$bits.TrailingZeros (param $x i32) (result i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$ret_0 i32)
		(local $$t0 i32)
		(local $$t1 i32)
		(local $$t2 i64)
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
							i32.const 1
							if
								br $$Block_0
							else
								br $$Block_1
							end
						end
						i32.const 1
						local.set $$current_block
						local.get $x
						local.set $$t0
						local.get $$t0
						call $math$bits.TrailingZeros32
						local.set $$t1
						local.get $$t1
						local.set $$ret_0
						br $$BlockFnBody
					end
					i32.const 2
					local.set $$current_block
					local.get $x
					i64.extend_i32_u
					local.set $$t2
					local.get $$t2
					call $math$bits.TrailingZeros64
					local.set $$t3
					local.get $$t3
					local.set $$ret_0
					br $$BlockFnBody
				end
			end
		end
		local.get $$ret_0
	)
	(func $math$bits.TrailingZeros32 (param $x i32) (result i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$ret_0 i32)
		(local $$t0 i32)
		(local $$t1 i32)
		(local $$t2 i32)
		(local $$t3 i32)
		(local $$t4 i32)
		(local $$t5 i32)
		(local $$t6.0 i32)
		(local $$t6.1 i32)
		(local $$t7 i32)
		(local $$t8 i32)
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
							i32.const 0
							i32.eq
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
						i32.const 32
						local.set $$ret_0
						br $$BlockFnBody
					end
					i32.const 2
					local.set $$current_block
					i32.const 0
					local.get $x
					i32.sub
					local.set $$t1
					local.get $x
					local.get $$t1
					i32.and
					local.set $$t2
					local.get $$t2
					i32.const 125613361
					i32.mul
					local.set $$t3
					local.get $$t3
					i64.const 27
					i32.wrap_i64
					i32.shr_u
					local.set $$t4
					local.get $$t4
					local.set $$t5
					i32.const 0
					i32.const 14904
					i32.const 1
					local.get $$t5
					i32.mul
					i32.add
					local.set $$t6.1
					local.get $$t6.0
					call $runtime.Block.Release
					local.set $$t6.0
					local.get $$t6.1
					i32.load8_u align=1
					local.set $$t7
					local.get $$t7
					local.set $$t8
					local.get $$t8
					local.set $$ret_0
					br $$BlockFnBody
				end
			end
		end
		local.get $$ret_0
		local.get $$t6.0
		call $runtime.Block.Release
	)
	(func $math$bits.TrailingZeros64 (param $x i64) (result i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$ret_0 i32)
		(local $$t0 i32)
		(local $$t1 i64)
		(local $$t2 i64)
		(local $$t3 i64)
		(local $$t4 i64)
		(local $$t5 i32)
		(local $$t6.0 i32)
		(local $$t6.1 i32)
		(local $$t7 i32)
		(local $$t8 i32)
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
							i64.const 0
							i64.eq
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
						i32.const 64
						local.set $$ret_0
						br $$BlockFnBody
					end
					i32.const 2
					local.set $$current_block
					i64.const 0
					local.get $x
					i64.sub
					local.set $$t1
					local.get $x
					local.get $$t1
					i64.and
					local.set $$t2
					local.get $$t2
					i64.const 285870213051353865
					i64.mul
					local.set $$t3
					local.get $$t3
					i64.const 58
					i64.shr_u
					local.set $$t4
					local.get $$t4
					i32.wrap_i64
					local.set $$t5
					i32.const 0
					i32.const 14936
					i32.const 1
					local.get $$t5
					i32.mul
					i32.add
					local.set $$t6.1
					local.get $$t6.0
					call $runtime.Block.Release
					local.set $$t6.0
					local.get $$t6.1
					i32.load8_u align=1
					local.set $$t7
					local.get $$t7
					local.set $$t8
					local.get $$t8
					local.set $$ret_0
					br $$BlockFnBody
				end
			end
		end
		local.get $$ret_0
		local.get $$t6.0
		call $runtime.Block.Release
	)
	(func $math$bits.init
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
							global.get $math$bits.init$guard
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
						global.set $math$bits.init$guard
						br $$Block_1
					end
					i32.const 2
					local.set $$current_block
					br $$BlockFnBody
				end
			end
		end
	)
	(func $$errors.errorString.$$onFree (param $$ptr i32)
		local.get $$ptr
		i32.const 3
		call_indirect 0 (type $$onFree)
	)
	(func $errors.New (param $text.0 i32) (param $text.1 i32) (param $text.2 i32) (result i32 i32 i32 i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$ret_0.0.0 i32)
		(local $$ret_0.0.1 i32)
		(local $$ret_0.1 i32)
		(local $$ret_0.2 i32)
		(local $$t0.0 i32)
		(local $$t0.1 i32)
		(local $$t1.0 i32)
		(local $$t1.1 i32)
		(local $$t2.0.0 i32)
		(local $$t2.0.1 i32)
		(local $$t2.1 i32)
		(local $$t2.2 i32)
		block $$BlockFnBody
			loop $$BlockDisp
				block $$Block_0
					block $$BlockSel
						local.get $$block_selector
						br_table  0 0
					end
					i32.const 0
					local.set $$current_block
					i32.const 28
					call $runtime.HeapAlloc
					i32.const 1
					i32.const 23
					i32.const 12
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
					local.get $$t1.1
					local.get $text.0
					call $runtime.Block.Retain
					local.get $$t1.1
					i32.load align=1
					call $runtime.Block.Release
					i32.store align=1
					local.get $$t1.1
					local.get $text.1
					i32.store offset=4
					local.get $$t1.1
					local.get $text.2
					i32.store offset=8
					local.get $$t0.0
					call $runtime.Block.Retain
					local.get $$t0.1
					i32.const 2
					i32.const -2
					i32.const 0
					call $runtime.getItab
					i32.const 0
					local.set $$t2.2
					local.set $$t2.1
					local.set $$t2.0.1
					local.get $$t2.0.0
					call $runtime.Block.Release
					local.set $$t2.0.0
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
		local.get $$t1.0
		call $runtime.Block.Release
		local.get $$t2.0.0
		call $runtime.Block.Release
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
		if(result i32 i32 i32 i32)
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
					br $loop1
				end
			end
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
		end
	)
	(func $$.error.underlying.$$onFree (param $$ptr i32)
		local.get $$ptr
		i32.const 10
		call_indirect 0 (type $$onFree)
	)
	(func $$strconv.NumError.$$onFree (param $$ptr i32)
		local.get $$ptr
		i32.const 3
		call_indirect 0 (type $$onFree)
		local.get $$ptr
		i32.const 12
		i32.add
		i32.const 3
		call_indirect 0 (type $$onFree)
		local.get $$ptr
		i32.const 24
		i32.add
		i32.const 24
		call_indirect 0 (type $$onFree)
	)
	(func $strconv.FormatInt (param $i i64) (param $base i32) (result i32 i32 i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$ret_0.0 i32)
		(local $$ret_0.1 i32)
		(local $$ret_0.2 i32)
		(local $$t0 i32)
		(local $$t1.0 i32)
		(local $$t1.1 i32)
		(local $$t1.2 i32)
		(local $$t2 i64)
		(local $$t3 i32)
		(local $$t4.0.0 i32)
		(local $$t4.0.1 i32)
		(local $$t4.0.2 i32)
		(local $$t4.0.3 i32)
		(local $$t4.1.0 i32)
		(local $$t4.1.1 i32)
		(local $$t4.1.2 i32)
		(local $$t5.0 i32)
		(local $$t5.1 i32)
		(local $$t5.2 i32)
		(local $$t5.3 i32)
		(local $$t6.0 i32)
		(local $$t6.1 i32)
		(local $$t6.2 i32)
		(local $$t7 i32)
		(local $$t8 i32)
		(local $$t9 i32)
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
										i32.const 1
										if
											br $$Block_4
										else
											br $$Block_1
										end
									end
									i32.const 1
									local.set $$current_block
									local.get $i
									i32.wrap_i64
									local.set $$t0
									local.get $$t0
									call $strconv.small
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
								i32.const 2
								local.set $$current_block
								local.get $i
								local.set $$t2
								local.get $i
								i64.const 0
								i64.lt_s
								local.set $$t3
								i32.const 0
								i32.const 0
								i32.const 0
								i32.const 0
								local.get $$t2
								local.get $base
								local.get $$t3
								i32.const 0
								call $strconv.formatBits
								local.set $$t4.1.2
								local.set $$t4.1.1
								local.get $$t4.1.0
								call $runtime.Block.Release
								local.set $$t4.1.0
								local.set $$t4.0.3
								local.set $$t4.0.2
								local.set $$t4.0.1
								local.get $$t4.0.0
								call $runtime.Block.Release
								local.set $$t4.0.0
								local.get $$t4.0.0
								call $runtime.Block.Retain
								local.get $$t4.0.1
								local.get $$t4.0.2
								local.get $$t4.0.3
								local.set $$t5.3
								local.set $$t5.2
								local.set $$t5.1
								local.get $$t5.0
								call $runtime.Block.Release
								local.set $$t5.0
								local.get $$t4.1.0
								call $runtime.Block.Retain
								local.get $$t4.1.1
								local.get $$t4.1.2
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
							i32.const 3
							local.set $$current_block
							local.get $base
							i32.const 10
							i32.eq
							local.set $$t7
							local.get $$t7
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
						local.get $i
						i64.const 100
						i64.lt_s
						local.set $$t8
						local.get $$t8
						if
							i32.const 3
							local.set $$block_selector
							br $$BlockDisp
						else
							i32.const 2
							local.set $$block_selector
							br $$BlockDisp
						end
					end
					i32.const 5
					local.set $$current_block
					i64.const 0
					local.get $i
					i64.le_s
					local.set $$t9
					local.get $$t9
					if
						i32.const 4
						local.set $$block_selector
						br $$BlockDisp
					else
						i32.const 2
						local.set $$block_selector
						br $$BlockDisp
					end
				end
			end
		end
		local.get $$ret_0.0
		call $runtime.Block.Retain
		local.get $$ret_0.1
		local.get $$ret_0.2
		local.get $$ret_0.0
		call $runtime.Block.Release
		local.get $$t1.0
		call $runtime.Block.Release
		local.get $$t4.1.0
		call $runtime.Block.Release
		local.get $$t4.0.0
		call $runtime.Block.Release
		local.get $$t5.0
		call $runtime.Block.Release
		local.get $$t6.0
		call $runtime.Block.Release
	)
	(func $strconv.IsPrint (param $r i32) (result i32)
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
		(local $$t7.0 i32)
		(local $$t7.1 i32)
		(local $$t8.0 i32)
		(local $$t8.1 i32)
		(local $$t9 i32)
		(local $$t10.0 i32)
		(local $$t10.1 i32)
		(local $$t10.2 i32)
		(local $$t10.3 i32)
		(local $$t10.4 i32)
		(local $$t10.5 i32)
		(local $$t10.6 i32)
		(local $$t10.7 i32)
		(local $$t10.8 i32)
		(local $$t10.9 i32)
		(local $$t10.10 i32)
		(local $$t10.11 i32)
		(local $$t10.12 i32)
		(local $$t10.13 i32)
		(local $$t10.14 i32)
		(local $$t10.15 i32)
		(local $$t10.16 i32)
		(local $$t10.17 i32)
		(local $$t10.18 i32)
		(local $$t10.19 i32)
		(local $$t10.20 i32)
		(local $$t10.21 i32)
		(local $$t10.22 i32)
		(local $$t10.23 i32)
		(local $$t10.24 i32)
		(local $$t10.25 i32)
		(local $$t10.26 i32)
		(local $$t10.27 i32)
		(local $$t10.28 i32)
		(local $$t10.29 i32)
		(local $$t10.30 i32)
		(local $$t10.31 i32)
		(local $$t10.32 i32)
		(local $$t10.33 i32)
		(local $$t10.34 i32)
		(local $$t10.35 i32)
		(local $$t10.36 i32)
		(local $$t10.37 i32)
		(local $$t10.38 i32)
		(local $$t10.39 i32)
		(local $$t10.40 i32)
		(local $$t10.41 i32)
		(local $$t10.42 i32)
		(local $$t10.43 i32)
		(local $$t10.44 i32)
		(local $$t10.45 i32)
		(local $$t10.46 i32)
		(local $$t10.47 i32)
		(local $$t10.48 i32)
		(local $$t10.49 i32)
		(local $$t10.50 i32)
		(local $$t10.51 i32)
		(local $$t10.52 i32)
		(local $$t10.53 i32)
		(local $$t10.54 i32)
		(local $$t10.55 i32)
		(local $$t10.56 i32)
		(local $$t10.57 i32)
		(local $$t10.58 i32)
		(local $$t10.59 i32)
		(local $$t10.60 i32)
		(local $$t10.61 i32)
		(local $$t10.62 i32)
		(local $$t10.63 i32)
		(local $$t10.64 i32)
		(local $$t10.65 i32)
		(local $$t10.66 i32)
		(local $$t10.67 i32)
		(local $$t10.68 i32)
		(local $$t10.69 i32)
		(local $$t10.70 i32)
		(local $$t10.71 i32)
		(local $$t10.72 i32)
		(local $$t10.73 i32)
		(local $$t10.74 i32)
		(local $$t10.75 i32)
		(local $$t10.76 i32)
		(local $$t10.77 i32)
		(local $$t10.78 i32)
		(local $$t10.79 i32)
		(local $$t10.80 i32)
		(local $$t10.81 i32)
		(local $$t10.82 i32)
		(local $$t10.83 i32)
		(local $$t10.84 i32)
		(local $$t10.85 i32)
		(local $$t10.86 i32)
		(local $$t10.87 i32)
		(local $$t10.88 i32)
		(local $$t10.89 i32)
		(local $$t10.90 i32)
		(local $$t10.91 i32)
		(local $$t10.92 i32)
		(local $$t10.93 i32)
		(local $$t10.94 i32)
		(local $$t10.95 i32)
		(local $$t10.96 i32)
		(local $$t10.97 i32)
		(local $$t10.98 i32)
		(local $$t10.99 i32)
		(local $$t10.100 i32)
		(local $$t10.101 i32)
		(local $$t10.102 i32)
		(local $$t10.103 i32)
		(local $$t10.104 i32)
		(local $$t10.105 i32)
		(local $$t10.106 i32)
		(local $$t10.107 i32)
		(local $$t10.108 i32)
		(local $$t10.109 i32)
		(local $$t10.110 i32)
		(local $$t10.111 i32)
		(local $$t10.112 i32)
		(local $$t10.113 i32)
		(local $$t10.114 i32)
		(local $$t10.115 i32)
		(local $$t10.116 i32)
		(local $$t10.117 i32)
		(local $$t10.118 i32)
		(local $$t10.119 i32)
		(local $$t10.120 i32)
		(local $$t10.121 i32)
		(local $$t10.122 i32)
		(local $$t10.123 i32)
		(local $$t10.124 i32)
		(local $$t10.125 i32)
		(local $$t10.126 i32)
		(local $$t10.127 i32)
		(local $$t10.128 i32)
		(local $$t10.129 i32)
		(local $$t10.130 i32)
		(local $$t10.131 i32)
		(local $$t10.132 i32)
		(local $$t10.133 i32)
		(local $$t10.134 i32)
		(local $$t10.135 i32)
		(local $$t10.136 i32)
		(local $$t10.137 i32)
		(local $$t10.138 i32)
		(local $$t10.139 i32)
		(local $$t10.140 i32)
		(local $$t10.141 i32)
		(local $$t10.142 i32)
		(local $$t10.143 i32)
		(local $$t10.144 i32)
		(local $$t10.145 i32)
		(local $$t10.146 i32)
		(local $$t10.147 i32)
		(local $$t10.148 i32)
		(local $$t10.149 i32)
		(local $$t10.150 i32)
		(local $$t10.151 i32)
		(local $$t10.152 i32)
		(local $$t10.153 i32)
		(local $$t10.154 i32)
		(local $$t10.155 i32)
		(local $$t10.156 i32)
		(local $$t10.157 i32)
		(local $$t10.158 i32)
		(local $$t10.159 i32)
		(local $$t10.160 i32)
		(local $$t10.161 i32)
		(local $$t10.162 i32)
		(local $$t10.163 i32)
		(local $$t10.164 i32)
		(local $$t10.165 i32)
		(local $$t10.166 i32)
		(local $$t10.167 i32)
		(local $$t10.168 i32)
		(local $$t10.169 i32)
		(local $$t10.170 i32)
		(local $$t10.171 i32)
		(local $$t10.172 i32)
		(local $$t10.173 i32)
		(local $$t10.174 i32)
		(local $$t10.175 i32)
		(local $$t10.176 i32)
		(local $$t10.177 i32)
		(local $$t10.178 i32)
		(local $$t10.179 i32)
		(local $$t10.180 i32)
		(local $$t10.181 i32)
		(local $$t10.182 i32)
		(local $$t10.183 i32)
		(local $$t10.184 i32)
		(local $$t10.185 i32)
		(local $$t10.186 i32)
		(local $$t10.187 i32)
		(local $$t10.188 i32)
		(local $$t10.189 i32)
		(local $$t10.190 i32)
		(local $$t10.191 i32)
		(local $$t10.192 i32)
		(local $$t10.193 i32)
		(local $$t10.194 i32)
		(local $$t10.195 i32)
		(local $$t10.196 i32)
		(local $$t10.197 i32)
		(local $$t10.198 i32)
		(local $$t10.199 i32)
		(local $$t10.200 i32)
		(local $$t10.201 i32)
		(local $$t10.202 i32)
		(local $$t10.203 i32)
		(local $$t10.204 i32)
		(local $$t10.205 i32)
		(local $$t10.206 i32)
		(local $$t10.207 i32)
		(local $$t10.208 i32)
		(local $$t10.209 i32)
		(local $$t10.210 i32)
		(local $$t10.211 i32)
		(local $$t10.212 i32)
		(local $$t10.213 i32)
		(local $$t10.214 i32)
		(local $$t10.215 i32)
		(local $$t10.216 i32)
		(local $$t10.217 i32)
		(local $$t10.218 i32)
		(local $$t10.219 i32)
		(local $$t10.220 i32)
		(local $$t10.221 i32)
		(local $$t10.222 i32)
		(local $$t10.223 i32)
		(local $$t10.224 i32)
		(local $$t10.225 i32)
		(local $$t10.226 i32)
		(local $$t10.227 i32)
		(local $$t10.228 i32)
		(local $$t10.229 i32)
		(local $$t10.230 i32)
		(local $$t10.231 i32)
		(local $$t10.232 i32)
		(local $$t10.233 i32)
		(local $$t10.234 i32)
		(local $$t10.235 i32)
		(local $$t10.236 i32)
		(local $$t10.237 i32)
		(local $$t10.238 i32)
		(local $$t10.239 i32)
		(local $$t10.240 i32)
		(local $$t10.241 i32)
		(local $$t10.242 i32)
		(local $$t10.243 i32)
		(local $$t10.244 i32)
		(local $$t10.245 i32)
		(local $$t10.246 i32)
		(local $$t10.247 i32)
		(local $$t10.248 i32)
		(local $$t10.249 i32)
		(local $$t10.250 i32)
		(local $$t10.251 i32)
		(local $$t10.252 i32)
		(local $$t10.253 i32)
		(local $$t10.254 i32)
		(local $$t10.255 i32)
		(local $$t10.256 i32)
		(local $$t10.257 i32)
		(local $$t10.258 i32)
		(local $$t10.259 i32)
		(local $$t10.260 i32)
		(local $$t10.261 i32)
		(local $$t10.262 i32)
		(local $$t10.263 i32)
		(local $$t10.264 i32)
		(local $$t10.265 i32)
		(local $$t10.266 i32)
		(local $$t10.267 i32)
		(local $$t10.268 i32)
		(local $$t10.269 i32)
		(local $$t10.270 i32)
		(local $$t10.271 i32)
		(local $$t10.272 i32)
		(local $$t10.273 i32)
		(local $$t10.274 i32)
		(local $$t10.275 i32)
		(local $$t10.276 i32)
		(local $$t10.277 i32)
		(local $$t10.278 i32)
		(local $$t10.279 i32)
		(local $$t10.280 i32)
		(local $$t10.281 i32)
		(local $$t10.282 i32)
		(local $$t10.283 i32)
		(local $$t10.284 i32)
		(local $$t10.285 i32)
		(local $$t10.286 i32)
		(local $$t10.287 i32)
		(local $$t10.288 i32)
		(local $$t10.289 i32)
		(local $$t10.290 i32)
		(local $$t10.291 i32)
		(local $$t10.292 i32)
		(local $$t10.293 i32)
		(local $$t10.294 i32)
		(local $$t10.295 i32)
		(local $$t10.296 i32)
		(local $$t10.297 i32)
		(local $$t10.298 i32)
		(local $$t10.299 i32)
		(local $$t10.300 i32)
		(local $$t10.301 i32)
		(local $$t10.302 i32)
		(local $$t10.303 i32)
		(local $$t10.304 i32)
		(local $$t10.305 i32)
		(local $$t10.306 i32)
		(local $$t10.307 i32)
		(local $$t10.308 i32)
		(local $$t10.309 i32)
		(local $$t10.310 i32)
		(local $$t10.311 i32)
		(local $$t10.312 i32)
		(local $$t10.313 i32)
		(local $$t10.314 i32)
		(local $$t10.315 i32)
		(local $$t10.316 i32)
		(local $$t10.317 i32)
		(local $$t10.318 i32)
		(local $$t10.319 i32)
		(local $$t10.320 i32)
		(local $$t10.321 i32)
		(local $$t10.322 i32)
		(local $$t10.323 i32)
		(local $$t10.324 i32)
		(local $$t10.325 i32)
		(local $$t10.326 i32)
		(local $$t10.327 i32)
		(local $$t10.328 i32)
		(local $$t10.329 i32)
		(local $$t10.330 i32)
		(local $$t10.331 i32)
		(local $$t10.332 i32)
		(local $$t10.333 i32)
		(local $$t10.334 i32)
		(local $$t10.335 i32)
		(local $$t10.336 i32)
		(local $$t10.337 i32)
		(local $$t10.338 i32)
		(local $$t10.339 i32)
		(local $$t10.340 i32)
		(local $$t10.341 i32)
		(local $$t10.342 i32)
		(local $$t10.343 i32)
		(local $$t10.344 i32)
		(local $$t10.345 i32)
		(local $$t10.346 i32)
		(local $$t10.347 i32)
		(local $$t10.348 i32)
		(local $$t10.349 i32)
		(local $$t10.350 i32)
		(local $$t10.351 i32)
		(local $$t10.352 i32)
		(local $$t10.353 i32)
		(local $$t10.354 i32)
		(local $$t10.355 i32)
		(local $$t10.356 i32)
		(local $$t10.357 i32)
		(local $$t10.358 i32)
		(local $$t10.359 i32)
		(local $$t10.360 i32)
		(local $$t10.361 i32)
		(local $$t10.362 i32)
		(local $$t10.363 i32)
		(local $$t10.364 i32)
		(local $$t10.365 i32)
		(local $$t10.366 i32)
		(local $$t10.367 i32)
		(local $$t10.368 i32)
		(local $$t10.369 i32)
		(local $$t10.370 i32)
		(local $$t10.371 i32)
		(local $$t10.372 i32)
		(local $$t10.373 i32)
		(local $$t10.374 i32)
		(local $$t10.375 i32)
		(local $$t10.376 i32)
		(local $$t10.377 i32)
		(local $$t10.378 i32)
		(local $$t10.379 i32)
		(local $$t10.380 i32)
		(local $$t10.381 i32)
		(local $$t10.382 i32)
		(local $$t10.383 i32)
		(local $$t10.384 i32)
		(local $$t10.385 i32)
		(local $$t10.386 i32)
		(local $$t10.387 i32)
		(local $$t10.388 i32)
		(local $$t10.389 i32)
		(local $$t10.390 i32)
		(local $$t10.391 i32)
		(local $$t10.392 i32)
		(local $$t10.393 i32)
		(local $$t10.394 i32)
		(local $$t10.395 i32)
		(local $$t10.396 i32)
		(local $$t10.397 i32)
		(local $$t10.398 i32)
		(local $$t10.399 i32)
		(local $$t10.400 i32)
		(local $$t10.401 i32)
		(local $$t10.402 i32)
		(local $$t10.403 i32)
		(local $$t10.404 i32)
		(local $$t10.405 i32)
		(local $$t10.406 i32)
		(local $$t10.407 i32)
		(local $$t10.408 i32)
		(local $$t10.409 i32)
		(local $$t10.410 i32)
		(local $$t10.411 i32)
		(local $$t10.412 i32)
		(local $$t10.413 i32)
		(local $$t10.414 i32)
		(local $$t10.415 i32)
		(local $$t10.416 i32)
		(local $$t10.417 i32)
		(local $$t10.418 i32)
		(local $$t10.419 i32)
		(local $$t10.420 i32)
		(local $$t10.421 i32)
		(local $$t10.422 i32)
		(local $$t10.423 i32)
		(local $$t10.424 i32)
		(local $$t10.425 i32)
		(local $$t10.426 i32)
		(local $$t10.427 i32)
		(local $$t10.428 i32)
		(local $$t10.429 i32)
		(local $$t10.430 i32)
		(local $$t10.431 i32)
		(local $$t10.432 i32)
		(local $$t10.433 i32)
		(local $$t11.0 i32)
		(local $$t11.1 i32)
		(local $$t11.2 i32)
		(local $$t11.3 i32)
		(local $$t11.4 i32)
		(local $$t11.5 i32)
		(local $$t11.6 i32)
		(local $$t11.7 i32)
		(local $$t11.8 i32)
		(local $$t11.9 i32)
		(local $$t11.10 i32)
		(local $$t11.11 i32)
		(local $$t11.12 i32)
		(local $$t11.13 i32)
		(local $$t11.14 i32)
		(local $$t11.15 i32)
		(local $$t11.16 i32)
		(local $$t11.17 i32)
		(local $$t11.18 i32)
		(local $$t11.19 i32)
		(local $$t11.20 i32)
		(local $$t11.21 i32)
		(local $$t11.22 i32)
		(local $$t11.23 i32)
		(local $$t11.24 i32)
		(local $$t11.25 i32)
		(local $$t11.26 i32)
		(local $$t11.27 i32)
		(local $$t11.28 i32)
		(local $$t11.29 i32)
		(local $$t11.30 i32)
		(local $$t11.31 i32)
		(local $$t11.32 i32)
		(local $$t11.33 i32)
		(local $$t11.34 i32)
		(local $$t11.35 i32)
		(local $$t11.36 i32)
		(local $$t11.37 i32)
		(local $$t11.38 i32)
		(local $$t11.39 i32)
		(local $$t11.40 i32)
		(local $$t11.41 i32)
		(local $$t11.42 i32)
		(local $$t11.43 i32)
		(local $$t11.44 i32)
		(local $$t11.45 i32)
		(local $$t11.46 i32)
		(local $$t11.47 i32)
		(local $$t11.48 i32)
		(local $$t11.49 i32)
		(local $$t11.50 i32)
		(local $$t11.51 i32)
		(local $$t11.52 i32)
		(local $$t11.53 i32)
		(local $$t11.54 i32)
		(local $$t11.55 i32)
		(local $$t11.56 i32)
		(local $$t11.57 i32)
		(local $$t11.58 i32)
		(local $$t11.59 i32)
		(local $$t11.60 i32)
		(local $$t11.61 i32)
		(local $$t11.62 i32)
		(local $$t11.63 i32)
		(local $$t11.64 i32)
		(local $$t11.65 i32)
		(local $$t11.66 i32)
		(local $$t11.67 i32)
		(local $$t11.68 i32)
		(local $$t11.69 i32)
		(local $$t11.70 i32)
		(local $$t11.71 i32)
		(local $$t11.72 i32)
		(local $$t11.73 i32)
		(local $$t11.74 i32)
		(local $$t11.75 i32)
		(local $$t11.76 i32)
		(local $$t11.77 i32)
		(local $$t11.78 i32)
		(local $$t11.79 i32)
		(local $$t11.80 i32)
		(local $$t11.81 i32)
		(local $$t11.82 i32)
		(local $$t11.83 i32)
		(local $$t11.84 i32)
		(local $$t11.85 i32)
		(local $$t11.86 i32)
		(local $$t11.87 i32)
		(local $$t11.88 i32)
		(local $$t11.89 i32)
		(local $$t11.90 i32)
		(local $$t11.91 i32)
		(local $$t11.92 i32)
		(local $$t11.93 i32)
		(local $$t11.94 i32)
		(local $$t11.95 i32)
		(local $$t11.96 i32)
		(local $$t11.97 i32)
		(local $$t11.98 i32)
		(local $$t11.99 i32)
		(local $$t11.100 i32)
		(local $$t11.101 i32)
		(local $$t11.102 i32)
		(local $$t11.103 i32)
		(local $$t11.104 i32)
		(local $$t11.105 i32)
		(local $$t11.106 i32)
		(local $$t11.107 i32)
		(local $$t11.108 i32)
		(local $$t11.109 i32)
		(local $$t11.110 i32)
		(local $$t11.111 i32)
		(local $$t11.112 i32)
		(local $$t11.113 i32)
		(local $$t11.114 i32)
		(local $$t11.115 i32)
		(local $$t11.116 i32)
		(local $$t11.117 i32)
		(local $$t11.118 i32)
		(local $$t11.119 i32)
		(local $$t11.120 i32)
		(local $$t11.121 i32)
		(local $$t11.122 i32)
		(local $$t11.123 i32)
		(local $$t11.124 i32)
		(local $$t11.125 i32)
		(local $$t11.126 i32)
		(local $$t11.127 i32)
		(local $$t11.128 i32)
		(local $$t11.129 i32)
		(local $$t11.130 i32)
		(local $$t11.131 i32)
		(local $$t12.0 i32)
		(local $$t12.1 i32)
		(local $$t12.2 i32)
		(local $$t12.3 i32)
		(local $$t13 i32)
		(local $$t14 i32)
		(local $$t15.0 i32)
		(local $$t15.1 i32)
		(local $$t16.0 i32)
		(local $$t16.1 i32)
		(local $$t17 i32)
		(local $$t18.0 i32)
		(local $$t18.1 i32)
		(local $$t18.2 i32)
		(local $$t18.3 i32)
		(local $$t18.4 i32)
		(local $$t18.5 i32)
		(local $$t18.6 i32)
		(local $$t18.7 i32)
		(local $$t18.8 i32)
		(local $$t18.9 i32)
		(local $$t18.10 i32)
		(local $$t18.11 i32)
		(local $$t18.12 i32)
		(local $$t18.13 i32)
		(local $$t18.14 i32)
		(local $$t18.15 i32)
		(local $$t18.16 i32)
		(local $$t18.17 i32)
		(local $$t18.18 i32)
		(local $$t18.19 i32)
		(local $$t18.20 i32)
		(local $$t18.21 i32)
		(local $$t18.22 i32)
		(local $$t18.23 i32)
		(local $$t18.24 i32)
		(local $$t18.25 i32)
		(local $$t18.26 i32)
		(local $$t18.27 i32)
		(local $$t18.28 i32)
		(local $$t18.29 i32)
		(local $$t18.30 i32)
		(local $$t18.31 i32)
		(local $$t18.32 i32)
		(local $$t18.33 i32)
		(local $$t18.34 i32)
		(local $$t18.35 i32)
		(local $$t18.36 i32)
		(local $$t18.37 i32)
		(local $$t18.38 i32)
		(local $$t18.39 i32)
		(local $$t18.40 i32)
		(local $$t18.41 i32)
		(local $$t18.42 i32)
		(local $$t18.43 i32)
		(local $$t18.44 i32)
		(local $$t18.45 i32)
		(local $$t18.46 i32)
		(local $$t18.47 i32)
		(local $$t18.48 i32)
		(local $$t18.49 i32)
		(local $$t18.50 i32)
		(local $$t18.51 i32)
		(local $$t18.52 i32)
		(local $$t18.53 i32)
		(local $$t18.54 i32)
		(local $$t18.55 i32)
		(local $$t18.56 i32)
		(local $$t18.57 i32)
		(local $$t18.58 i32)
		(local $$t18.59 i32)
		(local $$t18.60 i32)
		(local $$t18.61 i32)
		(local $$t18.62 i32)
		(local $$t18.63 i32)
		(local $$t18.64 i32)
		(local $$t18.65 i32)
		(local $$t18.66 i32)
		(local $$t18.67 i32)
		(local $$t18.68 i32)
		(local $$t18.69 i32)
		(local $$t18.70 i32)
		(local $$t18.71 i32)
		(local $$t18.72 i32)
		(local $$t18.73 i32)
		(local $$t18.74 i32)
		(local $$t18.75 i32)
		(local $$t18.76 i32)
		(local $$t18.77 i32)
		(local $$t18.78 i32)
		(local $$t18.79 i32)
		(local $$t18.80 i32)
		(local $$t18.81 i32)
		(local $$t18.82 i32)
		(local $$t18.83 i32)
		(local $$t18.84 i32)
		(local $$t18.85 i32)
		(local $$t18.86 i32)
		(local $$t18.87 i32)
		(local $$t18.88 i32)
		(local $$t18.89 i32)
		(local $$t18.90 i32)
		(local $$t18.91 i32)
		(local $$t18.92 i32)
		(local $$t18.93 i32)
		(local $$t18.94 i32)
		(local $$t18.95 i32)
		(local $$t18.96 i32)
		(local $$t18.97 i32)
		(local $$t18.98 i32)
		(local $$t18.99 i32)
		(local $$t18.100 i32)
		(local $$t18.101 i32)
		(local $$t18.102 i32)
		(local $$t18.103 i32)
		(local $$t18.104 i32)
		(local $$t18.105 i32)
		(local $$t18.106 i32)
		(local $$t18.107 i32)
		(local $$t18.108 i32)
		(local $$t18.109 i32)
		(local $$t18.110 i32)
		(local $$t18.111 i32)
		(local $$t18.112 i32)
		(local $$t18.113 i32)
		(local $$t18.114 i32)
		(local $$t18.115 i32)
		(local $$t18.116 i32)
		(local $$t18.117 i32)
		(local $$t18.118 i32)
		(local $$t18.119 i32)
		(local $$t18.120 i32)
		(local $$t18.121 i32)
		(local $$t18.122 i32)
		(local $$t18.123 i32)
		(local $$t18.124 i32)
		(local $$t18.125 i32)
		(local $$t18.126 i32)
		(local $$t18.127 i32)
		(local $$t18.128 i32)
		(local $$t18.129 i32)
		(local $$t18.130 i32)
		(local $$t18.131 i32)
		(local $$t18.132 i32)
		(local $$t18.133 i32)
		(local $$t18.134 i32)
		(local $$t18.135 i32)
		(local $$t18.136 i32)
		(local $$t18.137 i32)
		(local $$t18.138 i32)
		(local $$t18.139 i32)
		(local $$t18.140 i32)
		(local $$t18.141 i32)
		(local $$t18.142 i32)
		(local $$t18.143 i32)
		(local $$t18.144 i32)
		(local $$t18.145 i32)
		(local $$t18.146 i32)
		(local $$t18.147 i32)
		(local $$t18.148 i32)
		(local $$t18.149 i32)
		(local $$t18.150 i32)
		(local $$t18.151 i32)
		(local $$t18.152 i32)
		(local $$t18.153 i32)
		(local $$t18.154 i32)
		(local $$t18.155 i32)
		(local $$t18.156 i32)
		(local $$t18.157 i32)
		(local $$t18.158 i32)
		(local $$t18.159 i32)
		(local $$t18.160 i32)
		(local $$t18.161 i32)
		(local $$t18.162 i32)
		(local $$t18.163 i32)
		(local $$t18.164 i32)
		(local $$t18.165 i32)
		(local $$t18.166 i32)
		(local $$t18.167 i32)
		(local $$t18.168 i32)
		(local $$t18.169 i32)
		(local $$t18.170 i32)
		(local $$t18.171 i32)
		(local $$t18.172 i32)
		(local $$t18.173 i32)
		(local $$t18.174 i32)
		(local $$t18.175 i32)
		(local $$t18.176 i32)
		(local $$t18.177 i32)
		(local $$t18.178 i32)
		(local $$t18.179 i32)
		(local $$t18.180 i32)
		(local $$t18.181 i32)
		(local $$t18.182 i32)
		(local $$t18.183 i32)
		(local $$t18.184 i32)
		(local $$t18.185 i32)
		(local $$t18.186 i32)
		(local $$t18.187 i32)
		(local $$t18.188 i32)
		(local $$t18.189 i32)
		(local $$t18.190 i32)
		(local $$t18.191 i32)
		(local $$t18.192 i32)
		(local $$t18.193 i32)
		(local $$t18.194 i32)
		(local $$t18.195 i32)
		(local $$t18.196 i32)
		(local $$t18.197 i32)
		(local $$t18.198 i32)
		(local $$t18.199 i32)
		(local $$t18.200 i32)
		(local $$t18.201 i32)
		(local $$t18.202 i32)
		(local $$t18.203 i32)
		(local $$t18.204 i32)
		(local $$t18.205 i32)
		(local $$t18.206 i32)
		(local $$t18.207 i32)
		(local $$t18.208 i32)
		(local $$t18.209 i32)
		(local $$t18.210 i32)
		(local $$t18.211 i32)
		(local $$t18.212 i32)
		(local $$t18.213 i32)
		(local $$t18.214 i32)
		(local $$t18.215 i32)
		(local $$t18.216 i32)
		(local $$t18.217 i32)
		(local $$t18.218 i32)
		(local $$t18.219 i32)
		(local $$t18.220 i32)
		(local $$t18.221 i32)
		(local $$t18.222 i32)
		(local $$t18.223 i32)
		(local $$t18.224 i32)
		(local $$t18.225 i32)
		(local $$t18.226 i32)
		(local $$t18.227 i32)
		(local $$t18.228 i32)
		(local $$t18.229 i32)
		(local $$t18.230 i32)
		(local $$t18.231 i32)
		(local $$t18.232 i32)
		(local $$t18.233 i32)
		(local $$t18.234 i32)
		(local $$t18.235 i32)
		(local $$t18.236 i32)
		(local $$t18.237 i32)
		(local $$t18.238 i32)
		(local $$t18.239 i32)
		(local $$t18.240 i32)
		(local $$t18.241 i32)
		(local $$t18.242 i32)
		(local $$t18.243 i32)
		(local $$t18.244 i32)
		(local $$t18.245 i32)
		(local $$t18.246 i32)
		(local $$t18.247 i32)
		(local $$t18.248 i32)
		(local $$t18.249 i32)
		(local $$t18.250 i32)
		(local $$t18.251 i32)
		(local $$t18.252 i32)
		(local $$t18.253 i32)
		(local $$t18.254 i32)
		(local $$t18.255 i32)
		(local $$t18.256 i32)
		(local $$t18.257 i32)
		(local $$t18.258 i32)
		(local $$t18.259 i32)
		(local $$t18.260 i32)
		(local $$t18.261 i32)
		(local $$t18.262 i32)
		(local $$t18.263 i32)
		(local $$t18.264 i32)
		(local $$t18.265 i32)
		(local $$t18.266 i32)
		(local $$t18.267 i32)
		(local $$t18.268 i32)
		(local $$t18.269 i32)
		(local $$t18.270 i32)
		(local $$t18.271 i32)
		(local $$t18.272 i32)
		(local $$t18.273 i32)
		(local $$t18.274 i32)
		(local $$t18.275 i32)
		(local $$t18.276 i32)
		(local $$t18.277 i32)
		(local $$t18.278 i32)
		(local $$t18.279 i32)
		(local $$t18.280 i32)
		(local $$t18.281 i32)
		(local $$t18.282 i32)
		(local $$t18.283 i32)
		(local $$t18.284 i32)
		(local $$t18.285 i32)
		(local $$t18.286 i32)
		(local $$t18.287 i32)
		(local $$t18.288 i32)
		(local $$t18.289 i32)
		(local $$t18.290 i32)
		(local $$t18.291 i32)
		(local $$t18.292 i32)
		(local $$t18.293 i32)
		(local $$t18.294 i32)
		(local $$t18.295 i32)
		(local $$t18.296 i32)
		(local $$t18.297 i32)
		(local $$t18.298 i32)
		(local $$t18.299 i32)
		(local $$t18.300 i32)
		(local $$t18.301 i32)
		(local $$t18.302 i32)
		(local $$t18.303 i32)
		(local $$t18.304 i32)
		(local $$t18.305 i32)
		(local $$t18.306 i32)
		(local $$t18.307 i32)
		(local $$t18.308 i32)
		(local $$t18.309 i32)
		(local $$t18.310 i32)
		(local $$t18.311 i32)
		(local $$t18.312 i32)
		(local $$t18.313 i32)
		(local $$t18.314 i32)
		(local $$t18.315 i32)
		(local $$t18.316 i32)
		(local $$t18.317 i32)
		(local $$t18.318 i32)
		(local $$t18.319 i32)
		(local $$t18.320 i32)
		(local $$t18.321 i32)
		(local $$t18.322 i32)
		(local $$t18.323 i32)
		(local $$t18.324 i32)
		(local $$t18.325 i32)
		(local $$t18.326 i32)
		(local $$t18.327 i32)
		(local $$t18.328 i32)
		(local $$t18.329 i32)
		(local $$t18.330 i32)
		(local $$t18.331 i32)
		(local $$t18.332 i32)
		(local $$t18.333 i32)
		(local $$t18.334 i32)
		(local $$t18.335 i32)
		(local $$t18.336 i32)
		(local $$t18.337 i32)
		(local $$t18.338 i32)
		(local $$t18.339 i32)
		(local $$t18.340 i32)
		(local $$t18.341 i32)
		(local $$t18.342 i32)
		(local $$t18.343 i32)
		(local $$t18.344 i32)
		(local $$t18.345 i32)
		(local $$t18.346 i32)
		(local $$t18.347 i32)
		(local $$t18.348 i32)
		(local $$t18.349 i32)
		(local $$t18.350 i32)
		(local $$t18.351 i32)
		(local $$t18.352 i32)
		(local $$t18.353 i32)
		(local $$t18.354 i32)
		(local $$t18.355 i32)
		(local $$t18.356 i32)
		(local $$t18.357 i32)
		(local $$t18.358 i32)
		(local $$t18.359 i32)
		(local $$t18.360 i32)
		(local $$t18.361 i32)
		(local $$t18.362 i32)
		(local $$t18.363 i32)
		(local $$t18.364 i32)
		(local $$t18.365 i32)
		(local $$t18.366 i32)
		(local $$t18.367 i32)
		(local $$t18.368 i32)
		(local $$t18.369 i32)
		(local $$t18.370 i32)
		(local $$t18.371 i32)
		(local $$t18.372 i32)
		(local $$t18.373 i32)
		(local $$t18.374 i32)
		(local $$t18.375 i32)
		(local $$t18.376 i32)
		(local $$t18.377 i32)
		(local $$t18.378 i32)
		(local $$t18.379 i32)
		(local $$t18.380 i32)
		(local $$t18.381 i32)
		(local $$t18.382 i32)
		(local $$t18.383 i32)
		(local $$t18.384 i32)
		(local $$t18.385 i32)
		(local $$t18.386 i32)
		(local $$t18.387 i32)
		(local $$t18.388 i32)
		(local $$t18.389 i32)
		(local $$t18.390 i32)
		(local $$t18.391 i32)
		(local $$t18.392 i32)
		(local $$t18.393 i32)
		(local $$t18.394 i32)
		(local $$t18.395 i32)
		(local $$t18.396 i32)
		(local $$t18.397 i32)
		(local $$t18.398 i32)
		(local $$t18.399 i32)
		(local $$t18.400 i32)
		(local $$t18.401 i32)
		(local $$t18.402 i32)
		(local $$t18.403 i32)
		(local $$t18.404 i32)
		(local $$t18.405 i32)
		(local $$t18.406 i32)
		(local $$t18.407 i32)
		(local $$t18.408 i32)
		(local $$t18.409 i32)
		(local $$t18.410 i32)
		(local $$t18.411 i32)
		(local $$t18.412 i32)
		(local $$t18.413 i32)
		(local $$t18.414 i32)
		(local $$t18.415 i32)
		(local $$t18.416 i32)
		(local $$t18.417 i32)
		(local $$t18.418 i32)
		(local $$t18.419 i32)
		(local $$t18.420 i32)
		(local $$t18.421 i32)
		(local $$t18.422 i32)
		(local $$t18.423 i32)
		(local $$t18.424 i32)
		(local $$t18.425 i32)
		(local $$t18.426 i32)
		(local $$t18.427 i32)
		(local $$t18.428 i32)
		(local $$t18.429 i32)
		(local $$t18.430 i32)
		(local $$t18.431 i32)
		(local $$t18.432 i32)
		(local $$t18.433 i32)
		(local $$t18.434 i32)
		(local $$t18.435 i32)
		(local $$t18.436 i32)
		(local $$t18.437 i32)
		(local $$t18.438 i32)
		(local $$t18.439 i32)
		(local $$t18.440 i32)
		(local $$t18.441 i32)
		(local $$t18.442 i32)
		(local $$t18.443 i32)
		(local $$t18.444 i32)
		(local $$t18.445 i32)
		(local $$t18.446 i32)
		(local $$t18.447 i32)
		(local $$t18.448 i32)
		(local $$t18.449 i32)
		(local $$t18.450 i32)
		(local $$t18.451 i32)
		(local $$t18.452 i32)
		(local $$t18.453 i32)
		(local $$t18.454 i32)
		(local $$t18.455 i32)
		(local $$t18.456 i32)
		(local $$t18.457 i32)
		(local $$t18.458 i32)
		(local $$t18.459 i32)
		(local $$t18.460 i32)
		(local $$t18.461 i32)
		(local $$t18.462 i32)
		(local $$t18.463 i32)
		(local $$t18.464 i32)
		(local $$t18.465 i32)
		(local $$t18.466 i32)
		(local $$t18.467 i32)
		(local $$t19.0 i32)
		(local $$t19.1 i32)
		(local $$t19.2 i32)
		(local $$t19.3 i32)
		(local $$t19.4 i32)
		(local $$t19.5 i32)
		(local $$t19.6 i32)
		(local $$t19.7 i32)
		(local $$t19.8 i32)
		(local $$t19.9 i32)
		(local $$t19.10 i32)
		(local $$t19.11 i32)
		(local $$t19.12 i32)
		(local $$t19.13 i32)
		(local $$t19.14 i32)
		(local $$t19.15 i32)
		(local $$t19.16 i32)
		(local $$t19.17 i32)
		(local $$t19.18 i32)
		(local $$t19.19 i32)
		(local $$t19.20 i32)
		(local $$t19.21 i32)
		(local $$t19.22 i32)
		(local $$t19.23 i32)
		(local $$t19.24 i32)
		(local $$t19.25 i32)
		(local $$t19.26 i32)
		(local $$t19.27 i32)
		(local $$t19.28 i32)
		(local $$t19.29 i32)
		(local $$t19.30 i32)
		(local $$t19.31 i32)
		(local $$t19.32 i32)
		(local $$t19.33 i32)
		(local $$t19.34 i32)
		(local $$t19.35 i32)
		(local $$t19.36 i32)
		(local $$t19.37 i32)
		(local $$t19.38 i32)
		(local $$t19.39 i32)
		(local $$t19.40 i32)
		(local $$t19.41 i32)
		(local $$t19.42 i32)
		(local $$t19.43 i32)
		(local $$t19.44 i32)
		(local $$t19.45 i32)
		(local $$t19.46 i32)
		(local $$t19.47 i32)
		(local $$t19.48 i32)
		(local $$t19.49 i32)
		(local $$t19.50 i32)
		(local $$t19.51 i32)
		(local $$t19.52 i32)
		(local $$t19.53 i32)
		(local $$t19.54 i32)
		(local $$t19.55 i32)
		(local $$t19.56 i32)
		(local $$t19.57 i32)
		(local $$t19.58 i32)
		(local $$t19.59 i32)
		(local $$t19.60 i32)
		(local $$t19.61 i32)
		(local $$t19.62 i32)
		(local $$t19.63 i32)
		(local $$t19.64 i32)
		(local $$t19.65 i32)
		(local $$t19.66 i32)
		(local $$t19.67 i32)
		(local $$t19.68 i32)
		(local $$t19.69 i32)
		(local $$t19.70 i32)
		(local $$t19.71 i32)
		(local $$t19.72 i32)
		(local $$t19.73 i32)
		(local $$t19.74 i32)
		(local $$t19.75 i32)
		(local $$t19.76 i32)
		(local $$t19.77 i32)
		(local $$t19.78 i32)
		(local $$t19.79 i32)
		(local $$t19.80 i32)
		(local $$t19.81 i32)
		(local $$t19.82 i32)
		(local $$t19.83 i32)
		(local $$t19.84 i32)
		(local $$t19.85 i32)
		(local $$t19.86 i32)
		(local $$t19.87 i32)
		(local $$t19.88 i32)
		(local $$t19.89 i32)
		(local $$t19.90 i32)
		(local $$t19.91 i32)
		(local $$t19.92 i32)
		(local $$t19.93 i32)
		(local $$t19.94 i32)
		(local $$t20.0 i32)
		(local $$t20.1 i32)
		(local $$t20.2 i32)
		(local $$t20.3 i32)
		(local $$t21 i32)
		(local $$t22 i32)
		(local $$t23 i32)
		(local $$t24.0 i32)
		(local $$t24.1 i32)
		(local $$t24.2 i32)
		(local $$t24.3 i32)
		(local $$t25 i32)
		(local $$t26 i32)
		(local $$t27 i32)
		(local $$t28.0 i32)
		(local $$t28.1 i32)
		(local $$t29 i32)
		(local $$t30 i32)
		(local $$t31 i32)
		(local $$t32.0 i32)
		(local $$t32.1 i32)
		(local $$t33 i32)
		(local $$t34 i32)
		(local $$t35.0 i32)
		(local $$t35.1 i32)
		(local $$t36 i32)
		(local $$t37 i32)
		(local $$t38 i32)
		(local $$t39 i32)
		(local $$t40 i32)
		(local $$t41.0 i32)
		(local $$t41.1 i32)
		(local $$t42 i32)
		(local $$t43 i32)
		(local $$t44 i32)
		(local $$t45.0 i32)
		(local $$t45.1 i32)
		(local $$t46 i32)
		(local $$t47 i32)
		(local $$t48 i32)
		(local $$t49.0 i32)
		(local $$t49.1 i32)
		(local $$t49.2 i32)
		(local $$t49.3 i32)
		(local $$t50 i32)
		(local $$t51 i32)
		(local $$t52 i32)
		(local $$t53.0 i32)
		(local $$t53.1 i32)
		(local $$t54 i32)
		(local $$t55 i32)
		(local $$t56 i32)
		(local $$t57 i32)
		block $$BlockFnBody
			loop $$BlockDisp
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
																															br_table  0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 0
																														end
																														i32.const 0
																														local.set $$current_block
																														local.get $r
																														i32.const 255
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
																													i32.const 32
																													local.get $r
																													i32.le_s
																													local.set $$t1
																													local.get $$t1
																													if
																														br $$Block_4
																													else
																														br $$Block_3
																													end
																												end
																												i32.const 2
																												local.set $$current_block
																												i32.const 0
																												local.get $r
																												i32.le_s
																												local.set $$t2
																												local.get $$t2
																												if
																													br $$Block_10
																												else
																													br $$Block_9
																												end
																											end
																											i32.const 3
																											local.set $$current_block
																											i32.const 1
																											local.set $$ret_0
																											br $$BlockFnBody
																										end
																										i32.const 4
																										local.set $$current_block
																										i32.const 161
																										local.get $r
																										i32.le_s
																										local.set $$t3
																										local.get $$t3
																										if
																											br $$Block_7
																										else
																											br $$Block_6
																										end
																									end
																									i32.const 5
																									local.set $$current_block
																									local.get $r
																									i32.const 126
																									i32.le_s
																									local.set $$t4
																									local.get $$t4
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
																								local.get $r
																								i32.const 173
																								i32.eq
																								i32.eqz
																								local.set $$t5
																								local.get $$t5
																								local.set $$ret_0
																								br $$BlockFnBody
																							end
																							i32.const 7
																							local.set $$current_block
																							i32.const 0
																							local.set $$ret_0
																							br $$BlockFnBody
																						end
																						i32.const 8
																						local.set $$current_block
																						local.get $r
																						i32.const 255
																						i32.le_s
																						local.set $$t6
																						local.get $$t6
																						if
																							i32.const 6
																							local.set $$block_selector
																							br $$BlockDisp
																						else
																							i32.const 7
																							local.set $$block_selector
																							br $$BlockDisp
																						end
																					end
																					i32.const 9
																					local.set $$current_block
																					i32.const 884
																					call $runtime.HeapAlloc
																					i32.const 1
																					i32.const 0
																					i32.const 868
																					call $runtime.Block.Init
																					call $runtime.DupI32
																					i32.const 16
																					i32.add
																					local.set $$t7.1
																					local.get $$t7.0
																					call $runtime.Block.Release
																					local.set $$t7.0
																					i32.const 280
																					call $runtime.HeapAlloc
																					i32.const 1
																					i32.const 0
																					i32.const 264
																					call $runtime.Block.Init
																					call $runtime.DupI32
																					i32.const 16
																					i32.add
																					local.set $$t8.1
																					local.get $$t8.0
																					call $runtime.Block.Release
																					local.set $$t8.0
																					local.get $r
																					i32.const 65535
																					i32.and
																					local.set $$t9
																					i32.const 26910
																					i32.load16_u
																					i32.const 26910
																					i32.load16_u offset=2
																					i32.const 26910
																					i32.load16_u offset=4
																					i32.const 26910
																					i32.load16_u offset=6
																					i32.const 26910
																					i32.load16_u offset=8
																					i32.const 26910
																					i32.load16_u offset=10
																					i32.const 26910
																					i32.load16_u offset=12
																					i32.const 26910
																					i32.load16_u offset=14
																					i32.const 26910
																					i32.load16_u offset=16
																					i32.const 26910
																					i32.load16_u offset=18
																					i32.const 26910
																					i32.load16_u offset=20
																					i32.const 26910
																					i32.load16_u offset=22
																					i32.const 26910
																					i32.load16_u offset=24
																					i32.const 26910
																					i32.load16_u offset=26
																					i32.const 26910
																					i32.load16_u offset=28
																					i32.const 26910
																					i32.load16_u offset=30
																					i32.const 26910
																					i32.load16_u offset=32
																					i32.const 26910
																					i32.load16_u offset=34
																					i32.const 26910
																					i32.load16_u offset=36
																					i32.const 26910
																					i32.load16_u offset=38
																					i32.const 26910
																					i32.load16_u offset=40
																					i32.const 26910
																					i32.load16_u offset=42
																					i32.const 26910
																					i32.load16_u offset=44
																					i32.const 26910
																					i32.load16_u offset=46
																					i32.const 26910
																					i32.load16_u offset=48
																					i32.const 26910
																					i32.load16_u offset=50
																					i32.const 26910
																					i32.load16_u offset=52
																					i32.const 26910
																					i32.load16_u offset=54
																					i32.const 26910
																					i32.load16_u offset=56
																					i32.const 26910
																					i32.load16_u offset=58
																					i32.const 26910
																					i32.load16_u offset=60
																					i32.const 26910
																					i32.load16_u offset=62
																					i32.const 26910
																					i32.load16_u offset=64
																					i32.const 26910
																					i32.load16_u offset=66
																					i32.const 26910
																					i32.load16_u offset=68
																					i32.const 26910
																					i32.load16_u offset=70
																					i32.const 26910
																					i32.load16_u offset=72
																					i32.const 26910
																					i32.load16_u offset=74
																					i32.const 26910
																					i32.load16_u offset=76
																					i32.const 26910
																					i32.load16_u offset=78
																					i32.const 26910
																					i32.load16_u offset=80
																					i32.const 26910
																					i32.load16_u offset=82
																					i32.const 26910
																					i32.load16_u offset=84
																					i32.const 26910
																					i32.load16_u offset=86
																					i32.const 26910
																					i32.load16_u offset=88
																					i32.const 26910
																					i32.load16_u offset=90
																					i32.const 26910
																					i32.load16_u offset=92
																					i32.const 26910
																					i32.load16_u offset=94
																					i32.const 26910
																					i32.load16_u offset=96
																					i32.const 26910
																					i32.load16_u offset=98
																					i32.const 26910
																					i32.load16_u offset=100
																					i32.const 26910
																					i32.load16_u offset=102
																					i32.const 26910
																					i32.load16_u offset=104
																					i32.const 26910
																					i32.load16_u offset=106
																					i32.const 26910
																					i32.load16_u offset=108
																					i32.const 26910
																					i32.load16_u offset=110
																					i32.const 26910
																					i32.load16_u offset=112
																					i32.const 26910
																					i32.load16_u offset=114
																					i32.const 26910
																					i32.load16_u offset=116
																					i32.const 26910
																					i32.load16_u offset=118
																					i32.const 26910
																					i32.load16_u offset=120
																					i32.const 26910
																					i32.load16_u offset=122
																					i32.const 26910
																					i32.load16_u offset=124
																					i32.const 26910
																					i32.load16_u offset=126
																					i32.const 26910
																					i32.load16_u offset=128
																					i32.const 26910
																					i32.load16_u offset=130
																					i32.const 26910
																					i32.load16_u offset=132
																					i32.const 26910
																					i32.load16_u offset=134
																					i32.const 26910
																					i32.load16_u offset=136
																					i32.const 26910
																					i32.load16_u offset=138
																					i32.const 26910
																					i32.load16_u offset=140
																					i32.const 26910
																					i32.load16_u offset=142
																					i32.const 26910
																					i32.load16_u offset=144
																					i32.const 26910
																					i32.load16_u offset=146
																					i32.const 26910
																					i32.load16_u offset=148
																					i32.const 26910
																					i32.load16_u offset=150
																					i32.const 26910
																					i32.load16_u offset=152
																					i32.const 26910
																					i32.load16_u offset=154
																					i32.const 26910
																					i32.load16_u offset=156
																					i32.const 26910
																					i32.load16_u offset=158
																					i32.const 26910
																					i32.load16_u offset=160
																					i32.const 26910
																					i32.load16_u offset=162
																					i32.const 26910
																					i32.load16_u offset=164
																					i32.const 26910
																					i32.load16_u offset=166
																					i32.const 26910
																					i32.load16_u offset=168
																					i32.const 26910
																					i32.load16_u offset=170
																					i32.const 26910
																					i32.load16_u offset=172
																					i32.const 26910
																					i32.load16_u offset=174
																					i32.const 26910
																					i32.load16_u offset=176
																					i32.const 26910
																					i32.load16_u offset=178
																					i32.const 26910
																					i32.load16_u offset=180
																					i32.const 26910
																					i32.load16_u offset=182
																					i32.const 26910
																					i32.load16_u offset=184
																					i32.const 26910
																					i32.load16_u offset=186
																					i32.const 26910
																					i32.load16_u offset=188
																					i32.const 26910
																					i32.load16_u offset=190
																					i32.const 26910
																					i32.load16_u offset=192
																					i32.const 26910
																					i32.load16_u offset=194
																					i32.const 26910
																					i32.load16_u offset=196
																					i32.const 26910
																					i32.load16_u offset=198
																					i32.const 26910
																					i32.load16_u offset=200
																					i32.const 26910
																					i32.load16_u offset=202
																					i32.const 26910
																					i32.load16_u offset=204
																					i32.const 26910
																					i32.load16_u offset=206
																					i32.const 26910
																					i32.load16_u offset=208
																					i32.const 26910
																					i32.load16_u offset=210
																					i32.const 26910
																					i32.load16_u offset=212
																					i32.const 26910
																					i32.load16_u offset=214
																					i32.const 26910
																					i32.load16_u offset=216
																					i32.const 26910
																					i32.load16_u offset=218
																					i32.const 26910
																					i32.load16_u offset=220
																					i32.const 26910
																					i32.load16_u offset=222
																					i32.const 26910
																					i32.load16_u offset=224
																					i32.const 26910
																					i32.load16_u offset=226
																					i32.const 26910
																					i32.load16_u offset=228
																					i32.const 26910
																					i32.load16_u offset=230
																					i32.const 26910
																					i32.load16_u offset=232
																					i32.const 26910
																					i32.load16_u offset=234
																					i32.const 26910
																					i32.load16_u offset=236
																					i32.const 26910
																					i32.load16_u offset=238
																					i32.const 26910
																					i32.load16_u offset=240
																					i32.const 26910
																					i32.load16_u offset=242
																					i32.const 26910
																					i32.load16_u offset=244
																					i32.const 26910
																					i32.load16_u offset=246
																					i32.const 26910
																					i32.load16_u offset=248
																					i32.const 26910
																					i32.load16_u offset=250
																					i32.const 26910
																					i32.load16_u offset=252
																					i32.const 26910
																					i32.load16_u offset=254
																					i32.const 26910
																					i32.load16_u offset=256
																					i32.const 26910
																					i32.load16_u offset=258
																					i32.const 26910
																					i32.load16_u offset=260
																					i32.const 26910
																					i32.load16_u offset=262
																					i32.const 26910
																					i32.load16_u offset=264
																					i32.const 26910
																					i32.load16_u offset=266
																					i32.const 26910
																					i32.load16_u offset=268
																					i32.const 26910
																					i32.load16_u offset=270
																					i32.const 26910
																					i32.load16_u offset=272
																					i32.const 26910
																					i32.load16_u offset=274
																					i32.const 26910
																					i32.load16_u offset=276
																					i32.const 26910
																					i32.load16_u offset=278
																					i32.const 26910
																					i32.load16_u offset=280
																					i32.const 26910
																					i32.load16_u offset=282
																					i32.const 26910
																					i32.load16_u offset=284
																					i32.const 26910
																					i32.load16_u offset=286
																					i32.const 26910
																					i32.load16_u offset=288
																					i32.const 26910
																					i32.load16_u offset=290
																					i32.const 26910
																					i32.load16_u offset=292
																					i32.const 26910
																					i32.load16_u offset=294
																					i32.const 26910
																					i32.load16_u offset=296
																					i32.const 26910
																					i32.load16_u offset=298
																					i32.const 26910
																					i32.load16_u offset=300
																					i32.const 26910
																					i32.load16_u offset=302
																					i32.const 26910
																					i32.load16_u offset=304
																					i32.const 26910
																					i32.load16_u offset=306
																					i32.const 26910
																					i32.load16_u offset=308
																					i32.const 26910
																					i32.load16_u offset=310
																					i32.const 26910
																					i32.load16_u offset=312
																					i32.const 26910
																					i32.load16_u offset=314
																					i32.const 26910
																					i32.load16_u offset=316
																					i32.const 26910
																					i32.load16_u offset=318
																					i32.const 26910
																					i32.load16_u offset=320
																					i32.const 26910
																					i32.load16_u offset=322
																					i32.const 26910
																					i32.load16_u offset=324
																					i32.const 26910
																					i32.load16_u offset=326
																					i32.const 26910
																					i32.load16_u offset=328
																					i32.const 26910
																					i32.load16_u offset=330
																					i32.const 26910
																					i32.load16_u offset=332
																					i32.const 26910
																					i32.load16_u offset=334
																					i32.const 26910
																					i32.load16_u offset=336
																					i32.const 26910
																					i32.load16_u offset=338
																					i32.const 26910
																					i32.load16_u offset=340
																					i32.const 26910
																					i32.load16_u offset=342
																					i32.const 26910
																					i32.load16_u offset=344
																					i32.const 26910
																					i32.load16_u offset=346
																					i32.const 26910
																					i32.load16_u offset=348
																					i32.const 26910
																					i32.load16_u offset=350
																					i32.const 26910
																					i32.load16_u offset=352
																					i32.const 26910
																					i32.load16_u offset=354
																					i32.const 26910
																					i32.load16_u offset=356
																					i32.const 26910
																					i32.load16_u offset=358
																					i32.const 26910
																					i32.load16_u offset=360
																					i32.const 26910
																					i32.load16_u offset=362
																					i32.const 26910
																					i32.load16_u offset=364
																					i32.const 26910
																					i32.load16_u offset=366
																					i32.const 26910
																					i32.load16_u offset=368
																					i32.const 26910
																					i32.load16_u offset=370
																					i32.const 26910
																					i32.load16_u offset=372
																					i32.const 26910
																					i32.load16_u offset=374
																					i32.const 26910
																					i32.load16_u offset=376
																					i32.const 26910
																					i32.load16_u offset=378
																					i32.const 26910
																					i32.load16_u offset=380
																					i32.const 26910
																					i32.load16_u offset=382
																					i32.const 26910
																					i32.load16_u offset=384
																					i32.const 26910
																					i32.load16_u offset=386
																					i32.const 26910
																					i32.load16_u offset=388
																					i32.const 26910
																					i32.load16_u offset=390
																					i32.const 26910
																					i32.load16_u offset=392
																					i32.const 26910
																					i32.load16_u offset=394
																					i32.const 26910
																					i32.load16_u offset=396
																					i32.const 26910
																					i32.load16_u offset=398
																					i32.const 26910
																					i32.load16_u offset=400
																					i32.const 26910
																					i32.load16_u offset=402
																					i32.const 26910
																					i32.load16_u offset=404
																					i32.const 26910
																					i32.load16_u offset=406
																					i32.const 26910
																					i32.load16_u offset=408
																					i32.const 26910
																					i32.load16_u offset=410
																					i32.const 26910
																					i32.load16_u offset=412
																					i32.const 26910
																					i32.load16_u offset=414
																					i32.const 26910
																					i32.load16_u offset=416
																					i32.const 26910
																					i32.load16_u offset=418
																					i32.const 26910
																					i32.load16_u offset=420
																					i32.const 26910
																					i32.load16_u offset=422
																					i32.const 26910
																					i32.load16_u offset=424
																					i32.const 26910
																					i32.load16_u offset=426
																					i32.const 26910
																					i32.load16_u offset=428
																					i32.const 26910
																					i32.load16_u offset=430
																					i32.const 26910
																					i32.load16_u offset=432
																					i32.const 26910
																					i32.load16_u offset=434
																					i32.const 26910
																					i32.load16_u offset=436
																					i32.const 26910
																					i32.load16_u offset=438
																					i32.const 26910
																					i32.load16_u offset=440
																					i32.const 26910
																					i32.load16_u offset=442
																					i32.const 26910
																					i32.load16_u offset=444
																					i32.const 26910
																					i32.load16_u offset=446
																					i32.const 26910
																					i32.load16_u offset=448
																					i32.const 26910
																					i32.load16_u offset=450
																					i32.const 26910
																					i32.load16_u offset=452
																					i32.const 26910
																					i32.load16_u offset=454
																					i32.const 26910
																					i32.load16_u offset=456
																					i32.const 26910
																					i32.load16_u offset=458
																					i32.const 26910
																					i32.load16_u offset=460
																					i32.const 26910
																					i32.load16_u offset=462
																					i32.const 26910
																					i32.load16_u offset=464
																					i32.const 26910
																					i32.load16_u offset=466
																					i32.const 26910
																					i32.load16_u offset=468
																					i32.const 26910
																					i32.load16_u offset=470
																					i32.const 26910
																					i32.load16_u offset=472
																					i32.const 26910
																					i32.load16_u offset=474
																					i32.const 26910
																					i32.load16_u offset=476
																					i32.const 26910
																					i32.load16_u offset=478
																					i32.const 26910
																					i32.load16_u offset=480
																					i32.const 26910
																					i32.load16_u offset=482
																					i32.const 26910
																					i32.load16_u offset=484
																					i32.const 26910
																					i32.load16_u offset=486
																					i32.const 26910
																					i32.load16_u offset=488
																					i32.const 26910
																					i32.load16_u offset=490
																					i32.const 26910
																					i32.load16_u offset=492
																					i32.const 26910
																					i32.load16_u offset=494
																					i32.const 26910
																					i32.load16_u offset=496
																					i32.const 26910
																					i32.load16_u offset=498
																					i32.const 26910
																					i32.load16_u offset=500
																					i32.const 26910
																					i32.load16_u offset=502
																					i32.const 26910
																					i32.load16_u offset=504
																					i32.const 26910
																					i32.load16_u offset=506
																					i32.const 26910
																					i32.load16_u offset=508
																					i32.const 26910
																					i32.load16_u offset=510
																					i32.const 26910
																					i32.load16_u offset=512
																					i32.const 26910
																					i32.load16_u offset=514
																					i32.const 26910
																					i32.load16_u offset=516
																					i32.const 26910
																					i32.load16_u offset=518
																					i32.const 26910
																					i32.load16_u offset=520
																					i32.const 26910
																					i32.load16_u offset=522
																					i32.const 26910
																					i32.load16_u offset=524
																					i32.const 26910
																					i32.load16_u offset=526
																					i32.const 26910
																					i32.load16_u offset=528
																					i32.const 26910
																					i32.load16_u offset=530
																					i32.const 26910
																					i32.load16_u offset=532
																					i32.const 26910
																					i32.load16_u offset=534
																					i32.const 26910
																					i32.load16_u offset=536
																					i32.const 26910
																					i32.load16_u offset=538
																					i32.const 26910
																					i32.load16_u offset=540
																					i32.const 26910
																					i32.load16_u offset=542
																					i32.const 26910
																					i32.load16_u offset=544
																					i32.const 26910
																					i32.load16_u offset=546
																					i32.const 26910
																					i32.load16_u offset=548
																					i32.const 26910
																					i32.load16_u offset=550
																					i32.const 26910
																					i32.load16_u offset=552
																					i32.const 26910
																					i32.load16_u offset=554
																					i32.const 26910
																					i32.load16_u offset=556
																					i32.const 26910
																					i32.load16_u offset=558
																					i32.const 26910
																					i32.load16_u offset=560
																					i32.const 26910
																					i32.load16_u offset=562
																					i32.const 26910
																					i32.load16_u offset=564
																					i32.const 26910
																					i32.load16_u offset=566
																					i32.const 26910
																					i32.load16_u offset=568
																					i32.const 26910
																					i32.load16_u offset=570
																					i32.const 26910
																					i32.load16_u offset=572
																					i32.const 26910
																					i32.load16_u offset=574
																					i32.const 26910
																					i32.load16_u offset=576
																					i32.const 26910
																					i32.load16_u offset=578
																					i32.const 26910
																					i32.load16_u offset=580
																					i32.const 26910
																					i32.load16_u offset=582
																					i32.const 26910
																					i32.load16_u offset=584
																					i32.const 26910
																					i32.load16_u offset=586
																					i32.const 26910
																					i32.load16_u offset=588
																					i32.const 26910
																					i32.load16_u offset=590
																					i32.const 26910
																					i32.load16_u offset=592
																					i32.const 26910
																					i32.load16_u offset=594
																					i32.const 26910
																					i32.load16_u offset=596
																					i32.const 26910
																					i32.load16_u offset=598
																					i32.const 26910
																					i32.load16_u offset=600
																					i32.const 26910
																					i32.load16_u offset=602
																					i32.const 26910
																					i32.load16_u offset=604
																					i32.const 26910
																					i32.load16_u offset=606
																					i32.const 26910
																					i32.load16_u offset=608
																					i32.const 26910
																					i32.load16_u offset=610
																					i32.const 26910
																					i32.load16_u offset=612
																					i32.const 26910
																					i32.load16_u offset=614
																					i32.const 26910
																					i32.load16_u offset=616
																					i32.const 26910
																					i32.load16_u offset=618
																					i32.const 26910
																					i32.load16_u offset=620
																					i32.const 26910
																					i32.load16_u offset=622
																					i32.const 26910
																					i32.load16_u offset=624
																					i32.const 26910
																					i32.load16_u offset=626
																					i32.const 26910
																					i32.load16_u offset=628
																					i32.const 26910
																					i32.load16_u offset=630
																					i32.const 26910
																					i32.load16_u offset=632
																					i32.const 26910
																					i32.load16_u offset=634
																					i32.const 26910
																					i32.load16_u offset=636
																					i32.const 26910
																					i32.load16_u offset=638
																					i32.const 26910
																					i32.load16_u offset=640
																					i32.const 26910
																					i32.load16_u offset=642
																					i32.const 26910
																					i32.load16_u offset=644
																					i32.const 26910
																					i32.load16_u offset=646
																					i32.const 26910
																					i32.load16_u offset=648
																					i32.const 26910
																					i32.load16_u offset=650
																					i32.const 26910
																					i32.load16_u offset=652
																					i32.const 26910
																					i32.load16_u offset=654
																					i32.const 26910
																					i32.load16_u offset=656
																					i32.const 26910
																					i32.load16_u offset=658
																					i32.const 26910
																					i32.load16_u offset=660
																					i32.const 26910
																					i32.load16_u offset=662
																					i32.const 26910
																					i32.load16_u offset=664
																					i32.const 26910
																					i32.load16_u offset=666
																					i32.const 26910
																					i32.load16_u offset=668
																					i32.const 26910
																					i32.load16_u offset=670
																					i32.const 26910
																					i32.load16_u offset=672
																					i32.const 26910
																					i32.load16_u offset=674
																					i32.const 26910
																					i32.load16_u offset=676
																					i32.const 26910
																					i32.load16_u offset=678
																					i32.const 26910
																					i32.load16_u offset=680
																					i32.const 26910
																					i32.load16_u offset=682
																					i32.const 26910
																					i32.load16_u offset=684
																					i32.const 26910
																					i32.load16_u offset=686
																					i32.const 26910
																					i32.load16_u offset=688
																					i32.const 26910
																					i32.load16_u offset=690
																					i32.const 26910
																					i32.load16_u offset=692
																					i32.const 26910
																					i32.load16_u offset=694
																					i32.const 26910
																					i32.load16_u offset=696
																					i32.const 26910
																					i32.load16_u offset=698
																					i32.const 26910
																					i32.load16_u offset=700
																					i32.const 26910
																					i32.load16_u offset=702
																					i32.const 26910
																					i32.load16_u offset=704
																					i32.const 26910
																					i32.load16_u offset=706
																					i32.const 26910
																					i32.load16_u offset=708
																					i32.const 26910
																					i32.load16_u offset=710
																					i32.const 26910
																					i32.load16_u offset=712
																					i32.const 26910
																					i32.load16_u offset=714
																					i32.const 26910
																					i32.load16_u offset=716
																					i32.const 26910
																					i32.load16_u offset=718
																					i32.const 26910
																					i32.load16_u offset=720
																					i32.const 26910
																					i32.load16_u offset=722
																					i32.const 26910
																					i32.load16_u offset=724
																					i32.const 26910
																					i32.load16_u offset=726
																					i32.const 26910
																					i32.load16_u offset=728
																					i32.const 26910
																					i32.load16_u offset=730
																					i32.const 26910
																					i32.load16_u offset=732
																					i32.const 26910
																					i32.load16_u offset=734
																					i32.const 26910
																					i32.load16_u offset=736
																					i32.const 26910
																					i32.load16_u offset=738
																					i32.const 26910
																					i32.load16_u offset=740
																					i32.const 26910
																					i32.load16_u offset=742
																					i32.const 26910
																					i32.load16_u offset=744
																					i32.const 26910
																					i32.load16_u offset=746
																					i32.const 26910
																					i32.load16_u offset=748
																					i32.const 26910
																					i32.load16_u offset=750
																					i32.const 26910
																					i32.load16_u offset=752
																					i32.const 26910
																					i32.load16_u offset=754
																					i32.const 26910
																					i32.load16_u offset=756
																					i32.const 26910
																					i32.load16_u offset=758
																					i32.const 26910
																					i32.load16_u offset=760
																					i32.const 26910
																					i32.load16_u offset=762
																					i32.const 26910
																					i32.load16_u offset=764
																					i32.const 26910
																					i32.load16_u offset=766
																					i32.const 26910
																					i32.load16_u offset=768
																					i32.const 26910
																					i32.load16_u offset=770
																					i32.const 26910
																					i32.load16_u offset=772
																					i32.const 26910
																					i32.load16_u offset=774
																					i32.const 26910
																					i32.load16_u offset=776
																					i32.const 26910
																					i32.load16_u offset=778
																					i32.const 26910
																					i32.load16_u offset=780
																					i32.const 26910
																					i32.load16_u offset=782
																					i32.const 26910
																					i32.load16_u offset=784
																					i32.const 26910
																					i32.load16_u offset=786
																					i32.const 26910
																					i32.load16_u offset=788
																					i32.const 26910
																					i32.load16_u offset=790
																					i32.const 26910
																					i32.load16_u offset=792
																					i32.const 26910
																					i32.load16_u offset=794
																					i32.const 26910
																					i32.load16_u offset=796
																					i32.const 26910
																					i32.load16_u offset=798
																					i32.const 26910
																					i32.load16_u offset=800
																					i32.const 26910
																					i32.load16_u offset=802
																					i32.const 26910
																					i32.load16_u offset=804
																					i32.const 26910
																					i32.load16_u offset=806
																					i32.const 26910
																					i32.load16_u offset=808
																					i32.const 26910
																					i32.load16_u offset=810
																					i32.const 26910
																					i32.load16_u offset=812
																					i32.const 26910
																					i32.load16_u offset=814
																					i32.const 26910
																					i32.load16_u offset=816
																					i32.const 26910
																					i32.load16_u offset=818
																					i32.const 26910
																					i32.load16_u offset=820
																					i32.const 26910
																					i32.load16_u offset=822
																					i32.const 26910
																					i32.load16_u offset=824
																					i32.const 26910
																					i32.load16_u offset=826
																					i32.const 26910
																					i32.load16_u offset=828
																					i32.const 26910
																					i32.load16_u offset=830
																					i32.const 26910
																					i32.load16_u offset=832
																					i32.const 26910
																					i32.load16_u offset=834
																					i32.const 26910
																					i32.load16_u offset=836
																					i32.const 26910
																					i32.load16_u offset=838
																					i32.const 26910
																					i32.load16_u offset=840
																					i32.const 26910
																					i32.load16_u offset=842
																					i32.const 26910
																					i32.load16_u offset=844
																					i32.const 26910
																					i32.load16_u offset=846
																					i32.const 26910
																					i32.load16_u offset=848
																					i32.const 26910
																					i32.load16_u offset=850
																					i32.const 26910
																					i32.load16_u offset=852
																					i32.const 26910
																					i32.load16_u offset=854
																					i32.const 26910
																					i32.load16_u offset=856
																					i32.const 26910
																					i32.load16_u offset=858
																					i32.const 26910
																					i32.load16_u offset=860
																					i32.const 26910
																					i32.load16_u offset=862
																					i32.const 26910
																					i32.load16_u offset=864
																					i32.const 26910
																					i32.load16_u offset=866
																					local.set $$t10.433
																					local.set $$t10.432
																					local.set $$t10.431
																					local.set $$t10.430
																					local.set $$t10.429
																					local.set $$t10.428
																					local.set $$t10.427
																					local.set $$t10.426
																					local.set $$t10.425
																					local.set $$t10.424
																					local.set $$t10.423
																					local.set $$t10.422
																					local.set $$t10.421
																					local.set $$t10.420
																					local.set $$t10.419
																					local.set $$t10.418
																					local.set $$t10.417
																					local.set $$t10.416
																					local.set $$t10.415
																					local.set $$t10.414
																					local.set $$t10.413
																					local.set $$t10.412
																					local.set $$t10.411
																					local.set $$t10.410
																					local.set $$t10.409
																					local.set $$t10.408
																					local.set $$t10.407
																					local.set $$t10.406
																					local.set $$t10.405
																					local.set $$t10.404
																					local.set $$t10.403
																					local.set $$t10.402
																					local.set $$t10.401
																					local.set $$t10.400
																					local.set $$t10.399
																					local.set $$t10.398
																					local.set $$t10.397
																					local.set $$t10.396
																					local.set $$t10.395
																					local.set $$t10.394
																					local.set $$t10.393
																					local.set $$t10.392
																					local.set $$t10.391
																					local.set $$t10.390
																					local.set $$t10.389
																					local.set $$t10.388
																					local.set $$t10.387
																					local.set $$t10.386
																					local.set $$t10.385
																					local.set $$t10.384
																					local.set $$t10.383
																					local.set $$t10.382
																					local.set $$t10.381
																					local.set $$t10.380
																					local.set $$t10.379
																					local.set $$t10.378
																					local.set $$t10.377
																					local.set $$t10.376
																					local.set $$t10.375
																					local.set $$t10.374
																					local.set $$t10.373
																					local.set $$t10.372
																					local.set $$t10.371
																					local.set $$t10.370
																					local.set $$t10.369
																					local.set $$t10.368
																					local.set $$t10.367
																					local.set $$t10.366
																					local.set $$t10.365
																					local.set $$t10.364
																					local.set $$t10.363
																					local.set $$t10.362
																					local.set $$t10.361
																					local.set $$t10.360
																					local.set $$t10.359
																					local.set $$t10.358
																					local.set $$t10.357
																					local.set $$t10.356
																					local.set $$t10.355
																					local.set $$t10.354
																					local.set $$t10.353
																					local.set $$t10.352
																					local.set $$t10.351
																					local.set $$t10.350
																					local.set $$t10.349
																					local.set $$t10.348
																					local.set $$t10.347
																					local.set $$t10.346
																					local.set $$t10.345
																					local.set $$t10.344
																					local.set $$t10.343
																					local.set $$t10.342
																					local.set $$t10.341
																					local.set $$t10.340
																					local.set $$t10.339
																					local.set $$t10.338
																					local.set $$t10.337
																					local.set $$t10.336
																					local.set $$t10.335
																					local.set $$t10.334
																					local.set $$t10.333
																					local.set $$t10.332
																					local.set $$t10.331
																					local.set $$t10.330
																					local.set $$t10.329
																					local.set $$t10.328
																					local.set $$t10.327
																					local.set $$t10.326
																					local.set $$t10.325
																					local.set $$t10.324
																					local.set $$t10.323
																					local.set $$t10.322
																					local.set $$t10.321
																					local.set $$t10.320
																					local.set $$t10.319
																					local.set $$t10.318
																					local.set $$t10.317
																					local.set $$t10.316
																					local.set $$t10.315
																					local.set $$t10.314
																					local.set $$t10.313
																					local.set $$t10.312
																					local.set $$t10.311
																					local.set $$t10.310
																					local.set $$t10.309
																					local.set $$t10.308
																					local.set $$t10.307
																					local.set $$t10.306
																					local.set $$t10.305
																					local.set $$t10.304
																					local.set $$t10.303
																					local.set $$t10.302
																					local.set $$t10.301
																					local.set $$t10.300
																					local.set $$t10.299
																					local.set $$t10.298
																					local.set $$t10.297
																					local.set $$t10.296
																					local.set $$t10.295
																					local.set $$t10.294
																					local.set $$t10.293
																					local.set $$t10.292
																					local.set $$t10.291
																					local.set $$t10.290
																					local.set $$t10.289
																					local.set $$t10.288
																					local.set $$t10.287
																					local.set $$t10.286
																					local.set $$t10.285
																					local.set $$t10.284
																					local.set $$t10.283
																					local.set $$t10.282
																					local.set $$t10.281
																					local.set $$t10.280
																					local.set $$t10.279
																					local.set $$t10.278
																					local.set $$t10.277
																					local.set $$t10.276
																					local.set $$t10.275
																					local.set $$t10.274
																					local.set $$t10.273
																					local.set $$t10.272
																					local.set $$t10.271
																					local.set $$t10.270
																					local.set $$t10.269
																					local.set $$t10.268
																					local.set $$t10.267
																					local.set $$t10.266
																					local.set $$t10.265
																					local.set $$t10.264
																					local.set $$t10.263
																					local.set $$t10.262
																					local.set $$t10.261
																					local.set $$t10.260
																					local.set $$t10.259
																					local.set $$t10.258
																					local.set $$t10.257
																					local.set $$t10.256
																					local.set $$t10.255
																					local.set $$t10.254
																					local.set $$t10.253
																					local.set $$t10.252
																					local.set $$t10.251
																					local.set $$t10.250
																					local.set $$t10.249
																					local.set $$t10.248
																					local.set $$t10.247
																					local.set $$t10.246
																					local.set $$t10.245
																					local.set $$t10.244
																					local.set $$t10.243
																					local.set $$t10.242
																					local.set $$t10.241
																					local.set $$t10.240
																					local.set $$t10.239
																					local.set $$t10.238
																					local.set $$t10.237
																					local.set $$t10.236
																					local.set $$t10.235
																					local.set $$t10.234
																					local.set $$t10.233
																					local.set $$t10.232
																					local.set $$t10.231
																					local.set $$t10.230
																					local.set $$t10.229
																					local.set $$t10.228
																					local.set $$t10.227
																					local.set $$t10.226
																					local.set $$t10.225
																					local.set $$t10.224
																					local.set $$t10.223
																					local.set $$t10.222
																					local.set $$t10.221
																					local.set $$t10.220
																					local.set $$t10.219
																					local.set $$t10.218
																					local.set $$t10.217
																					local.set $$t10.216
																					local.set $$t10.215
																					local.set $$t10.214
																					local.set $$t10.213
																					local.set $$t10.212
																					local.set $$t10.211
																					local.set $$t10.210
																					local.set $$t10.209
																					local.set $$t10.208
																					local.set $$t10.207
																					local.set $$t10.206
																					local.set $$t10.205
																					local.set $$t10.204
																					local.set $$t10.203
																					local.set $$t10.202
																					local.set $$t10.201
																					local.set $$t10.200
																					local.set $$t10.199
																					local.set $$t10.198
																					local.set $$t10.197
																					local.set $$t10.196
																					local.set $$t10.195
																					local.set $$t10.194
																					local.set $$t10.193
																					local.set $$t10.192
																					local.set $$t10.191
																					local.set $$t10.190
																					local.set $$t10.189
																					local.set $$t10.188
																					local.set $$t10.187
																					local.set $$t10.186
																					local.set $$t10.185
																					local.set $$t10.184
																					local.set $$t10.183
																					local.set $$t10.182
																					local.set $$t10.181
																					local.set $$t10.180
																					local.set $$t10.179
																					local.set $$t10.178
																					local.set $$t10.177
																					local.set $$t10.176
																					local.set $$t10.175
																					local.set $$t10.174
																					local.set $$t10.173
																					local.set $$t10.172
																					local.set $$t10.171
																					local.set $$t10.170
																					local.set $$t10.169
																					local.set $$t10.168
																					local.set $$t10.167
																					local.set $$t10.166
																					local.set $$t10.165
																					local.set $$t10.164
																					local.set $$t10.163
																					local.set $$t10.162
																					local.set $$t10.161
																					local.set $$t10.160
																					local.set $$t10.159
																					local.set $$t10.158
																					local.set $$t10.157
																					local.set $$t10.156
																					local.set $$t10.155
																					local.set $$t10.154
																					local.set $$t10.153
																					local.set $$t10.152
																					local.set $$t10.151
																					local.set $$t10.150
																					local.set $$t10.149
																					local.set $$t10.148
																					local.set $$t10.147
																					local.set $$t10.146
																					local.set $$t10.145
																					local.set $$t10.144
																					local.set $$t10.143
																					local.set $$t10.142
																					local.set $$t10.141
																					local.set $$t10.140
																					local.set $$t10.139
																					local.set $$t10.138
																					local.set $$t10.137
																					local.set $$t10.136
																					local.set $$t10.135
																					local.set $$t10.134
																					local.set $$t10.133
																					local.set $$t10.132
																					local.set $$t10.131
																					local.set $$t10.130
																					local.set $$t10.129
																					local.set $$t10.128
																					local.set $$t10.127
																					local.set $$t10.126
																					local.set $$t10.125
																					local.set $$t10.124
																					local.set $$t10.123
																					local.set $$t10.122
																					local.set $$t10.121
																					local.set $$t10.120
																					local.set $$t10.119
																					local.set $$t10.118
																					local.set $$t10.117
																					local.set $$t10.116
																					local.set $$t10.115
																					local.set $$t10.114
																					local.set $$t10.113
																					local.set $$t10.112
																					local.set $$t10.111
																					local.set $$t10.110
																					local.set $$t10.109
																					local.set $$t10.108
																					local.set $$t10.107
																					local.set $$t10.106
																					local.set $$t10.105
																					local.set $$t10.104
																					local.set $$t10.103
																					local.set $$t10.102
																					local.set $$t10.101
																					local.set $$t10.100
																					local.set $$t10.99
																					local.set $$t10.98
																					local.set $$t10.97
																					local.set $$t10.96
																					local.set $$t10.95
																					local.set $$t10.94
																					local.set $$t10.93
																					local.set $$t10.92
																					local.set $$t10.91
																					local.set $$t10.90
																					local.set $$t10.89
																					local.set $$t10.88
																					local.set $$t10.87
																					local.set $$t10.86
																					local.set $$t10.85
																					local.set $$t10.84
																					local.set $$t10.83
																					local.set $$t10.82
																					local.set $$t10.81
																					local.set $$t10.80
																					local.set $$t10.79
																					local.set $$t10.78
																					local.set $$t10.77
																					local.set $$t10.76
																					local.set $$t10.75
																					local.set $$t10.74
																					local.set $$t10.73
																					local.set $$t10.72
																					local.set $$t10.71
																					local.set $$t10.70
																					local.set $$t10.69
																					local.set $$t10.68
																					local.set $$t10.67
																					local.set $$t10.66
																					local.set $$t10.65
																					local.set $$t10.64
																					local.set $$t10.63
																					local.set $$t10.62
																					local.set $$t10.61
																					local.set $$t10.60
																					local.set $$t10.59
																					local.set $$t10.58
																					local.set $$t10.57
																					local.set $$t10.56
																					local.set $$t10.55
																					local.set $$t10.54
																					local.set $$t10.53
																					local.set $$t10.52
																					local.set $$t10.51
																					local.set $$t10.50
																					local.set $$t10.49
																					local.set $$t10.48
																					local.set $$t10.47
																					local.set $$t10.46
																					local.set $$t10.45
																					local.set $$t10.44
																					local.set $$t10.43
																					local.set $$t10.42
																					local.set $$t10.41
																					local.set $$t10.40
																					local.set $$t10.39
																					local.set $$t10.38
																					local.set $$t10.37
																					local.set $$t10.36
																					local.set $$t10.35
																					local.set $$t10.34
																					local.set $$t10.33
																					local.set $$t10.32
																					local.set $$t10.31
																					local.set $$t10.30
																					local.set $$t10.29
																					local.set $$t10.28
																					local.set $$t10.27
																					local.set $$t10.26
																					local.set $$t10.25
																					local.set $$t10.24
																					local.set $$t10.23
																					local.set $$t10.22
																					local.set $$t10.21
																					local.set $$t10.20
																					local.set $$t10.19
																					local.set $$t10.18
																					local.set $$t10.17
																					local.set $$t10.16
																					local.set $$t10.15
																					local.set $$t10.14
																					local.set $$t10.13
																					local.set $$t10.12
																					local.set $$t10.11
																					local.set $$t10.10
																					local.set $$t10.9
																					local.set $$t10.8
																					local.set $$t10.7
																					local.set $$t10.6
																					local.set $$t10.5
																					local.set $$t10.4
																					local.set $$t10.3
																					local.set $$t10.2
																					local.set $$t10.1
																					local.set $$t10.0
																					i32.const 26456
																					i32.load16_u
																					i32.const 26456
																					i32.load16_u offset=2
																					i32.const 26456
																					i32.load16_u offset=4
																					i32.const 26456
																					i32.load16_u offset=6
																					i32.const 26456
																					i32.load16_u offset=8
																					i32.const 26456
																					i32.load16_u offset=10
																					i32.const 26456
																					i32.load16_u offset=12
																					i32.const 26456
																					i32.load16_u offset=14
																					i32.const 26456
																					i32.load16_u offset=16
																					i32.const 26456
																					i32.load16_u offset=18
																					i32.const 26456
																					i32.load16_u offset=20
																					i32.const 26456
																					i32.load16_u offset=22
																					i32.const 26456
																					i32.load16_u offset=24
																					i32.const 26456
																					i32.load16_u offset=26
																					i32.const 26456
																					i32.load16_u offset=28
																					i32.const 26456
																					i32.load16_u offset=30
																					i32.const 26456
																					i32.load16_u offset=32
																					i32.const 26456
																					i32.load16_u offset=34
																					i32.const 26456
																					i32.load16_u offset=36
																					i32.const 26456
																					i32.load16_u offset=38
																					i32.const 26456
																					i32.load16_u offset=40
																					i32.const 26456
																					i32.load16_u offset=42
																					i32.const 26456
																					i32.load16_u offset=44
																					i32.const 26456
																					i32.load16_u offset=46
																					i32.const 26456
																					i32.load16_u offset=48
																					i32.const 26456
																					i32.load16_u offset=50
																					i32.const 26456
																					i32.load16_u offset=52
																					i32.const 26456
																					i32.load16_u offset=54
																					i32.const 26456
																					i32.load16_u offset=56
																					i32.const 26456
																					i32.load16_u offset=58
																					i32.const 26456
																					i32.load16_u offset=60
																					i32.const 26456
																					i32.load16_u offset=62
																					i32.const 26456
																					i32.load16_u offset=64
																					i32.const 26456
																					i32.load16_u offset=66
																					i32.const 26456
																					i32.load16_u offset=68
																					i32.const 26456
																					i32.load16_u offset=70
																					i32.const 26456
																					i32.load16_u offset=72
																					i32.const 26456
																					i32.load16_u offset=74
																					i32.const 26456
																					i32.load16_u offset=76
																					i32.const 26456
																					i32.load16_u offset=78
																					i32.const 26456
																					i32.load16_u offset=80
																					i32.const 26456
																					i32.load16_u offset=82
																					i32.const 26456
																					i32.load16_u offset=84
																					i32.const 26456
																					i32.load16_u offset=86
																					i32.const 26456
																					i32.load16_u offset=88
																					i32.const 26456
																					i32.load16_u offset=90
																					i32.const 26456
																					i32.load16_u offset=92
																					i32.const 26456
																					i32.load16_u offset=94
																					i32.const 26456
																					i32.load16_u offset=96
																					i32.const 26456
																					i32.load16_u offset=98
																					i32.const 26456
																					i32.load16_u offset=100
																					i32.const 26456
																					i32.load16_u offset=102
																					i32.const 26456
																					i32.load16_u offset=104
																					i32.const 26456
																					i32.load16_u offset=106
																					i32.const 26456
																					i32.load16_u offset=108
																					i32.const 26456
																					i32.load16_u offset=110
																					i32.const 26456
																					i32.load16_u offset=112
																					i32.const 26456
																					i32.load16_u offset=114
																					i32.const 26456
																					i32.load16_u offset=116
																					i32.const 26456
																					i32.load16_u offset=118
																					i32.const 26456
																					i32.load16_u offset=120
																					i32.const 26456
																					i32.load16_u offset=122
																					i32.const 26456
																					i32.load16_u offset=124
																					i32.const 26456
																					i32.load16_u offset=126
																					i32.const 26456
																					i32.load16_u offset=128
																					i32.const 26456
																					i32.load16_u offset=130
																					i32.const 26456
																					i32.load16_u offset=132
																					i32.const 26456
																					i32.load16_u offset=134
																					i32.const 26456
																					i32.load16_u offset=136
																					i32.const 26456
																					i32.load16_u offset=138
																					i32.const 26456
																					i32.load16_u offset=140
																					i32.const 26456
																					i32.load16_u offset=142
																					i32.const 26456
																					i32.load16_u offset=144
																					i32.const 26456
																					i32.load16_u offset=146
																					i32.const 26456
																					i32.load16_u offset=148
																					i32.const 26456
																					i32.load16_u offset=150
																					i32.const 26456
																					i32.load16_u offset=152
																					i32.const 26456
																					i32.load16_u offset=154
																					i32.const 26456
																					i32.load16_u offset=156
																					i32.const 26456
																					i32.load16_u offset=158
																					i32.const 26456
																					i32.load16_u offset=160
																					i32.const 26456
																					i32.load16_u offset=162
																					i32.const 26456
																					i32.load16_u offset=164
																					i32.const 26456
																					i32.load16_u offset=166
																					i32.const 26456
																					i32.load16_u offset=168
																					i32.const 26456
																					i32.load16_u offset=170
																					i32.const 26456
																					i32.load16_u offset=172
																					i32.const 26456
																					i32.load16_u offset=174
																					i32.const 26456
																					i32.load16_u offset=176
																					i32.const 26456
																					i32.load16_u offset=178
																					i32.const 26456
																					i32.load16_u offset=180
																					i32.const 26456
																					i32.load16_u offset=182
																					i32.const 26456
																					i32.load16_u offset=184
																					i32.const 26456
																					i32.load16_u offset=186
																					i32.const 26456
																					i32.load16_u offset=188
																					i32.const 26456
																					i32.load16_u offset=190
																					i32.const 26456
																					i32.load16_u offset=192
																					i32.const 26456
																					i32.load16_u offset=194
																					i32.const 26456
																					i32.load16_u offset=196
																					i32.const 26456
																					i32.load16_u offset=198
																					i32.const 26456
																					i32.load16_u offset=200
																					i32.const 26456
																					i32.load16_u offset=202
																					i32.const 26456
																					i32.load16_u offset=204
																					i32.const 26456
																					i32.load16_u offset=206
																					i32.const 26456
																					i32.load16_u offset=208
																					i32.const 26456
																					i32.load16_u offset=210
																					i32.const 26456
																					i32.load16_u offset=212
																					i32.const 26456
																					i32.load16_u offset=214
																					i32.const 26456
																					i32.load16_u offset=216
																					i32.const 26456
																					i32.load16_u offset=218
																					i32.const 26456
																					i32.load16_u offset=220
																					i32.const 26456
																					i32.load16_u offset=222
																					i32.const 26456
																					i32.load16_u offset=224
																					i32.const 26456
																					i32.load16_u offset=226
																					i32.const 26456
																					i32.load16_u offset=228
																					i32.const 26456
																					i32.load16_u offset=230
																					i32.const 26456
																					i32.load16_u offset=232
																					i32.const 26456
																					i32.load16_u offset=234
																					i32.const 26456
																					i32.load16_u offset=236
																					i32.const 26456
																					i32.load16_u offset=238
																					i32.const 26456
																					i32.load16_u offset=240
																					i32.const 26456
																					i32.load16_u offset=242
																					i32.const 26456
																					i32.load16_u offset=244
																					i32.const 26456
																					i32.load16_u offset=246
																					i32.const 26456
																					i32.load16_u offset=248
																					i32.const 26456
																					i32.load16_u offset=250
																					i32.const 26456
																					i32.load16_u offset=252
																					i32.const 26456
																					i32.load16_u offset=254
																					i32.const 26456
																					i32.load16_u offset=256
																					i32.const 26456
																					i32.load16_u offset=258
																					i32.const 26456
																					i32.load16_u offset=260
																					i32.const 26456
																					i32.load16_u offset=262
																					local.set $$t11.131
																					local.set $$t11.130
																					local.set $$t11.129
																					local.set $$t11.128
																					local.set $$t11.127
																					local.set $$t11.126
																					local.set $$t11.125
																					local.set $$t11.124
																					local.set $$t11.123
																					local.set $$t11.122
																					local.set $$t11.121
																					local.set $$t11.120
																					local.set $$t11.119
																					local.set $$t11.118
																					local.set $$t11.117
																					local.set $$t11.116
																					local.set $$t11.115
																					local.set $$t11.114
																					local.set $$t11.113
																					local.set $$t11.112
																					local.set $$t11.111
																					local.set $$t11.110
																					local.set $$t11.109
																					local.set $$t11.108
																					local.set $$t11.107
																					local.set $$t11.106
																					local.set $$t11.105
																					local.set $$t11.104
																					local.set $$t11.103
																					local.set $$t11.102
																					local.set $$t11.101
																					local.set $$t11.100
																					local.set $$t11.99
																					local.set $$t11.98
																					local.set $$t11.97
																					local.set $$t11.96
																					local.set $$t11.95
																					local.set $$t11.94
																					local.set $$t11.93
																					local.set $$t11.92
																					local.set $$t11.91
																					local.set $$t11.90
																					local.set $$t11.89
																					local.set $$t11.88
																					local.set $$t11.87
																					local.set $$t11.86
																					local.set $$t11.85
																					local.set $$t11.84
																					local.set $$t11.83
																					local.set $$t11.82
																					local.set $$t11.81
																					local.set $$t11.80
																					local.set $$t11.79
																					local.set $$t11.78
																					local.set $$t11.77
																					local.set $$t11.76
																					local.set $$t11.75
																					local.set $$t11.74
																					local.set $$t11.73
																					local.set $$t11.72
																					local.set $$t11.71
																					local.set $$t11.70
																					local.set $$t11.69
																					local.set $$t11.68
																					local.set $$t11.67
																					local.set $$t11.66
																					local.set $$t11.65
																					local.set $$t11.64
																					local.set $$t11.63
																					local.set $$t11.62
																					local.set $$t11.61
																					local.set $$t11.60
																					local.set $$t11.59
																					local.set $$t11.58
																					local.set $$t11.57
																					local.set $$t11.56
																					local.set $$t11.55
																					local.set $$t11.54
																					local.set $$t11.53
																					local.set $$t11.52
																					local.set $$t11.51
																					local.set $$t11.50
																					local.set $$t11.49
																					local.set $$t11.48
																					local.set $$t11.47
																					local.set $$t11.46
																					local.set $$t11.45
																					local.set $$t11.44
																					local.set $$t11.43
																					local.set $$t11.42
																					local.set $$t11.41
																					local.set $$t11.40
																					local.set $$t11.39
																					local.set $$t11.38
																					local.set $$t11.37
																					local.set $$t11.36
																					local.set $$t11.35
																					local.set $$t11.34
																					local.set $$t11.33
																					local.set $$t11.32
																					local.set $$t11.31
																					local.set $$t11.30
																					local.set $$t11.29
																					local.set $$t11.28
																					local.set $$t11.27
																					local.set $$t11.26
																					local.set $$t11.25
																					local.set $$t11.24
																					local.set $$t11.23
																					local.set $$t11.22
																					local.set $$t11.21
																					local.set $$t11.20
																					local.set $$t11.19
																					local.set $$t11.18
																					local.set $$t11.17
																					local.set $$t11.16
																					local.set $$t11.15
																					local.set $$t11.14
																					local.set $$t11.13
																					local.set $$t11.12
																					local.set $$t11.11
																					local.set $$t11.10
																					local.set $$t11.9
																					local.set $$t11.8
																					local.set $$t11.7
																					local.set $$t11.6
																					local.set $$t11.5
																					local.set $$t11.4
																					local.set $$t11.3
																					local.set $$t11.2
																					local.set $$t11.1
																					local.set $$t11.0
																					local.get $$t7.1
																					local.get $$t10.0
																					i32.store16
																					local.get $$t7.1
																					local.get $$t10.1
																					i32.store16 offset=2
																					local.get $$t7.1
																					local.get $$t10.2
																					i32.store16 offset=4
																					local.get $$t7.1
																					local.get $$t10.3
																					i32.store16 offset=6
																					local.get $$t7.1
																					local.get $$t10.4
																					i32.store16 offset=8
																					local.get $$t7.1
																					local.get $$t10.5
																					i32.store16 offset=10
																					local.get $$t7.1
																					local.get $$t10.6
																					i32.store16 offset=12
																					local.get $$t7.1
																					local.get $$t10.7
																					i32.store16 offset=14
																					local.get $$t7.1
																					local.get $$t10.8
																					i32.store16 offset=16
																					local.get $$t7.1
																					local.get $$t10.9
																					i32.store16 offset=18
																					local.get $$t7.1
																					local.get $$t10.10
																					i32.store16 offset=20
																					local.get $$t7.1
																					local.get $$t10.11
																					i32.store16 offset=22
																					local.get $$t7.1
																					local.get $$t10.12
																					i32.store16 offset=24
																					local.get $$t7.1
																					local.get $$t10.13
																					i32.store16 offset=26
																					local.get $$t7.1
																					local.get $$t10.14
																					i32.store16 offset=28
																					local.get $$t7.1
																					local.get $$t10.15
																					i32.store16 offset=30
																					local.get $$t7.1
																					local.get $$t10.16
																					i32.store16 offset=32
																					local.get $$t7.1
																					local.get $$t10.17
																					i32.store16 offset=34
																					local.get $$t7.1
																					local.get $$t10.18
																					i32.store16 offset=36
																					local.get $$t7.1
																					local.get $$t10.19
																					i32.store16 offset=38
																					local.get $$t7.1
																					local.get $$t10.20
																					i32.store16 offset=40
																					local.get $$t7.1
																					local.get $$t10.21
																					i32.store16 offset=42
																					local.get $$t7.1
																					local.get $$t10.22
																					i32.store16 offset=44
																					local.get $$t7.1
																					local.get $$t10.23
																					i32.store16 offset=46
																					local.get $$t7.1
																					local.get $$t10.24
																					i32.store16 offset=48
																					local.get $$t7.1
																					local.get $$t10.25
																					i32.store16 offset=50
																					local.get $$t7.1
																					local.get $$t10.26
																					i32.store16 offset=52
																					local.get $$t7.1
																					local.get $$t10.27
																					i32.store16 offset=54
																					local.get $$t7.1
																					local.get $$t10.28
																					i32.store16 offset=56
																					local.get $$t7.1
																					local.get $$t10.29
																					i32.store16 offset=58
																					local.get $$t7.1
																					local.get $$t10.30
																					i32.store16 offset=60
																					local.get $$t7.1
																					local.get $$t10.31
																					i32.store16 offset=62
																					local.get $$t7.1
																					local.get $$t10.32
																					i32.store16 offset=64
																					local.get $$t7.1
																					local.get $$t10.33
																					i32.store16 offset=66
																					local.get $$t7.1
																					local.get $$t10.34
																					i32.store16 offset=68
																					local.get $$t7.1
																					local.get $$t10.35
																					i32.store16 offset=70
																					local.get $$t7.1
																					local.get $$t10.36
																					i32.store16 offset=72
																					local.get $$t7.1
																					local.get $$t10.37
																					i32.store16 offset=74
																					local.get $$t7.1
																					local.get $$t10.38
																					i32.store16 offset=76
																					local.get $$t7.1
																					local.get $$t10.39
																					i32.store16 offset=78
																					local.get $$t7.1
																					local.get $$t10.40
																					i32.store16 offset=80
																					local.get $$t7.1
																					local.get $$t10.41
																					i32.store16 offset=82
																					local.get $$t7.1
																					local.get $$t10.42
																					i32.store16 offset=84
																					local.get $$t7.1
																					local.get $$t10.43
																					i32.store16 offset=86
																					local.get $$t7.1
																					local.get $$t10.44
																					i32.store16 offset=88
																					local.get $$t7.1
																					local.get $$t10.45
																					i32.store16 offset=90
																					local.get $$t7.1
																					local.get $$t10.46
																					i32.store16 offset=92
																					local.get $$t7.1
																					local.get $$t10.47
																					i32.store16 offset=94
																					local.get $$t7.1
																					local.get $$t10.48
																					i32.store16 offset=96
																					local.get $$t7.1
																					local.get $$t10.49
																					i32.store16 offset=98
																					local.get $$t7.1
																					local.get $$t10.50
																					i32.store16 offset=100
																					local.get $$t7.1
																					local.get $$t10.51
																					i32.store16 offset=102
																					local.get $$t7.1
																					local.get $$t10.52
																					i32.store16 offset=104
																					local.get $$t7.1
																					local.get $$t10.53
																					i32.store16 offset=106
																					local.get $$t7.1
																					local.get $$t10.54
																					i32.store16 offset=108
																					local.get $$t7.1
																					local.get $$t10.55
																					i32.store16 offset=110
																					local.get $$t7.1
																					local.get $$t10.56
																					i32.store16 offset=112
																					local.get $$t7.1
																					local.get $$t10.57
																					i32.store16 offset=114
																					local.get $$t7.1
																					local.get $$t10.58
																					i32.store16 offset=116
																					local.get $$t7.1
																					local.get $$t10.59
																					i32.store16 offset=118
																					local.get $$t7.1
																					local.get $$t10.60
																					i32.store16 offset=120
																					local.get $$t7.1
																					local.get $$t10.61
																					i32.store16 offset=122
																					local.get $$t7.1
																					local.get $$t10.62
																					i32.store16 offset=124
																					local.get $$t7.1
																					local.get $$t10.63
																					i32.store16 offset=126
																					local.get $$t7.1
																					local.get $$t10.64
																					i32.store16 offset=128
																					local.get $$t7.1
																					local.get $$t10.65
																					i32.store16 offset=130
																					local.get $$t7.1
																					local.get $$t10.66
																					i32.store16 offset=132
																					local.get $$t7.1
																					local.get $$t10.67
																					i32.store16 offset=134
																					local.get $$t7.1
																					local.get $$t10.68
																					i32.store16 offset=136
																					local.get $$t7.1
																					local.get $$t10.69
																					i32.store16 offset=138
																					local.get $$t7.1
																					local.get $$t10.70
																					i32.store16 offset=140
																					local.get $$t7.1
																					local.get $$t10.71
																					i32.store16 offset=142
																					local.get $$t7.1
																					local.get $$t10.72
																					i32.store16 offset=144
																					local.get $$t7.1
																					local.get $$t10.73
																					i32.store16 offset=146
																					local.get $$t7.1
																					local.get $$t10.74
																					i32.store16 offset=148
																					local.get $$t7.1
																					local.get $$t10.75
																					i32.store16 offset=150
																					local.get $$t7.1
																					local.get $$t10.76
																					i32.store16 offset=152
																					local.get $$t7.1
																					local.get $$t10.77
																					i32.store16 offset=154
																					local.get $$t7.1
																					local.get $$t10.78
																					i32.store16 offset=156
																					local.get $$t7.1
																					local.get $$t10.79
																					i32.store16 offset=158
																					local.get $$t7.1
																					local.get $$t10.80
																					i32.store16 offset=160
																					local.get $$t7.1
																					local.get $$t10.81
																					i32.store16 offset=162
																					local.get $$t7.1
																					local.get $$t10.82
																					i32.store16 offset=164
																					local.get $$t7.1
																					local.get $$t10.83
																					i32.store16 offset=166
																					local.get $$t7.1
																					local.get $$t10.84
																					i32.store16 offset=168
																					local.get $$t7.1
																					local.get $$t10.85
																					i32.store16 offset=170
																					local.get $$t7.1
																					local.get $$t10.86
																					i32.store16 offset=172
																					local.get $$t7.1
																					local.get $$t10.87
																					i32.store16 offset=174
																					local.get $$t7.1
																					local.get $$t10.88
																					i32.store16 offset=176
																					local.get $$t7.1
																					local.get $$t10.89
																					i32.store16 offset=178
																					local.get $$t7.1
																					local.get $$t10.90
																					i32.store16 offset=180
																					local.get $$t7.1
																					local.get $$t10.91
																					i32.store16 offset=182
																					local.get $$t7.1
																					local.get $$t10.92
																					i32.store16 offset=184
																					local.get $$t7.1
																					local.get $$t10.93
																					i32.store16 offset=186
																					local.get $$t7.1
																					local.get $$t10.94
																					i32.store16 offset=188
																					local.get $$t7.1
																					local.get $$t10.95
																					i32.store16 offset=190
																					local.get $$t7.1
																					local.get $$t10.96
																					i32.store16 offset=192
																					local.get $$t7.1
																					local.get $$t10.97
																					i32.store16 offset=194
																					local.get $$t7.1
																					local.get $$t10.98
																					i32.store16 offset=196
																					local.get $$t7.1
																					local.get $$t10.99
																					i32.store16 offset=198
																					local.get $$t7.1
																					local.get $$t10.100
																					i32.store16 offset=200
																					local.get $$t7.1
																					local.get $$t10.101
																					i32.store16 offset=202
																					local.get $$t7.1
																					local.get $$t10.102
																					i32.store16 offset=204
																					local.get $$t7.1
																					local.get $$t10.103
																					i32.store16 offset=206
																					local.get $$t7.1
																					local.get $$t10.104
																					i32.store16 offset=208
																					local.get $$t7.1
																					local.get $$t10.105
																					i32.store16 offset=210
																					local.get $$t7.1
																					local.get $$t10.106
																					i32.store16 offset=212
																					local.get $$t7.1
																					local.get $$t10.107
																					i32.store16 offset=214
																					local.get $$t7.1
																					local.get $$t10.108
																					i32.store16 offset=216
																					local.get $$t7.1
																					local.get $$t10.109
																					i32.store16 offset=218
																					local.get $$t7.1
																					local.get $$t10.110
																					i32.store16 offset=220
																					local.get $$t7.1
																					local.get $$t10.111
																					i32.store16 offset=222
																					local.get $$t7.1
																					local.get $$t10.112
																					i32.store16 offset=224
																					local.get $$t7.1
																					local.get $$t10.113
																					i32.store16 offset=226
																					local.get $$t7.1
																					local.get $$t10.114
																					i32.store16 offset=228
																					local.get $$t7.1
																					local.get $$t10.115
																					i32.store16 offset=230
																					local.get $$t7.1
																					local.get $$t10.116
																					i32.store16 offset=232
																					local.get $$t7.1
																					local.get $$t10.117
																					i32.store16 offset=234
																					local.get $$t7.1
																					local.get $$t10.118
																					i32.store16 offset=236
																					local.get $$t7.1
																					local.get $$t10.119
																					i32.store16 offset=238
																					local.get $$t7.1
																					local.get $$t10.120
																					i32.store16 offset=240
																					local.get $$t7.1
																					local.get $$t10.121
																					i32.store16 offset=242
																					local.get $$t7.1
																					local.get $$t10.122
																					i32.store16 offset=244
																					local.get $$t7.1
																					local.get $$t10.123
																					i32.store16 offset=246
																					local.get $$t7.1
																					local.get $$t10.124
																					i32.store16 offset=248
																					local.get $$t7.1
																					local.get $$t10.125
																					i32.store16 offset=250
																					local.get $$t7.1
																					local.get $$t10.126
																					i32.store16 offset=252
																					local.get $$t7.1
																					local.get $$t10.127
																					i32.store16 offset=254
																					local.get $$t7.1
																					local.get $$t10.128
																					i32.store16 offset=256
																					local.get $$t7.1
																					local.get $$t10.129
																					i32.store16 offset=258
																					local.get $$t7.1
																					local.get $$t10.130
																					i32.store16 offset=260
																					local.get $$t7.1
																					local.get $$t10.131
																					i32.store16 offset=262
																					local.get $$t7.1
																					local.get $$t10.132
																					i32.store16 offset=264
																					local.get $$t7.1
																					local.get $$t10.133
																					i32.store16 offset=266
																					local.get $$t7.1
																					local.get $$t10.134
																					i32.store16 offset=268
																					local.get $$t7.1
																					local.get $$t10.135
																					i32.store16 offset=270
																					local.get $$t7.1
																					local.get $$t10.136
																					i32.store16 offset=272
																					local.get $$t7.1
																					local.get $$t10.137
																					i32.store16 offset=274
																					local.get $$t7.1
																					local.get $$t10.138
																					i32.store16 offset=276
																					local.get $$t7.1
																					local.get $$t10.139
																					i32.store16 offset=278
																					local.get $$t7.1
																					local.get $$t10.140
																					i32.store16 offset=280
																					local.get $$t7.1
																					local.get $$t10.141
																					i32.store16 offset=282
																					local.get $$t7.1
																					local.get $$t10.142
																					i32.store16 offset=284
																					local.get $$t7.1
																					local.get $$t10.143
																					i32.store16 offset=286
																					local.get $$t7.1
																					local.get $$t10.144
																					i32.store16 offset=288
																					local.get $$t7.1
																					local.get $$t10.145
																					i32.store16 offset=290
																					local.get $$t7.1
																					local.get $$t10.146
																					i32.store16 offset=292
																					local.get $$t7.1
																					local.get $$t10.147
																					i32.store16 offset=294
																					local.get $$t7.1
																					local.get $$t10.148
																					i32.store16 offset=296
																					local.get $$t7.1
																					local.get $$t10.149
																					i32.store16 offset=298
																					local.get $$t7.1
																					local.get $$t10.150
																					i32.store16 offset=300
																					local.get $$t7.1
																					local.get $$t10.151
																					i32.store16 offset=302
																					local.get $$t7.1
																					local.get $$t10.152
																					i32.store16 offset=304
																					local.get $$t7.1
																					local.get $$t10.153
																					i32.store16 offset=306
																					local.get $$t7.1
																					local.get $$t10.154
																					i32.store16 offset=308
																					local.get $$t7.1
																					local.get $$t10.155
																					i32.store16 offset=310
																					local.get $$t7.1
																					local.get $$t10.156
																					i32.store16 offset=312
																					local.get $$t7.1
																					local.get $$t10.157
																					i32.store16 offset=314
																					local.get $$t7.1
																					local.get $$t10.158
																					i32.store16 offset=316
																					local.get $$t7.1
																					local.get $$t10.159
																					i32.store16 offset=318
																					local.get $$t7.1
																					local.get $$t10.160
																					i32.store16 offset=320
																					local.get $$t7.1
																					local.get $$t10.161
																					i32.store16 offset=322
																					local.get $$t7.1
																					local.get $$t10.162
																					i32.store16 offset=324
																					local.get $$t7.1
																					local.get $$t10.163
																					i32.store16 offset=326
																					local.get $$t7.1
																					local.get $$t10.164
																					i32.store16 offset=328
																					local.get $$t7.1
																					local.get $$t10.165
																					i32.store16 offset=330
																					local.get $$t7.1
																					local.get $$t10.166
																					i32.store16 offset=332
																					local.get $$t7.1
																					local.get $$t10.167
																					i32.store16 offset=334
																					local.get $$t7.1
																					local.get $$t10.168
																					i32.store16 offset=336
																					local.get $$t7.1
																					local.get $$t10.169
																					i32.store16 offset=338
																					local.get $$t7.1
																					local.get $$t10.170
																					i32.store16 offset=340
																					local.get $$t7.1
																					local.get $$t10.171
																					i32.store16 offset=342
																					local.get $$t7.1
																					local.get $$t10.172
																					i32.store16 offset=344
																					local.get $$t7.1
																					local.get $$t10.173
																					i32.store16 offset=346
																					local.get $$t7.1
																					local.get $$t10.174
																					i32.store16 offset=348
																					local.get $$t7.1
																					local.get $$t10.175
																					i32.store16 offset=350
																					local.get $$t7.1
																					local.get $$t10.176
																					i32.store16 offset=352
																					local.get $$t7.1
																					local.get $$t10.177
																					i32.store16 offset=354
																					local.get $$t7.1
																					local.get $$t10.178
																					i32.store16 offset=356
																					local.get $$t7.1
																					local.get $$t10.179
																					i32.store16 offset=358
																					local.get $$t7.1
																					local.get $$t10.180
																					i32.store16 offset=360
																					local.get $$t7.1
																					local.get $$t10.181
																					i32.store16 offset=362
																					local.get $$t7.1
																					local.get $$t10.182
																					i32.store16 offset=364
																					local.get $$t7.1
																					local.get $$t10.183
																					i32.store16 offset=366
																					local.get $$t7.1
																					local.get $$t10.184
																					i32.store16 offset=368
																					local.get $$t7.1
																					local.get $$t10.185
																					i32.store16 offset=370
																					local.get $$t7.1
																					local.get $$t10.186
																					i32.store16 offset=372
																					local.get $$t7.1
																					local.get $$t10.187
																					i32.store16 offset=374
																					local.get $$t7.1
																					local.get $$t10.188
																					i32.store16 offset=376
																					local.get $$t7.1
																					local.get $$t10.189
																					i32.store16 offset=378
																					local.get $$t7.1
																					local.get $$t10.190
																					i32.store16 offset=380
																					local.get $$t7.1
																					local.get $$t10.191
																					i32.store16 offset=382
																					local.get $$t7.1
																					local.get $$t10.192
																					i32.store16 offset=384
																					local.get $$t7.1
																					local.get $$t10.193
																					i32.store16 offset=386
																					local.get $$t7.1
																					local.get $$t10.194
																					i32.store16 offset=388
																					local.get $$t7.1
																					local.get $$t10.195
																					i32.store16 offset=390
																					local.get $$t7.1
																					local.get $$t10.196
																					i32.store16 offset=392
																					local.get $$t7.1
																					local.get $$t10.197
																					i32.store16 offset=394
																					local.get $$t7.1
																					local.get $$t10.198
																					i32.store16 offset=396
																					local.get $$t7.1
																					local.get $$t10.199
																					i32.store16 offset=398
																					local.get $$t7.1
																					local.get $$t10.200
																					i32.store16 offset=400
																					local.get $$t7.1
																					local.get $$t10.201
																					i32.store16 offset=402
																					local.get $$t7.1
																					local.get $$t10.202
																					i32.store16 offset=404
																					local.get $$t7.1
																					local.get $$t10.203
																					i32.store16 offset=406
																					local.get $$t7.1
																					local.get $$t10.204
																					i32.store16 offset=408
																					local.get $$t7.1
																					local.get $$t10.205
																					i32.store16 offset=410
																					local.get $$t7.1
																					local.get $$t10.206
																					i32.store16 offset=412
																					local.get $$t7.1
																					local.get $$t10.207
																					i32.store16 offset=414
																					local.get $$t7.1
																					local.get $$t10.208
																					i32.store16 offset=416
																					local.get $$t7.1
																					local.get $$t10.209
																					i32.store16 offset=418
																					local.get $$t7.1
																					local.get $$t10.210
																					i32.store16 offset=420
																					local.get $$t7.1
																					local.get $$t10.211
																					i32.store16 offset=422
																					local.get $$t7.1
																					local.get $$t10.212
																					i32.store16 offset=424
																					local.get $$t7.1
																					local.get $$t10.213
																					i32.store16 offset=426
																					local.get $$t7.1
																					local.get $$t10.214
																					i32.store16 offset=428
																					local.get $$t7.1
																					local.get $$t10.215
																					i32.store16 offset=430
																					local.get $$t7.1
																					local.get $$t10.216
																					i32.store16 offset=432
																					local.get $$t7.1
																					local.get $$t10.217
																					i32.store16 offset=434
																					local.get $$t7.1
																					local.get $$t10.218
																					i32.store16 offset=436
																					local.get $$t7.1
																					local.get $$t10.219
																					i32.store16 offset=438
																					local.get $$t7.1
																					local.get $$t10.220
																					i32.store16 offset=440
																					local.get $$t7.1
																					local.get $$t10.221
																					i32.store16 offset=442
																					local.get $$t7.1
																					local.get $$t10.222
																					i32.store16 offset=444
																					local.get $$t7.1
																					local.get $$t10.223
																					i32.store16 offset=446
																					local.get $$t7.1
																					local.get $$t10.224
																					i32.store16 offset=448
																					local.get $$t7.1
																					local.get $$t10.225
																					i32.store16 offset=450
																					local.get $$t7.1
																					local.get $$t10.226
																					i32.store16 offset=452
																					local.get $$t7.1
																					local.get $$t10.227
																					i32.store16 offset=454
																					local.get $$t7.1
																					local.get $$t10.228
																					i32.store16 offset=456
																					local.get $$t7.1
																					local.get $$t10.229
																					i32.store16 offset=458
																					local.get $$t7.1
																					local.get $$t10.230
																					i32.store16 offset=460
																					local.get $$t7.1
																					local.get $$t10.231
																					i32.store16 offset=462
																					local.get $$t7.1
																					local.get $$t10.232
																					i32.store16 offset=464
																					local.get $$t7.1
																					local.get $$t10.233
																					i32.store16 offset=466
																					local.get $$t7.1
																					local.get $$t10.234
																					i32.store16 offset=468
																					local.get $$t7.1
																					local.get $$t10.235
																					i32.store16 offset=470
																					local.get $$t7.1
																					local.get $$t10.236
																					i32.store16 offset=472
																					local.get $$t7.1
																					local.get $$t10.237
																					i32.store16 offset=474
																					local.get $$t7.1
																					local.get $$t10.238
																					i32.store16 offset=476
																					local.get $$t7.1
																					local.get $$t10.239
																					i32.store16 offset=478
																					local.get $$t7.1
																					local.get $$t10.240
																					i32.store16 offset=480
																					local.get $$t7.1
																					local.get $$t10.241
																					i32.store16 offset=482
																					local.get $$t7.1
																					local.get $$t10.242
																					i32.store16 offset=484
																					local.get $$t7.1
																					local.get $$t10.243
																					i32.store16 offset=486
																					local.get $$t7.1
																					local.get $$t10.244
																					i32.store16 offset=488
																					local.get $$t7.1
																					local.get $$t10.245
																					i32.store16 offset=490
																					local.get $$t7.1
																					local.get $$t10.246
																					i32.store16 offset=492
																					local.get $$t7.1
																					local.get $$t10.247
																					i32.store16 offset=494
																					local.get $$t7.1
																					local.get $$t10.248
																					i32.store16 offset=496
																					local.get $$t7.1
																					local.get $$t10.249
																					i32.store16 offset=498
																					local.get $$t7.1
																					local.get $$t10.250
																					i32.store16 offset=500
																					local.get $$t7.1
																					local.get $$t10.251
																					i32.store16 offset=502
																					local.get $$t7.1
																					local.get $$t10.252
																					i32.store16 offset=504
																					local.get $$t7.1
																					local.get $$t10.253
																					i32.store16 offset=506
																					local.get $$t7.1
																					local.get $$t10.254
																					i32.store16 offset=508
																					local.get $$t7.1
																					local.get $$t10.255
																					i32.store16 offset=510
																					local.get $$t7.1
																					local.get $$t10.256
																					i32.store16 offset=512
																					local.get $$t7.1
																					local.get $$t10.257
																					i32.store16 offset=514
																					local.get $$t7.1
																					local.get $$t10.258
																					i32.store16 offset=516
																					local.get $$t7.1
																					local.get $$t10.259
																					i32.store16 offset=518
																					local.get $$t7.1
																					local.get $$t10.260
																					i32.store16 offset=520
																					local.get $$t7.1
																					local.get $$t10.261
																					i32.store16 offset=522
																					local.get $$t7.1
																					local.get $$t10.262
																					i32.store16 offset=524
																					local.get $$t7.1
																					local.get $$t10.263
																					i32.store16 offset=526
																					local.get $$t7.1
																					local.get $$t10.264
																					i32.store16 offset=528
																					local.get $$t7.1
																					local.get $$t10.265
																					i32.store16 offset=530
																					local.get $$t7.1
																					local.get $$t10.266
																					i32.store16 offset=532
																					local.get $$t7.1
																					local.get $$t10.267
																					i32.store16 offset=534
																					local.get $$t7.1
																					local.get $$t10.268
																					i32.store16 offset=536
																					local.get $$t7.1
																					local.get $$t10.269
																					i32.store16 offset=538
																					local.get $$t7.1
																					local.get $$t10.270
																					i32.store16 offset=540
																					local.get $$t7.1
																					local.get $$t10.271
																					i32.store16 offset=542
																					local.get $$t7.1
																					local.get $$t10.272
																					i32.store16 offset=544
																					local.get $$t7.1
																					local.get $$t10.273
																					i32.store16 offset=546
																					local.get $$t7.1
																					local.get $$t10.274
																					i32.store16 offset=548
																					local.get $$t7.1
																					local.get $$t10.275
																					i32.store16 offset=550
																					local.get $$t7.1
																					local.get $$t10.276
																					i32.store16 offset=552
																					local.get $$t7.1
																					local.get $$t10.277
																					i32.store16 offset=554
																					local.get $$t7.1
																					local.get $$t10.278
																					i32.store16 offset=556
																					local.get $$t7.1
																					local.get $$t10.279
																					i32.store16 offset=558
																					local.get $$t7.1
																					local.get $$t10.280
																					i32.store16 offset=560
																					local.get $$t7.1
																					local.get $$t10.281
																					i32.store16 offset=562
																					local.get $$t7.1
																					local.get $$t10.282
																					i32.store16 offset=564
																					local.get $$t7.1
																					local.get $$t10.283
																					i32.store16 offset=566
																					local.get $$t7.1
																					local.get $$t10.284
																					i32.store16 offset=568
																					local.get $$t7.1
																					local.get $$t10.285
																					i32.store16 offset=570
																					local.get $$t7.1
																					local.get $$t10.286
																					i32.store16 offset=572
																					local.get $$t7.1
																					local.get $$t10.287
																					i32.store16 offset=574
																					local.get $$t7.1
																					local.get $$t10.288
																					i32.store16 offset=576
																					local.get $$t7.1
																					local.get $$t10.289
																					i32.store16 offset=578
																					local.get $$t7.1
																					local.get $$t10.290
																					i32.store16 offset=580
																					local.get $$t7.1
																					local.get $$t10.291
																					i32.store16 offset=582
																					local.get $$t7.1
																					local.get $$t10.292
																					i32.store16 offset=584
																					local.get $$t7.1
																					local.get $$t10.293
																					i32.store16 offset=586
																					local.get $$t7.1
																					local.get $$t10.294
																					i32.store16 offset=588
																					local.get $$t7.1
																					local.get $$t10.295
																					i32.store16 offset=590
																					local.get $$t7.1
																					local.get $$t10.296
																					i32.store16 offset=592
																					local.get $$t7.1
																					local.get $$t10.297
																					i32.store16 offset=594
																					local.get $$t7.1
																					local.get $$t10.298
																					i32.store16 offset=596
																					local.get $$t7.1
																					local.get $$t10.299
																					i32.store16 offset=598
																					local.get $$t7.1
																					local.get $$t10.300
																					i32.store16 offset=600
																					local.get $$t7.1
																					local.get $$t10.301
																					i32.store16 offset=602
																					local.get $$t7.1
																					local.get $$t10.302
																					i32.store16 offset=604
																					local.get $$t7.1
																					local.get $$t10.303
																					i32.store16 offset=606
																					local.get $$t7.1
																					local.get $$t10.304
																					i32.store16 offset=608
																					local.get $$t7.1
																					local.get $$t10.305
																					i32.store16 offset=610
																					local.get $$t7.1
																					local.get $$t10.306
																					i32.store16 offset=612
																					local.get $$t7.1
																					local.get $$t10.307
																					i32.store16 offset=614
																					local.get $$t7.1
																					local.get $$t10.308
																					i32.store16 offset=616
																					local.get $$t7.1
																					local.get $$t10.309
																					i32.store16 offset=618
																					local.get $$t7.1
																					local.get $$t10.310
																					i32.store16 offset=620
																					local.get $$t7.1
																					local.get $$t10.311
																					i32.store16 offset=622
																					local.get $$t7.1
																					local.get $$t10.312
																					i32.store16 offset=624
																					local.get $$t7.1
																					local.get $$t10.313
																					i32.store16 offset=626
																					local.get $$t7.1
																					local.get $$t10.314
																					i32.store16 offset=628
																					local.get $$t7.1
																					local.get $$t10.315
																					i32.store16 offset=630
																					local.get $$t7.1
																					local.get $$t10.316
																					i32.store16 offset=632
																					local.get $$t7.1
																					local.get $$t10.317
																					i32.store16 offset=634
																					local.get $$t7.1
																					local.get $$t10.318
																					i32.store16 offset=636
																					local.get $$t7.1
																					local.get $$t10.319
																					i32.store16 offset=638
																					local.get $$t7.1
																					local.get $$t10.320
																					i32.store16 offset=640
																					local.get $$t7.1
																					local.get $$t10.321
																					i32.store16 offset=642
																					local.get $$t7.1
																					local.get $$t10.322
																					i32.store16 offset=644
																					local.get $$t7.1
																					local.get $$t10.323
																					i32.store16 offset=646
																					local.get $$t7.1
																					local.get $$t10.324
																					i32.store16 offset=648
																					local.get $$t7.1
																					local.get $$t10.325
																					i32.store16 offset=650
																					local.get $$t7.1
																					local.get $$t10.326
																					i32.store16 offset=652
																					local.get $$t7.1
																					local.get $$t10.327
																					i32.store16 offset=654
																					local.get $$t7.1
																					local.get $$t10.328
																					i32.store16 offset=656
																					local.get $$t7.1
																					local.get $$t10.329
																					i32.store16 offset=658
																					local.get $$t7.1
																					local.get $$t10.330
																					i32.store16 offset=660
																					local.get $$t7.1
																					local.get $$t10.331
																					i32.store16 offset=662
																					local.get $$t7.1
																					local.get $$t10.332
																					i32.store16 offset=664
																					local.get $$t7.1
																					local.get $$t10.333
																					i32.store16 offset=666
																					local.get $$t7.1
																					local.get $$t10.334
																					i32.store16 offset=668
																					local.get $$t7.1
																					local.get $$t10.335
																					i32.store16 offset=670
																					local.get $$t7.1
																					local.get $$t10.336
																					i32.store16 offset=672
																					local.get $$t7.1
																					local.get $$t10.337
																					i32.store16 offset=674
																					local.get $$t7.1
																					local.get $$t10.338
																					i32.store16 offset=676
																					local.get $$t7.1
																					local.get $$t10.339
																					i32.store16 offset=678
																					local.get $$t7.1
																					local.get $$t10.340
																					i32.store16 offset=680
																					local.get $$t7.1
																					local.get $$t10.341
																					i32.store16 offset=682
																					local.get $$t7.1
																					local.get $$t10.342
																					i32.store16 offset=684
																					local.get $$t7.1
																					local.get $$t10.343
																					i32.store16 offset=686
																					local.get $$t7.1
																					local.get $$t10.344
																					i32.store16 offset=688
																					local.get $$t7.1
																					local.get $$t10.345
																					i32.store16 offset=690
																					local.get $$t7.1
																					local.get $$t10.346
																					i32.store16 offset=692
																					local.get $$t7.1
																					local.get $$t10.347
																					i32.store16 offset=694
																					local.get $$t7.1
																					local.get $$t10.348
																					i32.store16 offset=696
																					local.get $$t7.1
																					local.get $$t10.349
																					i32.store16 offset=698
																					local.get $$t7.1
																					local.get $$t10.350
																					i32.store16 offset=700
																					local.get $$t7.1
																					local.get $$t10.351
																					i32.store16 offset=702
																					local.get $$t7.1
																					local.get $$t10.352
																					i32.store16 offset=704
																					local.get $$t7.1
																					local.get $$t10.353
																					i32.store16 offset=706
																					local.get $$t7.1
																					local.get $$t10.354
																					i32.store16 offset=708
																					local.get $$t7.1
																					local.get $$t10.355
																					i32.store16 offset=710
																					local.get $$t7.1
																					local.get $$t10.356
																					i32.store16 offset=712
																					local.get $$t7.1
																					local.get $$t10.357
																					i32.store16 offset=714
																					local.get $$t7.1
																					local.get $$t10.358
																					i32.store16 offset=716
																					local.get $$t7.1
																					local.get $$t10.359
																					i32.store16 offset=718
																					local.get $$t7.1
																					local.get $$t10.360
																					i32.store16 offset=720
																					local.get $$t7.1
																					local.get $$t10.361
																					i32.store16 offset=722
																					local.get $$t7.1
																					local.get $$t10.362
																					i32.store16 offset=724
																					local.get $$t7.1
																					local.get $$t10.363
																					i32.store16 offset=726
																					local.get $$t7.1
																					local.get $$t10.364
																					i32.store16 offset=728
																					local.get $$t7.1
																					local.get $$t10.365
																					i32.store16 offset=730
																					local.get $$t7.1
																					local.get $$t10.366
																					i32.store16 offset=732
																					local.get $$t7.1
																					local.get $$t10.367
																					i32.store16 offset=734
																					local.get $$t7.1
																					local.get $$t10.368
																					i32.store16 offset=736
																					local.get $$t7.1
																					local.get $$t10.369
																					i32.store16 offset=738
																					local.get $$t7.1
																					local.get $$t10.370
																					i32.store16 offset=740
																					local.get $$t7.1
																					local.get $$t10.371
																					i32.store16 offset=742
																					local.get $$t7.1
																					local.get $$t10.372
																					i32.store16 offset=744
																					local.get $$t7.1
																					local.get $$t10.373
																					i32.store16 offset=746
																					local.get $$t7.1
																					local.get $$t10.374
																					i32.store16 offset=748
																					local.get $$t7.1
																					local.get $$t10.375
																					i32.store16 offset=750
																					local.get $$t7.1
																					local.get $$t10.376
																					i32.store16 offset=752
																					local.get $$t7.1
																					local.get $$t10.377
																					i32.store16 offset=754
																					local.get $$t7.1
																					local.get $$t10.378
																					i32.store16 offset=756
																					local.get $$t7.1
																					local.get $$t10.379
																					i32.store16 offset=758
																					local.get $$t7.1
																					local.get $$t10.380
																					i32.store16 offset=760
																					local.get $$t7.1
																					local.get $$t10.381
																					i32.store16 offset=762
																					local.get $$t7.1
																					local.get $$t10.382
																					i32.store16 offset=764
																					local.get $$t7.1
																					local.get $$t10.383
																					i32.store16 offset=766
																					local.get $$t7.1
																					local.get $$t10.384
																					i32.store16 offset=768
																					local.get $$t7.1
																					local.get $$t10.385
																					i32.store16 offset=770
																					local.get $$t7.1
																					local.get $$t10.386
																					i32.store16 offset=772
																					local.get $$t7.1
																					local.get $$t10.387
																					i32.store16 offset=774
																					local.get $$t7.1
																					local.get $$t10.388
																					i32.store16 offset=776
																					local.get $$t7.1
																					local.get $$t10.389
																					i32.store16 offset=778
																					local.get $$t7.1
																					local.get $$t10.390
																					i32.store16 offset=780
																					local.get $$t7.1
																					local.get $$t10.391
																					i32.store16 offset=782
																					local.get $$t7.1
																					local.get $$t10.392
																					i32.store16 offset=784
																					local.get $$t7.1
																					local.get $$t10.393
																					i32.store16 offset=786
																					local.get $$t7.1
																					local.get $$t10.394
																					i32.store16 offset=788
																					local.get $$t7.1
																					local.get $$t10.395
																					i32.store16 offset=790
																					local.get $$t7.1
																					local.get $$t10.396
																					i32.store16 offset=792
																					local.get $$t7.1
																					local.get $$t10.397
																					i32.store16 offset=794
																					local.get $$t7.1
																					local.get $$t10.398
																					i32.store16 offset=796
																					local.get $$t7.1
																					local.get $$t10.399
																					i32.store16 offset=798
																					local.get $$t7.1
																					local.get $$t10.400
																					i32.store16 offset=800
																					local.get $$t7.1
																					local.get $$t10.401
																					i32.store16 offset=802
																					local.get $$t7.1
																					local.get $$t10.402
																					i32.store16 offset=804
																					local.get $$t7.1
																					local.get $$t10.403
																					i32.store16 offset=806
																					local.get $$t7.1
																					local.get $$t10.404
																					i32.store16 offset=808
																					local.get $$t7.1
																					local.get $$t10.405
																					i32.store16 offset=810
																					local.get $$t7.1
																					local.get $$t10.406
																					i32.store16 offset=812
																					local.get $$t7.1
																					local.get $$t10.407
																					i32.store16 offset=814
																					local.get $$t7.1
																					local.get $$t10.408
																					i32.store16 offset=816
																					local.get $$t7.1
																					local.get $$t10.409
																					i32.store16 offset=818
																					local.get $$t7.1
																					local.get $$t10.410
																					i32.store16 offset=820
																					local.get $$t7.1
																					local.get $$t10.411
																					i32.store16 offset=822
																					local.get $$t7.1
																					local.get $$t10.412
																					i32.store16 offset=824
																					local.get $$t7.1
																					local.get $$t10.413
																					i32.store16 offset=826
																					local.get $$t7.1
																					local.get $$t10.414
																					i32.store16 offset=828
																					local.get $$t7.1
																					local.get $$t10.415
																					i32.store16 offset=830
																					local.get $$t7.1
																					local.get $$t10.416
																					i32.store16 offset=832
																					local.get $$t7.1
																					local.get $$t10.417
																					i32.store16 offset=834
																					local.get $$t7.1
																					local.get $$t10.418
																					i32.store16 offset=836
																					local.get $$t7.1
																					local.get $$t10.419
																					i32.store16 offset=838
																					local.get $$t7.1
																					local.get $$t10.420
																					i32.store16 offset=840
																					local.get $$t7.1
																					local.get $$t10.421
																					i32.store16 offset=842
																					local.get $$t7.1
																					local.get $$t10.422
																					i32.store16 offset=844
																					local.get $$t7.1
																					local.get $$t10.423
																					i32.store16 offset=846
																					local.get $$t7.1
																					local.get $$t10.424
																					i32.store16 offset=848
																					local.get $$t7.1
																					local.get $$t10.425
																					i32.store16 offset=850
																					local.get $$t7.1
																					local.get $$t10.426
																					i32.store16 offset=852
																					local.get $$t7.1
																					local.get $$t10.427
																					i32.store16 offset=854
																					local.get $$t7.1
																					local.get $$t10.428
																					i32.store16 offset=856
																					local.get $$t7.1
																					local.get $$t10.429
																					i32.store16 offset=858
																					local.get $$t7.1
																					local.get $$t10.430
																					i32.store16 offset=860
																					local.get $$t7.1
																					local.get $$t10.431
																					i32.store16 offset=862
																					local.get $$t7.1
																					local.get $$t10.432
																					i32.store16 offset=864
																					local.get $$t7.1
																					local.get $$t10.433
																					i32.store16 offset=866
																					local.get $$t8.1
																					local.get $$t11.0
																					i32.store16
																					local.get $$t8.1
																					local.get $$t11.1
																					i32.store16 offset=2
																					local.get $$t8.1
																					local.get $$t11.2
																					i32.store16 offset=4
																					local.get $$t8.1
																					local.get $$t11.3
																					i32.store16 offset=6
																					local.get $$t8.1
																					local.get $$t11.4
																					i32.store16 offset=8
																					local.get $$t8.1
																					local.get $$t11.5
																					i32.store16 offset=10
																					local.get $$t8.1
																					local.get $$t11.6
																					i32.store16 offset=12
																					local.get $$t8.1
																					local.get $$t11.7
																					i32.store16 offset=14
																					local.get $$t8.1
																					local.get $$t11.8
																					i32.store16 offset=16
																					local.get $$t8.1
																					local.get $$t11.9
																					i32.store16 offset=18
																					local.get $$t8.1
																					local.get $$t11.10
																					i32.store16 offset=20
																					local.get $$t8.1
																					local.get $$t11.11
																					i32.store16 offset=22
																					local.get $$t8.1
																					local.get $$t11.12
																					i32.store16 offset=24
																					local.get $$t8.1
																					local.get $$t11.13
																					i32.store16 offset=26
																					local.get $$t8.1
																					local.get $$t11.14
																					i32.store16 offset=28
																					local.get $$t8.1
																					local.get $$t11.15
																					i32.store16 offset=30
																					local.get $$t8.1
																					local.get $$t11.16
																					i32.store16 offset=32
																					local.get $$t8.1
																					local.get $$t11.17
																					i32.store16 offset=34
																					local.get $$t8.1
																					local.get $$t11.18
																					i32.store16 offset=36
																					local.get $$t8.1
																					local.get $$t11.19
																					i32.store16 offset=38
																					local.get $$t8.1
																					local.get $$t11.20
																					i32.store16 offset=40
																					local.get $$t8.1
																					local.get $$t11.21
																					i32.store16 offset=42
																					local.get $$t8.1
																					local.get $$t11.22
																					i32.store16 offset=44
																					local.get $$t8.1
																					local.get $$t11.23
																					i32.store16 offset=46
																					local.get $$t8.1
																					local.get $$t11.24
																					i32.store16 offset=48
																					local.get $$t8.1
																					local.get $$t11.25
																					i32.store16 offset=50
																					local.get $$t8.1
																					local.get $$t11.26
																					i32.store16 offset=52
																					local.get $$t8.1
																					local.get $$t11.27
																					i32.store16 offset=54
																					local.get $$t8.1
																					local.get $$t11.28
																					i32.store16 offset=56
																					local.get $$t8.1
																					local.get $$t11.29
																					i32.store16 offset=58
																					local.get $$t8.1
																					local.get $$t11.30
																					i32.store16 offset=60
																					local.get $$t8.1
																					local.get $$t11.31
																					i32.store16 offset=62
																					local.get $$t8.1
																					local.get $$t11.32
																					i32.store16 offset=64
																					local.get $$t8.1
																					local.get $$t11.33
																					i32.store16 offset=66
																					local.get $$t8.1
																					local.get $$t11.34
																					i32.store16 offset=68
																					local.get $$t8.1
																					local.get $$t11.35
																					i32.store16 offset=70
																					local.get $$t8.1
																					local.get $$t11.36
																					i32.store16 offset=72
																					local.get $$t8.1
																					local.get $$t11.37
																					i32.store16 offset=74
																					local.get $$t8.1
																					local.get $$t11.38
																					i32.store16 offset=76
																					local.get $$t8.1
																					local.get $$t11.39
																					i32.store16 offset=78
																					local.get $$t8.1
																					local.get $$t11.40
																					i32.store16 offset=80
																					local.get $$t8.1
																					local.get $$t11.41
																					i32.store16 offset=82
																					local.get $$t8.1
																					local.get $$t11.42
																					i32.store16 offset=84
																					local.get $$t8.1
																					local.get $$t11.43
																					i32.store16 offset=86
																					local.get $$t8.1
																					local.get $$t11.44
																					i32.store16 offset=88
																					local.get $$t8.1
																					local.get $$t11.45
																					i32.store16 offset=90
																					local.get $$t8.1
																					local.get $$t11.46
																					i32.store16 offset=92
																					local.get $$t8.1
																					local.get $$t11.47
																					i32.store16 offset=94
																					local.get $$t8.1
																					local.get $$t11.48
																					i32.store16 offset=96
																					local.get $$t8.1
																					local.get $$t11.49
																					i32.store16 offset=98
																					local.get $$t8.1
																					local.get $$t11.50
																					i32.store16 offset=100
																					local.get $$t8.1
																					local.get $$t11.51
																					i32.store16 offset=102
																					local.get $$t8.1
																					local.get $$t11.52
																					i32.store16 offset=104
																					local.get $$t8.1
																					local.get $$t11.53
																					i32.store16 offset=106
																					local.get $$t8.1
																					local.get $$t11.54
																					i32.store16 offset=108
																					local.get $$t8.1
																					local.get $$t11.55
																					i32.store16 offset=110
																					local.get $$t8.1
																					local.get $$t11.56
																					i32.store16 offset=112
																					local.get $$t8.1
																					local.get $$t11.57
																					i32.store16 offset=114
																					local.get $$t8.1
																					local.get $$t11.58
																					i32.store16 offset=116
																					local.get $$t8.1
																					local.get $$t11.59
																					i32.store16 offset=118
																					local.get $$t8.1
																					local.get $$t11.60
																					i32.store16 offset=120
																					local.get $$t8.1
																					local.get $$t11.61
																					i32.store16 offset=122
																					local.get $$t8.1
																					local.get $$t11.62
																					i32.store16 offset=124
																					local.get $$t8.1
																					local.get $$t11.63
																					i32.store16 offset=126
																					local.get $$t8.1
																					local.get $$t11.64
																					i32.store16 offset=128
																					local.get $$t8.1
																					local.get $$t11.65
																					i32.store16 offset=130
																					local.get $$t8.1
																					local.get $$t11.66
																					i32.store16 offset=132
																					local.get $$t8.1
																					local.get $$t11.67
																					i32.store16 offset=134
																					local.get $$t8.1
																					local.get $$t11.68
																					i32.store16 offset=136
																					local.get $$t8.1
																					local.get $$t11.69
																					i32.store16 offset=138
																					local.get $$t8.1
																					local.get $$t11.70
																					i32.store16 offset=140
																					local.get $$t8.1
																					local.get $$t11.71
																					i32.store16 offset=142
																					local.get $$t8.1
																					local.get $$t11.72
																					i32.store16 offset=144
																					local.get $$t8.1
																					local.get $$t11.73
																					i32.store16 offset=146
																					local.get $$t8.1
																					local.get $$t11.74
																					i32.store16 offset=148
																					local.get $$t8.1
																					local.get $$t11.75
																					i32.store16 offset=150
																					local.get $$t8.1
																					local.get $$t11.76
																					i32.store16 offset=152
																					local.get $$t8.1
																					local.get $$t11.77
																					i32.store16 offset=154
																					local.get $$t8.1
																					local.get $$t11.78
																					i32.store16 offset=156
																					local.get $$t8.1
																					local.get $$t11.79
																					i32.store16 offset=158
																					local.get $$t8.1
																					local.get $$t11.80
																					i32.store16 offset=160
																					local.get $$t8.1
																					local.get $$t11.81
																					i32.store16 offset=162
																					local.get $$t8.1
																					local.get $$t11.82
																					i32.store16 offset=164
																					local.get $$t8.1
																					local.get $$t11.83
																					i32.store16 offset=166
																					local.get $$t8.1
																					local.get $$t11.84
																					i32.store16 offset=168
																					local.get $$t8.1
																					local.get $$t11.85
																					i32.store16 offset=170
																					local.get $$t8.1
																					local.get $$t11.86
																					i32.store16 offset=172
																					local.get $$t8.1
																					local.get $$t11.87
																					i32.store16 offset=174
																					local.get $$t8.1
																					local.get $$t11.88
																					i32.store16 offset=176
																					local.get $$t8.1
																					local.get $$t11.89
																					i32.store16 offset=178
																					local.get $$t8.1
																					local.get $$t11.90
																					i32.store16 offset=180
																					local.get $$t8.1
																					local.get $$t11.91
																					i32.store16 offset=182
																					local.get $$t8.1
																					local.get $$t11.92
																					i32.store16 offset=184
																					local.get $$t8.1
																					local.get $$t11.93
																					i32.store16 offset=186
																					local.get $$t8.1
																					local.get $$t11.94
																					i32.store16 offset=188
																					local.get $$t8.1
																					local.get $$t11.95
																					i32.store16 offset=190
																					local.get $$t8.1
																					local.get $$t11.96
																					i32.store16 offset=192
																					local.get $$t8.1
																					local.get $$t11.97
																					i32.store16 offset=194
																					local.get $$t8.1
																					local.get $$t11.98
																					i32.store16 offset=196
																					local.get $$t8.1
																					local.get $$t11.99
																					i32.store16 offset=198
																					local.get $$t8.1
																					local.get $$t11.100
																					i32.store16 offset=200
																					local.get $$t8.1
																					local.get $$t11.101
																					i32.store16 offset=202
																					local.get $$t8.1
																					local.get $$t11.102
																					i32.store16 offset=204
																					local.get $$t8.1
																					local.get $$t11.103
																					i32.store16 offset=206
																					local.get $$t8.1
																					local.get $$t11.104
																					i32.store16 offset=208
																					local.get $$t8.1
																					local.get $$t11.105
																					i32.store16 offset=210
																					local.get $$t8.1
																					local.get $$t11.106
																					i32.store16 offset=212
																					local.get $$t8.1
																					local.get $$t11.107
																					i32.store16 offset=214
																					local.get $$t8.1
																					local.get $$t11.108
																					i32.store16 offset=216
																					local.get $$t8.1
																					local.get $$t11.109
																					i32.store16 offset=218
																					local.get $$t8.1
																					local.get $$t11.110
																					i32.store16 offset=220
																					local.get $$t8.1
																					local.get $$t11.111
																					i32.store16 offset=222
																					local.get $$t8.1
																					local.get $$t11.112
																					i32.store16 offset=224
																					local.get $$t8.1
																					local.get $$t11.113
																					i32.store16 offset=226
																					local.get $$t8.1
																					local.get $$t11.114
																					i32.store16 offset=228
																					local.get $$t8.1
																					local.get $$t11.115
																					i32.store16 offset=230
																					local.get $$t8.1
																					local.get $$t11.116
																					i32.store16 offset=232
																					local.get $$t8.1
																					local.get $$t11.117
																					i32.store16 offset=234
																					local.get $$t8.1
																					local.get $$t11.118
																					i32.store16 offset=236
																					local.get $$t8.1
																					local.get $$t11.119
																					i32.store16 offset=238
																					local.get $$t8.1
																					local.get $$t11.120
																					i32.store16 offset=240
																					local.get $$t8.1
																					local.get $$t11.121
																					i32.store16 offset=242
																					local.get $$t8.1
																					local.get $$t11.122
																					i32.store16 offset=244
																					local.get $$t8.1
																					local.get $$t11.123
																					i32.store16 offset=246
																					local.get $$t8.1
																					local.get $$t11.124
																					i32.store16 offset=248
																					local.get $$t8.1
																					local.get $$t11.125
																					i32.store16 offset=250
																					local.get $$t8.1
																					local.get $$t11.126
																					i32.store16 offset=252
																					local.get $$t8.1
																					local.get $$t11.127
																					i32.store16 offset=254
																					local.get $$t8.1
																					local.get $$t11.128
																					i32.store16 offset=256
																					local.get $$t8.1
																					local.get $$t11.129
																					i32.store16 offset=258
																					local.get $$t8.1
																					local.get $$t11.130
																					i32.store16 offset=260
																					local.get $$t8.1
																					local.get $$t11.131
																					i32.store16 offset=262
																					local.get $$t7.0
																					call $runtime.Block.Retain
																					local.get $$t7.1
																					i32.const 2
																					i32.const 0
																					i32.mul
																					i32.add
																					i32.const 434
																					i32.const 0
																					i32.sub
																					i32.const 434
																					i32.const 0
																					i32.sub
																					local.set $$t12.3
																					local.set $$t12.2
																					local.set $$t12.1
																					local.get $$t12.0
																					call $runtime.Block.Release
																					local.set $$t12.0
																					local.get $$t12.0
																					local.get $$t12.1
																					local.get $$t12.2
																					local.get $$t12.3
																					local.get $$t9
																					call $strconv.bsearch16
																					local.set $$t13
																					local.get $$t13
																					i32.const 434
																					i32.ge_s
																					local.set $$t14
																					local.get $$t14
																					if
																						br $$Block_11
																					else
																						br $$Block_14
																					end
																				end
																				i32.const 10
																				local.set $$current_block
																				i32.const 1888
																				call $runtime.HeapAlloc
																				i32.const 1
																				i32.const 0
																				i32.const 1872
																				call $runtime.Block.Init
																				call $runtime.DupI32
																				i32.const 16
																				i32.add
																				local.set $$t15.1
																				local.get $$t15.0
																				call $runtime.Block.Release
																				local.set $$t15.0
																				i32.const 206
																				call $runtime.HeapAlloc
																				i32.const 1
																				i32.const 0
																				i32.const 190
																				call $runtime.Block.Init
																				call $runtime.DupI32
																				i32.const 16
																				i32.add
																				local.set $$t16.1
																				local.get $$t16.0
																				call $runtime.Block.Release
																				local.set $$t16.0
																				local.get $r
																				local.set $$t17
																				i32.const 27780
																				i32.load
																				i32.const 27780
																				i32.load offset=4
																				i32.const 27780
																				i32.load offset=8
																				i32.const 27780
																				i32.load offset=12
																				i32.const 27780
																				i32.load offset=16
																				i32.const 27780
																				i32.load offset=20
																				i32.const 27780
																				i32.load offset=24
																				i32.const 27780
																				i32.load offset=28
																				i32.const 27780
																				i32.load offset=32
																				i32.const 27780
																				i32.load offset=36
																				i32.const 27780
																				i32.load offset=40
																				i32.const 27780
																				i32.load offset=44
																				i32.const 27780
																				i32.load offset=48
																				i32.const 27780
																				i32.load offset=52
																				i32.const 27780
																				i32.load offset=56
																				i32.const 27780
																				i32.load offset=60
																				i32.const 27780
																				i32.load offset=64
																				i32.const 27780
																				i32.load offset=68
																				i32.const 27780
																				i32.load offset=72
																				i32.const 27780
																				i32.load offset=76
																				i32.const 27780
																				i32.load offset=80
																				i32.const 27780
																				i32.load offset=84
																				i32.const 27780
																				i32.load offset=88
																				i32.const 27780
																				i32.load offset=92
																				i32.const 27780
																				i32.load offset=96
																				i32.const 27780
																				i32.load offset=100
																				i32.const 27780
																				i32.load offset=104
																				i32.const 27780
																				i32.load offset=108
																				i32.const 27780
																				i32.load offset=112
																				i32.const 27780
																				i32.load offset=116
																				i32.const 27780
																				i32.load offset=120
																				i32.const 27780
																				i32.load offset=124
																				i32.const 27780
																				i32.load offset=128
																				i32.const 27780
																				i32.load offset=132
																				i32.const 27780
																				i32.load offset=136
																				i32.const 27780
																				i32.load offset=140
																				i32.const 27780
																				i32.load offset=144
																				i32.const 27780
																				i32.load offset=148
																				i32.const 27780
																				i32.load offset=152
																				i32.const 27780
																				i32.load offset=156
																				i32.const 27780
																				i32.load offset=160
																				i32.const 27780
																				i32.load offset=164
																				i32.const 27780
																				i32.load offset=168
																				i32.const 27780
																				i32.load offset=172
																				i32.const 27780
																				i32.load offset=176
																				i32.const 27780
																				i32.load offset=180
																				i32.const 27780
																				i32.load offset=184
																				i32.const 27780
																				i32.load offset=188
																				i32.const 27780
																				i32.load offset=192
																				i32.const 27780
																				i32.load offset=196
																				i32.const 27780
																				i32.load offset=200
																				i32.const 27780
																				i32.load offset=204
																				i32.const 27780
																				i32.load offset=208
																				i32.const 27780
																				i32.load offset=212
																				i32.const 27780
																				i32.load offset=216
																				i32.const 27780
																				i32.load offset=220
																				i32.const 27780
																				i32.load offset=224
																				i32.const 27780
																				i32.load offset=228
																				i32.const 27780
																				i32.load offset=232
																				i32.const 27780
																				i32.load offset=236
																				i32.const 27780
																				i32.load offset=240
																				i32.const 27780
																				i32.load offset=244
																				i32.const 27780
																				i32.load offset=248
																				i32.const 27780
																				i32.load offset=252
																				i32.const 27780
																				i32.load offset=256
																				i32.const 27780
																				i32.load offset=260
																				i32.const 27780
																				i32.load offset=264
																				i32.const 27780
																				i32.load offset=268
																				i32.const 27780
																				i32.load offset=272
																				i32.const 27780
																				i32.load offset=276
																				i32.const 27780
																				i32.load offset=280
																				i32.const 27780
																				i32.load offset=284
																				i32.const 27780
																				i32.load offset=288
																				i32.const 27780
																				i32.load offset=292
																				i32.const 27780
																				i32.load offset=296
																				i32.const 27780
																				i32.load offset=300
																				i32.const 27780
																				i32.load offset=304
																				i32.const 27780
																				i32.load offset=308
																				i32.const 27780
																				i32.load offset=312
																				i32.const 27780
																				i32.load offset=316
																				i32.const 27780
																				i32.load offset=320
																				i32.const 27780
																				i32.load offset=324
																				i32.const 27780
																				i32.load offset=328
																				i32.const 27780
																				i32.load offset=332
																				i32.const 27780
																				i32.load offset=336
																				i32.const 27780
																				i32.load offset=340
																				i32.const 27780
																				i32.load offset=344
																				i32.const 27780
																				i32.load offset=348
																				i32.const 27780
																				i32.load offset=352
																				i32.const 27780
																				i32.load offset=356
																				i32.const 27780
																				i32.load offset=360
																				i32.const 27780
																				i32.load offset=364
																				i32.const 27780
																				i32.load offset=368
																				i32.const 27780
																				i32.load offset=372
																				i32.const 27780
																				i32.load offset=376
																				i32.const 27780
																				i32.load offset=380
																				i32.const 27780
																				i32.load offset=384
																				i32.const 27780
																				i32.load offset=388
																				i32.const 27780
																				i32.load offset=392
																				i32.const 27780
																				i32.load offset=396
																				i32.const 27780
																				i32.load offset=400
																				i32.const 27780
																				i32.load offset=404
																				i32.const 27780
																				i32.load offset=408
																				i32.const 27780
																				i32.load offset=412
																				i32.const 27780
																				i32.load offset=416
																				i32.const 27780
																				i32.load offset=420
																				i32.const 27780
																				i32.load offset=424
																				i32.const 27780
																				i32.load offset=428
																				i32.const 27780
																				i32.load offset=432
																				i32.const 27780
																				i32.load offset=436
																				i32.const 27780
																				i32.load offset=440
																				i32.const 27780
																				i32.load offset=444
																				i32.const 27780
																				i32.load offset=448
																				i32.const 27780
																				i32.load offset=452
																				i32.const 27780
																				i32.load offset=456
																				i32.const 27780
																				i32.load offset=460
																				i32.const 27780
																				i32.load offset=464
																				i32.const 27780
																				i32.load offset=468
																				i32.const 27780
																				i32.load offset=472
																				i32.const 27780
																				i32.load offset=476
																				i32.const 27780
																				i32.load offset=480
																				i32.const 27780
																				i32.load offset=484
																				i32.const 27780
																				i32.load offset=488
																				i32.const 27780
																				i32.load offset=492
																				i32.const 27780
																				i32.load offset=496
																				i32.const 27780
																				i32.load offset=500
																				i32.const 27780
																				i32.load offset=504
																				i32.const 27780
																				i32.load offset=508
																				i32.const 27780
																				i32.load offset=512
																				i32.const 27780
																				i32.load offset=516
																				i32.const 27780
																				i32.load offset=520
																				i32.const 27780
																				i32.load offset=524
																				i32.const 27780
																				i32.load offset=528
																				i32.const 27780
																				i32.load offset=532
																				i32.const 27780
																				i32.load offset=536
																				i32.const 27780
																				i32.load offset=540
																				i32.const 27780
																				i32.load offset=544
																				i32.const 27780
																				i32.load offset=548
																				i32.const 27780
																				i32.load offset=552
																				i32.const 27780
																				i32.load offset=556
																				i32.const 27780
																				i32.load offset=560
																				i32.const 27780
																				i32.load offset=564
																				i32.const 27780
																				i32.load offset=568
																				i32.const 27780
																				i32.load offset=572
																				i32.const 27780
																				i32.load offset=576
																				i32.const 27780
																				i32.load offset=580
																				i32.const 27780
																				i32.load offset=584
																				i32.const 27780
																				i32.load offset=588
																				i32.const 27780
																				i32.load offset=592
																				i32.const 27780
																				i32.load offset=596
																				i32.const 27780
																				i32.load offset=600
																				i32.const 27780
																				i32.load offset=604
																				i32.const 27780
																				i32.load offset=608
																				i32.const 27780
																				i32.load offset=612
																				i32.const 27780
																				i32.load offset=616
																				i32.const 27780
																				i32.load offset=620
																				i32.const 27780
																				i32.load offset=624
																				i32.const 27780
																				i32.load offset=628
																				i32.const 27780
																				i32.load offset=632
																				i32.const 27780
																				i32.load offset=636
																				i32.const 27780
																				i32.load offset=640
																				i32.const 27780
																				i32.load offset=644
																				i32.const 27780
																				i32.load offset=648
																				i32.const 27780
																				i32.load offset=652
																				i32.const 27780
																				i32.load offset=656
																				i32.const 27780
																				i32.load offset=660
																				i32.const 27780
																				i32.load offset=664
																				i32.const 27780
																				i32.load offset=668
																				i32.const 27780
																				i32.load offset=672
																				i32.const 27780
																				i32.load offset=676
																				i32.const 27780
																				i32.load offset=680
																				i32.const 27780
																				i32.load offset=684
																				i32.const 27780
																				i32.load offset=688
																				i32.const 27780
																				i32.load offset=692
																				i32.const 27780
																				i32.load offset=696
																				i32.const 27780
																				i32.load offset=700
																				i32.const 27780
																				i32.load offset=704
																				i32.const 27780
																				i32.load offset=708
																				i32.const 27780
																				i32.load offset=712
																				i32.const 27780
																				i32.load offset=716
																				i32.const 27780
																				i32.load offset=720
																				i32.const 27780
																				i32.load offset=724
																				i32.const 27780
																				i32.load offset=728
																				i32.const 27780
																				i32.load offset=732
																				i32.const 27780
																				i32.load offset=736
																				i32.const 27780
																				i32.load offset=740
																				i32.const 27780
																				i32.load offset=744
																				i32.const 27780
																				i32.load offset=748
																				i32.const 27780
																				i32.load offset=752
																				i32.const 27780
																				i32.load offset=756
																				i32.const 27780
																				i32.load offset=760
																				i32.const 27780
																				i32.load offset=764
																				i32.const 27780
																				i32.load offset=768
																				i32.const 27780
																				i32.load offset=772
																				i32.const 27780
																				i32.load offset=776
																				i32.const 27780
																				i32.load offset=780
																				i32.const 27780
																				i32.load offset=784
																				i32.const 27780
																				i32.load offset=788
																				i32.const 27780
																				i32.load offset=792
																				i32.const 27780
																				i32.load offset=796
																				i32.const 27780
																				i32.load offset=800
																				i32.const 27780
																				i32.load offset=804
																				i32.const 27780
																				i32.load offset=808
																				i32.const 27780
																				i32.load offset=812
																				i32.const 27780
																				i32.load offset=816
																				i32.const 27780
																				i32.load offset=820
																				i32.const 27780
																				i32.load offset=824
																				i32.const 27780
																				i32.load offset=828
																				i32.const 27780
																				i32.load offset=832
																				i32.const 27780
																				i32.load offset=836
																				i32.const 27780
																				i32.load offset=840
																				i32.const 27780
																				i32.load offset=844
																				i32.const 27780
																				i32.load offset=848
																				i32.const 27780
																				i32.load offset=852
																				i32.const 27780
																				i32.load offset=856
																				i32.const 27780
																				i32.load offset=860
																				i32.const 27780
																				i32.load offset=864
																				i32.const 27780
																				i32.load offset=868
																				i32.const 27780
																				i32.load offset=872
																				i32.const 27780
																				i32.load offset=876
																				i32.const 27780
																				i32.load offset=880
																				i32.const 27780
																				i32.load offset=884
																				i32.const 27780
																				i32.load offset=888
																				i32.const 27780
																				i32.load offset=892
																				i32.const 27780
																				i32.load offset=896
																				i32.const 27780
																				i32.load offset=900
																				i32.const 27780
																				i32.load offset=904
																				i32.const 27780
																				i32.load offset=908
																				i32.const 27780
																				i32.load offset=912
																				i32.const 27780
																				i32.load offset=916
																				i32.const 27780
																				i32.load offset=920
																				i32.const 27780
																				i32.load offset=924
																				i32.const 27780
																				i32.load offset=928
																				i32.const 27780
																				i32.load offset=932
																				i32.const 27780
																				i32.load offset=936
																				i32.const 27780
																				i32.load offset=940
																				i32.const 27780
																				i32.load offset=944
																				i32.const 27780
																				i32.load offset=948
																				i32.const 27780
																				i32.load offset=952
																				i32.const 27780
																				i32.load offset=956
																				i32.const 27780
																				i32.load offset=960
																				i32.const 27780
																				i32.load offset=964
																				i32.const 27780
																				i32.load offset=968
																				i32.const 27780
																				i32.load offset=972
																				i32.const 27780
																				i32.load offset=976
																				i32.const 27780
																				i32.load offset=980
																				i32.const 27780
																				i32.load offset=984
																				i32.const 27780
																				i32.load offset=988
																				i32.const 27780
																				i32.load offset=992
																				i32.const 27780
																				i32.load offset=996
																				i32.const 27780
																				i32.load offset=1000
																				i32.const 27780
																				i32.load offset=1004
																				i32.const 27780
																				i32.load offset=1008
																				i32.const 27780
																				i32.load offset=1012
																				i32.const 27780
																				i32.load offset=1016
																				i32.const 27780
																				i32.load offset=1020
																				i32.const 27780
																				i32.load offset=1024
																				i32.const 27780
																				i32.load offset=1028
																				i32.const 27780
																				i32.load offset=1032
																				i32.const 27780
																				i32.load offset=1036
																				i32.const 27780
																				i32.load offset=1040
																				i32.const 27780
																				i32.load offset=1044
																				i32.const 27780
																				i32.load offset=1048
																				i32.const 27780
																				i32.load offset=1052
																				i32.const 27780
																				i32.load offset=1056
																				i32.const 27780
																				i32.load offset=1060
																				i32.const 27780
																				i32.load offset=1064
																				i32.const 27780
																				i32.load offset=1068
																				i32.const 27780
																				i32.load offset=1072
																				i32.const 27780
																				i32.load offset=1076
																				i32.const 27780
																				i32.load offset=1080
																				i32.const 27780
																				i32.load offset=1084
																				i32.const 27780
																				i32.load offset=1088
																				i32.const 27780
																				i32.load offset=1092
																				i32.const 27780
																				i32.load offset=1096
																				i32.const 27780
																				i32.load offset=1100
																				i32.const 27780
																				i32.load offset=1104
																				i32.const 27780
																				i32.load offset=1108
																				i32.const 27780
																				i32.load offset=1112
																				i32.const 27780
																				i32.load offset=1116
																				i32.const 27780
																				i32.load offset=1120
																				i32.const 27780
																				i32.load offset=1124
																				i32.const 27780
																				i32.load offset=1128
																				i32.const 27780
																				i32.load offset=1132
																				i32.const 27780
																				i32.load offset=1136
																				i32.const 27780
																				i32.load offset=1140
																				i32.const 27780
																				i32.load offset=1144
																				i32.const 27780
																				i32.load offset=1148
																				i32.const 27780
																				i32.load offset=1152
																				i32.const 27780
																				i32.load offset=1156
																				i32.const 27780
																				i32.load offset=1160
																				i32.const 27780
																				i32.load offset=1164
																				i32.const 27780
																				i32.load offset=1168
																				i32.const 27780
																				i32.load offset=1172
																				i32.const 27780
																				i32.load offset=1176
																				i32.const 27780
																				i32.load offset=1180
																				i32.const 27780
																				i32.load offset=1184
																				i32.const 27780
																				i32.load offset=1188
																				i32.const 27780
																				i32.load offset=1192
																				i32.const 27780
																				i32.load offset=1196
																				i32.const 27780
																				i32.load offset=1200
																				i32.const 27780
																				i32.load offset=1204
																				i32.const 27780
																				i32.load offset=1208
																				i32.const 27780
																				i32.load offset=1212
																				i32.const 27780
																				i32.load offset=1216
																				i32.const 27780
																				i32.load offset=1220
																				i32.const 27780
																				i32.load offset=1224
																				i32.const 27780
																				i32.load offset=1228
																				i32.const 27780
																				i32.load offset=1232
																				i32.const 27780
																				i32.load offset=1236
																				i32.const 27780
																				i32.load offset=1240
																				i32.const 27780
																				i32.load offset=1244
																				i32.const 27780
																				i32.load offset=1248
																				i32.const 27780
																				i32.load offset=1252
																				i32.const 27780
																				i32.load offset=1256
																				i32.const 27780
																				i32.load offset=1260
																				i32.const 27780
																				i32.load offset=1264
																				i32.const 27780
																				i32.load offset=1268
																				i32.const 27780
																				i32.load offset=1272
																				i32.const 27780
																				i32.load offset=1276
																				i32.const 27780
																				i32.load offset=1280
																				i32.const 27780
																				i32.load offset=1284
																				i32.const 27780
																				i32.load offset=1288
																				i32.const 27780
																				i32.load offset=1292
																				i32.const 27780
																				i32.load offset=1296
																				i32.const 27780
																				i32.load offset=1300
																				i32.const 27780
																				i32.load offset=1304
																				i32.const 27780
																				i32.load offset=1308
																				i32.const 27780
																				i32.load offset=1312
																				i32.const 27780
																				i32.load offset=1316
																				i32.const 27780
																				i32.load offset=1320
																				i32.const 27780
																				i32.load offset=1324
																				i32.const 27780
																				i32.load offset=1328
																				i32.const 27780
																				i32.load offset=1332
																				i32.const 27780
																				i32.load offset=1336
																				i32.const 27780
																				i32.load offset=1340
																				i32.const 27780
																				i32.load offset=1344
																				i32.const 27780
																				i32.load offset=1348
																				i32.const 27780
																				i32.load offset=1352
																				i32.const 27780
																				i32.load offset=1356
																				i32.const 27780
																				i32.load offset=1360
																				i32.const 27780
																				i32.load offset=1364
																				i32.const 27780
																				i32.load offset=1368
																				i32.const 27780
																				i32.load offset=1372
																				i32.const 27780
																				i32.load offset=1376
																				i32.const 27780
																				i32.load offset=1380
																				i32.const 27780
																				i32.load offset=1384
																				i32.const 27780
																				i32.load offset=1388
																				i32.const 27780
																				i32.load offset=1392
																				i32.const 27780
																				i32.load offset=1396
																				i32.const 27780
																				i32.load offset=1400
																				i32.const 27780
																				i32.load offset=1404
																				i32.const 27780
																				i32.load offset=1408
																				i32.const 27780
																				i32.load offset=1412
																				i32.const 27780
																				i32.load offset=1416
																				i32.const 27780
																				i32.load offset=1420
																				i32.const 27780
																				i32.load offset=1424
																				i32.const 27780
																				i32.load offset=1428
																				i32.const 27780
																				i32.load offset=1432
																				i32.const 27780
																				i32.load offset=1436
																				i32.const 27780
																				i32.load offset=1440
																				i32.const 27780
																				i32.load offset=1444
																				i32.const 27780
																				i32.load offset=1448
																				i32.const 27780
																				i32.load offset=1452
																				i32.const 27780
																				i32.load offset=1456
																				i32.const 27780
																				i32.load offset=1460
																				i32.const 27780
																				i32.load offset=1464
																				i32.const 27780
																				i32.load offset=1468
																				i32.const 27780
																				i32.load offset=1472
																				i32.const 27780
																				i32.load offset=1476
																				i32.const 27780
																				i32.load offset=1480
																				i32.const 27780
																				i32.load offset=1484
																				i32.const 27780
																				i32.load offset=1488
																				i32.const 27780
																				i32.load offset=1492
																				i32.const 27780
																				i32.load offset=1496
																				i32.const 27780
																				i32.load offset=1500
																				i32.const 27780
																				i32.load offset=1504
																				i32.const 27780
																				i32.load offset=1508
																				i32.const 27780
																				i32.load offset=1512
																				i32.const 27780
																				i32.load offset=1516
																				i32.const 27780
																				i32.load offset=1520
																				i32.const 27780
																				i32.load offset=1524
																				i32.const 27780
																				i32.load offset=1528
																				i32.const 27780
																				i32.load offset=1532
																				i32.const 27780
																				i32.load offset=1536
																				i32.const 27780
																				i32.load offset=1540
																				i32.const 27780
																				i32.load offset=1544
																				i32.const 27780
																				i32.load offset=1548
																				i32.const 27780
																				i32.load offset=1552
																				i32.const 27780
																				i32.load offset=1556
																				i32.const 27780
																				i32.load offset=1560
																				i32.const 27780
																				i32.load offset=1564
																				i32.const 27780
																				i32.load offset=1568
																				i32.const 27780
																				i32.load offset=1572
																				i32.const 27780
																				i32.load offset=1576
																				i32.const 27780
																				i32.load offset=1580
																				i32.const 27780
																				i32.load offset=1584
																				i32.const 27780
																				i32.load offset=1588
																				i32.const 27780
																				i32.load offset=1592
																				i32.const 27780
																				i32.load offset=1596
																				i32.const 27780
																				i32.load offset=1600
																				i32.const 27780
																				i32.load offset=1604
																				i32.const 27780
																				i32.load offset=1608
																				i32.const 27780
																				i32.load offset=1612
																				i32.const 27780
																				i32.load offset=1616
																				i32.const 27780
																				i32.load offset=1620
																				i32.const 27780
																				i32.load offset=1624
																				i32.const 27780
																				i32.load offset=1628
																				i32.const 27780
																				i32.load offset=1632
																				i32.const 27780
																				i32.load offset=1636
																				i32.const 27780
																				i32.load offset=1640
																				i32.const 27780
																				i32.load offset=1644
																				i32.const 27780
																				i32.load offset=1648
																				i32.const 27780
																				i32.load offset=1652
																				i32.const 27780
																				i32.load offset=1656
																				i32.const 27780
																				i32.load offset=1660
																				i32.const 27780
																				i32.load offset=1664
																				i32.const 27780
																				i32.load offset=1668
																				i32.const 27780
																				i32.load offset=1672
																				i32.const 27780
																				i32.load offset=1676
																				i32.const 27780
																				i32.load offset=1680
																				i32.const 27780
																				i32.load offset=1684
																				i32.const 27780
																				i32.load offset=1688
																				i32.const 27780
																				i32.load offset=1692
																				i32.const 27780
																				i32.load offset=1696
																				i32.const 27780
																				i32.load offset=1700
																				i32.const 27780
																				i32.load offset=1704
																				i32.const 27780
																				i32.load offset=1708
																				i32.const 27780
																				i32.load offset=1712
																				i32.const 27780
																				i32.load offset=1716
																				i32.const 27780
																				i32.load offset=1720
																				i32.const 27780
																				i32.load offset=1724
																				i32.const 27780
																				i32.load offset=1728
																				i32.const 27780
																				i32.load offset=1732
																				i32.const 27780
																				i32.load offset=1736
																				i32.const 27780
																				i32.load offset=1740
																				i32.const 27780
																				i32.load offset=1744
																				i32.const 27780
																				i32.load offset=1748
																				i32.const 27780
																				i32.load offset=1752
																				i32.const 27780
																				i32.load offset=1756
																				i32.const 27780
																				i32.load offset=1760
																				i32.const 27780
																				i32.load offset=1764
																				i32.const 27780
																				i32.load offset=1768
																				i32.const 27780
																				i32.load offset=1772
																				i32.const 27780
																				i32.load offset=1776
																				i32.const 27780
																				i32.load offset=1780
																				i32.const 27780
																				i32.load offset=1784
																				i32.const 27780
																				i32.load offset=1788
																				i32.const 27780
																				i32.load offset=1792
																				i32.const 27780
																				i32.load offset=1796
																				i32.const 27780
																				i32.load offset=1800
																				i32.const 27780
																				i32.load offset=1804
																				i32.const 27780
																				i32.load offset=1808
																				i32.const 27780
																				i32.load offset=1812
																				i32.const 27780
																				i32.load offset=1816
																				i32.const 27780
																				i32.load offset=1820
																				i32.const 27780
																				i32.load offset=1824
																				i32.const 27780
																				i32.load offset=1828
																				i32.const 27780
																				i32.load offset=1832
																				i32.const 27780
																				i32.load offset=1836
																				i32.const 27780
																				i32.load offset=1840
																				i32.const 27780
																				i32.load offset=1844
																				i32.const 27780
																				i32.load offset=1848
																				i32.const 27780
																				i32.load offset=1852
																				i32.const 27780
																				i32.load offset=1856
																				i32.const 27780
																				i32.load offset=1860
																				i32.const 27780
																				i32.load offset=1864
																				i32.const 27780
																				i32.load offset=1868
																				local.set $$t18.467
																				local.set $$t18.466
																				local.set $$t18.465
																				local.set $$t18.464
																				local.set $$t18.463
																				local.set $$t18.462
																				local.set $$t18.461
																				local.set $$t18.460
																				local.set $$t18.459
																				local.set $$t18.458
																				local.set $$t18.457
																				local.set $$t18.456
																				local.set $$t18.455
																				local.set $$t18.454
																				local.set $$t18.453
																				local.set $$t18.452
																				local.set $$t18.451
																				local.set $$t18.450
																				local.set $$t18.449
																				local.set $$t18.448
																				local.set $$t18.447
																				local.set $$t18.446
																				local.set $$t18.445
																				local.set $$t18.444
																				local.set $$t18.443
																				local.set $$t18.442
																				local.set $$t18.441
																				local.set $$t18.440
																				local.set $$t18.439
																				local.set $$t18.438
																				local.set $$t18.437
																				local.set $$t18.436
																				local.set $$t18.435
																				local.set $$t18.434
																				local.set $$t18.433
																				local.set $$t18.432
																				local.set $$t18.431
																				local.set $$t18.430
																				local.set $$t18.429
																				local.set $$t18.428
																				local.set $$t18.427
																				local.set $$t18.426
																				local.set $$t18.425
																				local.set $$t18.424
																				local.set $$t18.423
																				local.set $$t18.422
																				local.set $$t18.421
																				local.set $$t18.420
																				local.set $$t18.419
																				local.set $$t18.418
																				local.set $$t18.417
																				local.set $$t18.416
																				local.set $$t18.415
																				local.set $$t18.414
																				local.set $$t18.413
																				local.set $$t18.412
																				local.set $$t18.411
																				local.set $$t18.410
																				local.set $$t18.409
																				local.set $$t18.408
																				local.set $$t18.407
																				local.set $$t18.406
																				local.set $$t18.405
																				local.set $$t18.404
																				local.set $$t18.403
																				local.set $$t18.402
																				local.set $$t18.401
																				local.set $$t18.400
																				local.set $$t18.399
																				local.set $$t18.398
																				local.set $$t18.397
																				local.set $$t18.396
																				local.set $$t18.395
																				local.set $$t18.394
																				local.set $$t18.393
																				local.set $$t18.392
																				local.set $$t18.391
																				local.set $$t18.390
																				local.set $$t18.389
																				local.set $$t18.388
																				local.set $$t18.387
																				local.set $$t18.386
																				local.set $$t18.385
																				local.set $$t18.384
																				local.set $$t18.383
																				local.set $$t18.382
																				local.set $$t18.381
																				local.set $$t18.380
																				local.set $$t18.379
																				local.set $$t18.378
																				local.set $$t18.377
																				local.set $$t18.376
																				local.set $$t18.375
																				local.set $$t18.374
																				local.set $$t18.373
																				local.set $$t18.372
																				local.set $$t18.371
																				local.set $$t18.370
																				local.set $$t18.369
																				local.set $$t18.368
																				local.set $$t18.367
																				local.set $$t18.366
																				local.set $$t18.365
																				local.set $$t18.364
																				local.set $$t18.363
																				local.set $$t18.362
																				local.set $$t18.361
																				local.set $$t18.360
																				local.set $$t18.359
																				local.set $$t18.358
																				local.set $$t18.357
																				local.set $$t18.356
																				local.set $$t18.355
																				local.set $$t18.354
																				local.set $$t18.353
																				local.set $$t18.352
																				local.set $$t18.351
																				local.set $$t18.350
																				local.set $$t18.349
																				local.set $$t18.348
																				local.set $$t18.347
																				local.set $$t18.346
																				local.set $$t18.345
																				local.set $$t18.344
																				local.set $$t18.343
																				local.set $$t18.342
																				local.set $$t18.341
																				local.set $$t18.340
																				local.set $$t18.339
																				local.set $$t18.338
																				local.set $$t18.337
																				local.set $$t18.336
																				local.set $$t18.335
																				local.set $$t18.334
																				local.set $$t18.333
																				local.set $$t18.332
																				local.set $$t18.331
																				local.set $$t18.330
																				local.set $$t18.329
																				local.set $$t18.328
																				local.set $$t18.327
																				local.set $$t18.326
																				local.set $$t18.325
																				local.set $$t18.324
																				local.set $$t18.323
																				local.set $$t18.322
																				local.set $$t18.321
																				local.set $$t18.320
																				local.set $$t18.319
																				local.set $$t18.318
																				local.set $$t18.317
																				local.set $$t18.316
																				local.set $$t18.315
																				local.set $$t18.314
																				local.set $$t18.313
																				local.set $$t18.312
																				local.set $$t18.311
																				local.set $$t18.310
																				local.set $$t18.309
																				local.set $$t18.308
																				local.set $$t18.307
																				local.set $$t18.306
																				local.set $$t18.305
																				local.set $$t18.304
																				local.set $$t18.303
																				local.set $$t18.302
																				local.set $$t18.301
																				local.set $$t18.300
																				local.set $$t18.299
																				local.set $$t18.298
																				local.set $$t18.297
																				local.set $$t18.296
																				local.set $$t18.295
																				local.set $$t18.294
																				local.set $$t18.293
																				local.set $$t18.292
																				local.set $$t18.291
																				local.set $$t18.290
																				local.set $$t18.289
																				local.set $$t18.288
																				local.set $$t18.287
																				local.set $$t18.286
																				local.set $$t18.285
																				local.set $$t18.284
																				local.set $$t18.283
																				local.set $$t18.282
																				local.set $$t18.281
																				local.set $$t18.280
																				local.set $$t18.279
																				local.set $$t18.278
																				local.set $$t18.277
																				local.set $$t18.276
																				local.set $$t18.275
																				local.set $$t18.274
																				local.set $$t18.273
																				local.set $$t18.272
																				local.set $$t18.271
																				local.set $$t18.270
																				local.set $$t18.269
																				local.set $$t18.268
																				local.set $$t18.267
																				local.set $$t18.266
																				local.set $$t18.265
																				local.set $$t18.264
																				local.set $$t18.263
																				local.set $$t18.262
																				local.set $$t18.261
																				local.set $$t18.260
																				local.set $$t18.259
																				local.set $$t18.258
																				local.set $$t18.257
																				local.set $$t18.256
																				local.set $$t18.255
																				local.set $$t18.254
																				local.set $$t18.253
																				local.set $$t18.252
																				local.set $$t18.251
																				local.set $$t18.250
																				local.set $$t18.249
																				local.set $$t18.248
																				local.set $$t18.247
																				local.set $$t18.246
																				local.set $$t18.245
																				local.set $$t18.244
																				local.set $$t18.243
																				local.set $$t18.242
																				local.set $$t18.241
																				local.set $$t18.240
																				local.set $$t18.239
																				local.set $$t18.238
																				local.set $$t18.237
																				local.set $$t18.236
																				local.set $$t18.235
																				local.set $$t18.234
																				local.set $$t18.233
																				local.set $$t18.232
																				local.set $$t18.231
																				local.set $$t18.230
																				local.set $$t18.229
																				local.set $$t18.228
																				local.set $$t18.227
																				local.set $$t18.226
																				local.set $$t18.225
																				local.set $$t18.224
																				local.set $$t18.223
																				local.set $$t18.222
																				local.set $$t18.221
																				local.set $$t18.220
																				local.set $$t18.219
																				local.set $$t18.218
																				local.set $$t18.217
																				local.set $$t18.216
																				local.set $$t18.215
																				local.set $$t18.214
																				local.set $$t18.213
																				local.set $$t18.212
																				local.set $$t18.211
																				local.set $$t18.210
																				local.set $$t18.209
																				local.set $$t18.208
																				local.set $$t18.207
																				local.set $$t18.206
																				local.set $$t18.205
																				local.set $$t18.204
																				local.set $$t18.203
																				local.set $$t18.202
																				local.set $$t18.201
																				local.set $$t18.200
																				local.set $$t18.199
																				local.set $$t18.198
																				local.set $$t18.197
																				local.set $$t18.196
																				local.set $$t18.195
																				local.set $$t18.194
																				local.set $$t18.193
																				local.set $$t18.192
																				local.set $$t18.191
																				local.set $$t18.190
																				local.set $$t18.189
																				local.set $$t18.188
																				local.set $$t18.187
																				local.set $$t18.186
																				local.set $$t18.185
																				local.set $$t18.184
																				local.set $$t18.183
																				local.set $$t18.182
																				local.set $$t18.181
																				local.set $$t18.180
																				local.set $$t18.179
																				local.set $$t18.178
																				local.set $$t18.177
																				local.set $$t18.176
																				local.set $$t18.175
																				local.set $$t18.174
																				local.set $$t18.173
																				local.set $$t18.172
																				local.set $$t18.171
																				local.set $$t18.170
																				local.set $$t18.169
																				local.set $$t18.168
																				local.set $$t18.167
																				local.set $$t18.166
																				local.set $$t18.165
																				local.set $$t18.164
																				local.set $$t18.163
																				local.set $$t18.162
																				local.set $$t18.161
																				local.set $$t18.160
																				local.set $$t18.159
																				local.set $$t18.158
																				local.set $$t18.157
																				local.set $$t18.156
																				local.set $$t18.155
																				local.set $$t18.154
																				local.set $$t18.153
																				local.set $$t18.152
																				local.set $$t18.151
																				local.set $$t18.150
																				local.set $$t18.149
																				local.set $$t18.148
																				local.set $$t18.147
																				local.set $$t18.146
																				local.set $$t18.145
																				local.set $$t18.144
																				local.set $$t18.143
																				local.set $$t18.142
																				local.set $$t18.141
																				local.set $$t18.140
																				local.set $$t18.139
																				local.set $$t18.138
																				local.set $$t18.137
																				local.set $$t18.136
																				local.set $$t18.135
																				local.set $$t18.134
																				local.set $$t18.133
																				local.set $$t18.132
																				local.set $$t18.131
																				local.set $$t18.130
																				local.set $$t18.129
																				local.set $$t18.128
																				local.set $$t18.127
																				local.set $$t18.126
																				local.set $$t18.125
																				local.set $$t18.124
																				local.set $$t18.123
																				local.set $$t18.122
																				local.set $$t18.121
																				local.set $$t18.120
																				local.set $$t18.119
																				local.set $$t18.118
																				local.set $$t18.117
																				local.set $$t18.116
																				local.set $$t18.115
																				local.set $$t18.114
																				local.set $$t18.113
																				local.set $$t18.112
																				local.set $$t18.111
																				local.set $$t18.110
																				local.set $$t18.109
																				local.set $$t18.108
																				local.set $$t18.107
																				local.set $$t18.106
																				local.set $$t18.105
																				local.set $$t18.104
																				local.set $$t18.103
																				local.set $$t18.102
																				local.set $$t18.101
																				local.set $$t18.100
																				local.set $$t18.99
																				local.set $$t18.98
																				local.set $$t18.97
																				local.set $$t18.96
																				local.set $$t18.95
																				local.set $$t18.94
																				local.set $$t18.93
																				local.set $$t18.92
																				local.set $$t18.91
																				local.set $$t18.90
																				local.set $$t18.89
																				local.set $$t18.88
																				local.set $$t18.87
																				local.set $$t18.86
																				local.set $$t18.85
																				local.set $$t18.84
																				local.set $$t18.83
																				local.set $$t18.82
																				local.set $$t18.81
																				local.set $$t18.80
																				local.set $$t18.79
																				local.set $$t18.78
																				local.set $$t18.77
																				local.set $$t18.76
																				local.set $$t18.75
																				local.set $$t18.74
																				local.set $$t18.73
																				local.set $$t18.72
																				local.set $$t18.71
																				local.set $$t18.70
																				local.set $$t18.69
																				local.set $$t18.68
																				local.set $$t18.67
																				local.set $$t18.66
																				local.set $$t18.65
																				local.set $$t18.64
																				local.set $$t18.63
																				local.set $$t18.62
																				local.set $$t18.61
																				local.set $$t18.60
																				local.set $$t18.59
																				local.set $$t18.58
																				local.set $$t18.57
																				local.set $$t18.56
																				local.set $$t18.55
																				local.set $$t18.54
																				local.set $$t18.53
																				local.set $$t18.52
																				local.set $$t18.51
																				local.set $$t18.50
																				local.set $$t18.49
																				local.set $$t18.48
																				local.set $$t18.47
																				local.set $$t18.46
																				local.set $$t18.45
																				local.set $$t18.44
																				local.set $$t18.43
																				local.set $$t18.42
																				local.set $$t18.41
																				local.set $$t18.40
																				local.set $$t18.39
																				local.set $$t18.38
																				local.set $$t18.37
																				local.set $$t18.36
																				local.set $$t18.35
																				local.set $$t18.34
																				local.set $$t18.33
																				local.set $$t18.32
																				local.set $$t18.31
																				local.set $$t18.30
																				local.set $$t18.29
																				local.set $$t18.28
																				local.set $$t18.27
																				local.set $$t18.26
																				local.set $$t18.25
																				local.set $$t18.24
																				local.set $$t18.23
																				local.set $$t18.22
																				local.set $$t18.21
																				local.set $$t18.20
																				local.set $$t18.19
																				local.set $$t18.18
																				local.set $$t18.17
																				local.set $$t18.16
																				local.set $$t18.15
																				local.set $$t18.14
																				local.set $$t18.13
																				local.set $$t18.12
																				local.set $$t18.11
																				local.set $$t18.10
																				local.set $$t18.9
																				local.set $$t18.8
																				local.set $$t18.7
																				local.set $$t18.6
																				local.set $$t18.5
																				local.set $$t18.4
																				local.set $$t18.3
																				local.set $$t18.2
																				local.set $$t18.1
																				local.set $$t18.0
																				i32.const 26720
																				i32.load16_u
																				i32.const 26720
																				i32.load16_u offset=2
																				i32.const 26720
																				i32.load16_u offset=4
																				i32.const 26720
																				i32.load16_u offset=6
																				i32.const 26720
																				i32.load16_u offset=8
																				i32.const 26720
																				i32.load16_u offset=10
																				i32.const 26720
																				i32.load16_u offset=12
																				i32.const 26720
																				i32.load16_u offset=14
																				i32.const 26720
																				i32.load16_u offset=16
																				i32.const 26720
																				i32.load16_u offset=18
																				i32.const 26720
																				i32.load16_u offset=20
																				i32.const 26720
																				i32.load16_u offset=22
																				i32.const 26720
																				i32.load16_u offset=24
																				i32.const 26720
																				i32.load16_u offset=26
																				i32.const 26720
																				i32.load16_u offset=28
																				i32.const 26720
																				i32.load16_u offset=30
																				i32.const 26720
																				i32.load16_u offset=32
																				i32.const 26720
																				i32.load16_u offset=34
																				i32.const 26720
																				i32.load16_u offset=36
																				i32.const 26720
																				i32.load16_u offset=38
																				i32.const 26720
																				i32.load16_u offset=40
																				i32.const 26720
																				i32.load16_u offset=42
																				i32.const 26720
																				i32.load16_u offset=44
																				i32.const 26720
																				i32.load16_u offset=46
																				i32.const 26720
																				i32.load16_u offset=48
																				i32.const 26720
																				i32.load16_u offset=50
																				i32.const 26720
																				i32.load16_u offset=52
																				i32.const 26720
																				i32.load16_u offset=54
																				i32.const 26720
																				i32.load16_u offset=56
																				i32.const 26720
																				i32.load16_u offset=58
																				i32.const 26720
																				i32.load16_u offset=60
																				i32.const 26720
																				i32.load16_u offset=62
																				i32.const 26720
																				i32.load16_u offset=64
																				i32.const 26720
																				i32.load16_u offset=66
																				i32.const 26720
																				i32.load16_u offset=68
																				i32.const 26720
																				i32.load16_u offset=70
																				i32.const 26720
																				i32.load16_u offset=72
																				i32.const 26720
																				i32.load16_u offset=74
																				i32.const 26720
																				i32.load16_u offset=76
																				i32.const 26720
																				i32.load16_u offset=78
																				i32.const 26720
																				i32.load16_u offset=80
																				i32.const 26720
																				i32.load16_u offset=82
																				i32.const 26720
																				i32.load16_u offset=84
																				i32.const 26720
																				i32.load16_u offset=86
																				i32.const 26720
																				i32.load16_u offset=88
																				i32.const 26720
																				i32.load16_u offset=90
																				i32.const 26720
																				i32.load16_u offset=92
																				i32.const 26720
																				i32.load16_u offset=94
																				i32.const 26720
																				i32.load16_u offset=96
																				i32.const 26720
																				i32.load16_u offset=98
																				i32.const 26720
																				i32.load16_u offset=100
																				i32.const 26720
																				i32.load16_u offset=102
																				i32.const 26720
																				i32.load16_u offset=104
																				i32.const 26720
																				i32.load16_u offset=106
																				i32.const 26720
																				i32.load16_u offset=108
																				i32.const 26720
																				i32.load16_u offset=110
																				i32.const 26720
																				i32.load16_u offset=112
																				i32.const 26720
																				i32.load16_u offset=114
																				i32.const 26720
																				i32.load16_u offset=116
																				i32.const 26720
																				i32.load16_u offset=118
																				i32.const 26720
																				i32.load16_u offset=120
																				i32.const 26720
																				i32.load16_u offset=122
																				i32.const 26720
																				i32.load16_u offset=124
																				i32.const 26720
																				i32.load16_u offset=126
																				i32.const 26720
																				i32.load16_u offset=128
																				i32.const 26720
																				i32.load16_u offset=130
																				i32.const 26720
																				i32.load16_u offset=132
																				i32.const 26720
																				i32.load16_u offset=134
																				i32.const 26720
																				i32.load16_u offset=136
																				i32.const 26720
																				i32.load16_u offset=138
																				i32.const 26720
																				i32.load16_u offset=140
																				i32.const 26720
																				i32.load16_u offset=142
																				i32.const 26720
																				i32.load16_u offset=144
																				i32.const 26720
																				i32.load16_u offset=146
																				i32.const 26720
																				i32.load16_u offset=148
																				i32.const 26720
																				i32.load16_u offset=150
																				i32.const 26720
																				i32.load16_u offset=152
																				i32.const 26720
																				i32.load16_u offset=154
																				i32.const 26720
																				i32.load16_u offset=156
																				i32.const 26720
																				i32.load16_u offset=158
																				i32.const 26720
																				i32.load16_u offset=160
																				i32.const 26720
																				i32.load16_u offset=162
																				i32.const 26720
																				i32.load16_u offset=164
																				i32.const 26720
																				i32.load16_u offset=166
																				i32.const 26720
																				i32.load16_u offset=168
																				i32.const 26720
																				i32.load16_u offset=170
																				i32.const 26720
																				i32.load16_u offset=172
																				i32.const 26720
																				i32.load16_u offset=174
																				i32.const 26720
																				i32.load16_u offset=176
																				i32.const 26720
																				i32.load16_u offset=178
																				i32.const 26720
																				i32.load16_u offset=180
																				i32.const 26720
																				i32.load16_u offset=182
																				i32.const 26720
																				i32.load16_u offset=184
																				i32.const 26720
																				i32.load16_u offset=186
																				i32.const 26720
																				i32.load16_u offset=188
																				local.set $$t19.94
																				local.set $$t19.93
																				local.set $$t19.92
																				local.set $$t19.91
																				local.set $$t19.90
																				local.set $$t19.89
																				local.set $$t19.88
																				local.set $$t19.87
																				local.set $$t19.86
																				local.set $$t19.85
																				local.set $$t19.84
																				local.set $$t19.83
																				local.set $$t19.82
																				local.set $$t19.81
																				local.set $$t19.80
																				local.set $$t19.79
																				local.set $$t19.78
																				local.set $$t19.77
																				local.set $$t19.76
																				local.set $$t19.75
																				local.set $$t19.74
																				local.set $$t19.73
																				local.set $$t19.72
																				local.set $$t19.71
																				local.set $$t19.70
																				local.set $$t19.69
																				local.set $$t19.68
																				local.set $$t19.67
																				local.set $$t19.66
																				local.set $$t19.65
																				local.set $$t19.64
																				local.set $$t19.63
																				local.set $$t19.62
																				local.set $$t19.61
																				local.set $$t19.60
																				local.set $$t19.59
																				local.set $$t19.58
																				local.set $$t19.57
																				local.set $$t19.56
																				local.set $$t19.55
																				local.set $$t19.54
																				local.set $$t19.53
																				local.set $$t19.52
																				local.set $$t19.51
																				local.set $$t19.50
																				local.set $$t19.49
																				local.set $$t19.48
																				local.set $$t19.47
																				local.set $$t19.46
																				local.set $$t19.45
																				local.set $$t19.44
																				local.set $$t19.43
																				local.set $$t19.42
																				local.set $$t19.41
																				local.set $$t19.40
																				local.set $$t19.39
																				local.set $$t19.38
																				local.set $$t19.37
																				local.set $$t19.36
																				local.set $$t19.35
																				local.set $$t19.34
																				local.set $$t19.33
																				local.set $$t19.32
																				local.set $$t19.31
																				local.set $$t19.30
																				local.set $$t19.29
																				local.set $$t19.28
																				local.set $$t19.27
																				local.set $$t19.26
																				local.set $$t19.25
																				local.set $$t19.24
																				local.set $$t19.23
																				local.set $$t19.22
																				local.set $$t19.21
																				local.set $$t19.20
																				local.set $$t19.19
																				local.set $$t19.18
																				local.set $$t19.17
																				local.set $$t19.16
																				local.set $$t19.15
																				local.set $$t19.14
																				local.set $$t19.13
																				local.set $$t19.12
																				local.set $$t19.11
																				local.set $$t19.10
																				local.set $$t19.9
																				local.set $$t19.8
																				local.set $$t19.7
																				local.set $$t19.6
																				local.set $$t19.5
																				local.set $$t19.4
																				local.set $$t19.3
																				local.set $$t19.2
																				local.set $$t19.1
																				local.set $$t19.0
																				local.get $$t15.1
																				local.get $$t18.0
																				i32.store
																				local.get $$t15.1
																				local.get $$t18.1
																				i32.store offset=4
																				local.get $$t15.1
																				local.get $$t18.2
																				i32.store offset=8
																				local.get $$t15.1
																				local.get $$t18.3
																				i32.store offset=12
																				local.get $$t15.1
																				local.get $$t18.4
																				i32.store offset=16
																				local.get $$t15.1
																				local.get $$t18.5
																				i32.store offset=20
																				local.get $$t15.1
																				local.get $$t18.6
																				i32.store offset=24
																				local.get $$t15.1
																				local.get $$t18.7
																				i32.store offset=28
																				local.get $$t15.1
																				local.get $$t18.8
																				i32.store offset=32
																				local.get $$t15.1
																				local.get $$t18.9
																				i32.store offset=36
																				local.get $$t15.1
																				local.get $$t18.10
																				i32.store offset=40
																				local.get $$t15.1
																				local.get $$t18.11
																				i32.store offset=44
																				local.get $$t15.1
																				local.get $$t18.12
																				i32.store offset=48
																				local.get $$t15.1
																				local.get $$t18.13
																				i32.store offset=52
																				local.get $$t15.1
																				local.get $$t18.14
																				i32.store offset=56
																				local.get $$t15.1
																				local.get $$t18.15
																				i32.store offset=60
																				local.get $$t15.1
																				local.get $$t18.16
																				i32.store offset=64
																				local.get $$t15.1
																				local.get $$t18.17
																				i32.store offset=68
																				local.get $$t15.1
																				local.get $$t18.18
																				i32.store offset=72
																				local.get $$t15.1
																				local.get $$t18.19
																				i32.store offset=76
																				local.get $$t15.1
																				local.get $$t18.20
																				i32.store offset=80
																				local.get $$t15.1
																				local.get $$t18.21
																				i32.store offset=84
																				local.get $$t15.1
																				local.get $$t18.22
																				i32.store offset=88
																				local.get $$t15.1
																				local.get $$t18.23
																				i32.store offset=92
																				local.get $$t15.1
																				local.get $$t18.24
																				i32.store offset=96
																				local.get $$t15.1
																				local.get $$t18.25
																				i32.store offset=100
																				local.get $$t15.1
																				local.get $$t18.26
																				i32.store offset=104
																				local.get $$t15.1
																				local.get $$t18.27
																				i32.store offset=108
																				local.get $$t15.1
																				local.get $$t18.28
																				i32.store offset=112
																				local.get $$t15.1
																				local.get $$t18.29
																				i32.store offset=116
																				local.get $$t15.1
																				local.get $$t18.30
																				i32.store offset=120
																				local.get $$t15.1
																				local.get $$t18.31
																				i32.store offset=124
																				local.get $$t15.1
																				local.get $$t18.32
																				i32.store offset=128
																				local.get $$t15.1
																				local.get $$t18.33
																				i32.store offset=132
																				local.get $$t15.1
																				local.get $$t18.34
																				i32.store offset=136
																				local.get $$t15.1
																				local.get $$t18.35
																				i32.store offset=140
																				local.get $$t15.1
																				local.get $$t18.36
																				i32.store offset=144
																				local.get $$t15.1
																				local.get $$t18.37
																				i32.store offset=148
																				local.get $$t15.1
																				local.get $$t18.38
																				i32.store offset=152
																				local.get $$t15.1
																				local.get $$t18.39
																				i32.store offset=156
																				local.get $$t15.1
																				local.get $$t18.40
																				i32.store offset=160
																				local.get $$t15.1
																				local.get $$t18.41
																				i32.store offset=164
																				local.get $$t15.1
																				local.get $$t18.42
																				i32.store offset=168
																				local.get $$t15.1
																				local.get $$t18.43
																				i32.store offset=172
																				local.get $$t15.1
																				local.get $$t18.44
																				i32.store offset=176
																				local.get $$t15.1
																				local.get $$t18.45
																				i32.store offset=180
																				local.get $$t15.1
																				local.get $$t18.46
																				i32.store offset=184
																				local.get $$t15.1
																				local.get $$t18.47
																				i32.store offset=188
																				local.get $$t15.1
																				local.get $$t18.48
																				i32.store offset=192
																				local.get $$t15.1
																				local.get $$t18.49
																				i32.store offset=196
																				local.get $$t15.1
																				local.get $$t18.50
																				i32.store offset=200
																				local.get $$t15.1
																				local.get $$t18.51
																				i32.store offset=204
																				local.get $$t15.1
																				local.get $$t18.52
																				i32.store offset=208
																				local.get $$t15.1
																				local.get $$t18.53
																				i32.store offset=212
																				local.get $$t15.1
																				local.get $$t18.54
																				i32.store offset=216
																				local.get $$t15.1
																				local.get $$t18.55
																				i32.store offset=220
																				local.get $$t15.1
																				local.get $$t18.56
																				i32.store offset=224
																				local.get $$t15.1
																				local.get $$t18.57
																				i32.store offset=228
																				local.get $$t15.1
																				local.get $$t18.58
																				i32.store offset=232
																				local.get $$t15.1
																				local.get $$t18.59
																				i32.store offset=236
																				local.get $$t15.1
																				local.get $$t18.60
																				i32.store offset=240
																				local.get $$t15.1
																				local.get $$t18.61
																				i32.store offset=244
																				local.get $$t15.1
																				local.get $$t18.62
																				i32.store offset=248
																				local.get $$t15.1
																				local.get $$t18.63
																				i32.store offset=252
																				local.get $$t15.1
																				local.get $$t18.64
																				i32.store offset=256
																				local.get $$t15.1
																				local.get $$t18.65
																				i32.store offset=260
																				local.get $$t15.1
																				local.get $$t18.66
																				i32.store offset=264
																				local.get $$t15.1
																				local.get $$t18.67
																				i32.store offset=268
																				local.get $$t15.1
																				local.get $$t18.68
																				i32.store offset=272
																				local.get $$t15.1
																				local.get $$t18.69
																				i32.store offset=276
																				local.get $$t15.1
																				local.get $$t18.70
																				i32.store offset=280
																				local.get $$t15.1
																				local.get $$t18.71
																				i32.store offset=284
																				local.get $$t15.1
																				local.get $$t18.72
																				i32.store offset=288
																				local.get $$t15.1
																				local.get $$t18.73
																				i32.store offset=292
																				local.get $$t15.1
																				local.get $$t18.74
																				i32.store offset=296
																				local.get $$t15.1
																				local.get $$t18.75
																				i32.store offset=300
																				local.get $$t15.1
																				local.get $$t18.76
																				i32.store offset=304
																				local.get $$t15.1
																				local.get $$t18.77
																				i32.store offset=308
																				local.get $$t15.1
																				local.get $$t18.78
																				i32.store offset=312
																				local.get $$t15.1
																				local.get $$t18.79
																				i32.store offset=316
																				local.get $$t15.1
																				local.get $$t18.80
																				i32.store offset=320
																				local.get $$t15.1
																				local.get $$t18.81
																				i32.store offset=324
																				local.get $$t15.1
																				local.get $$t18.82
																				i32.store offset=328
																				local.get $$t15.1
																				local.get $$t18.83
																				i32.store offset=332
																				local.get $$t15.1
																				local.get $$t18.84
																				i32.store offset=336
																				local.get $$t15.1
																				local.get $$t18.85
																				i32.store offset=340
																				local.get $$t15.1
																				local.get $$t18.86
																				i32.store offset=344
																				local.get $$t15.1
																				local.get $$t18.87
																				i32.store offset=348
																				local.get $$t15.1
																				local.get $$t18.88
																				i32.store offset=352
																				local.get $$t15.1
																				local.get $$t18.89
																				i32.store offset=356
																				local.get $$t15.1
																				local.get $$t18.90
																				i32.store offset=360
																				local.get $$t15.1
																				local.get $$t18.91
																				i32.store offset=364
																				local.get $$t15.1
																				local.get $$t18.92
																				i32.store offset=368
																				local.get $$t15.1
																				local.get $$t18.93
																				i32.store offset=372
																				local.get $$t15.1
																				local.get $$t18.94
																				i32.store offset=376
																				local.get $$t15.1
																				local.get $$t18.95
																				i32.store offset=380
																				local.get $$t15.1
																				local.get $$t18.96
																				i32.store offset=384
																				local.get $$t15.1
																				local.get $$t18.97
																				i32.store offset=388
																				local.get $$t15.1
																				local.get $$t18.98
																				i32.store offset=392
																				local.get $$t15.1
																				local.get $$t18.99
																				i32.store offset=396
																				local.get $$t15.1
																				local.get $$t18.100
																				i32.store offset=400
																				local.get $$t15.1
																				local.get $$t18.101
																				i32.store offset=404
																				local.get $$t15.1
																				local.get $$t18.102
																				i32.store offset=408
																				local.get $$t15.1
																				local.get $$t18.103
																				i32.store offset=412
																				local.get $$t15.1
																				local.get $$t18.104
																				i32.store offset=416
																				local.get $$t15.1
																				local.get $$t18.105
																				i32.store offset=420
																				local.get $$t15.1
																				local.get $$t18.106
																				i32.store offset=424
																				local.get $$t15.1
																				local.get $$t18.107
																				i32.store offset=428
																				local.get $$t15.1
																				local.get $$t18.108
																				i32.store offset=432
																				local.get $$t15.1
																				local.get $$t18.109
																				i32.store offset=436
																				local.get $$t15.1
																				local.get $$t18.110
																				i32.store offset=440
																				local.get $$t15.1
																				local.get $$t18.111
																				i32.store offset=444
																				local.get $$t15.1
																				local.get $$t18.112
																				i32.store offset=448
																				local.get $$t15.1
																				local.get $$t18.113
																				i32.store offset=452
																				local.get $$t15.1
																				local.get $$t18.114
																				i32.store offset=456
																				local.get $$t15.1
																				local.get $$t18.115
																				i32.store offset=460
																				local.get $$t15.1
																				local.get $$t18.116
																				i32.store offset=464
																				local.get $$t15.1
																				local.get $$t18.117
																				i32.store offset=468
																				local.get $$t15.1
																				local.get $$t18.118
																				i32.store offset=472
																				local.get $$t15.1
																				local.get $$t18.119
																				i32.store offset=476
																				local.get $$t15.1
																				local.get $$t18.120
																				i32.store offset=480
																				local.get $$t15.1
																				local.get $$t18.121
																				i32.store offset=484
																				local.get $$t15.1
																				local.get $$t18.122
																				i32.store offset=488
																				local.get $$t15.1
																				local.get $$t18.123
																				i32.store offset=492
																				local.get $$t15.1
																				local.get $$t18.124
																				i32.store offset=496
																				local.get $$t15.1
																				local.get $$t18.125
																				i32.store offset=500
																				local.get $$t15.1
																				local.get $$t18.126
																				i32.store offset=504
																				local.get $$t15.1
																				local.get $$t18.127
																				i32.store offset=508
																				local.get $$t15.1
																				local.get $$t18.128
																				i32.store offset=512
																				local.get $$t15.1
																				local.get $$t18.129
																				i32.store offset=516
																				local.get $$t15.1
																				local.get $$t18.130
																				i32.store offset=520
																				local.get $$t15.1
																				local.get $$t18.131
																				i32.store offset=524
																				local.get $$t15.1
																				local.get $$t18.132
																				i32.store offset=528
																				local.get $$t15.1
																				local.get $$t18.133
																				i32.store offset=532
																				local.get $$t15.1
																				local.get $$t18.134
																				i32.store offset=536
																				local.get $$t15.1
																				local.get $$t18.135
																				i32.store offset=540
																				local.get $$t15.1
																				local.get $$t18.136
																				i32.store offset=544
																				local.get $$t15.1
																				local.get $$t18.137
																				i32.store offset=548
																				local.get $$t15.1
																				local.get $$t18.138
																				i32.store offset=552
																				local.get $$t15.1
																				local.get $$t18.139
																				i32.store offset=556
																				local.get $$t15.1
																				local.get $$t18.140
																				i32.store offset=560
																				local.get $$t15.1
																				local.get $$t18.141
																				i32.store offset=564
																				local.get $$t15.1
																				local.get $$t18.142
																				i32.store offset=568
																				local.get $$t15.1
																				local.get $$t18.143
																				i32.store offset=572
																				local.get $$t15.1
																				local.get $$t18.144
																				i32.store offset=576
																				local.get $$t15.1
																				local.get $$t18.145
																				i32.store offset=580
																				local.get $$t15.1
																				local.get $$t18.146
																				i32.store offset=584
																				local.get $$t15.1
																				local.get $$t18.147
																				i32.store offset=588
																				local.get $$t15.1
																				local.get $$t18.148
																				i32.store offset=592
																				local.get $$t15.1
																				local.get $$t18.149
																				i32.store offset=596
																				local.get $$t15.1
																				local.get $$t18.150
																				i32.store offset=600
																				local.get $$t15.1
																				local.get $$t18.151
																				i32.store offset=604
																				local.get $$t15.1
																				local.get $$t18.152
																				i32.store offset=608
																				local.get $$t15.1
																				local.get $$t18.153
																				i32.store offset=612
																				local.get $$t15.1
																				local.get $$t18.154
																				i32.store offset=616
																				local.get $$t15.1
																				local.get $$t18.155
																				i32.store offset=620
																				local.get $$t15.1
																				local.get $$t18.156
																				i32.store offset=624
																				local.get $$t15.1
																				local.get $$t18.157
																				i32.store offset=628
																				local.get $$t15.1
																				local.get $$t18.158
																				i32.store offset=632
																				local.get $$t15.1
																				local.get $$t18.159
																				i32.store offset=636
																				local.get $$t15.1
																				local.get $$t18.160
																				i32.store offset=640
																				local.get $$t15.1
																				local.get $$t18.161
																				i32.store offset=644
																				local.get $$t15.1
																				local.get $$t18.162
																				i32.store offset=648
																				local.get $$t15.1
																				local.get $$t18.163
																				i32.store offset=652
																				local.get $$t15.1
																				local.get $$t18.164
																				i32.store offset=656
																				local.get $$t15.1
																				local.get $$t18.165
																				i32.store offset=660
																				local.get $$t15.1
																				local.get $$t18.166
																				i32.store offset=664
																				local.get $$t15.1
																				local.get $$t18.167
																				i32.store offset=668
																				local.get $$t15.1
																				local.get $$t18.168
																				i32.store offset=672
																				local.get $$t15.1
																				local.get $$t18.169
																				i32.store offset=676
																				local.get $$t15.1
																				local.get $$t18.170
																				i32.store offset=680
																				local.get $$t15.1
																				local.get $$t18.171
																				i32.store offset=684
																				local.get $$t15.1
																				local.get $$t18.172
																				i32.store offset=688
																				local.get $$t15.1
																				local.get $$t18.173
																				i32.store offset=692
																				local.get $$t15.1
																				local.get $$t18.174
																				i32.store offset=696
																				local.get $$t15.1
																				local.get $$t18.175
																				i32.store offset=700
																				local.get $$t15.1
																				local.get $$t18.176
																				i32.store offset=704
																				local.get $$t15.1
																				local.get $$t18.177
																				i32.store offset=708
																				local.get $$t15.1
																				local.get $$t18.178
																				i32.store offset=712
																				local.get $$t15.1
																				local.get $$t18.179
																				i32.store offset=716
																				local.get $$t15.1
																				local.get $$t18.180
																				i32.store offset=720
																				local.get $$t15.1
																				local.get $$t18.181
																				i32.store offset=724
																				local.get $$t15.1
																				local.get $$t18.182
																				i32.store offset=728
																				local.get $$t15.1
																				local.get $$t18.183
																				i32.store offset=732
																				local.get $$t15.1
																				local.get $$t18.184
																				i32.store offset=736
																				local.get $$t15.1
																				local.get $$t18.185
																				i32.store offset=740
																				local.get $$t15.1
																				local.get $$t18.186
																				i32.store offset=744
																				local.get $$t15.1
																				local.get $$t18.187
																				i32.store offset=748
																				local.get $$t15.1
																				local.get $$t18.188
																				i32.store offset=752
																				local.get $$t15.1
																				local.get $$t18.189
																				i32.store offset=756
																				local.get $$t15.1
																				local.get $$t18.190
																				i32.store offset=760
																				local.get $$t15.1
																				local.get $$t18.191
																				i32.store offset=764
																				local.get $$t15.1
																				local.get $$t18.192
																				i32.store offset=768
																				local.get $$t15.1
																				local.get $$t18.193
																				i32.store offset=772
																				local.get $$t15.1
																				local.get $$t18.194
																				i32.store offset=776
																				local.get $$t15.1
																				local.get $$t18.195
																				i32.store offset=780
																				local.get $$t15.1
																				local.get $$t18.196
																				i32.store offset=784
																				local.get $$t15.1
																				local.get $$t18.197
																				i32.store offset=788
																				local.get $$t15.1
																				local.get $$t18.198
																				i32.store offset=792
																				local.get $$t15.1
																				local.get $$t18.199
																				i32.store offset=796
																				local.get $$t15.1
																				local.get $$t18.200
																				i32.store offset=800
																				local.get $$t15.1
																				local.get $$t18.201
																				i32.store offset=804
																				local.get $$t15.1
																				local.get $$t18.202
																				i32.store offset=808
																				local.get $$t15.1
																				local.get $$t18.203
																				i32.store offset=812
																				local.get $$t15.1
																				local.get $$t18.204
																				i32.store offset=816
																				local.get $$t15.1
																				local.get $$t18.205
																				i32.store offset=820
																				local.get $$t15.1
																				local.get $$t18.206
																				i32.store offset=824
																				local.get $$t15.1
																				local.get $$t18.207
																				i32.store offset=828
																				local.get $$t15.1
																				local.get $$t18.208
																				i32.store offset=832
																				local.get $$t15.1
																				local.get $$t18.209
																				i32.store offset=836
																				local.get $$t15.1
																				local.get $$t18.210
																				i32.store offset=840
																				local.get $$t15.1
																				local.get $$t18.211
																				i32.store offset=844
																				local.get $$t15.1
																				local.get $$t18.212
																				i32.store offset=848
																				local.get $$t15.1
																				local.get $$t18.213
																				i32.store offset=852
																				local.get $$t15.1
																				local.get $$t18.214
																				i32.store offset=856
																				local.get $$t15.1
																				local.get $$t18.215
																				i32.store offset=860
																				local.get $$t15.1
																				local.get $$t18.216
																				i32.store offset=864
																				local.get $$t15.1
																				local.get $$t18.217
																				i32.store offset=868
																				local.get $$t15.1
																				local.get $$t18.218
																				i32.store offset=872
																				local.get $$t15.1
																				local.get $$t18.219
																				i32.store offset=876
																				local.get $$t15.1
																				local.get $$t18.220
																				i32.store offset=880
																				local.get $$t15.1
																				local.get $$t18.221
																				i32.store offset=884
																				local.get $$t15.1
																				local.get $$t18.222
																				i32.store offset=888
																				local.get $$t15.1
																				local.get $$t18.223
																				i32.store offset=892
																				local.get $$t15.1
																				local.get $$t18.224
																				i32.store offset=896
																				local.get $$t15.1
																				local.get $$t18.225
																				i32.store offset=900
																				local.get $$t15.1
																				local.get $$t18.226
																				i32.store offset=904
																				local.get $$t15.1
																				local.get $$t18.227
																				i32.store offset=908
																				local.get $$t15.1
																				local.get $$t18.228
																				i32.store offset=912
																				local.get $$t15.1
																				local.get $$t18.229
																				i32.store offset=916
																				local.get $$t15.1
																				local.get $$t18.230
																				i32.store offset=920
																				local.get $$t15.1
																				local.get $$t18.231
																				i32.store offset=924
																				local.get $$t15.1
																				local.get $$t18.232
																				i32.store offset=928
																				local.get $$t15.1
																				local.get $$t18.233
																				i32.store offset=932
																				local.get $$t15.1
																				local.get $$t18.234
																				i32.store offset=936
																				local.get $$t15.1
																				local.get $$t18.235
																				i32.store offset=940
																				local.get $$t15.1
																				local.get $$t18.236
																				i32.store offset=944
																				local.get $$t15.1
																				local.get $$t18.237
																				i32.store offset=948
																				local.get $$t15.1
																				local.get $$t18.238
																				i32.store offset=952
																				local.get $$t15.1
																				local.get $$t18.239
																				i32.store offset=956
																				local.get $$t15.1
																				local.get $$t18.240
																				i32.store offset=960
																				local.get $$t15.1
																				local.get $$t18.241
																				i32.store offset=964
																				local.get $$t15.1
																				local.get $$t18.242
																				i32.store offset=968
																				local.get $$t15.1
																				local.get $$t18.243
																				i32.store offset=972
																				local.get $$t15.1
																				local.get $$t18.244
																				i32.store offset=976
																				local.get $$t15.1
																				local.get $$t18.245
																				i32.store offset=980
																				local.get $$t15.1
																				local.get $$t18.246
																				i32.store offset=984
																				local.get $$t15.1
																				local.get $$t18.247
																				i32.store offset=988
																				local.get $$t15.1
																				local.get $$t18.248
																				i32.store offset=992
																				local.get $$t15.1
																				local.get $$t18.249
																				i32.store offset=996
																				local.get $$t15.1
																				local.get $$t18.250
																				i32.store offset=1000
																				local.get $$t15.1
																				local.get $$t18.251
																				i32.store offset=1004
																				local.get $$t15.1
																				local.get $$t18.252
																				i32.store offset=1008
																				local.get $$t15.1
																				local.get $$t18.253
																				i32.store offset=1012
																				local.get $$t15.1
																				local.get $$t18.254
																				i32.store offset=1016
																				local.get $$t15.1
																				local.get $$t18.255
																				i32.store offset=1020
																				local.get $$t15.1
																				local.get $$t18.256
																				i32.store offset=1024
																				local.get $$t15.1
																				local.get $$t18.257
																				i32.store offset=1028
																				local.get $$t15.1
																				local.get $$t18.258
																				i32.store offset=1032
																				local.get $$t15.1
																				local.get $$t18.259
																				i32.store offset=1036
																				local.get $$t15.1
																				local.get $$t18.260
																				i32.store offset=1040
																				local.get $$t15.1
																				local.get $$t18.261
																				i32.store offset=1044
																				local.get $$t15.1
																				local.get $$t18.262
																				i32.store offset=1048
																				local.get $$t15.1
																				local.get $$t18.263
																				i32.store offset=1052
																				local.get $$t15.1
																				local.get $$t18.264
																				i32.store offset=1056
																				local.get $$t15.1
																				local.get $$t18.265
																				i32.store offset=1060
																				local.get $$t15.1
																				local.get $$t18.266
																				i32.store offset=1064
																				local.get $$t15.1
																				local.get $$t18.267
																				i32.store offset=1068
																				local.get $$t15.1
																				local.get $$t18.268
																				i32.store offset=1072
																				local.get $$t15.1
																				local.get $$t18.269
																				i32.store offset=1076
																				local.get $$t15.1
																				local.get $$t18.270
																				i32.store offset=1080
																				local.get $$t15.1
																				local.get $$t18.271
																				i32.store offset=1084
																				local.get $$t15.1
																				local.get $$t18.272
																				i32.store offset=1088
																				local.get $$t15.1
																				local.get $$t18.273
																				i32.store offset=1092
																				local.get $$t15.1
																				local.get $$t18.274
																				i32.store offset=1096
																				local.get $$t15.1
																				local.get $$t18.275
																				i32.store offset=1100
																				local.get $$t15.1
																				local.get $$t18.276
																				i32.store offset=1104
																				local.get $$t15.1
																				local.get $$t18.277
																				i32.store offset=1108
																				local.get $$t15.1
																				local.get $$t18.278
																				i32.store offset=1112
																				local.get $$t15.1
																				local.get $$t18.279
																				i32.store offset=1116
																				local.get $$t15.1
																				local.get $$t18.280
																				i32.store offset=1120
																				local.get $$t15.1
																				local.get $$t18.281
																				i32.store offset=1124
																				local.get $$t15.1
																				local.get $$t18.282
																				i32.store offset=1128
																				local.get $$t15.1
																				local.get $$t18.283
																				i32.store offset=1132
																				local.get $$t15.1
																				local.get $$t18.284
																				i32.store offset=1136
																				local.get $$t15.1
																				local.get $$t18.285
																				i32.store offset=1140
																				local.get $$t15.1
																				local.get $$t18.286
																				i32.store offset=1144
																				local.get $$t15.1
																				local.get $$t18.287
																				i32.store offset=1148
																				local.get $$t15.1
																				local.get $$t18.288
																				i32.store offset=1152
																				local.get $$t15.1
																				local.get $$t18.289
																				i32.store offset=1156
																				local.get $$t15.1
																				local.get $$t18.290
																				i32.store offset=1160
																				local.get $$t15.1
																				local.get $$t18.291
																				i32.store offset=1164
																				local.get $$t15.1
																				local.get $$t18.292
																				i32.store offset=1168
																				local.get $$t15.1
																				local.get $$t18.293
																				i32.store offset=1172
																				local.get $$t15.1
																				local.get $$t18.294
																				i32.store offset=1176
																				local.get $$t15.1
																				local.get $$t18.295
																				i32.store offset=1180
																				local.get $$t15.1
																				local.get $$t18.296
																				i32.store offset=1184
																				local.get $$t15.1
																				local.get $$t18.297
																				i32.store offset=1188
																				local.get $$t15.1
																				local.get $$t18.298
																				i32.store offset=1192
																				local.get $$t15.1
																				local.get $$t18.299
																				i32.store offset=1196
																				local.get $$t15.1
																				local.get $$t18.300
																				i32.store offset=1200
																				local.get $$t15.1
																				local.get $$t18.301
																				i32.store offset=1204
																				local.get $$t15.1
																				local.get $$t18.302
																				i32.store offset=1208
																				local.get $$t15.1
																				local.get $$t18.303
																				i32.store offset=1212
																				local.get $$t15.1
																				local.get $$t18.304
																				i32.store offset=1216
																				local.get $$t15.1
																				local.get $$t18.305
																				i32.store offset=1220
																				local.get $$t15.1
																				local.get $$t18.306
																				i32.store offset=1224
																				local.get $$t15.1
																				local.get $$t18.307
																				i32.store offset=1228
																				local.get $$t15.1
																				local.get $$t18.308
																				i32.store offset=1232
																				local.get $$t15.1
																				local.get $$t18.309
																				i32.store offset=1236
																				local.get $$t15.1
																				local.get $$t18.310
																				i32.store offset=1240
																				local.get $$t15.1
																				local.get $$t18.311
																				i32.store offset=1244
																				local.get $$t15.1
																				local.get $$t18.312
																				i32.store offset=1248
																				local.get $$t15.1
																				local.get $$t18.313
																				i32.store offset=1252
																				local.get $$t15.1
																				local.get $$t18.314
																				i32.store offset=1256
																				local.get $$t15.1
																				local.get $$t18.315
																				i32.store offset=1260
																				local.get $$t15.1
																				local.get $$t18.316
																				i32.store offset=1264
																				local.get $$t15.1
																				local.get $$t18.317
																				i32.store offset=1268
																				local.get $$t15.1
																				local.get $$t18.318
																				i32.store offset=1272
																				local.get $$t15.1
																				local.get $$t18.319
																				i32.store offset=1276
																				local.get $$t15.1
																				local.get $$t18.320
																				i32.store offset=1280
																				local.get $$t15.1
																				local.get $$t18.321
																				i32.store offset=1284
																				local.get $$t15.1
																				local.get $$t18.322
																				i32.store offset=1288
																				local.get $$t15.1
																				local.get $$t18.323
																				i32.store offset=1292
																				local.get $$t15.1
																				local.get $$t18.324
																				i32.store offset=1296
																				local.get $$t15.1
																				local.get $$t18.325
																				i32.store offset=1300
																				local.get $$t15.1
																				local.get $$t18.326
																				i32.store offset=1304
																				local.get $$t15.1
																				local.get $$t18.327
																				i32.store offset=1308
																				local.get $$t15.1
																				local.get $$t18.328
																				i32.store offset=1312
																				local.get $$t15.1
																				local.get $$t18.329
																				i32.store offset=1316
																				local.get $$t15.1
																				local.get $$t18.330
																				i32.store offset=1320
																				local.get $$t15.1
																				local.get $$t18.331
																				i32.store offset=1324
																				local.get $$t15.1
																				local.get $$t18.332
																				i32.store offset=1328
																				local.get $$t15.1
																				local.get $$t18.333
																				i32.store offset=1332
																				local.get $$t15.1
																				local.get $$t18.334
																				i32.store offset=1336
																				local.get $$t15.1
																				local.get $$t18.335
																				i32.store offset=1340
																				local.get $$t15.1
																				local.get $$t18.336
																				i32.store offset=1344
																				local.get $$t15.1
																				local.get $$t18.337
																				i32.store offset=1348
																				local.get $$t15.1
																				local.get $$t18.338
																				i32.store offset=1352
																				local.get $$t15.1
																				local.get $$t18.339
																				i32.store offset=1356
																				local.get $$t15.1
																				local.get $$t18.340
																				i32.store offset=1360
																				local.get $$t15.1
																				local.get $$t18.341
																				i32.store offset=1364
																				local.get $$t15.1
																				local.get $$t18.342
																				i32.store offset=1368
																				local.get $$t15.1
																				local.get $$t18.343
																				i32.store offset=1372
																				local.get $$t15.1
																				local.get $$t18.344
																				i32.store offset=1376
																				local.get $$t15.1
																				local.get $$t18.345
																				i32.store offset=1380
																				local.get $$t15.1
																				local.get $$t18.346
																				i32.store offset=1384
																				local.get $$t15.1
																				local.get $$t18.347
																				i32.store offset=1388
																				local.get $$t15.1
																				local.get $$t18.348
																				i32.store offset=1392
																				local.get $$t15.1
																				local.get $$t18.349
																				i32.store offset=1396
																				local.get $$t15.1
																				local.get $$t18.350
																				i32.store offset=1400
																				local.get $$t15.1
																				local.get $$t18.351
																				i32.store offset=1404
																				local.get $$t15.1
																				local.get $$t18.352
																				i32.store offset=1408
																				local.get $$t15.1
																				local.get $$t18.353
																				i32.store offset=1412
																				local.get $$t15.1
																				local.get $$t18.354
																				i32.store offset=1416
																				local.get $$t15.1
																				local.get $$t18.355
																				i32.store offset=1420
																				local.get $$t15.1
																				local.get $$t18.356
																				i32.store offset=1424
																				local.get $$t15.1
																				local.get $$t18.357
																				i32.store offset=1428
																				local.get $$t15.1
																				local.get $$t18.358
																				i32.store offset=1432
																				local.get $$t15.1
																				local.get $$t18.359
																				i32.store offset=1436
																				local.get $$t15.1
																				local.get $$t18.360
																				i32.store offset=1440
																				local.get $$t15.1
																				local.get $$t18.361
																				i32.store offset=1444
																				local.get $$t15.1
																				local.get $$t18.362
																				i32.store offset=1448
																				local.get $$t15.1
																				local.get $$t18.363
																				i32.store offset=1452
																				local.get $$t15.1
																				local.get $$t18.364
																				i32.store offset=1456
																				local.get $$t15.1
																				local.get $$t18.365
																				i32.store offset=1460
																				local.get $$t15.1
																				local.get $$t18.366
																				i32.store offset=1464
																				local.get $$t15.1
																				local.get $$t18.367
																				i32.store offset=1468
																				local.get $$t15.1
																				local.get $$t18.368
																				i32.store offset=1472
																				local.get $$t15.1
																				local.get $$t18.369
																				i32.store offset=1476
																				local.get $$t15.1
																				local.get $$t18.370
																				i32.store offset=1480
																				local.get $$t15.1
																				local.get $$t18.371
																				i32.store offset=1484
																				local.get $$t15.1
																				local.get $$t18.372
																				i32.store offset=1488
																				local.get $$t15.1
																				local.get $$t18.373
																				i32.store offset=1492
																				local.get $$t15.1
																				local.get $$t18.374
																				i32.store offset=1496
																				local.get $$t15.1
																				local.get $$t18.375
																				i32.store offset=1500
																				local.get $$t15.1
																				local.get $$t18.376
																				i32.store offset=1504
																				local.get $$t15.1
																				local.get $$t18.377
																				i32.store offset=1508
																				local.get $$t15.1
																				local.get $$t18.378
																				i32.store offset=1512
																				local.get $$t15.1
																				local.get $$t18.379
																				i32.store offset=1516
																				local.get $$t15.1
																				local.get $$t18.380
																				i32.store offset=1520
																				local.get $$t15.1
																				local.get $$t18.381
																				i32.store offset=1524
																				local.get $$t15.1
																				local.get $$t18.382
																				i32.store offset=1528
																				local.get $$t15.1
																				local.get $$t18.383
																				i32.store offset=1532
																				local.get $$t15.1
																				local.get $$t18.384
																				i32.store offset=1536
																				local.get $$t15.1
																				local.get $$t18.385
																				i32.store offset=1540
																				local.get $$t15.1
																				local.get $$t18.386
																				i32.store offset=1544
																				local.get $$t15.1
																				local.get $$t18.387
																				i32.store offset=1548
																				local.get $$t15.1
																				local.get $$t18.388
																				i32.store offset=1552
																				local.get $$t15.1
																				local.get $$t18.389
																				i32.store offset=1556
																				local.get $$t15.1
																				local.get $$t18.390
																				i32.store offset=1560
																				local.get $$t15.1
																				local.get $$t18.391
																				i32.store offset=1564
																				local.get $$t15.1
																				local.get $$t18.392
																				i32.store offset=1568
																				local.get $$t15.1
																				local.get $$t18.393
																				i32.store offset=1572
																				local.get $$t15.1
																				local.get $$t18.394
																				i32.store offset=1576
																				local.get $$t15.1
																				local.get $$t18.395
																				i32.store offset=1580
																				local.get $$t15.1
																				local.get $$t18.396
																				i32.store offset=1584
																				local.get $$t15.1
																				local.get $$t18.397
																				i32.store offset=1588
																				local.get $$t15.1
																				local.get $$t18.398
																				i32.store offset=1592
																				local.get $$t15.1
																				local.get $$t18.399
																				i32.store offset=1596
																				local.get $$t15.1
																				local.get $$t18.400
																				i32.store offset=1600
																				local.get $$t15.1
																				local.get $$t18.401
																				i32.store offset=1604
																				local.get $$t15.1
																				local.get $$t18.402
																				i32.store offset=1608
																				local.get $$t15.1
																				local.get $$t18.403
																				i32.store offset=1612
																				local.get $$t15.1
																				local.get $$t18.404
																				i32.store offset=1616
																				local.get $$t15.1
																				local.get $$t18.405
																				i32.store offset=1620
																				local.get $$t15.1
																				local.get $$t18.406
																				i32.store offset=1624
																				local.get $$t15.1
																				local.get $$t18.407
																				i32.store offset=1628
																				local.get $$t15.1
																				local.get $$t18.408
																				i32.store offset=1632
																				local.get $$t15.1
																				local.get $$t18.409
																				i32.store offset=1636
																				local.get $$t15.1
																				local.get $$t18.410
																				i32.store offset=1640
																				local.get $$t15.1
																				local.get $$t18.411
																				i32.store offset=1644
																				local.get $$t15.1
																				local.get $$t18.412
																				i32.store offset=1648
																				local.get $$t15.1
																				local.get $$t18.413
																				i32.store offset=1652
																				local.get $$t15.1
																				local.get $$t18.414
																				i32.store offset=1656
																				local.get $$t15.1
																				local.get $$t18.415
																				i32.store offset=1660
																				local.get $$t15.1
																				local.get $$t18.416
																				i32.store offset=1664
																				local.get $$t15.1
																				local.get $$t18.417
																				i32.store offset=1668
																				local.get $$t15.1
																				local.get $$t18.418
																				i32.store offset=1672
																				local.get $$t15.1
																				local.get $$t18.419
																				i32.store offset=1676
																				local.get $$t15.1
																				local.get $$t18.420
																				i32.store offset=1680
																				local.get $$t15.1
																				local.get $$t18.421
																				i32.store offset=1684
																				local.get $$t15.1
																				local.get $$t18.422
																				i32.store offset=1688
																				local.get $$t15.1
																				local.get $$t18.423
																				i32.store offset=1692
																				local.get $$t15.1
																				local.get $$t18.424
																				i32.store offset=1696
																				local.get $$t15.1
																				local.get $$t18.425
																				i32.store offset=1700
																				local.get $$t15.1
																				local.get $$t18.426
																				i32.store offset=1704
																				local.get $$t15.1
																				local.get $$t18.427
																				i32.store offset=1708
																				local.get $$t15.1
																				local.get $$t18.428
																				i32.store offset=1712
																				local.get $$t15.1
																				local.get $$t18.429
																				i32.store offset=1716
																				local.get $$t15.1
																				local.get $$t18.430
																				i32.store offset=1720
																				local.get $$t15.1
																				local.get $$t18.431
																				i32.store offset=1724
																				local.get $$t15.1
																				local.get $$t18.432
																				i32.store offset=1728
																				local.get $$t15.1
																				local.get $$t18.433
																				i32.store offset=1732
																				local.get $$t15.1
																				local.get $$t18.434
																				i32.store offset=1736
																				local.get $$t15.1
																				local.get $$t18.435
																				i32.store offset=1740
																				local.get $$t15.1
																				local.get $$t18.436
																				i32.store offset=1744
																				local.get $$t15.1
																				local.get $$t18.437
																				i32.store offset=1748
																				local.get $$t15.1
																				local.get $$t18.438
																				i32.store offset=1752
																				local.get $$t15.1
																				local.get $$t18.439
																				i32.store offset=1756
																				local.get $$t15.1
																				local.get $$t18.440
																				i32.store offset=1760
																				local.get $$t15.1
																				local.get $$t18.441
																				i32.store offset=1764
																				local.get $$t15.1
																				local.get $$t18.442
																				i32.store offset=1768
																				local.get $$t15.1
																				local.get $$t18.443
																				i32.store offset=1772
																				local.get $$t15.1
																				local.get $$t18.444
																				i32.store offset=1776
																				local.get $$t15.1
																				local.get $$t18.445
																				i32.store offset=1780
																				local.get $$t15.1
																				local.get $$t18.446
																				i32.store offset=1784
																				local.get $$t15.1
																				local.get $$t18.447
																				i32.store offset=1788
																				local.get $$t15.1
																				local.get $$t18.448
																				i32.store offset=1792
																				local.get $$t15.1
																				local.get $$t18.449
																				i32.store offset=1796
																				local.get $$t15.1
																				local.get $$t18.450
																				i32.store offset=1800
																				local.get $$t15.1
																				local.get $$t18.451
																				i32.store offset=1804
																				local.get $$t15.1
																				local.get $$t18.452
																				i32.store offset=1808
																				local.get $$t15.1
																				local.get $$t18.453
																				i32.store offset=1812
																				local.get $$t15.1
																				local.get $$t18.454
																				i32.store offset=1816
																				local.get $$t15.1
																				local.get $$t18.455
																				i32.store offset=1820
																				local.get $$t15.1
																				local.get $$t18.456
																				i32.store offset=1824
																				local.get $$t15.1
																				local.get $$t18.457
																				i32.store offset=1828
																				local.get $$t15.1
																				local.get $$t18.458
																				i32.store offset=1832
																				local.get $$t15.1
																				local.get $$t18.459
																				i32.store offset=1836
																				local.get $$t15.1
																				local.get $$t18.460
																				i32.store offset=1840
																				local.get $$t15.1
																				local.get $$t18.461
																				i32.store offset=1844
																				local.get $$t15.1
																				local.get $$t18.462
																				i32.store offset=1848
																				local.get $$t15.1
																				local.get $$t18.463
																				i32.store offset=1852
																				local.get $$t15.1
																				local.get $$t18.464
																				i32.store offset=1856
																				local.get $$t15.1
																				local.get $$t18.465
																				i32.store offset=1860
																				local.get $$t15.1
																				local.get $$t18.466
																				i32.store offset=1864
																				local.get $$t15.1
																				local.get $$t18.467
																				i32.store offset=1868
																				local.get $$t16.1
																				local.get $$t19.0
																				i32.store16
																				local.get $$t16.1
																				local.get $$t19.1
																				i32.store16 offset=2
																				local.get $$t16.1
																				local.get $$t19.2
																				i32.store16 offset=4
																				local.get $$t16.1
																				local.get $$t19.3
																				i32.store16 offset=6
																				local.get $$t16.1
																				local.get $$t19.4
																				i32.store16 offset=8
																				local.get $$t16.1
																				local.get $$t19.5
																				i32.store16 offset=10
																				local.get $$t16.1
																				local.get $$t19.6
																				i32.store16 offset=12
																				local.get $$t16.1
																				local.get $$t19.7
																				i32.store16 offset=14
																				local.get $$t16.1
																				local.get $$t19.8
																				i32.store16 offset=16
																				local.get $$t16.1
																				local.get $$t19.9
																				i32.store16 offset=18
																				local.get $$t16.1
																				local.get $$t19.10
																				i32.store16 offset=20
																				local.get $$t16.1
																				local.get $$t19.11
																				i32.store16 offset=22
																				local.get $$t16.1
																				local.get $$t19.12
																				i32.store16 offset=24
																				local.get $$t16.1
																				local.get $$t19.13
																				i32.store16 offset=26
																				local.get $$t16.1
																				local.get $$t19.14
																				i32.store16 offset=28
																				local.get $$t16.1
																				local.get $$t19.15
																				i32.store16 offset=30
																				local.get $$t16.1
																				local.get $$t19.16
																				i32.store16 offset=32
																				local.get $$t16.1
																				local.get $$t19.17
																				i32.store16 offset=34
																				local.get $$t16.1
																				local.get $$t19.18
																				i32.store16 offset=36
																				local.get $$t16.1
																				local.get $$t19.19
																				i32.store16 offset=38
																				local.get $$t16.1
																				local.get $$t19.20
																				i32.store16 offset=40
																				local.get $$t16.1
																				local.get $$t19.21
																				i32.store16 offset=42
																				local.get $$t16.1
																				local.get $$t19.22
																				i32.store16 offset=44
																				local.get $$t16.1
																				local.get $$t19.23
																				i32.store16 offset=46
																				local.get $$t16.1
																				local.get $$t19.24
																				i32.store16 offset=48
																				local.get $$t16.1
																				local.get $$t19.25
																				i32.store16 offset=50
																				local.get $$t16.1
																				local.get $$t19.26
																				i32.store16 offset=52
																				local.get $$t16.1
																				local.get $$t19.27
																				i32.store16 offset=54
																				local.get $$t16.1
																				local.get $$t19.28
																				i32.store16 offset=56
																				local.get $$t16.1
																				local.get $$t19.29
																				i32.store16 offset=58
																				local.get $$t16.1
																				local.get $$t19.30
																				i32.store16 offset=60
																				local.get $$t16.1
																				local.get $$t19.31
																				i32.store16 offset=62
																				local.get $$t16.1
																				local.get $$t19.32
																				i32.store16 offset=64
																				local.get $$t16.1
																				local.get $$t19.33
																				i32.store16 offset=66
																				local.get $$t16.1
																				local.get $$t19.34
																				i32.store16 offset=68
																				local.get $$t16.1
																				local.get $$t19.35
																				i32.store16 offset=70
																				local.get $$t16.1
																				local.get $$t19.36
																				i32.store16 offset=72
																				local.get $$t16.1
																				local.get $$t19.37
																				i32.store16 offset=74
																				local.get $$t16.1
																				local.get $$t19.38
																				i32.store16 offset=76
																				local.get $$t16.1
																				local.get $$t19.39
																				i32.store16 offset=78
																				local.get $$t16.1
																				local.get $$t19.40
																				i32.store16 offset=80
																				local.get $$t16.1
																				local.get $$t19.41
																				i32.store16 offset=82
																				local.get $$t16.1
																				local.get $$t19.42
																				i32.store16 offset=84
																				local.get $$t16.1
																				local.get $$t19.43
																				i32.store16 offset=86
																				local.get $$t16.1
																				local.get $$t19.44
																				i32.store16 offset=88
																				local.get $$t16.1
																				local.get $$t19.45
																				i32.store16 offset=90
																				local.get $$t16.1
																				local.get $$t19.46
																				i32.store16 offset=92
																				local.get $$t16.1
																				local.get $$t19.47
																				i32.store16 offset=94
																				local.get $$t16.1
																				local.get $$t19.48
																				i32.store16 offset=96
																				local.get $$t16.1
																				local.get $$t19.49
																				i32.store16 offset=98
																				local.get $$t16.1
																				local.get $$t19.50
																				i32.store16 offset=100
																				local.get $$t16.1
																				local.get $$t19.51
																				i32.store16 offset=102
																				local.get $$t16.1
																				local.get $$t19.52
																				i32.store16 offset=104
																				local.get $$t16.1
																				local.get $$t19.53
																				i32.store16 offset=106
																				local.get $$t16.1
																				local.get $$t19.54
																				i32.store16 offset=108
																				local.get $$t16.1
																				local.get $$t19.55
																				i32.store16 offset=110
																				local.get $$t16.1
																				local.get $$t19.56
																				i32.store16 offset=112
																				local.get $$t16.1
																				local.get $$t19.57
																				i32.store16 offset=114
																				local.get $$t16.1
																				local.get $$t19.58
																				i32.store16 offset=116
																				local.get $$t16.1
																				local.get $$t19.59
																				i32.store16 offset=118
																				local.get $$t16.1
																				local.get $$t19.60
																				i32.store16 offset=120
																				local.get $$t16.1
																				local.get $$t19.61
																				i32.store16 offset=122
																				local.get $$t16.1
																				local.get $$t19.62
																				i32.store16 offset=124
																				local.get $$t16.1
																				local.get $$t19.63
																				i32.store16 offset=126
																				local.get $$t16.1
																				local.get $$t19.64
																				i32.store16 offset=128
																				local.get $$t16.1
																				local.get $$t19.65
																				i32.store16 offset=130
																				local.get $$t16.1
																				local.get $$t19.66
																				i32.store16 offset=132
																				local.get $$t16.1
																				local.get $$t19.67
																				i32.store16 offset=134
																				local.get $$t16.1
																				local.get $$t19.68
																				i32.store16 offset=136
																				local.get $$t16.1
																				local.get $$t19.69
																				i32.store16 offset=138
																				local.get $$t16.1
																				local.get $$t19.70
																				i32.store16 offset=140
																				local.get $$t16.1
																				local.get $$t19.71
																				i32.store16 offset=142
																				local.get $$t16.1
																				local.get $$t19.72
																				i32.store16 offset=144
																				local.get $$t16.1
																				local.get $$t19.73
																				i32.store16 offset=146
																				local.get $$t16.1
																				local.get $$t19.74
																				i32.store16 offset=148
																				local.get $$t16.1
																				local.get $$t19.75
																				i32.store16 offset=150
																				local.get $$t16.1
																				local.get $$t19.76
																				i32.store16 offset=152
																				local.get $$t16.1
																				local.get $$t19.77
																				i32.store16 offset=154
																				local.get $$t16.1
																				local.get $$t19.78
																				i32.store16 offset=156
																				local.get $$t16.1
																				local.get $$t19.79
																				i32.store16 offset=158
																				local.get $$t16.1
																				local.get $$t19.80
																				i32.store16 offset=160
																				local.get $$t16.1
																				local.get $$t19.81
																				i32.store16 offset=162
																				local.get $$t16.1
																				local.get $$t19.82
																				i32.store16 offset=164
																				local.get $$t16.1
																				local.get $$t19.83
																				i32.store16 offset=166
																				local.get $$t16.1
																				local.get $$t19.84
																				i32.store16 offset=168
																				local.get $$t16.1
																				local.get $$t19.85
																				i32.store16 offset=170
																				local.get $$t16.1
																				local.get $$t19.86
																				i32.store16 offset=172
																				local.get $$t16.1
																				local.get $$t19.87
																				i32.store16 offset=174
																				local.get $$t16.1
																				local.get $$t19.88
																				i32.store16 offset=176
																				local.get $$t16.1
																				local.get $$t19.89
																				i32.store16 offset=178
																				local.get $$t16.1
																				local.get $$t19.90
																				i32.store16 offset=180
																				local.get $$t16.1
																				local.get $$t19.91
																				i32.store16 offset=182
																				local.get $$t16.1
																				local.get $$t19.92
																				i32.store16 offset=184
																				local.get $$t16.1
																				local.get $$t19.93
																				i32.store16 offset=186
																				local.get $$t16.1
																				local.get $$t19.94
																				i32.store16 offset=188
																				local.get $$t15.0
																				call $runtime.Block.Retain
																				local.get $$t15.1
																				i32.const 4
																				i32.const 0
																				i32.mul
																				i32.add
																				i32.const 468
																				i32.const 0
																				i32.sub
																				i32.const 468
																				i32.const 0
																				i32.sub
																				local.set $$t20.3
																				local.set $$t20.2
																				local.set $$t20.1
																				local.get $$t20.0
																				call $runtime.Block.Release
																				local.set $$t20.0
																				local.get $$t20.0
																				local.get $$t20.1
																				local.get $$t20.2
																				local.get $$t20.3
																				local.get $$t17
																				call $strconv.bsearch32
																				local.set $$t21
																				local.get $$t21
																				i32.const 468
																				i32.ge_s
																				local.set $$t22
																				local.get $$t22
																				if
																					br $$Block_17
																				else
																					br $$Block_20
																				end
																			end
																			i32.const 11
																			local.set $$current_block
																			local.get $r
																			i32.const 65536
																			i32.lt_s
																			local.set $$t23
																			local.get $$t23
																			if
																				i32.const 9
																				local.set $$block_selector
																				br $$BlockDisp
																			else
																				i32.const 10
																				local.set $$block_selector
																				br $$BlockDisp
																			end
																		end
																		i32.const 12
																		local.set $$current_block
																		i32.const 0
																		local.set $$ret_0
																		br $$BlockFnBody
																	end
																	i32.const 13
																	local.set $$current_block
																	local.get $$t8.0
																	call $runtime.Block.Retain
																	local.get $$t8.1
																	i32.const 2
																	i32.const 0
																	i32.mul
																	i32.add
																	i32.const 132
																	i32.const 0
																	i32.sub
																	i32.const 132
																	i32.const 0
																	i32.sub
																	local.set $$t24.3
																	local.set $$t24.2
																	local.set $$t24.1
																	local.get $$t24.0
																	call $runtime.Block.Release
																	local.set $$t24.0
																	local.get $$t24.0
																	local.get $$t24.1
																	local.get $$t24.2
																	local.get $$t24.3
																	local.get $$t9
																	call $strconv.bsearch16
																	local.set $$t25
																	local.get $$t25
																	i32.const 132
																	i32.ge_s
																	local.set $$t26
																	local.get $$t26
																	if
																		br $$Block_16
																	else
																		br $$Block_15
																	end
																end
																i32.const 14
																local.set $$current_block
																local.get $$t13
																i32.const 1
																i32.or
																local.set $$t27
																local.get $$t7.0
																call $runtime.Block.Retain
																local.get $$t7.1
																i32.const 2
																local.get $$t27
																i32.mul
																i32.add
																local.set $$t28.1
																local.get $$t28.0
																call $runtime.Block.Release
																local.set $$t28.0
																local.get $$t28.1
																i32.load16_u
																local.set $$t29
																local.get $$t29
																local.get $$t9
																i32.lt_u
																local.set $$t30
																local.get $$t30
																if
																	i32.const 12
																	local.set $$block_selector
																	br $$BlockDisp
																else
																	i32.const 13
																	local.set $$block_selector
																	br $$BlockDisp
																end
															end
															i32.const 15
															local.set $$current_block
															local.get $$t13
															i32.const 1
															i32.const -1
															i32.xor
															i32.and
															local.set $$t31
															local.get $$t7.0
															call $runtime.Block.Retain
															local.get $$t7.1
															i32.const 2
															local.get $$t31
															i32.mul
															i32.add
															local.set $$t32.1
															local.get $$t32.0
															call $runtime.Block.Release
															local.set $$t32.0
															local.get $$t32.1
															i32.load16_u
															local.set $$t33
															local.get $$t9
															local.get $$t33
															i32.lt_u
															local.set $$t34
															local.get $$t34
															if
																i32.const 12
																local.set $$block_selector
																br $$BlockDisp
															else
																i32.const 14
																local.set $$block_selector
																br $$BlockDisp
															end
														end
														i32.const 16
														local.set $$current_block
														local.get $$t8.0
														call $runtime.Block.Retain
														local.get $$t8.1
														i32.const 2
														local.get $$t25
														i32.mul
														i32.add
														local.set $$t35.1
														local.get $$t35.0
														call $runtime.Block.Release
														local.set $$t35.0
														local.get $$t35.1
														i32.load16_u
														local.set $$t36
														local.get $$t36
														local.get $$t9
														i32.eq
														i32.eqz
														local.set $$t37
														br $$Block_16
													end
													local.get $$current_block
													i32.const 13
													i32.eq
													if(result i32)
														i32.const 1
													else
														local.get $$t37
													end
													local.set $$t38
													i32.const 17
													local.set $$current_block
													local.get $$t38
													local.set $$ret_0
													br $$BlockFnBody
												end
												i32.const 18
												local.set $$current_block
												i32.const 0
												local.set $$ret_0
												br $$BlockFnBody
											end
											i32.const 19
											local.set $$current_block
											local.get $r
											i32.const 131072
											i32.ge_s
											local.set $$t39
											local.get $$t39
											if
												br $$Block_21
											else
												br $$Block_22
											end
										end
										i32.const 20
										local.set $$current_block
										local.get $$t21
										i32.const 1
										i32.or
										local.set $$t40
										local.get $$t15.0
										call $runtime.Block.Retain
										local.get $$t15.1
										i32.const 4
										local.get $$t40
										i32.mul
										i32.add
										local.set $$t41.1
										local.get $$t41.0
										call $runtime.Block.Release
										local.set $$t41.0
										local.get $$t41.1
										i32.load
										local.set $$t42
										local.get $$t42
										local.get $$t17
										i32.lt_u
										local.set $$t43
										local.get $$t43
										if
											i32.const 18
											local.set $$block_selector
											br $$BlockDisp
										else
											i32.const 19
											local.set $$block_selector
											br $$BlockDisp
										end
									end
									i32.const 21
									local.set $$current_block
									local.get $$t21
									i32.const 1
									i32.const -1
									i32.xor
									i32.and
									local.set $$t44
									local.get $$t15.0
									call $runtime.Block.Retain
									local.get $$t15.1
									i32.const 4
									local.get $$t44
									i32.mul
									i32.add
									local.set $$t45.1
									local.get $$t45.0
									call $runtime.Block.Release
									local.set $$t45.0
									local.get $$t45.1
									i32.load
									local.set $$t46
									local.get $$t17
									local.get $$t46
									i32.lt_u
									local.set $$t47
									local.get $$t47
									if
										i32.const 18
										local.set $$block_selector
										br $$BlockDisp
									else
										i32.const 20
										local.set $$block_selector
										br $$BlockDisp
									end
								end
								i32.const 22
								local.set $$current_block
								i32.const 1
								local.set $$ret_0
								br $$BlockFnBody
							end
							i32.const 23
							local.set $$current_block
							local.get $r
							i32.const 65536
							i32.sub
							local.set $$t48
							local.get $$t16.0
							call $runtime.Block.Retain
							local.get $$t16.1
							i32.const 2
							i32.const 0
							i32.mul
							i32.add
							i32.const 95
							i32.const 0
							i32.sub
							i32.const 95
							i32.const 0
							i32.sub
							local.set $$t49.3
							local.set $$t49.2
							local.set $$t49.1
							local.get $$t49.0
							call $runtime.Block.Release
							local.set $$t49.0
							local.get $$t48
							i32.const 65535
							i32.and
							local.set $$t50
							local.get $$t49.0
							local.get $$t49.1
							local.get $$t49.2
							local.get $$t49.3
							local.get $$t50
							call $strconv.bsearch16
							local.set $$t51
							local.get $$t51
							i32.const 95
							i32.ge_s
							local.set $$t52
							local.get $$t52
							if
								br $$Block_24
							else
								br $$Block_23
							end
						end
						i32.const 24
						local.set $$current_block
						local.get $$t16.0
						call $runtime.Block.Retain
						local.get $$t16.1
						i32.const 2
						local.get $$t51
						i32.mul
						i32.add
						local.set $$t53.1
						local.get $$t53.0
						call $runtime.Block.Release
						local.set $$t53.0
						local.get $$t53.1
						i32.load16_u
						local.set $$t54
						local.get $$t48
						i32.const 65535
						i32.and
						local.set $$t55
						local.get $$t54
						local.get $$t55
						i32.eq
						i32.eqz
						local.set $$t56
						br $$Block_24
					end
					local.get $$current_block
					i32.const 23
					i32.eq
					if(result i32)
						i32.const 1
					else
						local.get $$t56
					end
					local.set $$t57
					i32.const 25
					local.set $$current_block
					local.get $$t57
					local.set $$ret_0
					br $$BlockFnBody
				end
			end
		end
		local.get $$ret_0
		local.get $$t7.0
		call $runtime.Block.Release
		local.get $$t8.0
		call $runtime.Block.Release
		local.get $$t12.0
		call $runtime.Block.Release
		local.get $$t15.0
		call $runtime.Block.Release
		local.get $$t16.0
		call $runtime.Block.Release
		local.get $$t20.0
		call $runtime.Block.Release
		local.get $$t24.0
		call $runtime.Block.Release
		local.get $$t28.0
		call $runtime.Block.Release
		local.get $$t32.0
		call $runtime.Block.Release
		local.get $$t35.0
		call $runtime.Block.Release
		local.get $$t41.0
		call $runtime.Block.Release
		local.get $$t45.0
		call $runtime.Block.Release
		local.get $$t49.0
		call $runtime.Block.Release
		local.get $$t53.0
		call $runtime.Block.Release
	)
	(func $strconv.Itoa (param $i i32) (result i32 i32 i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$ret_0.0 i32)
		(local $$ret_0.1 i32)
		(local $$ret_0.2 i32)
		(local $$t0 i64)
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
					local.get $i
					i64.extend_i32_s
					local.set $$t0
					local.get $$t0
					i32.const 10
					call $strconv.FormatInt
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
		local.get $$t1.0
		call $runtime.Block.Release
	)
	(func $strconv.Quote (param $s.0 i32) (param $s.1 i32) (param $s.2 i32) (result i32 i32 i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$ret_0.0 i32)
		(local $$ret_0.1 i32)
		(local $$ret_0.2 i32)
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
					local.get $s.0
					local.get $s.1
					local.get $s.2
					i32.const 34
					i32.const 0
					i32.const 0
					call $strconv.quoteWith
					local.set $$t0.2
					local.set $$t0.1
					local.get $$t0.0
					call $runtime.Block.Release
					local.set $$t0.0
					local.get $$t0.0
					call $runtime.Block.Retain
					local.get $$t0.1
					local.get $$t0.2
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
	)
	(func $strconv.appendEscapedRune (param $buf.0 i32) (param $buf.1 i32) (param $buf.2 i32) (param $buf.3 i32) (param $r i32) (param $quote i32) (param $ASCIIonly i32) (param $graphicOnly i32) (result i32 i32 i32 i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$ret_0.0 i32)
		(local $$ret_0.1 i32)
		(local $$ret_0.2 i32)
		(local $$ret_0.3 i32)
		(local $$t0.0 i32)
		(local $$t0.1 i32)
		(local $$t1 i32)
		(local $$t2 i32)
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
		(local $$t7 i32)
		(local $$t8.0 i32)
		(local $$t8.1 i32)
		(local $$t9.0 i32)
		(local $$t9.1 i32)
		(local $$t10.0 i32)
		(local $$t10.1 i32)
		(local $$t10.2 i32)
		(local $$t10.3 i32)
		(local $$t11.0 i32)
		(local $$t11.1 i32)
		(local $$t11.2 i32)
		(local $$t11.3 i32)
		(local $$t12 i32)
		(local $$t13 i32)
		(local $$t14 i32)
		(local $$t15 i32)
		(local $$t16 i32)
		(local $$t17.0 i32)
		(local $$t17.1 i32)
		(local $$t18.0 i32)
		(local $$t18.1 i32)
		(local $$t19.0 i32)
		(local $$t19.1 i32)
		(local $$t19.2 i32)
		(local $$t19.3 i32)
		(local $$t20.0 i32)
		(local $$t20.1 i32)
		(local $$t20.2 i32)
		(local $$t20.3 i32)
		(local $$t21 i32)
		(local $$t22.0 i32)
		(local $$t22.1 i32)
		(local $$t22.2 i32)
		(local $$t22.3 i32)
		(local $$t23 i32)
		(local $$t24.0 i32)
		(local $$t24.1 i32)
		(local $$t24.2 i32)
		(local $$t24.3 i32)
		(local $$t25.0 i32)
		(local $$t25.1 i32)
		(local $$t25.2 i32)
		(local $$t25.3 i32)
		(local $$t26 i32)
		(local $$t27.0 i32)
		(local $$t27.1 i32)
		(local $$t27.2 i32)
		(local $$t27.3 i32)
		(local $$t28.0 i32)
		(local $$t28.1 i32)
		(local $$t28.2 i32)
		(local $$t28.3 i32)
		(local $$t29.0 i32)
		(local $$t29.1 i32)
		(local $$t29.2 i32)
		(local $$t29.3 i32)
		(local $$t30.0 i32)
		(local $$t30.1 i32)
		(local $$t30.2 i32)
		(local $$t30.3 i32)
		(local $$t31.0 i32)
		(local $$t31.1 i32)
		(local $$t31.2 i32)
		(local $$t31.3 i32)
		(local $$t32.0 i32)
		(local $$t32.1 i32)
		(local $$t32.2 i32)
		(local $$t32.3 i32)
		(local $$t33.0 i32)
		(local $$t33.1 i32)
		(local $$t33.2 i32)
		(local $$t33.3 i32)
		(local $$t34.0 i32)
		(local $$t34.1 i32)
		(local $$t34.2 i32)
		(local $$t34.3 i32)
		(local $$t35.0 i32)
		(local $$t35.1 i32)
		(local $$t35.2 i32)
		(local $$t35.3 i32)
		(local $$t36.0 i32)
		(local $$t36.1 i32)
		(local $$t36.2 i32)
		(local $$t36.3 i32)
		(local $$t37.0 i32)
		(local $$t37.1 i32)
		(local $$t37.2 i32)
		(local $$t37.3 i32)
		(local $$t38.0 i32)
		(local $$t38.1 i32)
		(local $$t38.2 i32)
		(local $$t38.3 i32)
		(local $$t39 i32)
		(local $$t40 i32)
		(local $$t41 i32)
		(local $$t42 i32)
		(local $$t43 i32)
		(local $$t44 i32)
		(local $$t45 i32)
		(local $$t46.0 i32)
		(local $$t46.1 i32)
		(local $$t46.2 i32)
		(local $$t46.3 i32)
		(local $$t47 i32)
		(local $$t48 i32)
		(local $$t49 i32)
		(local $$t50.0 i32)
		(local $$t50.1 i32)
		(local $$t51.0 i32)
		(local $$t51.1 i32)
		(local $$t52.0 i32)
		(local $$t52.1 i32)
		(local $$t52.2 i32)
		(local $$t52.3 i32)
		(local $$t53.0 i32)
		(local $$t53.1 i32)
		(local $$t53.2 i32)
		(local $$t53.3 i32)
		(local $$t54 i32)
		(local $$t55 i32)
		(local $$t56 i32)
		(local $$t57.0 i32)
		(local $$t57.1 i32)
		(local $$t58.0 i32)
		(local $$t58.1 i32)
		(local $$t59.0 i32)
		(local $$t59.1 i32)
		(local $$t59.2 i32)
		(local $$t59.3 i32)
		(local $$t60.0 i32)
		(local $$t60.1 i32)
		(local $$t60.2 i32)
		(local $$t60.3 i32)
		(local $$t61 i32)
		(local $$t62.0 i32)
		(local $$t62.1 i32)
		(local $$t62.2 i32)
		(local $$t62.3 i32)
		(local $$t63 i32)
		(local $$t64 i32)
		(local $$t65 i32)
		(local $$t66 i32)
		(local $$t67 i32)
		(local $$t68 i32)
		(local $$t69.0 i32)
		(local $$t69.1 i32)
		(local $$t70.0 i32)
		(local $$t70.1 i32)
		(local $$t71.0 i32)
		(local $$t71.1 i32)
		(local $$t71.2 i32)
		(local $$t71.3 i32)
		(local $$t72.0 i32)
		(local $$t72.1 i32)
		(local $$t72.2 i32)
		(local $$t72.3 i32)
		(local $$t73 i32)
		(local $$t74 i32)
		(local $$t75.0 i32)
		(local $$t75.1 i32)
		(local $$t75.2 i32)
		(local $$t75.3 i32)
		(local $$t76 i32)
		(local $$t77 i32)
		(local $$t78 i32)
		(local $$t79 i32)
		(local $$t80 i32)
		(local $$t81.0 i32)
		(local $$t81.1 i32)
		(local $$t82.0 i32)
		(local $$t82.1 i32)
		(local $$t83.0 i32)
		(local $$t83.1 i32)
		(local $$t83.2 i32)
		(local $$t83.3 i32)
		(local $$t84.0 i32)
		(local $$t84.1 i32)
		(local $$t84.2 i32)
		(local $$t84.3 i32)
		(local $$t85 i32)
		(local $$t86 i32)
		(local $$t87 i32)
		(local $$t88 i32)
		(local $$t89 i32)
		(local $$t90 i32)
		(local $$t91 i32)
		(local $$t92.0 i32)
		(local $$t92.1 i32)
		(local $$t93.0 i32)
		(local $$t93.1 i32)
		(local $$t94.0 i32)
		(local $$t94.1 i32)
		(local $$t94.2 i32)
		(local $$t94.3 i32)
		(local $$t95.0 i32)
		(local $$t95.1 i32)
		(local $$t95.2 i32)
		(local $$t95.3 i32)
		(local $$t96 i32)
		(local $$t97 i32)
		block $$BlockFnBody
			loop $$BlockDisp
				block $$Block_38
					block $$Block_37
						block $$Block_36
							block $$Block_35
								block $$Block_34
									block $$Block_33
										block $$Block_32
											block $$Block_31
												block $$Block_30
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
																																												br_table  0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31 32 33 34 35 36 37 38 0
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
																																											local.get $quote
																																											local.set $$t1
																																											local.get $r
																																											local.get $$t1
																																											i32.eq
																																											local.set $$t2
																																											local.get $$t2
																																											if
																																												br $$Block_0
																																											else
																																												br $$Block_2
																																											end
																																										end
																																										i32.const 1
																																										local.set $$current_block
																																										i32.const 17
																																										call $runtime.HeapAlloc
																																										i32.const 1
																																										i32.const 0
																																										i32.const 1
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
																																										i32.const 1
																																										i32.const 0
																																										i32.mul
																																										i32.add
																																										local.set $$t4.1
																																										local.get $$t4.0
																																										call $runtime.Block.Release
																																										local.set $$t4.0
																																										local.get $$t4.1
																																										i32.const 92
																																										i32.store8 align=1
																																										local.get $$t3.0
																																										call $runtime.Block.Retain
																																										local.get $$t3.1
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
																																										local.set $$t5.3
																																										local.set $$t5.2
																																										local.set $$t5.1
																																										local.get $$t5.0
																																										call $runtime.Block.Release
																																										local.set $$t5.0
																																										local.get $buf.0
																																										local.get $buf.1
																																										local.get $buf.2
																																										local.get $buf.3
																																										local.get $$t5.0
																																										local.get $$t5.1
																																										local.get $$t5.2
																																										local.get $$t5.3
																																										call $$u8.$slice.append
																																										local.set $$t6.3
																																										local.set $$t6.2
																																										local.set $$t6.1
																																										local.get $$t6.0
																																										call $runtime.Block.Release
																																										local.set $$t6.0
																																										local.get $r
																																										i32.const 255
																																										i32.and
																																										local.set $$t7
																																										i32.const 17
																																										call $runtime.HeapAlloc
																																										i32.const 1
																																										i32.const 0
																																										i32.const 1
																																										call $runtime.Block.Init
																																										call $runtime.DupI32
																																										i32.const 16
																																										i32.add
																																										local.set $$t8.1
																																										local.get $$t8.0
																																										call $runtime.Block.Release
																																										local.set $$t8.0
																																										local.get $$t8.0
																																										call $runtime.Block.Retain
																																										local.get $$t8.1
																																										i32.const 1
																																										i32.const 0
																																										i32.mul
																																										i32.add
																																										local.set $$t9.1
																																										local.get $$t9.0
																																										call $runtime.Block.Release
																																										local.set $$t9.0
																																										local.get $$t9.1
																																										local.get $$t7
																																										i32.store8 align=1
																																										local.get $$t8.0
																																										call $runtime.Block.Retain
																																										local.get $$t8.1
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
																																										local.set $$t10.3
																																										local.set $$t10.2
																																										local.set $$t10.1
																																										local.get $$t10.0
																																										call $runtime.Block.Release
																																										local.set $$t10.0
																																										local.get $$t6.0
																																										local.get $$t6.1
																																										local.get $$t6.2
																																										local.get $$t6.3
																																										local.get $$t10.0
																																										local.get $$t10.1
																																										local.get $$t10.2
																																										local.get $$t10.3
																																										call $$u8.$slice.append
																																										local.set $$t11.3
																																										local.set $$t11.2
																																										local.set $$t11.1
																																										local.get $$t11.0
																																										call $runtime.Block.Release
																																										local.set $$t11.0
																																										local.get $$t11.0
																																										call $runtime.Block.Retain
																																										local.get $$t11.1
																																										local.get $$t11.2
																																										local.get $$t11.3
																																										local.set $$ret_0.3
																																										local.set $$ret_0.2
																																										local.set $$ret_0.1
																																										local.get $$ret_0.0
																																										call $runtime.Block.Release
																																										local.set $$ret_0.0
																																										br $$BlockFnBody
																																									end
																																									i32.const 2
																																									local.set $$current_block
																																									local.get $ASCIIonly
																																									if
																																										br $$Block_3
																																									else
																																										br $$Block_5
																																									end
																																								end
																																								i32.const 3
																																								local.set $$current_block
																																								local.get $r
																																								i32.const 92
																																								i32.eq
																																								local.set $$t12
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
																																							end
																																							i32.const 4
																																							local.set $$current_block
																																							local.get $r
																																							i32.const 128
																																							i32.lt_s
																																							local.set $$t13
																																							local.get $$t13
																																							if
																																								br $$Block_7
																																							else
																																								br $$Block_4
																																							end
																																						end
																																						i32.const 5
																																						local.set $$current_block
																																						local.get $r
																																						i32.const 7
																																						i32.eq
																																						local.set $$t14
																																						local.get $$t14
																																						if
																																							br $$Block_12
																																						else
																																							br $$Block_14
																																						end
																																					end
																																					i32.const 6
																																					local.set $$current_block
																																					local.get $r
																																					call $strconv.IsPrint
																																					local.set $$t15
																																					local.get $$t15
																																					if
																																						br $$Block_8
																																					else
																																						br $$Block_9
																																					end
																																				end
																																				i32.const 7
																																				local.set $$current_block
																																				local.get $r
																																				i32.const 255
																																				i32.and
																																				local.set $$t16
																																				i32.const 17
																																				call $runtime.HeapAlloc
																																				i32.const 1
																																				i32.const 0
																																				i32.const 1
																																				call $runtime.Block.Init
																																				call $runtime.DupI32
																																				i32.const 16
																																				i32.add
																																				local.set $$t17.1
																																				local.get $$t17.0
																																				call $runtime.Block.Release
																																				local.set $$t17.0
																																				local.get $$t17.0
																																				call $runtime.Block.Retain
																																				local.get $$t17.1
																																				i32.const 1
																																				i32.const 0
																																				i32.mul
																																				i32.add
																																				local.set $$t18.1
																																				local.get $$t18.0
																																				call $runtime.Block.Release
																																				local.set $$t18.0
																																				local.get $$t18.1
																																				local.get $$t16
																																				i32.store8 align=1
																																				local.get $$t17.0
																																				call $runtime.Block.Retain
																																				local.get $$t17.1
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
																																				local.set $$t19.3
																																				local.set $$t19.2
																																				local.set $$t19.1
																																				local.get $$t19.0
																																				call $runtime.Block.Release
																																				local.set $$t19.0
																																				local.get $buf.0
																																				local.get $buf.1
																																				local.get $buf.2
																																				local.get $buf.3
																																				local.get $$t19.0
																																				local.get $$t19.1
																																				local.get $$t19.2
																																				local.get $$t19.3
																																				call $$u8.$slice.append
																																				local.set $$t20.3
																																				local.set $$t20.2
																																				local.set $$t20.1
																																				local.get $$t20.0
																																				call $runtime.Block.Release
																																				local.set $$t20.0
																																				local.get $$t20.0
																																				call $runtime.Block.Retain
																																				local.get $$t20.1
																																				local.get $$t20.2
																																				local.get $$t20.3
																																				local.set $$ret_0.3
																																				local.set $$ret_0.2
																																				local.set $$ret_0.1
																																				local.get $$ret_0.0
																																				call $runtime.Block.Release
																																				local.set $$ret_0.0
																																				br $$BlockFnBody
																																			end
																																			i32.const 8
																																			local.set $$current_block
																																			local.get $r
																																			call $strconv.IsPrint
																																			local.set $$t21
																																			local.get $$t21
																																			if
																																				i32.const 7
																																				local.set $$block_selector
																																				br $$BlockDisp
																																			else
																																				i32.const 5
																																				local.set $$block_selector
																																				br $$BlockDisp
																																			end
																																		end
																																		i32.const 9
																																		local.set $$current_block
																																		local.get $$t0.0
																																		call $runtime.Block.Retain
																																		local.get $$t0.1
																																		i32.const 1
																																		i32.const 0
																																		i32.mul
																																		i32.add
																																		i32.const 4
																																		i32.const 0
																																		i32.sub
																																		i32.const 4
																																		i32.const 0
																																		i32.sub
																																		local.set $$t22.3
																																		local.set $$t22.2
																																		local.set $$t22.1
																																		local.get $$t22.0
																																		call $runtime.Block.Release
																																		local.set $$t22.0
																																		local.get $$t22.0
																																		local.get $$t22.1
																																		local.get $$t22.2
																																		local.get $$t22.3
																																		local.get $r
																																		call $unicode$utf8.EncodeRune
																																		local.set $$t23
																																		local.get $$t0.0
																																		call $runtime.Block.Retain
																																		local.get $$t0.1
																																		i32.const 1
																																		i32.const 0
																																		i32.mul
																																		i32.add
																																		local.get $$t23
																																		i32.const 0
																																		i32.sub
																																		i32.const 4
																																		i32.const 0
																																		i32.sub
																																		local.set $$t24.3
																																		local.set $$t24.2
																																		local.set $$t24.1
																																		local.get $$t24.0
																																		call $runtime.Block.Release
																																		local.set $$t24.0
																																		local.get $buf.0
																																		local.get $buf.1
																																		local.get $buf.2
																																		local.get $buf.3
																																		local.get $$t24.0
																																		local.get $$t24.1
																																		local.get $$t24.2
																																		local.get $$t24.3
																																		call $$u8.$slice.append
																																		local.set $$t25.3
																																		local.set $$t25.2
																																		local.set $$t25.1
																																		local.get $$t25.0
																																		call $runtime.Block.Release
																																		local.set $$t25.0
																																		local.get $$t25.0
																																		call $runtime.Block.Retain
																																		local.get $$t25.1
																																		local.get $$t25.2
																																		local.get $$t25.3
																																		local.set $$ret_0.3
																																		local.set $$ret_0.2
																																		local.set $$ret_0.1
																																		local.get $$ret_0.0
																																		call $runtime.Block.Release
																																		local.set $$ret_0.0
																																		br $$BlockFnBody
																																	end
																																	i32.const 10
																																	local.set $$current_block
																																	local.get $graphicOnly
																																	if
																																		br $$Block_10
																																	else
																																		i32.const 5
																																		local.set $$block_selector
																																		br $$BlockDisp
																																	end
																																end
																																i32.const 11
																																local.set $$current_block
																																local.get $r
																																call $strconv.isInGraphicList
																																local.set $$t26
																																local.get $$t26
																																if
																																	i32.const 9
																																	local.set $$block_selector
																																	br $$BlockDisp
																																else
																																	i32.const 5
																																	local.set $$block_selector
																																	br $$BlockDisp
																																end
																															end
																															local.get $$current_block
																															i32.const 13
																															i32.eq
																															if(result i32 i32 i32 i32)
																																local.get $$t27.0
																																call $runtime.Block.Retain
																																local.get $$t27.1
																																local.get $$t27.2
																																local.get $$t27.3
																															else
																																local.get $$current_block
																																i32.const 14
																																i32.eq
																																if(result i32 i32 i32 i32)
																																	local.get $$t28.0
																																	call $runtime.Block.Retain
																																	local.get $$t28.1
																																	local.get $$t28.2
																																	local.get $$t28.3
																																else
																																	local.get $$current_block
																																	i32.const 16
																																	i32.eq
																																	if(result i32 i32 i32 i32)
																																		local.get $$t29.0
																																		call $runtime.Block.Retain
																																		local.get $$t29.1
																																		local.get $$t29.2
																																		local.get $$t29.3
																																	else
																																		local.get $$current_block
																																		i32.const 18
																																		i32.eq
																																		if(result i32 i32 i32 i32)
																																			local.get $$t30.0
																																			call $runtime.Block.Retain
																																			local.get $$t30.1
																																			local.get $$t30.2
																																			local.get $$t30.3
																																		else
																																			local.get $$current_block
																																			i32.const 20
																																			i32.eq
																																			if(result i32 i32 i32 i32)
																																				local.get $$t31.0
																																				call $runtime.Block.Retain
																																				local.get $$t31.1
																																				local.get $$t31.2
																																				local.get $$t31.3
																																			else
																																				local.get $$current_block
																																				i32.const 22
																																				i32.eq
																																				if(result i32 i32 i32 i32)
																																					local.get $$t32.0
																																					call $runtime.Block.Retain
																																					local.get $$t32.1
																																					local.get $$t32.2
																																					local.get $$t32.3
																																				else
																																					local.get $$current_block
																																					i32.const 24
																																					i32.eq
																																					if(result i32 i32 i32 i32)
																																						local.get $$t33.0
																																						call $runtime.Block.Retain
																																						local.get $$t33.1
																																						local.get $$t33.2
																																						local.get $$t33.3
																																					else
																																						local.get $$current_block
																																						i32.const 27
																																						i32.eq
																																						if(result i32 i32 i32 i32)
																																							local.get $$t34.0
																																							call $runtime.Block.Retain
																																							local.get $$t34.1
																																							local.get $$t34.2
																																							local.get $$t34.3
																																						else
																																							local.get $$current_block
																																							i32.const 33
																																							i32.eq
																																							if(result i32 i32 i32 i32)
																																								local.get $$t35.0
																																								call $runtime.Block.Retain
																																								local.get $$t35.1
																																								local.get $$t35.2
																																								local.get $$t35.3
																																							else
																																								local.get $$current_block
																																								i32.const 36
																																								i32.eq
																																								if(result i32 i32 i32 i32)
																																									local.get $$t36.0
																																									call $runtime.Block.Retain
																																									local.get $$t36.1
																																									local.get $$t36.2
																																									local.get $$t36.3
																																								else
																																									local.get $$t37.0
																																									call $runtime.Block.Retain
																																									local.get $$t37.1
																																									local.get $$t37.2
																																									local.get $$t37.3
																																								end
																																							end
																																						end
																																					end
																																				end
																																			end
																																		end
																																	end
																																end
																															end
																															local.set $$t38.3
																															local.set $$t38.2
																															local.set $$t38.1
																															local.get $$t38.0
																															call $runtime.Block.Release
																															local.set $$t38.0
																															i32.const 12
																															local.set $$current_block
																															local.get $$t38.0
																															call $runtime.Block.Retain
																															local.get $$t38.1
																															local.get $$t38.2
																															local.get $$t38.3
																															local.set $$ret_0.3
																															local.set $$ret_0.2
																															local.set $$ret_0.1
																															local.get $$ret_0.0
																															call $runtime.Block.Release
																															local.set $$ret_0.0
																															br $$BlockFnBody
																														end
																														i32.const 13
																														local.set $$current_block
																														local.get $buf.0
																														local.get $buf.1
																														local.get $buf.2
																														local.get $buf.3
																														i32.const 0
																														i32.const 32720
																														i32.const 2
																														i32.const 2
																														call $$u8.$slice.append
																														local.set $$t27.3
																														local.set $$t27.2
																														local.set $$t27.1
																														local.get $$t27.0
																														call $runtime.Block.Release
																														local.set $$t27.0
																														i32.const 12
																														local.set $$block_selector
																														br $$BlockDisp
																													end
																													i32.const 14
																													local.set $$current_block
																													local.get $buf.0
																													local.get $buf.1
																													local.get $buf.2
																													local.get $buf.3
																													i32.const 0
																													i32.const 32722
																													i32.const 2
																													i32.const 2
																													call $$u8.$slice.append
																													local.set $$t28.3
																													local.set $$t28.2
																													local.set $$t28.1
																													local.get $$t28.0
																													call $runtime.Block.Release
																													local.set $$t28.0
																													i32.const 12
																													local.set $$block_selector
																													br $$BlockDisp
																												end
																												i32.const 15
																												local.set $$current_block
																												local.get $r
																												i32.const 8
																												i32.eq
																												local.set $$t39
																												local.get $$t39
																												if
																													i32.const 14
																													local.set $$block_selector
																													br $$BlockDisp
																												else
																													br $$Block_16
																												end
																											end
																											i32.const 16
																											local.set $$current_block
																											local.get $buf.0
																											local.get $buf.1
																											local.get $buf.2
																											local.get $buf.3
																											i32.const 0
																											i32.const 32724
																											i32.const 2
																											i32.const 2
																											call $$u8.$slice.append
																											local.set $$t29.3
																											local.set $$t29.2
																											local.set $$t29.1
																											local.get $$t29.0
																											call $runtime.Block.Release
																											local.set $$t29.0
																											i32.const 12
																											local.set $$block_selector
																											br $$BlockDisp
																										end
																										i32.const 17
																										local.set $$current_block
																										local.get $r
																										i32.const 12
																										i32.eq
																										local.set $$t40
																										local.get $$t40
																										if
																											i32.const 16
																											local.set $$block_selector
																											br $$BlockDisp
																										else
																											br $$Block_18
																										end
																									end
																									i32.const 18
																									local.set $$current_block
																									local.get $buf.0
																									local.get $buf.1
																									local.get $buf.2
																									local.get $buf.3
																									i32.const 0
																									i32.const 32726
																									i32.const 2
																									i32.const 2
																									call $$u8.$slice.append
																									local.set $$t30.3
																									local.set $$t30.2
																									local.set $$t30.1
																									local.get $$t30.0
																									call $runtime.Block.Release
																									local.set $$t30.0
																									i32.const 12
																									local.set $$block_selector
																									br $$BlockDisp
																								end
																								i32.const 19
																								local.set $$current_block
																								local.get $r
																								i32.const 10
																								i32.eq
																								local.set $$t41
																								local.get $$t41
																								if
																									i32.const 18
																									local.set $$block_selector
																									br $$BlockDisp
																								else
																									br $$Block_20
																								end
																							end
																							i32.const 20
																							local.set $$current_block
																							local.get $buf.0
																							local.get $buf.1
																							local.get $buf.2
																							local.get $buf.3
																							i32.const 0
																							i32.const 32728
																							i32.const 2
																							i32.const 2
																							call $$u8.$slice.append
																							local.set $$t31.3
																							local.set $$t31.2
																							local.set $$t31.1
																							local.get $$t31.0
																							call $runtime.Block.Release
																							local.set $$t31.0
																							i32.const 12
																							local.set $$block_selector
																							br $$BlockDisp
																						end
																						i32.const 21
																						local.set $$current_block
																						local.get $r
																						i32.const 13
																						i32.eq
																						local.set $$t42
																						local.get $$t42
																						if
																							i32.const 20
																							local.set $$block_selector
																							br $$BlockDisp
																						else
																							br $$Block_22
																						end
																					end
																					i32.const 22
																					local.set $$current_block
																					local.get $buf.0
																					local.get $buf.1
																					local.get $buf.2
																					local.get $buf.3
																					i32.const 0
																					i32.const 32730
																					i32.const 2
																					i32.const 2
																					call $$u8.$slice.append
																					local.set $$t32.3
																					local.set $$t32.2
																					local.set $$t32.1
																					local.get $$t32.0
																					call $runtime.Block.Release
																					local.set $$t32.0
																					i32.const 12
																					local.set $$block_selector
																					br $$BlockDisp
																				end
																				i32.const 23
																				local.set $$current_block
																				local.get $r
																				i32.const 9
																				i32.eq
																				local.set $$t43
																				local.get $$t43
																				if
																					i32.const 22
																					local.set $$block_selector
																					br $$BlockDisp
																				else
																					br $$Block_24
																				end
																			end
																			i32.const 24
																			local.set $$current_block
																			local.get $buf.0
																			local.get $buf.1
																			local.get $buf.2
																			local.get $buf.3
																			i32.const 0
																			i32.const 32732
																			i32.const 2
																			i32.const 2
																			call $$u8.$slice.append
																			local.set $$t33.3
																			local.set $$t33.2
																			local.set $$t33.1
																			local.get $$t33.0
																			call $runtime.Block.Release
																			local.set $$t33.0
																			i32.const 12
																			local.set $$block_selector
																			br $$BlockDisp
																		end
																		i32.const 25
																		local.set $$current_block
																		local.get $r
																		i32.const 11
																		i32.eq
																		local.set $$t44
																		local.get $$t44
																		if
																			i32.const 24
																			local.set $$block_selector
																			br $$BlockDisp
																		else
																			br $$Block_25
																		end
																	end
																	i32.const 26
																	local.set $$current_block
																	local.get $r
																	i32.const 32
																	i32.lt_s
																	local.set $$t45
																	local.get $$t45
																	if
																		br $$Block_26
																	else
																		br $$Block_28
																	end
																end
																i32.const 27
																local.set $$current_block
																local.get $buf.0
																local.get $buf.1
																local.get $buf.2
																local.get $buf.3
																i32.const 0
																i32.const 32734
																i32.const 2
																i32.const 2
																call $$u8.$slice.append
																local.set $$t46.3
																local.set $$t46.2
																local.set $$t46.1
																local.get $$t46.0
																call $runtime.Block.Release
																local.set $$t46.0
																local.get $r
																i32.const 255
																i32.and
																local.set $$t47
																local.get $$t47
																i64.const 4
																i32.wrap_i64
																i32.shr_u
																local.set $$t48
																i32.const 31528
																local.get $$t48
																i32.add
																i32.load8_u align=1
																local.set $$t49
																i32.const 17
																call $runtime.HeapAlloc
																i32.const 1
																i32.const 0
																i32.const 1
																call $runtime.Block.Init
																call $runtime.DupI32
																i32.const 16
																i32.add
																local.set $$t50.1
																local.get $$t50.0
																call $runtime.Block.Release
																local.set $$t50.0
																local.get $$t50.0
																call $runtime.Block.Retain
																local.get $$t50.1
																i32.const 1
																i32.const 0
																i32.mul
																i32.add
																local.set $$t51.1
																local.get $$t51.0
																call $runtime.Block.Release
																local.set $$t51.0
																local.get $$t51.1
																local.get $$t49
																i32.store8 align=1
																local.get $$t50.0
																call $runtime.Block.Retain
																local.get $$t50.1
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
																local.set $$t52.3
																local.set $$t52.2
																local.set $$t52.1
																local.get $$t52.0
																call $runtime.Block.Release
																local.set $$t52.0
																local.get $$t46.0
																local.get $$t46.1
																local.get $$t46.2
																local.get $$t46.3
																local.get $$t52.0
																local.get $$t52.1
																local.get $$t52.2
																local.get $$t52.3
																call $$u8.$slice.append
																local.set $$t53.3
																local.set $$t53.2
																local.set $$t53.1
																local.get $$t53.0
																call $runtime.Block.Release
																local.set $$t53.0
																local.get $r
																i32.const 255
																i32.and
																local.set $$t54
																local.get $$t54
																i32.const 15
																i32.and
																local.set $$t55
																i32.const 31528
																local.get $$t55
																i32.add
																i32.load8_u align=1
																local.set $$t56
																i32.const 17
																call $runtime.HeapAlloc
																i32.const 1
																i32.const 0
																i32.const 1
																call $runtime.Block.Init
																call $runtime.DupI32
																i32.const 16
																i32.add
																local.set $$t57.1
																local.get $$t57.0
																call $runtime.Block.Release
																local.set $$t57.0
																local.get $$t57.0
																call $runtime.Block.Retain
																local.get $$t57.1
																i32.const 1
																i32.const 0
																i32.mul
																i32.add
																local.set $$t58.1
																local.get $$t58.0
																call $runtime.Block.Release
																local.set $$t58.0
																local.get $$t58.1
																local.get $$t56
																i32.store8 align=1
																local.get $$t57.0
																call $runtime.Block.Retain
																local.get $$t57.1
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
																local.set $$t59.3
																local.set $$t59.2
																local.set $$t59.1
																local.get $$t59.0
																call $runtime.Block.Release
																local.set $$t59.0
																local.get $$t53.0
																local.get $$t53.1
																local.get $$t53.2
																local.get $$t53.3
																local.get $$t59.0
																local.get $$t59.1
																local.get $$t59.2
																local.get $$t59.3
																call $$u8.$slice.append
																local.set $$t34.3
																local.set $$t34.2
																local.set $$t34.1
																local.get $$t34.0
																call $runtime.Block.Release
																local.set $$t34.0
																i32.const 12
																local.set $$block_selector
																br $$BlockDisp
															end
															i32.const 28
															local.set $$current_block
															local.get $buf.0
															local.get $buf.1
															local.get $buf.2
															local.get $buf.3
															i32.const 0
															i32.const 32736
															i32.const 2
															i32.const 2
															call $$u8.$slice.append
															local.set $$t60.3
															local.set $$t60.2
															local.set $$t60.1
															local.get $$t60.0
															call $runtime.Block.Release
															local.set $$t60.0
															br $$Block_32
														end
														i32.const 29
														local.set $$current_block
														local.get $r
														i32.const 1114111
														i32.gt_s
														local.set $$t61
														local.get $$t61
														if
															i32.const 28
															local.set $$block_selector
															br $$BlockDisp
														else
															br $$Block_30
														end
													end
													i32.const 30
													local.set $$current_block
													local.get $buf.0
													local.get $buf.1
													local.get $buf.2
													local.get $buf.3
													i32.const 0
													i32.const 32736
													i32.const 2
													i32.const 2
													call $$u8.$slice.append
													local.set $$t62.3
													local.set $$t62.2
													local.set $$t62.1
													local.get $$t62.0
													call $runtime.Block.Release
													local.set $$t62.0
													br $$Block_35
												end
												i32.const 31
												local.set $$current_block
												local.get $r
												i32.const 65536
												i32.lt_s
												local.set $$t63
												local.get $$t63
												if
													i32.const 30
													local.set $$block_selector
													br $$BlockDisp
												else
													br $$Block_33
												end
											end
											i32.const 32
											local.set $$current_block
											local.get $$t64
											local.set $$t65
											i32.const 65533
											local.get $$t65
											i32.shr_s
											local.set $$t66
											local.get $$t66
											i32.const 15
											i32.and
											local.set $$t67
											i32.const 31528
											local.get $$t67
											i32.add
											i32.load8_u align=1
											local.set $$t68
											i32.const 17
											call $runtime.HeapAlloc
											i32.const 1
											i32.const 0
											i32.const 1
											call $runtime.Block.Init
											call $runtime.DupI32
											i32.const 16
											i32.add
											local.set $$t69.1
											local.get $$t69.0
											call $runtime.Block.Release
											local.set $$t69.0
											local.get $$t69.0
											call $runtime.Block.Retain
											local.get $$t69.1
											i32.const 1
											i32.const 0
											i32.mul
											i32.add
											local.set $$t70.1
											local.get $$t70.0
											call $runtime.Block.Release
											local.set $$t70.0
											local.get $$t70.1
											local.get $$t68
											i32.store8 align=1
											local.get $$t69.0
											call $runtime.Block.Retain
											local.get $$t69.1
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
											local.set $$t71.3
											local.set $$t71.2
											local.set $$t71.1
											local.get $$t71.0
											call $runtime.Block.Release
											local.set $$t71.0
											local.get $$t35.0
											local.get $$t35.1
											local.get $$t35.2
											local.get $$t35.3
											local.get $$t71.0
											local.get $$t71.1
											local.get $$t71.2
											local.get $$t71.3
											call $$u8.$slice.append
											local.set $$t72.3
											local.set $$t72.2
											local.set $$t72.1
											local.get $$t72.0
											call $runtime.Block.Release
											local.set $$t72.0
											local.get $$t64
											i32.const 4
											i32.sub
											local.set $$t73
											br $$Block_32
										end
										local.get $$current_block
										i32.const 28
										i32.eq
										if(result i32 i32 i32 i32)
											local.get $$t60.0
											call $runtime.Block.Retain
											local.get $$t60.1
											local.get $$t60.2
											local.get $$t60.3
										else
											local.get $$t72.0
											call $runtime.Block.Retain
											local.get $$t72.1
											local.get $$t72.2
											local.get $$t72.3
										end
										local.get $$current_block
										i32.const 28
										i32.eq
										if(result i32)
											i32.const 12
										else
											local.get $$t73
										end
										local.set $$t64
										local.set $$t35.3
										local.set $$t35.2
										local.set $$t35.1
										local.get $$t35.0
										call $runtime.Block.Release
										local.set $$t35.0
										i32.const 33
										local.set $$current_block
										local.get $$t64
										i32.const 0
										i32.ge_s
										local.set $$t74
										local.get $$t74
										if
											i32.const 32
											local.set $$block_selector
											br $$BlockDisp
										else
											i32.const 12
											local.set $$block_selector
											br $$BlockDisp
										end
									end
									i32.const 34
									local.set $$current_block
									local.get $buf.0
									local.get $buf.1
									local.get $buf.2
									local.get $buf.3
									i32.const 0
									i32.const 32738
									i32.const 2
									i32.const 2
									call $$u8.$slice.append
									local.set $$t75.3
									local.set $$t75.2
									local.set $$t75.1
									local.get $$t75.0
									call $runtime.Block.Release
									local.set $$t75.0
									br $$Block_37
								end
								i32.const 35
								local.set $$current_block
								local.get $$t76
								local.set $$t77
								local.get $r
								local.get $$t77
								i32.shr_s
								local.set $$t78
								local.get $$t78
								i32.const 15
								i32.and
								local.set $$t79
								i32.const 31528
								local.get $$t79
								i32.add
								i32.load8_u align=1
								local.set $$t80
								i32.const 17
								call $runtime.HeapAlloc
								i32.const 1
								i32.const 0
								i32.const 1
								call $runtime.Block.Init
								call $runtime.DupI32
								i32.const 16
								i32.add
								local.set $$t81.1
								local.get $$t81.0
								call $runtime.Block.Release
								local.set $$t81.0
								local.get $$t81.0
								call $runtime.Block.Retain
								local.get $$t81.1
								i32.const 1
								i32.const 0
								i32.mul
								i32.add
								local.set $$t82.1
								local.get $$t82.0
								call $runtime.Block.Release
								local.set $$t82.0
								local.get $$t82.1
								local.get $$t80
								i32.store8 align=1
								local.get $$t81.0
								call $runtime.Block.Retain
								local.get $$t81.1
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
								local.set $$t83.3
								local.set $$t83.2
								local.set $$t83.1
								local.get $$t83.0
								call $runtime.Block.Release
								local.set $$t83.0
								local.get $$t36.0
								local.get $$t36.1
								local.get $$t36.2
								local.get $$t36.3
								local.get $$t83.0
								local.get $$t83.1
								local.get $$t83.2
								local.get $$t83.3
								call $$u8.$slice.append
								local.set $$t84.3
								local.set $$t84.2
								local.set $$t84.1
								local.get $$t84.0
								call $runtime.Block.Release
								local.set $$t84.0
								local.get $$t76
								i32.const 4
								i32.sub
								local.set $$t85
								br $$Block_35
							end
							local.get $$current_block
							i32.const 30
							i32.eq
							if(result i32 i32 i32 i32)
								local.get $$t62.0
								call $runtime.Block.Retain
								local.get $$t62.1
								local.get $$t62.2
								local.get $$t62.3
							else
								local.get $$t84.0
								call $runtime.Block.Retain
								local.get $$t84.1
								local.get $$t84.2
								local.get $$t84.3
							end
							local.get $$current_block
							i32.const 30
							i32.eq
							if(result i32)
								i32.const 12
							else
								local.get $$t85
							end
							local.set $$t76
							local.set $$t36.3
							local.set $$t36.2
							local.set $$t36.1
							local.get $$t36.0
							call $runtime.Block.Release
							local.set $$t36.0
							i32.const 36
							local.set $$current_block
							local.get $$t76
							i32.const 0
							i32.ge_s
							local.set $$t86
							local.get $$t86
							if
								i32.const 35
								local.set $$block_selector
								br $$BlockDisp
							else
								i32.const 12
								local.set $$block_selector
								br $$BlockDisp
							end
						end
						i32.const 37
						local.set $$current_block
						local.get $$t87
						local.set $$t88
						local.get $r
						local.get $$t88
						i32.shr_s
						local.set $$t89
						local.get $$t89
						i32.const 15
						i32.and
						local.set $$t90
						i32.const 31528
						local.get $$t90
						i32.add
						i32.load8_u align=1
						local.set $$t91
						i32.const 17
						call $runtime.HeapAlloc
						i32.const 1
						i32.const 0
						i32.const 1
						call $runtime.Block.Init
						call $runtime.DupI32
						i32.const 16
						i32.add
						local.set $$t92.1
						local.get $$t92.0
						call $runtime.Block.Release
						local.set $$t92.0
						local.get $$t92.0
						call $runtime.Block.Retain
						local.get $$t92.1
						i32.const 1
						i32.const 0
						i32.mul
						i32.add
						local.set $$t93.1
						local.get $$t93.0
						call $runtime.Block.Release
						local.set $$t93.0
						local.get $$t93.1
						local.get $$t91
						i32.store8 align=1
						local.get $$t92.0
						call $runtime.Block.Retain
						local.get $$t92.1
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
						local.set $$t94.3
						local.set $$t94.2
						local.set $$t94.1
						local.get $$t94.0
						call $runtime.Block.Release
						local.set $$t94.0
						local.get $$t37.0
						local.get $$t37.1
						local.get $$t37.2
						local.get $$t37.3
						local.get $$t94.0
						local.get $$t94.1
						local.get $$t94.2
						local.get $$t94.3
						call $$u8.$slice.append
						local.set $$t95.3
						local.set $$t95.2
						local.set $$t95.1
						local.get $$t95.0
						call $runtime.Block.Release
						local.set $$t95.0
						local.get $$t87
						i32.const 4
						i32.sub
						local.set $$t96
						br $$Block_37
					end
					local.get $$current_block
					i32.const 34
					i32.eq
					if(result i32 i32 i32 i32)
						local.get $$t75.0
						call $runtime.Block.Retain
						local.get $$t75.1
						local.get $$t75.2
						local.get $$t75.3
					else
						local.get $$t95.0
						call $runtime.Block.Retain
						local.get $$t95.1
						local.get $$t95.2
						local.get $$t95.3
					end
					local.get $$current_block
					i32.const 34
					i32.eq
					if(result i32)
						i32.const 28
					else
						local.get $$t96
					end
					local.set $$t87
					local.set $$t37.3
					local.set $$t37.2
					local.set $$t37.1
					local.get $$t37.0
					call $runtime.Block.Release
					local.set $$t37.0
					i32.const 38
					local.set $$current_block
					local.get $$t87
					i32.const 0
					i32.ge_s
					local.set $$t97
					local.get $$t97
					if
						i32.const 37
						local.set $$block_selector
						br $$BlockDisp
					else
						i32.const 12
						local.set $$block_selector
						br $$BlockDisp
					end
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
		local.get $$t3.0
		call $runtime.Block.Release
		local.get $$t4.0
		call $runtime.Block.Release
		local.get $$t5.0
		call $runtime.Block.Release
		local.get $$t6.0
		call $runtime.Block.Release
		local.get $$t8.0
		call $runtime.Block.Release
		local.get $$t9.0
		call $runtime.Block.Release
		local.get $$t10.0
		call $runtime.Block.Release
		local.get $$t11.0
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
		local.get $$t24.0
		call $runtime.Block.Release
		local.get $$t25.0
		call $runtime.Block.Release
		local.get $$t27.0
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
		local.get $$t46.0
		call $runtime.Block.Release
		local.get $$t50.0
		call $runtime.Block.Release
		local.get $$t51.0
		call $runtime.Block.Release
		local.get $$t52.0
		call $runtime.Block.Release
		local.get $$t53.0
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
		local.get $$t69.0
		call $runtime.Block.Release
		local.get $$t70.0
		call $runtime.Block.Release
		local.get $$t71.0
		call $runtime.Block.Release
		local.get $$t72.0
		call $runtime.Block.Release
		local.get $$t75.0
		call $runtime.Block.Release
		local.get $$t81.0
		call $runtime.Block.Release
		local.get $$t82.0
		call $runtime.Block.Release
		local.get $$t83.0
		call $runtime.Block.Release
		local.get $$t84.0
		call $runtime.Block.Release
		local.get $$t92.0
		call $runtime.Block.Release
		local.get $$t93.0
		call $runtime.Block.Release
		local.get $$t94.0
		call $runtime.Block.Release
		local.get $$t95.0
		call $runtime.Block.Release
	)
	(func $$u8.$slice.copy (param $d.0 i32) (param $d.1 i32) (param $d.2 i32) (param $d.3 i32) (param $s.0 i32) (param $s.1 i32) (param $s.2 i32) (param $s.3 i32) (result i32)
		(local $item i32)
		(local $count i32)
		(local $dp i32)
		(local $sp i32)
		(local $item_size i32)
		local.get $d.2
		local.get $s.2
		i32.gt_u
		if
			local.get $s.2
			local.set $count
		else
			local.get $d.2
			local.set $count
		end
		local.get $count
		local.get $d.1
		local.get $s.1
		i32.lt_u
		if
			local.get $d.1
			local.set $dp
			local.get $s.1
			local.set $sp
			i32.const 1
			local.set $item_size
		else
			local.get $count
			i32.const 1
			i32.sub
			i32.const 1
			i32.mul
			local.set $item_size
			local.get $d.1
			local.get $item_size
			i32.add
			local.set $dp
			local.get $s.1
			local.get $item_size
			i32.add
			local.set $sp
			i32.const 0
			i32.const 1
			i32.sub
			local.set $item_size
		end
		block $b0
			loop $l0
				local.get $count
				i32.eqz
				if
					br $b0
				else
					local.get $sp
					i32.load8_u align=1
					local.set $item
					local.get $dp
					local.get $item
					i32.store8 align=1
					local.get $sp
					local.get $item_size
					i32.add
					local.set $sp
					local.get $dp
					local.get $item_size
					i32.add
					local.set $dp
					local.get $count
					i32.const 1
					i32.sub
					local.set $count
					br $l0
				end
			end
		end
	)
	(func $strconv.appendQuotedWith (param $buf.0 i32) (param $buf.1 i32) (param $buf.2 i32) (param $buf.3 i32) (param $s.0 i32) (param $s.1 i32) (param $s.2 i32) (param $quote i32) (param $ASCIIonly i32) (param $graphicOnly i32) (result i32 i32 i32 i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$ret_0.0 i32)
		(local $$ret_0.1 i32)
		(local $$ret_0.2 i32)
		(local $$ret_0.3 i32)
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
		(local $$t11.0 i32)
		(local $$t11.1 i32)
		(local $$t11.2 i32)
		(local $$t11.3 i32)
		(local $$t12 i32)
		(local $$t13.0 i32)
		(local $$t13.1 i32)
		(local $$t13.2 i32)
		(local $$t13.3 i32)
		(local $$t14.0 i32)
		(local $$t14.1 i32)
		(local $$t15.0 i32)
		(local $$t15.1 i32)
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
		(local $$t19 i32)
		(local $$t20 i32)
		(local $$t21 i32)
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
		(local $$t26.2 i32)
		(local $$t26.3 i32)
		(local $$t27.0 i32)
		(local $$t27.1 i32)
		(local $$t27.2 i32)
		(local $$t27.3 i32)
		(local $$t28.0 i32)
		(local $$t28.1 i32)
		(local $$t28.2 i32)
		(local $$t29 i32)
		(local $$t30 i32)
		(local $$t31.0 i32)
		(local $$t31.1 i32)
		(local $$t31.2 i32)
		(local $$t31.3 i32)
		(local $$t32.0 i32)
		(local $$t32.1 i32)
		(local $$t32.2 i32)
		(local $$t32.3 i32)
		(local $$t33 i32)
		(local $$t34.0 i32)
		(local $$t34.1 i32)
		(local $$t35 i32)
		(local $$t36 i32)
		(local $$t37 i32)
		(local $$t38 i32)
		(local $$t39.0 i32)
		(local $$t39.1 i32)
		(local $$t39.2 i32)
		(local $$t39.3 i32)
		(local $$t40 i32)
		(local $$t41 i32)
		(local $$t42 i32)
		(local $$t43.0 i32)
		(local $$t43.1 i32)
		(local $$t44.0 i32)
		(local $$t44.1 i32)
		(local $$t45.0 i32)
		(local $$t45.1 i32)
		(local $$t45.2 i32)
		(local $$t45.3 i32)
		(local $$t46.0 i32)
		(local $$t46.1 i32)
		(local $$t46.2 i32)
		(local $$t46.3 i32)
		(local $$t47 i32)
		(local $$t48 i32)
		(local $$t49 i32)
		(local $$t50.0 i32)
		(local $$t50.1 i32)
		(local $$t51.0 i32)
		(local $$t51.1 i32)
		(local $$t52.0 i32)
		(local $$t52.1 i32)
		(local $$t52.2 i32)
		(local $$t52.3 i32)
		(local $$t53 i32)
		block $$BlockFnBody
			loop $$BlockDisp
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
																	br_table  0 1 2 3 4 5 6 7 8 9 10 11 0
																end
																i32.const 0
																local.set $$current_block
																local.get $buf.3
																local.set $$t0
																local.get $buf.2
																local.set $$t1
																local.get $$t0
																local.get $$t1
																i32.sub
																local.set $$t2
																local.get $s.2
																local.set $$t3
																local.get $$t2
																local.get $$t3
																i32.lt_s
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
															local.get $buf.2
															local.set $$t5
															local.get $buf.2
															local.set $$t6
															local.get $$t6
															i32.const 1
															i32.add
															local.set $$t7
															local.get $s.2
															local.set $$t8
															local.get $$t7
															local.get $$t8
															i32.add
															local.set $$t9
															local.get $$t9
															i32.const 1
															i32.add
															local.set $$t10
															local.get $$t10
															i32.const 1
															i32.mul
															i32.const 16
															i32.add
															call $runtime.HeapAlloc
															local.get $$t10
															i32.const 0
															i32.const 1
															call $runtime.Block.Init
															call $runtime.DupI32
															i32.const 16
															i32.add
															local.get $$t5
															local.get $$t10
															local.set $$t11.3
															local.set $$t11.2
															local.set $$t11.1
															local.get $$t11.0
															call $runtime.Block.Release
															local.set $$t11.0
															local.get $$t11.0
															local.get $$t11.1
															local.get $$t11.2
															local.get $$t11.3
															local.get $buf.0
															local.get $buf.1
															local.get $buf.2
															local.get $buf.3
															call $$u8.$slice.copy
															local.set $$t12
															br $$Block_1
														end
														local.get $$current_block
														i32.const 0
														i32.eq
														if(result i32 i32 i32 i32)
															local.get $buf.0
															call $runtime.Block.Retain
															local.get $buf.1
															local.get $buf.2
															local.get $buf.3
														else
															local.get $$t11.0
															call $runtime.Block.Retain
															local.get $$t11.1
															local.get $$t11.2
															local.get $$t11.3
														end
														local.set $$t13.3
														local.set $$t13.2
														local.set $$t13.1
														local.get $$t13.0
														call $runtime.Block.Release
														local.set $$t13.0
														i32.const 2
														local.set $$current_block
														i32.const 17
														call $runtime.HeapAlloc
														i32.const 1
														i32.const 0
														i32.const 1
														call $runtime.Block.Init
														call $runtime.DupI32
														i32.const 16
														i32.add
														local.set $$t14.1
														local.get $$t14.0
														call $runtime.Block.Release
														local.set $$t14.0
														local.get $$t14.0
														call $runtime.Block.Retain
														local.get $$t14.1
														i32.const 1
														i32.const 0
														i32.mul
														i32.add
														local.set $$t15.1
														local.get $$t15.0
														call $runtime.Block.Release
														local.set $$t15.0
														local.get $$t15.1
														local.get $quote
														i32.store8 align=1
														local.get $$t14.0
														call $runtime.Block.Retain
														local.get $$t14.1
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
														local.set $$t16.3
														local.set $$t16.2
														local.set $$t16.1
														local.get $$t16.0
														call $runtime.Block.Release
														local.set $$t16.0
														local.get $$t13.0
														local.get $$t13.1
														local.get $$t13.2
														local.get $$t13.3
														local.get $$t16.0
														local.get $$t16.1
														local.get $$t16.2
														local.get $$t16.3
														call $$u8.$slice.append
														local.set $$t17.3
														local.set $$t17.2
														local.set $$t17.1
														local.get $$t17.0
														call $runtime.Block.Release
														local.set $$t17.0
														br $$Block_4
													end
													i32.const 3
													local.set $$current_block
													local.get $$t18.1
													i32.const 0
													i32.add
													i32.load8_u align=1
													local.set $$t19
													local.get $$t19
													local.set $$t20
													local.get $$t20
													i32.const 128
													i32.ge_s
													local.set $$t21
													local.get $$t21
													if
														br $$Block_6
													else
														br $$Block_7
													end
												end
												i32.const 4
												local.set $$current_block
												i32.const 17
												call $runtime.HeapAlloc
												i32.const 1
												i32.const 0
												i32.const 1
												call $runtime.Block.Init
												call $runtime.DupI32
												i32.const 16
												i32.add
												local.set $$t22.1
												local.get $$t22.0
												call $runtime.Block.Release
												local.set $$t22.0
												local.get $$t22.0
												call $runtime.Block.Retain
												local.get $$t22.1
												i32.const 1
												i32.const 0
												i32.mul
												i32.add
												local.set $$t23.1
												local.get $$t23.0
												call $runtime.Block.Release
												local.set $$t23.0
												local.get $$t23.1
												local.get $quote
												i32.store8 align=1
												local.get $$t22.0
												call $runtime.Block.Retain
												local.get $$t22.1
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
												local.set $$t24.3
												local.set $$t24.2
												local.set $$t24.1
												local.get $$t24.0
												call $runtime.Block.Release
												local.set $$t24.0
												local.get $$t25.0
												local.get $$t25.1
												local.get $$t25.2
												local.get $$t25.3
												local.get $$t24.0
												local.get $$t24.1
												local.get $$t24.2
												local.get $$t24.3
												call $$u8.$slice.append
												local.set $$t26.3
												local.set $$t26.2
												local.set $$t26.1
												local.get $$t26.0
												call $runtime.Block.Release
												local.set $$t26.0
												local.get $$t26.0
												call $runtime.Block.Retain
												local.get $$t26.1
												local.get $$t26.2
												local.get $$t26.3
												local.set $$ret_0.3
												local.set $$ret_0.2
												local.set $$ret_0.1
												local.get $$ret_0.0
												call $runtime.Block.Release
												local.set $$ret_0.0
												br $$BlockFnBody
											end
											local.get $$current_block
											i32.const 2
											i32.eq
											if(result i32 i32 i32 i32)
												local.get $$t17.0
												call $runtime.Block.Retain
												local.get $$t17.1
												local.get $$t17.2
												local.get $$t17.3
											else
												local.get $$t27.0
												call $runtime.Block.Retain
												local.get $$t27.1
												local.get $$t27.2
												local.get $$t27.3
											end
											local.get $$current_block
											i32.const 2
											i32.eq
											if(result i32 i32 i32)
												local.get $s.0
												call $runtime.Block.Retain
												local.get $s.1
												local.get $s.2
											else
												local.get $$t28.0
												call $runtime.Block.Retain
												local.get $$t28.1
												local.get $$t28.2
											end
											local.set $$t18.2
											local.set $$t18.1
											local.get $$t18.0
											call $runtime.Block.Release
											local.set $$t18.0
											local.set $$t25.3
											local.set $$t25.2
											local.set $$t25.1
											local.get $$t25.0
											call $runtime.Block.Release
											local.set $$t25.0
											i32.const 5
											local.set $$current_block
											local.get $$t18.2
											local.set $$t29
											local.get $$t29
											i32.const 0
											i32.gt_s
											local.set $$t30
											local.get $$t30
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
										local.get $$current_block
										i32.const 9
										i32.eq
										if(result i32 i32 i32 i32)
											local.get $$t31.0
											call $runtime.Block.Retain
											local.get $$t31.1
											local.get $$t31.2
											local.get $$t31.3
										else
											local.get $$t32.0
											call $runtime.Block.Retain
											local.get $$t32.1
											local.get $$t32.2
											local.get $$t32.3
										end
										local.set $$t27.3
										local.set $$t27.2
										local.set $$t27.1
										local.get $$t27.0
										call $runtime.Block.Release
										local.set $$t27.0
										i32.const 6
										local.set $$current_block
										local.get $$t18.0
										call $runtime.Block.Retain
										local.get $$t18.1
										local.get $$t33
										i32.add
										local.get $$t18.2
										local.get $$t33
										i32.sub
										local.set $$t28.2
										local.set $$t28.1
										local.get $$t28.0
										call $runtime.Block.Release
										local.set $$t28.0
										i32.const 5
										local.set $$block_selector
										br $$BlockDisp
									end
									i32.const 7
									local.set $$current_block
									local.get $$t18.0
									local.get $$t18.1
									local.get $$t18.2
									call $unicode$utf8.DecodeRuneInString
									local.set $$t34.1
									local.set $$t34.0
									local.get $$t34.0
									local.set $$t35
									local.get $$t34.1
									local.set $$t36
									br $$Block_7
								end
								local.get $$current_block
								i32.const 3
								i32.eq
								if(result i32)
									i32.const 1
								else
									local.get $$t36
								end
								local.get $$current_block
								i32.const 3
								i32.eq
								if(result i32)
									local.get $$t20
								else
									local.get $$t35
								end
								local.set $$t37
								local.set $$t33
								i32.const 8
								local.set $$current_block
								local.get $$t33
								i32.const 1
								i32.eq
								local.set $$t38
								local.get $$t38
								if
									br $$Block_10
								else
									br $$Block_9
								end
							end
							i32.const 9
							local.set $$current_block
							local.get $$t25.0
							local.get $$t25.1
							local.get $$t25.2
							local.get $$t25.3
							i32.const 0
							i32.const 32734
							i32.const 2
							i32.const 2
							call $$u8.$slice.append
							local.set $$t39.3
							local.set $$t39.2
							local.set $$t39.1
							local.get $$t39.0
							call $runtime.Block.Release
							local.set $$t39.0
							local.get $$t18.1
							i32.const 0
							i32.add
							i32.load8_u align=1
							local.set $$t40
							local.get $$t40
							i64.const 4
							i32.wrap_i64
							i32.shr_u
							local.set $$t41
							i32.const 31528
							local.get $$t41
							i32.add
							i32.load8_u align=1
							local.set $$t42
							i32.const 17
							call $runtime.HeapAlloc
							i32.const 1
							i32.const 0
							i32.const 1
							call $runtime.Block.Init
							call $runtime.DupI32
							i32.const 16
							i32.add
							local.set $$t43.1
							local.get $$t43.0
							call $runtime.Block.Release
							local.set $$t43.0
							local.get $$t43.0
							call $runtime.Block.Retain
							local.get $$t43.1
							i32.const 1
							i32.const 0
							i32.mul
							i32.add
							local.set $$t44.1
							local.get $$t44.0
							call $runtime.Block.Release
							local.set $$t44.0
							local.get $$t44.1
							local.get $$t42
							i32.store8 align=1
							local.get $$t43.0
							call $runtime.Block.Retain
							local.get $$t43.1
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
							local.set $$t45.3
							local.set $$t45.2
							local.set $$t45.1
							local.get $$t45.0
							call $runtime.Block.Release
							local.set $$t45.0
							local.get $$t39.0
							local.get $$t39.1
							local.get $$t39.2
							local.get $$t39.3
							local.get $$t45.0
							local.get $$t45.1
							local.get $$t45.2
							local.get $$t45.3
							call $$u8.$slice.append
							local.set $$t46.3
							local.set $$t46.2
							local.set $$t46.1
							local.get $$t46.0
							call $runtime.Block.Release
							local.set $$t46.0
							local.get $$t18.1
							i32.const 0
							i32.add
							i32.load8_u align=1
							local.set $$t47
							local.get $$t47
							i32.const 15
							i32.and
							local.set $$t48
							i32.const 31528
							local.get $$t48
							i32.add
							i32.load8_u align=1
							local.set $$t49
							i32.const 17
							call $runtime.HeapAlloc
							i32.const 1
							i32.const 0
							i32.const 1
							call $runtime.Block.Init
							call $runtime.DupI32
							i32.const 16
							i32.add
							local.set $$t50.1
							local.get $$t50.0
							call $runtime.Block.Release
							local.set $$t50.0
							local.get $$t50.0
							call $runtime.Block.Retain
							local.get $$t50.1
							i32.const 1
							i32.const 0
							i32.mul
							i32.add
							local.set $$t51.1
							local.get $$t51.0
							call $runtime.Block.Release
							local.set $$t51.0
							local.get $$t51.1
							local.get $$t49
							i32.store8 align=1
							local.get $$t50.0
							call $runtime.Block.Retain
							local.get $$t50.1
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
							local.set $$t52.3
							local.set $$t52.2
							local.set $$t52.1
							local.get $$t52.0
							call $runtime.Block.Release
							local.set $$t52.0
							local.get $$t46.0
							local.get $$t46.1
							local.get $$t46.2
							local.get $$t46.3
							local.get $$t52.0
							local.get $$t52.1
							local.get $$t52.2
							local.get $$t52.3
							call $$u8.$slice.append
							local.set $$t31.3
							local.set $$t31.2
							local.set $$t31.1
							local.get $$t31.0
							call $runtime.Block.Release
							local.set $$t31.0
							i32.const 6
							local.set $$block_selector
							br $$BlockDisp
						end
						i32.const 10
						local.set $$current_block
						local.get $$t25.0
						local.get $$t25.1
						local.get $$t25.2
						local.get $$t25.3
						local.get $$t37
						local.get $quote
						local.get $ASCIIonly
						local.get $graphicOnly
						call $strconv.appendEscapedRune
						local.set $$t32.3
						local.set $$t32.2
						local.set $$t32.1
						local.get $$t32.0
						call $runtime.Block.Release
						local.set $$t32.0
						i32.const 6
						local.set $$block_selector
						br $$BlockDisp
					end
					i32.const 11
					local.set $$current_block
					local.get $$t37
					i32.const 65533
					i32.eq
					local.set $$t53
					local.get $$t53
					if
						i32.const 9
						local.set $$block_selector
						br $$BlockDisp
					else
						i32.const 10
						local.set $$block_selector
						br $$BlockDisp
					end
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
		local.get $$t11.0
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
		local.get $$t28.0
		call $runtime.Block.Release
		local.get $$t31.0
		call $runtime.Block.Release
		local.get $$t32.0
		call $runtime.Block.Release
		local.get $$t39.0
		call $runtime.Block.Release
		local.get $$t43.0
		call $runtime.Block.Release
		local.get $$t44.0
		call $runtime.Block.Release
		local.get $$t45.0
		call $runtime.Block.Release
		local.get $$t46.0
		call $runtime.Block.Release
		local.get $$t50.0
		call $runtime.Block.Release
		local.get $$t51.0
		call $runtime.Block.Release
		local.get $$t52.0
		call $runtime.Block.Release
	)
	(func $$u8.$slice.underlying.$$onFree (param $$ptr i32)
		local.get $$ptr
		i32.const 2
		call_indirect 0 (type $$onFree)
	)
	(func $$strconv.decimalSlice.$$onFree (param $$ptr i32)
		local.get $$ptr
		i32.const 26
		call_indirect 0 (type $$onFree)
	)
	(func $strconv.bsearch16 (param $a.0 i32) (param $a.1 i32) (param $a.2 i32) (param $a.3 i32) (param $x i32) (result i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$ret_0 i32)
		(local $$t0 i32)
		(local $$t1 i32)
		(local $$t2 i32)
		(local $$t3 i32)
		(local $$t4 i32)
		(local $$t5 i32)
		(local $$t6.0 i32)
		(local $$t6.1 i32)
		(local $$t7 i32)
		(local $$t8 i32)
		(local $$t9 i32)
		(local $$t10 i32)
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
										local.get $a.2
										local.set $$t0
										br $$Block_2
									end
									i32.const 1
									local.set $$current_block
									local.get $$t1
									local.get $$t2
									i32.sub
									local.set $$t3
									local.get $$t3
									i64.const 1
									i32.wrap_i64
									i32.shr_s
									local.set $$t4
									local.get $$t2
									local.get $$t4
									i32.add
									local.set $$t5
									local.get $a.0
									call $runtime.Block.Retain
									local.get $a.1
									i32.const 2
									local.get $$t5
									i32.mul
									i32.add
									local.set $$t6.1
									local.get $$t6.0
									call $runtime.Block.Release
									local.set $$t6.0
									local.get $$t6.1
									i32.load16_u
									local.set $$t7
									local.get $$t7
									local.get $x
									i32.lt_u
									local.set $$t8
									local.get $$t8
									if
										br $$Block_3
									else
										br $$Block_4
									end
								end
								i32.const 2
								local.set $$current_block
								local.get $$t2
								local.set $$ret_0
								br $$BlockFnBody
							end
							local.get $$current_block
							i32.const 0
							i32.eq
							if(result i32)
								i32.const 0
							else
								local.get $$current_block
								i32.const 4
								i32.eq
								if(result i32)
									local.get $$t9
								else
									local.get $$t2
								end
							end
							local.get $$current_block
							i32.const 0
							i32.eq
							if(result i32)
								local.get $$t0
							else
								local.get $$current_block
								i32.const 4
								i32.eq
								if(result i32)
									local.get $$t1
								else
									local.get $$t5
								end
							end
							local.set $$t1
							local.set $$t2
							i32.const 3
							local.set $$current_block
							local.get $$t2
							local.get $$t1
							i32.lt_s
							local.set $$t10
							local.get $$t10
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
						local.get $$t5
						i32.const 1
						i32.add
						local.set $$t9
						i32.const 3
						local.set $$block_selector
						br $$BlockDisp
					end
					i32.const 5
					local.set $$current_block
					i32.const 3
					local.set $$block_selector
					br $$BlockDisp
				end
			end
		end
		local.get $$ret_0
		local.get $$t6.0
		call $runtime.Block.Release
	)
	(func $strconv.bsearch32 (param $a.0 i32) (param $a.1 i32) (param $a.2 i32) (param $a.3 i32) (param $x i32) (result i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$ret_0 i32)
		(local $$t0 i32)
		(local $$t1 i32)
		(local $$t2 i32)
		(local $$t3 i32)
		(local $$t4 i32)
		(local $$t5 i32)
		(local $$t6.0 i32)
		(local $$t6.1 i32)
		(local $$t7 i32)
		(local $$t8 i32)
		(local $$t9 i32)
		(local $$t10 i32)
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
										local.get $a.2
										local.set $$t0
										br $$Block_2
									end
									i32.const 1
									local.set $$current_block
									local.get $$t1
									local.get $$t2
									i32.sub
									local.set $$t3
									local.get $$t3
									i64.const 1
									i32.wrap_i64
									i32.shr_s
									local.set $$t4
									local.get $$t2
									local.get $$t4
									i32.add
									local.set $$t5
									local.get $a.0
									call $runtime.Block.Retain
									local.get $a.1
									i32.const 4
									local.get $$t5
									i32.mul
									i32.add
									local.set $$t6.1
									local.get $$t6.0
									call $runtime.Block.Release
									local.set $$t6.0
									local.get $$t6.1
									i32.load
									local.set $$t7
									local.get $$t7
									local.get $x
									i32.lt_u
									local.set $$t8
									local.get $$t8
									if
										br $$Block_3
									else
										br $$Block_4
									end
								end
								i32.const 2
								local.set $$current_block
								local.get $$t2
								local.set $$ret_0
								br $$BlockFnBody
							end
							local.get $$current_block
							i32.const 0
							i32.eq
							if(result i32)
								i32.const 0
							else
								local.get $$current_block
								i32.const 4
								i32.eq
								if(result i32)
									local.get $$t9
								else
									local.get $$t2
								end
							end
							local.get $$current_block
							i32.const 0
							i32.eq
							if(result i32)
								local.get $$t0
							else
								local.get $$current_block
								i32.const 4
								i32.eq
								if(result i32)
									local.get $$t1
								else
									local.get $$t5
								end
							end
							local.set $$t1
							local.set $$t2
							i32.const 3
							local.set $$current_block
							local.get $$t2
							local.get $$t1
							i32.lt_s
							local.set $$t10
							local.get $$t10
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
						local.get $$t5
						i32.const 1
						i32.add
						local.set $$t9
						i32.const 3
						local.set $$block_selector
						br $$BlockDisp
					end
					i32.const 5
					local.set $$current_block
					i32.const 3
					local.set $$block_selector
					br $$BlockDisp
				end
			end
		end
		local.get $$ret_0
		local.get $$t6.0
		call $runtime.Block.Release
	)
	(func $strconv.formatBits (param $dst.0 i32) (param $dst.1 i32) (param $dst.2 i32) (param $dst.3 i32) (param $u i64) (param $base i32) (param $neg i32) (param $append_ i32) (result i32 i32 i32 i32 i32 i32 i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$ret_0.0 i32)
		(local $$ret_0.1 i32)
		(local $$ret_0.2 i32)
		(local $$ret_0.3 i32)
		(local $$ret_1.0 i32)
		(local $$ret_1.1 i32)
		(local $$ret_1.2 i32)
		(local $$t0 i32)
		(local $$t1.0 i32)
		(local $$t1.1 i32)
		(local $$t2 i32)
		(local $$t3 i64)
		(local $$t4 i64)
		(local $$t5 i32)
		(local $$t6 i32)
		(local $$t7 i32)
		(local $$t8 i32)
		(local $$t9 i32)
		(local $$t10 i32)
		(local $$t11 i32)
		(local $$t12 i64)
		(local $$t13 i64)
		(local $$t14 i32)
		(local $$t15 i32)
		(local $$t16 i32)
		(local $$t17 i64)
		(local $$t18 i64)
		(local $$t19 i64)
		(local $$t20 i32)
		(local $$t21 i32)
		(local $$t22 i32)
		(local $$t23 i32)
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
		(local $$t33 i32)
		(local $$t34.0 i32)
		(local $$t34.1 i32)
		(local $$t35 i32)
		(local $$t36 i32)
		(local $$t37 i32)
		(local $$t38 i32)
		(local $$t39.0 i32)
		(local $$t39.1 i32)
		(local $$t40 i32)
		(local $$t41 i32)
		(local $$t42 i32)
		(local $$t43 i32)
		(local $$t44 i32)
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
		(local $$t54 i32)
		(local $$t55.0 i32)
		(local $$t55.1 i32)
		(local $$t56 i32)
		(local $$t57 i32)
		(local $$t58 i32)
		(local $$t59.0 i32)
		(local $$t59.1 i32)
		(local $$t60 i32)
		(local $$t61 i32)
		(local $$t62 i32)
		(local $$t63 i32)
		(local $$t64.0 i32)
		(local $$t64.1 i32)
		(local $$t65 i32)
		(local $$t66 i32)
		(local $$t67 i32)
		(local $$t68 i32)
		(local $$t69 i32)
		(local $$t70 i64)
		(local $$t71 i32)
		(local $$t72 i32)
		(local $$t73 i64)
		(local $$t74 i32)
		(local $$t75 i32)
		(local $$t76.0 i32)
		(local $$t76.1 i32)
		(local $$t77 i64)
		(local $$t78 i32)
		(local $$t79 i32)
		(local $$t80 i32)
		(local $$t81 i64)
		(local $$t82 i64)
		(local $$t83.0 i32)
		(local $$t83.1 i32)
		(local $$t84 i32)
		(local $$t85 i32)
		(local $$t86 i32)
		(local $$t87 i32)
		(local $$t88 i32)
		(local $$t89 i64)
		(local $$t90 i64)
		(local $$t91.0 i32)
		(local $$t91.1 i32)
		(local $$t92 i64)
		(local $$t93 i64)
		(local $$t94 i32)
		(local $$t95 i32)
		(local $$t96.0 i32)
		(local $$t96.1 i32)
		(local $$t97 i32)
		(local $$t98 i32)
		(local $$t99 i32)
		(local $$t100 i32)
		(local $$t101.0 i32)
		(local $$t101.1 i32)
		(local $$t102 i32)
		(local $$t103.0 i32)
		(local $$t103.1 i32)
		(local $$t103.2 i32)
		(local $$t103.3 i32)
		(local $$t104.0 i32)
		(local $$t104.1 i32)
		(local $$t104.2 i32)
		(local $$t104.3 i32)
		(local $$t105.0 i32)
		(local $$t105.1 i32)
		(local $$t105.2 i32)
		(local $$t105.3 i32)
		(local $$t106.0 i32)
		(local $$t106.1 i32)
		(local $$t106.2 i32)
		block $$BlockFnBody
			loop $$BlockDisp
				block $$Block_30
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
																																				br_table  0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 0
																																			end
																																			i32.const 0
																																			local.set $$current_block
																																			local.get $base
																																			i32.const 2
																																			i32.lt_s
																																			local.set $$t0
																																			local.get $$t0
																																			if
																																				br $$Block_0
																																			else
																																				br $$Block_2
																																			end
																																		end
																																		i32.const 1
																																		local.set $$current_block
																																		i32.const 32786
																																		i32.const 41
																																		i32.const 32827
																																		i32.const 12
																																		call $$runtime.panic_
																																	end
																																	i32.const 2
																																	local.set $$current_block
																																	i32.const 81
																																	call $runtime.HeapAlloc
																																	i32.const 1
																																	i32.const 0
																																	i32.const 65
																																	call $runtime.Block.Init
																																	call $runtime.DupI32
																																	i32.const 16
																																	i32.add
																																	local.set $$t1.1
																																	local.get $$t1.0
																																	call $runtime.Block.Release
																																	local.set $$t1.0
																																	local.get $neg
																																	if
																																		br $$Block_3
																																	else
																																		br $$Block_4
																																	end
																																end
																																i32.const 3
																																local.set $$current_block
																																local.get $base
																																i32.const 36
																																i32.gt_s
																																local.set $$t2
																																local.get $$t2
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
																															i64.const 0
																															local.get $u
																															i64.sub
																															local.set $$t3
																															br $$Block_4
																														end
																														local.get $$current_block
																														i32.const 2
																														i32.eq
																														if(result i64)
																															local.get $u
																														else
																															local.get $$t3
																														end
																														local.set $$t4
																														i32.const 5
																														local.set $$current_block
																														local.get $base
																														i32.const 10
																														i32.eq
																														local.set $$t5
																														local.get $$t5
																														if
																															br $$Block_5
																														else
																															br $$Block_7
																														end
																													end
																													i32.const 6
																													local.set $$current_block
																													i32.const 1
																													if
																														br $$Block_10
																													else
																														br $$Block_8
																													end
																												end
																												local.get $$current_block
																												i32.const 16
																												i32.eq
																												if(result i32)
																													local.get $$t6
																												else
																													local.get $$current_block
																													i32.const 22
																													i32.eq
																													if(result i32)
																														local.get $$t7
																													else
																														local.get $$current_block
																														i32.const 18
																														i32.eq
																														if(result i32)
																															local.get $$t8
																														else
																															local.get $$t9
																														end
																													end
																												end
																												local.set $$t10
																												i32.const 7
																												local.set $$current_block
																												local.get $neg
																												if
																													br $$Block_26
																												else
																													br $$Block_27
																												end
																											end
																											i32.const 8
																											local.set $$current_block
																											local.get $base
																											call $strconv.isPowerOfTwo
																											local.set $$t11
																											local.get $$t11
																											if
																												br $$Block_18
																											else
																												br $$Block_19
																											end
																										end
																										local.get $$current_block
																										i32.const 6
																										i32.eq
																										if(result i64)
																											local.get $$t4
																										else
																											local.get $$t12
																										end
																										local.get $$current_block
																										i32.const 6
																										i32.eq
																										if(result i32)
																											i32.const 65
																										else
																											local.get $$t14
																										end
																										local.set $$t15
																										local.set $$t13
																										i32.const 9
																										local.set $$current_block
																										local.get $$t13
																										i32.wrap_i64
																										local.set $$t16
																										br $$Block_16
																									end
																									i32.const 10
																									local.set $$current_block
																									local.get $$t12
																									i64.const 1000000000
																									i64.div_u
																									local.set $$t17
																									local.get $$t17
																									i64.const 1000000000
																									i64.mul
																									local.set $$t18
																									local.get $$t12
																									local.get $$t18
																									i64.sub
																									local.set $$t19
																									local.get $$t19
																									i32.wrap_i64
																									local.set $$t20
																									br $$Block_13
																								end
																								local.get $$current_block
																								i32.const 6
																								i32.eq
																								if(result i64)
																									local.get $$t4
																								else
																									local.get $$t17
																								end
																								local.get $$current_block
																								i32.const 6
																								i32.eq
																								if(result i32)
																									i32.const 65
																								else
																									local.get $$t21
																								end
																								local.set $$t14
																								local.set $$t12
																								i32.const 11
																								local.set $$current_block
																								local.get $$t12
																								i64.const 1000000000
																								i64.ge_u
																								local.set $$t22
																								local.get $$t22
																								if
																									i32.const 10
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
																							local.get $$t23
																							i32.const 100
																							i32.rem_u
																							local.set $$t24
																							local.get $$t24
																							i32.const 2
																							i32.mul
																							local.set $$t25
																							local.get $$t23
																							i32.const 100
																							i32.div_u
																							local.set $$t26
																							local.get $$t27
																							i32.const 2
																							i32.sub
																							local.set $$t28
																							local.get $$t28
																							i32.const 1
																							i32.add
																							local.set $$t29
																							local.get $$t1.0
																							call $runtime.Block.Retain
																							local.get $$t1.1
																							i32.const 1
																							local.get $$t29
																							i32.mul
																							i32.add
																							local.set $$t30.1
																							local.get $$t30.0
																							call $runtime.Block.Release
																							local.set $$t30.0
																							local.get $$t25
																							i32.const 1
																							i32.add
																							local.set $$t31
																							i32.const 32839
																							local.get $$t31
																							i32.add
																							i32.load8_u align=1
																							local.set $$t32
																							local.get $$t30.1
																							local.get $$t32
																							i32.store8 align=1
																							local.get $$t28
																							i32.const 0
																							i32.add
																							local.set $$t33
																							local.get $$t1.0
																							call $runtime.Block.Retain
																							local.get $$t1.1
																							i32.const 1
																							local.get $$t33
																							i32.mul
																							i32.add
																							local.set $$t34.1
																							local.get $$t34.0
																							call $runtime.Block.Release
																							local.set $$t34.0
																							local.get $$t25
																							i32.const 0
																							i32.add
																							local.set $$t35
																							i32.const 32839
																							local.get $$t35
																							i32.add
																							i32.load8_u align=1
																							local.set $$t36
																							local.get $$t34.1
																							local.get $$t36
																							i32.store8 align=1
																							local.get $$t37
																							i32.const 1
																							i32.sub
																							local.set $$t38
																							br $$Block_13
																						end
																						i32.const 13
																						local.set $$current_block
																						local.get $$t27
																						i32.const 1
																						i32.sub
																						local.set $$t21
																						local.get $$t1.0
																						call $runtime.Block.Retain
																						local.get $$t1.1
																						i32.const 1
																						local.get $$t21
																						i32.mul
																						i32.add
																						local.set $$t39.1
																						local.get $$t39.0
																						call $runtime.Block.Release
																						local.set $$t39.0
																						local.get $$t23
																						i32.const 2
																						i32.mul
																						local.set $$t40
																						local.get $$t40
																						i32.const 1
																						i32.add
																						local.set $$t41
																						i32.const 32839
																						local.get $$t41
																						i32.add
																						i32.load8_u align=1
																						local.set $$t42
																						local.get $$t39.1
																						local.get $$t42
																						i32.store8 align=1
																						i32.const 11
																						local.set $$block_selector
																						br $$BlockDisp
																					end
																					local.get $$current_block
																					i32.const 10
																					i32.eq
																					if(result i32)
																						local.get $$t14
																					else
																						local.get $$t28
																					end
																					local.get $$current_block
																					i32.const 10
																					i32.eq
																					if(result i32)
																						local.get $$t20
																					else
																						local.get $$t26
																					end
																					local.get $$current_block
																					i32.const 10
																					i32.eq
																					if(result i32)
																						i32.const 4
																					else
																						local.get $$t38
																					end
																					local.set $$t37
																					local.set $$t23
																					local.set $$t27
																					i32.const 14
																					local.set $$current_block
																					local.get $$t37
																					i32.const 0
																					i32.gt_s
																					local.set $$t43
																					local.get $$t43
																					if
																						i32.const 12
																						local.set $$block_selector
																						br $$BlockDisp
																					else
																						i32.const 13
																						local.set $$block_selector
																						br $$BlockDisp
																					end
																				end
																				i32.const 15
																				local.set $$current_block
																				local.get $$t44
																				i32.const 100
																				i32.rem_u
																				local.set $$t45
																				local.get $$t45
																				i32.const 2
																				i32.mul
																				local.set $$t46
																				local.get $$t44
																				i32.const 100
																				i32.div_u
																				local.set $$t47
																				local.get $$t48
																				i32.const 2
																				i32.sub
																				local.set $$t49
																				local.get $$t49
																				i32.const 1
																				i32.add
																				local.set $$t50
																				local.get $$t1.0
																				call $runtime.Block.Retain
																				local.get $$t1.1
																				i32.const 1
																				local.get $$t50
																				i32.mul
																				i32.add
																				local.set $$t51.1
																				local.get $$t51.0
																				call $runtime.Block.Release
																				local.set $$t51.0
																				local.get $$t46
																				i32.const 1
																				i32.add
																				local.set $$t52
																				i32.const 32839
																				local.get $$t52
																				i32.add
																				i32.load8_u align=1
																				local.set $$t53
																				local.get $$t51.1
																				local.get $$t53
																				i32.store8 align=1
																				local.get $$t49
																				i32.const 0
																				i32.add
																				local.set $$t54
																				local.get $$t1.0
																				call $runtime.Block.Retain
																				local.get $$t1.1
																				i32.const 1
																				local.get $$t54
																				i32.mul
																				i32.add
																				local.set $$t55.1
																				local.get $$t55.0
																				call $runtime.Block.Release
																				local.set $$t55.0
																				local.get $$t46
																				i32.const 0
																				i32.add
																				local.set $$t56
																				i32.const 32839
																				local.get $$t56
																				i32.add
																				i32.load8_u align=1
																				local.set $$t57
																				local.get $$t55.1
																				local.get $$t57
																				i32.store8 align=1
																				br $$Block_16
																			end
																			i32.const 16
																			local.set $$current_block
																			local.get $$t44
																			i32.const 2
																			i32.mul
																			local.set $$t58
																			local.get $$t48
																			i32.const 1
																			i32.sub
																			local.set $$t6
																			local.get $$t1.0
																			call $runtime.Block.Retain
																			local.get $$t1.1
																			i32.const 1
																			local.get $$t6
																			i32.mul
																			i32.add
																			local.set $$t59.1
																			local.get $$t59.0
																			call $runtime.Block.Release
																			local.set $$t59.0
																			local.get $$t58
																			i32.const 1
																			i32.add
																			local.set $$t60
																			i32.const 32839
																			local.get $$t60
																			i32.add
																			i32.load8_u align=1
																			local.set $$t61
																			local.get $$t59.1
																			local.get $$t61
																			i32.store8 align=1
																			local.get $$t44
																			i32.const 10
																			i32.ge_u
																			local.set $$t62
																			local.get $$t62
																			if
																				br $$Block_17
																			else
																				i32.const 7
																				local.set $$block_selector
																				br $$BlockDisp
																			end
																		end
																		local.get $$current_block
																		i32.const 9
																		i32.eq
																		if(result i32)
																			local.get $$t15
																		else
																			local.get $$t49
																		end
																		local.get $$current_block
																		i32.const 9
																		i32.eq
																		if(result i32)
																			local.get $$t16
																		else
																			local.get $$t47
																		end
																		local.set $$t44
																		local.set $$t48
																		i32.const 17
																		local.set $$current_block
																		local.get $$t44
																		i32.const 100
																		i32.ge_u
																		local.set $$t63
																		local.get $$t63
																		if
																			i32.const 15
																			local.set $$block_selector
																			br $$BlockDisp
																		else
																			i32.const 16
																			local.set $$block_selector
																			br $$BlockDisp
																		end
																	end
																	i32.const 18
																	local.set $$current_block
																	local.get $$t6
																	i32.const 1
																	i32.sub
																	local.set $$t8
																	local.get $$t1.0
																	call $runtime.Block.Retain
																	local.get $$t1.1
																	i32.const 1
																	local.get $$t8
																	i32.mul
																	i32.add
																	local.set $$t64.1
																	local.get $$t64.0
																	call $runtime.Block.Release
																	local.set $$t64.0
																	i32.const 32839
																	local.get $$t58
																	i32.add
																	i32.load8_u align=1
																	local.set $$t65
																	local.get $$t64.1
																	local.get $$t65
																	i32.store8 align=1
																	i32.const 7
																	local.set $$block_selector
																	br $$BlockDisp
																end
																i32.const 19
																local.set $$current_block
																local.get $base
																local.set $$t66
																local.get $$t66
																call $math$bits.TrailingZeros
																local.set $$t67
																local.get $$t67
																local.set $$t68
																local.get $$t68
																i32.const 7
																i32.and
																local.set $$t69
																local.get $base
																i64.extend_i32_u
																local.set $$t70
																local.get $base
																local.set $$t71
																local.get $$t71
																i32.const 1
																i32.sub
																local.set $$t72
																br $$Block_22
															end
															i32.const 20
															local.set $$current_block
															local.get $base
															i64.extend_i32_u
															local.set $$t73
															br $$Block_25
														end
														i32.const 21
														local.set $$current_block
														local.get $$t74
														i32.const 1
														i32.sub
														local.set $$t75
														local.get $$t1.0
														call $runtime.Block.Retain
														local.get $$t1.1
														i32.const 1
														local.get $$t75
														i32.mul
														i32.add
														local.set $$t76.1
														local.get $$t76.0
														call $runtime.Block.Release
														local.set $$t76.0
														local.get $$t77
														i32.wrap_i64
														local.set $$t78
														local.get $$t78
														local.get $$t72
														i32.and
														local.set $$t79
														i32.const 33039
														local.get $$t79
														i32.add
														i32.load8_u align=1
														local.set $$t80
														local.get $$t76.1
														local.get $$t80
														i32.store8 align=1
														local.get $$t69
														i64.extend_i32_u
														local.set $$t81
														local.get $$t77
														local.get $$t81
														i64.shr_u
														local.set $$t82
														br $$Block_22
													end
													i32.const 22
													local.set $$current_block
													local.get $$t74
													i32.const 1
													i32.sub
													local.set $$t7
													local.get $$t1.0
													call $runtime.Block.Retain
													local.get $$t1.1
													i32.const 1
													local.get $$t7
													i32.mul
													i32.add
													local.set $$t83.1
													local.get $$t83.0
													call $runtime.Block.Release
													local.set $$t83.0
													local.get $$t77
													i32.wrap_i64
													local.set $$t84
													i32.const 33039
													local.get $$t84
													i32.add
													i32.load8_u align=1
													local.set $$t85
													local.get $$t83.1
													local.get $$t85
													i32.store8 align=1
													i32.const 7
													local.set $$block_selector
													br $$BlockDisp
												end
												local.get $$current_block
												i32.const 19
												i32.eq
												if(result i64)
													local.get $$t4
												else
													local.get $$t82
												end
												local.get $$current_block
												i32.const 19
												i32.eq
												if(result i32)
													i32.const 65
												else
													local.get $$t75
												end
												local.set $$t74
												local.set $$t77
												i32.const 23
												local.set $$current_block
												local.get $$t77
												local.get $$t70
												i64.ge_u
												local.set $$t86
												local.get $$t86
												if
													i32.const 21
													local.set $$block_selector
													br $$BlockDisp
												else
													i32.const 22
													local.set $$block_selector
													br $$BlockDisp
												end
											end
											i32.const 24
											local.set $$current_block
											local.get $$t87
											i32.const 1
											i32.sub
											local.set $$t88
											local.get $$t89
											local.get $$t73
											i64.div_u
											local.set $$t90
											local.get $$t1.0
											call $runtime.Block.Retain
											local.get $$t1.1
											i32.const 1
											local.get $$t88
											i32.mul
											i32.add
											local.set $$t91.1
											local.get $$t91.0
											call $runtime.Block.Release
											local.set $$t91.0
											local.get $$t90
											local.get $$t73
											i64.mul
											local.set $$t92
											local.get $$t89
											local.get $$t92
											i64.sub
											local.set $$t93
											local.get $$t93
											i32.wrap_i64
											local.set $$t94
											i32.const 33039
											local.get $$t94
											i32.add
											i32.load8_u align=1
											local.set $$t95
											local.get $$t91.1
											local.get $$t95
											i32.store8 align=1
											br $$Block_25
										end
										i32.const 25
										local.set $$current_block
										local.get $$t87
										i32.const 1
										i32.sub
										local.set $$t9
										local.get $$t1.0
										call $runtime.Block.Retain
										local.get $$t1.1
										i32.const 1
										local.get $$t9
										i32.mul
										i32.add
										local.set $$t96.1
										local.get $$t96.0
										call $runtime.Block.Release
										local.set $$t96.0
										local.get $$t89
										i32.wrap_i64
										local.set $$t97
										i32.const 33039
										local.get $$t97
										i32.add
										i32.load8_u align=1
										local.set $$t98
										local.get $$t96.1
										local.get $$t98
										i32.store8 align=1
										i32.const 7
										local.set $$block_selector
										br $$BlockDisp
									end
									local.get $$current_block
									i32.const 20
									i32.eq
									if(result i64)
										local.get $$t4
									else
										local.get $$t90
									end
									local.get $$current_block
									i32.const 20
									i32.eq
									if(result i32)
										i32.const 65
									else
										local.get $$t88
									end
									local.set $$t87
									local.set $$t89
									i32.const 26
									local.set $$current_block
									local.get $$t89
									local.get $$t73
									i64.ge_u
									local.set $$t99
									local.get $$t99
									if
										i32.const 24
										local.set $$block_selector
										br $$BlockDisp
									else
										i32.const 25
										local.set $$block_selector
										br $$BlockDisp
									end
								end
								i32.const 27
								local.set $$current_block
								local.get $$t10
								i32.const 1
								i32.sub
								local.set $$t100
								local.get $$t1.0
								call $runtime.Block.Retain
								local.get $$t1.1
								i32.const 1
								local.get $$t100
								i32.mul
								i32.add
								local.set $$t101.1
								local.get $$t101.0
								call $runtime.Block.Release
								local.set $$t101.0
								local.get $$t101.1
								i32.const 45
								i32.store8 align=1
								br $$Block_27
							end
							local.get $$current_block
							i32.const 7
							i32.eq
							if(result i32)
								local.get $$t10
							else
								local.get $$t100
							end
							local.set $$t102
							i32.const 28
							local.set $$current_block
							local.get $append_
							if
								br $$Block_28
							else
								br $$Block_29
							end
						end
						i32.const 29
						local.set $$current_block
						local.get $$t1.0
						call $runtime.Block.Retain
						local.get $$t1.1
						i32.const 1
						local.get $$t102
						i32.mul
						i32.add
						i32.const 65
						local.get $$t102
						i32.sub
						i32.const 65
						local.get $$t102
						i32.sub
						local.set $$t103.3
						local.set $$t103.2
						local.set $$t103.1
						local.get $$t103.0
						call $runtime.Block.Release
						local.set $$t103.0
						local.get $dst.0
						local.get $dst.1
						local.get $dst.2
						local.get $dst.3
						local.get $$t103.0
						local.get $$t103.1
						local.get $$t103.2
						local.get $$t103.3
						call $$u8.$slice.append
						local.set $$t104.3
						local.set $$t104.2
						local.set $$t104.1
						local.get $$t104.0
						call $runtime.Block.Release
						local.set $$t104.0
						local.get $$t104.0
						call $runtime.Block.Retain
						local.get $$t104.1
						local.get $$t104.2
						local.get $$t104.3
						local.set $$ret_0.3
						local.set $$ret_0.2
						local.set $$ret_0.1
						local.get $$ret_0.0
						call $runtime.Block.Release
						local.set $$ret_0.0
						i32.const 0
						i32.const 14784
						i32.const 0
						local.set $$ret_1.2
						local.set $$ret_1.1
						local.get $$ret_1.0
						call $runtime.Block.Release
						local.set $$ret_1.0
						br $$BlockFnBody
					end
					i32.const 30
					local.set $$current_block
					local.get $$t1.0
					call $runtime.Block.Retain
					local.get $$t1.1
					i32.const 1
					local.get $$t102
					i32.mul
					i32.add
					i32.const 65
					local.get $$t102
					i32.sub
					i32.const 65
					local.get $$t102
					i32.sub
					local.set $$t105.3
					local.set $$t105.2
					local.set $$t105.1
					local.get $$t105.0
					call $runtime.Block.Release
					local.set $$t105.0
					i32.const 0
					i32.const 14784
					i32.const 0
					local.get $$t105.0
					local.get $$t105.1
					local.get $$t105.2
					call $$string.appendstr
					local.set $$t106.2
					local.set $$t106.1
					local.get $$t106.0
					call $runtime.Block.Release
					local.set $$t106.0
					i32.const 0
					i32.const 0
					i32.const 0
					i32.const 0
					local.set $$ret_0.3
					local.set $$ret_0.2
					local.set $$ret_0.1
					local.get $$ret_0.0
					call $runtime.Block.Release
					local.set $$ret_0.0
					local.get $$t106.0
					call $runtime.Block.Retain
					local.get $$t106.1
					local.get $$t106.2
					local.set $$ret_1.2
					local.set $$ret_1.1
					local.get $$ret_1.0
					call $runtime.Block.Release
					local.set $$ret_1.0
					br $$BlockFnBody
				end
			end
		end
		local.get $$ret_0.0
		call $runtime.Block.Retain
		local.get $$ret_0.1
		local.get $$ret_0.2
		local.get $$ret_0.3
		local.get $$ret_1.0
		call $runtime.Block.Retain
		local.get $$ret_1.1
		local.get $$ret_1.2
		local.get $$ret_0.0
		call $runtime.Block.Release
		local.get $$ret_1.0
		call $runtime.Block.Release
		local.get $$t1.0
		call $runtime.Block.Release
		local.get $$t30.0
		call $runtime.Block.Release
		local.get $$t34.0
		call $runtime.Block.Release
		local.get $$t39.0
		call $runtime.Block.Release
		local.get $$t51.0
		call $runtime.Block.Release
		local.get $$t55.0
		call $runtime.Block.Release
		local.get $$t59.0
		call $runtime.Block.Release
		local.get $$t64.0
		call $runtime.Block.Release
		local.get $$t76.0
		call $runtime.Block.Release
		local.get $$t83.0
		call $runtime.Block.Release
		local.get $$t91.0
		call $runtime.Block.Release
		local.get $$t96.0
		call $runtime.Block.Release
		local.get $$t101.0
		call $runtime.Block.Release
		local.get $$t103.0
		call $runtime.Block.Release
		local.get $$t104.0
		call $runtime.Block.Release
		local.get $$t105.0
		call $runtime.Block.Release
		local.get $$t106.0
		call $runtime.Block.Release
	)
	(func $strconv.init
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$t0 i32)
		(local $$t1.0.0 i32)
		(local $$t1.0.1 i32)
		(local $$t1.1 i32)
		(local $$t1.2 i32)
		(local $$t2.0.0 i32)
		(local $$t2.0.1 i32)
		(local $$t2.1 i32)
		(local $$t2.2 i32)
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
							global.get $strconv.init$guard
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
						global.set $strconv.init$guard
						call $math.init
						call $errors.init
						call $math$bits.init
						call $unicode$utf8.init
						i32.const 0
						i32.const 33135
						i32.const 18
						call $errors.New
						local.set $$t1.2
						local.set $$t1.1
						local.set $$t1.0.1
						local.get $$t1.0.0
						call $runtime.Block.Release
						local.set $$t1.0.0
						i32.const 15000
						local.get $$t1.0.0
						call $runtime.Block.Retain
						i32.const 15000
						i32.load align=1
						call $runtime.Block.Release
						i32.store align=1
						i32.const 15000
						local.get $$t1.0.1
						i32.store offset=4
						i32.const 15000
						local.get $$t1.1
						i32.store offset=8
						i32.const 15000
						local.get $$t1.2
						i32.store offset=12
						i32.const 0
						i32.const 33153
						i32.const 14
						call $errors.New
						local.set $$t2.2
						local.set $$t2.1
						local.set $$t2.0.1
						local.get $$t2.0.0
						call $runtime.Block.Release
						local.set $$t2.0.0
						i32.const 15016
						local.get $$t2.0.0
						call $runtime.Block.Retain
						i32.const 15016
						i32.load align=1
						call $runtime.Block.Release
						i32.store align=1
						i32.const 15016
						local.get $$t2.0.1
						i32.store offset=4
						i32.const 15016
						local.get $$t2.1
						i32.store offset=8
						i32.const 15016
						local.get $$t2.2
						i32.store offset=12
						br $$Block_1
					end
					i32.const 2
					local.set $$current_block
					br $$BlockFnBody
				end
			end
		end
		local.get $$t1.0.0
		call $runtime.Block.Release
		local.get $$t2.0.0
		call $runtime.Block.Release
	)
	(func $strconv.isInGraphicList (param $r i32) (result i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$ret_0 i32)
		(local $$t0 i32)
		(local $$t1 i32)
		(local $$t2.0 i32)
		(local $$t2.1 i32)
		(local $$t2.2 i32)
		(local $$t2.3 i32)
		(local $$t3 i32)
		(local $$t4 i32)
		(local $$t5.0 i32)
		(local $$t5.1 i32)
		(local $$t6 i32)
		(local $$t7 i32)
		(local $$t8 i32)
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
									local.get $r
									i32.const 65535
									i32.gt_s
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
								i32.const 0
								local.set $$ret_0
								br $$BlockFnBody
							end
							i32.const 2
							local.set $$current_block
							local.get $r
							i32.const 65535
							i32.and
							local.set $$t1
							i32.const 0
							i32.const 26424
							i32.const 2
							i32.const 0
							i32.mul
							i32.add
							i32.const 16
							i32.const 0
							i32.sub
							i32.const 16
							i32.const 0
							i32.sub
							local.set $$t2.3
							local.set $$t2.2
							local.set $$t2.1
							local.get $$t2.0
							call $runtime.Block.Release
							local.set $$t2.0
							local.get $$t2.0
							local.get $$t2.1
							local.get $$t2.2
							local.get $$t2.3
							local.get $$t1
							call $strconv.bsearch16
							local.set $$t3
							local.get $$t3
							i32.const 16
							i32.lt_s
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
						i32.const 0
						i32.const 26424
						i32.const 2
						local.get $$t3
						i32.mul
						i32.add
						local.set $$t5.1
						local.get $$t5.0
						call $runtime.Block.Release
						local.set $$t5.0
						local.get $$t5.1
						i32.load16_u
						local.set $$t6
						local.get $$t1
						local.get $$t6
						i32.eq
						local.set $$t7
						br $$Block_3
					end
					local.get $$current_block
					i32.const 2
					i32.eq
					if(result i32)
						i32.const 0
					else
						local.get $$t7
					end
					local.set $$t8
					i32.const 4
					local.set $$current_block
					local.get $$t8
					local.set $$ret_0
					br $$BlockFnBody
				end
			end
		end
		local.get $$ret_0
		local.get $$t2.0
		call $runtime.Block.Release
		local.get $$t5.0
		call $runtime.Block.Release
	)
	(func $strconv.isPowerOfTwo (param $x i32) (result i32)
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
					local.get $x
					i32.const 1
					i32.sub
					local.set $$t0
					local.get $x
					local.get $$t0
					i32.and
					local.set $$t1
					local.get $$t1
					i32.const 0
					i32.eq
					local.set $$t2
					local.get $$t2
					local.set $$ret_0
					br $$BlockFnBody
				end
			end
		end
		local.get $$ret_0
	)
	(func $strconv.quoteWith (param $s.0 i32) (param $s.1 i32) (param $s.2 i32) (param $quote i32) (param $ASCIIonly i32) (param $graphicOnly i32) (result i32 i32 i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$ret_0.0 i32)
		(local $$ret_0.1 i32)
		(local $$ret_0.2 i32)
		(local $$t0 i32)
		(local $$t1 i32)
		(local $$t2 i32)
		(local $$t3.0 i32)
		(local $$t3.1 i32)
		(local $$t3.2 i32)
		(local $$t3.3 i32)
		(local $$t4.0 i32)
		(local $$t4.1 i32)
		(local $$t4.2 i32)
		(local $$t4.3 i32)
		(local $$t5.0 i32)
		(local $$t5.1 i32)
		(local $$t5.2 i32)
		block $$BlockFnBody
			loop $$BlockDisp
				block $$Block_0
					block $$BlockSel
						local.get $$block_selector
						br_table  0 0
					end
					i32.const 0
					local.set $$current_block
					local.get $s.2
					local.set $$t0
					i32.const 3
					local.get $$t0
					i32.mul
					local.set $$t1
					local.get $$t1
					i32.const 2
					i32.div_s
					local.set $$t2
					local.get $$t2
					i32.const 1
					i32.mul
					i32.const 16
					i32.add
					call $runtime.HeapAlloc
					local.get $$t2
					i32.const 0
					i32.const 1
					call $runtime.Block.Init
					call $runtime.DupI32
					i32.const 16
					i32.add
					i32.const 0
					local.get $$t2
					local.set $$t3.3
					local.set $$t3.2
					local.set $$t3.1
					local.get $$t3.0
					call $runtime.Block.Release
					local.set $$t3.0
					local.get $$t3.0
					local.get $$t3.1
					local.get $$t3.2
					local.get $$t3.3
					local.get $s.0
					local.get $s.1
					local.get $s.2
					local.get $quote
					local.get $ASCIIonly
					local.get $graphicOnly
					call $strconv.appendQuotedWith
					local.set $$t4.3
					local.set $$t4.2
					local.set $$t4.1
					local.get $$t4.0
					call $runtime.Block.Release
					local.set $$t4.0
					i32.const 0
					i32.const 14784
					i32.const 0
					local.get $$t4.0
					local.get $$t4.1
					local.get $$t4.2
					call $$string.appendstr
					local.set $$t5.2
					local.set $$t5.1
					local.get $$t5.0
					call $runtime.Block.Release
					local.set $$t5.0
					local.get $$t5.0
					call $runtime.Block.Retain
					local.get $$t5.1
					local.get $$t5.2
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
		local.get $$t3.0
		call $runtime.Block.Release
		local.get $$t4.0
		call $runtime.Block.Release
		local.get $$t5.0
		call $runtime.Block.Release
	)
	(func $strconv.small (param $i i32) (result i32 i32 i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$ret_0.0 i32)
		(local $$ret_0.1 i32)
		(local $$ret_0.2 i32)
		(local $$t0 i32)
		(local $$t1 i32)
		(local $$t2.0 i32)
		(local $$t2.1 i32)
		(local $$t2.2 i32)
		(local $$t3 i32)
		(local $$t4 i32)
		(local $$t5 i32)
		(local $$t6.0 i32)
		(local $$t6.1 i32)
		(local $$t6.2 i32)
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
							local.get $i
							i32.const 10
							i32.lt_s
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
						local.get $i
						i32.const 1
						i32.add
						local.set $$t1
						i32.const 0
						i32.const 33039
						local.get $i
						i32.add
						local.get $$t1
						local.get $i
						i32.sub
						local.set $$t2.2
						local.set $$t2.1
						local.get $$t2.0
						call $runtime.Block.Release
						local.set $$t2.0
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
					end
					i32.const 2
					local.set $$current_block
					local.get $i
					i32.const 2
					i32.mul
					local.set $$t3
					local.get $$t3
					i32.const 2
					i32.add
					local.set $$t4
					local.get $i
					i32.const 2
					i32.mul
					local.set $$t5
					i32.const 0
					i32.const 32839
					local.get $$t5
					i32.add
					local.get $$t4
					local.get $$t5
					i32.sub
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
			end
		end
		local.get $$ret_0.0
		call $runtime.Block.Retain
		local.get $$ret_0.1
		local.get $$ret_0.2
		local.get $$ret_0.0
		call $runtime.Block.Release
		local.get $$t2.0
		call $runtime.Block.Release
		local.get $$t6.0
		call $runtime.Block.Release
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
	(func $syscall$wasm4.LineI32 (param $x1 i32) (param $y1 i32) (param $x2 i32) (param $y2 i32)
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
					local.get $x1
					local.get $y1
					local.get $x2
					local.get $y2
					call $syscall$wasm4.__import__line
					br $$BlockFnBody
				end
			end
		end
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
	(func $syscall$wasm4.SetDrawColorsU16 (param $x i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$t0.0 i32)
		(local $$t0.1 i32)
		(local $$t0.2 i32)
		(local $$t0.3 i32)
		(local $$t1.0 i32)
		(local $$t1.1 i32)
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
					local.get $x
					i32.store16
					br $$BlockFnBody
				end
			end
		end
		local.get $$t0.0
		call $runtime.Block.Release
		local.get $$t1.0
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
	(func $syscall$wasm4.Text (param $s.0 i32) (param $s.1 i32) (param $s.2 i32) (param $x i32) (param $y i32)
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
						br_table  0 0
					end
					i32.const 0
					local.set $$current_block
					local.get $s.0
					local.get $s.1
					local.get $s.2
					call $$syscall/wasm4.__linkname__string_data_ptr
					local.set $$t0
					local.get $s.2
					local.set $$t1
					local.get $$t1
					local.set $$t2
					local.get $$t0
					local.get $$t2
					local.get $x
					local.get $y
					call $syscall$wasm4.__import__textUtf8
					br $$BlockFnBody
				end
			end
		end
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
	(func $syscall$wasm4.Trace (param $msg.0 i32) (param $msg.1 i32) (param $msg.2 i32)
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
						br_table  0 0
					end
					i32.const 0
					local.set $$current_block
					local.get $msg.0
					local.get $msg.1
					local.get $msg.2
					call $$syscall/wasm4.__linkname__string_data_ptr
					local.set $$t0
					local.get $msg.2
					local.set $$t1
					local.get $$t1
					local.set $$t2
					local.get $$t0
					local.get $$t2
					call $syscall$wasm4.__import__traceUtf8
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
	(func $unicode$utf8.DecodeRuneInString (param $s.0 i32) (param $s.1 i32) (param $s.2 i32) (result i32 i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$ret_0 i32)
		(local $$ret_1 i32)
		(local $$t0 i32)
		(local $$t1 i32)
		(local $$t2 i32)
		(local $$t3 i32)
		(local $$t4.0 i32)
		(local $$t4.1 i32)
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
		(local $$t16 i32)
		(local $$t17.0 i32)
		(local $$t17.1 i32)
		(local $$t18 i32)
		(local $$t19 i32)
		(local $$t20.0 i32)
		(local $$t20.1 i32)
		(local $$t21.0 i32)
		(local $$t21.1 i32)
		(local $$t22 i32)
		(local $$t23 i32)
		(local $$t24.0 i32)
		(local $$t24.1 i32)
		(local $$t25 i32)
		(local $$t26 i32)
		(local $$t27 i32)
		(local $$t28.0 i32)
		(local $$t28.1 i32)
		(local $$t29 i32)
		(local $$t30 i32)
		(local $$t31 i32)
		(local $$t32 i32)
		(local $$t33 i32)
		(local $$t34 i32)
		(local $$t35 i32)
		(local $$t36 i32)
		(local $$t37 i32)
		(local $$t38 i32)
		(local $$t39 i32)
		(local $$t40 i32)
		(local $$t41 i32)
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
		(local $$t54 i32)
		(local $$t55 i32)
		(local $$t56 i32)
		(local $$t57 i32)
		(local $$t58 i32)
		(local $$t59 i32)
		(local $$t60 i32)
		(local $$t61 i32)
		(local $$t62 i32)
		(local $$t63 i32)
		(local $$t64 i32)
		(local $$t65 i32)
		(local $$t66 i32)
		(local $$t67 i32)
		block $$BlockFnBody
			loop $$BlockDisp
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
																									br_table  0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 0
																								end
																								i32.const 0
																								local.set $$current_block
																								local.get $s.2
																								local.set $$t0
																								local.get $$t0
																								i32.const 1
																								i32.lt_s
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
																							i32.const 65533
																							local.set $$ret_0
																							i32.const 0
																							local.set $$ret_1
																							br $$BlockFnBody
																						end
																						i32.const 2
																						local.set $$current_block
																						local.get $s.1
																						i32.const 0
																						i32.add
																						i32.load8_u align=1
																						local.set $$t2
																						local.get $$t2
																						local.set $$t3
																						i32.const 0
																						i32.const 30864
																						i32.const 1
																						local.get $$t3
																						i32.mul
																						i32.add
																						local.set $$t4.1
																						local.get $$t4.0
																						call $runtime.Block.Release
																						local.set $$t4.0
																						local.get $$t4.1
																						i32.load8_u align=1
																						local.set $$t5
																						local.get $$t5
																						i32.const 240
																						i32.ge_u
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
																					local.get $$t5
																					local.set $$t7
																					local.get $$t7
																					i64.const 31
																					i32.wrap_i64
																					i32.shl
																					local.set $$t8
																					local.get $$t8
																					i64.const 31
																					i32.wrap_i64
																					i32.shr_s
																					local.set $$t9
																					local.get $s.1
																					i32.const 0
																					i32.add
																					i32.load8_u align=1
																					local.set $$t10
																					local.get $$t10
																					local.set $$t11
																					local.get $$t11
																					local.get $$t9
																					i32.const -1
																					i32.xor
																					i32.and
																					local.set $$t12
																					i32.const 65533
																					local.get $$t9
																					i32.and
																					local.set $$t13
																					local.get $$t12
																					local.get $$t13
																					i32.or
																					local.set $$t14
																					local.get $$t14
																					local.set $$ret_0
																					i32.const 1
																					local.set $$ret_1
																					br $$BlockFnBody
																				end
																				i32.const 4
																				local.set $$current_block
																				local.get $$t5
																				i32.const 7
																				i32.and
																				local.set $$t15
																				local.get $$t15
																				local.set $$t16
																				i32.const 18
																				call $runtime.HeapAlloc
																				i32.const 1
																				i32.const 0
																				i32.const 2
																				call $runtime.Block.Init
																				call $runtime.DupI32
																				i32.const 16
																				i32.add
																				local.set $$t17.1
																				local.get $$t17.0
																				call $runtime.Block.Release
																				local.set $$t17.0
																				local.get $$t5
																				i64.const 4
																				i32.wrap_i64
																				i32.shr_u
																				local.set $$t18
																				local.get $$t18
																				local.set $$t19
																				i32.const 0
																				i32.const 30832
																				i32.const 2
																				local.get $$t19
																				i32.mul
																				i32.add
																				local.set $$t20.1
																				local.get $$t20.0
																				call $runtime.Block.Release
																				local.set $$t20.0
																				local.get $$t20.1
																				i32.load8_u align=1
																				local.get $$t20.1
																				i32.load8_u offset=1 align=1
																				local.set $$t21.1
																				local.set $$t21.0
																				local.get $$t17.1
																				local.get $$t21.0
																				i32.store8 align=1
																				local.get $$t17.1
																				local.get $$t21.1
																				i32.store8 offset=1 align=1
																				local.get $$t0
																				local.get $$t16
																				i32.lt_s
																				local.set $$t22
																				local.get $$t22
																				if
																					br $$Block_4
																				else
																					br $$Block_5
																				end
																			end
																			i32.const 5
																			local.set $$current_block
																			i32.const 65533
																			local.set $$ret_0
																			i32.const 1
																			local.set $$ret_1
																			br $$BlockFnBody
																		end
																		i32.const 6
																		local.set $$current_block
																		local.get $s.1
																		i32.const 1
																		i32.add
																		i32.load8_u align=1
																		local.set $$t23
																		local.get $$t17.0
																		call $runtime.Block.Retain
																		local.get $$t17.1
																		i32.const 0
																		i32.add
																		local.set $$t24.1
																		local.get $$t24.0
																		call $runtime.Block.Release
																		local.set $$t24.0
																		local.get $$t24.1
																		i32.load8_u align=1
																		local.set $$t25
																		local.get $$t23
																		local.get $$t25
																		i32.lt_u
																		local.set $$t26
																		local.get $$t26
																		if
																			br $$Block_6
																		else
																			br $$Block_8
																		end
																	end
																	i32.const 7
																	local.set $$current_block
																	i32.const 65533
																	local.set $$ret_0
																	i32.const 1
																	local.set $$ret_1
																	br $$BlockFnBody
																end
																i32.const 8
																local.set $$current_block
																local.get $$t16
																i32.const 2
																i32.le_s
																local.set $$t27
																local.get $$t27
																if
																	br $$Block_9
																else
																	br $$Block_10
																end
															end
															i32.const 9
															local.set $$current_block
															local.get $$t17.0
															call $runtime.Block.Retain
															local.get $$t17.1
															i32.const 1
															i32.add
															local.set $$t28.1
															local.get $$t28.0
															call $runtime.Block.Release
															local.set $$t28.0
															local.get $$t28.1
															i32.load8_u align=1
															local.set $$t29
															local.get $$t29
															local.get $$t23
															i32.lt_u
															local.set $$t30
															local.get $$t30
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
														local.get $$t2
														i32.const 31
														i32.and
														local.set $$t31
														local.get $$t31
														local.set $$t32
														local.get $$t32
														i64.const 6
														i32.wrap_i64
														i32.shl
														local.set $$t33
														local.get $$t23
														i32.const 63
														i32.and
														local.set $$t34
														local.get $$t34
														local.set $$t35
														local.get $$t33
														local.get $$t35
														i32.or
														local.set $$t36
														local.get $$t36
														local.set $$ret_0
														i32.const 2
														local.set $$ret_1
														br $$BlockFnBody
													end
													i32.const 11
													local.set $$current_block
													local.get $s.1
													i32.const 2
													i32.add
													i32.load8_u align=1
													local.set $$t37
													local.get $$t37
													i32.const 128
													i32.lt_u
													local.set $$t38
													local.get $$t38
													if
														br $$Block_11
													else
														br $$Block_13
													end
												end
												i32.const 12
												local.set $$current_block
												i32.const 65533
												local.set $$ret_0
												i32.const 1
												local.set $$ret_1
												br $$BlockFnBody
											end
											i32.const 13
											local.set $$current_block
											local.get $$t16
											i32.const 3
											i32.le_s
											local.set $$t39
											local.get $$t39
											if
												br $$Block_14
											else
												br $$Block_15
											end
										end
										i32.const 14
										local.set $$current_block
										i32.const 191
										local.get $$t37
										i32.lt_u
										local.set $$t40
										local.get $$t40
										if
											i32.const 12
											local.set $$block_selector
											br $$BlockDisp
										else
											i32.const 13
											local.set $$block_selector
											br $$BlockDisp
										end
									end
									i32.const 15
									local.set $$current_block
									local.get $$t2
									i32.const 15
									i32.and
									local.set $$t41
									local.get $$t41
									local.set $$t42
									local.get $$t42
									i64.const 12
									i32.wrap_i64
									i32.shl
									local.set $$t43
									local.get $$t23
									i32.const 63
									i32.and
									local.set $$t44
									local.get $$t44
									local.set $$t45
									local.get $$t45
									i64.const 6
									i32.wrap_i64
									i32.shl
									local.set $$t46
									local.get $$t43
									local.get $$t46
									i32.or
									local.set $$t47
									local.get $$t37
									i32.const 63
									i32.and
									local.set $$t48
									local.get $$t48
									local.set $$t49
									local.get $$t47
									local.get $$t49
									i32.or
									local.set $$t50
									local.get $$t50
									local.set $$ret_0
									i32.const 3
									local.set $$ret_1
									br $$BlockFnBody
								end
								i32.const 16
								local.set $$current_block
								local.get $s.1
								i32.const 3
								i32.add
								i32.load8_u align=1
								local.set $$t51
								local.get $$t51
								i32.const 128
								i32.lt_u
								local.set $$t52
								local.get $$t52
								if
									br $$Block_16
								else
									br $$Block_18
								end
							end
							i32.const 17
							local.set $$current_block
							i32.const 65533
							local.set $$ret_0
							i32.const 1
							local.set $$ret_1
							br $$BlockFnBody
						end
						i32.const 18
						local.set $$current_block
						local.get $$t2
						i32.const 7
						i32.and
						local.set $$t53
						local.get $$t53
						local.set $$t54
						local.get $$t54
						i64.const 18
						i32.wrap_i64
						i32.shl
						local.set $$t55
						local.get $$t23
						i32.const 63
						i32.and
						local.set $$t56
						local.get $$t56
						local.set $$t57
						local.get $$t57
						i64.const 12
						i32.wrap_i64
						i32.shl
						local.set $$t58
						local.get $$t55
						local.get $$t58
						i32.or
						local.set $$t59
						local.get $$t37
						i32.const 63
						i32.and
						local.set $$t60
						local.get $$t60
						local.set $$t61
						local.get $$t61
						i64.const 6
						i32.wrap_i64
						i32.shl
						local.set $$t62
						local.get $$t59
						local.get $$t62
						i32.or
						local.set $$t63
						local.get $$t51
						i32.const 63
						i32.and
						local.set $$t64
						local.get $$t64
						local.set $$t65
						local.get $$t63
						local.get $$t65
						i32.or
						local.set $$t66
						local.get $$t66
						local.set $$ret_0
						i32.const 4
						local.set $$ret_1
						br $$BlockFnBody
					end
					i32.const 19
					local.set $$current_block
					i32.const 191
					local.get $$t51
					i32.lt_u
					local.set $$t67
					local.get $$t67
					if
						i32.const 17
						local.set $$block_selector
						br $$BlockDisp
					else
						i32.const 18
						local.set $$block_selector
						br $$BlockDisp
					end
				end
			end
		end
		local.get $$ret_0
		local.get $$ret_1
		local.get $$t4.0
		call $runtime.Block.Release
		local.get $$t17.0
		call $runtime.Block.Release
		local.get $$t20.0
		call $runtime.Block.Release
		local.get $$t24.0
		call $runtime.Block.Release
		local.get $$t28.0
		call $runtime.Block.Release
	)
	(func $unicode$utf8.EncodeRune (param $p.0 i32) (param $p.1 i32) (param $p.2 i32) (param $p.3 i32) (param $r i32) (result i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$ret_0 i32)
		(local $$t0 i32)
		(local $$t1 i32)
		(local $$t2.0 i32)
		(local $$t2.1 i32)
		(local $$t3 i32)
		(local $$t4.0 i32)
		(local $$t4.1 i32)
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
		(local $$t13 i32)
		(local $$t14 i32)
		(local $$t15.0 i32)
		(local $$t15.1 i32)
		(local $$t16 i32)
		(local $$t17.0 i32)
		(local $$t17.1 i32)
		(local $$t18 i32)
		(local $$t19 i32)
		(local $$t20 i32)
		(local $$t21.0 i32)
		(local $$t21.1 i32)
		(local $$t22 i32)
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
		(local $$t33.0 i32)
		(local $$t33.1 i32)
		(local $$t34 i32)
		(local $$t35 i32)
		(local $$t36 i32)
		(local $$t37.0 i32)
		(local $$t37.1 i32)
		(local $$t38 i32)
		(local $$t39 i32)
		(local $$t40 i32)
		(local $$t41 i32)
		(local $$t42.0 i32)
		(local $$t42.1 i32)
		(local $$t43 i32)
		(local $$t44 i32)
		(local $$t45 i32)
		(local $$t46 i32)
		(local $$t47 i32)
		(local $$t48 i32)
		(local $$t49 i32)
		(local $$t50.0 i32)
		(local $$t50.1 i32)
		(local $$t51 i32)
		(local $$t52.0 i32)
		(local $$t52.1 i32)
		(local $$t53 i32)
		(local $$t54 i32)
		(local $$t55 i32)
		(local $$t56.0 i32)
		(local $$t56.1 i32)
		(local $$t57 i32)
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
		block $$BlockFnBody
			loop $$BlockDisp
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
																	br_table  0 1 2 3 4 5 6 7 8 9 10 11 0
																end
																i32.const 0
																local.set $$current_block
																local.get $r
																local.set $$t0
																local.get $$t0
																i32.const 127
																i32.le_u
																local.set $$t1
																local.get $$t1
																if
																	br $$Block_0
																else
																	br $$Block_2
																end
															end
															i32.const 1
															local.set $$current_block
															local.get $p.0
															call $runtime.Block.Retain
															local.get $p.1
															i32.const 1
															i32.const 0
															i32.mul
															i32.add
															local.set $$t2.1
															local.get $$t2.0
															call $runtime.Block.Release
															local.set $$t2.0
															local.get $r
															i32.const 255
															i32.and
															local.set $$t3
															local.get $$t2.1
															local.get $$t3
															i32.store8 align=1
															i32.const 1
															local.set $$ret_0
															br $$BlockFnBody
														end
														i32.const 2
														local.set $$current_block
														local.get $p.0
														call $runtime.Block.Retain
														local.get $p.1
														i32.const 1
														i32.const 1
														i32.mul
														i32.add
														local.set $$t4.1
														local.get $$t4.0
														call $runtime.Block.Release
														local.set $$t4.0
														local.get $$t4.1
														i32.load8_u align=1
														local.set $$t5
														local.get $p.0
														call $runtime.Block.Retain
														local.get $p.1
														i32.const 1
														i32.const 0
														i32.mul
														i32.add
														local.set $$t6.1
														local.get $$t6.0
														call $runtime.Block.Release
														local.set $$t6.0
														local.get $r
														i64.const 6
														i32.wrap_i64
														i32.shr_s
														local.set $$t7
														local.get $$t7
														i32.const 255
														i32.and
														local.set $$t8
														i32.const 192
														local.get $$t8
														i32.or
														local.set $$t9
														local.get $$t6.1
														local.get $$t9
														i32.store8 align=1
														local.get $p.0
														call $runtime.Block.Retain
														local.get $p.1
														i32.const 1
														i32.const 1
														i32.mul
														i32.add
														local.set $$t10.1
														local.get $$t10.0
														call $runtime.Block.Release
														local.set $$t10.0
														local.get $r
														i32.const 255
														i32.and
														local.set $$t11
														local.get $$t11
														i32.const 63
														i32.and
														local.set $$t12
														i32.const 128
														local.get $$t12
														i32.or
														local.set $$t13
														local.get $$t10.1
														local.get $$t13
														i32.store8 align=1
														i32.const 2
														local.set $$ret_0
														br $$BlockFnBody
													end
													i32.const 3
													local.set $$current_block
													local.get $$t0
													i32.const 2047
													i32.le_u
													local.set $$t14
													local.get $$t14
													if
														i32.const 2
														local.set $$block_selector
														br $$BlockDisp
													else
														br $$Block_4
													end
												end
												i32.const 4
												local.set $$current_block
												local.get $p.0
												call $runtime.Block.Retain
												local.get $p.1
												i32.const 1
												i32.const 2
												i32.mul
												i32.add
												local.set $$t15.1
												local.get $$t15.0
												call $runtime.Block.Release
												local.set $$t15.0
												local.get $$t15.1
												i32.load8_u align=1
												local.set $$t16
												local.get $p.0
												call $runtime.Block.Retain
												local.get $p.1
												i32.const 1
												i32.const 0
												i32.mul
												i32.add
												local.set $$t17.1
												local.get $$t17.0
												call $runtime.Block.Release
												local.set $$t17.0
												i32.const 65533
												i64.const 12
												i32.wrap_i64
												i32.shr_s
												local.set $$t18
												local.get $$t18
												i32.const 255
												i32.and
												local.set $$t19
												i32.const 224
												local.get $$t19
												i32.or
												local.set $$t20
												local.get $$t17.1
												local.get $$t20
												i32.store8 align=1
												local.get $p.0
												call $runtime.Block.Retain
												local.get $p.1
												i32.const 1
												i32.const 1
												i32.mul
												i32.add
												local.set $$t21.1
												local.get $$t21.0
												call $runtime.Block.Release
												local.set $$t21.0
												i32.const 65533
												i64.const 6
												i32.wrap_i64
												i32.shr_s
												local.set $$t22
												local.get $$t22
												i32.const 255
												i32.and
												local.set $$t23
												local.get $$t23
												i32.const 63
												i32.and
												local.set $$t24
												i32.const 128
												local.get $$t24
												i32.or
												local.set $$t25
												local.get $$t21.1
												local.get $$t25
												i32.store8 align=1
												local.get $p.0
												call $runtime.Block.Retain
												local.get $p.1
												i32.const 1
												i32.const 2
												i32.mul
												i32.add
												local.set $$t26.1
												local.get $$t26.0
												call $runtime.Block.Release
												local.set $$t26.0
												i32.const 65533
												i32.const 255
												i32.and
												local.set $$t27
												local.get $$t27
												i32.const 63
												i32.and
												local.set $$t28
												i32.const 128
												local.get $$t28
												i32.or
												local.set $$t29
												local.get $$t26.1
												local.get $$t29
												i32.store8 align=1
												i32.const 3
												local.set $$ret_0
												br $$BlockFnBody
											end
											i32.const 5
											local.set $$current_block
											local.get $$t0
											i32.const 1114111
											i32.gt_u
											local.set $$t30
											local.get $$t30
											if
												i32.const 4
												local.set $$block_selector
												br $$BlockDisp
											else
												br $$Block_6
											end
										end
										i32.const 6
										local.set $$current_block
										local.get $p.0
										call $runtime.Block.Retain
										local.get $p.1
										i32.const 1
										i32.const 2
										i32.mul
										i32.add
										local.set $$t31.1
										local.get $$t31.0
										call $runtime.Block.Release
										local.set $$t31.0
										local.get $$t31.1
										i32.load8_u align=1
										local.set $$t32
										local.get $p.0
										call $runtime.Block.Retain
										local.get $p.1
										i32.const 1
										i32.const 0
										i32.mul
										i32.add
										local.set $$t33.1
										local.get $$t33.0
										call $runtime.Block.Release
										local.set $$t33.0
										local.get $r
										i64.const 12
										i32.wrap_i64
										i32.shr_s
										local.set $$t34
										local.get $$t34
										i32.const 255
										i32.and
										local.set $$t35
										i32.const 224
										local.get $$t35
										i32.or
										local.set $$t36
										local.get $$t33.1
										local.get $$t36
										i32.store8 align=1
										local.get $p.0
										call $runtime.Block.Retain
										local.get $p.1
										i32.const 1
										i32.const 1
										i32.mul
										i32.add
										local.set $$t37.1
										local.get $$t37.0
										call $runtime.Block.Release
										local.set $$t37.0
										local.get $r
										i64.const 6
										i32.wrap_i64
										i32.shr_s
										local.set $$t38
										local.get $$t38
										i32.const 255
										i32.and
										local.set $$t39
										local.get $$t39
										i32.const 63
										i32.and
										local.set $$t40
										i32.const 128
										local.get $$t40
										i32.or
										local.set $$t41
										local.get $$t37.1
										local.get $$t41
										i32.store8 align=1
										local.get $p.0
										call $runtime.Block.Retain
										local.get $p.1
										i32.const 1
										i32.const 2
										i32.mul
										i32.add
										local.set $$t42.1
										local.get $$t42.0
										call $runtime.Block.Release
										local.set $$t42.0
										local.get $r
										i32.const 255
										i32.and
										local.set $$t43
										local.get $$t43
										i32.const 63
										i32.and
										local.set $$t44
										i32.const 128
										local.get $$t44
										i32.or
										local.set $$t45
										local.get $$t42.1
										local.get $$t45
										i32.store8 align=1
										i32.const 3
										local.set $$ret_0
										br $$BlockFnBody
									end
									i32.const 7
									local.set $$current_block
									i32.const 55296
									local.get $$t0
									i32.le_u
									local.set $$t46
									local.get $$t46
									if
										br $$Block_8
									else
										br $$Block_9
									end
								end
								i32.const 8
								local.set $$current_block
								local.get $$t0
								i32.const 65535
								i32.le_u
								local.set $$t47
								local.get $$t47
								if
									i32.const 6
									local.set $$block_selector
									br $$BlockDisp
								else
									br $$Block_10
								end
							end
							i32.const 9
							local.set $$current_block
							local.get $$t0
							i32.const 57343
							i32.le_u
							local.set $$t48
							br $$Block_9
						end
						local.get $$current_block
						i32.const 7
						i32.eq
						if(result i32)
							i32.const 0
						else
							local.get $$t48
						end
						local.set $$t49
						i32.const 10
						local.set $$current_block
						local.get $$t49
						if
							i32.const 4
							local.set $$block_selector
							br $$BlockDisp
						else
							i32.const 8
							local.set $$block_selector
							br $$BlockDisp
						end
					end
					i32.const 11
					local.set $$current_block
					local.get $p.0
					call $runtime.Block.Retain
					local.get $p.1
					i32.const 1
					i32.const 3
					i32.mul
					i32.add
					local.set $$t50.1
					local.get $$t50.0
					call $runtime.Block.Release
					local.set $$t50.0
					local.get $$t50.1
					i32.load8_u align=1
					local.set $$t51
					local.get $p.0
					call $runtime.Block.Retain
					local.get $p.1
					i32.const 1
					i32.const 0
					i32.mul
					i32.add
					local.set $$t52.1
					local.get $$t52.0
					call $runtime.Block.Release
					local.set $$t52.0
					local.get $r
					i64.const 18
					i32.wrap_i64
					i32.shr_s
					local.set $$t53
					local.get $$t53
					i32.const 255
					i32.and
					local.set $$t54
					i32.const 240
					local.get $$t54
					i32.or
					local.set $$t55
					local.get $$t52.1
					local.get $$t55
					i32.store8 align=1
					local.get $p.0
					call $runtime.Block.Retain
					local.get $p.1
					i32.const 1
					i32.const 1
					i32.mul
					i32.add
					local.set $$t56.1
					local.get $$t56.0
					call $runtime.Block.Release
					local.set $$t56.0
					local.get $r
					i64.const 12
					i32.wrap_i64
					i32.shr_s
					local.set $$t57
					local.get $$t57
					i32.const 255
					i32.and
					local.set $$t58
					local.get $$t58
					i32.const 63
					i32.and
					local.set $$t59
					i32.const 128
					local.get $$t59
					i32.or
					local.set $$t60
					local.get $$t56.1
					local.get $$t60
					i32.store8 align=1
					local.get $p.0
					call $runtime.Block.Retain
					local.get $p.1
					i32.const 1
					i32.const 2
					i32.mul
					i32.add
					local.set $$t61.1
					local.get $$t61.0
					call $runtime.Block.Release
					local.set $$t61.0
					local.get $r
					i64.const 6
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
					local.get $p.0
					call $runtime.Block.Retain
					local.get $p.1
					i32.const 1
					i32.const 3
					i32.mul
					i32.add
					local.set $$t66.1
					local.get $$t66.0
					call $runtime.Block.Release
					local.set $$t66.0
					local.get $r
					i32.const 255
					i32.and
					local.set $$t67
					local.get $$t67
					i32.const 63
					i32.and
					local.set $$t68
					i32.const 128
					local.get $$t68
					i32.or
					local.set $$t69
					local.get $$t66.1
					local.get $$t69
					i32.store8 align=1
					i32.const 4
					local.set $$ret_0
					br $$BlockFnBody
				end
			end
		end
		local.get $$ret_0
		local.get $$t2.0
		call $runtime.Block.Release
		local.get $$t4.0
		call $runtime.Block.Release
		local.get $$t6.0
		call $runtime.Block.Release
		local.get $$t10.0
		call $runtime.Block.Release
		local.get $$t15.0
		call $runtime.Block.Release
		local.get $$t17.0
		call $runtime.Block.Release
		local.get $$t21.0
		call $runtime.Block.Release
		local.get $$t26.0
		call $runtime.Block.Release
		local.get $$t31.0
		call $runtime.Block.Release
		local.get $$t33.0
		call $runtime.Block.Release
		local.get $$t37.0
		call $runtime.Block.Release
		local.get $$t42.0
		call $runtime.Block.Release
		local.get $$t50.0
		call $runtime.Block.Release
		local.get $$t52.0
		call $runtime.Block.Release
		local.get $$t56.0
		call $runtime.Block.Release
		local.get $$t61.0
		call $runtime.Block.Release
		local.get $$t66.0
		call $runtime.Block.Release
	)
	(func $unicode$utf8.init
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
							global.get $unicode$utf8.init$guard
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
						global.set $unicode$utf8.init$guard
						br $$Block_1
					end
					i32.const 2
					local.set $$current_block
					br $$BlockFnBody
				end
			end
		end
	)
	(func $$w42048$game.Board.underlying.$$onFree (param $$ptr i32)
		local.get $$ptr
		i32.const 10
		call_indirect 0 (type $$onFree)
	)
	(func $$w42048.UI.$$onFree (param $$ptr i32)
		local.get $$ptr
		i32.const 28
		call_indirect 0 (type $$onFree)
	)
	(func $w42048.NewUI (export "w42048.NewUI") (param $board.0.0 i32) (param $board.0.1 i32) (param $board.1 i32) (param $board.2 i32) (param $pal.0 i32) (param $pal.1 i32) (param $pal.2 i32) (param $pal.3 i32) (result i32 i32)
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
		(local $$t4.2 i32)
		(local $$t4.3 i32)
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
					i32.const 0
					i32.const 16
					call $runtime.Block.Init
					call $runtime.DupI32
					i32.const 16
					i32.add
					local.set $$t0.1
					local.get $$t0.0
					call $runtime.Block.Release
					local.set $$t0.0
					local.get $$t0.1
					local.get $pal.0
					i32.store
					local.get $$t0.1
					local.get $pal.1
					i32.store offset=4
					local.get $$t0.1
					local.get $pal.2
					i32.store offset=8
					local.get $$t0.1
					local.get $pal.3
					i32.store offset=12
					i32.const 56
					call $runtime.HeapAlloc
					i32.const 1
					i32.const 29
					i32.const 40
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
					local.get $$t1.0
					call $runtime.Block.Retain
					local.get $$t1.1
					i32.const 24
					i32.add
					local.set $$t3.1
					local.get $$t3.0
					call $runtime.Block.Release
					local.set $$t3.0
					local.get $$t0.1
					i32.load
					local.get $$t0.1
					i32.load offset=4
					local.get $$t0.1
					i32.load offset=8
					local.get $$t0.1
					i32.load offset=12
					local.set $$t4.3
					local.set $$t4.2
					local.set $$t4.1
					local.set $$t4.0
					local.get $$t2.1
					local.get $board.0.0
					call $runtime.Block.Retain
					local.get $$t2.1
					i32.load align=1
					call $runtime.Block.Release
					i32.store align=1
					local.get $$t2.1
					local.get $board.0.1
					i32.store offset=4
					local.get $$t2.1
					local.get $board.1
					i32.store offset=8
					local.get $$t2.1
					local.get $board.2
					i32.store offset=12
					local.get $$t3.1
					local.get $$t4.0
					i32.store
					local.get $$t3.1
					local.get $$t4.1
					i32.store offset=4
					local.get $$t3.1
					local.get $$t4.2
					i32.store offset=8
					local.get $$t3.1
					local.get $$t4.3
					i32.store offset=12
					local.get $$t1.0
					call $runtime.Block.Retain
					local.get $$t1.1
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
	)
	(func $w42048.Start (export "start")
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$t0.0.0 i32)
		(local $$t0.0.1 i32)
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
		block $$BlockFnBody
			loop $$BlockDisp
				block $$Block_0
					block $$BlockSel
						local.get $$block_selector
						br_table  0 0
					end
					i32.const 0
					local.set $$current_block
					i32.const 478194671
					call $w42048$game.New
					local.set $$t0.2
					local.set $$t0.1
					local.set $$t0.0.1
					local.get $$t0.0.0
					call $runtime.Block.Release
					local.set $$t0.0.0
					i32.const 31368
					i32.load
					i32.const 31368
					i32.load offset=4
					i32.const 31368
					i32.load offset=8
					i32.const 31368
					i32.load offset=12
					local.set $$t1.3
					local.set $$t1.2
					local.set $$t1.1
					local.set $$t1.0
					local.get $$t1.0
					local.get $$t1.1
					local.get $$t1.2
					local.get $$t1.3
					local.set $$t2.3
					local.set $$t2.2
					local.set $$t2.1
					local.set $$t2.0
					local.get $$t0.0.0
					local.get $$t0.0.1
					local.get $$t0.1
					local.get $$t0.2
					local.get $$t2.0
					local.get $$t2.1
					local.get $$t2.2
					local.get $$t2.3
					call $w42048.NewUI
					local.set $$t3.1
					local.get $$t3.0
					call $runtime.Block.Release
					local.set $$t3.0
					i32.const 31264
					local.get $$t3.0
					call $runtime.Block.Retain
					i32.const 31264
					i32.load align=1
					call $runtime.Block.Release
					i32.store align=1
					i32.const 31264
					local.get $$t3.1
					i32.store offset=4
					br $$BlockFnBody
				end
			end
		end
		local.get $$t0.0.0
		call $runtime.Block.Release
		local.get $$t3.0
		call $runtime.Block.Release
	)
	(func $w42048.Update (export "update")
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$t0.0 i32)
		(local $$t0.1 i32)
		(local $$t1.0 i32)
		(local $$t1.1 i32)
		(local $$t2 i32)
		(local $$t3 i32)
		(local $$t4.0 i32)
		(local $$t4.1 i32)
		(local $$t5.0 i32)
		(local $$t5.1 i32)
		(local $$t6.0 i32)
		(local $$t6.1 i32)
		(local $$t7.0 i32)
		(local $$t7.1 i32)
		block $$BlockFnBody
			loop $$BlockDisp
				block $$Block_0
					block $$BlockSel
						local.get $$block_selector
						br_table  0 0
					end
					i32.const 0
					local.set $$current_block
					i32.const 31264
					i32.load
					call $runtime.Block.Retain
					i32.const 31264
					i32.load offset=4
					local.set $$t0.1
					local.get $$t0.0
					call $runtime.Block.Release
					local.set $$t0.0
					local.get $$t0.0
					call $runtime.Block.Retain
					local.get $$t0.1
					i32.const 16
					i32.add
					local.set $$t1.1
					local.get $$t1.0
					call $runtime.Block.Release
					local.set $$t1.0
					local.get $$t1.1
					i32.load
					local.set $$t2
					local.get $$t2
					i32.const 1
					i32.add
					local.set $$t3
					local.get $$t1.1
					local.get $$t3
					i32.store
					i32.const 31264
					i32.load
					call $runtime.Block.Retain
					i32.const 31264
					i32.load offset=4
					local.set $$t4.1
					local.get $$t4.0
					call $runtime.Block.Release
					local.set $$t4.0
					local.get $$t4.0
					local.get $$t4.1
					call $w42048.UI.show
					i32.const 31264
					i32.load
					call $runtime.Block.Retain
					i32.const 31264
					i32.load offset=4
					local.set $$t5.1
					local.get $$t5.0
					call $runtime.Block.Release
					local.set $$t5.0
					local.get $$t5.0
					local.get $$t5.1
					call $w42048.UI.sound
					i32.const 31264
					i32.load
					call $runtime.Block.Retain
					i32.const 31264
					i32.load offset=4
					local.set $$t6.1
					local.get $$t6.0
					call $runtime.Block.Release
					local.set $$t6.0
					local.get $$t6.0
					local.get $$t6.1
					call $w42048.UI.music
					i32.const 31264
					i32.load
					call $runtime.Block.Retain
					i32.const 31264
					i32.load offset=4
					local.set $$t7.1
					local.get $$t7.0
					call $runtime.Block.Release
					local.set $$t7.0
					local.get $$t7.0
					local.get $$t7.1
					call $w42048.UI.input
					br $$BlockFnBody
				end
			end
		end
		local.get $$t0.0
		call $runtime.Block.Release
		local.get $$t1.0
		call $runtime.Block.Release
		local.get $$t4.0
		call $runtime.Block.Release
		local.get $$t5.0
		call $runtime.Block.Release
		local.get $$t6.0
		call $runtime.Block.Release
		local.get $$t7.0
		call $runtime.Block.Release
	)
	(func $w42048.dotbg (param $x1 i32) (param $y1 i32) (param $w i32) (param $h i32) (param $s i32) (param $dc i32) (param $bg i32)
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
													local.get $x1
												else
													local.get $$t0
												end
												local.set $$t1
												i32.const 3
												local.set $$current_block
												local.get $$t1
												local.get $w
												i32.lt_s
												local.set $$t2
												local.get $$t2
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
											local.get $$t1
											local.get $s
											i32.rem_s
											local.set $$t3
											local.get $$t3
											i32.const 0
											i32.eq
											local.set $$t4
											local.get $$t4
											if
												br $$Block_9
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
										local.get $y1
									else
										local.get $$t5
									end
									local.set $$t6
									i32.const 6
									local.set $$current_block
									local.get $$t6
									local.get $h
									i32.lt_s
									local.set $$t7
									local.get $$t7
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
								local.get $$t1
								local.get $$t6
								local.get $dc
								call $w42048.set
								br $$Block_7
							end
							i32.const 8
							local.set $$current_block
							local.get $$t6
							i32.const 1
							i32.add
							local.set $$t5
							i32.const 6
							local.set $$block_selector
							br $$BlockDisp
						end
						i32.const 9
						local.set $$current_block
						local.get $$t1
						local.get $$t6
						local.get $bg
						call $w42048.set
						i32.const 8
						local.set $$block_selector
						br $$BlockDisp
					end
					i32.const 10
					local.set $$current_block
					local.get $$t6
					local.get $s
					i32.rem_s
					local.set $$t8
					local.get $$t8
					i32.const 0
					i32.eq
					local.set $$t9
					local.get $$t9
					if
						i32.const 7
						local.set $$block_selector
						br $$BlockDisp
					else
						i32.const 9
						local.set $$block_selector
						br $$BlockDisp
					end
				end
			end
		end
	)
	(func $w42048.init
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
							global.get $w42048.init$guard
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
						global.set $w42048.init$guard
						call $runtime.init
						call $strconv.init
						call $syscall$wasm4.init
						call $w42048$game.init
						call $w42048$palettes.init
						br $$Block_1
					end
					i32.const 2
					local.set $$current_block
					br $$BlockFnBody
				end
			end
		end
	)
	(func $w42048.leftpad (param $s.0 i32) (param $s.1 i32) (param $s.2 i32) (param $c.0 i32) (param $c.1 i32) (param $c.2 i32) (param $w i32) (result i32 i32 i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$ret_0.0 i32)
		(local $$ret_0.1 i32)
		(local $$ret_0.2 i32)
		(local $$t0 i32)
		(local $$t1 i32)
		(local $$t2 i32)
		(local $$t3.0 i32)
		(local $$t3.1 i32)
		(local $$t3.2 i32)
		(local $$t4.0 i32)
		(local $$t4.1 i32)
		(local $$t4.2 i32)
		(local $$t5 i32)
		(local $$t6 i32)
		(local $$t7.0 i32)
		(local $$t7.1 i32)
		(local $$t7.2 i32)
		(local $$t8 i32)
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
										local.get $s.2
										local.set $$t0
										local.get $w
										local.get $$t0
										i32.sub
										local.set $$t1
										local.get $$t1
										i32.const 0
										i32.le_s
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
									local.get $s.0
									call $runtime.Block.Retain
									local.get $s.1
									local.get $s.2
									local.set $$ret_0.2
									local.set $$ret_0.1
									local.get $$ret_0.0
									call $runtime.Block.Release
									local.set $$ret_0.0
									br $$BlockFnBody
								end
								i32.const 2
								local.set $$current_block
								br $$Block_4
							end
							i32.const 3
							local.set $$current_block
							local.get $$t3.0
							local.get $$t3.1
							local.get $$t3.2
							local.get $c.0
							local.get $c.1
							local.get $c.2
							call $$string.appendstr
							local.set $$t4.2
							local.set $$t4.1
							local.get $$t4.0
							call $runtime.Block.Release
							local.set $$t4.0
							local.get $$t5
							i32.const 1
							i32.add
							local.set $$t6
							br $$Block_4
						end
						i32.const 4
						local.set $$current_block
						local.get $$t3.0
						local.get $$t3.1
						local.get $$t3.2
						local.get $s.0
						local.get $s.1
						local.get $s.2
						call $$string.appendstr
						local.set $$t7.2
						local.set $$t7.1
						local.get $$t7.0
						call $runtime.Block.Release
						local.set $$t7.0
						local.get $$t7.0
						call $runtime.Block.Retain
						local.get $$t7.1
						local.get $$t7.2
						local.set $$ret_0.2
						local.set $$ret_0.1
						local.get $$ret_0.0
						call $runtime.Block.Release
						local.set $$ret_0.0
						br $$BlockFnBody
					end
					local.get $$current_block
					i32.const 2
					i32.eq
					if(result i32 i32 i32)
						i32.const 0
						i32.const 14784
						i32.const 0
					else
						local.get $$t4.0
						call $runtime.Block.Retain
						local.get $$t4.1
						local.get $$t4.2
					end
					local.get $$current_block
					i32.const 2
					i32.eq
					if(result i32)
						i32.const 0
					else
						local.get $$t6
					end
					local.set $$t5
					local.set $$t3.2
					local.set $$t3.1
					local.get $$t3.0
					call $runtime.Block.Release
					local.set $$t3.0
					i32.const 5
					local.set $$current_block
					local.get $$t5
					local.get $$t1
					i32.lt_s
					local.set $$t8
					local.get $$t8
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
			end
		end
		local.get $$ret_0.0
		call $runtime.Block.Retain
		local.get $$ret_0.1
		local.get $$ret_0.2
		local.get $$ret_0.0
		call $runtime.Block.Release
		local.get $$t3.0
		call $runtime.Block.Release
		local.get $$t4.0
		call $runtime.Block.Release
		local.get $$t7.0
		call $runtime.Block.Release
	)
	(func $w42048.log (param $s.0 i32) (param $s.1 i32) (param $s.2 i32)
		(local $$block_selector i32)
		(local $$current_block i32)
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
							i32.const 0
							if
								br $$Block_1
							else
								br $$Block_0
							end
						end
						i32.const 1
						local.set $$current_block
						br $$BlockFnBody
					end
					i32.const 2
					local.set $$current_block
					local.get $s.0
					local.get $s.1
					local.get $s.2
					call $syscall$wasm4.Trace
					br $$BlockFnBody
				end
			end
		end
	)
	(func $w42048.play (param $s.0 i32) (param $s.1 i32) (param $s.2 i32) (param $s.3 i32) (param $s.4 i32) (param $s.5 i32) (param $s.6 i32) (param $s.7 i32) (param $s.8 i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$t0.0 i32)
		(local $$t0.1 i32)
		(local $$t1.0 i32)
		(local $$t1.1 i32)
		(local $$t2 i32)
		(local $$t3.0 i32)
		(local $$t3.1 i32)
		(local $$t4 i32)
		(local $$t5 i32)
		(local $$t6 i32)
		(local $$t7.0 i32)
		(local $$t7.1 i32)
		(local $$t8 i32)
		(local $$t9 i32)
		(local $$t10.0 i32)
		(local $$t10.1 i32)
		(local $$t11 i32)
		(local $$t12 i32)
		(local $$t13 i32)
		(local $$t14.0 i32)
		(local $$t14.1 i32)
		(local $$t15 i32)
		(local $$t16 i32)
		(local $$t17.0 i32)
		(local $$t17.1 i32)
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
		(local $$t29.0 i32)
		(local $$t29.1 i32)
		(local $$t30 i32)
		(local $$t31 i32)
		(local $$t32 i32)
		block $$BlockFnBody
			loop $$BlockDisp
				block $$Block_0
					block $$BlockSel
						local.get $$block_selector
						br_table  0 0
					end
					i32.const 0
					local.set $$current_block
					i32.const 52
					call $runtime.HeapAlloc
					i32.const 1
					i32.const 0
					i32.const 36
					call $runtime.Block.Init
					call $runtime.DupI32
					i32.const 16
					i32.add
					local.set $$t0.1
					local.get $$t0.0
					call $runtime.Block.Release
					local.set $$t0.0
					local.get $$t0.1
					local.get $s.0
					i32.store
					local.get $$t0.1
					local.get $s.1
					i32.store offset=4
					local.get $$t0.1
					local.get $s.2
					i32.store offset=8
					local.get $$t0.1
					local.get $s.3
					i32.store offset=12
					local.get $$t0.1
					local.get $s.4
					i32.store offset=16
					local.get $$t0.1
					local.get $s.5
					i32.store offset=20
					local.get $$t0.1
					local.get $s.6
					i32.store offset=24
					local.get $$t0.1
					local.get $s.7
					i32.store offset=28
					local.get $$t0.1
					local.get $s.8
					i32.store offset=32
					local.get $$t0.0
					call $runtime.Block.Retain
					local.get $$t0.1
					i32.const 0
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
					local.get $$t4
					i64.const 16
					i32.wrap_i64
					i32.shl
					local.set $$t5
					local.get $$t2
					local.get $$t5
					i32.or
					local.set $$t6
					local.get $$t0.0
					call $runtime.Block.Retain
					local.get $$t0.1
					i32.const 8
					i32.add
					local.set $$t7.1
					local.get $$t7.0
					call $runtime.Block.Release
					local.set $$t7.0
					local.get $$t7.1
					i32.load
					local.set $$t8
					local.get $$t8
					i64.const 24
					i32.wrap_i64
					i32.shl
					local.set $$t9
					local.get $$t0.0
					call $runtime.Block.Retain
					local.get $$t0.1
					i32.const 12
					i32.add
					local.set $$t10.1
					local.get $$t10.0
					call $runtime.Block.Release
					local.set $$t10.0
					local.get $$t10.1
					i32.load
					local.set $$t11
					local.get $$t11
					i64.const 16
					i32.wrap_i64
					i32.shl
					local.set $$t12
					local.get $$t9
					local.get $$t12
					i32.or
					local.set $$t13
					local.get $$t0.0
					call $runtime.Block.Retain
					local.get $$t0.1
					i32.const 16
					i32.add
					local.set $$t14.1
					local.get $$t14.0
					call $runtime.Block.Release
					local.set $$t14.0
					local.get $$t14.1
					i32.load
					local.set $$t15
					local.get $$t13
					local.get $$t15
					i32.or
					local.set $$t16
					local.get $$t0.0
					call $runtime.Block.Retain
					local.get $$t0.1
					i32.const 20
					i32.add
					local.set $$t17.1
					local.get $$t17.0
					call $runtime.Block.Release
					local.set $$t17.0
					local.get $$t17.1
					i32.load
					local.set $$t18
					local.get $$t18
					i64.const 8
					i32.wrap_i64
					i32.shl
					local.set $$t19
					local.get $$t16
					local.get $$t19
					i32.or
					local.set $$t20
					local.get $$t0.0
					call $runtime.Block.Retain
					local.get $$t0.1
					i32.const 28
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
					i32.const 32
					i32.add
					local.set $$t23.1
					local.get $$t23.0
					call $runtime.Block.Release
					local.set $$t23.0
					local.get $$t23.1
					i32.load
					local.set $$t24
					local.get $$t24
					i64.const 2
					i32.wrap_i64
					i32.shl
					local.set $$t25
					local.get $$t22
					local.get $$t25
					i32.or
					local.set $$t26
					local.get $$t6
					local.set $$t27
					local.get $$t20
					local.set $$t28
					local.get $$t0.0
					call $runtime.Block.Retain
					local.get $$t0.1
					i32.const 24
					i32.add
					local.set $$t29.1
					local.get $$t29.0
					call $runtime.Block.Release
					local.set $$t29.0
					local.get $$t29.1
					i32.load
					local.set $$t30
					local.get $$t30
					local.set $$t31
					local.get $$t26
					local.set $$t32
					local.get $$t27
					local.get $$t28
					local.get $$t31
					local.get $$t32
					call $syscall$wasm4.Tone
					br $$BlockFnBody
				end
			end
		end
		local.get $$t0.0
		call $runtime.Block.Release
		local.get $$t1.0
		call $runtime.Block.Release
		local.get $$t3.0
		call $runtime.Block.Release
		local.get $$t7.0
		call $runtime.Block.Release
		local.get $$t10.0
		call $runtime.Block.Release
		local.get $$t14.0
		call $runtime.Block.Release
		local.get $$t17.0
		call $runtime.Block.Release
		local.get $$t21.0
		call $runtime.Block.Release
		local.get $$t23.0
		call $runtime.Block.Release
		local.get $$t29.0
		call $runtime.Block.Release
	)
	(func $w42048.set (param $x i32) (param $y i32) (param $c i32)
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
					local.get $c
					call $syscall$wasm4.SetDrawColorsU16
					local.get $x
					local.get $y
					local.get $x
					local.get $y
					call $syscall$wasm4.LineI32
					br $$BlockFnBody
				end
			end
		end
	)
	(func $w42048.showTile (param $col i32) (param $row i32) (param $val i32)
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
		(local $$t23 i32)
		(local $$t24 i32)
		(local $$t25 i32)
		(local $$t26.0 i32)
		(local $$t26.1 i32)
		(local $$t26.2 i32)
		(local $$t27.0 i32)
		(local $$t27.1 i32)
		(local $$t27.2 i32)
		(local $$t28 i32)
		(local $$t29 i32)
		(local $$t30 i32)
		block $$BlockFnBody
			loop $$BlockDisp
				block $$Block_0
					block $$BlockSel
						local.get $$block_selector
						br_table  0 0
					end
					i32.const 0
					local.set $$current_block
					local.get $col
					i32.const 41
					i32.mul
					local.set $$t0
					i32.const 2
					local.get $$t0
					i32.add
					local.set $$t1
					local.get $row
					i32.const 32
					i32.mul
					local.set $$t2
					i32.const 32
					local.get $$t2
					i32.add
					local.set $$t3
					local.get $val
					call $w42048.tileShadow
					local.get $$t1
					i32.const 1
					i32.sub
					local.set $$t4
					local.get $$t3
					i32.const 1
					i32.sub
					local.set $$t5
					local.get $$t4
					local.get $$t5
					i32.const 35
					i32.const 27
					call $syscall$wasm4.RectI32
					local.get $$t1
					local.get $$t3
					i32.const 67
					call $w42048.set
					i32.const 1
					call $syscall$wasm4.SetDrawColorsU16
					local.get $$t1
					i32.const 1
					i32.sub
					local.set $$t6
					local.get $$t3
					i32.const 1
					i32.sub
					local.set $$t7
					local.get $$t6
					local.get $$t3
					local.get $$t1
					local.get $$t7
					call $syscall$wasm4.LineI32
					i32.const 35
					i32.const 3
					i32.sub
					local.set $$t8
					local.get $$t1
					local.get $$t8
					i32.add
					local.set $$t9
					local.get $$t3
					i32.const 1
					i32.sub
					local.set $$t10
					i32.const 35
					i32.const 2
					i32.sub
					local.set $$t11
					local.get $$t1
					local.get $$t11
					i32.add
					local.set $$t12
					local.get $$t9
					local.get $$t10
					local.get $$t12
					local.get $$t3
					call $syscall$wasm4.LineI32
					local.get $$t1
					i32.const 1
					i32.sub
					local.set $$t13
					local.get $$t3
					i32.const 1
					i32.sub
					local.set $$t14
					local.get $$t13
					local.get $$t14
					i32.const 1
					call $w42048.set
					i32.const 35
					i32.const 2
					i32.sub
					local.set $$t15
					local.get $$t1
					local.get $$t15
					i32.add
					local.set $$t16
					local.get $$t3
					i32.const 1
					i32.sub
					local.set $$t17
					local.get $$t16
					local.get $$t17
					i32.const 1
					call $w42048.set
					local.get $$t1
					i32.const 1
					i32.sub
					local.set $$t18
					local.get $$t3
					i32.const 27
					i32.add
					local.set $$t19
					local.get $$t19
					i32.const 2
					i32.sub
					local.set $$t20
					local.get $$t18
					local.get $$t20
					i32.const 1
					call $w42048.set
					i32.const 35
					i32.const 2
					i32.sub
					local.set $$t21
					local.get $$t1
					local.get $$t21
					i32.add
					local.set $$t22
					local.get $$t3
					i32.const 27
					i32.add
					local.set $$t23
					local.get $$t23
					i32.const 2
					i32.sub
					local.set $$t24
					local.get $$t22
					local.get $$t24
					i32.const 1
					call $w42048.set
					local.get $val
					call $w42048.tileColor
					i32.const 0
					i32.const 34902
					i32.const 4
					local.get $$t1
					local.get $$t3
					call $syscall$wasm4.Text
					local.get $val
					local.set $$t25
					local.get $$t25
					call $strconv.Itoa
					local.set $$t26.2
					local.set $$t26.1
					local.get $$t26.0
					call $runtime.Block.Release
					local.set $$t26.0
					local.get $$t26.0
					local.get $$t26.1
					local.get $$t26.2
					i32.const 0
					i32.const 14981
					i32.const 1
					i32.const 4
					call $w42048.leftpad
					local.set $$t27.2
					local.set $$t27.1
					local.get $$t27.0
					call $runtime.Block.Release
					local.set $$t27.0
					local.get $$t3
					i32.const 8
					i32.add
					local.set $$t28
					local.get $$t27.0
					local.get $$t27.1
					local.get $$t27.2
					local.get $$t1
					local.get $$t28
					call $syscall$wasm4.Text
					i32.const 8
					i32.const 2
					i32.mul
					local.set $$t29
					local.get $$t3
					local.get $$t29
					i32.add
					local.set $$t30
					i32.const 0
					i32.const 34902
					i32.const 4
					local.get $$t1
					local.get $$t30
					call $syscall$wasm4.Text
					br $$BlockFnBody
				end
			end
		end
		local.get $$t26.0
		call $runtime.Block.Release
		local.get $$t27.0
		call $runtime.Block.Release
	)
	(func $w42048.tileColor (param $val i32)
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
		block $$BlockFnBody
			loop $$BlockDisp
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
																												br_table  0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 0
																											end
																											i32.const 0
																											local.set $$current_block
																											local.get $val
																											i32.const 2
																											i32.eq
																											local.set $$t0
																											local.get $$t0
																											if
																												br $$Block_1
																											else
																												br $$Block_3
																											end
																										end
																										i32.const 1
																										local.set $$current_block
																										br $$BlockFnBody
																									end
																									i32.const 2
																									local.set $$current_block
																									i32.const 33
																									call $syscall$wasm4.SetDrawColorsU16
																									i32.const 1
																									local.set $$block_selector
																									br $$BlockDisp
																								end
																								i32.const 3
																								local.set $$current_block
																								i32.const 50
																								call $syscall$wasm4.SetDrawColorsU16
																								i32.const 1
																								local.set $$block_selector
																								br $$BlockDisp
																							end
																							i32.const 4
																							local.set $$current_block
																							local.get $val
																							i32.const 4
																							i32.eq
																							local.set $$t1
																							local.get $$t1
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
																						i32.const 36
																						call $syscall$wasm4.SetDrawColorsU16
																						i32.const 1
																						local.set $$block_selector
																						br $$BlockDisp
																					end
																					i32.const 6
																					local.set $$current_block
																					local.get $val
																					i32.const 8
																					i32.eq
																					local.set $$t2
																					local.get $$t2
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
																				i32.const 50
																				call $syscall$wasm4.SetDrawColorsU16
																				i32.const 1
																				local.set $$block_selector
																				br $$BlockDisp
																			end
																			i32.const 8
																			local.set $$current_block
																			local.get $val
																			i32.const 16
																			i32.eq
																			local.set $$t3
																			local.get $$t3
																			if
																				i32.const 7
																				local.set $$block_selector
																				br $$BlockDisp
																			else
																				br $$Block_9
																			end
																		end
																		i32.const 9
																		local.set $$current_block
																		i32.const 52
																		call $syscall$wasm4.SetDrawColorsU16
																		i32.const 1
																		local.set $$block_selector
																		br $$BlockDisp
																	end
																	i32.const 10
																	local.set $$current_block
																	local.get $val
																	i32.const 32
																	i32.eq
																	local.set $$t4
																	local.get $$t4
																	if
																		i32.const 9
																		local.set $$block_selector
																		br $$BlockDisp
																	else
																		br $$Block_11
																	end
																end
																i32.const 11
																local.set $$current_block
																i32.const 50
																call $syscall$wasm4.SetDrawColorsU16
																i32.const 1
																local.set $$block_selector
																br $$BlockDisp
															end
															i32.const 12
															local.set $$current_block
															local.get $val
															i32.const 64
															i32.eq
															local.set $$t5
															local.get $$t5
															if
																i32.const 11
																local.set $$block_selector
																br $$BlockDisp
															else
																br $$Block_13
															end
														end
														i32.const 13
														local.set $$current_block
														i32.const 64
														call $syscall$wasm4.SetDrawColorsU16
														i32.const 1
														local.set $$block_selector
														br $$BlockDisp
													end
													i32.const 14
													local.set $$current_block
													local.get $val
													i32.const 128
													i32.eq
													local.set $$t6
													local.get $$t6
													if
														i32.const 13
														local.set $$block_selector
														br $$BlockDisp
													else
														br $$Block_15
													end
												end
												i32.const 15
												local.set $$current_block
												i32.const 66
												call $syscall$wasm4.SetDrawColorsU16
												i32.const 1
												local.set $$block_selector
												br $$BlockDisp
											end
											i32.const 16
											local.set $$current_block
											local.get $val
											i32.const 256
											i32.eq
											local.set $$t7
											local.get $$t7
											if
												i32.const 15
												local.set $$block_selector
												br $$BlockDisp
											else
												br $$Block_17
											end
										end
										i32.const 17
										local.set $$current_block
										i32.const 67
										call $syscall$wasm4.SetDrawColorsU16
										i32.const 1
										local.set $$block_selector
										br $$BlockDisp
									end
									i32.const 18
									local.set $$current_block
									local.get $val
									i32.const 512
									i32.eq
									local.set $$t8
									local.get $$t8
									if
										i32.const 17
										local.set $$block_selector
										br $$BlockDisp
									else
										br $$Block_19
									end
								end
								i32.const 19
								local.set $$current_block
								i32.const 52
								call $syscall$wasm4.SetDrawColorsU16
								i32.const 1
								local.set $$block_selector
								br $$BlockDisp
							end
							i32.const 20
							local.set $$current_block
							local.get $val
							i32.const 1024
							i32.eq
							local.set $$t9
							local.get $$t9
							if
								i32.const 19
								local.set $$block_selector
								br $$BlockDisp
							else
								br $$Block_21
							end
						end
						i32.const 21
						local.set $$current_block
						i32.const 36
						call $syscall$wasm4.SetDrawColorsU16
						i32.const 1
						local.set $$block_selector
						br $$BlockDisp
					end
					i32.const 22
					local.set $$current_block
					local.get $val
					i32.const 2048
					i32.eq
					local.set $$t10
					local.get $$t10
					if
						i32.const 21
						local.set $$block_selector
						br $$BlockDisp
					else
						i32.const 1
						local.set $$block_selector
						br $$BlockDisp
					end
				end
			end
		end
	)
	(func $w42048.tileShadow (param $val i32)
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
		block $$BlockFnBody
			loop $$BlockDisp
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
																												br_table  0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 0
																											end
																											i32.const 0
																											local.set $$current_block
																											local.get $val
																											i32.const 2
																											i32.eq
																											local.set $$t0
																											local.get $$t0
																											if
																												br $$Block_1
																											else
																												br $$Block_3
																											end
																										end
																										i32.const 1
																										local.set $$current_block
																										br $$BlockFnBody
																									end
																									i32.const 2
																									local.set $$current_block
																									i32.const 67
																									call $syscall$wasm4.SetDrawColorsU16
																									i32.const 1
																									local.set $$block_selector
																									br $$BlockDisp
																								end
																								i32.const 3
																								local.set $$current_block
																								i32.const 67
																								call $syscall$wasm4.SetDrawColorsU16
																								i32.const 1
																								local.set $$block_selector
																								br $$BlockDisp
																							end
																							i32.const 4
																							local.set $$current_block
																							local.get $val
																							i32.const 4
																							i32.eq
																							local.set $$t1
																							local.get $$t1
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
																						i32.const 66
																						call $syscall$wasm4.SetDrawColorsU16
																						i32.const 1
																						local.set $$block_selector
																						br $$BlockDisp
																					end
																					i32.const 6
																					local.set $$current_block
																					local.get $val
																					i32.const 8
																					i32.eq
																					local.set $$t2
																					local.get $$t2
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
																				i32.const 50
																				call $syscall$wasm4.SetDrawColorsU16
																				i32.const 1
																				local.set $$block_selector
																				br $$BlockDisp
																			end
																			i32.const 8
																			local.set $$current_block
																			local.get $val
																			i32.const 16
																			i32.eq
																			local.set $$t3
																			local.get $$t3
																			if
																				i32.const 7
																				local.set $$block_selector
																				br $$BlockDisp
																			else
																				br $$Block_9
																			end
																		end
																		i32.const 9
																		local.set $$current_block
																		i32.const 52
																		call $syscall$wasm4.SetDrawColorsU16
																		i32.const 1
																		local.set $$block_selector
																		br $$BlockDisp
																	end
																	i32.const 10
																	local.set $$current_block
																	local.get $val
																	i32.const 32
																	i32.eq
																	local.set $$t4
																	local.get $$t4
																	if
																		i32.const 9
																		local.set $$block_selector
																		br $$BlockDisp
																	else
																		br $$Block_11
																	end
																end
																i32.const 11
																local.set $$current_block
																i32.const 50
																call $syscall$wasm4.SetDrawColorsU16
																i32.const 1
																local.set $$block_selector
																br $$BlockDisp
															end
															i32.const 12
															local.set $$current_block
															local.get $val
															i32.const 64
															i32.eq
															local.set $$t5
															local.get $$t5
															if
																i32.const 11
																local.set $$block_selector
																br $$BlockDisp
															else
																br $$Block_13
															end
														end
														i32.const 13
														local.set $$current_block
														i32.const 67
														call $syscall$wasm4.SetDrawColorsU16
														i32.const 1
														local.set $$block_selector
														br $$BlockDisp
													end
													i32.const 14
													local.set $$current_block
													local.get $val
													i32.const 128
													i32.eq
													local.set $$t6
													local.get $$t6
													if
														i32.const 13
														local.set $$block_selector
														br $$BlockDisp
													else
														br $$Block_15
													end
												end
												i32.const 15
												local.set $$current_block
												i32.const 66
												call $syscall$wasm4.SetDrawColorsU16
												i32.const 1
												local.set $$block_selector
												br $$BlockDisp
											end
											i32.const 16
											local.set $$current_block
											local.get $val
											i32.const 256
											i32.eq
											local.set $$t7
											local.get $$t7
											if
												i32.const 15
												local.set $$block_selector
												br $$BlockDisp
											else
												br $$Block_17
											end
										end
										i32.const 17
										local.set $$current_block
										i32.const 36
										call $syscall$wasm4.SetDrawColorsU16
										i32.const 1
										local.set $$block_selector
										br $$BlockDisp
									end
									i32.const 18
									local.set $$current_block
									local.get $val
									i32.const 512
									i32.eq
									local.set $$t8
									local.get $$t8
									if
										i32.const 17
										local.set $$block_selector
										br $$BlockDisp
									else
										br $$Block_19
									end
								end
								i32.const 19
								local.set $$current_block
								i32.const 68
								call $syscall$wasm4.SetDrawColorsU16
								i32.const 1
								local.set $$block_selector
								br $$BlockDisp
							end
							i32.const 20
							local.set $$current_block
							local.get $val
							i32.const 1024
							i32.eq
							local.set $$t9
							local.get $$t9
							if
								i32.const 19
								local.set $$block_selector
								br $$BlockDisp
							else
								br $$Block_21
							end
						end
						i32.const 21
						local.set $$current_block
						i32.const 68
						call $syscall$wasm4.SetDrawColorsU16
						i32.const 1
						local.set $$block_selector
						br $$BlockDisp
					end
					i32.const 22
					local.set $$current_block
					local.get $val
					i32.const 2048
					i32.eq
					local.set $$t10
					local.get $$t10
					if
						i32.const 21
						local.set $$block_selector
						br $$BlockDisp
					else
						i32.const 1
						local.set $$block_selector
						br $$BlockDisp
					end
				end
			end
		end
	)
	(func $$$$$u32$$.underlying.$$onFree (param $$ptr i32)
		local.get $$ptr
		i32.const 4
		i32.add
		i32.const 10
		call_indirect 0 (type $$onFree)
	)
	(func $$i32.$slice.$$block.$$onFree (param $ptr i32)
		local.get $ptr
		i32.load align=1
		call $runtime.Block.Release
		local.get $ptr
		i32.const 0
		i32.store align=1
	)
	(func $$i32.$slice.$slice.underlying.$$onFree (param $$ptr i32)
		local.get $$ptr
		i32.const 31
		call_indirect 0 (type $$onFree)
	)
	(func $$w42048$game.board.$$onFree (param $$ptr i32)
		local.get $$ptr
		i32.const 30
		call_indirect 0 (type $$onFree)
		local.get $$ptr
		i32.const 12
		i32.add
		i32.const 32
		call_indirect 0 (type $$onFree)
	)
	(func $w42048$game.New (param $seed i32) (result i32 i32 i32 i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$ret_0.0.0 i32)
		(local $$ret_0.0.1 i32)
		(local $$ret_0.1 i32)
		(local $$ret_0.2 i32)
		(local $$t0.0 i32)
		(local $$t0.1 i32)
		(local $$t1.0 i32)
		(local $$t1.1 i32)
		(local $$t2.0 i32)
		(local $$t2.1 i32)
		(local $$t2.2 i32)
		(local $$t2.3 i32)
		(local $$t3.0 i32)
		(local $$t3.1 i32)
		(local $$t4.0 i32)
		(local $$t4.1.0 i32)
		(local $$t4.1.1 i32)
		(local $$t5.0.0 i32)
		(local $$t5.0.1 i32)
		(local $$t5.1 i32)
		(local $$t5.2 i32)
		block $$BlockFnBody
			loop $$BlockDisp
				block $$Block_0
					block $$BlockSel
						local.get $$block_selector
						br_table  0 0
					end
					i32.const 0
					local.set $$current_block
					i32.const 56
					call $runtime.HeapAlloc
					i32.const 1
					i32.const 33
					i32.const 40
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
					i32.const 12
					i32.add
					local.set $$t1.1
					local.get $$t1.0
					call $runtime.Block.Release
					local.set $$t1.0
					call $w42048$game.newMatrix
					local.set $$t2.3
					local.set $$t2.2
					local.set $$t2.1
					local.get $$t2.0
					call $runtime.Block.Release
					local.set $$t2.0
					local.get $$t0.0
					call $runtime.Block.Retain
					local.get $$t0.1
					i32.const 0
					i32.add
					local.set $$t3.1
					local.get $$t3.0
					call $runtime.Block.Release
					local.set $$t3.0
					i32.const 1103515245
					i32.const 12345
					i32.const -2147483648
					local.get $seed
					call $w42048$game.lcg
					local.set $$t4.1.1
					local.get $$t4.1.0
					call $runtime.Block.Release
					local.set $$t4.1.0
					local.set $$t4.0
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
					local.get $$t3.1
					local.get $$t4.0
					i32.store
					local.get $$t3.1
					local.get $$t4.1.0
					call $runtime.Block.Retain
					local.get $$t3.1
					i32.load offset=4 align=1
					call $runtime.Block.Release
					i32.store offset=4 align=1
					local.get $$t3.1
					local.get $$t4.1.1
					i32.store offset=8
					local.get $$t0.0
					local.get $$t0.1
					call $w42048$game.board.Add
					local.get $$t0.0
					local.get $$t0.1
					call $w42048$game.board.Add
					local.get $$t0.0
					call $runtime.Block.Retain
					local.get $$t0.1
					i32.const 4
					i32.const -3
					i32.const 0
					call $runtime.getItab
					i32.const 0
					local.set $$t5.2
					local.set $$t5.1
					local.set $$t5.0.1
					local.get $$t5.0.0
					call $runtime.Block.Release
					local.set $$t5.0.0
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
		local.get $$t1.0
		call $runtime.Block.Release
		local.get $$t2.0
		call $runtime.Block.Release
		local.get $$t3.0
		call $runtime.Block.Release
		local.get $$t4.1.0
		call $runtime.Block.Release
		local.get $$t5.0.0
		call $runtime.Block.Release
	)
	(func $w42048$game.init
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
							global.get $w42048$game.init$guard
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
						global.set $w42048$game.init$guard
						br $$Block_1
					end
					i32.const 2
					local.set $$current_block
					br $$BlockFnBody
				end
			end
		end
	)
	(func $w42048$game.lcg$1 (param $r.0 i32) (param $r.1 i32) (param $a.0 i32) (param $a.1 i32) (param $c.0 i32) (param $c.1 i32) (param $m.0 i32) (param $m.1 i32) (result i32)
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
				block $$Block_0
					block $$BlockSel
						local.get $$block_selector
						br_table  0 0
					end
					i32.const 0
					local.set $$current_block
					local.get $a.1
					i32.load
					local.set $$t0
					local.get $r.1
					i32.load
					local.set $$t1
					local.get $$t0
					local.get $$t1
					i32.mul
					local.set $$t2
					local.get $c.1
					i32.load
					local.set $$t3
					local.get $$t2
					local.get $$t3
					i32.add
					local.set $$t4
					local.get $m.1
					i32.load
					local.set $$t5
					local.get $$t4
					local.get $$t5
					i32.rem_u
					local.set $$t6
					local.get $r.1
					local.get $$t6
					i32.store
					local.get $r.1
					i32.load
					local.set $$t7
					local.get $$t7
					local.set $$ret_0
					br $$BlockFnBody
				end
			end
		end
		local.get $$ret_0
	)
	(func $w42048$game.lcg$1.$warpfn (result i32)
		global.get $$wa.runtime.closure_data
		i32.load
		global.get $$wa.runtime.closure_data
		i32.load offset=4
		global.get $$wa.runtime.closure_data
		i32.load offset=8
		global.get $$wa.runtime.closure_data
		i32.load offset=12
		global.get $$wa.runtime.closure_data
		i32.load offset=16
		global.get $$wa.runtime.closure_data
		i32.load offset=20
		global.get $$wa.runtime.closure_data
		i32.load offset=24
		global.get $$wa.runtime.closure_data
		i32.load offset=28
		i32.const 0
		global.set $$wa.runtime.closure_data
		call $w42048$game.lcg$1
	)
	(func $$u32.$$block.$$onFree (param $ptr i32)
		local.get $ptr
		i32.load align=1
		call $runtime.Block.Release
		local.get $ptr
		i32.const 0
		i32.store align=1
	)
	(func $$u32.$ref.underlying.$$onFree (param $$ptr i32)
		local.get $$ptr
		i32.const 35
		call_indirect 0 (type $$onFree)
	)
	(func $$w42048$game.lcg$1.$warpdata.$$onFree (param $$ptr i32)
		local.get $$ptr
		i32.const 36
		call_indirect 0 (type $$onFree)
		local.get $$ptr
		i32.const 8
		i32.add
		i32.const 36
		call_indirect 0 (type $$onFree)
		local.get $$ptr
		i32.const 16
		i32.add
		i32.const 36
		call_indirect 0 (type $$onFree)
		local.get $$ptr
		i32.const 24
		i32.add
		i32.const 36
		call_indirect 0 (type $$onFree)
	)
	(func $w42048$game.lcg (param $a i32) (param $c i32) (param $m i32) (param $seed i32) (result i32 i32 i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$ret_0.0 i32)
		(local $$ret_0.1.0 i32)
		(local $$ret_0.1.1 i32)
		(local $$t0.0 i32)
		(local $$t0.1 i32)
		(local $$t1.0 i32)
		(local $$t1.1 i32)
		(local $$t2.0 i32)
		(local $$t2.1 i32)
		(local $$t3.0 i32)
		(local $$t3.1 i32)
		(local $$t4.0 i32)
		(local $$t4.1.0 i32)
		(local $$t4.1.1 i32)
		(local $$t5.0.0 i32)
		(local $$t5.0.1 i32)
		(local $$t5.1.0 i32)
		(local $$t5.1.1 i32)
		(local $$t5.2.0 i32)
		(local $$t5.2.1 i32)
		(local $$t5.3.0 i32)
		(local $$t5.3.1 i32)
		(local $$t6.0 i32)
		(local $$t6.1.0 i32)
		(local $$t6.1.1 i32)
		block $$BlockFnBody
			loop $$BlockDisp
				block $$Block_0
					block $$BlockSel
						local.get $$block_selector
						br_table  0 0
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
					local.get $$t0.1
					local.get $a
					i32.store
					i32.const 20
					call $runtime.HeapAlloc
					i32.const 1
					i32.const 0
					i32.const 4
					call $runtime.Block.Init
					call $runtime.DupI32
					i32.const 16
					i32.add
					local.set $$t1.1
					local.get $$t1.0
					call $runtime.Block.Release
					local.set $$t1.0
					local.get $$t1.1
					local.get $c
					i32.store
					i32.const 20
					call $runtime.HeapAlloc
					i32.const 1
					i32.const 0
					i32.const 4
					call $runtime.Block.Init
					call $runtime.DupI32
					i32.const 16
					i32.add
					local.set $$t2.1
					local.get $$t2.0
					call $runtime.Block.Release
					local.set $$t2.0
					local.get $$t2.1
					local.get $m
					i32.store
					i32.const 20
					call $runtime.HeapAlloc
					i32.const 1
					i32.const 0
					i32.const 4
					call $runtime.Block.Init
					call $runtime.DupI32
					i32.const 16
					i32.add
					local.set $$t3.1
					local.get $$t3.0
					call $runtime.Block.Release
					local.set $$t3.0
					local.get $$t3.1
					local.get $seed
					i32.store
					local.get $$t3.0
					call $runtime.Block.Retain
					local.get $$t3.1
					local.set $$t5.0.1
					local.get $$t5.0.0
					call $runtime.Block.Release
					local.set $$t5.0.0
					local.get $$t0.0
					call $runtime.Block.Retain
					local.get $$t0.1
					local.set $$t5.1.1
					local.get $$t5.1.0
					call $runtime.Block.Release
					local.set $$t5.1.0
					local.get $$t1.0
					call $runtime.Block.Retain
					local.get $$t1.1
					local.set $$t5.2.1
					local.get $$t5.2.0
					call $runtime.Block.Release
					local.set $$t5.2.0
					local.get $$t2.0
					call $runtime.Block.Retain
					local.get $$t2.1
					local.set $$t5.3.1
					local.get $$t5.3.0
					call $runtime.Block.Release
					local.set $$t5.3.0
					i32.const 34
					local.set $$t4.0
					i32.const 48
					call $runtime.HeapAlloc
					i32.const 1
					i32.const 37
					i32.const 32
					call $runtime.Block.Init
					call $runtime.DupI32
					i32.const 16
					i32.add
					local.set $$t4.1.1
					local.get $$t4.1.0
					call $runtime.Block.Release
					local.set $$t4.1.0
					local.get $$t4.1.1
					local.get $$t5.0.0
					call $runtime.Block.Retain
					local.get $$t4.1.1
					i32.load align=1
					call $runtime.Block.Release
					i32.store align=1
					local.get $$t4.1.1
					local.get $$t5.0.1
					i32.store offset=4
					local.get $$t4.1.1
					local.get $$t5.1.0
					call $runtime.Block.Retain
					local.get $$t4.1.1
					i32.load offset=8 align=1
					call $runtime.Block.Release
					i32.store offset=8 align=1
					local.get $$t4.1.1
					local.get $$t5.1.1
					i32.store offset=12
					local.get $$t4.1.1
					local.get $$t5.2.0
					call $runtime.Block.Retain
					local.get $$t4.1.1
					i32.load offset=16 align=1
					call $runtime.Block.Release
					i32.store offset=16 align=1
					local.get $$t4.1.1
					local.get $$t5.2.1
					i32.store offset=20
					local.get $$t4.1.1
					local.get $$t5.3.0
					call $runtime.Block.Retain
					local.get $$t4.1.1
					i32.load offset=24 align=1
					call $runtime.Block.Release
					i32.store offset=24 align=1
					local.get $$t4.1.1
					local.get $$t5.3.1
					i32.store offset=28
					local.get $$t5.3.0
					call $runtime.Block.Release
					local.get $$t5.2.0
					call $runtime.Block.Release
					local.get $$t5.1.0
					call $runtime.Block.Release
					local.get $$t5.0.0
					call $runtime.Block.Release
					i32.const 0
					local.set $$t5.0.0
					i32.const 0
					local.set $$t5.0.1
					i32.const 0
					local.set $$t5.1.0
					i32.const 0
					local.set $$t5.1.1
					i32.const 0
					local.set $$t5.2.0
					i32.const 0
					local.set $$t5.2.1
					i32.const 0
					local.set $$t5.3.0
					i32.const 0
					local.set $$t5.3.1
					local.get $$t4.0
					local.get $$t4.1.0
					call $runtime.Block.Retain
					local.get $$t4.1.1
					local.set $$t6.1.1
					local.get $$t6.1.0
					call $runtime.Block.Release
					local.set $$t6.1.0
					local.set $$t6.0
					local.get $$t6.0
					local.get $$t6.1.0
					call $runtime.Block.Retain
					local.get $$t6.1.1
					local.set $$ret_0.1.1
					local.get $$ret_0.1.0
					call $runtime.Block.Release
					local.set $$ret_0.1.0
					local.set $$ret_0.0
					br $$BlockFnBody
				end
			end
		end
		local.get $$ret_0.0
		local.get $$ret_0.1.0
		call $runtime.Block.Retain
		local.get $$ret_0.1.1
		local.get $$ret_0.1.0
		call $runtime.Block.Release
		local.get $$t0.0
		call $runtime.Block.Release
		local.get $$t1.0
		call $runtime.Block.Release
		local.get $$t2.0
		call $runtime.Block.Release
		local.get $$t3.0
		call $runtime.Block.Release
		local.get $$t4.1.0
		call $runtime.Block.Release
		local.get $$t5.3.0
		call $runtime.Block.Release
		local.get $$t5.2.0
		call $runtime.Block.Release
		local.get $$t5.1.0
		call $runtime.Block.Release
		local.get $$t5.0.0
		call $runtime.Block.Release
		local.get $$t6.1.0
		call $runtime.Block.Release
	)
	(func $w42048$game.mergeElements (param $old.0 i32) (param $old.1 i32) (param $old.2 i32) (param $old.3 i32) (result i32 i32 i32 i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$ret_0.0 i32)
		(local $$ret_0.1 i32)
		(local $$ret_0.2 i32)
		(local $$ret_0.3 i32)
		(local $$t0 i32)
		(local $$t1.0 i32)
		(local $$t1.1 i32)
		(local $$t1.2 i32)
		(local $$t1.3 i32)
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
		(local $$t10 i32)
		(local $$t11 i32)
		(local $$t12 i32)
		(local $$t13 i32)
		(local $$t14 i32)
		(local $$t15 i32)
		(local $$t16.0 i32)
		(local $$t16.1 i32)
		(local $$t17.0 i32)
		(local $$t17.1 i32)
		(local $$t18 i32)
		(local $$t19 i32)
		(local $$t20 i32)
		(local $$t21 i32)
		(local $$t22.0 i32)
		(local $$t22.1 i32)
		(local $$t23.0 i32)
		(local $$t23.1 i32)
		(local $$t24 i32)
		block $$BlockFnBody
			loop $$BlockDisp
				block $$Block_6
					block $$Block_5
						block $$Block_4
							block $$Block_3
								block $$Block_2
									block $$Block_1
										block $$Block_0
											block $$BlockSel
												local.get $$block_selector
												br_table  0 1 2 3 4 5 6 0
											end
											i32.const 0
											local.set $$current_block
											local.get $old.2
											local.set $$t0
											local.get $$t0
											i32.const 4
											i32.mul
											i32.const 16
											i32.add
											call $runtime.HeapAlloc
											local.get $$t0
											i32.const 0
											i32.const 4
											call $runtime.Block.Init
											call $runtime.DupI32
											i32.const 16
											i32.add
											local.get $$t0
											local.get $$t0
											local.set $$t1.3
											local.set $$t1.2
											local.set $$t1.1
											local.get $$t1.0
											call $runtime.Block.Release
											local.set $$t1.0
											local.get $$t1.0
											call $runtime.Block.Retain
											local.get $$t1.1
											i32.const 4
											i32.const 0
											i32.mul
											i32.add
											local.set $$t2.1
											local.get $$t2.0
											call $runtime.Block.Release
											local.set $$t2.0
											local.get $old.0
											call $runtime.Block.Retain
											local.get $old.1
											i32.const 4
											i32.const 0
											i32.mul
											i32.add
											local.set $$t3.1
											local.get $$t3.0
											call $runtime.Block.Release
											local.set $$t3.0
											local.get $$t3.1
											i32.load
											local.set $$t4
											local.get $$t2.1
											local.get $$t4
											i32.store
											br $$Block_2
										end
										i32.const 1
										local.set $$current_block
										local.get $old.0
										call $runtime.Block.Retain
										local.get $old.1
										i32.const 4
										local.get $$t5
										i32.mul
										i32.add
										local.set $$t6.1
										local.get $$t6.0
										call $runtime.Block.Release
										local.set $$t6.0
										local.get $$t6.1
										i32.load
										local.set $$t7
										local.get $$t1.0
										call $runtime.Block.Retain
										local.get $$t1.1
										i32.const 4
										local.get $$t8
										i32.mul
										i32.add
										local.set $$t9.1
										local.get $$t9.0
										call $runtime.Block.Release
										local.set $$t9.0
										local.get $$t9.1
										i32.load
										local.set $$t10
										local.get $$t7
										local.get $$t10
										i32.eq
										local.set $$t11
										local.get $$t11
										if
											br $$Block_3
										else
											br $$Block_5
										end
									end
									i32.const 2
									local.set $$current_block
									local.get $$t1.0
									call $runtime.Block.Retain
									local.get $$t1.1
									local.get $$t1.2
									local.get $$t1.3
									local.set $$ret_0.3
									local.set $$ret_0.2
									local.set $$ret_0.1
									local.get $$ret_0.0
									call $runtime.Block.Release
									local.set $$ret_0.0
									br $$BlockFnBody
								end
								local.get $$current_block
								i32.const 0
								i32.eq
								if(result i32)
									i32.const 0
								else
									local.get $$t12
								end
								local.get $$current_block
								i32.const 0
								i32.eq
								if(result i32)
									i32.const 1
								else
									local.get $$t13
								end
								local.set $$t5
								local.set $$t8
								i32.const 3
								local.set $$current_block
								local.get $old.2
								local.set $$t14
								local.get $$t5
								local.get $$t14
								i32.lt_s
								local.set $$t15
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
							end
							i32.const 4
							local.set $$current_block
							local.get $$t1.0
							call $runtime.Block.Retain
							local.get $$t1.1
							i32.const 4
							local.get $$t8
							i32.mul
							i32.add
							local.set $$t16.1
							local.get $$t16.0
							call $runtime.Block.Release
							local.set $$t16.0
							local.get $old.0
							call $runtime.Block.Retain
							local.get $old.1
							i32.const 4
							local.get $$t5
							i32.mul
							i32.add
							local.set $$t17.1
							local.get $$t17.0
							call $runtime.Block.Release
							local.set $$t17.0
							local.get $$t17.1
							i32.load
							local.set $$t18
							local.get $$t16.1
							i32.load
							local.set $$t19
							local.get $$t19
							local.get $$t18
							i32.add
							local.set $$t20
							local.get $$t16.1
							local.get $$t20
							i32.store
							br $$Block_4
						end
						local.get $$current_block
						i32.const 4
						i32.eq
						if(result i32)
							local.get $$t8
						else
							local.get $$t21
						end
						local.set $$t12
						i32.const 5
						local.set $$current_block
						local.get $$t5
						i32.const 1
						i32.add
						local.set $$t13
						i32.const 3
						local.set $$block_selector
						br $$BlockDisp
					end
					i32.const 6
					local.set $$current_block
					local.get $$t8
					i32.const 1
					i32.add
					local.set $$t21
					local.get $$t1.0
					call $runtime.Block.Retain
					local.get $$t1.1
					i32.const 4
					local.get $$t21
					i32.mul
					i32.add
					local.set $$t22.1
					local.get $$t22.0
					call $runtime.Block.Release
					local.set $$t22.0
					local.get $old.0
					call $runtime.Block.Retain
					local.get $old.1
					i32.const 4
					local.get $$t5
					i32.mul
					i32.add
					local.set $$t23.1
					local.get $$t23.0
					call $runtime.Block.Release
					local.set $$t23.0
					local.get $$t23.1
					i32.load
					local.set $$t24
					local.get $$t22.1
					local.get $$t24
					i32.store
					i32.const 5
					local.set $$block_selector
					br $$BlockDisp
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
		local.get $$t1.0
		call $runtime.Block.Release
		local.get $$t2.0
		call $runtime.Block.Release
		local.get $$t3.0
		call $runtime.Block.Release
		local.get $$t6.0
		call $runtime.Block.Release
		local.get $$t9.0
		call $runtime.Block.Release
		local.get $$t16.0
		call $runtime.Block.Release
		local.get $$t17.0
		call $runtime.Block.Release
		local.get $$t22.0
		call $runtime.Block.Release
		local.get $$t23.0
		call $runtime.Block.Release
	)
	(func $$i32.$slice.append (param $x.0 i32) (param $x.1 i32) (param $x.2 i32) (param $x.3 i32) (param $y.0 i32) (param $y.1 i32) (param $y.2 i32) (param $y.3 i32) (result i32 i32 i32 i32)
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
		if(result i32 i32 i32 i32)
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
					end
					local.get $src
					i32.load
					local.set $item
					local.get $dest
					local.get $item
					i32.store
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
				end
			end
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
					end
					local.get $src
					i32.load
					local.set $item
					local.get $dest
					local.get $item
					i32.store
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
					i32.load
					local.set $item
					local.get $dest
					local.get $item
					i32.store
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
				end
			end
		end
	)
	(func $w42048$game.movedRow (param $elems.0 i32) (param $elems.1 i32) (param $elems.2 i32) (param $elems.3 i32) (result i32 i32 i32 i32)
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
		(local $$t2 i32)
		(local $$t3.0 i32)
		(local $$t3.1 i32)
		(local $$t4 i32)
		(local $$t5 i32)
		(local $$t6.0 i32)
		(local $$t6.1 i32)
		(local $$t6.2 i32)
		(local $$t6.3 i32)
		(local $$t7 i32)
		(local $$t8 i32)
		(local $$t9.0 i32)
		(local $$t9.1 i32)
		(local $$t9.2 i32)
		(local $$t9.3 i32)
		(local $$t10 i32)
		(local $$t11 i32)
		(local $$t12.0 i32)
		(local $$t12.1 i32)
		(local $$t13 i32)
		(local $$t14.0 i32)
		(local $$t14.1 i32)
		(local $$t15.0 i32)
		(local $$t15.1 i32)
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
		(local $$t19.0 i32)
		(local $$t19.1 i32)
		(local $$t20.0 i32)
		(local $$t20.1 i32)
		(local $$t20.2 i32)
		(local $$t20.3 i32)
		(local $$t21.0 i32)
		(local $$t21.1 i32)
		(local $$t21.2 i32)
		(local $$t21.3 i32)
		(local $$t22.0 i32)
		(local $$t22.1 i32)
		(local $$t22.2 i32)
		(local $$t22.3 i32)
		(local $$t23 i32)
		(local $$t24 i32)
		(local $$t25.0 i32)
		(local $$t25.1 i32)
		(local $$t25.2 i32)
		(local $$t25.3 i32)
		(local $$t26 i32)
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
													i32.const 16
													call $runtime.HeapAlloc
													i32.const 1
													i32.const 0
													i32.const 0
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
													i32.const 4
													i32.const 0
													i32.mul
													i32.add
													i32.const 0
													i32.const 0
													i32.sub
													i32.const 0
													i32.const 0
													i32.sub
													local.set $$t1.3
													local.set $$t1.2
													local.set $$t1.1
													local.get $$t1.0
													call $runtime.Block.Release
													local.set $$t1.0
													br $$Block_2
												end
												i32.const 1
												local.set $$current_block
												local.get $elems.0
												call $runtime.Block.Retain
												local.get $elems.1
												i32.const 4
												local.get $$t2
												i32.mul
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
												i32.eqz
												local.set $$t5
												local.get $$t5
												if
													br $$Block_3
												else
													br $$Block_4
												end
											end
											i32.const 2
											local.set $$current_block
											local.get $$t6.2
											local.set $$t7
											i32.const 4
											local.get $$t7
											i32.sub
											local.set $$t8
											br $$Block_7
										end
										local.get $$current_block
										i32.const 0
										i32.eq
										if(result i32 i32 i32 i32)
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
										local.get $$current_block
										i32.const 0
										i32.eq
										if(result i32)
											i32.const 0
										else
											local.get $$t10
										end
										local.set $$t2
										local.set $$t6.3
										local.set $$t6.2
										local.set $$t6.1
										local.get $$t6.0
										call $runtime.Block.Release
										local.set $$t6.0
										i32.const 3
										local.set $$current_block
										local.get $$t2
										i32.const 4
										i32.lt_s
										local.set $$t11
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
									end
									i32.const 4
									local.set $$current_block
									local.get $elems.0
									call $runtime.Block.Retain
									local.get $elems.1
									i32.const 4
									local.get $$t2
									i32.mul
									i32.add
									local.set $$t12.1
									local.get $$t12.0
									call $runtime.Block.Release
									local.set $$t12.0
									local.get $$t12.1
									i32.load
									local.set $$t13
									i32.const 20
									call $runtime.HeapAlloc
									i32.const 1
									i32.const 0
									i32.const 4
									call $runtime.Block.Init
									call $runtime.DupI32
									i32.const 16
									i32.add
									local.set $$t14.1
									local.get $$t14.0
									call $runtime.Block.Release
									local.set $$t14.0
									local.get $$t14.0
									call $runtime.Block.Retain
									local.get $$t14.1
									i32.const 4
									i32.const 0
									i32.mul
									i32.add
									local.set $$t15.1
									local.get $$t15.0
									call $runtime.Block.Release
									local.set $$t15.0
									local.get $$t15.1
									local.get $$t13
									i32.store
									local.get $$t14.0
									call $runtime.Block.Retain
									local.get $$t14.1
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
									local.set $$t16.3
									local.set $$t16.2
									local.set $$t16.1
									local.get $$t16.0
									call $runtime.Block.Release
									local.set $$t16.0
									local.get $$t6.0
									local.get $$t6.1
									local.get $$t6.2
									local.get $$t6.3
									local.get $$t16.0
									local.get $$t16.1
									local.get $$t16.2
									local.get $$t16.3
									call $$i32.$slice.append
									local.set $$t17.3
									local.set $$t17.2
									local.set $$t17.1
									local.get $$t17.0
									call $runtime.Block.Release
									local.set $$t17.0
									br $$Block_4
								end
								local.get $$current_block
								i32.const 1
								i32.eq
								if(result i32 i32 i32 i32)
									local.get $$t6.0
									call $runtime.Block.Retain
									local.get $$t6.1
									local.get $$t6.2
									local.get $$t6.3
								else
									local.get $$t17.0
									call $runtime.Block.Retain
									local.get $$t17.1
									local.get $$t17.2
									local.get $$t17.3
								end
								local.set $$t9.3
								local.set $$t9.2
								local.set $$t9.1
								local.get $$t9.0
								call $runtime.Block.Release
								local.set $$t9.0
								i32.const 5
								local.set $$current_block
								local.get $$t2
								i32.const 1
								i32.add
								local.set $$t10
								i32.const 3
								local.set $$block_selector
								br $$BlockDisp
							end
							i32.const 6
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
							local.set $$t18.1
							local.get $$t18.0
							call $runtime.Block.Release
							local.set $$t18.0
							local.get $$t18.0
							call $runtime.Block.Retain
							local.get $$t18.1
							i32.const 4
							i32.const 0
							i32.mul
							i32.add
							local.set $$t19.1
							local.get $$t19.0
							call $runtime.Block.Release
							local.set $$t19.0
							local.get $$t19.1
							i32.const 0
							i32.store
							local.get $$t18.0
							call $runtime.Block.Retain
							local.get $$t18.1
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
							local.set $$t20.3
							local.set $$t20.2
							local.set $$t20.1
							local.get $$t20.0
							call $runtime.Block.Release
							local.set $$t20.0
							local.get $$t21.0
							local.get $$t21.1
							local.get $$t21.2
							local.get $$t21.3
							local.get $$t20.0
							local.get $$t20.1
							local.get $$t20.2
							local.get $$t20.3
							call $$i32.$slice.append
							local.set $$t22.3
							local.set $$t22.2
							local.set $$t22.1
							local.get $$t22.0
							call $runtime.Block.Release
							local.set $$t22.0
							local.get $$t23
							i32.const 1
							i32.add
							local.set $$t24
							br $$Block_7
						end
						i32.const 7
						local.set $$current_block
						local.get $$t21.0
						local.get $$t21.1
						local.get $$t21.2
						local.get $$t21.3
						call $w42048$game.mergeElements
						local.set $$t25.3
						local.set $$t25.2
						local.set $$t25.1
						local.get $$t25.0
						call $runtime.Block.Release
						local.set $$t25.0
						local.get $$t25.0
						call $runtime.Block.Retain
						local.get $$t25.1
						local.get $$t25.2
						local.get $$t25.3
						local.set $$ret_0.3
						local.set $$ret_0.2
						local.set $$ret_0.1
						local.get $$ret_0.0
						call $runtime.Block.Release
						local.set $$ret_0.0
						br $$BlockFnBody
					end
					local.get $$current_block
					i32.const 2
					i32.eq
					if(result i32 i32 i32 i32)
						local.get $$t6.0
						call $runtime.Block.Retain
						local.get $$t6.1
						local.get $$t6.2
						local.get $$t6.3
					else
						local.get $$t22.0
						call $runtime.Block.Retain
						local.get $$t22.1
						local.get $$t22.2
						local.get $$t22.3
					end
					local.get $$current_block
					i32.const 2
					i32.eq
					if(result i32)
						i32.const 0
					else
						local.get $$t24
					end
					local.set $$t23
					local.set $$t21.3
					local.set $$t21.2
					local.set $$t21.1
					local.get $$t21.0
					call $runtime.Block.Release
					local.set $$t21.0
					i32.const 8
					local.set $$current_block
					local.get $$t23
					local.get $$t8
					i32.lt_s
					local.set $$t26
					local.get $$t26
					if
						i32.const 6
						local.set $$block_selector
						br $$BlockDisp
					else
						i32.const 7
						local.set $$block_selector
						br $$BlockDisp
					end
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
		local.get $$t1.0
		call $runtime.Block.Release
		local.get $$t3.0
		call $runtime.Block.Release
		local.get $$t6.0
		call $runtime.Block.Release
		local.get $$t9.0
		call $runtime.Block.Release
		local.get $$t12.0
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
		local.get $$t25.0
		call $runtime.Block.Release
	)
	(func $$i32.$$block.$$onFree (param $ptr i32)
		local.get $ptr
		i32.load align=1
		call $runtime.Block.Release
		local.get $ptr
		i32.const 0
		i32.store align=1
	)
	(func $$i32.$slice.underlying.$$onFree (param $$ptr i32)
		local.get $$ptr
		i32.const 38
		call_indirect 0 (type $$onFree)
	)
	(func $$i32.$slice.$array1.underlying.$$onFree (param $$ptr i32)
		local.get $$ptr
		i32.const 39
		call_indirect 0 (type $$onFree)
	)
	(func $$i32.$slice.$slice.append (param $x.0 i32) (param $x.1 i32) (param $x.2 i32) (param $x.3 i32) (param $y.0 i32) (param $y.1 i32) (param $y.2 i32) (param $y.3 i32) (result i32 i32 i32 i32)
		(local $item.0 i32)
		(local $item.1 i32)
		(local $item.2 i32)
		(local $item.3 i32)
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
					end
					local.get $src
					i32.load
					call $runtime.Block.Retain
					local.get $src
					i32.load offset=4
					local.get $src
					i32.load offset=8
					local.get $src
					i32.load offset=12
					local.set $item.3
					local.set $item.2
					local.set $item.1
					local.get $item.0
					call $runtime.Block.Release
					local.set $item.0
					local.get $dest
					local.get $item.0
					call $runtime.Block.Retain
					local.get $dest
					i32.load align=1
					call $runtime.Block.Release
					i32.store align=1
					local.get $dest
					local.get $item.1
					i32.store offset=4
					local.get $dest
					local.get $item.2
					i32.store offset=8
					local.get $dest
					local.get $item.3
					i32.store offset=12
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
				end
			end
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
			i32.const 39
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
					end
					local.get $src
					i32.load
					call $runtime.Block.Retain
					local.get $src
					i32.load offset=4
					local.get $src
					i32.load offset=8
					local.get $src
					i32.load offset=12
					local.set $item.3
					local.set $item.2
					local.set $item.1
					local.get $item.0
					call $runtime.Block.Release
					local.set $item.0
					local.get $dest
					local.get $item.0
					call $runtime.Block.Retain
					local.get $dest
					i32.load align=1
					call $runtime.Block.Release
					i32.store align=1
					local.get $dest
					local.get $item.1
					i32.store offset=4
					local.get $dest
					local.get $item.2
					i32.store offset=8
					local.get $dest
					local.get $item.3
					i32.store offset=12
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
					i32.load
					call $runtime.Block.Retain
					local.get $src
					i32.load offset=4
					local.get $src
					i32.load offset=8
					local.get $src
					i32.load offset=12
					local.set $item.3
					local.set $item.2
					local.set $item.1
					local.get $item.0
					call $runtime.Block.Release
					local.set $item.0
					local.get $dest
					local.get $item.0
					call $runtime.Block.Retain
					local.get $dest
					i32.load align=1
					call $runtime.Block.Release
					i32.store align=1
					local.get $dest
					local.get $item.1
					i32.store offset=4
					local.get $dest
					local.get $item.2
					i32.store offset=8
					local.get $dest
					local.get $item.3
					i32.store offset=12
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
				end
			end
		end
		local.get $item.0
		call $runtime.Block.Release
	)
	(func $w42048$game.newMatrix (result i32 i32 i32 i32)
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
		(local $$t8.0 i32)
		(local $$t8.1 i32)
		(local $$t8.2 i32)
		(local $$t8.3 i32)
		(local $$t9 i32)
		(local $$t10 i32)
		(local $$t11 i32)
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
								i32.const 16
								call $runtime.HeapAlloc
								i32.const 1
								i32.const 0
								i32.const 0
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
								i32.const 16
								i32.const 0
								i32.mul
								i32.add
								i32.const 0
								i32.const 0
								i32.sub
								i32.const 0
								i32.const 0
								i32.sub
								local.set $$t1.3
								local.set $$t1.2
								local.set $$t1.1
								local.get $$t1.0
								call $runtime.Block.Release
								local.set $$t1.0
								br $$Block_2
							end
							i32.const 1
							local.set $$current_block
							i32.const 32
							call $runtime.HeapAlloc
							i32.const 1
							i32.const 0
							i32.const 16
							call $runtime.Block.Init
							call $runtime.DupI32
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
							i32.const 0
							i32.mul
							i32.add
							i32.const 4
							i32.const 0
							i32.sub
							i32.const 4
							i32.const 0
							i32.sub
							local.set $$t3.3
							local.set $$t3.2
							local.set $$t3.1
							local.get $$t3.0
							call $runtime.Block.Release
							local.set $$t3.0
							i32.const 32
							call $runtime.HeapAlloc
							i32.const 1
							i32.const 40
							i32.const 16
							call $runtime.Block.Init
							call $runtime.DupI32
							i32.const 16
							i32.add
							local.set $$t4.1
							local.get $$t4.0
							call $runtime.Block.Release
							local.set $$t4.0
							local.get $$t4.0
							call $runtime.Block.Retain
							local.get $$t4.1
							i32.const 16
							i32.const 0
							i32.mul
							i32.add
							local.set $$t5.1
							local.get $$t5.0
							call $runtime.Block.Release
							local.set $$t5.0
							local.get $$t5.1
							local.get $$t3.0
							call $runtime.Block.Retain
							local.get $$t5.1
							i32.load align=1
							call $runtime.Block.Release
							i32.store align=1
							local.get $$t5.1
							local.get $$t3.1
							i32.store offset=4
							local.get $$t5.1
							local.get $$t3.2
							i32.store offset=8
							local.get $$t5.1
							local.get $$t3.3
							i32.store offset=12
							local.get $$t4.0
							call $runtime.Block.Retain
							local.get $$t4.1
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
							local.set $$t6.3
							local.set $$t6.2
							local.set $$t6.1
							local.get $$t6.0
							call $runtime.Block.Release
							local.set $$t6.0
							local.get $$t7.0
							local.get $$t7.1
							local.get $$t7.2
							local.get $$t7.3
							local.get $$t6.0
							local.get $$t6.1
							local.get $$t6.2
							local.get $$t6.3
							call $$i32.$slice.$slice.append
							local.set $$t8.3
							local.set $$t8.2
							local.set $$t8.1
							local.get $$t8.0
							call $runtime.Block.Release
							local.set $$t8.0
							local.get $$t9
							i32.const 1
							i32.add
							local.set $$t10
							br $$Block_2
						end
						i32.const 2
						local.set $$current_block
						local.get $$t7.0
						call $runtime.Block.Retain
						local.get $$t7.1
						local.get $$t7.2
						local.get $$t7.3
						local.set $$ret_0.3
						local.set $$ret_0.2
						local.set $$ret_0.1
						local.get $$ret_0.0
						call $runtime.Block.Release
						local.set $$ret_0.0
						br $$BlockFnBody
					end
					local.get $$current_block
					i32.const 0
					i32.eq
					if(result i32 i32 i32 i32)
						local.get $$t1.0
						call $runtime.Block.Retain
						local.get $$t1.1
						local.get $$t1.2
						local.get $$t1.3
					else
						local.get $$t8.0
						call $runtime.Block.Retain
						local.get $$t8.1
						local.get $$t8.2
						local.get $$t8.3
					end
					local.get $$current_block
					i32.const 0
					i32.eq
					if(result i32)
						i32.const 0
					else
						local.get $$t10
					end
					local.set $$t9
					local.set $$t7.3
					local.set $$t7.2
					local.set $$t7.1
					local.get $$t7.0
					call $runtime.Block.Release
					local.set $$t7.0
					i32.const 3
					local.set $$current_block
					local.get $$t9
					i32.const 4
					i32.lt_s
					local.set $$t11
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
	)
	(func $w42048$game.reverseRow (param $arr.0 i32) (param $arr.1 i32) (param $arr.2 i32) (param $arr.3 i32) (result i32 i32 i32 i32)
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
		(local $$t2 i32)
		(local $$t3 i32)
		(local $$t4 i32)
		(local $$t5.0 i32)
		(local $$t5.1 i32)
		(local $$t6 i32)
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
		(local $$t10.2 i32)
		(local $$t10.3 i32)
		(local $$t11.0 i32)
		(local $$t11.1 i32)
		(local $$t11.2 i32)
		(local $$t11.3 i32)
		(local $$t12 i32)
		(local $$t13 i32)
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
								i32.const 16
								call $runtime.HeapAlloc
								i32.const 1
								i32.const 0
								i32.const 0
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
								i32.const 4
								i32.const 0
								i32.mul
								i32.add
								i32.const 0
								i32.const 0
								i32.sub
								i32.const 0
								i32.const 0
								i32.sub
								local.set $$t1.3
								local.set $$t1.2
								local.set $$t1.1
								local.get $$t1.0
								call $runtime.Block.Release
								local.set $$t1.0
								local.get $arr.2
								local.set $$t2
								local.get $$t2
								i32.const 1
								i32.sub
								local.set $$t3
								br $$Block_2
							end
							i32.const 1
							local.set $$current_block
							local.get $arr.0
							call $runtime.Block.Retain
							local.get $arr.1
							i32.const 4
							local.get $$t4
							i32.mul
							i32.add
							local.set $$t5.1
							local.get $$t5.0
							call $runtime.Block.Release
							local.set $$t5.0
							local.get $$t5.1
							i32.load
							local.set $$t6
							i32.const 20
							call $runtime.HeapAlloc
							i32.const 1
							i32.const 0
							i32.const 4
							call $runtime.Block.Init
							call $runtime.DupI32
							i32.const 16
							i32.add
							local.set $$t7.1
							local.get $$t7.0
							call $runtime.Block.Release
							local.set $$t7.0
							local.get $$t7.0
							call $runtime.Block.Retain
							local.get $$t7.1
							i32.const 4
							i32.const 0
							i32.mul
							i32.add
							local.set $$t8.1
							local.get $$t8.0
							call $runtime.Block.Release
							local.set $$t8.0
							local.get $$t8.1
							local.get $$t6
							i32.store
							local.get $$t7.0
							call $runtime.Block.Retain
							local.get $$t7.1
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
							local.set $$t9.3
							local.set $$t9.2
							local.set $$t9.1
							local.get $$t9.0
							call $runtime.Block.Release
							local.set $$t9.0
							local.get $$t10.0
							local.get $$t10.1
							local.get $$t10.2
							local.get $$t10.3
							local.get $$t9.0
							local.get $$t9.1
							local.get $$t9.2
							local.get $$t9.3
							call $$i32.$slice.append
							local.set $$t11.3
							local.set $$t11.2
							local.set $$t11.1
							local.get $$t11.0
							call $runtime.Block.Release
							local.set $$t11.0
							local.get $$t4
							i32.const 1
							i32.sub
							local.set $$t12
							br $$Block_2
						end
						i32.const 2
						local.set $$current_block
						local.get $$t10.0
						call $runtime.Block.Retain
						local.get $$t10.1
						local.get $$t10.2
						local.get $$t10.3
						local.set $$ret_0.3
						local.set $$ret_0.2
						local.set $$ret_0.1
						local.get $$ret_0.0
						call $runtime.Block.Release
						local.set $$ret_0.0
						br $$BlockFnBody
					end
					local.get $$current_block
					i32.const 0
					i32.eq
					if(result i32 i32 i32 i32)
						local.get $$t1.0
						call $runtime.Block.Retain
						local.get $$t1.1
						local.get $$t1.2
						local.get $$t1.3
					else
						local.get $$t11.0
						call $runtime.Block.Retain
						local.get $$t11.1
						local.get $$t11.2
						local.get $$t11.3
					end
					local.get $$current_block
					i32.const 0
					i32.eq
					if(result i32)
						local.get $$t3
					else
						local.get $$t12
					end
					local.set $$t4
					local.set $$t10.3
					local.set $$t10.2
					local.set $$t10.1
					local.get $$t10.0
					call $runtime.Block.Release
					local.set $$t10.0
					i32.const 3
					local.set $$current_block
					local.get $$t4
					i32.const 0
					i32.ge_s
					local.set $$t13
					local.get $$t13
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
		local.get $$t1.0
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
	)
	(func $w42048$palettes.init
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$t0 i32)
		(local $$t1.0 i32)
		(local $$t1.1 i32)
		(local $$t2.0 i32)
		(local $$t2.1 i32)
		(local $$t3.0 i32)
		(local $$t3.1 i32)
		(local $$t3.2 i32)
		(local $$t3.3 i32)
		(local $$t4.0 i32)
		(local $$t4.1 i32)
		(local $$t5.0 i32)
		(local $$t5.1 i32)
		(local $$t5.2 i32)
		(local $$t5.3 i32)
		(local $$t6.0 i32)
		(local $$t6.1 i32)
		(local $$t7.0 i32)
		(local $$t7.1 i32)
		(local $$t7.2 i32)
		(local $$t7.3 i32)
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
		(local $$t11.2 i32)
		(local $$t11.3 i32)
		(local $$t12.0 i32)
		(local $$t12.1 i32)
		(local $$t13.0 i32)
		(local $$t13.1 i32)
		(local $$t13.2 i32)
		(local $$t13.3 i32)
		(local $$t14.0 i32)
		(local $$t14.1 i32)
		(local $$t15.0 i32)
		(local $$t15.1 i32)
		(local $$t15.2 i32)
		(local $$t15.3 i32)
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
		(local $$t19.2 i32)
		(local $$t19.3 i32)
		(local $$t20.0 i32)
		(local $$t20.1 i32)
		(local $$t21.0 i32)
		(local $$t21.1 i32)
		(local $$t21.2 i32)
		(local $$t21.3 i32)
		(local $$t22.0 i32)
		(local $$t22.1 i32)
		(local $$t22.2 i32)
		(local $$t22.3 i32)
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
							global.get $w42048$palettes.init$guard
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
						global.set $w42048$palettes.init$guard
						i32.const 176
						call $runtime.HeapAlloc
						i32.const 1
						i32.const 0
						i32.const 160
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
						i32.const 16
						i32.const 0
						i32.mul
						i32.add
						local.set $$t2.1
						local.get $$t2.0
						call $runtime.Block.Release
						local.set $$t2.0
						i32.const 31288
						i32.load
						i32.const 31288
						i32.load offset=4
						i32.const 31288
						i32.load offset=8
						i32.const 31288
						i32.load offset=12
						local.set $$t3.3
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
						local.get $$t2.1
						local.get $$t3.3
						i32.store offset=12
						local.get $$t1.0
						call $runtime.Block.Retain
						local.get $$t1.1
						i32.const 16
						i32.const 1
						i32.mul
						i32.add
						local.set $$t4.1
						local.get $$t4.0
						call $runtime.Block.Release
						local.set $$t4.0
						i32.const 31304
						i32.load
						i32.const 31304
						i32.load offset=4
						i32.const 31304
						i32.load offset=8
						i32.const 31304
						i32.load offset=12
						local.set $$t5.3
						local.set $$t5.2
						local.set $$t5.1
						local.set $$t5.0
						local.get $$t4.1
						local.get $$t5.0
						i32.store
						local.get $$t4.1
						local.get $$t5.1
						i32.store offset=4
						local.get $$t4.1
						local.get $$t5.2
						i32.store offset=8
						local.get $$t4.1
						local.get $$t5.3
						i32.store offset=12
						local.get $$t1.0
						call $runtime.Block.Retain
						local.get $$t1.1
						i32.const 16
						i32.const 2
						i32.mul
						i32.add
						local.set $$t6.1
						local.get $$t6.0
						call $runtime.Block.Release
						local.set $$t6.0
						i32.const 31320
						i32.load
						i32.const 31320
						i32.load offset=4
						i32.const 31320
						i32.load offset=8
						i32.const 31320
						i32.load offset=12
						local.set $$t7.3
						local.set $$t7.2
						local.set $$t7.1
						local.set $$t7.0
						local.get $$t6.1
						local.get $$t7.0
						i32.store
						local.get $$t6.1
						local.get $$t7.1
						i32.store offset=4
						local.get $$t6.1
						local.get $$t7.2
						i32.store offset=8
						local.get $$t6.1
						local.get $$t7.3
						i32.store offset=12
						local.get $$t1.0
						call $runtime.Block.Retain
						local.get $$t1.1
						i32.const 16
						i32.const 3
						i32.mul
						i32.add
						local.set $$t8.1
						local.get $$t8.0
						call $runtime.Block.Release
						local.set $$t8.0
						i32.const 31336
						i32.load
						i32.const 31336
						i32.load offset=4
						i32.const 31336
						i32.load offset=8
						i32.const 31336
						i32.load offset=12
						local.set $$t9.3
						local.set $$t9.2
						local.set $$t9.1
						local.set $$t9.0
						local.get $$t8.1
						local.get $$t9.0
						i32.store
						local.get $$t8.1
						local.get $$t9.1
						i32.store offset=4
						local.get $$t8.1
						local.get $$t9.2
						i32.store offset=8
						local.get $$t8.1
						local.get $$t9.3
						i32.store offset=12
						local.get $$t1.0
						call $runtime.Block.Retain
						local.get $$t1.1
						i32.const 16
						i32.const 4
						i32.mul
						i32.add
						local.set $$t10.1
						local.get $$t10.0
						call $runtime.Block.Release
						local.set $$t10.0
						i32.const 31352
						i32.load
						i32.const 31352
						i32.load offset=4
						i32.const 31352
						i32.load offset=8
						i32.const 31352
						i32.load offset=12
						local.set $$t11.3
						local.set $$t11.2
						local.set $$t11.1
						local.set $$t11.0
						local.get $$t10.1
						local.get $$t11.0
						i32.store
						local.get $$t10.1
						local.get $$t11.1
						i32.store offset=4
						local.get $$t10.1
						local.get $$t11.2
						i32.store offset=8
						local.get $$t10.1
						local.get $$t11.3
						i32.store offset=12
						local.get $$t1.0
						call $runtime.Block.Retain
						local.get $$t1.1
						i32.const 16
						i32.const 5
						i32.mul
						i32.add
						local.set $$t12.1
						local.get $$t12.0
						call $runtime.Block.Release
						local.set $$t12.0
						i32.const 31368
						i32.load
						i32.const 31368
						i32.load offset=4
						i32.const 31368
						i32.load offset=8
						i32.const 31368
						i32.load offset=12
						local.set $$t13.3
						local.set $$t13.2
						local.set $$t13.1
						local.set $$t13.0
						local.get $$t12.1
						local.get $$t13.0
						i32.store
						local.get $$t12.1
						local.get $$t13.1
						i32.store offset=4
						local.get $$t12.1
						local.get $$t13.2
						i32.store offset=8
						local.get $$t12.1
						local.get $$t13.3
						i32.store offset=12
						local.get $$t1.0
						call $runtime.Block.Retain
						local.get $$t1.1
						i32.const 16
						i32.const 6
						i32.mul
						i32.add
						local.set $$t14.1
						local.get $$t14.0
						call $runtime.Block.Release
						local.set $$t14.0
						i32.const 31384
						i32.load
						i32.const 31384
						i32.load offset=4
						i32.const 31384
						i32.load offset=8
						i32.const 31384
						i32.load offset=12
						local.set $$t15.3
						local.set $$t15.2
						local.set $$t15.1
						local.set $$t15.0
						local.get $$t14.1
						local.get $$t15.0
						i32.store
						local.get $$t14.1
						local.get $$t15.1
						i32.store offset=4
						local.get $$t14.1
						local.get $$t15.2
						i32.store offset=8
						local.get $$t14.1
						local.get $$t15.3
						i32.store offset=12
						local.get $$t1.0
						call $runtime.Block.Retain
						local.get $$t1.1
						i32.const 16
						i32.const 7
						i32.mul
						i32.add
						local.set $$t16.1
						local.get $$t16.0
						call $runtime.Block.Release
						local.set $$t16.0
						i32.const 31400
						i32.load
						i32.const 31400
						i32.load offset=4
						i32.const 31400
						i32.load offset=8
						i32.const 31400
						i32.load offset=12
						local.set $$t17.3
						local.set $$t17.2
						local.set $$t17.1
						local.set $$t17.0
						local.get $$t16.1
						local.get $$t17.0
						i32.store
						local.get $$t16.1
						local.get $$t17.1
						i32.store offset=4
						local.get $$t16.1
						local.get $$t17.2
						i32.store offset=8
						local.get $$t16.1
						local.get $$t17.3
						i32.store offset=12
						local.get $$t1.0
						call $runtime.Block.Retain
						local.get $$t1.1
						i32.const 16
						i32.const 8
						i32.mul
						i32.add
						local.set $$t18.1
						local.get $$t18.0
						call $runtime.Block.Release
						local.set $$t18.0
						i32.const 31416
						i32.load
						i32.const 31416
						i32.load offset=4
						i32.const 31416
						i32.load offset=8
						i32.const 31416
						i32.load offset=12
						local.set $$t19.3
						local.set $$t19.2
						local.set $$t19.1
						local.set $$t19.0
						local.get $$t18.1
						local.get $$t19.0
						i32.store
						local.get $$t18.1
						local.get $$t19.1
						i32.store offset=4
						local.get $$t18.1
						local.get $$t19.2
						i32.store offset=8
						local.get $$t18.1
						local.get $$t19.3
						i32.store offset=12
						local.get $$t1.0
						call $runtime.Block.Retain
						local.get $$t1.1
						i32.const 16
						i32.const 9
						i32.mul
						i32.add
						local.set $$t20.1
						local.get $$t20.0
						call $runtime.Block.Release
						local.set $$t20.0
						i32.const 31432
						i32.load
						i32.const 31432
						i32.load offset=4
						i32.const 31432
						i32.load offset=8
						i32.const 31432
						i32.load offset=12
						local.set $$t21.3
						local.set $$t21.2
						local.set $$t21.1
						local.set $$t21.0
						local.get $$t20.1
						local.get $$t21.0
						i32.store
						local.get $$t20.1
						local.get $$t21.1
						i32.store offset=4
						local.get $$t20.1
						local.get $$t21.2
						i32.store offset=8
						local.get $$t20.1
						local.get $$t21.3
						i32.store offset=12
						local.get $$t1.0
						call $runtime.Block.Retain
						local.get $$t1.1
						i32.const 16
						i32.const 0
						i32.mul
						i32.add
						i32.const 10
						i32.const 0
						i32.sub
						i32.const 10
						i32.const 0
						i32.sub
						local.set $$t22.3
						local.set $$t22.2
						local.set $$t22.1
						local.get $$t22.0
						call $runtime.Block.Release
						local.set $$t22.0
						i32.const 31272
						local.get $$t22.0
						call $runtime.Block.Retain
						i32.const 31272
						i32.load align=1
						call $runtime.Block.Release
						i32.store align=1
						i32.const 31272
						local.get $$t22.1
						i32.store offset=4
						i32.const 31272
						local.get $$t22.2
						i32.store offset=8
						i32.const 31272
						local.get $$t22.3
						i32.store offset=12
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
		local.get $$t4.0
		call $runtime.Block.Release
		local.get $$t6.0
		call $runtime.Block.Release
		local.get $$t8.0
		call $runtime.Block.Release
		local.get $$t10.0
		call $runtime.Block.Release
		local.get $$t12.0
		call $runtime.Block.Release
		local.get $$t14.0
		call $runtime.Block.Release
		local.get $$t16.0
		call $runtime.Block.Release
		local.get $$t18.0
		call $runtime.Block.Release
		local.get $$t20.0
		call $runtime.Block.Release
		local.get $$t22.0
		call $runtime.Block.Release
	)
	(func $w42048.UI.GetPad (param $this.0 i32) (param $this.1 i32) (result i32)
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
					call $syscall$wasm4.GetGamePad1
					local.set $$t0
					local.get $$t0
					local.set $$ret_0
					br $$BlockFnBody
				end
			end
		end
		local.get $$ret_0
	)
	(func $w42048.UI.beat (param $this.0 i32) (param $this.1 i32) (param $n i32) (result i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$ret_0 i32)
		(local $$t0.0 i32)
		(local $$t0.1 i32)
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
					local.get $this.0
					call $runtime.Block.Retain
					local.get $this.1
					i32.const 16
					i32.add
					local.set $$t0.1
					local.get $$t0.0
					call $runtime.Block.Release
					local.set $$t0.0
					local.get $$t0.1
					i32.load
					local.set $$t1
					local.get $$t1
					local.get $n
					i32.rem_u
					local.set $$t2
					local.get $$t2
					i32.const 0
					i32.eq
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
	)
	(func $w42048.UI.btn1 (param $this.0 i32) (param $this.1 i32) (result i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$ret_0 i32)
		(local $$t0 i32)
		(local $$t1 i32)
		(local $$t2 i32)
		(local $$t3.0 i32)
		(local $$t3.1 i32)
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
							local.get $this.0
							local.get $this.1
							call $w42048.UI.GetPad
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
						local.get $this.0
						call $runtime.Block.Retain
						local.get $this.1
						i32.const 20
						i32.add
						local.set $$t3.1
						local.get $$t3.0
						call $runtime.Block.Release
						local.set $$t3.0
						local.get $$t3.1
						i32.load8_u align=1
						local.set $$t4
						local.get $$t4
						i32.const 1
						i32.and
						local.set $$t5
						local.get $$t5
						i32.const 0
						i32.eq
						local.set $$t6
						br $$Block_1
					end
					local.get $$current_block
					i32.const 0
					i32.eq
					if(result i32)
						i32.const 0
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
		local.get $$t3.0
		call $runtime.Block.Release
	)
	(func $w42048.UI.btn2 (param $this.0 i32) (param $this.1 i32) (result i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$ret_0 i32)
		(local $$t0 i32)
		(local $$t1 i32)
		(local $$t2 i32)
		(local $$t3.0 i32)
		(local $$t3.1 i32)
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
							local.get $this.0
							local.get $this.1
							call $w42048.UI.GetPad
							local.set $$t0
							local.get $$t0
							i32.const 2
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
						local.get $this.0
						call $runtime.Block.Retain
						local.get $this.1
						i32.const 20
						i32.add
						local.set $$t3.1
						local.get $$t3.0
						call $runtime.Block.Release
						local.set $$t3.0
						local.get $$t3.1
						i32.load8_u align=1
						local.set $$t4
						local.get $$t4
						i32.const 2
						i32.and
						local.set $$t5
						local.get $$t5
						i32.const 0
						i32.eq
						local.set $$t6
						br $$Block_1
					end
					local.get $$current_block
					i32.const 0
					i32.eq
					if(result i32)
						i32.const 0
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
		local.get $$t3.0
		call $runtime.Block.Release
	)
	(func $w42048.UI.down (param $this.0 i32) (param $this.1 i32) (result i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$ret_0 i32)
		(local $$t0 i32)
		(local $$t1 i32)
		(local $$t2 i32)
		(local $$t3.0 i32)
		(local $$t3.1 i32)
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
							local.get $this.0
							local.get $this.1
							call $w42048.UI.GetPad
							local.set $$t0
							local.get $$t0
							i32.const 128
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
						local.get $this.0
						call $runtime.Block.Retain
						local.get $this.1
						i32.const 20
						i32.add
						local.set $$t3.1
						local.get $$t3.0
						call $runtime.Block.Release
						local.set $$t3.0
						local.get $$t3.1
						i32.load8_u align=1
						local.set $$t4
						local.get $$t4
						i32.const 128
						i32.and
						local.set $$t5
						local.get $$t5
						i32.const 0
						i32.eq
						local.set $$t6
						br $$Block_1
					end
					local.get $$current_block
					i32.const 0
					i32.eq
					if(result i32)
						i32.const 0
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
		local.get $$t3.0
		call $runtime.Block.Release
	)
	(func $w42048.UI.input (param $this.0 i32) (param $this.1 i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$t0 i32)
		(local $$t1.0 i32)
		(local $$t1.1 i32)
		(local $$t2 i32)
		(local $$t3.0 i32)
		(local $$t3.1 i32)
		(local $$t4.0.0 i32)
		(local $$t4.0.1 i32)
		(local $$t4.1 i32)
		(local $$t4.2 i32)
		(local $$t5.0 i32)
		(local $$t5.1 i32)
		(local $$t6.0.0 i32)
		(local $$t6.0.1 i32)
		(local $$t6.1 i32)
		(local $$t6.2 i32)
		(local $$t7 i32)
		(local $$t8.0 i32)
		(local $$t8.1 i32)
		(local $$t9.0.0 i32)
		(local $$t9.0.1 i32)
		(local $$t9.1 i32)
		(local $$t9.2 i32)
		(local $$t10 i32)
		(local $$t11.0 i32)
		(local $$t11.1 i32)
		(local $$t12.0.0 i32)
		(local $$t12.0.1 i32)
		(local $$t12.1 i32)
		(local $$t12.2 i32)
		(local $$t13 i32)
		(local $$t14.0 i32)
		(local $$t14.1 i32)
		(local $$t15.0.0 i32)
		(local $$t15.0.1 i32)
		(local $$t15.1 i32)
		(local $$t15.2 i32)
		(local $$t16 i32)
		(local $$t17.0 i32)
		(local $$t17.1 i32)
		(local $$t18.0.0 i32)
		(local $$t18.0.1 i32)
		(local $$t18.1 i32)
		(local $$t18.2 i32)
		(local $$t19 i32)
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
																	local.get $this.0
																	local.get $this.1
																	call $w42048.UI.btn1
																	local.set $$t0
																	local.get $$t0
																	if
																		br $$Block_1
																	else
																		br $$Block_3
																	end
																end
																i32.const 1
																local.set $$current_block
																local.get $this.0
																call $runtime.Block.Retain
																local.get $this.1
																i32.const 20
																i32.add
																local.set $$t1.1
																local.get $$t1.0
																call $runtime.Block.Release
																local.set $$t1.0
																call $syscall$wasm4.GetGamePad1
																local.set $$t2
																local.get $$t1.1
																local.get $$t2
																i32.store8 align=1
																br $$BlockFnBody
															end
															i32.const 2
															local.set $$current_block
															i32.const 0
															i32.const 34906
															i32.const 4
															call $w42048.log
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
															call $runtime.Block.Retain
															local.get $$t3.1
															i32.load offset=4
															local.get $$t3.1
															i32.load offset=8
															local.get $$t3.1
															i32.load offset=12
															local.set $$t4.2
															local.set $$t4.1
															local.set $$t4.0.1
															local.get $$t4.0.0
															call $runtime.Block.Release
															local.set $$t4.0.0
															local.get $$t4.0.0
															local.get $$t4.0.1
															i32.const 4
															local.get $$t4.1
															i32.load offset=16
															call_indirect 0 (type $$$fnSig5)
															local.get $this.0
															local.get $this.1
															call $w42048.UI.randomPalette
															i32.const 1
															local.set $$block_selector
															br $$BlockDisp
														end
														i32.const 3
														local.set $$current_block
														i32.const 0
														i32.const 34910
														i32.const 4
														call $w42048.log
														local.get $this.0
														call $runtime.Block.Retain
														local.get $this.1
														i32.const 0
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
														local.set $$t6.2
														local.set $$t6.1
														local.set $$t6.0.1
														local.get $$t6.0.0
														call $runtime.Block.Release
														local.set $$t6.0.0
														local.get $$t6.0.0
														local.get $$t6.0.1
														i32.const 5
														local.get $$t6.1
														i32.load offset=16
														call_indirect 0 (type $$$fnSig5)
														i32.const 1
														local.set $$block_selector
														br $$BlockDisp
													end
													i32.const 4
													local.set $$current_block
													local.get $this.0
													local.get $this.1
													call $w42048.UI.btn2
													local.set $$t7
													local.get $$t7
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
												i32.const 0
												i32.const 34914
												i32.const 4
												call $w42048.log
												local.get $this.0
												call $runtime.Block.Retain
												local.get $this.1
												i32.const 0
												i32.add
												local.set $$t8.1
												local.get $$t8.0
												call $runtime.Block.Release
												local.set $$t8.0
												local.get $$t8.1
												i32.load
												call $runtime.Block.Retain
												local.get $$t8.1
												i32.load offset=4
												local.get $$t8.1
												i32.load offset=8
												local.get $$t8.1
												i32.load offset=12
												local.set $$t9.2
												local.set $$t9.1
												local.set $$t9.0.1
												local.get $$t9.0.0
												call $runtime.Block.Release
												local.set $$t9.0.0
												local.get $$t9.0.0
												local.get $$t9.0.1
												i32.const 0
												local.get $$t9.1
												i32.load offset=16
												call_indirect 0 (type $$$fnSig5)
												i32.const 1
												local.set $$block_selector
												br $$BlockDisp
											end
											i32.const 6
											local.set $$current_block
											local.get $this.0
											local.get $this.1
											call $w42048.UI.up
											local.set $$t10
											local.get $$t10
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
										i32.const 0
										i32.const 34918
										i32.const 4
										call $w42048.log
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
										local.set $$t12.2
										local.set $$t12.1
										local.set $$t12.0.1
										local.get $$t12.0.0
										call $runtime.Block.Release
										local.set $$t12.0.0
										local.get $$t12.0.0
										local.get $$t12.0.1
										i32.const 1
										local.get $$t12.1
										i32.load offset=16
										call_indirect 0 (type $$$fnSig5)
										i32.const 1
										local.set $$block_selector
										br $$BlockDisp
									end
									i32.const 8
									local.set $$current_block
									local.get $this.0
									local.get $this.1
									call $w42048.UI.down
									local.set $$t13
									local.get $$t13
									if
										i32.const 7
										local.set $$block_selector
										br $$BlockDisp
									else
										br $$Block_9
									end
								end
								i32.const 9
								local.set $$current_block
								i32.const 0
								i32.const 34922
								i32.const 4
								call $w42048.log
								local.get $this.0
								call $runtime.Block.Retain
								local.get $this.1
								i32.const 0
								i32.add
								local.set $$t14.1
								local.get $$t14.0
								call $runtime.Block.Release
								local.set $$t14.0
								local.get $$t14.1
								i32.load
								call $runtime.Block.Retain
								local.get $$t14.1
								i32.load offset=4
								local.get $$t14.1
								i32.load offset=8
								local.get $$t14.1
								i32.load offset=12
								local.set $$t15.2
								local.set $$t15.1
								local.set $$t15.0.1
								local.get $$t15.0.0
								call $runtime.Block.Release
								local.set $$t15.0.0
								local.get $$t15.0.0
								local.get $$t15.0.1
								i32.const 2
								local.get $$t15.1
								i32.load offset=16
								call_indirect 0 (type $$$fnSig5)
								i32.const 1
								local.set $$block_selector
								br $$BlockDisp
							end
							i32.const 10
							local.set $$current_block
							local.get $this.0
							local.get $this.1
							call $w42048.UI.right
							local.set $$t16
							local.get $$t16
							if
								i32.const 9
								local.set $$block_selector
								br $$BlockDisp
							else
								br $$Block_11
							end
						end
						i32.const 11
						local.set $$current_block
						i32.const 0
						i32.const 34926
						i32.const 4
						call $w42048.log
						local.get $this.0
						call $runtime.Block.Retain
						local.get $this.1
						i32.const 0
						i32.add
						local.set $$t17.1
						local.get $$t17.0
						call $runtime.Block.Release
						local.set $$t17.0
						local.get $$t17.1
						i32.load
						call $runtime.Block.Retain
						local.get $$t17.1
						i32.load offset=4
						local.get $$t17.1
						i32.load offset=8
						local.get $$t17.1
						i32.load offset=12
						local.set $$t18.2
						local.set $$t18.1
						local.set $$t18.0.1
						local.get $$t18.0.0
						call $runtime.Block.Release
						local.set $$t18.0.0
						local.get $$t18.0.0
						local.get $$t18.0.1
						i32.const 3
						local.get $$t18.1
						i32.load offset=16
						call_indirect 0 (type $$$fnSig5)
						i32.const 1
						local.set $$block_selector
						br $$BlockDisp
					end
					i32.const 12
					local.set $$current_block
					local.get $this.0
					local.get $this.1
					call $w42048.UI.left
					local.set $$t19
					local.get $$t19
					if
						i32.const 11
						local.set $$block_selector
						br $$BlockDisp
					else
						i32.const 1
						local.set $$block_selector
						br $$BlockDisp
					end
				end
			end
		end
		local.get $$t1.0
		call $runtime.Block.Release
		local.get $$t3.0
		call $runtime.Block.Release
		local.get $$t4.0.0
		call $runtime.Block.Release
		local.get $$t5.0
		call $runtime.Block.Release
		local.get $$t6.0.0
		call $runtime.Block.Release
		local.get $$t8.0
		call $runtime.Block.Release
		local.get $$t9.0.0
		call $runtime.Block.Release
		local.get $$t11.0
		call $runtime.Block.Release
		local.get $$t12.0.0
		call $runtime.Block.Release
		local.get $$t14.0
		call $runtime.Block.Release
		local.get $$t15.0.0
		call $runtime.Block.Release
		local.get $$t17.0
		call $runtime.Block.Release
		local.get $$t18.0.0
		call $runtime.Block.Release
	)
	(func $w42048.UI.left (param $this.0 i32) (param $this.1 i32) (result i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$ret_0 i32)
		(local $$t0 i32)
		(local $$t1 i32)
		(local $$t2 i32)
		(local $$t3.0 i32)
		(local $$t3.1 i32)
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
							local.get $this.0
							local.get $this.1
							call $w42048.UI.GetPad
							local.set $$t0
							local.get $$t0
							i32.const 16
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
						local.get $this.0
						call $runtime.Block.Retain
						local.get $this.1
						i32.const 20
						i32.add
						local.set $$t3.1
						local.get $$t3.0
						call $runtime.Block.Release
						local.set $$t3.0
						local.get $$t3.1
						i32.load8_u align=1
						local.set $$t4
						local.get $$t4
						i32.const 16
						i32.and
						local.set $$t5
						local.get $$t5
						i32.const 0
						i32.eq
						local.set $$t6
						br $$Block_1
					end
					local.get $$current_block
					i32.const 0
					i32.eq
					if(result i32)
						i32.const 0
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
		local.get $$t3.0
		call $runtime.Block.Release
	)
	(func $w42048.UI.music (param $this.0 i32) (param $this.1 i32)
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
		(local $$t11.2 i32)
		(local $$t11.3 i32)
		(local $$t11.4 i32)
		(local $$t11.5 i32)
		(local $$t11.6 i32)
		(local $$t11.7 i32)
		(local $$t11.8 i32)
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
		(local $$t22.0 i32)
		(local $$t22.1 i32)
		(local $$t22.2 i32)
		(local $$t22.3 i32)
		(local $$t22.4 i32)
		(local $$t22.5 i32)
		(local $$t22.6 i32)
		(local $$t22.7 i32)
		(local $$t22.8 i32)
		(local $$t23 i32)
		(local $$t24.0 i32)
		(local $$t24.1 i32)
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
		(local $$t34.2 i32)
		(local $$t34.3 i32)
		(local $$t34.4 i32)
		(local $$t34.5 i32)
		(local $$t34.6 i32)
		(local $$t34.7 i32)
		(local $$t34.8 i32)
		(local $$t35 i32)
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
													i32.const 1
													if
														br $$Block_1
													else
														br $$Block_0
													end
												end
												i32.const 1
												local.set $$current_block
												br $$BlockFnBody
											end
											i32.const 2
											local.set $$current_block
											local.get $this.0
											local.get $this.1
											i32.const 24
											call $w42048.UI.beat
											local.set $$t0
											local.get $$t0
											if
												br $$Block_3
											else
												br $$Block_5
											end
										end
										i32.const 3
										local.set $$current_block
										br $$BlockFnBody
									end
									i32.const 4
									local.set $$current_block
									i32.const 52
									call $runtime.HeapAlloc
									i32.const 1
									i32.const 0
									i32.const 36
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
									local.get $$t1.0
									call $runtime.Block.Retain
									local.get $$t1.1
									i32.const 4
									i32.add
									local.set $$t3.1
									local.get $$t3.0
									call $runtime.Block.Release
									local.set $$t3.0
									local.get $$t1.0
									call $runtime.Block.Retain
									local.get $$t1.1
									i32.const 8
									i32.add
									local.set $$t4.1
									local.get $$t4.0
									call $runtime.Block.Release
									local.set $$t4.0
									local.get $$t1.0
									call $runtime.Block.Retain
									local.get $$t1.1
									i32.const 12
									i32.add
									local.set $$t5.1
									local.get $$t5.0
									call $runtime.Block.Release
									local.set $$t5.0
									local.get $$t1.0
									call $runtime.Block.Retain
									local.get $$t1.1
									i32.const 16
									i32.add
									local.set $$t6.1
									local.get $$t6.0
									call $runtime.Block.Release
									local.set $$t6.0
									local.get $$t1.0
									call $runtime.Block.Retain
									local.get $$t1.1
									i32.const 20
									i32.add
									local.set $$t7.1
									local.get $$t7.0
									call $runtime.Block.Release
									local.set $$t7.0
									local.get $$t1.0
									call $runtime.Block.Retain
									local.get $$t1.1
									i32.const 24
									i32.add
									local.set $$t8.1
									local.get $$t8.0
									call $runtime.Block.Release
									local.set $$t8.0
									local.get $$t1.0
									call $runtime.Block.Retain
									local.get $$t1.1
									i32.const 28
									i32.add
									local.set $$t9.1
									local.get $$t9.0
									call $runtime.Block.Release
									local.set $$t9.0
									local.get $$t1.0
									call $runtime.Block.Retain
									local.get $$t1.1
									i32.const 32
									i32.add
									local.set $$t10.1
									local.get $$t10.0
									call $runtime.Block.Release
									local.set $$t10.0
									local.get $$t2.1
									i32.const 40
									i32.store
									local.get $$t3.1
									i32.const 0
									i32.store
									local.get $$t4.1
									i32.const 0
									i32.store
									local.get $$t5.1
									i32.const 8
									i32.store
									local.get $$t6.1
									i32.const 0
									i32.store
									local.get $$t7.1
									i32.const 0
									i32.store
									local.get $$t8.1
									i32.const 0
									i32.store
									local.get $$t9.1
									i32.const 2
									i32.store
									local.get $$t10.1
									i32.const 0
									i32.store
									local.get $$t1.1
									i32.load
									local.get $$t1.1
									i32.load offset=4
									local.get $$t1.1
									i32.load offset=8
									local.get $$t1.1
									i32.load offset=12
									local.get $$t1.1
									i32.load offset=16
									local.get $$t1.1
									i32.load offset=20
									local.get $$t1.1
									i32.load offset=24
									local.get $$t1.1
									i32.load offset=28
									local.get $$t1.1
									i32.load offset=32
									local.set $$t11.8
									local.set $$t11.7
									local.set $$t11.6
									local.set $$t11.5
									local.set $$t11.4
									local.set $$t11.3
									local.set $$t11.2
									local.set $$t11.1
									local.set $$t11.0
									local.get $$t11.0
									local.get $$t11.1
									local.get $$t11.2
									local.get $$t11.3
									local.get $$t11.4
									local.get $$t11.5
									local.get $$t11.6
									local.get $$t11.7
									local.get $$t11.8
									call $w42048.play
									i32.const 3
									local.set $$block_selector
									br $$BlockDisp
								end
								i32.const 5
								local.set $$current_block
								i32.const 52
								call $runtime.HeapAlloc
								i32.const 1
								i32.const 0
								i32.const 36
								call $runtime.Block.Init
								call $runtime.DupI32
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
								local.get $$t12.0
								call $runtime.Block.Retain
								local.get $$t12.1
								i32.const 8
								i32.add
								local.set $$t15.1
								local.get $$t15.0
								call $runtime.Block.Release
								local.set $$t15.0
								local.get $$t12.0
								call $runtime.Block.Retain
								local.get $$t12.1
								i32.const 12
								i32.add
								local.set $$t16.1
								local.get $$t16.0
								call $runtime.Block.Release
								local.set $$t16.0
								local.get $$t12.0
								call $runtime.Block.Retain
								local.get $$t12.1
								i32.const 16
								i32.add
								local.set $$t17.1
								local.get $$t17.0
								call $runtime.Block.Release
								local.set $$t17.0
								local.get $$t12.0
								call $runtime.Block.Retain
								local.get $$t12.1
								i32.const 20
								i32.add
								local.set $$t18.1
								local.get $$t18.0
								call $runtime.Block.Release
								local.set $$t18.0
								local.get $$t12.0
								call $runtime.Block.Retain
								local.get $$t12.1
								i32.const 24
								i32.add
								local.set $$t19.1
								local.get $$t19.0
								call $runtime.Block.Release
								local.set $$t19.0
								local.get $$t12.0
								call $runtime.Block.Retain
								local.get $$t12.1
								i32.const 28
								i32.add
								local.set $$t20.1
								local.get $$t20.0
								call $runtime.Block.Release
								local.set $$t20.0
								local.get $$t12.0
								call $runtime.Block.Retain
								local.get $$t12.1
								i32.const 32
								i32.add
								local.set $$t21.1
								local.get $$t21.0
								call $runtime.Block.Release
								local.set $$t21.0
								local.get $$t13.1
								i32.const 80
								i32.store
								local.get $$t14.1
								i32.const 0
								i32.store
								local.get $$t15.1
								i32.const 0
								i32.store
								local.get $$t16.1
								i32.const 8
								i32.store
								local.get $$t17.1
								i32.const 0
								i32.store
								local.get $$t18.1
								i32.const 0
								i32.store
								local.get $$t19.1
								i32.const 0
								i32.store
								local.get $$t20.1
								i32.const 2
								i32.store
								local.get $$t21.1
								i32.const 0
								i32.store
								local.get $$t12.1
								i32.load
								local.get $$t12.1
								i32.load offset=4
								local.get $$t12.1
								i32.load offset=8
								local.get $$t12.1
								i32.load offset=12
								local.get $$t12.1
								i32.load offset=16
								local.get $$t12.1
								i32.load offset=20
								local.get $$t12.1
								i32.load offset=24
								local.get $$t12.1
								i32.load offset=28
								local.get $$t12.1
								i32.load offset=32
								local.set $$t22.8
								local.set $$t22.7
								local.set $$t22.6
								local.set $$t22.5
								local.set $$t22.4
								local.set $$t22.3
								local.set $$t22.2
								local.set $$t22.1
								local.set $$t22.0
								local.get $$t22.0
								local.get $$t22.1
								local.get $$t22.2
								local.get $$t22.3
								local.get $$t22.4
								local.get $$t22.5
								local.get $$t22.6
								local.get $$t22.7
								local.get $$t22.8
								call $w42048.play
								i32.const 3
								local.set $$block_selector
								br $$BlockDisp
							end
							i32.const 6
							local.set $$current_block
							local.get $this.0
							local.get $this.1
							i32.const 32
							call $w42048.UI.beat
							local.set $$t23
							local.get $$t23
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
						i32.const 52
						call $runtime.HeapAlloc
						i32.const 1
						i32.const 0
						i32.const 36
						call $runtime.Block.Init
						call $runtime.DupI32
						i32.const 16
						i32.add
						local.set $$t24.1
						local.get $$t24.0
						call $runtime.Block.Release
						local.set $$t24.0
						local.get $$t24.0
						call $runtime.Block.Retain
						local.get $$t24.1
						i32.const 0
						i32.add
						local.set $$t25.1
						local.get $$t25.0
						call $runtime.Block.Release
						local.set $$t25.0
						local.get $$t24.0
						call $runtime.Block.Retain
						local.get $$t24.1
						i32.const 4
						i32.add
						local.set $$t26.1
						local.get $$t26.0
						call $runtime.Block.Release
						local.set $$t26.0
						local.get $$t24.0
						call $runtime.Block.Retain
						local.get $$t24.1
						i32.const 8
						i32.add
						local.set $$t27.1
						local.get $$t27.0
						call $runtime.Block.Release
						local.set $$t27.0
						local.get $$t24.0
						call $runtime.Block.Retain
						local.get $$t24.1
						i32.const 12
						i32.add
						local.set $$t28.1
						local.get $$t28.0
						call $runtime.Block.Release
						local.set $$t28.0
						local.get $$t24.0
						call $runtime.Block.Retain
						local.get $$t24.1
						i32.const 16
						i32.add
						local.set $$t29.1
						local.get $$t29.0
						call $runtime.Block.Release
						local.set $$t29.0
						local.get $$t24.0
						call $runtime.Block.Retain
						local.get $$t24.1
						i32.const 20
						i32.add
						local.set $$t30.1
						local.get $$t30.0
						call $runtime.Block.Release
						local.set $$t30.0
						local.get $$t24.0
						call $runtime.Block.Retain
						local.get $$t24.1
						i32.const 24
						i32.add
						local.set $$t31.1
						local.get $$t31.0
						call $runtime.Block.Release
						local.set $$t31.0
						local.get $$t24.0
						call $runtime.Block.Retain
						local.get $$t24.1
						i32.const 28
						i32.add
						local.set $$t32.1
						local.get $$t32.0
						call $runtime.Block.Release
						local.set $$t32.0
						local.get $$t24.0
						call $runtime.Block.Retain
						local.get $$t24.1
						i32.const 32
						i32.add
						local.set $$t33.1
						local.get $$t33.0
						call $runtime.Block.Release
						local.set $$t33.0
						local.get $$t25.1
						i32.const 160
						i32.store
						local.get $$t26.1
						i32.const 0
						i32.store
						local.get $$t27.1
						i32.const 0
						i32.store
						local.get $$t28.1
						i32.const 8
						i32.store
						local.get $$t29.1
						i32.const 0
						i32.store
						local.get $$t30.1
						i32.const 0
						i32.store
						local.get $$t31.1
						i32.const 0
						i32.store
						local.get $$t32.1
						i32.const 2
						i32.store
						local.get $$t33.1
						i32.const 0
						i32.store
						local.get $$t24.1
						i32.load
						local.get $$t24.1
						i32.load offset=4
						local.get $$t24.1
						i32.load offset=8
						local.get $$t24.1
						i32.load offset=12
						local.get $$t24.1
						i32.load offset=16
						local.get $$t24.1
						i32.load offset=20
						local.get $$t24.1
						i32.load offset=24
						local.get $$t24.1
						i32.load offset=28
						local.get $$t24.1
						i32.load offset=32
						local.set $$t34.8
						local.set $$t34.7
						local.set $$t34.6
						local.set $$t34.5
						local.set $$t34.4
						local.set $$t34.3
						local.set $$t34.2
						local.set $$t34.1
						local.set $$t34.0
						local.get $$t34.0
						local.get $$t34.1
						local.get $$t34.2
						local.get $$t34.3
						local.get $$t34.4
						local.get $$t34.5
						local.get $$t34.6
						local.get $$t34.7
						local.get $$t34.8
						call $w42048.play
						i32.const 3
						local.set $$block_selector
						br $$BlockDisp
					end
					i32.const 8
					local.set $$current_block
					local.get $this.0
					local.get $this.1
					i32.const 96
					call $w42048.UI.beat
					local.set $$t35
					local.get $$t35
					if
						i32.const 7
						local.set $$block_selector
						br $$BlockDisp
					else
						i32.const 3
						local.set $$block_selector
						br $$BlockDisp
					end
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
		local.get $$t24.0
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
		local.get $$t31.0
		call $runtime.Block.Release
		local.get $$t32.0
		call $runtime.Block.Release
		local.get $$t33.0
		call $runtime.Block.Release
	)
	(func $w42048.UI.randomPalette (param $this.0 i32) (param $this.1 i32)
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
		(local $$t3 i32)
		(local $$t4 i32)
		(local $$t5.0 i32)
		(local $$t5.1 i32)
		(local $$t5.2 i32)
		(local $$t5.3 i32)
		(local $$t6 i32)
		(local $$t7 i32)
		(local $$t8.0 i32)
		(local $$t8.1 i32)
		(local $$t9.0 i32)
		(local $$t9.1 i32)
		(local $$t9.2 i32)
		(local $$t9.3 i32)
		(local $$t10.0 i32)
		(local $$t10.1 i32)
		(local $$t10.2 i32)
		(local $$t10.3 i32)
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
					i32.const 24
					i32.add
					local.set $$t0.1
					local.get $$t0.0
					call $runtime.Block.Release
					local.set $$t0.0
					i32.const 31272
					i32.load
					call $runtime.Block.Retain
					i32.const 31272
					i32.load offset=4
					i32.const 31272
					i32.load offset=8
					i32.const 31272
					i32.load offset=12
					local.set $$t1.3
					local.set $$t1.2
					local.set $$t1.1
					local.get $$t1.0
					call $runtime.Block.Release
					local.set $$t1.0
					local.get $this.0
					call $runtime.Block.Retain
					local.get $this.1
					i32.const 16
					i32.add
					local.set $$t2.1
					local.get $$t2.0
					call $runtime.Block.Release
					local.set $$t2.0
					local.get $$t2.1
					i32.load
					local.set $$t3
					local.get $$t3
					local.set $$t4
					i32.const 31272
					i32.load
					call $runtime.Block.Retain
					i32.const 31272
					i32.load offset=4
					i32.const 31272
					i32.load offset=8
					i32.const 31272
					i32.load offset=12
					local.set $$t5.3
					local.set $$t5.2
					local.set $$t5.1
					local.get $$t5.0
					call $runtime.Block.Release
					local.set $$t5.0
					local.get $$t5.2
					local.set $$t6
					local.get $$t4
					local.get $$t6
					i32.rem_s
					local.set $$t7
					local.get $$t1.0
					call $runtime.Block.Retain
					local.get $$t1.1
					i32.const 16
					local.get $$t7
					i32.mul
					i32.add
					local.set $$t8.1
					local.get $$t8.0
					call $runtime.Block.Release
					local.set $$t8.0
					local.get $$t8.1
					i32.load
					local.get $$t8.1
					i32.load offset=4
					local.get $$t8.1
					i32.load offset=8
					local.get $$t8.1
					i32.load offset=12
					local.set $$t9.3
					local.set $$t9.2
					local.set $$t9.1
					local.set $$t9.0
					local.get $$t9.0
					local.get $$t9.1
					local.get $$t9.2
					local.get $$t9.3
					local.set $$t10.3
					local.set $$t10.2
					local.set $$t10.1
					local.set $$t10.0
					local.get $$t0.1
					local.get $$t10.0
					i32.store
					local.get $$t0.1
					local.get $$t10.1
					i32.store offset=4
					local.get $$t0.1
					local.get $$t10.2
					i32.store offset=8
					local.get $$t0.1
					local.get $$t10.3
					i32.store offset=12
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
		local.get $$t5.0
		call $runtime.Block.Release
		local.get $$t8.0
		call $runtime.Block.Release
	)
	(func $w42048.UI.right (param $this.0 i32) (param $this.1 i32) (result i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$ret_0 i32)
		(local $$t0 i32)
		(local $$t1 i32)
		(local $$t2 i32)
		(local $$t3.0 i32)
		(local $$t3.1 i32)
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
							local.get $this.0
							local.get $this.1
							call $w42048.UI.GetPad
							local.set $$t0
							local.get $$t0
							i32.const 32
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
						local.get $this.0
						call $runtime.Block.Retain
						local.get $this.1
						i32.const 20
						i32.add
						local.set $$t3.1
						local.get $$t3.0
						call $runtime.Block.Release
						local.set $$t3.0
						local.get $$t3.1
						i32.load8_u align=1
						local.set $$t4
						local.get $$t4
						i32.const 32
						i32.and
						local.set $$t5
						local.get $$t5
						i32.const 0
						i32.eq
						local.set $$t6
						br $$Block_1
					end
					local.get $$current_block
					i32.const 0
					i32.eq
					if(result i32)
						i32.const 0
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
		local.get $$t3.0
		call $runtime.Block.Release
	)
	(func $w42048.UI.show (param $this.0 i32) (param $this.1 i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$t0.0 i32)
		(local $$t0.1 i32)
		(local $$t1.0 i32)
		(local $$t1.1 i32)
		(local $$t2 i32)
		(local $$t3.0 i32)
		(local $$t3.1 i32)
		(local $$t4.0 i32)
		(local $$t4.1 i32)
		(local $$t5 i32)
		(local $$t6.0 i32)
		(local $$t6.1 i32)
		(local $$t7.0 i32)
		(local $$t7.1 i32)
		(local $$t8 i32)
		(local $$t9.0 i32)
		(local $$t9.1 i32)
		(local $$t10.0 i32)
		(local $$t10.1 i32)
		(local $$t11 i32)
		(local $$t12.0 i32)
		(local $$t12.1 i32)
		(local $$t13.0.0 i32)
		(local $$t13.0.1 i32)
		(local $$t13.1 i32)
		(local $$t13.2 i32)
		(local $$t14 i32)
		(local $$t15.0 i32)
		(local $$t15.1 i32)
		(local $$t15.2 i32)
		(local $$t16.0 i32)
		(local $$t16.1 i32)
		(local $$t16.2 i32)
		(local $$t17.0 i32)
		(local $$t17.1 i32)
		(local $$t17.2 i32)
		(local $$t18.0 i32)
		(local $$t18.1 i32)
		(local $$t19 i32)
		(local $$t20 i32)
		(local $$t21 i32)
		(local $$t22 i32)
		(local $$t23 i32)
		(local $$t24.0 i32)
		(local $$t24.1 i32)
		(local $$t25.0.0 i32)
		(local $$t25.0.1 i32)
		(local $$t25.1 i32)
		(local $$t25.2 i32)
		(local $$t26 i32)
		(local $$t27 i32)
		(local $$t28 i32)
		(local $$t29 i32)
		(local $$t30 i32)
		(local $$t31 i32)
		(local $$t32 i32)
		(local $$t33 i32)
		(local $$t34.0 i32)
		(local $$t34.1 i32)
		(local $$t35 i32)
		(local $$t36 i32)
		(local $$t37.0 i32)
		(local $$t37.1 i32)
		(local $$t37.2 i32)
		(local $$t38 i32)
		(local $$t39 i32)
		(local $$t40 i32)
		(local $$t41 i32)
		(local $$t42 i32)
		(local $$t43 i32)
		(local $$t44 i32)
		(local $$t45 i32)
		(local $$t46 i32)
		(local $$t47 i32)
		(local $$t48 i32)
		(local $$t49 i32)
		(local $$t50.0 i32)
		(local $$t50.1 i32)
		(local $$t51 i32)
		(local $$t52 i32)
		(local $$t53.0 i32)
		(local $$t53.1 i32)
		(local $$t54 i32)
		(local $$t55 i32)
		(local $$t56 i32)
		(local $$t57.0 i32)
		(local $$t57.1 i32)
		(local $$t58 i32)
		(local $$t59 i32)
		(local $$t60.0 i32)
		(local $$t60.1 i32)
		(local $$t61 i32)
		(local $$t62 i32)
		(local $$t63 i32)
		(local $$t64.0 i32)
		(local $$t64.1 i32)
		(local $$t65 i32)
		(local $$t66 i32)
		(local $$t67.0 i32)
		(local $$t67.1 i32)
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
																																	br_table  0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 0
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
																																local.get $$t1.1
																																i32.load
																																local.set $$t2
																																local.get $this.0
																																call $runtime.Block.Retain
																																local.get $this.1
																																i32.const 24
																																i32.add
																																local.set $$t3.1
																																local.get $$t3.0
																																call $runtime.Block.Release
																																local.set $$t3.0
																																local.get $$t3.0
																																call $runtime.Block.Retain
																																local.get $$t3.1
																																i32.const 4
																																i32.const 1
																																i32.mul
																																i32.add
																																local.set $$t4.1
																																local.get $$t4.0
																																call $runtime.Block.Release
																																local.set $$t4.0
																																local.get $$t4.1
																																i32.load
																																local.set $$t5
																																local.get $this.0
																																call $runtime.Block.Retain
																																local.get $this.1
																																i32.const 24
																																i32.add
																																local.set $$t6.1
																																local.get $$t6.0
																																call $runtime.Block.Release
																																local.set $$t6.0
																																local.get $$t6.0
																																call $runtime.Block.Retain
																																local.get $$t6.1
																																i32.const 4
																																i32.const 2
																																i32.mul
																																i32.add
																																local.set $$t7.1
																																local.get $$t7.0
																																call $runtime.Block.Release
																																local.set $$t7.0
																																local.get $$t7.1
																																i32.load
																																local.set $$t8
																																local.get $this.0
																																call $runtime.Block.Retain
																																local.get $this.1
																																i32.const 24
																																i32.add
																																local.set $$t9.1
																																local.get $$t9.0
																																call $runtime.Block.Release
																																local.set $$t9.0
																																local.get $$t9.0
																																call $runtime.Block.Retain
																																local.get $$t9.1
																																i32.const 4
																																i32.const 3
																																i32.mul
																																i32.add
																																local.set $$t10.1
																																local.get $$t10.0
																																call $runtime.Block.Release
																																local.set $$t10.0
																																local.get $$t10.1
																																i32.load
																																local.set $$t11
																																local.get $$t2
																																local.get $$t5
																																local.get $$t8
																																local.get $$t11
																																call $syscall$wasm4.SetPalette
																																i32.const 0
																																i32.const 0
																																i32.const 160
																																i32.const 27
																																i32.const 3
																																i32.const 3
																																i32.const 4
																																call $w42048.dotbg
																																i32.const 50
																																call $syscall$wasm4.SetDrawColorsU16
																																i32.const 0
																																i32.const 6
																																i32.const 84
																																i32.const 12
																																call $syscall$wasm4.RectI32
																																i32.const 4
																																call $syscall$wasm4.SetDrawColorsU16
																																i32.const 0
																																i32.const 157
																																i32.const 160
																																i32.const 3
																																call $syscall$wasm4.RectI32
																																i32.const 36
																																call $syscall$wasm4.SetDrawColorsU16
																																local.get $this.0
																																call $runtime.Block.Retain
																																local.get $this.1
																																i32.const 0
																																i32.add
																																local.set $$t12.1
																																local.get $$t12.0
																																call $runtime.Block.Release
																																local.set $$t12.0
																																local.get $$t12.1
																																i32.load
																																call $runtime.Block.Retain
																																local.get $$t12.1
																																i32.load offset=4
																																local.get $$t12.1
																																i32.load offset=8
																																local.get $$t12.1
																																i32.load offset=12
																																local.set $$t13.2
																																local.set $$t13.1
																																local.set $$t13.0.1
																																local.get $$t13.0.0
																																call $runtime.Block.Release
																																local.set $$t13.0.0
																																local.get $$t13.0.0
																																local.get $$t13.0.1
																																local.get $$t13.1
																																i32.load offset=28
																																call_indirect 0 (type $$$fnSig7)
																																local.set $$t14
																																local.get $$t14
																																call $strconv.Itoa
																																local.set $$t15.2
																																local.set $$t15.1
																																local.get $$t15.0
																																call $runtime.Block.Release
																																local.set $$t15.0
																																local.get $$t15.0
																																local.get $$t15.1
																																local.get $$t15.2
																																i32.const 0
																																i32.const 14981
																																i32.const 1
																																i32.const 5
																																call $w42048.leftpad
																																local.set $$t16.2
																																local.set $$t16.1
																																local.get $$t16.0
																																call $runtime.Block.Release
																																local.set $$t16.0
																																i32.const 0
																																i32.const 34933
																																i32.const 5
																																local.get $$t16.0
																																local.get $$t16.1
																																local.get $$t16.2
																																call $$string.appendstr
																																local.set $$t17.2
																																local.set $$t17.1
																																local.get $$t17.0
																																call $runtime.Block.Release
																																local.set $$t17.0
																																local.get $$t17.0
																																local.get $$t17.1
																																local.get $$t17.2
																																i32.const 2
																																i32.const 8
																																call $syscall$wasm4.Text
																																i32.const 67
																																call $syscall$wasm4.SetDrawColorsU16
																																i32.const 0
																																i32.const 34938
																																i32.const 11
																																i32.const 0
																																i32.const 19
																																call $syscall$wasm4.Text
																																i32.const 65
																																call $syscall$wasm4.SetDrawColorsU16
																																i32.const 0
																																i32.const 34949
																																i32.const 9
																																i32.const 88
																																i32.const 0
																																call $syscall$wasm4.Text
																																i32.const 66
																																call $syscall$wasm4.SetDrawColorsU16
																																i32.const 0
																																i32.const 34958
																																i32.const 9
																																i32.const 88
																																i32.const 8
																																call $syscall$wasm4.Text
																																i32.const 67
																																call $syscall$wasm4.SetDrawColorsU16
																																i32.const 0
																																i32.const 34949
																																i32.const 9
																																i32.const 88
																																i32.const 16
																																call $syscall$wasm4.Text
																																br $$Block_2
																															end
																															i32.const 1
																															local.set $$current_block
																															br $$Block_5
																														end
																														i32.const 2
																														local.set $$current_block
																														i32.const 3
																														call $syscall$wasm4.SetDrawColorsU16
																														i32.const 0
																														i32.const 28
																														i32.const 160
																														i32.const 28
																														call $syscall$wasm4.LineI32
																														i32.const 0
																														i32.const 60
																														i32.const 160
																														i32.const 60
																														call $syscall$wasm4.LineI32
																														i32.const 0
																														i32.const 92
																														i32.const 160
																														i32.const 92
																														call $syscall$wasm4.LineI32
																														i32.const 0
																														i32.const 124
																														i32.const 160
																														i32.const 124
																														call $syscall$wasm4.LineI32
																														i32.const 0
																														i32.const 156
																														i32.const 160
																														i32.const 156
																														call $syscall$wasm4.LineI32
																														i32.const 2
																														call $syscall$wasm4.SetDrawColorsU16
																														i32.const 0
																														i32.const 27
																														i32.const 160
																														i32.const 27
																														call $syscall$wasm4.LineI32
																														i32.const 0
																														i32.const 59
																														i32.const 160
																														i32.const 59
																														call $syscall$wasm4.LineI32
																														i32.const 0
																														i32.const 91
																														i32.const 160
																														i32.const 91
																														call $syscall$wasm4.LineI32
																														i32.const 0
																														i32.const 123
																														i32.const 160
																														i32.const 123
																														call $syscall$wasm4.LineI32
																														i32.const 0
																														i32.const 155
																														i32.const 160
																														i32.const 155
																														call $syscall$wasm4.LineI32
																														i32.const 39
																														i32.const 28
																														i32.const 39
																														i32.const 155
																														call $syscall$wasm4.LineI32
																														i32.const 80
																														i32.const 28
																														i32.const 80
																														i32.const 155
																														call $syscall$wasm4.LineI32
																														i32.const 120
																														i32.const 28
																														i32.const 120
																														i32.const 155
																														call $syscall$wasm4.LineI32
																														local.get $this.0
																														call $runtime.Block.Retain
																														local.get $this.1
																														i32.const 16
																														i32.add
																														local.set $$t18.1
																														local.get $$t18.0
																														call $runtime.Block.Release
																														local.set $$t18.0
																														local.get $$t18.1
																														i32.load
																														local.set $$t19
																														local.get $$t19
																														i32.const 124
																														i32.lt_u
																														local.set $$t20
																														local.get $$t20
																														if
																															br $$Block_8
																														else
																															br $$Block_9
																														end
																													end
																													local.get $$current_block
																													i32.const 0
																													i32.eq
																													if(result i32)
																														i32.const 0
																													else
																														local.get $$t21
																													end
																													local.set $$t22
																													i32.const 3
																													local.set $$current_block
																													local.get $$t22
																													i32.const 4
																													i32.lt_s
																													local.set $$t23
																													local.get $$t23
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
																												call $runtime.Block.Retain
																												local.get $$t24.1
																												i32.load offset=4
																												local.get $$t24.1
																												i32.load offset=8
																												local.get $$t24.1
																												i32.load offset=12
																												local.set $$t25.2
																												local.set $$t25.1
																												local.set $$t25.0.1
																												local.get $$t25.0.0
																												call $runtime.Block.Release
																												local.set $$t25.0.0
																												local.get $$t22
																												local.set $$t26
																												local.get $$t27
																												local.set $$t28
																												local.get $$t25.0.0
																												local.get $$t25.0.1
																												local.get $$t26
																												local.get $$t28
																												local.get $$t25.1
																												i32.load offset=12
																												call_indirect 0 (type $$$fnSig4)
																												local.set $$t29
																												local.get $$t29
																												i32.const 0
																												i32.gt_s
																												local.set $$t30
																												local.get $$t30
																												if
																													br $$Block_6
																												else
																													br $$Block_7
																												end
																											end
																											i32.const 5
																											local.set $$current_block
																											local.get $$t22
																											i32.const 1
																											i32.add
																											local.set $$t21
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
																											local.get $$t31
																										end
																										local.set $$t27
																										i32.const 6
																										local.set $$current_block
																										local.get $$t27
																										i32.const 4
																										i32.lt_s
																										local.set $$t32
																										local.get $$t32
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
																									local.get $$t29
																									local.set $$t33
																									local.get $$t22
																									local.get $$t27
																									local.get $$t33
																									call $w42048.showTile
																									br $$Block_7
																								end
																								i32.const 8
																								local.set $$current_block
																								local.get $$t27
																								i32.const 1
																								i32.add
																								local.set $$t31
																								i32.const 6
																								local.set $$block_selector
																								br $$BlockDisp
																							end
																							i32.const 9
																							local.set $$current_block
																							local.get $this.0
																							call $runtime.Block.Retain
																							local.get $this.1
																							i32.const 16
																							i32.add
																							local.set $$t34.1
																							local.get $$t34.0
																							call $runtime.Block.Release
																							local.set $$t34.0
																							local.get $$t34.1
																							i32.load
																							local.set $$t35
																							local.get $$t35
																							i32.const 6
																							i32.lt_u
																							local.set $$t36
																							local.get $$t36
																							if
																								br $$Block_15
																							else
																								br $$Block_14
																							end
																						end
																						i32.const 10
																						local.set $$current_block
																						br $$BlockFnBody
																					end
																					local.get $$current_block
																					i32.const 12
																					i32.eq
																					if(result i32 i32 i32)
																						i32.const 0
																						i32.const 34967
																						i32.const 25
																					else
																						local.get $$current_block
																						i32.const 13
																						i32.eq
																						if(result i32 i32 i32)
																							i32.const 0
																							i32.const 34992
																							i32.const 19
																						else
																							local.get $$current_block
																							i32.const 17
																							i32.eq
																							if(result i32 i32 i32)
																								i32.const 0
																								i32.const 35011
																								i32.const 13
																							else
																								local.get $$current_block
																								i32.const 21
																								i32.eq
																								if(result i32 i32 i32)
																									i32.const 0
																									i32.const 35024
																									i32.const 8
																								else
																									i32.const 0
																									i32.const 35024
																									i32.const 8
																								end
																							end
																						end
																					end
																					local.get $$current_block
																					i32.const 12
																					i32.eq
																					if(result i32)
																						i32.const 2
																					else
																						local.get $$current_block
																						i32.const 13
																						i32.eq
																						if(result i32)
																							i32.const 3
																						else
																							local.get $$current_block
																							i32.const 17
																							i32.eq
																							if(result i32)
																								i32.const 3
																							else
																								local.get $$current_block
																								i32.const 21
																								i32.eq
																								if(result i32)
																									i32.const 4
																								else
																									i32.const 1
																								end
																							end
																						end
																					end
																					local.get $$current_block
																					i32.const 12
																					i32.eq
																					if(result i32)
																						i32.const 33
																					else
																						local.get $$current_block
																						i32.const 13
																						i32.eq
																						if(result i32)
																							i32.const 33
																						else
																							local.get $$current_block
																							i32.const 17
																							i32.eq
																							if(result i32)
																								i32.const 50
																							else
																								local.get $$current_block
																								i32.const 21
																								i32.eq
																								if(result i32)
																									i32.const 35
																								else
																									i32.const 52
																								end
																							end
																						end
																					end
																					local.get $$current_block
																					i32.const 12
																					i32.eq
																					if(result i32)
																						i32.const 120
																					else
																						local.get $$current_block
																						i32.const 13
																						i32.eq
																						if(result i32)
																							i32.const 115
																						else
																							local.get $$current_block
																							i32.const 17
																							i32.eq
																							if(result i32)
																								i32.const 110
																							else
																								local.get $$current_block
																								i32.const 21
																								i32.eq
																								if(result i32)
																									i32.const 105
																								else
																									i32.const 100
																								end
																							end
																						end
																					end
																					local.get $$current_block
																					i32.const 12
																					i32.eq
																					if(result i32)
																						i32.const 10
																					else
																						local.get $$current_block
																						i32.const 13
																						i32.eq
																						if(result i32)
																							i32.const 15
																						else
																							local.get $$current_block
																							i32.const 17
																							i32.eq
																							if(result i32)
																								i32.const 20
																							else
																								local.get $$current_block
																								i32.const 21
																								i32.eq
																								if(result i32)
																									i32.const 25
																								else
																									i32.const 30
																								end
																							end
																						end
																					end
																					local.set $$t41
																					local.set $$t40
																					local.set $$t39
																					local.set $$t38
																					local.set $$t37.2
																					local.set $$t37.1
																					local.get $$t37.0
																					call $runtime.Block.Release
																					local.set $$t37.0
																					i32.const 11
																					local.set $$current_block
																					local.get $$t39
																					call $syscall$wasm4.SetDrawColorsU16
																					i32.const 8
																					local.get $$t40
																					i32.const 144
																					local.get $$t41
																					call $syscall$wasm4.RectI32
																					local.get $$t38
																					call $syscall$wasm4.SetDrawColorsU16
																					local.get $$t37.2
																					local.set $$t42
																					local.get $$t42
																					local.set $$t43
																					local.get $$t43
																					i32.const 8
																					i32.mul
																					local.set $$t44
																					local.get $$t44
																					i32.const 2
																					i32.div_s
																					local.set $$t45
																					i32.const 80
																					local.get $$t45
																					i32.sub
																					local.set $$t46
																					local.get $$t41
																					i32.const 2
																					i32.div_s
																					local.set $$t47
																					local.get $$t40
																					local.get $$t47
																					i32.add
																					local.set $$t48
																					local.get $$t48
																					i32.const 5
																					i32.sub
																					local.set $$t49
																					local.get $$t37.0
																					local.get $$t37.1
																					local.get $$t37.2
																					local.get $$t46
																					local.get $$t49
																					call $syscall$wasm4.Text
																					i32.const 10
																					local.set $$block_selector
																					br $$BlockDisp
																				end
																				i32.const 12
																				local.set $$current_block
																				i32.const 11
																				local.set $$block_selector
																				br $$BlockDisp
																			end
																			i32.const 13
																			local.set $$current_block
																			i32.const 11
																			local.set $$block_selector
																			br $$BlockDisp
																		end
																		i32.const 14
																		local.set $$current_block
																		local.get $this.0
																		call $runtime.Block.Retain
																		local.get $this.1
																		i32.const 16
																		i32.add
																		local.set $$t50.1
																		local.get $$t50.0
																		call $runtime.Block.Release
																		local.set $$t50.0
																		local.get $$t50.1
																		i32.load
																		local.set $$t51
																		local.get $$t51
																		i32.const 12
																		i32.lt_u
																		local.set $$t52
																		local.get $$t52
																		if
																			br $$Block_19
																		else
																			br $$Block_18
																		end
																	end
																	i32.const 15
																	local.set $$current_block
																	local.get $this.0
																	call $runtime.Block.Retain
																	local.get $this.1
																	i32.const 16
																	i32.add
																	local.set $$t53.1
																	local.get $$t53.0
																	call $runtime.Block.Release
																	local.set $$t53.0
																	local.get $$t53.1
																	i32.load
																	local.set $$t54
																	local.get $$t54
																	i32.const 118
																	i32.gt_u
																	local.set $$t55
																	br $$Block_15
																end
																local.get $$current_block
																i32.const 9
																i32.eq
																if(result i32)
																	i32.const 1
																else
																	local.get $$t55
																end
																local.set $$t56
																i32.const 16
																local.set $$current_block
																local.get $$t56
																if
																	i32.const 12
																	local.set $$block_selector
																	br $$BlockDisp
																else
																	i32.const 14
																	local.set $$block_selector
																	br $$BlockDisp
																end
															end
															i32.const 17
															local.set $$current_block
															i32.const 11
															local.set $$block_selector
															br $$BlockDisp
														end
														i32.const 18
														local.set $$current_block
														local.get $this.0
														call $runtime.Block.Retain
														local.get $this.1
														i32.const 16
														i32.add
														local.set $$t57.1
														local.get $$t57.0
														call $runtime.Block.Release
														local.set $$t57.0
														local.get $$t57.1
														i32.load
														local.set $$t58
														local.get $$t58
														i32.const 18
														i32.lt_u
														local.set $$t59
														local.get $$t59
														if
															br $$Block_23
														else
															br $$Block_22
														end
													end
													i32.const 19
													local.set $$current_block
													local.get $this.0
													call $runtime.Block.Retain
													local.get $this.1
													i32.const 16
													i32.add
													local.set $$t60.1
													local.get $$t60.0
													call $runtime.Block.Release
													local.set $$t60.0
													local.get $$t60.1
													i32.load
													local.set $$t61
													local.get $$t61
													i32.const 112
													i32.gt_u
													local.set $$t62
													br $$Block_19
												end
												local.get $$current_block
												i32.const 14
												i32.eq
												if(result i32)
													i32.const 1
												else
													local.get $$t62
												end
												local.set $$t63
												i32.const 20
												local.set $$current_block
												local.get $$t63
												if
													i32.const 13
													local.set $$block_selector
													br $$BlockDisp
												else
													i32.const 18
													local.set $$block_selector
													br $$BlockDisp
												end
											end
											i32.const 21
											local.set $$current_block
											i32.const 11
											local.set $$block_selector
											br $$BlockDisp
										end
										i32.const 22
										local.set $$current_block
										local.get $this.0
										call $runtime.Block.Retain
										local.get $this.1
										i32.const 16
										i32.add
										local.set $$t64.1
										local.get $$t64.0
										call $runtime.Block.Release
										local.set $$t64.0
										local.get $$t64.1
										i32.load
										local.set $$t65
										local.get $$t65
										i32.const 24
										i32.lt_u
										local.set $$t66
										local.get $$t66
										if
											br $$Block_26
										else
											br $$Block_25
										end
									end
									i32.const 23
									local.set $$current_block
									local.get $this.0
									call $runtime.Block.Retain
									local.get $this.1
									i32.const 16
									i32.add
									local.set $$t67.1
									local.get $$t67.0
									call $runtime.Block.Release
									local.set $$t67.0
									local.get $$t67.1
									i32.load
									local.set $$t68
									local.get $$t68
									i32.const 106
									i32.gt_u
									local.set $$t69
									br $$Block_23
								end
								local.get $$current_block
								i32.const 18
								i32.eq
								if(result i32)
									i32.const 1
								else
									local.get $$t69
								end
								local.set $$t70
								i32.const 24
								local.set $$current_block
								local.get $$t70
								if
									i32.const 17
									local.set $$block_selector
									br $$BlockDisp
								else
									i32.const 22
									local.set $$block_selector
									br $$BlockDisp
								end
							end
							i32.const 25
							local.set $$current_block
							i32.const 11
							local.set $$block_selector
							br $$BlockDisp
						end
						i32.const 26
						local.set $$current_block
						local.get $this.0
						call $runtime.Block.Retain
						local.get $this.1
						i32.const 16
						i32.add
						local.set $$t71.1
						local.get $$t71.0
						call $runtime.Block.Release
						local.set $$t71.0
						local.get $$t71.1
						i32.load
						local.set $$t72
						local.get $$t72
						i32.const 100
						i32.gt_u
						local.set $$t73
						br $$Block_26
					end
					local.get $$current_block
					i32.const 22
					i32.eq
					if(result i32)
						i32.const 1
					else
						local.get $$t73
					end
					local.set $$t74
					i32.const 27
					local.set $$current_block
					local.get $$t74
					if
						i32.const 21
						local.set $$block_selector
						br $$BlockDisp
					else
						i32.const 25
						local.set $$block_selector
						br $$BlockDisp
					end
				end
			end
		end
		local.get $$t0.0
		call $runtime.Block.Release
		local.get $$t1.0
		call $runtime.Block.Release
		local.get $$t3.0
		call $runtime.Block.Release
		local.get $$t4.0
		call $runtime.Block.Release
		local.get $$t6.0
		call $runtime.Block.Release
		local.get $$t7.0
		call $runtime.Block.Release
		local.get $$t9.0
		call $runtime.Block.Release
		local.get $$t10.0
		call $runtime.Block.Release
		local.get $$t12.0
		call $runtime.Block.Release
		local.get $$t13.0.0
		call $runtime.Block.Release
		local.get $$t15.0
		call $runtime.Block.Release
		local.get $$t16.0
		call $runtime.Block.Release
		local.get $$t17.0
		call $runtime.Block.Release
		local.get $$t18.0
		call $runtime.Block.Release
		local.get $$t24.0
		call $runtime.Block.Release
		local.get $$t25.0.0
		call $runtime.Block.Release
		local.get $$t34.0
		call $runtime.Block.Release
		local.get $$t37.0
		call $runtime.Block.Release
		local.get $$t50.0
		call $runtime.Block.Release
		local.get $$t53.0
		call $runtime.Block.Release
		local.get $$t57.0
		call $runtime.Block.Release
		local.get $$t60.0
		call $runtime.Block.Release
		local.get $$t64.0
		call $runtime.Block.Release
		local.get $$t67.0
		call $runtime.Block.Release
		local.get $$t71.0
		call $runtime.Block.Release
	)
	(func $w42048.UI.sound (param $this.0 i32) (param $this.1 i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$t0 i32)
		(local $$t1 i32)
		(local $$t2.0 i32)
		(local $$t2.1 i32)
		(local $$t2.2 i32)
		(local $$t2.3 i32)
		(local $$t2.4 i32)
		(local $$t2.5 i32)
		(local $$t2.6 i32)
		(local $$t2.7 i32)
		(local $$t2.8 i32)
		(local $$t3 i32)
		(local $$t4.0 i32)
		(local $$t4.1 i32)
		(local $$t4.2 i32)
		(local $$t4.3 i32)
		(local $$t4.4 i32)
		(local $$t4.5 i32)
		(local $$t4.6 i32)
		(local $$t4.7 i32)
		(local $$t4.8 i32)
		(local $$t5 i32)
		(local $$t6.0 i32)
		(local $$t6.1 i32)
		(local $$t6.2 i32)
		(local $$t6.3 i32)
		(local $$t6.4 i32)
		(local $$t6.5 i32)
		(local $$t6.6 i32)
		(local $$t6.7 i32)
		(local $$t6.8 i32)
		(local $$t7 i32)
		(local $$t8.0 i32)
		(local $$t8.1 i32)
		(local $$t8.2 i32)
		(local $$t8.3 i32)
		(local $$t8.4 i32)
		(local $$t8.5 i32)
		(local $$t8.6 i32)
		(local $$t8.7 i32)
		(local $$t8.8 i32)
		(local $$t9 i32)
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
																	i32.const 1
																	if
																		br $$Block_1
																	else
																		br $$Block_0
																	end
																end
																i32.const 1
																local.set $$current_block
																br $$BlockFnBody
															end
															i32.const 2
															local.set $$current_block
															local.get $this.0
															local.get $this.1
															call $w42048.UI.btn1
															local.set $$t0
															local.get $$t0
															if
																br $$Block_2
															else
																br $$Block_3
															end
														end
														i32.const 3
														local.set $$current_block
														br $$BlockFnBody
													end
													i32.const 4
													local.set $$current_block
													local.get $this.0
													local.get $this.1
													call $w42048.UI.btn2
													local.set $$t1
													local.get $$t1
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
												i32.const 31228
												i32.load
												i32.const 31228
												i32.load offset=4
												i32.const 31228
												i32.load offset=8
												i32.const 31228
												i32.load offset=12
												i32.const 31228
												i32.load offset=16
												i32.const 31228
												i32.load offset=20
												i32.const 31228
												i32.load offset=24
												i32.const 31228
												i32.load offset=28
												i32.const 31228
												i32.load offset=32
												local.set $$t2.8
												local.set $$t2.7
												local.set $$t2.6
												local.set $$t2.5
												local.set $$t2.4
												local.set $$t2.3
												local.set $$t2.2
												local.set $$t2.1
												local.set $$t2.0
												local.get $$t2.0
												local.get $$t2.1
												local.get $$t2.2
												local.get $$t2.3
												local.get $$t2.4
												local.get $$t2.5
												local.get $$t2.6
												local.get $$t2.7
												local.get $$t2.8
												call $w42048.play
												i32.const 3
												local.set $$block_selector
												br $$BlockDisp
											end
											i32.const 6
											local.set $$current_block
											local.get $this.0
											local.get $this.1
											call $w42048.UI.up
											local.set $$t3
											local.get $$t3
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
										i32.const 31120
										i32.load
										i32.const 31120
										i32.load offset=4
										i32.const 31120
										i32.load offset=8
										i32.const 31120
										i32.load offset=12
										i32.const 31120
										i32.load offset=16
										i32.const 31120
										i32.load offset=20
										i32.const 31120
										i32.load offset=24
										i32.const 31120
										i32.load offset=28
										i32.const 31120
										i32.load offset=32
										local.set $$t4.8
										local.set $$t4.7
										local.set $$t4.6
										local.set $$t4.5
										local.set $$t4.4
										local.set $$t4.3
										local.set $$t4.2
										local.set $$t4.1
										local.set $$t4.0
										local.get $$t4.0
										local.get $$t4.1
										local.get $$t4.2
										local.get $$t4.3
										local.get $$t4.4
										local.get $$t4.5
										local.get $$t4.6
										local.get $$t4.7
										local.get $$t4.8
										call $w42048.play
										i32.const 3
										local.set $$block_selector
										br $$BlockDisp
									end
									i32.const 8
									local.set $$current_block
									local.get $this.0
									local.get $this.1
									call $w42048.UI.down
									local.set $$t5
									local.get $$t5
									if
										i32.const 7
										local.set $$block_selector
										br $$BlockDisp
									else
										br $$Block_9
									end
								end
								i32.const 9
								local.set $$current_block
								i32.const 31192
								i32.load
								i32.const 31192
								i32.load offset=4
								i32.const 31192
								i32.load offset=8
								i32.const 31192
								i32.load offset=12
								i32.const 31192
								i32.load offset=16
								i32.const 31192
								i32.load offset=20
								i32.const 31192
								i32.load offset=24
								i32.const 31192
								i32.load offset=28
								i32.const 31192
								i32.load offset=32
								local.set $$t6.8
								local.set $$t6.7
								local.set $$t6.6
								local.set $$t6.5
								local.set $$t6.4
								local.set $$t6.3
								local.set $$t6.2
								local.set $$t6.1
								local.set $$t6.0
								local.get $$t6.0
								local.get $$t6.1
								local.get $$t6.2
								local.get $$t6.3
								local.get $$t6.4
								local.get $$t6.5
								local.get $$t6.6
								local.get $$t6.7
								local.get $$t6.8
								call $w42048.play
								i32.const 3
								local.set $$block_selector
								br $$BlockDisp
							end
							i32.const 10
							local.set $$current_block
							local.get $this.0
							local.get $this.1
							call $w42048.UI.right
							local.set $$t7
							local.get $$t7
							if
								i32.const 9
								local.set $$block_selector
								br $$BlockDisp
							else
								br $$Block_11
							end
						end
						i32.const 11
						local.set $$current_block
						i32.const 31156
						i32.load
						i32.const 31156
						i32.load offset=4
						i32.const 31156
						i32.load offset=8
						i32.const 31156
						i32.load offset=12
						i32.const 31156
						i32.load offset=16
						i32.const 31156
						i32.load offset=20
						i32.const 31156
						i32.load offset=24
						i32.const 31156
						i32.load offset=28
						i32.const 31156
						i32.load offset=32
						local.set $$t8.8
						local.set $$t8.7
						local.set $$t8.6
						local.set $$t8.5
						local.set $$t8.4
						local.set $$t8.3
						local.set $$t8.2
						local.set $$t8.1
						local.set $$t8.0
						local.get $$t8.0
						local.get $$t8.1
						local.get $$t8.2
						local.get $$t8.3
						local.get $$t8.4
						local.get $$t8.5
						local.get $$t8.6
						local.get $$t8.7
						local.get $$t8.8
						call $w42048.play
						i32.const 3
						local.set $$block_selector
						br $$BlockDisp
					end
					i32.const 12
					local.set $$current_block
					local.get $this.0
					local.get $this.1
					call $w42048.UI.left
					local.set $$t9
					local.get $$t9
					if
						i32.const 11
						local.set $$block_selector
						br $$BlockDisp
					else
						i32.const 3
						local.set $$block_selector
						br $$BlockDisp
					end
				end
			end
		end
	)
	(func $w42048.UI.up (param $this.0 i32) (param $this.1 i32) (result i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$ret_0 i32)
		(local $$t0 i32)
		(local $$t1 i32)
		(local $$t2 i32)
		(local $$t3.0 i32)
		(local $$t3.1 i32)
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
							local.get $this.0
							local.get $this.1
							call $w42048.UI.GetPad
							local.set $$t0
							local.get $$t0
							i32.const 64
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
						local.get $this.0
						call $runtime.Block.Retain
						local.get $this.1
						i32.const 20
						i32.add
						local.set $$t3.1
						local.get $$t3.0
						call $runtime.Block.Release
						local.set $$t3.0
						local.get $$t3.1
						i32.load8_u align=1
						local.set $$t4
						local.get $$t4
						i32.const 64
						i32.and
						local.set $$t5
						local.get $$t5
						i32.const 0
						i32.eq
						local.set $$t6
						br $$Block_1
					end
					local.get $$current_block
					i32.const 0
					i32.eq
					if(result i32)
						i32.const 0
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
		local.get $$t3.0
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
	(func $strconv.NumError.Error (param $this.0 i32) (param $this.1 i32) (result i32 i32 i32)
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
		(local $$t2.0 i32)
		(local $$t2.1 i32)
		(local $$t2.2 i32)
		(local $$t3.0 i32)
		(local $$t3.1 i32)
		(local $$t3.2 i32)
		(local $$t4.0 i32)
		(local $$t4.1 i32)
		(local $$t4.2 i32)
		(local $$t5.0 i32)
		(local $$t5.1 i32)
		(local $$t6.0 i32)
		(local $$t6.1 i32)
		(local $$t6.2 i32)
		(local $$t7.0 i32)
		(local $$t7.1 i32)
		(local $$t7.2 i32)
		(local $$t8.0 i32)
		(local $$t8.1 i32)
		(local $$t8.2 i32)
		(local $$t9.0 i32)
		(local $$t9.1 i32)
		(local $$t9.2 i32)
		(local $$t10.0 i32)
		(local $$t10.1 i32)
		(local $$t11.0.0 i32)
		(local $$t11.0.1 i32)
		(local $$t11.1 i32)
		(local $$t11.2 i32)
		(local $$t12.0 i32)
		(local $$t12.1 i32)
		(local $$t12.2 i32)
		(local $$t13.0 i32)
		(local $$t13.1 i32)
		(local $$t13.2 i32)
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
					i32.const 0
					i32.const 35032
					i32.const 8
					local.get $$t1.0
					local.get $$t1.1
					local.get $$t1.2
					call $$string.appendstr
					local.set $$t2.2
					local.set $$t2.1
					local.get $$t2.0
					call $runtime.Block.Release
					local.set $$t2.0
					local.get $$t2.0
					local.get $$t2.1
					local.get $$t2.2
					i32.const 0
					i32.const 31477
					i32.const 2
					call $$string.appendstr
					local.set $$t3.2
					local.set $$t3.1
					local.get $$t3.0
					call $runtime.Block.Release
					local.set $$t3.0
					local.get $$t3.0
					local.get $$t3.1
					local.get $$t3.2
					i32.const 0
					i32.const 35040
					i32.const 8
					call $$string.appendstr
					local.set $$t4.2
					local.set $$t4.1
					local.get $$t4.0
					call $runtime.Block.Release
					local.set $$t4.0
					local.get $this.0
					call $runtime.Block.Retain
					local.get $this.1
					i32.const 12
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
					local.set $$t6.2
					local.set $$t6.1
					local.get $$t6.0
					call $runtime.Block.Release
					local.set $$t6.0
					local.get $$t6.0
					local.get $$t6.1
					local.get $$t6.2
					call $strconv.Quote
					local.set $$t7.2
					local.set $$t7.1
					local.get $$t7.0
					call $runtime.Block.Release
					local.set $$t7.0
					local.get $$t4.0
					local.get $$t4.1
					local.get $$t4.2
					local.get $$t7.0
					local.get $$t7.1
					local.get $$t7.2
					call $$string.appendstr
					local.set $$t8.2
					local.set $$t8.1
					local.get $$t8.0
					call $runtime.Block.Release
					local.set $$t8.0
					local.get $$t8.0
					local.get $$t8.1
					local.get $$t8.2
					i32.const 0
					i32.const 31477
					i32.const 2
					call $$string.appendstr
					local.set $$t9.2
					local.set $$t9.1
					local.get $$t9.0
					call $runtime.Block.Release
					local.set $$t9.0
					local.get $this.0
					call $runtime.Block.Retain
					local.get $this.1
					i32.const 24
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
					local.set $$t11.2
					local.set $$t11.1
					local.set $$t11.0.1
					local.get $$t11.0.0
					call $runtime.Block.Release
					local.set $$t11.0.0
					local.get $$t11.0.0
					local.get $$t11.0.1
					local.get $$t11.1
					i32.load offset=8
					call_indirect 0 (type $$$fnSig2)
					local.set $$t12.2
					local.set $$t12.1
					local.get $$t12.0
					call $runtime.Block.Release
					local.set $$t12.0
					local.get $$t9.0
					local.get $$t9.1
					local.get $$t9.2
					local.get $$t12.0
					local.get $$t12.1
					local.get $$t12.2
					call $$string.appendstr
					local.set $$t13.2
					local.set $$t13.1
					local.get $$t13.0
					call $runtime.Block.Release
					local.set $$t13.0
					local.get $$t13.0
					call $runtime.Block.Retain
					local.get $$t13.1
					local.get $$t13.2
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
		local.get $$t11.0.0
		call $runtime.Block.Release
		local.get $$t12.0
		call $runtime.Block.Release
		local.get $$t13.0
		call $runtime.Block.Release
	)
	(func $w42048$game.board.Add (param $this.0 i32) (param $this.1 i32)
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
		(local $$t11.0 i32)
		(local $$t11.1 i32)
		(local $$t11.2 i32)
		(local $$t11.3 i32)
		(local $$t12.0 i32)
		(local $$t12.1 i32)
		(local $$t13.0 i32)
		(local $$t13.1 i32)
		(local $$t13.2 i32)
		(local $$t13.3 i32)
		(local $$t14 i32)
		(local $$t15.0 i32)
		(local $$t15.1 i32)
		(local $$t16 i32)
		(local $$t17 i32)
		(local $$t18 i32)
		(local $$t19 i32)
		(local $$t20 i32)
		(local $$t21 i32)
		(local $$t22 i32)
		(local $$t23 i32)
		(local $$t24 i32)
		(local $$t25 i32)
		(local $$t26 i32)
		(local $$t27.0 i32)
		(local $$t27.1 i32)
		(local $$t28.0 i32)
		(local $$t28.1 i32)
		(local $$t28.2 i32)
		(local $$t28.3 i32)
		(local $$t29.0 i32)
		(local $$t29.1 i32)
		(local $$t30.0 i32)
		(local $$t30.1 i32)
		(local $$t30.2 i32)
		(local $$t30.3 i32)
		(local $$t31 i32)
		(local $$t32.0 i32)
		(local $$t32.1 i32)
		(local $$t33 i32)
		(local $$t34 i32)
		(local $$t35 i32)
		(local $$t36 i32)
		(local $$t37 i32)
		(local $$t38 i32)
		(local $$t39 i32)
		(local $$t40.0 i32)
		(local $$t40.1 i32)
		(local $$t41.0 i32)
		(local $$t41.1 i32)
		(local $$t42.0 i32)
		(local $$t42.1 i32)
		(local $$t43.0 i32)
		(local $$t43.1 i32)
		(local $$t43.2 i32)
		(local $$t43.3 i32)
		(local $$t44.0 i32)
		(local $$t44.1 i32)
		(local $$t45.0 i32)
		(local $$t45.1 i32)
		(local $$t45.2 i32)
		(local $$t45.3 i32)
		(local $$t46.0 i32)
		(local $$t46.1 i32)
		block $$BlockFnBody
			loop $$BlockDisp
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
																									br_table  0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 0
																								end
																								i32.const 0
																								local.set $$current_block
																								local.get $this.0
																								local.get $this.1
																								i32.const 10
																								call $w42048$game.board.intn
																								local.set $$t0
																								local.get $$t0
																								i32.const 8
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
																							br $$Block_1
																						end
																						local.get $$current_block
																						i32.const 0
																						i32.eq
																						if(result i32)
																							i32.const 2
																						else
																							i32.const 4
																						end
																						local.set $$t2
																						i32.const 2
																						local.set $$current_block
																						br $$Block_4
																					end
																					i32.const 3
																					local.set $$current_block
																					br $$Block_7
																				end
																				i32.const 4
																				local.set $$current_block
																				local.get $this.0
																				local.get $this.1
																				local.get $$t3
																				call $w42048$game.board.intn
																				local.set $$t4
																				local.get $$t4
																				i32.const 1
																				i32.add
																				local.set $$t5
																				br $$Block_12
																			end
																			local.get $$current_block
																			i32.const 2
																			i32.eq
																			if(result i32)
																				i32.const 0
																			else
																				local.get $$t6
																			end
																			local.get $$current_block
																			i32.const 2
																			i32.eq
																			if(result i32)
																				i32.const 0
																			else
																				local.get $$t7
																			end
																			local.set $$t8
																			local.set $$t3
																			i32.const 5
																			local.set $$current_block
																			local.get $$t8
																			i32.const 4
																			i32.lt_s
																			local.set $$t9
																			local.get $$t9
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
																		i32.const 12
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
																		local.get $$t11.0
																		call $runtime.Block.Retain
																		local.get $$t11.1
																		i32.const 16
																		local.get $$t8
																		i32.mul
																		i32.add
																		local.set $$t12.1
																		local.get $$t12.0
																		call $runtime.Block.Release
																		local.set $$t12.0
																		local.get $$t12.1
																		i32.load
																		call $runtime.Block.Retain
																		local.get $$t12.1
																		i32.load offset=4
																		local.get $$t12.1
																		i32.load offset=8
																		local.get $$t12.1
																		i32.load offset=12
																		local.set $$t13.3
																		local.set $$t13.2
																		local.set $$t13.1
																		local.get $$t13.0
																		call $runtime.Block.Release
																		local.set $$t13.0
																		local.get $$t13.0
																		call $runtime.Block.Retain
																		local.get $$t13.1
																		i32.const 4
																		local.get $$t14
																		i32.mul
																		i32.add
																		local.set $$t15.1
																		local.get $$t15.0
																		call $runtime.Block.Release
																		local.set $$t15.0
																		local.get $$t15.1
																		i32.load
																		local.set $$t16
																		local.get $$t16
																		i32.const 0
																		i32.eq
																		local.set $$t17
																		local.get $$t17
																		if
																			br $$Block_8
																		else
																			br $$Block_9
																		end
																	end
																	i32.const 7
																	local.set $$current_block
																	local.get $$t8
																	i32.const 1
																	i32.add
																	local.set $$t7
																	i32.const 5
																	local.set $$block_selector
																	br $$BlockDisp
																end
																local.get $$current_block
																i32.const 3
																i32.eq
																if(result i32)
																	local.get $$t3
																else
																	local.get $$t18
																end
																local.get $$current_block
																i32.const 3
																i32.eq
																if(result i32)
																	i32.const 0
																else
																	local.get $$t19
																end
																local.set $$t14
																local.set $$t6
																i32.const 8
																local.set $$current_block
																local.get $$t14
																i32.const 4
																i32.lt_s
																local.set $$t20
																local.get $$t20
																if
																	i32.const 6
																	local.set $$block_selector
																	br $$BlockDisp
																else
																	i32.const 7
																	local.set $$block_selector
																	br $$BlockDisp
																end
															end
															i32.const 9
															local.set $$current_block
															local.get $$t6
															i32.const 1
															i32.add
															local.set $$t21
															br $$Block_9
														end
														local.get $$current_block
														i32.const 6
														i32.eq
														if(result i32)
															local.get $$t6
														else
															local.get $$t21
														end
														local.set $$t18
														i32.const 10
														local.set $$current_block
														local.get $$t14
														i32.const 1
														i32.add
														local.set $$t19
														i32.const 8
														local.set $$block_selector
														br $$BlockDisp
													end
													i32.const 11
													local.set $$current_block
													br $$Block_15
												end
												i32.const 12
												local.set $$current_block
												br $$BlockFnBody
											end
											local.get $$current_block
											i32.const 4
											i32.eq
											if(result i32)
												i32.const 0
											else
												local.get $$t22
											end
											local.get $$current_block
											i32.const 4
											i32.eq
											if(result i32)
												i32.const 0
											else
												local.get $$t24
											end
											local.set $$t25
											local.set $$t23
											i32.const 13
											local.set $$current_block
											local.get $$t25
											i32.const 4
											i32.lt_s
											local.set $$t26
											local.get $$t26
											if
												i32.const 11
												local.set $$block_selector
												br $$BlockDisp
											else
												i32.const 12
												local.set $$block_selector
												br $$BlockDisp
											end
										end
										i32.const 14
										local.set $$current_block
										local.get $this.0
										call $runtime.Block.Retain
										local.get $this.1
										i32.const 12
										i32.add
										local.set $$t27.1
										local.get $$t27.0
										call $runtime.Block.Release
										local.set $$t27.0
										local.get $$t27.1
										i32.load
										call $runtime.Block.Retain
										local.get $$t27.1
										i32.load offset=4
										local.get $$t27.1
										i32.load offset=8
										local.get $$t27.1
										i32.load offset=12
										local.set $$t28.3
										local.set $$t28.2
										local.set $$t28.1
										local.get $$t28.0
										call $runtime.Block.Release
										local.set $$t28.0
										local.get $$t28.0
										call $runtime.Block.Retain
										local.get $$t28.1
										i32.const 16
										local.get $$t25
										i32.mul
										i32.add
										local.set $$t29.1
										local.get $$t29.0
										call $runtime.Block.Release
										local.set $$t29.0
										local.get $$t29.1
										i32.load
										call $runtime.Block.Retain
										local.get $$t29.1
										i32.load offset=4
										local.get $$t29.1
										i32.load offset=8
										local.get $$t29.1
										i32.load offset=12
										local.set $$t30.3
										local.set $$t30.2
										local.set $$t30.1
										local.get $$t30.0
										call $runtime.Block.Release
										local.set $$t30.0
										local.get $$t30.0
										call $runtime.Block.Retain
										local.get $$t30.1
										i32.const 4
										local.get $$t31
										i32.mul
										i32.add
										local.set $$t32.1
										local.get $$t32.0
										call $runtime.Block.Release
										local.set $$t32.0
										local.get $$t32.1
										i32.load
										local.set $$t33
										local.get $$t33
										i32.const 0
										i32.eq
										local.set $$t34
										local.get $$t34
										if
											br $$Block_16
										else
											br $$Block_17
										end
									end
									i32.const 15
									local.set $$current_block
									local.get $$t25
									i32.const 1
									i32.add
									local.set $$t24
									i32.const 13
									local.set $$block_selector
									br $$BlockDisp
								end
								local.get $$current_block
								i32.const 11
								i32.eq
								if(result i32)
									local.get $$t23
								else
									local.get $$t35
								end
								local.get $$current_block
								i32.const 11
								i32.eq
								if(result i32)
									i32.const 0
								else
									local.get $$t36
								end
								local.set $$t31
								local.set $$t22
								i32.const 16
								local.set $$current_block
								local.get $$t31
								i32.const 4
								i32.lt_s
								local.set $$t37
								local.get $$t37
								if
									i32.const 14
									local.set $$block_selector
									br $$BlockDisp
								else
									i32.const 15
									local.set $$block_selector
									br $$BlockDisp
								end
							end
							i32.const 17
							local.set $$current_block
							local.get $$t22
							i32.const 1
							i32.add
							local.set $$t38
							local.get $$t38
							local.get $$t5
							i32.eq
							local.set $$t39
							local.get $$t39
							if
								br $$Block_18
							else
								br $$Block_17
							end
						end
						local.get $$current_block
						i32.const 14
						i32.eq
						if(result i32)
							local.get $$t22
						else
							local.get $$t38
						end
						local.set $$t35
						i32.const 18
						local.set $$current_block
						local.get $$t31
						i32.const 1
						i32.add
						local.set $$t36
						i32.const 16
						local.set $$block_selector
						br $$BlockDisp
					end
					i32.const 19
					local.set $$current_block
					local.get $this.0
					call $runtime.Block.Retain
					local.get $this.1
					i32.const 32
					i32.add
					local.set $$t40.1
					local.get $$t40.0
					call $runtime.Block.Release
					local.set $$t40.0
					local.get $$t40.1
					local.get $$t25
					i32.store
					local.get $this.0
					call $runtime.Block.Retain
					local.get $this.1
					i32.const 36
					i32.add
					local.set $$t41.1
					local.get $$t41.0
					call $runtime.Block.Release
					local.set $$t41.0
					local.get $$t41.1
					local.get $$t31
					i32.store
					local.get $this.0
					call $runtime.Block.Retain
					local.get $this.1
					i32.const 12
					i32.add
					local.set $$t42.1
					local.get $$t42.0
					call $runtime.Block.Release
					local.set $$t42.0
					local.get $$t42.1
					i32.load
					call $runtime.Block.Retain
					local.get $$t42.1
					i32.load offset=4
					local.get $$t42.1
					i32.load offset=8
					local.get $$t42.1
					i32.load offset=12
					local.set $$t43.3
					local.set $$t43.2
					local.set $$t43.1
					local.get $$t43.0
					call $runtime.Block.Release
					local.set $$t43.0
					local.get $$t43.0
					call $runtime.Block.Retain
					local.get $$t43.1
					i32.const 16
					local.get $$t25
					i32.mul
					i32.add
					local.set $$t44.1
					local.get $$t44.0
					call $runtime.Block.Release
					local.set $$t44.0
					local.get $$t44.1
					i32.load
					call $runtime.Block.Retain
					local.get $$t44.1
					i32.load offset=4
					local.get $$t44.1
					i32.load offset=8
					local.get $$t44.1
					i32.load offset=12
					local.set $$t45.3
					local.set $$t45.2
					local.set $$t45.1
					local.get $$t45.0
					call $runtime.Block.Release
					local.set $$t45.0
					local.get $$t45.0
					call $runtime.Block.Retain
					local.get $$t45.1
					i32.const 4
					local.get $$t31
					i32.mul
					i32.add
					local.set $$t46.1
					local.get $$t46.0
					call $runtime.Block.Release
					local.set $$t46.0
					local.get $$t46.1
					local.get $$t2
					i32.store
					br $$BlockFnBody
				end
			end
		end
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
		local.get $$t27.0
		call $runtime.Block.Release
		local.get $$t28.0
		call $runtime.Block.Release
		local.get $$t29.0
		call $runtime.Block.Release
		local.get $$t30.0
		call $runtime.Block.Release
		local.get $$t32.0
		call $runtime.Block.Release
		local.get $$t40.0
		call $runtime.Block.Release
		local.get $$t41.0
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
	)
	(func $w42048$game.board.Get (param $this.0 i32) (param $this.1 i32) (param $col i32) (param $row i32) (result i32)
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
		(local $$t3.2 i32)
		(local $$t3.3 i32)
		(local $$t4.0 i32)
		(local $$t4.1 i32)
		(local $$t5 i32)
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
					i32.const 12
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
					i32.const 16
					local.get $row
					i32.mul
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
					local.get $$t3.0
					call $runtime.Block.Retain
					local.get $$t3.1
					i32.const 4
					local.get $col
					i32.mul
					i32.add
					local.set $$t4.1
					local.get $$t4.0
					call $runtime.Block.Release
					local.set $$t4.0
					local.get $$t4.1
					i32.load
					local.set $$t5
					local.get $$t5
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
		local.get $$t2.0
		call $runtime.Block.Release
		local.get $$t3.0
		call $runtime.Block.Release
		local.get $$t4.0
		call $runtime.Block.Release
	)
	(func $w42048$game.board.Input (param $this.0 i32) (param $this.1 i32) (param $key i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$t0 i32)
		(local $$t1 i32)
		(local $$t2 i32)
		(local $$t3 i32)
		(local $$t4 i32)
		(local $$t5 i32)
		block $$BlockFnBody
			loop $$BlockDisp
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
																	br_table  0 1 2 3 4 5 6 7 8 9 10 11 0
																end
																i32.const 0
																local.set $$current_block
																local.get $key
																i32.const 0
																i32.eq
																local.set $$t0
																local.get $$t0
																if
																	br $$Block_1
																else
																	br $$Block_3
																end
															end
															i32.const 1
															local.set $$current_block
															br $$BlockFnBody
														end
														i32.const 2
														local.set $$current_block
														local.get $this.0
														local.get $this.1
														call $w42048$game.board.up
														local.get $this.0
														local.get $this.1
														call $w42048$game.board.Add
														i32.const 1
														local.set $$block_selector
														br $$BlockDisp
													end
													i32.const 3
													local.set $$current_block
													local.get $this.0
													local.get $this.1
													call $w42048$game.board.down
													local.get $this.0
													local.get $this.1
													call $w42048$game.board.Add
													i32.const 1
													local.set $$block_selector
													br $$BlockDisp
												end
												i32.const 4
												local.set $$current_block
												local.get $key
												i32.const 1
												i32.eq
												local.set $$t1
												local.get $$t1
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
											local.get $this.0
											local.get $this.1
											call $w42048$game.board.right
											local.get $this.0
											local.get $this.1
											call $w42048$game.board.Add
											i32.const 1
											local.set $$block_selector
											br $$BlockDisp
										end
										i32.const 6
										local.set $$current_block
										local.get $key
										i32.const 2
										i32.eq
										local.set $$t2
										local.get $$t2
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
									local.get $this.0
									local.get $this.1
									call $w42048$game.board.left
									local.get $this.0
									local.get $this.1
									call $w42048$game.board.Add
									i32.const 1
									local.set $$block_selector
									br $$BlockDisp
								end
								i32.const 8
								local.set $$current_block
								local.get $key
								i32.const 3
								i32.eq
								local.set $$t3
								local.get $$t3
								if
									i32.const 7
									local.set $$block_selector
									br $$BlockDisp
								else
									br $$Block_8
								end
							end
							i32.const 9
							local.set $$current_block
							local.get $key
							i32.const 4
							i32.eq
							local.set $$t4
							local.get $$t4
							if
								i32.const 1
								local.set $$block_selector
								br $$BlockDisp
							else
								br $$Block_10
							end
						end
						i32.const 10
						local.set $$current_block
						local.get $this.0
						local.get $this.1
						call $w42048$game.board.Restart
						i32.const 1
						local.set $$block_selector
						br $$BlockDisp
					end
					i32.const 11
					local.set $$current_block
					local.get $key
					i32.const 5
					i32.eq
					local.set $$t5
					local.get $$t5
					if
						i32.const 10
						local.set $$block_selector
						br $$BlockDisp
					else
						i32.const 1
						local.set $$block_selector
						br $$BlockDisp
					end
				end
			end
		end
	)
	(func $w42048$game.board.IsOver (param $this.0 i32) (param $this.1 i32) (result i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$ret_0 i32)
		(local $$t0 i32)
		(local $$t1 i32)
		(local $$t2 i32)
		(local $$t3 i32)
		(local $$t4 i32)
		(local $$t5 i32)
		(local $$t6.0 i32)
		(local $$t6.1 i32)
		(local $$t7.0 i32)
		(local $$t7.1 i32)
		(local $$t7.2 i32)
		(local $$t7.3 i32)
		(local $$t8.0 i32)
		(local $$t8.1 i32)
		(local $$t9.0 i32)
		(local $$t9.1 i32)
		(local $$t9.2 i32)
		(local $$t9.3 i32)
		(local $$t10 i32)
		(local $$t11.0 i32)
		(local $$t11.1 i32)
		(local $$t12 i32)
		(local $$t13 i32)
		(local $$t14 i32)
		(local $$t15 i32)
		(local $$t16 i32)
		(local $$t17 i32)
		(local $$t18.0 i32)
		(local $$t18.1 i32)
		(local $$t19 i32)
		(local $$t20 i32)
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
															br $$Block_2
														end
														i32.const 1
														local.set $$current_block
														br $$Block_5
													end
													i32.const 2
													local.set $$current_block
													local.get $$t0
													i32.const 0
													i32.eq
													local.set $$t1
													local.get $$t1
													if
														br $$Block_9
													else
														br $$Block_8
													end
												end
												local.get $$current_block
												i32.const 0
												i32.eq
												if(result i32)
													i32.const 0
												else
													local.get $$t2
												end
												local.get $$current_block
												i32.const 0
												i32.eq
												if(result i32)
													i32.const 0
												else
													local.get $$t3
												end
												local.set $$t4
												local.set $$t0
												i32.const 3
												local.set $$current_block
												local.get $$t4
												i32.const 4
												i32.lt_s
												local.set $$t5
												local.get $$t5
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
											local.get $this.0
											call $runtime.Block.Retain
											local.get $this.1
											i32.const 12
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
											i32.const 16
											local.get $$t4
											i32.mul
											i32.add
											local.set $$t8.1
											local.get $$t8.0
											call $runtime.Block.Release
											local.set $$t8.0
											local.get $$t8.1
											i32.load
											call $runtime.Block.Retain
											local.get $$t8.1
											i32.load offset=4
											local.get $$t8.1
											i32.load offset=8
											local.get $$t8.1
											i32.load offset=12
											local.set $$t9.3
											local.set $$t9.2
											local.set $$t9.1
											local.get $$t9.0
											call $runtime.Block.Release
											local.set $$t9.0
											local.get $$t9.0
											call $runtime.Block.Retain
											local.get $$t9.1
											i32.const 4
											local.get $$t10
											i32.mul
											i32.add
											local.set $$t11.1
											local.get $$t11.0
											call $runtime.Block.Release
											local.set $$t11.0
											local.get $$t11.1
											i32.load
											local.set $$t12
											local.get $$t12
											i32.const 0
											i32.eq
											local.set $$t13
											local.get $$t13
											if
												br $$Block_6
											else
												br $$Block_7
											end
										end
										i32.const 5
										local.set $$current_block
										local.get $$t4
										i32.const 1
										i32.add
										local.set $$t3
										i32.const 3
										local.set $$block_selector
										br $$BlockDisp
									end
									local.get $$current_block
									i32.const 1
									i32.eq
									if(result i32)
										local.get $$t0
									else
										local.get $$t14
									end
									local.get $$current_block
									i32.const 1
									i32.eq
									if(result i32)
										i32.const 0
									else
										local.get $$t15
									end
									local.set $$t10
									local.set $$t2
									i32.const 6
									local.set $$current_block
									local.get $$t10
									i32.const 4
									i32.lt_s
									local.set $$t16
									local.get $$t16
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
								local.get $$t2
								i32.const 1
								i32.add
								local.set $$t17
								br $$Block_7
							end
							local.get $$current_block
							i32.const 4
							i32.eq
							if(result i32)
								local.get $$t2
							else
								local.get $$t17
							end
							local.set $$t14
							i32.const 8
							local.set $$current_block
							local.get $$t10
							i32.const 1
							i32.add
							local.set $$t15
							i32.const 6
							local.set $$block_selector
							br $$BlockDisp
						end
						i32.const 9
						local.set $$current_block
						local.get $this.0
						call $runtime.Block.Retain
						local.get $this.1
						i32.const 28
						i32.add
						local.set $$t18.1
						local.get $$t18.0
						call $runtime.Block.Release
						local.set $$t18.0
						local.get $$t18.1
						i32.load8_u align=1
						local.set $$t19
						br $$Block_9
					end
					local.get $$current_block
					i32.const 2
					i32.eq
					if(result i32)
						i32.const 1
					else
						local.get $$t19
					end
					local.set $$t20
					i32.const 10
					local.set $$current_block
					local.get $$t20
					local.set $$ret_0
					br $$BlockFnBody
				end
			end
		end
		local.get $$ret_0
		local.get $$t6.0
		call $runtime.Block.Release
		local.get $$t7.0
		call $runtime.Block.Release
		local.get $$t8.0
		call $runtime.Block.Release
		local.get $$t9.0
		call $runtime.Block.Release
		local.get $$t11.0
		call $runtime.Block.Release
		local.get $$t18.0
		call $runtime.Block.Release
	)
	(func $w42048$game.board.Restart (param $this.0 i32) (param $this.1 i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$t0.0 i32)
		(local $$t0.1 i32)
		(local $$t1.0 i32)
		(local $$t1.1 i32)
		(local $$t1.2 i32)
		(local $$t1.3 i32)
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
					i32.const 12
					i32.add
					local.set $$t0.1
					local.get $$t0.0
					call $runtime.Block.Release
					local.set $$t0.0
					call $w42048$game.newMatrix
					local.set $$t1.3
					local.set $$t1.2
					local.set $$t1.1
					local.get $$t1.0
					call $runtime.Block.Release
					local.set $$t1.0
					local.get $$t0.1
					local.get $$t1.0
					call $runtime.Block.Retain
					local.get $$t0.1
					i32.load align=1
					call $runtime.Block.Release
					i32.store align=1
					local.get $$t0.1
					local.get $$t1.1
					i32.store offset=4
					local.get $$t0.1
					local.get $$t1.2
					i32.store offset=8
					local.get $$t0.1
					local.get $$t1.3
					i32.store offset=12
					local.get $this.0
					local.get $this.1
					call $w42048$game.board.Add
					local.get $this.0
					local.get $this.1
					call $w42048$game.board.Add
					br $$BlockFnBody
				end
			end
		end
		local.get $$t0.0
		call $runtime.Block.Release
		local.get $$t1.0
		call $runtime.Block.Release
	)
	(func $w42048$game.board.Total (param $this.0 i32) (param $this.1 i32) (result i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$ret_0 i32)
		(local $$t0 i32)
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
		(local $$t8.0 i32)
		(local $$t8.1 i32)
		(local $$t8.2 i32)
		(local $$t8.3 i32)
		(local $$t9 i32)
		(local $$t10.0 i32)
		(local $$t10.1 i32)
		(local $$t11 i32)
		(local $$t12 i32)
		(local $$t13 i32)
		(local $$t14 i32)
		block $$BlockFnBody
			loop $$BlockDisp
				block $$Block_6
					block $$Block_5
						block $$Block_4
							block $$Block_3
								block $$Block_2
									block $$Block_1
										block $$Block_0
											block $$BlockSel
												local.get $$block_selector
												br_table  0 1 2 3 4 5 6 0
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
									local.get $$t0
									local.set $$ret_0
									br $$BlockFnBody
								end
								local.get $$current_block
								i32.const 0
								i32.eq
								if(result i32)
									i32.const 0
								else
									local.get $$t1
								end
								local.get $$current_block
								i32.const 0
								i32.eq
								if(result i32)
									i32.const 0
								else
									local.get $$t2
								end
								local.set $$t3
								local.set $$t0
								i32.const 3
								local.set $$current_block
								local.get $$t3
								i32.const 4
								i32.lt_s
								local.set $$t4
								local.get $$t4
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
							local.get $this.0
							call $runtime.Block.Retain
							local.get $this.1
							i32.const 12
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
							i32.const 16
							local.get $$t3
							i32.mul
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
							i32.const 4
							local.get $$t9
							i32.mul
							i32.add
							local.set $$t10.1
							local.get $$t10.0
							call $runtime.Block.Release
							local.set $$t10.0
							local.get $$t10.1
							i32.load
							local.set $$t11
							local.get $$t1
							local.get $$t11
							i32.add
							local.set $$t12
							local.get $$t9
							i32.const 1
							i32.add
							local.set $$t13
							br $$Block_5
						end
						i32.const 5
						local.set $$current_block
						local.get $$t3
						i32.const 1
						i32.add
						local.set $$t2
						i32.const 3
						local.set $$block_selector
						br $$BlockDisp
					end
					local.get $$current_block
					i32.const 1
					i32.eq
					if(result i32)
						local.get $$t0
					else
						local.get $$t12
					end
					local.get $$current_block
					i32.const 1
					i32.eq
					if(result i32)
						i32.const 0
					else
						local.get $$t13
					end
					local.set $$t9
					local.set $$t1
					i32.const 6
					local.set $$current_block
					local.get $$t9
					i32.const 4
					i32.lt_s
					local.set $$t14
					local.get $$t14
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
			end
		end
		local.get $$ret_0
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
	)
	(func $w42048$game.board.down (param $this.0 i32) (param $this.1 i32)
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
					local.get $this.0
					local.get $this.1
					call $w42048$game.board.transpose
					local.get $this.0
					local.get $this.1
					call $w42048$game.board.left
					local.get $this.0
					local.get $this.1
					call $w42048$game.board.transpose
					local.get $this.0
					local.get $this.1
					call $w42048$game.board.transpose
					local.get $this.0
					local.get $this.1
					call $w42048$game.board.transpose
					br $$BlockFnBody
				end
			end
		end
	)
	(func $w42048$game.board.intn (param $this.0 i32) (param $this.1 i32) (param $n i32) (result i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$ret_0 i32)
		(local $$t0 i32)
		(local $$t1.0 i32)
		(local $$t1.1 i32)
		(local $$t2.0 i32)
		(local $$t2.1.0 i32)
		(local $$t2.1.1 i32)
		(local $$t3 i32)
		(local $$t4 i32)
		(local $$t5 i32)
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
						i32.const -1
						local.set $$ret_0
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
					local.get $$t1.1
					i32.load offset=4
					call $runtime.Block.Retain
					local.get $$t1.1
					i32.load offset=8
					local.set $$t2.1.1
					local.get $$t2.1.0
					call $runtime.Block.Release
					local.set $$t2.1.0
					local.set $$t2.0
					local.get $$t2.0
					local.get $$t2.1.1
					global.set $$wa.runtime.closure_data
					call_indirect 0 (type $$$fnSig8)
					local.set $$t3
					local.get $$t3
					local.set $$t4
					local.get $$t4
					local.get $n
					i32.rem_s
					local.set $$t5
					local.get $$t5
					local.set $$ret_0
					br $$BlockFnBody
				end
			end
		end
		local.get $$ret_0
		local.get $$t1.0
		call $runtime.Block.Release
		local.get $$t2.1.0
		call $runtime.Block.Release
	)
	(func $w42048$game.board.left (param $this.0 i32) (param $this.1 i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$t0.0 i32)
		(local $$t0.1 i32)
		(local $$t1.0 i32)
		(local $$t1.1 i32)
		(local $$t1.2 i32)
		(local $$t1.3 i32)
		(local $$t2 i32)
		(local $$t3.0 i32)
		(local $$t3.1 i32)
		(local $$t4.0 i32)
		(local $$t4.1 i32)
		(local $$t4.2 i32)
		(local $$t4.3 i32)
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
		(local $$t8.2 i32)
		(local $$t8.3 i32)
		(local $$t9 i32)
		(local $$t10 i32)
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
								br $$Block_2
							end
							i32.const 1
							local.set $$current_block
							local.get $this.0
							call $runtime.Block.Retain
							local.get $this.1
							i32.const 12
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
							i32.const 16
							local.get $$t2
							i32.mul
							i32.add
							local.set $$t3.1
							local.get $$t3.0
							call $runtime.Block.Release
							local.set $$t3.0
							local.get $$t3.1
							i32.load
							call $runtime.Block.Retain
							local.get $$t3.1
							i32.load offset=4
							local.get $$t3.1
							i32.load offset=8
							local.get $$t3.1
							i32.load offset=12
							local.set $$t4.3
							local.set $$t4.2
							local.set $$t4.1
							local.get $$t4.0
							call $runtime.Block.Release
							local.set $$t4.0
							local.get $this.0
							call $runtime.Block.Retain
							local.get $this.1
							i32.const 12
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
							i32.const 16
							local.get $$t2
							i32.mul
							i32.add
							local.set $$t7.1
							local.get $$t7.0
							call $runtime.Block.Release
							local.set $$t7.0
							local.get $$t4.0
							local.get $$t4.1
							local.get $$t4.2
							local.get $$t4.3
							call $w42048$game.movedRow
							local.set $$t8.3
							local.set $$t8.2
							local.set $$t8.1
							local.get $$t8.0
							call $runtime.Block.Release
							local.set $$t8.0
							local.get $$t7.1
							local.get $$t8.0
							call $runtime.Block.Retain
							local.get $$t7.1
							i32.load align=1
							call $runtime.Block.Release
							i32.store align=1
							local.get $$t7.1
							local.get $$t8.1
							i32.store offset=4
							local.get $$t7.1
							local.get $$t8.2
							i32.store offset=8
							local.get $$t7.1
							local.get $$t8.3
							i32.store offset=12
							local.get $$t2
							i32.const 1
							i32.add
							local.set $$t9
							br $$Block_2
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
						local.get $$t9
					end
					local.set $$t2
					i32.const 3
					local.set $$current_block
					local.get $$t2
					i32.const 4
					i32.lt_s
					local.set $$t10
					local.get $$t10
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
			end
		end
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
	)
	(func $w42048$game.board.reverse (param $this.0 i32) (param $this.1 i32)
		(local $$block_selector i32)
		(local $$current_block i32)
		(local $$t0.0 i32)
		(local $$t0.1 i32)
		(local $$t1.0 i32)
		(local $$t1.1 i32)
		(local $$t1.2 i32)
		(local $$t1.3 i32)
		(local $$t2 i32)
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
		(local $$t7.0 i32)
		(local $$t7.1 i32)
		(local $$t7.2 i32)
		(local $$t7.3 i32)
		(local $$t8.0 i32)
		(local $$t8.1 i32)
		(local $$t8.2 i32)
		(local $$t8.3 i32)
		(local $$t9 i32)
		(local $$t10 i32)
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
								br $$Block_2
							end
							i32.const 1
							local.set $$current_block
							local.get $this.0
							call $runtime.Block.Retain
							local.get $this.1
							i32.const 12
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
							i32.const 16
							local.get $$t2
							i32.mul
							i32.add
							local.set $$t3.1
							local.get $$t3.0
							call $runtime.Block.Release
							local.set $$t3.0
							local.get $this.0
							call $runtime.Block.Retain
							local.get $this.1
							i32.const 12
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
							i32.const 16
							local.get $$t2
							i32.mul
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
							local.get $$t7.1
							local.get $$t7.2
							local.get $$t7.3
							call $w42048$game.reverseRow
							local.set $$t8.3
							local.set $$t8.2
							local.set $$t8.1
							local.get $$t8.0
							call $runtime.Block.Release
							local.set $$t8.0
							local.get $$t3.1
							local.get $$t8.0
							call $runtime.Block.Retain
							local.get $$t3.1
							i32.load align=1
							call $runtime.Block.Release
							i32.store align=1
							local.get $$t3.1
							local.get $$t8.1
							i32.store offset=4
							local.get $$t3.1
							local.get $$t8.2
							i32.store offset=8
							local.get $$t3.1
							local.get $$t8.3
							i32.store offset=12
							local.get $$t2
							i32.const 1
							i32.add
							local.set $$t9
							br $$Block_2
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
						local.get $$t9
					end
					local.set $$t2
					i32.const 3
					local.set $$current_block
					local.get $$t2
					i32.const 4
					i32.lt_s
					local.set $$t10
					local.get $$t10
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
			end
		end
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
	)
	(func $w42048$game.board.reverseRows (param $this.0 i32) (param $this.1 i32)
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
		(local $$t8.0 i32)
		(local $$t8.1 i32)
		(local $$t8.2 i32)
		(local $$t8.3 i32)
		(local $$t9 i32)
		(local $$t10 i32)
		(local $$t11 i32)
		(local $$t12.0 i32)
		(local $$t12.1 i32)
		(local $$t13 i32)
		(local $$t14 i32)
		(local $$t15 i32)
		(local $$t16 i32)
		(local $$t17 i32)
		(local $$t18.0 i32)
		(local $$t18.1 i32)
		(local $$t19.0 i32)
		(local $$t19.1 i32)
		(local $$t19.2 i32)
		(local $$t19.3 i32)
		(local $$t20 i32)
		(local $$t21.0 i32)
		(local $$t21.1 i32)
		(local $$t22.0 i32)
		(local $$t22.1 i32)
		(local $$t23.0 i32)
		(local $$t23.1 i32)
		(local $$t23.2 i32)
		(local $$t23.3 i32)
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
		(local $$t29 i32)
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
														i32.const 16
														call $runtime.HeapAlloc
														i32.const 1
														i32.const 0
														i32.const 0
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
														i32.const 16
														i32.const 0
														i32.mul
														i32.add
														i32.const 0
														i32.const 0
														i32.sub
														i32.const 0
														i32.const 0
														i32.sub
														local.set $$t1.3
														local.set $$t1.2
														local.set $$t1.1
														local.get $$t1.0
														call $runtime.Block.Release
														local.set $$t1.0
														br $$Block_2
													end
													i32.const 1
													local.set $$current_block
													i32.const 32
													call $runtime.HeapAlloc
													i32.const 1
													i32.const 0
													i32.const 16
													call $runtime.Block.Init
													call $runtime.DupI32
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
													i32.const 0
													i32.mul
													i32.add
													i32.const 4
													i32.const 0
													i32.sub
													i32.const 4
													i32.const 0
													i32.sub
													local.set $$t3.3
													local.set $$t3.2
													local.set $$t3.1
													local.get $$t3.0
													call $runtime.Block.Release
													local.set $$t3.0
													i32.const 32
													call $runtime.HeapAlloc
													i32.const 1
													i32.const 40
													i32.const 16
													call $runtime.Block.Init
													call $runtime.DupI32
													i32.const 16
													i32.add
													local.set $$t4.1
													local.get $$t4.0
													call $runtime.Block.Release
													local.set $$t4.0
													local.get $$t4.0
													call $runtime.Block.Retain
													local.get $$t4.1
													i32.const 16
													i32.const 0
													i32.mul
													i32.add
													local.set $$t5.1
													local.get $$t5.0
													call $runtime.Block.Release
													local.set $$t5.0
													local.get $$t5.1
													local.get $$t3.0
													call $runtime.Block.Retain
													local.get $$t5.1
													i32.load align=1
													call $runtime.Block.Release
													i32.store align=1
													local.get $$t5.1
													local.get $$t3.1
													i32.store offset=4
													local.get $$t5.1
													local.get $$t3.2
													i32.store offset=8
													local.get $$t5.1
													local.get $$t3.3
													i32.store offset=12
													local.get $$t4.0
													call $runtime.Block.Retain
													local.get $$t4.1
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
													local.set $$t6.3
													local.set $$t6.2
													local.set $$t6.1
													local.get $$t6.0
													call $runtime.Block.Release
													local.set $$t6.0
													local.get $$t7.0
													local.get $$t7.1
													local.get $$t7.2
													local.get $$t7.3
													local.get $$t6.0
													local.get $$t6.1
													local.get $$t6.2
													local.get $$t6.3
													call $$i32.$slice.$slice.append
													local.set $$t8.3
													local.set $$t8.2
													local.set $$t8.1
													local.get $$t8.0
													call $runtime.Block.Release
													local.set $$t8.0
													local.get $$t9
													i32.const 1
													i32.add
													local.set $$t10
													br $$Block_2
												end
												i32.const 2
												local.set $$current_block
												br $$Block_5
											end
											local.get $$current_block
											i32.const 0
											i32.eq
											if(result i32 i32 i32 i32)
												local.get $$t1.0
												call $runtime.Block.Retain
												local.get $$t1.1
												local.get $$t1.2
												local.get $$t1.3
											else
												local.get $$t8.0
												call $runtime.Block.Retain
												local.get $$t8.1
												local.get $$t8.2
												local.get $$t8.3
											end
											local.get $$current_block
											i32.const 0
											i32.eq
											if(result i32)
												i32.const 0
											else
												local.get $$t10
											end
											local.set $$t9
											local.set $$t7.3
											local.set $$t7.2
											local.set $$t7.1
											local.get $$t7.0
											call $runtime.Block.Release
											local.set $$t7.0
											i32.const 3
											local.set $$current_block
											local.get $$t9
											i32.const 4
											i32.lt_s
											local.set $$t11
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
										end
										i32.const 4
										local.set $$current_block
										br $$Block_8
									end
									i32.const 5
									local.set $$current_block
									local.get $this.0
									call $runtime.Block.Retain
									local.get $this.1
									i32.const 12
									i32.add
									local.set $$t12.1
									local.get $$t12.0
									call $runtime.Block.Release
									local.set $$t12.0
									local.get $$t12.1
									local.get $$t7.0
									call $runtime.Block.Retain
									local.get $$t12.1
									i32.load align=1
									call $runtime.Block.Release
									i32.store align=1
									local.get $$t12.1
									local.get $$t7.1
									i32.store offset=4
									local.get $$t12.1
									local.get $$t7.2
									i32.store offset=8
									local.get $$t12.1
									local.get $$t7.3
									i32.store offset=12
									br $$BlockFnBody
								end
								local.get $$current_block
								i32.const 2
								i32.eq
								if(result i32)
									i32.const 0
								else
									local.get $$t13
								end
								local.set $$t14
								i32.const 6
								local.set $$current_block
								local.get $$t14
								i32.const 4
								i32.lt_s
								local.set $$t15
								local.get $$t15
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
							i32.const 4
							local.get $$t14
							i32.sub
							local.set $$t16
							local.get $$t16
							i32.const 1
							i32.sub
							local.set $$t17
							local.get $$t7.0
							call $runtime.Block.Retain
							local.get $$t7.1
							i32.const 16
							local.get $$t17
							i32.mul
							i32.add
							local.set $$t18.1
							local.get $$t18.0
							call $runtime.Block.Release
							local.set $$t18.0
							local.get $$t18.1
							i32.load
							call $runtime.Block.Retain
							local.get $$t18.1
							i32.load offset=4
							local.get $$t18.1
							i32.load offset=8
							local.get $$t18.1
							i32.load offset=12
							local.set $$t19.3
							local.set $$t19.2
							local.set $$t19.1
							local.get $$t19.0
							call $runtime.Block.Release
							local.set $$t19.0
							local.get $$t19.0
							call $runtime.Block.Retain
							local.get $$t19.1
							i32.const 4
							local.get $$t20
							i32.mul
							i32.add
							local.set $$t21.1
							local.get $$t21.0
							call $runtime.Block.Release
							local.set $$t21.0
							local.get $this.0
							call $runtime.Block.Retain
							local.get $this.1
							i32.const 12
							i32.add
							local.set $$t22.1
							local.get $$t22.0
							call $runtime.Block.Release
							local.set $$t22.0
							local.get $$t22.1
							i32.load
							call $runtime.Block.Retain
							local.get $$t22.1
							i32.load offset=4
							local.get $$t22.1
							i32.load offset=8
							local.get $$t22.1
							i32.load offset=12
							local.set $$t23.3
							local.set $$t23.2
							local.set $$t23.1
							local.get $$t23.0
							call $runtime.Block.Release
							local.set $$t23.0
							local.get $$t23.0
							call $runtime.Block.Retain
							local.get $$t23.1
							i32.const 16
							local.get $$t14
							i32.mul
							i32.add
							local.set $$t24.1
							local.get $$t24.0
							call $runtime.Block.Release
							local.set $$t24.0
							local.get $$t24.1
							i32.load
							call $runtime.Block.Retain
							local.get $$t24.1
							i32.load offset=4
							local.get $$t24.1
							i32.load offset=8
							local.get $$t24.1
							i32.load offset=12
							local.set $$t25.3
							local.set $$t25.2
							local.set $$t25.1
							local.get $$t25.0
							call $runtime.Block.Release
							local.set $$t25.0
							local.get $$t25.0
							call $runtime.Block.Retain
							local.get $$t25.1
							i32.const 4
							local.get $$t20
							i32.mul
							i32.add
							local.set $$t26.1
							local.get $$t26.0
							call $runtime.Block.Release
							local.set $$t26.0
							local.get $$t26.1
							i32.load
							local.set $$t27
							local.get $$t21.1
							local.get $$t27
							i32.store
							local.get $$t20
							i32.const 1
							i32.add
							local.set $$t28
							br $$Block_8
						end
						i32.const 8
						local.set $$current_block
						local.get $$t14
						i32.const 1
						i32.add
						local.set $$t13
						i32.const 6
						local.set $$block_selector
						br $$BlockDisp
					end
					local.get $$current_block
					i32.const 4
					i32.eq
					if(result i32)
						i32.const 0
					else
						local.get $$t28
					end
					local.set $$t20
					i32.const 9
					local.set $$current_block
					local.get $$t20
					i32.const 4
					i32.lt_s
					local.set $$t29
					local.get $$t29
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
		local.get $$t12.0
		call $runtime.Block.Release
		local.get $$t18.0
		call $runtime.Block.Release
		local.get $$t19.0
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
	)
	(func $w42048$game.board.right (param $this.0 i32) (param $this.1 i32)
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
					local.get $this.0
					local.get $this.1
					call $w42048$game.board.reverse
					local.get $this.0
					local.get $this.1
					call $w42048$game.board.left
					local.get $this.0
					local.get $this.1
					call $w42048$game.board.reverse
					br $$BlockFnBody
				end
			end
		end
	)
	(func $w42048$game.board.transpose (param $this.0 i32) (param $this.1 i32)
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
		(local $$t8.0 i32)
		(local $$t8.1 i32)
		(local $$t8.2 i32)
		(local $$t8.3 i32)
		(local $$t9 i32)
		(local $$t10 i32)
		(local $$t11 i32)
		(local $$t12.0 i32)
		(local $$t12.1 i32)
		(local $$t13 i32)
		(local $$t14 i32)
		(local $$t15 i32)
		(local $$t16.0 i32)
		(local $$t16.1 i32)
		(local $$t17.0 i32)
		(local $$t17.1 i32)
		(local $$t17.2 i32)
		(local $$t17.3 i32)
		(local $$t18 i32)
		(local $$t19.0 i32)
		(local $$t19.1 i32)
		(local $$t20.0 i32)
		(local $$t20.1 i32)
		(local $$t21.0 i32)
		(local $$t21.1 i32)
		(local $$t21.2 i32)
		(local $$t21.3 i32)
		(local $$t22 i32)
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
		(local $$t29 i32)
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
														i32.const 16
														call $runtime.HeapAlloc
														i32.const 1
														i32.const 0
														i32.const 0
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
														i32.const 16
														i32.const 0
														i32.mul
														i32.add
														i32.const 0
														i32.const 0
														i32.sub
														i32.const 0
														i32.const 0
														i32.sub
														local.set $$t1.3
														local.set $$t1.2
														local.set $$t1.1
														local.get $$t1.0
														call $runtime.Block.Release
														local.set $$t1.0
														br $$Block_2
													end
													i32.const 1
													local.set $$current_block
													i32.const 32
													call $runtime.HeapAlloc
													i32.const 1
													i32.const 0
													i32.const 16
													call $runtime.Block.Init
													call $runtime.DupI32
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
													i32.const 0
													i32.mul
													i32.add
													i32.const 4
													i32.const 0
													i32.sub
													i32.const 4
													i32.const 0
													i32.sub
													local.set $$t3.3
													local.set $$t3.2
													local.set $$t3.1
													local.get $$t3.0
													call $runtime.Block.Release
													local.set $$t3.0
													i32.const 32
													call $runtime.HeapAlloc
													i32.const 1
													i32.const 40
													i32.const 16
													call $runtime.Block.Init
													call $runtime.DupI32
													i32.const 16
													i32.add
													local.set $$t4.1
													local.get $$t4.0
													call $runtime.Block.Release
													local.set $$t4.0
													local.get $$t4.0
													call $runtime.Block.Retain
													local.get $$t4.1
													i32.const 16
													i32.const 0
													i32.mul
													i32.add
													local.set $$t5.1
													local.get $$t5.0
													call $runtime.Block.Release
													local.set $$t5.0
													local.get $$t5.1
													local.get $$t3.0
													call $runtime.Block.Retain
													local.get $$t5.1
													i32.load align=1
													call $runtime.Block.Release
													i32.store align=1
													local.get $$t5.1
													local.get $$t3.1
													i32.store offset=4
													local.get $$t5.1
													local.get $$t3.2
													i32.store offset=8
													local.get $$t5.1
													local.get $$t3.3
													i32.store offset=12
													local.get $$t4.0
													call $runtime.Block.Retain
													local.get $$t4.1
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
													local.set $$t6.3
													local.set $$t6.2
													local.set $$t6.1
													local.get $$t6.0
													call $runtime.Block.Release
													local.set $$t6.0
													local.get $$t7.0
													local.get $$t7.1
													local.get $$t7.2
													local.get $$t7.3
													local.get $$t6.0
													local.get $$t6.1
													local.get $$t6.2
													local.get $$t6.3
													call $$i32.$slice.$slice.append
													local.set $$t8.3
													local.set $$t8.2
													local.set $$t8.1
													local.get $$t8.0
													call $runtime.Block.Release
													local.set $$t8.0
													local.get $$t9
													i32.const 1
													i32.add
													local.set $$t10
													br $$Block_2
												end
												i32.const 2
												local.set $$current_block
												br $$Block_5
											end
											local.get $$current_block
											i32.const 0
											i32.eq
											if(result i32 i32 i32 i32)
												local.get $$t1.0
												call $runtime.Block.Retain
												local.get $$t1.1
												local.get $$t1.2
												local.get $$t1.3
											else
												local.get $$t8.0
												call $runtime.Block.Retain
												local.get $$t8.1
												local.get $$t8.2
												local.get $$t8.3
											end
											local.get $$current_block
											i32.const 0
											i32.eq
											if(result i32)
												i32.const 0
											else
												local.get $$t10
											end
											local.set $$t9
											local.set $$t7.3
											local.set $$t7.2
											local.set $$t7.1
											local.get $$t7.0
											call $runtime.Block.Release
											local.set $$t7.0
											i32.const 3
											local.set $$current_block
											local.get $$t9
											i32.const 4
											i32.lt_s
											local.set $$t11
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
										end
										i32.const 4
										local.set $$current_block
										br $$Block_8
									end
									i32.const 5
									local.set $$current_block
									local.get $this.0
									call $runtime.Block.Retain
									local.get $this.1
									i32.const 12
									i32.add
									local.set $$t12.1
									local.get $$t12.0
									call $runtime.Block.Release
									local.set $$t12.0
									local.get $$t12.1
									local.get $$t7.0
									call $runtime.Block.Retain
									local.get $$t12.1
									i32.load align=1
									call $runtime.Block.Release
									i32.store align=1
									local.get $$t12.1
									local.get $$t7.1
									i32.store offset=4
									local.get $$t12.1
									local.get $$t7.2
									i32.store offset=8
									local.get $$t12.1
									local.get $$t7.3
									i32.store offset=12
									br $$BlockFnBody
								end
								local.get $$current_block
								i32.const 2
								i32.eq
								if(result i32)
									i32.const 0
								else
									local.get $$t13
								end
								local.set $$t14
								i32.const 6
								local.set $$current_block
								local.get $$t14
								i32.const 4
								i32.lt_s
								local.set $$t15
								local.get $$t15
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
							local.get $$t7.0
							call $runtime.Block.Retain
							local.get $$t7.1
							i32.const 16
							local.get $$t14
							i32.mul
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
							i32.const 4
							local.get $$t18
							i32.mul
							i32.add
							local.set $$t19.1
							local.get $$t19.0
							call $runtime.Block.Release
							local.set $$t19.0
							local.get $this.0
							call $runtime.Block.Retain
							local.get $this.1
							i32.const 12
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
							i32.const 4
							local.get $$t18
							i32.sub
							local.set $$t22
							local.get $$t22
							i32.const 1
							i32.sub
							local.set $$t23
							local.get $$t21.0
							call $runtime.Block.Retain
							local.get $$t21.1
							i32.const 16
							local.get $$t23
							i32.mul
							i32.add
							local.set $$t24.1
							local.get $$t24.0
							call $runtime.Block.Release
							local.set $$t24.0
							local.get $$t24.1
							i32.load
							call $runtime.Block.Retain
							local.get $$t24.1
							i32.load offset=4
							local.get $$t24.1
							i32.load offset=8
							local.get $$t24.1
							i32.load offset=12
							local.set $$t25.3
							local.set $$t25.2
							local.set $$t25.1
							local.get $$t25.0
							call $runtime.Block.Release
							local.set $$t25.0
							local.get $$t25.0
							call $runtime.Block.Retain
							local.get $$t25.1
							i32.const 4
							local.get $$t14
							i32.mul
							i32.add
							local.set $$t26.1
							local.get $$t26.0
							call $runtime.Block.Release
							local.set $$t26.0
							local.get $$t26.1
							i32.load
							local.set $$t27
							local.get $$t19.1
							local.get $$t27
							i32.store
							local.get $$t18
							i32.const 1
							i32.add
							local.set $$t28
							br $$Block_8
						end
						i32.const 8
						local.set $$current_block
						local.get $$t14
						i32.const 1
						i32.add
						local.set $$t13
						i32.const 6
						local.set $$block_selector
						br $$BlockDisp
					end
					local.get $$current_block
					i32.const 4
					i32.eq
					if(result i32)
						i32.const 0
					else
						local.get $$t28
					end
					local.set $$t18
					i32.const 9
					local.set $$current_block
					local.get $$t18
					i32.const 4
					i32.lt_s
					local.set $$t29
					local.get $$t29
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
		local.get $$t12.0
		call $runtime.Block.Release
		local.get $$t16.0
		call $runtime.Block.Release
		local.get $$t17.0
		call $runtime.Block.Release
		local.get $$t19.0
		call $runtime.Block.Release
		local.get $$t20.0
		call $runtime.Block.Release
		local.get $$t21.0
		call $runtime.Block.Release
		local.get $$t24.0
		call $runtime.Block.Release
		local.get $$t25.0
		call $runtime.Block.Release
		local.get $$t26.0
		call $runtime.Block.Release
	)
	(func $w42048$game.board.up (param $this.0 i32) (param $this.1 i32)
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
					local.get $this.0
					local.get $this.1
					call $w42048$game.board.reverseRows
					local.get $this.0
					local.get $this.1
					call $w42048$game.board.down
					local.get $this.0
					local.get $this.1
					call $w42048$game.board.reverseRows
					br $$BlockFnBody
				end
			end
		end
	)
	(func $_start (export "_start")
		call $w42048.init
	)
	(func $_main (export "_main"))
	(elem (i32.const 1) $$string.$$compAddr)
	(elem (i32.const 2) $$u8.$$block.$$onFree)
	(elem (i32.const 3) $$string.underlying.$$onFree)
	(elem (i32.const 4) $$runtime.mapImp.$$block.$$onFree)
	(elem (i32.const 5) $$runtime.mapImp.$ref.underlying.$$onFree)
	(elem (i32.const 6) $$runtime.mapIter.$$onFree)
	(elem (i32.const 7) $$runtime.mapNode.$$block.$$onFree)
	(elem (i32.const 8) $$runtime.mapNode.$ref.underlying.$$onFree)
	(elem (i32.const 9) $$void.$$block.$$onFree)
	(elem (i32.const 10) $$void.$ref.underlying.$$onFree)
	(elem (i32.const 11) $$i`0`.underlying.$$onFree)
	(elem (i32.const 12) $$runtime.mapNode.$$onFree)
	(elem (i32.const 13) $$runtime.mapNode.$ref.$$block.$$onFree)
	(elem (i32.const 14) $$runtime.mapNode.$ref.$slice.underlying.$$onFree)
	(elem (i32.const 15) $$runtime.mapImp.$$onFree)
	(elem (i32.const 16) $$runtime.mapNode.$ref.$array1.underlying.$$onFree)
	(elem (i32.const 17) $$$$$$.underlying.$$onFree)
	(elem (i32.const 18) $$$$$$.$array1.underlying.$$onFree)
	(elem (i32.const 19) $$$$$$.$$block.$$onFree)
	(elem (i32.const 20) $$$$$$.$slice.underlying.$$onFree)
	(elem (i32.const 21) $$runtime.defers.$$onFree)
	(elem (i32.const 22) $$runtime.defers.$array1.underlying.$$onFree)
	(elem (i32.const 23) $$errors.errorString.$$onFree)
	(elem (i32.const 24) $$.error.underlying.$$onFree)
	(elem (i32.const 25) $$strconv.NumError.$$onFree)
	(elem (i32.const 26) $$u8.$slice.underlying.$$onFree)
	(elem (i32.const 27) $$strconv.decimalSlice.$$onFree)
	(elem (i32.const 28) $$w42048$game.Board.underlying.$$onFree)
	(elem (i32.const 29) $$w42048.UI.$$onFree)
	(elem (i32.const 30) $$$$$u32$$.underlying.$$onFree)
	(elem (i32.const 31) $$i32.$slice.$$block.$$onFree)
	(elem (i32.const 32) $$i32.$slice.$slice.underlying.$$onFree)
	(elem (i32.const 33) $$w42048$game.board.$$onFree)
	(elem (i32.const 34) $w42048$game.lcg$1.$warpfn)
	(elem (i32.const 35) $$u32.$$block.$$onFree)
	(elem (i32.const 36) $$u32.$ref.underlying.$$onFree)
	(elem (i32.const 37) $$w42048$game.lcg$1.$warpdata.$$onFree)
	(elem (i32.const 38) $$i32.$$block.$$onFree)
	(elem (i32.const 39) $$i32.$slice.underlying.$$onFree)
	(elem (i32.const 40) $$i32.$slice.$array1.underlying.$$onFree)
	(elem (i32.const 41) $errors.errorString.Error)
	(elem (i32.const 42) $strconv.NumError.Error)
	(elem (i32.const 43) $w42048$game.board.Add)
	(elem (i32.const 44) $w42048$game.board.Get)
	(elem (i32.const 45) $w42048$game.board.Input)
	(elem (i32.const 46) $w42048$game.board.IsOver)
	(elem (i32.const 47) $w42048$game.board.Restart)
	(elem (i32.const 48) $w42048$game.board.Total)
)