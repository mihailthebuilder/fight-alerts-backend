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
                    cd functions
                    make test
                """
            }
        }

        stage("build") {
            steps {
                sh """
                    cd functions
                    make build
                """
            }
        }

        stage("deploy") {
            steps {
                sh """
                    printenv AWS_SHARED_CREDENTIALS_FILE
                    cd terraform
                    terraform init
                    terraform apply
                """
            }
        }
    }
}
