# Declaration of variables here - which can be used later in make rules (each rule is basically a command that can be invoked by make)
DB_CONN = "host=localhost port=5432 user=baloo password=junglebook dbname=lenslocked sslmode=disable" 
MIGRATION_DIR = "./migrations"


# runs the status 
goose_status:
	goose postgres -dir $(MIGRATION_DIR) $(DB_CONN) status

# runs the migration up
goose_up:
	goose postgres -dir $(MIGRATION_DIR) $(DB_CONN) up

# runs the migration down
goose_down:
	goose postgres -dir $(MIGRATION_DIR) $(DB_CONN) down

# resets all migrations to 0
goose_reset:
	goose postgres -dir $(MIGRATION_DIR) $(DB_CONN) reset