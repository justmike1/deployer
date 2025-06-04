release:
	DOCKER_BUILDKIT=1 docker buildx build --platform=linux/amd64 \
		--push \
		--progress=plain \
		-t mikeengineering/deployer:latest \
		-f Dockerfile .