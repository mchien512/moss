PROTO_SRC := protobuf

PROTO_DST := go/internal/genproto

PROTO_FILES := $(shell find $(PROTO_SRC) -name "*.proto")

PATH := $(PATH):$(shell go env GOPATH)/bin

# ─────────────────────────────────────────────────────────────────────────────
# “make proto”: regenerate all .proto → .pb.go + _grpc.pb.go under $(PROTO_DST)
# ─────────────────────────────────────────────────────────────────────────────
.PHONY: proto
proto:
	@echo "Regenerating Go stubs from all .proto files…"
	@for proto_file in $(PROTO_FILES); do \
		rel_path=$${proto_file#$(PROTO_SRC)/}; \
		mkdir -p $(PROTO_DST); \
		protoc \
		  --proto_path=$(PROTO_SRC) \
		  --go_out=$(PROTO_DST) --go_opt=paths=source_relative \
		  --go-grpc_out=$(PROTO_DST) --go-grpc_opt=paths=source_relative \
		  $$proto_file; \
		\
		echo "  • Generated: $$proto_file → $(PROTO_DST)/$${rel_path%.proto}.pb.go"; \
	done
	@echo "Done generating all .proto → Go files."

# ─────────────────────────────────────────────────────────────────────────────
# “make proto-clean”: remove all generated code under $(PROTO_DST)
# ─────────────────────────────────────────────────────────────────────────────
.PHONY: proto-clean
proto-clean:
	@echo "Removing generated code under $(PROTO_DST)/…"
	rm -rf $(PROTO_DST)
	@echo "$(PROTO_DST)/ cleaned."

.PHONY: buf-generate
buf-generate:
	@echo "Generating code with buf..."
	buf generate
	@echo "Done generating code with buf."


# Database and SQLC related commands
.PHONY: db-start db-stop db-reset db-psql sqlc-gen

# Start the database container
db-start:
	@echo "Starting PostgreSQL container..."
	docker-compose up -d postgres

# Stop the database container
db-stop:
	@echo "Stopping PostgreSQL container..."
	docker-compose stop postgres

# Reset database (delete volume and restart)
db-reset:
	@echo "Removing PostgreSQL container and volume..."
	docker-compose down -v
	@echo "Starting fresh PostgreSQL container..."
	docker-compose up -d postgres

# Connect to PostgreSQL using psql
db-psql:
	@echo "🔌 Connecting to PostgreSQL..."
	docker-compose exec postgres psql -U postgres -d moss_db

# Database down command - removes containers and optionally volumes
.PHONY: db-down db-down-v

# Stop and remove containers but preserve the data
db-down:
	@echo "Stopping and removing PostgreSQL container..."
	docker-compose down
	@echo "Database container removed. Data volume is preserved."

# Stop and remove containers AND delete volume data
db-down-v:
	@echo "Stopping PostgreSQL container and removing all data..."
	docker-compose down -v
	@echo "Database container and volume removed completely."




