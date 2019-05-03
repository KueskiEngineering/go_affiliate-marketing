build: build_osx build_linux

build_osx: test
	GOOS=darwin go build

build_linux: test
	GOOS=linux go build -o go-affiliate-marketing-linux

coverage: test
	go tool cover -html=coverage.out

test: dep
	go test -v -cover -coverprofile=coverage.out ./...

dep:
	dep ensure
