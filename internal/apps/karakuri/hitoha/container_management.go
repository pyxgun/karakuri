package hitoha

import (
	"encoding/json"
	"karakuripkgs"
	"os"
)

// ----------------
// container info
type PortInfo struct {
	HostPort      int    `json:"host_port"`
	ContainerPort int    `json:"container_port"`
	Protocol      string `json:"protocol"`
}

type ContainerInfo struct {
	Id        string     `json:"container_id"`
	Name      string     `json:"name"`
	Namespace string     `json:"namespace"`
	Image     string     `json:"image"`
	Command   string     `json:"command"`
	Status    string     `json:"status"`
	Port      []PortInfo `json:"port"`
	Pid       int        `json:"pid"`
	Restart   string     `json:"restart"`
}

type ContainerList struct {
	List []ContainerInfo `json:"list"`
}

// ----------------

func newContainerList() {
	var container_list_data ContainerList
	data, _ := json.MarshalIndent(container_list_data, "", "  ")
	if err := os.WriteFile(karakuripkgs.HITOHA_CONTAINER_LIST, data, os.ModePerm); err != nil {
		panic(err)
	}
}

func addNewContainer(config_spec karakuripkgs.ConfigSpec, image string, name string, namespace string) ContainerInfo {
	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.HITOHA_CONTAINER_LIST)
	if err != nil {
		newContainerList()
		bytes, _ = os.ReadFile(karakuripkgs.HITOHA_CONTAINER_LIST)
	}

	var container_list_data ContainerList
	if err := json.Unmarshal(bytes, &container_list_data); err != nil {
		panic(err)
	}

	// command
	var command string = ""
	for _, cmd := range config_spec.Process.Args {
		command += cmd + " "
	}

	// port
	var ports []PortInfo
	for _, entry := range config_spec.Network.Port {
		ports = append(ports, PortInfo{
			HostPort:      entry.HostPort,
			ContainerPort: entry.TargetPort,
			Protocol:      entry.Protocol,
		})
	}

	// set container info
	var container_info = ContainerInfo{
		Id:        config_spec.Hostname,
		Name:      name,
		Namespace: namespace,
		Image:     image,
		Command:   command,
		Status:    "created",
		Port:      ports,
		Pid:       config_spec.Process.Pid,
		Restart:   config_spec.Restart,
	}
	container_list_data.List = append(container_list_data.List, container_info)

	data, _ := json.MarshalIndent(container_list_data, "", "  ")
	if err := os.WriteFile(karakuripkgs.HITOHA_CONTAINER_LIST, data, os.ModePerm); err != nil {
		panic(err)
	}

	return container_info
}

func UpdateContainerStatus(id string, status string) ContainerInfo {
	// read spec file
	config_spec := karakuripkgs.ReadSpecFile(karakuripkgs.FUTABA_ROOT + "/" + id)

	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.HITOHA_CONTAINER_LIST)
	if err != nil {
		panic(err)
	}

	var (
		container_list_data     ContainerList
		new_container_list_data ContainerList
		new_container_info      ContainerInfo
	)
	if err := json.Unmarshal(bytes, &container_list_data); err != nil {
		panic(err)
	}
	for _, entry := range container_list_data.List {
		if entry.Id == id {
			new_container_info = ContainerInfo{
				Id:        entry.Id,
				Name:      entry.Name,
				Namespace: entry.Namespace,
				Image:     entry.Image,
				Command:   entry.Command,
				Status:    status,
				Port:      entry.Port,
				Pid:       config_spec.Process.Pid,
				Restart:   entry.Restart,
			}
			new_container_list_data.List = append(new_container_list_data.List, new_container_info)
		} else {
			new_container_list_data.List = append(new_container_list_data.List, entry)
		}
	}

	data, _ := json.MarshalIndent(new_container_list_data, "", "  ")
	if err := os.WriteFile(karakuripkgs.HITOHA_CONTAINER_LIST, data, os.ModePerm); err != nil {
		panic(err)
	}

	return new_container_info
}

func deleteContainerList(id string) {
	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.HITOHA_CONTAINER_LIST)
	if err != nil {
		panic(err)
	}

	var (
		container_list_data     ContainerList
		new_container_list_data ContainerList
	)
	if err := json.Unmarshal(bytes, &container_list_data); err != nil {
		panic(err)
	}

	for _, entry := range container_list_data.List {
		if entry.Id != id {
			new_container_list_data.List = append(new_container_list_data.List, entry)
		}
	}

	data, _ := json.MarshalIndent(new_container_list_data, "", "  ")
	if err := os.WriteFile(karakuripkgs.HITOHA_CONTAINER_LIST, data, os.ModePerm); err != nil {
		panic(err)
	}
}

func syncContainerList() {
	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.HITOHA_CONTAINER_LIST)
	if err != nil {
		panic(err)
	}

	var (
		container_list_data     ContainerList
		new_container_list_data ContainerList
	)
	if err := json.Unmarshal(bytes, &container_list_data); err != nil {
		panic(err)
	}

	for _, entry := range container_list_data.List {
		config_spec := karakuripkgs.ReadSpecFile(karakuripkgs.FUTABA_ROOT + "/" + entry.Id)
		// command
		cmd_str := ""
		for _, cmd := range config_spec.Process.Args {
			cmd_str += cmd + " "
		}
		new_container_info := ContainerInfo{
			Id:        entry.Id,
			Name:      entry.Name,
			Namespace: entry.Namespace,
			Image:     entry.Image,
			Command:   cmd_str,
			Status:    entry.Status,
			Port:      entry.Port,
			Pid:       entry.Pid,
			Restart:   entry.Restart,
		}
		new_container_list_data.List = append(new_container_list_data.List, new_container_info)
	}

	data, _ := json.MarshalIndent(new_container_list_data, "", "  ")
	if err := os.WriteFile(karakuripkgs.HITOHA_CONTAINER_LIST, data, os.ModePerm); err != nil {
		panic(err)
	}
}

func ShowContainerList(namespace string) ResponseContainerList {
	syncContainerList()

	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.HITOHA_CONTAINER_LIST)
	if err != nil {
		return createResponseContainerList("error", ContainerList{})
	}

	var container_list_data ContainerList
	if err := json.Unmarshal(bytes, &container_list_data); err != nil {
		panic(err)
	}

	var target_container_list ContainerList

	if namespace == "all" {
		target_container_list = container_list_data
	} else {
		// check if namespace exists
		if !isNamespaceExist(namespace) {
			return createResponseContainerList("error", ContainerList{})
		}

		// retrieve target namespace container
		for _, entry := range container_list_data.List {
			if entry.Namespace == namespace {
				target_container_list.List = append(target_container_list.List, entry)
			}
		}
	}

	return createResponseContainerList("success", target_container_list)
}

func ShowContainerSpec(id string) ResponseContainerSpec {
	// retrieve config spec
	config_spec := karakuripkgs.ReadSpecFile(karakuripkgs.FUTABA_ROOT + "/" + id)

	return createResponseSpecContainer("success", config_spec)
}

func isContainerNameExists(name string) bool {
	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.HITOHA_CONTAINER_LIST)
	if err != nil {
		newContainerList()
		bytes, _ = os.ReadFile(karakuripkgs.HITOHA_CONTAINER_LIST)
	}

	var container_list_data ContainerList
	if err := json.Unmarshal(bytes, &container_list_data); err != nil {
		panic(err)
	}

	for _, entry := range container_list_data.List {
		if entry.Name == name {
			return true
		}
	}
	return false
}

func checkContainerStatus(id string) string {
	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.HITOHA_CONTAINER_LIST)
	if err != nil {
		panic(err)
	}

	var container_list_data ContainerList
	if err := json.Unmarshal(bytes, &container_list_data); err != nil {
		panic(err)
	}

	for _, entry := range container_list_data.List {
		if entry.Id == id {
			return entry.Status
		}
	}

	return ""
}

func retrieveContainerId(name string) ResponseContainerId {
	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.HITOHA_CONTAINER_LIST)
	if err != nil {
		panic(err)
	}

	var container_list_data ContainerList
	if err := json.Unmarshal(bytes, &container_list_data); err != nil {
		panic(err)
	}

	for _, entry := range container_list_data.List {
		if entry.Name == name {
			return createResponseContainerId("success", entry.Id, "retrieve container id success")
		}
	}
	return createResponseContainerId("error", "", "no such container, name: "+name)
}

func autoStartContainer() {
	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.HITOHA_CONTAINER_LIST)
	if err != nil {
		panic(err)
	}

	var container_list_data ContainerList
	if err := json.Unmarshal(bytes, &container_list_data); err != nil {
		panic(err)
	}

	for _, entry := range container_list_data.List {
		if entry.Restart == "on-boot" {
			config_spec := karakuripkgs.ReadSpecFile(karakuripkgs.FUTABA_ROOT + "/" + entry.Id)
			karakuripkgs.RuntimeStart(entry.Id, false)
			UpdateContainerStatus(entry.Id, "running")
			if config_spec.Network.Port != nil {
				SetupPortForwarding("add", config_spec.Network)
			}
		}
	}
}
