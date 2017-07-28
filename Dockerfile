FROM golang:1.8-alpine3.6 as builder
COPY * /go/src/
WORKDIR /go/src/
RUN apk --no-cache add make git \
    && make golang

FROM alpine:3.6
MAINTAINER Martin Buchleitner "martin@nitram.at"
RUN apk --no-cache add ca-certificates && mkdir /.credentials
COPY --from=builder /go/src/loxonegogooglecalendar /loxonegogooglecalendar
WORKDIR  /
EXPOSE 8080
ENTRYPOINT  ["/loxonegogooglecalendar"]

