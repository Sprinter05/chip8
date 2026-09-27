# Environment
CC=go
BUILD=build

# Set the extension and compiler options
# We need to set CGO if we want to compile for windows from linux

# Running under windows
ifeq ($(OS),Windows_NT)
	OS=windows
	PREFIX=.exe
	CGO=CGO_ENABLED=1 CXX=x86_64-w64-mingw32-g++ CC=x86_64-w64-mingw32-gcc
endif

# Compiling for windows
ifeq ($(OS), windows)
	PREFIX=.exe
	CGO=CGO_ENABLED=1 CXX=x86_64-w64-mingw32-g++ CC=x86_64-w64-mingw32-gcc
endif

# If compiling for linux we remove
ifeq ($(OS), linux)
	undefine PREFIX
	undefine CGO
endif

# Executable names
NAME=chip8go

.PHONY: clean
default: $(BUILD)/$(NAME)

# Create build folder if it doesn't exist
$(BUILD):
	if ! [ -d "./$(BUILD)" ]; then mkdir $(BUILD); fi

# We check the OS environment varible for the .exe extension
$(BUILD)/$(NAME): $(BUILD)
	GOOS=$(OS) GOARCH=$(ARCH) $(CGO) \
	$(CC) build -o $(BUILD)/$(NAME) .

# Clean build folder
clean: $(BUILD)
	rm -r $(BUILD)

