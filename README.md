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
- [Postgres](https://www.postgresql.org/) with `username` user, `password` password and `test` table

[makefile](./functions/makefile) has all the instructions for developing locally.

Deployment is handled by local Jenkins server according to instructions in [Jenkinsfile](./Jenkinsfile).

# TODO

Set up AWS Aurora Postgres db to write the data to
- ~~write Cucumber test~~
- write code the TDD way
  - ~~get `TestInsertFightRecords` passing~~
  - ~~get `service_test` passing~~
- write terraform & deploy
- tidy up...
  - ~~export common code from `datastore_test` and `service_test`/`aurora_client`~~
  - service test
    - replace `GetHostName()` with setting the localhost name in the `context`
    - get Colly to connect to site in first service test

Figure out how to do the notificiation sender
- maybe a lambda that continuosly checks the db and if it's close to event, it gets triggered

# Technical debt

Improve test coverage in `scraper.go`
- one way is to create a mock html page and run `getResultsFromUrl` against it
    - but see why it doesn't get triggered with `espn.co.uk` test in `scraper_test.go`

Set up right access policies for AWS resources

Split service vs unit test logging clearer

Figure out how to aggregate coverage results for unit tests

Consider building the url into the Scraper as opposed to having it declared directly

Test [main.go](functions/main.go)

Separate Go source code from the rest (e.g. `bin` and `test_results` folder)

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