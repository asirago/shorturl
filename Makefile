# ==================================================
#					     BUILD
# ==================================================

.PHONY run: 
run:
	@go run ./cmd/api/...

.PHONY docker/rebuild:
docker/rebuild:
	@echo "rebuilding..."
	@docker compose up --build
 	
.PHONY docker/run:
docker/run:
	@docker compose up

.PHONY docker/purge:
docker/purge:
	@docker stop $$(docker ps -aq)
	@docker rm $$(docker ps -aq)
	@docker rmi $$(docker images -q)

