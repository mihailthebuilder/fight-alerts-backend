# Fight Alerts Backend

Backend service to track upcoming MMA fights, written in Go and deployed using AWS serverless architecture.

# Features completed

Service that scrapes a site for upcoming fights

AWS lambda that runs the service

Terraform deployment instructions

Jenkins test & deployment pipeline, with local Jenkins server hosted in a separate repository

Unit tests & Cucumber service tests

# Development

Pre-requisites:
- [Go](https://go.dev/)
- [Make](https://www.gnu.org/software/make/)
- [Terraform](https://www.terraform.io/)
- [Jenkins](https://www.jenkins.io/)
- [AWS CLI](https://aws.amazon.com/cli/)

[makefile](./functions/makefile) has all the instructions for developing locally.

Deployment is handled by local Jenkins server according to instructions in [Jenkinsfile](./Jenkinsfile).

# TODO

~~Better file org~~
- ~~move `mmaUrl` to separate package for reuse across other packages~~
- ~~move `scraper` to separate package~~
- ~~get all tests to still run~~
- ~~consider moving `resources.MmaUrl` to `scraper` package~~
- ~~move service test to separate folder so you can run the test separate from unit test~~
- ~~move `handler` to separate package~~

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

Set up right access policies for AWS resources

Split service vs unit test logging clearer

Figure out how to aggregate coverage results for unit tests

Consider building the url into the Scraper as opposed to having it declared directly

Test [main.go](functions/main.go)

# Log

[Official guide](https://learn.hashicorp.com/tutorials/terraform/lambda-api-gateway) on how to use Terraform with lambda.

[Guide](https://levelup.gitconnected.com/setup-your-go-lambda-and-deploy-with-terraform-9105bda2bd18) on how to use Go with AWS Lambda & Terraform.

You don't need `go build` to ignore test files; it [already does so](https://stackoverflow.com/a/65844817/7874516).

I got Jenkins running locally using a [Docker server](https://www.jenkins.io/doc/book/installing/docker/). Tried setting up on Windows 10, but the security certs were [blocking the download of plugins](https://stackoverflow.com/questions/24563694/jenkins-unable-to-find-valid-certification-path-to-requested-target-error-whil#:~:text=That%20error%20is%20a%20common,is%20a%20Self-Signed%20Certificate).

The terraform backend is stored in an S3 bucket so the local Jenkins server can access it. I then use the `-force-copy` option with `terraform init` in order to avoid Terraform asking me how to manage the new state in the Jenkins server vs the existing state in the S3 bucket. See [Terraform docs](https://www.terraform.io/cli/commands/init#backend-initialization) for more.

You can't get coverage results from service test because...
a. You're using the `bin/scraper` binary instead of the source code
b. The binary is placed in a lambda that runs it
c. The lambda is ran inside a container