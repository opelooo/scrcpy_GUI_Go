.PHONY: all build build-windows test clean

BINARY_NAME=scrcpy-gui
DIST_DIR=dist

all: test build

test:
	go test -v ./...
	go vet ./...

build:
	mkdir -p $(DIST_DIR)
	go build -o $(DIST_DIR)/$(BINARY_NAME) .

build-windows:
	mkdir -p $(DIST_DIR)
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC="zig cc -target x86_64-windows-gnu" go build -ldflags="-H=windowsgui -s -w" -o $(DIST_DIR)/$(BINARY_NAME).exe .

clean:
	rm -rf $(DIST_DIR) $(BINARY_NAME) $(BINARY_NAME).exe
