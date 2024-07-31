(module $__walang__

(func $runtime.printI64 (param $v i64))

(func $$runtime.waPrintI32 (param $i i32)
  (local $$block_selector i32)
  (local $$current_block i32)
  (local $$t0 i64)
  block $$BlockFnBody
    loop $$BlockDisp
      block $$Block_0
        block $$BlockSel
          local.get $$block_selector
          br_table 0 0
        end ;;$BlockSel
        i32.const 0
        local.set $$current_block

        ;;convert i64 <- i32 (i)
        local.get $i
        i64.extend_i32_s
        local.set $$t0

        ;;printI64(t0)
        local.get $$t0
        call $runtime.printI64

        ;;return
        br $$BlockFnBody

      end ;;$Block_0
    end ;;$BlockDisp
  end ;;$BlockFnBody
) ;;$runtime.waPrintI32

)
