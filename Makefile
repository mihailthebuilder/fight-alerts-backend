.PHONY test:
	mockgen --source=scraper.go --destination=./scraper_mocks.go --package=main
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out -o ./coverage.html