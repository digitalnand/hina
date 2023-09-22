FROM debian:stable-slim

WORKDIR /hina
COPY . .

RUN apt-get update
RUN apt-get install -y wget make

RUN wget https://go.dev/dl/go1.21.1.linux-amd64.tar.gz
RUN mv go1.21.1.linux-amd64.tar.gz /tmp/
RUN tar -xf /tmp/go1.21.1.linux-amd64.tar.gz -C /usr/local/
RUN ln -s /usr/local/go/bin/go /usr/local/bin/go

RUN make
COPY files/fib.json /var/rinha/source.rinha.json

ENTRYPOINT ["bin/hina/hina", "/var/rinha/source.rinha.json"]
