test-app:
	@docker compose -p fetch-system-worker -f ./build/docker-compose.test.yml up --attach server --build --abort-on-container-exit
	@docker compose -f ./build/docker-compose.test.yml down

build-image:
	docker build -t fetch-system-worker -f ./build/Dockerfile .