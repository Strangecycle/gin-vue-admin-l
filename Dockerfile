FROM golang:alpine

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.io,direct
WORKDIR /work/gva_server
COPY . .

RUN go env && go build -o server .

# alpine:latest 轻量级 linux 操作系统
FROM alpine:latest
LABEL MAINTAINER="xiziangup@163.com"

WORKDIR /work/gva_server

# 将上一阶段编译的可执行文件拷贝到当前工作目录
COPY --from=0 /work/gva_server/server ./
COPY --from=0 /work/gva_server/config.yaml ./
COPY --from=0 /work/gva_server/resource ./resource

# docker run 时运行的可执行文件
ENTRYPOINT ./server
