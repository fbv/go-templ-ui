projects=\
	showcase

all: build test

clean:
	@rm -rf bin/*

build: $(projects)

templ:
	templ generate -include-version=false

$(projects): templ tailwind
	go build -o bin/$@ cmd/$@/main.go

test:
	go test ./...

tailwind-watch:
	./tailwindcss -i ./view/input.css -o ./view/style.css --watch

tailwind:
	./tailwindcss -i ./view/input.css -o ./view/style.css --minify

.PHONY: all clean build test templ tailwind-watch tailwind $(projects)
