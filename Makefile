lint:
	golangci-lint run ./... -c .golangci.yaml

dev-infrastructure:
	docker compose -f deploy/dev/infrastructure/docker-compose.yaml -p project-layout-dev-infrastructure up -d

dev-build:
	cd deploy/dev && docker compose -p project-layout-dev build

dev-run:
	cd deploy/dev && docker compose -p project-layout-dev up -d