export CGO_ENABLED := 0
export GOPROXY := https://goproxy.io,direct

dep:
	go get -u -v ./... && \
	go mod tidy && \
	go fmt ./...

test:
	go test -v -count=1 ./...
