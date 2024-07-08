# 编译 Linux amd64
FROM --platform=linux/amd64 golang:1.16.15 AS builder-linux-amd64
RUN apt-get update && \
    apt-get install -y libpcap-dev
WORKDIR /app
COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w " -trimpath -o dist/tcpdump_go-linux-amd64 main.go

# 编译 Linux 386
FROM --platform=linux/amd64 golang:1.16.15 AS builder-linux-386
RUN dpkg --add-architecture i386 && \
    apt-get update && apt-get install -y \
    libpcap-dev \
    libpcap0.8-dev:i386 \
    gcc-multilib \
    g++-multilib
WORKDIR /app
COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=386 CC="gcc -m32" CXX="g++ -m32" \
    PKG_CONFIG_PATH=/usr/lib/i386-linux-gnu/pkgconfig \
    go build -ldflags="-s -w " -trimpath -o dist/tcpdump_go-linux-386 main.go

# 编译 Linux arm64
FROM --platform=linux/arm64 golang:1.16.15 AS builder-linux-arm64
RUN apt-get update && \
    apt-get install -y libpcap-dev gcc-aarch64-linux-gnu
WORKDIR /app
COPY . .
RUN CGO_ENABLED=1 CC=aarch64-linux-gnu-gcc GOOS=linux GOARCH=arm64 \
    go build -ldflags="-s -w " -trimpath -o dist/tcpdump_go-linux-arm64 main.go

# 编译 Windows amd64
FROM --platform=linux/amd64 golang:1.16.15 AS builder-windows-amd64
RUN apt-get update && \
    apt-get install -y libpcap-dev mingw-w64
WORKDIR /app
COPY . .
RUN CGO_ENABLED=1 GOOS=windows GOARCH=amd64 \
    go build -ldflags="-s -w " -trimpath -o dist/tcpdump_go-windows-amd64.exe .

# 编译 Windows 386
FROM --platform=linux/amd64 golang:1.16.15 AS builder-windows-386
RUN apt-get update && \
    apt-get install -y libpcap-dev mingw-w64
WORKDIR /app
COPY . .
RUN CGO_ENABLED=1 GOOS=windows GOARCH=386 \
    go build -ldflags="-s -w " -trimpath -o dist/tcpdump_go-windows-386.exe main.go

# 收集所有编译结果
FROM alpine:latest AS final
WORKDIR /app
COPY --from=builder-linux-amd64 /app/dist/tcpdump_go-linux-amd64 ./dist/tcpdump_go-linux-amd64
COPY --from=builder-linux-386 /app/dist/tcpdump_go-linux-386 ./dist/tcpdump_go-linux-386
COPY --from=builder-linux-arm64 /app/dist/tcpdump_go-linux-arm64 ./dist/tcpdump_go-linux-arm64
COPY --from=builder-windows-amd64 /app/dist/tcpdump_go-windows-amd64.exe ./dist/tcpdump_go-windows-amd64.exe
COPY --from=builder-windows-386 /app/dist/tcpdump_go-windows-386.exe ./dist/tcpdump_go-windows-386.exe
CMD ["sh"]
