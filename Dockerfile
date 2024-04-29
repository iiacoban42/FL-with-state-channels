FROM ubuntu:latest
RUN apt update

FROM golang:1.20

# add the code to the docker image
RUN mkdir /app-channel
COPY ./ /app-channel

# build the code
RUN chmod +x /app-channel/*
WORKDIR /app-channel
RUN go build
