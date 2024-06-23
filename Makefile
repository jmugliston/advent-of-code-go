BINARY=aoc

VERSION=`git describe --tags`
BUILD=`date -u +%Y%m%d.%H%M%S`

LDFLAGS=-ldflags "-w -s -X main.Version=${VERSION} -X main.Build=${BUILD}"

build:
	@echo "Building: ${BINARY} ${VERSION} ${BUILD}"
	@go build ${LDFLAGS} -o ${BINARY}

install:
	@echo "Installing: ${BINARY} ${VERSION} ${BUILD}"
	@go install ${LDFLAGS}

clean:
	@echo "Cleaning: ${BINARY}"
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

.PHONY: clean install
