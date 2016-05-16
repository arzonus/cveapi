FROM golang:latest

WORKDIR /opt/src/github.com/arzonus/cveapi
ADD . /opt/src/github.com/arzonus/cveapi
ENV GOPATH=/opt

RUN go get
RUN go build -o /opt/bin/cveapi
ENTRYPOINT ["/opt/bin/cveapi"]