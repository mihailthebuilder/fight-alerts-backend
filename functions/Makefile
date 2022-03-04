BASE_DIR_LINUX=$(shell pwd)
BASE_DIR_WINDOWS=$(shell cygpath -m $(BASE_DIR_LINUX))
TEST_RESULTS_DIR=$(BASE_DIR_WINDOWS)/test_results

build:
	GOOS=linux go build -o ./bin/scraper

test:
#	mockgen --source=scraper.go --destination=./scraper_mocks.go --package=main
	go test -coverprofile=$(TEST_RESULTS_DIR)/coverage.out
	go tool cover -html=$(TEST_RESULTS_DIR)/coverage.out -o $(TEST_RESULTS_DIR)/coverage.html

open-coverage:
	start chrome "$(TEST_RESULTS_DIR)/coverage.html"