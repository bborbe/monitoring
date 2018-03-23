REGISTRY ?= docker.io
ifeq ($(VERSION),)
	VERSION := $(shell git fetch --tags; git describe --tags `git rev-list --tags --max-count=1`)
endif

all: test install run

install:
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install bin/monitoring_check/*.go
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install bin/monitoring_cron/*.go
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install bin/monitoring_server/*.go

glide:
	go get github.com/Masterminds/glide

test: glide
	GO15VENDOREXPERIMENT=1 go test -cover `glide novendor`

unittest: glide
	GO15VENDOREXPERIMENT=1 go test -short -cover `glide novendor`

vet:
	go tool vet .
	go tool vet --shadow .

lint:
	golint -min_confidence 1 ./...

errcheck:
	errcheck -ignore '(Close|Write)' ./...

check: lint vet errcheck

run:
	monitoring_check \
	-logtostderr \
	-v=2 \
	-config=sample_config.xml

goimports:
	go get golang.org/x/tools/cmd/goimports

format: goimports
	find . -type f -name '*.go' -not -path './vendor/*' -exec gofmt -w "{}" +
	find . -type f -name '*.go' -not -path './vendor/*' -exec goimports -w "{}" +

prepare:
	go get -u golang.org/x/tools/cmd/goimports
	go get -u github.com/Masterminds/glide
	go get -u github.com/golang/lint/golint
	go get -u github.com/kisielk/errcheck

clean:
	docker rmi $(REGISTRY)/bborbe/monitoring-build:$(VERSION)
	docker rmi $(REGISTRY)/bborbe/monitoring:$(VERSION)

buildgo:
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-s" -a -installsuffix cgo -o monitoring_server ./go/src/github.com/bborbe/monitoring/bin/monitoring_server

build:
	docker build --build-arg VERSION=$(VERSION) --no-cache --rm=true -t $(REGISTRY)/bborbe/monitoring-build:$(VERSION) -f ./Dockerfile.build .
	docker run -t $(REGISTRY)/bborbe/monitoring-build:$(VERSION) /bin/true
	docker cp `docker ps -q -n=1 -f ancestor=$(REGISTRY)/bborbe/monitoring-build:$(VERSION) -f status=exited`:/monitoring_server .
	docker rm `docker ps -q -n=1 -f ancestor=$(REGISTRY)/bborbe/monitoring-build:$(VERSION) -f status=exited`
	docker build --no-cache --rm=true --tag=$(REGISTRY)/bborbe/monitoring:$(VERSION) -f Dockerfile.static .
	rm monitoring_server

upload:
	docker push $(REGISTRY)/bborbe/monitoring:$(VERSION)

rundocker:
	docker run \
	--publish 8080:8080 \
	--env PORT=8080 \
	--env CONFIG=/data/config.xml \
	--volume `pwd`/example:/data \
	$(REGISTRY)/bborbe/monitoring:$(VERSION) \
	-logtostderr \
	-v=0


