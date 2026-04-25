[parallel]
dev: db-dev backend-dev frontend-dev

db-dev:
    docker compose up -d

backend-dev:
	air

frontend-dev:
	cd wails/frontend && bun run dev


build: frontend-build backend-build

backend-build:
    go build ./cmd/server

frontend-build:
    cd wails/frontend && bun run build

gen:
    go run "internal/gen/typescript.go" && cd wails/frontend && bun run fmt

seed:
    go run ./cmd/server seed

app-build:
    cd wails && wails build

app-dev:
    cd wails && wails dev
