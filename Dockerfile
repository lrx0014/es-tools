FROM golang AS build-env

WORKDIR /gopath/golang/src/github.com/lrx0014/log-tools
ENV GOPATH /gopath/golang
ENV GOBIN /gopath/golang/bin
ADD . /gopath/golang/src/github.com/lrx0014/log-tools
RUN go install cmd/apiserver/apiserver.go

FROM daocloud.io/daocloud/go-busybox:glibc
COPY --from=build-env /gopath/golang/bin/apiserver /apiserver
ENV TZ Asia/Shanghai
EXPOSE 8080
CMD ["/apiserver"]


