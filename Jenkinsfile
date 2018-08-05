pipeline {
    agent none

    stages {
	stage ('build') {
	    agent {
		docker 'golang:latest'
	    }
	    steps{
		sh 'echo $GOROOT'
		sh 'echo $GOPATH'
		sh 'go build' 
	    }
	}
    }
}
