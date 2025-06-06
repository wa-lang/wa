#include "wa-app.h"

// 配置初始内存
static uint8_t app_host_memory[{{.MemoryBytes}}];
static int32_t app_host_memory_page_size = 0;

// 初始化内存
extern "C" void app_memory_init(uint8_t** pp_memory, int32_t* page_size) {
    *pp_memory = app_host_memory;
    *page_size = app_host_memory_page_size;
}

// 内存增长
extern "C" int32_t app_memory_grow(uint8_t** pp_memory, int32_t* page_size, int32_t new_size) {
    return -1;
}

extern "C" void app_syscall_js_print_bool(bool i) {
    printf(i? "true": "false");
}

extern "C" void app_syscall_js_print_i32(int32_t i) {
    printf("%d", i);
}

extern "C" void app_syscall_js_print_i32(uint32_t i) {
    printf("%d", i);
}

extern "C" void app_syscall_js_print_ptr(in32_t i) {
    printf("%x", i);
}

extern "C" void app_syscall_js_print_i64(int64_t i) {
    printf("%lld", i);
}

extern "C" void app_syscall_js_print_u64(uint64_t i) {
    printf("%lld", i);
}

extern "C" void app_syscall_js_print_f32(float i) {
    printf("%f", i);
}

extern "C" void app_syscall_js_print_f64(double i) {
    printf("%llf", i);
}

extern "C" void app_syscall_js_print_rune(in32_t i) {
    printf("%c", i);
}

extern "C" void app_syscall_js_print_str(in32_t ptr, int32_t len) {
    fwrite(stdout, &app_host_memory[ptr], len);
}
