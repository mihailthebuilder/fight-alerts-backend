pipeline {
    agent any

    tools {
        go 'Go 1.17.1'
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
