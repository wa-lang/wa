default:
	go run main.go > a.out.ll
	clang -Wno-override-module a.out.ll
	./a.out

	-@rm *.out

clean:
	-rm *.out*
