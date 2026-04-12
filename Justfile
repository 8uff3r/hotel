[parallel]
dev: db-dev backend-dev frontend-dev

db-dev:
    docker compose up -d

backend-dev:
	air

frontend-dev:
	cd web && bun run dev


build: frontend-build backend-build

backend-build:
    go build ./cmd/server

frontend-build:
    cd web && bun run build

gen:
    go run "internal/gen/typescript.go" && cd web && bun run fmt

seed:
    go run ./cmd/server seed
