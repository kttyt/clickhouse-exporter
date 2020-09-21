pipeline {
    agent any
    environment {
        IMAGE_NAME = ''
    }
    stages {
        stage('Build image'){
            steps {
                script {
                    tag = sh(returnStdout: true, script: "git tag --contains").trim()
                    
                    IMAGE_NAME = env.BRANCH_NAME == "master" ? ${IMAGE_NAME} : "${IMAGE_NAME}-${env.BRANCH_NAME}"
                    IMAGE_NAME = IMAGE_NAME.toLowerCase()
                    VERSION    = tag ?: env.BRANCH_NAME
                    DATE       = new Date().format("yyyy-MM-dd.HHmm", TimeZone.getTimeZone('UTC'))
                    GIT_SHA    = sh(returnStdout: true, script: "git rev-parse --verify HEAD").trim()

                    sh "docker build -t ${IMAGE_NAME}:latest --build-arg VERSION=${VERSION} --build-arg GIT_SHA=${GIT_SHA} --build-arg NOW='${DATE}' ."
                    sh "docker push ${IMAGE_NAME}:latest"

                    if (!tag) { return }
                    sh "docker tag ${IMAGE_NAME}:latest ${IMAGE_NAME}:${tag}"
                    sh "docker push ${IMAGE_NAME}:${tag}"
                    sh "docker rmi ${IMAGE_NAME}:${tag}"
                }
            }
        }
    }
}
