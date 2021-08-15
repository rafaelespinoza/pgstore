BASENAME=pgstore_test
DOCKER_DIR=.docker

.PHONY: test-up test-down
test-up: build
	docker-compose -f $(DOCKER_DIR)/compose.yml -- up --exit-code-from client
test-down:
	docker-compose -f $(DOCKER_DIR)/compose.yml -- down --rmi all --volumes

.PHONY: build
build: $(eval BUILD_DIR=$(shell mktemp -d -p /tmp $(BASENAME)_XXXXXX))
build:
	git clone --depth=1 file://$(PWD) $(BUILD_DIR)
	docker image build -f $(DOCKER_DIR)/client.Dockerfile -t $(BASENAME)/client $(BUILD_DIR)
	docker image build -f $(DOCKER_DIR)/server.Dockerfile -t $(BASENAME)/server $(BUILD_DIR)
	rm -rf $(BUILD_DIR)

.PHONY: clean kill rm rmi
clean: kill rm rmi
kill:
	docker kill $(BASENAME)/client $(BASENAME)/server
rm:
	docker container rm $(shell docker container ls -aq --filter ancestor=$(BASENAME)/client)
	docker container rm $(shell docker container ls -aq --filter ancestor=$(BASENAME)/server)
rmi:
	docker image rmi $(shell docker image ls -aq $(BASENAME)/client)
	docker image rmi $(shell docker image ls -aq $(BASENAME)/server)
