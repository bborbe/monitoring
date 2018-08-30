REGISTRY ?= docker.io
IMAGE ?= bborbe/ip
ifeq ($(VERSION),)
	VERSION := $(shell git fetch --tags; git describe --tags `git rev-list --tags --max-count=1`)
endif

all: test install

install:
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install cmd/ip-server/*.go

test:
	go test -cover -race $(shell go list ./... | grep -v /vendor/)

vet:
	go tool vet .
	go tool vet --shadow .

lint:
	golint -min_confidence 1 ./...

errcheck:
	errcheck -ignore '(Close|Write)' ./...

check: lint vet errcheck

format:
	go get golang.org/x/tools/cmd/goimports
	find . -type f -name '*.go' -not -path './vendor/*' -exec gofmt -w "{}" +
	find . -type f -name '*.go' -not -path './vendor/*' -exec goimports -w "{}" +

prepare:
	go get -u golang.org/x/tools/cmd/goimports
	go get -u github.com/golang/lint/golint
	go get -u github.com/kisielk/errcheck
	go get -u github.com/bborbe/docker-utils/cmd/docker-remote-tag-exists
	go get -u github.com/golang/dep/cmd/dep

clean:
	docker rmi $(REGISTRY)/$(IMAGE)-build:$(VERSION)
	docker rmi $(REGISTRY)/$(IMAGE):$(VERSION)

buildgo:
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-s" -a -installsuffix cgo -o ip-server ./go/src/github.com/$(IMAGE)/cmd/ip-server

build:
	docker build --build-arg VERSION=$(VERSION) --no-cache --rm=true -t $(REGISTRY)/$(IMAGE)-build:$(VERSION) -f ./Dockerfile.build .
	docker run -t $(REGISTRY)/$(IMAGE)-build:$(VERSION) /bin/true
	docker cp `docker ps -q -n=1 -f ancestor=$(REGISTRY)/$(IMAGE)-build:$(VERSION) -f status=exited`:/ip-server .
	docker rm `docker ps -q -n=1 -f ancestor=$(REGISTRY)/$(IMAGE)-build:$(VERSION) -f status=exited`
	docker build --no-cache --rm=true --tag=$(REGISTRY)/$(IMAGE):$(VERSION) -f Dockerfile.static .
	rm ip-server

upload:
	docker push $(REGISTRY)/$(IMAGE):$(VERSION)

trigger:
	@go get github.com/bborbe/docker-utils/cmd/docker-remote-tag-exists
	@exists=`docker-remote-tag-exists \
		-registry=${REGISTRY} \
		-repository="${IMAGE}" \
		-credentialsfromfile \
		-tag="${VERSION}" \
		-logtostderr \
		-v=0`; \
	trigger="build"; \
	if [ "$${exists}" = "true" ]; then \
		trigger="skip"; \
	fi; \
	echo $${trigger}

run:
	docker run \
	-p 9090:9090 \
	-e PORT=9090 \
	$(REGISTRY)/bborbe/ip:$(VERSION) \
	-logtostderr \
	-v=0

