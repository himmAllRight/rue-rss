pipeline {
    agent none

    stages {
	stage ('build') {
	    agent {
		docker 'golang:latest'
	    }
	    steps{
		sh 'ls -la'
		sh 'go build' 
	    }
	}
    }
}
