#! /bin/bash

wa native --target=avr arduino_blink.wa
wa native array_0.wa
wa native array_1.wa
wa native array_2.wa
wa native bitwise_logic.wa
wa native compare.wa
wa native convert_0.wa
wa native convert_1.wa
wa native convert_2.wa
wa native convert_3.wa
wa native float32.wa
wa native global_constant.wa
wa native global_variable_0.wa
wa native global_variable_1.wa
wa native heart.wa
wa native loop_0.wa
wa native loop_1.wa
wa native multi_ret.wa
wa native pointer.wa
wa native prime.wa
wa native print_0.wa
wa native print_1.wa
wa native print_2.wa
wa native print_3.wa
wa native shift.wa
wa native struct_0.wa
wa native struct_1.wa
wa native struct_2.wa
wa native struct_3.wa
wa native unconditional_jump.wa

rm -f *.ll *.s
