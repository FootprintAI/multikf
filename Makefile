BUILDDIR ?= "./build"
BUILDTIME=$(shell date --rfc-3339=seconds)
GITCOMMITID=$(shell git rev-parse HEAD)

tidy: 
	./gomodtidy.sh

builddocker:
	docker build -t footprintai/multikf:v1 \
		--no-cache -f Dockerfile .
	docker push footprintai/multikf:v1

windows: ## Build for Windows
	env GOOS=windows GOARCH=amd64 CGO_ENABLED=0 \
	go build -ldflags '-X "github.com/footprintai/multikf/pkg/version.BuildTime='"${BUILDTIME}"'" -X "github.com/footprintai/multikf/pkg/version.GitCommitId='"${GITCOMMITID}"'" -extldflags "-static"' --tags windows -o ${BUILDDIR}/multikf.windows.exe main.go

linux: ## Build for Linux
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
    go build -ldflags '-X "github.com/footprintai/multikf/pkg/version.BuildTime='"${BUILDTIME}"'" -X "github.com/footprintai/multikf/pkg/version.GitCommitId='"${GITCOMMITID}"'" -extldflags "-static"' --tags linux -o ${BUILDDIR}/multikf.linux main.go

darwin: ## Build for Darwin (macOS)
	env GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 \
    go build -ldflags '-X "github.com/footprintai/multikf/pkg/version.BuildTime='"${BUILDTIME}"'" -X "github.com/footprintai/multikf/pkg/version.GitCommitId='"${GITCOMMITID}"'" -extldflags "-static"' --tags darwin -o ${BUILDDIR}/multikf.darwin main.go

darwinSilicon: ## Build for Darwin Silicon (macOS M1)
	env GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 \
    go build -ldflags '-X "github.com/footprintai/multikf/pkg/version.BuildTime='"${BUILDTIME}"'" -X "github.com/footprintai/multikf/pkg/version.GitCommitId='"${GITCOMMITID}"'" -extldflags "-static"' --tags darwin -o ${BUILDDIR}/multikf.darwin-arm64 main.go


build: windows linux darwin darwinSilicon ## Build binaries
	@echo commitid: $(GITCOMMITID)
	@echo buildtime: $(BUILDTIME)

help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
