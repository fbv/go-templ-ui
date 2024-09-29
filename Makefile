projects=\
	showcase

all: build test

clean:
	@rm -rf bin/*

build: $(projects)

$(projects):
	go build -o bin/$@ cmd/$@/main.go

test:
	go test ./...

.PHONY: all clean build test $(projects)
