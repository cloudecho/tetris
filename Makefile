default: build

build:
	go build -o bin/tetris  ./main
tidy:
	go mod tidy
