FROM alpine:3.3

ENV GOROOT=/usr/lib/go \
    GOPATH=/go \
    GOBIN=/go/bin \
    PATH=$PATH:$GOROOT/bin:$GOPATH/bin

WORKDIR /go/src/github.com/arzonus/cveapi
ADD . /go/src/github.com/arzonus/cveapi

RUN apk add -U git go && \
    go get  && \
    apk del git go && \
    rm -rf /go/pkg && \
    rm -rf /go/src && \
    rm -rf /var/cache/apk/*

ENTRYPOINT ["/go/bin/cveapi"]