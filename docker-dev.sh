set -e -x

imeage_name=kfs-dev
image_tag=`date '+%Y%m%d_%H%M%S'`
docker build -f Dockerfile.dev -t ${imeage_name}:${image_tag} .
docker run --privileged=true -p 2223:22 -p 2123:1123 -p 2124:1124 --name ${imeage_name}_${image_tag} \
    -v ~/code:/root/code -w /root/code/kfs -d \
    ${imeage_name}:${image_tag} tail -f /dev/null
docker cp ~/.ssh ${imeage_name}_${image_tag}:/root/.ssh
# git config --global --add url."git@github.com:".insteadOf "https://github.com/"
docker cp ~/.gitconfig ${imeage_name}_${image_tag}:/root/.gitconfig
docker exec -it ${imeage_name}_${image_tag} /bin/bash
