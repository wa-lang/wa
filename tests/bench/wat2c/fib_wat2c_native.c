// Auto Generated by https://wa-lang.org/wa/wat2c. DONOT EDIT!!!

#include <stdint.h>
#include <string.h>
#include <math.h>

typedef uint8_t   u8_t;
typedef int8_t    i8_t;
typedef uint16_t  u16_t;
typedef int16_t   i16_t;
typedef uint32_t  u32_t;
typedef int32_t   i32_t;
typedef uint64_t  u64_t;
typedef int64_t   i64_t;
typedef float     f32_t;
typedef double    f64_t;
typedef uintptr_t ref_t;

typedef union val_t {
  i64_t i64;
  f64_t f64;
  i32_t i32;
  f32_t f32;
  ref_t ref;
} val_t;

// func $fib (param $n i64) (result i64)
typedef struct { val_t $R0; } fn_fib_ret_t;
static fn_fib_ret_t fn_fib(val_t n);

// func fib (param $n i64) (result i64)
static fn_fib_ret_t fn_fib(val_t n) {
  fn_fib_ret_t $result;
  u32_t $R_u32;
  u16_t $R_u16;
  u8_t  $R_u8;
  val_t $R0, $R1, $R2;

  $R0 = n;
  $R1.i64 = 2;
  $R0.i32 = ((u64_t)($R0.i64)<=(u64_t)($R1.i64))? 1: 0;
  if($R0.i32) {
    $R0.i64 = 1;
    $result.$R0 = $R0;
    return $result;
  }
  $R0 = n;
  $R1.i64 = 1;
  $R0.i64 = $R0.i64 - $R1.i64;
  {
    fn_fib_ret_t $ret = fn_fib($R0);
    $R0 = $ret.$R0;
  }
  $R1 = n;
  $R2.i64 = 2;
  $R1.i64 = $R1.i64 - $R2.i64;
  {
    fn_fib_ret_t $ret = fn_fib($R1);
    $R1 = $ret.$R0;
  }
  $R0.i64 = $R0.i64 + $R1.i64;
  $result.$R0 = $R0;
  return $result;
}
