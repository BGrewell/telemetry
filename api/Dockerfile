# Dockerfile.protogen
FROM golang:latest

LABEL maintainer="Benjamin Grewell <benjamin.grewell@intel.com>"

ENV PROTOC_VERSION 3.6.1
ENV PROTOC_GEN_GO_VERSION v1.2.0

WORKDIR /go/src/github.com/BGrewell/tgams/api

RUN apt update
RUN apt install -y protobuf-compiler python3 python3-pip
RUN go get -u github.com/golang/protobuf/protoc-gen-go
RUN pip3 install grpcio-tools
RUN export PATH=$PATH:$GOPATH/bin
RUN echo $PATH

COPY proto proto

RUN mkdir go
RUN protoc -I proto/ --go_out=plugins=grpc:go proto/tgams.proto

RUN mkdir csharp
RUN protoc -I proto/ --csharp_out=csharp proto/tgams.proto

RUN mkdir python
RUN python3 -m grpc_tools.protoc -I proto/ --python_out=python/. --grpc_python_out=python/. proto/tgams.proto