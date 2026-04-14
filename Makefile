.PHONY: all build build-frontend build-server run dev dev-frontend dev-backend test clean install help

# --- quick reference ---------------------------------------------------
# make dev          # full dev env — Go backend + Nuxt frontend, Ctrl+C kills both
# make dev-backend  # just the Go API on :8080
# make dev-frontend # just Nuxt dev on :3000 (proxies /api → :8080)
# make build        # production build: static frontend + single Go binary
# make run          # production build + run (serves everything on :8080)
# ---------------------------------------------------------------------

all: install build

install:
	npm install

help:
	@grep '^# make ' Makefile | sed 's/^# //'

# ---- Production build (what the release ships) -----------------------
build: build-frontend build-server

build-frontend:
	npx nuxt generate
	rm -rf frontend/dist
	cp -R .output/public frontend/dist

build-server:
	go build -o haus ./cmd/server

run: build
	./haus

# ---- Dev mode --------------------------------------------------------
# `make dev` runs both servers together. The Go backend goes in the
# background, Nuxt in the foreground — Ctrl+C stops Nuxt and the trap
# kills the Go server too. Visit http://localhost:3000 (Nuxt proxies
# /api and /api/ws to the Go backend on :8080).
dev:
	@trap 'kill 0' INT TERM EXIT; \
		( go run ./cmd/server 2>&1 | sed -u 's/^/[go] /' ) & \
		( npx nuxt dev --port 3000 2>&1 | sed -u 's/^/[nuxt] /' ) & \
		wait

dev-backend:
	go run ./cmd/server

dev-frontend:
	npx nuxt dev --port 3000

# ---- Misc -----------------------------------------------------------
test:
	go test ./...

clean:
	rm -f haus
	rm -rf frontend/dist .nuxt .output
