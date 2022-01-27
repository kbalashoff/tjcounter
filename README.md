TJ in Go
========

TJ Counter (Sever GO + Client HTML+JS + SSE)

Run this code like:
 > go run tjcounter.go

 Then open up your browser to http://localhost:8181


Get image from Docker repo:
 > docker pull kbalashoff/tjcounter


Run Docker container:
 > docker run -d -p 8181:8181 kbalashoff/tjcounter


Deploy in Kubernetes:

 > docker tag kbalashoff/tjcounter:latest <private repo>/tjcounter:1.0
 > docker push <private repo>/tjcounter:1.0

 > kubectl create deployment kba-tj --image=<private repo>/tjcounter:1.0
 > kubectl expose deployment kba-tj --port=8181 --type=LoadBalancer --name=kba-tj-svc

 Then open up your browser to http://<exposed ip>:8181

