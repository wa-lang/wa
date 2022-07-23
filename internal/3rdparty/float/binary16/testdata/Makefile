C=$(wildcard *.c)
LL=$(C:.c=.ll)

all: $(LL)

%.ll: %.c sar
	clang -S -emit-llvm -o $@ $<
	sar -i "(?m:^[^@\n][^\n]*[\n])" "" $@
	sar -i ", align[^\n]*" "" $@
	sar -i "[\n]+" "\n" $@
	sar -i "^[\n]+" "" $@

sar:
	@if ! which $@ &> /dev/null ; then \
		echo "Please install the \"sar\" tool"; \
		echo; \
		echo "   go get -u github.com/mewkiz/cmd/sar"; \
		echo; \
		exit 1; \
	fi

.PHONY: all
