
_TARGET = TelcomSim
_SOURCE = $(shell find . -name '*.go')

_BINDIR = assets/
_ASSETS = $(shell find $(_BINDIR))
_BINDATA = bindata.go

all: $(_TARGET)

$(_BINDATA): $(_ASSETS)
	@go get github.com/shuLhan/go-bindata/cmd/go-bindata
	@go-bindata -o $(_BINDATA) $(_ASSETS)

$(_TARGET): $(_SOURCE) $(_BINDATA)
	go build .

run: $(_TARGET)
	./$(_TARGET)
