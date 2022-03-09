pipeline {
    agent any

    tools {
        go 'Go 1.17.8'
    }

    stages {
        stage("prepare environment") {
            steps {
                sh """
                    cd functions
                    make test-jenkins
                """
            }
        }

        stage("test") {
            steps {
                sh """
                    cd functions
                    ls
                    make test
                """
            }
        }

        stage('build') {
            steps {
                sh """
                    make build
                    cd bin
                    ls
                """
            }
        }
    }
}
