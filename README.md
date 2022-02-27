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
~~Using dynamic folder paths in `Makefile`~~
- Move coverage, mocks, unit tests and service tests to separate folders

~~Use mockgen~~
~~Convert `main_test.go` to table tests~~
Consider using own mocks instead of mockgen in `main_test.go`

Look into how to test `scraper.go`