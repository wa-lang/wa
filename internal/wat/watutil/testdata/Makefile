# 版权 @2024 凹语言 作者。保留所有权利。

# wabt-1.0.29

WAT_FILES=$(wildcard ./*.wat)
WAT_WASM_FILES=$(patsubst %.wat,%.wat.wasm,$(WAT_FILES))
WAT_WASM_NONAME_FILES=$(patsubst %.wat,%.wat.noname.wasm,$(WAT_FILES))

default: $(WAT_WASM_FILES) $(WAT_WASM_NONAME_FILES)
	@echo "ok"

clean:
	-rm *.wat.wasm *.wasm

%.wat.wasm: %.wat
	wat2wasm --debug-names $< -o $@

%.wat.noname.wasm: %.wat
	wat2wasm $< -o $@
