(module $__walang__

;; Copyright 2023 The Wa Authors. All rights reserved.

(memory $memory 1024)

(export "memory" (memory $memory))

;; +-----------------+---------------------+--------------+
;; | 0 <-- stack --> | <-- static-data --> | <-- heap --> |
;; +-----------------+---------------------+--------------+

(global $__stack_ptr (mut i32) (i32.const 1024))     ;; index=0
(global $__heap_base i32 (i32.const 1048576))     ;; index=1
(global $__heap_max  i32       (i32.const 67108864)) ;; 64MB, 1024 page

;; Copyright 2023 The Wa Authors. All rights reserved.

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

(func $runtime.heapMax(result i32)
	global.get $__heap_max
)

(global $$knr_basep (mut i32) (i32.const 0))
(global $$knr_freep (mut i32) (i32.const 0))

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
                                    br_table 0 1 2 3 4 5 6 7 8 9 10 11 12 13 0
                                  end ;;$BlockSel
                                  i32.const 0
                                  local.set $$current_block

                                  ;;nbytes == 0:u32
                                  local.get $nbytes
                                  i32.const 0
                                  i32.eq
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

                                ;;return 0:u32
                                i32.const 0
                                local.set $$ret_0
                                br $$BlockFnBody

                              end ;;$Block_1
                              i32.const 2
                              local.set $$current_block

                              ;;*knr_basep
                              global.get $$knr_basep
                              local.set $$t1

                              ;;t1 == 0:u32
                              local.get $$t1
                              i32.const 0
                              i32.eq
                              local.set $$t2

                              ;;if t2 goto 4 else 5
                              local.get $$t2
                              if
                                br $$Block_3
                              else
                                br $$Block_4
                              end

                            end ;;$Block_2
                            i32.const 3
                            local.set $$current_block

                            ;;heapMax()
                            call $runtime.heapMax
                            local.set $$t3

                            ;;*knr_basep
                            global.get $$knr_basep
                            local.set $$t4

                            ;;t3 - t4
                            local.get $$t3
                            local.get $$t4
                            i32.sub
                            local.set $$t5

                            ;;nbytes >= t5
                            local.get $nbytes
                            local.get $$t5
                            i32.ge_u
                            local.set $$t6

                            ;;if t6 goto 1 else 2
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

                          end ;;$Block_3
                          i32.const 4
                          local.set $$current_block

                          ;;heapBase()
                          call $runtime.heapBase
                          local.set $$t7

                          ;;*knr_basep = t7
                          local.get $$t7
                          global.set $$knr_basep

                          ;;*knr_basep
                          global.get $$knr_basep
                          local.set $$t8

                          ;;*knr_freep = t8
                          local.get $$t8
                          global.set $$knr_freep

                          ;;local knr_Header (base)
                          i32.const 0
                          local.set $$t9.0
                          i32.const 0
                          local.set $$t9.1

                          ;;&t9.ptr [#0]

                          ;;*knr_basep
                          global.get $$knr_basep
                          local.set $$t10

                          ;;&t9.size [#1]

                          ;;heapMax()
                          call $runtime.heapMax
                          local.set $$t11

                          ;;*knr_basep
                          global.get $$knr_basep
                          local.set $$t12

                          ;;t13 - t14
                          local.get $$t11
                          local.get $$t12
                          i32.sub
                          local.set $$t13

                          ;;t15 / 8:u32
                          local.get $$t13
                          i32.const 8
                          i32.div_u
                          local.set $$t14

                          ;;t16 - 1:u32
                          local.get $$t14
                          i32.const 1
                          i32.sub
                          local.set $$t15

                          ;;*t10 = t11
                          local.get $$t10
                          local.set $$t9.0

                          ;;*t12 = t17
                          local.get $$t15
                          local.set $$t9.1

                          ;;*knr_basep
                          global.get $$knr_basep
                          local.set $$t16

                          ;;*t9
                          local.get $$t9.0
                          local.get $$t9.1
                          local.set $$t17.1
                          local.set $$t17.0

                          ;;knr_setBlockHeader(t18, t19)
                          local.get $$t16
                          local.get $$t17.0
                          local.get $$t17.1
                          call $runtime.knr_setBlockHeader

                          ;;jump 5
                          br $$Block_4

                        end ;;$Block_4
                        i32.const 5
                        local.set $$current_block

                        ;;nbytes + 8:u32
                        local.get $nbytes
                        i32.const 8
                        i32.add
                        local.set $$t18

                        ;;t21 - 1:u32
                        local.get $$t18
                        i32.const 1
                        i32.sub
                        local.set $$t19

                        ;;t22 / 8:u32
                        local.get $$t19
                        i32.const 8
                        i32.div_u
                        local.set $$t20

                        ;;t23 + 1:u32
                        local.get $$t20
                        i32.const 1
                        i32.add
                        local.set $$t21

                        ;;*knr_freep
                        global.get $$knr_freep
                        local.set $$t22

                        ;;local knr_Header (prevp)
                        i32.const 0
                        local.set $$t23.0
                        i32.const 0
                        local.set $$t23.1

                        ;;knr_getBlockHeader(t25)
                        local.get $$t22
                        call $runtime.knr_getBlockHeader
                        local.set $$t24.1
                        local.set $$t24.0

                        ;;*t26 = t27
                        local.get $$t24.0
                        local.get $$t24.1
                        local.set $$t23.1
                        local.set $$t23.0

                        ;;&t26.ptr [#0]

                        ;;*t28
                        local.get $$t23.0
                        local.set $$t25

                        ;;local knr_Header (p)
                        i32.const 0
                        local.set $$t26.0
                        i32.const 0
                        local.set $$t26.1

                        ;;knr_getBlockHeader(t29)
                        local.get $$t25
                        call $runtime.knr_getBlockHeader
                        local.set $$t27.1
                        local.set $$t27.0

                        ;;*t30 = t31
                        local.get $$t27.0
                        local.get $$t27.1
                        local.set $$t26.1
                        local.set $$t26.0

                        ;;jump 6
                        br $$Block_5

                      end ;;$Block_5
                      ;;phi [5: t25, 13: t33] #prevp_addr
                      local.get $$current_block
                      i32.const 5
                      i32.eq
                      if (result i32)
                        local.get $$t22
                      else
                        local.get $$t28
                      end
                      local.set $$t29

                      ;;phi [5: t29, 13: t64] #p_addr
                      local.get $$current_block
                      i32.const 5
                      i32.eq
                      if (result i32)
                        local.get $$t25
                      else
                        local.get $$t30
                      end
                      local.set $$t28

                      i32.const 6
                      local.set $$current_block

                      ;;&t30.size [#1]

                      ;;*t34
                      local.get $$t26.1
                      local.set $$t31

                      ;;t35 >= t24
                      local.get $$t31
                      local.get $$t21
                      i32.ge_u
                      local.set $$t32

                      ;;if t36 goto 7 else 8
                      local.get $$t32
                      if
                        br $$Block_6
                      else
                        br $$Block_7
                      end

                    end ;;$Block_6
                    i32.const 7
                    local.set $$current_block

                    ;;&t30.size [#1]

                    ;;*t37
                    local.get $$t26.1
                    local.set $$t33

                    ;;t38 == t24
                    local.get $$t33
                    local.get $$t21
                    i32.eq
                    local.set $$t34

                    ;;if t39 goto 9 else 11
                    local.get $$t34
                    if
                      br $$Block_8
                    else
                      br $$Block_10
                    end

                  end ;;$Block_7
                  i32.const 8
                  local.set $$current_block

                  ;;*knr_freep
                  global.get $$knr_freep
                  local.set $$t35

                  ;;t33 == t40
                  local.get $$t28
                  local.get $$t35
                  i32.eq
                  local.set $$t36

                  ;;if t41 goto 12 else 13
                  local.get $$t36
                  if
                    br $$Block_11
                  else
                    br $$Block_12
                  end

                end ;;$Block_8
                i32.const 9
                local.set $$current_block

                ;;&t26.ptr [#0]

                ;;&t30.ptr [#0]

                ;;*t43
                local.get $$t26.0
                local.set $$t37

                ;;*t42 = t44
                local.get $$t37
                local.set $$t23.0

                ;;*t26
                local.get $$t23.0
                local.get $$t23.1
                local.set $$t38.1
                local.set $$t38.0

                ;;knr_setBlockHeader(t32, t45)
                local.get $$t29
                local.get $$t38.0
                local.get $$t38.1
                call $runtime.knr_setBlockHeader

                ;;jump 10
                br $$Block_9

              end ;;$Block_9
              ;;phi [9: t33, 11: t57] #p_addr
              local.get $$current_block
              i32.const 9
              i32.eq
              if (result i32)
                local.get $$t28
              else
                local.get $$t39
              end
              local.set $$t40

              i32.const 10
              local.set $$current_block

              ;;*knr_freep = t32
              local.get $$t29
              global.set $$knr_freep

              ;;t47 + 8:u32
              local.get $$t40
              i32.const 8
              i32.add
              local.set $$t41

              ;;return t48
              local.get $$t41
              local.set $$ret_0
              br $$BlockFnBody

            end ;;$Block_10
            i32.const 11
            local.set $$current_block

            ;;&t30.size [#1]

            ;;*t49
            local.get $$t26.1
            local.set $$t42

            ;;t50 - t24
            local.get $$t42
            local.get $$t21
            i32.sub
            local.set $$t43

            ;;*t49 = t51
            local.get $$t43
            local.set $$t26.1

            ;;*t30
            local.get $$t26.0
            local.get $$t26.1
            local.set $$t44.1
            local.set $$t44.0

            ;;knr_setBlockHeader(t33, t52)
            local.get $$t28
            local.get $$t44.0
            local.get $$t44.1
            call $runtime.knr_setBlockHeader

            ;;&t30.size [#1]

            ;;*t54
            local.get $$t26.1
            local.set $$t45

            ;;t55 * 8:u32
            local.get $$t45
            i32.const 8
            i32.mul
            local.set $$t46

            ;;t33 + t56
            local.get $$t28
            local.get $$t46
            i32.add
            local.set $$t39

            ;;knr_getBlockHeader(t57)
            local.get $$t39
            call $runtime.knr_getBlockHeader
            local.set $$t47.1
            local.set $$t47.0

            ;;*t30 = t58
            local.get $$t47.0
            local.get $$t47.1
            local.set $$t26.1
            local.set $$t26.0

            ;;&t30.size [#1]

            ;;*t59 = t24
            local.get $$t21
            local.set $$t26.1

            ;;*t30
            local.get $$t26.0
            local.get $$t26.1
            local.set $$t48.1
            local.set $$t48.0

            ;;knr_setBlockHeader(t57, t60)
            local.get $$t39
            local.get $$t48.0
            local.get $$t48.1
            call $runtime.knr_setBlockHeader

            ;;jump 10
            i32.const 10
            local.set $$block_selector
            br $$BlockDisp

          end ;;$Block_11
          i32.const 12
          local.set $$current_block

          ;;return 0:u32
          i32.const 0
          local.set $$ret_0
          br $$BlockFnBody

        end ;;$Block_12
        i32.const 13
        local.set $$current_block

        ;;knr_getBlockHeader(t33)
        local.get $$t28
        call $runtime.knr_getBlockHeader
        local.set $$t49.1
        local.set $$t49.0

        ;;*t26 = t62
        local.get $$t49.0
        local.get $$t49.1
        local.set $$t23.1
        local.set $$t23.0

        ;;&t30.ptr [#0]

        ;;*t63
        local.get $$t26.0
        local.set $$t30

        ;;knr_getBlockHeader(t64)
        local.get $$t30
        call $runtime.knr_getBlockHeader
        local.set $$t50.1
        local.set $$t50.0

        ;;*t30 = t65
        local.get $$t50.0
        local.get $$t50.1
        local.set $$t26.1
        local.set $$t26.0

        ;;jump 6
        i32.const 6
        local.set $$block_selector
        br $$BlockDisp

      end ;;$Block_13
    end ;;$BlockDisp
  end ;;$BlockFnBody
  local.get $$ret_0
) ;;runtime.malloc

) ;; module