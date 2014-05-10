.PHONY: all test fmt benchmark git-add-hook

all: test fmt vet

vet:
	go vet ./...

%.go: %.rl
	ragel -Z $<
	gofmt -w $@

build: query/graphiteParser/parser.go
	go build ./bin/MetricBase

clean:
	rm -f MetricBase
	go clean ./...

test: query/graphiteParser/parser.go
	go test ./...

fmt:
	go fmt ./...

benchmark: query/graphiteParser/parser.go
	go test ./... -bench=".*"

git-pre-commit-hook:
	curl -s 'http://tip.golang.org/misc/git/pre-commit?m=text' > .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit
