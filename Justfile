[parallel]
dev: db-dev backend-dev frontend-dev

[parallel]
dev-app: db-dev backend-dev app-dev

db-dev:
    docker compose up -d

backend-dev:
	air

frontend-dev:
	cd frontend && bun run dev


build: frontend-build backend-build

backend-build:
    go build ./cmd/server

frontend-build:
    cd frontend && bun run build

gen-perms:
    go run "cmd/codegen/permissions/main.go" && cd frontend && bun run fmt

seed:
    go run ./cmd/server seed

app-build:
    cd wails && wails build

app-dev:
    cd wails && wails dev
