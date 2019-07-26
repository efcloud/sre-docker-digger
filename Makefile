NAME       := digger
IMAGE_NAME := 479788333518.dkr.ecr.eu-west-1.amazonaws.com/efcloud/sre/digger
VERSION    :=$(shell git describe --abbrev=0 --tags --exact-match 2>/dev/null || git rev-parse --short HEAD)
LDFLAGS    := -w -extldflags "-static" -X 'main.version=$(VERSION)'

ifndef DRONE_TAG
	VERSION :=$(shell git describe --abbrev=0 --tags --exact-match 2>/dev/null || git rev-parse --short HEAD)
else
	VERSION := $(DRONE_TAG)
endif

.PHONY: in-docker-lint
in-docker-lint:
	golint -set_exit_status .
	go vet ./...

.PHONY: in-docker-test
in-docker-test:
	go test -coverprofile=coverage.out -v ./...

.PHONY: in-docker-build-app
in-docker-build-app:
	CGO_ENABLED=0 GO111MODULE=on go build -o $(NAME) -ldflags "$(LDFLAGS)" .
	strip $(NAME)

.PHONY: setup
setup:
	docker build \
		--tag="$(IMAGE_NAME):$(VERSION)_setup" \
		--file Dockerfile_base .

.PHONY: lint
lint:
	docker build \
		--build-arg SOURCE="$(IMAGE_NAME):$(VERSION)_setup" \
		--tag="$(IMAGE_NAME):$(VERSION)_lint" \
		--target=lint .

.PHONY: test
test:
	docker build \
		--build-arg SOURCE="$(IMAGE_NAME):$(VERSION)_setup" \
		--tag="$(IMAGE_NAME):$(VERSION)_lint" \
		--target=test .

.PHONY: build-app
build-app:
	docker build \
	--build-arg SOURCE="$(IMAGE_NAME):$(VERSION)_setup" \
	--tag="$(IMAGE_NAME):$(VERSION)" \
	--target=builder .

.PHONY: build
build:
	docker build \
	--build-arg SOURCE="$(IMAGE_NAME):$(VERSION)_setup" \
	--tag="$(IMAGE_NAME):$(VERSION)" \
	--target=final .

.PHONY: tag
tag:
	docker tag "$(IMAGE_NAME):$(VERSION)" "$(IMAGE_NAME):$(DRONE_BRANCH)"

.PHONY: tag_release
tag_release:
	docker tag "$(IMAGE_NAME):$(VERSION)" "$(IMAGE_NAME):$(DRONE_TAG)"

.PHONY: publish
publish:
	docker push "$(IMAGE_NAME):$(VERSION)"
	docker push "$(IMAGE_NAME):$(DRONE_BRANCH)"

.PHONY: publish_release
publish_release:
	docker push "$(IMAGE_NAME):$(DRONE_TAG)"

build_all: setup lint test build-app build

.PHONY: run_dev_vault
run_vault:
	docker run -d --cap-add=IPC_LOCK -e 'VAULT_DEV_ROOT_TOKEN_ID=myroot' -p 8200:8200 -e 'VAULT_DEV_LISTEN_ADDRESS=0.0.0.0:8200' vault
