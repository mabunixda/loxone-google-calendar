FROM alpine:latest 
MAINTAINER Martin Buchleitner "martin@nitram.at"
RUN apk --no-cache add ca-certificates
COPY loxonegogooglecalendar /opt/loxonegogooglecalendar
RUN chmod 755 /opt/loxonegogooglecalendar

ENTRYPOINT ["/opt/loxonegogooglecalendar"]
EXPOSE 8080
CMD [""]
