pipeline {
    agent any

    tools {
        go 'Go 1.17.8'
        terraform 'Terraform 1.1.7'
    }

    stages {
        stage("test") {
            steps {
                sh """
                    cd functions/scraper-lambda
                    JENKINS="true" make test
                """
            }
        }

        stage("build") {
            steps {
                sh """
                    cd functions/scraper-lambda
                    make build
                """
            }
        }

        stage("deploy") {
            steps {
                sh """
                    cd terraform
                    terraform init -force-copy
                    terraform apply -auto-approve
                """
            }
        }
    }
}
