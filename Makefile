GO := go

MAIN_FILE := cmd/hina/main.go
BIN_FOLDER := bin/hina

all: build

build: | ${BIN_FOLDER}
	$(GO) build -o ${BIN_FOLDER}/hina $(MAIN_FILE)
${BIN_FOLDER}:
	mkdir -p ${BIN_FOLDER}