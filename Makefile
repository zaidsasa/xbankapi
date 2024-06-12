TOOL_MIGRATE=go run -tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate
TOOL_GOLANGCI_LINT=go run github.com/golangci/golangci-lint/cmd/golangci-lint

.PHONY: generate
generate: 
	go generate

.PHONY: vet
vet:
	go vet ./...

.PHONY: lint
lint:
	$(TOOL_GOLANGCI_LINT) run

.PHONY: test
test:
	go test ./... --race

.PHONY: build
build:
	go build -o ./out/

.PHONY: migrate-up
migrate-up:
	$(call migrate, up)

.PHONY: migrate-down
migrate-down:
	$(call migrate, down)

.PHONY: migrate-create
migrate-create:
	$(TOOL_MIGRATE) create -dir ./db/migrations/ -ext sql ${name}

define migrate
	$(if $(filter postgres://%,	$(DATABASE_URL)),, $(error DATABASE_URL is not set))
	$(TOOL_MIGRATE) -source file://./db/migrations/ -database ${DATABASE_URL} $1
endef
