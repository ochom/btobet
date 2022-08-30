SHELL=/bin/bash

test:
	source env.sh && go test -v ./...
