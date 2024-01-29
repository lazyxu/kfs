set -e -x

imeage_name=kfs-server
docker rm -f ${imeage_name} || true
docker images | grep ${imeage_name} | awk '{printf "%s:%s ",$1,$2}' | xargs docker rmi || true
image_tag=`date '+%Y%m%d_%H%M%S'`
docker build -t ${imeage_name}:${image_tag} .
docker run --privileged=true -p 1123:1123 -p 1124:1124 --name ${imeage_name} \
    -v ~/code:/root/code -w /root/code/kfs -d \
    ${imeage_name}:${image_tag}
