pipeline {
    agent none

    stages {
	stage ('build') {
	    agent {
		docker 'golang:latest'
	    }
	    steps{
		sh 'go build' 
	    }
	}
    }
}
