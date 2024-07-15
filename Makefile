init: init-ci
init-ci:
	docker-pull docker-build docker-up
up: docker-up
down: docker-down
restart: down up

test:
	go test -v ./...

lint:
	golangci-lint run -v

docker-up:
	docker compose up -d
docker-down:
	docker compose down --remove-orphans

docker-down-clear:
	docker compose down -v --remove-orphans

docker-pull:
	docker compose pull

docker-build:
	docker compose build --pull

dev-docker-build:
	REGISTRY=localhost IMAGE_TAG=main-1 make docker-build

docker-build: docker-build-service docker-build-parser

docker-build-service:
	DOCKER_BUILDKIT=1 docker --log-level=debug build --pull --build-arg BUILDKIT_INLINE_CACHE=1 \
    	--tag ${REGISTRY}/svodd-server-service:${IMAGE_TAG} \
    	--file ./docker/service/Dockerfile .

docker-build-parser:
	DOCKER_BUILDKIT=1 docker --log-level=debug build --pull --build-arg BUILDKIT_INLINE_CACHE=1 \
    	--tag ${REGISTRY}/svodd-server-parser:${IMAGE_TAG} \
    	--file ./docker/parser/Dockerfile .

push:
	docker push ${REGISTRY}/svodd-server-service:${IMAGE_TAG}
	docker push ${REGISTRY}/svodd-server-parser:${IMAGE_TAG}

deploy:
	ssh -o StrictHostKeyChecking=no deploy@${HOST} -p ${PORT} 'docker network create --driver=overlay traefik-public || true'
	ssh -o StrictHostKeyChecking=no deploy@${HOST} -p ${PORT} 'docker network create --driver=overlay svodd-network || true'
	ssh -o StrictHostKeyChecking=no deploy@${HOST} -p ${PORT} 'rm -rf svodd-server-service_${BUILD_NUMBER} && mkdir svodd-server-service_${BUILD_NUMBER}'

	envsubst < docker-compose-production.yml > docker-compose-production-env.yml
	scp -o StrictHostKeyChecking=no -P ${PORT} docker-compose-production-env.yml deploy@${HOST}:svodd-server-service_${BUILD_NUMBER}/docker-compose.yml
	rm -f docker-compose-production-env.yml

	ssh -o StrictHostKeyChecking=no deploy@${HOST} -p ${PORT} 'cd svodd-server-service_${BUILD_NUMBER} && docker stack deploy --compose-file docker-compose.yml svodd-server-service --with-registry-auth --prune'
