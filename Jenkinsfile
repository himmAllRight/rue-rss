pipeline {
    agent none

    stages {
	stage ('build') {
	    agent {
		docker 'golang:latest'
	    }
	    steps{
		sh 'cd ..'
		sh 'ls -la'
		sh 'go build' 
	    }
	}
    }
}
