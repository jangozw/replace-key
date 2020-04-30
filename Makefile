BINARY=replace-key
build-env:
	go build ${LDFLAGS} -o ${BINARY}
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BINARY}
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
.PHONY:  clean install
