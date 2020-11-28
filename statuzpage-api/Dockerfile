FROM golang:1.11 as builder

RUN git clone https://github.com/nopp/statuzpage-api.git

WORKDIR statuzpage-api/

RUN go get -d -v github.com/go-sql-driver/mysql \
	&& go get -d -v github.com/gorilla/mux \
	&& go build -o statuzpage-api .

FROM alpine:latest

LABEL maintainer "Carlos Augusto Malucelli <malucellicarlos@gmail.com>"

ENV mysqlhost localhost
ENV mysqluser user
ENV mysqlpass pass
ENV token 123

RUN apk --no-cache add ca-certificates openssl curl \
	&& curl -o /etc/apk/keys/sgerrand.rsa.pub https://alpine-pkgs.sgerrand.com/sgerrand.rsa.pub \
	&& curl -LO https://github.com/sgerrand/alpine-pkg-glibc/releases/download/2.28-r0/glibc-2.28-r0.apk \
	&& apk add glibc-2.28-r0.apk \
	&& mkdir /etc/statuzpage-api

COPY  config.json /etc/statuzpage-api/config.json

RUN sed -i s/xxxmysqlhostxxx/"$mysqlhost"/ /etc/statuzpage-api/config.json \
	&& sed -i s/xxxmysqluserxxx/"$mysqluser"/ /etc/statuzpage-api/config.json \
	&& sed -i s/xxxmysqlpassxxx/"$mysqlpass"/ /etc/statuzpage-api/config.json \
	&& sed -i s/xxxtokenxxx/"$token"/ /etc/statuzpage-api/config.json

WORKDIR /root/

COPY --from=builder /go/statuzpage-api/statuzpage-api .

EXPOSE 8000 

CMD ["./statuzpage-api"]
