#include "wa-app.h"

void setup() {
    Serial.begin(9600);
    app_init();
}

void loop() {
    app_loop();
}
