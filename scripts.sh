#!/bin/bash

set -e

root=$(cd "$(dirname "$0")"; pwd)

cp pb/fs.proto ui/public

is_command_exist () {
  which $1 >/dev/null 2>&1
}

if ! is_command_exist protoc-gen-go; then
  GOOS= GOARCH= go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
  GOOS= GOARCH= go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
fi

export GO111MODULE=on
export PATH="$PATH:$(go env GOPATH)/bin"

protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. pb/fs.proto

usage () {
  echo 'Usage:
  bash scripts.sh start [web|desktop]
  bash scripts.sh build [server|cli|desktop]'
}

case $1 in
  start)
    case $2 in
      web)
        cd $root/ui
        yarn start
        ;;

      electron)
        trap "kill 0" EXIT
        cd $root/ui
        tempfile=$(mktemp)
        yarn watch > $tempfile 2>&1 &
        cnt=1
        while IFS= read -r line; do
          echo $line
          if [[ $line =~ "Start webpack watch" ]]; then
            cnt=`expr $cnt + 1`
          fi
          if [[ $cnt == 2 ]]; then
            break;
          fi
        done < <(tail -f $tempfile)
        yarn start:electron
        ;;

      *)
        usage
        ;;
    esac
    ;;

  build)
    case $2 in
      server)
        cd $root/ui
        yarn
        yarn build
        cd $root/cmd/kfs-server
        echo "GOOS=$GOOS GOARCH=$GOARCH"
        if [[ $GOOS != '' && $GOARCH != '' ]]; then
          GOOS=$GOOS GOARCH=$GOARCH go build -o kfs-server-$GOOS-$GOARCH
        else
          go build -o kfs-cli
        fi
        ;;

      cli)
        cd $root/cmd/kfs-cli
        echo "GOOS=$GOOS GOARCH=$GOARCH"
        if [[ $GOOS != '' && $GOARCH != '' ]]; then
          GOOS=$GOOS GOARCH=$GOARCH go build -o kfs-cli-$GOOS-$GOARCH
        else
          go build -o kfs-cli
        fi
        ;;

      electron)
        cd $root/ui
        yarn
        yarn build
        yarn build:electron
        ;;

      *)
        usage
        ;;
    esac
    ;;

  *)
    usage
    ;;
esac
