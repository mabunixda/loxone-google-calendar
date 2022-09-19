
CONTAINER=mabunixda/loxcalendar
TODAY=`date +'%Y%m%d'`

container:
	docker build -t ${CONTAINER}:${TODAY} .

container-publish:
	docker tag ${CONTAINER}:${TODAY} ${CONTAINER}:latest
	docker push ${CONTAINER}:latest

golang: goreq
	go build -o loxonegogooglecalendar

all: container

goreq:
	go mod download
