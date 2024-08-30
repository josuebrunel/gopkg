test:
	go test -failfast ./... -v -p=1 -count=1 -coverprofile .coverage.txt
	go tool cover -func .coverage.txt

all: test
