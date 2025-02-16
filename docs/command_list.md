# `karakuri` Command List
* [Container operation](#container-operations)
* [Image operation](#image-operations)
* [Namespace operation](#namespace-operations)
* [Module operation](#module-operations)
* [Registry Controller operation](#registry-controller-operation)

## Container operations
The operations related to containers are as follows:
| Operation | Description |
| --------- | ----------- |
| [`create`](#create) | Create container |
| [`start`](#start) | Start container |
| [`run`](#run) | Run container (`create`+`start`) |
| [`exec`](#exec) | Execute command in container |
| [`stop`](#stop) | Stop container |
| [`restart`](#restart) | Restart container (`stop`+`start`) |
| [`rm`](#rm) | Delete container |
| [`ls`](#ls) | Show container list |
| [`spec`](#spec) | Show container spec |
| [`logs`](#logs) | Show container logs |

### `create` 
Create container
| Option | Value | Required | Description | Example |
| ------ | ----- | -------- | ----------- | ------- |
| --image | [repositry]:[tag] | yes | Specify name of a image | --image alpine:latest |
| --name | [container_name] | no | Specify container name | --name my_container |
| --port | [host_port]:[container_port]:[protocol] | no | Map port `[host_port]` to `[container_port]/[protocol]` in the container (*1) | --port 8080:80:tcp |
| --mount | [host_path]:[container_path] | no | Mount `[host_path]` to `[container_path]` in the container (*2) | --mount /mnt/data:/data |
| --cmd | [arg_1],[arg2],... | no | Override entrypoint command | --cmd sleep,100 |
| --ns | [namespace] | no | Specify the namespace to which the container belongs | --ns sandbox |
| --registry | [registry]:[port] | no | Specify registry | --registry my.registry.local:5000 |

*1: If you want to map multiple port, enter the port information separated by commas. `--port 8080:80:tcp,2222:22:tcp`  
*2: If you want to mount multiple directory, enter the directories separated by commas. `--mount /mnt/data:/data,/home/user:/home`

### `start`
Start container
| Option | Value | Required | Description | Example |
| ------ | ----- | -------- | ----------- | ------- |
| --id | [container_id] | yes (*1)| Specify container ID | --id f62a4eb388bf |
| --name | [container_name] | yes (*1) | Specify container name | --name my_container |
| --it | n/a | no | Enable standard output | --it |

*1: Either `--id` or `--name` must be specified.

### `run`
Run container (`create`+`start`)
| Option | Value | Required | Description | Example |
| ------ | ----- | -------- | ----------- | ------- |
| --image | [repositry]:[tag] | yes | Specify name of a image | --image alpine:latest |
| --name | [container_name] | no | Specify container name | --name my_container |
| --rm | n/a | no | Automatically delete at exit of container | --rm |
| --port | [host_port]:[container_port]:[protocol] | no | Map port `[host_port]` to `[container_port]/[protocol]` in the container (*1) | --port 8080:80:tcp |
| --mount | [host_path]:[container_path] | no | Mount `[host_path]` to `[container_path]` in the container (*2) | --mount /mnt/data:/data |
| --cmd | [arg_1],[arg2],... | no | Override entrypoint command | --cmd sleep,100 |
| --ns | [namespace] | no | Specify the namespace to which the container belongs | --ns sandbox |
| --registry | [registry]:[port] | no | Specify registry | --registry my.registry.local:5000 |

*1: If you want to map multiple port, enter the port information separated by commas. `--port 8080:80:tcp,2222:22:tcp`  
*2: If you want to mount multiple directory, enter the directories separated by commas. `--mount /mnt/data:/data,/home/user:/home`


### `exec`
Execute command in container
| Option | Value | Required | Description | Example |
| ------ | ----- | -------- | ----------- | ------- |
| --id | [container_id] | yes (*1)| Specify container ID | --id f62a4eb388bf |
| --name | [container_name] | yes (*1) | Specify container name | --name my_container |
| --cmd | [arg_1],[arg2],... | yes | Specify command to be executed in the container | --cmd /bin/bash |
| --it | n/a | no | Enable standard output | --it |

*1: Either `--id` or `--name` must be specified.

### `stop`
Stop container
| Option | Value | Required | Description | Example |
| ------ | ----- | -------- | ----------- | ------- |
| --id | [container_id] | yes (*1)| Specify container ID | --id f62a4eb388bf |
| --name | [container_name] | yes (*1) | Specify container name | --name my_container |

*1: Either `--id` or `--name` must be specified.

### `restart`
Restart container (`stop`+`start`)
| Option | Value | Required | Description | Example |
| ------ | ----- | -------- | ----------- | ------- |
| --id | [container_id] | yes (*1)| Specify container ID | --id f62a4eb388bf |
| --name | [container_name] | yes (*1) | Specify container name | --name my_container |

*1: Either `--id` or `--name` must be specified.

### `rm`
Delete container
| Option | Value | Required | Description | Example |
| ------ | ----- | -------- | ----------- | ------- |
| --id | [container_id] | yes (*1)| Specify container ID | --id f62a4eb388bf |
| --name | [container_name] | yes (*1) | Specify container name | --name my_container |

*1: Either `--id` or `--name` must be specified.

### `ls`
Show container list
| Option | Value | Required | Description | Example |
| ------ | ----- | -------- | ----------- | ------- |
| --ns | [namespace] | no | Specify the namespace  | --ns sandbox |

### `spec`
Show container spec
| Option | Value | Required | Description | Example |
| ------ | ----- | -------- | ----------- | ------- |
| --id | [container_id] | yes (*1)| Specify container ID | --id f62a4eb388bf |
| --name | [container_name] | yes (*1) | Specify container name | --name my_container |

*1: Either `--id` or `--name` must be specified.

### `logs`
Show container logs
| Option | Value | Required | Description | Example |
| ------ | ----- | -------- | ----------- | ------- |
| --id | [container_id] | yes (*1)| Specify container ID | --id f62a4eb388bf |
| --name | [container_name] | yes (*1) | Specify container name | --name my_container |

*1: Either `--id` or `--name` must be specified.


## Image operations
The operations related to image are as follows:
| Operation | Description |
| --------- | ----------- |
| [`images`](#images) | Show images |
| [`pull`](#pull) | Pull image |
| [`push`](#push) | Push image |
| [`rmi`](#pull) | Delete image |
| [`build`](#build) | Build image |

### `images`
Show images
* No options available

### `pull`
Pull image
| Option | Value | Required | Description | Example |
| ------ | ----- | -------- | ----------- | ------- |
| --image | [repositry]:[tag] | yes | Specify name of a image | --image alpine:latest |
| --os | [os]:[arch] | no | Specify image's os/archtecture | --os linux:amd64 |
| --registry | [registry]:[port] | no | Specify registry | --registry my.registry.local:5000 |

### `push`
Push image
| Option | Value | Required | Description | Example |
| ------ | ----- | -------- | ----------- | ------- |
| --image | [repositry]:[tag] | yes | Specify name of a image | --image alpine:latest |
| --registry | [registry]:[port] | no (*1) | Specify registry | --registry my.registry.local:5000 |

(*1) If you have already connected to the registry with this command, this option is not necessary.


### `rmi`
Delete image
| Option | Value | Required | Description | Example |
| ------ | ----- | -------- | ----------- | ------- |
| --id | [container_id] | yes | Specify image ID | --id b0c9d60fc5e3 |

### `build`
Build image
| Option | Value | Required | Description | Example |
| ------ | ----- | -------- | ----------- | ------- |
| --name | [repositry]:[tag] | yes | Specify name of a image | --name my_app:local |


## Namespace operations
The operations related to Namespace are as follows:
| Operation | Description |
| --------- | ----------- |
| [`ns`](#ns) | Show namespace list |
| [`createns`](#createns) | Create namespace |
| [`rmns`](#rmns) | Delete namespace |

### `ns`
Show namespace list
* No options available

### `createns`
Create namespace
| Option | Value | Required | Description | Example |
| ------ | ----- | -------- | ----------- | ------- |
| --name | [namespace] | yes | Specify name of namespace | --name sandbox |

### `rmns`
Delete namespace
| Option | Value | Required | Description | Example |
| ------ | ----- | -------- | ----------- | ------- |
| --name | [namespace] | yes | Specify name of namespace | --name sandbox |

## Module operations
The operations related to Module are as follows:
| Operation | Description |
| --------- | ----------- |
| [`mod`](#mod) | Show module list |
| [`mod enable`](#mod-enable) | Enable module |
| [`mod disable`](#mod-disable) | Disable module |

### `mod`
Show module list
* No options available

### `mod enable`
Enable module
| Option | Value | Required | Description | Example |
| ------ | ----- | -------- | ----------- | ------- |
| --name | [module_name] | yes | module name | --name dns |

### `mod disable`
Disable module
| Option | Value | Required | Description | Example |
| ------ | ----- | -------- | ----------- | ------- |
| --name | [module_name] | yes | module name | --name dns |

## Node Controller operation
The operations related to node controller are as follows:
| Operation | Description |
| --------- | ----------- |
| [`nodectl info`](#nodectl-info) | Show node info |
| [`nodectl connect`](#nodectl-connect) | Connect node |
| [`nodectl disconnect`](#nodectl-disconnect) | Disconnect node |
| [`nodectl mode stand-alone`](#nodectl-mode-stand-alone) | Change mode to stand-alone |
| [`nodectl mode remote-control`](#nodectl-mode-remote-control) | Change mode to remote-control |

### `nodectl info`
Show node info
* No options available

### `nodectl connect`
Connect node
| Option | Value | Required | Description | Example |
| ------ | ----- | -------- | ----------- | ------- |
| --node | [node] | yes | node ip/fqdn | --node 192.168.1.1 |
| --auth | [auth_code] | no | auth code | --auth 0123456789abcdef |

### `nodectl disconnect`
Disconnect node
* No options available

### `nodectl mode stand-alone`
Change mode to stand-alone
* No options available

### `nodectl mode remote-control`
Change mode to remote-control
* No options available

## Registry Controller operation
The operations related to registry controller are as follows:
| Operation | Description |
| --------- | ----------- |
| [`regctl info`](#regctl-info) | Show registry info |
| [`regctl connect`](#regctl-connect) | Connect registry |
| [`regctl disconnect`](#regctl-disconnect) | Disconnect registry |
| [`regctl get repository`](#regctl-get-repository) | Get repository list |
| [`regctl get tag`](#regctl-get-tag) | Get tag list |
| [`regctl delete`](#regctl-delete) | Delete image |

### `regctl info`
Show registry info
* No options available

### `regctl connect`
Connect registry
| Option | Value | Required | Description | Example |
| ------ | ----- | -------- | ----------- | ------- |
| --registry | [registry] | yes | registry ip/fqdn | --registry 192.168.1.10:5000 |

### `regctl disconnect`
Disconnect registry
* No options available

### `regctl get repository`
Get repository list
* No options available

### `regctl get tag`
Get tag list
| Option | Value | Required | Description | Example |
| ------ | ----- | -------- | ----------- | ------- |
| --repository | [repository] | yes | repository name | --repository ubuntu |

### `regctl delete`
Delete image
| Option | Value | Required | Description | Example |
| ------ | ----- | -------- | ----------- | ------- |
| --image | [repositry]:[tag] | yes | Specify name of a image | --image alpine:latest |
