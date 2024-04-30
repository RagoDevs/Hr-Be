## help: print this help message
.PHONY: help
help:
	@echo "Usage:"
	@sed -n "s/^##//p" ${MAKEFILE_LIST} | column -t -s ":" | sed -e "s/^/ /"

.PHONY: confirm
confirm:
	@echo -n "Are you sure? [y/N] " && read ans && [ $${ans:-N} = y ]

## migrations/up: apply all up database migrations
.PHONY: migrations/up
migrations/up: confirm
	migrate -path ./db/migrations -database ${HR_DB_DSN} --verbose up


## migrations/down: apply all down database migrations
.PHONY: migrations/down
migrations/down: confirm
	migrate -path ./db/migrations -database ${HR_DB_DSN} --verbose down


## migrations/new name=$1: create a new database migration
.PHONY: migrations/new
migrations/new:
	@echo "Creating migration files for ${name}..."
	migrate create -ext sql -dir db/migrations -seq  ${name}


## sqli: sqlc init 
.PHONY: sqli
sqli:
	sqlc init

## sqlg: sqlc generate
.PHONY: sqlg
sqlg:		
	sqlc generate


## server: Run the fucking server 
.PHONY: server
server:
	go run ./cmd/api

## test: Run the fucking tests
test:
	go test -v -count=1 ./...



## audit: tidy dependencies and format, vet and test all code
.PHONY: audit
audit: vendor
	@echo "Formatting code..."
	go fmt ./...
	@echo "Vetting code..."
	go vet ./...
	staticcheck ./...
	@echo "Running tests..."
	go test -race -vet=off ./...


## vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo "Tidying and verifying module dependencies..."
	go mod tidy
	go mod verify
	@echo "Vendoring dependencies..."
	go mod vendor





