.PHONY: all build build-frontend build-server run dev dev-server test clean install

all: install build

install:
	npm install

build: build-frontend build-server

build-frontend:
	npx nuxt generate
	rm -rf frontend/dist
	cp -R .output/public frontend/dist

build-server:
	go build -o haus ./cmd/server

run: build
	./haus

dev:
	npx nuxt dev

dev-server:
	go run ./cmd/server

test:
	go test ./...

clean:
	rm -f haus
	rm -rf frontend/dist .nuxt .output
