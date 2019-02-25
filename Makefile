.PHONY : dep lint test integration coverage doc build run deploy

dep:
	dep ensure

lint:
	golangci-lint run --config .golangci.yaml ./...

test:
	mkdir -p .coverage
	go test -v -cover -coverpkg=./... -coverprofile=.coverage/unit.cover.out ./...
	gocov convert .coverage/unit.cover.out | gocov-xml > .coverage/unit.xml

integration: ;

coverage:
	mkdir -p .coverage
	gocovmerge .coverage/*.cover.out > .coverage/combined.cover.out
	gocov convert .coverage/combined.cover.out | gocov-xml > .coverage/combined.xml

doc:
	godoc -http ':9090'

build: ;

run: ;

deploy: ;
