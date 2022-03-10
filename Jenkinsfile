pipeline {
    agent any

    tools {
        go 'Go 1.17.8'
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
    }
}
