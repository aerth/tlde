CGO_ENABLED?=0

all: build

build:
	go get -v -d .
	env CGO_ENABLED=${CGO_ENABLED} go build -ldflags='-w -s' -o tlde github.com/aerth/tlde

test:
	go test -v ./...

run: build
	env PORT=8080 ADMIN=tl.de ./tlde

run2: build
	env PORT=8080 ADMIN=tl.de FORMATPATH='/usr/home/%s/%s' ./tlde
