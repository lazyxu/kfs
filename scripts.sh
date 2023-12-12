#!/bin/bash

set -e

root=$(cd "$(dirname "$0")"; pwd)

cp pb/fs.proto ui/web/public
cp pb/fs.proto ui

is_command_exist () {
  which $1 >/dev/null 2>&1
}

# https://go.dev/dl/
# https://goproxy.cn/
if ! is_command_exist protoc-gen-go; then
  GOOS= GOARCH= go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
fi

if ! is_command_exist protoc-gen-go-grpc; then
  GOOS= GOARCH= go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
fi

export GO111MODULE=on
export PATH="$PATH:$(go env GOPATH)/bin"

rm -f pb/*.g

# https://github.com/protocolbuffers/protobuf/releases
protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. pb/fs.proto

export PATH="$PATH:/c/Users/dell/lib"

#OUTPUT_DIR=./ui/src/pb
#protoc -I=pb fs.proto \
#  --js_out=import_style=commonjs:. \
#  --grpc-web_out=import_style=commonjs,mode=grpcwebtext:"${OUTPUT_DIR}"
# sed -i "" "1i\\"$'\n'" /* eslint-disable */"$'\n' ${OUTPUT_DIR}/fs_pb.js
# sed -i "" "1i\\"$'\n'" /* eslint-disable */"$'\n' ${OUTPUT_DIR}/fs_pb_service.js

usage () {
  echo 'Usage:
  bash scripts.sh start [web|electron]
  bash scripts.sh build [server|cli|electron]
  bash scripts.sh unittest [go|js]
  bash scripts.sh benchmark [go]'
}

cliTest () {
  echo "-------- kfs-cli-test: storage $1, database: $2 --------"
  export kfs_test_storage_type=$1
  export kfs_test_database_type=$2
  cd $root/cmd/kfs-cli && go test -v ./...
}

# https://nodejs.org/en/download/current
# yarn add node-sass --sass_binary_site=https://npm.taobao.org/mirrors/node-sass/
# kfs-electron: -p 11234
# T3: CGO_ENABLED=1
# T3: --storage-dir "F://test-0807" --data-source-name "F://test-0807/kfs.db" --thumbnail-dir "F://test-0807/thumbnail" --transcode-dir "F://test-0807/transcode"

case $1 in
  start)
    case $2 in
      web)
        cd $root/ui/web
        NODE_ENV=development REACT_APP_PLATFORM=web yarn dev
        ;;

      electron)
        trap "kill 0" EXIT
        cd $root/ui/web
        tempfile=$(mktemp)
        NODE_ENV=development REACT_APP_PLATFORM=$(go env GOOS) yarn watch > $tempfile 2>&1 &
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
        cd $root/ui/web
        yarn
        DISABLE_ESLINT_PLUGIN='true' NODE_ENV=production REACT_APP_PLATFORM=web BUILD_PATH=$root/cmd/kfs-server/build homepage=/ yarn build
        cd $root/cmd/kfs-server
        echo "GOOS=$GOOS GOARCH=$GOARCH"
        if [[ $GOOS != '' && $GOARCH != '' ]]; then
          CGO_ENABLED=1 GOOS=$GOOS GOARCH=$GOARCH go build -o kfs-server-$GOOS-$GOARCH
        else
          CGO_ENABLED=1 go build -o kfs-server
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
        cd $root/cmd/kfs-electron
        echo "GOOS=$GOOS GOARCH=$GOARCH"
        GOOS=$GOOS GOARCH=$GOARCH go build -o kfs-electron
        cp kfs-electron $root/ui/web
        cd $root/ui/web
        yarn
        DISABLE_ESLINT_PLUGIN='true' NODE_ENV=production REACT_APP_PLATFORM=$(go env GOOS) BUILD_PATH=electron-production PUBLIC_URL=. yarn build
        yarn build:electron
        ;;

      *)
        usage
        ;;
    esac
    ;;

  unittest)
    case $2 in
      go)
        cd $root/storage/local && go test -v ./...
        cd $root/db/gosqlite && go test -v ./...
        cd $root/db/mysql && go test -v ./...
        cd $root/core && go test -v ./...
        cliTest 0 sqlite
        cliTest 1 sqlite
        cliTest 2 sqlite
        cliTest 3 sqlite
        cliTest 0 mysql
        cliTest 1 mysql
        cliTest 2 mysql
        cliTest 3 mysql
        ;;

      js)
        bash scripts.sh build server
        ./cmd/kfs-server/kfs-server tmp &
        backend_process=$!
        cd $root/ui/web
        yarn
        yarn test
        kill $backend_process
        ;;

      *)
        usage
        ;;
    esac
    ;;

  benchmark)
    case $2 in
      go)
        cd $root/core && go test -bench . | tee output.txt
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
