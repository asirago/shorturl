# ==================================================
#					     BUILD
# ==================================================

run: 
	@go run ./cmd/api/...

docker/rebuild:
	@echo "rebuilding..."
	@-docker stop shorturl > /dev/null 2>&1 || true
	@-docker rm shorturl > /dev/null 2>&1 || true
	@docker build --tag shorturl . > /dev/null 2>&1 || true
 	
docker/run:
	@docker run --name shorturl shorturl

docker/run-and-rebuild:
	@echo "rebuilding and running"
	@make docker/rebuild
	@make docker/run
