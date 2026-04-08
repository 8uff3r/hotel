[parallel]
run: backend frontend

backend:
	cd backend && go run ./cmd/server

frontend:
	bun run dev
