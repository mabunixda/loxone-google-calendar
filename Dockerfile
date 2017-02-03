FROM alpine:latest 
MAINTAINER Martin Buchleitner "martin@nitram.at"
ENV HTTP_PROXY http://squid.avl.com:8080
ENV HTTPS_PROXY http://squid.avl.com:8080

RUN apk --no-cache add ca-certificates
COPY loxonegogooglecalendar /opt/loxonegogooglecalendar
RUN chmod 755 /opt/loxonegogooglecalendar
EXPOSE 8080
ENTRYPOINT ["/opt/loxonegogooglecalendar"]

