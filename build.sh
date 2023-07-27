#! /bin/sh

CONFIG_MOD="github.com/3elDU/bamboo/config"

export GOARCH=$1
echo "Binary name: 'bamboo-$GOARCH'"

GIT_COMMIT=$(git rev-parse --short HEAD)
GIT_TAG=$(git describe --tags --abbrev=0)
BUILD_MACHINE=$(uname -s -n)
BUILD_DATE=$(date)
go build -ldflags "-X \"$CONFIG_MOD.GitCommit=$GIT_COMMIT\" \
                   -X \"$CONFIG_MOD.GitTag=$GIT_TAG\" \
                   -X \"$CONFIG_MOD.BuildMachine=$BUILD_MACHINE\" \
                   -X \"$CONFIG_MOD.BuildDate=$BUILD_DATE\""  \
                   -o "bamboo-$GOARCH"