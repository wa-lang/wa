# 版权 @2022 凹语言 作者。保留所有权利。

default: base.wasm
	cd ../../../tools/wa-wasmer-run && go install
	wa-wasmer-run -file=base.wasm

base.wasm: base.wat
	wat2wasm base.wat -o base.wasm

clean:
	rm base.wasm
