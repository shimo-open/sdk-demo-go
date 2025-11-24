APP_NAME:=sdk-demo-go
APP_PATH:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
SCRIPT_PATH:=$(APP_PATH)/scripts
COMPILE_OUT:=$(APP_PATH)/bin/$(APP_NAME)

server:export EGO_DEBUG=true
server:export EGO_MODE=dev
server:
	@cd $(APP_PATH) && go run main.go server --config=config/local.toml

test:export EGO_DEBUG=true
test:
	@cd $(APP_PATH) && go test -covermode count `go list ./... | grep -v tests`

build:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making build app<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@chmod +x $(SCRIPT_PATH)/build/*.sh
	@cd $(APP_PATH) && $(SCRIPT_PATH)/build/gobuild.sh $(APP_NAME) $(COMPILE_OUT)
	@echo -e "\n"

api-test-all:
	go run main.go sdk-ctl api-test batch-test all all