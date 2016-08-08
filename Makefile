# Name of our binary file
BINARY=logstats

VERSION=1.0.0
BUILD=`git rev-parse HEAD`

LDFLAGS=-ldflags "-X main.Version ${VERSION} -X main.BuildTime ${BUILD}"

# Default Target
.DEFAULT_GOAL: ${BINARY}

# Build the project
build:
	go build ${LDFLAGS} -o ${BINARY}

run:
	go build ${LDFLAGS} -o ${BINARY} && ./${BINARY}

# Install the project and copy binary
install:
	go install ${LDFLAGS} -o ${BINARY}

# Cleans our project and deleted binary
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

.PHONY: clean install
