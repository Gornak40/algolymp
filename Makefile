UTILS = $(shell ls ./cmd)

build:
	go build -v ./...

clean:
	rm ${UTILS}

test:
	echo "No tests"
