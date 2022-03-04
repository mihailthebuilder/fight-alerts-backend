# Fight Alerts Backend

Backend service to track upcoming MMA fights, written in Go and deployed using AWS serverless architecture.

# Commands

In the `functions` folder:

```
# Run tests
make test

# Check coverage results after running tests
make open-coverage

# Run script and print out results
make run
```

In the `terraform` folder:

```
# Log in to AWS CLI
../../../aws-adfs-cli/aws-adfs

# Initialise terraform
terraform init

# Apply terraform
terraform apply
```

# TODO

~~Amend code to work in a lambda~~
- use [this](https://levelup.gitconnected.com/setup-your-go-lambda-and-deploy-with-terraform-9105bda2bd18)

~~Set up S3 bucket in terraform to store the code~~
- used [this guide](https://learn.hashicorp.com/tutorials/terraform/lambda-api-gateway)

Write instructions to...
- ~~compile the code~~
- archive it
- send it to the S3 bucket

Set up lambda in terraform

Check if you can test lambda. If not, set up CloudWatch

Move `scraper` and `handler` to a separate package

Update logging to something that's recognised by a lambda

Set up Jenkins deployment

Set up AWS RDS db to write the data to
- write terraform
- amend code to use db (with tests)

Figure out how to do the notificiation sender
- maybe a lambda that continuosly checks the db and if it's close to event, it gets triggered

# Technical debt

Improve test coverage in `scraper.go`
- one way is to create a mock html page and run `getResultsFromUrl` against it
    - but see why it doesn't get triggered with `espn.co.uk` test in `scraper_integration_test.go`

Organise terraform files in a similar way to ST backend

Restrict access to S3 bucket

Check if you're compiling test files into binary. If so, remove them as we've done in submission tracker.