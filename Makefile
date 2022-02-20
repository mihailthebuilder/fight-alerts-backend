.PHONY build:
	go build

.PHONY run: build
	./fight-alerts-backend.exe

.PHONY test:
	mockgen --source=scraper.go --destination=./scraper_mocks.go --package=main
	go test -coverprofile=./test_results/coverage.out
	go tool cover -html=./test_results/coverage.out -o ./test_results/coverage.html