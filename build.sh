#! /bin/sh
echo "start build capture packets Tools ..."

echo "build Dockerfile"

docker build -t tcpdump_go-builder .

echo "run docker images"

docker run --rm -d --name tcpdump_go -it tcpdump_go-builder sh

echo "copy app file"

docker cp tcpdump_go:/app/dist/ ./

echo "clean docker"

docker stop tcpdump_go

echo "build macos arm64"

CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 \
    go build -ldflags="-s -w " -trimpath -o dist/tcpdump_go-darwin-arm64 main.go

echo "build macos amd64"

CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 \
    go build -ldflags="-s -w " -trimpath -o dist/tcpdump_go-darwin-amd64 main.go

echo "all build OK"
