node {
	def commit="null"
	stage('fetch') {
		checkout scm
		sh "git rev-parse HEAD > commit"
		commit=readFile('commit').trim()
	}
	stage('build') {
		sh "docker build -f Dockerfile.test -t ${env.JOB_NAME}:${commit} ."
	}
	stage('test') {
		ansiColor('xterm') {
			sh "docker run --privileged -t ${env.JOB_NAME}:${commit} make test"
		}
	}
}
