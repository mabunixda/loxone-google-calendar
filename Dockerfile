FROM golang:alpine as builder
COPY * /go/src/
WORKDIR /go/src/
RUN apk --no-cache add make git \
    && make golang

FROM alpine:latest
LABEL "Maintainer" "Martin Buchleitner <martin@nitram.at>"
RUN apk --no-cache add ca-certificates && mkdir /.credentials
COPY --from=builder /go/src/loxonegogooglecalendar /loxonegogooglecalendar
WORKDIR  /
EXPOSE 8080
ENTRYPOINT  ["/loxonegogooglecalendar"]

