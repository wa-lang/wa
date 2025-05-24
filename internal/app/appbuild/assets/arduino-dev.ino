#include "wa-app.h"

void setup() {
    pinMode(LED_BUILTIN, OUTPUT);
    // wasm_init();
}

void loop() {
    digitalWrite(LED_BUILTIN, HIGH);
    delay(1000);
    digitalWrite(LED_BUILTIN, LOW);
    delay(1000);
}
