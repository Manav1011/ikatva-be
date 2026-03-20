DB_URL=postgres://postgres:postgres@localhost:54322/ikatva?sslmode=disable

# --- Migrations ---
migrate-up:
	migrate -path ./migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path ./migrations -database "$(DB_URL)" down 1

migrate-down-all:
	migrate -path ./migrations -database "$(DB_URL)" down

migrate-force:
	migrate -path ./migrations -database "$(DB_URL)" force $(VERSION)

migrate-version:
	migrate -path ./migrations -database "$(DB_URL)" version

# --- Create Migration ---
create-migration:
	migrate create -ext sql -dir ./migrations -seq $(name)
