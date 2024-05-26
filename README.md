# README

## Overview
The app.py is a Flask application that interacts with a Kubernetes cluster to manage images on different nodes. 
It provides several HTTP endpoints to list, delete, and pull images on specific nodes. 
It also provides endpoints to get a list of nodes and to check and create deployments and services.

## Installation and Running
1. Clone the repository:
```shell
git clone https://github.com/your_username/your_microservice.git
```
2. Navigate to the project directory:
```shell
cd imageservice
```
3. Build the project using Golang
```shell
go build
```
4. Build the Docker image:
```shell
docker build -t <you_image_url> .
```
5. Deploying to Kubernetes:
   Make sure you have a Kubernetes cluster set up
```shell
kubectl apply -f deployment.yaml
```

## API Endpoints

### GET /update/monitor
This endpoint checks and creates deployments and services for each node in the Kubernetes cluster. 
It does not require any input parameters and returns a list of nodes for which deployments and services were created.\
e.g.
```shell
["ices04-x11dai-n", "icespve01-standard-pc-i440fx-piix-1996"]
```
### GET /get/nodes
This endpoint returns a list of all nodes in the Kubernetes cluster. It does not require any input parameters. 
The response is a JSON array of node names.  \
e.g.
```shell
["ices04-x11dai-n", "icespve01-standard-pc-i440fx-piix-1996", "icespve02-standard-pc-i440fx-piix-1996",
"icespve03-standard-pc-i440fx-piix-1996", "icespve04--standard-pc-i440fx-piix-1996"]
```
### POST /list
This endpoint returns a list of all images on a specific node. 
It requires a JSON object with a node_name field in the request body. 
The response is a JSON array of image details including image name, version, and size.  \
e.g.\
request:
```shell
{
    "node_name": "icespve01-standard-pc-i440fx-piix-1996"
}
```
response:
```shell
[
    {
        "imageName":"192.168.1.104:5000/cloud-collaboration-platform/depd-analysis-service",
        "imageVersion":"0.1.4",
        "imageSize":95589785
    },
    {
        "imageName":"192.168.1.104:5000/cloud-collaboration-platform/depd-analysis-service",
        "imageVersion":"0.1.8",
        "imageSize":95654224
    }
]
```
### POST /delete
This endpoint deletes specific images on a specific node. 
It requires a JSON object with a node_name field and an image field in the request body. 
The response is a JSON object with the result of the deletion operation.\
e.g.\
request:
```shell
{
  "node_name": "icespve01-standard-pc-i440fx-piix-1996",
  "image": "192.168.1.104:5000/cloud-collaboration-platform/depd-analysis-service:0.1.4"
}
```
response:
```shell
{
  "image": "192.168.1.104:5000/cloud-collaboration-platform/depd-analysis-service:0.1.4",
  "success": "delete successfully"
}
```
### POST /pull
This endpoint pulls specific images on a specific node. 
It requires a JSON object with a node_name field and an image field in the request body. 
The response is a JSON object with the result of the pull operation.\
e.g.\
request:
```shell
{
  "node_name": "icespve01-standard-pc-i440fx-piix-1996",
  "image": "192.168.1.104:5000/cloud-collaboration-platform/depd-analysis-service:0.1.4"
}
```
response:
```shell
{
  "image": "192.168.1.104:5000/cloud-collaboration-platform/depd-analysis-service:0.1.4",
  "success": "pull successfully"
}
```