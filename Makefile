UTILS = $(shell ls ./cmd)

build:
	$(foreach util, ${UTILS}, go build ./cmd/${util};)

clean:
	rm ${UTILS}

test:
	echo "No tests"
