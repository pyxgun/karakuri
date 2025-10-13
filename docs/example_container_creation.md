# Example Container Creation
`karakuri create` command to create a container has various options. Here are some use cases and the commands to do so.

## Allow external access
For applications that are accessed not only locally but also from external networks, such as web applications,  
external access can be enabled by specifying the following options.  

- `--port [host_port]:[container_port]:[protocol(tcp|udp)]`

To forward TCP traffic to port 8080 of the host to port 80 of the container, create a container using the following command.
```
$ sudo karakuri create --name web-app --image nginx:alpine --port 8080:80:tcp
```

Multiple ports can be forwarded by entering comma-separated entries.  
```
$ sudo karakuri create --name web-app --image nginx:alpine --port 8080:80:tcp,4443:443:tcp
```

## Volume Persistence
Normally, files created in a container are deleted along with the container when it is deleted.  
If you need to persist data such as storage services or databases so that only the data is not deleted when the container is deleted, use the following options.

- `--mount [host_path]:[container_path]`

To mount the host's /mnt/data directory to the container's /data directory, create a container using the following command.
```
$ sudo karakuri create --name web-app --image myapp --mount /mnt/data:/data
```

Multiple directory can be mounted by entering comma-separated entries.  
```
$ sudo karakuri create --name web-app --image nginx:alpine --mount /mnt/data:data,/mnt/logs:/log
```

## Automatic container start
When the host device is rebooted for maintenance, etc., the container will not start by default and must be started manually.  
If you want to automatically start the container when the host device is rebooted, use the following options.

- `--restart on-boot`

For a container to be started automatically, create a container using the following command.
```
$ sudo karakuri create --name web-app --image nginx:alpine --restart on-boot
```

## Pull image from private registry
If you do not have an image locally, the default is to pull an image from Docker Hub.
If you want to pull images from a private registry instead of DockerHub, use the following options.

- `--registry [registry_address]`

To pull an image from a private registry listening on 192.168.1.1:5000, create a container using the following command.
```
$ sudo karakuri create --name my-app --image my-app:priv --registry 192.168.1.1:5000
```

## Override entrypoint
The entry point at container startup is set by default to the entry point listed in the image.  
If you wish to change the entry point for testing or for any reason, use the following options.

- `--cmd [command]`

In images where `/bin/sh` is invoked by default, to invoke it with `cat /etc/logs`, execute the following command.
```
$ sudo karakuri create --name web-app --image alpine --cmd cat,/etc/logs
```
please enter comma separated command.

## 
By combining the options introduced, you can create a container that fits your use case.

1. Allow external access
2. Pull image from private registry
3. Application database persistence
4. Automatic container startup at host device startup

If you want to create a container like the one above, the command would be the following
```
$ sudo karakuri create \
--name my-app \
--image my-app:priv \
--registry 192.168.1.1:5000 \
--port 8080:80:tcp,4443:443:tcp \
--mount /mnt/data:/data,/mnt/logs:/logs
--restart on-boot
```
