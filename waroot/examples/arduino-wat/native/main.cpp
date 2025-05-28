#include <stdio.h>

extern "C" void fn_main();
extern "C" void host_print_pins();

int main() {
    host_print_pins();
    printf("-----------\n");
    fn_main();
    host_print_pins();
    printf("-----------\n");
}
