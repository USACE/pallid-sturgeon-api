FROM golang:1.24.6-alpine AS builder
RUN apk add build-base
# Install Git
RUN apk update && apk add --no-cache git

# Copy In Source Code
WORKDIR /go/src/app
COPY . .

# Install Dependencies
RUN go get -d -v
# Build
RUN go get -d -v \
  && GOOS=linux GOARCH=amd64 CGO_ENABLED=1 \
    go build -ldflags="-w -s" -o /go/bin/pallid_sturgeon_api

# not SCRATCH IMAGE
# remove curl -o instantclient-basiclite.zip https://download.oracle.com/otn_software/linux/instantclient/instantclient-basiclite-linuxx64.zip -SL && \
FROM alpine:latest
RUN apk add build-base
# RUN apk add --no-cache bash
RUN apk update
RUN apk upgrade 
RUN apk --no-cache add libaio libnsl libc6-compat curl && \
    cd /tmp && \
    curl -o instantclient-basiclite.zip https://download.oracle.com/otn_software/linux/instantclient/2120000/instantclient-basic-linux.x64-21.20.0.0.0dbru.zip -SL && \
    unzip instantclient-basiclite.zip && \
    mv instantclient*/ /usr/lib/instantclient && \
    rm instantclient-basiclite.zip && \
    ln -s /usr/lib/instantclient/libclntsh.so.21.20 /usr/lib/libclntsh.so && \
    ln -s /usr/lib/instantclient/libocci.so.21.20 /usr/lib/libocci.so && \
    ln -s /usr/lib/instantclient/libociicus.so /usr/lib/libociicus.so && \
    ln -s /usr/lib/instantclient/libnnz19.so /usr/lib/libnnz19.so && \
    ln -s /usr/lib/libnsl.so.2 /usr/lib/libnsl.so.1 && \
    ln -s /lib/libc.so.6 /usr/lib/libresolv.so.2 && \
    ln -s /lib64/ld-linux-x86-64.so.2 /usr/lib/ld-linux-x86-64.so.2

ENV ORACLE_BASE /usr/lib/instantclient
ENV LD_LIBRARY_PATH /usr/lib/instantclient
ENV TNS_ADMIN /usr/lib/instantclient
ENV ORACLE_HOME /usr/lib/instantclient

COPY --from=builder /go/bin/pallid_sturgeon_api /go/bin/pallid_sturgeon_api
ENTRYPOINT ["/go/bin/pallid_sturgeon_api"]