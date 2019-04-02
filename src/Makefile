VERSION=0.0.1

SERVER=ccServer
NODE=ccNode
BUILD=${VERSION}
BIN=data/bin
DIR=data/temp/v${VERSION}/${BUILD}
ANDROIDVERSION=26

L=Linux-x64
A=Android
M=Linux-mips

default: ccserver ccnode-linux

ccserver:
	export GOOS=linux;export GOARCH=amd64;go build -o ${DIR}/${SERVER} main.go

ccnode-linux:
	export GOOS=linux;export GOARCH=amd64;go build -o ${DIR}/${NODE} cmd/ccnode/main.go

ccnode-android:
	export GOOS=android;export GOARCH=arm64;export CC=~/Android/Sdk/ndk-bundle/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android${ANDROIDVERSION}-clang;export CCX=~/Android/Sdk/ndk-bundle/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android${ANDROIDVERSION}-clang++;export CGO_ENABLED=1; export GOARM=7;go build -p=8 -o ${DIR}/${NODE}-${A} cmd/ccnode/main.go