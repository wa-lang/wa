// 导入 Wasm3 头文件
//
// Wasm3 是 Volodymyr Shymanskyy, Steven Massey 开发的 WebAssembly 解释器,
// 可以在 Arduino Nano 33 平台执行 wasm 程序.
#include <wasm3.h>
#include <m3_env.h>

// 定义 WASM 和 本地栈大小
#define WASM_STACK_SLOTS  1024
#define NATIVE_STACK_SIZE (32*1024)

// WASM 最大内存限制, 一般不得超过 64KB 大小
#define WASM_MEMORY_LIMIT (32*1024)

// 导入 凹语言 生成的 WASM 文件对应的二进制头文件
#include "app.wasm.h"

// m3ApiRawFunction 宏用于 Arduino 宿主 API 函数绑定
//
// 每个函数必须通过以下 3 种方式之一返回:
//	m3ApiReturn(val) - 返回一个值
//	m3ApiSuccess()   - 返回, void
//	m3ApiTrap(trap)  - 返回, 异常

// fn arduino.millis() => u32
m3ApiRawFunction(m3_arduino_millis) {
	m3ApiReturnType(uint32_t)
	m3ApiReturn(millis());
}

// fn arduino.delay(ms: u32)
m3ApiRawFunction(m3_arduino_delay) {
	m3ApiGetArg(uint32_t, ms)
	delay(ms);
	m3ApiSuccess();
}

// fn arduino.pinMode(pin: u32, mode: u32)
m3ApiRawFunction(m3_arduino_pinMode) {
	m3ApiGetArg(uint32_t, pin)
	m3ApiGetArg(uint32_t, mode)

	switch(mode) {
	case 0: pinMode(pin, INPUT); break;
	case 1: pinMode(pin, OUTPUT); break;
	case 2: pinMode(pin, INPUT_PULLUP); break;
	default: pinMode(pin, INPUT); break;
	}

	m3ApiSuccess();
}

// fn arduino.digitalWrite(pin: u32, value: u32)
m3ApiRawFunction(m3_arduino_digitalWrite) {
	m3ApiGetArg(uint32_t, pin)
	m3ApiGetArg(uint32_t, value)

	digitalWrite(pin, value);
	m3ApiSuccess();
}

// fn arduino.getPinLED() => u32
m3ApiRawFunction(m3_arduino_getPinLED) {
	m3ApiReturnType(uint32_t)
	m3ApiReturn(LED_BUILTIN);
}

// fn arduino.print(ptr: *u8, len: u32)
m3ApiRawFunction(m3_arduino_print) {
	m3ApiGetArgMem(const uint8_t*, buf)
	m3ApiGetArg(uint32_t, len)

	Serial.write(buf, len);
	m3ApiSuccess();
}

// 初始化 Arduino 宿主 API
M3Result LinkArduino(IM3Runtime rt) {
	m3_LinkRawFunction (rt->modules, "arduino", "millis",       "i()",   &m3_arduino_millis);
	m3_LinkRawFunction (rt->modules, "arduino", "delay",        "v(i)",  &m3_arduino_delay);
	m3_LinkRawFunction (rt->modules, "arduino", "pinMode",      "v(ii)", &m3_arduino_pinMode);
	m3_LinkRawFunction (rt->modules, "arduino", "digitalWrite", "v(ii)", &m3_arduino_digitalWrite);
	m3_LinkRawFunction (rt->modules, "arduino", "getPinLED",    "i()",   &m3_arduino_getPinLED);
	m3_LinkRawFunction (rt->modules, "arduino", "print",        "v(*i)", &m3_arduino_print);
	return m3Err_none;
}

// 阻塞执行 wasm 程序
void wasm_task(void* ctx) {
	M3Result result = m3Err_none;

	IM3Environment env = m3_NewEnvironment ();
	if(!env) {
		Serial.print("Fatal: NewEnvironment failed");
		return;
	}
	
	IM3Runtime runtime = m3_NewRuntime (env, WASM_STACK_SLOTS, NULL);
	if(!runtime) {
		Serial.print("Fatal: NewRuntime failed");
		return;
	}

	if(WASM_MEMORY_LIMIT > 0) {
		runtime->memoryLimit = WASM_MEMORY_LIMIT;
	}

	IM3Module module;
	result = m3_ParseModule (env, &module, app_wasm, app_wasm_len);
	if (result) {
		Serial.print("Fatal: ParseModule ");
		Serial.print(result);
		return;
	}

	result = m3_LoadModule (runtime, module);
	if(result) {
		Serial.print("Fatal: LoadModule ");
		Serial.print(result);
		return;
	}

	result = LinkArduino (runtime);
	if(result) {
		Serial.print("Fatal: LinkArduino ");
		Serial.print(result);
		return;
	}

	IM3Function fn_start;
	result = m3_FindFunction (&fn_start, runtime, "_start");
	if (result)  {
		Serial.print("Fatal: FindFunction ");
		Serial.print(result);
		return;
	}

	Serial.println("Running Wa/WebAssembly...");
	result = m3_CallV(fn_start); // loop

	// 执行到这里说明有错误
	if (result) {
		M3ErrorInfo info;
		m3_GetErrorInfo (runtime, &info);
		Serial.print("Error44: ");
		Serial.print(result);
		Serial.print(" (");
		Serial.print(info.message);
		Serial.println(")");
		if (info.file && strlen(info.file) && info.line) {
			Serial.print("At ");
			Serial.print(info.file);
			Serial.print(":");
			Serial.println(info.line);
		}
	}
}

// setup 作为 main 函数用户
void setup() {
	// 串口初始化
	Serial.begin(115200);
	delay(100);

	// 等待串口初始化完成, 必须是 USB 串口
	while(!Serial) {}

	// 打印提示信息
	Serial.println("\nWasm3 v" M3_VERSION " (" M3_ARCH "), build " __DATE__ " " __TIME__);

	// 阻塞执行 wasm 程序, 不会返回
	wasm_task(NULL);
}

// 该函数不会被执行
// 定义该函数只是为了确保 Arduino 编译通过
void loop() {
	delay(100);
}
