.PHONY: help build build-local up down logs ps test
.DEFAULT_GOAL := help

DOCKER_TAG := latest
build: ## 배포용 도커 이미지 빌드
	docker build -t bsj1209/gotodo:${DOCKER_TAG} \
	--target deploy ./

build-local: ## 로컬 환경용 도커 이미지 빌드
	docker compose build --no-cache


up: ## 자동 새로고침을 사용한 도커 이미지 빌드
	docker compose up -d

down: ## 도커 구성기 종료
	docker compose down

logs: ## 도커 구성기 로그 출력
	docker compose logs -f

ps: ## 컨테이너 상태 확인
	docker compose ps

test: ## 테스트 실행
	go test -race -shuffle=on ./...

dry-migrate: ## Try migration
	mysqldef -u todo -p todo -h 127.0.0.1 -P 3306 todo --dry-run < ./_tools/mysql/schema.sql

migrate:  ## Execute migration
	mysqldef -u todo -p todo -h 127.0.0.1 -P 3306 todo < ./_tools/mysql/schema.sql

generate: ## Generate codes
	go generate ./...

help: ## Show options
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

