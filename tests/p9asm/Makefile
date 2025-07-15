WA:=go run ../../main.go

default:
	$(WA) p9asm -o hello.o hello_windows.s

link:
	$(WA) p9link -H windows -o hello.exe hello.o

clean:
