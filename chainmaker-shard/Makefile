ifeq ($(OS),Windows_NT)
  PLATFORM="Windows"
else
  ifeq ($(shell uname),Darwin)
    PLATFORM="MacOS"
  else
    PLATFORM="Linux"
  endif
endif
DATETIME=$(shell date "+%Y%m%d%H%M%S")
VERSION=V1.0.0

PROJECT_PATH=${shell pwd}
COVERAGE_PATH=${shell pwd}/docs/coverage
CI_PATH=${shell pwd}/docs/ci

cross:
	@echo "build cross-chain binary execute file..."
	@cd main && go build -o ./cross-chain && cd ..
	@mv main/cross-chain ./release/lib/
	@echo "make cross finished"

clean:
	@rm -rf ./release

ci:
	@rm -rf ${CI_PATH}
	@mkdir ${CI_PATH}
	@echo module ci path: ${CI_PATH}
	@cd module/adapter  && golangci-lint run ./... > ${CI_PATH}/adapter.ci.out && cd -
	@cd module/channel  && golangci-lint run ./... > ${CI_PATH}/channel.ci.out && cd -
	@cd module/conf     && golangci-lint run ./... > ${CI_PATH}/conf.ci.out && cd -
	@cd module/event    && golangci-lint run ./... > ${CI_PATH}/event.ci.out && cd -
	@cd module/listener && golangci-lint run ./... > ${CI_PATH}/listener.ci.out && cd -
	@cd module/logger   && golangci-lint run ./... > ${CI_PATH}/logger.ci.out && cd -
	@cd module/net      && golangci-lint run ./... > ${CI_PATH}/net.ci.out && cd -
	@cd module/prover   && golangci-lint run ./... > ${CI_PATH}/prover.ci.out && cd -
	@cd module/router   && golangci-lint run ./... > ${CI_PATH}/router.ci.out && cd -
	@cd module/server   && golangci-lint run ./... > ${CI_PATH}/server.ci.out && cd -
	@cd module/store    && golangci-lint run ./... > ${CI_PATH}/store.ci.out && cd -
	@cd module/utils    && golangci-lint run ./... > ${CI_PATH}/utils.ci.out && cd -
	@cd module/transaction && golangci-lint run ./... > ${CI_PATH}/transaction.ci.out && cd -


cover:
	@rm -rf ${COVERAGE_PATH}
	@mkdir ${COVERAGE_PATH}
	@echo module coverage path: ${COVERAGE_PATH}
	@cd module/adapter  && go test ./... -coverprofile=adapter.out   && go tool cover -func=adapter.out > ${COVERAGE_PATH}/adapter.out   && rm adapter.out  && find . -name 'logs'| xargs rm -r && cd -
	@cd module/channel  && go test ./... -coverprofile=channel.out   && go tool cover -func=channel.out > ${COVERAGE_PATH}/channel.out   && rm channel.out  && find . -name 'logs'| xargs rm -r && cd -
	@cd module/conf     && go test ./... -coverprofile=conf.out      && go tool cover -func=conf.out > ${COVERAGE_PATH}/conf.out         && rm conf.out     && find . -name 'logs'| xargs rm -r && cd -
	@cd module/event    && go test ./... -coverprofile=event.out     && go tool cover -func=event.out > ${COVERAGE_PATH}/event.out       && rm event.out    && find . -name 'logs'| xargs rm -r && cd -
	@cd module/listener && go test ./... -coverprofile=listener.out  && go tool cover -func=listener.out > ${COVERAGE_PATH}/listener.out && rm listener.out && find . -name 'logs'| xargs rm -r && cd -
	@cd module/logger   && go test ./... -coverprofile=logger.out    && go tool cover -func=logger.out > ${COVERAGE_PATH}/logger.out     && rm logger.out   && find . -name 'logs'| xargs rm -r && cd -
	@cd module/net      && go test ./... -coverprofile=net.out       && go tool cover -func=net.out > ${COVERAGE_PATH}/net.out           && rm net.out      && find . -name 'logs'| xargs rm -r && cd -
	@cd module/prover   && go test ./... -coverprofile=prover.out    && go tool cover -func=prover.out > ${COVERAGE_PATH}/prover.out     && rm prover.out   && find . -name 'logs'| xargs rm -r && cd -
	@cd module/router   && go test ./... -coverprofile=router.out    && go tool cover -func=router.out > ${COVERAGE_PATH}/router.out     && rm router.out   && find . -name 'logs'| xargs rm -r && cd -
	@cd module/server   && go test ./... -coverprofile=server.out    && go tool cover -func=server.out > ${COVERAGE_PATH}/server.out     && rm server.out   && find . -name 'logs'| xargs rm -r && cd -
	@cd module/store    && go test ./... -coverprofile=store.out     && go tool cover -func=store.out > ${COVERAGE_PATH}/store.out       && rm store.out    && find . -name 'logs' -or -name 'testdata'| xargs rm -r && cd -
	@cd module/utils    && go test ./... -coverprofile=utils.out     && go tool cover -func=utils.out > ${COVERAGE_PATH}/utils.out       && rm utils.out    && find . -name 'logs'| xargs rm -r && cd -
	@cd module/transaction && go test ./... -coverprofile=transaction.out && go tool cover -func=transaction.out > ${COVERAGE_PATH}/transaction.out && rm transaction.out && find . -name 'logs'| xargs rm -r && cd -
	@cd tools/sdk && go test ./... -coverprofile=sdk.out && go tool cover -func=sdk.out > ${COVERAGE_PATH}/sdk.out && rm sdk.out && find . -name 'logs'| xargs rm -r && cd -

