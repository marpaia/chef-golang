FROM ubuntu:12.04

RUN apt-get update
RUN DEBIAN_FRONTEND=noninteractive apt-get install -y wget curl
RUN wget --progress=bar https://opscode-omnibus-packages.s3.amazonaws.com/ubuntu/12.04/x86_64/chef-server_11.0.10-1.ubuntu.12.04_amd64.deb

RUN dpkg -i chef-server_11.0.10-1.ubuntu.12.04_amd64.deb

# super janky dance to get chef-server fully configured
RUN echo "echo 'started'" > /sbin/initctl
RUN timeout 5 chef-server-ctl reconfigure || echo 'ok'
RUN sh -c "/opt/chef-server/embedded/bin/runsvdir-start & timeout 15 chef-server-ctl reconfigure || echo 'ok'"
RUN echo 'shared_buffers = 20MB' >> /var/opt/chef-server/postgresql/data/postgresql.conf
RUN sh -c "/opt/chef-server/embedded/bin/runsvdir-start & timeout 45 chef-server-ctl reconfigure || echo 'ok'"

RUN curl -L https://go.googlecode.com/files/go1.2.linux-amd64.tar.gz | tar xzf -

ADD . /src/github.com/marpaia/chef-golang
WORKDIR /src/github.com/marpaia/chef-golang

ENV GOPATH /
ENV GOROOT /go
ENV PATH $PATH:/go/bin

RUN sh test.sh
