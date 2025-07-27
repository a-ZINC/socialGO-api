include .env

MIGRATION_DIR = cmd/migrate/migrations

.PHONY: migration-create
migration-create:
	@echo "Creating migration file..."
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir $(MIGRATION_DIR) "$$name"; \
	echo "Migration file created in $(MIGRATION_DIR)"

.PHONY: migrate-up
migrate-up:
	@echo "Applying migrations..."
	@read -p "noumber of migrations to apply (default: all): " count; \
	if [ -z "$$count" ]; then \
		migrate -path $(MIGRATION_DIR) -database $(DB_ADDR) up; \
	else \
		migrate -path $(MIGRATION_DIR) -database $(DB_ADDR) up "$$count"; \
	fi
.PHONY: migrate-down
migrate-down:
	@echo "Rolling back migrations..."
	@read -p "Number of migrations to roll back (default: 1): " count; \
	if [ -z "$$count" ]; then \
		migrate -path $(MIGRATION_DIR) -database $(DB_ADDR) down 1; \
	elif [ "$$count" = "all" ]; then \
		migrate -path $(MIGRATION_DIR) -database $(DB_ADDR) down; \
	else \
		migrate -path $(MIGRATION_DIR) -database $(DB_ADDR) down "$$count"; \
	fi
.PHONY: migrate-force
migrate-force:
	@echo "Forcing migration to version..."
	@read -p "Enter version to force: " version; \
	migrate -path $(MIGRATION_DIR) -database $(DB_ADDR) force "$$version"; \
	echo "Migration forced to version $$version"
.PHONY: seed
seed:
	@echo "Seeding database..."
	@go run ./cmd/migrate/seed/main.go
.PHONY: gen-docs
gen-docs:
	@echo "Generating API documentation..."
	@swag init -g ./cmd/api/main.go -d ./ --parseDependency --parseInternal
	@echo "API documentation generated at ./cmd/api/swagger/doc.json"

