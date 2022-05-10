# Fight Alerts Backend

Backend service to track upcoming MMA fights, written in Go and deployed using AWS serverless architecture.

# Features completed

Service that scrapes a site for upcoming fights

AWS Lambda that runs the service

Terraform deployment instructions

Jenkins test & deployment pipeline, with local Jenkins server hosted in a separate repository

Unit tests & Cucumber service tests

AWS Aurora Postgres instance to host the scraped data

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

- ~~Lambda should replace data in db with new scraped records~~
- ~~Figure out how to do the notificiation sender~~
- Implement notification service as per [proposal](#proposal-for-notification-service)
  - scraper service to insert message in cloudwatch events
    - ~~Update `.feature` file to reflect there's a message inserted into cloudwatch event~~
    - implement the feature in TDD
      - You need to start by deleting all the rules in place for that target. [This API](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/eventbridge#Client.ListTargetsByRule) lists all rules for given resources, and [This](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/eventbridge#Client.DeleteRule) deletes the resource.
      - Once you've deleted the previous rule, you need to use [this API](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/eventbridge#Client.PutRule) to create the new rule. it's simpler to start off with this, then do the deletion
        - steps
          - ~~write the implementation~~
          - ~~write the unit tests~~
          - ~~write the service tests~~
          - deploy in jenkins to make sure everything still works
  - notification service to send SMS
    - add notification service in the `.feature` file 
- can't get consistent internet access with Lambda. start with [this article](https://stackoverflow.com/questions/51380018/aws-lambda-in-vpc-sometimes-doesnt-have-internet-access)
- Create architecture diagram
- Set up cloudwatch event trigger to run the scraper every day
- Look into [flyway](https://flywaydb.org/) for your db

# Technical debt

## Source code

When the scraper service is triggered, we're currently deleting existing records and inserting new ones. But if we use that table when triggering the notification service, we can get a situation where the service can't find the content it's supposed to publish because the scraper service is midway through the delete-and-insert process.
- The likelihood of that happening is extremely low in the current state. The table updates should complete very quickly as the table only has about 15 rows.

Separate Go source code from the rest (e.g. `bin` and `test_results` folder)

## Testing

Use testify mocks in the `datastore` and `scraper` packages. It'll enable you to solve the 2 issues below.

Differentiate between different panics in the assertions and check whether the methods were called. It's pretty easy to spot when we miss out panic/method calls though because the coverage results will show a gap.

Improve test coverage in `scraper.go`
- one way is to create a mock html page and run `getResultsFromUrl` against it
    - but see why it doesn't get triggered with `espn.co.uk` test in `scraper_test.go`

Split service vs unit test logging clearer

Aggregate coverage results for unit tests

Test [main.go](functions/main.go)

# Log

## Proposal for notification service

The scraper lambda creates a cloudwatch event rule that will trigger a notification lambda. The event will be scheduled for the day when the next fight event happens. [See this](https://stackoverflow.com/a/41655661/7874516).

The notification lambda just sends a message to SNS. The SNS, in turn, is configured to send an SMS message to the phone number in question.

## Tests

I'm consciously allowing the Sherdog website to be a dependency in my unit and service tests. It enables a much tighter feedback loop and it hasn't caused any issues so far other than with the firewall (see below). 

You'll need to disable your firewall in order to run the service test. Otherwise the scraper lambda won't be able to connect to the internet and it'll return an empty slice.

You don't need `go build` to ignore test files; it [already does so](https://stackoverflow.com/a/65844817/7874516).

You can't get coverage results from service test because...
1. You're using the `bin/scraper` binary instead of the source code
2. The binary is placed in a lambda that runs it
3. The lambda is ran inside a container

## AWS

You [can't modify](https://serverfault.com/questions/816820/aws-can-not-change-db-subnet-group-for-aws-rds) DB subnet groups for RDS, so you'll need to `terraform destroy` every time you make changes there.

It takes a long time to destroy all resources because of the ENI interfaces attached to the lambda's security group. You can't delete these interfaces manually. They usually get deleted automatically after 15-30 minutes. See a few mentions:
- https://stackoverflow.com/questions/41299662/aws-lambda-created-eni-not-deleting-while-deletion-of-stack
- https://stackoverflow.com/questions/58276376/deleting-orphaned-aws-eni-sg-currently-in-use-and-has-a-dependent-object
- https://www.reddit.com/r/aws/comments/dytfmy/unable_to_delete_network_interface_likely_due_to/

I used [this guide](https://aws.amazon.com/premiumsupport/knowledge-center/internet-access-lambda-function/) to get scraping Lambda to connect to the internet while inside a VPC. I also created 2 public & private subnets as per [this guide](https://jasonwatmore.com/post/2021/05/30/aws-create-a-vpc-with-public-and-private-subnets-and-a-nat-gateway).
- When I initially set up the route tables, I needed to manually make the 2nd public subnet association. But I didn't need to do it afterwards.

When you update the lambda, it takes a fwe hours until it regains internet connectivity.

The RDS is currently publicly accessible. You can't switch it to private because you'd need to establish a VPN connection into the VPC - [source](https://stackoverflow.com/a/69320090/7874516).

The lambda security groups take a really long time to destroy. You could look into creating it outside of terraform, then take its id into the app; but you won't be updating them anyway so there's no value for now.

## Terraform

[Official guide](https://learn.hashicorp.com/tutorials/terraform/lambda-api-gateway) on how to use Terraform with lambda.

[Guide](https://levelup.gitconnected.com/setup-your-go-lambda-and-deploy-with-terraform-9105bda2bd18) on how to use Go with AWS Lambda & Terraform.

The terraform backend is stored in an S3 bucket so the local Jenkins server can access it. I then use the `-force-copy` option with `terraform init` in order to avoid Terraform asking me how to manage the new state in the Jenkins server vs the existing state in the S3 bucket. See [Terraform docs](https://www.terraform.io/cli/commands/init#backend-initialization) for more.

You'll get timeout errors when you're creating/updating `aws_route_table_association` resources; it'll work if you just run the terraform script again. You can't customize the length of timeout.

## Other

I got Jenkins running locally using a [Docker server](https://www.jenkins.io/doc/book/installing/docker/). Tried setting up on Windows 10, but the security certs were [blocking the download of plugins](https://stackoverflow.com/questions/24563694/jenkins-unable-to-find-valid-certification-path-to-requested-target-error-whil#:~:text=That%20error%20is%20a%20common,is%20a%20Self-Signed%20Certificate).