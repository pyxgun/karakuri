package hitoha

import (
	"encoding/json"
	"karakuripkgs"
	"os"
	"strconv"
)

type NamespaceNetwork struct {
	DevType string `json:"type"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Subnet  string `json:"subnet"`
}

type NamespaceInfo struct {
	Name      string           `json:"name"`
	Network   NamespaceNetwork `json:"network"`
	Container []string         `json:"container"`
}

type NamespaceList struct {
	Namespaces []NamespaceInfo `json:"namespace"`
}

func createNamespaceList() {
	namespace_list_data := NamespaceList{
		Namespaces: []NamespaceInfo{
			{
				Name: "system",
				Network: NamespaceNetwork{
					DevType: "bridge",
					Name:    "karakuri0",
					Address: "10.157.0.1",
					Subnet:  "/24",
				},
				Container: nil,
			},
			{
				Name: "system-mod",
				Network: NamespaceNetwork{
					DevType: "bridge",
					Name:    "karakuri1",
					Address: "10.157.1.1",
					Subnet:  "/24",
				},
				Container: nil,
			},
			{
				Name: "default",
				Network: NamespaceNetwork{
					DevType: "bridge",
					Name:    "karakuri2",
					Address: "10.157.2.1",
					Subnet:  "/24",
				},
				Container: nil,
			},
		},
	}
	data, _ := json.MarshalIndent(namespace_list_data, "", "  ")
	if err := os.WriteFile(karakuripkgs.HITOHA_NAMESPACE_LIST, data, os.ModePerm); err != nil {
		panic(err)
	}
}

func isNamespaceExist(namespace string) bool {
	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.HITOHA_NAMESPACE_LIST)
	if err != nil {
		panic(err)
	}

	var namespace_list_data NamespaceList
	if err := json.Unmarshal(bytes, &namespace_list_data); err != nil {
		panic(err)
	}

	for _, entry := range namespace_list_data.Namespaces {
		if entry.Name == namespace {
			return true
		}
	}
	return false
}

func isNamespaceHasContainer(namespace_info NamespaceInfo) bool {
	return len(namespace_info.Container) != 0
}

func showNamespace(namespace string) NamespaceInfo {
	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.HITOHA_NAMESPACE_LIST)
	if err != nil {
		createNamespaceList()
		bytes, _ = os.ReadFile(karakuripkgs.HITOHA_NAMESPACE_LIST)
	}

	var namespace_list_data NamespaceList
	if err := json.Unmarshal(bytes, &namespace_list_data); err != nil {
		panic(err)
	}

	for _, entry := range namespace_list_data.Namespaces {
		if entry.Name == namespace {
			return entry
		}
	}

	return NamespaceInfo{}
}

func showNamespaceList() ResponseNamespaceList {
	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.HITOHA_NAMESPACE_LIST)
	if err != nil {
		createNamespaceList()
		bytes, _ = os.ReadFile(karakuripkgs.HITOHA_NAMESPACE_LIST)
	}

	var namespace_list_data NamespaceList
	if err := json.Unmarshal(bytes, &namespace_list_data); err != nil {
		panic(err)
	}

	return createResponseNamespaceList("success", namespace_list_data)
}

func createNewNamespace(namespace string) ResponseCreateNamespace {
	// check if namespace is reserved name
	if namespace == "system" || namespace == "system-mod" || namespace == "default" || namespace == "all" {
		return createResponseCreateNamespace("error", "Namespace: \""+namespace+"\" can not create because that namespace is reserved for system usage.")
	}
	// check if namespace is exist
	if isNamespaceExist(namespace) {
		return createResponseCreateNamespace("error", "Namespace: \""+namespace+"\" is already exists.")
	}

	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.HITOHA_NAMESPACE_LIST)
	if err != nil {
		panic(err)
	}

	var namespace_list_data NamespaceList
	if err := json.Unmarshal(bytes, &namespace_list_data); err != nil {
		panic(err)
	}

	// network
	num_namespace := len(namespace_list_data.Namespaces)
	network_index := strconv.Itoa(num_namespace)
	network := NamespaceNetwork{
		DevType: "bridge",
		Name:    "karakuri" + network_index,
		Address: "10.157." + network_index + ".1",
		Subnet:  "/24",
	}

	// create network interface
	createNetworkInterface(network.DevType, network.Name, network.Address, network.Subnet)

	// allow namespace traffic
	allowContainerTrafficRule(network.Name)

	// allow system-mod traffic
	allowContainerModTrafficRule(network.Name)

	// create namespace to networklist
	createNamespaceToNetworkList(namespace)

	namespace_list_data.Namespaces = append(
		namespace_list_data.Namespaces,
		NamespaceInfo{
			Name:      namespace,
			Network:   network,
			Container: nil,
		},
	)

	data, _ := json.MarshalIndent(namespace_list_data, "", "  ")
	if err := os.WriteFile(karakuripkgs.HITOHA_NAMESPACE_LIST, data, os.ModePerm); err != nil {
		panic(err)
	}

	return createResponseCreateNamespace("success", "Namespace: \""+namespace+"\" create success.")
}

func deleteNamespace(target_namespace string) ResponseDeleteNamespace {
	// check if namespace is reserved name
	if target_namespace == "system" || target_namespace == "system-mod" || target_namespace == "default" {
		return createResponseDeleteNamespace("error", "Namespace: \""+target_namespace+"\" can not delete because that namespace is reserved for system usage.")
	}
	// check if namespace exists
	if !isNamespaceExist(target_namespace) {
		return createResponseDeleteNamespace("error", "Namespace: \""+target_namespace+"\" is not exists.")
	}

	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.HITOHA_NAMESPACE_LIST)
	if err != nil {
		panic(err)
	}

	var namespace_list_data NamespaceList
	if err := json.Unmarshal(bytes, &namespace_list_data); err != nil {
		panic(err)
	}

	var new_namespace_list NamespaceList
	for _, entry := range namespace_list_data.Namespaces {
		if entry.Name == target_namespace {
			// check if namespace has container
			if isNamespaceHasContainer(entry) {
				return createResponseDeleteNamespace("error", "Namespace: \""+target_namespace+"\" can not delete because container(s) are still attached to this namespace.")
			}
			// delete network interface
			deleteNetworkInterface(entry.Network.Name)
			// delete namespace traffic rule
			deleteContainerTrafficRule(entry.Network.Name)
			// delete system-mod traffic rule
			deleteContainerModTrafficRule(entry.Network.Name)
		} else {
			new_namespace_list.Namespaces = append(new_namespace_list.Namespaces, entry)
		}
	}

	// delete namespace from network list
	deleteNamespaceFromNetworkList(target_namespace)

	data, _ := json.MarshalIndent(new_namespace_list, "", "  ")
	if err := os.WriteFile(karakuripkgs.HITOHA_NAMESPACE_LIST, data, os.ModePerm); err != nil {
		panic(err)
	}

	return createResponseDeleteNamespace("success", "Namespace: "+target_namespace+" delete success.")
}

func addContainerToNamespace(namespace string, id string) {
	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.HITOHA_NAMESPACE_LIST)
	if err != nil {
		panic(err)
	}

	var namespace_list_data NamespaceList
	if err := json.Unmarshal(bytes, &namespace_list_data); err != nil {
		panic(err)
	}

	var new_namespace_list NamespaceList
	for _, entry := range namespace_list_data.Namespaces {
		if entry.Name == namespace {
			entry.Container = append(entry.Container, id)
		}
		new_namespace_list.Namespaces = append(new_namespace_list.Namespaces, entry)
	}

	data, _ := json.MarshalIndent(new_namespace_list, "", "  ")
	if err := os.WriteFile(karakuripkgs.HITOHA_NAMESPACE_LIST, data, os.ModePerm); err != nil {
		panic(err)
	}
}

func deleteContainerFromNamespace(target_id string) {
	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.HITOHA_NAMESPACE_LIST)
	if err != nil {
		panic(err)
	}

	var namespace_list_data NamespaceList
	if err := json.Unmarshal(bytes, &namespace_list_data); err != nil {
		panic(err)
	}

	var new_namespace_list NamespaceList
	for _, entry := range namespace_list_data.Namespaces {
		var tmp NamespaceInfo
		tmp.Name = entry.Name
		tmp.Network = entry.Network
		for _, id := range entry.Container {
			if id != target_id {
				tmp.Container = append(tmp.Container, id)
			}
		}
		new_namespace_list.Namespaces = append(new_namespace_list.Namespaces, tmp)
	}

	data, _ := json.MarshalIndent(new_namespace_list, "", "  ")
	if err := os.WriteFile(karakuripkgs.HITOHA_NAMESPACE_LIST, data, os.ModePerm); err != nil {
		panic(err)
	}
}
