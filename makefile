CGO_ENABLED?=0

all: build

build:
	env CGO_ENABLED=${CGO_ENABLED} go build -ldflags='-w -s' -o tlde github.com/aerth/tlde

test:
	go test -v ./...

run: build
	env PORT=8080 ADMIN=tl.de ./tlde
