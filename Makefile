.PHONY: all build build-windows test vet clean install uninstall

BINARY_NAME=scrcpy-gui
CMD_DIR=./cmd/scrcpy-gui
DIST_DIR=dist

all: test vet build

test:
	go test -v ./...

vet:
	go vet ./...

build:
	mkdir -p $(DIST_DIR)
	go build -o $(DIST_DIR)/$(BINARY_NAME) $(CMD_DIR)

build-windows:
	mkdir -p $(DIST_DIR)
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC="zig cc -target x86_64-windows-gnu" go build -ldflags="-H=windowsgui -s -w" -o $(DIST_DIR)/$(BINARY_NAME).exe $(CMD_DIR)

install:
	./install.sh

uninstall:
	./uninstall.sh

clean:
	rm -rf $(DIST_DIR)
