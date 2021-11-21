all: push

ifeq ($(VERSION),)
    VERSION = 0.1
endif

TAG = $(VERSION)
REPO = camelinx
IMAGE = l7bench
BINNAME = build/l7bench
CURDIR = $(shell pwd)

PREFIX = $(REPO)/$(IMAGE)

DOCKER_RUN = docker run --rm -v $(CURDIR)/../:/go/src/github.com -w /go/src/github.com/l7bench/
GOLANG_CONTAINER = golang:1.16
DOCKERFILE = build/Dockerfile

l7bench:
	$(DOCKER_RUN) -e CGO_ENABLED=0 $(GOLANG_CONTAINER) go build -ldflags "-w -X main.version=${VERSION}" -o $(BINNAME) github.com/l7bench/cmd/l7bench

test:
	$(DOCKER_RUN) $(GOLANG_CONTAINER) go test ./...

container: l7bench
	docker build -f $(DOCKERFILE) -t $(PREFIX):$(TAG) .

push: container
	docker push $(PREFIX):$(TAG)

clean:
	rm -f $(BINNAME)
	docker rmi $(PREFIX):$(TAG)
