set -e -x

imeage_name=kfs-server
docker rm -f ${imeage_name} || true
docker images | grep ${imeage_name} | awk '{printf "%s:%s ",$1,$2}' | xargs docker rmi || true
image_tag=`date '+%Y%m%d_%H%M%S'`
mkdir -p .kfs-server
docker build -t ${imeage_name}:${image_tag} .
docker tag ${imeage_name}:${image_tag} ${imeage_name}:latest
docker run --privileged=true -p 1123:1123 -p 1124:1124 --name ${imeage_name} \
    -v `pwd`/.kfs:/root/kfs -w /root/kfs --restart=always -d \
    ${imeage_name}:${image_tag}

# docker run --privileged=true -p 1125:1123 -p 1126:1124 --name kfs-server-test -v /.kfs:/root/kfs -w /root/kfs --restart=always -d kfs-server:latest
