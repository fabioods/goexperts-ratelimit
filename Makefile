build:
	@go build -ldflags="-w -s" -o ./bin/api ./cmd/api/main.go

run:
	docker-compose up -d

test:
	@./scripts/test.sh

test_k6_smoke:
	@echo "==> Running smoke test with k6 via Docker"
	docker-compose run k6 run /scripts/smoke/smoke.test.js

install-deps:
	go mod tidy

setup: install-deps

clean:
	rm -rf ./bin ./tmp ./coverage.txt
