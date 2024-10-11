#include "./fib_wat2c_native.c"

#include <stdio.h>	

int main() {
  i64_t i = 46;
  fn_fib_ret_t $result = fn_fib(i);
  printf("fib(%d) = %lld\n", 46, $result.$R0);
  return 0;
}
