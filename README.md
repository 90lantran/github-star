# github-star

Future work:
    - connect o github with go-github https://github.com/google/go-github, no authentication for now
    - no authentication to use github API because I dont want to commit by credentails -> rate limit

    - rate limit: how to avoid it need to log in

unit test: 
    - make unit-test (need to install go convey first)
    - show code coverage: make code-coverage

build dokcer image: make create-image

run locally: make local-server


    




deployment in minikube: 
    - update the new image to dockerhub first 
    - install hypervisor: hyperkit (brew install hyperkit) or install virtualbox
    - install minikube : brew install minikube
    - start minikuke with driver virtualbox : minikube start --vm-driver=virtualbox
    - kubectl apply -f deployment/server.yaml
    - minikube service my-server-service
    - minikube delete

generate swagger api: go-swagger, install : brew install go-swagger



