FROM ubuntu:22.04

RUN apt-get update
RUN apt-get install -y wget curl zip
RUN apt-get install -y vim

RUN apt-get install -y gcc python3 make
# g++
RUN apt-get install -y build-essential --fix-missing
RUN apt-get install -y git

# nodejs
RUN curl -fsSL https://deb.nodesource.com/setup_21.x | bash -
RUN apt-get install -y libatomic1 nodejs && npm install --global yarn

RUN yarn config set registry https://registry.npm.taobao.org
RUN yarn config set sass_binary_site https://npm.taobao.org/mirrors/node-sass/
RUN yarn config set electron_mirror https://npmmirror.com/mirrors/electron/

# golang
RUN wget https://dl.google.com/go/go1.21.6.linux-armv6l.tar.gz \
    && tar -xf go1.21.6.linux-armv6l.tar.gz
RUN /usr/bin/echo -e "\n\
export PATH=/go/bin:$PATH\n\
export GOPROXY=https://proxy.golang.com.cn,direct\n\
export GOROOT=/go\n\
export GOFLAGS=-buildvcs=false\n\
" >> ~/.bashrc

# protobuf
RUN apt-get install -y protobuf-compiler

# protoc-gen-go 
RUN export PATH=/go/bin:$PATH; \
    export GOPROXY=https://proxy.golang.com.cn,direct; \
    export GOROOT=/go; \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28; \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2;

# github
RUN /usr/bin/echo -e "\n\
curl https://hosts.gitcdn.top/hosts.txt > /etc/hosts\n\
chown -R root:root ~/.ssh\n\
" >> ~/.bashrc

# ssh
RUN apt-get install -y supervisor openssh-server
RUN echo "root:root" | chpasswd
RUN sed -i "s/#PermitRootLogin no/PermitRootLogin yes/g" /etc/ssh/sshd_config
RUN mkdir /var/run/sshd
EXPOSE 22

RUN echo -e "\n\
[supervisord]\n\
\n\
nodaemon=true\n\
[program:sshd]\n\
command=/usr/sbin/sshd -D\n\
" >> /etc/supervisord.conf

CMD ["/usr/bin/supervisord"]
