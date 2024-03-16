# Written by yijian on 2024/03/03

all: sql2struct

sql2struct: main.go
ifeq ($(OS),Windows_NT)
	set GOOS=windows
	set GOARCH=amd64
	go build -o sql2struct.exe main.go
else
	go build -o sql2struct main.go
endif

install: sql2struct
ifeq ($(OS),Windows_NT)
	copy sql2struct.exe %GOPATH%\bin\
else
	cp sql2struct $$GOPATH/bin/
endif

.PHONY: clean tidy

clean:
ifeq ($(OS),Windows_NT)
	del sql2struct.exe
else
	rm -f sql2struct
endif

tidy:
	go mod tidy
