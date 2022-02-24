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
- Using dynamic folder paths in `Makefile`
    - Move coverage, mocks, unit tests and service tests to separate folders
- Set up `fight_record.go` in an OOP way so that you can use mocks 
    - Use mockgen