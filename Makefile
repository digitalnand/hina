GO := go

MAIN_FILE := cmd/hina/main.go

all: main

main: 
	$(GO) run $(MAIN_FILE)