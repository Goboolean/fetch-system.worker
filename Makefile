test-app:
	@docker compose -p fetch-system-worker -f ./deploy/docker-compose.test.yml up --attach server --build --abort-on-container-exit
	@docker compose -f ./deploy/docker-compose.test.yml down

build-app:
	docker build -t registry.mulmuri.dev/fetch-system-worker:latest -f ./deploy/Dockerfile .
	docker push registry.mulmuri.dev/fetch-system-worker:latest
	helm upgrade fetch-system ~/lab -n goboolean

generate-wire:
	wire cmd/wire/wire_setup.go