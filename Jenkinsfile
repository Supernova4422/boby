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
                sh 'export PATH=$PATH:/usr/local/go/bin && export GOROOT=/usr/local/go && cd src/main && go build'
            }
        }

        stage('Test') {
            steps {
                withEnv(["PATH+GO=${GOPATH}/bin"]){
                    echo 'Running test'
                    sh 'export PATH=$PATH:/usr/local/go/bin && export GOROOT=/usr/local/go && cd src/test && go test -v'
                }
            }
        }
    }
}
