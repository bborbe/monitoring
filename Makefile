install:
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install bin/monitoring_check/monitoring_check.go
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install bin/monitoring_cron/monitoring_cron.go
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install bin/monitoring_server/monitoring_server.go
test:
	GO15VENDOREXPERIMENT=1 go test -cover `glide novendor`
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
format:
	find . -name "*.go" -exec gofmt -w "{}" \;
	goimports -w=true .
prepare:
	go get -u golang.org/x/tools/cmd/goimports
	go get -u github.com/Masterminds/glide
	go get -u github.com/golang/lint/golint
	go get -u github.com/kisielk/errcheck
	glide install
update:
	glide up
clean:
	rm -rf vendor phantomJsdriver.log phantomJsOutput.log
