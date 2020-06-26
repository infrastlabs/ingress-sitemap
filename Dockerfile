# base build image
FROM registry.cn-shenzhen.aliyuncs.com/infrastlabs/golang:1.13.9-alpine3.10 as gomod
ARG dir=/go/src/devcn.fun/infrastlabs/ingsitemap

# use go modules
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn

# only copy go.mod and go.sum
WORKDIR ${dir}
COPY go.mod .
COPY go.sum .

RUN go mod download

# Build the manager binary
# FROM registry.cn-shenzhen.aliyuncs.com/infrastlabs/golang:1.13.9-alpine3.10 as builder
FROM gomod AS builder
ARG dir=/go/src/devcn.fun/infrastlabs/ingsitemap

# use go modules
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn

# Build
RUN domain="mirrors.aliyun.com" \
&& echo "http://$domain/alpine/v3.8/main" > /etc/apk/repositories \
&& echo "http://$domain/alpine/v3.8/community" >> /etc/apk/repositories \
&& apk add curl tree bash git && go get -u github.com/go-bindata/go-bindata/...

# Copy in the go src
WORKDIR ${dir}
COPY . .

RUN pwd && bash go-build.sh
# RUN pwd && ls -h && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o ingsitemap ${dir}/

# Copy data into a empty image
FROM registry.cn-shenzhen.aliyuncs.com/infrastlabs/alpine-ext
MAINTAINER sam <sldevsir@126.com>
ARG dir=/go/src/devcn.fun/infrastlabs/ingsitemap

USER root
WORKDIR /app
COPY --from=builder ${dir}/ingsitemap /app

# Configure Docker Container
# VOLUME ["/data"]
# EXPOSE 8000
ENTRYPOINT ["./ingsitemap"]
