#include "wa-app.h"

#include <Arduino.h>

// 内存
static uint8_t* app_host_memory = NULL;
static int32_t app_host_memory_page_size = 0;

// 初始化内存
extern "C" void app_memory_init(uint8_t** pp_memory, int32_t* page_size) {
    *pp_memory = app_host_memory;
    *page_size = app_host_memory_page_size;
}

// 内存增长
extern "C" int32_t app_memory_grow(uint8_t** pp_memory, int32_t* page_size, int32_t new_size) {
    return -1; // 不支持扩容
}

extern "C" int32_t host_fn_arduino_millis() {
    return millis();
}

extern "C" void host_fn_arduino_delay(int32_t ms) {
    delay(ms);
}

extern "C" void host_fn_arduino_delayMicroseconds(int32_t us) {
    delayMicroseconds(us);
}

extern "C" void host_fn_arduino_pinMode(int32_t pin, int32_t mode) {
    pinMode(pin, mode);
}

extern "C" int32_t host_fn_arduino_digitalRead(int32_t pin) {
    return digitalRead(pin);
}

extern "C" void host_fn_arduino_digitalWrite(int32_t pin, int32_t value) {
    digitalWrite(pin, value);
}

extern "C" int32_t host_fn_arduino_analogRead(int32_t pin) {
    return analogRead(pin);
}

extern "C" void host_fn_arduino_analogWrite(int32_t pin, int32_t value) {
    analogWrite(pin, value);
}

extern "C" int32_t host_fn_arduino_getPinLED() {
    return LED_BUILTIN;
}

extern "C" void host_fn_arduino_print(int32_t ptr, int32_t len) {
    if(app_host_memory != NULL) {
        Serial.write(&app_host_memory[ptr], len);
    }
}
