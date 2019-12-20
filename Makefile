NAME       := digger
IMAGE_NAME := 479788333518.dkr.ecr.eu-west-1.amazonaws.com/efcloud/sre/digger
VERSION    :=$(shell git describe --abbrev=0 --tags --exact-match 2>/dev/null || git rev-parse --short HEAD)
LDFLAGS    := -w -extldflags "-static" -X 'main.version=$(VERSION)'

ifdef DRONE_BRANCH
	IMAGE_VERSION = $(DRONE_BRANCH)_$(DRONE_BUILD_NUMBER)
else
	IMAGE_VERSION = $(shell git describe --abbrev=0 --tags --exact-match 2>/dev/null || git rev-parse --short HEAD)
endif

ifdef GIT_TAG
	GIT_TAG = $(shell git describe --abbrev=0 --tags --exact-match 2>/dev/null)
endif

.PHONY: in-docker-lint
in-docker-lint:
	golint -set_exit_status .
	go vet ./...

.PHONY: in-docker-test
in-docker-test:
	go test -coverprofile=/tmp/coverage.out -v ./...

.PHONY: in-docker-build-app
in-docker-build-app:
	CGO_ENABLED=0 GO111MODULE=on go build -o $(NAME) -ldflags "$(LDFLAGS)" .
	strip $(NAME)

.PHONY: setup
setup:
	docker build \
		--tag="$(IMAGE_NAME):$(IMAGE_VERSION)_setup" \
		--file Dockerfile_base .

.PHONY: lint
lint:
	docker run --rm \
		"$(IMAGE_NAME):$(IMAGE_VERSION)_setup" \
		make in-docker-lint

.PHONY: test
test:
	docker run \
		--name "test-$(DRONE_BUILD_NUMBER)" \
		"$(IMAGE_NAME):$(IMAGE_VERSION)_setup" \
		make in-docker-test

	docker cp "test-$(DRONE_BUILD_NUMBER)":/tmp/coverage.out .
	docker rm "test-$(DRONE_BUILD_NUMBER)"

.PHONY: build-app
build-app:
	docker build \
	--build-arg SOURCE="$(IMAGE_NAME):$(IMAGE_VERSION)_setup" \
	--tag="$(IMAGE_NAME):$(IMAGE_VERSION)" \
	--target=builder .

.PHONY: build
build:
	docker build \
	--build-arg SOURCE="$(IMAGE_NAME):$(IMAGE_VERSION)_setup" \
	--tag="$(IMAGE_NAME):$(IMAGE_VERSION)" \
	--target=final .

.PHONY: tag
tag:
	docker tag "$(IMAGE_NAME):$(IMAGE_VERSION)" "$(IMAGE_NAME):$(GIT_TAG)"

.PHONY: tag_release
tag_release:
	docker tag "$(IMAGE_NAME):$(IMAGE_VERSION)" "$(IMAGE_NAME):$(DRONE_TAG)"

.PHONY: publish
publish:
	docker push "$(IMAGE_NAME):$(IMAGE_VERSION)"
	docker push "$(IMAGE_NAME):$(GIT_TAG)"

.PHONY: publish_release
publish_release:
	docker push "$(IMAGE_NAME):$(DRONE_TAG)"

build_all: setup lint test build-app build
