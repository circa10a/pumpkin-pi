PROJECT=circa10a/pumpkin-pi
BINARY=pumpkin-pi
VERSION=0.1.0
GOBUILDFLAGS=-ldflags="-s -w"
DOCKERBUILDX=docker buildx build
DOCKERBUILDARM64=$(DOCKERBUILDX) --platform linux/arm64 -t $(PROJECT) .
DOCKERBUILDARMv7=$(DOCKERBUILDX) --platform linux/arm/v7 -t $(PROJECT) .

.PHONY: build

build:
	go build $(GOBUILDFLAGS) -o $(BINARY)

build-docker:
	docker build -t $(PROJECT) .

build-docker-arm64:
	$(DOCKERBUILDARM64)

build-docker-armv7:
	$(DOCKERBUILDARMv7)

push-docker:
	docker push $(PROJECT)

run:
	go run .

lint:
	 golangci-lint run -v

compile:
	GOOS=linux GOARCH=amd64 $(GOGENERATE) && go build $(GOBUILDFLAGS) -o bin/$(BINARY)-linux-amd64
	GOOS=linux GOARCH=arm $(GOGENERATE) && go build $(GOBUILDFLAGS) -o bin/$(BINARY)-linux-arm
	GOOS=linux GOARCH=arm64$(GOGENERATE) && go build $(GOBUILDFLAGS) -o bin/$(BINARY)-linux-arm64
	GOOS=darwin GOARCH=amd64 $(GOGENERATE) && go build $(GOBUILDFLAGS) -o bin/$(BINARY)-darwin-amd64
