pipeline {
    agent any
    tools {
        go 'Go-1.15'
    }
    environment {
        GO114MODULE = 'on'
        CGO_ENABLED = 0 
        GOPATH = "${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}"
    }
    stages {
        stage('Build') {
            steps {
                echo 'Compiling and building'
                sh 'export PATH=$PATH:/usr/local/go/bin && cd src/main && echo `pwd` && go build'
            }
        }

        stage('Test') {
            steps {
                withEnv(["PATH+GO=${GOPATH}/bin"]){
                    echo 'Running test'
                    sh 'cd src/test && go test -v'
                }
            }
        }
    }
}
