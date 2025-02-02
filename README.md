# Karakuri
## Introduction
`karakuri` is a platform for managing and running containers.  
This platform consists of three components.
1. `karakuri`  
   A CLI tool that provides the ability to send requests to high-level container runtime.
1. `hitoha`  
   High-level container runtime.  
   Runs as a daemon process and performs container networking, lifecycle management and image management.  
   `hitoha` provides a REST API as an interface and executes low-level container runtime in response to received requests.
1. `futaba`  
   Low-level container runtime.  
   `futaba` provides actual container operation, including namespace isolation, mounts, root filesystem changes, etc.

## Building
### Pre-requisites
Before building `karakuri`, must install `cgroup-tools`.  
On Ubuntu/Debian:  
```
$ apt update && apt install -y cgroup-tools
```

### Build
```
$ git clone https://github.com/pyxgun/karakuri.git
$ cd karakuri

$ sh scripts/build.sh
```

### Running the test
```
$ sudo karakuri run --rm --it --image=hello-world

Hello from Docker!
This message shows that your installation appears to be working correctly.
  :
```

## Quick start
The following is how to run a new container and delete container.  
### Pull image (Optional)
Retrieve an image from registry.
This operation is optional because if an image that doesn't exist local is specified when `create` command is executed, the image is automatically retrieved.
```
$ sudo karakuri pull --image=alpine
```
Check the local image list.
```
$ sudo karakuri images

REPOSITORY      | TAG           | ID
----------------+---------------+-----------------
alpine          | latest        | b0c9d60fc5e3
```

### Create container
Create a container.
```
$ sudo karakuri create --name=mycontainer --image=alpine
```
Check container list.
```
$ sudo karakuri ls

CONTAINER ID | Name        | IMAGE  | STATUS  | PORT | COMMAND
-------------+-------------+--------+---------+------+---------------
52e579b39313 | mycontainer | alpine | created |      | /bin/sh                         
```

### Start container
Start the container.
```
$ sudo karakuri start --name=mycontainer --it

/ # 
```
`--it` option is specified when executing commands that require standard output, such as shell and other interactive tools.

### Delete container
Delete the container
```
$ sudo karakuri rm --name=mycontainer
```

## More documentation
* [Command List](./docs/command_list.md)

## License
The code and docs are released under the [MIT license](LICENSE).