# ─────────────────────────────────────────────────────────────────────────────
# Makefile (place this at the project root, alongside go.mod and main.go)
# ─────────────────────────────────────────────────────────────────────────────

# 1) All your .proto files live under:
PROTO_SRC := protobuf

# 2) Base directory where generated Go code will be written:
PROTO_DST := internal/genproto

# 3) Find every .proto file under $(PROTO_SRC):
PROTO_FILES := $(shell find $(PROTO_SRC) -name "*.proto")

# 4) Ensure protoc-gen-go & protoc-gen-go-grpc are on your PATH:
PATH := $(PATH):$(shell go env GOPATH)/bin

# ─────────────────────────────────────────────────────────────────────────────
# “make proto”: regenerate all .proto → .pb.go + _grpc.pb.go under $(PROTO_DST)
# ─────────────────────────────────────────────────────────────────────────────
.PHONY: proto
proto:
	@echo "🔄 Regenerating Go stubs from all .proto files…"
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
	@echo "✅ Done generating all .proto → Go files."

# ─────────────────────────────────────────────────────────────────────────────
# “make proto-clean”: remove all generated code under $(PROTO_DST)
# ─────────────────────────────────────────────────────────────────────────────
.PHONY: proto-clean
proto-clean:
	@echo "🧹 Removing generated code under $(PROTO_DST)/…"
	rm -rf $(PROTO_DST)
	@echo "✅ $(PROTO_DST)/ cleaned."



# Start postgres container
up:
	docker-compose up -d postgres
	@echo "Waiting for postgres to be ready..."
	@sleep 3

# Stop postgres container
down:
	docker-compose down

# Connect to postgres - using docker exec instead of docker-compose exec
psql:
	@if [ "$$(docker ps -q -f name=lumo-postgres-1)" ]; then \
		docker exec -it lumo-postgres-1 psql -U postgres -d lumo_db; \
	else \
		echo "Postgres is not running. Starting it now..."; \
		$(MAKE) up; \
		docker exec -it lumo_db psql -U postgres -d lumo_db; \
	fi




