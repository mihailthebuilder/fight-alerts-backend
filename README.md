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

~~Set up lambda v1 in S3 bucket with cloudwatch~~

~~Reorganise `main.tf`...~~
- ~~move out outputs to `outputs.tf`~~
- ~~create variables~~
    - ~~resource tags~~
    - ~~function names~~

Set up Jenkins deployment
- ~~get Jenkins running locally using Docker server~~
- make generic pipeline work using the repo hosted on GitHub
  - you can't use the local repo because the files aren't being copied into the container
  - use [this guide](https://www.jenkins.io/doc/pipeline/tour/hello-world/#examples)
  - I'm struggling to connect the Jenkins pipeline to the repo on GitHub. I have some access rights issues, [this](https://stackoverflow.com/questions/61105368/how-to-use-github-personal-access-token-in-jenkins/61105369#61105369) might offer an explanation. A good way might be to somehow set up the access using Docker CLI, then have a go at it again.
- customise pipeline

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

[Guide]((https://levelup.gitconnected.com/setup-your-go-lambda-and-deploy-with-terraform-9105bda2bd18)) on how to use Go with AWS Lambda & Terraform.

You don't need `go build` to ignore test files; it [already does so](https://stackoverflow.com/a/65844817/7874516).

I got Jenkins running locally using a [Docker server](https://www.jenkins.io/doc/book/installing/docker/). Tried setting up on Windows 10, but the security certs were [blocking the download of plugins](https://stackoverflow.com/questions/24563694/jenkins-unable-to-find-valid-certification-path-to-requested-target-error-whil#:~:text=That%20error%20is%20a%20common,is%20a%20Self-Signed%20Certificate).