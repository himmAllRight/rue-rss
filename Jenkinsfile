pipeline {
    agent {
	label 'mr-mime'
    }

    stages {
	stage ('build') {
	    steps{
		sh 'cd ..'
		sh 'ls -la'
		sh 'go build' 
	    }
	}
    }
}
