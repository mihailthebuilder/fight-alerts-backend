BASE_DIR_LINUX:=$(shell pwd)

build:
	GOOS=linux go build -o ./bin/scraper

test: build
	go test ./ ./scraper -coverprofile=./test_results/coverage.out
	go tool cover -html=./test_results/coverage.out -o ./test_results/coverage.html

open-coverage-win:
	start chrome "$(shell cygpath -m $(BASE_DIR_LINUX))/test_results/coverage.html"