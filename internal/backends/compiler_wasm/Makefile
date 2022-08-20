# 版权 @2022 凹语言 作者。保留所有权利。

default: base.wasm
	cd ../../../tools/wa-wasmer-run && go install
	wa-wasmer-run -file=base.wasm

base.wasm: base.wat
	wat2wasm --debug-names -o base.wasm base.wat

clean:
	rm base.wasm
