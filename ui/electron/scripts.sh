#!/bin/bash

copy_kfs_client () {
  mkdir -p public/extraResources
  cd ../../cmd/client
  go build -o kfs-client
  cp kfs-client ../../ui/electron/public/extraResources/
  cp *.pem ../../ui/electron/public/extraResources/
  cd ../../ui/electron
}

electron_dev () {
  export ELECTRON_START_URL=http://localhost:3005
  electron .
}

react_dist () {
  cd ../desktop
  mkdir -p public/extraResources
  cp ../../cmd/client/*.pem public/extraResources/
  yarn build
  cd ../electron
}

electron_dist () {
  electron-builder -c electron-builder.yml
}

case $1 in
  start)
    copy_kfs_client
    electron_dev
    ;;

  dist)
    react_dist
    copy_kfs_client
    electron_dist
    ;;

  *)
    echo "invalid arg: ", $1
    ;;
esac