pipeline {
    agent { docker { image 'golang:1.17.5-alpine' } }
    stages {
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
