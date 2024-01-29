FROM ubuntu:22.04

ADD cmd/kfs-server/kfs-server /kfs-server

WORKDIR /

ENV KFS_ROOT /root/kfs
ENV WebServer 1123
ENV SocketServer 1124
ENV DatabaseType sqlite
ENV DataSourceNameStr $KFS_ROOT/sqlite.db
ENV StorageType 1
ENV StorageDir $KFS_ROOT
ENV ThumbnailDir $KFS_ROOT/thumbnail
ENV TransCodeDir $KFS_ROOT/transcode

CMD /kfs-server \
    --web-server $WebServer \
    --socket-server $SocketServer \
    --database-type $DatabaseType \
    --data-source-name $DataSourceNameStr \
    --storage-type $StorageType \
    --storage-dir $StorageDir \
    --thumbnail-dir $ThumbnailDir \
    --transcode-dir $TransCodeDir
