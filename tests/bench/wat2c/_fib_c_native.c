
#include <stdio.h>	
#include <stdint.h>

int64_t fn_fib(int64_t n) {
  if (n <= 2) {
	return 1;
  }
  return fn_fib(n - 1) + fn_fib(n - 2);
}

int main() {
  int64_t v = fn_fib(46);
  printf("fib(%d) = %lld\n", 46, v);
  return 0;
}
