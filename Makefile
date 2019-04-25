GOOS ?= darwin

.PHONY: test cover

test:
	go test -v ./... -cover -coverprofile=cov.out

cover: 
	go tool cover -html=cov.out

.PHONY: clean build

clean:
	go clean
	rm -rf target
	rm -rf cov.out

build: clean
	@GOOS=$(GOOS) go build -o target/gopen
	chmod +x target/gopen
