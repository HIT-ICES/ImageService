# README

## Overview
app.py是一个Flask应用程序，它与Kubernetes集群交互以管理不同节点上的映像。它提供了几个HTTP端点来列出、删除和提取特定节点上的图像。
它还提供端点以获取节点列表，并检查和创建部署和服务。
## API

### GET /update/monitor
检查并为Kubernetes集群中的每个节点创建部署和服务。它不需要任何输入参数，也不返回任何数据。
### GET /get/nodes
返回Kubernetes集群中所有节点的列表。它不需要任何输入参数。响应是一个节点名称的JSON数组。\
形如:
```shell
["ices04-x11dai-n", "icespve01-standard-pc-i440fx-piix-1996", "icespve02-standard-pc-i440fx-piix-1996",
"icespve03-standard-pc-i440fx-piix-1996", "icespve04--standard-pc-i440fx-piix-1996"]
```
### POST /list
返回特定节点上所有镜像的列表。它需要一个JSON对象，该对象在请求正文中具有node_name字段。\
请求:\
```shell
{
    "node_name": "icespve01-standard-pc-i440fx-piix-1996"
}
```
返回:
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
删除特定节点上的特定图像。它需要一个JSON对象，该对象在请求体中具有node_name字段和image字段。\
请求:\
```shell
{
  "node_name": "icespve01-standard-pc-i440fx-piix-1996",
  "image": "192.168.1.104:5000/cloud-collaboration-platform/depd-analysis-service:0.1.4"
}
```
响应是一个JSON对象，带有删除操作的结果。\
```shell
{
  "image": "192.168.1.104:5000/cloud-collaboration-platform/depd-analysis-service:0.1.4",
  "success": "delete successfully"
}
```
### POST /pull
在特定节点上提取特定图像。它需要一个JSON对象，该对象在请求体中具有node_name字段和image字段。\
请求:\
```shell
{
  "node_name": "icespve01-standard-pc-i440fx-piix-1996",
  "image": "192.168.1.104:5000/cloud-collaboration-platform/depd-analysis-service:0.1.4"
}
```
响应是一个JSON对象，带有pull操作的结果。\
```shell
{
  "image": "192.168.1.104:5000/cloud-collaboration-platform/depd-analysis-service:0.1.4",
  "success": "pull successfully"
}
```
## 部署
运行deployment即可运行
```shell
kubectl apply -f deployment.yaml
```