package hitoha

import (
	"encoding/json"
	"karakuripkgs"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type ContainerNetwork struct {
	ContainerId string `json:"container_id"`
	Address     string `json:"address"`
}

type NetworkInfo struct {
	Namespace string             `json:"namespace"`
	Network   []ContainerNetwork `json:"container"`
}

type NetworkList struct {
	List []NetworkInfo `json:"list"`
}

// ----------
// interface and iptables
func createNetworkInterface(dev_type string, dev_name string, address string, subnet string) {
	// wait for network setup
	time.Sleep(500 * time.Millisecond)
	// create bridge interface
	cmd1 := exec.Command("ip", "link", "add", dev_name, "type", dev_type)
	if err := cmd1.Run(); err != nil {
		panic(err)
	}
	// set address
	cmd2 := exec.Command("ip", "address", "add", address+subnet, "dev", dev_name)
	if err := cmd2.Run(); err != nil {
		panic(err)
	}
	// link up
	cmd3 := exec.Command("ip", "link", "set", "up", dev_name)
	if err := cmd3.Run(); err != nil {
		panic(err)
	}
}

func deleteNetworkInterface(dev_name string) {
	cmd := exec.Command("ip", "link", "del", dev_name)
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

func allowContainerModTrafficRule(dev_name string) {
	// accept input to bridge
	cmd2 := exec.Command("iptables", "-A", "FORWARD", "-i", "karakuri1", "-o", dev_name, "-j", "ACCEPT")
	if err := cmd2.Run(); err != nil {
		panic(err)
	}
	// accept output to bridge
	cmd3 := exec.Command("iptables", "-A", "FORWARD", "-o", "karakuri1", "-i", dev_name, "-j", "ACCEPT")
	if err := cmd3.Run(); err != nil {
		panic(err)
	}
}

func deleteContainerModTrafficRule(dev_name string) {
	// accept input to bridge
	cmd2 := exec.Command("iptables", "-D", "FORWARD", "-i", "karakuri1", "-o", dev_name, "-j", "ACCEPT")
	if err := cmd2.Run(); err != nil {
		panic(err)
	}
	// accept output to bridge
	cmd3 := exec.Command("iptables", "-D", "FORWARD", "-o", "karakuri1", "-i", dev_name, "-j", "ACCEPT")
	if err := cmd3.Run(); err != nil {
		panic(err)
	}
}

func allowContainerTrafficRule(dev_name string) {
	// accept input to bridge
	cmd2 := exec.Command("iptables", "-A", "FORWARD", "-i", karakuripkgs.HOST_NIC, "-o", dev_name, "-j", "ACCEPT")
	if err := cmd2.Run(); err != nil {
		panic(err)
	}
	// accept output to bridge
	cmd3 := exec.Command("iptables", "-A", "FORWARD", "-o", karakuripkgs.HOST_NIC, "-i", dev_name, "-j", "ACCEPT")
	if err := cmd3.Run(); err != nil {
		panic(err)
	}
}

func deleteContainerTrafficRule(dev_name string) {
	// delete rule input to bridge
	cmd2 := exec.Command("iptables", "-D", "FORWARD", "-i", karakuripkgs.HOST_NIC, "-o", dev_name, "-j", "ACCEPT")
	if err := cmd2.Run(); err != nil {
		panic(err)
	}
	// delete rule output to bridge
	cmd3 := exec.Command("iptables", "-D", "FORWARD", "-o", karakuripkgs.HOST_NIC, "-i", dev_name, "-j", "ACCEPT")
	if err := cmd3.Run(); err != nil {
		panic(err)
	}
}

// ----------

func createNewNetworkList() {
	var network_list_data NetworkList

	// namespace: system
	network_list_data.List = append(network_list_data.List,
		NetworkInfo{
			Namespace: "system",
			Network:   nil,
		},
	)
	// namespace: system-mod
	network_list_data.List = append(network_list_data.List,
		NetworkInfo{
			Namespace: "system-mod",
			Network:   nil,
		},
	)
	// namespace: default
	network_list_data.List = append(network_list_data.List,
		NetworkInfo{
			Namespace: "default",
			Network:   nil,
		},
	)

	data, _ := json.MarshalIndent(network_list_data, "", "  ")
	if err := os.WriteFile(karakuripkgs.HITOHA_NETWORK_LIST, data, os.ModePerm); err != nil {
		panic(err)
	}
}

func createNamespaceToNetworkList(namespace string) {
	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.HITOHA_NETWORK_LIST)
	if err != nil {
		createNewNetworkList()
		bytes, _ = os.ReadFile(karakuripkgs.HITOHA_NETWORK_LIST)
	}

	var network_list_data NetworkList
	if err := json.Unmarshal(bytes, &network_list_data); err != nil {
		panic(err)
	}

	network_list_data.List = append(network_list_data.List,
		NetworkInfo{
			Namespace: namespace,
			Network:   nil,
		},
	)

	data, _ := json.MarshalIndent(network_list_data, "", "  ")
	if err := os.WriteFile(karakuripkgs.HITOHA_NETWORK_LIST, data, os.ModePerm); err != nil {
		panic(err)
	}
}

func deleteNamespaceFromNetworkList(namespace string) {
	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.HITOHA_NETWORK_LIST)
	if err != nil {
		createNewNetworkList()
		bytes, _ = os.ReadFile(karakuripkgs.HITOHA_NETWORK_LIST)
	}

	var network_list_data NetworkList
	if err := json.Unmarshal(bytes, &network_list_data); err != nil {
		panic(err)
	}

	var new_network_list_data NetworkList
	for _, entry := range network_list_data.List {
		if entry.Namespace != namespace {
			new_network_list_data.List = append(new_network_list_data.List, entry)
		}
	}

	data, _ := json.MarshalIndent(new_network_list_data, "", "  ")
	if err := os.WriteFile(karakuripkgs.HITOHA_NETWORK_LIST, data, os.ModePerm); err != nil {
		panic(err)
	}
}

func leaseNewAddress(used_address string) (address string, result bool) {
	address_info := strings.Split(used_address, "/")
	network := address_info[0]
	octets := strings.Split(network, ".")
	subnet := address_info[1]

	octet_4, _ := strconv.Atoi(octets[3])
	if octet_4 == 254 {
		return "", false
	}
	new_octet_4 := strconv.Itoa(octet_4 + 1)

	new_address := octets[0] + "." + octets[1] + "." + octets[2] + "." + new_octet_4 + "/" + subnet

	return new_address, true
}

func assignNewAddress(namespace string) (address string, result bool) {
	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.HITOHA_NETWORK_LIST)
	if err != nil {
		createNewNetworkList()
		bytes, _ = os.ReadFile(karakuripkgs.HITOHA_NETWORK_LIST)
	}

	var network_list_data NetworkList
	if err := json.Unmarshal(bytes, &network_list_data); err != nil {
		panic(err)
	}

	// retrieve namespace's base address
	namespace_info := showNamespace(namespace)
	base_address := namespace_info.Network.Address
	subnet := namespace_info.Network.Subnet

	var lease_address string
	for i, entry := range network_list_data.List {
		if entry.Namespace == namespace {
			leased := len(entry.Network)
			var last_address string
			if leased == 0 {
				last_address = base_address + subnet
			} else {
				last_address = entry.Network[leased-1].Address
			}

			new_address, res := leaseNewAddress(last_address)
			if !res {
				return "", false
			}

			// append to list
			network_list_data.List[i].Network = append(network_list_data.List[i].Network,
				ContainerNetwork{
					ContainerId: "",
					Address:     new_address,
				},
			)

			lease_address = new_address
		}
	}

	data, _ := json.MarshalIndent(network_list_data, "", "  ")
	if err := os.WriteFile(karakuripkgs.HITOHA_NETWORK_LIST, data, os.ModePerm); err != nil {
		panic(err)
	}

	return lease_address, true
}

func bindAddressToContainerId(namespace string, id string, address string) {
	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.HITOHA_NETWORK_LIST)
	if err != nil {
		panic(err)
	}

	var network_list_data NetworkList
	if err := json.Unmarshal(bytes, &network_list_data); err != nil {
		panic(err)
	}

	for i, entry := range network_list_data.List {
		if entry.Namespace == namespace {
			for j, container := range entry.Network {
				if container.Address == address {
					network_list_data.List[i].Network[j].ContainerId = id
				}
			}
		}
	}

	data, _ := json.MarshalIndent(network_list_data, "", "  ")
	if err := os.WriteFile(karakuripkgs.HITOHA_NETWORK_LIST, data, os.ModePerm); err != nil {
		panic(err)
	}
}

func freeAddress(id string) {
	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.HITOHA_NETWORK_LIST)
	if err != nil {
		panic(err)
	}

	var network_list_data NetworkList
	if err := json.Unmarshal(bytes, &network_list_data); err != nil {
		panic(err)
	}

	var new_network_list_data NetworkList
	for _, entry := range network_list_data.List {
		var tmp []ContainerNetwork
		for _, network := range entry.Network {
			if network.ContainerId != id {
				tmp = append(tmp, network)
			}
		}
		new_network_list_data.List = append(new_network_list_data.List,
			NetworkInfo{
				Namespace: entry.Namespace,
				Network:   tmp,
			},
		)
	}

	data, _ := json.MarshalIndent(new_network_list_data, "", "  ")
	if err := os.WriteFile(karakuripkgs.HITOHA_NETWORK_LIST, data, os.ModePerm); err != nil {
		panic(err)
	}
}

func setIpTables(mode string, host_port int, target_ip string, target_port int, protocol string) {
	source_port := strconv.Itoa(host_port)
	container_ip := (strings.Split(target_ip, "/"))[0]
	container_port := strconv.Itoa(target_port)

	var mode_option string
	if mode == "add" {
		mode_option = "-A"
	} else {
		mode_option = "-D"
	}

	prerouting := exec.Command("iptables", "-t", "nat", mode_option, "PREROUTING", "-p", protocol, "--dport", source_port, "-j", "DNAT", "--to-destination", container_ip+":"+container_port)
	if err := prerouting.Run(); err != nil {
		panic(err)
	}
	postrouting := exec.Command("iptables", "-t", "nat", mode_option, "POSTROUTING", "-p", protocol, "-d", container_ip, "--dport", container_port, "-j", "MASQUERADE")
	if err := postrouting.Run(); err != nil {
		panic(err)
	}
	container_in := exec.Command("iptables", mode_option, "FORWARD", "-p", protocol, "-d", container_ip, "--dport", container_port, "-j", "ACCEPT")
	if err := container_in.Run(); err != nil {
		panic(err)
	}
	container_out := exec.Command("iptables", mode_option, "FORWARD", "-p", protocol, "-s", container_ip, "--sport", container_port, "-j", "ACCEPT")
	if err := container_out.Run(); err != nil {
		panic(err)
	}
}

func SetupPortForwarding(mode string, network karakuripkgs.SpecNetwork) {
	container_ip := network.Address
	for _, entry := range network.Port {
		host_port := entry.HostPort
		container_port := entry.TargetPort
		protocol := entry.Protocol
		setIpTables(mode, host_port, container_ip, container_port, protocol)
	}
}
