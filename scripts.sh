#!/bin/bash

trap "kill 0" EXIT

set -e

root=$(cd "$(dirname "$0")"; pwd)

cp pb/fs.proto ui/public

is_command_exist () {
  which $1 >/dev/null 2>&1
}

if ! is_command_exist protoc-gen-go; then
  echo "install protoc-gen-go and protoc-gen-go-grpc"
  export GO111MODULE=on
  go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
  export PATH="$PATH:$(go env GOPATH)/bin"
fi

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
        yarn build
        cd $root/cmd/kfs-server
        GOOS=$GOOS GOARCH=$GOARCH go build -o kfs-server
        ;;

      cli)
        cd $root/cmd/kfs-cli
        GOOS=$GOOS GOARCH=$GOARCH go build -o kfs-cli
        ;;

      electron)
        cd $root/ui
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
