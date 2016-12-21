# Reference Logging Application
Golang application that generates several different types of logs. This was used to test a centralized logging system for a Kubernetes cluster. 


To run locally: 
$ go build log_app.go

$ ./log_app

To run in Kubernetes cluster:
$ kubectl create -f pod.yaml

$ kubectl create -f service.yaml
