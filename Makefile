VERSION=0.0.1

SERVER=ccServer
NODE=ccNode
BUILD=${VERSION}
BIN=data/bin
DIR=data/temp/v${VERSION}/${BUILD}

L=Linux-x64
A=Android
M=Linux-mips

default: ccserver ccnode-linux

ccserver:
	export GOOS=linux;export GOARCH=amd64;go build -o ${DIR}/${SERVER} main.go

ccnode-linux:
	export GOOS=linux;export GOARCH=amd64;go build -o ${DIR}/${NODE} cmd/ccnode/main.go

ccnode-android:
	export GOOS=android;export GOARCH=arm;export CGO_ENABLED=0;go build -o ${DIR}/${NODE}-${A} cmd/ccnode/main.go