PACKAGE = github.com/xuender/comic

tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/cespare/reflex@latest
	go install github.com/rakyll/gotest@latest
	go install github.com/psampaz/go-mod-outdated@latest
	go install github.com/jondot/goweight@latest
	go install github.com/sonatype-nexus-community/nancy@latest
	go install golang.org/x/tools/cmd/cover@latest
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

outdated: tools
	go list -u -m -json all | go-mod-outdated -update -direct

audit: tools
	go list -json -m all | nancy sleuth

coverage:
	go test -v -gcflags=all=-l -coverprofile=cover.out -covermode=atomic ./...
	go tool cover -html=cover.out -o cover.html

test:
	go test -race -v ./... -gcflags=all=-l

watch-test:
	reflex -t 50ms -s -- sh -c 'gotest -v ./...'

bench:
	go test -benchmem -count 3 -bench ./...

weight: tools
	goweight
