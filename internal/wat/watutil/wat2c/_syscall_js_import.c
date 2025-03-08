#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>

extern uint8_t wasm_memory[];

void host_fn_syscall_js_print_bool(int32_t v) {
  printf(v? "true": "false");
}
void host_fn_syscall_js_print_f32(float v) {
  printf("%f", v);
}

void host_fn_syscall_js_print_f64(double v) {
  printf("%lf", v);
}
void host_fn_syscall_js_print_i32(int32_t v) {
  printf("%d", v);
}
void host_fn_syscall_js_print_i64(int64_t v) {
  printf("%lld", v);
}
void host_fn_syscall_js_print_ptr(int32_t v) {
	printf("%p", &wasm_memory[v]);
}
void host_fn_syscall_js_print_rune(int32_t n) {
  printf("%c", n);
}
void host_fn_syscall_js_print_str(int32_t ptr, int32_t len) {
	int i;
	for (i = 0; i < len; i++) {
		printf("%c", (int)(*(char*)&wasm_memory[ptr+i]));
	}
}

void host_fn_syscall_js_print_u32(int32_t v) {
  printf("%u", v);
}
void host_fn_syscall_js_print_u64(int64_t v) {
  printf("%llu", v);
}
void host_fn_syscall_js_proc_exit(int32_t v) {
  exit(v);
}
