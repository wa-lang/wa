;; 版权 @2025 arduino-wat 作者。保留所有权利。

(module $wa-arduino-hello
	(import "arduino" "getPinLED" (func $getPinLED (result i32)))
	(import "arduino" "pinMode" (func $pinMode (param $pin i32) (param $mode i32)))
	(import "arduino" "digitalWrite" (func $digitalWrite (param $pin i32) (param $value i32)))
	(import "arduino" "delay" (func $delay (param $ms i32)))
	(import "arduino" "delayMicroseconds" (func $delayMicroseconds (param $us i32)))

	(memory (;0;) 0)
	(export "memory" (memory 0))

	(global $HIGH i32 (i32.const 1))
	(global $LOW  i32 (i32.const 0))

	(global $INPUT  i32 (i32.const 0))
	(global $OUTPUT i32 (i32.const 1))

	(global $LED_BUILTIN i32 (i32.const 13))

	;; +---------------------------------------------------------------------+
	;; | LCD1602 Module                                                      |
	;; |-----+-----------------------+---------------------+-----------------+
	;; | Pin | Label | Connected To  | Description         | In Code         |
	;; |-----+-------+---------------+---------------------+-----------------+
	;; |  1  | GND   | GND           | Ground              | -               |
	;; |  2  | VCC   | 5V            | Power Supply        | -               |
	;; |  3  | VO    | Potentiometer | Contrast Control    | -               |
	;; |  4  | RS    | D7            | Register Select     | digitalWrite(RS)|
	;; |  5  | RW    | GND           | Write Only (GND)    | always LOW      |
	;; |  6  | E     | D6            | Enable Signal       | pulseEnable()   |
	;; | 11  | D4    | D5            | Data Bit 4          | write4bits()    |
	;; | 12  | D5    | D4            | Data Bit 5          | write4bits()    |
	;; | 13  | D6    | D3            | Data Bit 6          | write4bits()    |
	;; | 14  | D7    | D2            | Data Bit 7          | write4bits()    |
	;; +-----+-------+---------------+---------------------+-----------------+

	;; 引脚定义
	(global $RS i32 (i32.const 7))
	(global $E  i32 (i32.const 6))
	(global $D4 i32 (i32.const 5))
	(global $D5 i32 (i32.const 4))
	(global $D6 i32 (i32.const 3))
	(global $D7 i32 (i32.const 2))

	(func $lcdPulseEnable
		;; digitalWrite(E, LOW);
		global.get $E
		global.get $LOW
		call $digitalWrite

		;; delayMicroseconds(1);
		i32.const 1
		call $delayMicroseconds

		;; digitalWrite(E, HIGH);
		global.get $E
		global.get $HIGH
		call $digitalWrite
		
		;; delayMicroseconds(1);
		i32.const 1
		call $delayMicroseconds

		;; digitalWrite(E, LOW);
		global.get $E
		global.get $LOW
		call $digitalWrite

		;; delayMicroseconds(100); // 等待命令执行
		i32.const 100
		call $delayMicroseconds
	)

	;; 取 1bit
	(func $getBitAt (param $v i32) (param $i i32) (result i32)
		;; result = (v >> i) & 1
		local.get $v
		local.get $i
		i32.shr_u
		i32.const 1
		i32.and
		return
	)

	;; 取 4bit
	(func $get4BitAt (param $v i32) (param $i i32) (result i32)
		;; result = (v >> i) & 1
		local.get $v
		local.get $i
		i32.shr_u
		i32.const 0xF
		i32.and
		return
	)

	(func $lcdWrite4bits (param $byteValue i32)
		;; digitalWrite(D4, (value >> 0) & 0x01);
		global.get $D4
		local.get $byteValue
		i32.const 0
		call $getBitAt
		call $digitalWrite
		
		;; digitalWrite(D5, (value >> 1) & 0x01);
		global.get $D5
		local.get $byteValue
		i32.const 1
		call $getBitAt
		call $digitalWrite

		;; digitalWrite(D6, (value >> 2) & 0x01);
		global.get $D6
		local.get $byteValue
		i32.const 2
		call $getBitAt
		call $digitalWrite

		;; digitalWrite(D7, (value >> 3) & 0x01);
		global.get $D7
		local.get $byteValue
		i32.const 3
		call $getBitAt
		call $digitalWrite

		;; lcdPulseEnable();
		call $lcdPulseEnable
	)

	(func $lcdSend (param $value i32) (param $mode i32)
		;; digitalWrite(RS, mode);
		global.get $RS
		local.get $mode
		call $digitalWrite
		
		;; lcdWrite4bits(value >> 4);   // 高四位
		local.get $value
		i32.const 4
		call $get4BitAt
		call $lcdWrite4bits

		;; lcdWrite4bits(value & 0x0F); // 低四位
		local.get $value
		i32.const 0
		call $get4BitAt
		call $lcdWrite4bits
	)

	(func $lcdCommand (param $value i32)
		;; lcdSend(value, LOW);
		local.get $value
		global.get $LOW
		call $lcdSend
	)

	(func $LCDSetCursor (param $row i32) (param $col i32)
		;; if (row == 0) {
		;;   lcdCommand(0x80 + col);
		;; } else {
		;;   lcdCommand(0xC0 + col);
		;; }
		local.get $row
		i32.eqz
		if (result i32)
			i32.const 0x80
		else
			i32.const 0xC0
		end
		local.get $col
		i32.add
		call $lcdCommand
	)

	(func $LCDWriteChar (param $value i32)
		;; lcdSend(value, HIGH);
		local.get $value
		global.get $HIGH
		call $lcdSend
	)

	(func $LCDClear
		;; lcdCommand(0x01);
		i32.const 0x01
		call $lcdCommand

		;; delay(2);
		i32.const 2
		call $delay
	)

	(func $LCDInit
		;; pinMode(RS, OUTPUT);
		global.get $RS
		global.get $OUTPUT
		call $pinMode

		;; pinMode(E, OUTPUT);
		global.get $E
		global.get $OUTPUT
		call $pinMode

		;; pinMode(D4, OUTPUT);
		global.get $D4
		global.get $OUTPUT
		call $pinMode

		;; pinMode(D5, OUTPUT);
		global.get $D5
		global.get $OUTPUT
		call $pinMode

		;; pinMode(D6, OUTPUT);
		global.get $D6
		global.get $OUTPUT
		call $pinMode

		;; pinMode(D7, OUTPUT);
		global.get $D7
		global.get $OUTPUT
		call $pinMode

		;; delay(50); // 等待LCD启动
		i32.const 50
		call $delay

		;; 初始化到4-bit模式

		;; lcdWrite4bits(0x03);
		i32.const 0x03
		call $lcdWrite4bits

		;; delay(5);
		i32.const 5
		call $delay

		;; lcdWrite4bits(0x03);
		i32.const 0x03
		call $lcdWrite4bits

		;; delayMicroseconds(150);
		i32.const 150
		call $delayMicroseconds

		;; lcdWrite4bits(0x03);
		i32.const 0x03
		call $lcdWrite4bits

		;; write4bits(0x02); // 设置4-bit模式
		i32.const 0x02
		call $lcdWrite4bits

		;; 几个基本设置

		;; lcdCommand(0x28); // 4-bit, 2行, 5x8 点阵
		i32.const 0x28
		call $lcdCommand

		;; lcdCommand(0x08); // 显示关闭
		i32.const 0x08
		call $lcdCommand

		;; lcdCommand(0x01); // 清屏
		i32.const 0x01
		call $lcdCommand

		;; delay(2);
		i32.const 2
		call $delay

		;; lcdCommand(0x06); // 输入模式：写入后光标右移
		i32.const 0x06
		call $lcdCommand

		;; lcdCommand(0x0C); // 显示开启，光标关闭
		i32.const 0x0C
		call $lcdCommand
	)

	(func $_start (start)
		;; LCD
		call $LCDInit

		;; LCDSetCursor(0, 1)
		i32.const 0
		i32.const 1
		call $LCDSetCursor

		;; LCDWriteChar('A')
		i32.const 'h'
		call $LCDWriteChar
		i32.const 'e'
		call $LCDWriteChar
		i32.const 'l'
		call $LCDWriteChar
		i32.const 'l'
		call $LCDWriteChar
		i32.const 'o'
		call $LCDWriteChar
		i32.const ' '
		call $LCDWriteChar
		i32.const 'w'
		call $LCDWriteChar
		i32.const 'a'
		call $LCDWriteChar
		i32.const '-'
		call $LCDWriteChar
		i32.const 'l'
		call $LCDWriteChar
		i32.const 'a'
		call $LCDWriteChar
		i32.const 'n'
		call $LCDWriteChar
		i32.const 'g'
		call $LCDWriteChar
		i32.const '!'
		call $LCDWriteChar
	)

	(func $loop (export "loop")
		;; delay(100)
		i32.const 100
		call $delay
	)
)
