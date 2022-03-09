pipeline {
    agent any

    tools {
        go 'Go 1.17.8'
    }

    stages {
        stage("prepare environment") {
            steps {
                sh """
                    make --version
                """
            }
        }

        stage("test") {
            steps {
                sh """
                    cd functions
                    ls
                    make test-jenkins
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
