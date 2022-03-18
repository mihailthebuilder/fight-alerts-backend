# Fight Alerts Backend

Backend service to track upcoming MMA fights, written in Go and deployed using AWS serverless architecture.

# Commands

In the `functions` folder:

```
# Run tests
make test

# Check coverage results after running tests (in Windows only)
make open-coverage-win

# Run script and print out results
make run
```

In the `terraform` folder:

```
# Log in to AWS CLI
../../../aws-adfs-cli/aws-adfs

# Initialise terraform
terraform init -force-copy

# Apply terraform
terraform apply -auto-approve
```

# TODO

~~Set up lambda v1 in S3 bucket with cloudwatch~~

~~Set up Jenkins deployment~~

Move `scraper` and `handler` to a separate package

Set up Cucumber for service tests

Set up AWS RDS db to write the data to
- write terraform
- amend code to use db (with tests)
- add unit & service tests

Figure out how to do the notificiation sender
- maybe a lambda that continuosly checks the db and if it's close to event, it gets triggered

# Technical debt

Improve test coverage in `scraper.go`
- one way is to create a mock html page and run `getResultsFromUrl` against it
    - but see why it doesn't get triggered with `espn.co.uk` test in `scraper_integration_test.go`

Organise terraform files in a similar way to ST backend

Set up right access policies for AWS resources

# Log

[Official guide](https://learn.hashicorp.com/tutorials/terraform/lambda-api-gateway) on how to use Terraform with lambda.

[Guide](https://levelup.gitconnected.com/setup-your-go-lambda-and-deploy-with-terraform-9105bda2bd18) on how to use Go with AWS Lambda & Terraform.

You don't need `go build` to ignore test files; it [already does so](https://stackoverflow.com/a/65844817/7874516).

I got Jenkins running locally using a [Docker server](https://www.jenkins.io/doc/book/installing/docker/). Tried setting up on Windows 10, but the security certs were [blocking the download of plugins](https://stackoverflow.com/questions/24563694/jenkins-unable-to-find-valid-certification-path-to-requested-target-error-whil#:~:text=That%20error%20is%20a%20common,is%20a%20Self-Signed%20Certificate).

The terraform backend is stored in an S3 bucket so the local Jenkins server can access it. I then use the `-force-copy` option with `terraform init` in order to avoid Terraform asking me how to manage the new state in the Jenkins server vs the existing state in the S3 bucket. See [Terraform docs](https://www.terraform.io/cli/commands/init#backend-initialization) for more.