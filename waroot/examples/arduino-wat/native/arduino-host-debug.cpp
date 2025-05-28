#include <assert.h>
#include <stdio.h>
#include <stdint.h>

typedef int32_t i32_t;

#define HIGH 0x1
#define LOW  0x0

#define INPUT 0x0
#define OUTPUT 0x1
#define INPUT_PULLUP 0x2

static i32_t kPinsNum = 16;

static i32_t pinsMode[16];
static i32_t pinsValue[16];

extern "C" void host_print_pins() {
    for(int i = 0; i < kPinsNum; i++) {
        printf("%01d: val = %d, mode = [%s]\n", i, pinsValue[i], pinsMode[i]?"OUTPUT":"INPUT");
    }
}

extern "C" i32_t host_fn_arduino_millis() {
    return 0;
}

extern "C" void host_fn_arduino_delay(i32_t ms) {
    //
}

extern "C" void host_fn_arduino_delayMicroseconds(i32_t us) {
    //
}

extern "C" void host_fn_arduino_pinMode(i32_t pin, i32_t mode) {
    assert(pin >= 0 && pin < kPinsNum);
    pinsMode[pin] = mode;
}

extern "C" i32_t host_fn_arduino_digitalRead(i32_t pin) {
    assert(pin >= 0 && pin < kPinsNum);
    return pinsValue[pin]? 1: 0;
}

extern "C" void host_fn_arduino_digitalWrite(i32_t pin, i32_t value) {
    assert(pin >= 0 && pin < kPinsNum);
    pinsValue[pin] = value? 1: 0;
}

extern "C" i32_t host_fn_arduino_analogRead(i32_t pin) {
    assert(pin >= 0 && pin < kPinsNum);
    return pinsValue[pin];
}

extern "C" void host_fn_arduino_analogWrite(i32_t pin, i32_t value) {
    assert(pin >= 0 && pin < kPinsNum);
    pinsValue[pin] = value;
}

extern "C" i32_t host_fn_arduino_getPinLED() {
    return 13;
}

extern "C" void host_fn_arduino_print(i32_t ptr, i32_t len) {
    extern uint8_t wasm_memory[];
    for(int i = 0; i < len; i++) {
        putchar(wasm_memory[ptr+i]);
    }
}
