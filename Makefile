# dbus-snotify Makefile
#
GOMOD=on
default: all

phony := all
all: build

phony += build
build:
	go build -o build/dbus-snotify

.PHONY: $(phony)
