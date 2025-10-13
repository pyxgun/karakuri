# Registry Controller
If you have enabled `registry` module or have your own private registry, you can connect to it.  
You can check the repository and tags of the connected registry, push/pull images, and delete images.  
  
It is not connected to any registry by default. The following commands can be used to connect to the registry.
```
$ sudo karakuri regctl connect --registry [REGISTRY]
```
The registry to be connected and the connection status can be checked with the following command.
```
$ sudo karakuri regctl target

Registry : 172.17.20.150:5000
Status   : connected
```

## Get Repositories/Tags list
Get repositories:
```
$ sudo karakuri regctl get repository

REPOSITORY
--------------------------
alpine
ubuntu
nginx
```
Get tags:
```
$ sudo karakuri regctl get tag --repository alpine

REPOSITORY: alpine
TAG
--------------
3.20
3.21
latest
```

## Push Image
If you are connected to a registry, `karakuri push` command will automatically push the image to the connected registry.  
```
$ karakuri push --image ubuntu:24.04

Pushing image, ubuntu:24.04 ...
Push completed
```
(Optional) If you are not connected to a registry, you must specify a registry to push with `--registry [REGISTRY]` option.

## Delete Image
```
$ sudo karakuri regctl delete --image ubuntu:24.04

delete alpine:3.20 success.
```

## Change the connection registry
Before changing the registry to which you are connecting, you must disconnect from the currently connected registry with the following command:
```
$ sudo karakuri regctl disconnect

registry dissconnected.

$ sudo karakuri regctl info
Registry : 172.17.20.150:5000
Status   : disconnected
```
Then connect to new registry with the following command:
```
$ sudo karakuri regctl connect --registry [NEW_REGISTRY]
```
