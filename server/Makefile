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

.PHONY integration/test:
integration/test:
	@-docker stop redis > /dev/null 
	@-docker rm redis > /dev/null
	@docker run --detach --name redis -p 6379:6379 redis > /dev/null
	INTEGRATION=1 go test -count=1 -v ./...


