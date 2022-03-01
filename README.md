# Fight Alerts Backend

Backend service to track upcoming MMA fights, written in Go and deployed using AWS serverless architecture.

# Commands

```
# Run tests
make test

# Run script and print out results
make run
```

# TODO

# Tech debt
Improve test coverage in `scraper.go`
- one way is to create a mock html page and run `getResultsFromUrl` against it