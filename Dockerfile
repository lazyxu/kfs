FROM ubuntu:22.04

RUN apt-get update
RUN apt-get install -y wget xz-utils

ARG arch

WORKDIR /root

RUN wget https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-${arch}-static.tar.xz && \
    tar -xvf ffmpeg-release-${arch}-static.tar.xz && \
    rm -rf ffmpeg-release-${arch}-static.tar.xz 

ADD cmd/kfs-server/kfs-server /root/kfs-server

ENV KFS_ROOT /root/kfs-root
ENV WebServer 1123
ENV SocketServer 1124
ENV DatabaseType sqlite
ENV DataSourceNameStr $KFS_ROOT/sqlite.db
ENV StorageType 1
ENV StorageDir $KFS_ROOT
ENV ThumbnailDir $KFS_ROOT/thumbnail
ENV TransCodeDir $KFS_ROOT/transcode

ENV PATH="/root/ffmpeg-6.1-${arch}-static:${PATH}"

CMD /root/kfs-server \
    --web-server $WebServer \
    --socket-server $SocketServer \
    --database-type $DatabaseType \
    --data-source-name $DataSourceNameStr \
    --storage-type $StorageType \
    --storage-dir $StorageDir \
    --thumbnail-dir $ThumbnailDir \
    --transcode-dir $TransCodeDir
