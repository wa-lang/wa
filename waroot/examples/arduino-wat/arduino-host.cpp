// 版权 @2025 arduino-wat 作者。保留所有权利。

#include "wa-app.h"

#include <Arduino.h>

// 常量定义
extern "C" int32_t app_arduino_HIGH = HIGH;
extern "C" int32_t app_arduino_LOW = LOW;

extern "C" int32_t app_arduino_INPUT = INPUT;
extern "C" int32_t app_arduino_OUTPUT = OUTPUT;

extern "C" int32_t app_arduino_LED_BUILTIN = LED_BUILTIN;

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

extern "C" int32_t app_arduino_millis() {
    return millis();
}

extern "C" void app_arduino_delay(int32_t ms) {
    delay(ms);
}

extern "C" void app_arduino_delayMicroseconds(int32_t us) {
    delayMicroseconds(us);
}

extern "C" void app_arduino_pinMode(int32_t pin, int32_t mode) {
    pinMode(pin, mode);
}

extern "C" int32_t app_arduino_digitalRead(int32_t pin) {
    return digitalRead(pin);
}

extern "C" void app_arduino_digitalWrite(int32_t pin, int32_t value) {
    digitalWrite(pin, value);
}

extern "C" int32_t app_arduino_analogRead(int32_t pin) {
    return analogRead(pin);
}

extern "C" void app_arduino_analogWrite(int32_t pin, int32_t value) {
    analogWrite(pin, value);
}

extern "C" void app_arduino_print(int32_t ptr, int32_t len) {
    if(app_host_memory != NULL) {
        Serial.write(&app_host_memory[ptr], len);
    }
}
