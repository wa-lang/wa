#include "./fib_wat2c_native.c"

#include <stdio.h>	

int main() {
  i64_t i = 46;
  i64_t $result = fn_fib(i);
  printf("fib(%d) = %lld\n", 46, $result);
  return 0;
}
