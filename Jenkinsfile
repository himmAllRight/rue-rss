pipeline {
    agent {
	label 'mr-mime'
    }

    stages {
	stage ('build') {
	    steps{
		withEnv(["GOROOT=${root}", "PATH+GO=${root}/bin:${HOME}/go/bin"]) {
		    sh "go get github.com/himmallright/rue-rss.git"

		}
	    }
	}
    }
}
