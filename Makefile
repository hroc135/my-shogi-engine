EXECUTABLE_FILE = "my-shogi-engine.exe"

build:
	go build ./

run: build
	./$(EXECUTABLE_FILE)
