# ==================================================
#					     BUILD
# ==================================================

.PHONY run: 
run:
	@go run ./cmd/api/...

.PHONY docker/rebuild:
docker/rebuild:
	@echo "purging..."
	-make docker/purge
	@echo "rebuilding..."
	docker compose up --build
 	
.PHONY docker/run:
docker/run:
	@docker compose up

.PHONY docker/purge:
docker/purge:
	@docker stop shorturl-api redis
	@docker rm shorturl-api redis
	@docker rmi shorturl-api redis

.PHONY purge/all:
purge/all:
	@-docker stop $$(docker ps -aq)
	@-docker rm $$(docker ps -aq)
	@-docker rmi $$(docker images)
