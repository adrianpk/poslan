# Vars
STAGE_TAG=stage
PROD_TAG=v0.0.1
IMAGE_NAME=poslan
# Accounts
DOCKERHUB_USER=adrianpksw
# Go
MAKE_CMD=make
# Go
GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_CLEAN=$(GO_CMD) clean
GO_TEST=$(GO_CMD) test
GO_GET=$(GO_CMD) get
## Docker
DOCKER_CMD=docker
DOCKER_LOGIN=$(DOCKER_CMD) login
DOCKER_BUILD=$(DOCKER_CMD) build
DOCKER_PUSH=$(DOCKER_CMD) push
## Kubernetes
KUBECTL_CMD=kubectl
## Helm
HELM_CMD=helm
HELM_DEL=$(HELM_CMD) del
HELM_INSTALL=$(HELM_CMD) install
# Gcloud
GCLOUD_CMD=gcloud
# Misc
BINARY_NAME=main
BINARY_UNIX=$(BINARY_NAME)_unix

all: test build

build:
	$(GO_BUILD) -o ./bin/$(BINARY_NAME) ./main.go

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO_BUILD) -o $(BINARY_UNIX) -v

test:
	$(GO_TEST) -v ./...

clean:
	$(GO_CLEAN)
	rm -f ./bin/$(BINARY_NAME)
	rm -f ./bin/$(BINARY_UNIX)

run:
	make build
	./scripts/run.sh

deps:
	$(GO_GET) -u github.com/BurntSushi/toml
	$(GO_GET) -u github.com/VividCortex/gohistogram
	$(GO_GET) -u github.com/aws/aws-sdk-go
	$(GO_GET) -u github.com/go-kit/kit
	$(GO_GET) -u github.com/go-logfmt/logfmt
	$(GO_GET) -u github.com/go-stack/stack
	$(GO_GET) -u github.com/google/uuid
	$(GO_GET) -u github.com/openzipkin/zipkin-go
	$(GO_GET) -u github.com/sendgrid/rest
	$(GO_GET) -u github.com/sendgrid/sendgrid-go
	$(GO_GET) -u golang.org/x/crypto
	$(GO_GET) -u gopkg.in/yaml.v2
	$(GO_GET) -u honnef.co/go/tools
	$(GO_GET) -u github.com/dgrijalva/jwt-go
	$(GO_GET) -u github.com/go-kit/kit/auth/jwt

build-stage:
	$(MAKE_CMD) build
	$(DOCKER_LOGIN)
	$(DOCKER_BUILD) -t $(DOCKERHUB_USER)/$(IMAGE_NAME):$(STAGE_TAG) .
	$(DOCKER_PUSH) $(DOCKERHUB_USER)/$(IMAGE_NAME):$(STAGE_TAG)


