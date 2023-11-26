make test-app:
	@docker compose -f ./build/docker-compose.test.yml up --build --abort-on-container-exit
	@docker compose -f ./build/docker-compose.test.yml down