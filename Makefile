VERSION=0.0.1

SERVER=ccServer
NODE=ccNode
BUILD=${VERSION}
BIN=data/bin
DIR=data/temp/v${VERSION}/${BUILD}
ANDROIDVERSION=21

L=Linux-x64
AA=aarch64
A=armv7a
M=Linux-mips

default: ccserver ccnode-linux ccnode-aarch64 ccnode-armv7a

ccserver:
	export GO111MODULE=on GOOS=linux;export GOARCH=amd64;go build -o ${DIR}/${SERVER} main.go

ccnode-linux:
	export GO111MODULE=on GOOS=linux;export GOARCH=amd64;go build -o ${DIR}/${NODE} cmd/ccnode/main.go

ccnode-aarch64:
	export GO111MODULE=on GOOS=android;export GOARCH=arm64;export CC=~/Android/Sdk/ndk-bundle/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android${ANDROIDVERSION}-clang;export CCX=~/Android/Sdk/ndk-bundle/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android${ANDROIDVERSION}-clang++;export CGO_ENABLED=1; go build -p=8 -o ${DIR}/${NODE}-${AA} cmd/ccnode/main.go

ccnode-armv7a:
	export GO111MODULE=on GOOS=android;export GOARCH=arm;export CC=~/Android/Sdk/ndk-bundle/toolchains/llvm/prebuilt/linux-x86_64/bin/armv7a-linux-androideabi${ANDROIDVERSION}-clang;export CCX=~/Android/Sdk/ndk-bundle/toolchains/llvm/prebuilt/linux-x86_64/bin/armv7a-linux-androideabi${ANDROIDVERSION}-clang++;export CGO_ENABLED=1; export GOARM=7; go build -p=8 -o ${DIR}/${NODE}-${A} cmd/ccnode/main.go