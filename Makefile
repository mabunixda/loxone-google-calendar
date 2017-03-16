
CONTAINER=r.nitram.at/loxcalendar
TODAY=`date +'%Y%m%d'`

container: golang
	docker build -t ${CONTAINER}:${TODAY} . 
	docker tag ${CONTAINER}:${TODAY} ${CONTAINER}:latest

golang: goreq
	export GOPATH=${PWD}
	go build

all: container
	
goreq:
	export GOTPATH=${PWD}
	go get golang.org/x/net/context
	go get golang.org/x/oauth2
	go get golang.org/x/oauth2/google
	go get google.golang.org/api/calendar/v3        
	go get github.com/Sirupsen/logrus
