# Name of our binary file
BINARY=logstats

VERSION=1.0.0
BUILD=`git rev-parse HEAD`

LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"

# Default Target
.DEFAULT_GOAL: ${BINARY}

# Build the project
build:
	go build ${LDFLAGS} -o ${BINARY}

# Run the project
run:
	go build ${LDFLAGS} -o ${BINARY} && ./${BINARY}

# Clean up and run
clean-run:
	make clean && make run

# Install the project and copy binary
install:
	go install ${LDFLAGS} -o ${BINARY}

# Cleans our project and deleted binary
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
	if [ -f "db.bin" ] ; then rm "db.bin" ; fi

.PHONY: clean install
