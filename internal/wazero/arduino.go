// 版权 @2022 凹语言 作者。保留所有权利。

package wazero

import (
	"context"
	"fmt"
	"time"

	"wa-lang.org/wa/internal/3rdparty/wazero"
	"wa-lang.org/wa/internal/3rdparty/wazero/api"
	"wa-lang.org/wa/internal/config"
)

func ArduinoInstantiate(ctx context.Context, rt wazero.Runtime) (api.Closer, error) {
	startTime := time.Now()
	return rt.NewHostModuleBuilder(config.WaOS_arduino).
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, m api.Module) int32 {
			t := time.Now().Sub(startTime).Milliseconds()
			fmt.Printf("arduino.millis(): %v\n", t)
			return int32(t)
		}).
		WithParameterNames().
		Export("millis").
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, ms uint32) {
			fmt.Printf("arduino.delay(%d)...\n", ms)
			time.Sleep(time.Millisecond * time.Duration(ms))
		}).
		WithParameterNames("ms").
		Export("delay").
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, pin, mode int32) {
			switch mode {
			case 0:
				fmt.Printf("arduino.pinMode(%d, %s)\n", pin, "INPUT")
			case 1:
				fmt.Printf("arduino.pinMode(%d, %s)\n", pin, "OUTPUT")
			default:
				fmt.Printf("arduino.pinMode(%d, %s)\n", pin, "INPUT_PULLUP")
			}
		}).
		WithParameterNames("pin", "mode").
		Export("pinMode").
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, pin, value int32) {
			if value == 0 {
				fmt.Printf("arduino.digitalWrite(%d, %s)\n", pin, "LOW")
			} else {
				fmt.Printf("arduino.digitalWrite(%d, %s)\n", pin, "HIGH")
			}
		}).
		WithParameterNames("pin", "value").
		Export("digitalWrite").
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context) int32 {
			const pin = 13
			fmt.Printf("arduino.getPinLED(): %v\n", pin)
			return pin
		}).
		WithParameterNames().
		Export("getPinLED").
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, m api.Module, ptr, len uint32) {
			bytes, _ := m.Memory().Read(ctx, ptr, len)
			fmt.Printf("arduino.print(%q)\n", string(bytes))
		}).
		WithParameterNames("ptr", "len").
		Export("print").
		Instantiate(ctx, rt)
}
