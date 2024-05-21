TOOLS_DIR:=$(CURDIR)/tools/

PHONY: generate
generate: install-tool-sqlc install-tool-mockery
	./tools/sqlc generate
	mockery

PHONY: lint
lint:
	go vet
	golangci-lint run

PHONY: test
test:
	go test ./...

PHONY: build
build:
	go build -o ./out/

PHONY: migrate-%
migrate-%: install-tools-migrate
ifndef DATABASE_URL
	$(error DATABASE_URL is not set)
endif
	$(call migrate, $*)

PHONY: migrate-create
migrate-create: install-tools-migrate
	$(TOOL_MIGRATE) create -dir ./db/migrations/ -ext sql ${name}

TOOL_MIGRATE="${TOOLS_DIR}migrate"
PHONY: install-tools-migrate
install-tools-migrate:
	$(call install_go_tool, ${TOOL_MIGRATE}, -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17)

TOOL_SQLC="${TOOLS_DIR}sqlc"
PHONY: install-tool-sqlc
install-tool-sqlc:
	$(call install_go_tool, ${TOOL_SQLC}, github.com/sqlc-dev/sqlc/cmd/sqlc@v1.26)

TOOL_MOCKERY="${TOOLS_DIR}mockery"
PHONY: install-tool-mockery
install-tool-mockery:
	$(call install_go_tool, ${TOOL_MOCKERY}, github.com/vektra/mockery/v2@v2.42.3)

define install_go_tool
@which $1 >/dev/null 2>&1 || echo "installing go tool: $1" && GOBIN=${CURDIR}/tools go install $2
endef

define migrate
	$(TOOL_MIGRATE) -source file://./db/migrations/ -database ${DATABASE_URL} $1
endef
