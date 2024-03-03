# Written by yijian on 2024/03/03

all: sql2struct

sql2struct: main.go
	go build -o $@ $<

install: sql2struct
	cp sql2struct $$GOPATH/bin/

.PHONY: clean tidy

clean:
	rm -f sql2struct

tidy:
	go mod tidy
