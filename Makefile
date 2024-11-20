projects=\
	showcase

all: build test

clean:
	@rm -rf bin/*

build: $(projects)

templ:
	templ generate -include-version=false

$(projects): templ
	go build -o bin/$@ cmd/$@/main.go

test:
	go test ./...

.PHONY: all clean build test templ $(projects)
