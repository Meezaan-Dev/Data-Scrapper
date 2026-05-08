.PHONY: sandbox scrape stop reset logs

SHELL := /bin/sh

sandbox:
	./scripts/sandbox.sh

install-sandbox-runtime:
	./scripts/install-sandbox-runtime.sh

scrape:
	curl -X POST http://localhost:8080/scrape

stop:
	docker compose down

reset:
	docker compose down -v

logs:
	docker compose logs -f
