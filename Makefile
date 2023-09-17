GO := go

MAIN_FILE := cmd/hina/main.go

all: build

build: 
	$(GO) build -o hina $(MAIN_FILE)