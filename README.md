# kfs

[![UnitTest](https://github.com/lazyxu/kfs/actions/workflows/UnitTest.yml/badge.svg)](https://github.com/lazyxu/kfs/actions/workflows/UnitTest.yml)

Koala file system is a network disk used to store personal files. User-friendly clients can be obtained on various platforms. It is specially optimized for the storage and display of photos and videos.

## project goals

There are many network disk applications now, but I haven’t found one that suits my needs. Currently Synology Photos is the closest to these goals.

I want it to be:

1. There are no storage size and traffic restrictions, which requires it to be open source and be able to be deployed on your own server
2. It can be accessed on various platforms, including web pages, desktops, mobile phones, command lines, and can also be mounted to local file systems.
3. There are some useful apps, such as photo albums, notes, etc.
4. Only one copy of duplicate files is saved to save storage space.
5. Sync automatically and quickly.
6. Able to synchronize the contents of other network disks.
7. Can use other network disks to store files to prevent file loss and damage.

## environment configuration

### development environment

You can run `bash docker-dev.sh` to quickly build a docker-based development environment.

### running environment

You can run `bash docker-dev.sh` to quickly build a docker-based running environment.

## running guidelines

All running scripts are concentrated in script.sh, you can run `./scripts.sh` to get usage.

### kfs-server

The server is used to maintain files and their metadata, and responsible for communicating with various clients.

### kfs-cli

Synchronize files with the server through the command line.

### kfs-web

Obtaining services through web pages, some capabilities are missing.

### kfs-electron

The desktop application uses the electron framework.

### kfs-mobile

The mobile application is developed using React Native and Expo SDK.


