include .env

MIGRATION_DIR = cmd/migrate/migrations

.PHONY: migration-create
migration-create:
	@echo "Creating migration file..."
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir $(MIGRATION_DIR) -seq "$$name"; \
	echo "Migration file created in $(MIGRATION_DIR)"

.PHONY: migrate-up
migrate-up:
	@echo "Applying migrations..."
	@read -p "noumber of migrations to apply (default: all): " count; \
	if [ -z "$$count" ]; then \
		migrate -path $(MIGRATION_DIR) -database $(DB_ADDR) up; \
	else \
		migrate -path $(MIGRATION_DIR) -database $(DB_ADDR) up $$count; \
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

