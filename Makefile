GOCMD=CGO_ENABLED=0 go
GOFLAGS=-a -ldflags '-w -extldflags "-static"' -mod vendor

.PHONY: all ui sds sds-ui docker \
	fmt vet lint test test-race test-fuzz check clean

# --- Build targets ---

all: ui sds-ui

ui:
	npm --prefix ./ui/webapp ci
	npm --prefix ./ui/webapp run build

sds:
	$(GOCMD) build $(GOFLAGS) cmd/sds/sigplot_data_service.go

sds-ui: ui
	$(GOCMD) build $(GOFLAGS) -tags ui cmd/sds/sigplot_data_service.go

docker:
	docker build -t sds:0.7 .

# --- Code quality targets ---

fmt:
	@echo "==> Running gofmt..."
	@gofmt -l -w $$(find . -name '*.go' -not -path './vendor/*')

vet:
	@echo "==> Running go vet..."
	go vet ./...

lint:
	@echo "==> Running golangci-lint..."
	golangci-lint run ./...

# --- Test targets ---

test:
	@echo "==> Running tests..."
	go test ./... -count=1

test-race:
	@echo "==> Running tests with race detector..."
	go test ./... -race -count=1

test-fuzz:
	@echo "==> Running fuzz tests (10s each)..."
	go test ./internal/bluefile/... -fuzz=FuzzConvertFileData -fuzztime=10s
	go test ./internal/image/... -fuzz=FuzzTransform -fuzztime=10s

# --- Combined targets ---

check: fmt vet lint test
	@echo "==> All checks passed."

clean:
	@rm -f sigplot_data_service
	@rm -rf sdscache/
