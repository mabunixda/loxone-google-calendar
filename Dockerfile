FROM debian:jessie
MAINTAINER Martin Buchleitner "martin@nitram.at"
RUN apt-get update
RUN apt-get install -y ca-certificates

COPY loxonegogooglecalendar /opt/loxonegogooglecalendar
RUN chmod 755 /opt/loxonegogooglecalendar

ENTRYPOINT ["/opt/loxonegogooglecalendar"]
EXPOSE 8080
CMD [""]
