import json
from http import HTTPStatus

from flask import Flask, request, jsonify
from kubernetes import client, config
import requests

app = Flask(__name__)
try:
    config.load_incluster_config()
except config.config_exception.ConfigException:
    print("Not running in a cluster. Switching to KubeConfig.")
    config.load_kube_config()


def send_post_request(url, data):
    headers = {'Content-Type': 'application/json'}
    response = requests.post(url, data=data, headers=headers)
    return response


def create_deployment(api_instance, node_name):
    deployment = client.V1Deployment(
        api_version='apps/v1',
        kind='Deployment',
        metadata=client.V1ObjectMeta(name=f'{node_name}-image-service-deployment'),
        spec=client.V1DeploymentSpec(
            replicas=1,
            selector=client.V1LabelSelector(
                match_labels={'app': f'{node_name}-image-service-app'}
            ),
            template=client.V1PodTemplateSpec(
                metadata=client.V1ObjectMeta(labels={'app': f'{node_name}-image-service-app'}),
                spec=client.V1PodSpec(
                    node_selector={"kubernetes.io/hostname": node_name},
                    containers=[
                        client.V1Container(
                            name=f'{node_name}-image-service',
                            image='docker.io/microyui/containerd:3.0',
                            ports=[client.V1ContainerPort(container_port=80)],
                            volume_mounts=[
                                client.V1VolumeMount(
                                    name="containerd-socket",
                                    mount_path="/run/containerd/containerd.sock",
                                    sub_path=None,
                                    read_only=True
                                )
                            ]
                        )
                    ],
                    volumes=[
                        client.V1Volume(
                            name="containerd-socket",
                            host_path=client.V1HostPathVolumeSource(
                                path="/run/containerd/containerd.sock"
                            )
                        )
                    ]
                )
            )
        )
    )
    api_instance.create_namespaced_deployment(namespace='default', body=deployment)


def create_service(api_instance, node_name):
    service = client.V1Service(
        api_version='v1',
        kind='Service',
        metadata=client.V1ObjectMeta(name=f'{node_name}-image-service'),
        spec=client.V1ServiceSpec(
            selector={'app': f'{node_name}-image-service-app'},
            ports=[client.V1ServicePort(port=8080, target_port=8080)]
        )
    )
    api_instance.create_namespaced_service(namespace='default', body=service)


@app.route('update/monitor', methods=['GET'])
def check_and_create_deployment_and_service():
    print("check and create deployment and service")
    core_v1 = client.CoreV1Api()
    app_v1 = client.AppsV1Api()
    ret = core_v1.list_namespaced_service(namespace='default')
    service_name_list = []
    for i in ret.items:
        service_name_list.append(i.metadata.name)
    ret = core_v1.list_node()
    for i in ret.items:
        node_name = i.metadata.name
        image_service_name = node_name + '-image-service'
        if image_service_name not in service_name_list:
            create_deployment(app_v1, node_name)
            create_service(core_v1, node_name)
            print(f'create deployment and service on {node_name}')


@app.route('/get/nodes', methods=['GET'])
def get_nodes():
    core_v1 = client.CoreV1Api()
    ret = core_v1.list_node()
    nodes = []
    for i in ret.items:
        nodes.append(i.metadata.name)
    return json.dumps(nodes)


@app.route('/list', methods=['POST'])
def get_images_of_node():
    data = request.json
    node_name = data['node_name']
    response = requests.get(f'http://{node_name}-image-service:8080/listImages')
    if response.status_code == HTTPStatus.OK:
        return response.text
    else:
        return jsonify({'error': 'Node not found'}), HTTPStatus.NOT_FOUND


@app.route('/delete', methods=['POST'])
def delete_images_of_node():
    data = request.json
    node_name = data['node_name']
    image = data['image']
    data = {"image", image}
    json_data = json.dumps(data)
    response = send_post_request(f'http://{node_name}-image-service:8080/deleteImages', json_data)
    if response.status_code == HTTPStatus.OK:
        return response.text
    else:
        return jsonify({'error': 'Node not found'}), HTTPStatus.NOT_FOUND


@app.route('/pull', methods=['POST'])
def pull_images_of_node():
    data = request.json
    node_name = data['node_name']
    image = data['image']
    data = {"image", image}
    json_data = json.dumps(data)
    response = send_post_request(f'http://{node_name}-image-service:8080/pullImages', json_data)
    if response.status_code == HTTPStatus.OK:
        return response.text
    else:
        return jsonify({'error': 'Node not found'}), HTTPStatus.NOT_FOUND


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8080)
