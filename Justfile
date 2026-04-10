[parallel]
dev: backend frontend

backend:
	go run ./cmd/server

frontend:
	cd web && bun run dev
