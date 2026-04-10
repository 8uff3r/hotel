[parallel]
dev: backend-dev frontend-dev

backend-dev:
	air

frontend-dev:
	cd web && bun run dev


build: frontend-build backend-build

backend-build:
    go build ./cmd/server

frontend-build:
    cd web && bun run build
