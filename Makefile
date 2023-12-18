test-app:
	@docker compose -p fetch-system-worker -f ./deploy/docker-compose.test.yml up --attach server --build --abort-on-container-exit
	@docker compose -f ./deploy/docker-compose.test.yml down

build-image:
	docker build -t fetch-system-worker -f ./deploy/Dockerfile .

generate-wire:
	wire cmd/wire/wire_setup.go