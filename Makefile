.PHONY: all dev env-clean build env-up env-down run cleanfiles

all: env-clean build env-up run



dev: build run

##### BUILD
build:
	@echo "Build ..."
#	@dep ensure
	@go build
	@echo "Build done"

##### ENV
env-up:
	@echo "Start environment ..."
	@docker-compose -f fixtures/docker-compose-cli.yaml up --force-recreate -d
	@echo "Environment up"
	@docker ps

env-down:
	@echo "Stop environment ..."
	@docker-compose -f fixtures/docker-compose-cli.yaml down
	@echo "Environment down"
	@docker ps

##### RUN
run:
	@echo "Start app ..."
	@./FabricSDKGoExample

##### CLEAN
env-clean: env-down
	@echo "Clean up ..."
	@rm -rf /tmp/exp-* FabricSDKGoExample
	@docker volume prune -f  # 清理挂载卷
	@docker network prune -f # 来清理没有再被任何容器引用的networks
	@docker rm -f -v `docker ps -a --no-trunc | grep "mycc" | cut -d ' ' -f 1` 2>/dev/null || true  #mycc 是链码的名字
	@docker rmi `docker images --no-trunc | grep "mycc" | cut -d ' ' -f 1` 2>/dev/null || true
	@echo "Clean up done"

into-cli:
	@docker exec -it cli bash


cleanfiles:
	@rm -rf ./fixtures/crypto-config
	@rm -rf ./fixtures/channel-artifacts/*

test: env-clean env-up into-cli
