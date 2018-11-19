_CONFIG ?= debug
_FLAGS  += -tags $(_CONFIG)

ifeq ($(_CONFIG),release)
	_BINDIR = assets
	_BINDATA = bindata.go
	ifeq ($(OS),Windows_NT)
		_FLAGS += -ldflags -H=windowsgui
		_ASSETS = $(shell powershell -Command "Get-ChildItem -File -Recurse $(_BINDIR) | Select -exp FullName")
	else
		_ASSETS = $(shell find $(_BINDIR))
	endif
endif

ifeq ($(OS),Windows_NT)
	_TARGET = TelcomSim.exe
	_SOURCE = $(shell powershell -Command "Get-ChildItem -Filter '*.go' . | Select -exp FullName")
else
	_TARGET = TelcomSim
	_SOURCE = $(shell find . -name '*.go')
endif

all: $(_TARGET)

ifeq ($(_CONFIG),release)
$(_BINDATA): $(_ASSETS)
	@go get github.com/shuLhan/go-bindata/cmd/go-bindata
	@go-bindata -o $(_BINDATA) $(_BINDIR)/...
endif

$(_TARGET): $(_SOURCE) $(_BINDATA)
	go build $(_FLAGS) -o $(_TARGET) .

run: $(_TARGET)
	./$(_TARGET)
