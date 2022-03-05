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

~~Set up S3 bucket in terraform to store the code~~

~~Write instructions to...~~
- ~~compile the code~~
- ~~archive it~~
- ~~send it to the S3 bucket~~

~~Set up lambda in terraform~~

~~Set up CloudWatch~~

~~Test lambda~~

Move `scraper` and `handler` to a separate package

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

Set up right access policies for AWS resources

Check if you're compiling test files into binary. If so, remove them as we've done in submission tracker.

# Resources

[Official guide](https://learn.hashicorp.com/tutorials/terraform/lambda-api-gateway) on how to use Terraform with lambda.

[Guide]((https://levelup.gitconnected.com/setup-your-go-lambda-and-deploy-with-terraform-9105bda2bd18)) on how to use Go with AWS Lambda & Terraform.