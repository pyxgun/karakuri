package hitoha

import (
	"karakuri_mod"
	"karakuripkgs"
	"strings"
)

type ParamsCreateContainer struct {
	ImageInfo string
	Name      string
	Namespace string
	Port      string
	Mount     string
	Cmd       string
	Repositry string
	Restart   string
}

type ParamsRunContainer struct {
	ImageInfo string
	Name      string
	Namespace string
	Port      string
	Mount     string
	Cmd       string
	Repositry string
	Restart   string
}

func CreateContainer(params ParamsCreateContainer) ResponseContainerInfo {
	// check if namespace exists
	if !isNamespaceExist(params.Namespace) {
		return createResponseContainerInfo("error", ContainerInfo{}, "Namespace: \""+params.Namespace+"\" is not exists.")
	}

	// check if name is already used
	if isContainerNameExists(params.Name) {
		return createResponseContainerInfo("error", ContainerInfo{}, "Container name: \""+params.Name+"\" is already used.")
	}

	// parse image tag
	image_info := strings.Split(params.ImageInfo, ":")
	image := image_info[0]
	tag := "latest"
	if len(image_info) == 2 {
		tag = image_info[1]
	}

	if !isImageExists(image, tag) {
		PullImage(params.ImageInfo, "linux:amd64", params.Repositry)
	}

	// retrieve rootfs
	rootfs := getImageRootfs(image, tag)

	// retrieve blob file
	image_id := GetImageId(image, tag)
	blob_file := readBlobFile(image_id)

	// retrieve command from blob
	command := params.Cmd
	if command == "none" {
		var new_command string = ""
		entrypoints := blob_file.Config.Entrypoint
		cmds := blob_file.Config.Cmd

		for _, entry := range entrypoints {
			new_command += entry + ","
		}

		for _, entry := range cmds {
			new_command += entry + ","
		}
		command = strings.TrimRight(new_command, ",")
	}

	// retrieve environment vars from blob
	var envs string = ""
	for _, entry := range blob_file.Config.Env {
		new_env := strings.Replace(entry, " ", "*", -1)
		new_env = strings.Replace(new_env, "\t", "", -1)
		envs += new_env + "!"
	}
	envs = strings.TrimRight(envs, "!")

	// lease address
	address, res := assignNewAddress(params.Namespace)
	if !res {
		return createResponseContainerInfo("error", ContainerInfo{}, "no ip address available.")
	}

	namespace_info := showNamespace(params.Namespace)
	// hostdevice
	hostdevice := namespace_info.Network.Name
	// gateway
	gateway := namespace_info.Network.Address
	// nameserver
	nameserver := "none"
	// if DNS module is enabled, set DNS module address to container nameserver
	if karakuri_mod.IsModuleEnabled("dns") {
		dns_module_id := karakuripkgs.RetrieveContainerId("", "dns")
		dns_spec := karakuripkgs.ReadSpecFile(karakuripkgs.FUTABA_ROOT + "/" + dns_module_id)
		nameserver = (strings.Split(dns_spec.Network.Address, "/"))[0]
	}

	// execute runtime: spec
	//karakuripkgs.RuntimeSpec(rootfs, port, mount, address, cmd)
	karakuripkgs.RuntimeSpec(karakuripkgs.ParamsRuntimeSpec{
		ImagePath:  rootfs,
		Port:       params.Port,
		Mount:      params.Mount,
		HostDevice: hostdevice,
		Address:    address,
		Gateway:    gateway,
		Nameserver: nameserver,
		Command:    command,
		EnvVars:    envs,
		Restart:    params.Restart,
	})

	// add container list
	config_spec := karakuripkgs.ReadSpecFile(".")

	// bind address to container
	bindAddressToContainerId(params.Namespace, config_spec.Hostname, address)

	// add new container
	container_info := addNewContainer(config_spec, image, params.Name, params.Namespace)

	// add container to namespace
	addContainerToNamespace(params.Namespace, config_spec.Hostname)

	// if DNS module is enabled, add container entry to dns
	if karakuri_mod.IsModuleEnabled("dns") {
		karakuri_mod.AddDnsRecord(config_spec.Hostname, config_spec.Network.Address)
	}

	// execute runtime: create
	if err := karakuripkgs.RuntimeCreate(); err != nil {
		return createResponseContainerInfo("error", container_info, "failed to create container.")
	}
	return createResponseContainerInfo("success", container_info, "container create success.")
}

func StartContainer(id string) ResponseContainerInfo {
	// check container status
	container_status := checkContainerStatus(id)
	if container_status == "created" || container_status == "stopped" {
		// update status: running
		container_info := UpdateContainerStatus(id, "running")
		return createResponseContainerInfo("success", container_info, "container: "+id+" start success.")
	} else {
		return createResponseContainerInfo("error", ContainerInfo{}, "container: "+id+" is already up and running.")
	}
}

func ExecContainer(id string) ResponseContainerInfo {
	container_status := checkContainerStatus(id)
	if container_status == "running" {
		return createResponseContainerInfo("success", ContainerInfo{}, "container exec success.")
	} else {
		return createResponseContainerInfo("error", ContainerInfo{}, "container: "+id+" is not running.")
	}
}

func RunContainer(params ParamsRunContainer) ResponseRunContainer {
	//resp := CreateContainer(image_tag, port, mount, cmd, repositry)
	resp := CreateContainer(ParamsCreateContainer{
		ImageInfo: params.ImageInfo,
		Port:      params.Port,
		Mount:     params.Port,
		Cmd:       params.Cmd,
		Repositry: params.Repositry,
		Name:      params.Name,
		Namespace: params.Namespace,
	})
	if resp.Result != "success" {
		return createResponseRunContainer(resp.Result, "", resp.Message)
	}
	id := resp.Container.Id

	resp = StartContainer(id)

	return createResponseRunContainer("success", id, resp.Message)
}

func KillContainer(id string) ResponseStopContainer {
	container_status := checkContainerStatus(id)
	if container_status == "running" {
		// execute runtime: kill
		karakuripkgs.RuntimeKill(id)

		// update status
		UpdateContainerStatus(id, "stopped")

		return createResponseStopContainer("success", id, "container: "+id+" stop success.")
	} else {
		return createResponseStopContainer("error", id, "container: "+id+" is not running.")
	}
}

func DeleteContainer(id string) ResponseDeleteContainer {
	container_status := checkContainerStatus(id)
	if container_status == "created" || container_status == "stopped" {
		// execute runtime: ddelete
		karakuripkgs.RuntimeDelete(id)

		// delete from container list
		deleteContainerList(id)

		// delete from namespace
		deleteContainerFromNamespace(id)

		// free address
		freeAddress(id)

		// if DNS module is enabled, delete entry
		if karakuri_mod.IsModuleEnabled("dns") {
			karakuri_mod.DeleteDnsRecord(id)
		}

		return createResponseDeleteContainer("success", id, "container: "+id+" delete success.")
	} else {
		return createResponseDeleteContainer("error", id, "container: "+id+" is still running.")
	}
}
