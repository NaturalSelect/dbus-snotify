# dbus-snotify Makefile
#
GOMOD=on
default: all

phony := all
all: build

phony += build
build:
	go build -o build/dbus-snotify

phony += test
test:
	go test -v ./...

phony += image
image:
	docker build -t naturalselect/dbus-snotifypod:latest .

.PHONY: $(phony)
