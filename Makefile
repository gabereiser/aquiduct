PLATFORMS := linux/amd64/ linux/arm64/ windows/amd64/.exe windows/arm64/.exe darwin/amd64 darwin/arm64

temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))
ext = $(word 3, $(temp))

docker:
	GOOS=linux GOARCH=amd64 go build -o 'build/linux/server' cmd/server.go

release: $(PLATFORMS)

clean:
	@go clean
	@rm -rf build

deps:
	@go mod tidy

$(PLATFORMS):
	GOOS=$(os) GOARCH=$(arch) go build -o 'build/$(os)/$(os)-$(arch)$(ext)' cmd/server.go

.PHONY: release $(PLATFORMS)
