// 版权 @2025 arduino-wat 作者。保留所有权利。

#include "wa-app.h"

void setup() {
    Serial.begin(9600);
    app_init();
}

void loop() {
    app_loop();
}
