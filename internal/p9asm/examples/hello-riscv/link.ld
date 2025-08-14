ENTRY(_start)

SECTIONS
{
  . = 0x80000000;    /* QEMU virt 默认从0x80000000开始执行 */

  .text : {
    *(.text*)
  }

  .rodata : {
    *(.rodata*)
  }

  .data : {
    *(.data*)
  }

  .bss : {
    *(.bss*)
    *(COMMON)
  }
}
