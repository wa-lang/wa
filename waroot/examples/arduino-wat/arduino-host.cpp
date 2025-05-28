#include <Arduino.h>

#include <stdint.h>

typedef int32_t i32_t;

extern "C" i32_t host_fn_arduino_millis() {
    return millis();
}

extern "C" void host_fn_arduino_delay(i32_t ms) {
    delay(ms);
}

extern "C" void host_fn_arduino_delayMicroseconds(i32_t us) {
    delayMicroseconds(us);
}

extern "C" void host_fn_arduino_pinMode(i32_t pin, i32_t mode) {
    pinMode(pin, mode);
}

extern "C" i32_t host_fn_arduino_digitalRead(i32_t pin) {
    return digitalRead(pin);
}

extern "C" void host_fn_arduino_digitalWrite(i32_t pin, i32_t value) {
    digitalWrite(pin, value);
}

extern "C" i32_t host_fn_arduino_analogRead(i32_t pin) {
    return analogRead(pin);
}

extern "C" void host_fn_arduino_analogWrite(i32_t pin, i32_t value) {
    analogWrite(pin, value);
}

extern "C" i32_t host_fn_arduino_getPinLED() {
    return LED_BUILTIN;
}

extern "C" void host_fn_arduino_print(i32_t ptr, i32_t len) {
    extern uint8_t wasm_memory[];
    Serial.write(&wasm_memory[ptr], len);
}
