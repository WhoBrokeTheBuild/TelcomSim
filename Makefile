_CONFIG ?= debug
_FLAGS  += -tags $(_CONFIG)

_ASSET_DIR = data/assets/
_BINDATA  = data/bindata.go
_BINFLAGS = -pkg data -prefix $(_ASSET_DIR)

ifeq ($(_CONFIG),debug)
	_BINFLAGS += -debug
endif

ifeq ($(OS),Windows_NT)
	ifeq ($(_CONFIG),release)
		_FLAGS += -ldflags -H=windowsgui
	endif
	_TARGET = TelcomSim.exe
	_SOURCE = $(shell powershell -Command "Get-ChildItem -Filter '*.go' . | Select -exp FullName")
	_ASSETS = $(shell powershell -Command "Get-ChildItem -File -Recurse $(_ASSET_DIR) | Select -exp FullName")
	_CLEAN_CMD = del /f $(_TARGET) $(_BINDATA) 2>NUL
else
	_TARGET = TelcomSim
	_SOURCE = $(shell find . -name '*.go')
	_ASSETS = $(shell find $(_ASSET_DIR))
	_CLEAN_CMD = rm -f $(_TARGET) $(_BINDATA)
endif

all: $(_TARGET)

clean:
	$(_CLEAN_CMD)

$(_BINDATA): $(_ASSETS)
	@go get github.com/shuLhan/go-bindata/cmd/go-bindata
	go-bindata $(_BINFLAGS) -o $(_BINDATA) $(_ASSET_DIR)...

$(_TARGET): $(_SOURCE) $(_BINDATA)
	go build -v $(_FLAGS) -o $(_TARGET) .

run: $(_TARGET)
	./$(_TARGET)
