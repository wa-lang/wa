WA:=go run ../../../../main.go

# todo: gen host func

default:
	$(WA) build code.wa
	$(WA) wat2c -o wa-app.c code.wat
	go run .

clean:
