#include "./fib_wat2c_native.c"

#include <stdio.h>	

int main() {
  val_t $args[1];
  $args[0].i64 = 46;
  fn_fib_ret_t $result = fn_fib($args[0]);
  printf("fib(%d) = %lld\n", 46, $result.$R0.i64);
  return 0;
}
