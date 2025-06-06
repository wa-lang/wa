#include "wa-app.h"

// 配置初始内存
static uint8_t app_host_memory[{{.MemoryBytes}}];
static int32_t app_host_memory_page_size = ({{.MemoryBytes}})/(1<<16);

// 初始化内存
extern "C" void app_memory_init(uint8_t** pp_memory, int32_t* page_size) {
    *pp_memory = app_host_memory;
    *page_size = app_host_memory_page_size;
}

// 内存增长
extern "C" int32_t app_memory_grow(uint8_t** pp_memory, int32_t* page_size, int32_t new_size) {
    return -1;
}

extern "C" void app_env_blit(
    int32_t sprite,
    int32_t x, int32_t y, int32_t width, int32_t height,
    int32_t flags
) {
    // TODO
}

extern "C" void app_env_blitSub(
    int32_t sprite,
    int32_t x, int32_t y, int32_t width, int32_t height,
    int32_t srcX, int32_t srcY, int32_t stride,
    int32_t flags
) {
    // TODO
}

extern "C" void app_env_line(int32_t x1, int32_t y1, int32_t x2, int32_t y2) {
    // TODO
}

extern "C" void app_env_hline(int32_t x1, int32_t y1, int32_t x2, int32_t y2) {
    // TODO
}

extern "C" void app_env_vline(int32_t x1, int32_t y1, int32_t x2, int32_t y2) {
    // TODO
}

extern "C" void app_env_oval(int32_t x, int32_t y, int32_t width, int32_t height) {
    // TODO
}

extern "C" void app_env_rect(int32_t x, int32_t y, int32_t width, int32_t height) {
    // TODO
}

extern "C" void app_env_textUtf8(int32_t textPtr, int32_t textLen, int32_t x, int32_t y) {
    // TODO
}

extern "C" void app_env_tone(int32_t frequency, int32_t duration, int32_t volume, int32_t flags) {
    // TODO
}

extern "C" void app_env_diskr(int32_t ptr, int32_t count) {
    // TODO
}

extern "C" void app_env_diskw(int32_t ptr, int32_t count) {
    // TODO
}

extern "C" void app_env_traceUtf8(int32_t msgPtr, int32_t msgLen) {
    // TODO
}

