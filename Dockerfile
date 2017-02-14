FROM alpine:latest 
MAINTAINER Martin Buchleitner "martin@nitram.at"

RUN apk --no-cache add ca-certificates

WORKDIR  /opt
COPY loxonegogooglecalendar /opt/loxonegogooglecalendar
RUN chmod 755 /opt/loxonegogooglecalendar 
RUN mkdir /opt/.credentials
EXPOSE 8080

ENTRYPOINT  ["/opt/loxonegogooglecalendar"]

