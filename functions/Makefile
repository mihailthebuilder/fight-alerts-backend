BASE_DIR_LINUX:=$(shell pwd)

build:
	GOOS=linux go build -o ./bin/scraper

unit-test:
	go test ./ ./scraper ./lambda_handler -coverprofile=./test_results/coverage.out
	go tool cover -html=./test_results/coverage.out -o ./test_results/coverage.html

int-test: build
	cp ./bin/scraper ./service_tests
	cd service_tests && go test

test: unit-test int-test

open-coverage-win:
	start chrome "$(shell cygpath -m $(BASE_DIR_LINUX))/test_results/coverage.html"