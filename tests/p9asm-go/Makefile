linux:
	go tool asm -o main_linux_amd64.o main_linux_amd64.s
	go tool nm main_linux_amd64.o | grep main.main
	go tool link -H linux -o a.out main_linux_amd64.o

	./hello
	echo $?

windows:

clean:
