// 版权 @2025 arduino-wat 作者。保留所有权利。

#include "wa-app.h"

#include <Arduino.h>

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

extern "C" int32_t app_arduino_getPinLED() {
    return LED_BUILTIN;
}

extern "C" void app_arduino_print(int32_t ptr, int32_t len) {
    Serial.write(&app_memory[ptr], len);
}
