#include "wa-app.h"

extern "C" void host_print_pins();

int main() {
    wasm_init();
    host_print_pins();
    return 0;
}
