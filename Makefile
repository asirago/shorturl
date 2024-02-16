# ==================================================
#					     BUILD
# ==================================================

run: 
	@go run ./cmd/api/...

docker/rebuild:
	@echo "rebuilding..."
	@docker compose up --build
 	
docker/run:
	@docker compose up

docker/purge:
	@docker stop $$(docker ps -aq)
	@docker rm $$(docker ps -aq)
	@docker rmi $$(docker images -q)

