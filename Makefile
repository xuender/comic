PACKAGE = github.com/xuender/comic

tools:
	go install github.com/google/wire/cmd/wire@latest
	go install fyne.io/fyne/v2/cmd/fyne@latest
	go install github.com/spf13/cobra-cli@latest
	go install github.com/cosmtrek/air@latest

wire: tools
	go mod tidy
	wire gen ${PACKAGE}/cmd

dev: tools
	air

build: tools
	fyne package -os linux

lint:
	golangci-lint run --timeout 60s --max-same-issues 50 ./...

lint-fix:
	golangci-lint run --timeout 60s --max-same-issues 50 --fix ./...
