#!/bin/bash -e
cd $(dirname $0)

function announce
{
  echo
  echo $@
}

PATH=$HOME/go/bin:$PATH

unset GOPATH

go mod download

# delete artefacts from previous build (if any)
mkdir -p reports
rm -f reports/*.out reports/*.html */*.txt demo/*_sql.go

gofmt -l -w *.go
go vet ./...
go test ./...

echo .
go test . -covermode=count -coverprofile=reports/dot.out .
go tool cover -func=reports/dot.out
[ -z "$COVERALLS_TOKEN" ] || goveralls -coverprofile=$d.out -service=travis-ci -repotoken $COVERALLS_TOKEN
